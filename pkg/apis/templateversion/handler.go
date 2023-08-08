package templateversion

import (
	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/templateversion/view"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/template"
	"github.com/seal-io/seal/pkg/dao/model/templateversion"
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
	return "TemplateVersion"
}

func (h Handler) Validating() any {
	return h.modelClient
}

// Basic APIs.

func (h Handler) Get(ctx *gin.Context, req view.GetRequest) (view.GetResponse, error) {
	entity, err := h.modelClient.TemplateVersions().Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return model.ExposeTemplateVersion(entity), nil
}

// Batch APIs.

var (
	queryFields = []string{
		templateversion.FieldVersion,
	}
	getFields  = templateversion.Columns
	sortFields = []string{
		templateversion.FieldVersion,
	}
)

func (h Handler) CollectionGet(
	ctx *gin.Context,
	req view.CollectionGetRequest,
) (view.CollectionGetResponse, int, error) {
	query := h.modelClient.TemplateVersions().Query().
		Where(templateversion.TemplateNameIn(req.TemplateNames...)).
		WithTemplate(func(tq *model.TemplateQuery) {
			tq.Select(template.FieldID, template.FieldName)
		})

	if queries, ok := req.Querying(queryFields); ok {
		query.Where(queries)
	}

	// Get count.
	cnt, err := query.Clone().Count(ctx)
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

	if orders, ok := req.Sorting(sortFields, model.Desc(templateversion.FieldCreateTime)); ok {
		query.Order(orders...)
	}

	entities, err := query.
		// Allow returning without sorting keys.
		Unique(false).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeTemplateVersions(entities), cnt, nil
}

// Extensional APIs.
