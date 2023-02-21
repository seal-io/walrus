// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package model

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/model/applicationrevision"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/types"
)

// EnvironmentUpdate is the builder for updating Environment entities.
type EnvironmentUpdate struct {
	config
	hooks     []Hook
	mutation  *EnvironmentMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the EnvironmentUpdate builder.
func (eu *EnvironmentUpdate) Where(ps ...predicate.Environment) *EnvironmentUpdate {
	eu.mutation.Where(ps...)
	return eu
}

// SetName sets the "name" field.
func (eu *EnvironmentUpdate) SetName(s string) *EnvironmentUpdate {
	eu.mutation.SetName(s)
	return eu
}

// SetDescription sets the "description" field.
func (eu *EnvironmentUpdate) SetDescription(s string) *EnvironmentUpdate {
	eu.mutation.SetDescription(s)
	return eu
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (eu *EnvironmentUpdate) SetNillableDescription(s *string) *EnvironmentUpdate {
	if s != nil {
		eu.SetDescription(*s)
	}
	return eu
}

// ClearDescription clears the value of the "description" field.
func (eu *EnvironmentUpdate) ClearDescription() *EnvironmentUpdate {
	eu.mutation.ClearDescription()
	return eu
}

// SetLabels sets the "labels" field.
func (eu *EnvironmentUpdate) SetLabels(m map[string]string) *EnvironmentUpdate {
	eu.mutation.SetLabels(m)
	return eu
}

// ClearLabels clears the value of the "labels" field.
func (eu *EnvironmentUpdate) ClearLabels() *EnvironmentUpdate {
	eu.mutation.ClearLabels()
	return eu
}

// SetUpdateTime sets the "updateTime" field.
func (eu *EnvironmentUpdate) SetUpdateTime(t time.Time) *EnvironmentUpdate {
	eu.mutation.SetUpdateTime(t)
	return eu
}

// SetVariables sets the "variables" field.
func (eu *EnvironmentUpdate) SetVariables(m map[string]interface{}) *EnvironmentUpdate {
	eu.mutation.SetVariables(m)
	return eu
}

// ClearVariables clears the value of the "variables" field.
func (eu *EnvironmentUpdate) ClearVariables() *EnvironmentUpdate {
	eu.mutation.ClearVariables()
	return eu
}

// AddConnectorIDs adds the "connectors" edge to the Connector entity by IDs.
func (eu *EnvironmentUpdate) AddConnectorIDs(ids ...types.ID) *EnvironmentUpdate {
	eu.mutation.AddConnectorIDs(ids...)
	return eu
}

// AddConnectors adds the "connectors" edges to the Connector entity.
func (eu *EnvironmentUpdate) AddConnectors(c ...*Connector) *EnvironmentUpdate {
	ids := make([]types.ID, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return eu.AddConnectorIDs(ids...)
}

// AddApplicationIDs adds the "applications" edge to the Application entity by IDs.
func (eu *EnvironmentUpdate) AddApplicationIDs(ids ...types.ID) *EnvironmentUpdate {
	eu.mutation.AddApplicationIDs(ids...)
	return eu
}

// AddApplications adds the "applications" edges to the Application entity.
func (eu *EnvironmentUpdate) AddApplications(a ...*Application) *EnvironmentUpdate {
	ids := make([]types.ID, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return eu.AddApplicationIDs(ids...)
}

// AddRevisionIDs adds the "revisions" edge to the ApplicationRevision entity by IDs.
func (eu *EnvironmentUpdate) AddRevisionIDs(ids ...types.ID) *EnvironmentUpdate {
	eu.mutation.AddRevisionIDs(ids...)
	return eu
}

// AddRevisions adds the "revisions" edges to the ApplicationRevision entity.
func (eu *EnvironmentUpdate) AddRevisions(a ...*ApplicationRevision) *EnvironmentUpdate {
	ids := make([]types.ID, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return eu.AddRevisionIDs(ids...)
}

// Mutation returns the EnvironmentMutation object of the builder.
func (eu *EnvironmentUpdate) Mutation() *EnvironmentMutation {
	return eu.mutation
}

// ClearConnectors clears all "connectors" edges to the Connector entity.
func (eu *EnvironmentUpdate) ClearConnectors() *EnvironmentUpdate {
	eu.mutation.ClearConnectors()
	return eu
}

// RemoveConnectorIDs removes the "connectors" edge to Connector entities by IDs.
func (eu *EnvironmentUpdate) RemoveConnectorIDs(ids ...types.ID) *EnvironmentUpdate {
	eu.mutation.RemoveConnectorIDs(ids...)
	return eu
}

// RemoveConnectors removes "connectors" edges to Connector entities.
func (eu *EnvironmentUpdate) RemoveConnectors(c ...*Connector) *EnvironmentUpdate {
	ids := make([]types.ID, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return eu.RemoveConnectorIDs(ids...)
}

// ClearApplications clears all "applications" edges to the Application entity.
func (eu *EnvironmentUpdate) ClearApplications() *EnvironmentUpdate {
	eu.mutation.ClearApplications()
	return eu
}

// RemoveApplicationIDs removes the "applications" edge to Application entities by IDs.
func (eu *EnvironmentUpdate) RemoveApplicationIDs(ids ...types.ID) *EnvironmentUpdate {
	eu.mutation.RemoveApplicationIDs(ids...)
	return eu
}

// RemoveApplications removes "applications" edges to Application entities.
func (eu *EnvironmentUpdate) RemoveApplications(a ...*Application) *EnvironmentUpdate {
	ids := make([]types.ID, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return eu.RemoveApplicationIDs(ids...)
}

// ClearRevisions clears all "revisions" edges to the ApplicationRevision entity.
func (eu *EnvironmentUpdate) ClearRevisions() *EnvironmentUpdate {
	eu.mutation.ClearRevisions()
	return eu
}

// RemoveRevisionIDs removes the "revisions" edge to ApplicationRevision entities by IDs.
func (eu *EnvironmentUpdate) RemoveRevisionIDs(ids ...types.ID) *EnvironmentUpdate {
	eu.mutation.RemoveRevisionIDs(ids...)
	return eu
}

// RemoveRevisions removes "revisions" edges to ApplicationRevision entities.
func (eu *EnvironmentUpdate) RemoveRevisions(a ...*ApplicationRevision) *EnvironmentUpdate {
	ids := make([]types.ID, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return eu.RemoveRevisionIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (eu *EnvironmentUpdate) Save(ctx context.Context) (int, error) {
	if err := eu.defaults(); err != nil {
		return 0, err
	}
	return withHooks[int, EnvironmentMutation](ctx, eu.sqlSave, eu.mutation, eu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (eu *EnvironmentUpdate) SaveX(ctx context.Context) int {
	affected, err := eu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (eu *EnvironmentUpdate) Exec(ctx context.Context) error {
	_, err := eu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (eu *EnvironmentUpdate) ExecX(ctx context.Context) {
	if err := eu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (eu *EnvironmentUpdate) defaults() error {
	if _, ok := eu.mutation.UpdateTime(); !ok {
		if environment.UpdateDefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized environment.UpdateDefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := environment.UpdateDefaultUpdateTime()
		eu.mutation.SetUpdateTime(v)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (eu *EnvironmentUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *EnvironmentUpdate {
	eu.modifiers = append(eu.modifiers, modifiers...)
	return eu
}

func (eu *EnvironmentUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   environment.Table,
			Columns: environment.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: environment.FieldID,
			},
		},
	}
	if ps := eu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := eu.mutation.Name(); ok {
		_spec.SetField(environment.FieldName, field.TypeString, value)
	}
	if value, ok := eu.mutation.Description(); ok {
		_spec.SetField(environment.FieldDescription, field.TypeString, value)
	}
	if eu.mutation.DescriptionCleared() {
		_spec.ClearField(environment.FieldDescription, field.TypeString)
	}
	if value, ok := eu.mutation.Labels(); ok {
		_spec.SetField(environment.FieldLabels, field.TypeJSON, value)
	}
	if eu.mutation.LabelsCleared() {
		_spec.ClearField(environment.FieldLabels, field.TypeJSON)
	}
	if value, ok := eu.mutation.UpdateTime(); ok {
		_spec.SetField(environment.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := eu.mutation.Variables(); ok {
		_spec.SetField(environment.FieldVariables, field.TypeJSON, value)
	}
	if eu.mutation.VariablesCleared() {
		_spec.ClearField(environment.FieldVariables, field.TypeJSON)
	}
	if eu.mutation.ConnectorsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   environment.ConnectorsTable,
			Columns: environment.ConnectorsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: connector.FieldID,
				},
			},
		}
		edge.Schema = eu.schemaConfig.EnvironmentConnectorRelationship
		createE := &EnvironmentConnectorRelationshipCreate{config: eu.config, mutation: newEnvironmentConnectorRelationshipMutation(eu.config, OpCreate)}
		createE.defaults()
		_, specE := createE.createSpec()
		edge.Target.Fields = specE.Fields
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := eu.mutation.RemovedConnectorsIDs(); len(nodes) > 0 && !eu.mutation.ConnectorsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   environment.ConnectorsTable,
			Columns: environment.ConnectorsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: connector.FieldID,
				},
			},
		}
		edge.Schema = eu.schemaConfig.EnvironmentConnectorRelationship
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		createE := &EnvironmentConnectorRelationshipCreate{config: eu.config, mutation: newEnvironmentConnectorRelationshipMutation(eu.config, OpCreate)}
		createE.defaults()
		_, specE := createE.createSpec()
		edge.Target.Fields = specE.Fields
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := eu.mutation.ConnectorsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   environment.ConnectorsTable,
			Columns: environment.ConnectorsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: connector.FieldID,
				},
			},
		}
		edge.Schema = eu.schemaConfig.EnvironmentConnectorRelationship
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		createE := &EnvironmentConnectorRelationshipCreate{config: eu.config, mutation: newEnvironmentConnectorRelationshipMutation(eu.config, OpCreate)}
		createE.defaults()
		_, specE := createE.createSpec()
		edge.Target.Fields = specE.Fields
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if eu.mutation.ApplicationsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   environment.ApplicationsTable,
			Columns: []string{environment.ApplicationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: application.FieldID,
				},
			},
		}
		edge.Schema = eu.schemaConfig.Application
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := eu.mutation.RemovedApplicationsIDs(); len(nodes) > 0 && !eu.mutation.ApplicationsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   environment.ApplicationsTable,
			Columns: []string{environment.ApplicationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: application.FieldID,
				},
			},
		}
		edge.Schema = eu.schemaConfig.Application
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := eu.mutation.ApplicationsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   environment.ApplicationsTable,
			Columns: []string{environment.ApplicationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: application.FieldID,
				},
			},
		}
		edge.Schema = eu.schemaConfig.Application
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if eu.mutation.RevisionsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   environment.RevisionsTable,
			Columns: []string{environment.RevisionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: applicationrevision.FieldID,
				},
			},
		}
		edge.Schema = eu.schemaConfig.ApplicationRevision
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := eu.mutation.RemovedRevisionsIDs(); len(nodes) > 0 && !eu.mutation.RevisionsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   environment.RevisionsTable,
			Columns: []string{environment.RevisionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: applicationrevision.FieldID,
				},
			},
		}
		edge.Schema = eu.schemaConfig.ApplicationRevision
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := eu.mutation.RevisionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   environment.RevisionsTable,
			Columns: []string{environment.RevisionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: applicationrevision.FieldID,
				},
			},
		}
		edge.Schema = eu.schemaConfig.ApplicationRevision
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.Node.Schema = eu.schemaConfig.Environment
	ctx = internal.NewSchemaConfigContext(ctx, eu.schemaConfig)
	_spec.AddModifiers(eu.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, eu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{environment.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	eu.mutation.done = true
	return n, nil
}

// EnvironmentUpdateOne is the builder for updating a single Environment entity.
type EnvironmentUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *EnvironmentMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetName sets the "name" field.
func (euo *EnvironmentUpdateOne) SetName(s string) *EnvironmentUpdateOne {
	euo.mutation.SetName(s)
	return euo
}

// SetDescription sets the "description" field.
func (euo *EnvironmentUpdateOne) SetDescription(s string) *EnvironmentUpdateOne {
	euo.mutation.SetDescription(s)
	return euo
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (euo *EnvironmentUpdateOne) SetNillableDescription(s *string) *EnvironmentUpdateOne {
	if s != nil {
		euo.SetDescription(*s)
	}
	return euo
}

// ClearDescription clears the value of the "description" field.
func (euo *EnvironmentUpdateOne) ClearDescription() *EnvironmentUpdateOne {
	euo.mutation.ClearDescription()
	return euo
}

// SetLabels sets the "labels" field.
func (euo *EnvironmentUpdateOne) SetLabels(m map[string]string) *EnvironmentUpdateOne {
	euo.mutation.SetLabels(m)
	return euo
}

// ClearLabels clears the value of the "labels" field.
func (euo *EnvironmentUpdateOne) ClearLabels() *EnvironmentUpdateOne {
	euo.mutation.ClearLabels()
	return euo
}

// SetUpdateTime sets the "updateTime" field.
func (euo *EnvironmentUpdateOne) SetUpdateTime(t time.Time) *EnvironmentUpdateOne {
	euo.mutation.SetUpdateTime(t)
	return euo
}

// SetVariables sets the "variables" field.
func (euo *EnvironmentUpdateOne) SetVariables(m map[string]interface{}) *EnvironmentUpdateOne {
	euo.mutation.SetVariables(m)
	return euo
}

// ClearVariables clears the value of the "variables" field.
func (euo *EnvironmentUpdateOne) ClearVariables() *EnvironmentUpdateOne {
	euo.mutation.ClearVariables()
	return euo
}

// AddConnectorIDs adds the "connectors" edge to the Connector entity by IDs.
func (euo *EnvironmentUpdateOne) AddConnectorIDs(ids ...types.ID) *EnvironmentUpdateOne {
	euo.mutation.AddConnectorIDs(ids...)
	return euo
}

// AddConnectors adds the "connectors" edges to the Connector entity.
func (euo *EnvironmentUpdateOne) AddConnectors(c ...*Connector) *EnvironmentUpdateOne {
	ids := make([]types.ID, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return euo.AddConnectorIDs(ids...)
}

// AddApplicationIDs adds the "applications" edge to the Application entity by IDs.
func (euo *EnvironmentUpdateOne) AddApplicationIDs(ids ...types.ID) *EnvironmentUpdateOne {
	euo.mutation.AddApplicationIDs(ids...)
	return euo
}

// AddApplications adds the "applications" edges to the Application entity.
func (euo *EnvironmentUpdateOne) AddApplications(a ...*Application) *EnvironmentUpdateOne {
	ids := make([]types.ID, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return euo.AddApplicationIDs(ids...)
}

// AddRevisionIDs adds the "revisions" edge to the ApplicationRevision entity by IDs.
func (euo *EnvironmentUpdateOne) AddRevisionIDs(ids ...types.ID) *EnvironmentUpdateOne {
	euo.mutation.AddRevisionIDs(ids...)
	return euo
}

// AddRevisions adds the "revisions" edges to the ApplicationRevision entity.
func (euo *EnvironmentUpdateOne) AddRevisions(a ...*ApplicationRevision) *EnvironmentUpdateOne {
	ids := make([]types.ID, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return euo.AddRevisionIDs(ids...)
}

// Mutation returns the EnvironmentMutation object of the builder.
func (euo *EnvironmentUpdateOne) Mutation() *EnvironmentMutation {
	return euo.mutation
}

// ClearConnectors clears all "connectors" edges to the Connector entity.
func (euo *EnvironmentUpdateOne) ClearConnectors() *EnvironmentUpdateOne {
	euo.mutation.ClearConnectors()
	return euo
}

// RemoveConnectorIDs removes the "connectors" edge to Connector entities by IDs.
func (euo *EnvironmentUpdateOne) RemoveConnectorIDs(ids ...types.ID) *EnvironmentUpdateOne {
	euo.mutation.RemoveConnectorIDs(ids...)
	return euo
}

// RemoveConnectors removes "connectors" edges to Connector entities.
func (euo *EnvironmentUpdateOne) RemoveConnectors(c ...*Connector) *EnvironmentUpdateOne {
	ids := make([]types.ID, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return euo.RemoveConnectorIDs(ids...)
}

// ClearApplications clears all "applications" edges to the Application entity.
func (euo *EnvironmentUpdateOne) ClearApplications() *EnvironmentUpdateOne {
	euo.mutation.ClearApplications()
	return euo
}

// RemoveApplicationIDs removes the "applications" edge to Application entities by IDs.
func (euo *EnvironmentUpdateOne) RemoveApplicationIDs(ids ...types.ID) *EnvironmentUpdateOne {
	euo.mutation.RemoveApplicationIDs(ids...)
	return euo
}

// RemoveApplications removes "applications" edges to Application entities.
func (euo *EnvironmentUpdateOne) RemoveApplications(a ...*Application) *EnvironmentUpdateOne {
	ids := make([]types.ID, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return euo.RemoveApplicationIDs(ids...)
}

// ClearRevisions clears all "revisions" edges to the ApplicationRevision entity.
func (euo *EnvironmentUpdateOne) ClearRevisions() *EnvironmentUpdateOne {
	euo.mutation.ClearRevisions()
	return euo
}

// RemoveRevisionIDs removes the "revisions" edge to ApplicationRevision entities by IDs.
func (euo *EnvironmentUpdateOne) RemoveRevisionIDs(ids ...types.ID) *EnvironmentUpdateOne {
	euo.mutation.RemoveRevisionIDs(ids...)
	return euo
}

// RemoveRevisions removes "revisions" edges to ApplicationRevision entities.
func (euo *EnvironmentUpdateOne) RemoveRevisions(a ...*ApplicationRevision) *EnvironmentUpdateOne {
	ids := make([]types.ID, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return euo.RemoveRevisionIDs(ids...)
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (euo *EnvironmentUpdateOne) Select(field string, fields ...string) *EnvironmentUpdateOne {
	euo.fields = append([]string{field}, fields...)
	return euo
}

// Save executes the query and returns the updated Environment entity.
func (euo *EnvironmentUpdateOne) Save(ctx context.Context) (*Environment, error) {
	if err := euo.defaults(); err != nil {
		return nil, err
	}
	return withHooks[*Environment, EnvironmentMutation](ctx, euo.sqlSave, euo.mutation, euo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (euo *EnvironmentUpdateOne) SaveX(ctx context.Context) *Environment {
	node, err := euo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (euo *EnvironmentUpdateOne) Exec(ctx context.Context) error {
	_, err := euo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (euo *EnvironmentUpdateOne) ExecX(ctx context.Context) {
	if err := euo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (euo *EnvironmentUpdateOne) defaults() error {
	if _, ok := euo.mutation.UpdateTime(); !ok {
		if environment.UpdateDefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized environment.UpdateDefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := environment.UpdateDefaultUpdateTime()
		euo.mutation.SetUpdateTime(v)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (euo *EnvironmentUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *EnvironmentUpdateOne {
	euo.modifiers = append(euo.modifiers, modifiers...)
	return euo
}

func (euo *EnvironmentUpdateOne) sqlSave(ctx context.Context) (_node *Environment, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   environment.Table,
			Columns: environment.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: environment.FieldID,
			},
		},
	}
	id, ok := euo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`model: missing "Environment.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := euo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, environment.FieldID)
		for _, f := range fields {
			if !environment.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("model: invalid field %q for query", f)}
			}
			if f != environment.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := euo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := euo.mutation.Name(); ok {
		_spec.SetField(environment.FieldName, field.TypeString, value)
	}
	if value, ok := euo.mutation.Description(); ok {
		_spec.SetField(environment.FieldDescription, field.TypeString, value)
	}
	if euo.mutation.DescriptionCleared() {
		_spec.ClearField(environment.FieldDescription, field.TypeString)
	}
	if value, ok := euo.mutation.Labels(); ok {
		_spec.SetField(environment.FieldLabels, field.TypeJSON, value)
	}
	if euo.mutation.LabelsCleared() {
		_spec.ClearField(environment.FieldLabels, field.TypeJSON)
	}
	if value, ok := euo.mutation.UpdateTime(); ok {
		_spec.SetField(environment.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := euo.mutation.Variables(); ok {
		_spec.SetField(environment.FieldVariables, field.TypeJSON, value)
	}
	if euo.mutation.VariablesCleared() {
		_spec.ClearField(environment.FieldVariables, field.TypeJSON)
	}
	if euo.mutation.ConnectorsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   environment.ConnectorsTable,
			Columns: environment.ConnectorsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: connector.FieldID,
				},
			},
		}
		edge.Schema = euo.schemaConfig.EnvironmentConnectorRelationship
		createE := &EnvironmentConnectorRelationshipCreate{config: euo.config, mutation: newEnvironmentConnectorRelationshipMutation(euo.config, OpCreate)}
		createE.defaults()
		_, specE := createE.createSpec()
		edge.Target.Fields = specE.Fields
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := euo.mutation.RemovedConnectorsIDs(); len(nodes) > 0 && !euo.mutation.ConnectorsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   environment.ConnectorsTable,
			Columns: environment.ConnectorsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: connector.FieldID,
				},
			},
		}
		edge.Schema = euo.schemaConfig.EnvironmentConnectorRelationship
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		createE := &EnvironmentConnectorRelationshipCreate{config: euo.config, mutation: newEnvironmentConnectorRelationshipMutation(euo.config, OpCreate)}
		createE.defaults()
		_, specE := createE.createSpec()
		edge.Target.Fields = specE.Fields
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := euo.mutation.ConnectorsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   environment.ConnectorsTable,
			Columns: environment.ConnectorsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: connector.FieldID,
				},
			},
		}
		edge.Schema = euo.schemaConfig.EnvironmentConnectorRelationship
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		createE := &EnvironmentConnectorRelationshipCreate{config: euo.config, mutation: newEnvironmentConnectorRelationshipMutation(euo.config, OpCreate)}
		createE.defaults()
		_, specE := createE.createSpec()
		edge.Target.Fields = specE.Fields
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if euo.mutation.ApplicationsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   environment.ApplicationsTable,
			Columns: []string{environment.ApplicationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: application.FieldID,
				},
			},
		}
		edge.Schema = euo.schemaConfig.Application
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := euo.mutation.RemovedApplicationsIDs(); len(nodes) > 0 && !euo.mutation.ApplicationsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   environment.ApplicationsTable,
			Columns: []string{environment.ApplicationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: application.FieldID,
				},
			},
		}
		edge.Schema = euo.schemaConfig.Application
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := euo.mutation.ApplicationsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   environment.ApplicationsTable,
			Columns: []string{environment.ApplicationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: application.FieldID,
				},
			},
		}
		edge.Schema = euo.schemaConfig.Application
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if euo.mutation.RevisionsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   environment.RevisionsTable,
			Columns: []string{environment.RevisionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: applicationrevision.FieldID,
				},
			},
		}
		edge.Schema = euo.schemaConfig.ApplicationRevision
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := euo.mutation.RemovedRevisionsIDs(); len(nodes) > 0 && !euo.mutation.RevisionsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   environment.RevisionsTable,
			Columns: []string{environment.RevisionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: applicationrevision.FieldID,
				},
			},
		}
		edge.Schema = euo.schemaConfig.ApplicationRevision
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := euo.mutation.RevisionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   environment.RevisionsTable,
			Columns: []string{environment.RevisionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: applicationrevision.FieldID,
				},
			},
		}
		edge.Schema = euo.schemaConfig.ApplicationRevision
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.Node.Schema = euo.schemaConfig.Environment
	ctx = internal.NewSchemaConfigContext(ctx, euo.schemaConfig)
	_spec.AddModifiers(euo.modifiers...)
	_node = &Environment{config: euo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, euo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{environment.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	euo.mutation.done = true
	return _node, nil
}
