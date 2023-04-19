// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package model

import (
	"time"

	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/crypto"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/pkg/dao/types/status"
)

// ConnectorQueryInput is the input for the Connector query.
type ConnectorQueryInput struct {
	// ID holds the value of the "id" field.
	ID oid.ID `uri:"id,omitempty" json:"id,omitempty"`
}

// Model converts the ConnectorQueryInput to Connector.
func (in ConnectorQueryInput) Model() *Connector {
	return &Connector{
		ID: in.ID,
	}
}

// ConnectorCreateInput is the input for the Connector creation.
type ConnectorCreateInput struct {
	// Name of the resource.
	Name string `json:"name"`
	// Description of the resource.
	Description string `json:"description,omitempty"`
	// Labels of the resource.
	Labels map[string]string `json:"labels,omitempty"`
	// Status of the object.
	Status status.Status `json:"status,omitempty"`
	// Type of the connector.
	Type string `json:"type"`
	// Connector config version.
	ConfigVersion string `json:"configVersion"`
	// Connector config data.
	ConfigData crypto.Properties `json:"configData,omitempty"`
	// Config whether enable finOps, will install prometheus and opencost while enable.
	EnableFinOps bool `json:"enableFinOps,omitempty"`
	// Custom pricing user defined.
	FinOpsCustomPricing types.FinOpsCustomPricing `json:"finOpsCustomPricing,omitempty"`
	// Category of the connector.
	Category string `json:"category"`
}

// Model converts the ConnectorCreateInput to Connector.
func (in ConnectorCreateInput) Model() *Connector {
	var entity = &Connector{
		Name:                in.Name,
		Description:         in.Description,
		Labels:              in.Labels,
		Status:              in.Status,
		Type:                in.Type,
		ConfigVersion:       in.ConfigVersion,
		ConfigData:          in.ConfigData,
		EnableFinOps:        in.EnableFinOps,
		FinOpsCustomPricing: in.FinOpsCustomPricing,
		Category:            in.Category,
	}
	return entity
}

// ConnectorUpdateInput is the input for the Connector modification.
type ConnectorUpdateInput struct {
	// ID holds the value of the "id" field.
	ID oid.ID `uri:"id" json:"-"`
	// Name of the resource.
	Name string `json:"name,omitempty"`
	// Description of the resource.
	Description string `json:"description,omitempty"`
	// Labels of the resource.
	Labels map[string]string `json:"labels,omitempty"`
	// Status of the object.
	Status status.Status `json:"status,omitempty"`
	// Connector config version.
	ConfigVersion string `json:"configVersion,omitempty"`
	// Connector config data.
	ConfigData crypto.Properties `json:"configData,omitempty"`
	// Config whether enable finOps, will install prometheus and opencost while enable.
	EnableFinOps bool `json:"enableFinOps,omitempty"`
	// Custom pricing user defined.
	FinOpsCustomPricing types.FinOpsCustomPricing `json:"finOpsCustomPricing,omitempty"`
	// Category of the connector.
	Category string `json:"category,omitempty"`
}

// Model converts the ConnectorUpdateInput to Connector.
func (in ConnectorUpdateInput) Model() *Connector {
	var entity = &Connector{
		ID:                  in.ID,
		Name:                in.Name,
		Description:         in.Description,
		Labels:              in.Labels,
		Status:              in.Status,
		ConfigVersion:       in.ConfigVersion,
		ConfigData:          in.ConfigData,
		EnableFinOps:        in.EnableFinOps,
		FinOpsCustomPricing: in.FinOpsCustomPricing,
		Category:            in.Category,
	}
	return entity
}

// ConnectorOutput is the output for the Connector.
type ConnectorOutput struct {
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
	// Status of the object.
	Status status.Status `json:"status,omitempty"`
	// Type of the connector.
	Type string `json:"type,omitempty"`
	// Connector config version.
	ConfigVersion string `json:"configVersion,omitempty"`
	// Connector config data.
	ConfigData crypto.Properties `json:"configData,omitempty"`
	// Config whether enable finOps, will install prometheus and opencost while enable.
	EnableFinOps bool `json:"enableFinOps,omitempty"`
	// Custom pricing user defined.
	FinOpsCustomPricing types.FinOpsCustomPricing `json:"finOpsCustomPricing,omitempty"`
	// Category of the connector.
	Category string `json:"category,omitempty"`
	// Environments holds the value of the environments edge.
	Environments []*EnvironmentConnectorRelationshipOutput `json:"environments,omitempty"`
	// Resources that belong to the application.
	Resources []*ApplicationResourceOutput `json:"resources,omitempty"`
	// Cluster costs that linked to the connection
	ClusterCosts []*ClusterCostOutput `json:"clusterCosts,omitempty"`
	// Cluster allocation resource costs that linked to the connection.
	AllocationCosts []*AllocationCostOutput `json:"allocationCosts,omitempty"`
}

// ExposeConnector converts the Connector to ConnectorOutput.
func ExposeConnector(in *Connector) *ConnectorOutput {
	if in == nil {
		return nil
	}
	var entity = &ConnectorOutput{
		ID:                  in.ID,
		Name:                in.Name,
		Description:         in.Description,
		Labels:              in.Labels,
		CreateTime:          in.CreateTime,
		UpdateTime:          in.UpdateTime,
		Status:              in.Status,
		Type:                in.Type,
		ConfigVersion:       in.ConfigVersion,
		ConfigData:          in.ConfigData,
		EnableFinOps:        in.EnableFinOps,
		FinOpsCustomPricing: in.FinOpsCustomPricing,
		Category:            in.Category,
		Environments:        ExposeEnvironmentConnectorRelationships(in.Edges.Environments),
		Resources:           ExposeApplicationResources(in.Edges.Resources),
		ClusterCosts:        ExposeClusterCosts(in.Edges.ClusterCosts),
		AllocationCosts:     ExposeAllocationCosts(in.Edges.AllocationCosts),
	}
	return entity
}

// ExposeConnectors converts the Connector slice to ConnectorOutput pointer slice.
func ExposeConnectors(in []*Connector) []*ConnectorOutput {
	var out = make([]*ConnectorOutput, 0, len(in))
	for i := 0; i < len(in); i++ {
		var o = ExposeConnector(in[i])
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
