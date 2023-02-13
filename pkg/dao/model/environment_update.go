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
	"entgo.io/ent/dialect/sql/sqljson"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/oid"
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

// SetUpdateTime sets the "updateTime" field.
func (eu *EnvironmentUpdate) SetUpdateTime(t time.Time) *EnvironmentUpdate {
	eu.mutation.SetUpdateTime(t)
	return eu
}

// SetConnectorIDs sets the "connectorIDs" field.
func (eu *EnvironmentUpdate) SetConnectorIDs(o []oid.ID) *EnvironmentUpdate {
	eu.mutation.SetConnectorIDs(o)
	return eu
}

// AppendConnectorIDs appends o to the "connectorIDs" field.
func (eu *EnvironmentUpdate) AppendConnectorIDs(o []oid.ID) *EnvironmentUpdate {
	eu.mutation.AppendConnectorIDs(o)
	return eu
}

// ClearConnectorIDs clears the value of the "connectorIDs" field.
func (eu *EnvironmentUpdate) ClearConnectorIDs() *EnvironmentUpdate {
	eu.mutation.ClearConnectorIDs()
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

// Mutation returns the EnvironmentMutation object of the builder.
func (eu *EnvironmentUpdate) Mutation() *EnvironmentMutation {
	return eu.mutation
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
				Type:   field.TypeOther,
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
	if value, ok := eu.mutation.UpdateTime(); ok {
		_spec.SetField(environment.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := eu.mutation.ConnectorIDs(); ok {
		_spec.SetField(environment.FieldConnectorIDs, field.TypeJSON, value)
	}
	if value, ok := eu.mutation.AppendedConnectorIDs(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, environment.FieldConnectorIDs, value)
		})
	}
	if eu.mutation.ConnectorIDsCleared() {
		_spec.ClearField(environment.FieldConnectorIDs, field.TypeJSON)
	}
	if value, ok := eu.mutation.Variables(); ok {
		_spec.SetField(environment.FieldVariables, field.TypeJSON, value)
	}
	if eu.mutation.VariablesCleared() {
		_spec.ClearField(environment.FieldVariables, field.TypeJSON)
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

// SetUpdateTime sets the "updateTime" field.
func (euo *EnvironmentUpdateOne) SetUpdateTime(t time.Time) *EnvironmentUpdateOne {
	euo.mutation.SetUpdateTime(t)
	return euo
}

// SetConnectorIDs sets the "connectorIDs" field.
func (euo *EnvironmentUpdateOne) SetConnectorIDs(o []oid.ID) *EnvironmentUpdateOne {
	euo.mutation.SetConnectorIDs(o)
	return euo
}

// AppendConnectorIDs appends o to the "connectorIDs" field.
func (euo *EnvironmentUpdateOne) AppendConnectorIDs(o []oid.ID) *EnvironmentUpdateOne {
	euo.mutation.AppendConnectorIDs(o)
	return euo
}

// ClearConnectorIDs clears the value of the "connectorIDs" field.
func (euo *EnvironmentUpdateOne) ClearConnectorIDs() *EnvironmentUpdateOne {
	euo.mutation.ClearConnectorIDs()
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

// Mutation returns the EnvironmentMutation object of the builder.
func (euo *EnvironmentUpdateOne) Mutation() *EnvironmentMutation {
	return euo.mutation
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
				Type:   field.TypeOther,
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
	if value, ok := euo.mutation.UpdateTime(); ok {
		_spec.SetField(environment.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := euo.mutation.ConnectorIDs(); ok {
		_spec.SetField(environment.FieldConnectorIDs, field.TypeJSON, value)
	}
	if value, ok := euo.mutation.AppendedConnectorIDs(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, environment.FieldConnectorIDs, value)
		})
	}
	if euo.mutation.ConnectorIDsCleared() {
		_spec.ClearField(environment.FieldConnectorIDs, field.TypeJSON)
	}
	if value, ok := euo.mutation.Variables(); ok {
		_spec.SetField(environment.FieldVariables, field.TypeJSON, value)
	}
	if euo.mutation.VariablesCleared() {
		_spec.ClearField(environment.FieldVariables, field.TypeJSON)
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
