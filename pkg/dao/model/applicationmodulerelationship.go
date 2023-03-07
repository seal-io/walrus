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
	"github.com/seal-io/seal/pkg/dao/model/applicationmodulerelationship"
	"github.com/seal-io/seal/pkg/dao/model/module"
	"github.com/seal-io/seal/pkg/dao/types"
)

// ApplicationModuleRelationship is the model entity for the ApplicationModuleRelationship schema.
type ApplicationModuleRelationship struct {
	config `json:"-"`
	// Describe creation time.
	CreateTime *time.Time `json:"createTime,omitempty"`
	// Describe modification time.
	UpdateTime *time.Time `json:"updateTime,omitempty"`
	// ID of the application to which the relationship connects.
	ApplicationID types.ID `json:"applicationID"`
	// ID of the module to which the relationship connects.
	ModuleID string `json:"moduleID"`
	// Name of the module customized to the application.
	Name string `json:"name,omitempty"`
	// Attributes to configure the module.
	Attributes map[string]interface{} `json:"attributes,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ApplicationModuleRelationshipQuery when eager-loading is set.
	Edges ApplicationModuleRelationshipEdges `json:"edges,omitempty"`
}

// ApplicationModuleRelationshipEdges holds the relations/edges for other nodes in the graph.
type ApplicationModuleRelationshipEdges struct {
	// Applications that connect to the relationship.
	Application *Application `json:"application,omitempty"`
	// Modules that connect to the relationship.
	Module *Module `json:"module,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// ApplicationOrErr returns the Application value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ApplicationModuleRelationshipEdges) ApplicationOrErr() (*Application, error) {
	if e.loadedTypes[0] {
		if e.Application == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: application.Label}
		}
		return e.Application, nil
	}
	return nil, &NotLoadedError{edge: "application"}
}

// ModuleOrErr returns the Module value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ApplicationModuleRelationshipEdges) ModuleOrErr() (*Module, error) {
	if e.loadedTypes[1] {
		if e.Module == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: module.Label}
		}
		return e.Module, nil
	}
	return nil, &NotLoadedError{edge: "module"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*ApplicationModuleRelationship) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case applicationmodulerelationship.FieldAttributes:
			values[i] = new([]byte)
		case applicationmodulerelationship.FieldModuleID, applicationmodulerelationship.FieldName:
			values[i] = new(sql.NullString)
		case applicationmodulerelationship.FieldCreateTime, applicationmodulerelationship.FieldUpdateTime:
			values[i] = new(sql.NullTime)
		case applicationmodulerelationship.FieldApplicationID:
			values[i] = new(types.ID)
		default:
			return nil, fmt.Errorf("unexpected column %q for type ApplicationModuleRelationship", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the ApplicationModuleRelationship fields.
func (amr *ApplicationModuleRelationship) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case applicationmodulerelationship.FieldCreateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field createTime", values[i])
			} else if value.Valid {
				amr.CreateTime = new(time.Time)
				*amr.CreateTime = value.Time
			}
		case applicationmodulerelationship.FieldUpdateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updateTime", values[i])
			} else if value.Valid {
				amr.UpdateTime = new(time.Time)
				*amr.UpdateTime = value.Time
			}
		case applicationmodulerelationship.FieldApplicationID:
			if value, ok := values[i].(*types.ID); !ok {
				return fmt.Errorf("unexpected type %T for field application_id", values[i])
			} else if value != nil {
				amr.ApplicationID = *value
			}
		case applicationmodulerelationship.FieldModuleID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field module_id", values[i])
			} else if value.Valid {
				amr.ModuleID = value.String
			}
		case applicationmodulerelationship.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				amr.Name = value.String
			}
		case applicationmodulerelationship.FieldAttributes:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field attributes", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &amr.Attributes); err != nil {
					return fmt.Errorf("unmarshal field attributes: %w", err)
				}
			}
		}
	}
	return nil
}

// QueryApplication queries the "application" edge of the ApplicationModuleRelationship entity.
func (amr *ApplicationModuleRelationship) QueryApplication() *ApplicationQuery {
	return NewApplicationModuleRelationshipClient(amr.config).QueryApplication(amr)
}

// QueryModule queries the "module" edge of the ApplicationModuleRelationship entity.
func (amr *ApplicationModuleRelationship) QueryModule() *ModuleQuery {
	return NewApplicationModuleRelationshipClient(amr.config).QueryModule(amr)
}

// Update returns a builder for updating this ApplicationModuleRelationship.
// Note that you need to call ApplicationModuleRelationship.Unwrap() before calling this method if this ApplicationModuleRelationship
// was returned from a transaction, and the transaction was committed or rolled back.
func (amr *ApplicationModuleRelationship) Update() *ApplicationModuleRelationshipUpdateOne {
	return NewApplicationModuleRelationshipClient(amr.config).UpdateOne(amr)
}

// Unwrap unwraps the ApplicationModuleRelationship entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (amr *ApplicationModuleRelationship) Unwrap() *ApplicationModuleRelationship {
	_tx, ok := amr.config.driver.(*txDriver)
	if !ok {
		panic("model: ApplicationModuleRelationship is not a transactional entity")
	}
	amr.config.driver = _tx.drv
	return amr
}

// String implements the fmt.Stringer.
func (amr *ApplicationModuleRelationship) String() string {
	var builder strings.Builder
	builder.WriteString("ApplicationModuleRelationship(")
	if v := amr.CreateTime; v != nil {
		builder.WriteString("createTime=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	if v := amr.UpdateTime; v != nil {
		builder.WriteString("updateTime=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("application_id=")
	builder.WriteString(fmt.Sprintf("%v", amr.ApplicationID))
	builder.WriteString(", ")
	builder.WriteString("module_id=")
	builder.WriteString(amr.ModuleID)
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(amr.Name)
	builder.WriteString(", ")
	builder.WriteString("attributes=")
	builder.WriteString(fmt.Sprintf("%v", amr.Attributes))
	builder.WriteByte(')')
	return builder.String()
}

// ApplicationModuleRelationships is a parsable slice of ApplicationModuleRelationship.
type ApplicationModuleRelationships []*ApplicationModuleRelationship
