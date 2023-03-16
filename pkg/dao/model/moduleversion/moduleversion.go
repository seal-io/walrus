// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package moduleversion

import (
	"time"

	"entgo.io/ent"

	"github.com/seal-io/seal/pkg/dao/types"
)

const (
	// Label holds the string label denoting the moduleversion type in the database.
	Label = "module_version"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the createtime field in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the updatetime field in the database.
	FieldUpdateTime = "update_time"
	// FieldModuleID holds the string denoting the moduleid field in the database.
	FieldModuleID = "module_id"
	// FieldVersion holds the string denoting the version field in the database.
	FieldVersion = "version"
	// FieldSource holds the string denoting the source field in the database.
	FieldSource = "source"
	// FieldSchema holds the string denoting the schema field in the database.
	FieldSchema = "schema"
	// EdgeModule holds the string denoting the module edge name in mutations.
	EdgeModule = "module"
	// Table holds the table name of the moduleversion in the database.
	Table = "module_versions"
	// ModuleTable is the table that holds the module relation/edge.
	ModuleTable = "module_versions"
	// ModuleInverseTable is the table name for the Module entity.
	// It exists in this package in order to avoid circular dependency with the "module" package.
	ModuleInverseTable = "modules"
	// ModuleColumn is the table column denoting the module relation/edge.
	ModuleColumn = "module_id"
)

// Columns holds all SQL columns for moduleversion fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldModuleID,
	FieldVersion,
	FieldSource,
	FieldSchema,
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
	// ModuleIDValidator is a validator for the "moduleID" field. It is called by the builders before save.
	ModuleIDValidator func(string) error
	// VersionValidator is a validator for the "version" field. It is called by the builders before save.
	VersionValidator func(string) error
	// SourceValidator is a validator for the "source" field. It is called by the builders before save.
	SourceValidator func(string) error
	// DefaultSchema holds the default value on creation for the "schema" field.
	DefaultSchema *types.ModuleSchema
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
