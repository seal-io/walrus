package token

import (
	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/token/view"
	"github.com/seal-io/seal/pkg/auths"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/token"
	"github.com/seal-io/seal/pkg/dao/types"
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
	return "Token"
}

// Basic APIs.

func (h Handler) Create(ctx *gin.Context, req view.CreateRequest) (*view.CreateResponse, error) {
	at, err := auths.CreateAccessToken(ctx,
		h.modelClient, types.TokenKindAPI, req.Name, req.ExpirationSeconds)
	if err != nil {
		return nil, err
	}

	return &view.CreateResponse{
		TokenOutput: model.ExposeToken(at.Raw),
		AccessToken: at.Value,
	}, nil
}

func (h Handler) Delete(ctx *gin.Context, req view.DeleteRequest) error {
	return h.modelClient.Tokens().DeleteOne(req.Model()).Exec(ctx)
}

// Batch APIs.

var (
	queryFields = []string{
		token.FieldName,
	}
	getFields = token.WithoutFields(
		token.FieldValue)
)

func (h Handler) CollectionGet(
	ctx *gin.Context,
	req view.CollectionGetRequest,
) (view.CollectionGetResponse, int, error) {
	query := h.modelClient.Tokens().Query().
		Where(token.Kind(types.TokenKindAPI))
	if queries, ok := req.Querying(queryFields); ok {
		query.Where(queries)
	}

	entities, err := query.
		Order(model.Desc(token.FieldCreateTime)).
		Select(getFields...).
		// Allow returning without sorting keys.
		Unique(false).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeTokens(entities), len(entities), nil
}

// Extensional APIs.
