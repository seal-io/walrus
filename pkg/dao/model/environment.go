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

	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/types"
)

// Environment is the model entity for the Environment schema.
type Environment struct {
	config `json:"-"`
	// ID of the ent.
	ID types.ID `json:"id,omitempty"`
	// Name of the resource.
	Name string `json:"name"`
	// Description of the resource.
	Description string `json:"description,omitempty"`
	// Labels of the resource.
	Labels map[string]string `json:"labels,omitempty"`
	// Describe creation time.
	CreateTime *time.Time `json:"createTime,omitempty"`
	// Describe modification time.
	UpdateTime *time.Time `json:"updateTime,omitempty"`
	// Variables of the environment.
	Variables map[string]interface{} `json:"variables,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the EnvironmentQuery when eager-loading is set.
	Edges EnvironmentEdges `json:"edges,omitempty"`
	// [EXTENSION] Connectors is the collection of the related connectors.
	// It does not store in the database and only uses for creating or updating.
	Connectors []*Connector `json:"connectors,omitempty"`
}

// EnvironmentEdges holds the relations/edges for other nodes in the graph.
type EnvironmentEdges struct {
	// Connectors that configure to the environment.
	Connectors []*Connector `json:"connectors,omitempty"`
	// Applications that belong to the environment.
	Applications []*Application `json:"applications,omitempty"`
	// Revisions that belong to the environment.
	Revisions []*ApplicationRevision `json:"revisions,omitempty"`
	// EnvironmentConnectorRelationships holds the value of the environmentConnectorRelationships edge.
	EnvironmentConnectorRelationships []*EnvironmentConnectorRelationship `json:"environmentConnectorRelationships,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes                            [4]bool
	namedConnectors                        map[string][]*Connector
	namedApplications                      map[string][]*Application
	namedRevisions                         map[string][]*ApplicationRevision
	namedEnvironmentConnectorRelationships map[string][]*EnvironmentConnectorRelationship
}

// ConnectorsOrErr returns the Connectors value or an error if the edge
// was not loaded in eager-loading.
func (e EnvironmentEdges) ConnectorsOrErr() ([]*Connector, error) {
	if e.loadedTypes[0] {
		return e.Connectors, nil
	}
	return nil, &NotLoadedError{edge: "connectors"}
}

// ApplicationsOrErr returns the Applications value or an error if the edge
// was not loaded in eager-loading.
func (e EnvironmentEdges) ApplicationsOrErr() ([]*Application, error) {
	if e.loadedTypes[1] {
		return e.Applications, nil
	}
	return nil, &NotLoadedError{edge: "applications"}
}

// RevisionsOrErr returns the Revisions value or an error if the edge
// was not loaded in eager-loading.
func (e EnvironmentEdges) RevisionsOrErr() ([]*ApplicationRevision, error) {
	if e.loadedTypes[2] {
		return e.Revisions, nil
	}
	return nil, &NotLoadedError{edge: "revisions"}
}

// EnvironmentConnectorRelationshipsOrErr returns the EnvironmentConnectorRelationships value or an error if the edge
// was not loaded in eager-loading.
func (e EnvironmentEdges) EnvironmentConnectorRelationshipsOrErr() ([]*EnvironmentConnectorRelationship, error) {
	if e.loadedTypes[3] {
		return e.EnvironmentConnectorRelationships, nil
	}
	return nil, &NotLoadedError{edge: "environmentConnectorRelationships"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Environment) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case environment.FieldLabels, environment.FieldVariables:
			values[i] = new([]byte)
		case environment.FieldName, environment.FieldDescription:
			values[i] = new(sql.NullString)
		case environment.FieldCreateTime, environment.FieldUpdateTime:
			values[i] = new(sql.NullTime)
		case environment.FieldID:
			values[i] = new(types.ID)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Environment", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Environment fields.
func (e *Environment) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case environment.FieldID:
			if value, ok := values[i].(*types.ID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				e.ID = *value
			}
		case environment.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				e.Name = value.String
			}
		case environment.FieldDescription:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field description", values[i])
			} else if value.Valid {
				e.Description = value.String
			}
		case environment.FieldLabels:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field labels", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &e.Labels); err != nil {
					return fmt.Errorf("unmarshal field labels: %w", err)
				}
			}
		case environment.FieldCreateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field createTime", values[i])
			} else if value.Valid {
				e.CreateTime = new(time.Time)
				*e.CreateTime = value.Time
			}
		case environment.FieldUpdateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updateTime", values[i])
			} else if value.Valid {
				e.UpdateTime = new(time.Time)
				*e.UpdateTime = value.Time
			}
		case environment.FieldVariables:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field variables", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &e.Variables); err != nil {
					return fmt.Errorf("unmarshal field variables: %w", err)
				}
			}
		}
	}
	return nil
}

// QueryConnectors queries the "connectors" edge of the Environment entity.
func (e *Environment) QueryConnectors() *ConnectorQuery {
	return NewEnvironmentClient(e.config).QueryConnectors(e)
}

// QueryApplications queries the "applications" edge of the Environment entity.
func (e *Environment) QueryApplications() *ApplicationQuery {
	return NewEnvironmentClient(e.config).QueryApplications(e)
}

// QueryRevisions queries the "revisions" edge of the Environment entity.
func (e *Environment) QueryRevisions() *ApplicationRevisionQuery {
	return NewEnvironmentClient(e.config).QueryRevisions(e)
}

// QueryEnvironmentConnectorRelationships queries the "environmentConnectorRelationships" edge of the Environment entity.
func (e *Environment) QueryEnvironmentConnectorRelationships() *EnvironmentConnectorRelationshipQuery {
	return NewEnvironmentClient(e.config).QueryEnvironmentConnectorRelationships(e)
}

// Update returns a builder for updating this Environment.
// Note that you need to call Environment.Unwrap() before calling this method if this Environment
// was returned from a transaction, and the transaction was committed or rolled back.
func (e *Environment) Update() *EnvironmentUpdateOne {
	return NewEnvironmentClient(e.config).UpdateOne(e)
}

// Unwrap unwraps the Environment entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (e *Environment) Unwrap() *Environment {
	_tx, ok := e.config.driver.(*txDriver)
	if !ok {
		panic("model: Environment is not a transactional entity")
	}
	e.config.driver = _tx.drv
	return e
}

// String implements the fmt.Stringer.
func (e *Environment) String() string {
	var builder strings.Builder
	builder.WriteString("Environment(")
	builder.WriteString(fmt.Sprintf("id=%v, ", e.ID))
	builder.WriteString("name=")
	builder.WriteString(e.Name)
	builder.WriteString(", ")
	builder.WriteString("description=")
	builder.WriteString(e.Description)
	builder.WriteString(", ")
	builder.WriteString("labels=")
	builder.WriteString(fmt.Sprintf("%v", e.Labels))
	builder.WriteString(", ")
	if v := e.CreateTime; v != nil {
		builder.WriteString("createTime=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	if v := e.UpdateTime; v != nil {
		builder.WriteString("updateTime=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("variables=")
	builder.WriteString(fmt.Sprintf("%v", e.Variables))
	builder.WriteByte(')')
	return builder.String()
}

// MarshalJSON implements the json.Marshaler interface.
func (e *Environment) MarshalJSON() ([]byte, error) {
	type Alias Environment
	// mutate `.Edges.EnvironmentConnectorRelationships` to `.Connectors`.
	if len(e.Edges.EnvironmentConnectorRelationships) != 0 {
		for _, r := range e.Edges.EnvironmentConnectorRelationships {
			if r == nil {
				continue
			}
			e.Connectors = append(e.Connectors,
				&Connector{
					ID: r.ConnectorID,
				})
		}
		e.Edges.EnvironmentConnectorRelationships = nil // release
	}
	// mutate `.Edges.Connectors` to `.Connectors`.
	if len(e.Edges.Connectors) != 0 {
		e.Connectors = e.Edges.Connectors
		e.Edges.Connectors = nil // release
	}
	return json.Marshal(&struct {
		*Alias `json:",inline"`
	}{
		Alias: (*Alias)(e),
	})
}

// NamedConnectors returns the Connectors named value or an error if the edge was not
// loaded in eager-loading with this name.
func (e *Environment) NamedConnectors(name string) ([]*Connector, error) {
	if e.Edges.namedConnectors == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := e.Edges.namedConnectors[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (e *Environment) appendNamedConnectors(name string, edges ...*Connector) {
	if e.Edges.namedConnectors == nil {
		e.Edges.namedConnectors = make(map[string][]*Connector)
	}
	if len(edges) == 0 {
		e.Edges.namedConnectors[name] = []*Connector{}
	} else {
		e.Edges.namedConnectors[name] = append(e.Edges.namedConnectors[name], edges...)
	}
}

// NamedApplications returns the Applications named value or an error if the edge was not
// loaded in eager-loading with this name.
func (e *Environment) NamedApplications(name string) ([]*Application, error) {
	if e.Edges.namedApplications == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := e.Edges.namedApplications[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (e *Environment) appendNamedApplications(name string, edges ...*Application) {
	if e.Edges.namedApplications == nil {
		e.Edges.namedApplications = make(map[string][]*Application)
	}
	if len(edges) == 0 {
		e.Edges.namedApplications[name] = []*Application{}
	} else {
		e.Edges.namedApplications[name] = append(e.Edges.namedApplications[name], edges...)
	}
}

// NamedRevisions returns the Revisions named value or an error if the edge was not
// loaded in eager-loading with this name.
func (e *Environment) NamedRevisions(name string) ([]*ApplicationRevision, error) {
	if e.Edges.namedRevisions == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := e.Edges.namedRevisions[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (e *Environment) appendNamedRevisions(name string, edges ...*ApplicationRevision) {
	if e.Edges.namedRevisions == nil {
		e.Edges.namedRevisions = make(map[string][]*ApplicationRevision)
	}
	if len(edges) == 0 {
		e.Edges.namedRevisions[name] = []*ApplicationRevision{}
	} else {
		e.Edges.namedRevisions[name] = append(e.Edges.namedRevisions[name], edges...)
	}
}

// NamedEnvironmentConnectorRelationships returns the EnvironmentConnectorRelationships named value or an error if the edge was not
// loaded in eager-loading with this name.
func (e *Environment) NamedEnvironmentConnectorRelationships(name string) ([]*EnvironmentConnectorRelationship, error) {
	if e.Edges.namedEnvironmentConnectorRelationships == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := e.Edges.namedEnvironmentConnectorRelationships[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (e *Environment) appendNamedEnvironmentConnectorRelationships(name string, edges ...*EnvironmentConnectorRelationship) {
	if e.Edges.namedEnvironmentConnectorRelationships == nil {
		e.Edges.namedEnvironmentConnectorRelationships = make(map[string][]*EnvironmentConnectorRelationship)
	}
	if len(edges) == 0 {
		e.Edges.namedEnvironmentConnectorRelationships[name] = []*EnvironmentConnectorRelationship{}
	} else {
		e.Edges.namedEnvironmentConnectorRelationships[name] = append(e.Edges.namedEnvironmentConnectorRelationships[name], edges...)
	}
}

// Environments is a parsable slice of Environment.
type Environments []*Environment

func (e Environments) config(cfg config) {
	for _i := range e {
		e[_i].config = cfg
	}
}
