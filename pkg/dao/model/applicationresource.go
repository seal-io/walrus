// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package model

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"

	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/oid"
)

// ApplicationResource is the model entity for the ApplicationResource schema.
type ApplicationResource struct {
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
	// Module that generates the resource
	Module string `json:"module"`
	// Resource type
	Type string `json:"type"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*ApplicationResource) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case applicationresource.FieldID, applicationresource.FieldApplicationID:
			values[i] = new(oid.ID)
		case applicationresource.FieldStatus, applicationresource.FieldStatusMessage, applicationresource.FieldModule, applicationresource.FieldType:
			values[i] = new(sql.NullString)
		case applicationresource.FieldCreateTime, applicationresource.FieldUpdateTime:
			values[i] = new(sql.NullTime)
		default:
			return nil, fmt.Errorf("unexpected column %q for type ApplicationResource", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the ApplicationResource fields.
func (ar *ApplicationResource) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case applicationresource.FieldID:
			if value, ok := values[i].(*oid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				ar.ID = *value
			}
		case applicationresource.FieldStatus:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field status", values[i])
			} else if value.Valid {
				ar.Status = value.String
			}
		case applicationresource.FieldStatusMessage:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field statusMessage", values[i])
			} else if value.Valid {
				ar.StatusMessage = value.String
			}
		case applicationresource.FieldCreateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field createTime", values[i])
			} else if value.Valid {
				ar.CreateTime = new(time.Time)
				*ar.CreateTime = value.Time
			}
		case applicationresource.FieldUpdateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updateTime", values[i])
			} else if value.Valid {
				ar.UpdateTime = new(time.Time)
				*ar.UpdateTime = value.Time
			}
		case applicationresource.FieldApplicationID:
			if value, ok := values[i].(*oid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field applicationID", values[i])
			} else if value != nil {
				ar.ApplicationID = *value
			}
		case applicationresource.FieldModule:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field module", values[i])
			} else if value.Valid {
				ar.Module = value.String
			}
		case applicationresource.FieldType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field type", values[i])
			} else if value.Valid {
				ar.Type = value.String
			}
		}
	}
	return nil
}

// Update returns a builder for updating this ApplicationResource.
// Note that you need to call ApplicationResource.Unwrap() before calling this method if this ApplicationResource
// was returned from a transaction, and the transaction was committed or rolled back.
func (ar *ApplicationResource) Update() *ApplicationResourceUpdateOne {
	return NewApplicationResourceClient(ar.config).UpdateOne(ar)
}

// Unwrap unwraps the ApplicationResource entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (ar *ApplicationResource) Unwrap() *ApplicationResource {
	_tx, ok := ar.config.driver.(*txDriver)
	if !ok {
		panic("model: ApplicationResource is not a transactional entity")
	}
	ar.config.driver = _tx.drv
	return ar
}

// String implements the fmt.Stringer.
func (ar *ApplicationResource) String() string {
	var builder strings.Builder
	builder.WriteString("ApplicationResource(")
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
	builder.WriteString("module=")
	builder.WriteString(ar.Module)
	builder.WriteString(", ")
	builder.WriteString("type=")
	builder.WriteString(ar.Type)
	builder.WriteByte(')')
	return builder.String()
}

// ApplicationResources is a parsable slice of ApplicationResource.
type ApplicationResources []*ApplicationResource

func (ar ApplicationResources) config(cfg config) {
	for _i := range ar {
		ar[_i].config = cfg
	}
}
