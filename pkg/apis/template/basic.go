package template

import (
	"github.com/seal-io/walrus/pkg/apis/runtime"
	modbus "github.com/seal-io/walrus/pkg/bus/template"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/template"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/topic/datamessage"
	"github.com/seal-io/walrus/utils/topic"
)

func (h Handler) Create(req CreateRequest) (CreateResponse, error) {
	entity := req.Model()

	err := h.modelClient.WithTx(req.Context, func(tx *model.Tx) error {
		created, err := tx.Templates().Create().
			Set(entity).
			Save(req.Context)
		if err != nil {
			return err
		}

		if err = modbus.Notify(req.Context, created); err != nil {
			return err
		}

		entity = created

		return nil
	})
	if err != nil {
		return nil, err
	}

	return model.ExposeTemplate(entity), nil
}

func (h Handler) Get(req GetRequest) (GetResponse, error) {
	entity, err := h.modelClient.Templates().
		Get(req.Context, req.ID)
	if err != nil {
		return nil, err
	}

	return model.ExposeTemplate(entity), nil
}

func (h Handler) Update(req UpdateRequest) error {
	entity := req.Model()

	status.TemplateStatusInitialized.Unknown(entity, "Initializing template")
	entity.Status.SetSummary(status.WalkTemplate(&entity.Status))

	updated, err := h.modelClient.Templates().UpdateOne(entity).
		Set(entity).
		Save(req.Context)
	if err != nil {
		return err
	}

	return modbus.Notify(req.Context, updated)
}

func (h Handler) Delete(req DeleteRequest) error {
	return h.modelClient.Templates().DeleteOneID(req.ID).
		Exec(req.Context)
}

var (
	queryFields = []string{
		template.FieldName,
	}
	getFields  = template.WithoutFields()
	sortFields = []string{
		template.FieldName,
		template.FieldStatus,
		template.FieldCreateTime,
	}
)

func (h Handler) CollectionGet(req CollectionGetRequest) (CollectionGetResponse, int, error) {
	query := h.modelClient.Templates().Query()

	if len(req.CatalogIDs) != 0 {
		query.Where(template.CatalogIDIn(req.CatalogIDs...))
	}

	if queries, ok := req.Querying(queryFields); ok {
		query.Where(queries)
	}

	if stream := req.Stream; stream != nil {
		// Handle stream request.
		if fields, ok := req.Extracting(getFields, getFields...); ok {
			query.Select(fields...)
		}

		if orders, ok := req.Sorting(sortFields, model.Desc(template.FieldCreateTime)); ok {
			query.Order(orders...)
		}

		t, err := topic.Subscribe(datamessage.Template)
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

			var items []*model.TemplateOutput

			switch dm.Type {
			case datamessage.EventCreate, datamessage.EventUpdate:
				entities, err := query.Clone().
					// Must extract catalog ID.
					Select(template.FieldCatalogID).
					Where(template.IDIn(dm.Data...)).
					Unique(false).
					All(stream)
				if err != nil {
					return nil, 0, err
				}

				items = model.ExposeTemplates(entities)
			case datamessage.EventDelete:
				items = make([]*model.TemplateOutput, len(dm.Data))
				for i := range dm.Data {
					items[i] = &model.TemplateOutput{
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

	if orders, ok := req.Sorting(sortFields, model.Desc(template.FieldCreateTime)); ok {
		query.Order(orders...)
	}

	entities, err := query.
		// Must extract catalog ID.
		Select(template.FieldCatalogID).
		Unique(false).
		All(req.Context)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeTemplates(entities), cnt, nil
}

func (h Handler) CollectionDelete(req CollectionDeleteRequest) error {
	ids := req.IDs()

	return h.modelClient.WithTx(req.Context, func(tx *model.Tx) error {
		_, err := tx.Templates().Delete().
			Where(template.IDIn(ids...)).
			Exec(req.Context)

		return err
	})
}
