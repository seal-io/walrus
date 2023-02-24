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
	// FieldType holds the string denoting the type field in the database.
	FieldType = "type"
	// FieldConfigVersion holds the string denoting the configversion field in the database.
	FieldConfigVersion = "config_version"
	// FieldConfigData holds the string denoting the configdata field in the database.
	FieldConfigData = "config_data"
	// FieldEnableFinOps holds the string denoting the enablefinops field in the database.
	FieldEnableFinOps = "enable_fin_ops"
	// FieldFinOpsStatus holds the string denoting the finopsstatus field in the database.
	FieldFinOpsStatus = "fin_ops_status"
	// FieldFinOpsStatusMessage holds the string denoting the finopsstatusmessage field in the database.
	FieldFinOpsStatusMessage = "fin_ops_status_message"
	// EdgeEnvironments holds the string denoting the environments edge name in mutations.
	EdgeEnvironments = "environments"
	// EdgeResources holds the string denoting the resources edge name in mutations.
	EdgeResources = "resources"
	// EdgeEnvironmentConnectorRelationships holds the string denoting the environmentconnectorrelationships edge name in mutations.
	EdgeEnvironmentConnectorRelationships = "environmentConnectorRelationships"
	// Table holds the table name of the connector in the database.
	Table = "connectors"
	// EnvironmentsTable is the table that holds the environments relation/edge. The primary key declared below.
	EnvironmentsTable = "environment_connector_relationships"
	// EnvironmentsInverseTable is the table name for the Environment entity.
	// It exists in this package in order to avoid circular dependency with the "environment" package.
	EnvironmentsInverseTable = "environments"
	// ResourcesTable is the table that holds the resources relation/edge.
	ResourcesTable = "application_resources"
	// ResourcesInverseTable is the table name for the ApplicationResource entity.
	// It exists in this package in order to avoid circular dependency with the "applicationresource" package.
	ResourcesInverseTable = "application_resources"
	// ResourcesColumn is the table column denoting the resources relation/edge.
	ResourcesColumn = "connector_id"
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
	FieldType,
	FieldConfigVersion,
	FieldConfigData,
	FieldEnableFinOps,
	FieldFinOpsStatus,
	FieldFinOpsStatusMessage,
}

var (
	// EnvironmentsPrimaryKey and EnvironmentsColumn2 are the table columns denoting the
	// primary key for the environments relation (M2M).
	EnvironmentsPrimaryKey = []string{"environment_id", "connector_id"}
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
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
	// DefaultLabels holds the default value on creation for the "labels" field.
	DefaultLabels map[string]string
	// DefaultCreateTime holds the default value on creation for the "createTime" field.
	DefaultCreateTime func() time.Time
	// DefaultUpdateTime holds the default value on creation for the "updateTime" field.
	DefaultUpdateTime func() time.Time
	// UpdateDefaultUpdateTime holds the default value on update for the "updateTime" field.
	UpdateDefaultUpdateTime func() time.Time
	// TypeValidator is a validator for the "type" field. It is called by the builders before save.
	TypeValidator func(string) error
	// ConfigVersionValidator is a validator for the "configVersion" field. It is called by the builders before save.
	ConfigVersionValidator func(string) error
	// DefaultConfigData holds the default value on creation for the "configData" field.
	DefaultConfigData map[string]interface{}
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
