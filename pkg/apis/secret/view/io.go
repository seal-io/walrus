package view

import (
	"context"
	"errors"
	"net/http"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/project"
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
	runtime.RequestPagination `query:",inline"`

	ProjectID *types.ID `query:"projectID,omitempty"`
}

func (r *CollectionGetRequest) ValidateWith(ctx context.Context, input any) error {
	var modelClient = input.(model.ClientSet)

	if r.ProjectID != nil {
		if !r.ProjectID.Valid(0) {
			return errors.New("invalid project id: blank")
		}
		_, err := modelClient.Projects().Query().
			Where(project.ID(*r.ProjectID)).
			OnlyID(ctx)
		if err != nil {
			return runtime.Error(http.StatusNotFound, "invalid project id: not found")
		}
	}
	return nil
}

type CollectionGetResponse = []*model.SecretOutput

// Extensional APIs
