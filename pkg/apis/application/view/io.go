package view

import (
	"context"
	"errors"
	"net/http"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/platform/operator"
)

// Basic APIs

type CreateRequest struct {
	*model.ApplicationCreateInput `json:",inline"`
}

func (r *CreateRequest) Validate() error {
	if r.Project.ID == "" {
		return errors.New("invalid project id: blank")
	}
	if r.Environment.ID == "" {
		return errors.New("invalid environment id: blank")
	}
	if r.Name == "" {
		return errors.New("invalid name: blank")
	}
	if len(r.Modules) != 0 {
		for i := 0; i < len(r.Modules); i++ {
			if r.Modules[i].Module.ID == "" {
				return errors.New("invalid module id: blank")
			}
			if r.Modules[i].Name == "" {
				return errors.New("invalid module name: blank")
			}
		}
	}
	return nil
}

type CreateResponse = *model.ApplicationOutput

type DeleteRequest = GetRequest

type UpdateRequest struct {
	*model.ApplicationUpdateInput `uri:",inline" json:",inline"`
}

func (r *UpdateRequest) Validate() error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}
	if len(r.Modules) != 0 {
		for i := 0; i < len(r.Modules); i++ {
			if r.Modules[i].Module.ID == "" {
				return errors.New("invalid module id: blank")
			}
			if r.Modules[i].Name == "" {
				return errors.New("invalid module name: blank")
			}
		}
	}
	return nil
}

type GetRequest struct {
	*model.ApplicationQueryInput `uri:",inline"`
}

func (r *GetRequest) Validate() error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	return nil
}

type GetResponse = *model.ApplicationOutput

// Batch APIs

// Extensional APIs

type GetResourcesRequest struct {
	*model.ApplicationQueryInput `uri:",inline"`
	runtime.RequestPagination    `query:",inline"`
	runtime.RequestExtracting    `query:",inline"`
	runtime.RequestSorting       `query:",inline"`

	WithoutKeys bool `query:"withoutKeys,omitempty"`
}

func (r *GetResourcesRequest) ValidateWith(ctx context.Context, input any) error {
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

type ApplicationResource struct {
	*model.ApplicationResourceOutput `json:",inline"`

	OperatorKeys  *operator.Keys `json:"operatorKeys"`
	ConnectorName string         `json:"connectorName"`
}

type GetResourcesResponse = []ApplicationResource

type GetRevisionsRequest struct {
	*model.ApplicationQueryInput `uri:",inline"`
	runtime.RequestPagination    `query:",inline"`
	runtime.RequestExtracting    `query:",inline"`
	runtime.RequestSorting       `query:",inline"`
}

func (r *GetRevisionsRequest) ValidateWith(ctx context.Context, input any) error {
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

type GetRevisionsResponse = []*model.ApplicationRevisionOutput
