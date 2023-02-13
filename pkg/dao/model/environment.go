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

	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/oid"
)

// Environment is the model entity for the Environment schema.
type Environment struct {
	config `json:"-"`
	// ID of the ent.
	// ID of the resource.
	ID oid.ID `json:"id,omitempty"`
	// Describe creation time.
	CreateTime *time.Time `json:"createTime,omitempty"`
	// Describe modification time.
	UpdateTime *time.Time `json:"updateTime,omitempty"`
	// ID of connectors of the environment
	ConnectorIDs []oid.ID `json:"connectorIDs,omitempty"`
	// Variables of the environment
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Environment) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case environment.FieldConnectorIDs, environment.FieldVariables:
			values[i] = new([]byte)
		case environment.FieldID:
			values[i] = new(oid.ID)
		case environment.FieldCreateTime, environment.FieldUpdateTime:
			values[i] = new(sql.NullTime)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Environment", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Environment fields.
func (e *Environment) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case environment.FieldID:
			if value, ok := values[i].(*oid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				e.ID = *value
			}
		case environment.FieldCreateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field createTime", values[i])
			} else if value.Valid {
				e.CreateTime = new(time.Time)
				*e.CreateTime = value.Time
			}
		case environment.FieldUpdateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updateTime", values[i])
			} else if value.Valid {
				e.UpdateTime = new(time.Time)
				*e.UpdateTime = value.Time
			}
		case environment.FieldConnectorIDs:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field connectorIDs", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &e.ConnectorIDs); err != nil {
					return fmt.Errorf("unmarshal field connectorIDs: %w", err)
				}
			}
		case environment.FieldVariables:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field variables", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &e.Variables); err != nil {
					return fmt.Errorf("unmarshal field variables: %w", err)
				}
			}
		}
	}
	return nil
}

// Update returns a builder for updating this Environment.
// Note that you need to call Environment.Unwrap() before calling this method if this Environment
// was returned from a transaction, and the transaction was committed or rolled back.
func (e *Environment) Update() *EnvironmentUpdateOne {
	return NewEnvironmentClient(e.config).UpdateOne(e)
}

// Unwrap unwraps the Environment entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (e *Environment) Unwrap() *Environment {
	_tx, ok := e.config.driver.(*txDriver)
	if !ok {
		panic("model: Environment is not a transactional entity")
	}
	e.config.driver = _tx.drv
	return e
}

// String implements the fmt.Stringer.
func (e *Environment) String() string {
	var builder strings.Builder
	builder.WriteString("Environment(")
	builder.WriteString(fmt.Sprintf("id=%v, ", e.ID))
	if v := e.CreateTime; v != nil {
		builder.WriteString("createTime=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	if v := e.UpdateTime; v != nil {
		builder.WriteString("updateTime=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("connectorIDs=")
	builder.WriteString(fmt.Sprintf("%v", e.ConnectorIDs))
	builder.WriteString(", ")
	builder.WriteString("variables=")
	builder.WriteString(fmt.Sprintf("%v", e.Variables))
	builder.WriteByte(')')
	return builder.String()
}

// Environments is a parsable slice of Environment.
type Environments []*Environment

func (e Environments) config(cfg config) {
	for _i := range e {
		e[_i].config = cfg
	}
}
