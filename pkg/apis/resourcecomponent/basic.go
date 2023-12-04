package resourcecomponent

import (
	"context"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/resourcecomponent"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/datalisten/modelchange"
	pkgresource "github.com/seal-io/walrus/pkg/resourcecomponents"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/topic"
)

var (
	queryFields = []string{
		resourcecomponent.FieldName,
	}
	getFields = resourcecomponent.WithoutFields(
		resourcecomponent.FieldUpdateTime)
	sortFields = []string{
		resourcecomponent.FieldMode,
		resourcecomponent.FieldType,
		resourcecomponent.FieldName,
		resourcecomponent.FieldCreateTime,
	}
)

func (h Handler) CollectionGet(req CollectionGetRequest) (CollectionGetResponse, int, error) {
	query := h.modelClient.ResourceComponents().Query().
		Where(resourcecomponent.ResourceID(req.Resource.ID))

	if queries, ok := req.Querying(queryFields); ok {
		query.Where(queries)
	}

	if stream := req.Stream; stream != nil {
		// Handle stream request.
		if fields, ok := req.Extracting(getFields, getFields...); ok {
			query.Select(fields...)
		}

		if orders, ok := req.Sorting(sortFields, model.Desc(resourcecomponent.FieldCreateTime)); ok {
			query.Order(orders...)
		}

		// Exclude "data" mode resources.
		query.Where(resourcecomponent.ModeNEQ(types.ResourceComponentModeData))

		t, err := topic.Subscribe(modelchange.ResourceComponent)
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

			var items []*model.ResourceComponentOutput

			ids := dm.IDs()

			switch dm.Type {
			case modelchange.EventTypeCreate, modelchange.EventTypeUpdate:
				entities, err := getCollection(
					stream, h.modelClient, query.Clone().Where(resourcecomponent.IDIn(ids...)), req.WithoutKeys)
				if err != nil {
					return nil, 0, err
				}

				items = model.ExposeResourceComponents(entities)
			case modelchange.EventTypeDelete:
				items = make([]*model.ResourceComponentOutput, len(ids))
				for i := range ids {
					items[i] = &model.ResourceComponentOutput{
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

	if orders, ok := req.Sorting(sortFields, model.Desc(resourcecomponent.FieldCreateTime)); ok {
		query.Order(orders...)
	}

	// Only "managed" mode resource.
	query.Where(resourcecomponent.ModeEQ(types.ResourceComponentModeManaged))

	entities, err := getCollection(
		req.Context, h.modelClient, query, req.WithoutKeys)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeResourceComponents(entities), cnt, nil
}

func getCollection(
	ctx context.Context,
	modelClient model.ClientSet,
	query *model.ResourceComponentQuery,
	withoutKeys bool,
) (model.ResourceComponents, error) {
	// Query resource component with its components.
	entities, err := query.
		// Only "instance" type resources.
		Where(resourcecomponent.Shape(types.ResourceComponentShapeInstance)).
		// Must append the following IDs.
		Select(
			resourcecomponent.FieldResourceID,
			resourcecomponent.FieldConnectorID).
		// Must extract components.
		WithComponents(func(rq *model.ResourceComponentQuery) {
			rq.Select(getFields...).
				Order(model.Desc(resourcecomponent.FieldCreateTime)).
				Where(resourcecomponent.Mode(types.ResourceComponentModeDiscovered))
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
			log.WithName("api").WithName("resource-component"),
			modelClient,
			entities,
			nil)
	}

	return entities, nil
}
