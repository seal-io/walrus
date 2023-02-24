// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package application

import (
	"time"

	"entgo.io/ent"
)

const (
	// Label holds the string label denoting the application type in the database.
	Label = "application"
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
	// FieldProjectID holds the string denoting the projectid field in the database.
	FieldProjectID = "project_id"
	// FieldEnvironmentID holds the string denoting the environmentid field in the database.
	FieldEnvironmentID = "environment_id"
	// EdgeProject holds the string denoting the project edge name in mutations.
	EdgeProject = "project"
	// EdgeEnvironment holds the string denoting the environment edge name in mutations.
	EdgeEnvironment = "environment"
	// EdgeResources holds the string denoting the resources edge name in mutations.
	EdgeResources = "resources"
	// EdgeRevisions holds the string denoting the revisions edge name in mutations.
	EdgeRevisions = "revisions"
	// EdgeModules holds the string denoting the modules edge name in mutations.
	EdgeModules = "modules"
	// EdgeApplicationModuleRelationships holds the string denoting the applicationmodulerelationships edge name in mutations.
	EdgeApplicationModuleRelationships = "applicationModuleRelationships"
	// Table holds the table name of the application in the database.
	Table = "applications"
	// ProjectTable is the table that holds the project relation/edge.
	ProjectTable = "applications"
	// ProjectInverseTable is the table name for the Project entity.
	// It exists in this package in order to avoid circular dependency with the "project" package.
	ProjectInverseTable = "projects"
	// ProjectColumn is the table column denoting the project relation/edge.
	ProjectColumn = "project_id"
	// EnvironmentTable is the table that holds the environment relation/edge.
	EnvironmentTable = "applications"
	// EnvironmentInverseTable is the table name for the Environment entity.
	// It exists in this package in order to avoid circular dependency with the "environment" package.
	EnvironmentInverseTable = "environments"
	// EnvironmentColumn is the table column denoting the environment relation/edge.
	EnvironmentColumn = "environment_id"
	// ResourcesTable is the table that holds the resources relation/edge.
	ResourcesTable = "application_resources"
	// ResourcesInverseTable is the table name for the ApplicationResource entity.
	// It exists in this package in order to avoid circular dependency with the "applicationresource" package.
	ResourcesInverseTable = "application_resources"
	// ResourcesColumn is the table column denoting the resources relation/edge.
	ResourcesColumn = "application_id"
	// RevisionsTable is the table that holds the revisions relation/edge.
	RevisionsTable = "application_revisions"
	// RevisionsInverseTable is the table name for the ApplicationRevision entity.
	// It exists in this package in order to avoid circular dependency with the "applicationrevision" package.
	RevisionsInverseTable = "application_revisions"
	// RevisionsColumn is the table column denoting the revisions relation/edge.
	RevisionsColumn = "application_id"
	// ModulesTable is the table that holds the modules relation/edge. The primary key declared below.
	ModulesTable = "application_module_relationships"
	// ModulesInverseTable is the table name for the Module entity.
	// It exists in this package in order to avoid circular dependency with the "module" package.
	ModulesInverseTable = "modules"
	// ApplicationModuleRelationshipsTable is the table that holds the applicationModuleRelationships relation/edge.
	ApplicationModuleRelationshipsTable = "application_module_relationships"
	// ApplicationModuleRelationshipsInverseTable is the table name for the ApplicationModuleRelationship entity.
	// It exists in this package in order to avoid circular dependency with the "applicationmodulerelationship" package.
	ApplicationModuleRelationshipsInverseTable = "application_module_relationships"
	// ApplicationModuleRelationshipsColumn is the table column denoting the applicationModuleRelationships relation/edge.
	ApplicationModuleRelationshipsColumn = "application_id"
)

// Columns holds all SQL columns for application fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldDescription,
	FieldLabels,
	FieldCreateTime,
	FieldUpdateTime,
	FieldProjectID,
	FieldEnvironmentID,
}

var (
	// ModulesPrimaryKey and ModulesColumn2 are the table columns denoting the
	// primary key for the modules relation (M2M).
	ModulesPrimaryKey = []string{"application_id", "module_id"}
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
	// ProjectIDValidator is a validator for the "projectID" field. It is called by the builders before save.
	ProjectIDValidator func(string) error
	// EnvironmentIDValidator is a validator for the "environmentID" field. It is called by the builders before save.
	EnvironmentIDValidator func(string) error
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
