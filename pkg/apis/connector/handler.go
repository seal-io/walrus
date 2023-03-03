package connector

import (
	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/connector/view"
	pkgconn "github.com/seal-io/seal/pkg/connectors"
	"github.com/seal-io/seal/pkg/costs/scheduler"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/types/status"
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

	err = pkgconn.Notify(ctx, h.modelClient, entity, false)
	if err != nil {
		return nil, err
	}
	return model.ExposeConnector(entity), nil
}

func (h Handler) Delete(ctx *gin.Context, req view.DeleteRequest) error {
	return h.modelClient.Connectors().DeleteOne(req.Model()).Exec(ctx)
}

func (h Handler) Update(ctx *gin.Context, req view.UpdateRequest) (view.UpdateResponse, error) {
	var entity = req.Model()

	var update, err = dao.ConnectorUpdate(h.modelClient, entity)
	if err != nil {
		return nil, err
	}
	entity, err = update.Save(ctx)
	if err != nil {
		return nil, err
	}

	err = pkgconn.Notify(ctx, h.modelClient, entity, false)
	if err != nil {
		return nil, err
	}

	return model.ExposeConnector(entity), nil
}

func (h Handler) Get(ctx *gin.Context, req view.GetRequest) (view.UpdateResponse, error) {
	var entity, err = h.modelClient.Connectors().Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return model.ExposeConnector(entity), nil
}

// Batch APIs

var (
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
	entities, err := query.
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeConnectors(entities), cnt, nil
}

// Extensional APIs

func (h Handler) RouteApplyCostTools(ctx *gin.Context, req view.ApplyCostToolsRequest) error {
	o, err := h.modelClient.Connectors().Get(ctx, req.ID)
	if err != nil {
		return err
	}

	o.FinOpsStatus = status.Initializing
	o.StatusMessage = ""
	update, err := dao.ConnectorUpdate(h.modelClient, o)
	if err != nil {
		return err
	}
	if err = update.Exec(ctx); err != nil {
		return err
	}

	return pkgconn.Notify(ctx, h.modelClient, o, true)
}

func (h Handler) RouteSyncCostOpsData(ctx *gin.Context, req view.SyncCostDataRequest) error {
	o, err := h.modelClient.Connectors().Get(ctx, req.ID)
	if err != nil {
		return err
	}

	task, err := scheduler.NewCostSyncTask(h.modelClient)
	if err != nil {
		return err
	}

	return task.SyncK8sCost(ctx, o, nil)
}
