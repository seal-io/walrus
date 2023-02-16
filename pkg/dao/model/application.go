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
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/project"
	"github.com/seal-io/seal/pkg/dao/types"
)

// Application is the model entity for the Application schema.
type Application struct {
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
	// ID of the project to which the application belongs.
	ProjectID types.ID `json:"projectID"`
	// ID of the environment to which the application deploys.
	EnvironmentID types.ID `json:"environmentID"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ApplicationQuery when eager-loading is set.
	Edges ApplicationEdges `json:"edges,omitempty"`
	// [EXTENSION] Modules is the collection of the related modules.
	// It does not store in the database and only uses for creating or updating.
	Modules []types.ApplicationModule `json:"modules,omitempty"`
}

// ApplicationEdges holds the relations/edges for other nodes in the graph.
type ApplicationEdges struct {
	// Project to which this application belongs.
	Project *Project `json:"project,omitempty"`
	// Environment to which the application belongs.
	Environment *Environment `json:"environment,omitempty"`
	// Resources that belong to the application.
	Resources []*ApplicationResource `json:"resources,omitempty"`
	// Revisions that belong to this application.
	Revisions []*ApplicationRevision `json:"revisions,omitempty"`
	// Modules that configure to the application.
	Modules []*Module `json:"modules,omitempty"`
	// ApplicationModuleRelationships holds the value of the applicationModuleRelationships edge.
	ApplicationModuleRelationships []*ApplicationModuleRelationship `json:"applicationModuleRelationships,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes                         [6]bool
	namedResources                      map[string][]*ApplicationResource
	namedRevisions                      map[string][]*ApplicationRevision
	namedModules                        map[string][]*Module
	namedApplicationModuleRelationships map[string][]*ApplicationModuleRelationship
}

// ProjectOrErr returns the Project value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ApplicationEdges) ProjectOrErr() (*Project, error) {
	if e.loadedTypes[0] {
		if e.Project == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: project.Label}
		}
		return e.Project, nil
	}
	return nil, &NotLoadedError{edge: "project"}
}

// EnvironmentOrErr returns the Environment value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ApplicationEdges) EnvironmentOrErr() (*Environment, error) {
	if e.loadedTypes[1] {
		if e.Environment == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: environment.Label}
		}
		return e.Environment, nil
	}
	return nil, &NotLoadedError{edge: "environment"}
}

// ResourcesOrErr returns the Resources value or an error if the edge
// was not loaded in eager-loading.
func (e ApplicationEdges) ResourcesOrErr() ([]*ApplicationResource, error) {
	if e.loadedTypes[2] {
		return e.Resources, nil
	}
	return nil, &NotLoadedError{edge: "resources"}
}

// RevisionsOrErr returns the Revisions value or an error if the edge
// was not loaded in eager-loading.
func (e ApplicationEdges) RevisionsOrErr() ([]*ApplicationRevision, error) {
	if e.loadedTypes[3] {
		return e.Revisions, nil
	}
	return nil, &NotLoadedError{edge: "revisions"}
}

// ModulesOrErr returns the Modules value or an error if the edge
// was not loaded in eager-loading.
func (e ApplicationEdges) ModulesOrErr() ([]*Module, error) {
	if e.loadedTypes[4] {
		return e.Modules, nil
	}
	return nil, &NotLoadedError{edge: "modules"}
}

// ApplicationModuleRelationshipsOrErr returns the ApplicationModuleRelationships value or an error if the edge
// was not loaded in eager-loading.
func (e ApplicationEdges) ApplicationModuleRelationshipsOrErr() ([]*ApplicationModuleRelationship, error) {
	if e.loadedTypes[5] {
		return e.ApplicationModuleRelationships, nil
	}
	return nil, &NotLoadedError{edge: "applicationModuleRelationships"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Application) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case application.FieldLabels:
			values[i] = new([]byte)
		case application.FieldName, application.FieldDescription:
			values[i] = new(sql.NullString)
		case application.FieldCreateTime, application.FieldUpdateTime:
			values[i] = new(sql.NullTime)
		case application.FieldID, application.FieldProjectID, application.FieldEnvironmentID:
			values[i] = new(types.ID)
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
			if value, ok := values[i].(*types.ID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				a.ID = *value
			}
		case application.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				a.Name = value.String
			}
		case application.FieldDescription:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field description", values[i])
			} else if value.Valid {
				a.Description = value.String
			}
		case application.FieldLabels:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field labels", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &a.Labels); err != nil {
					return fmt.Errorf("unmarshal field labels: %w", err)
				}
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
			if value, ok := values[i].(*types.ID); !ok {
				return fmt.Errorf("unexpected type %T for field projectID", values[i])
			} else if value != nil {
				a.ProjectID = *value
			}
		case application.FieldEnvironmentID:
			if value, ok := values[i].(*types.ID); !ok {
				return fmt.Errorf("unexpected type %T for field environmentID", values[i])
			} else if value != nil {
				a.EnvironmentID = *value
			}
		}
	}
	return nil
}

// QueryProject queries the "project" edge of the Application entity.
func (a *Application) QueryProject() *ProjectQuery {
	return NewApplicationClient(a.config).QueryProject(a)
}

// QueryEnvironment queries the "environment" edge of the Application entity.
func (a *Application) QueryEnvironment() *EnvironmentQuery {
	return NewApplicationClient(a.config).QueryEnvironment(a)
}

// QueryResources queries the "resources" edge of the Application entity.
func (a *Application) QueryResources() *ApplicationResourceQuery {
	return NewApplicationClient(a.config).QueryResources(a)
}

// QueryRevisions queries the "revisions" edge of the Application entity.
func (a *Application) QueryRevisions() *ApplicationRevisionQuery {
	return NewApplicationClient(a.config).QueryRevisions(a)
}

// QueryModules queries the "modules" edge of the Application entity.
func (a *Application) QueryModules() *ModuleQuery {
	return NewApplicationClient(a.config).QueryModules(a)
}

// QueryApplicationModuleRelationships queries the "applicationModuleRelationships" edge of the Application entity.
func (a *Application) QueryApplicationModuleRelationships() *ApplicationModuleRelationshipQuery {
	return NewApplicationClient(a.config).QueryApplicationModuleRelationships(a)
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
	builder.WriteString("name=")
	builder.WriteString(a.Name)
	builder.WriteString(", ")
	builder.WriteString("description=")
	builder.WriteString(a.Description)
	builder.WriteString(", ")
	builder.WriteString("labels=")
	builder.WriteString(fmt.Sprintf("%v", a.Labels))
	builder.WriteString(", ")
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
	builder.WriteByte(')')
	return builder.String()
}

// MarshalJSON implements the json.Marshaler interface.
func (a *Application) MarshalJSON() ([]byte, error) {
	type Alias Application
	// mutate `.Edges.ApplicationModuleRelationships` to `.Modules`.
	if len(a.Edges.ApplicationModuleRelationships) != 0 {
		for _, r := range a.Edges.ApplicationModuleRelationships {
			if r == nil {
				continue
			}
			a.Modules = append(a.Modules,
				types.ApplicationModule{
					ModuleID:  r.ModuleID,
					Name:      r.Name,
					Variables: r.Variables,
				})
		}
		a.Edges.ApplicationModuleRelationships = nil // release
	}
	return json.Marshal(&struct {
		*Alias `json:",inline"`
	}{
		Alias: (*Alias)(a),
	})
}

// NamedResources returns the Resources named value or an error if the edge was not
// loaded in eager-loading with this name.
func (a *Application) NamedResources(name string) ([]*ApplicationResource, error) {
	if a.Edges.namedResources == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := a.Edges.namedResources[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (a *Application) appendNamedResources(name string, edges ...*ApplicationResource) {
	if a.Edges.namedResources == nil {
		a.Edges.namedResources = make(map[string][]*ApplicationResource)
	}
	if len(edges) == 0 {
		a.Edges.namedResources[name] = []*ApplicationResource{}
	} else {
		a.Edges.namedResources[name] = append(a.Edges.namedResources[name], edges...)
	}
}

// NamedRevisions returns the Revisions named value or an error if the edge was not
// loaded in eager-loading with this name.
func (a *Application) NamedRevisions(name string) ([]*ApplicationRevision, error) {
	if a.Edges.namedRevisions == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := a.Edges.namedRevisions[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (a *Application) appendNamedRevisions(name string, edges ...*ApplicationRevision) {
	if a.Edges.namedRevisions == nil {
		a.Edges.namedRevisions = make(map[string][]*ApplicationRevision)
	}
	if len(edges) == 0 {
		a.Edges.namedRevisions[name] = []*ApplicationRevision{}
	} else {
		a.Edges.namedRevisions[name] = append(a.Edges.namedRevisions[name], edges...)
	}
}

// NamedModules returns the Modules named value or an error if the edge was not
// loaded in eager-loading with this name.
func (a *Application) NamedModules(name string) ([]*Module, error) {
	if a.Edges.namedModules == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := a.Edges.namedModules[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (a *Application) appendNamedModules(name string, edges ...*Module) {
	if a.Edges.namedModules == nil {
		a.Edges.namedModules = make(map[string][]*Module)
	}
	if len(edges) == 0 {
		a.Edges.namedModules[name] = []*Module{}
	} else {
		a.Edges.namedModules[name] = append(a.Edges.namedModules[name], edges...)
	}
}

// NamedApplicationModuleRelationships returns the ApplicationModuleRelationships named value or an error if the edge was not
// loaded in eager-loading with this name.
func (a *Application) NamedApplicationModuleRelationships(name string) ([]*ApplicationModuleRelationship, error) {
	if a.Edges.namedApplicationModuleRelationships == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := a.Edges.namedApplicationModuleRelationships[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (a *Application) appendNamedApplicationModuleRelationships(name string, edges ...*ApplicationModuleRelationship) {
	if a.Edges.namedApplicationModuleRelationships == nil {
		a.Edges.namedApplicationModuleRelationships = make(map[string][]*ApplicationModuleRelationship)
	}
	if len(edges) == 0 {
		a.Edges.namedApplicationModuleRelationships[name] = []*ApplicationModuleRelationship{}
	} else {
		a.Edges.namedApplicationModuleRelationships[name] = append(a.Edges.namedApplicationModuleRelationships[name], edges...)
	}
}

// Applications is a parsable slice of Application.
type Applications []*Application

func (a Applications) config(cfg config) {
	for _i := range a {
		a[_i].config = cfg
	}
}
