package view

import (
	"context"
	"errors"
	"net/http"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/project"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/utils/json"
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

func (r CreateRequest) Model() *model.Project {
	return &model.Project{
		Name:        r.Name,
		Description: r.Description,
		Labels:      r.Labels,
	}
}

type CreateResponse = model.Project

type DeleteRequest = GetRequest

type UpdateRequest struct {
	ID          types.ID          `uri:"id"`
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
}

func (r UpdateRequest) Validate() error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}
	return nil
}

func (r UpdateRequest) Model() *model.Project {
	return &model.Project{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Labels:      r.Labels,
	}
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

func (r GetApplicationsRequest) ValidateWith(ctx context.Context, input any) error {
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

type GetApplicationResponse struct {
	*model.Application `json:",inline"`

	EnvironmentName string `json:"environmentName"`
}

func (o *GetApplicationResponse) MarshalJSON() ([]byte, error) {
	type Alias GetApplicationResponse

	// move `.Edges.Environment.Name` to `.EnvironmentName`.
	if o.Edges.Environment != nil {
		o.EnvironmentName = o.Edges.Environment.Name
		o.Edges.Environment = nil // release
	}

	return json.Marshal(&struct {
		*Alias `json:",inline"`
	}{
		Alias: (*Alias)(o),
	})
}

type GetApplicationsResponse = []GetApplicationResponse
