// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package module

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the module type in the database.
	Label = "module"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldStatus holds the string denoting the status field in the database.
	FieldStatus = "status"
	// FieldStatusMessage holds the string denoting the statusmessage field in the database.
	FieldStatusMessage = "status_message"
	// FieldCreateTime holds the string denoting the createtime field in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the updatetime field in the database.
	FieldUpdateTime = "update_time"
	// FieldDescription holds the string denoting the description field in the database.
	FieldDescription = "description"
	// FieldIcon holds the string denoting the icon field in the database.
	FieldIcon = "icon"
	// FieldLabels holds the string denoting the labels field in the database.
	FieldLabels = "labels"
	// FieldSource holds the string denoting the source field in the database.
	FieldSource = "source"
	// EdgeApplications holds the string denoting the applications edge name in mutations.
	EdgeApplications = "applications"
	// EdgeVersions holds the string denoting the versions edge name in mutations.
	EdgeVersions = "versions"
	// Table holds the table name of the module in the database.
	Table = "modules"
	// ApplicationsTable is the table that holds the applications relation/edge.
	ApplicationsTable = "application_module_relationships"
	// ApplicationsInverseTable is the table name for the ApplicationModuleRelationship entity.
	// It exists in this package in order to avoid circular dependency with the "applicationmodulerelationship" package.
	ApplicationsInverseTable = "application_module_relationships"
	// ApplicationsColumn is the table column denoting the applications relation/edge.
	ApplicationsColumn = "module_id"
	// VersionsTable is the table that holds the versions relation/edge.
	VersionsTable = "module_versions"
	// VersionsInverseTable is the table name for the ModuleVersion entity.
	// It exists in this package in order to avoid circular dependency with the "moduleversion" package.
	VersionsInverseTable = "module_versions"
	// VersionsColumn is the table column denoting the versions relation/edge.
	VersionsColumn = "module_id"
)

// Columns holds all SQL columns for module fields.
var Columns = []string{
	FieldID,
	FieldStatus,
	FieldStatusMessage,
	FieldCreateTime,
	FieldUpdateTime,
	FieldDescription,
	FieldIcon,
	FieldLabels,
	FieldSource,
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
	// DefaultLabels holds the default value on creation for the "labels" field.
	DefaultLabels map[string]string
	// SourceValidator is a validator for the "source" field. It is called by the builders before save.
	SourceValidator func(string) error
	// IDValidator is a validator for the "id" field. It is called by the builders before save.
	IDValidator func(string) error
)

// OrderOption defines the ordering options for the Module queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByStatus orders the results by the status field.
func ByStatus(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldStatus, opts...).ToFunc()
}

// ByStatusMessage orders the results by the statusMessage field.
func ByStatusMessage(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldStatusMessage, opts...).ToFunc()
}

// ByCreateTime orders the results by the createTime field.
func ByCreateTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreateTime, opts...).ToFunc()
}

// ByUpdateTime orders the results by the updateTime field.
func ByUpdateTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdateTime, opts...).ToFunc()
}

// ByDescription orders the results by the description field.
func ByDescription(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDescription, opts...).ToFunc()
}

// ByIcon orders the results by the icon field.
func ByIcon(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldIcon, opts...).ToFunc()
}

// BySource orders the results by the source field.
func BySource(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSource, opts...).ToFunc()
}

// ByApplicationsCount orders the results by applications count.
func ByApplicationsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newApplicationsStep(), opts...)
	}
}

// ByApplications orders the results by applications terms.
func ByApplications(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newApplicationsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByVersionsCount orders the results by versions count.
func ByVersionsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newVersionsStep(), opts...)
	}
}

// ByVersions orders the results by versions terms.
func ByVersions(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newVersionsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newApplicationsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ApplicationsInverseTable, ApplicationsColumn),
		sqlgraph.Edge(sqlgraph.O2M, true, ApplicationsTable, ApplicationsColumn),
	)
}
func newVersionsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(VersionsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, VersionsTable, VersionsColumn),
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
