// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "seal". DO NOT EDIT.

package serviceresourcerelationship

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"golang.org/x/exp/slices"
)

const (
	// Label holds the string label denoting the serviceresourcerelationship type in the database.
	Label = "service_resource_relationship"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the create_time field in the database.
	FieldCreateTime = "create_time"
	// FieldServiceResourceID holds the string denoting the service_resource_id field in the database.
	FieldServiceResourceID = "service_resource_id"
	// FieldDependencyID holds the string denoting the dependency_id field in the database.
	FieldDependencyID = "dependency_id"
	// FieldType holds the string denoting the type field in the database.
	FieldType = "type"
	// EdgeServiceResource holds the string denoting the serviceresource edge name in mutations.
	EdgeServiceResource = "serviceResource"
	// EdgeDependency holds the string denoting the dependency edge name in mutations.
	EdgeDependency = "dependency"
	// Table holds the table name of the serviceresourcerelationship in the database.
	Table = "service_resource_relationships"
	// ServiceResourceTable is the table that holds the serviceResource relation/edge.
	ServiceResourceTable = "service_resource_relationships"
	// ServiceResourceInverseTable is the table name for the ServiceResource entity.
	// It exists in this package in order to avoid circular dependency with the "serviceresource" package.
	ServiceResourceInverseTable = "service_resources"
	// ServiceResourceColumn is the table column denoting the serviceResource relation/edge.
	ServiceResourceColumn = "service_resource_id"
	// DependencyTable is the table that holds the dependency relation/edge.
	DependencyTable = "service_resource_relationships"
	// DependencyInverseTable is the table name for the ServiceResource entity.
	// It exists in this package in order to avoid circular dependency with the "serviceresource" package.
	DependencyInverseTable = "service_resources"
	// DependencyColumn is the table column denoting the dependency relation/edge.
	DependencyColumn = "dependency_id"
)

// Columns holds all SQL columns for serviceresourcerelationship fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldServiceResourceID,
	FieldDependencyID,
	FieldType,
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
	// DefaultCreateTime holds the default value on creation for the "create_time" field.
	DefaultCreateTime func() time.Time
	// ServiceResourceIDValidator is a validator for the "service_resource_id" field. It is called by the builders before save.
	ServiceResourceIDValidator func(string) error
	// DependencyIDValidator is a validator for the "dependency_id" field. It is called by the builders before save.
	DependencyIDValidator func(string) error
	// TypeValidator is a validator for the "type" field. It is called by the builders before save.
	TypeValidator func(string) error
)

// OrderOption defines the ordering options for the ServiceResourceRelationship queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByCreateTime orders the results by the create_time field.
func ByCreateTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreateTime, opts...).ToFunc()
}

// ByServiceResourceID orders the results by the service_resource_id field.
func ByServiceResourceID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldServiceResourceID, opts...).ToFunc()
}

// ByDependencyID orders the results by the dependency_id field.
func ByDependencyID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDependencyID, opts...).ToFunc()
}

// ByType orders the results by the type field.
func ByType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldType, opts...).ToFunc()
}

// ByServiceResourceField orders the results by serviceResource field.
func ByServiceResourceField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newServiceResourceStep(), sql.OrderByField(field, opts...))
	}
}

// ByDependencyField orders the results by dependency field.
func ByDependencyField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newDependencyStep(), sql.OrderByField(field, opts...))
	}
}
func newServiceResourceStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ServiceResourceInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, ServiceResourceTable, ServiceResourceColumn),
	)
}
func newDependencyStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(DependencyInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, DependencyTable, DependencyColumn),
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
