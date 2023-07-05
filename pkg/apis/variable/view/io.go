package view

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/project"
	"github.com/seal-io/seal/pkg/dao/model/variable"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/utils/validation"
)

// Basic APIs.

type CreateRequest struct {
	model.VariableCreateInput `json:",inline"`

	ProjectID       oid.ID `query:"projectID,omitempty"`
	ProjectName     string `query:"projectName,omitempty"`
	EnvironmentID   oid.ID `query:"environmentID,omitempty"`
	EnvironmentName string `query:"environmentName,omitempty"`
}

func (r *CreateRequest) ValidateWith(ctx context.Context, input any) error {
	modelClient := input.(model.ClientSet)

	switch {
	case r.ProjectID != "":
		if !r.ProjectID.Valid(0) {
			return errors.New("invalid project id: blank")
		}

		r.Project = &model.ProjectQueryInput{
			ID: r.ProjectID,
		}
	case r.ProjectName != "":
		projectID, err := modelClient.Projects().Query().
			Where(project.Name(r.ProjectName)).
			OnlyID(ctx)
		if err != nil {
			return runtime.Errorw(err, "failed to get project")
		}

		r.ProjectID = projectID
		r.Project = &model.ProjectQueryInput{
			ID: projectID,
		}
	}

	switch {
	case r.EnvironmentID != "":
		if r.ProjectID == "" {
			return errors.New("invalid project id: blank")
		}

		if !r.EnvironmentID.Valid(0) {
			return errors.New("invalid environment id: blank")
		}

		r.Environment = &model.EnvironmentQueryInput{
			ID: r.EnvironmentID,
		}
	case r.EnvironmentName != "":
		if r.ProjectID == "" {
			return errors.New("invalid project id: blank")
		}

		envID, err := modelClient.Environments().Query().
			Where(
				environment.ProjectID(r.ProjectID),
				environment.Name(r.EnvironmentName),
			).
			OnlyID(ctx)
		if err != nil {
			return runtime.Errorw(err, "failed to get environment")
		}

		r.EnvironmentID = envID
		r.Environment = &model.EnvironmentQueryInput{
			ID: envID,
		}
	}

	if err := validation.IsDNSLabel(r.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	if r.Value == "" {
		return errors.New("invalid value: blank")
	}

	return nil
}

type CreateResponse = *model.VariableOutput

func ExposeVariable(in *model.Variable) *model.VariableOutput {
	if in.Sensitive {
		in.Value = ""
	}

	return model.ExposeVariable(in)
}

type DeleteRequest struct {
	model.VariableQueryInput `uri:",inline"`

	ProjectID   oid.ID `query:"projectID,omitempty"`
	ProjectName string `query:"projectName,omitempty"`
}

func (r *DeleteRequest) ValidateWith(ctx context.Context, input any) error {
	modelClient := input.(model.ClientSet)

	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

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

	// FIXME(thxCode): a workaround to protect general user deleting global variable,
	//   returns a not found error instead of forbidden.
	if r.ProjectID != "" {
		exist, err := modelClient.Variables().Query().
			Where(
				variable.ID(r.ID),
				variable.ProjectID(r.ProjectID)).
			Exist(ctx)
		if err != nil {
			return err
		}

		if !exist {
			return runtime.Errorc(http.StatusNotFound)
		}
	}

	return nil
}

type UpdateRequest struct {
	model.VariableUpdateInput `uri:",inline" json:",inline"`

	ProjectID   oid.ID `query:"projectID,omitempty"`
	ProjectName string `query:"projectName,omitempty"`
}

func (r *UpdateRequest) ValidateWith(ctx context.Context, input any) error {
	modelClient := input.(model.ClientSet)

	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	if r.Value == "" {
		return errors.New("invalid value: blank")
	}

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

	// FIXME(thxCode): a workaround to protect general user deleting global variable,
	//   returns a not found error instead of forbidden.
	if r.ProjectID != "" {
		exist, err := modelClient.Variables().Query().
			Where(
				variable.ID(r.ID),
				variable.ProjectID(r.ProjectID)).
			Exist(ctx)
		if err != nil {
			return err
		}

		if !exist {
			return runtime.Errorc(http.StatusNotFound)
		}
	}

	return nil
}

// Batch APIs.

type CollectionDeleteRequest struct {
	Items []*model.VariableQueryInput `json:"items"`

	ProjectID   oid.ID `query:"projectID,omitempty"`
	ProjectName string `query:"projectName,omitempty"`
}

func (r CollectionDeleteRequest) ValidateWith(ctx context.Context, input any) error {
	modelClient := input.(model.ClientSet)

	if len(r.Items) == 0 {
		return errors.New("invalid input: empty")
	}

	for _, i := range r.Items {
		if !i.ID.Valid(0) {
			return errors.New("invalid id: blank")
		}
	}

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

	// FIXME(thxCode): a workaround to protect general user deleting global variable,
	//   returns a not found error instead of forbidden.
	if r.ProjectID != "" {
		ids := make([]oid.ID, len(r.Items))
		for i := range r.Items {
			ids[i] = r.Items[i].ID
		}

		cnt, err := modelClient.Variables().Query().
			Where(
				variable.IDIn(ids...),
				variable.ProjectID(r.ProjectID)).
			Count(ctx)
		if err != nil {
			return err
		}

		if cnt != len(ids) {
			return runtime.Errorc(http.StatusNotFound)
		}
	}

	return nil
}

type CollectionGetRequest struct {
	runtime.RequestCollection[predicate.Variable, variable.OrderOption] `query:",inline"`

	ProjectID       oid.ID `query:"projectID,omitempty"`
	ProjectName     string `query:"projectName,omitempty"`
	EnvironmentID   oid.ID `query:"environmentID,omitempty"`
	EnvironmentName string `query:"environmentName,omitempty"`
	WithGlobal      bool   `query:"withGlobal,omitempty"`
	WithProject     bool   `query:"withProject,omitempty"`
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

	switch {
	case r.EnvironmentID != "":
		if r.ProjectID == "" {
			return errors.New("invalid project id: blank")
		}

		if !r.EnvironmentID.Valid(0) {
			return errors.New("invalid environment id: blank")
		}
	case r.EnvironmentName != "":
		if r.ProjectID == "" {
			return errors.New("invalid project id: blank")
		}

		envID, err := modelClient.Environments().Query().
			Where(
				environment.ProjectID(r.ProjectID),
				environment.Name(r.EnvironmentName),
			).
			OnlyID(ctx)
		if err != nil {
			return runtime.Errorw(err, "failed to get environment")
		}

		r.EnvironmentID = envID
	}

	return nil
}

type CollectionGetResponse = []*model.VariableOutput

// ExposeVariables converts the Variable slice to VariableOutput pointer slice.
func ExposeVariables(in []*model.Variable) []*model.VariableOutput {
	out := make([]*model.VariableOutput, 0, len(in))

	for i := 0; i < len(in); i++ {
		o := ExposeVariable(in[i])
		if o == nil {
			continue
		}

		out = append(out, o)
	}

	if len(out) == 0 {
		return nil
	}

	return out
}

// Extensional APIs.
