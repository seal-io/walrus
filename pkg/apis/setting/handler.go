package setting

import (
	"github.com/gin-gonic/gin"

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
		var changed, err = settings.Index(req.Name).Set(ctx, tx, *req.Value)
		if err != nil {
			return err
		}
		if !changed {
			return nil
		}
		return settingbus.Notify(ctx, tx, model.Settings{
			&model.Setting{
				Name:  req.Name,
				Value: *req.Value,
			}})
	})
}

func (h Handler) Get(ctx *gin.Context, req view.GetRequest) (view.GetResponse, error) {
	var input = []predicate.Setting{
		setting.Private(false),
	}
	if req.ID.IsNaive() {
		input = append(input, setting.ID(req.ID))
	} else {
		var keys = req.ID.Split()
		input = append(input, setting.Name(keys[0]))
	}

	var entity, err = h.modelClient.Settings().Query().
		Where(input...).
		Select(setting.WithoutFields(
			setting.FieldCreateTime, setting.FieldUpdateTime, setting.FieldPrivate)...).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	sanitize(model.Settings{entity})

	return model.ExposeSetting(entity), nil
}

// Batch APIs.

func (h Handler) CollectionUpdate(ctx *gin.Context, req view.CollectionUpdateRequest) error {
	// Bypass the validations or cascade works to settings definition.
	return h.modelClient.WithTx(ctx, func(tx *model.Tx) error {
		var list = make(model.Settings, 0, len(req))
		for i := range req {
			var changed, err = settings.Index(req[i].Name).Set(ctx, tx, *req[i].Value)
			if err != nil {
				return err
			}
			if !changed {
				continue
			}
			list = append(list, &model.Setting{
				Name:  req[i].Name,
				Value: *req[i].Value,
			})
		}
		if len(list) == 0 {
			return nil
		}
		return settingbus.Notify(ctx, tx, list)
	})
}

var (
	getFields = setting.WithoutFields(
		setting.FieldCreateTime,
		setting.FieldPrivate)
)

func (h Handler) CollectionGet(ctx *gin.Context, req view.CollectionGetRequest) (view.CollectionGetResponse, int, error) {
	var input = []predicate.Setting{
		setting.Private(false),
	}
	{
		var sps []predicate.Setting
		for i := 0; i < len(req.IDs); i++ {
			if req.IDs[i].IsNaive() {
				sps = append(sps, setting.ID(req.IDs[i]))
			} else {
				var keys = req.IDs[i].Split()
				sps = append(sps, setting.Name(keys[0]))
			}
		}
		if len(sps) != 0 {
			input = append(input, setting.Or(sps...))
		}
	}

	var query = h.modelClient.Settings().Query().
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

	sanitize(entities)

	return model.ExposeSettings(entities), cnt, nil
}

// FIXME This is a temporary way to protect openai API token. Refactor when setting value encryption is ready.
func sanitize(ss model.Settings) {
	for _, s := range ss {
		if s.Name == settings.OpenAiApiToken.Name() {
			s.Value = ""
		}
	}
}

// Extensional APIs.
