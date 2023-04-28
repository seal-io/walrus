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

	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/model/project"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/pkg/dao/types/property"
	"github.com/seal-io/seal/utils/json"
)

// Application is the model entity for the Application schema.
type Application struct {
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
	// ID of the project to which the application belongs.
	ProjectID oid.ID `json:"projectID,omitempty" sql:"projectID"`
	// Variables definition of the application, the variables of instance derived by this definition
	Variables property.Schemas `json:"variables,omitempty" sql:"variables"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ApplicationQuery when eager-loading is set.
	Edges        ApplicationEdges `json:"edges,omitempty"`
	selectValues sql.SelectValues
}

// ApplicationEdges holds the relations/edges for other nodes in the graph.
type ApplicationEdges struct {
	// Project to which this application belongs.
	Project *Project `json:"project,omitempty" sql:"project"`
	// Application instances that belong to this application.
	Instances []*ApplicationInstance `json:"instances,omitempty" sql:"instances"`
	// Modules holds the value of the modules edge.
	Modules []*ApplicationModuleRelationship `json:"modules,omitempty" sql:"modules"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [3]bool
}

// ProjectOrErr returns the Project value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ApplicationEdges) ProjectOrErr() (*Project, error) {
	if e.loadedTypes[0] {
		if e.Project == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: project.Label}
		}
		return e.Project, nil
	}
	return nil, &NotLoadedError{edge: "project"}
}

// InstancesOrErr returns the Instances value or an error if the edge
// was not loaded in eager-loading.
func (e ApplicationEdges) InstancesOrErr() ([]*ApplicationInstance, error) {
	if e.loadedTypes[1] {
		return e.Instances, nil
	}
	return nil, &NotLoadedError{edge: "instances"}
}

// ModulesOrErr returns the Modules value or an error if the edge
// was not loaded in eager-loading.
func (e ApplicationEdges) ModulesOrErr() ([]*ApplicationModuleRelationship, error) {
	if e.loadedTypes[2] {
		return e.Modules, nil
	}
	return nil, &NotLoadedError{edge: "modules"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Application) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case application.FieldLabels:
			values[i] = new([]byte)
		case application.FieldID, application.FieldProjectID:
			values[i] = new(oid.ID)
		case application.FieldVariables:
			values[i] = new(property.Schemas)
		case application.FieldName, application.FieldDescription:
			values[i] = new(sql.NullString)
		case application.FieldCreateTime, application.FieldUpdateTime:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Application fields.
func (a *Application) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case application.FieldID:
			if value, ok := values[i].(*oid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				a.ID = *value
			}
		case application.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				a.Name = value.String
			}
		case application.FieldDescription:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field description", values[i])
			} else if value.Valid {
				a.Description = value.String
			}
		case application.FieldLabels:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field labels", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &a.Labels); err != nil {
					return fmt.Errorf("unmarshal field labels: %w", err)
				}
			}
		case application.FieldCreateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field createTime", values[i])
			} else if value.Valid {
				a.CreateTime = new(time.Time)
				*a.CreateTime = value.Time
			}
		case application.FieldUpdateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updateTime", values[i])
			} else if value.Valid {
				a.UpdateTime = new(time.Time)
				*a.UpdateTime = value.Time
			}
		case application.FieldProjectID:
			if value, ok := values[i].(*oid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field projectID", values[i])
			} else if value != nil {
				a.ProjectID = *value
			}
		case application.FieldVariables:
			if value, ok := values[i].(*property.Schemas); !ok {
				return fmt.Errorf("unexpected type %T for field variables", values[i])
			} else if value != nil {
				a.Variables = *value
			}
		default:
			a.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Application.
// This includes values selected through modifiers, order, etc.
func (a *Application) Value(name string) (ent.Value, error) {
	return a.selectValues.Get(name)
}

// QueryProject queries the "project" edge of the Application entity.
func (a *Application) QueryProject() *ProjectQuery {
	return NewApplicationClient(a.config).QueryProject(a)
}

// QueryInstances queries the "instances" edge of the Application entity.
func (a *Application) QueryInstances() *ApplicationInstanceQuery {
	return NewApplicationClient(a.config).QueryInstances(a)
}

// QueryModules queries the "modules" edge of the Application entity.
func (a *Application) QueryModules() *ApplicationModuleRelationshipQuery {
	return NewApplicationClient(a.config).QueryModules(a)
}

// Update returns a builder for updating this Application.
// Note that you need to call Application.Unwrap() before calling this method if this Application
// was returned from a transaction, and the transaction was committed or rolled back.
func (a *Application) Update() *ApplicationUpdateOne {
	return NewApplicationClient(a.config).UpdateOne(a)
}

// Unwrap unwraps the Application entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (a *Application) Unwrap() *Application {
	_tx, ok := a.config.driver.(*txDriver)
	if !ok {
		panic("model: Application is not a transactional entity")
	}
	a.config.driver = _tx.drv
	return a
}

// String implements the fmt.Stringer.
func (a *Application) String() string {
	var builder strings.Builder
	builder.WriteString("Application(")
	builder.WriteString(fmt.Sprintf("id=%v, ", a.ID))
	builder.WriteString("name=")
	builder.WriteString(a.Name)
	builder.WriteString(", ")
	builder.WriteString("description=")
	builder.WriteString(a.Description)
	builder.WriteString(", ")
	builder.WriteString("labels=")
	builder.WriteString(fmt.Sprintf("%v", a.Labels))
	builder.WriteString(", ")
	if v := a.CreateTime; v != nil {
		builder.WriteString("createTime=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	if v := a.UpdateTime; v != nil {
		builder.WriteString("updateTime=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("projectID=")
	builder.WriteString(fmt.Sprintf("%v", a.ProjectID))
	builder.WriteString(", ")
	builder.WriteString("variables=")
	builder.WriteString(fmt.Sprintf("%v", a.Variables))
	builder.WriteByte(')')
	return builder.String()
}

// Applications is a parsable slice of Application.
type Applications []*Application
