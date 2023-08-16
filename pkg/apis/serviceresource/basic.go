package serviceresource

import (
	"context"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/model/serviceresource"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/object"
	optypes "github.com/seal-io/seal/pkg/operator/types"
	pkgresource "github.com/seal-io/seal/pkg/serviceresources"
	"github.com/seal-io/seal/pkg/topic/datamessage"
	"github.com/seal-io/seal/utils/topic"
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

		t, err := topic.Subscribe(datamessage.ServiceResource)
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

			var items []*model.ServiceResourceOutput

			switch dm.Type {
			case datamessage.EventCreate, datamessage.EventUpdate:
				entities, err := getCollection(
					stream, query.Clone().Where(serviceresource.IDIn(dm.Data...)), req.WithoutKeys)
				if err != nil {
					return nil, 0, err
				}

				items = model.ExposeServiceResources(entities)
			case datamessage.EventDelete:
				items = make([]*model.ServiceResourceOutput, len(dm.Data))
				for i := range dm.Data {
					items[i] = &model.ServiceResourceOutput{
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
		req.Context, query, req.WithoutKeys)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeServiceResources(entities), cnt, nil
}

func getCollection(
	ctx context.Context,
	query *model.ServiceResourceQuery,
	withoutKeys bool,
) (model.ServiceResources, error) {
	wcOpts := func(cq *model.ConnectorQuery) {
		cq.Select(
			connector.FieldName,
			connector.FieldType,
			connector.FieldCategory,
			connector.FieldConfigVersion,
			connector.FieldConfigData)
	}

	// Query service resource with its components.
	entities, err := query.
		// Only "instance" type resources.
		Where(serviceresource.Shape(types.ServiceResourceShapeInstance)).
		// Must append service ID.
		Select(serviceresource.FieldServiceID).
		// Must extract connector.
		Select(serviceresource.FieldConnectorID).
		WithConnector(wcOpts).
		// Must extract components.
		WithComponents(func(rq *model.ServiceResourceQuery) {
			rq.Select(getFields...).
				Order(model.Desc(serviceresource.FieldCreateTime)).
				Where(serviceresource.Mode(types.ServiceResourceModeDiscovered)).
				WithConnector(wcOpts)
		}).
		Unique(false).
		All(ctx)
	if err != nil {
		return nil, err
	}

	// Return directly if no need next operations, e.g. Log, Exec and so on.
	if withoutKeys {
		return entities, nil
	}

	operators := make(map[object.ID]optypes.Operator)
	entities = pkgresource.SetKeys(ctx, entities, operators)

	return entities, nil
}
