// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package model

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"

	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/utils/json"
)

// Environment is the model entity for the Environment schema.
type Environment struct {
	config `json:"-"`
	// ID of the ent.
	ID oid.ID `json:"id,omitempty" sql:"id"`
	// Name of the resource.
	Name string `json:"name,omitempty" sql:"name"`
	// Description of the resource.
	Description string `json:"description,omitempty" sql:"description"`
	// Labels of the resource.
	Labels map[string]string `json:"labels,omitempty" sql:"labels"`
	// Describe creation time.
	CreateTime *time.Time `json:"createTime,omitempty" sql:"createTime"`
	// Describe modification time.
	UpdateTime *time.Time `json:"updateTime,omitempty" sql:"updateTime"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the EnvironmentQuery when eager-loading is set.
	Edges        EnvironmentEdges `json:"edges,omitempty"`
	selectValues sql.SelectValues
}

// EnvironmentEdges holds the relations/edges for other nodes in the graph.
type EnvironmentEdges struct {
	// Connectors holds the value of the connectors edge.
	Connectors []*EnvironmentConnectorRelationship `json:"connectors,omitempty" sql:"connectors"`
	// Application instances that belong to the environment.
	Instances []*ApplicationInstance `json:"instances,omitempty" sql:"instances"`
	// Application revisions that belong to the environment.
	Revisions []*ApplicationRevision `json:"revisions,omitempty" sql:"revisions"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [3]bool
}

// ConnectorsOrErr returns the Connectors value or an error if the edge
// was not loaded in eager-loading.
func (e EnvironmentEdges) ConnectorsOrErr() ([]*EnvironmentConnectorRelationship, error) {
	if e.loadedTypes[0] {
		return e.Connectors, nil
	}
	return nil, &NotLoadedError{edge: "connectors"}
}

// InstancesOrErr returns the Instances value or an error if the edge
// was not loaded in eager-loading.
func (e EnvironmentEdges) InstancesOrErr() ([]*ApplicationInstance, error) {
	if e.loadedTypes[1] {
		return e.Instances, nil
	}
	return nil, &NotLoadedError{edge: "instances"}
}

// RevisionsOrErr returns the Revisions value or an error if the edge
// was not loaded in eager-loading.
func (e EnvironmentEdges) RevisionsOrErr() ([]*ApplicationRevision, error) {
	if e.loadedTypes[2] {
		return e.Revisions, nil
	}
	return nil, &NotLoadedError{edge: "revisions"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Environment) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case environment.FieldLabels:
			values[i] = new([]byte)
		case environment.FieldID:
			values[i] = new(oid.ID)
		case environment.FieldName, environment.FieldDescription:
			values[i] = new(sql.NullString)
		case environment.FieldCreateTime, environment.FieldUpdateTime:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
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
		case environment.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				e.Name = value.String
			}
		case environment.FieldDescription:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field description", values[i])
			} else if value.Valid {
				e.Description = value.String
			}
		case environment.FieldLabels:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field labels", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &e.Labels); err != nil {
					return fmt.Errorf("unmarshal field labels: %w", err)
				}
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
		default:
			e.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Environment.
// This includes values selected through modifiers, order, etc.
func (e *Environment) Value(name string) (ent.Value, error) {
	return e.selectValues.Get(name)
}

// QueryConnectors queries the "connectors" edge of the Environment entity.
func (e *Environment) QueryConnectors() *EnvironmentConnectorRelationshipQuery {
	return NewEnvironmentClient(e.config).QueryConnectors(e)
}

// QueryInstances queries the "instances" edge of the Environment entity.
func (e *Environment) QueryInstances() *ApplicationInstanceQuery {
	return NewEnvironmentClient(e.config).QueryInstances(e)
}

// QueryRevisions queries the "revisions" edge of the Environment entity.
func (e *Environment) QueryRevisions() *ApplicationRevisionQuery {
	return NewEnvironmentClient(e.config).QueryRevisions(e)
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
	builder.WriteString("name=")
	builder.WriteString(e.Name)
	builder.WriteString(", ")
	builder.WriteString("description=")
	builder.WriteString(e.Description)
	builder.WriteString(", ")
	builder.WriteString("labels=")
	builder.WriteString(fmt.Sprintf("%v", e.Labels))
	builder.WriteString(", ")
	if v := e.CreateTime; v != nil {
		builder.WriteString("createTime=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	if v := e.UpdateTime; v != nil {
		builder.WriteString("updateTime=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteByte(')')
	return builder.String()
}

// Environments is a parsable slice of Environment.
type Environments []*Environment
