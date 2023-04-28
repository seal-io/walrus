package view

import (
	"errors"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/secret"
	"github.com/seal-io/seal/pkg/dao/types"
)

// Basic APIs

type CreateRequest struct {
	*model.SecretCreateInput `json:",inline"`
}

func (r *CreateRequest) Validate() error {
	if r.Project != nil && !r.Project.ID.Valid(0) {
		return errors.New("invalid project id: blank")
	}
	if r.Name == "" {
		return errors.New("invalid name: blank")
	}
	if r.Value == "" {
		return errors.New("invalid value: blank")
	}
	return nil
}

type CreateResponse = *model.SecretOutput

type DeleteRequest struct {
	*model.SecretQueryInput `uri:",inline"`
}

type UpdateRequest struct {
	*model.SecretUpdateInput `uri:",inline" json:",inline"`
}

func (r *UpdateRequest) Validate() error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}
	if r.Value == "" {
		return errors.New("invalid value: blank")
	}
	return nil
}

// Batch APIs

type CollectionDeleteRequest []*model.SecretQueryInput

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
	runtime.RequestCollection[predicate.Secret, secret.OrderOption] `query:",inline"`

	ProjectIDs []types.ID `query:"projectID,omitempty"`
}

func (r *CollectionGetRequest) Validate() error {
	// query global scope secret if the given `ProjectIDs` is empty,
	// otherwise, query project scope secret.
	for i := range r.ProjectIDs {
		if !r.ProjectIDs[i].Valid(0) {
			return errors.New("invalid project id: blank")
		}
	}
	return nil
}

type CollectionGetResponse = []*model.SecretOutput

// Extensional APIs
