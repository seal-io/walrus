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

	"github.com/seal-io/seal/pkg/dao/model/module"
	"github.com/seal-io/seal/pkg/dao/oid"
)

// Module is the model entity for the Module schema.
type Module struct {
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
	// Source of the module
	Source string `json:"source"`
	// Version of the module
	Version string `json:"version"`
	// Input schema of the module
	InputSchema map[string]interface{} `json:"inputSchema,omitempty"`
	// Output schema of the module
	OutputSchema map[string]interface{} `json:"outputSchema,omitempty"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Module) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case module.FieldInputSchema, module.FieldOutputSchema:
			values[i] = new([]byte)
		case module.FieldID:
			values[i] = new(oid.ID)
		case module.FieldStatus, module.FieldStatusMessage, module.FieldSource, module.FieldVersion:
			values[i] = new(sql.NullString)
		case module.FieldCreateTime, module.FieldUpdateTime:
			values[i] = new(sql.NullTime)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Module", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Module fields.
func (m *Module) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case module.FieldID:
			if value, ok := values[i].(*oid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				m.ID = *value
			}
		case module.FieldStatus:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field status", values[i])
			} else if value.Valid {
				m.Status = value.String
			}
		case module.FieldStatusMessage:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field statusMessage", values[i])
			} else if value.Valid {
				m.StatusMessage = value.String
			}
		case module.FieldCreateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field createTime", values[i])
			} else if value.Valid {
				m.CreateTime = new(time.Time)
				*m.CreateTime = value.Time
			}
		case module.FieldUpdateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updateTime", values[i])
			} else if value.Valid {
				m.UpdateTime = new(time.Time)
				*m.UpdateTime = value.Time
			}
		case module.FieldSource:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field source", values[i])
			} else if value.Valid {
				m.Source = value.String
			}
		case module.FieldVersion:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field version", values[i])
			} else if value.Valid {
				m.Version = value.String
			}
		case module.FieldInputSchema:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field inputSchema", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &m.InputSchema); err != nil {
					return fmt.Errorf("unmarshal field inputSchema: %w", err)
				}
			}
		case module.FieldOutputSchema:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field outputSchema", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &m.OutputSchema); err != nil {
					return fmt.Errorf("unmarshal field outputSchema: %w", err)
				}
			}
		}
	}
	return nil
}

// Update returns a builder for updating this Module.
// Note that you need to call Module.Unwrap() before calling this method if this Module
// was returned from a transaction, and the transaction was committed or rolled back.
func (m *Module) Update() *ModuleUpdateOne {
	return NewModuleClient(m.config).UpdateOne(m)
}

// Unwrap unwraps the Module entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (m *Module) Unwrap() *Module {
	_tx, ok := m.config.driver.(*txDriver)
	if !ok {
		panic("model: Module is not a transactional entity")
	}
	m.config.driver = _tx.drv
	return m
}

// String implements the fmt.Stringer.
func (m *Module) String() string {
	var builder strings.Builder
	builder.WriteString("Module(")
	builder.WriteString(fmt.Sprintf("id=%v, ", m.ID))
	builder.WriteString("status=")
	builder.WriteString(m.Status)
	builder.WriteString(", ")
	builder.WriteString("statusMessage=")
	builder.WriteString(m.StatusMessage)
	builder.WriteString(", ")
	if v := m.CreateTime; v != nil {
		builder.WriteString("createTime=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	if v := m.UpdateTime; v != nil {
		builder.WriteString("updateTime=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("source=")
	builder.WriteString(m.Source)
	builder.WriteString(", ")
	builder.WriteString("version=")
	builder.WriteString(m.Version)
	builder.WriteString(", ")
	builder.WriteString("inputSchema=")
	builder.WriteString(fmt.Sprintf("%v", m.InputSchema))
	builder.WriteString(", ")
	builder.WriteString("outputSchema=")
	builder.WriteString(fmt.Sprintf("%v", m.OutputSchema))
	builder.WriteByte(')')
	return builder.String()
}

// Modules is a parsable slice of Module.
type Modules []*Module

func (m Modules) config(cfg config) {
	for _i := range m {
		m[_i].config = cfg
	}
}
