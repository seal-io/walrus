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

	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/oid"
	"github.com/seal-io/seal/pkg/dao/schema"
)

// Application is the model entity for the Application schema.
type Application struct {
	config `json:"-"`
	// ID of the ent.
	// ID of the resource.
	ID oid.ID `json:"id,omitempty"`
	// Describe creation time.
	CreateTime *time.Time `json:"createTime,omitempty"`
	// Describe modification time.
	UpdateTime *time.Time `json:"updateTime,omitempty"`
	// ID of the project to which the application belongs
	ProjectID oid.ID `json:"projectID"`
	// ID of the environment to which the application deploys
	EnvironmentID oid.ID `json:"environmentID"`
	// Application modules
	Modules []schema.ApplicationModule `json:"modules"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Application) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case application.FieldModules:
			values[i] = new([]byte)
		case application.FieldID, application.FieldProjectID, application.FieldEnvironmentID:
			values[i] = new(oid.ID)
		case application.FieldCreateTime, application.FieldUpdateTime:
			values[i] = new(sql.NullTime)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Application", columns[i])
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
		case application.FieldEnvironmentID:
			if value, ok := values[i].(*oid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field environmentID", values[i])
			} else if value != nil {
				a.EnvironmentID = *value
			}
		case application.FieldModules:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field modules", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &a.Modules); err != nil {
					return fmt.Errorf("unmarshal field modules: %w", err)
				}
			}
		}
	}
	return nil
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
	builder.WriteString("environmentID=")
	builder.WriteString(fmt.Sprintf("%v", a.EnvironmentID))
	builder.WriteString(", ")
	builder.WriteString("modules=")
	builder.WriteString(fmt.Sprintf("%v", a.Modules))
	builder.WriteByte(')')
	return builder.String()
}

// Applications is a parsable slice of Application.
type Applications []*Application

func (a Applications) config(cfg config) {
	for _i := range a {
		a[_i].config = cfg
	}
}
