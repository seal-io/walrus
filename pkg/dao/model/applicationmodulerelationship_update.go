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

	"github.com/seal-io/seal/pkg/dao/model/applicationmodulerelationship"
	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
)

// ApplicationModuleRelationshipUpdate is the builder for updating ApplicationModuleRelationship entities.
type ApplicationModuleRelationshipUpdate struct {
	config
	hooks     []Hook
	mutation  *ApplicationModuleRelationshipMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the ApplicationModuleRelationshipUpdate builder.
func (amru *ApplicationModuleRelationshipUpdate) Where(ps ...predicate.ApplicationModuleRelationship) *ApplicationModuleRelationshipUpdate {
	amru.mutation.Where(ps...)
	return amru
}

// SetUpdateTime sets the "updateTime" field.
func (amru *ApplicationModuleRelationshipUpdate) SetUpdateTime(t time.Time) *ApplicationModuleRelationshipUpdate {
	amru.mutation.SetUpdateTime(t)
	return amru
}

// SetVariables sets the "variables" field.
func (amru *ApplicationModuleRelationshipUpdate) SetVariables(m map[string]interface{}) *ApplicationModuleRelationshipUpdate {
	amru.mutation.SetVariables(m)
	return amru
}

// ClearVariables clears the value of the "variables" field.
func (amru *ApplicationModuleRelationshipUpdate) ClearVariables() *ApplicationModuleRelationshipUpdate {
	amru.mutation.ClearVariables()
	return amru
}

// Mutation returns the ApplicationModuleRelationshipMutation object of the builder.
func (amru *ApplicationModuleRelationshipUpdate) Mutation() *ApplicationModuleRelationshipMutation {
	return amru.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (amru *ApplicationModuleRelationshipUpdate) Save(ctx context.Context) (int, error) {
	amru.defaults()
	return withHooks[int, ApplicationModuleRelationshipMutation](ctx, amru.sqlSave, amru.mutation, amru.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (amru *ApplicationModuleRelationshipUpdate) SaveX(ctx context.Context) int {
	affected, err := amru.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (amru *ApplicationModuleRelationshipUpdate) Exec(ctx context.Context) error {
	_, err := amru.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (amru *ApplicationModuleRelationshipUpdate) ExecX(ctx context.Context) {
	if err := amru.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (amru *ApplicationModuleRelationshipUpdate) defaults() {
	if _, ok := amru.mutation.UpdateTime(); !ok {
		v := applicationmodulerelationship.UpdateDefaultUpdateTime()
		amru.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (amru *ApplicationModuleRelationshipUpdate) check() error {
	if _, ok := amru.mutation.ApplicationID(); amru.mutation.ApplicationCleared() && !ok {
		return errors.New(`model: clearing a required unique edge "ApplicationModuleRelationship.application"`)
	}
	if _, ok := amru.mutation.ModuleID(); amru.mutation.ModuleCleared() && !ok {
		return errors.New(`model: clearing a required unique edge "ApplicationModuleRelationship.module"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (amru *ApplicationModuleRelationshipUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ApplicationModuleRelationshipUpdate {
	amru.modifiers = append(amru.modifiers, modifiers...)
	return amru
}

func (amru *ApplicationModuleRelationshipUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := amru.check(); err != nil {
		return n, err
	}
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   applicationmodulerelationship.Table,
			Columns: applicationmodulerelationship.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: applicationmodulerelationship.FieldID,
			},
		},
	}
	if ps := amru.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := amru.mutation.UpdateTime(); ok {
		_spec.SetField(applicationmodulerelationship.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := amru.mutation.Variables(); ok {
		_spec.SetField(applicationmodulerelationship.FieldVariables, field.TypeJSON, value)
	}
	if amru.mutation.VariablesCleared() {
		_spec.ClearField(applicationmodulerelationship.FieldVariables, field.TypeJSON)
	}
	_spec.Node.Schema = amru.schemaConfig.ApplicationModuleRelationship
	ctx = internal.NewSchemaConfigContext(ctx, amru.schemaConfig)
	_spec.AddModifiers(amru.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, amru.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{applicationmodulerelationship.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	amru.mutation.done = true
	return n, nil
}

// ApplicationModuleRelationshipUpdateOne is the builder for updating a single ApplicationModuleRelationship entity.
type ApplicationModuleRelationshipUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *ApplicationModuleRelationshipMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetUpdateTime sets the "updateTime" field.
func (amruo *ApplicationModuleRelationshipUpdateOne) SetUpdateTime(t time.Time) *ApplicationModuleRelationshipUpdateOne {
	amruo.mutation.SetUpdateTime(t)
	return amruo
}

// SetVariables sets the "variables" field.
func (amruo *ApplicationModuleRelationshipUpdateOne) SetVariables(m map[string]interface{}) *ApplicationModuleRelationshipUpdateOne {
	amruo.mutation.SetVariables(m)
	return amruo
}

// ClearVariables clears the value of the "variables" field.
func (amruo *ApplicationModuleRelationshipUpdateOne) ClearVariables() *ApplicationModuleRelationshipUpdateOne {
	amruo.mutation.ClearVariables()
	return amruo
}

// Mutation returns the ApplicationModuleRelationshipMutation object of the builder.
func (amruo *ApplicationModuleRelationshipUpdateOne) Mutation() *ApplicationModuleRelationshipMutation {
	return amruo.mutation
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (amruo *ApplicationModuleRelationshipUpdateOne) Select(field string, fields ...string) *ApplicationModuleRelationshipUpdateOne {
	amruo.fields = append([]string{field}, fields...)
	return amruo
}

// Save executes the query and returns the updated ApplicationModuleRelationship entity.
func (amruo *ApplicationModuleRelationshipUpdateOne) Save(ctx context.Context) (*ApplicationModuleRelationship, error) {
	amruo.defaults()
	return withHooks[*ApplicationModuleRelationship, ApplicationModuleRelationshipMutation](ctx, amruo.sqlSave, amruo.mutation, amruo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (amruo *ApplicationModuleRelationshipUpdateOne) SaveX(ctx context.Context) *ApplicationModuleRelationship {
	node, err := amruo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (amruo *ApplicationModuleRelationshipUpdateOne) Exec(ctx context.Context) error {
	_, err := amruo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (amruo *ApplicationModuleRelationshipUpdateOne) ExecX(ctx context.Context) {
	if err := amruo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (amruo *ApplicationModuleRelationshipUpdateOne) defaults() {
	if _, ok := amruo.mutation.UpdateTime(); !ok {
		v := applicationmodulerelationship.UpdateDefaultUpdateTime()
		amruo.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (amruo *ApplicationModuleRelationshipUpdateOne) check() error {
	if _, ok := amruo.mutation.ApplicationID(); amruo.mutation.ApplicationCleared() && !ok {
		return errors.New(`model: clearing a required unique edge "ApplicationModuleRelationship.application"`)
	}
	if _, ok := amruo.mutation.ModuleID(); amruo.mutation.ModuleCleared() && !ok {
		return errors.New(`model: clearing a required unique edge "ApplicationModuleRelationship.module"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (amruo *ApplicationModuleRelationshipUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ApplicationModuleRelationshipUpdateOne {
	amruo.modifiers = append(amruo.modifiers, modifiers...)
	return amruo
}

func (amruo *ApplicationModuleRelationshipUpdateOne) sqlSave(ctx context.Context) (_node *ApplicationModuleRelationship, err error) {
	if err := amruo.check(); err != nil {
		return _node, err
	}
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   applicationmodulerelationship.Table,
			Columns: applicationmodulerelationship.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: applicationmodulerelationship.FieldID,
			},
		},
	}
	id, ok := amruo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`model: missing "ApplicationModuleRelationship.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := amruo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, applicationmodulerelationship.FieldID)
		for _, f := range fields {
			if !applicationmodulerelationship.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("model: invalid field %q for query", f)}
			}
			if f != applicationmodulerelationship.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := amruo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := amruo.mutation.UpdateTime(); ok {
		_spec.SetField(applicationmodulerelationship.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := amruo.mutation.Variables(); ok {
		_spec.SetField(applicationmodulerelationship.FieldVariables, field.TypeJSON, value)
	}
	if amruo.mutation.VariablesCleared() {
		_spec.ClearField(applicationmodulerelationship.FieldVariables, field.TypeJSON)
	}
	_spec.Node.Schema = amruo.schemaConfig.ApplicationModuleRelationship
	ctx = internal.NewSchemaConfigContext(ctx, amruo.schemaConfig)
	_spec.AddModifiers(amruo.modifiers...)
	_node = &ApplicationModuleRelationship{config: amruo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, amruo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{applicationmodulerelationship.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	amruo.mutation.done = true
	return _node, nil
}
