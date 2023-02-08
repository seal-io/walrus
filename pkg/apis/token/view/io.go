package view

import (
	"errors"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/oid"
)

// Basic APIs

type CreateRequest struct {
	Name       string `json:"name"`
	Expiration *int   `json:"expiration,omitempty"`
}

func (r CreateRequest) Validate() error {
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

type CreateResponse = model.Token

type DeleteRequest struct {
	ID oid.ID `uri:"id"`
}

func (r DeleteRequest) Validate() error {
	if !r.ID.Valid(1) {
		return errors.New("invalid id: blank")
	}
	return nil
}

type GetRequest struct {
	ID oid.ID `uri:"id"`
}

func (r GetRequest) Validate() error {
	if !r.ID.Valid(1) {
		return errors.New("invalid id: blank")
	}
	return nil
}

type GetResponse = model.Token

// Batch APIs

type CollectionGetRequest struct {
	runtime.RequestPagination `query:",inline"`
}

type CollectionGetResponse = model.Tokens

// Extensional APIs
