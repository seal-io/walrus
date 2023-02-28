package module

import (
	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/module/view"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/module"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/modules"
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
	var creates, err = dao.ModuleCreates(h.modelClient, req.Module)
	if err != nil {
		return nil, err
	}

	m, err := creates[0].Save(ctx)
	if err != nil {
		return nil, err
	}

	modules.Notify(ctx, h.modelClient, req.Module)

	return m, nil
}

func (h Handler) Delete(ctx *gin.Context, req view.DeleteRequest) error {
	return h.modelClient.Modules().DeleteOneID(req.ID).Exec(ctx)
}

func (h Handler) Update(ctx *gin.Context, req view.UpdateRequest) error {
	prev, err := h.modelClient.Modules().Get(ctx, req.ID)
	if err != nil {
		return err
	}

	var (
		module = req.Module
		// Sync schema on source/version updates
		shouldSyncSchema = prev.Source != req.Source || prev.Version != req.Version
	)

	if shouldSyncSchema {
		module.Schema = nil
		module.Status = status.Initializing
		module.StatusMessage = ""
	}

	update, err := dao.ModuleUpdate(h.modelClient, req.Module)
	if err != nil {
		return err
	}

	if _, err = update.Save(ctx); err != nil {
		return err
	}

	if shouldSyncSchema {
		modules.Notify(ctx, h.modelClient, module)
	}

	return nil
}

func (h Handler) Get(ctx *gin.Context, req view.GetRequest) (view.GetResponse, error) {
	return h.modelClient.Modules().Get(ctx, req.ID)
}

// Batch APIs

var (
	getFields  = module.Columns
	sortFields = []string{module.FieldID, module.FieldCreateTime, module.FieldUpdateTime}
)

func (h Handler) CollectionGet(ctx *gin.Context, req view.CollectionGetRequest) (view.CollectionGetResponse, int, error) {

	var query = h.modelClient.Modules().Query()

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
	if orders, ok := req.Sorting(sortFields); ok {
		query.Order(orders...)
	}
	entities, err := query.
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return entities, cnt, nil
}

// Extensional APIs

func (h Handler) RouteRefresh(ctx *gin.Context, req view.RefreshRequest) error {
	m, err := h.modelClient.Modules().Get(ctx, req.ID)
	if err != nil {
		return err
	}
	m.Schema = nil
	m.Status = status.Initializing
	m.StatusMessage = ""
	update, err := dao.ModuleUpdate(h.modelClient, m)
	if err != nil {
		return err
	}
	if err = update.Exec(ctx); err != nil {
		return err
	}

	modules.Notify(ctx, h.modelClient, m)

	return nil
}
