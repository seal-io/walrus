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

	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/secret"
	"github.com/seal-io/seal/pkg/dao/types/crypto"
)

// SecretUpdate is the builder for updating Secret entities.
type SecretUpdate struct {
	config
	hooks     []Hook
	mutation  *SecretMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the SecretUpdate builder.
func (su *SecretUpdate) Where(ps ...predicate.Secret) *SecretUpdate {
	su.mutation.Where(ps...)
	return su
}

// SetUpdateTime sets the "updateTime" field.
func (su *SecretUpdate) SetUpdateTime(t time.Time) *SecretUpdate {
	su.mutation.SetUpdateTime(t)
	return su
}

// SetValue sets the "value" field.
func (su *SecretUpdate) SetValue(c crypto.String) *SecretUpdate {
	su.mutation.SetValue(c)
	return su
}

// Mutation returns the SecretMutation object of the builder.
func (su *SecretUpdate) Mutation() *SecretMutation {
	return su.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (su *SecretUpdate) Save(ctx context.Context) (int, error) {
	if err := su.defaults(); err != nil {
		return 0, err
	}
	return withHooks[int, SecretMutation](ctx, su.sqlSave, su.mutation, su.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (su *SecretUpdate) SaveX(ctx context.Context) int {
	affected, err := su.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (su *SecretUpdate) Exec(ctx context.Context) error {
	_, err := su.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (su *SecretUpdate) ExecX(ctx context.Context) {
	if err := su.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (su *SecretUpdate) defaults() error {
	if _, ok := su.mutation.UpdateTime(); !ok {
		if secret.UpdateDefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized secret.UpdateDefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := secret.UpdateDefaultUpdateTime()
		su.mutation.SetUpdateTime(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (su *SecretUpdate) check() error {
	if v, ok := su.mutation.Value(); ok {
		if err := secret.ValueValidator(string(v)); err != nil {
			return &ValidationError{Name: "value", err: fmt.Errorf(`model: validator failed for field "Secret.value": %w`, err)}
		}
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (su *SecretUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *SecretUpdate {
	su.modifiers = append(su.modifiers, modifiers...)
	return su
}

func (su *SecretUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := su.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(secret.Table, secret.Columns, sqlgraph.NewFieldSpec(secret.FieldID, field.TypeString))
	if ps := su.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := su.mutation.UpdateTime(); ok {
		_spec.SetField(secret.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := su.mutation.Value(); ok {
		_spec.SetField(secret.FieldValue, field.TypeString, value)
	}
	_spec.Node.Schema = su.schemaConfig.Secret
	ctx = internal.NewSchemaConfigContext(ctx, su.schemaConfig)
	_spec.AddModifiers(su.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, su.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{secret.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	su.mutation.done = true
	return n, nil
}

// SecretUpdateOne is the builder for updating a single Secret entity.
type SecretUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *SecretMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetUpdateTime sets the "updateTime" field.
func (suo *SecretUpdateOne) SetUpdateTime(t time.Time) *SecretUpdateOne {
	suo.mutation.SetUpdateTime(t)
	return suo
}

// SetValue sets the "value" field.
func (suo *SecretUpdateOne) SetValue(c crypto.String) *SecretUpdateOne {
	suo.mutation.SetValue(c)
	return suo
}

// Mutation returns the SecretMutation object of the builder.
func (suo *SecretUpdateOne) Mutation() *SecretMutation {
	return suo.mutation
}

// Where appends a list predicates to the SecretUpdate builder.
func (suo *SecretUpdateOne) Where(ps ...predicate.Secret) *SecretUpdateOne {
	suo.mutation.Where(ps...)
	return suo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (suo *SecretUpdateOne) Select(field string, fields ...string) *SecretUpdateOne {
	suo.fields = append([]string{field}, fields...)
	return suo
}

// Save executes the query and returns the updated Secret entity.
func (suo *SecretUpdateOne) Save(ctx context.Context) (*Secret, error) {
	if err := suo.defaults(); err != nil {
		return nil, err
	}
	return withHooks[*Secret, SecretMutation](ctx, suo.sqlSave, suo.mutation, suo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (suo *SecretUpdateOne) SaveX(ctx context.Context) *Secret {
	node, err := suo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (suo *SecretUpdateOne) Exec(ctx context.Context) error {
	_, err := suo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (suo *SecretUpdateOne) ExecX(ctx context.Context) {
	if err := suo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (suo *SecretUpdateOne) defaults() error {
	if _, ok := suo.mutation.UpdateTime(); !ok {
		if secret.UpdateDefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized secret.UpdateDefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := secret.UpdateDefaultUpdateTime()
		suo.mutation.SetUpdateTime(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (suo *SecretUpdateOne) check() error {
	if v, ok := suo.mutation.Value(); ok {
		if err := secret.ValueValidator(string(v)); err != nil {
			return &ValidationError{Name: "value", err: fmt.Errorf(`model: validator failed for field "Secret.value": %w`, err)}
		}
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (suo *SecretUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *SecretUpdateOne {
	suo.modifiers = append(suo.modifiers, modifiers...)
	return suo
}

func (suo *SecretUpdateOne) sqlSave(ctx context.Context) (_node *Secret, err error) {
	if err := suo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(secret.Table, secret.Columns, sqlgraph.NewFieldSpec(secret.FieldID, field.TypeString))
	id, ok := suo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`model: missing "Secret.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := suo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, secret.FieldID)
		for _, f := range fields {
			if !secret.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("model: invalid field %q for query", f)}
			}
			if f != secret.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := suo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := suo.mutation.UpdateTime(); ok {
		_spec.SetField(secret.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := suo.mutation.Value(); ok {
		_spec.SetField(secret.FieldValue, field.TypeString, value)
	}
	_spec.Node.Schema = suo.schemaConfig.Secret
	ctx = internal.NewSchemaConfigContext(ctx, suo.schemaConfig)
	_spec.AddModifiers(suo.modifiers...)
	_node = &Secret{config: suo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, suo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{secret.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	suo.mutation.done = true
	return _node, nil
}
