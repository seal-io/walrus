// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "ent". DO NOT EDIT.

package model

import (
	"time"

	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/pkg/dao/types/property"
	"github.com/seal-io/seal/pkg/dao/types/status"
)

// ServiceQueryInput is the input for the Service query.
type ServiceQueryInput struct {
	// ID holds the value of the "id" field.
	ID oid.ID `uri:"id,omitempty" json:"id,omitempty"`
}

// Model converts the ServiceQueryInput to Service.
func (in ServiceQueryInput) Model() *Service {
	return &Service{
		ID: in.ID,
	}
}

// ServiceCreateInput is the input for the Service creation.
type ServiceCreateInput struct {
	// Name holds the value of the "name" field.
	Name string `json:"name"`
	// Description holds the value of the "description" field.
	Description string `json:"description,omitempty"`
	// Labels holds the value of the "labels" field.
	Labels map[string]string `json:"labels,omitempty"`
	// Template ID and version.
	Template types.TemplateVersionRef `json:"template,omitempty"`
	// Attributes to configure the template.
	Attributes property.Values `json:"attributes,omitempty"`
	// Status of the service.
	Status status.Status `json:"status,omitempty"`
	// Environment to which the service belongs.
	Environment EnvironmentQueryInput `json:"environment"`
	// Dependencies holds the value of the dependencies edge.
	Dependencies []*ServiceRelationshipCreateInput `json:"dependencies,omitempty"`
}

// Model converts the ServiceCreateInput to Service.
func (in ServiceCreateInput) Model() *Service {
	var entity = &Service{
		Name:        in.Name,
		Description: in.Description,
		Labels:      in.Labels,
		Template:    in.Template,
		Attributes:  in.Attributes,
		Status:      in.Status,
	}
	entity.EnvironmentID = in.Environment.ID
	for i := 0; i < len(in.Dependencies); i++ {
		if in.Dependencies[i] == nil {
			continue
		}
		entity.Edges.Dependencies = append(entity.Edges.Dependencies, in.Dependencies[i].Model())
	}
	return entity
}

// ServiceUpdateInput is the input for the Service modification.
type ServiceUpdateInput struct {
	// ID holds the value of the "id" field.
	ID oid.ID `uri:"id" json:"-"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Description holds the value of the "description" field.
	Description string `json:"description,omitempty"`
	// Labels holds the value of the "labels" field.
	Labels map[string]string `json:"labels,omitempty"`
	// Template ID and version.
	Template types.TemplateVersionRef `json:"template,omitempty"`
	// Attributes to configure the template.
	Attributes property.Values `json:"attributes,omitempty"`
	// Status of the service.
	Status status.Status `json:"status,omitempty"`
	// Dependencies holds the value of the dependencies edge.
	Dependencies []*ServiceRelationshipUpdateInput `json:"dependencies,omitempty"`
}

// Model converts the ServiceUpdateInput to Service.
func (in ServiceUpdateInput) Model() *Service {
	var entity = &Service{
		ID:          in.ID,
		Name:        in.Name,
		Description: in.Description,
		Labels:      in.Labels,
		Template:    in.Template,
		Attributes:  in.Attributes,
		Status:      in.Status,
	}
	for i := 0; i < len(in.Dependencies); i++ {
		if in.Dependencies[i] == nil {
			continue
		}
		entity.Edges.Dependencies = append(entity.Edges.Dependencies, in.Dependencies[i].Model())
	}
	return entity
}

// ServiceOutput is the output for the Service.
type ServiceOutput struct {
	// ID holds the value of the "id" field.
	ID oid.ID `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Description holds the value of the "description" field.
	Description string `json:"description,omitempty"`
	// Labels holds the value of the "labels" field.
	Labels map[string]string `json:"labels,omitempty"`
	// CreateTime holds the value of the "createTime" field.
	CreateTime *time.Time `json:"createTime,omitempty"`
	// UpdateTime holds the value of the "updateTime" field.
	UpdateTime *time.Time `json:"updateTime,omitempty"`
	// Template ID and version.
	Template types.TemplateVersionRef `json:"template,omitempty"`
	// Attributes to configure the template.
	Attributes property.Values `json:"attributes,omitempty"`
	// Status of the service.
	Status status.Status `json:"status,omitempty"`
	// Project to which the service belongs.
	Project *ProjectOutput `json:"project,omitempty"`
	// Environment to which the service belongs.
	Environment *EnvironmentOutput `json:"environment,omitempty"`
	// Dependencies holds the value of the dependencies edge.
	Dependencies []*ServiceRelationshipOutput `json:"dependencies,omitempty"`
}

// ExposeService converts the Service to ServiceOutput.
func ExposeService(in *Service) *ServiceOutput {
	if in == nil {
		return nil
	}
	var entity = &ServiceOutput{
		ID:           in.ID,
		Name:         in.Name,
		Description:  in.Description,
		Labels:       in.Labels,
		CreateTime:   in.CreateTime,
		UpdateTime:   in.UpdateTime,
		Template:     in.Template,
		Attributes:   in.Attributes,
		Status:       in.Status,
		Project:      ExposeProject(in.Edges.Project),
		Environment:  ExposeEnvironment(in.Edges.Environment),
		Dependencies: ExposeServiceRelationships(in.Edges.Dependencies),
	}
	if in.ProjectID != "" {
		if entity.Project == nil {
			entity.Project = &ProjectOutput{}
		}
		entity.Project.ID = in.ProjectID
	}
	if in.EnvironmentID != "" {
		if entity.Environment == nil {
			entity.Environment = &EnvironmentOutput{}
		}
		entity.Environment.ID = in.EnvironmentID
	}
	return entity
}

// ExposeServices converts the Service slice to ServiceOutput pointer slice.
func ExposeServices(in []*Service) []*ServiceOutput {
	var out = make([]*ServiceOutput, 0, len(in))
	for i := 0; i < len(in); i++ {
		var o = ExposeService(in[i])
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
