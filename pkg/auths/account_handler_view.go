package auths

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/auths/session"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/token"
)

type (
	LoginRequest struct {
		_ struct{} `route:"POST=/login"`

		Username string `json:"username"`
		Password string `json:"password"`

		Context *gin.Context
	}

	LoginResponse = *session.Subject
)

func (r *LoginRequest) SetGinContext(ctx *gin.Context) {
	r.Context = ctx
}

func (r *LoginRequest) Validate() error {
	if r.Username == "" {
		return errors.New("invalid username: blank")
	}

	if r.Password == "" {
		return errors.New("invalid password: blank")
	}

	return nil
}

type LogoutRequest struct {
	_ struct{} `route:"DELETE=/logout"`

	Context *gin.Context
}

func (r *LogoutRequest) SetGinContext(ctx *gin.Context) {
	r.Context = ctx
}

type (
	GetInfoRequest struct {
		_ struct{} `route:"GET=/info"`

		Context *gin.Context
	}

	GetInfoResponse = *session.Subject
)

func (r *GetInfoRequest) SetGinContext(ctx *gin.Context) {
	r.Context = ctx
}

type UpdateInfoRequest struct {
	_ struct{} `route:"PUT=/info"`

	Password    string `json:"password,omitempty"`
	OldPassword string `json:"oldPassword,omitempty"`

	Context *gin.Context
}

func (r *UpdateInfoRequest) SetGinContext(ctx *gin.Context) {
	r.Context = ctx
}

func (r *UpdateInfoRequest) Validate() error {
	if r.Password == "" {
		return errors.New("invalid password: blank")
	}

	if r.OldPassword == "" {
		return errors.New("invalid old password: blank")
	}

	if r.Password == r.OldPassword {
		return errors.New("invalid password: the same")
	}

	return nil
}

type (
	CreateTokenRequest struct {
		_ struct{} `route:"POST=/tokens"`

		Name              string `json:"name"`
		ExpirationSeconds *int   `json:"expirationSeconds,omitempty"`

		Context *gin.Context
	}

	CreateTokenResponse = *model.TokenOutput
)

func (r *CreateTokenRequest) SetGinContext(ctx *gin.Context) {
	r.Context = ctx
}

func (r *CreateTokenRequest) Validate() error {
	if r.Name == "" {
		return errors.New("invalid name: blank")
	}

	if r.ExpirationSeconds != nil && *r.ExpirationSeconds < 0 {
		return errors.New("invalid expiration seconds: negative")
	}

	return nil
}

type DeleteTokenRequest struct {
	_ struct{} `route:"DELETE=/tokens/:token"`

	model.TokenDeleteInput `path:",inline"`
}

type (
	GetTokensRequest struct {
		_ struct{} `route:"GET=/tokens"`

		model.TokenQueryInputs `path:",inline" query:",inline"`

		runtime.RequestCollection[
			predicate.Token, token.OrderOption,
		] `query:",inline"`
	}

	GetTokensResponse = []*model.TokenOutput
)
