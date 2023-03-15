// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package model

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"

	"github.com/seal-io/seal/pkg/dao/model/setting"
	"github.com/seal-io/seal/pkg/dao/types/oid"
)

// Setting is the model entity for the Setting schema.
type Setting struct {
	config `json:"-"`
	// ID of the ent.
	ID oid.ID `json:"id,omitempty"`
	// Describe creation time.
	CreateTime *time.Time `json:"createTime,omitempty"`
	// Describe modification time.
	UpdateTime *time.Time `json:"updateTime,omitempty"`
	// The name of system setting.
	Name string `json:"name,omitempty"`
	// The value of system setting, store in string.
	Value string `json:"value,omitempty"`
	// Indicate the system setting should be hidden or not, default is visible.
	Hidden *bool `json:"hidden,omitempty"`
	// Indicate the system setting should be edited or not, default is readonly.
	Editable *bool `json:"editable,omitempty"`
	// Indicate the system setting should be exposed or not, default is exposed.
	Private *bool `json:"private,omitempty"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Setting) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case setting.FieldID:
			values[i] = new(oid.ID)
		case setting.FieldHidden, setting.FieldEditable, setting.FieldPrivate:
			values[i] = new(sql.NullBool)
		case setting.FieldName, setting.FieldValue:
			values[i] = new(sql.NullString)
		case setting.FieldCreateTime, setting.FieldUpdateTime:
			values[i] = new(sql.NullTime)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Setting", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Setting fields.
func (s *Setting) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case setting.FieldID:
			if value, ok := values[i].(*oid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				s.ID = *value
			}
		case setting.FieldCreateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field createTime", values[i])
			} else if value.Valid {
				s.CreateTime = new(time.Time)
				*s.CreateTime = value.Time
			}
		case setting.FieldUpdateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updateTime", values[i])
			} else if value.Valid {
				s.UpdateTime = new(time.Time)
				*s.UpdateTime = value.Time
			}
		case setting.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				s.Name = value.String
			}
		case setting.FieldValue:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field value", values[i])
			} else if value.Valid {
				s.Value = value.String
			}
		case setting.FieldHidden:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field hidden", values[i])
			} else if value.Valid {
				s.Hidden = new(bool)
				*s.Hidden = value.Bool
			}
		case setting.FieldEditable:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field editable", values[i])
			} else if value.Valid {
				s.Editable = new(bool)
				*s.Editable = value.Bool
			}
		case setting.FieldPrivate:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field private", values[i])
			} else if value.Valid {
				s.Private = new(bool)
				*s.Private = value.Bool
			}
		}
	}
	return nil
}

// Update returns a builder for updating this Setting.
// Note that you need to call Setting.Unwrap() before calling this method if this Setting
// was returned from a transaction, and the transaction was committed or rolled back.
func (s *Setting) Update() *SettingUpdateOne {
	return NewSettingClient(s.config).UpdateOne(s)
}

// Unwrap unwraps the Setting entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (s *Setting) Unwrap() *Setting {
	_tx, ok := s.config.driver.(*txDriver)
	if !ok {
		panic("model: Setting is not a transactional entity")
	}
	s.config.driver = _tx.drv
	return s
}

// String implements the fmt.Stringer.
func (s *Setting) String() string {
	var builder strings.Builder
	builder.WriteString("Setting(")
	builder.WriteString(fmt.Sprintf("id=%v, ", s.ID))
	if v := s.CreateTime; v != nil {
		builder.WriteString("createTime=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	if v := s.UpdateTime; v != nil {
		builder.WriteString("updateTime=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(s.Name)
	builder.WriteString(", ")
	builder.WriteString("value=")
	builder.WriteString(s.Value)
	builder.WriteString(", ")
	if v := s.Hidden; v != nil {
		builder.WriteString("hidden=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteString(", ")
	if v := s.Editable; v != nil {
		builder.WriteString("editable=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteString(", ")
	if v := s.Private; v != nil {
		builder.WriteString("private=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteByte(')')
	return builder.String()
}

// Settings is a parsable slice of Setting.
type Settings []*Setting
