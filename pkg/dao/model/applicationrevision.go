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

	"github.com/seal-io/seal/pkg/dao/model/applicationrevision"
	"github.com/seal-io/seal/pkg/dao/oid"
	"github.com/seal-io/seal/pkg/dao/schema"
)

// ApplicationRevision is the model entity for the ApplicationRevision schema.
type ApplicationRevision struct {
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
	// ID of the application to which the revision belongs
	ApplicationID oid.ID `json:"applicationID"`
	// ID of the environment to which the application deploys
	EnvironmentID oid.ID `json:"environmentID"`
	// Application modules
	Modules []schema.ApplicationModule `json:"modules"`
	// Input variables of the revision
	InputVariables map[string]interface{} `json:"inputVariables"`
	// Input plan of the revision
	InputPlan string `json:"inputPlan"`
	// Output of the revision
	Output string `json:"output"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*ApplicationRevision) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case applicationrevision.FieldModules, applicationrevision.FieldInputVariables:
			values[i] = new([]byte)
		case applicationrevision.FieldID, applicationrevision.FieldApplicationID, applicationrevision.FieldEnvironmentID:
			values[i] = new(oid.ID)
		case applicationrevision.FieldStatus, applicationrevision.FieldStatusMessage, applicationrevision.FieldInputPlan, applicationrevision.FieldOutput:
			values[i] = new(sql.NullString)
		case applicationrevision.FieldCreateTime, applicationrevision.FieldUpdateTime:
			values[i] = new(sql.NullTime)
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
			if value, ok := values[i].(*oid.ID); !ok {
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
		case applicationrevision.FieldUpdateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updateTime", values[i])
			} else if value.Valid {
				ar.UpdateTime = new(time.Time)
				*ar.UpdateTime = value.Time
			}
		case applicationrevision.FieldApplicationID:
			if value, ok := values[i].(*oid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field applicationID", values[i])
			} else if value != nil {
				ar.ApplicationID = *value
			}
		case applicationrevision.FieldEnvironmentID:
			if value, ok := values[i].(*oid.ID); !ok {
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
		}
	}
	return nil
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
	if v := ar.UpdateTime; v != nil {
		builder.WriteString("updateTime=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("applicationID=")
	builder.WriteString(fmt.Sprintf("%v", ar.ApplicationID))
	builder.WriteString(", ")
	builder.WriteString("environmentID=")
	builder.WriteString(fmt.Sprintf("%v", ar.EnvironmentID))
	builder.WriteString(", ")
	builder.WriteString("modules=")
	builder.WriteString(fmt.Sprintf("%v", ar.Modules))
	builder.WriteString(", ")
	builder.WriteString("inputVariables=")
	builder.WriteString(fmt.Sprintf("%v", ar.InputVariables))
	builder.WriteString(", ")
	builder.WriteString("inputPlan=")
	builder.WriteString(ar.InputPlan)
	builder.WriteString(", ")
	builder.WriteString("output=")
	builder.WriteString(ar.Output)
	builder.WriteByte(')')
	return builder.String()
}

// ApplicationRevisions is a parsable slice of ApplicationRevision.
type ApplicationRevisions []*ApplicationRevision

func (ar ApplicationRevisions) config(cfg config) {
	for _i := range ar {
		ar[_i].config = cfg
	}
}
