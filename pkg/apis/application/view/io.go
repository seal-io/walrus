package view

import (
	"context"
	"errors"
	"net/http"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/platform/operator"
	"github.com/seal-io/seal/utils/json"
)

// Basic APIs

type CreateRequest struct {
	Name          string                    `json:"name"`
	Description   string                    `json:"description,omitempty"`
	Labels        map[string]string         `json:"labels,omitempty"`
	ProjectID     types.ID                  `json:"projectID"`
	EnvironmentID types.ID                  `json:"environmentID"`
	Modules       []types.ApplicationModule `json:"modules,omitempty"`
}

func (r CreateRequest) Validate() error {
	if r.ProjectID == "" {
		return errors.New("invalid project id: blank")
	}
	if r.EnvironmentID == "" {
		return errors.New("invalid environment id: blank")
	}
	if r.Name == "" {
		return errors.New("invalid name: blank")
	}
	if len(r.Modules) != 0 {
		for i := 0; i < len(r.Modules); i++ {
			if r.Modules[i].ModuleID == "" {
				return errors.New("invalid module id: blank")
			}
			if r.Modules[i].Name == "" {
				return errors.New("invalid module name: blank")
			}
		}
	}
	return nil
}

func (r CreateRequest) Model() *model.Application {
	var input = &model.Application{
		Name:          r.Name,
		Description:   r.Description,
		Labels:        r.Labels,
		ProjectID:     r.ProjectID,
		EnvironmentID: r.EnvironmentID,
	}
	for _, m := range r.Modules {
		input.Edges.Modules = append(input.Edges.Modules,
			&model.ApplicationModuleRelationship{
				ModuleID:  m.ModuleID,
				Name:      m.Name,
				Variables: m.Variables,
			})
	}
	return input
}

type CreateResponse = GetResponse

type DeleteRequest = GetRequest

type UpdateRequest struct {
	ID          types.ID                  `uri:"id"`
	Name        string                    `json:"name,omitempty"`
	Description string                    `json:"description,omitempty"`
	Labels      map[string]string         `json:"labels,omitempty"`
	Modules     []types.ApplicationModule `json:"modules,omitempty"`
}

func (r UpdateRequest) Validate() error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}
	if len(r.Modules) != 0 {
		for i := 0; i < len(r.Modules); i++ {
			if r.Modules[i].ModuleID == "" {
				return errors.New("invalid module id: blank")
			}
			if r.Modules[i].Name == "" {
				return errors.New("invalid module name: blank")
			}
		}
	}
	return nil
}

func (r UpdateRequest) Model() *model.Application {
	var input = &model.Application{
		Name:        r.Name,
		Description: r.Description,
		Labels:      r.Labels,
	}
	for _, m := range r.Modules {
		input.Edges.Modules = append(input.Edges.Modules,
			&model.ApplicationModuleRelationship{
				ModuleID:  m.ModuleID,
				Name:      m.Name,
				Variables: m.Variables,
			})
	}
	return input
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

type GetResponse struct {
	*model.Application `json:",inline"`

	Modules []types.ApplicationModule `json:"modules,omitempty"`
}

func (r GetResponse) MarshalJSON() ([]byte, error) {
	type Alias GetResponse

	// mutate `.Edges.ApplicationModuleRelationships` to `.Modules`.
	if len(r.Edges.Modules) != 0 {
		for _, s := range r.Edges.Modules {
			if s == nil {
				continue
			}
			r.Modules = append(r.Modules,
				types.ApplicationModule{
					ModuleID:  s.ModuleID,
					Name:      s.Name,
					Variables: s.Variables,
				})
		}
		r.Edges.Modules = nil // release
	}

	return json.Marshal(&struct {
		*Alias `json:",inline"`
	}{
		Alias: (*Alias)(&r),
	})
}

// Batch APIs

// Extensional APIs

type GetResourcesRequest struct {
	runtime.RequestPagination `query:",inline"`
	runtime.RequestSorting    `query:",inline"`

	ID          types.ID `uri:"id"`
	WithoutKeys bool     `query:"withoutKeys,omitempty"`
}

func (r GetResourcesRequest) ValidateWith(ctx context.Context, input any) error {
	var modelClient = input.(model.ClientSet)

	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}
	_, err := modelClient.Applications().Query().
		Where(application.ID(r.ID)).
		OnlyID(ctx)
	if err != nil {
		return runtime.Error(http.StatusNotFound, "invalid id: not found")
	}

	return nil
}

type GetResourceResponse struct {
	*model.ApplicationResource `json:",inline"`

	OperatorKeys  *operator.Keys `json:"operatorKeys"`
	ConnectorName string         `json:"connectorName"`
}

func (o *GetResourceResponse) MarshalJSON() ([]byte, error) {
	type Alias GetResourceResponse

	// move `.Edges.Connector.Name` to `.ConnectorName`.
	if o.Edges.Connector != nil {
		o.ConnectorName = o.Edges.Connector.Name
		o.Edges.Connector = nil // release
	}

	return json.Marshal(&struct {
		*Alias `json:",inline"`
	}{
		Alias: (*Alias)(o),
	})
}

type GetResourcesResponse = []GetResourceResponse

type GetRevisionsRequest struct {
	runtime.RequestPagination `query:",inline"`
	runtime.RequestSorting    `query:",inline"`

	ID types.ID `uri:"id"`
}

func (r GetRevisionsRequest) ValidateWith(ctx context.Context, input any) error {
	var modelClient = input.(model.ClientSet)

	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}
	_, err := modelClient.Applications().Query().
		Where(application.ID(r.ID)).
		OnlyID(ctx)
	if err != nil {
		return runtime.Error(http.StatusNotFound, "invalid id: not found")
	}

	return nil
}

type GetRevisionsResponse = []*model.ApplicationRevision
