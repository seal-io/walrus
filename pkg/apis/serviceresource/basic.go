package serviceresource

import (
	"context"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/serviceresource"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/datalisten/modelchange"
	pkgresource "github.com/seal-io/walrus/pkg/serviceresources"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/topic"
)

var (
	queryFields = []string{
		serviceresource.FieldName,
	}
	getFields = serviceresource.WithoutFields(
		serviceresource.FieldUpdateTime)
	sortFields = []string{
		serviceresource.FieldMode,
		serviceresource.FieldType,
		serviceresource.FieldName,
		serviceresource.FieldCreateTime,
	}
)

func (h Handler) CollectionGet(req CollectionGetRequest) (CollectionGetResponse, int, error) {
	query := h.modelClient.ServiceResources().Query().
		Where(serviceresource.ServiceID(req.Service.ID))

	if queries, ok := req.Querying(queryFields); ok {
		query.Where(queries)
	}

	if stream := req.Stream; stream != nil {
		// Handle stream request.
		if fields, ok := req.Extracting(getFields, getFields...); ok {
			query.Select(fields...)
		}

		if orders, ok := req.Sorting(sortFields, model.Desc(serviceresource.FieldCreateTime)); ok {
			query.Order(orders...)
		}

		// Exclude "data" mode resources.
		query.Where(serviceresource.ModeNEQ(types.ServiceResourceModeData))

		t, err := topic.Subscribe(modelchange.ServiceResource)
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

			var items []*model.ServiceResourceOutput

			switch dm.Type {
			case modelchange.EventTypeCreate, modelchange.EventTypeUpdate:
				entities, err := getCollection(
					stream, h.modelClient, query.Clone().Where(serviceresource.IDIn(dm.IDs...)), req.WithoutKeys)
				if err != nil {
					return nil, 0, err
				}

				items = model.ExposeServiceResources(entities)
			case modelchange.EventTypeDelete:
				items = make([]*model.ServiceResourceOutput, len(dm.IDs))
				for i := range dm.IDs {
					items[i] = &model.ServiceResourceOutput{
						ID: dm.IDs[i],
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

	if orders, ok := req.Sorting(sortFields, model.Desc(serviceresource.FieldCreateTime)); ok {
		query.Order(orders...)
	}

	// Only "managed" mode resource.
	query.Where(serviceresource.ModeEQ(types.ServiceResourceModeManaged))

	entities, err := getCollection(
		req.Context, h.modelClient, query, req.WithoutKeys)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeServiceResources(entities), cnt, nil
}

func getCollection(
	ctx context.Context,
	modelClient model.ClientSet,
	query *model.ServiceResourceQuery,
	withoutKeys bool,
) (model.ServiceResources, error) {
	// Query service resource with its components.
	entities, err := query.
		// Only "instance" type resources.
		Where(serviceresource.Shape(types.ServiceResourceShapeInstance)).
		// Must append the following IDs.
		Select(
			serviceresource.FieldServiceID,
			serviceresource.FieldConnectorID).
		// Must extract components.
		WithComponents(func(rq *model.ServiceResourceQuery) {
			rq.Select(getFields...).
				Order(model.Desc(serviceresource.FieldCreateTime)).
				Where(serviceresource.Mode(types.ServiceResourceModeDiscovered))
		}).
		Unique(false).
		All(ctx)
	if err != nil {
		return nil, err
	}

	// Set keys for next operations, e.g. Log, Exec and so on.
	if !withoutKeys {
		pkgresource.SetKeys(
			ctx,
			log.WithName("api").WithName("service-resource"),
			modelClient,
			entities,
			nil)
	}

	return entities, nil
}
