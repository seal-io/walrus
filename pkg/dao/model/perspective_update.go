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

	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/perspective"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/types"
)

// PerspectiveUpdate is the builder for updating Perspective entities.
type PerspectiveUpdate struct {
	config
	hooks     []Hook
	mutation  *PerspectiveMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the PerspectiveUpdate builder.
func (pu *PerspectiveUpdate) Where(ps ...predicate.Perspective) *PerspectiveUpdate {
	pu.mutation.Where(ps...)
	return pu
}

// SetUpdateTime sets the "updateTime" field.
func (pu *PerspectiveUpdate) SetUpdateTime(t time.Time) *PerspectiveUpdate {
	pu.mutation.SetUpdateTime(t)
	return pu
}

// SetStartTime sets the "startTime" field.
func (pu *PerspectiveUpdate) SetStartTime(s string) *PerspectiveUpdate {
	pu.mutation.SetStartTime(s)
	return pu
}

// SetEndTime sets the "endTime" field.
func (pu *PerspectiveUpdate) SetEndTime(s string) *PerspectiveUpdate {
	pu.mutation.SetEndTime(s)
	return pu
}

// SetBuiltin sets the "builtin" field.
func (pu *PerspectiveUpdate) SetBuiltin(b bool) *PerspectiveUpdate {
	pu.mutation.SetBuiltin(b)
	return pu
}

// SetNillableBuiltin sets the "builtin" field if the given value is not nil.
func (pu *PerspectiveUpdate) SetNillableBuiltin(b *bool) *PerspectiveUpdate {
	if b != nil {
		pu.SetBuiltin(*b)
	}
	return pu
}

// SetAllocationQueries sets the "allocationQueries" field.
func (pu *PerspectiveUpdate) SetAllocationQueries(tc []types.QueryCondition) *PerspectiveUpdate {
	pu.mutation.SetAllocationQueries(tc)
	return pu
}

// AppendAllocationQueries appends tc to the "allocationQueries" field.
func (pu *PerspectiveUpdate) AppendAllocationQueries(tc []types.QueryCondition) *PerspectiveUpdate {
	pu.mutation.AppendAllocationQueries(tc)
	return pu
}

// Mutation returns the PerspectiveMutation object of the builder.
func (pu *PerspectiveUpdate) Mutation() *PerspectiveMutation {
	return pu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (pu *PerspectiveUpdate) Save(ctx context.Context) (int, error) {
	if err := pu.defaults(); err != nil {
		return 0, err
	}
	return withHooks[int, PerspectiveMutation](ctx, pu.sqlSave, pu.mutation, pu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (pu *PerspectiveUpdate) SaveX(ctx context.Context) int {
	affected, err := pu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (pu *PerspectiveUpdate) Exec(ctx context.Context) error {
	_, err := pu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pu *PerspectiveUpdate) ExecX(ctx context.Context) {
	if err := pu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (pu *PerspectiveUpdate) defaults() error {
	if _, ok := pu.mutation.UpdateTime(); !ok {
		if perspective.UpdateDefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized perspective.UpdateDefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := perspective.UpdateDefaultUpdateTime()
		pu.mutation.SetUpdateTime(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (pu *PerspectiveUpdate) check() error {
	if v, ok := pu.mutation.StartTime(); ok {
		if err := perspective.StartTimeValidator(v); err != nil {
			return &ValidationError{Name: "startTime", err: fmt.Errorf(`model: validator failed for field "Perspective.startTime": %w`, err)}
		}
	}
	if v, ok := pu.mutation.EndTime(); ok {
		if err := perspective.EndTimeValidator(v); err != nil {
			return &ValidationError{Name: "endTime", err: fmt.Errorf(`model: validator failed for field "Perspective.endTime": %w`, err)}
		}
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (pu *PerspectiveUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *PerspectiveUpdate {
	pu.modifiers = append(pu.modifiers, modifiers...)
	return pu
}

func (pu *PerspectiveUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := pu.check(); err != nil {
		return n, err
	}
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   perspective.Table,
			Columns: perspective.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: perspective.FieldID,
			},
		},
	}
	if ps := pu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := pu.mutation.UpdateTime(); ok {
		_spec.SetField(perspective.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := pu.mutation.StartTime(); ok {
		_spec.SetField(perspective.FieldStartTime, field.TypeString, value)
	}
	if value, ok := pu.mutation.EndTime(); ok {
		_spec.SetField(perspective.FieldEndTime, field.TypeString, value)
	}
	if value, ok := pu.mutation.Builtin(); ok {
		_spec.SetField(perspective.FieldBuiltin, field.TypeBool, value)
	}
	if value, ok := pu.mutation.AllocationQueries(); ok {
		_spec.SetField(perspective.FieldAllocationQueries, field.TypeJSON, value)
	}
	if value, ok := pu.mutation.AppendedAllocationQueries(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, perspective.FieldAllocationQueries, value)
		})
	}
	_spec.Node.Schema = pu.schemaConfig.Perspective
	ctx = internal.NewSchemaConfigContext(ctx, pu.schemaConfig)
	_spec.AddModifiers(pu.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, pu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{perspective.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	pu.mutation.done = true
	return n, nil
}

// PerspectiveUpdateOne is the builder for updating a single Perspective entity.
type PerspectiveUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *PerspectiveMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetUpdateTime sets the "updateTime" field.
func (puo *PerspectiveUpdateOne) SetUpdateTime(t time.Time) *PerspectiveUpdateOne {
	puo.mutation.SetUpdateTime(t)
	return puo
}

// SetStartTime sets the "startTime" field.
func (puo *PerspectiveUpdateOne) SetStartTime(s string) *PerspectiveUpdateOne {
	puo.mutation.SetStartTime(s)
	return puo
}

// SetEndTime sets the "endTime" field.
func (puo *PerspectiveUpdateOne) SetEndTime(s string) *PerspectiveUpdateOne {
	puo.mutation.SetEndTime(s)
	return puo
}

// SetBuiltin sets the "builtin" field.
func (puo *PerspectiveUpdateOne) SetBuiltin(b bool) *PerspectiveUpdateOne {
	puo.mutation.SetBuiltin(b)
	return puo
}

// SetNillableBuiltin sets the "builtin" field if the given value is not nil.
func (puo *PerspectiveUpdateOne) SetNillableBuiltin(b *bool) *PerspectiveUpdateOne {
	if b != nil {
		puo.SetBuiltin(*b)
	}
	return puo
}

// SetAllocationQueries sets the "allocationQueries" field.
func (puo *PerspectiveUpdateOne) SetAllocationQueries(tc []types.QueryCondition) *PerspectiveUpdateOne {
	puo.mutation.SetAllocationQueries(tc)
	return puo
}

// AppendAllocationQueries appends tc to the "allocationQueries" field.
func (puo *PerspectiveUpdateOne) AppendAllocationQueries(tc []types.QueryCondition) *PerspectiveUpdateOne {
	puo.mutation.AppendAllocationQueries(tc)
	return puo
}

// Mutation returns the PerspectiveMutation object of the builder.
func (puo *PerspectiveUpdateOne) Mutation() *PerspectiveMutation {
	return puo.mutation
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (puo *PerspectiveUpdateOne) Select(field string, fields ...string) *PerspectiveUpdateOne {
	puo.fields = append([]string{field}, fields...)
	return puo
}

// Save executes the query and returns the updated Perspective entity.
func (puo *PerspectiveUpdateOne) Save(ctx context.Context) (*Perspective, error) {
	if err := puo.defaults(); err != nil {
		return nil, err
	}
	return withHooks[*Perspective, PerspectiveMutation](ctx, puo.sqlSave, puo.mutation, puo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (puo *PerspectiveUpdateOne) SaveX(ctx context.Context) *Perspective {
	node, err := puo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (puo *PerspectiveUpdateOne) Exec(ctx context.Context) error {
	_, err := puo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (puo *PerspectiveUpdateOne) ExecX(ctx context.Context) {
	if err := puo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (puo *PerspectiveUpdateOne) defaults() error {
	if _, ok := puo.mutation.UpdateTime(); !ok {
		if perspective.UpdateDefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized perspective.UpdateDefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := perspective.UpdateDefaultUpdateTime()
		puo.mutation.SetUpdateTime(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (puo *PerspectiveUpdateOne) check() error {
	if v, ok := puo.mutation.StartTime(); ok {
		if err := perspective.StartTimeValidator(v); err != nil {
			return &ValidationError{Name: "startTime", err: fmt.Errorf(`model: validator failed for field "Perspective.startTime": %w`, err)}
		}
	}
	if v, ok := puo.mutation.EndTime(); ok {
		if err := perspective.EndTimeValidator(v); err != nil {
			return &ValidationError{Name: "endTime", err: fmt.Errorf(`model: validator failed for field "Perspective.endTime": %w`, err)}
		}
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (puo *PerspectiveUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *PerspectiveUpdateOne {
	puo.modifiers = append(puo.modifiers, modifiers...)
	return puo
}

func (puo *PerspectiveUpdateOne) sqlSave(ctx context.Context) (_node *Perspective, err error) {
	if err := puo.check(); err != nil {
		return _node, err
	}
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   perspective.Table,
			Columns: perspective.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: perspective.FieldID,
			},
		},
	}
	id, ok := puo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`model: missing "Perspective.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := puo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, perspective.FieldID)
		for _, f := range fields {
			if !perspective.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("model: invalid field %q for query", f)}
			}
			if f != perspective.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := puo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := puo.mutation.UpdateTime(); ok {
		_spec.SetField(perspective.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := puo.mutation.StartTime(); ok {
		_spec.SetField(perspective.FieldStartTime, field.TypeString, value)
	}
	if value, ok := puo.mutation.EndTime(); ok {
		_spec.SetField(perspective.FieldEndTime, field.TypeString, value)
	}
	if value, ok := puo.mutation.Builtin(); ok {
		_spec.SetField(perspective.FieldBuiltin, field.TypeBool, value)
	}
	if value, ok := puo.mutation.AllocationQueries(); ok {
		_spec.SetField(perspective.FieldAllocationQueries, field.TypeJSON, value)
	}
	if value, ok := puo.mutation.AppendedAllocationQueries(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, perspective.FieldAllocationQueries, value)
		})
	}
	_spec.Node.Schema = puo.schemaConfig.Perspective
	ctx = internal.NewSchemaConfigContext(ctx, puo.schemaConfig)
	_spec.AddModifiers(puo.modifiers...)
	_node = &Perspective{config: puo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, puo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{perspective.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	puo.mutation.done = true
	return _node, nil
}
