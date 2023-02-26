// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package applicationmodulerelationship

import (
	"time"
)

const (
	// Label holds the string label denoting the applicationmodulerelationship type in the database.
	Label = "application_module_relationship"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the createtime field in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the updatetime field in the database.
	FieldUpdateTime = "update_time"
	// FieldApplicationID holds the string denoting the application_id field in the database.
	FieldApplicationID = "application_id"
	// FieldModuleID holds the string denoting the module_id field in the database.
	FieldModuleID = "module_id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldVariables holds the string denoting the variables field in the database.
	FieldVariables = "variables"
	// EdgeApplication holds the string denoting the application edge name in mutations.
	EdgeApplication = "application"
	// EdgeModule holds the string denoting the module edge name in mutations.
	EdgeModule = "module"
	// Table holds the table name of the applicationmodulerelationship in the database.
	Table = "application_module_relationships"
	// ApplicationTable is the table that holds the application relation/edge.
	ApplicationTable = "application_module_relationships"
	// ApplicationInverseTable is the table name for the Application entity.
	// It exists in this package in order to avoid circular dependency with the "application" package.
	ApplicationInverseTable = "applications"
	// ApplicationColumn is the table column denoting the application relation/edge.
	ApplicationColumn = "application_id"
	// ModuleTable is the table that holds the module relation/edge.
	ModuleTable = "application_module_relationships"
	// ModuleInverseTable is the table name for the Module entity.
	// It exists in this package in order to avoid circular dependency with the "module" package.
	ModuleInverseTable = "modules"
	// ModuleColumn is the table column denoting the module relation/edge.
	ModuleColumn = "module_id"
)

// Columns holds all SQL columns for applicationmodulerelationship fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldApplicationID,
	FieldModuleID,
	FieldName,
	FieldVariables,
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
	// DefaultUpdateTime holds the default value on creation for the "updateTime" field.
	DefaultUpdateTime func() time.Time
	// UpdateDefaultUpdateTime holds the default value on update for the "updateTime" field.
	UpdateDefaultUpdateTime func() time.Time
	// ApplicationIDValidator is a validator for the "application_id" field. It is called by the builders before save.
	ApplicationIDValidator func(string) error
	// ModuleIDValidator is a validator for the "module_id" field. It is called by the builders before save.
	ModuleIDValidator func(string) error
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
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
