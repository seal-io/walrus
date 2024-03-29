// SPDX-FileCopyrightText: 2024 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "walrus". DO NOT EDIT.

package resourcestate

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"golang.org/x/exp/slices"
)

const (
	// Label holds the string label denoting the resourcestate type in the database.
	Label = "resource_state"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldData holds the string denoting the data field in the database.
	FieldData = "data"
	// FieldResourceID holds the string denoting the resource_id field in the database.
	FieldResourceID = "resource_id"
	// EdgeResource holds the string denoting the resource edge name in mutations.
	EdgeResource = "resource"
	// Table holds the table name of the resourcestate in the database.
	Table = "resource_states"
	// ResourceTable is the table that holds the resource relation/edge.
	ResourceTable = "resource_states"
	// ResourceInverseTable is the table name for the Resource entity.
	// It exists in this package in order to avoid circular dependency with the "resource" package.
	ResourceInverseTable = "resources"
	// ResourceColumn is the table column denoting the resource relation/edge.
	ResourceColumn = "resource_id"
)

// Columns holds all SQL columns for resourcestate fields.
var Columns = []string{
	FieldID,
	FieldData,
	FieldResourceID,
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
	Hooks [1]ent.Hook
	// DefaultData holds the default value on creation for the "data" field.
	DefaultData string
	// ResourceIDValidator is a validator for the "resource_id" field. It is called by the builders before save.
	ResourceIDValidator func(string) error
)

// OrderOption defines the ordering options for the ResourceState queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByData orders the results by the data field.
func ByData(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldData, opts...).ToFunc()
}

// ByResourceID orders the results by the resource_id field.
func ByResourceID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldResourceID, opts...).ToFunc()
}

// ByResourceField orders the results by resource field.
func ByResourceField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newResourceStep(), sql.OrderByField(field, opts...))
	}
}
func newResourceStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ResourceInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2O, true, ResourceTable, ResourceColumn),
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
