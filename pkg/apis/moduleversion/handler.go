package moduleversion

import (
	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/moduleversion/view"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/moduleversion"
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
	return "ModuleVersion"
}

func (h Handler) Validating() any {
	return h.modelClient
}

// Basic APIs

func (h Handler) Get(ctx *gin.Context, req view.GetRequest) (view.GetResponse, error) {
	var entity, err = h.modelClient.ModuleVersions().Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return model.ExposeModuleVersion(entity), nil
}

// Batch APIs

var (
	queryFields = []string{
		moduleversion.FieldVersion,
	}
	getFields  = moduleversion.Columns
	sortFields = []string{
		moduleversion.FieldVersion}
)

func (h Handler) CollectionGet(ctx *gin.Context, req view.CollectionGetRequest) (view.CollectionGetResponse, int, error) {
	var query = h.modelClient.ModuleVersions().Query().
		Where(moduleversion.ModuleIDIn(req.ModuleIDs...))
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
	if orders, ok := req.Sorting(sortFields, model.Desc(moduleversion.FieldCreateTime)); ok {
		query.Order(orders...)
	}
	entities, err := query.
		// allow returning without sorting keys.
		Unique(false).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeModuleVersions(entities), cnt, nil
}

// Extensional APIs
