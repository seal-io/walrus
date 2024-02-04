// SPDX-FileCopyrightText: 2024 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "walrus". DO NOT EDIT.

package model

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"

	"github.com/seal-io/walrus/pkg/dao/model/project"
	"github.com/seal-io/walrus/pkg/dao/model/workflow"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/utils/json"
)

// Workflow is the model entity for the Workflow schema.
type Workflow struct {
	config `json:"-"`
	// ID of the ent.
	ID object.ID `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Description holds the value of the "description" field.
	Description string `json:"description,omitempty"`
	// Labels holds the value of the "labels" field.
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations holds the value of the "annotations" field.
	Annotations map[string]string `json:"annotations,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime *time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime *time.Time `json:"update_time,omitempty"`
	// ID of the project that this workflow belongs to.
	ProjectID object.ID `json:"project_id,omitempty"`
	// ID of the environment that this workflow belongs to.
	EnvironmentID object.ID `json:"environment_id,omitempty"`
	// Type of the workflow.
	Type string `json:"type,omitempty"`
	// Number of task pods that can be executed in parallel of workflow.
	Parallelism int `json:"parallelism,omitempty"`
	// Timeout seconds of the workflow.
	Timeout int `json:"timeout,omitempty"`
	// Execution version of the workflow.
	Version int `json:"version,omitempty"`
	// Configs of workflow variables.
	Variables types.WorkflowVariables `json:"variables,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the WorkflowQuery when eager-loading is set.
	Edges        WorkflowEdges `json:"edges,omitempty"`
	selectValues sql.SelectValues
}

// WorkflowEdges holds the relations/edges for other nodes in the graph.
type WorkflowEdges struct {
	// Project to which the workflow belongs.
	Project *Project `json:"project,omitempty"`
	// Stages that belong to this workflow.
	Stages []*WorkflowStage `json:"stages,omitempty"`
	// Workflow executions that belong to this workflow.
	Executions []*WorkflowExecution `json:"executions,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [3]bool
}

// ProjectOrErr returns the Project value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e WorkflowEdges) ProjectOrErr() (*Project, error) {
	if e.loadedTypes[0] {
		if e.Project == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: project.Label}
		}
		return e.Project, nil
	}
	return nil, &NotLoadedError{edge: "project"}
}

// StagesOrErr returns the Stages value or an error if the edge
// was not loaded in eager-loading.
func (e WorkflowEdges) StagesOrErr() ([]*WorkflowStage, error) {
	if e.loadedTypes[1] {
		return e.Stages, nil
	}
	return nil, &NotLoadedError{edge: "stages"}
}

// ExecutionsOrErr returns the Executions value or an error if the edge
// was not loaded in eager-loading.
func (e WorkflowEdges) ExecutionsOrErr() ([]*WorkflowExecution, error) {
	if e.loadedTypes[2] {
		return e.Executions, nil
	}
	return nil, &NotLoadedError{edge: "executions"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Workflow) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case workflow.FieldLabels, workflow.FieldAnnotations, workflow.FieldVariables:
			values[i] = new([]byte)
		case workflow.FieldID, workflow.FieldProjectID, workflow.FieldEnvironmentID:
			values[i] = new(object.ID)
		case workflow.FieldParallelism, workflow.FieldTimeout, workflow.FieldVersion:
			values[i] = new(sql.NullInt64)
		case workflow.FieldName, workflow.FieldDescription, workflow.FieldType:
			values[i] = new(sql.NullString)
		case workflow.FieldCreateTime, workflow.FieldUpdateTime:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Workflow fields.
func (w *Workflow) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case workflow.FieldID:
			if value, ok := values[i].(*object.ID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				w.ID = *value
			}
		case workflow.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				w.Name = value.String
			}
		case workflow.FieldDescription:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field description", values[i])
			} else if value.Valid {
				w.Description = value.String
			}
		case workflow.FieldLabels:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field labels", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &w.Labels); err != nil {
					return fmt.Errorf("unmarshal field labels: %w", err)
				}
			}
		case workflow.FieldAnnotations:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field annotations", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &w.Annotations); err != nil {
					return fmt.Errorf("unmarshal field annotations: %w", err)
				}
			}
		case workflow.FieldCreateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field create_time", values[i])
			} else if value.Valid {
				w.CreateTime = new(time.Time)
				*w.CreateTime = value.Time
			}
		case workflow.FieldUpdateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field update_time", values[i])
			} else if value.Valid {
				w.UpdateTime = new(time.Time)
				*w.UpdateTime = value.Time
			}
		case workflow.FieldProjectID:
			if value, ok := values[i].(*object.ID); !ok {
				return fmt.Errorf("unexpected type %T for field project_id", values[i])
			} else if value != nil {
				w.ProjectID = *value
			}
		case workflow.FieldEnvironmentID:
			if value, ok := values[i].(*object.ID); !ok {
				return fmt.Errorf("unexpected type %T for field environment_id", values[i])
			} else if value != nil {
				w.EnvironmentID = *value
			}
		case workflow.FieldType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field type", values[i])
			} else if value.Valid {
				w.Type = value.String
			}
		case workflow.FieldParallelism:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field parallelism", values[i])
			} else if value.Valid {
				w.Parallelism = int(value.Int64)
			}
		case workflow.FieldTimeout:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field timeout", values[i])
			} else if value.Valid {
				w.Timeout = int(value.Int64)
			}
		case workflow.FieldVersion:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field version", values[i])
			} else if value.Valid {
				w.Version = int(value.Int64)
			}
		case workflow.FieldVariables:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field variables", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &w.Variables); err != nil {
					return fmt.Errorf("unmarshal field variables: %w", err)
				}
			}
		default:
			w.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Workflow.
// This includes values selected through modifiers, order, etc.
func (w *Workflow) Value(name string) (ent.Value, error) {
	return w.selectValues.Get(name)
}

// QueryProject queries the "project" edge of the Workflow entity.
func (w *Workflow) QueryProject() *ProjectQuery {
	return NewWorkflowClient(w.config).QueryProject(w)
}

// QueryStages queries the "stages" edge of the Workflow entity.
func (w *Workflow) QueryStages() *WorkflowStageQuery {
	return NewWorkflowClient(w.config).QueryStages(w)
}

// QueryExecutions queries the "executions" edge of the Workflow entity.
func (w *Workflow) QueryExecutions() *WorkflowExecutionQuery {
	return NewWorkflowClient(w.config).QueryExecutions(w)
}

// Update returns a builder for updating this Workflow.
// Note that you need to call Workflow.Unwrap() before calling this method if this Workflow
// was returned from a transaction, and the transaction was committed or rolled back.
func (w *Workflow) Update() *WorkflowUpdateOne {
	return NewWorkflowClient(w.config).UpdateOne(w)
}

// Unwrap unwraps the Workflow entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (w *Workflow) Unwrap() *Workflow {
	_tx, ok := w.config.driver.(*txDriver)
	if !ok {
		panic("model: Workflow is not a transactional entity")
	}
	w.config.driver = _tx.drv
	return w
}

// String implements the fmt.Stringer.
func (w *Workflow) String() string {
	var builder strings.Builder
	builder.WriteString("Workflow(")
	builder.WriteString(fmt.Sprintf("id=%v, ", w.ID))
	builder.WriteString("name=")
	builder.WriteString(w.Name)
	builder.WriteString(", ")
	builder.WriteString("description=")
	builder.WriteString(w.Description)
	builder.WriteString(", ")
	builder.WriteString("labels=")
	builder.WriteString(fmt.Sprintf("%v", w.Labels))
	builder.WriteString(", ")
	builder.WriteString("annotations=")
	builder.WriteString(fmt.Sprintf("%v", w.Annotations))
	builder.WriteString(", ")
	if v := w.CreateTime; v != nil {
		builder.WriteString("create_time=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	if v := w.UpdateTime; v != nil {
		builder.WriteString("update_time=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("project_id=")
	builder.WriteString(fmt.Sprintf("%v", w.ProjectID))
	builder.WriteString(", ")
	builder.WriteString("environment_id=")
	builder.WriteString(fmt.Sprintf("%v", w.EnvironmentID))
	builder.WriteString(", ")
	builder.WriteString("type=")
	builder.WriteString(w.Type)
	builder.WriteString(", ")
	builder.WriteString("parallelism=")
	builder.WriteString(fmt.Sprintf("%v", w.Parallelism))
	builder.WriteString(", ")
	builder.WriteString("timeout=")
	builder.WriteString(fmt.Sprintf("%v", w.Timeout))
	builder.WriteString(", ")
	builder.WriteString("version=")
	builder.WriteString(fmt.Sprintf("%v", w.Version))
	builder.WriteString(", ")
	builder.WriteString("variables=")
	builder.WriteString(fmt.Sprintf("%v", w.Variables))
	builder.WriteByte(')')
	return builder.String()
}

// Workflows is a parsable slice of Workflow.
type Workflows []*Workflow
