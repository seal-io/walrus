// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package model

import (
	"time"

	"github.com/seal-io/seal/pkg/dao/types"
)

// ApplicationQueryInput is the input for the Application query.
type ApplicationQueryInput struct {
	// ID holds the value of the "id" field.
	ID types.ID `uri:"id,omitempty" json:"id,omitempty"`
}

// Model converts the ApplicationQueryInput to Application.
func (in ApplicationQueryInput) Model() *Application {
	return &Application{
		ID: in.ID,
	}
}

// ApplicationCreateInput is the input for the Application creation.
type ApplicationCreateInput struct {
	// Name of the resource.
	Name string `json:"name"`
	// Description of the resource.
	Description string `json:"description,omitempty"`
	// Labels of the resource.
	Labels map[string]string `json:"labels,omitempty"`
	// Project to which this application belongs.
	Project ProjectQueryInput `json:"project"`
	// Environment to which the application belongs.
	Environment EnvironmentQueryInput `json:"environment"`
	// Modules holds the value of the modules edge.
	Modules []*ApplicationModuleRelationshipCreateInput `json:"modules,omitempty"`
}

// Model converts the ApplicationCreateInput to Application.
func (in ApplicationCreateInput) Model() *Application {
	var entity = &Application{
		Name:        in.Name,
		Description: in.Description,
		Labels:      in.Labels,
	}
	entity.ProjectID = in.Project.ID
	entity.EnvironmentID = in.Environment.ID
	for i := 0; i < len(in.Modules); i++ {
		if in.Modules[i] == nil {
			continue
		}
		entity.Edges.Modules = append(entity.Edges.Modules, in.Modules[i].Model())
	}
	return entity
}

// ApplicationUpdateInput is the input for the Application modification.
type ApplicationUpdateInput struct {
	// ID holds the value of the "id" field.
	ID types.ID `uri:"id" json:"-"`
	// Name of the resource.
	Name string `json:"name,omitempty"`
	// Description of the resource.
	Description string `json:"description,omitempty"`
	// Labels of the resource.
	Labels map[string]string `json:"labels,omitempty"`
	// Modules holds the value of the modules edge.
	Modules []*ApplicationModuleRelationshipUpdateInput `json:"modules,omitempty"`
}

// Model converts the ApplicationUpdateInput to Application.
func (in ApplicationUpdateInput) Model() *Application {
	var entity = &Application{
		ID:          in.ID,
		Name:        in.Name,
		Description: in.Description,
		Labels:      in.Labels,
	}
	for i := 0; i < len(in.Modules); i++ {
		if in.Modules[i] == nil {
			continue
		}
		entity.Edges.Modules = append(entity.Edges.Modules, in.Modules[i].Model())
	}
	return entity
}

// ApplicationOutput is the output for the Application.
type ApplicationOutput struct {
	// ID holds the value of the "id" field.
	ID types.ID `json:"id,omitempty"`
	// Name of the resource.
	Name string `json:"name,omitempty"`
	// Description of the resource.
	Description string `json:"description,omitempty"`
	// Labels of the resource.
	Labels map[string]string `json:"labels,omitempty"`
	// Describe creation time.
	CreateTime *time.Time `json:"createTime,omitempty"`
	// Describe modification time.
	UpdateTime *time.Time `json:"updateTime,omitempty"`
	// Project to which this application belongs.
	Project *ProjectOutput `json:"project,omitempty"`
	// Environment to which the application belongs.
	Environment *EnvironmentOutput `json:"environment,omitempty"`
	// Resources that belong to the application.
	Resources []*ApplicationResourceOutput `json:"resources,omitempty"`
	// Revisions that belong to this application.
	Revisions []*ApplicationRevisionOutput `json:"revisions,omitempty"`
	// Modules holds the value of the modules edge.
	Modules []*ApplicationModuleRelationshipOutput `json:"modules,omitempty"`
}

// ExposeApplication converts the Application to ApplicationOutput.
func ExposeApplication(in *Application) *ApplicationOutput {
	if in == nil {
		return nil
	}
	var entity = &ApplicationOutput{
		ID:          in.ID,
		Name:        in.Name,
		Description: in.Description,
		Labels:      in.Labels,
		CreateTime:  in.CreateTime,
		UpdateTime:  in.UpdateTime,
		Project:     ExposeProject(in.Edges.Project),
		Environment: ExposeEnvironment(in.Edges.Environment),
		Resources:   ExposeApplicationResources(in.Edges.Resources),
		Revisions:   ExposeApplicationRevisions(in.Edges.Revisions),
		Modules:     ExposeApplicationModuleRelationships(in.Edges.Modules),
	}
	if entity.Project == nil {
		entity.Project = &ProjectOutput{}
	}
	entity.Project.ID = in.ProjectID
	if entity.Environment == nil {
		entity.Environment = &EnvironmentOutput{}
	}
	entity.Environment.ID = in.EnvironmentID
	return entity
}

// ExposeApplications converts the Application slice to ApplicationOutput pointer slice.
func ExposeApplications(in []*Application) []*ApplicationOutput {
	var out = make([]*ApplicationOutput, 0, len(in))
	for i := 0; i < len(in); i++ {
		var o = ExposeApplication(in[i])
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
