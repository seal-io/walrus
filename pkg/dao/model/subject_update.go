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
	"github.com/seal-io/seal/pkg/dao/model/subject"
	"github.com/seal-io/seal/pkg/dao/types"
)

// SubjectUpdate is the builder for updating Subject entities.
type SubjectUpdate struct {
	config
	hooks     []Hook
	mutation  *SubjectMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the SubjectUpdate builder.
func (su *SubjectUpdate) Where(ps ...predicate.Subject) *SubjectUpdate {
	su.mutation.Where(ps...)
	return su
}

// SetUpdateTime sets the "updateTime" field.
func (su *SubjectUpdate) SetUpdateTime(t time.Time) *SubjectUpdate {
	su.mutation.SetUpdateTime(t)
	return su
}

// SetGroup sets the "group" field.
func (su *SubjectUpdate) SetGroup(s string) *SubjectUpdate {
	su.mutation.SetGroup(s)
	return su
}

// SetNillableGroup sets the "group" field if the given value is not nil.
func (su *SubjectUpdate) SetNillableGroup(s *string) *SubjectUpdate {
	if s != nil {
		su.SetGroup(*s)
	}
	return su
}

// SetDescription sets the "description" field.
func (su *SubjectUpdate) SetDescription(s string) *SubjectUpdate {
	su.mutation.SetDescription(s)
	return su
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (su *SubjectUpdate) SetNillableDescription(s *string) *SubjectUpdate {
	if s != nil {
		su.SetDescription(*s)
	}
	return su
}

// ClearDescription clears the value of the "description" field.
func (su *SubjectUpdate) ClearDescription() *SubjectUpdate {
	su.mutation.ClearDescription()
	return su
}

// SetMountTo sets the "mountTo" field.
func (su *SubjectUpdate) SetMountTo(b bool) *SubjectUpdate {
	su.mutation.SetMountTo(b)
	return su
}

// SetNillableMountTo sets the "mountTo" field if the given value is not nil.
func (su *SubjectUpdate) SetNillableMountTo(b *bool) *SubjectUpdate {
	if b != nil {
		su.SetMountTo(*b)
	}
	return su
}

// SetLoginTo sets the "loginTo" field.
func (su *SubjectUpdate) SetLoginTo(b bool) *SubjectUpdate {
	su.mutation.SetLoginTo(b)
	return su
}

// SetNillableLoginTo sets the "loginTo" field if the given value is not nil.
func (su *SubjectUpdate) SetNillableLoginTo(b *bool) *SubjectUpdate {
	if b != nil {
		su.SetLoginTo(*b)
	}
	return su
}

// SetRoles sets the "roles" field.
func (su *SubjectUpdate) SetRoles(tr types.SubjectRoles) *SubjectUpdate {
	su.mutation.SetRoles(tr)
	return su
}

// AppendRoles appends tr to the "roles" field.
func (su *SubjectUpdate) AppendRoles(tr types.SubjectRoles) *SubjectUpdate {
	su.mutation.AppendRoles(tr)
	return su
}

// SetPaths sets the "paths" field.
func (su *SubjectUpdate) SetPaths(s []string) *SubjectUpdate {
	su.mutation.SetPaths(s)
	return su
}

// AppendPaths appends s to the "paths" field.
func (su *SubjectUpdate) AppendPaths(s []string) *SubjectUpdate {
	su.mutation.AppendPaths(s)
	return su
}

// SetBuiltin sets the "builtin" field.
func (su *SubjectUpdate) SetBuiltin(b bool) *SubjectUpdate {
	su.mutation.SetBuiltin(b)
	return su
}

// SetNillableBuiltin sets the "builtin" field if the given value is not nil.
func (su *SubjectUpdate) SetNillableBuiltin(b *bool) *SubjectUpdate {
	if b != nil {
		su.SetBuiltin(*b)
	}
	return su
}

// Mutation returns the SubjectMutation object of the builder.
func (su *SubjectUpdate) Mutation() *SubjectMutation {
	return su.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (su *SubjectUpdate) Save(ctx context.Context) (int, error) {
	if err := su.defaults(); err != nil {
		return 0, err
	}
	return withHooks[int, SubjectMutation](ctx, su.sqlSave, su.mutation, su.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (su *SubjectUpdate) SaveX(ctx context.Context) int {
	affected, err := su.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (su *SubjectUpdate) Exec(ctx context.Context) error {
	_, err := su.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (su *SubjectUpdate) ExecX(ctx context.Context) {
	if err := su.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (su *SubjectUpdate) defaults() error {
	if _, ok := su.mutation.UpdateTime(); !ok {
		if subject.UpdateDefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized subject.UpdateDefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := subject.UpdateDefaultUpdateTime()
		su.mutation.SetUpdateTime(v)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (su *SubjectUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *SubjectUpdate {
	su.modifiers = append(su.modifiers, modifiers...)
	return su
}

func (su *SubjectUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(subject.Table, subject.Columns, sqlgraph.NewFieldSpec(subject.FieldID, field.TypeString))
	if ps := su.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := su.mutation.UpdateTime(); ok {
		_spec.SetField(subject.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := su.mutation.Group(); ok {
		_spec.SetField(subject.FieldGroup, field.TypeString, value)
	}
	if value, ok := su.mutation.Description(); ok {
		_spec.SetField(subject.FieldDescription, field.TypeString, value)
	}
	if su.mutation.DescriptionCleared() {
		_spec.ClearField(subject.FieldDescription, field.TypeString)
	}
	if value, ok := su.mutation.MountTo(); ok {
		_spec.SetField(subject.FieldMountTo, field.TypeBool, value)
	}
	if value, ok := su.mutation.LoginTo(); ok {
		_spec.SetField(subject.FieldLoginTo, field.TypeBool, value)
	}
	if value, ok := su.mutation.Roles(); ok {
		_spec.SetField(subject.FieldRoles, field.TypeJSON, value)
	}
	if value, ok := su.mutation.AppendedRoles(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, subject.FieldRoles, value)
		})
	}
	if value, ok := su.mutation.Paths(); ok {
		_spec.SetField(subject.FieldPaths, field.TypeJSON, value)
	}
	if value, ok := su.mutation.AppendedPaths(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, subject.FieldPaths, value)
		})
	}
	if value, ok := su.mutation.Builtin(); ok {
		_spec.SetField(subject.FieldBuiltin, field.TypeBool, value)
	}
	_spec.Node.Schema = su.schemaConfig.Subject
	ctx = internal.NewSchemaConfigContext(ctx, su.schemaConfig)
	_spec.AddModifiers(su.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, su.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{subject.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	su.mutation.done = true
	return n, nil
}

// SubjectUpdateOne is the builder for updating a single Subject entity.
type SubjectUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *SubjectMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetUpdateTime sets the "updateTime" field.
func (suo *SubjectUpdateOne) SetUpdateTime(t time.Time) *SubjectUpdateOne {
	suo.mutation.SetUpdateTime(t)
	return suo
}

// SetGroup sets the "group" field.
func (suo *SubjectUpdateOne) SetGroup(s string) *SubjectUpdateOne {
	suo.mutation.SetGroup(s)
	return suo
}

// SetNillableGroup sets the "group" field if the given value is not nil.
func (suo *SubjectUpdateOne) SetNillableGroup(s *string) *SubjectUpdateOne {
	if s != nil {
		suo.SetGroup(*s)
	}
	return suo
}

// SetDescription sets the "description" field.
func (suo *SubjectUpdateOne) SetDescription(s string) *SubjectUpdateOne {
	suo.mutation.SetDescription(s)
	return suo
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (suo *SubjectUpdateOne) SetNillableDescription(s *string) *SubjectUpdateOne {
	if s != nil {
		suo.SetDescription(*s)
	}
	return suo
}

// ClearDescription clears the value of the "description" field.
func (suo *SubjectUpdateOne) ClearDescription() *SubjectUpdateOne {
	suo.mutation.ClearDescription()
	return suo
}

// SetMountTo sets the "mountTo" field.
func (suo *SubjectUpdateOne) SetMountTo(b bool) *SubjectUpdateOne {
	suo.mutation.SetMountTo(b)
	return suo
}

// SetNillableMountTo sets the "mountTo" field if the given value is not nil.
func (suo *SubjectUpdateOne) SetNillableMountTo(b *bool) *SubjectUpdateOne {
	if b != nil {
		suo.SetMountTo(*b)
	}
	return suo
}

// SetLoginTo sets the "loginTo" field.
func (suo *SubjectUpdateOne) SetLoginTo(b bool) *SubjectUpdateOne {
	suo.mutation.SetLoginTo(b)
	return suo
}

// SetNillableLoginTo sets the "loginTo" field if the given value is not nil.
func (suo *SubjectUpdateOne) SetNillableLoginTo(b *bool) *SubjectUpdateOne {
	if b != nil {
		suo.SetLoginTo(*b)
	}
	return suo
}

// SetRoles sets the "roles" field.
func (suo *SubjectUpdateOne) SetRoles(tr types.SubjectRoles) *SubjectUpdateOne {
	suo.mutation.SetRoles(tr)
	return suo
}

// AppendRoles appends tr to the "roles" field.
func (suo *SubjectUpdateOne) AppendRoles(tr types.SubjectRoles) *SubjectUpdateOne {
	suo.mutation.AppendRoles(tr)
	return suo
}

// SetPaths sets the "paths" field.
func (suo *SubjectUpdateOne) SetPaths(s []string) *SubjectUpdateOne {
	suo.mutation.SetPaths(s)
	return suo
}

// AppendPaths appends s to the "paths" field.
func (suo *SubjectUpdateOne) AppendPaths(s []string) *SubjectUpdateOne {
	suo.mutation.AppendPaths(s)
	return suo
}

// SetBuiltin sets the "builtin" field.
func (suo *SubjectUpdateOne) SetBuiltin(b bool) *SubjectUpdateOne {
	suo.mutation.SetBuiltin(b)
	return suo
}

// SetNillableBuiltin sets the "builtin" field if the given value is not nil.
func (suo *SubjectUpdateOne) SetNillableBuiltin(b *bool) *SubjectUpdateOne {
	if b != nil {
		suo.SetBuiltin(*b)
	}
	return suo
}

// Mutation returns the SubjectMutation object of the builder.
func (suo *SubjectUpdateOne) Mutation() *SubjectMutation {
	return suo.mutation
}

// Where appends a list predicates to the SubjectUpdate builder.
func (suo *SubjectUpdateOne) Where(ps ...predicate.Subject) *SubjectUpdateOne {
	suo.mutation.Where(ps...)
	return suo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (suo *SubjectUpdateOne) Select(field string, fields ...string) *SubjectUpdateOne {
	suo.fields = append([]string{field}, fields...)
	return suo
}

// Save executes the query and returns the updated Subject entity.
func (suo *SubjectUpdateOne) Save(ctx context.Context) (*Subject, error) {
	if err := suo.defaults(); err != nil {
		return nil, err
	}
	return withHooks[*Subject, SubjectMutation](ctx, suo.sqlSave, suo.mutation, suo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (suo *SubjectUpdateOne) SaveX(ctx context.Context) *Subject {
	node, err := suo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (suo *SubjectUpdateOne) Exec(ctx context.Context) error {
	_, err := suo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (suo *SubjectUpdateOne) ExecX(ctx context.Context) {
	if err := suo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (suo *SubjectUpdateOne) defaults() error {
	if _, ok := suo.mutation.UpdateTime(); !ok {
		if subject.UpdateDefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized subject.UpdateDefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := subject.UpdateDefaultUpdateTime()
		suo.mutation.SetUpdateTime(v)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (suo *SubjectUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *SubjectUpdateOne {
	suo.modifiers = append(suo.modifiers, modifiers...)
	return suo
}

func (suo *SubjectUpdateOne) sqlSave(ctx context.Context) (_node *Subject, err error) {
	_spec := sqlgraph.NewUpdateSpec(subject.Table, subject.Columns, sqlgraph.NewFieldSpec(subject.FieldID, field.TypeString))
	id, ok := suo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`model: missing "Subject.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := suo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, subject.FieldID)
		for _, f := range fields {
			if !subject.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("model: invalid field %q for query", f)}
			}
			if f != subject.FieldID {
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
		_spec.SetField(subject.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := suo.mutation.Group(); ok {
		_spec.SetField(subject.FieldGroup, field.TypeString, value)
	}
	if value, ok := suo.mutation.Description(); ok {
		_spec.SetField(subject.FieldDescription, field.TypeString, value)
	}
	if suo.mutation.DescriptionCleared() {
		_spec.ClearField(subject.FieldDescription, field.TypeString)
	}
	if value, ok := suo.mutation.MountTo(); ok {
		_spec.SetField(subject.FieldMountTo, field.TypeBool, value)
	}
	if value, ok := suo.mutation.LoginTo(); ok {
		_spec.SetField(subject.FieldLoginTo, field.TypeBool, value)
	}
	if value, ok := suo.mutation.Roles(); ok {
		_spec.SetField(subject.FieldRoles, field.TypeJSON, value)
	}
	if value, ok := suo.mutation.AppendedRoles(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, subject.FieldRoles, value)
		})
	}
	if value, ok := suo.mutation.Paths(); ok {
		_spec.SetField(subject.FieldPaths, field.TypeJSON, value)
	}
	if value, ok := suo.mutation.AppendedPaths(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, subject.FieldPaths, value)
		})
	}
	if value, ok := suo.mutation.Builtin(); ok {
		_spec.SetField(subject.FieldBuiltin, field.TypeBool, value)
	}
	_spec.Node.Schema = suo.schemaConfig.Subject
	ctx = internal.NewSchemaConfigContext(ctx, suo.schemaConfig)
	_spec.AddModifiers(suo.modifiers...)
	_node = &Subject{config: suo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, suo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{subject.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	suo.mutation.done = true
	return _node, nil
}
