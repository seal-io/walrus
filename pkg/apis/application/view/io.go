package view

import (
	"context"
	"errors"
	"net/http"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/project"
	"github.com/seal-io/seal/pkg/dao/types"
)

// Basic APIs

type CreateRequest struct {
	*model.ApplicationCreateInput `json:",inline"`
}

func (r *CreateRequest) Validate() error {
	if r.Project.ID == "" {
		return errors.New("invalid project id: blank")
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

type CollectionDeleteRequest []*model.ApplicationQueryInput

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
	runtime.RequestCollection[predicate.Application] `query:",inline"`

	ProjectID types.ID `query:"projectID"`
}

func (r *CollectionGetRequest) ValidateWith(ctx context.Context, input any) error {
	var modelClient = input.(model.ClientSet)

	if !r.ProjectID.Valid(0) {
		return errors.New("invalid project id: blank")
	}
	_, err := modelClient.Projects().Query().
		Where(project.ID(r.ProjectID)).
		OnlyID(ctx)
	if err != nil {
		return runtime.Error(http.StatusNotFound, "invalid project id: not found")
	}
	return nil
}

type CollectionGetResponse = []*model.ApplicationOutput

// Extensional APIs

type RouteDeployRequest struct {
	_ struct{} `route:"POST=/deploy"`

	*model.ApplicationQueryInput          `uri:",inline"`
	*model.ApplicationInstanceCreateInput `json:",inline"`
}

func (r *RouteDeployRequest) ValidateWith(ctx context.Context, input any) error {
	var modelClient = input.(model.ClientSet)

	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}
	if !r.Environment.ID.Valid(0) {
		return errors.New("invalid environment id: blank")
	}
	if r.Name == "" {
		return errors.New("invalid name: blank")
	}

	_, err := modelClient.Applications().Query().
		Where(application.ID(r.ID)).
		OnlyID(ctx)
	if err != nil {
		return runtime.Error(http.StatusNotFound, "invalid id: not found")
	}
	_, err = modelClient.Environments().Query().
		Where(environment.ID(r.Environment.ID)).
		OnlyID(ctx)
	if err != nil {
		return runtime.Error(http.StatusNotFound, "invalid environment id: not found")
	}

	return nil
}

func (r *RouteDeployRequest) Model() *model.ApplicationInstance {
	var entity = r.ApplicationInstanceCreateInput.Model()
	entity.ApplicationID = r.ID
	return entity
}

type RouteDeployResponse = *model.ApplicationInstanceOutput
