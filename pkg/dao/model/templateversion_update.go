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
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/templateversion"
	"github.com/seal-io/seal/pkg/dao/types"
)

// TemplateVersionUpdate is the builder for updating TemplateVersion entities.
type TemplateVersionUpdate struct {
	config
	hooks     []Hook
	mutation  *TemplateVersionMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the TemplateVersionUpdate builder.
func (tvu *TemplateVersionUpdate) Where(ps ...predicate.TemplateVersion) *TemplateVersionUpdate {
	tvu.mutation.Where(ps...)
	return tvu
}

// SetUpdateTime sets the "updateTime" field.
func (tvu *TemplateVersionUpdate) SetUpdateTime(t time.Time) *TemplateVersionUpdate {
	tvu.mutation.SetUpdateTime(t)
	return tvu
}

// SetSchema sets the "schema" field.
func (tvu *TemplateVersionUpdate) SetSchema(ts *types.TemplateSchema) *TemplateVersionUpdate {
	tvu.mutation.SetSchema(ts)
	return tvu
}

// Mutation returns the TemplateVersionMutation object of the builder.
func (tvu *TemplateVersionUpdate) Mutation() *TemplateVersionMutation {
	return tvu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (tvu *TemplateVersionUpdate) Save(ctx context.Context) (int, error) {
	if err := tvu.defaults(); err != nil {
		return 0, err
	}
	return withHooks[int, TemplateVersionMutation](ctx, tvu.sqlSave, tvu.mutation, tvu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (tvu *TemplateVersionUpdate) SaveX(ctx context.Context) int {
	affected, err := tvu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (tvu *TemplateVersionUpdate) Exec(ctx context.Context) error {
	_, err := tvu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tvu *TemplateVersionUpdate) ExecX(ctx context.Context) {
	if err := tvu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (tvu *TemplateVersionUpdate) defaults() error {
	if _, ok := tvu.mutation.UpdateTime(); !ok {
		if templateversion.UpdateDefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized templateversion.UpdateDefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := templateversion.UpdateDefaultUpdateTime()
		tvu.mutation.SetUpdateTime(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (tvu *TemplateVersionUpdate) check() error {
	if _, ok := tvu.mutation.TemplateID(); tvu.mutation.TemplateCleared() && !ok {
		return errors.New(`model: clearing a required unique edge "TemplateVersion.template"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (tvu *TemplateVersionUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *TemplateVersionUpdate {
	tvu.modifiers = append(tvu.modifiers, modifiers...)
	return tvu
}

func (tvu *TemplateVersionUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := tvu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(templateversion.Table, templateversion.Columns, sqlgraph.NewFieldSpec(templateversion.FieldID, field.TypeString))
	if ps := tvu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := tvu.mutation.UpdateTime(); ok {
		_spec.SetField(templateversion.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := tvu.mutation.Schema(); ok {
		_spec.SetField(templateversion.FieldSchema, field.TypeJSON, value)
	}
	_spec.Node.Schema = tvu.schemaConfig.TemplateVersion
	ctx = internal.NewSchemaConfigContext(ctx, tvu.schemaConfig)
	_spec.AddModifiers(tvu.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, tvu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{templateversion.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	tvu.mutation.done = true
	return n, nil
}

// TemplateVersionUpdateOne is the builder for updating a single TemplateVersion entity.
type TemplateVersionUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *TemplateVersionMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetUpdateTime sets the "updateTime" field.
func (tvuo *TemplateVersionUpdateOne) SetUpdateTime(t time.Time) *TemplateVersionUpdateOne {
	tvuo.mutation.SetUpdateTime(t)
	return tvuo
}

// SetSchema sets the "schema" field.
func (tvuo *TemplateVersionUpdateOne) SetSchema(ts *types.TemplateSchema) *TemplateVersionUpdateOne {
	tvuo.mutation.SetSchema(ts)
	return tvuo
}

// Mutation returns the TemplateVersionMutation object of the builder.
func (tvuo *TemplateVersionUpdateOne) Mutation() *TemplateVersionMutation {
	return tvuo.mutation
}

// Where appends a list predicates to the TemplateVersionUpdate builder.
func (tvuo *TemplateVersionUpdateOne) Where(ps ...predicate.TemplateVersion) *TemplateVersionUpdateOne {
	tvuo.mutation.Where(ps...)
	return tvuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (tvuo *TemplateVersionUpdateOne) Select(field string, fields ...string) *TemplateVersionUpdateOne {
	tvuo.fields = append([]string{field}, fields...)
	return tvuo
}

// Save executes the query and returns the updated TemplateVersion entity.
func (tvuo *TemplateVersionUpdateOne) Save(ctx context.Context) (*TemplateVersion, error) {
	if err := tvuo.defaults(); err != nil {
		return nil, err
	}
	return withHooks[*TemplateVersion, TemplateVersionMutation](ctx, tvuo.sqlSave, tvuo.mutation, tvuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (tvuo *TemplateVersionUpdateOne) SaveX(ctx context.Context) *TemplateVersion {
	node, err := tvuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (tvuo *TemplateVersionUpdateOne) Exec(ctx context.Context) error {
	_, err := tvuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tvuo *TemplateVersionUpdateOne) ExecX(ctx context.Context) {
	if err := tvuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (tvuo *TemplateVersionUpdateOne) defaults() error {
	if _, ok := tvuo.mutation.UpdateTime(); !ok {
		if templateversion.UpdateDefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized templateversion.UpdateDefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := templateversion.UpdateDefaultUpdateTime()
		tvuo.mutation.SetUpdateTime(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (tvuo *TemplateVersionUpdateOne) check() error {
	if _, ok := tvuo.mutation.TemplateID(); tvuo.mutation.TemplateCleared() && !ok {
		return errors.New(`model: clearing a required unique edge "TemplateVersion.template"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (tvuo *TemplateVersionUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *TemplateVersionUpdateOne {
	tvuo.modifiers = append(tvuo.modifiers, modifiers...)
	return tvuo
}

func (tvuo *TemplateVersionUpdateOne) sqlSave(ctx context.Context) (_node *TemplateVersion, err error) {
	if err := tvuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(templateversion.Table, templateversion.Columns, sqlgraph.NewFieldSpec(templateversion.FieldID, field.TypeString))
	id, ok := tvuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`model: missing "TemplateVersion.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := tvuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, templateversion.FieldID)
		for _, f := range fields {
			if !templateversion.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("model: invalid field %q for query", f)}
			}
			if f != templateversion.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := tvuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := tvuo.mutation.UpdateTime(); ok {
		_spec.SetField(templateversion.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := tvuo.mutation.Schema(); ok {
		_spec.SetField(templateversion.FieldSchema, field.TypeJSON, value)
	}
	_spec.Node.Schema = tvuo.schemaConfig.TemplateVersion
	ctx = internal.NewSchemaConfigContext(ctx, tvuo.schemaConfig)
	_spec.AddModifiers(tvuo.modifiers...)
	_node = &TemplateVersion{config: tvuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, tvuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{templateversion.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	tvuo.mutation.done = true
	return _node, nil
}
