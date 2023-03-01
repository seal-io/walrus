// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package model

import (
	"time"

	"github.com/seal-io/seal/pkg/dao/types"
)

// ApplicationResourceQueryInput is the input for the ApplicationResource query.
type ApplicationResourceQueryInput struct {
	// ID holds the value of the "id" field.
	ID types.ID `uri:"id,omitempty" json:"id,omitempty"`
}

// Model converts the ApplicationResourceQueryInput to ApplicationResource.
func (in ApplicationResourceQueryInput) Model() *ApplicationResource {
	return &ApplicationResource{
		ID: in.ID,
	}
}

// ApplicationResourceCreateInput is the input for the ApplicationResource creation.
type ApplicationResourceCreateInput struct {
	// Status of the resource.
	Status string `json:"status,omitempty"`
	// Extra message for status, like error details.
	StatusMessage string `json:"statusMessage,omitempty"`
	// Name of the module that generates the resource.
	Module string `json:"module"`
	// Mode that manages the generated resource, it is the management way of the deployer to the resource, which provides by deployer.
	Mode string `json:"mode"`
	// Type of the generated resource, it is the type of the resource which the deployer observes, which provides by deployer.
	Type string `json:"type"`
	// Name of the generated resource, it is the real identifier of the resource, which provides by deployer.
	Name string `json:"name"`
	// Application to which the resource belongs.
	Application ApplicationQueryInput `json:"application"`
	// Connector to which the resource deploys.
	Connector ConnectorQueryInput `json:"connector"`
}

// Model converts the ApplicationResourceCreateInput to ApplicationResource.
func (in ApplicationResourceCreateInput) Model() *ApplicationResource {
	var entity = &ApplicationResource{
		Status:        in.Status,
		StatusMessage: in.StatusMessage,
		Module:        in.Module,
		Mode:          in.Mode,
		Type:          in.Type,
		Name:          in.Name,
	}
	entity.ApplicationID = in.Application.ID
	entity.ConnectorID = in.Connector.ID
	return entity
}

// ApplicationResourceUpdateInput is the input for the ApplicationResource modification.
type ApplicationResourceUpdateInput struct {
	// ID holds the value of the "id" field.
	ID types.ID `uri:"id" json:"-"`
	// Status of the resource.
	Status string `json:"status,omitempty"`
	// Extra message for status, like error details.
	StatusMessage string `json:"statusMessage,omitempty"`
}

// Model converts the ApplicationResourceUpdateInput to ApplicationResource.
func (in ApplicationResourceUpdateInput) Model() *ApplicationResource {
	var entity = &ApplicationResource{
		ID:            in.ID,
		Status:        in.Status,
		StatusMessage: in.StatusMessage,
	}
	return entity
}

// ApplicationResourceOutput is the output for the ApplicationResource.
type ApplicationResourceOutput struct {
	// ID holds the value of the "id" field.
	ID types.ID `json:"id,omitempty"`
	// Status of the resource.
	Status string `json:"status,omitempty"`
	// Extra message for status, like error details.
	StatusMessage string `json:"statusMessage,omitempty"`
	// Describe creation time.
	CreateTime *time.Time `json:"createTime,omitempty"`
	// Describe modification time.
	UpdateTime *time.Time `json:"updateTime,omitempty"`
	// Name of the module that generates the resource.
	Module string `json:"module,omitempty"`
	// Mode that manages the generated resource, it is the management way of the deployer to the resource, which provides by deployer.
	Mode string `json:"mode,omitempty"`
	// Type of the generated resource, it is the type of the resource which the deployer observes, which provides by deployer.
	Type string `json:"type,omitempty"`
	// Name of the generated resource, it is the real identifier of the resource, which provides by deployer.
	Name string `json:"name,omitempty"`
	// Application to which the resource belongs.
	Application *ApplicationOutput `json:"application,omitempty"`
	// Connector to which the resource deploys.
	Connector *ConnectorOutput `json:"connector,omitempty"`
}

// ExposeApplicationResource converts the ApplicationResource to ApplicationResourceOutput.
func ExposeApplicationResource(in *ApplicationResource) *ApplicationResourceOutput {
	if in == nil {
		return nil
	}
	var entity = &ApplicationResourceOutput{
		ID:            in.ID,
		Status:        in.Status,
		StatusMessage: in.StatusMessage,
		CreateTime:    in.CreateTime,
		UpdateTime:    in.UpdateTime,
		Module:        in.Module,
		Mode:          in.Mode,
		Type:          in.Type,
		Name:          in.Name,
		Application:   ExposeApplication(in.Edges.Application),
		Connector:     ExposeConnector(in.Edges.Connector),
	}
	if entity.Application == nil {
		entity.Application = &ApplicationOutput{}
	}
	entity.Application.ID = in.ApplicationID
	if entity.Connector == nil {
		entity.Connector = &ConnectorOutput{}
	}
	entity.Connector.ID = in.ConnectorID
	return entity
}

// ExposeApplicationResources converts the ApplicationResource slice to ApplicationResourceOutput pointer slice.
func ExposeApplicationResources(in []*ApplicationResource) []*ApplicationResourceOutput {
	var out = make([]*ApplicationResourceOutput, 0, len(in))
	for i := 0; i < len(in); i++ {
		var o = ExposeApplicationResource(in[i])
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
