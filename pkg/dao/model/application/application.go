// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package application

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
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
	// FieldVariables holds the string denoting the variables field in the database.
	FieldVariables = "variables"
	// EdgeProject holds the string denoting the project edge name in mutations.
	EdgeProject = "project"
	// EdgeInstances holds the string denoting the instances edge name in mutations.
	EdgeInstances = "instances"
	// EdgeModules holds the string denoting the modules edge name in mutations.
	EdgeModules = "modules"
	// Table holds the table name of the application in the database.
	Table = "applications"
	// ProjectTable is the table that holds the project relation/edge.
	ProjectTable = "applications"
	// ProjectInverseTable is the table name for the Project entity.
	// It exists in this package in order to avoid circular dependency with the "project" package.
	ProjectInverseTable = "projects"
	// ProjectColumn is the table column denoting the project relation/edge.
	ProjectColumn = "project_id"
	// InstancesTable is the table that holds the instances relation/edge.
	InstancesTable = "application_instances"
	// InstancesInverseTable is the table name for the ApplicationInstance entity.
	// It exists in this package in order to avoid circular dependency with the "applicationinstance" package.
	InstancesInverseTable = "application_instances"
	// InstancesColumn is the table column denoting the instances relation/edge.
	InstancesColumn = "application_id"
	// ModulesTable is the table that holds the modules relation/edge.
	ModulesTable = "application_module_relationships"
	// ModulesInverseTable is the table name for the ApplicationModuleRelationship entity.
	// It exists in this package in order to avoid circular dependency with the "applicationmodulerelationship" package.
	ModulesInverseTable = "application_module_relationships"
	// ModulesColumn is the table column denoting the modules relation/edge.
	ModulesColumn = "application_id"
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
)

// OrderOption defines the ordering options for the Application queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByDescription orders the results by the description field.
func ByDescription(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDescription, opts...).ToFunc()
}

// ByCreateTime orders the results by the createTime field.
func ByCreateTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreateTime, opts...).ToFunc()
}

// ByUpdateTime orders the results by the updateTime field.
func ByUpdateTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdateTime, opts...).ToFunc()
}

// ByProjectID orders the results by the projectID field.
func ByProjectID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldProjectID, opts...).ToFunc()
}

// ByVariables orders the results by the variables field.
func ByVariables(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldVariables, opts...).ToFunc()
}

// ByProjectField orders the results by project field.
func ByProjectField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newProjectStep(), sql.OrderByField(field, opts...))
	}
}

// ByInstancesCount orders the results by instances count.
func ByInstancesCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newInstancesStep(), opts...)
	}
}

// ByInstances orders the results by instances terms.
func ByInstances(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newInstancesStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByModulesCount orders the results by modules count.
func ByModulesCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newModulesStep(), opts...)
	}
}

// ByModules orders the results by modules terms.
func ByModules(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newModulesStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newProjectStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ProjectInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, ProjectTable, ProjectColumn),
	)
}
func newInstancesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(InstancesInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, InstancesTable, InstancesColumn),
	)
}
func newModulesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ModulesInverseTable, ModulesColumn),
		sqlgraph.Edge(sqlgraph.O2M, true, ModulesTable, ModulesColumn),
	)
}

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
