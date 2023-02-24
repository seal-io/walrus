// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package model

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"

	"github.com/seal-io/seal/pkg/dao/model/clustercost"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/types"
)

// ClusterCost is the model entity for the ClusterCost schema.
type ClusterCost struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Usage start time for current cost
	StartTime time.Time `json:"startTime"`
	// Usage end time for current cost
	EndTime time.Time `json:"endTime"`
	// Usage minutes from start time to end time
	Minutes float64 `json:"minutes"`
	// ID of the connector
	ConnectorID types.ID `json:"connectorID"`
	// Cluster name for current cost
	ClusterName string `json:"clusterName"`
	// Cost number
	TotalCost float64 `json:"totalCost,omitempty"`
	// Cost currency
	Currency int `json:"currency,omitempty"`
	// CPU cost for current cost
	CpuCost float64 `json:"cpuCost,omitempty"`
	// GPU cost for current cost
	GpuCost float64 `json:"gpuCost,omitempty"`
	// Ram cost for current cost
	RamCost float64 `json:"ramCost,omitempty"`
	// Storage cost for current cost
	StorageCost float64 `json:"storageCost,omitempty"`
	// Allocation cost for current cost
	AllocationCost float64 `json:"allocationCost,omitempty"`
	// Idle cost for current cost
	IdleCost float64 `json:"idleCost,omitempty"`
	// Storage cost for current cost
	ManagementCost float64 `json:"managementCost,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ClusterCostQuery when eager-loading is set.
	Edges ClusterCostEdges `json:"edges"`
}

// ClusterCostEdges holds the relations/edges for other nodes in the graph.
type ClusterCostEdges struct {
	// Connector current cost linked
	Connector *Connector `json:"connector,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// ConnectorOrErr returns the Connector value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ClusterCostEdges) ConnectorOrErr() (*Connector, error) {
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
func (*ClusterCost) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case clustercost.FieldMinutes, clustercost.FieldTotalCost, clustercost.FieldCpuCost, clustercost.FieldGpuCost, clustercost.FieldRamCost, clustercost.FieldStorageCost, clustercost.FieldAllocationCost, clustercost.FieldIdleCost, clustercost.FieldManagementCost:
			values[i] = new(sql.NullFloat64)
		case clustercost.FieldID, clustercost.FieldCurrency:
			values[i] = new(sql.NullInt64)
		case clustercost.FieldClusterName:
			values[i] = new(sql.NullString)
		case clustercost.FieldStartTime, clustercost.FieldEndTime:
			values[i] = new(sql.NullTime)
		case clustercost.FieldConnectorID:
			values[i] = new(types.ID)
		default:
			return nil, fmt.Errorf("unexpected column %q for type ClusterCost", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the ClusterCost fields.
func (cc *ClusterCost) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case clustercost.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			cc.ID = int(value.Int64)
		case clustercost.FieldStartTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field startTime", values[i])
			} else if value.Valid {
				cc.StartTime = value.Time
			}
		case clustercost.FieldEndTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field endTime", values[i])
			} else if value.Valid {
				cc.EndTime = value.Time
			}
		case clustercost.FieldMinutes:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field minutes", values[i])
			} else if value.Valid {
				cc.Minutes = value.Float64
			}
		case clustercost.FieldConnectorID:
			if value, ok := values[i].(*types.ID); !ok {
				return fmt.Errorf("unexpected type %T for field connectorID", values[i])
			} else if value != nil {
				cc.ConnectorID = *value
			}
		case clustercost.FieldClusterName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field clusterName", values[i])
			} else if value.Valid {
				cc.ClusterName = value.String
			}
		case clustercost.FieldTotalCost:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field totalCost", values[i])
			} else if value.Valid {
				cc.TotalCost = value.Float64
			}
		case clustercost.FieldCurrency:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field currency", values[i])
			} else if value.Valid {
				cc.Currency = int(value.Int64)
			}
		case clustercost.FieldCpuCost:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field cpuCost", values[i])
			} else if value.Valid {
				cc.CpuCost = value.Float64
			}
		case clustercost.FieldGpuCost:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field gpuCost", values[i])
			} else if value.Valid {
				cc.GpuCost = value.Float64
			}
		case clustercost.FieldRamCost:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field ramCost", values[i])
			} else if value.Valid {
				cc.RamCost = value.Float64
			}
		case clustercost.FieldStorageCost:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field storageCost", values[i])
			} else if value.Valid {
				cc.StorageCost = value.Float64
			}
		case clustercost.FieldAllocationCost:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field allocationCost", values[i])
			} else if value.Valid {
				cc.AllocationCost = value.Float64
			}
		case clustercost.FieldIdleCost:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field idleCost", values[i])
			} else if value.Valid {
				cc.IdleCost = value.Float64
			}
		case clustercost.FieldManagementCost:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field managementCost", values[i])
			} else if value.Valid {
				cc.ManagementCost = value.Float64
			}
		}
	}
	return nil
}

// QueryConnector queries the "connector" edge of the ClusterCost entity.
func (cc *ClusterCost) QueryConnector() *ConnectorQuery {
	return NewClusterCostClient(cc.config).QueryConnector(cc)
}

// Update returns a builder for updating this ClusterCost.
// Note that you need to call ClusterCost.Unwrap() before calling this method if this ClusterCost
// was returned from a transaction, and the transaction was committed or rolled back.
func (cc *ClusterCost) Update() *ClusterCostUpdateOne {
	return NewClusterCostClient(cc.config).UpdateOne(cc)
}

// Unwrap unwraps the ClusterCost entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (cc *ClusterCost) Unwrap() *ClusterCost {
	_tx, ok := cc.config.driver.(*txDriver)
	if !ok {
		panic("model: ClusterCost is not a transactional entity")
	}
	cc.config.driver = _tx.drv
	return cc
}

// String implements the fmt.Stringer.
func (cc *ClusterCost) String() string {
	var builder strings.Builder
	builder.WriteString("ClusterCost(")
	builder.WriteString(fmt.Sprintf("id=%v, ", cc.ID))
	builder.WriteString("startTime=")
	builder.WriteString(cc.StartTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("endTime=")
	builder.WriteString(cc.EndTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("minutes=")
	builder.WriteString(fmt.Sprintf("%v", cc.Minutes))
	builder.WriteString(", ")
	builder.WriteString("connectorID=")
	builder.WriteString(fmt.Sprintf("%v", cc.ConnectorID))
	builder.WriteString(", ")
	builder.WriteString("clusterName=")
	builder.WriteString(cc.ClusterName)
	builder.WriteString(", ")
	builder.WriteString("totalCost=")
	builder.WriteString(fmt.Sprintf("%v", cc.TotalCost))
	builder.WriteString(", ")
	builder.WriteString("currency=")
	builder.WriteString(fmt.Sprintf("%v", cc.Currency))
	builder.WriteString(", ")
	builder.WriteString("cpuCost=")
	builder.WriteString(fmt.Sprintf("%v", cc.CpuCost))
	builder.WriteString(", ")
	builder.WriteString("gpuCost=")
	builder.WriteString(fmt.Sprintf("%v", cc.GpuCost))
	builder.WriteString(", ")
	builder.WriteString("ramCost=")
	builder.WriteString(fmt.Sprintf("%v", cc.RamCost))
	builder.WriteString(", ")
	builder.WriteString("storageCost=")
	builder.WriteString(fmt.Sprintf("%v", cc.StorageCost))
	builder.WriteString(", ")
	builder.WriteString("allocationCost=")
	builder.WriteString(fmt.Sprintf("%v", cc.AllocationCost))
	builder.WriteString(", ")
	builder.WriteString("idleCost=")
	builder.WriteString(fmt.Sprintf("%v", cc.IdleCost))
	builder.WriteString(", ")
	builder.WriteString("managementCost=")
	builder.WriteString(fmt.Sprintf("%v", cc.ManagementCost))
	builder.WriteByte(')')
	return builder.String()
}

// ClusterCosts is a parsable slice of ClusterCost.
type ClusterCosts []*ClusterCost

func (cc ClusterCosts) config(cfg config) {
	for _i := range cc {
		cc[_i].config = cfg
	}
}
