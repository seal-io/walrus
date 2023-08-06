package perspective

import (
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/perspective"
)

func (h Handler) Create(req CreateRequest) (CreateResponse, error) {
	entity := req.Model()

	entity, err := h.modelClient.Perspectives().Create().
		Set(entity).
		Save(req.Context)
	if err != nil {
		return nil, err
	}

	return model.ExposePerspective(entity), nil
}

func (h Handler) Get(req GetRequest) (GetResponse, error) {
	entity, err := h.modelClient.Perspectives().Get(req.Context, req.ID)
	if err != nil {
		return nil, err
	}

	return model.ExposePerspective(entity), nil
}

func (h Handler) Update(req UpdateRequest) error {
	entity := req.Model()

	return h.modelClient.Perspectives().UpdateOne(entity).
		Set(entity).
		Exec(req.Context)
}

func (h Handler) Delete(req DeleteRequest) error {
	return h.modelClient.Perspectives().DeleteOneID(req.ID).
		Exec(req.Context)
}

var (
	queryFields = []string{
		perspective.FieldName,
	}
	getFields  = perspective.WithoutFields()
	sortFields = []string{
		perspective.FieldName,
		perspective.FieldCreateTime,
	}
)

func (h Handler) CollectionGet(
	req CollectionGetRequest,
) (CollectionGetResponse, int, error) {
	query := h.modelClient.Perspectives().Query()

	if req.Name != "" {
		query.Where(perspective.NameContains(req.Name))
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

	if orders, ok := req.Sorting(sortFields, model.Desc(perspective.FieldCreateTime)); ok {
		query.Order(orders...)
	}

	entities, err := query.
		Unique(false).
		All(req.Context)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposePerspectives(entities), cnt, nil
}

func (h Handler) CollectionDelete(req CollectionDeleteRequest) error {
	ids := req.IDs()

	return h.modelClient.WithTx(req.Context, func(tx *model.Tx) error {
		_, err := tx.Perspectives().Delete().
			Where(perspective.IDIn(ids...)).
			Exec(req.Context)

		return err
	})
}
