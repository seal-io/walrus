// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "ent". DO NOT EDIT.

package service

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the service type in the database.
	Label = "service"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldDescription holds the string denoting the description field in the database.
	FieldDescription = "description"
	// FieldLabels holds the string denoting the labels field in the database.
	FieldLabels = "labels"
	// FieldAnnotations holds the string denoting the annotations field in the database.
	FieldAnnotations = "annotations"
	// FieldCreateTime holds the string denoting the createtime field in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the updatetime field in the database.
	FieldUpdateTime = "update_time"
	// FieldProjectID holds the string denoting the projectid field in the database.
	FieldProjectID = "project_id"
	// FieldStatus holds the string denoting the status field in the database.
	FieldStatus = "status"
	// FieldEnvironmentID holds the string denoting the environmentid field in the database.
	FieldEnvironmentID = "environment_id"
	// FieldTemplate holds the string denoting the template field in the database.
	FieldTemplate = "template"
	// FieldAttributes holds the string denoting the attributes field in the database.
	FieldAttributes = "attributes"
	// FieldDriftResult holds the string denoting the driftresult field in the database.
	FieldDriftResult = "drift_result"
	// EdgeProject holds the string denoting the project edge name in mutations.
	EdgeProject = "project"
	// EdgeEnvironment holds the string denoting the environment edge name in mutations.
	EdgeEnvironment = "environment"
	// EdgeRevisions holds the string denoting the revisions edge name in mutations.
	EdgeRevisions = "revisions"
	// EdgeResources holds the string denoting the resources edge name in mutations.
	EdgeResources = "resources"
	// EdgeDependencies holds the string denoting the dependencies edge name in mutations.
	EdgeDependencies = "dependencies"
	// Table holds the table name of the service in the database.
	Table = "services"
	// ProjectTable is the table that holds the project relation/edge.
	ProjectTable = "services"
	// ProjectInverseTable is the table name for the Project entity.
	// It exists in this package in order to avoid circular dependency with the "project" package.
	ProjectInverseTable = "projects"
	// ProjectColumn is the table column denoting the project relation/edge.
	ProjectColumn = "project_id"
	// EnvironmentTable is the table that holds the environment relation/edge.
	EnvironmentTable = "services"
	// EnvironmentInverseTable is the table name for the Environment entity.
	// It exists in this package in order to avoid circular dependency with the "environment" package.
	EnvironmentInverseTable = "environments"
	// EnvironmentColumn is the table column denoting the environment relation/edge.
	EnvironmentColumn = "environment_id"
	// RevisionsTable is the table that holds the revisions relation/edge.
	RevisionsTable = "service_revisions"
	// RevisionsInverseTable is the table name for the ServiceRevision entity.
	// It exists in this package in order to avoid circular dependency with the "servicerevision" package.
	RevisionsInverseTable = "service_revisions"
	// RevisionsColumn is the table column denoting the revisions relation/edge.
	RevisionsColumn = "service_id"
	// ResourcesTable is the table that holds the resources relation/edge.
	ResourcesTable = "service_resources"
	// ResourcesInverseTable is the table name for the ServiceResource entity.
	// It exists in this package in order to avoid circular dependency with the "serviceresource" package.
	ResourcesInverseTable = "service_resources"
	// ResourcesColumn is the table column denoting the resources relation/edge.
	ResourcesColumn = "service_id"
	// DependenciesTable is the table that holds the dependencies relation/edge.
	DependenciesTable = "service_relationships"
	// DependenciesInverseTable is the table name for the ServiceRelationship entity.
	// It exists in this package in order to avoid circular dependency with the "servicerelationship" package.
	DependenciesInverseTable = "service_relationships"
	// DependenciesColumn is the table column denoting the dependencies relation/edge.
	DependenciesColumn = "service_id"
)

// Columns holds all SQL columns for service fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldDescription,
	FieldLabels,
	FieldAnnotations,
	FieldCreateTime,
	FieldUpdateTime,
	FieldProjectID,
	FieldStatus,
	FieldEnvironmentID,
	FieldTemplate,
	FieldAttributes,
	FieldDriftResult,
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
	Hooks        [3]ent.Hook
	Interceptors [1]ent.Interceptor
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
	// DefaultLabels holds the default value on creation for the "labels" field.
	DefaultLabels map[string]string
	// DefaultAnnotations holds the default value on creation for the "annotations" field.
	DefaultAnnotations map[string]string
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

// OrderOption defines the ordering options for the Service queries.
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

// ByEnvironmentID orders the results by the environmentID field.
func ByEnvironmentID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldEnvironmentID, opts...).ToFunc()
}

// ByAttributes orders the results by the attributes field.
func ByAttributes(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldAttributes, opts...).ToFunc()
}

// ByProjectField orders the results by project field.
func ByProjectField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newProjectStep(), sql.OrderByField(field, opts...))
	}
}

// ByEnvironmentField orders the results by environment field.
func ByEnvironmentField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newEnvironmentStep(), sql.OrderByField(field, opts...))
	}
}

// ByRevisionsCount orders the results by revisions count.
func ByRevisionsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newRevisionsStep(), opts...)
	}
}

// ByRevisions orders the results by revisions terms.
func ByRevisions(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newRevisionsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByResourcesCount orders the results by resources count.
func ByResourcesCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newResourcesStep(), opts...)
	}
}

// ByResources orders the results by resources terms.
func ByResources(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newResourcesStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByDependenciesCount orders the results by dependencies count.
func ByDependenciesCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newDependenciesStep(), opts...)
	}
}

// ByDependencies orders the results by dependencies terms.
func ByDependencies(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newDependenciesStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newProjectStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ProjectInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, ProjectTable, ProjectColumn),
	)
}
func newEnvironmentStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(EnvironmentInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, EnvironmentTable, EnvironmentColumn),
	)
}
func newRevisionsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(RevisionsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, RevisionsTable, RevisionsColumn),
	)
}
func newResourcesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ResourcesInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, ResourcesTable, ResourcesColumn),
	)
}
func newDependenciesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(DependenciesInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, DependenciesTable, DependenciesColumn),
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
