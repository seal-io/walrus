// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "ent". DO NOT EDIT.

package project

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the project type in the database.
	Label = "project"
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
	// EdgeEnvironments holds the string denoting the environments edge name in mutations.
	EdgeEnvironments = "environments"
	// EdgeConnectors holds the string denoting the connectors edge name in mutations.
	EdgeConnectors = "connectors"
	// EdgeSecrets holds the string denoting the secrets edge name in mutations.
	EdgeSecrets = "secrets"
	// EdgeSubjectRoles holds the string denoting the subjectroles edge name in mutations.
	EdgeSubjectRoles = "subjectRoles"
	// EdgeServices holds the string denoting the services edge name in mutations.
	EdgeServices = "services"
	// EdgeServiceRevisions holds the string denoting the servicerevisions edge name in mutations.
	EdgeServiceRevisions = "serviceRevisions"
	// Table holds the table name of the project in the database.
	Table = "projects"
	// EnvironmentsTable is the table that holds the environments relation/edge.
	EnvironmentsTable = "environments"
	// EnvironmentsInverseTable is the table name for the Environment entity.
	// It exists in this package in order to avoid circular dependency with the "environment" package.
	EnvironmentsInverseTable = "environments"
	// EnvironmentsColumn is the table column denoting the environments relation/edge.
	EnvironmentsColumn = "project_id"
	// ConnectorsTable is the table that holds the connectors relation/edge.
	ConnectorsTable = "connectors"
	// ConnectorsInverseTable is the table name for the Connector entity.
	// It exists in this package in order to avoid circular dependency with the "connector" package.
	ConnectorsInverseTable = "connectors"
	// ConnectorsColumn is the table column denoting the connectors relation/edge.
	ConnectorsColumn = "project_id"
	// SecretsTable is the table that holds the secrets relation/edge.
	SecretsTable = "secrets"
	// SecretsInverseTable is the table name for the Secret entity.
	// It exists in this package in order to avoid circular dependency with the "secret" package.
	SecretsInverseTable = "secrets"
	// SecretsColumn is the table column denoting the secrets relation/edge.
	SecretsColumn = "project_id"
	// SubjectRolesTable is the table that holds the subjectRoles relation/edge.
	SubjectRolesTable = "subject_role_relationships"
	// SubjectRolesInverseTable is the table name for the SubjectRoleRelationship entity.
	// It exists in this package in order to avoid circular dependency with the "subjectrolerelationship" package.
	SubjectRolesInverseTable = "subject_role_relationships"
	// SubjectRolesColumn is the table column denoting the subjectRoles relation/edge.
	SubjectRolesColumn = "project_id"
	// ServicesTable is the table that holds the services relation/edge.
	ServicesTable = "services"
	// ServicesInverseTable is the table name for the Service entity.
	// It exists in this package in order to avoid circular dependency with the "service" package.
	ServicesInverseTable = "services"
	// ServicesColumn is the table column denoting the services relation/edge.
	ServicesColumn = "project_id"
	// ServiceRevisionsTable is the table that holds the serviceRevisions relation/edge.
	ServiceRevisionsTable = "service_revisions"
	// ServiceRevisionsInverseTable is the table name for the ServiceRevision entity.
	// It exists in this package in order to avoid circular dependency with the "servicerevision" package.
	ServiceRevisionsInverseTable = "service_revisions"
	// ServiceRevisionsColumn is the table column denoting the serviceRevisions relation/edge.
	ServiceRevisionsColumn = "project_id"
)

// Columns holds all SQL columns for project fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldDescription,
	FieldLabels,
	FieldAnnotations,
	FieldCreateTime,
	FieldUpdateTime,
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
	// DefaultAnnotations holds the default value on creation for the "annotations" field.
	DefaultAnnotations map[string]string
	// DefaultCreateTime holds the default value on creation for the "createTime" field.
	DefaultCreateTime func() time.Time
	// DefaultUpdateTime holds the default value on creation for the "updateTime" field.
	DefaultUpdateTime func() time.Time
	// UpdateDefaultUpdateTime holds the default value on update for the "updateTime" field.
	UpdateDefaultUpdateTime func() time.Time
)

// OrderOption defines the ordering options for the Project queries.
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

// ByEnvironmentsCount orders the results by environments count.
func ByEnvironmentsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newEnvironmentsStep(), opts...)
	}
}

// ByEnvironments orders the results by environments terms.
func ByEnvironments(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newEnvironmentsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByConnectorsCount orders the results by connectors count.
func ByConnectorsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newConnectorsStep(), opts...)
	}
}

// ByConnectors orders the results by connectors terms.
func ByConnectors(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newConnectorsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// BySecretsCount orders the results by secrets count.
func BySecretsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newSecretsStep(), opts...)
	}
}

// BySecrets orders the results by secrets terms.
func BySecrets(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newSecretsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// BySubjectRolesCount orders the results by subjectRoles count.
func BySubjectRolesCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newSubjectRolesStep(), opts...)
	}
}

// BySubjectRoles orders the results by subjectRoles terms.
func BySubjectRoles(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newSubjectRolesStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByServicesCount orders the results by services count.
func ByServicesCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newServicesStep(), opts...)
	}
}

// ByServices orders the results by services terms.
func ByServices(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newServicesStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByServiceRevisionsCount orders the results by serviceRevisions count.
func ByServiceRevisionsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newServiceRevisionsStep(), opts...)
	}
}

// ByServiceRevisions orders the results by serviceRevisions terms.
func ByServiceRevisions(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newServiceRevisionsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newEnvironmentsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(EnvironmentsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, EnvironmentsTable, EnvironmentsColumn),
	)
}
func newConnectorsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ConnectorsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, ConnectorsTable, ConnectorsColumn),
	)
}
func newSecretsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(SecretsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, SecretsTable, SecretsColumn),
	)
}
func newSubjectRolesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(SubjectRolesInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, SubjectRolesTable, SubjectRolesColumn),
	)
}
func newServicesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ServicesInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, ServicesTable, ServicesColumn),
	)
}
func newServiceRevisionsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ServiceRevisionsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, ServiceRevisionsTable, ServiceRevisionsColumn),
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
