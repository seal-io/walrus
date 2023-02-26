package environment

import (
	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/environment/view"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/environment"
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
	return "Environment"
}

func (h Handler) Validating() any {
	return h.modelClient
}

// Basic APIs

func (h Handler) Create(ctx *gin.Context, req view.EnvironmentCreateRequest) (*view.EnvironmentCreateResponse, error) {
	var creates, err = dao.EnvironmentCreates(h.modelClient, req.Model())
	if err != nil {
		return nil, err
	}

	result, err := creates[0].Save(ctx)
	if err != nil {
		return nil, err
	}

	return &view.EnvironmentCreateResponse{Environment: result}, nil
}

func (h Handler) Delete(ctx *gin.Context, req view.EnvironmentDeleteRequest) error {
	return h.modelClient.Environments().DeleteOneID(req.ID).Exec(ctx)
}

func (h Handler) Update(ctx *gin.Context, req view.EnvironmentUpdateRequest) error {
	var updates, err = dao.EnvironmentUpdates(h.modelClient, req.EnvironmentVO.Model())
	if err != nil {
		return err
	}
	return updates[0].Exec(ctx)
}

func (h Handler) Get(ctx *gin.Context, req view.EnvironmentGetRequest) (*view.EnvironmentGetResponse, error) {
	env, err := h.modelClient.Environments().Query().Where(environment.ID(req.ID)).WithConnectors().Only(ctx)
	if err != nil {
		return nil, err
	}
	return &view.EnvironmentGetResponse{Environment: env}, nil
}

// Batch APIs

var (
	getFields  = environment.Columns
	sortFields = []string{environment.FieldName, environment.FieldCreateTime, environment.FieldUpdateTime}
)

func (h Handler) CollectionGet(ctx *gin.Context, req view.CollectionGetRequest) (view.CollectionGetResponse, int, error) {
	var query = h.modelClient.Environments().Query().WithConnectors()

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

	var result = make([]*view.EnvironmentVO, len(entities))
	for i, e := range entities {
		result[i] = &view.EnvironmentVO{Environment: e}
	}

	return result, cnt, nil
}

// Extensional APIs
