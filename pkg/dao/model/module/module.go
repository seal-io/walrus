// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package module

import (
	"time"

	"github.com/seal-io/seal/pkg/dao/types"
)

const (
	// Label holds the string label denoting the module type in the database.
	Label = "module"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldStatus holds the string denoting the status field in the database.
	FieldStatus = "status"
	// FieldStatusMessage holds the string denoting the statusmessage field in the database.
	FieldStatusMessage = "status_message"
	// FieldCreateTime holds the string denoting the createtime field in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the updatetime field in the database.
	FieldUpdateTime = "update_time"
	// FieldDescription holds the string denoting the description field in the database.
	FieldDescription = "description"
	// FieldLabels holds the string denoting the labels field in the database.
	FieldLabels = "labels"
	// FieldSource holds the string denoting the source field in the database.
	FieldSource = "source"
	// FieldVersion holds the string denoting the version field in the database.
	FieldVersion = "version"
	// FieldSchema holds the string denoting the schema field in the database.
	FieldSchema = "schema"
	// EdgeApplication holds the string denoting the application edge name in mutations.
	EdgeApplication = "application"
	// EdgeApplicationModuleRelationships holds the string denoting the applicationmodulerelationships edge name in mutations.
	EdgeApplicationModuleRelationships = "applicationModuleRelationships"
	// Table holds the table name of the module in the database.
	Table = "modules"
	// ApplicationTable is the table that holds the application relation/edge. The primary key declared below.
	ApplicationTable = "application_module_relationships"
	// ApplicationInverseTable is the table name for the Application entity.
	// It exists in this package in order to avoid circular dependency with the "application" package.
	ApplicationInverseTable = "applications"
	// ApplicationModuleRelationshipsTable is the table that holds the applicationModuleRelationships relation/edge.
	ApplicationModuleRelationshipsTable = "application_module_relationships"
	// ApplicationModuleRelationshipsInverseTable is the table name for the ApplicationModuleRelationship entity.
	// It exists in this package in order to avoid circular dependency with the "applicationmodulerelationship" package.
	ApplicationModuleRelationshipsInverseTable = "application_module_relationships"
	// ApplicationModuleRelationshipsColumn is the table column denoting the applicationModuleRelationships relation/edge.
	ApplicationModuleRelationshipsColumn = "module_id"
)

// Columns holds all SQL columns for module fields.
var Columns = []string{
	FieldID,
	FieldStatus,
	FieldStatusMessage,
	FieldCreateTime,
	FieldUpdateTime,
	FieldDescription,
	FieldLabels,
	FieldSource,
	FieldVersion,
	FieldSchema,
}

var (
	// ApplicationPrimaryKey and ApplicationColumn2 are the table columns denoting the
	// primary key for the application relation (M2M).
	ApplicationPrimaryKey = []string{"application_id", "module_id"}
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

var (
	// DefaultCreateTime holds the default value on creation for the "createTime" field.
	DefaultCreateTime func() time.Time
	// DefaultUpdateTime holds the default value on creation for the "updateTime" field.
	DefaultUpdateTime func() time.Time
	// UpdateDefaultUpdateTime holds the default value on update for the "updateTime" field.
	UpdateDefaultUpdateTime func() time.Time
	// DefaultLabels holds the default value on creation for the "labels" field.
	DefaultLabels map[string]string
	// SourceValidator is a validator for the "source" field. It is called by the builders before save.
	SourceValidator func(string) error
	// DefaultSchema holds the default value on creation for the "schema" field.
	DefaultSchema *types.ModuleSchema
	// IDValidator is a validator for the "id" field. It is called by the builders before save.
	IDValidator func(string) error
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
