package view

import (
	"errors"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
)

// Basic APIs

type CreateRequest struct {
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
}

func (r CreateRequest) Validate() error {
	if r.Name == "" {
		return errors.New("invalid name: blank")
	}
	return nil
}

type CreateResponse = model.Project

type DeleteRequest = GetRequest

type UpdateRequest struct {
	ID          types.ID          `uri:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
}

func (r UpdateRequest) Validate() error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}
	return nil
}

type GetRequest struct {
	ID types.ID `uri:"id"`
}

func (r GetRequest) Validate() error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	return nil
}

type GetResponse = model.Project

// Batch APIs

type CollectionGetRequest struct {
	runtime.RequestPagination `query:",inline"`
	runtime.RequestSorting    `query:",inline"`
}

type CollectionGetResponse = []*model.Project

// Extensional APIs

type GetApplicationsRequest struct {
	runtime.RequestPagination `query:",inline"`
	runtime.RequestSorting    `query:",inline"`

	ID types.ID `uri:"id"`
}

func (r GetApplicationsRequest) Validate() error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	return nil
}

type GetApplicationResponse struct {
	*model.Application `json:",inline"`

	EnvironmentName string `json:"environmentName"`
}

type GetApplicationsResponse = []GetApplicationResponse
