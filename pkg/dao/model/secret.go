// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package model

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"

	"github.com/seal-io/seal/pkg/dao/model/project"
	"github.com/seal-io/seal/pkg/dao/model/secret"
	"github.com/seal-io/seal/pkg/dao/types/crypto"
	"github.com/seal-io/seal/pkg/dao/types/oid"
)

// Secret is the model entity for the Secret schema.
type Secret struct {
	config `json:"-"`
	// ID of the ent.
	ID oid.ID `json:"id,omitempty"`
	// Describe creation time.
	CreateTime *time.Time `json:"createTime,omitempty"`
	// Describe modification time.
	UpdateTime *time.Time `json:"updateTime,omitempty"`
	// ID of the project to which the secret belongs, empty means sharing for all projects.
	ProjectID oid.ID `json:"projectID,omitempty"`
	// The name of secret.
	Name string `json:"name,omitempty"`
	// The value of secret, store in string.
	Value crypto.String `json:"-"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the SecretQuery when eager-loading is set.
	Edges SecretEdges `json:"edges,omitempty"`
}

// SecretEdges holds the relations/edges for other nodes in the graph.
type SecretEdges struct {
	// Project to which this secret belongs.
	Project *Project `json:"project,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// ProjectOrErr returns the Project value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e SecretEdges) ProjectOrErr() (*Project, error) {
	if e.loadedTypes[0] {
		if e.Project == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: project.Label}
		}
		return e.Project, nil
	}
	return nil, &NotLoadedError{edge: "project"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Secret) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case secret.FieldValue:
			values[i] = new(crypto.String)
		case secret.FieldID, secret.FieldProjectID:
			values[i] = new(oid.ID)
		case secret.FieldName:
			values[i] = new(sql.NullString)
		case secret.FieldCreateTime, secret.FieldUpdateTime:
			values[i] = new(sql.NullTime)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Secret", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Secret fields.
func (s *Secret) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case secret.FieldID:
			if value, ok := values[i].(*oid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				s.ID = *value
			}
		case secret.FieldCreateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field createTime", values[i])
			} else if value.Valid {
				s.CreateTime = new(time.Time)
				*s.CreateTime = value.Time
			}
		case secret.FieldUpdateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updateTime", values[i])
			} else if value.Valid {
				s.UpdateTime = new(time.Time)
				*s.UpdateTime = value.Time
			}
		case secret.FieldProjectID:
			if value, ok := values[i].(*oid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field projectID", values[i])
			} else if value != nil {
				s.ProjectID = *value
			}
		case secret.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				s.Name = value.String
			}
		case secret.FieldValue:
			if value, ok := values[i].(*crypto.String); !ok {
				return fmt.Errorf("unexpected type %T for field value", values[i])
			} else if value != nil {
				s.Value = *value
			}
		}
	}
	return nil
}

// QueryProject queries the "project" edge of the Secret entity.
func (s *Secret) QueryProject() *ProjectQuery {
	return NewSecretClient(s.config).QueryProject(s)
}

// Update returns a builder for updating this Secret.
// Note that you need to call Secret.Unwrap() before calling this method if this Secret
// was returned from a transaction, and the transaction was committed or rolled back.
func (s *Secret) Update() *SecretUpdateOne {
	return NewSecretClient(s.config).UpdateOne(s)
}

// Unwrap unwraps the Secret entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (s *Secret) Unwrap() *Secret {
	_tx, ok := s.config.driver.(*txDriver)
	if !ok {
		panic("model: Secret is not a transactional entity")
	}
	s.config.driver = _tx.drv
	return s
}

// String implements the fmt.Stringer.
func (s *Secret) String() string {
	var builder strings.Builder
	builder.WriteString("Secret(")
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
	builder.WriteString("projectID=")
	builder.WriteString(fmt.Sprintf("%v", s.ProjectID))
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(s.Name)
	builder.WriteString(", ")
	builder.WriteString("value=<sensitive>")
	builder.WriteByte(')')
	return builder.String()
}

// Secrets is a parsable slice of Secret.
type Secrets []*Secret
