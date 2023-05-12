// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "ent". DO NOT EDIT.

package model

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/moduleversion"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/types"
)

// ModuleVersionUpdate is the builder for updating ModuleVersion entities.
type ModuleVersionUpdate struct {
	config
	hooks     []Hook
	mutation  *ModuleVersionMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the ModuleVersionUpdate builder.
func (mvu *ModuleVersionUpdate) Where(ps ...predicate.ModuleVersion) *ModuleVersionUpdate {
	mvu.mutation.Where(ps...)
	return mvu
}

// SetUpdateTime sets the "updateTime" field.
func (mvu *ModuleVersionUpdate) SetUpdateTime(t time.Time) *ModuleVersionUpdate {
	mvu.mutation.SetUpdateTime(t)
	return mvu
}

// SetSchema sets the "schema" field.
func (mvu *ModuleVersionUpdate) SetSchema(ts *types.ModuleSchema) *ModuleVersionUpdate {
	mvu.mutation.SetSchema(ts)
	return mvu
}

// Mutation returns the ModuleVersionMutation object of the builder.
func (mvu *ModuleVersionUpdate) Mutation() *ModuleVersionMutation {
	return mvu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (mvu *ModuleVersionUpdate) Save(ctx context.Context) (int, error) {
	if err := mvu.defaults(); err != nil {
		return 0, err
	}
	return withHooks[int, ModuleVersionMutation](ctx, mvu.sqlSave, mvu.mutation, mvu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (mvu *ModuleVersionUpdate) SaveX(ctx context.Context) int {
	affected, err := mvu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (mvu *ModuleVersionUpdate) Exec(ctx context.Context) error {
	_, err := mvu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mvu *ModuleVersionUpdate) ExecX(ctx context.Context) {
	if err := mvu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (mvu *ModuleVersionUpdate) defaults() error {
	if _, ok := mvu.mutation.UpdateTime(); !ok {
		if moduleversion.UpdateDefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized moduleversion.UpdateDefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := moduleversion.UpdateDefaultUpdateTime()
		mvu.mutation.SetUpdateTime(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (mvu *ModuleVersionUpdate) check() error {
	if _, ok := mvu.mutation.ModuleID(); mvu.mutation.ModuleCleared() && !ok {
		return errors.New(`model: clearing a required unique edge "ModuleVersion.module"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (mvu *ModuleVersionUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ModuleVersionUpdate {
	mvu.modifiers = append(mvu.modifiers, modifiers...)
	return mvu
}

func (mvu *ModuleVersionUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := mvu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(moduleversion.Table, moduleversion.Columns, sqlgraph.NewFieldSpec(moduleversion.FieldID, field.TypeString))
	if ps := mvu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := mvu.mutation.UpdateTime(); ok {
		_spec.SetField(moduleversion.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := mvu.mutation.Schema(); ok {
		_spec.SetField(moduleversion.FieldSchema, field.TypeJSON, value)
	}
	_spec.Node.Schema = mvu.schemaConfig.ModuleVersion
	ctx = internal.NewSchemaConfigContext(ctx, mvu.schemaConfig)
	_spec.AddModifiers(mvu.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, mvu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{moduleversion.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	mvu.mutation.done = true
	return n, nil
}

// ModuleVersionUpdateOne is the builder for updating a single ModuleVersion entity.
type ModuleVersionUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *ModuleVersionMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetUpdateTime sets the "updateTime" field.
func (mvuo *ModuleVersionUpdateOne) SetUpdateTime(t time.Time) *ModuleVersionUpdateOne {
	mvuo.mutation.SetUpdateTime(t)
	return mvuo
}

// SetSchema sets the "schema" field.
func (mvuo *ModuleVersionUpdateOne) SetSchema(ts *types.ModuleSchema) *ModuleVersionUpdateOne {
	mvuo.mutation.SetSchema(ts)
	return mvuo
}

// Mutation returns the ModuleVersionMutation object of the builder.
func (mvuo *ModuleVersionUpdateOne) Mutation() *ModuleVersionMutation {
	return mvuo.mutation
}

// Where appends a list predicates to the ModuleVersionUpdate builder.
func (mvuo *ModuleVersionUpdateOne) Where(ps ...predicate.ModuleVersion) *ModuleVersionUpdateOne {
	mvuo.mutation.Where(ps...)
	return mvuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (mvuo *ModuleVersionUpdateOne) Select(field string, fields ...string) *ModuleVersionUpdateOne {
	mvuo.fields = append([]string{field}, fields...)
	return mvuo
}

// Save executes the query and returns the updated ModuleVersion entity.
func (mvuo *ModuleVersionUpdateOne) Save(ctx context.Context) (*ModuleVersion, error) {
	if err := mvuo.defaults(); err != nil {
		return nil, err
	}
	return withHooks[*ModuleVersion, ModuleVersionMutation](ctx, mvuo.sqlSave, mvuo.mutation, mvuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (mvuo *ModuleVersionUpdateOne) SaveX(ctx context.Context) *ModuleVersion {
	node, err := mvuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (mvuo *ModuleVersionUpdateOne) Exec(ctx context.Context) error {
	_, err := mvuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mvuo *ModuleVersionUpdateOne) ExecX(ctx context.Context) {
	if err := mvuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (mvuo *ModuleVersionUpdateOne) defaults() error {
	if _, ok := mvuo.mutation.UpdateTime(); !ok {
		if moduleversion.UpdateDefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized moduleversion.UpdateDefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := moduleversion.UpdateDefaultUpdateTime()
		mvuo.mutation.SetUpdateTime(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (mvuo *ModuleVersionUpdateOne) check() error {
	if _, ok := mvuo.mutation.ModuleID(); mvuo.mutation.ModuleCleared() && !ok {
		return errors.New(`model: clearing a required unique edge "ModuleVersion.module"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (mvuo *ModuleVersionUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ModuleVersionUpdateOne {
	mvuo.modifiers = append(mvuo.modifiers, modifiers...)
	return mvuo
}

func (mvuo *ModuleVersionUpdateOne) sqlSave(ctx context.Context) (_node *ModuleVersion, err error) {
	if err := mvuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(moduleversion.Table, moduleversion.Columns, sqlgraph.NewFieldSpec(moduleversion.FieldID, field.TypeString))
	id, ok := mvuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`model: missing "ModuleVersion.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := mvuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, moduleversion.FieldID)
		for _, f := range fields {
			if !moduleversion.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("model: invalid field %q for query", f)}
			}
			if f != moduleversion.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := mvuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := mvuo.mutation.UpdateTime(); ok {
		_spec.SetField(moduleversion.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := mvuo.mutation.Schema(); ok {
		_spec.SetField(moduleversion.FieldSchema, field.TypeJSON, value)
	}
	_spec.Node.Schema = mvuo.schemaConfig.ModuleVersion
	ctx = internal.NewSchemaConfigContext(ctx, mvuo.schemaConfig)
	_spec.AddModifiers(mvuo.modifiers...)
	_node = &ModuleVersion{config: mvuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, mvuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{moduleversion.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	mvuo.mutation.done = true
	return _node, nil
}
