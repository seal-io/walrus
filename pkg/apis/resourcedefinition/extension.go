package resourcedefinition

import (
	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/model/templateversion"
	"github.com/seal-io/walrus/pkg/datalisten/modelchange"
	"github.com/seal-io/walrus/utils/topic"
)

var (
	queryResourceFields = []string{
		resource.FieldName,
	}
	getResourceFields = resource.WithoutFields(
		resource.FieldUpdateTime)
	sortResourceFields = []string{
		resource.FieldName,
		resource.FieldCreateTime,
	}
)

func (h Handler) RouteGetResources(req RouteGetResourcesRequest) (RouteGetResourcesResponse, int, error) {
	query := h.modelClient.Resources().Query().
		Where(resource.ResourceDefinitionID(req.ID))

	if queries, ok := req.Querying(queryResourceFields); ok {
		query.Where(queries)
	}

	if stream := req.Stream; stream != nil {
		// Handle stream request.
		if fields, ok := req.Extracting(getResourceFields, getResourceFields...); ok {
			query.Select(fields...)
		}

		if orders, ok := req.Sorting(sortResourceFields, model.Desc(resource.FieldCreateTime)); ok {
			query.Order(orders...)
		}

		t, err := topic.Subscribe(modelchange.Resource)
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

			var items []*model.ResourceOutput

			ids := dm.IDs()

			switch dm.Type {
			case modelchange.EventTypeCreate, modelchange.EventTypeUpdate:
				entities, err := query.Clone().
					Where(resource.IDIn(ids...)).
					// Must append environment ID.
					Select(resource.FieldEnvironmentID).
					// Must extract template.
					Select(resource.FieldTemplateID).
					WithTemplate(func(tvq *model.TemplateVersionQuery) {
						tvq.Select(
							templateversion.FieldID,
							templateversion.FieldName,
							templateversion.FieldVersion)
					}).
					Unique(false).
					All(stream)
				if err != nil {
					return nil, 0, err
				}

				items = model.ExposeResources(entities)
			case modelchange.EventTypeDelete:
				items = make([]*model.ResourceOutput, len(ids))
				for i := range ids {
					items[i] = &model.ResourceOutput{
						ID:   ids[i],
						Name: dm.Data[i].Name,
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

	if fields, ok := req.Extracting(getResourceFields, getResourceFields...); ok {
		query.Select(fields...)
	}

	if orders, ok := req.Sorting(sortResourceFields, model.Desc(resource.FieldCreateTime)); ok {
		query.Order(orders...)
	}

	entities, err := query.
		// Must append environment ID.
		Select(resource.FieldEnvironmentID).
		// Must extract template.
		Select(resource.FieldTemplateID).
		WithTemplate(func(tvq *model.TemplateVersionQuery) {
			tvq.Select(
				templateversion.FieldID,
				templateversion.FieldName,
				templateversion.FieldVersion)
		}).
		Unique(false).
		All(req.Context)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeResources(entities), cnt, nil
}
