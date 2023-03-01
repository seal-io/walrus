package view

import (
	"context"
	"errors"
	"net/http"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/project"
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

type CollectionGetRequest struct {
	runtime.RequestPagination `query:",inline"`
	runtime.RequestSorting    `query:",inline"`
}

type CollectionGetResponse = []*model.ProjectOutput

// Extensional APIs

type GetApplicationsRequest struct {
	*model.ProjectQueryInput  `uri:",inline"`
	runtime.RequestPagination `query:",inline"`
	runtime.RequestExtracting `query:",inline"`
	runtime.RequestSorting    `query:",inline"`
}

func (r *GetApplicationsRequest) ValidateWith(ctx context.Context, input any) error {
	var modelClient = input.(model.ClientSet)

	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}
	_, err := modelClient.Projects().Query().
		Where(project.ID(r.ID)).
		OnlyID(ctx)
	if err != nil {
		return runtime.Error(http.StatusNotFound, "invalid id: not found")
	}
	return nil
}

type GetApplicationsResponse = []*model.ApplicationOutput
