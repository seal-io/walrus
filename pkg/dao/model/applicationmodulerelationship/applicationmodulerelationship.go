// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package applicationmodulerelationship

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the applicationmodulerelationship type in the database.
	Label = "application_module_relationship"
	// FieldCreateTime holds the string denoting the createtime field in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the updatetime field in the database.
	FieldUpdateTime = "update_time"
	// FieldApplicationID holds the string denoting the application_id field in the database.
	FieldApplicationID = "application_id"
	// FieldModuleID holds the string denoting the module_id field in the database.
	FieldModuleID = "module_id"
	// FieldVersion holds the string denoting the version field in the database.
	FieldVersion = "version"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldAttributes holds the string denoting the attributes field in the database.
	FieldAttributes = "attributes"
	// EdgeApplication holds the string denoting the application edge name in mutations.
	EdgeApplication = "application"
	// EdgeModule holds the string denoting the module edge name in mutations.
	EdgeModule = "module"
	// ApplicationFieldID holds the string denoting the ID field of the Application.
	ApplicationFieldID = "id"
	// ModuleFieldID holds the string denoting the ID field of the Module.
	ModuleFieldID = "id"
	// Table holds the table name of the applicationmodulerelationship in the database.
	Table = "application_module_relationships"
	// ApplicationTable is the table that holds the application relation/edge.
	ApplicationTable = "application_module_relationships"
	// ApplicationInverseTable is the table name for the Application entity.
	// It exists in this package in order to avoid circular dependency with the "application" package.
	ApplicationInverseTable = "applications"
	// ApplicationColumn is the table column denoting the application relation/edge.
	ApplicationColumn = "application_id"
	// ModuleTable is the table that holds the module relation/edge.
	ModuleTable = "application_module_relationships"
	// ModuleInverseTable is the table name for the Module entity.
	// It exists in this package in order to avoid circular dependency with the "module" package.
	ModuleInverseTable = "modules"
	// ModuleColumn is the table column denoting the module relation/edge.
	ModuleColumn = "module_id"
)

// Columns holds all SQL columns for applicationmodulerelationship fields.
var Columns = []string{
	FieldCreateTime,
	FieldUpdateTime,
	FieldApplicationID,
	FieldModuleID,
	FieldVersion,
	FieldName,
	FieldAttributes,
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

var (
	// DefaultCreateTime holds the default value on creation for the "createTime" field.
	DefaultCreateTime func() time.Time
	// DefaultUpdateTime holds the default value on creation for the "updateTime" field.
	DefaultUpdateTime func() time.Time
	// UpdateDefaultUpdateTime holds the default value on update for the "updateTime" field.
	UpdateDefaultUpdateTime func() time.Time
	// ApplicationIDValidator is a validator for the "application_id" field. It is called by the builders before save.
	ApplicationIDValidator func(string) error
	// ModuleIDValidator is a validator for the "module_id" field. It is called by the builders before save.
	ModuleIDValidator func(string) error
	// VersionValidator is a validator for the "version" field. It is called by the builders before save.
	VersionValidator func(string) error
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
)

// OrderOption defines the ordering options for the ApplicationModuleRelationship queries.
type OrderOption func(*sql.Selector)

// ByCreateTime orders the results by the createTime field.
func ByCreateTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreateTime, opts...).ToFunc()
}

// ByUpdateTime orders the results by the updateTime field.
func ByUpdateTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdateTime, opts...).ToFunc()
}

// ByApplicationID orders the results by the application_id field.
func ByApplicationID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldApplicationID, opts...).ToFunc()
}

// ByModuleID orders the results by the module_id field.
func ByModuleID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldModuleID, opts...).ToFunc()
}

// ByVersion orders the results by the version field.
func ByVersion(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldVersion, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByAttributes orders the results by the attributes field.
func ByAttributes(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldAttributes, opts...).ToFunc()
}

// ByApplicationField orders the results by application field.
func ByApplicationField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newApplicationStep(), sql.OrderByField(field, opts...))
	}
}

// ByModuleField orders the results by module field.
func ByModuleField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newModuleStep(), sql.OrderByField(field, opts...))
	}
}
func newApplicationStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, ApplicationColumn),
		sqlgraph.To(ApplicationInverseTable, ApplicationFieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, ApplicationTable, ApplicationColumn),
	)
}
func newModuleStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, ModuleColumn),
		sqlgraph.To(ModuleInverseTable, ModuleFieldID),
		sqlgraph.Edge(sqlgraph.M2O, false, ModuleTable, ModuleColumn),
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
