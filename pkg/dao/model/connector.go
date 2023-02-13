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
	"github.com/seal-io/seal/pkg/dao/oid"
)

// Connector is the model entity for the Connector schema.
type Connector struct {
	config `json:"-"`
	// ID of the ent.
	// ID of the resource.
	ID oid.ID `json:"id,omitempty"`
	// Status of the resource
	Status string `json:"status,omitempty"`
	// extra message for status, like error details
	StatusMessage string `json:"statusMessage,omitempty"`
	// Describe creation time.
	CreateTime *time.Time `json:"createTime,omitempty"`
	// Describe modification time.
	UpdateTime *time.Time `json:"updateTime,omitempty"`
	// Driver type of the connector
	Driver string `json:"driver"`
	// Connector config version
	ConfigVersion string `json:"configVersion"`
	// Connector config data
	ConfigData map[string]interface{} `json:"configData,omitempty"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Connector) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case connector.FieldConfigData:
			values[i] = new([]byte)
		case connector.FieldID:
			values[i] = new(oid.ID)
		case connector.FieldStatus, connector.FieldStatusMessage, connector.FieldDriver, connector.FieldConfigVersion:
			values[i] = new(sql.NullString)
		case connector.FieldCreateTime, connector.FieldUpdateTime:
			values[i] = new(sql.NullTime)
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
			if value, ok := values[i].(*oid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				c.ID = *value
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
		case connector.FieldDriver:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field driver", values[i])
			} else if value.Valid {
				c.Driver = value.String
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
		}
	}
	return nil
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
	builder.WriteString("driver=")
	builder.WriteString(c.Driver)
	builder.WriteString(", ")
	builder.WriteString("configVersion=")
	builder.WriteString(c.ConfigVersion)
	builder.WriteString(", ")
	builder.WriteString("configData=")
	builder.WriteString(fmt.Sprintf("%v", c.ConfigData))
	builder.WriteByte(')')
	return builder.String()
}

// Connectors is a parsable slice of Connector.
type Connectors []*Connector

func (c Connectors) config(cfg config) {
	for _i := range c {
		c[_i].config = cfg
	}
}
