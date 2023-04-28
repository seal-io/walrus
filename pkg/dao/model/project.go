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

	"github.com/seal-io/seal/pkg/dao/model/project"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/utils/json"
)

// Project is the model entity for the Project schema.
type Project struct {
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
	// The values are being populated by the ProjectQuery when eager-loading is set.
	Edges        ProjectEdges `json:"edges,omitempty"`
	selectValues sql.SelectValues
}

// ProjectEdges holds the relations/edges for other nodes in the graph.
type ProjectEdges struct {
	// Applications that belong to the project.
	Applications []*Application `json:"applications,omitempty" sql:"applications"`
	// Secrets that belong to the project.
	Secrets []*Secret `json:"secrets,omitempty" sql:"secrets"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// ApplicationsOrErr returns the Applications value or an error if the edge
// was not loaded in eager-loading.
func (e ProjectEdges) ApplicationsOrErr() ([]*Application, error) {
	if e.loadedTypes[0] {
		return e.Applications, nil
	}
	return nil, &NotLoadedError{edge: "applications"}
}

// SecretsOrErr returns the Secrets value or an error if the edge
// was not loaded in eager-loading.
func (e ProjectEdges) SecretsOrErr() ([]*Secret, error) {
	if e.loadedTypes[1] {
		return e.Secrets, nil
	}
	return nil, &NotLoadedError{edge: "secrets"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Project) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case project.FieldLabels:
			values[i] = new([]byte)
		case project.FieldID:
			values[i] = new(oid.ID)
		case project.FieldName, project.FieldDescription:
			values[i] = new(sql.NullString)
		case project.FieldCreateTime, project.FieldUpdateTime:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Project fields.
func (pr *Project) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case project.FieldID:
			if value, ok := values[i].(*oid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				pr.ID = *value
			}
		case project.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				pr.Name = value.String
			}
		case project.FieldDescription:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field description", values[i])
			} else if value.Valid {
				pr.Description = value.String
			}
		case project.FieldLabels:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field labels", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &pr.Labels); err != nil {
					return fmt.Errorf("unmarshal field labels: %w", err)
				}
			}
		case project.FieldCreateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field createTime", values[i])
			} else if value.Valid {
				pr.CreateTime = new(time.Time)
				*pr.CreateTime = value.Time
			}
		case project.FieldUpdateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updateTime", values[i])
			} else if value.Valid {
				pr.UpdateTime = new(time.Time)
				*pr.UpdateTime = value.Time
			}
		default:
			pr.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Project.
// This includes values selected through modifiers, order, etc.
func (pr *Project) Value(name string) (ent.Value, error) {
	return pr.selectValues.Get(name)
}

// QueryApplications queries the "applications" edge of the Project entity.
func (pr *Project) QueryApplications() *ApplicationQuery {
	return NewProjectClient(pr.config).QueryApplications(pr)
}

// QuerySecrets queries the "secrets" edge of the Project entity.
func (pr *Project) QuerySecrets() *SecretQuery {
	return NewProjectClient(pr.config).QuerySecrets(pr)
}

// Update returns a builder for updating this Project.
// Note that you need to call Project.Unwrap() before calling this method if this Project
// was returned from a transaction, and the transaction was committed or rolled back.
func (pr *Project) Update() *ProjectUpdateOne {
	return NewProjectClient(pr.config).UpdateOne(pr)
}

// Unwrap unwraps the Project entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (pr *Project) Unwrap() *Project {
	_tx, ok := pr.config.driver.(*txDriver)
	if !ok {
		panic("model: Project is not a transactional entity")
	}
	pr.config.driver = _tx.drv
	return pr
}

// String implements the fmt.Stringer.
func (pr *Project) String() string {
	var builder strings.Builder
	builder.WriteString("Project(")
	builder.WriteString(fmt.Sprintf("id=%v, ", pr.ID))
	builder.WriteString("name=")
	builder.WriteString(pr.Name)
	builder.WriteString(", ")
	builder.WriteString("description=")
	builder.WriteString(pr.Description)
	builder.WriteString(", ")
	builder.WriteString("labels=")
	builder.WriteString(fmt.Sprintf("%v", pr.Labels))
	builder.WriteString(", ")
	if v := pr.CreateTime; v != nil {
		builder.WriteString("createTime=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	if v := pr.UpdateTime; v != nil {
		builder.WriteString("updateTime=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteByte(')')
	return builder.String()
}

// Projects is a parsable slice of Project.
type Projects []*Project
