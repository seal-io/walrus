// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package model

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"

	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/types"
)

// ApplicationResource is the model entity for the ApplicationResource schema.
type ApplicationResource struct {
	config `json:"-"`
	// ID of the ent.
	ID types.ID `json:"id,omitempty"`
	// Status of the resource.
	Status string `json:"status,omitempty"`
	// Extra message for status, like error details.
	StatusMessage string `json:"statusMessage,omitempty"`
	// Describe creation time.
	CreateTime *time.Time `json:"createTime,omitempty"`
	// Describe modification time.
	UpdateTime *time.Time `json:"updateTime,omitempty"`
	// ID of the application to which the resource belongs.
	ApplicationID types.ID `json:"applicationID"`
	// ID of the connector to which the resource deploys.
	ConnectorID types.ID `json:"connectorID"`
	// Name of the module that generates the resource.
	Module string `json:"module"`
	// Mode that manages the generated resource, it is the management way of the deployer to the resource, which provides by deployer.
	Mode string `json:"mode"`
	// Type of the generated resource, it is the type of the resource which the deployer observes, which provides by deployer.
	Type string `json:"type"`
	// Name of the generated resource, it is the real identifier of the resource, which provides by deployer.
	Name string `json:"name"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ApplicationResourceQuery when eager-loading is set.
	Edges ApplicationResourceEdges `json:"edges,omitempty"`
}

// ApplicationResourceEdges holds the relations/edges for other nodes in the graph.
type ApplicationResourceEdges struct {
	// Application to which the resource belongs.
	Application *Application `json:"application,omitempty"`
	// Connector to which the resource deploys.
	Connector *Connector `json:"connector,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// ApplicationOrErr returns the Application value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ApplicationResourceEdges) ApplicationOrErr() (*Application, error) {
	if e.loadedTypes[0] {
		if e.Application == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: application.Label}
		}
		return e.Application, nil
	}
	return nil, &NotLoadedError{edge: "application"}
}

// ConnectorOrErr returns the Connector value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ApplicationResourceEdges) ConnectorOrErr() (*Connector, error) {
	if e.loadedTypes[1] {
		if e.Connector == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: connector.Label}
		}
		return e.Connector, nil
	}
	return nil, &NotLoadedError{edge: "connector"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*ApplicationResource) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case applicationresource.FieldStatus, applicationresource.FieldStatusMessage, applicationresource.FieldModule, applicationresource.FieldMode, applicationresource.FieldType, applicationresource.FieldName:
			values[i] = new(sql.NullString)
		case applicationresource.FieldCreateTime, applicationresource.FieldUpdateTime:
			values[i] = new(sql.NullTime)
		case applicationresource.FieldID, applicationresource.FieldApplicationID, applicationresource.FieldConnectorID:
			values[i] = new(types.ID)
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
			if value, ok := values[i].(*types.ID); !ok {
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
			if value, ok := values[i].(*types.ID); !ok {
				return fmt.Errorf("unexpected type %T for field applicationID", values[i])
			} else if value != nil {
				ar.ApplicationID = *value
			}
		case applicationresource.FieldConnectorID:
			if value, ok := values[i].(*types.ID); !ok {
				return fmt.Errorf("unexpected type %T for field connectorID", values[i])
			} else if value != nil {
				ar.ConnectorID = *value
			}
		case applicationresource.FieldModule:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field module", values[i])
			} else if value.Valid {
				ar.Module = value.String
			}
		case applicationresource.FieldMode:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field mode", values[i])
			} else if value.Valid {
				ar.Mode = value.String
			}
		case applicationresource.FieldType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field type", values[i])
			} else if value.Valid {
				ar.Type = value.String
			}
		case applicationresource.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				ar.Name = value.String
			}
		}
	}
	return nil
}

// QueryApplication queries the "application" edge of the ApplicationResource entity.
func (ar *ApplicationResource) QueryApplication() *ApplicationQuery {
	return NewApplicationResourceClient(ar.config).QueryApplication(ar)
}

// QueryConnector queries the "connector" edge of the ApplicationResource entity.
func (ar *ApplicationResource) QueryConnector() *ConnectorQuery {
	return NewApplicationResourceClient(ar.config).QueryConnector(ar)
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
	builder.WriteString("connectorID=")
	builder.WriteString(fmt.Sprintf("%v", ar.ConnectorID))
	builder.WriteString(", ")
	builder.WriteString("module=")
	builder.WriteString(ar.Module)
	builder.WriteString(", ")
	builder.WriteString("mode=")
	builder.WriteString(ar.Mode)
	builder.WriteString(", ")
	builder.WriteString("type=")
	builder.WriteString(ar.Type)
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(ar.Name)
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
