package view

import (
	"context"
	"errors"
	"fmt"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/project"
	"github.com/seal-io/seal/pkg/dao/model/templateversion"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/utils/strs"
	"github.com/seal-io/seal/utils/validation"
)

// Basic APIs.

type CreateRequest struct {
	model.EnvironmentCreateInput `json:",inline"`

	Services []model.ServiceCreateInput `json:"services"`
}

func (r *CreateRequest) ValidateWith(ctx context.Context, input any) error {
	if r.Name == "" {
		return errors.New("invalid name: blank")
	}

	modelClient := input.(model.ClientSet)

	// Get template versions.
	templateVersionKeys := sets.NewString()
	templateVersionPredicates := make([]predicate.TemplateVersion, 0)

	for _, s := range r.Services {
		key := strs.Join("/", s.Template.ID, s.Template.Version)
		if templateVersionKeys.Has(key) {
			continue
		}

		templateVersionKeys.Insert(key)

		templateVersionPredicates = append(templateVersionPredicates, templateversion.And(
			templateversion.TemplateID(s.Template.ID),
			templateversion.Version(s.Template.Version),
		))
	}

	templateVersions, err := modelClient.TemplateVersions().Query().
		Select(
			templateversion.FieldTemplateID,
			templateversion.FieldVersion,
			templateversion.FieldSchema,
		).
		Where(templateversion.Or(
			templateVersionPredicates...,
		)).
		All(ctx)
	if err != nil {
		return runtime.Errorw(err, "failed to get template version")
	}
	templateVersionMap := make(map[string]*model.TemplateVersion, len(templateVersions))

	for _, tv := range templateVersions {
		key := strs.Join("/", tv.TemplateID, tv.Version)
		if _, ok := templateVersionMap[key]; !ok {
			templateVersionMap[key] = tv
		}
	}

	for _, s := range r.Services {
		if s.Name == "" {
			return errors.New("invalid service name: blank")
		}

		if err := validation.IsDNSSubdomainName(s.Name); err != nil {
			return fmt.Errorf("invalid name: %w", err)
		}

		// Verify template version.
		key := strs.Join("/", s.Template.ID, s.Template.Version)

		templateVersion, ok := templateVersionMap[key]
		if !ok {
			return runtime.Errorw(err, "failed to get template version")
		}

		// Verify variables with variables schema that defined on the template version.
		err = s.Attributes.ValidateWith(templateVersion.Schema.Variables)
		if err != nil {
			return fmt.Errorf("invalid variables: %w", err)
		}
	}

	return nil
}

type CreateResponse = *model.EnvironmentOutput

type DeleteRequest = GetRequest

type UpdateRequest struct {
	model.EnvironmentUpdateInput `uri:",inline" json:",inline"`
}

func (r *UpdateRequest) Validate() error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	return nil
}

type GetRequest struct {
	model.EnvironmentQueryInput `uri:",inline"`
}

func (r *GetRequest) Validate() error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	return nil
}

type GetResponse = *model.EnvironmentOutput

// Batch APIs.

type CollectionDeleteRequest []*model.EnvironmentQueryInput

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
	runtime.RequestCollection[predicate.Environment, environment.OrderOption] `query:",inline"`
	ProjectID                                                                 oid.ID `query:"projectID,omitempty"`
	ProjectName                                                               string `query:"projectName,omitempty"`
}

func (r *CollectionGetRequest) ValidateWith(ctx context.Context, input any) error {
	modelClient := input.(model.ClientSet)

	switch {
	case r.ProjectID != "":
		if !r.ProjectID.Valid(0) {
			return errors.New("invalid project id: blank")
		}
	case r.ProjectName != "":
		projectID, err := modelClient.Projects().Query().
			Where(project.Name(r.ProjectName)).
			OnlyID(ctx)
		if err != nil {
			return runtime.Errorw(err, "failed to get project")
		}

		r.ProjectID = projectID
	}

	return nil
}

type CollectionGetResponse = []*model.EnvironmentOutput

// Extensional APIs.
