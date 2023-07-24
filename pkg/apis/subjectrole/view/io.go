package view

import (
	"errors"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/subjectrolerelationship"
	"github.com/seal-io/seal/pkg/dao/types/object"
)

// Basic APIs.

type CreateRequest struct {
	model.SubjectRoleRelationshipCreateInput `json:",inline"`

	ProjectID object.ID `query:"projectID,omitempty"`
}

func (r *CreateRequest) Validate() error {
	if r.ProjectID != "" && !r.ProjectID.Valid(0) {
		return errors.New("invalid project id: blank")
	}

	if !r.Subject.ID.Valid(0) {
		return errors.New("invalid subject id: blank")
	}

	if r.Role.ID == "" {
		return errors.New("invalid role id: blank")
	}

	return nil
}

type CreateResponse = *model.SubjectRoleRelationshipOutput

type DeleteRequest struct {
	model.SubjectRoleRelationshipQueryInput `uri:",inline"`

	ProjectID object.ID `query:"projectID,omitempty"`
}

func (r *DeleteRequest) Validate() error {
	if r.ProjectID != "" && !r.ProjectID.Valid(0) {
		return errors.New("invalid project id: blank")
	}

	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	return nil
}

// Batch APIs.

type CollectionDeleteRequest []*model.SubjectRoleRelationshipQueryInput

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
	runtime.RequestPagination                                   `query:",inline"`
	runtime.RequestSorting[subjectrolerelationship.OrderOption] `query:",inline"`

	ProjectID object.ID `query:"projectID,omitempty"`
}

func (r *CollectionGetRequest) Validate() error {
	// Query global scope subject roles if the given `ProjectID` is empty,
	// otherwise, query project scope subject roles.
	if r.ProjectID != "" && !r.ProjectID.Valid(0) {
		return errors.New("invalid project id")
	}

	return nil
}

type CollectionGetResponse = []*model.SubjectRoleRelationshipOutput

// Extensional APIs.
