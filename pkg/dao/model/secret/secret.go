// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package secret

import (
	"time"

	"entgo.io/ent"
)

const (
	// Label holds the string label denoting the secret type in the database.
	Label = "secret"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the createtime field in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the updatetime field in the database.
	FieldUpdateTime = "update_time"
	// FieldProjectID holds the string denoting the projectid field in the database.
	FieldProjectID = "project_id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldValue holds the string denoting the value field in the database.
	FieldValue = "value"
	// EdgeProject holds the string denoting the project edge name in mutations.
	EdgeProject = "project"
	// Table holds the table name of the secret in the database.
	Table = "secrets"
	// ProjectTable is the table that holds the project relation/edge.
	ProjectTable = "secrets"
	// ProjectInverseTable is the table name for the Project entity.
	// It exists in this package in order to avoid circular dependency with the "project" package.
	ProjectInverseTable = "projects"
	// ProjectColumn is the table column denoting the project relation/edge.
	ProjectColumn = "project_id"
)

// Columns holds all SQL columns for secret fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldProjectID,
	FieldName,
	FieldValue,
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
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
	// ValueValidator is a validator for the "value" field. It is called by the builders before save.
	ValueValidator func(string) error
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
