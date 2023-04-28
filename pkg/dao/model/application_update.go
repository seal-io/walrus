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
	"github.com/seal-io/seal/pkg/dao/model/applicationinstance"
	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/pkg/dao/types/property"
)

// ApplicationUpdate is the builder for updating Application entities.
type ApplicationUpdate struct {
	config
	hooks     []Hook
	mutation  *ApplicationMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the ApplicationUpdate builder.
func (au *ApplicationUpdate) Where(ps ...predicate.Application) *ApplicationUpdate {
	au.mutation.Where(ps...)
	return au
}

// SetName sets the "name" field.
func (au *ApplicationUpdate) SetName(s string) *ApplicationUpdate {
	au.mutation.SetName(s)
	return au
}

// SetDescription sets the "description" field.
func (au *ApplicationUpdate) SetDescription(s string) *ApplicationUpdate {
	au.mutation.SetDescription(s)
	return au
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (au *ApplicationUpdate) SetNillableDescription(s *string) *ApplicationUpdate {
	if s != nil {
		au.SetDescription(*s)
	}
	return au
}

// ClearDescription clears the value of the "description" field.
func (au *ApplicationUpdate) ClearDescription() *ApplicationUpdate {
	au.mutation.ClearDescription()
	return au
}

// SetLabels sets the "labels" field.
func (au *ApplicationUpdate) SetLabels(m map[string]string) *ApplicationUpdate {
	au.mutation.SetLabels(m)
	return au
}

// SetUpdateTime sets the "updateTime" field.
func (au *ApplicationUpdate) SetUpdateTime(t time.Time) *ApplicationUpdate {
	au.mutation.SetUpdateTime(t)
	return au
}

// SetVariables sets the "variables" field.
func (au *ApplicationUpdate) SetVariables(pr property.Schemas) *ApplicationUpdate {
	au.mutation.SetVariables(pr)
	return au
}

// ClearVariables clears the value of the "variables" field.
func (au *ApplicationUpdate) ClearVariables() *ApplicationUpdate {
	au.mutation.ClearVariables()
	return au
}

// AddInstanceIDs adds the "instances" edge to the ApplicationInstance entity by IDs.
func (au *ApplicationUpdate) AddInstanceIDs(ids ...oid.ID) *ApplicationUpdate {
	au.mutation.AddInstanceIDs(ids...)
	return au
}

// AddInstances adds the "instances" edges to the ApplicationInstance entity.
func (au *ApplicationUpdate) AddInstances(a ...*ApplicationInstance) *ApplicationUpdate {
	ids := make([]oid.ID, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return au.AddInstanceIDs(ids...)
}

// Mutation returns the ApplicationMutation object of the builder.
func (au *ApplicationUpdate) Mutation() *ApplicationMutation {
	return au.mutation
}

// ClearInstances clears all "instances" edges to the ApplicationInstance entity.
func (au *ApplicationUpdate) ClearInstances() *ApplicationUpdate {
	au.mutation.ClearInstances()
	return au
}

// RemoveInstanceIDs removes the "instances" edge to ApplicationInstance entities by IDs.
func (au *ApplicationUpdate) RemoveInstanceIDs(ids ...oid.ID) *ApplicationUpdate {
	au.mutation.RemoveInstanceIDs(ids...)
	return au
}

// RemoveInstances removes "instances" edges to ApplicationInstance entities.
func (au *ApplicationUpdate) RemoveInstances(a ...*ApplicationInstance) *ApplicationUpdate {
	ids := make([]oid.ID, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return au.RemoveInstanceIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (au *ApplicationUpdate) Save(ctx context.Context) (int, error) {
	if err := au.defaults(); err != nil {
		return 0, err
	}
	return withHooks[int, ApplicationMutation](ctx, au.sqlSave, au.mutation, au.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (au *ApplicationUpdate) SaveX(ctx context.Context) int {
	affected, err := au.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (au *ApplicationUpdate) Exec(ctx context.Context) error {
	_, err := au.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (au *ApplicationUpdate) ExecX(ctx context.Context) {
	if err := au.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (au *ApplicationUpdate) defaults() error {
	if _, ok := au.mutation.UpdateTime(); !ok {
		if application.UpdateDefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized application.UpdateDefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := application.UpdateDefaultUpdateTime()
		au.mutation.SetUpdateTime(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (au *ApplicationUpdate) check() error {
	if v, ok := au.mutation.Name(); ok {
		if err := application.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`model: validator failed for field "Application.name": %w`, err)}
		}
	}
	if _, ok := au.mutation.ProjectID(); au.mutation.ProjectCleared() && !ok {
		return errors.New(`model: clearing a required unique edge "Application.project"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (au *ApplicationUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ApplicationUpdate {
	au.modifiers = append(au.modifiers, modifiers...)
	return au
}

func (au *ApplicationUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := au.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(application.Table, application.Columns, sqlgraph.NewFieldSpec(application.FieldID, field.TypeString))
	if ps := au.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := au.mutation.Name(); ok {
		_spec.SetField(application.FieldName, field.TypeString, value)
	}
	if value, ok := au.mutation.Description(); ok {
		_spec.SetField(application.FieldDescription, field.TypeString, value)
	}
	if au.mutation.DescriptionCleared() {
		_spec.ClearField(application.FieldDescription, field.TypeString)
	}
	if value, ok := au.mutation.Labels(); ok {
		_spec.SetField(application.FieldLabels, field.TypeJSON, value)
	}
	if value, ok := au.mutation.UpdateTime(); ok {
		_spec.SetField(application.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := au.mutation.Variables(); ok {
		_spec.SetField(application.FieldVariables, field.TypeOther, value)
	}
	if au.mutation.VariablesCleared() {
		_spec.ClearField(application.FieldVariables, field.TypeOther)
	}
	if au.mutation.InstancesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   application.InstancesTable,
			Columns: []string{application.InstancesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(applicationinstance.FieldID, field.TypeString),
			},
		}
		edge.Schema = au.schemaConfig.ApplicationInstance
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := au.mutation.RemovedInstancesIDs(); len(nodes) > 0 && !au.mutation.InstancesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   application.InstancesTable,
			Columns: []string{application.InstancesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(applicationinstance.FieldID, field.TypeString),
			},
		}
		edge.Schema = au.schemaConfig.ApplicationInstance
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := au.mutation.InstancesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   application.InstancesTable,
			Columns: []string{application.InstancesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(applicationinstance.FieldID, field.TypeString),
			},
		}
		edge.Schema = au.schemaConfig.ApplicationInstance
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.Node.Schema = au.schemaConfig.Application
	ctx = internal.NewSchemaConfigContext(ctx, au.schemaConfig)
	_spec.AddModifiers(au.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, au.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{application.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	au.mutation.done = true
	return n, nil
}

// ApplicationUpdateOne is the builder for updating a single Application entity.
type ApplicationUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *ApplicationMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetName sets the "name" field.
func (auo *ApplicationUpdateOne) SetName(s string) *ApplicationUpdateOne {
	auo.mutation.SetName(s)
	return auo
}

// SetDescription sets the "description" field.
func (auo *ApplicationUpdateOne) SetDescription(s string) *ApplicationUpdateOne {
	auo.mutation.SetDescription(s)
	return auo
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (auo *ApplicationUpdateOne) SetNillableDescription(s *string) *ApplicationUpdateOne {
	if s != nil {
		auo.SetDescription(*s)
	}
	return auo
}

// ClearDescription clears the value of the "description" field.
func (auo *ApplicationUpdateOne) ClearDescription() *ApplicationUpdateOne {
	auo.mutation.ClearDescription()
	return auo
}

// SetLabels sets the "labels" field.
func (auo *ApplicationUpdateOne) SetLabels(m map[string]string) *ApplicationUpdateOne {
	auo.mutation.SetLabels(m)
	return auo
}

// SetUpdateTime sets the "updateTime" field.
func (auo *ApplicationUpdateOne) SetUpdateTime(t time.Time) *ApplicationUpdateOne {
	auo.mutation.SetUpdateTime(t)
	return auo
}

// SetVariables sets the "variables" field.
func (auo *ApplicationUpdateOne) SetVariables(pr property.Schemas) *ApplicationUpdateOne {
	auo.mutation.SetVariables(pr)
	return auo
}

// ClearVariables clears the value of the "variables" field.
func (auo *ApplicationUpdateOne) ClearVariables() *ApplicationUpdateOne {
	auo.mutation.ClearVariables()
	return auo
}

// AddInstanceIDs adds the "instances" edge to the ApplicationInstance entity by IDs.
func (auo *ApplicationUpdateOne) AddInstanceIDs(ids ...oid.ID) *ApplicationUpdateOne {
	auo.mutation.AddInstanceIDs(ids...)
	return auo
}

// AddInstances adds the "instances" edges to the ApplicationInstance entity.
func (auo *ApplicationUpdateOne) AddInstances(a ...*ApplicationInstance) *ApplicationUpdateOne {
	ids := make([]oid.ID, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return auo.AddInstanceIDs(ids...)
}

// Mutation returns the ApplicationMutation object of the builder.
func (auo *ApplicationUpdateOne) Mutation() *ApplicationMutation {
	return auo.mutation
}

// ClearInstances clears all "instances" edges to the ApplicationInstance entity.
func (auo *ApplicationUpdateOne) ClearInstances() *ApplicationUpdateOne {
	auo.mutation.ClearInstances()
	return auo
}

// RemoveInstanceIDs removes the "instances" edge to ApplicationInstance entities by IDs.
func (auo *ApplicationUpdateOne) RemoveInstanceIDs(ids ...oid.ID) *ApplicationUpdateOne {
	auo.mutation.RemoveInstanceIDs(ids...)
	return auo
}

// RemoveInstances removes "instances" edges to ApplicationInstance entities.
func (auo *ApplicationUpdateOne) RemoveInstances(a ...*ApplicationInstance) *ApplicationUpdateOne {
	ids := make([]oid.ID, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return auo.RemoveInstanceIDs(ids...)
}

// Where appends a list predicates to the ApplicationUpdate builder.
func (auo *ApplicationUpdateOne) Where(ps ...predicate.Application) *ApplicationUpdateOne {
	auo.mutation.Where(ps...)
	return auo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (auo *ApplicationUpdateOne) Select(field string, fields ...string) *ApplicationUpdateOne {
	auo.fields = append([]string{field}, fields...)
	return auo
}

// Save executes the query and returns the updated Application entity.
func (auo *ApplicationUpdateOne) Save(ctx context.Context) (*Application, error) {
	if err := auo.defaults(); err != nil {
		return nil, err
	}
	return withHooks[*Application, ApplicationMutation](ctx, auo.sqlSave, auo.mutation, auo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (auo *ApplicationUpdateOne) SaveX(ctx context.Context) *Application {
	node, err := auo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (auo *ApplicationUpdateOne) Exec(ctx context.Context) error {
	_, err := auo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (auo *ApplicationUpdateOne) ExecX(ctx context.Context) {
	if err := auo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (auo *ApplicationUpdateOne) defaults() error {
	if _, ok := auo.mutation.UpdateTime(); !ok {
		if application.UpdateDefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized application.UpdateDefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := application.UpdateDefaultUpdateTime()
		auo.mutation.SetUpdateTime(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (auo *ApplicationUpdateOne) check() error {
	if v, ok := auo.mutation.Name(); ok {
		if err := application.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`model: validator failed for field "Application.name": %w`, err)}
		}
	}
	if _, ok := auo.mutation.ProjectID(); auo.mutation.ProjectCleared() && !ok {
		return errors.New(`model: clearing a required unique edge "Application.project"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (auo *ApplicationUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ApplicationUpdateOne {
	auo.modifiers = append(auo.modifiers, modifiers...)
	return auo
}

func (auo *ApplicationUpdateOne) sqlSave(ctx context.Context) (_node *Application, err error) {
	if err := auo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(application.Table, application.Columns, sqlgraph.NewFieldSpec(application.FieldID, field.TypeString))
	id, ok := auo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`model: missing "Application.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := auo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, application.FieldID)
		for _, f := range fields {
			if !application.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("model: invalid field %q for query", f)}
			}
			if f != application.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := auo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := auo.mutation.Name(); ok {
		_spec.SetField(application.FieldName, field.TypeString, value)
	}
	if value, ok := auo.mutation.Description(); ok {
		_spec.SetField(application.FieldDescription, field.TypeString, value)
	}
	if auo.mutation.DescriptionCleared() {
		_spec.ClearField(application.FieldDescription, field.TypeString)
	}
	if value, ok := auo.mutation.Labels(); ok {
		_spec.SetField(application.FieldLabels, field.TypeJSON, value)
	}
	if value, ok := auo.mutation.UpdateTime(); ok {
		_spec.SetField(application.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := auo.mutation.Variables(); ok {
		_spec.SetField(application.FieldVariables, field.TypeOther, value)
	}
	if auo.mutation.VariablesCleared() {
		_spec.ClearField(application.FieldVariables, field.TypeOther)
	}
	if auo.mutation.InstancesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   application.InstancesTable,
			Columns: []string{application.InstancesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(applicationinstance.FieldID, field.TypeString),
			},
		}
		edge.Schema = auo.schemaConfig.ApplicationInstance
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := auo.mutation.RemovedInstancesIDs(); len(nodes) > 0 && !auo.mutation.InstancesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   application.InstancesTable,
			Columns: []string{application.InstancesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(applicationinstance.FieldID, field.TypeString),
			},
		}
		edge.Schema = auo.schemaConfig.ApplicationInstance
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := auo.mutation.InstancesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   application.InstancesTable,
			Columns: []string{application.InstancesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(applicationinstance.FieldID, field.TypeString),
			},
		}
		edge.Schema = auo.schemaConfig.ApplicationInstance
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.Node.Schema = auo.schemaConfig.Application
	ctx = internal.NewSchemaConfigContext(ctx, auo.schemaConfig)
	_spec.AddModifiers(auo.modifiers...)
	_node = &Application{config: auo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, auo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{application.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	auo.mutation.done = true
	return _node, nil
}
