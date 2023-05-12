package token

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/auth/cache"
	"github.com/seal-io/seal/pkg/apis/auth/session"
	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/apis/token/view"
	"github.com/seal-io/seal/pkg/casdoor"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/token"
	"github.com/seal-io/seal/pkg/settings"
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
	entity := req.Model()
	s := session.LoadSubject(ctx)

	var cred casdoor.ApplicationCredential

	err := settings.CasdoorCred.ValueJSONUnmarshal(ctx, h.modelClient, &cred)
	if err != nil {
		return nil, err
	}
	// Create token value from casdoor.
	t, err := casdoor.CreateToken(ctx, cred.ClientID, cred.ClientSecret, s.Name, req.Expiration)
	if err != nil {
		return nil, fmt.Errorf("failed to create token to casdoor: %w", err)
	}
	entity.CasdoorTokenName, entity.CasdoorTokenOwner = t.Name, t.Owner

	// Create token.
	var cerr error
	defer func() {
		// Revert token if any error occurs.
		if cerr == nil {
			return
		}
		_ = casdoor.DeleteToken(ctx, cred.ClientID, cred.ClientSecret, t.Owner, t.Name)
	}()

	creates, cerr := dao.TokenCreates(h.modelClient, entity)
	if cerr != nil {
		return nil, cerr
	}

	entity, cerr = creates[0].Save(ctx)
	if cerr != nil {
		return nil, cerr
	}

	return &view.CreateResponse{
		TokenOutput: model.ExposeToken(entity),
		AccessToken: t.AccessToken,
	}, nil
}

func (h Handler) Delete(ctx *gin.Context, req view.DeleteRequest) error {
	var input []predicate.Token
	if req.ID.IsNaive() {
		input = append(input, token.ID(req.ID))
	} else {
		keys := req.ID.Split()
		input = append(input, token.CasdoorTokenName(keys[0]))
	}

	// Delete token.
	entity, err := h.modelClient.Tokens().Query().
		Where(input...).
		Select(token.FieldCasdoorTokenOwner, token.FieldCasdoorTokenName).
		Only(ctx)
	if err != nil {
		return runtime.Errorw(err, "failed to get token")
	}

	err = h.modelClient.WithTx(ctx, func(tx *model.Tx) error {
		_, err := tx.Tokens().Delete().
			Where(input...).
			Exec(ctx)
		if err != nil {
			return err
		}
		// Remove token value from casdoor.
		var cred casdoor.ApplicationCredential

		err = settings.CasdoorCred.ValueJSONUnmarshal(ctx, h.modelClient, &cred)
		if err != nil {
			return err
		}

		err = casdoor.DeleteToken(ctx, cred.ClientID, cred.ClientSecret,
			entity.CasdoorTokenOwner, entity.CasdoorTokenName)
		if err != nil {
			return fmt.Errorf("failed to delete token from casdoor: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	// Clean cache.
	cache.CleanTokenSubjects()

	return nil
}

// Batch APIs.

var (
	queryFields = []string{
		token.FieldName,
	}
	getFields = token.WithoutFields(
		token.FieldUpdateTime,
		token.FieldCasdoorTokenName,
		token.FieldCasdoorTokenOwner)
)

func (h Handler) CollectionGet(
	ctx *gin.Context,
	req view.CollectionGetRequest,
) (view.CollectionGetResponse, int, error) {
	query := h.modelClient.Tokens().Query()
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
