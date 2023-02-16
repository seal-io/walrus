// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package connector

import (
	"time"

	"entgo.io/ent"
)

const (
	// Label holds the string label denoting the connector type in the database.
	Label = "connector"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldDescription holds the string denoting the description field in the database.
	FieldDescription = "description"
	// FieldLabels holds the string denoting the labels field in the database.
	FieldLabels = "labels"
	// FieldStatus holds the string denoting the status field in the database.
	FieldStatus = "status"
	// FieldStatusMessage holds the string denoting the statusmessage field in the database.
	FieldStatusMessage = "status_message"
	// FieldCreateTime holds the string denoting the createtime field in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the updatetime field in the database.
	FieldUpdateTime = "update_time"
	// FieldDriver holds the string denoting the driver field in the database.
	FieldDriver = "driver"
	// FieldConfigVersion holds the string denoting the configversion field in the database.
	FieldConfigVersion = "config_version"
	// FieldConfigData holds the string denoting the configdata field in the database.
	FieldConfigData = "config_data"
	// EdgeEnvironment holds the string denoting the environment edge name in mutations.
	EdgeEnvironment = "environment"
	// EdgeEnvironmentConnectorRelationships holds the string denoting the environmentconnectorrelationships edge name in mutations.
	EdgeEnvironmentConnectorRelationships = "environmentConnectorRelationships"
	// Table holds the table name of the connector in the database.
	Table = "connectors"
	// EnvironmentTable is the table that holds the environment relation/edge. The primary key declared below.
	EnvironmentTable = "environment_connector_relationships"
	// EnvironmentInverseTable is the table name for the Environment entity.
	// It exists in this package in order to avoid circular dependency with the "environment" package.
	EnvironmentInverseTable = "environments"
	// EnvironmentConnectorRelationshipsTable is the table that holds the environmentConnectorRelationships relation/edge.
	EnvironmentConnectorRelationshipsTable = "environment_connector_relationships"
	// EnvironmentConnectorRelationshipsInverseTable is the table name for the EnvironmentConnectorRelationship entity.
	// It exists in this package in order to avoid circular dependency with the "environmentconnectorrelationship" package.
	EnvironmentConnectorRelationshipsInverseTable = "environment_connector_relationships"
	// EnvironmentConnectorRelationshipsColumn is the table column denoting the environmentConnectorRelationships relation/edge.
	EnvironmentConnectorRelationshipsColumn = "connector_id"
)

// Columns holds all SQL columns for connector fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldDescription,
	FieldLabels,
	FieldStatus,
	FieldStatusMessage,
	FieldCreateTime,
	FieldUpdateTime,
	FieldDriver,
	FieldConfigVersion,
	FieldConfigData,
}

var (
	// EnvironmentPrimaryKey and EnvironmentColumn2 are the table columns denoting the
	// primary key for the environment relation (M2M).
	EnvironmentPrimaryKey = []string{"environment_id", "connector_id"}
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
