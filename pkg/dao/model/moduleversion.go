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

	"github.com/seal-io/seal/pkg/dao/model/module"
	"github.com/seal-io/seal/pkg/dao/model/moduleversion"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/oid"
)

// ModuleVersion is the model entity for the ModuleVersion schema.
type ModuleVersion struct {
	config `json:"-"`
	// ID of the ent.
	ID oid.ID `json:"id,omitempty" sql:"id"`
	// Describe creation time.
	CreateTime *time.Time `json:"createTime,omitempty" sql:"createTime"`
	// Describe modification time.
	UpdateTime *time.Time `json:"updateTime,omitempty" sql:"updateTime"`
	// ID of the module.
	ModuleID string `json:"moduleID,omitempty" sql:"moduleID"`
	// Module version.
	Version string `json:"version,omitempty" sql:"version"`
	// Module version source.
	Source string `json:"source,omitempty" sql:"source"`
	// Schema of the module.
	Schema *types.ModuleSchema `json:"schema,omitempty" sql:"schema"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ModuleVersionQuery when eager-loading is set.
	Edges ModuleVersionEdges `json:"edges,omitempty"`
}

// ModuleVersionEdges holds the relations/edges for other nodes in the graph.
type ModuleVersionEdges struct {
	// Module holds the value of the module edge.
	Module *Module `json:"module,omitempty" sql:"module"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// ModuleOrErr returns the Module value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ModuleVersionEdges) ModuleOrErr() (*Module, error) {
	if e.loadedTypes[0] {
		if e.Module == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: module.Label}
		}
		return e.Module, nil
	}
	return nil, &NotLoadedError{edge: "module"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*ModuleVersion) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case moduleversion.FieldSchema:
			values[i] = new([]byte)
		case moduleversion.FieldID:
			values[i] = new(oid.ID)
		case moduleversion.FieldModuleID, moduleversion.FieldVersion, moduleversion.FieldSource:
			values[i] = new(sql.NullString)
		case moduleversion.FieldCreateTime, moduleversion.FieldUpdateTime:
			values[i] = new(sql.NullTime)
		default:
			return nil, fmt.Errorf("unexpected column %q for type ModuleVersion", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the ModuleVersion fields.
func (mv *ModuleVersion) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case moduleversion.FieldID:
			if value, ok := values[i].(*oid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				mv.ID = *value
			}
		case moduleversion.FieldCreateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field createTime", values[i])
			} else if value.Valid {
				mv.CreateTime = new(time.Time)
				*mv.CreateTime = value.Time
			}
		case moduleversion.FieldUpdateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updateTime", values[i])
			} else if value.Valid {
				mv.UpdateTime = new(time.Time)
				*mv.UpdateTime = value.Time
			}
		case moduleversion.FieldModuleID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field moduleID", values[i])
			} else if value.Valid {
				mv.ModuleID = value.String
			}
		case moduleversion.FieldVersion:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field version", values[i])
			} else if value.Valid {
				mv.Version = value.String
			}
		case moduleversion.FieldSource:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field source", values[i])
			} else if value.Valid {
				mv.Source = value.String
			}
		case moduleversion.FieldSchema:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field schema", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &mv.Schema); err != nil {
					return fmt.Errorf("unmarshal field schema: %w", err)
				}
			}
		}
	}
	return nil
}

// QueryModule queries the "module" edge of the ModuleVersion entity.
func (mv *ModuleVersion) QueryModule() *ModuleQuery {
	return NewModuleVersionClient(mv.config).QueryModule(mv)
}

// Update returns a builder for updating this ModuleVersion.
// Note that you need to call ModuleVersion.Unwrap() before calling this method if this ModuleVersion
// was returned from a transaction, and the transaction was committed or rolled back.
func (mv *ModuleVersion) Update() *ModuleVersionUpdateOne {
	return NewModuleVersionClient(mv.config).UpdateOne(mv)
}

// Unwrap unwraps the ModuleVersion entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (mv *ModuleVersion) Unwrap() *ModuleVersion {
	_tx, ok := mv.config.driver.(*txDriver)
	if !ok {
		panic("model: ModuleVersion is not a transactional entity")
	}
	mv.config.driver = _tx.drv
	return mv
}

// String implements the fmt.Stringer.
func (mv *ModuleVersion) String() string {
	var builder strings.Builder
	builder.WriteString("ModuleVersion(")
	builder.WriteString(fmt.Sprintf("id=%v, ", mv.ID))
	if v := mv.CreateTime; v != nil {
		builder.WriteString("createTime=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	if v := mv.UpdateTime; v != nil {
		builder.WriteString("updateTime=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("moduleID=")
	builder.WriteString(mv.ModuleID)
	builder.WriteString(", ")
	builder.WriteString("version=")
	builder.WriteString(mv.Version)
	builder.WriteString(", ")
	builder.WriteString("source=")
	builder.WriteString(mv.Source)
	builder.WriteString(", ")
	builder.WriteString("schema=")
	builder.WriteString(fmt.Sprintf("%v", mv.Schema))
	builder.WriteByte(')')
	return builder.String()
}

// ModuleVersions is a parsable slice of ModuleVersion.
type ModuleVersions []*ModuleVersion
