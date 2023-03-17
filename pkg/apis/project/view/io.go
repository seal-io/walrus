package view

import (
	"errors"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
)

// Basic APIs

type CreateRequest struct {
	*model.ProjectCreateInput `json:",inline"`
}

func (r *CreateRequest) Validate() error {
	if r.Name == "" {
		return errors.New("invalid name: blank")
	}
	return nil
}

type CreateResponse = *model.ProjectOutput

type DeleteRequest = GetRequest

type UpdateRequest struct {
	*model.ProjectUpdateInput `uri:",inline" json:",inline"`
}

func (r *UpdateRequest) Validate() error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}
	return nil
}

type GetRequest struct {
	*model.ProjectQueryInput `uri:",inline"`
}

func (r *GetRequest) Validate() error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}
	return nil
}

type GetResponse = *model.ProjectOutput

// Batch APIs

type CollectionDeleteRequest []*model.ProjectQueryInput

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
	runtime.RequestCollection[predicate.Project] `query:",inline"`
}

type CollectionGetResponse = []*model.ProjectOutput

// Extensional APIs

type GetSecretsRequest struct {
	*model.ProjectQueryInput                    `uri:",inline"`
	runtime.RequestCollection[predicate.Secret] `query:",inline"`
}

func (r *GetSecretsRequest) Validate() error {
	if !r.ID.Valid(0) {
		return errors.New("invalid project id: blank")
	}
	return nil
}

type GetSecretsResponse = []*model.SecretOutput
