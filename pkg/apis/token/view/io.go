package view

import (
	"errors"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
)

// Basic APIs

type CreateRequest struct {
	*model.TokenCreateInput `json:",inline"`
}

func (r *CreateRequest) Validate() error {
	if r.Name == "" {
		return errors.New("invalid name: blank")
	}
	if r.Expiration != nil {
		if *r.Expiration < 0 {
			return errors.New("invalid expiration: negative")
		}
	}
	return nil
}

type CreateResponse struct {
	*model.TokenOutput `json:",inline"`

	// AccessToken is the token used for authentication.
	// It does not store in the database and only shows up after created.
	AccessToken string `json:"accessToken,omitempty"`
}

type DeleteRequest = GetRequest

type GetRequest struct {
	*model.TokenQueryInput `uri:",inline"`
}

func (r *GetRequest) Validate() error {
	if !r.ID.Valid(1) {
		return errors.New("invalid id: blank")
	}
	return nil
}

type GetResponse = *model.TokenOutput

// Batch APIs

type CollectionGetRequest struct {
	runtime.RequestQuerying[predicate.Token] `query:",inline"`
	runtime.RequestPagination                `query:",inline"`
}

type CollectionGetResponse = []*model.TokenOutput

// Extensional APIs
