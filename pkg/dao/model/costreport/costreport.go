// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "seal". DO NOT EDIT.

package costreport

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"

	"github.com/seal-io/seal/pkg/dao/types"
)

const (
	// Label holds the string label denoting the costreport type in the database.
	Label = "cost_report"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldStartTime holds the string denoting the start_time field in the database.
	FieldStartTime = "start_time"
	// FieldEndTime holds the string denoting the end_time field in the database.
	FieldEndTime = "end_time"
	// FieldMinutes holds the string denoting the minutes field in the database.
	FieldMinutes = "minutes"
	// FieldConnectorID holds the string denoting the connector_id field in the database.
	FieldConnectorID = "connector_id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldFingerprint holds the string denoting the fingerprint field in the database.
	FieldFingerprint = "fingerprint"
	// FieldClusterName holds the string denoting the cluster_name field in the database.
	FieldClusterName = "cluster_name"
	// FieldNamespace holds the string denoting the namespace field in the database.
	FieldNamespace = "namespace"
	// FieldNode holds the string denoting the node field in the database.
	FieldNode = "node"
	// FieldController holds the string denoting the controller field in the database.
	FieldController = "controller"
	// FieldControllerKind holds the string denoting the controller_kind field in the database.
	FieldControllerKind = "controller_kind"
	// FieldPod holds the string denoting the pod field in the database.
	FieldPod = "pod"
	// FieldContainer holds the string denoting the container field in the database.
	FieldContainer = "container"
	// FieldPvs holds the string denoting the pvs field in the database.
	FieldPvs = "pvs"
	// FieldLabels holds the string denoting the labels field in the database.
	FieldLabels = "labels"
	// FieldTotalCost holds the string denoting the totalcost field in the database.
	FieldTotalCost = "total_cost"
	// FieldCurrency holds the string denoting the currency field in the database.
	FieldCurrency = "currency"
	// FieldCPUCost holds the string denoting the cpu_cost field in the database.
	FieldCPUCost = "cpu_cost"
	// FieldCPUCoreRequest holds the string denoting the cpu_core_request field in the database.
	FieldCPUCoreRequest = "cpu_core_request"
	// FieldGpuCost holds the string denoting the gpu_cost field in the database.
	FieldGpuCost = "gpu_cost"
	// FieldGpuCount holds the string denoting the gpu_count field in the database.
	FieldGpuCount = "gpu_count"
	// FieldRAMCost holds the string denoting the ram_cost field in the database.
	FieldRAMCost = "ram_cost"
	// FieldRAMByteRequest holds the string denoting the ram_byte_request field in the database.
	FieldRAMByteRequest = "ram_byte_request"
	// FieldPvCost holds the string denoting the pv_cost field in the database.
	FieldPvCost = "pv_cost"
	// FieldPvBytes holds the string denoting the pv_bytes field in the database.
	FieldPvBytes = "pv_bytes"
	// FieldLoadBalancerCost holds the string denoting the load_balancer_cost field in the database.
	FieldLoadBalancerCost = "load_balancer_cost"
	// FieldCPUCoreUsageAverage holds the string denoting the cpu_core_usage_average field in the database.
	FieldCPUCoreUsageAverage = "cpu_core_usage_average"
	// FieldCPUCoreUsageMax holds the string denoting the cpu_core_usage_max field in the database.
	FieldCPUCoreUsageMax = "cpu_core_usage_max"
	// FieldRAMByteUsageAverage holds the string denoting the ram_byte_usage_average field in the database.
	FieldRAMByteUsageAverage = "ram_byte_usage_average"
	// FieldRAMByteUsageMax holds the string denoting the ram_byte_usage_max field in the database.
	FieldRAMByteUsageMax = "ram_byte_usage_max"
	// EdgeConnector holds the string denoting the connector edge name in mutations.
	EdgeConnector = "connector"
	// Table holds the table name of the costreport in the database.
	Table = "cost_reports"
	// ConnectorTable is the table that holds the connector relation/edge.
	ConnectorTable = "cost_reports"
	// ConnectorInverseTable is the table name for the Connector entity.
	// It exists in this package in order to avoid circular dependency with the "connector" package.
	ConnectorInverseTable = "connectors"
	// ConnectorColumn is the table column denoting the connector relation/edge.
	ConnectorColumn = "connector_id"
)

// Columns holds all SQL columns for costreport fields.
var Columns = []string{
	FieldID,
	FieldStartTime,
	FieldEndTime,
	FieldMinutes,
	FieldConnectorID,
	FieldName,
	FieldFingerprint,
	FieldClusterName,
	FieldNamespace,
	FieldNode,
	FieldController,
	FieldControllerKind,
	FieldPod,
	FieldContainer,
	FieldPvs,
	FieldLabels,
	FieldTotalCost,
	FieldCurrency,
	FieldCPUCost,
	FieldCPUCoreRequest,
	FieldGpuCost,
	FieldGpuCount,
	FieldRAMCost,
	FieldRAMByteRequest,
	FieldPvCost,
	FieldPvBytes,
	FieldLoadBalancerCost,
	FieldCPUCoreUsageAverage,
	FieldCPUCoreUsageMax,
	FieldRAMByteUsageAverage,
	FieldRAMByteUsageMax,
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
	// ConnectorIDValidator is a validator for the "connector_id" field. It is called by the builders before save.
	ConnectorIDValidator func(string) error
	// DefaultPvs holds the default value on creation for the "pvs" field.
	DefaultPvs map[string]types.PVCost
	// DefaultLabels holds the default value on creation for the "labels" field.
	DefaultLabels map[string]string
	// DefaultTotalCost holds the default value on creation for the "totalCost" field.
	DefaultTotalCost float64
	// TotalCostValidator is a validator for the "totalCost" field. It is called by the builders before save.
	TotalCostValidator func(float64) error
	// DefaultCPUCost holds the default value on creation for the "cpu_cost" field.
	DefaultCPUCost float64
	// CPUCostValidator is a validator for the "cpu_cost" field. It is called by the builders before save.
	CPUCostValidator func(float64) error
	// DefaultCPUCoreRequest holds the default value on creation for the "cpu_core_request" field.
	DefaultCPUCoreRequest float64
	// CPUCoreRequestValidator is a validator for the "cpu_core_request" field. It is called by the builders before save.
	CPUCoreRequestValidator func(float64) error
	// DefaultGpuCost holds the default value on creation for the "gpu_cost" field.
	DefaultGpuCost float64
	// GpuCostValidator is a validator for the "gpu_cost" field. It is called by the builders before save.
	GpuCostValidator func(float64) error
	// DefaultGpuCount holds the default value on creation for the "gpu_count" field.
	DefaultGpuCount float64
	// GpuCountValidator is a validator for the "gpu_count" field. It is called by the builders before save.
	GpuCountValidator func(float64) error
	// DefaultRAMCost holds the default value on creation for the "ram_cost" field.
	DefaultRAMCost float64
	// RAMCostValidator is a validator for the "ram_cost" field. It is called by the builders before save.
	RAMCostValidator func(float64) error
	// DefaultRAMByteRequest holds the default value on creation for the "ram_byte_request" field.
	DefaultRAMByteRequest float64
	// RAMByteRequestValidator is a validator for the "ram_byte_request" field. It is called by the builders before save.
	RAMByteRequestValidator func(float64) error
	// DefaultPvCost holds the default value on creation for the "pv_cost" field.
	DefaultPvCost float64
	// PvCostValidator is a validator for the "pv_cost" field. It is called by the builders before save.
	PvCostValidator func(float64) error
	// DefaultPvBytes holds the default value on creation for the "pv_bytes" field.
	DefaultPvBytes float64
	// PvBytesValidator is a validator for the "pv_bytes" field. It is called by the builders before save.
	PvBytesValidator func(float64) error
	// DefaultLoadBalancerCost holds the default value on creation for the "load_balancer_cost" field.
	DefaultLoadBalancerCost float64
	// LoadBalancerCostValidator is a validator for the "load_balancer_cost" field. It is called by the builders before save.
	LoadBalancerCostValidator func(float64) error
	// DefaultCPUCoreUsageAverage holds the default value on creation for the "cpu_core_usage_average" field.
	DefaultCPUCoreUsageAverage float64
	// CPUCoreUsageAverageValidator is a validator for the "cpu_core_usage_average" field. It is called by the builders before save.
	CPUCoreUsageAverageValidator func(float64) error
	// DefaultCPUCoreUsageMax holds the default value on creation for the "cpu_core_usage_max" field.
	DefaultCPUCoreUsageMax float64
	// CPUCoreUsageMaxValidator is a validator for the "cpu_core_usage_max" field. It is called by the builders before save.
	CPUCoreUsageMaxValidator func(float64) error
	// DefaultRAMByteUsageAverage holds the default value on creation for the "ram_byte_usage_average" field.
	DefaultRAMByteUsageAverage float64
	// RAMByteUsageAverageValidator is a validator for the "ram_byte_usage_average" field. It is called by the builders before save.
	RAMByteUsageAverageValidator func(float64) error
	// DefaultRAMByteUsageMax holds the default value on creation for the "ram_byte_usage_max" field.
	DefaultRAMByteUsageMax float64
	// RAMByteUsageMaxValidator is a validator for the "ram_byte_usage_max" field. It is called by the builders before save.
	RAMByteUsageMaxValidator func(float64) error
)

// OrderOption defines the ordering options for the CostReport queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByStartTime orders the results by the start_time field.
func ByStartTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldStartTime, opts...).ToFunc()
}

// ByEndTime orders the results by the end_time field.
func ByEndTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldEndTime, opts...).ToFunc()
}

// ByMinutes orders the results by the minutes field.
func ByMinutes(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldMinutes, opts...).ToFunc()
}

// ByConnectorID orders the results by the connector_id field.
func ByConnectorID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldConnectorID, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByFingerprint orders the results by the fingerprint field.
func ByFingerprint(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldFingerprint, opts...).ToFunc()
}

// ByClusterName orders the results by the cluster_name field.
func ByClusterName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldClusterName, opts...).ToFunc()
}

// ByNamespace orders the results by the namespace field.
func ByNamespace(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldNamespace, opts...).ToFunc()
}

// ByNode orders the results by the node field.
func ByNode(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldNode, opts...).ToFunc()
}

// ByController orders the results by the controller field.
func ByController(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldController, opts...).ToFunc()
}

// ByControllerKind orders the results by the controller_kind field.
func ByControllerKind(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldControllerKind, opts...).ToFunc()
}

// ByPod orders the results by the pod field.
func ByPod(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPod, opts...).ToFunc()
}

// ByContainer orders the results by the container field.
func ByContainer(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldContainer, opts...).ToFunc()
}

// ByTotalCost orders the results by the totalCost field.
func ByTotalCost(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldTotalCost, opts...).ToFunc()
}

// ByCurrency orders the results by the currency field.
func ByCurrency(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCurrency, opts...).ToFunc()
}

// ByCPUCost orders the results by the cpu_cost field.
func ByCPUCost(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCPUCost, opts...).ToFunc()
}

// ByCPUCoreRequest orders the results by the cpu_core_request field.
func ByCPUCoreRequest(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCPUCoreRequest, opts...).ToFunc()
}

// ByGpuCost orders the results by the gpu_cost field.
func ByGpuCost(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldGpuCost, opts...).ToFunc()
}

// ByGpuCount orders the results by the gpu_count field.
func ByGpuCount(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldGpuCount, opts...).ToFunc()
}

// ByRAMCost orders the results by the ram_cost field.
func ByRAMCost(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRAMCost, opts...).ToFunc()
}

// ByRAMByteRequest orders the results by the ram_byte_request field.
func ByRAMByteRequest(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRAMByteRequest, opts...).ToFunc()
}

// ByPvCost orders the results by the pv_cost field.
func ByPvCost(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPvCost, opts...).ToFunc()
}

// ByPvBytes orders the results by the pv_bytes field.
func ByPvBytes(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPvBytes, opts...).ToFunc()
}

// ByLoadBalancerCost orders the results by the load_balancer_cost field.
func ByLoadBalancerCost(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldLoadBalancerCost, opts...).ToFunc()
}

// ByCPUCoreUsageAverage orders the results by the cpu_core_usage_average field.
func ByCPUCoreUsageAverage(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCPUCoreUsageAverage, opts...).ToFunc()
}

// ByCPUCoreUsageMax orders the results by the cpu_core_usage_max field.
func ByCPUCoreUsageMax(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCPUCoreUsageMax, opts...).ToFunc()
}

// ByRAMByteUsageAverage orders the results by the ram_byte_usage_average field.
func ByRAMByteUsageAverage(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRAMByteUsageAverage, opts...).ToFunc()
}

// ByRAMByteUsageMax orders the results by the ram_byte_usage_max field.
func ByRAMByteUsageMax(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRAMByteUsageMax, opts...).ToFunc()
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