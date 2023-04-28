// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package applicationresource

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the applicationresource type in the database.
	Label = "application_resource"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the createtime field in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the updatetime field in the database.
	FieldUpdateTime = "update_time"
	// FieldInstanceID holds the string denoting the instanceid field in the database.
	FieldInstanceID = "instance_id"
	// FieldConnectorID holds the string denoting the connectorid field in the database.
	FieldConnectorID = "connector_id"
	// FieldCompositionID holds the string denoting the compositionid field in the database.
	FieldCompositionID = "composition_id"
	// FieldModule holds the string denoting the module field in the database.
	FieldModule = "module"
	// FieldMode holds the string denoting the mode field in the database.
	FieldMode = "mode"
	// FieldType holds the string denoting the type field in the database.
	FieldType = "type"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldDeployerType holds the string denoting the deployertype field in the database.
	FieldDeployerType = "deployer_type"
	// FieldStatus holds the string denoting the status field in the database.
	FieldStatus = "status"
	// EdgeInstance holds the string denoting the instance edge name in mutations.
	EdgeInstance = "instance"
	// EdgeConnector holds the string denoting the connector edge name in mutations.
	EdgeConnector = "connector"
	// EdgeComposition holds the string denoting the composition edge name in mutations.
	EdgeComposition = "composition"
	// EdgeComponents holds the string denoting the components edge name in mutations.
	EdgeComponents = "components"
	// Table holds the table name of the applicationresource in the database.
	Table = "application_resources"
	// InstanceTable is the table that holds the instance relation/edge.
	InstanceTable = "application_resources"
	// InstanceInverseTable is the table name for the ApplicationInstance entity.
	// It exists in this package in order to avoid circular dependency with the "applicationinstance" package.
	InstanceInverseTable = "application_instances"
	// InstanceColumn is the table column denoting the instance relation/edge.
	InstanceColumn = "instance_id"
	// ConnectorTable is the table that holds the connector relation/edge.
	ConnectorTable = "application_resources"
	// ConnectorInverseTable is the table name for the Connector entity.
	// It exists in this package in order to avoid circular dependency with the "connector" package.
	ConnectorInverseTable = "connectors"
	// ConnectorColumn is the table column denoting the connector relation/edge.
	ConnectorColumn = "connector_id"
	// CompositionTable is the table that holds the composition relation/edge.
	CompositionTable = "application_resources"
	// CompositionColumn is the table column denoting the composition relation/edge.
	CompositionColumn = "composition_id"
	// ComponentsTable is the table that holds the components relation/edge.
	ComponentsTable = "application_resources"
	// ComponentsColumn is the table column denoting the components relation/edge.
	ComponentsColumn = "composition_id"
)

// Columns holds all SQL columns for applicationresource fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldInstanceID,
	FieldConnectorID,
	FieldCompositionID,
	FieldModule,
	FieldMode,
	FieldType,
	FieldName,
	FieldDeployerType,
	FieldStatus,
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
	Hooks        [1]ent.Hook
	Interceptors [1]ent.Interceptor
	// DefaultCreateTime holds the default value on creation for the "createTime" field.
	DefaultCreateTime func() time.Time
	// DefaultUpdateTime holds the default value on creation for the "updateTime" field.
	DefaultUpdateTime func() time.Time
	// UpdateDefaultUpdateTime holds the default value on update for the "updateTime" field.
	UpdateDefaultUpdateTime func() time.Time
	// InstanceIDValidator is a validator for the "instanceID" field. It is called by the builders before save.
	InstanceIDValidator func(string) error
	// ConnectorIDValidator is a validator for the "connectorID" field. It is called by the builders before save.
	ConnectorIDValidator func(string) error
	// ModuleValidator is a validator for the "module" field. It is called by the builders before save.
	ModuleValidator func(string) error
	// ModeValidator is a validator for the "mode" field. It is called by the builders before save.
	ModeValidator func(string) error
	// TypeValidator is a validator for the "type" field. It is called by the builders before save.
	TypeValidator func(string) error
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
	// DeployerTypeValidator is a validator for the "deployerType" field. It is called by the builders before save.
	DeployerTypeValidator func(string) error
)

// OrderOption defines the ordering options for the ApplicationResource queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByCreateTime orders the results by the createTime field.
func ByCreateTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreateTime, opts...).ToFunc()
}

// ByUpdateTime orders the results by the updateTime field.
func ByUpdateTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdateTime, opts...).ToFunc()
}

// ByInstanceID orders the results by the instanceID field.
func ByInstanceID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldInstanceID, opts...).ToFunc()
}

// ByConnectorID orders the results by the connectorID field.
func ByConnectorID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldConnectorID, opts...).ToFunc()
}

// ByCompositionID orders the results by the compositionID field.
func ByCompositionID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCompositionID, opts...).ToFunc()
}

// ByModule orders the results by the module field.
func ByModule(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldModule, opts...).ToFunc()
}

// ByMode orders the results by the mode field.
func ByMode(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldMode, opts...).ToFunc()
}

// ByType orders the results by the type field.
func ByType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldType, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByDeployerType orders the results by the deployerType field.
func ByDeployerType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDeployerType, opts...).ToFunc()
}

// ByInstanceField orders the results by instance field.
func ByInstanceField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newInstanceStep(), sql.OrderByField(field, opts...))
	}
}

// ByConnectorField orders the results by connector field.
func ByConnectorField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newConnectorStep(), sql.OrderByField(field, opts...))
	}
}

// ByCompositionField orders the results by composition field.
func ByCompositionField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newCompositionStep(), sql.OrderByField(field, opts...))
	}
}

// ByComponentsCount orders the results by components count.
func ByComponentsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newComponentsStep(), opts...)
	}
}

// ByComponents orders the results by components terms.
func ByComponents(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newComponentsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newInstanceStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(InstanceInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, InstanceTable, InstanceColumn),
	)
}
func newConnectorStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ConnectorInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, ConnectorTable, ConnectorColumn),
	)
}
func newCompositionStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(Table, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, CompositionTable, CompositionColumn),
	)
}
func newComponentsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(Table, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, ComponentsTable, ComponentsColumn),
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
