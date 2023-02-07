// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package subject

import (
	"time"

	"entgo.io/ent"

	"github.com/seal-io/seal/pkg/dao/schema"
)

const (
	// Label holds the string label denoting the subject type in the database.
	Label = "subject"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the createtime field in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the updatetime field in the database.
	FieldUpdateTime = "update_time"
	// FieldKind holds the string denoting the kind field in the database.
	FieldKind = "kind"
	// FieldGroup holds the string denoting the group field in the database.
	FieldGroup = "group"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldDescription holds the string denoting the description field in the database.
	FieldDescription = "description"
	// FieldMountTo holds the string denoting the mountto field in the database.
	FieldMountTo = "mount_to"
	// FieldLoginTo holds the string denoting the loginto field in the database.
	FieldLoginTo = "login_to"
	// FieldRoles holds the string denoting the roles field in the database.
	FieldRoles = "roles"
	// FieldPaths holds the string denoting the paths field in the database.
	FieldPaths = "paths"
	// FieldBuiltin holds the string denoting the builtin field in the database.
	FieldBuiltin = "builtin"
	// Table holds the table name of the subject in the database.
	Table = "subjects"
)

// Columns holds all SQL columns for subject fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldKind,
	FieldGroup,
	FieldName,
	FieldDescription,
	FieldMountTo,
	FieldLoginTo,
	FieldRoles,
	FieldPaths,
	FieldBuiltin,
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
	// DefaultKind holds the default value on creation for the "kind" field.
	DefaultKind string
	// DefaultGroup holds the default value on creation for the "group" field.
	DefaultGroup string
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
	// DefaultDescription holds the default value on creation for the "description" field.
	DefaultDescription string
	// DefaultMountTo holds the default value on creation for the "mountTo" field.
	DefaultMountTo bool
	// DefaultLoginTo holds the default value on creation for the "loginTo" field.
	DefaultLoginTo bool
	// DefaultRoles holds the default value on creation for the "roles" field.
	DefaultRoles schema.SubjectRoles
	// DefaultPaths holds the default value on creation for the "paths" field.
	DefaultPaths []string
	// DefaultBuiltin holds the default value on creation for the "builtin" field.
	DefaultBuiltin bool
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
