// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package model

import (
	"time"

	"github.com/seal-io/seal/pkg/dao/types/oid"
)

// EnvironmentConnectorRelationshipQueryInput is the input for the EnvironmentConnectorRelationship query.
type EnvironmentConnectorRelationshipQueryInput struct {
	// ID of the environment to which the relationship connects.
	EnvironmentID oid.ID `json:"environmentId"`
	// ID of the connector to which the relationship connects.
	ConnectorID oid.ID `json:"connectorId"`
}

// Model converts the EnvironmentConnectorRelationshipQueryInput to EnvironmentConnectorRelationship.
func (in EnvironmentConnectorRelationshipQueryInput) Model() *EnvironmentConnectorRelationship {
	return &EnvironmentConnectorRelationship{
		EnvironmentID: in.EnvironmentID,
		ConnectorID:   in.ConnectorID,
	}
}

// EnvironmentConnectorRelationshipCreateInput is the input for the EnvironmentConnectorRelationship creation.
type EnvironmentConnectorRelationshipCreateInput struct {
	// Environments that connect to the relationship.
	Environment EnvironmentQueryInput `json:"environment"`
	// Connectors that connect to the relationship.
	Connector ConnectorQueryInput `json:"connector"`
}

// Model converts the EnvironmentConnectorRelationshipCreateInput to EnvironmentConnectorRelationship.
func (in EnvironmentConnectorRelationshipCreateInput) Model() *EnvironmentConnectorRelationship {
	var entity = &EnvironmentConnectorRelationship{}
	entity.EnvironmentID = in.Environment.ID
	entity.ConnectorID = in.Connector.ID
	return entity
}

// EnvironmentConnectorRelationshipUpdateInput is the input for the EnvironmentConnectorRelationship modification.
type EnvironmentConnectorRelationshipUpdateInput struct {
	// Environments that connect to the relationship.
	Environment EnvironmentQueryInput `json:"environment,omitempty"`
	// Connectors that connect to the relationship.
	Connector ConnectorQueryInput `json:"connector,omitempty"`
}

// Model converts the EnvironmentConnectorRelationshipUpdateInput to EnvironmentConnectorRelationship.
func (in EnvironmentConnectorRelationshipUpdateInput) Model() *EnvironmentConnectorRelationship {
	var entity = &EnvironmentConnectorRelationship{}
	entity.EnvironmentID = in.Environment.ID
	entity.ConnectorID = in.Connector.ID
	return entity
}

// EnvironmentConnectorRelationshipOutput is the output for the EnvironmentConnectorRelationship.
type EnvironmentConnectorRelationshipOutput struct {
	// Describe creation time.
	CreateTime *time.Time `json:"createTime,omitempty"`
	// Environments that connect to the relationship.
	Environment *EnvironmentOutput `json:"environment,omitempty"`
	// Connectors that connect to the relationship.
	Connector *ConnectorOutput `json:"connector,omitempty"`
}

// ExposeEnvironmentConnectorRelationship converts the EnvironmentConnectorRelationship to EnvironmentConnectorRelationshipOutput.
func ExposeEnvironmentConnectorRelationship(in *EnvironmentConnectorRelationship) *EnvironmentConnectorRelationshipOutput {
	if in == nil {
		return nil
	}
	var entity = &EnvironmentConnectorRelationshipOutput{
		CreateTime:  in.CreateTime,
		Environment: ExposeEnvironment(in.Edges.Environment),
		Connector:   ExposeConnector(in.Edges.Connector),
	}
	if in.EnvironmentID != "" {
		if entity.Environment == nil {
			entity.Environment = &EnvironmentOutput{}
		}
		entity.Environment.ID = in.EnvironmentID
	}
	if in.ConnectorID != "" {
		if entity.Connector == nil {
			entity.Connector = &ConnectorOutput{}
		}
		entity.Connector.ID = in.ConnectorID
	}
	return entity
}

// ExposeEnvironmentConnectorRelationships converts the EnvironmentConnectorRelationship slice to EnvironmentConnectorRelationshipOutput pointer slice.
func ExposeEnvironmentConnectorRelationships(in []*EnvironmentConnectorRelationship) []*EnvironmentConnectorRelationshipOutput {
	var out = make([]*EnvironmentConnectorRelationshipOutput, 0, len(in))
	for i := 0; i < len(in); i++ {
		var o = ExposeEnvironmentConnectorRelationship(in[i])
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
