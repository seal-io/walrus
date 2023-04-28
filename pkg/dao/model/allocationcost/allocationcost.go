// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package allocationcost

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"

	"github.com/seal-io/seal/pkg/dao/types"
)

const (
	// Label holds the string label denoting the allocationcost type in the database.
	Label = "allocation_cost"
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
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldFingerprint holds the string denoting the fingerprint field in the database.
	FieldFingerprint = "fingerprint"
	// FieldClusterName holds the string denoting the clustername field in the database.
	FieldClusterName = "cluster_name"
	// FieldNamespace holds the string denoting the namespace field in the database.
	FieldNamespace = "namespace"
	// FieldNode holds the string denoting the node field in the database.
	FieldNode = "node"
	// FieldController holds the string denoting the controller field in the database.
	FieldController = "controller"
	// FieldControllerKind holds the string denoting the controllerkind field in the database.
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
	// FieldCpuCost holds the string denoting the cpucost field in the database.
	FieldCpuCost = "cpu_cost"
	// FieldCpuCoreRequest holds the string denoting the cpucorerequest field in the database.
	FieldCpuCoreRequest = "cpu_core_request"
	// FieldGpuCost holds the string denoting the gpucost field in the database.
	FieldGpuCost = "gpu_cost"
	// FieldGpuCount holds the string denoting the gpucount field in the database.
	FieldGpuCount = "gpu_count"
	// FieldRamCost holds the string denoting the ramcost field in the database.
	FieldRamCost = "ram_cost"
	// FieldRamByteRequest holds the string denoting the rambyterequest field in the database.
	FieldRamByteRequest = "ram_byte_request"
	// FieldPvCost holds the string denoting the pvcost field in the database.
	FieldPvCost = "pv_cost"
	// FieldPvBytes holds the string denoting the pvbytes field in the database.
	FieldPvBytes = "pv_bytes"
	// FieldLoadBalancerCost holds the string denoting the loadbalancercost field in the database.
	FieldLoadBalancerCost = "load_balancer_cost"
	// FieldCpuCoreUsageAverage holds the string denoting the cpucoreusageaverage field in the database.
	FieldCpuCoreUsageAverage = "cpu_core_usage_average"
	// FieldCpuCoreUsageMax holds the string denoting the cpucoreusagemax field in the database.
	FieldCpuCoreUsageMax = "cpu_core_usage_max"
	// FieldRamByteUsageAverage holds the string denoting the rambyteusageaverage field in the database.
	FieldRamByteUsageAverage = "ram_byte_usage_average"
	// FieldRamByteUsageMax holds the string denoting the rambyteusagemax field in the database.
	FieldRamByteUsageMax = "ram_byte_usage_max"
	// EdgeConnector holds the string denoting the connector edge name in mutations.
	EdgeConnector = "connector"
	// Table holds the table name of the allocationcost in the database.
	Table = "allocation_costs"
	// ConnectorTable is the table that holds the connector relation/edge.
	ConnectorTable = "allocation_costs"
	// ConnectorInverseTable is the table name for the Connector entity.
	// It exists in this package in order to avoid circular dependency with the "connector" package.
	ConnectorInverseTable = "connectors"
	// ConnectorColumn is the table column denoting the connector relation/edge.
	ConnectorColumn = "connector_id"
)

// Columns holds all SQL columns for allocationcost fields.
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
	FieldCpuCost,
	FieldCpuCoreRequest,
	FieldGpuCost,
	FieldGpuCount,
	FieldRamCost,
	FieldRamByteRequest,
	FieldPvCost,
	FieldPvBytes,
	FieldLoadBalancerCost,
	FieldCpuCoreUsageAverage,
	FieldCpuCoreUsageMax,
	FieldRamByteUsageAverage,
	FieldRamByteUsageMax,
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
	// DefaultPvs holds the default value on creation for the "pvs" field.
	DefaultPvs map[string]types.PVCost
	// DefaultLabels holds the default value on creation for the "labels" field.
	DefaultLabels map[string]string
	// DefaultTotalCost holds the default value on creation for the "totalCost" field.
	DefaultTotalCost float64
	// TotalCostValidator is a validator for the "totalCost" field. It is called by the builders before save.
	TotalCostValidator func(float64) error
	// DefaultCpuCost holds the default value on creation for the "cpuCost" field.
	DefaultCpuCost float64
	// CpuCostValidator is a validator for the "cpuCost" field. It is called by the builders before save.
	CpuCostValidator func(float64) error
	// DefaultCpuCoreRequest holds the default value on creation for the "cpuCoreRequest" field.
	DefaultCpuCoreRequest float64
	// CpuCoreRequestValidator is a validator for the "cpuCoreRequest" field. It is called by the builders before save.
	CpuCoreRequestValidator func(float64) error
	// DefaultGpuCost holds the default value on creation for the "gpuCost" field.
	DefaultGpuCost float64
	// GpuCostValidator is a validator for the "gpuCost" field. It is called by the builders before save.
	GpuCostValidator func(float64) error
	// DefaultGpuCount holds the default value on creation for the "gpuCount" field.
	DefaultGpuCount float64
	// GpuCountValidator is a validator for the "gpuCount" field. It is called by the builders before save.
	GpuCountValidator func(float64) error
	// DefaultRamCost holds the default value on creation for the "ramCost" field.
	DefaultRamCost float64
	// RamCostValidator is a validator for the "ramCost" field. It is called by the builders before save.
	RamCostValidator func(float64) error
	// DefaultRamByteRequest holds the default value on creation for the "ramByteRequest" field.
	DefaultRamByteRequest float64
	// RamByteRequestValidator is a validator for the "ramByteRequest" field. It is called by the builders before save.
	RamByteRequestValidator func(float64) error
	// DefaultPvCost holds the default value on creation for the "pvCost" field.
	DefaultPvCost float64
	// PvCostValidator is a validator for the "pvCost" field. It is called by the builders before save.
	PvCostValidator func(float64) error
	// DefaultPvBytes holds the default value on creation for the "pvBytes" field.
	DefaultPvBytes float64
	// PvBytesValidator is a validator for the "pvBytes" field. It is called by the builders before save.
	PvBytesValidator func(float64) error
	// DefaultLoadBalancerCost holds the default value on creation for the "loadBalancerCost" field.
	DefaultLoadBalancerCost float64
	// LoadBalancerCostValidator is a validator for the "loadBalancerCost" field. It is called by the builders before save.
	LoadBalancerCostValidator func(float64) error
	// DefaultCpuCoreUsageAverage holds the default value on creation for the "cpuCoreUsageAverage" field.
	DefaultCpuCoreUsageAverage float64
	// CpuCoreUsageAverageValidator is a validator for the "cpuCoreUsageAverage" field. It is called by the builders before save.
	CpuCoreUsageAverageValidator func(float64) error
	// DefaultCpuCoreUsageMax holds the default value on creation for the "cpuCoreUsageMax" field.
	DefaultCpuCoreUsageMax float64
	// CpuCoreUsageMaxValidator is a validator for the "cpuCoreUsageMax" field. It is called by the builders before save.
	CpuCoreUsageMaxValidator func(float64) error
	// DefaultRamByteUsageAverage holds the default value on creation for the "ramByteUsageAverage" field.
	DefaultRamByteUsageAverage float64
	// RamByteUsageAverageValidator is a validator for the "ramByteUsageAverage" field. It is called by the builders before save.
	RamByteUsageAverageValidator func(float64) error
	// DefaultRamByteUsageMax holds the default value on creation for the "ramByteUsageMax" field.
	DefaultRamByteUsageMax float64
	// RamByteUsageMaxValidator is a validator for the "ramByteUsageMax" field. It is called by the builders before save.
	RamByteUsageMaxValidator func(float64) error
)

// OrderOption defines the ordering options for the AllocationCost queries.
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

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByFingerprint orders the results by the fingerprint field.
func ByFingerprint(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldFingerprint, opts...).ToFunc()
}

// ByClusterName orders the results by the clusterName field.
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

// ByControllerKind orders the results by the controllerKind field.
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

// ByCpuCost orders the results by the cpuCost field.
func ByCpuCost(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCpuCost, opts...).ToFunc()
}

// ByCpuCoreRequest orders the results by the cpuCoreRequest field.
func ByCpuCoreRequest(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCpuCoreRequest, opts...).ToFunc()
}

// ByGpuCost orders the results by the gpuCost field.
func ByGpuCost(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldGpuCost, opts...).ToFunc()
}

// ByGpuCount orders the results by the gpuCount field.
func ByGpuCount(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldGpuCount, opts...).ToFunc()
}

// ByRamCost orders the results by the ramCost field.
func ByRamCost(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRamCost, opts...).ToFunc()
}

// ByRamByteRequest orders the results by the ramByteRequest field.
func ByRamByteRequest(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRamByteRequest, opts...).ToFunc()
}

// ByPvCost orders the results by the pvCost field.
func ByPvCost(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPvCost, opts...).ToFunc()
}

// ByPvBytes orders the results by the pvBytes field.
func ByPvBytes(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPvBytes, opts...).ToFunc()
}

// ByLoadBalancerCost orders the results by the loadBalancerCost field.
func ByLoadBalancerCost(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldLoadBalancerCost, opts...).ToFunc()
}

// ByCpuCoreUsageAverage orders the results by the cpuCoreUsageAverage field.
func ByCpuCoreUsageAverage(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCpuCoreUsageAverage, opts...).ToFunc()
}

// ByCpuCoreUsageMax orders the results by the cpuCoreUsageMax field.
func ByCpuCoreUsageMax(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCpuCoreUsageMax, opts...).ToFunc()
}

// ByRamByteUsageAverage orders the results by the ramByteUsageAverage field.
func ByRamByteUsageAverage(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRamByteUsageAverage, opts...).ToFunc()
}

// ByRamByteUsageMax orders the results by the ramByteUsageMax field.
func ByRamByteUsageMax(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRamByteUsageMax, opts...).ToFunc()
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
