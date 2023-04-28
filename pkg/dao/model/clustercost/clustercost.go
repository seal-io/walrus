// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package clustercost

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the clustercost type in the database.
	Label = "cluster_cost"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldStartTime holds the string denoting the starttime field in the database.
	FieldStartTime = "start_time"
	// FieldEndTime holds the string denoting the endtime field in the database.
	FieldEndTime = "end_time"
	// FieldMinutes holds the string denoting the minutes field in the database.
	FieldMinutes = "minutes"
	// FieldConnectorID holds the string denoting the connectorid field in the database.
	FieldConnectorID = "connector_id"
	// FieldClusterName holds the string denoting the clustername field in the database.
	FieldClusterName = "cluster_name"
	// FieldTotalCost holds the string denoting the totalcost field in the database.
	FieldTotalCost = "total_cost"
	// FieldCurrency holds the string denoting the currency field in the database.
	FieldCurrency = "currency"
	// FieldAllocationCost holds the string denoting the allocationcost field in the database.
	FieldAllocationCost = "allocation_cost"
	// FieldIdleCost holds the string denoting the idlecost field in the database.
	FieldIdleCost = "idle_cost"
	// FieldManagementCost holds the string denoting the managementcost field in the database.
	FieldManagementCost = "management_cost"
	// EdgeConnector holds the string denoting the connector edge name in mutations.
	EdgeConnector = "connector"
	// Table holds the table name of the clustercost in the database.
	Table = "cluster_costs"
	// ConnectorTable is the table that holds the connector relation/edge.
	ConnectorTable = "cluster_costs"
	// ConnectorInverseTable is the table name for the Connector entity.
	// It exists in this package in order to avoid circular dependency with the "connector" package.
	ConnectorInverseTable = "connectors"
	// ConnectorColumn is the table column denoting the connector relation/edge.
	ConnectorColumn = "connector_id"
)

// Columns holds all SQL columns for clustercost fields.
var Columns = []string{
	FieldID,
	FieldStartTime,
	FieldEndTime,
	FieldMinutes,
	FieldConnectorID,
	FieldClusterName,
	FieldTotalCost,
	FieldCurrency,
	FieldAllocationCost,
	FieldIdleCost,
	FieldManagementCost,
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
	// ConnectorIDValidator is a validator for the "connectorID" field. It is called by the builders before save.
	ConnectorIDValidator func(string) error
	// ClusterNameValidator is a validator for the "clusterName" field. It is called by the builders before save.
	ClusterNameValidator func(string) error
	// DefaultTotalCost holds the default value on creation for the "totalCost" field.
	DefaultTotalCost float64
	// TotalCostValidator is a validator for the "totalCost" field. It is called by the builders before save.
	TotalCostValidator func(float64) error
	// DefaultAllocationCost holds the default value on creation for the "allocationCost" field.
	DefaultAllocationCost float64
	// AllocationCostValidator is a validator for the "allocationCost" field. It is called by the builders before save.
	AllocationCostValidator func(float64) error
	// DefaultIdleCost holds the default value on creation for the "idleCost" field.
	DefaultIdleCost float64
	// IdleCostValidator is a validator for the "idleCost" field. It is called by the builders before save.
	IdleCostValidator func(float64) error
	// DefaultManagementCost holds the default value on creation for the "managementCost" field.
	DefaultManagementCost float64
	// ManagementCostValidator is a validator for the "managementCost" field. It is called by the builders before save.
	ManagementCostValidator func(float64) error
)

// OrderOption defines the ordering options for the ClusterCost queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByStartTime orders the results by the startTime field.
func ByStartTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldStartTime, opts...).ToFunc()
}

// ByEndTime orders the results by the endTime field.
func ByEndTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldEndTime, opts...).ToFunc()
}

// ByMinutes orders the results by the minutes field.
func ByMinutes(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldMinutes, opts...).ToFunc()
}

// ByConnectorID orders the results by the connectorID field.
func ByConnectorID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldConnectorID, opts...).ToFunc()
}

// ByClusterName orders the results by the clusterName field.
func ByClusterName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldClusterName, opts...).ToFunc()
}

// ByTotalCost orders the results by the totalCost field.
func ByTotalCost(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldTotalCost, opts...).ToFunc()
}

// ByCurrency orders the results by the currency field.
func ByCurrency(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCurrency, opts...).ToFunc()
}

// ByAllocationCost orders the results by the allocationCost field.
func ByAllocationCost(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldAllocationCost, opts...).ToFunc()
}

// ByIdleCost orders the results by the idleCost field.
func ByIdleCost(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldIdleCost, opts...).ToFunc()
}

// ByManagementCost orders the results by the managementCost field.
func ByManagementCost(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldManagementCost, opts...).ToFunc()
}

// ByConnectorField orders the results by connector field.
func ByConnectorField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newConnectorStep(), sql.OrderByField(field, opts...))
	}
}
func newConnectorStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ConnectorInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, ConnectorTable, ConnectorColumn),
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
