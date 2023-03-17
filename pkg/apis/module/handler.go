package module

import (
	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/module/view"
	modbus "github.com/seal-io/seal/pkg/bus/module"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/module"
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
	return "Module"
}

func (h Handler) Validating() any {
	return h.modelClient
}

// Basic APIs

func (h Handler) Create(ctx *gin.Context, req view.CreateRequest) (view.CreateResponse, error) {
	var entity = req.Model()

	var creates, err = dao.ModuleCreates(h.modelClient, entity)
	if err != nil {
		return nil, err
	}
	entity, err = creates[0].Save(ctx)
	if err != nil {
		return nil, err
	}

	if err = modbus.Notify(ctx, entity); err != nil {
		return nil, err
	}

	return model.ExposeModule(entity), nil
}

func (h Handler) Delete(ctx *gin.Context, req view.DeleteRequest) error {
	return h.modelClient.Modules().DeleteOne(req.Model()).Exec(ctx)
}

func (h Handler) Update(ctx *gin.Context, req view.UpdateRequest) error {
	prev, err := h.modelClient.Modules().Get(ctx, req.ID)
	if err != nil {
		return err
	}

	var (
		entity = req.Model()
		// sync schema on source/version updates
		shouldSyncSchema = prev.Source != entity.Source
	)

	if shouldSyncSchema {
		entity.Status = status.ModuleStatusInitializing
		entity.StatusMessage = ""
	}

	update, err := dao.ModuleUpdate(h.modelClient, entity)
	if err != nil {
		return err
	}
	if _, err = update.Save(ctx); err != nil {
		return err
	}

	if !shouldSyncSchema {
		return nil
	}

	return modbus.Notify(ctx, entity)
}

func (h Handler) Get(ctx *gin.Context, req view.GetRequest) (view.GetResponse, error) {
	var entity, err = h.modelClient.Modules().Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return model.ExposeModule(entity), nil
}

// Batch APIs

func (h Handler) CollectionDelete(ctx *gin.Context, req view.CollectionDeleteRequest) error {
	return h.modelClient.WithTx(ctx, func(tx *model.Tx) (err error) {
		for i := range req {
			err = tx.Modules().DeleteOne(req[i].Model()).
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
		module.FieldID,
	}
	getFields  = module.Columns
	sortFields = []string{
		module.FieldID,
		module.FieldStatus,
		module.FieldCreateTime}
)

func (h Handler) CollectionGet(ctx *gin.Context, req view.CollectionGetRequest) (view.CollectionGetResponse, int, error) {
	var query = h.modelClient.Modules().Query()
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
	if orders, ok := req.Sorting(sortFields, model.Desc(module.FieldCreateTime)); ok {
		query.Order(orders...)
	}
	entities, err := query.
		// allow returning without sorting keys.
		Unique(false).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeModules(entities), cnt, nil
}

// Extensional APIs

func (h Handler) RouteRefresh(ctx *gin.Context, req view.RefreshRequest) error {
	m, err := h.modelClient.Modules().Get(ctx, req.ID)
	if err != nil {
		return err
	}
	m.Status = status.ModuleStatusInitializing
	m.StatusMessage = ""
	update, err := dao.ModuleUpdate(h.modelClient, m)
	if err != nil {
		return err
	}
	if err = update.Exec(ctx); err != nil {
		return err
	}

	return modbus.Notify(ctx, m)
}
