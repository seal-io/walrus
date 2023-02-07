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
	"github.com/seal-io/seal/pkg/dao/model/setting"
)

// SettingUpdate is the builder for updating Setting entities.
type SettingUpdate struct {
	config
	hooks     []Hook
	mutation  *SettingMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the SettingUpdate builder.
func (su *SettingUpdate) Where(ps ...predicate.Setting) *SettingUpdate {
	su.mutation.Where(ps...)
	return su
}

// SetUpdateTime sets the "updateTime" field.
func (su *SettingUpdate) SetUpdateTime(t time.Time) *SettingUpdate {
	su.mutation.SetUpdateTime(t)
	return su
}

// SetName sets the "name" field.
func (su *SettingUpdate) SetName(s string) *SettingUpdate {
	su.mutation.SetName(s)
	return su
}

// SetValue sets the "value" field.
func (su *SettingUpdate) SetValue(s string) *SettingUpdate {
	su.mutation.SetValue(s)
	return su
}

// SetNillableValue sets the "value" field if the given value is not nil.
func (su *SettingUpdate) SetNillableValue(s *string) *SettingUpdate {
	if s != nil {
		su.SetValue(*s)
	}
	return su
}

// SetHidden sets the "hidden" field.
func (su *SettingUpdate) SetHidden(b bool) *SettingUpdate {
	su.mutation.SetHidden(b)
	return su
}

// SetNillableHidden sets the "hidden" field if the given value is not nil.
func (su *SettingUpdate) SetNillableHidden(b *bool) *SettingUpdate {
	if b != nil {
		su.SetHidden(*b)
	}
	return su
}

// SetEditable sets the "editable" field.
func (su *SettingUpdate) SetEditable(b bool) *SettingUpdate {
	su.mutation.SetEditable(b)
	return su
}

// SetNillableEditable sets the "editable" field if the given value is not nil.
func (su *SettingUpdate) SetNillableEditable(b *bool) *SettingUpdate {
	if b != nil {
		su.SetEditable(*b)
	}
	return su
}

// SetPrivate sets the "private" field.
func (su *SettingUpdate) SetPrivate(b bool) *SettingUpdate {
	su.mutation.SetPrivate(b)
	return su
}

// SetNillablePrivate sets the "private" field if the given value is not nil.
func (su *SettingUpdate) SetNillablePrivate(b *bool) *SettingUpdate {
	if b != nil {
		su.SetPrivate(*b)
	}
	return su
}

// Mutation returns the SettingMutation object of the builder.
func (su *SettingUpdate) Mutation() *SettingMutation {
	return su.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (su *SettingUpdate) Save(ctx context.Context) (int, error) {
	if err := su.defaults(); err != nil {
		return 0, err
	}
	return withHooks[int, SettingMutation](ctx, su.sqlSave, su.mutation, su.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (su *SettingUpdate) SaveX(ctx context.Context) int {
	affected, err := su.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (su *SettingUpdate) Exec(ctx context.Context) error {
	_, err := su.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (su *SettingUpdate) ExecX(ctx context.Context) {
	if err := su.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (su *SettingUpdate) defaults() error {
	if _, ok := su.mutation.UpdateTime(); !ok {
		if setting.UpdateDefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized setting.UpdateDefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := setting.UpdateDefaultUpdateTime()
		su.mutation.SetUpdateTime(v)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (su *SettingUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *SettingUpdate {
	su.modifiers = append(su.modifiers, modifiers...)
	return su
}

func (su *SettingUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   setting.Table,
			Columns: setting.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeOther,
				Column: setting.FieldID,
			},
		},
	}
	if ps := su.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := su.mutation.UpdateTime(); ok {
		_spec.SetField(setting.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := su.mutation.Name(); ok {
		_spec.SetField(setting.FieldName, field.TypeString, value)
	}
	if value, ok := su.mutation.Value(); ok {
		_spec.SetField(setting.FieldValue, field.TypeString, value)
	}
	if value, ok := su.mutation.Hidden(); ok {
		_spec.SetField(setting.FieldHidden, field.TypeBool, value)
	}
	if value, ok := su.mutation.Editable(); ok {
		_spec.SetField(setting.FieldEditable, field.TypeBool, value)
	}
	if value, ok := su.mutation.Private(); ok {
		_spec.SetField(setting.FieldPrivate, field.TypeBool, value)
	}
	_spec.Node.Schema = su.schemaConfig.Setting
	ctx = internal.NewSchemaConfigContext(ctx, su.schemaConfig)
	_spec.AddModifiers(su.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, su.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{setting.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	su.mutation.done = true
	return n, nil
}

// SettingUpdateOne is the builder for updating a single Setting entity.
type SettingUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *SettingMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetUpdateTime sets the "updateTime" field.
func (suo *SettingUpdateOne) SetUpdateTime(t time.Time) *SettingUpdateOne {
	suo.mutation.SetUpdateTime(t)
	return suo
}

// SetName sets the "name" field.
func (suo *SettingUpdateOne) SetName(s string) *SettingUpdateOne {
	suo.mutation.SetName(s)
	return suo
}

// SetValue sets the "value" field.
func (suo *SettingUpdateOne) SetValue(s string) *SettingUpdateOne {
	suo.mutation.SetValue(s)
	return suo
}

// SetNillableValue sets the "value" field if the given value is not nil.
func (suo *SettingUpdateOne) SetNillableValue(s *string) *SettingUpdateOne {
	if s != nil {
		suo.SetValue(*s)
	}
	return suo
}

// SetHidden sets the "hidden" field.
func (suo *SettingUpdateOne) SetHidden(b bool) *SettingUpdateOne {
	suo.mutation.SetHidden(b)
	return suo
}

// SetNillableHidden sets the "hidden" field if the given value is not nil.
func (suo *SettingUpdateOne) SetNillableHidden(b *bool) *SettingUpdateOne {
	if b != nil {
		suo.SetHidden(*b)
	}
	return suo
}

// SetEditable sets the "editable" field.
func (suo *SettingUpdateOne) SetEditable(b bool) *SettingUpdateOne {
	suo.mutation.SetEditable(b)
	return suo
}

// SetNillableEditable sets the "editable" field if the given value is not nil.
func (suo *SettingUpdateOne) SetNillableEditable(b *bool) *SettingUpdateOne {
	if b != nil {
		suo.SetEditable(*b)
	}
	return suo
}

// SetPrivate sets the "private" field.
func (suo *SettingUpdateOne) SetPrivate(b bool) *SettingUpdateOne {
	suo.mutation.SetPrivate(b)
	return suo
}

// SetNillablePrivate sets the "private" field if the given value is not nil.
func (suo *SettingUpdateOne) SetNillablePrivate(b *bool) *SettingUpdateOne {
	if b != nil {
		suo.SetPrivate(*b)
	}
	return suo
}

// Mutation returns the SettingMutation object of the builder.
func (suo *SettingUpdateOne) Mutation() *SettingMutation {
	return suo.mutation
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (suo *SettingUpdateOne) Select(field string, fields ...string) *SettingUpdateOne {
	suo.fields = append([]string{field}, fields...)
	return suo
}

// Save executes the query and returns the updated Setting entity.
func (suo *SettingUpdateOne) Save(ctx context.Context) (*Setting, error) {
	if err := suo.defaults(); err != nil {
		return nil, err
	}
	return withHooks[*Setting, SettingMutation](ctx, suo.sqlSave, suo.mutation, suo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (suo *SettingUpdateOne) SaveX(ctx context.Context) *Setting {
	node, err := suo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (suo *SettingUpdateOne) Exec(ctx context.Context) error {
	_, err := suo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (suo *SettingUpdateOne) ExecX(ctx context.Context) {
	if err := suo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (suo *SettingUpdateOne) defaults() error {
	if _, ok := suo.mutation.UpdateTime(); !ok {
		if setting.UpdateDefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized setting.UpdateDefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := setting.UpdateDefaultUpdateTime()
		suo.mutation.SetUpdateTime(v)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (suo *SettingUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *SettingUpdateOne {
	suo.modifiers = append(suo.modifiers, modifiers...)
	return suo
}

func (suo *SettingUpdateOne) sqlSave(ctx context.Context) (_node *Setting, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   setting.Table,
			Columns: setting.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeOther,
				Column: setting.FieldID,
			},
		},
	}
	id, ok := suo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`model: missing "Setting.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := suo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, setting.FieldID)
		for _, f := range fields {
			if !setting.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("model: invalid field %q for query", f)}
			}
			if f != setting.FieldID {
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
		_spec.SetField(setting.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := suo.mutation.Name(); ok {
		_spec.SetField(setting.FieldName, field.TypeString, value)
	}
	if value, ok := suo.mutation.Value(); ok {
		_spec.SetField(setting.FieldValue, field.TypeString, value)
	}
	if value, ok := suo.mutation.Hidden(); ok {
		_spec.SetField(setting.FieldHidden, field.TypeBool, value)
	}
	if value, ok := suo.mutation.Editable(); ok {
		_spec.SetField(setting.FieldEditable, field.TypeBool, value)
	}
	if value, ok := suo.mutation.Private(); ok {
		_spec.SetField(setting.FieldPrivate, field.TypeBool, value)
	}
	_spec.Node.Schema = suo.schemaConfig.Setting
	ctx = internal.NewSchemaConfigContext(ctx, suo.schemaConfig)
	_spec.AddModifiers(suo.modifiers...)
	_node = &Setting{config: suo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, suo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{setting.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	suo.mutation.done = true
	return _node, nil
}
