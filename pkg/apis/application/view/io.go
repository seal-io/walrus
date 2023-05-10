package view

import (
	"context"
	"errors"
	"fmt"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/model/moduleversion"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/property"
	"github.com/seal-io/seal/pkg/topic/datamessage"
	"github.com/seal-io/seal/utils/validation"
)

// Basic APIs.

type CreateRequest struct {
	*model.ApplicationCreateInput `json:",inline"`
}

func (r *CreateRequest) ValidateWith(ctx context.Context, input any) error {
	modelClient := input.(model.ClientSet)

	if r.Project.ID == "" {
		return errors.New("invalid project id: blank")
	}
	if err := validation.IsDNSSubdomainName(r.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}
	if len(r.Modules) != 0 {
		return validateModules(ctx, modelClient, r.Model().Edges.Modules)
	}
	return nil
}

func validateModules(
	ctx context.Context,
	modelClient model.ClientSet,
	inputModules []*model.ApplicationModuleRelationship,
) error {
	moduleVersionKey := func(moduleID, version string) string {
		return fmt.Sprintf("%s/%s", moduleID, version)
	}

	moduleNames := sets.Set[string]{}
	ps := make([]predicate.ModuleVersion, len(inputModules))
	for i, m := range inputModules {
		if m.ModuleID == "" {
			return errors.New("invalid module id: blank")
		}
		if err := validation.IsDNSSubdomainName(m.Name); err != nil {
			return fmt.Errorf("invalid module name %s: %w", m.Name, err)
		}
		if moduleNames.Has(m.Name) {
			return fmt.Errorf("invalid module name %s: duplicated", m.Name)
		}
		moduleNames.Insert(m.Name)
		ps[i] = moduleversion.And(moduleversion.ModuleID(m.ModuleID), moduleversion.Version(m.Version))
	}

	moduleVersions, err := modelClient.ModuleVersions().
		Query().
		Select(
			moduleversion.FieldModuleID,
			moduleversion.FieldVersion,
			moduleversion.FieldSchema,
		).
		Where(moduleversion.Or(ps...)).
		All(ctx)
	if err != nil {
		return err
	}
	moduleSchemas := make(map[string]property.Schemas, len(moduleVersions))
	for _, m := range moduleVersions {
		if m.Schema == nil {
			continue
		}
		moduleSchemas[moduleVersionKey(m.ModuleID, m.Version)] = m.Schema.Variables
	}

	for _, v := range inputModules {
		moduleSchema, ok := moduleSchemas[moduleVersionKey(v.ModuleID, v.Version)]
		if !ok {
			return fmt.Errorf("invalid module %s: empty schema", v.Name)
		}
		// Verify attributes with attributes schema that defined on versioned module.
		err = v.Attributes.ValidateWith(moduleSchema)
		if err != nil {
			return fmt.Errorf("invalid module %s, "+
				"please refresh related module %q to load the latest attributes: %w",
				v.Name, v.ModuleID, err)
		}
	}
	return nil
}

type CreateResponse = *model.ApplicationOutput

type DeleteRequest = GetRequest

type UpdateRequest struct {
	*model.ApplicationUpdateInput `uri:",inline" json:",inline"`
}

func (r *UpdateRequest) ValidateWith(ctx context.Context, input any) error {
	modelClient := input.(model.ClientSet)

	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}
	if err := validation.IsDNSSubdomainName(r.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}
	if len(r.Modules) != 0 {
		return validateModules(ctx, modelClient, r.Model().Edges.Modules)
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

type StreamResponse struct {
	Type       datamessage.EventType      `json:"type"`
	IDs        []types.ID                 `json:"ids,omitempty"`
	Collection []*model.ApplicationOutput `json:"collection,omitempty"`
}

type StreamRequest struct {
	ID types.ID `uri:"id"`
}

func (r *StreamRequest) ValidateWith(ctx context.Context, input any) error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}
	client := input.(model.ClientSet)
	exist, err := client.Applications().Query().
		Where(application.ID(r.ID)).
		Exist(ctx)
	if err != nil || !exist {
		return runtime.Errorw(err, "invalid id: not found")
	}

	return nil
}

// Batch APIs.

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
	runtime.RequestCollection[predicate.Application, application.OrderOption] `query:",inline"`

	ProjectIDs []types.ID `query:"projectID"`
}

func (r *CollectionGetRequest) Validate() error {
	if len(r.ProjectIDs) == 0 {
		return errors.New("invalid input: missing project id")
	}

	for i := range r.ProjectIDs {
		if !r.ProjectIDs[i].Valid(0) {
			return errors.New("invalid project id: blank")
		}
	}
	return nil
}

type CollectionGetResponse = []*model.ApplicationOutput

type CollectionStreamRequest struct {
	runtime.RequestExtracting `query:",inline"`

	ProjectIDs []types.ID `query:"projectID,omitempty"`
}

func (r *CollectionStreamRequest) Validate() error {
	for i := range r.ProjectIDs {
		if !r.ProjectIDs[i].Valid(0) {
			return errors.New("invalid project id: blank")
		}
	}
	return nil
}

// Extensional APIs.
