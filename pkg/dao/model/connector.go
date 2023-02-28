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

	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/types"
)

// Connector is the model entity for the Connector schema.
type Connector struct {
	config `json:"-"`
	// ID of the ent.
	ID types.ID `json:"id,omitempty"`
	// Name of the resource.
	Name string `json:"name,omitempty"`
	// Description of the resource.
	Description string `json:"description,omitempty"`
	// Labels of the resource.
	Labels map[string]string `json:"labels,omitempty"`
	// Status of the resource.
	Status string `json:"status,omitempty"`
	// Extra message for status, like error details.
	StatusMessage string `json:"statusMessage,omitempty"`
	// Describe creation time.
	CreateTime *time.Time `json:"createTime,omitempty"`
	// Describe modification time.
	UpdateTime *time.Time `json:"updateTime,omitempty"`
	// Type of the connector.
	Type string `json:"type,omitempty"`
	// Connector config version.
	ConfigVersion string `json:"configVersion,omitempty"`
	// Connector config data.
	ConfigData map[string]interface{} `json:"configData,omitempty"`
	// Config whether enable finOps, will install prometheus and opencost while enable.
	EnableFinOps bool `json:"enableFinOps,omitempty"`
	// Status of the finOps tools.
	FinOpsStatus string `json:"finOpsStatus,omitempty"`
	// Extra message for finOps tools status, like error details.
	FinOpsStatusMessage string `json:"finOpsStatusMessage,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ConnectorQuery when eager-loading is set.
	Edges ConnectorEdges `json:"edges,omitempty"`
}

// ConnectorEdges holds the relations/edges for other nodes in the graph.
type ConnectorEdges struct {
	// Environments holds the value of the environments edge.
	Environments []*EnvironmentConnectorRelationship `json:"environments,omitempty"`
	// Resources that belong to the application.
	Resources []*ApplicationResource `json:"resources,omitempty"`
	// Cluster costs that linked to the connection
	ClusterCosts []*ClusterCost `json:"clusterCosts,omitempty"`
	// Cluster allocation resource costs that linked to the connection
	AllocationCosts []*AllocationCost `json:"allocationCosts,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [4]bool
}

// EnvironmentsOrErr returns the Environments value or an error if the edge
// was not loaded in eager-loading.
func (e ConnectorEdges) EnvironmentsOrErr() ([]*EnvironmentConnectorRelationship, error) {
	if e.loadedTypes[0] {
		return e.Environments, nil
	}
	return nil, &NotLoadedError{edge: "environments"}
}

// ResourcesOrErr returns the Resources value or an error if the edge
// was not loaded in eager-loading.
func (e ConnectorEdges) ResourcesOrErr() ([]*ApplicationResource, error) {
	if e.loadedTypes[1] {
		return e.Resources, nil
	}
	return nil, &NotLoadedError{edge: "resources"}
}

// ClusterCostsOrErr returns the ClusterCosts value or an error if the edge
// was not loaded in eager-loading.
func (e ConnectorEdges) ClusterCostsOrErr() ([]*ClusterCost, error) {
	if e.loadedTypes[2] {
		return e.ClusterCosts, nil
	}
	return nil, &NotLoadedError{edge: "clusterCosts"}
}

// AllocationCostsOrErr returns the AllocationCosts value or an error if the edge
// was not loaded in eager-loading.
func (e ConnectorEdges) AllocationCostsOrErr() ([]*AllocationCost, error) {
	if e.loadedTypes[3] {
		return e.AllocationCosts, nil
	}
	return nil, &NotLoadedError{edge: "allocationCosts"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Connector) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case connector.FieldLabels, connector.FieldConfigData:
			values[i] = new([]byte)
		case connector.FieldEnableFinOps:
			values[i] = new(sql.NullBool)
		case connector.FieldName, connector.FieldDescription, connector.FieldStatus, connector.FieldStatusMessage, connector.FieldType, connector.FieldConfigVersion, connector.FieldFinOpsStatus, connector.FieldFinOpsStatusMessage:
			values[i] = new(sql.NullString)
		case connector.FieldCreateTime, connector.FieldUpdateTime:
			values[i] = new(sql.NullTime)
		case connector.FieldID:
			values[i] = new(types.ID)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Connector", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Connector fields.
func (c *Connector) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case connector.FieldID:
			if value, ok := values[i].(*types.ID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				c.ID = *value
			}
		case connector.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				c.Name = value.String
			}
		case connector.FieldDescription:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field description", values[i])
			} else if value.Valid {
				c.Description = value.String
			}
		case connector.FieldLabels:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field labels", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &c.Labels); err != nil {
					return fmt.Errorf("unmarshal field labels: %w", err)
				}
			}
		case connector.FieldStatus:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field status", values[i])
			} else if value.Valid {
				c.Status = value.String
			}
		case connector.FieldStatusMessage:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field statusMessage", values[i])
			} else if value.Valid {
				c.StatusMessage = value.String
			}
		case connector.FieldCreateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field createTime", values[i])
			} else if value.Valid {
				c.CreateTime = new(time.Time)
				*c.CreateTime = value.Time
			}
		case connector.FieldUpdateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updateTime", values[i])
			} else if value.Valid {
				c.UpdateTime = new(time.Time)
				*c.UpdateTime = value.Time
			}
		case connector.FieldType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field type", values[i])
			} else if value.Valid {
				c.Type = value.String
			}
		case connector.FieldConfigVersion:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field configVersion", values[i])
			} else if value.Valid {
				c.ConfigVersion = value.String
			}
		case connector.FieldConfigData:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field configData", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &c.ConfigData); err != nil {
					return fmt.Errorf("unmarshal field configData: %w", err)
				}
			}
		case connector.FieldEnableFinOps:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field enableFinOps", values[i])
			} else if value.Valid {
				c.EnableFinOps = value.Bool
			}
		case connector.FieldFinOpsStatus:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field finOpsStatus", values[i])
			} else if value.Valid {
				c.FinOpsStatus = value.String
			}
		case connector.FieldFinOpsStatusMessage:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field finOpsStatusMessage", values[i])
			} else if value.Valid {
				c.FinOpsStatusMessage = value.String
			}
		}
	}
	return nil
}

// QueryEnvironments queries the "environments" edge of the Connector entity.
func (c *Connector) QueryEnvironments() *EnvironmentConnectorRelationshipQuery {
	return NewConnectorClient(c.config).QueryEnvironments(c)
}

// QueryResources queries the "resources" edge of the Connector entity.
func (c *Connector) QueryResources() *ApplicationResourceQuery {
	return NewConnectorClient(c.config).QueryResources(c)
}

// QueryClusterCosts queries the "clusterCosts" edge of the Connector entity.
func (c *Connector) QueryClusterCosts() *ClusterCostQuery {
	return NewConnectorClient(c.config).QueryClusterCosts(c)
}

// QueryAllocationCosts queries the "allocationCosts" edge of the Connector entity.
func (c *Connector) QueryAllocationCosts() *AllocationCostQuery {
	return NewConnectorClient(c.config).QueryAllocationCosts(c)
}

// Update returns a builder for updating this Connector.
// Note that you need to call Connector.Unwrap() before calling this method if this Connector
// was returned from a transaction, and the transaction was committed or rolled back.
func (c *Connector) Update() *ConnectorUpdateOne {
	return NewConnectorClient(c.config).UpdateOne(c)
}

// Unwrap unwraps the Connector entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (c *Connector) Unwrap() *Connector {
	_tx, ok := c.config.driver.(*txDriver)
	if !ok {
		panic("model: Connector is not a transactional entity")
	}
	c.config.driver = _tx.drv
	return c
}

// String implements the fmt.Stringer.
func (c *Connector) String() string {
	var builder strings.Builder
	builder.WriteString("Connector(")
	builder.WriteString(fmt.Sprintf("id=%v, ", c.ID))
	builder.WriteString("name=")
	builder.WriteString(c.Name)
	builder.WriteString(", ")
	builder.WriteString("description=")
	builder.WriteString(c.Description)
	builder.WriteString(", ")
	builder.WriteString("labels=")
	builder.WriteString(fmt.Sprintf("%v", c.Labels))
	builder.WriteString(", ")
	builder.WriteString("status=")
	builder.WriteString(c.Status)
	builder.WriteString(", ")
	builder.WriteString("statusMessage=")
	builder.WriteString(c.StatusMessage)
	builder.WriteString(", ")
	if v := c.CreateTime; v != nil {
		builder.WriteString("createTime=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	if v := c.UpdateTime; v != nil {
		builder.WriteString("updateTime=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("type=")
	builder.WriteString(c.Type)
	builder.WriteString(", ")
	builder.WriteString("configVersion=")
	builder.WriteString(c.ConfigVersion)
	builder.WriteString(", ")
	builder.WriteString("configData=")
	builder.WriteString(fmt.Sprintf("%v", c.ConfigData))
	builder.WriteString(", ")
	builder.WriteString("enableFinOps=")
	builder.WriteString(fmt.Sprintf("%v", c.EnableFinOps))
	builder.WriteString(", ")
	builder.WriteString("finOpsStatus=")
	builder.WriteString(c.FinOpsStatus)
	builder.WriteString(", ")
	builder.WriteString("finOpsStatusMessage=")
	builder.WriteString(c.FinOpsStatusMessage)
	builder.WriteByte(')')
	return builder.String()
}

// Connectors is a parsable slice of Connector.
type Connectors []*Connector
