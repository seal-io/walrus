package connector

import (
	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/connector/view"
	"github.com/seal-io/seal/pkg/apis/runtime"
	connbus "github.com/seal-io/seal/pkg/bus/connector"
	pkgconn "github.com/seal-io/seal/pkg/connectors"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/topic/datamessage"
	"github.com/seal-io/seal/utils/topic"
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
	return "Connector"
}

func (h Handler) Validating() any {
	return h.modelClient
}

// Basic APIs

func (h Handler) Create(ctx *gin.Context, req view.CreateRequest) (view.CreateResponse, error) {
	var entity = req.Model()

	var creates, err = dao.ConnectorCreates(h.modelClient, entity)
	if err != nil {
		return nil, err
	}
	entity, err = creates[0].Save(ctx)
	if err != nil {
		return nil, err
	}

	err = connbus.Notify(ctx, entity, false)
	if err != nil {
		return nil, err
	}

	return model.ExposeConnector(entity), nil
}

func (h Handler) Delete(ctx *gin.Context, req view.DeleteRequest) error {
	return h.modelClient.Connectors().DeleteOne(req.Model()).Exec(ctx)
}

func (h Handler) Update(ctx *gin.Context, req view.UpdateRequest) error {
	var entity = req.Model()

	var update, err = dao.ConnectorUpdate(h.modelClient, entity)
	if err != nil {
		return err
	}
	entity, err = update.Save(ctx)
	if err != nil {
		return err
	}

	return connbus.Notify(ctx, entity, false)
}

func (h Handler) Get(ctx *gin.Context, req view.GetRequest) (view.GetResponse, error) {
	var entity, err = h.modelClient.Connectors().Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return model.ExposeConnector(entity), nil
}

func (h Handler) Stream(ctx runtime.RequestStream, req view.StreamRequest) error {
	var t, err = topic.Subscribe(datamessage.Connector)
	if err != nil {
		return err
	}

	defer func() { t.Unsubscribe() }()
	for {
		var event topic.Event
		event, err = t.Receive(ctx)
		if err != nil {
			return err
		}
		dm, ok := event.Data.(datamessage.Message)
		if !ok {
			continue
		}

		var streamData view.StreamResponse
		for _, id := range dm.Data {
			if id != req.ID {
				continue
			}
			switch dm.Type {
			case datamessage.EventCreate, datamessage.EventUpdate:
				entity, err := h.modelClient.Connectors().Get(ctx, id)
				if err != nil {
					return err
				}
				streamData = view.StreamResponse{
					Type:       dm.Type,
					Collection: []*model.ConnectorOutput{model.ExposeConnector(entity)},
				}
			case datamessage.EventDelete:
				streamData = view.StreamResponse{
					Type: dm.Type,
					IDs:  dm.Data,
				}
			}
		}
		err = ctx.SendJSON(streamData)
		if err != nil {
			return err
		}
	}
}

// Batch APIs

func (h Handler) CollectionDelete(ctx *gin.Context, req view.CollectionDeleteRequest) error {
	return h.modelClient.WithTx(ctx, func(tx *model.Tx) (err error) {
		for i := range req {
			err = tx.Connectors().DeleteOne(req[i].Model()).
				Exec(ctx)
			if err != nil {
				return err
			}
		}
		return
	})
}

var (
	queryFields = []string{
		connector.FieldName,
	}
	getFields = connector.WithoutFields(
		connector.FieldConfigVersion,
		connector.FieldConfigData,
		connector.FieldUpdateTime)
	sortFields = []string{
		connector.FieldName,
		connector.FieldType,
		connector.FieldCreateTime}
)

func (h Handler) CollectionGet(ctx *gin.Context, req view.CollectionGetRequest) (view.CollectionGetResponse, int, error) {
	var query = h.modelClient.Connectors().Query()
	if queries, ok := req.Querying(queryFields); ok {
		query.Where(queries)
	}

	// get count.
	cnt, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// get entities.
	if limit, offset, ok := req.Paging(); ok {
		query.Limit(limit).Offset(offset)
	}
	if fields, ok := req.Extracting(getFields, getFields...); ok {
		query.Select(fields...)
	}
	if orders, ok := req.Sorting(sortFields, model.Desc(connector.FieldCreateTime)); ok {
		query.Order(orders...)
	}

	if req.Category != "" {
		query.Where(connector.Category(req.Category))
	}

	if req.Type != "" {
		query.Where(connector.Type(req.Type))
	}

	entities, err := query.
		// allow returning without sorting keys.
		Unique(false).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeConnectors(entities), cnt, nil
}

func (h Handler) CollectionStream(ctx runtime.RequestStream, req view.CollectionStreamRequest) error {
	var t, err = topic.Subscribe(datamessage.Connector)
	if err != nil {
		return err
	}

	query := h.modelClient.Connectors().Query()
	if fields, ok := req.Extracting(getFields, getFields...); ok {
		query.Select(fields...)
	}

	defer func() { t.Unsubscribe() }()
	for {
		var event topic.Event
		event, err = t.Receive(ctx)
		if err != nil {
			return err
		}
		dm, ok := event.Data.(datamessage.Message)
		if !ok {
			continue
		}

		var streamData view.StreamResponse
		switch dm.Type {
		case datamessage.EventCreate, datamessage.EventUpdate:
			connectors, err := query.Clone().
				Where(connector.IDIn(dm.Data...)).
				Unique(false).
				All(ctx)

			if err != nil {
				return err
			}
			streamData = view.StreamResponse{
				Type:       dm.Type,
				Collection: model.ExposeConnectors(connectors),
			}
		case datamessage.EventDelete:
			streamData = view.StreamResponse{
				Type: dm.Type,
				IDs:  dm.Data,
			}
		}
		if len(streamData.IDs) == 0 && len(streamData.Collection) == 0 {
			continue
		}
		err = ctx.SendJSON(streamData)
		if err != nil {
			return err
		}
	}
}

// Extensional APIs

func (h Handler) RouteApplyCostTools(ctx *gin.Context, req view.ApplyCostToolsRequest) error {
	o, err := h.modelClient.Connectors().Get(ctx, req.ID)
	if err != nil {
		return err
	}

	status.ConnectorStatusToolsDeployed.Unknown(o, "")
	if err = pkgconn.UpdateStatus(ctx, h.modelClient, o); err != nil {
		return err
	}

	return connbus.Notify(ctx, o, true)
}

func (h Handler) RouteSyncCostOpsData(ctx *gin.Context, req view.SyncCostDataRequest) error {
	o, err := h.modelClient.Connectors().Get(ctx, req.ID)
	if err != nil {
		return err
	}

	status.ConnectorStatusCostSynced.Unknown(o, "")
	if err = pkgconn.UpdateStatus(ctx, h.modelClient, o); err != nil {
		return err
	}

	syncer := pkgconn.NewStatusSyncer(h.modelClient)
	return syncer.SyncFinOpsStatus(ctx, o)
}
