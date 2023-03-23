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
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/role"
	"github.com/seal-io/seal/pkg/dao/types"
)

// RoleUpdate is the builder for updating Role entities.
type RoleUpdate struct {
	config
	hooks     []Hook
	mutation  *RoleMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the RoleUpdate builder.
func (ru *RoleUpdate) Where(ps ...predicate.Role) *RoleUpdate {
	ru.mutation.Where(ps...)
	return ru
}

// SetUpdateTime sets the "updateTime" field.
func (ru *RoleUpdate) SetUpdateTime(t time.Time) *RoleUpdate {
	ru.mutation.SetUpdateTime(t)
	return ru
}

// SetDescription sets the "description" field.
func (ru *RoleUpdate) SetDescription(s string) *RoleUpdate {
	ru.mutation.SetDescription(s)
	return ru
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (ru *RoleUpdate) SetNillableDescription(s *string) *RoleUpdate {
	if s != nil {
		ru.SetDescription(*s)
	}
	return ru
}

// ClearDescription clears the value of the "description" field.
func (ru *RoleUpdate) ClearDescription() *RoleUpdate {
	ru.mutation.ClearDescription()
	return ru
}

// SetPolicies sets the "policies" field.
func (ru *RoleUpdate) SetPolicies(tp types.RolePolicies) *RoleUpdate {
	ru.mutation.SetPolicies(tp)
	return ru
}

// AppendPolicies appends tp to the "policies" field.
func (ru *RoleUpdate) AppendPolicies(tp types.RolePolicies) *RoleUpdate {
	ru.mutation.AppendPolicies(tp)
	return ru
}

// Mutation returns the RoleMutation object of the builder.
func (ru *RoleUpdate) Mutation() *RoleMutation {
	return ru.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (ru *RoleUpdate) Save(ctx context.Context) (int, error) {
	if err := ru.defaults(); err != nil {
		return 0, err
	}
	return withHooks[int, RoleMutation](ctx, ru.sqlSave, ru.mutation, ru.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ru *RoleUpdate) SaveX(ctx context.Context) int {
	affected, err := ru.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (ru *RoleUpdate) Exec(ctx context.Context) error {
	_, err := ru.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ru *RoleUpdate) ExecX(ctx context.Context) {
	if err := ru.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ru *RoleUpdate) defaults() error {
	if _, ok := ru.mutation.UpdateTime(); !ok {
		if role.UpdateDefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized role.UpdateDefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := role.UpdateDefaultUpdateTime()
		ru.mutation.SetUpdateTime(v)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (ru *RoleUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *RoleUpdate {
	ru.modifiers = append(ru.modifiers, modifiers...)
	return ru
}

func (ru *RoleUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(role.Table, role.Columns, sqlgraph.NewFieldSpec(role.FieldID, field.TypeString))
	if ps := ru.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ru.mutation.UpdateTime(); ok {
		_spec.SetField(role.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := ru.mutation.Description(); ok {
		_spec.SetField(role.FieldDescription, field.TypeString, value)
	}
	if ru.mutation.DescriptionCleared() {
		_spec.ClearField(role.FieldDescription, field.TypeString)
	}
	if value, ok := ru.mutation.Policies(); ok {
		_spec.SetField(role.FieldPolicies, field.TypeJSON, value)
	}
	if value, ok := ru.mutation.AppendedPolicies(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, role.FieldPolicies, value)
		})
	}
	_spec.Node.Schema = ru.schemaConfig.Role
	ctx = internal.NewSchemaConfigContext(ctx, ru.schemaConfig)
	_spec.AddModifiers(ru.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, ru.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{role.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	ru.mutation.done = true
	return n, nil
}

// RoleUpdateOne is the builder for updating a single Role entity.
type RoleUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *RoleMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetUpdateTime sets the "updateTime" field.
func (ruo *RoleUpdateOne) SetUpdateTime(t time.Time) *RoleUpdateOne {
	ruo.mutation.SetUpdateTime(t)
	return ruo
}

// SetDescription sets the "description" field.
func (ruo *RoleUpdateOne) SetDescription(s string) *RoleUpdateOne {
	ruo.mutation.SetDescription(s)
	return ruo
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (ruo *RoleUpdateOne) SetNillableDescription(s *string) *RoleUpdateOne {
	if s != nil {
		ruo.SetDescription(*s)
	}
	return ruo
}

// ClearDescription clears the value of the "description" field.
func (ruo *RoleUpdateOne) ClearDescription() *RoleUpdateOne {
	ruo.mutation.ClearDescription()
	return ruo
}

// SetPolicies sets the "policies" field.
func (ruo *RoleUpdateOne) SetPolicies(tp types.RolePolicies) *RoleUpdateOne {
	ruo.mutation.SetPolicies(tp)
	return ruo
}

// AppendPolicies appends tp to the "policies" field.
func (ruo *RoleUpdateOne) AppendPolicies(tp types.RolePolicies) *RoleUpdateOne {
	ruo.mutation.AppendPolicies(tp)
	return ruo
}

// Mutation returns the RoleMutation object of the builder.
func (ruo *RoleUpdateOne) Mutation() *RoleMutation {
	return ruo.mutation
}

// Where appends a list predicates to the RoleUpdate builder.
func (ruo *RoleUpdateOne) Where(ps ...predicate.Role) *RoleUpdateOne {
	ruo.mutation.Where(ps...)
	return ruo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (ruo *RoleUpdateOne) Select(field string, fields ...string) *RoleUpdateOne {
	ruo.fields = append([]string{field}, fields...)
	return ruo
}

// Save executes the query and returns the updated Role entity.
func (ruo *RoleUpdateOne) Save(ctx context.Context) (*Role, error) {
	if err := ruo.defaults(); err != nil {
		return nil, err
	}
	return withHooks[*Role, RoleMutation](ctx, ruo.sqlSave, ruo.mutation, ruo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ruo *RoleUpdateOne) SaveX(ctx context.Context) *Role {
	node, err := ruo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (ruo *RoleUpdateOne) Exec(ctx context.Context) error {
	_, err := ruo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ruo *RoleUpdateOne) ExecX(ctx context.Context) {
	if err := ruo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ruo *RoleUpdateOne) defaults() error {
	if _, ok := ruo.mutation.UpdateTime(); !ok {
		if role.UpdateDefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized role.UpdateDefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := role.UpdateDefaultUpdateTime()
		ruo.mutation.SetUpdateTime(v)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (ruo *RoleUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *RoleUpdateOne {
	ruo.modifiers = append(ruo.modifiers, modifiers...)
	return ruo
}

func (ruo *RoleUpdateOne) sqlSave(ctx context.Context) (_node *Role, err error) {
	_spec := sqlgraph.NewUpdateSpec(role.Table, role.Columns, sqlgraph.NewFieldSpec(role.FieldID, field.TypeString))
	id, ok := ruo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`model: missing "Role.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := ruo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, role.FieldID)
		for _, f := range fields {
			if !role.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("model: invalid field %q for query", f)}
			}
			if f != role.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := ruo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ruo.mutation.UpdateTime(); ok {
		_spec.SetField(role.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := ruo.mutation.Description(); ok {
		_spec.SetField(role.FieldDescription, field.TypeString, value)
	}
	if ruo.mutation.DescriptionCleared() {
		_spec.ClearField(role.FieldDescription, field.TypeString)
	}
	if value, ok := ruo.mutation.Policies(); ok {
		_spec.SetField(role.FieldPolicies, field.TypeJSON, value)
	}
	if value, ok := ruo.mutation.AppendedPolicies(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, role.FieldPolicies, value)
		})
	}
	_spec.Node.Schema = ruo.schemaConfig.Role
	ctx = internal.NewSchemaConfigContext(ctx, ruo.schemaConfig)
	_spec.AddModifiers(ruo.modifiers...)
	_node = &Role{config: ruo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, ruo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{role.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	ruo.mutation.done = true
	return _node, nil
}
