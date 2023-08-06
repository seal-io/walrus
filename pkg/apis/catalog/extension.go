package catalog

import (
	catalogbus "github.com/seal-io/seal/pkg/bus/catalog"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/template"
	"github.com/seal-io/seal/pkg/dao/types/status"
)

func (h Handler) RouteRefresh(req RouteSyncCatalogRequest) error {
	catalog, err := h.modelClient.Catalogs().
		Get(req.Context, req.ID)
	if err != nil {
		return err
	}

	status.CatalogStatusInitialized.Unknown(catalog, "Initializing catalog templates")
	catalog.Status.SetSummary(status.WalkCatalog(&catalog.Status))

	catalog, err = h.modelClient.Catalogs().UpdateOne(catalog).
		SetStatus(catalog.Status).
		Save(req.Context)
	if err != nil {
		return err
	}

	return catalogbus.Notify(req.Context, h.modelClient, catalog)
}

var (
	queryTemplatesFields = []string{
		template.FieldID,
	}
	getTemplatesFields  = template.WithoutFields()
	sortTemplatesFields = []string{
		template.FieldID,
		template.FieldStatus,
		template.FieldCreateTime,
	}
)

func (h Handler) RouteGetTemplates(req RouteGetTemplatesRequest) (RouteGetTemplatesResponse, int, error) {
	query := h.modelClient.Templates().Query().
		Where(template.CatalogID(req.ID))

	if queries, ok := req.Querying(queryTemplatesFields); ok {
		query.Where(queries)
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

	if fields, ok := req.Extracting(getTemplatesFields, getTemplatesFields...); ok {
		query.Select(fields...)
	}

	if orders, ok := req.Sorting(sortTemplatesFields, model.Desc(template.FieldCreateTime)); ok {
		query.Order(orders...)
	}

	entities, err := query.
		Unique(false).
		All(req.Context)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeTemplates(entities), cnt, nil
}
