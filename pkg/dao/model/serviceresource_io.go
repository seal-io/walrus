// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "ent". DO NOT EDIT.

package model

import (
	"time"

	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/oid"
)

// ServiceResourceQueryInput is the input for the ServiceResource query.
type ServiceResourceQueryInput struct {
	// ID holds the value of the "id" field.
	ID oid.ID `uri:"id,omitempty" json:"id,omitempty"`
}

// Model converts the ServiceResourceQueryInput to ServiceResource.
func (in ServiceResourceQueryInput) Model() *ServiceResource {
	return &ServiceResource{
		ID: in.ID,
	}
}

// ServiceResourceCreateInput is the input for the ServiceResource creation.
type ServiceResourceCreateInput struct {
	// ID of the project to belong.
	ProjectID oid.ID `json:"projectID"`
	// Mode that manages the generated resource, it is the management way of the deployer to the resource, which provides by deployer.
	Mode string `json:"mode"`
	// Type of the generated resource, it is the type of the resource which the deployer observes, which provides by deployer.
	Type string `json:"type"`
	// Name of the generated resource, it is the real identifier of the resource, which provides by deployer.
	Name string `json:"name"`
	// Type of deployer.
	DeployerType string `json:"deployerType"`
	// Status of the resource.
	Status types.ServiceResourceStatus `json:"status,omitempty"`
	// Service resource to which the resource makes up.
	Composition *ServiceResourceQueryInput `json:"composition,omitempty"`
}

// Model converts the ServiceResourceCreateInput to ServiceResource.
func (in ServiceResourceCreateInput) Model() *ServiceResource {
	var entity = &ServiceResource{
		ProjectID:    in.ProjectID,
		Mode:         in.Mode,
		Type:         in.Type,
		Name:         in.Name,
		DeployerType: in.DeployerType,
		Status:       in.Status,
	}
	if in.Composition != nil {
		entity.CompositionID = in.Composition.ID
	}
	return entity
}

// ServiceResourceUpdateInput is the input for the ServiceResource modification.
type ServiceResourceUpdateInput struct {
	// ID holds the value of the "id" field.
	ID oid.ID `uri:"id" json:"-"`
	// Status of the resource.
	Status types.ServiceResourceStatus `json:"status,omitempty"`
}

// Model converts the ServiceResourceUpdateInput to ServiceResource.
func (in ServiceResourceUpdateInput) Model() *ServiceResource {
	var entity = &ServiceResource{
		ID:     in.ID,
		Status: in.Status,
	}
	return entity
}

// ServiceResourceOutput is the output for the ServiceResource.
type ServiceResourceOutput struct {
	// ID holds the value of the "id" field.
	ID oid.ID `json:"id,omitempty"`
	// CreateTime holds the value of the "createTime" field.
	CreateTime *time.Time `json:"createTime,omitempty"`
	// UpdateTime holds the value of the "updateTime" field.
	UpdateTime *time.Time `json:"updateTime,omitempty"`
	// ID of the project to belong.
	ProjectID oid.ID `json:"projectID,omitempty"`
	// Mode that manages the generated resource, it is the management way of the deployer to the resource, which provides by deployer.
	Mode string `json:"mode,omitempty"`
	// Type of the generated resource, it is the type of the resource which the deployer observes, which provides by deployer.
	Type string `json:"type,omitempty"`
	// Name of the generated resource, it is the real identifier of the resource, which provides by deployer.
	Name string `json:"name,omitempty"`
	// Type of deployer.
	DeployerType string `json:"deployerType,omitempty"`
	// Status of the resource.
	Status types.ServiceResourceStatus `json:"status,omitempty"`
	// Service to which the resource belongs.
	Service *ServiceOutput `json:"service,omitempty"`
	// Connector to which the resource deploys.
	Connector *ConnectorOutput `json:"connector,omitempty"`
	// Service resource to which the resource makes up.
	Composition *ServiceResourceOutput `json:"composition,omitempty"`
	// Sub-resources that make up the resource.
	Components []*ServiceResourceOutput `json:"components,omitempty"`
	// Keys is the list of key used for operating the service resource,
	// it does not store in the database and only records for additional operations.
	Keys *types.ServiceResourceOperationKeys `json:"keys,omitempty"`
}

// ExposeServiceResource converts the ServiceResource to ServiceResourceOutput.
func ExposeServiceResource(in *ServiceResource) *ServiceResourceOutput {
	if in == nil {
		return nil
	}
	var entity = &ServiceResourceOutput{
		ID:           in.ID,
		CreateTime:   in.CreateTime,
		UpdateTime:   in.UpdateTime,
		ProjectID:    in.ProjectID,
		Mode:         in.Mode,
		Type:         in.Type,
		Name:         in.Name,
		DeployerType: in.DeployerType,
		Status:       in.Status,
		Service:      ExposeService(in.Edges.Service),
		Connector:    ExposeConnector(in.Edges.Connector),
		Composition:  ExposeServiceResource(in.Edges.Composition),
		Components:   ExposeServiceResources(in.Edges.Components),
	}
	if in.ServiceID != "" {
		if entity.Service == nil {
			entity.Service = &ServiceOutput{}
		}
		entity.Service.ID = in.ServiceID
	}
	if in.ConnectorID != "" {
		if entity.Connector == nil {
			entity.Connector = &ConnectorOutput{}
		}
		entity.Connector.ID = in.ConnectorID
	}
	if in.CompositionID != "" {
		if entity.Composition == nil {
			entity.Composition = &ServiceResourceOutput{}
		}
		entity.Composition.ID = in.CompositionID
	}
	entity.Keys = in.Keys

	return entity
}

// ExposeServiceResources converts the ServiceResource slice to ServiceResourceOutput pointer slice.
func ExposeServiceResources(in []*ServiceResource) []*ServiceResourceOutput {
	var out = make([]*ServiceResourceOutput, 0, len(in))
	for i := 0; i < len(in); i++ {
		var o = ExposeServiceResource(in[i])
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
