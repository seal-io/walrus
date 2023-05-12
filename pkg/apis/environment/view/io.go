package view

import (
	"errors"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
)

// Basic APIs.

type CreateRequest struct {
	*model.EnvironmentCreateInput `json:",inline"`
}

func (r *CreateRequest) Validate() error {
	if r.Name == "" {
		return errors.New("invalid name: blank")
	}

	return nil
}

type CreateResponse = *model.EnvironmentOutput

type DeleteRequest = GetRequest

type UpdateRequest struct {
	*model.EnvironmentUpdateInput `uri:",inline" json:",inline"`
}

func (r *UpdateRequest) Validate() error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	return nil
}

type GetRequest struct {
	*model.EnvironmentQueryInput `uri:",inline"`
}

func (r *GetRequest) Validate() error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	return nil
}

type GetResponse = *model.EnvironmentOutput

// Batch APIs.

type CollectionDeleteRequest []*model.EnvironmentQueryInput

func (r CollectionDeleteRequest) Validate() error {
	if len(r) == 0 {
		return errors.New("invalid input: empty")
	}

	for _, i := range r {
		if !i.ID.Valid(0) {
			return errors.New("invalid id: blank")
		}
	}

	return nil
}

type CollectionGetRequest struct {
	runtime.RequestCollection[predicate.Environment, environment.OrderOption] `query:",inline"`
}

type CollectionGetResponse = []*model.EnvironmentOutput

// Extensional APIs.
