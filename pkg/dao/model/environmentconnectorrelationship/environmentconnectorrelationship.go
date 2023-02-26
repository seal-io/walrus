// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package environmentconnectorrelationship

import (
	"time"
)

const (
	// Label holds the string label denoting the environmentconnectorrelationship type in the database.
	Label = "environment_connector_relationship"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the createtime field in the database.
	FieldCreateTime = "create_time"
	// FieldEnvironmentID holds the string denoting the environment_id field in the database.
	FieldEnvironmentID = "environment_id"
	// FieldConnectorID holds the string denoting the connector_id field in the database.
	FieldConnectorID = "connector_id"
	// EdgeEnvironment holds the string denoting the environment edge name in mutations.
	EdgeEnvironment = "environment"
	// EdgeConnector holds the string denoting the connector edge name in mutations.
	EdgeConnector = "connector"
	// Table holds the table name of the environmentconnectorrelationship in the database.
	Table = "environment_connector_relationships"
	// EnvironmentTable is the table that holds the environment relation/edge.
	EnvironmentTable = "environment_connector_relationships"
	// EnvironmentInverseTable is the table name for the Environment entity.
	// It exists in this package in order to avoid circular dependency with the "environment" package.
	EnvironmentInverseTable = "environments"
	// EnvironmentColumn is the table column denoting the environment relation/edge.
	EnvironmentColumn = "environment_id"
	// ConnectorTable is the table that holds the connector relation/edge.
	ConnectorTable = "environment_connector_relationships"
	// ConnectorInverseTable is the table name for the Connector entity.
	// It exists in this package in order to avoid circular dependency with the "connector" package.
	ConnectorInverseTable = "connectors"
	// ConnectorColumn is the table column denoting the connector relation/edge.
	ConnectorColumn = "connector_id"
)

// Columns holds all SQL columns for environmentconnectorrelationship fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldEnvironmentID,
	FieldConnectorID,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultCreateTime holds the default value on creation for the "createTime" field.
	DefaultCreateTime func() time.Time
	// EnvironmentIDValidator is a validator for the "environment_id" field. It is called by the builders before save.
	EnvironmentIDValidator func(string) error
	// ConnectorIDValidator is a validator for the "connector_id" field. It is called by the builders before save.
	ConnectorIDValidator func(string) error
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
