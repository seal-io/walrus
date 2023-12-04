package resourcerevision

import (
	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/resourcerevision"
	"github.com/seal-io/walrus/pkg/datalisten/modelchange"
	"github.com/seal-io/walrus/utils/topic"
)

func (h Handler) Get(req GetRequest) (GetResponse, error) {
	entity, err := h.modelClient.ResourceRevisions().Get(req.Context, req.ID)
	if err != nil {
		return nil, err
	}

	return model.ExposeResourceRevision(entity), nil
}

var (
	getFields = resourcerevision.WithoutFields(
		resourcerevision.FieldRecord,
		resourcerevision.FieldInputPlan,
		resourcerevision.FieldOutput,
		resourcerevision.FieldTemplateName,
		resourcerevision.FieldTemplateVersion,
		resourcerevision.FieldAttributes,
		resourcerevision.FieldVariables,
	)
	sortFields = []string{
		resourcerevision.FieldCreateTime,
	}
)

func (h Handler) CollectionGet(req CollectionGetRequest) (CollectionGetResponse, int, error) {
	query := h.modelClient.ResourceRevisions().Query()

	if req.Resource != nil && req.Resource.ID != "" {
		query.Where(resourcerevision.ResourceID(req.Resource.ID))
	}

	if stream := req.Stream; stream != nil {
		// Handle stream request.
		if fields, ok := req.Extracting(getFields, getFields...); ok {
			query.Select(fields...)
		}

		if orders, ok := req.Sorting(sortFields, model.Desc(resourcerevision.FieldCreateTime)); ok {
			query.Order(orders...)
		}

		t, err := topic.Subscribe(modelchange.ResourceRevision)
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

			dm, ok := event.Data.(modelchange.Event)
			if !ok {
				continue
			}

			var items []*model.ResourceRevisionOutput

			ids := dm.IDs()

			switch dm.Type {
			case modelchange.EventTypeCreate, modelchange.EventTypeUpdate:
				revisions, err := query.Clone().
					Where(resourcerevision.IDIn(ids...)).
					// Must append service ID.
					Select(resourcerevision.FieldResourceID).
					Unique(false).
					All(stream)
				if err != nil {
					return nil, 0, err
				}

				items = model.ExposeResourceRevisions(revisions)
			case modelchange.EventTypeDelete:
				items = make([]*model.ResourceRevisionOutput, len(ids))
				for i := range ids {
					items[i] = &model.ResourceRevisionOutput{
						ID: ids[i],
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

	// Handler normal request.

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

	if orders, ok := req.Sorting(sortFields, model.Desc(resourcerevision.FieldCreateTime)); ok {
		query.Order(orders...)
	}

	entities, err := query.
		// Must append service ID.
		Select(resourcerevision.FieldResourceID).
		Unique(false).
		All(req.Context)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeResourceRevisions(entities), cnt, nil
}

func (h Handler) CollectionDelete(req CollectionDeleteRequest) error {
	ids := req.IDs()

	return h.modelClient.WithTx(req.Context, func(tx *model.Tx) error {
		_, err := tx.ResourceRevisions().Delete().
			Where(resourcerevision.IDIn(ids...)).
			Exec(req.Context)

		return err
	})
}
