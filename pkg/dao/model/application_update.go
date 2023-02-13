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

	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/schema"
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

// SetUpdateTime sets the "updateTime" field.
func (au *ApplicationUpdate) SetUpdateTime(t time.Time) *ApplicationUpdate {
	au.mutation.SetUpdateTime(t)
	return au
}

// SetModules sets the "modules" field.
func (au *ApplicationUpdate) SetModules(sm []schema.ApplicationModule) *ApplicationUpdate {
	au.mutation.SetModules(sm)
	return au
}

// AppendModules appends sm to the "modules" field.
func (au *ApplicationUpdate) AppendModules(sm []schema.ApplicationModule) *ApplicationUpdate {
	au.mutation.AppendModules(sm)
	return au
}

// Mutation returns the ApplicationMutation object of the builder.
func (au *ApplicationUpdate) Mutation() *ApplicationMutation {
	return au.mutation
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

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (au *ApplicationUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ApplicationUpdate {
	au.modifiers = append(au.modifiers, modifiers...)
	return au
}

func (au *ApplicationUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   application.Table,
			Columns: application.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeOther,
				Column: application.FieldID,
			},
		},
	}
	if ps := au.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := au.mutation.UpdateTime(); ok {
		_spec.SetField(application.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := au.mutation.Modules(); ok {
		_spec.SetField(application.FieldModules, field.TypeJSON, value)
	}
	if value, ok := au.mutation.AppendedModules(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, application.FieldModules, value)
		})
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

// SetUpdateTime sets the "updateTime" field.
func (auo *ApplicationUpdateOne) SetUpdateTime(t time.Time) *ApplicationUpdateOne {
	auo.mutation.SetUpdateTime(t)
	return auo
}

// SetModules sets the "modules" field.
func (auo *ApplicationUpdateOne) SetModules(sm []schema.ApplicationModule) *ApplicationUpdateOne {
	auo.mutation.SetModules(sm)
	return auo
}

// AppendModules appends sm to the "modules" field.
func (auo *ApplicationUpdateOne) AppendModules(sm []schema.ApplicationModule) *ApplicationUpdateOne {
	auo.mutation.AppendModules(sm)
	return auo
}

// Mutation returns the ApplicationMutation object of the builder.
func (auo *ApplicationUpdateOne) Mutation() *ApplicationMutation {
	return auo.mutation
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

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (auo *ApplicationUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ApplicationUpdateOne {
	auo.modifiers = append(auo.modifiers, modifiers...)
	return auo
}

func (auo *ApplicationUpdateOne) sqlSave(ctx context.Context) (_node *Application, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   application.Table,
			Columns: application.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeOther,
				Column: application.FieldID,
			},
		},
	}
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
	if value, ok := auo.mutation.UpdateTime(); ok {
		_spec.SetField(application.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := auo.mutation.Modules(); ok {
		_spec.SetField(application.FieldModules, field.TypeJSON, value)
	}
	if value, ok := auo.mutation.AppendedModules(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, application.FieldModules, value)
		})
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
