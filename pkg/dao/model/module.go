// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package model

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"

	"github.com/seal-io/seal/pkg/dao/model/module"
	"github.com/seal-io/seal/utils/json"
)

// Module is the model entity for the Module schema.
type Module struct {
	config `json:"-"`
	// ID of the ent.
	// It is also the name of the module.
	ID string `json:"id,omitempty" sql:"id"`
	// Status of the resource.
	Status string `json:"status,omitempty" sql:"status"`
	// Extra message for status, like error details.
	StatusMessage string `json:"statusMessage,omitempty" sql:"statusMessage"`
	// Describe creation time.
	CreateTime *time.Time `json:"createTime,omitempty" sql:"createTime"`
	// Describe modification time.
	UpdateTime *time.Time `json:"updateTime,omitempty" sql:"updateTime"`
	// Description of the module.
	Description string `json:"description,omitempty" sql:"description"`
	// A URL to an SVG or PNG image to be used as an icon.
	Icon string `json:"icon,omitempty" sql:"icon"`
	// Labels of the module.
	Labels map[string]string `json:"labels,omitempty" sql:"labels"`
	// Source of the module.
	Source string `json:"source,omitempty" sql:"source"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ModuleQuery when eager-loading is set.
	Edges ModuleEdges `json:"edges,omitempty"`
}

// ModuleEdges holds the relations/edges for other nodes in the graph.
type ModuleEdges struct {
	// Applications holds the value of the applications edge.
	Applications []*ApplicationModuleRelationship `json:"applications,omitempty" sql:"applications"`
	// versions of the module.
	Versions []*ModuleVersion `json:"versions,omitempty" sql:"versions"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// ApplicationsOrErr returns the Applications value or an error if the edge
// was not loaded in eager-loading.
func (e ModuleEdges) ApplicationsOrErr() ([]*ApplicationModuleRelationship, error) {
	if e.loadedTypes[0] {
		return e.Applications, nil
	}
	return nil, &NotLoadedError{edge: "applications"}
}

// VersionsOrErr returns the Versions value or an error if the edge
// was not loaded in eager-loading.
func (e ModuleEdges) VersionsOrErr() ([]*ModuleVersion, error) {
	if e.loadedTypes[1] {
		return e.Versions, nil
	}
	return nil, &NotLoadedError{edge: "versions"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Module) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case module.FieldLabels:
			values[i] = new([]byte)
		case module.FieldID, module.FieldStatus, module.FieldStatusMessage, module.FieldDescription, module.FieldIcon, module.FieldSource:
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
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value.Valid {
				m.ID = value.String
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
		case module.FieldDescription:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field description", values[i])
			} else if value.Valid {
				m.Description = value.String
			}
		case module.FieldIcon:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field icon", values[i])
			} else if value.Valid {
				m.Icon = value.String
			}
		case module.FieldLabels:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field labels", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &m.Labels); err != nil {
					return fmt.Errorf("unmarshal field labels: %w", err)
				}
			}
		case module.FieldSource:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field source", values[i])
			} else if value.Valid {
				m.Source = value.String
			}
		}
	}
	return nil
}

// QueryApplications queries the "applications" edge of the Module entity.
func (m *Module) QueryApplications() *ApplicationModuleRelationshipQuery {
	return NewModuleClient(m.config).QueryApplications(m)
}

// QueryVersions queries the "versions" edge of the Module entity.
func (m *Module) QueryVersions() *ModuleVersionQuery {
	return NewModuleClient(m.config).QueryVersions(m)
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
	builder.WriteString("description=")
	builder.WriteString(m.Description)
	builder.WriteString(", ")
	builder.WriteString("icon=")
	builder.WriteString(m.Icon)
	builder.WriteString(", ")
	builder.WriteString("labels=")
	builder.WriteString(fmt.Sprintf("%v", m.Labels))
	builder.WriteString(", ")
	builder.WriteString("source=")
	builder.WriteString(m.Source)
	builder.WriteByte(')')
	return builder.String()
}

// Modules is a parsable slice of Module.
type Modules []*Module
