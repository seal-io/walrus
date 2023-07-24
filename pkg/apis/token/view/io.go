package view

import (
	"errors"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
)

// Basic APIs.

type CreateRequest struct {
	Name              string `json:"name"`
	ExpirationSeconds *int   `json:"expirationSeconds,omitempty"`
}

func (r *CreateRequest) Validate() error {
	if r.Name == "" {
		return errors.New("invalid name: blank")
	}

	if r.ExpirationSeconds != nil {
		if *r.ExpirationSeconds < 0 {
			return errors.New("invalid expiration seconds: negative")
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

type DeleteRequest struct {
	model.TokenQueryInput `uri:",inline"`
}

func (r *DeleteRequest) Validate() error {
	if !r.ID.Valid() {
		return errors.New("invalid id: blank")
	}

	return nil
}

// Batch APIs.

type CollectionGetRequest struct {
	runtime.RequestQuerying[predicate.Token] `query:",inline"`
	runtime.RequestPagination                `query:",inline"`
}

type CollectionGetResponse = []*model.TokenOutput

// Extensional APIs.
