package token

import (
	"net/http"

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

// Basic APIs

func (h Handler) Create(ctx *gin.Context, req view.CreateRequest) (*view.CreateResponse, error) {
	var input = &model.Token{
		Name:       req.Name,
		Expiration: req.Expiration,
	}
	var s = session.LoadSubject(ctx)
	var cred casdoor.ApplicationCredential
	var err = settings.CasdoorCred.ValueJSONUnmarshal(ctx, h.modelClient, &cred)
	if err != nil {
		return nil, err
	}
	// create token value from casdoor.
	t, err := casdoor.CreateToken(ctx, cred.ClientID, cred.ClientSecret, s.Name, req.Expiration)
	if err != nil {
		return nil, runtime.ErrorfP(http.StatusBadRequest, "failed to create token to casdoor: %w", err)
	}
	input.CasdoorTokenName, input.CasdoorTokenOwner = t.Name, t.Owner

	// create token.
	var cerr error
	defer func() {
		// revert token if any error occurs.
		if cerr == nil {
			return
		}
		_ = casdoor.DeleteToken(ctx, cred.ClientID, cred.ClientSecret, t.Owner, t.Name)
	}()
	creates, cerr := dao.TokenCreates(h.modelClient, input)
	if cerr != nil {
		return nil, cerr
	}
	entity, cerr := creates[0].Save(ctx)
	if cerr != nil {
		return nil, cerr
	}

	var resp = entity
	resp.AccessToken = t.AccessToken
	return resp, nil
}

func (h Handler) Delete(ctx *gin.Context, req view.DeleteRequest) error {
	var input []predicate.Token
	if req.ID.IsNaive() {
		input = append(input, token.ID(req.ID))
	} else {
		var keys = req.ID.Split()
		input = append(input, token.CasdoorTokenName(keys[0]))
	}

	// delete token.
	var entity, err = h.modelClient.Tokens().Query().
		Where(input...).
		Select(token.FieldCasdoorTokenOwner, token.FieldCasdoorTokenName).
		Only(ctx)
	if err != nil {
		if model.IsNotFound(err) {
			return runtime.Error(http.StatusBadRequest, "invalid token: not found")
		}
		return runtime.ErrorfP(http.StatusInternalServerError, "failed to get requesting token: %w", err)
	}
	err = h.modelClient.WithTx(ctx, func(tx *model.Tx) error {
		var _, err = tx.Tokens().Delete().
			Where(input...).
			Exec(ctx)
		if err != nil {
			return err
		}
		// remove token value from casdoor.
		var cred casdoor.ApplicationCredential
		err = settings.CasdoorCred.ValueJSONUnmarshal(ctx, h.modelClient, &cred)
		if err != nil {
			return err
		}
		err = casdoor.DeleteToken(ctx, cred.ClientID, cred.ClientSecret, entity.CasdoorTokenOwner, entity.CasdoorTokenName)
		if err != nil {
			return runtime.ErrorfP(http.StatusBadRequest, "failed to delete token from casdoor: %w", err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// clean cache.
	cache.CleanTokenSubjects()
	return nil
}

// Batch APIs

func (h Handler) CollectionGet(ctx *gin.Context, _ view.CollectionGetRequest) (view.CollectionGetResponse, int, error) {
	var query = h.modelClient.Tokens().Query()

	var entities, err = query.Select(token.WithoutFields(
		token.FieldUpdateTime, token.FieldCasdoorTokenName, token.FieldCasdoorTokenOwner)...).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return entities, len(entities), nil
}

// Extensional APIs
