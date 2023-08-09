package catalog

import (
	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/catalog/view"
	buscatalog "github.com/seal-io/seal/pkg/bus/catalog"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/catalog"
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
	return "Catalog"
}

func (h Handler) Validating() any {
	return h.modelClient
}

// Basic APIs.

func (h Handler) Create(ctx *gin.Context, req view.CreateRequest) (view.CreateResponse, error) {
	entity := req.Model()

	status.CatalogStatusInitialized.Unknown(entity, "")
	entity.Status.SetSummary(status.WalkCatalog(&entity.Status))

	entity, err := h.modelClient.Catalogs().Create().
		Set(entity).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	err = buscatalog.Notify(ctx, h.modelClient, entity)
	if err != nil {
		return nil, err
	}

	return model.ExposeCatalog(entity), nil
}

func (h Handler) Update(ctx *gin.Context, req view.UpdateRequest) error {
	entity := req.Model()

	status.CatalogStatusInitialized.Reset(entity, "Initializing catalog template")
	entity.Status.SetSummary(status.WalkCatalog(&entity.Status))

	entity, err := h.modelClient.Catalogs().UpdateOne(entity).
		Set(entity).
		Save(ctx)
	if err != nil {
		return err
	}

	return buscatalog.Notify(ctx, h.modelClient, entity)
}

func (h Handler) Get(ctx *gin.Context, req view.GetRequest) (view.GetResponse, error) {
	entity, err := h.modelClient.Catalogs().Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return model.ExposeCatalog(entity), nil
}

func (h Handler) Delete(ctx *gin.Context, req view.DeleteRequest) error {
	return h.modelClient.Catalogs().DeleteOne(req.Model()).Exec(ctx)
}

// Batch APIs.

var (
	queryFields = []string{
		catalog.FieldName,
		catalog.FieldType,
	}
	getFields  = catalog.Columns
	sortFields = []string{
		catalog.FieldID,
		catalog.FieldName,
		catalog.FieldType,
		catalog.FieldSource,
	}
)

func (h Handler) CollectionGet(
	ctx *gin.Context,
	req view.CollectionGetRequest,
) (view.CollectionGetResponse, int, error) {
	query := h.modelClient.Catalogs().Query()
	if queries, ok := req.Querying(queryFields); ok {
		query = query.Where(queries)
	}

	if sorts, ok := req.Sorting(sortFields); ok {
		query = query.Order(sorts...)
	}

	// Get count.
	count, err := query.Clone().Count(ctx)
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

	if orders, ok := req.Sorting(sortFields, model.Desc(catalog.FieldCreateTime)); ok {
		query.Order(orders...)
	}

	entities, err := query.
		// Allow returning without sorting keys.
		Unique(false).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeCatalogs(entities), count, nil
}

func (h Handler) CollectionDelete(ctx *gin.Context, req view.CollectionDeleteRequest) error {
	return h.modelClient.WithTx(ctx, func(tx *model.Tx) (err error) {
		for i := range req {
			err = tx.Catalogs().DeleteOne(req[i].Model()).
				Exec(ctx)
			if err != nil {
				return err
			}
		}

		return
	})
}

// Extensional APIs.

func (h Handler) RouteRefresh(ctx *gin.Context, req view.RouteSyncCatalogRequest) error {
	catalog, err := h.modelClient.Catalogs().Get(ctx, req.ID)
	if err != nil {
		return err
	}

	status.CatalogStatusInitialized.Unknown(catalog, "Initializing catalog templates")
	catalog.Status.SetSummary(status.WalkCatalog(&catalog.Status))

	catalog, err = h.modelClient.Catalogs().UpdateOne(catalog).
		SetStatus(catalog.Status).
		Save(ctx)
	if err != nil {
		return err
	}

	return buscatalog.Notify(ctx, h.modelClient, catalog)
}
