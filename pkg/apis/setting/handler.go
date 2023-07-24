package setting

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/apis/setting/view"
	settingbus "github.com/seal-io/seal/pkg/bus/setting"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/setting"
	"github.com/seal-io/seal/pkg/settings"
)

func Handle(mc model.ClientSet) Handler {
	return Handler{
		modelClient: mc,
	}
}

type Handler struct {
	modelClient model.ClientSet
}

func (h Handler) Kind() string {
	return "Setting"
}

func (h Handler) Validating() any {
	return h.modelClient
}

// Basic APIs.

func (h Handler) Update(ctx *gin.Context, req view.UpdateRequest) error {
	// Bypass the validations or cascade works to settings definition.
	return h.modelClient.WithTx(ctx, func(tx *model.Tx) error {
		s := settings.Index(req.Name)
		if s == nil {
			return runtime.Errorc(http.StatusNotFound)
		}

		changed, err := s.Set(ctx, tx, req.Value)
		if err != nil {
			return err
		}

		if !changed {
			return nil
		}

		return settingbus.Notify(ctx, tx, model.Settings{
			&model.Setting{
				Name:  req.Name,
				Value: req.Value,
			},
		})
	})
}

func (h Handler) Get(ctx *gin.Context, req view.GetRequest) (*view.GetResponse, error) {
	var p predicate.Setting
	if req.ID.Valid() {
		p = setting.ID(req.ID)
	} else {
		p = setting.Name(string(req.ID))
	}

	entity, err := h.modelClient.Settings().Query().
		Where(
			p,
			setting.Private(false)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return expose(entity)[0], nil
}

// Batch APIs.

func (h Handler) CollectionUpdate(ctx *gin.Context, req view.CollectionUpdateRequest) error {
	// Bypass the validations or cascade works to settings definition.
	return h.modelClient.WithTx(ctx, func(tx *model.Tx) error {
		list := make(model.Settings, 0, len(req))

		for i := range req {
			s := settings.Index(req[i].Name)
			if s == nil {
				return runtime.Errorc(http.StatusNotFound)
			}

			changed, err := s.Set(ctx, tx, req[i].Value)
			if err != nil {
				return err
			}

			if !changed {
				continue
			}

			list = append(list, &model.Setting{
				Name:  req[i].Name,
				Value: req[i].Value,
			})
		}

		if len(list) == 0 {
			return nil
		}

		return settingbus.Notify(ctx, tx, list)
	})
}

var getFields = setting.WithoutFields(
	setting.FieldCreateTime,
	setting.FieldPrivate,
)

func (h Handler) CollectionGet(
	ctx *gin.Context,
	req view.CollectionGetRequest,
) (view.CollectionGetResponse, int, error) {
	input := []predicate.Setting{
		setting.Private(false),
	}
	{
		var sps []predicate.Setting

		for i := 0; i < len(req.IDs); i++ {
			if req.IDs[i].Valid() {
				sps = append(sps, setting.ID(req.IDs[i]))
			} else {
				sps = append(sps, setting.Name(string(req.IDs[i])))
			}
		}

		if len(sps) != 0 {
			input = append(input, setting.Or(sps...))
		}
	}

	query := h.modelClient.Settings().Query().
		Where(input...)

	// Get count.
	cnt, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Get entities.
	if limit, offset, ok := req.Paging(); !ok {
		query.Limit(limit).Offset(offset)
	}

	entities, err := query.
		Order(model.Desc(setting.FieldCreateTime)).
		Select(getFields...).
		// Allow returning without sorting keys.
		Unique(false).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return expose(entities...), cnt, nil
}

func expose(ss ...*model.Setting) view.CollectionGetResponse {
	rs := make([]*view.GetResponse, 0, len(ss))

	for _, s := range ss {
		r := view.GetResponse{
			SettingOutput: model.ExposeSetting(s),
			Configured:    s.Value != "",
		}

		// Erases the value of sensitive setting.
		if s.Sensitive != nil && *s.Sensitive {
			r.Value = ""
		}

		rs = append(rs, &r)
	}

	return rs
}

// Extensional APIs.
