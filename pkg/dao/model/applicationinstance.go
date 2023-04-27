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
	"github.com/seal-io/seal/pkg/dao/model/applicationinstance"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/pkg/dao/types/property"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/utils/json"
)

// ApplicationInstance is the model entity for the ApplicationInstance schema.
type ApplicationInstance struct {
	config `json:"-"`
	// ID of the ent.
	ID oid.ID `json:"id,omitempty" sql:"id"`
	// Describe creation time.
	CreateTime *time.Time `json:"createTime,omitempty" sql:"createTime"`
	// Describe modification time.
	UpdateTime *time.Time `json:"updateTime,omitempty" sql:"updateTime"`
	// ID of the application to which the instance belongs.
	ApplicationID oid.ID `json:"applicationID,omitempty" sql:"applicationID"`
	// ID of the environment to which the instance deploys.
	EnvironmentID oid.ID `json:"environmentID,omitempty" sql:"environmentID"`
	// Name of the instance.
	Name string `json:"name,omitempty" sql:"name"`
	// Variables of the instance.
	Variables property.Values `json:"variables,omitempty" sql:"variables"`
	// Status of the instance.
	Status status.Status `json:"status,omitempty" sql:"status"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ApplicationInstanceQuery when eager-loading is set.
	Edges ApplicationInstanceEdges `json:"edges,omitempty"`
}

// ApplicationInstanceEdges holds the relations/edges for other nodes in the graph.
type ApplicationInstanceEdges struct {
	// Application to which the instance belongs.
	Application *Application `json:"application,omitempty" sql:"application"`
	// Environment to which the instance belongs.
	Environment *Environment `json:"environment,omitempty" sql:"environment"`
	// Application revisions that belong to this instance.
	Revisions []*ApplicationRevision `json:"revisions,omitempty" sql:"revisions"`
	// Application resources that belong to the instance.
	Resources []*ApplicationResource `json:"resources,omitempty" sql:"resources"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [4]bool
}

// ApplicationOrErr returns the Application value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ApplicationInstanceEdges) ApplicationOrErr() (*Application, error) {
	if e.loadedTypes[0] {
		if e.Application == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: application.Label}
		}
		return e.Application, nil
	}
	return nil, &NotLoadedError{edge: "application"}
}

// EnvironmentOrErr returns the Environment value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ApplicationInstanceEdges) EnvironmentOrErr() (*Environment, error) {
	if e.loadedTypes[1] {
		if e.Environment == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: environment.Label}
		}
		return e.Environment, nil
	}
	return nil, &NotLoadedError{edge: "environment"}
}

// RevisionsOrErr returns the Revisions value or an error if the edge
// was not loaded in eager-loading.
func (e ApplicationInstanceEdges) RevisionsOrErr() ([]*ApplicationRevision, error) {
	if e.loadedTypes[2] {
		return e.Revisions, nil
	}
	return nil, &NotLoadedError{edge: "revisions"}
}

// ResourcesOrErr returns the Resources value or an error if the edge
// was not loaded in eager-loading.
func (e ApplicationInstanceEdges) ResourcesOrErr() ([]*ApplicationResource, error) {
	if e.loadedTypes[3] {
		return e.Resources, nil
	}
	return nil, &NotLoadedError{edge: "resources"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*ApplicationInstance) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case applicationinstance.FieldStatus:
			values[i] = new([]byte)
		case applicationinstance.FieldID, applicationinstance.FieldApplicationID, applicationinstance.FieldEnvironmentID:
			values[i] = new(oid.ID)
		case applicationinstance.FieldVariables:
			values[i] = new(property.Values)
		case applicationinstance.FieldName:
			values[i] = new(sql.NullString)
		case applicationinstance.FieldCreateTime, applicationinstance.FieldUpdateTime:
			values[i] = new(sql.NullTime)
		default:
			return nil, fmt.Errorf("unexpected column %q for type ApplicationInstance", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the ApplicationInstance fields.
func (ai *ApplicationInstance) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case applicationinstance.FieldID:
			if value, ok := values[i].(*oid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				ai.ID = *value
			}
		case applicationinstance.FieldCreateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field createTime", values[i])
			} else if value.Valid {
				ai.CreateTime = new(time.Time)
				*ai.CreateTime = value.Time
			}
		case applicationinstance.FieldUpdateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updateTime", values[i])
			} else if value.Valid {
				ai.UpdateTime = new(time.Time)
				*ai.UpdateTime = value.Time
			}
		case applicationinstance.FieldApplicationID:
			if value, ok := values[i].(*oid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field applicationID", values[i])
			} else if value != nil {
				ai.ApplicationID = *value
			}
		case applicationinstance.FieldEnvironmentID:
			if value, ok := values[i].(*oid.ID); !ok {
				return fmt.Errorf("unexpected type %T for field environmentID", values[i])
			} else if value != nil {
				ai.EnvironmentID = *value
			}
		case applicationinstance.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				ai.Name = value.String
			}
		case applicationinstance.FieldVariables:
			if value, ok := values[i].(*property.Values); !ok {
				return fmt.Errorf("unexpected type %T for field variables", values[i])
			} else if value != nil {
				ai.Variables = *value
			}
		case applicationinstance.FieldStatus:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field status", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &ai.Status); err != nil {
					return fmt.Errorf("unmarshal field status: %w", err)
				}
			}
		}
	}
	return nil
}

// QueryApplication queries the "application" edge of the ApplicationInstance entity.
func (ai *ApplicationInstance) QueryApplication() *ApplicationQuery {
	return NewApplicationInstanceClient(ai.config).QueryApplication(ai)
}

// QueryEnvironment queries the "environment" edge of the ApplicationInstance entity.
func (ai *ApplicationInstance) QueryEnvironment() *EnvironmentQuery {
	return NewApplicationInstanceClient(ai.config).QueryEnvironment(ai)
}

// QueryRevisions queries the "revisions" edge of the ApplicationInstance entity.
func (ai *ApplicationInstance) QueryRevisions() *ApplicationRevisionQuery {
	return NewApplicationInstanceClient(ai.config).QueryRevisions(ai)
}

// QueryResources queries the "resources" edge of the ApplicationInstance entity.
func (ai *ApplicationInstance) QueryResources() *ApplicationResourceQuery {
	return NewApplicationInstanceClient(ai.config).QueryResources(ai)
}

// Update returns a builder for updating this ApplicationInstance.
// Note that you need to call ApplicationInstance.Unwrap() before calling this method if this ApplicationInstance
// was returned from a transaction, and the transaction was committed or rolled back.
func (ai *ApplicationInstance) Update() *ApplicationInstanceUpdateOne {
	return NewApplicationInstanceClient(ai.config).UpdateOne(ai)
}

// Unwrap unwraps the ApplicationInstance entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (ai *ApplicationInstance) Unwrap() *ApplicationInstance {
	_tx, ok := ai.config.driver.(*txDriver)
	if !ok {
		panic("model: ApplicationInstance is not a transactional entity")
	}
	ai.config.driver = _tx.drv
	return ai
}

// String implements the fmt.Stringer.
func (ai *ApplicationInstance) String() string {
	var builder strings.Builder
	builder.WriteString("ApplicationInstance(")
	builder.WriteString(fmt.Sprintf("id=%v, ", ai.ID))
	if v := ai.CreateTime; v != nil {
		builder.WriteString("createTime=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	if v := ai.UpdateTime; v != nil {
		builder.WriteString("updateTime=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("applicationID=")
	builder.WriteString(fmt.Sprintf("%v", ai.ApplicationID))
	builder.WriteString(", ")
	builder.WriteString("environmentID=")
	builder.WriteString(fmt.Sprintf("%v", ai.EnvironmentID))
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(ai.Name)
	builder.WriteString(", ")
	builder.WriteString("variables=")
	builder.WriteString(fmt.Sprintf("%v", ai.Variables))
	builder.WriteString(", ")
	builder.WriteString("status=")
	builder.WriteString(fmt.Sprintf("%v", ai.Status))
	builder.WriteByte(')')
	return builder.String()
}

// ApplicationInstances is a parsable slice of ApplicationInstance.
type ApplicationInstances []*ApplicationInstance
