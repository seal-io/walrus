package catalog

import (
	"github.com/seal-io/seal/pkg/apis/runtime"
	catalogbus "github.com/seal-io/seal/pkg/bus/catalog"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/catalog"
	"github.com/seal-io/seal/pkg/dao/types/object"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/topic/datamessage"
	"github.com/seal-io/seal/utils/topic"
)

func (h Handler) Create(req CreateRequest) (CreateResponse, error) {
	entity := req.Model()

	status.CatalogStatusInitialized.Unknown(entity, "")
	entity.Status.SetSummary(status.WalkCatalog(&entity.Status))

	entity, err := h.modelClient.Catalogs().Create().
		Set(entity).
		Save(req.Context)
	if err != nil {
		return nil, err
	}

	err = catalogbus.Notify(req.Context, h.modelClient, entity)
	if err != nil {
		return nil, err
	}

	return model.ExposeCatalog(entity), nil
}

func (h Handler) Get(req GetRequest) (GetResponse, error) {
	entity, err := h.modelClient.Catalogs().Get(req.Context, req.ID)
	if err != nil {
		return nil, err
	}

	return model.ExposeCatalog(entity), nil
}

func (h Handler) Update(req UpdateRequest) error {
	entity := req.Model()

	status.CatalogStatusInitialized.Reset(entity, "Initializing catalog template")
	entity.Status.SetSummary(status.WalkCatalog(&entity.Status))

	entity, err := h.modelClient.Catalogs().UpdateOne(entity).
		Set(entity).
		Save(req.Context)
	if err != nil {
		return err
	}

	return catalogbus.Notify(req.Context, h.modelClient, entity)
}

func (h Handler) Delete(req DeleteRequest) error {
	return h.modelClient.Catalogs().DeleteOneID(req.ID).
		Exec(req.Context)
}

var (
	queryFields = []string{
		catalog.FieldName,
		catalog.FieldType,
	}
	getFields  = catalog.WithoutFields()
	sortFields = []string{
		catalog.FieldID,
		catalog.FieldName,
		catalog.FieldType,
		catalog.FieldSource,
	}
)

func (h Handler) CollectionGet(req CollectionGetRequest) (CollectionGetResponse, int, error) {
	query := h.modelClient.Catalogs().Query()

	if queries, ok := req.Querying(queryFields); ok {
		query = query.Where(queries)
	}

	if stream := req.Stream; stream != nil {
		// Handle stream request.
		if fields, ok := req.Extracting(getFields, getFields...); ok {
			query.Select(fields...)
		}

		if orders, ok := req.Sorting(sortFields, model.Desc(catalog.FieldCreateTime)); ok {
			query.Order(orders...)
		}

		t, err := topic.Subscribe(datamessage.Catalog)
		if err != nil {
			return nil, 0, err
		}

		defer func() { t.Unsubscribe() }()

		for {
			var event topic.Event

			event, err = t.Receive(stream)
			if err != nil {
				return nil, 0, err
			}

			dm, ok := event.Data.(datamessage.Message[object.ID])
			if !ok {
				continue
			}

			var items []*model.CatalogOutput

			switch dm.Type {
			case datamessage.EventCreate, datamessage.EventUpdate:
				entities, err := query.Clone().
					Where(catalog.IDIn(dm.Data...)).
					Unique(false).
					All(stream)
				if err != nil {
					return nil, 0, err
				}

				items = model.ExposeCatalogs(entities)
			case datamessage.EventDelete:
				items = make([]*model.CatalogOutput, len(dm.Data))
				for i := range dm.Data {
					items[i] = &model.CatalogOutput{
						ID: dm.Data[i],
					}
				}
			}

			if len(items) == 0 {
				continue
			}

			resp := runtime.TypedResponse(dm.Type.String(), items)
			if err = stream.SendJSON(resp); err != nil {
				return nil, 0, err
			}
		}
	}

	// Handle normal request.

	// Get count.
	count, err := query.Clone().Count(req.Context)
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

	if orders, ok := req.Sorting(sortFields, model.Desc(catalog.FieldCreateTime)); ok {
		query.Order(orders...)
	}

	entities, err := query.
		Unique(false).
		All(req.Context)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeCatalogs(entities), count, nil
}

func (h Handler) CollectionDelete(req CollectionDeleteRequest) error {
	ids := req.IDs()

	return h.modelClient.WithTx(req.Context, func(tx *model.Tx) error {
		_, err := tx.Catalogs().Delete().
			Where(catalog.IDIn(ids...)).
			Exec(req.Context)

		return err
	})
}
