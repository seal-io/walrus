package setting

import (
	"net/http"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/predicate"
	"github.com/seal-io/walrus/pkg/dao/model/setting"
	"github.com/seal-io/walrus/pkg/settings"
	"github.com/seal-io/walrus/utils/errorx"
)

func (h Handler) Get(req GetRequest) (GetResponse, error) {
	entity, err := h.modelClient.Settings().Query().
		Where(
			setting.ID(req.ID),
			setting.Private(false)).
		Only(req.Context)
	if err != nil {
		return nil, err
	}

	return exposeSetting(entity), nil
}

func (h Handler) Update(req UpdateRequest) error {
	// Bypass the validations or cascade works to settings definition.
	return h.modelClient.WithTx(req.Context, func(tx *model.Tx) error {
		s := settings.Index(req.Name)
		if s == nil {
			return errorx.HttpErrorf(http.StatusNotFound, "setting %s not found", req.Name)
		}

		return s.Set(req.Context, tx, req.Value)
	})
}

var (
	queryFields = []string{
		setting.FieldName,
	}
	getFields = setting.WithoutFields(
		setting.FieldPrivate,
	)
	sortFields = []string{
		setting.FieldName,
		setting.FieldCreateTime,
	}
)

func (h Handler) CollectionGet(req CollectionGetRequest) (CollectionGetResponse, int, error) {
	query := h.modelClient.Settings().Query().
		Where(setting.Private(false))

	{
		ps := make([]predicate.Setting, 0, len(req.IDs)+len(req.Names))

		for i := range req.IDs {
			ps = append(ps, setting.ID(req.IDs[i]))
		}

		for i := range req.Names {
			ps = append(ps, setting.Name(req.Names[i]))
		}

		if len(ps) != 0 {
			query.Where(setting.Or(ps...))
		}
	}

	if queries, ok := req.Querying(queryFields); ok {
		query.Where(queries)
	}

	// Get count.
	cnt, err := query.Clone().Count(req.Context)
	if err != nil {
		return nil, 0, err
	}

	// Get entities.
	if limit, offset, ok := req.Paging(); ok {
		query.Limit(limit).Offset(offset)
	}

	if fields, ok := req.Extracting(getFields, getFields...); ok {
		query.Select(fields...)
	}

	if orders, ok := req.Sorting(sortFields, model.Desc(setting.FieldCreateTime)); ok {
		query.Order(orders...)
	}

	entities, err := query.
		Unique(false).
		All(req.Context)
	if err != nil {
		return nil, 0, err
	}

	return exposeSettings(entities), cnt, nil
}

func (h Handler) CollectionUpdate(req CollectionUpdateRequest) error {
	// Bypass the validations or cascade works to settings definition.
	return h.modelClient.WithTx(req.Context, func(tx *model.Tx) error {
		for i := range req.Items {
			s := settings.Index(req.Items[i].Name)
			if s == nil {
				return errorx.HttpErrorf(http.StatusNotFound, "setting %s not found", req.Items[i].Name)
			}

			err := s.Set(req.Context, tx, req.Items[i].Value)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
