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

	"github.com/seal-io/seal/pkg/dao/model/applicationinstance"
	"github.com/seal-io/seal/pkg/dao/model/applicationrevision"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/types"
)

// ApplicationRevision is the model entity for the ApplicationRevision schema.
type ApplicationRevision struct {
	config `json:"-"`
	// ID of the ent.
	ID types.ID `json:"id,omitempty"`
	// Status of the resource.
	Status string `json:"status,omitempty"`
	// Extra message for status, like error details.
	StatusMessage string `json:"statusMessage,omitempty"`
	// Describe creation time.
	CreateTime *time.Time `json:"createTime,omitempty"`
	// ID of the application instance to which the revision belongs.
	InstanceID types.ID `json:"instanceID,omitempty"`
	// ID of the environment to which the application deploys.
	EnvironmentID types.ID `json:"environmentID,omitempty"`
	// Application modules.
	Modules []types.ApplicationModule `json:"modules,omitempty"`
	// Input variables of the revision.
	InputVariables map[string]interface{} `json:"-"`
	// Input plan of the revision.
	InputPlan string `json:"-"`
	// Output of the revision.
	Output string `json:"-"`
	// Type of deployer.
	DeployerType string `json:"deployerType,omitempty"`
	// Duration in seconds of the revision deploying.
	Duration int `json:"duration,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ApplicationRevisionQuery when eager-loading is set.
	Edges ApplicationRevisionEdges `json:"edges,omitempty"`
}

// ApplicationRevisionEdges holds the relations/edges for other nodes in the graph.
type ApplicationRevisionEdges struct {
	// Application instance to which the revision belongs.
	Instance *ApplicationInstance `json:"instance,omitempty"`
	// Environment to which the revision deploys.
	Environment *Environment `json:"environment,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// InstanceOrErr returns the Instance value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ApplicationRevisionEdges) InstanceOrErr() (*ApplicationInstance, error) {
	if e.loadedTypes[0] {
		if e.Instance == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: applicationinstance.Label}
		}
		return e.Instance, nil
	}
	return nil, &NotLoadedError{edge: "instance"}
}

// EnvironmentOrErr returns the Environment value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ApplicationRevisionEdges) EnvironmentOrErr() (*Environment, error) {
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
func (*ApplicationRevision) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case applicationrevision.FieldModules, applicationrevision.FieldInputVariables:
			values[i] = new([]byte)
		case applicationrevision.FieldDuration:
			values[i] = new(sql.NullInt64)
		case applicationrevision.FieldStatus, applicationrevision.FieldStatusMessage, applicationrevision.FieldInputPlan, applicationrevision.FieldOutput, applicationrevision.FieldDeployerType:
			values[i] = new(sql.NullString)
		case applicationrevision.FieldCreateTime:
			values[i] = new(sql.NullTime)
		case applicationrevision.FieldID, applicationrevision.FieldInstanceID, applicationrevision.FieldEnvironmentID:
			values[i] = new(types.ID)
		default:
			return nil, fmt.Errorf("unexpected column %q for type ApplicationRevision", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the ApplicationRevision fields.
func (ar *ApplicationRevision) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case applicationrevision.FieldID:
			if value, ok := values[i].(*types.ID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				ar.ID = *value
			}
		case applicationrevision.FieldStatus:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field status", values[i])
			} else if value.Valid {
				ar.Status = value.String
			}
		case applicationrevision.FieldStatusMessage:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field statusMessage", values[i])
			} else if value.Valid {
				ar.StatusMessage = value.String
			}
		case applicationrevision.FieldCreateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field createTime", values[i])
			} else if value.Valid {
				ar.CreateTime = new(time.Time)
				*ar.CreateTime = value.Time
			}
		case applicationrevision.FieldInstanceID:
			if value, ok := values[i].(*types.ID); !ok {
				return fmt.Errorf("unexpected type %T for field instanceID", values[i])
			} else if value != nil {
				ar.InstanceID = *value
			}
		case applicationrevision.FieldEnvironmentID:
			if value, ok := values[i].(*types.ID); !ok {
				return fmt.Errorf("unexpected type %T for field environmentID", values[i])
			} else if value != nil {
				ar.EnvironmentID = *value
			}
		case applicationrevision.FieldModules:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field modules", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &ar.Modules); err != nil {
					return fmt.Errorf("unmarshal field modules: %w", err)
				}
			}
		case applicationrevision.FieldInputVariables:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field inputVariables", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &ar.InputVariables); err != nil {
					return fmt.Errorf("unmarshal field inputVariables: %w", err)
				}
			}
		case applicationrevision.FieldInputPlan:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field inputPlan", values[i])
			} else if value.Valid {
				ar.InputPlan = value.String
			}
		case applicationrevision.FieldOutput:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field output", values[i])
			} else if value.Valid {
				ar.Output = value.String
			}
		case applicationrevision.FieldDeployerType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field deployerType", values[i])
			} else if value.Valid {
				ar.DeployerType = value.String
			}
		case applicationrevision.FieldDuration:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field duration", values[i])
			} else if value.Valid {
				ar.Duration = int(value.Int64)
			}
		}
	}
	return nil
}

// QueryInstance queries the "instance" edge of the ApplicationRevision entity.
func (ar *ApplicationRevision) QueryInstance() *ApplicationInstanceQuery {
	return NewApplicationRevisionClient(ar.config).QueryInstance(ar)
}

// QueryEnvironment queries the "environment" edge of the ApplicationRevision entity.
func (ar *ApplicationRevision) QueryEnvironment() *EnvironmentQuery {
	return NewApplicationRevisionClient(ar.config).QueryEnvironment(ar)
}

// Update returns a builder for updating this ApplicationRevision.
// Note that you need to call ApplicationRevision.Unwrap() before calling this method if this ApplicationRevision
// was returned from a transaction, and the transaction was committed or rolled back.
func (ar *ApplicationRevision) Update() *ApplicationRevisionUpdateOne {
	return NewApplicationRevisionClient(ar.config).UpdateOne(ar)
}

// Unwrap unwraps the ApplicationRevision entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (ar *ApplicationRevision) Unwrap() *ApplicationRevision {
	_tx, ok := ar.config.driver.(*txDriver)
	if !ok {
		panic("model: ApplicationRevision is not a transactional entity")
	}
	ar.config.driver = _tx.drv
	return ar
}

// String implements the fmt.Stringer.
func (ar *ApplicationRevision) String() string {
	var builder strings.Builder
	builder.WriteString("ApplicationRevision(")
	builder.WriteString(fmt.Sprintf("id=%v, ", ar.ID))
	builder.WriteString("status=")
	builder.WriteString(ar.Status)
	builder.WriteString(", ")
	builder.WriteString("statusMessage=")
	builder.WriteString(ar.StatusMessage)
	builder.WriteString(", ")
	if v := ar.CreateTime; v != nil {
		builder.WriteString("createTime=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("instanceID=")
	builder.WriteString(fmt.Sprintf("%v", ar.InstanceID))
	builder.WriteString(", ")
	builder.WriteString("environmentID=")
	builder.WriteString(fmt.Sprintf("%v", ar.EnvironmentID))
	builder.WriteString(", ")
	builder.WriteString("modules=")
	builder.WriteString(fmt.Sprintf("%v", ar.Modules))
	builder.WriteString(", ")
	builder.WriteString("inputVariables=<sensitive>")
	builder.WriteString(", ")
	builder.WriteString("inputPlan=<sensitive>")
	builder.WriteString(", ")
	builder.WriteString("output=<sensitive>")
	builder.WriteString(", ")
	builder.WriteString("deployerType=")
	builder.WriteString(ar.DeployerType)
	builder.WriteString(", ")
	builder.WriteString("duration=")
	builder.WriteString(fmt.Sprintf("%v", ar.Duration))
	builder.WriteByte(')')
	return builder.String()
}

// ApplicationRevisions is a parsable slice of ApplicationRevision.
type ApplicationRevisions []*ApplicationRevision
