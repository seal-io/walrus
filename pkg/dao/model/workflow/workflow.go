// SPDX-FileCopyrightText: 2024 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "walrus". DO NOT EDIT.

package workflow

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"golang.org/x/exp/slices"
)

const (
	// Label holds the string label denoting the workflow type in the database.
	Label = "workflow"
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
	// FieldCreateTime holds the string denoting the create_time field in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the update_time field in the database.
	FieldUpdateTime = "update_time"
	// FieldProjectID holds the string denoting the project_id field in the database.
	FieldProjectID = "project_id"
	// FieldEnvironmentID holds the string denoting the environment_id field in the database.
	FieldEnvironmentID = "environment_id"
	// FieldType holds the string denoting the type field in the database.
	FieldType = "type"
	// FieldParallelism holds the string denoting the parallelism field in the database.
	FieldParallelism = "parallelism"
	// FieldTimeout holds the string denoting the timeout field in the database.
	FieldTimeout = "timeout"
	// FieldVersion holds the string denoting the version field in the database.
	FieldVersion = "version"
	// FieldVariables holds the string denoting the variables field in the database.
	FieldVariables = "variables"
	// EdgeProject holds the string denoting the project edge name in mutations.
	EdgeProject = "project"
	// EdgeStages holds the string denoting the stages edge name in mutations.
	EdgeStages = "stages"
	// EdgeExecutions holds the string denoting the executions edge name in mutations.
	EdgeExecutions = "executions"
	// Table holds the table name of the workflow in the database.
	Table = "workflows"
	// ProjectTable is the table that holds the project relation/edge.
	ProjectTable = "workflows"
	// ProjectInverseTable is the table name for the Project entity.
	// It exists in this package in order to avoid circular dependency with the "project" package.
	ProjectInverseTable = "projects"
	// ProjectColumn is the table column denoting the project relation/edge.
	ProjectColumn = "project_id"
	// StagesTable is the table that holds the stages relation/edge.
	StagesTable = "workflow_stages"
	// StagesInverseTable is the table name for the WorkflowStage entity.
	// It exists in this package in order to avoid circular dependency with the "workflowstage" package.
	StagesInverseTable = "workflow_stages"
	// StagesColumn is the table column denoting the stages relation/edge.
	StagesColumn = "workflow_id"
	// ExecutionsTable is the table that holds the executions relation/edge.
	ExecutionsTable = "workflow_executions"
	// ExecutionsInverseTable is the table name for the WorkflowExecution entity.
	// It exists in this package in order to avoid circular dependency with the "workflowexecution" package.
	ExecutionsInverseTable = "workflow_executions"
	// ExecutionsColumn is the table column denoting the executions relation/edge.
	ExecutionsColumn = "workflow_id"
)

// Columns holds all SQL columns for workflow fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldDescription,
	FieldLabels,
	FieldAnnotations,
	FieldCreateTime,
	FieldUpdateTime,
	FieldProjectID,
	FieldEnvironmentID,
	FieldType,
	FieldParallelism,
	FieldTimeout,
	FieldVersion,
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
//	import _ "github.com/seal-io/walrus/pkg/dao/model/runtime"
var (
	Hooks        [1]ent.Hook
	Interceptors [1]ent.Interceptor
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
	// DefaultLabels holds the default value on creation for the "labels" field.
	DefaultLabels map[string]string
	// DefaultAnnotations holds the default value on creation for the "annotations" field.
	DefaultAnnotations map[string]string
	// DefaultCreateTime holds the default value on creation for the "create_time" field.
	DefaultCreateTime func() time.Time
	// DefaultUpdateTime holds the default value on creation for the "update_time" field.
	DefaultUpdateTime func() time.Time
	// UpdateDefaultUpdateTime holds the default value on update for the "update_time" field.
	UpdateDefaultUpdateTime func() time.Time
	// ProjectIDValidator is a validator for the "project_id" field. It is called by the builders before save.
	ProjectIDValidator func(string) error
	// TypeValidator is a validator for the "type" field. It is called by the builders before save.
	TypeValidator func(string) error
	// DefaultParallelism holds the default value on creation for the "parallelism" field.
	DefaultParallelism int
	// ParallelismValidator is a validator for the "parallelism" field. It is called by the builders before save.
	ParallelismValidator func(int) error
	// DefaultTimeout holds the default value on creation for the "timeout" field.
	DefaultTimeout int
	// TimeoutValidator is a validator for the "timeout" field. It is called by the builders before save.
	TimeoutValidator func(int) error
	// DefaultVersion holds the default value on creation for the "version" field.
	DefaultVersion int
	// VersionValidator is a validator for the "version" field. It is called by the builders before save.
	VersionValidator func(int) error
)

// OrderOption defines the ordering options for the Workflow queries.
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

// ByCreateTime orders the results by the create_time field.
func ByCreateTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreateTime, opts...).ToFunc()
}

// ByUpdateTime orders the results by the update_time field.
func ByUpdateTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdateTime, opts...).ToFunc()
}

// ByProjectID orders the results by the project_id field.
func ByProjectID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldProjectID, opts...).ToFunc()
}

// ByEnvironmentID orders the results by the environment_id field.
func ByEnvironmentID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldEnvironmentID, opts...).ToFunc()
}

// ByType orders the results by the type field.
func ByType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldType, opts...).ToFunc()
}

// ByParallelism orders the results by the parallelism field.
func ByParallelism(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldParallelism, opts...).ToFunc()
}

// ByTimeout orders the results by the timeout field.
func ByTimeout(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldTimeout, opts...).ToFunc()
}

// ByVersion orders the results by the version field.
func ByVersion(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldVersion, opts...).ToFunc()
}

// ByProjectField orders the results by project field.
func ByProjectField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newProjectStep(), sql.OrderByField(field, opts...))
	}
}

// ByStagesCount orders the results by stages count.
func ByStagesCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newStagesStep(), opts...)
	}
}

// ByStages orders the results by stages terms.
func ByStages(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newStagesStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByExecutionsCount orders the results by executions count.
func ByExecutionsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newExecutionsStep(), opts...)
	}
}

// ByExecutions orders the results by executions terms.
func ByExecutions(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newExecutionsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newProjectStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ProjectInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, ProjectTable, ProjectColumn),
	)
}
func newStagesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(StagesInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, StagesTable, StagesColumn),
	)
}
func newExecutionsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ExecutionsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, ExecutionsTable, ExecutionsColumn),
	)
}

// WithoutFields returns the fields ignored the given list.
func WithoutFields(ignores ...string) []string {
	if len(ignores) == 0 {
		return slices.Clone(Columns)
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
