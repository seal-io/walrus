package role

import (
	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/role/view"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/role"
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
	return "Role"
}

// Basic APIs

// Batch APIs

func (h Handler) CollectionGet(ctx *gin.Context, req view.CollectionGetRequest) (view.CollectionGetResponse, int, error) {
	// do not export session level roles.
	var input = []predicate.Role{
		role.Session(false),
	}
	if req.Domain != "" {
		input = append(input, role.Domain(req.Domain))
	}

	var query = h.modelClient.Roles().Query().
		Where(input...)

	// get count.
	cnt, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// get entities.
	if limit, offset, ok := req.Paging(); ok {
		query.Limit(limit).Offset(offset)
	}
	entities, err := query.
		Select(role.WithoutFields(
			role.FieldCreateTime, role.FieldUpdateTime)...).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return entities, cnt, nil
}

// Extensional APIs
