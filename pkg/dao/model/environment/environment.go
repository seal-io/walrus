// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package environment

import (
	"time"

	"entgo.io/ent"
)

const (
	// Label holds the string label denoting the environment type in the database.
	Label = "environment"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldDescription holds the string denoting the description field in the database.
	FieldDescription = "description"
	// FieldLabels holds the string denoting the labels field in the database.
	FieldLabels = "labels"
	// FieldCreateTime holds the string denoting the createtime field in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the updatetime field in the database.
	FieldUpdateTime = "update_time"
	// FieldVariables holds the string denoting the variables field in the database.
	FieldVariables = "variables"
	// EdgeConnectors holds the string denoting the connectors edge name in mutations.
	EdgeConnectors = "connectors"
	// EdgeApplications holds the string denoting the applications edge name in mutations.
	EdgeApplications = "applications"
	// EdgeRevisions holds the string denoting the revisions edge name in mutations.
	EdgeRevisions = "revisions"
	// EdgeEnvironmentConnectorRelationships holds the string denoting the environmentconnectorrelationships edge name in mutations.
	EdgeEnvironmentConnectorRelationships = "environmentConnectorRelationships"
	// Table holds the table name of the environment in the database.
	Table = "environments"
	// ConnectorsTable is the table that holds the connectors relation/edge. The primary key declared below.
	ConnectorsTable = "environment_connector_relationships"
	// ConnectorsInverseTable is the table name for the Connector entity.
	// It exists in this package in order to avoid circular dependency with the "connector" package.
	ConnectorsInverseTable = "connectors"
	// ApplicationsTable is the table that holds the applications relation/edge.
	ApplicationsTable = "applications"
	// ApplicationsInverseTable is the table name for the Application entity.
	// It exists in this package in order to avoid circular dependency with the "application" package.
	ApplicationsInverseTable = "applications"
	// ApplicationsColumn is the table column denoting the applications relation/edge.
	ApplicationsColumn = "environment_id"
	// RevisionsTable is the table that holds the revisions relation/edge.
	RevisionsTable = "application_revisions"
	// RevisionsInverseTable is the table name for the ApplicationRevision entity.
	// It exists in this package in order to avoid circular dependency with the "applicationrevision" package.
	RevisionsInverseTable = "application_revisions"
	// RevisionsColumn is the table column denoting the revisions relation/edge.
	RevisionsColumn = "environment_id"
	// EnvironmentConnectorRelationshipsTable is the table that holds the environmentConnectorRelationships relation/edge.
	EnvironmentConnectorRelationshipsTable = "environment_connector_relationships"
	// EnvironmentConnectorRelationshipsInverseTable is the table name for the EnvironmentConnectorRelationship entity.
	// It exists in this package in order to avoid circular dependency with the "environmentconnectorrelationship" package.
	EnvironmentConnectorRelationshipsInverseTable = "environment_connector_relationships"
	// EnvironmentConnectorRelationshipsColumn is the table column denoting the environmentConnectorRelationships relation/edge.
	EnvironmentConnectorRelationshipsColumn = "environment_id"
)

// Columns holds all SQL columns for environment fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldDescription,
	FieldLabels,
	FieldCreateTime,
	FieldUpdateTime,
	FieldVariables,
}

var (
	// ConnectorsPrimaryKey and ConnectorsColumn2 are the table columns denoting the
	// primary key for the connectors relation (M2M).
	ConnectorsPrimaryKey = []string{"environment_id", "connector_id"}
)

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

// Note that the variables below are initialized by the runtime
// package on the initialization of the application. Therefore,
// it should be imported in the main as follows:
//
//	import _ "github.com/seal-io/seal/pkg/dao/model/runtime"
var (
	Hooks [1]ent.Hook
	// DefaultCreateTime holds the default value on creation for the "createTime" field.
	DefaultCreateTime func() time.Time
	// DefaultUpdateTime holds the default value on creation for the "updateTime" field.
	DefaultUpdateTime func() time.Time
	// UpdateDefaultUpdateTime holds the default value on update for the "updateTime" field.
	UpdateDefaultUpdateTime func() time.Time
)

// WithoutFields returns the fields ignored the given list.
func WithoutFields(ignores ...string) []string {
	if len(ignores) == 0 {
		return Columns
	}

	var s = make(map[string]bool, len(ignores))
	for i := range ignores {
		s[ignores[i]] = true
	}

	var r = make([]string, 0, len(Columns)-len(s))
	for i := range Columns {
		if s[Columns[i]] {
			continue
		}
		r = append(r, Columns[i])
	}
	return r
}
