// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"

	"github.com/seal-io/seal/pkg/dao/model/allocationcost"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/types"
)

// AllocationCost is the model entity for the AllocationCost schema.
type AllocationCost struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Usage start time for current cost
	StartTime time.Time `json:"startTime,omitempty"`
	// Usage end time for current cost
	EndTime time.Time `json:"endTime,omitempty"`
	// Usage minutes from start time to end time
	Minutes float64 `json:"minutes,omitempty"`
	// ID of the connector
	ConnectorID types.ID `json:"connectorID,omitempty"`
	// Resource name for current cost, could be __unmounted__
	Name string `json:"name,omitempty"`
	// String generated from resource properties, used to identify this cost
	Fingerprint string `json:"fingerprint,omitempty"`
	// Cluster name for current cost
	ClusterName string `json:"clusterName,omitempty"`
	// Namespace for current cost
	Namespace string `json:"namespace,omitempty"`
	// Node for current cost
	Node string `json:"node,omitempty"`
	// Controller name for the cost linked resource
	Controller string `json:"controller,omitempty"`
	// Controller kind for the cost linked resource, deployment, statefulSet etc.
	ControllerKind string `json:"controllerKind,omitempty"`
	// Pod name for current cost
	Pod string `json:"pod,omitempty"`
	// Container name for current cost
	Container string `json:"container,omitempty"`
	// PV list for current cost linked
	Pvs map[string]types.PVCost `json:"pvs,omitempty"`
	// Labels for the cost linked resource
	Labels map[string]string `json:"labels,omitempty"`
	// Cost number
	TotalCost float64 `json:"totalCost,omitempty"`
	// Cost currency
	Currency int `json:"currency,omitempty"`
	// Cpu cost for current cost
	CpuCost float64 `json:"cpuCost,omitempty"`
	// Cpu core requested
	CpuCoreRequest float64 `json:"cpuCoreRequest,omitempty"`
	// GPU cost for current cost
	GpuCost float64 `json:"gpuCost,omitempty"`
	// GPU core count
	GpuCount float64 `json:"gpuCount,omitempty"`
	// Ram cost for current cost
	RamCost float64 `json:"ramCost,omitempty"`
	// Ram requested in byte
	RamByteRequest float64 `json:"ramByteRequest,omitempty"`
	// PV cost for current cost linked
	PvCost float64 `json:"pvCost,omitempty"`
	// PV bytes for current cost linked
	PvBytes float64 `json:"pvBytes,omitempty"`
	// CPU core average usage
	CpuCoreUsageAverage float64 `json:"cpuCoreUsageAverage,omitempty"`
	// CPU core max usage
	CpuCoreUsageMax float64 `json:"cpuCoreUsageMax,omitempty"`
	// Ram average usage in byte
	RamByteUsageAverage float64 `json:"ramByteUsageAverage,omitempty"`
	// Ram max usage in byte
	RamByteUsageMax float64 `json:"ramByteUsageMax,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the AllocationCostQuery when eager-loading is set.
	Edges AllocationCostEdges `json:"edges"`
}

// AllocationCostEdges holds the relations/edges for other nodes in the graph.
type AllocationCostEdges struct {
	// Connector current cost linked
	Connector *Connector `json:"connector,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// ConnectorOrErr returns the Connector value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e AllocationCostEdges) ConnectorOrErr() (*Connector, error) {
	if e.loadedTypes[0] {
		if e.Connector == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: connector.Label}
		}
		return e.Connector, nil
	}
	return nil, &NotLoadedError{edge: "connector"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*AllocationCost) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case allocationcost.FieldPvs, allocationcost.FieldLabels:
			values[i] = new([]byte)
		case allocationcost.FieldMinutes, allocationcost.FieldTotalCost, allocationcost.FieldCpuCost, allocationcost.FieldCpuCoreRequest, allocationcost.FieldGpuCost, allocationcost.FieldGpuCount, allocationcost.FieldRamCost, allocationcost.FieldRamByteRequest, allocationcost.FieldPvCost, allocationcost.FieldPvBytes, allocationcost.FieldCpuCoreUsageAverage, allocationcost.FieldCpuCoreUsageMax, allocationcost.FieldRamByteUsageAverage, allocationcost.FieldRamByteUsageMax:
			values[i] = new(sql.NullFloat64)
		case allocationcost.FieldID, allocationcost.FieldCurrency:
			values[i] = new(sql.NullInt64)
		case allocationcost.FieldName, allocationcost.FieldFingerprint, allocationcost.FieldClusterName, allocationcost.FieldNamespace, allocationcost.FieldNode, allocationcost.FieldController, allocationcost.FieldControllerKind, allocationcost.FieldPod, allocationcost.FieldContainer:
			values[i] = new(sql.NullString)
		case allocationcost.FieldStartTime, allocationcost.FieldEndTime:
			values[i] = new(sql.NullTime)
		case allocationcost.FieldConnectorID:
			values[i] = new(types.ID)
		default:
			return nil, fmt.Errorf("unexpected column %q for type AllocationCost", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the AllocationCost fields.
func (ac *AllocationCost) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case allocationcost.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			ac.ID = int(value.Int64)
		case allocationcost.FieldStartTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field startTime", values[i])
			} else if value.Valid {
				ac.StartTime = value.Time
			}
		case allocationcost.FieldEndTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field endTime", values[i])
			} else if value.Valid {
				ac.EndTime = value.Time
			}
		case allocationcost.FieldMinutes:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field minutes", values[i])
			} else if value.Valid {
				ac.Minutes = value.Float64
			}
		case allocationcost.FieldConnectorID:
			if value, ok := values[i].(*types.ID); !ok {
				return fmt.Errorf("unexpected type %T for field connectorID", values[i])
			} else if value != nil {
				ac.ConnectorID = *value
			}
		case allocationcost.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				ac.Name = value.String
			}
		case allocationcost.FieldFingerprint:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field fingerprint", values[i])
			} else if value.Valid {
				ac.Fingerprint = value.String
			}
		case allocationcost.FieldClusterName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field clusterName", values[i])
			} else if value.Valid {
				ac.ClusterName = value.String
			}
		case allocationcost.FieldNamespace:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field namespace", values[i])
			} else if value.Valid {
				ac.Namespace = value.String
			}
		case allocationcost.FieldNode:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field node", values[i])
			} else if value.Valid {
				ac.Node = value.String
			}
		case allocationcost.FieldController:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field controller", values[i])
			} else if value.Valid {
				ac.Controller = value.String
			}
		case allocationcost.FieldControllerKind:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field controllerKind", values[i])
			} else if value.Valid {
				ac.ControllerKind = value.String
			}
		case allocationcost.FieldPod:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field pod", values[i])
			} else if value.Valid {
				ac.Pod = value.String
			}
		case allocationcost.FieldContainer:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field container", values[i])
			} else if value.Valid {
				ac.Container = value.String
			}
		case allocationcost.FieldPvs:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field pvs", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &ac.Pvs); err != nil {
					return fmt.Errorf("unmarshal field pvs: %w", err)
				}
			}
		case allocationcost.FieldLabels:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field labels", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &ac.Labels); err != nil {
					return fmt.Errorf("unmarshal field labels: %w", err)
				}
			}
		case allocationcost.FieldTotalCost:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field totalCost", values[i])
			} else if value.Valid {
				ac.TotalCost = value.Float64
			}
		case allocationcost.FieldCurrency:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field currency", values[i])
			} else if value.Valid {
				ac.Currency = int(value.Int64)
			}
		case allocationcost.FieldCpuCost:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field cpuCost", values[i])
			} else if value.Valid {
				ac.CpuCost = value.Float64
			}
		case allocationcost.FieldCpuCoreRequest:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field cpuCoreRequest", values[i])
			} else if value.Valid {
				ac.CpuCoreRequest = value.Float64
			}
		case allocationcost.FieldGpuCost:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field gpuCost", values[i])
			} else if value.Valid {
				ac.GpuCost = value.Float64
			}
		case allocationcost.FieldGpuCount:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field gpuCount", values[i])
			} else if value.Valid {
				ac.GpuCount = value.Float64
			}
		case allocationcost.FieldRamCost:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field ramCost", values[i])
			} else if value.Valid {
				ac.RamCost = value.Float64
			}
		case allocationcost.FieldRamByteRequest:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field ramByteRequest", values[i])
			} else if value.Valid {
				ac.RamByteRequest = value.Float64
			}
		case allocationcost.FieldPvCost:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field pvCost", values[i])
			} else if value.Valid {
				ac.PvCost = value.Float64
			}
		case allocationcost.FieldPvBytes:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field pvBytes", values[i])
			} else if value.Valid {
				ac.PvBytes = value.Float64
			}
		case allocationcost.FieldCpuCoreUsageAverage:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field cpuCoreUsageAverage", values[i])
			} else if value.Valid {
				ac.CpuCoreUsageAverage = value.Float64
			}
		case allocationcost.FieldCpuCoreUsageMax:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field cpuCoreUsageMax", values[i])
			} else if value.Valid {
				ac.CpuCoreUsageMax = value.Float64
			}
		case allocationcost.FieldRamByteUsageAverage:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field ramByteUsageAverage", values[i])
			} else if value.Valid {
				ac.RamByteUsageAverage = value.Float64
			}
		case allocationcost.FieldRamByteUsageMax:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field ramByteUsageMax", values[i])
			} else if value.Valid {
				ac.RamByteUsageMax = value.Float64
			}
		}
	}
	return nil
}

// QueryConnector queries the "connector" edge of the AllocationCost entity.
func (ac *AllocationCost) QueryConnector() *ConnectorQuery {
	return NewAllocationCostClient(ac.config).QueryConnector(ac)
}

// Update returns a builder for updating this AllocationCost.
// Note that you need to call AllocationCost.Unwrap() before calling this method if this AllocationCost
// was returned from a transaction, and the transaction was committed or rolled back.
func (ac *AllocationCost) Update() *AllocationCostUpdateOne {
	return NewAllocationCostClient(ac.config).UpdateOne(ac)
}

// Unwrap unwraps the AllocationCost entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (ac *AllocationCost) Unwrap() *AllocationCost {
	_tx, ok := ac.config.driver.(*txDriver)
	if !ok {
		panic("model: AllocationCost is not a transactional entity")
	}
	ac.config.driver = _tx.drv
	return ac
}

// String implements the fmt.Stringer.
func (ac *AllocationCost) String() string {
	var builder strings.Builder
	builder.WriteString("AllocationCost(")
	builder.WriteString(fmt.Sprintf("id=%v, ", ac.ID))
	builder.WriteString("startTime=")
	builder.WriteString(ac.StartTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("endTime=")
	builder.WriteString(ac.EndTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("minutes=")
	builder.WriteString(fmt.Sprintf("%v", ac.Minutes))
	builder.WriteString(", ")
	builder.WriteString("connectorID=")
	builder.WriteString(fmt.Sprintf("%v", ac.ConnectorID))
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(ac.Name)
	builder.WriteString(", ")
	builder.WriteString("fingerprint=")
	builder.WriteString(ac.Fingerprint)
	builder.WriteString(", ")
	builder.WriteString("clusterName=")
	builder.WriteString(ac.ClusterName)
	builder.WriteString(", ")
	builder.WriteString("namespace=")
	builder.WriteString(ac.Namespace)
	builder.WriteString(", ")
	builder.WriteString("node=")
	builder.WriteString(ac.Node)
	builder.WriteString(", ")
	builder.WriteString("controller=")
	builder.WriteString(ac.Controller)
	builder.WriteString(", ")
	builder.WriteString("controllerKind=")
	builder.WriteString(ac.ControllerKind)
	builder.WriteString(", ")
	builder.WriteString("pod=")
	builder.WriteString(ac.Pod)
	builder.WriteString(", ")
	builder.WriteString("container=")
	builder.WriteString(ac.Container)
	builder.WriteString(", ")
	builder.WriteString("pvs=")
	builder.WriteString(fmt.Sprintf("%v", ac.Pvs))
	builder.WriteString(", ")
	builder.WriteString("labels=")
	builder.WriteString(fmt.Sprintf("%v", ac.Labels))
	builder.WriteString(", ")
	builder.WriteString("totalCost=")
	builder.WriteString(fmt.Sprintf("%v", ac.TotalCost))
	builder.WriteString(", ")
	builder.WriteString("currency=")
	builder.WriteString(fmt.Sprintf("%v", ac.Currency))
	builder.WriteString(", ")
	builder.WriteString("cpuCost=")
	builder.WriteString(fmt.Sprintf("%v", ac.CpuCost))
	builder.WriteString(", ")
	builder.WriteString("cpuCoreRequest=")
	builder.WriteString(fmt.Sprintf("%v", ac.CpuCoreRequest))
	builder.WriteString(", ")
	builder.WriteString("gpuCost=")
	builder.WriteString(fmt.Sprintf("%v", ac.GpuCost))
	builder.WriteString(", ")
	builder.WriteString("gpuCount=")
	builder.WriteString(fmt.Sprintf("%v", ac.GpuCount))
	builder.WriteString(", ")
	builder.WriteString("ramCost=")
	builder.WriteString(fmt.Sprintf("%v", ac.RamCost))
	builder.WriteString(", ")
	builder.WriteString("ramByteRequest=")
	builder.WriteString(fmt.Sprintf("%v", ac.RamByteRequest))
	builder.WriteString(", ")
	builder.WriteString("pvCost=")
	builder.WriteString(fmt.Sprintf("%v", ac.PvCost))
	builder.WriteString(", ")
	builder.WriteString("pvBytes=")
	builder.WriteString(fmt.Sprintf("%v", ac.PvBytes))
	builder.WriteString(", ")
	builder.WriteString("cpuCoreUsageAverage=")
	builder.WriteString(fmt.Sprintf("%v", ac.CpuCoreUsageAverage))
	builder.WriteString(", ")
	builder.WriteString("cpuCoreUsageMax=")
	builder.WriteString(fmt.Sprintf("%v", ac.CpuCoreUsageMax))
	builder.WriteString(", ")
	builder.WriteString("ramByteUsageAverage=")
	builder.WriteString(fmt.Sprintf("%v", ac.RamByteUsageAverage))
	builder.WriteString(", ")
	builder.WriteString("ramByteUsageMax=")
	builder.WriteString(fmt.Sprintf("%v", ac.RamByteUsageMax))
	builder.WriteByte(')')
	return builder.String()
}

// AllocationCosts is a parsable slice of AllocationCost.
type AllocationCosts []*AllocationCost
