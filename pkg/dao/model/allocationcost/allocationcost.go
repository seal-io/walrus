// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package allocationcost

import (
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
