// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package model

import (
	"time"

	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/pkg/dao/types/property"
)

// ApplicationQueryInput is the input for the Application query.
type ApplicationQueryInput struct {
	// ID holds the value of the "id" field.
	ID oid.ID `uri:"id,omitempty" json:"id,omitempty"`
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
	// Variables definition of the application, the variables of instance derived by this definition
	Variables property.Schemas `json:"variables,omitempty"`
	// Project to which this application belongs.
	Project ProjectQueryInput `json:"project"`
	// Modules holds the value of the modules edge.
	Modules []*ApplicationModuleRelationshipCreateInput `json:"modules,omitempty"`
}

// Model converts the ApplicationCreateInput to Application.
func (in ApplicationCreateInput) Model() *Application {
	var entity = &Application{
		Name:        in.Name,
		Description: in.Description,
		Labels:      in.Labels,
		Variables:   in.Variables,
	}
	entity.ProjectID = in.Project.ID
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
	ID oid.ID `uri:"id" json:"-"`
	// Name of the resource.
	Name string `json:"name,omitempty"`
	// Description of the resource.
	Description string `json:"description,omitempty"`
	// Labels of the resource.
	Labels map[string]string `json:"labels,omitempty"`
	// Variables definition of the application, the variables of instance derived by this definition
	Variables property.Schemas `json:"variables,omitempty"`
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
		Variables:   in.Variables,
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
	ID oid.ID `json:"id,omitempty"`
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
	// Variables definition of the application, the variables of instance derived by this definition
	Variables property.Schemas `json:"variables,omitempty"`
	// Project to which this application belongs.
	Project *ProjectOutput `json:"project,omitempty"`
	// Application instances that belong to this application.
	Instances []*ApplicationInstanceOutput `json:"instances,omitempty"`
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
		Variables:   in.Variables,
		Project:     ExposeProject(in.Edges.Project),
		Instances:   ExposeApplicationInstances(in.Edges.Instances),
		Modules:     ExposeApplicationModuleRelationships(in.Edges.Modules),
	}
	if entity.Project == nil {
		entity.Project = &ProjectOutput{}
	}
	entity.Project.ID = in.ProjectID
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
