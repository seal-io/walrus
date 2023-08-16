// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "seal". DO NOT EDIT.

package model

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"

	"github.com/seal-io/walrus/pkg/dao/model/environment"
	"github.com/seal-io/walrus/pkg/dao/model/project"
	"github.com/seal-io/walrus/pkg/dao/model/variable"
	"github.com/seal-io/walrus/pkg/dao/types/crypto"
	"github.com/seal-io/walrus/pkg/dao/types/object"
)

// Variable is the model entity for the Variable schema.
type Variable struct {
	config `json:"-"`
	// ID of the ent.
	ID object.ID `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime *time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime *time.Time `json:"update_time,omitempty"`
	// ID of the project to belong, empty means for all projects.
	ProjectID object.ID `json:"project_id,omitempty"`
	// ID of the environment to which the variable belongs to.
	EnvironmentID object.ID `json:"environment_id,omitempty"`
	// The name of variable.
	Name string `json:"name,omitempty"`
	// The value of variable, store in string.
	Value crypto.String `json:"value,omitempty"`
	// The value is sensitive or not.
	Sensitive bool `json:"sensitive,omitempty"`
	// Description of the variable.
	Description string `json:"description,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the VariableQuery when eager-loading is set.
	Edges        VariableEdges `json:"edges,omitempty"`
	selectValues sql.SelectValues
}

// VariableEdges holds the relations/edges for other nodes in the graph.
type VariableEdges struct {
	// Project to which the variable belongs.
	Project *Project `json:"project,omitempty"`
	// Environment to which the variable belongs.
	Environment *Environment `json:"environment,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// ProjectOrErr returns the Project value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e VariableEdges) ProjectOrErr() (*Project, error) {
	if e.loadedTypes[0] {
		if e.Project == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: project.Label}
		}
		return e.Project, nil
	}
	return nil, &NotLoadedError{edge: "project"}
}

// EnvironmentOrErr returns the Environment value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e VariableEdges) EnvironmentOrErr() (*Environment, error) {
	if e.loadedTypes[1] {
		if e.Environment == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: environment.Label}
		}
		return e.Environment, nil
	}
	return nil, &NotLoadedError{edge: "environment"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Variable) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case variable.FieldValue:
			values[i] = new(crypto.String)
		case variable.FieldID, variable.FieldProjectID, variable.FieldEnvironmentID:
			values[i] = new(object.ID)
		case variable.FieldSensitive:
			values[i] = new(sql.NullBool)
		case variable.FieldName, variable.FieldDescription:
			values[i] = new(sql.NullString)
		case variable.FieldCreateTime, variable.FieldUpdateTime:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Variable fields.
func (v *Variable) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case variable.FieldID:
			if value, ok := values[i].(*object.ID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				v.ID = *value
			}
		case variable.FieldCreateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field create_time", values[i])
			} else if value.Valid {
				v.CreateTime = new(time.Time)
				*v.CreateTime = value.Time
			}
		case variable.FieldUpdateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field update_time", values[i])
			} else if value.Valid {
				v.UpdateTime = new(time.Time)
				*v.UpdateTime = value.Time
			}
		case variable.FieldProjectID:
			if value, ok := values[i].(*object.ID); !ok {
				return fmt.Errorf("unexpected type %T for field project_id", values[i])
			} else if value != nil {
				v.ProjectID = *value
			}
		case variable.FieldEnvironmentID:
			if value, ok := values[i].(*object.ID); !ok {
				return fmt.Errorf("unexpected type %T for field environment_id", values[i])
			} else if value != nil {
				v.EnvironmentID = *value
			}
		case variable.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				v.Name = value.String
			}
		case variable.FieldValue:
			if value, ok := values[i].(*crypto.String); !ok {
				return fmt.Errorf("unexpected type %T for field value", values[i])
			} else if value != nil {
				v.Value = *value
			}
		case variable.FieldSensitive:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field sensitive", values[i])
			} else if value.Valid {
				v.Sensitive = value.Bool
			}
		case variable.FieldDescription:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field description", values[i])
			} else if value.Valid {
				v.Description = value.String
			}
		default:
			v.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// GetValue returns the ent.Value that was dynamically selected and assigned to the Variable.
// This includes values selected through modifiers, order, etc.
func (v *Variable) GetValue(name string) (ent.Value, error) {
	return v.selectValues.Get(name)
}

// QueryProject queries the "project" edge of the Variable entity.
func (v *Variable) QueryProject() *ProjectQuery {
	return NewVariableClient(v.config).QueryProject(v)
}

// QueryEnvironment queries the "environment" edge of the Variable entity.
func (v *Variable) QueryEnvironment() *EnvironmentQuery {
	return NewVariableClient(v.config).QueryEnvironment(v)
}

// Update returns a builder for updating this Variable.
// Note that you need to call Variable.Unwrap() before calling this method if this Variable
// was returned from a transaction, and the transaction was committed or rolled back.
func (v *Variable) Update() *VariableUpdateOne {
	return NewVariableClient(v.config).UpdateOne(v)
}

// Unwrap unwraps the Variable entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (v *Variable) Unwrap() *Variable {
	_tx, ok := v.config.driver.(*txDriver)
	if !ok {
		panic("model: Variable is not a transactional entity")
	}
	v.config.driver = _tx.drv
	return v
}

// String implements the fmt.Stringer.
func (v *Variable) String() string {
	var builder strings.Builder
	builder.WriteString("Variable(")
	builder.WriteString(fmt.Sprintf("id=%v, ", v.ID))
	if v := v.CreateTime; v != nil {
		builder.WriteString("create_time=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	if v := v.UpdateTime; v != nil {
		builder.WriteString("update_time=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("project_id=")
	builder.WriteString(fmt.Sprintf("%v", v.ProjectID))
	builder.WriteString(", ")
	builder.WriteString("environment_id=")
	builder.WriteString(fmt.Sprintf("%v", v.EnvironmentID))
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(v.Name)
	builder.WriteString(", ")
	builder.WriteString("value=")
	builder.WriteString(fmt.Sprintf("%v", v.Value))
	builder.WriteString(", ")
	builder.WriteString("sensitive=")
	builder.WriteString(fmt.Sprintf("%v", v.Sensitive))
	builder.WriteString(", ")
	builder.WriteString("description=")
	builder.WriteString(v.Description)
	builder.WriteByte(')')
	return builder.String()
}

// Variables is a parsable slice of Variable.
type Variables []*Variable
