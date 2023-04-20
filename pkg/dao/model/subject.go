// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package model

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"

	"github.com/seal-io/seal/pkg/dao/model/subject"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/utils/json"
)

// Subject is the model entity for the Subject schema.
type Subject struct {
	config `json:"-"`
	// ID of the ent.
	ID oid.ID `json:"id,omitempty" sql:"id"`
	// Describe creation time.
	CreateTime *time.Time `json:"createTime,omitempty" sql:"createTime"`
	// Describe modification time.
	UpdateTime *time.Time `json:"updateTime,omitempty" sql:"updateTime"`
	// The kind of the subject.
	Kind string `json:"kind,omitempty" sql:"kind"`
	// The group of the subject.
	Group string `json:"group,omitempty" sql:"group"`
	// The name of the subject.
	Name string `json:"name,omitempty" sql:"name"`
	// The detail of the subject.
	Description string `json:"description,omitempty" sql:"description"`
	// Indicate whether the user mount to the group.
	MountTo *bool `json:"mountTo,omitempty" sql:"mountTo"`
	// Indicate whether the user login to the group.
	LoginTo *bool `json:"loginTo,omitempty" sql:"loginTo"`
	// The role list of the subject.
	Roles types.SubjectRoles `json:"roles,omitempty" sql:"roles"`
	// The path of the subject from the root group to itself.
	Paths []string `json:"paths,omitempty" sql:"paths"`
	// Indicate whether the subject is builtin.
	Builtin bool `json:"builtin,omitempty" sql:"builtin"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Subject) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case subject.FieldRoles, subject.FieldPaths:
			values[i] = new([]byte)
		case subject.FieldID:
			values[i] = new(oid.ID)
		case subject.FieldMountTo, subject.FieldLoginTo, subject.FieldBuiltin:
			values[i] = new(sql.NullBool)
		case subject.FieldKind, subject.FieldGroup, subject.FieldName, subject.FieldDescription:
			values[i] = new(sql.NullString)
		case subject.FieldCreateTime, subject.FieldUpdateTime:
			values[i] = new(sql.NullTime)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Subject", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Subject fields.
func (s *Subject) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case subject.FieldID:
			if value, ok := values[i].(*oid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				s.ID = *value
			}
		case subject.FieldCreateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field createTime", values[i])
			} else if value.Valid {
				s.CreateTime = new(time.Time)
				*s.CreateTime = value.Time
			}
		case subject.FieldUpdateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updateTime", values[i])
			} else if value.Valid {
				s.UpdateTime = new(time.Time)
				*s.UpdateTime = value.Time
			}
		case subject.FieldKind:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field kind", values[i])
			} else if value.Valid {
				s.Kind = value.String
			}
		case subject.FieldGroup:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field group", values[i])
			} else if value.Valid {
				s.Group = value.String
			}
		case subject.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				s.Name = value.String
			}
		case subject.FieldDescription:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field description", values[i])
			} else if value.Valid {
				s.Description = value.String
			}
		case subject.FieldMountTo:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field mountTo", values[i])
			} else if value.Valid {
				s.MountTo = new(bool)
				*s.MountTo = value.Bool
			}
		case subject.FieldLoginTo:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field loginTo", values[i])
			} else if value.Valid {
				s.LoginTo = new(bool)
				*s.LoginTo = value.Bool
			}
		case subject.FieldRoles:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field roles", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &s.Roles); err != nil {
					return fmt.Errorf("unmarshal field roles: %w", err)
				}
			}
		case subject.FieldPaths:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field paths", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &s.Paths); err != nil {
					return fmt.Errorf("unmarshal field paths: %w", err)
				}
			}
		case subject.FieldBuiltin:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field builtin", values[i])
			} else if value.Valid {
				s.Builtin = value.Bool
			}
		}
	}
	return nil
}

// Update returns a builder for updating this Subject.
// Note that you need to call Subject.Unwrap() before calling this method if this Subject
// was returned from a transaction, and the transaction was committed or rolled back.
func (s *Subject) Update() *SubjectUpdateOne {
	return NewSubjectClient(s.config).UpdateOne(s)
}

// Unwrap unwraps the Subject entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (s *Subject) Unwrap() *Subject {
	_tx, ok := s.config.driver.(*txDriver)
	if !ok {
		panic("model: Subject is not a transactional entity")
	}
	s.config.driver = _tx.drv
	return s
}

// String implements the fmt.Stringer.
func (s *Subject) String() string {
	var builder strings.Builder
	builder.WriteString("Subject(")
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
	builder.WriteString("kind=")
	builder.WriteString(s.Kind)
	builder.WriteString(", ")
	builder.WriteString("group=")
	builder.WriteString(s.Group)
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(s.Name)
	builder.WriteString(", ")
	builder.WriteString("description=")
	builder.WriteString(s.Description)
	builder.WriteString(", ")
	if v := s.MountTo; v != nil {
		builder.WriteString("mountTo=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteString(", ")
	if v := s.LoginTo; v != nil {
		builder.WriteString("loginTo=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteString(", ")
	builder.WriteString("roles=")
	builder.WriteString(fmt.Sprintf("%v", s.Roles))
	builder.WriteString(", ")
	builder.WriteString("paths=")
	builder.WriteString(fmt.Sprintf("%v", s.Paths))
	builder.WriteString(", ")
	builder.WriteString("builtin=")
	builder.WriteString(fmt.Sprintf("%v", s.Builtin))
	builder.WriteByte(')')
	return builder.String()
}

// Subjects is a parsable slice of Subject.
type Subjects []*Subject
