// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package model

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/model/subject"
	"github.com/seal-io/seal/pkg/dao/oid"
	"github.com/seal-io/seal/pkg/dao/schema"
)

// SubjectCreate is the builder for creating a Subject entity.
type SubjectCreate struct {
	config
	mutation *SubjectMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreateTime sets the "createTime" field.
func (sc *SubjectCreate) SetCreateTime(t time.Time) *SubjectCreate {
	sc.mutation.SetCreateTime(t)
	return sc
}

// SetNillableCreateTime sets the "createTime" field if the given value is not nil.
func (sc *SubjectCreate) SetNillableCreateTime(t *time.Time) *SubjectCreate {
	if t != nil {
		sc.SetCreateTime(*t)
	}
	return sc
}

// SetUpdateTime sets the "updateTime" field.
func (sc *SubjectCreate) SetUpdateTime(t time.Time) *SubjectCreate {
	sc.mutation.SetUpdateTime(t)
	return sc
}

// SetNillableUpdateTime sets the "updateTime" field if the given value is not nil.
func (sc *SubjectCreate) SetNillableUpdateTime(t *time.Time) *SubjectCreate {
	if t != nil {
		sc.SetUpdateTime(*t)
	}
	return sc
}

// SetKind sets the "kind" field.
func (sc *SubjectCreate) SetKind(s string) *SubjectCreate {
	sc.mutation.SetKind(s)
	return sc
}

// SetNillableKind sets the "kind" field if the given value is not nil.
func (sc *SubjectCreate) SetNillableKind(s *string) *SubjectCreate {
	if s != nil {
		sc.SetKind(*s)
	}
	return sc
}

// SetGroup sets the "group" field.
func (sc *SubjectCreate) SetGroup(s string) *SubjectCreate {
	sc.mutation.SetGroup(s)
	return sc
}

// SetNillableGroup sets the "group" field if the given value is not nil.
func (sc *SubjectCreate) SetNillableGroup(s *string) *SubjectCreate {
	if s != nil {
		sc.SetGroup(*s)
	}
	return sc
}

// SetName sets the "name" field.
func (sc *SubjectCreate) SetName(s string) *SubjectCreate {
	sc.mutation.SetName(s)
	return sc
}

// SetDescription sets the "description" field.
func (sc *SubjectCreate) SetDescription(s string) *SubjectCreate {
	sc.mutation.SetDescription(s)
	return sc
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (sc *SubjectCreate) SetNillableDescription(s *string) *SubjectCreate {
	if s != nil {
		sc.SetDescription(*s)
	}
	return sc
}

// SetMountTo sets the "mountTo" field.
func (sc *SubjectCreate) SetMountTo(b bool) *SubjectCreate {
	sc.mutation.SetMountTo(b)
	return sc
}

// SetNillableMountTo sets the "mountTo" field if the given value is not nil.
func (sc *SubjectCreate) SetNillableMountTo(b *bool) *SubjectCreate {
	if b != nil {
		sc.SetMountTo(*b)
	}
	return sc
}

// SetLoginTo sets the "loginTo" field.
func (sc *SubjectCreate) SetLoginTo(b bool) *SubjectCreate {
	sc.mutation.SetLoginTo(b)
	return sc
}

// SetNillableLoginTo sets the "loginTo" field if the given value is not nil.
func (sc *SubjectCreate) SetNillableLoginTo(b *bool) *SubjectCreate {
	if b != nil {
		sc.SetLoginTo(*b)
	}
	return sc
}

// SetRoles sets the "roles" field.
func (sc *SubjectCreate) SetRoles(sr schema.SubjectRoles) *SubjectCreate {
	sc.mutation.SetRoles(sr)
	return sc
}

// SetPaths sets the "paths" field.
func (sc *SubjectCreate) SetPaths(s []string) *SubjectCreate {
	sc.mutation.SetPaths(s)
	return sc
}

// SetBuiltin sets the "builtin" field.
func (sc *SubjectCreate) SetBuiltin(b bool) *SubjectCreate {
	sc.mutation.SetBuiltin(b)
	return sc
}

// SetNillableBuiltin sets the "builtin" field if the given value is not nil.
func (sc *SubjectCreate) SetNillableBuiltin(b *bool) *SubjectCreate {
	if b != nil {
		sc.SetBuiltin(*b)
	}
	return sc
}

// SetID sets the "id" field.
func (sc *SubjectCreate) SetID(o oid.ID) *SubjectCreate {
	sc.mutation.SetID(o)
	return sc
}

// Mutation returns the SubjectMutation object of the builder.
func (sc *SubjectCreate) Mutation() *SubjectMutation {
	return sc.mutation
}

// Save creates the Subject in the database.
func (sc *SubjectCreate) Save(ctx context.Context) (*Subject, error) {
	if err := sc.defaults(); err != nil {
		return nil, err
	}
	return withHooks[*Subject, SubjectMutation](ctx, sc.sqlSave, sc.mutation, sc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (sc *SubjectCreate) SaveX(ctx context.Context) *Subject {
	v, err := sc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (sc *SubjectCreate) Exec(ctx context.Context) error {
	_, err := sc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (sc *SubjectCreate) ExecX(ctx context.Context) {
	if err := sc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (sc *SubjectCreate) defaults() error {
	if _, ok := sc.mutation.CreateTime(); !ok {
		if subject.DefaultCreateTime == nil {
			return fmt.Errorf("model: uninitialized subject.DefaultCreateTime (forgotten import model/runtime?)")
		}
		v := subject.DefaultCreateTime()
		sc.mutation.SetCreateTime(v)
	}
	if _, ok := sc.mutation.UpdateTime(); !ok {
		if subject.DefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized subject.DefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := subject.DefaultUpdateTime()
		sc.mutation.SetUpdateTime(v)
	}
	if _, ok := sc.mutation.Kind(); !ok {
		v := subject.DefaultKind
		sc.mutation.SetKind(v)
	}
	if _, ok := sc.mutation.Group(); !ok {
		v := subject.DefaultGroup
		sc.mutation.SetGroup(v)
	}
	if _, ok := sc.mutation.Description(); !ok {
		v := subject.DefaultDescription
		sc.mutation.SetDescription(v)
	}
	if _, ok := sc.mutation.MountTo(); !ok {
		v := subject.DefaultMountTo
		sc.mutation.SetMountTo(v)
	}
	if _, ok := sc.mutation.LoginTo(); !ok {
		v := subject.DefaultLoginTo
		sc.mutation.SetLoginTo(v)
	}
	if _, ok := sc.mutation.Roles(); !ok {
		v := subject.DefaultRoles
		sc.mutation.SetRoles(v)
	}
	if _, ok := sc.mutation.Paths(); !ok {
		v := subject.DefaultPaths
		sc.mutation.SetPaths(v)
	}
	if _, ok := sc.mutation.Builtin(); !ok {
		v := subject.DefaultBuiltin
		sc.mutation.SetBuiltin(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (sc *SubjectCreate) check() error {
	if _, ok := sc.mutation.CreateTime(); !ok {
		return &ValidationError{Name: "createTime", err: errors.New(`model: missing required field "Subject.createTime"`)}
	}
	if _, ok := sc.mutation.UpdateTime(); !ok {
		return &ValidationError{Name: "updateTime", err: errors.New(`model: missing required field "Subject.updateTime"`)}
	}
	if _, ok := sc.mutation.Kind(); !ok {
		return &ValidationError{Name: "kind", err: errors.New(`model: missing required field "Subject.kind"`)}
	}
	if _, ok := sc.mutation.Group(); !ok {
		return &ValidationError{Name: "group", err: errors.New(`model: missing required field "Subject.group"`)}
	}
	if _, ok := sc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`model: missing required field "Subject.name"`)}
	}
	if v, ok := sc.mutation.Name(); ok {
		if err := subject.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`model: validator failed for field "Subject.name": %w`, err)}
		}
	}
	if _, ok := sc.mutation.Description(); !ok {
		return &ValidationError{Name: "description", err: errors.New(`model: missing required field "Subject.description"`)}
	}
	if _, ok := sc.mutation.MountTo(); !ok {
		return &ValidationError{Name: "mountTo", err: errors.New(`model: missing required field "Subject.mountTo"`)}
	}
	if _, ok := sc.mutation.LoginTo(); !ok {
		return &ValidationError{Name: "loginTo", err: errors.New(`model: missing required field "Subject.loginTo"`)}
	}
	if _, ok := sc.mutation.Roles(); !ok {
		return &ValidationError{Name: "roles", err: errors.New(`model: missing required field "Subject.roles"`)}
	}
	if _, ok := sc.mutation.Paths(); !ok {
		return &ValidationError{Name: "paths", err: errors.New(`model: missing required field "Subject.paths"`)}
	}
	if _, ok := sc.mutation.Builtin(); !ok {
		return &ValidationError{Name: "builtin", err: errors.New(`model: missing required field "Subject.builtin"`)}
	}
	return nil
}

func (sc *SubjectCreate) sqlSave(ctx context.Context) (*Subject, error) {
	if err := sc.check(); err != nil {
		return nil, err
	}
	_node, _spec := sc.createSpec()
	if err := sqlgraph.CreateNode(ctx, sc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*oid.ID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	sc.mutation.id = &_node.ID
	sc.mutation.done = true
	return _node, nil
}

func (sc *SubjectCreate) createSpec() (*Subject, *sqlgraph.CreateSpec) {
	var (
		_node = &Subject{config: sc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: subject.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeOther,
				Column: subject.FieldID,
			},
		}
	)
	_spec.Schema = sc.schemaConfig.Subject
	_spec.OnConflict = sc.conflict
	if id, ok := sc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := sc.mutation.CreateTime(); ok {
		_spec.SetField(subject.FieldCreateTime, field.TypeTime, value)
		_node.CreateTime = &value
	}
	if value, ok := sc.mutation.UpdateTime(); ok {
		_spec.SetField(subject.FieldUpdateTime, field.TypeTime, value)
		_node.UpdateTime = &value
	}
	if value, ok := sc.mutation.Kind(); ok {
		_spec.SetField(subject.FieldKind, field.TypeString, value)
		_node.Kind = value
	}
	if value, ok := sc.mutation.Group(); ok {
		_spec.SetField(subject.FieldGroup, field.TypeString, value)
		_node.Group = value
	}
	if value, ok := sc.mutation.Name(); ok {
		_spec.SetField(subject.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := sc.mutation.Description(); ok {
		_spec.SetField(subject.FieldDescription, field.TypeString, value)
		_node.Description = value
	}
	if value, ok := sc.mutation.MountTo(); ok {
		_spec.SetField(subject.FieldMountTo, field.TypeBool, value)
		_node.MountTo = &value
	}
	if value, ok := sc.mutation.LoginTo(); ok {
		_spec.SetField(subject.FieldLoginTo, field.TypeBool, value)
		_node.LoginTo = &value
	}
	if value, ok := sc.mutation.Roles(); ok {
		_spec.SetField(subject.FieldRoles, field.TypeJSON, value)
		_node.Roles = value
	}
	if value, ok := sc.mutation.Paths(); ok {
		_spec.SetField(subject.FieldPaths, field.TypeJSON, value)
		_node.Paths = value
	}
	if value, ok := sc.mutation.Builtin(); ok {
		_spec.SetField(subject.FieldBuiltin, field.TypeBool, value)
		_node.Builtin = value
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Subject.Create().
//		SetCreateTime(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.SubjectUpsert) {
//			SetCreateTime(v+v).
//		}).
//		Exec(ctx)
func (sc *SubjectCreate) OnConflict(opts ...sql.ConflictOption) *SubjectUpsertOne {
	sc.conflict = opts
	return &SubjectUpsertOne{
		create: sc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Subject.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (sc *SubjectCreate) OnConflictColumns(columns ...string) *SubjectUpsertOne {
	sc.conflict = append(sc.conflict, sql.ConflictColumns(columns...))
	return &SubjectUpsertOne{
		create: sc,
	}
}

type (
	// SubjectUpsertOne is the builder for "upsert"-ing
	//  one Subject node.
	SubjectUpsertOne struct {
		create *SubjectCreate
	}

	// SubjectUpsert is the "OnConflict" setter.
	SubjectUpsert struct {
		*sql.UpdateSet
	}
)

// SetUpdateTime sets the "updateTime" field.
func (u *SubjectUpsert) SetUpdateTime(v time.Time) *SubjectUpsert {
	u.Set(subject.FieldUpdateTime, v)
	return u
}

// UpdateUpdateTime sets the "updateTime" field to the value that was provided on create.
func (u *SubjectUpsert) UpdateUpdateTime() *SubjectUpsert {
	u.SetExcluded(subject.FieldUpdateTime)
	return u
}

// SetGroup sets the "group" field.
func (u *SubjectUpsert) SetGroup(v string) *SubjectUpsert {
	u.Set(subject.FieldGroup, v)
	return u
}

// UpdateGroup sets the "group" field to the value that was provided on create.
func (u *SubjectUpsert) UpdateGroup() *SubjectUpsert {
	u.SetExcluded(subject.FieldGroup)
	return u
}

// SetDescription sets the "description" field.
func (u *SubjectUpsert) SetDescription(v string) *SubjectUpsert {
	u.Set(subject.FieldDescription, v)
	return u
}

// UpdateDescription sets the "description" field to the value that was provided on create.
func (u *SubjectUpsert) UpdateDescription() *SubjectUpsert {
	u.SetExcluded(subject.FieldDescription)
	return u
}

// SetMountTo sets the "mountTo" field.
func (u *SubjectUpsert) SetMountTo(v bool) *SubjectUpsert {
	u.Set(subject.FieldMountTo, v)
	return u
}

// UpdateMountTo sets the "mountTo" field to the value that was provided on create.
func (u *SubjectUpsert) UpdateMountTo() *SubjectUpsert {
	u.SetExcluded(subject.FieldMountTo)
	return u
}

// SetLoginTo sets the "loginTo" field.
func (u *SubjectUpsert) SetLoginTo(v bool) *SubjectUpsert {
	u.Set(subject.FieldLoginTo, v)
	return u
}

// UpdateLoginTo sets the "loginTo" field to the value that was provided on create.
func (u *SubjectUpsert) UpdateLoginTo() *SubjectUpsert {
	u.SetExcluded(subject.FieldLoginTo)
	return u
}

// SetRoles sets the "roles" field.
func (u *SubjectUpsert) SetRoles(v schema.SubjectRoles) *SubjectUpsert {
	u.Set(subject.FieldRoles, v)
	return u
}

// UpdateRoles sets the "roles" field to the value that was provided on create.
func (u *SubjectUpsert) UpdateRoles() *SubjectUpsert {
	u.SetExcluded(subject.FieldRoles)
	return u
}

// SetPaths sets the "paths" field.
func (u *SubjectUpsert) SetPaths(v []string) *SubjectUpsert {
	u.Set(subject.FieldPaths, v)
	return u
}

// UpdatePaths sets the "paths" field to the value that was provided on create.
func (u *SubjectUpsert) UpdatePaths() *SubjectUpsert {
	u.SetExcluded(subject.FieldPaths)
	return u
}

// SetBuiltin sets the "builtin" field.
func (u *SubjectUpsert) SetBuiltin(v bool) *SubjectUpsert {
	u.Set(subject.FieldBuiltin, v)
	return u
}

// UpdateBuiltin sets the "builtin" field to the value that was provided on create.
func (u *SubjectUpsert) UpdateBuiltin() *SubjectUpsert {
	u.SetExcluded(subject.FieldBuiltin)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.Subject.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(subject.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *SubjectUpsertOne) UpdateNewValues() *SubjectUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(subject.FieldID)
		}
		if _, exists := u.create.mutation.CreateTime(); exists {
			s.SetIgnore(subject.FieldCreateTime)
		}
		if _, exists := u.create.mutation.Kind(); exists {
			s.SetIgnore(subject.FieldKind)
		}
		if _, exists := u.create.mutation.Name(); exists {
			s.SetIgnore(subject.FieldName)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Subject.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *SubjectUpsertOne) Ignore() *SubjectUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *SubjectUpsertOne) DoNothing() *SubjectUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the SubjectCreate.OnConflict
// documentation for more info.
func (u *SubjectUpsertOne) Update(set func(*SubjectUpsert)) *SubjectUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&SubjectUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdateTime sets the "updateTime" field.
func (u *SubjectUpsertOne) SetUpdateTime(v time.Time) *SubjectUpsertOne {
	return u.Update(func(s *SubjectUpsert) {
		s.SetUpdateTime(v)
	})
}

// UpdateUpdateTime sets the "updateTime" field to the value that was provided on create.
func (u *SubjectUpsertOne) UpdateUpdateTime() *SubjectUpsertOne {
	return u.Update(func(s *SubjectUpsert) {
		s.UpdateUpdateTime()
	})
}

// SetGroup sets the "group" field.
func (u *SubjectUpsertOne) SetGroup(v string) *SubjectUpsertOne {
	return u.Update(func(s *SubjectUpsert) {
		s.SetGroup(v)
	})
}

// UpdateGroup sets the "group" field to the value that was provided on create.
func (u *SubjectUpsertOne) UpdateGroup() *SubjectUpsertOne {
	return u.Update(func(s *SubjectUpsert) {
		s.UpdateGroup()
	})
}

// SetDescription sets the "description" field.
func (u *SubjectUpsertOne) SetDescription(v string) *SubjectUpsertOne {
	return u.Update(func(s *SubjectUpsert) {
		s.SetDescription(v)
	})
}

// UpdateDescription sets the "description" field to the value that was provided on create.
func (u *SubjectUpsertOne) UpdateDescription() *SubjectUpsertOne {
	return u.Update(func(s *SubjectUpsert) {
		s.UpdateDescription()
	})
}

// SetMountTo sets the "mountTo" field.
func (u *SubjectUpsertOne) SetMountTo(v bool) *SubjectUpsertOne {
	return u.Update(func(s *SubjectUpsert) {
		s.SetMountTo(v)
	})
}

// UpdateMountTo sets the "mountTo" field to the value that was provided on create.
func (u *SubjectUpsertOne) UpdateMountTo() *SubjectUpsertOne {
	return u.Update(func(s *SubjectUpsert) {
		s.UpdateMountTo()
	})
}

// SetLoginTo sets the "loginTo" field.
func (u *SubjectUpsertOne) SetLoginTo(v bool) *SubjectUpsertOne {
	return u.Update(func(s *SubjectUpsert) {
		s.SetLoginTo(v)
	})
}

// UpdateLoginTo sets the "loginTo" field to the value that was provided on create.
func (u *SubjectUpsertOne) UpdateLoginTo() *SubjectUpsertOne {
	return u.Update(func(s *SubjectUpsert) {
		s.UpdateLoginTo()
	})
}

// SetRoles sets the "roles" field.
func (u *SubjectUpsertOne) SetRoles(v schema.SubjectRoles) *SubjectUpsertOne {
	return u.Update(func(s *SubjectUpsert) {
		s.SetRoles(v)
	})
}

// UpdateRoles sets the "roles" field to the value that was provided on create.
func (u *SubjectUpsertOne) UpdateRoles() *SubjectUpsertOne {
	return u.Update(func(s *SubjectUpsert) {
		s.UpdateRoles()
	})
}

// SetPaths sets the "paths" field.
func (u *SubjectUpsertOne) SetPaths(v []string) *SubjectUpsertOne {
	return u.Update(func(s *SubjectUpsert) {
		s.SetPaths(v)
	})
}

// UpdatePaths sets the "paths" field to the value that was provided on create.
func (u *SubjectUpsertOne) UpdatePaths() *SubjectUpsertOne {
	return u.Update(func(s *SubjectUpsert) {
		s.UpdatePaths()
	})
}

// SetBuiltin sets the "builtin" field.
func (u *SubjectUpsertOne) SetBuiltin(v bool) *SubjectUpsertOne {
	return u.Update(func(s *SubjectUpsert) {
		s.SetBuiltin(v)
	})
}

// UpdateBuiltin sets the "builtin" field to the value that was provided on create.
func (u *SubjectUpsertOne) UpdateBuiltin() *SubjectUpsertOne {
	return u.Update(func(s *SubjectUpsert) {
		s.UpdateBuiltin()
	})
}

// Exec executes the query.
func (u *SubjectUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("model: missing options for SubjectCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *SubjectUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *SubjectUpsertOne) ID(ctx context.Context) (id oid.ID, err error) {
	if u.create.driver.Dialect() == dialect.MySQL {
		// In case of "ON CONFLICT", there is no way to get back non-numeric ID
		// fields from the database since MySQL does not support the RETURNING clause.
		return id, errors.New("model: SubjectUpsertOne.ID is not supported by MySQL driver. Use SubjectUpsertOne.Exec instead")
	}
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *SubjectUpsertOne) IDX(ctx context.Context) oid.ID {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// SubjectCreateBulk is the builder for creating many Subject entities in bulk.
type SubjectCreateBulk struct {
	config
	builders []*SubjectCreate
	conflict []sql.ConflictOption
}

// Save creates the Subject entities in the database.
func (scb *SubjectCreateBulk) Save(ctx context.Context) ([]*Subject, error) {
	specs := make([]*sqlgraph.CreateSpec, len(scb.builders))
	nodes := make([]*Subject, len(scb.builders))
	mutators := make([]Mutator, len(scb.builders))
	for i := range scb.builders {
		func(i int, root context.Context) {
			builder := scb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*SubjectMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, scb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = scb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, scb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, scb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (scb *SubjectCreateBulk) SaveX(ctx context.Context) []*Subject {
	v, err := scb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (scb *SubjectCreateBulk) Exec(ctx context.Context) error {
	_, err := scb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (scb *SubjectCreateBulk) ExecX(ctx context.Context) {
	if err := scb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Subject.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.SubjectUpsert) {
//			SetCreateTime(v+v).
//		}).
//		Exec(ctx)
func (scb *SubjectCreateBulk) OnConflict(opts ...sql.ConflictOption) *SubjectUpsertBulk {
	scb.conflict = opts
	return &SubjectUpsertBulk{
		create: scb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Subject.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (scb *SubjectCreateBulk) OnConflictColumns(columns ...string) *SubjectUpsertBulk {
	scb.conflict = append(scb.conflict, sql.ConflictColumns(columns...))
	return &SubjectUpsertBulk{
		create: scb,
	}
}

// SubjectUpsertBulk is the builder for "upsert"-ing
// a bulk of Subject nodes.
type SubjectUpsertBulk struct {
	create *SubjectCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Subject.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(subject.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *SubjectUpsertBulk) UpdateNewValues() *SubjectUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(subject.FieldID)
			}
			if _, exists := b.mutation.CreateTime(); exists {
				s.SetIgnore(subject.FieldCreateTime)
			}
			if _, exists := b.mutation.Kind(); exists {
				s.SetIgnore(subject.FieldKind)
			}
			if _, exists := b.mutation.Name(); exists {
				s.SetIgnore(subject.FieldName)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Subject.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *SubjectUpsertBulk) Ignore() *SubjectUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *SubjectUpsertBulk) DoNothing() *SubjectUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the SubjectCreateBulk.OnConflict
// documentation for more info.
func (u *SubjectUpsertBulk) Update(set func(*SubjectUpsert)) *SubjectUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&SubjectUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdateTime sets the "updateTime" field.
func (u *SubjectUpsertBulk) SetUpdateTime(v time.Time) *SubjectUpsertBulk {
	return u.Update(func(s *SubjectUpsert) {
		s.SetUpdateTime(v)
	})
}

// UpdateUpdateTime sets the "updateTime" field to the value that was provided on create.
func (u *SubjectUpsertBulk) UpdateUpdateTime() *SubjectUpsertBulk {
	return u.Update(func(s *SubjectUpsert) {
		s.UpdateUpdateTime()
	})
}

// SetGroup sets the "group" field.
func (u *SubjectUpsertBulk) SetGroup(v string) *SubjectUpsertBulk {
	return u.Update(func(s *SubjectUpsert) {
		s.SetGroup(v)
	})
}

// UpdateGroup sets the "group" field to the value that was provided on create.
func (u *SubjectUpsertBulk) UpdateGroup() *SubjectUpsertBulk {
	return u.Update(func(s *SubjectUpsert) {
		s.UpdateGroup()
	})
}

// SetDescription sets the "description" field.
func (u *SubjectUpsertBulk) SetDescription(v string) *SubjectUpsertBulk {
	return u.Update(func(s *SubjectUpsert) {
		s.SetDescription(v)
	})
}

// UpdateDescription sets the "description" field to the value that was provided on create.
func (u *SubjectUpsertBulk) UpdateDescription() *SubjectUpsertBulk {
	return u.Update(func(s *SubjectUpsert) {
		s.UpdateDescription()
	})
}

// SetMountTo sets the "mountTo" field.
func (u *SubjectUpsertBulk) SetMountTo(v bool) *SubjectUpsertBulk {
	return u.Update(func(s *SubjectUpsert) {
		s.SetMountTo(v)
	})
}

// UpdateMountTo sets the "mountTo" field to the value that was provided on create.
func (u *SubjectUpsertBulk) UpdateMountTo() *SubjectUpsertBulk {
	return u.Update(func(s *SubjectUpsert) {
		s.UpdateMountTo()
	})
}

// SetLoginTo sets the "loginTo" field.
func (u *SubjectUpsertBulk) SetLoginTo(v bool) *SubjectUpsertBulk {
	return u.Update(func(s *SubjectUpsert) {
		s.SetLoginTo(v)
	})
}

// UpdateLoginTo sets the "loginTo" field to the value that was provided on create.
func (u *SubjectUpsertBulk) UpdateLoginTo() *SubjectUpsertBulk {
	return u.Update(func(s *SubjectUpsert) {
		s.UpdateLoginTo()
	})
}

// SetRoles sets the "roles" field.
func (u *SubjectUpsertBulk) SetRoles(v schema.SubjectRoles) *SubjectUpsertBulk {
	return u.Update(func(s *SubjectUpsert) {
		s.SetRoles(v)
	})
}

// UpdateRoles sets the "roles" field to the value that was provided on create.
func (u *SubjectUpsertBulk) UpdateRoles() *SubjectUpsertBulk {
	return u.Update(func(s *SubjectUpsert) {
		s.UpdateRoles()
	})
}

// SetPaths sets the "paths" field.
func (u *SubjectUpsertBulk) SetPaths(v []string) *SubjectUpsertBulk {
	return u.Update(func(s *SubjectUpsert) {
		s.SetPaths(v)
	})
}

// UpdatePaths sets the "paths" field to the value that was provided on create.
func (u *SubjectUpsertBulk) UpdatePaths() *SubjectUpsertBulk {
	return u.Update(func(s *SubjectUpsert) {
		s.UpdatePaths()
	})
}

// SetBuiltin sets the "builtin" field.
func (u *SubjectUpsertBulk) SetBuiltin(v bool) *SubjectUpsertBulk {
	return u.Update(func(s *SubjectUpsert) {
		s.SetBuiltin(v)
	})
}

// UpdateBuiltin sets the "builtin" field to the value that was provided on create.
func (u *SubjectUpsertBulk) UpdateBuiltin() *SubjectUpsertBulk {
	return u.Update(func(s *SubjectUpsert) {
		s.UpdateBuiltin()
	})
}

// Exec executes the query.
func (u *SubjectUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("model: OnConflict was set for builder %d. Set it on the SubjectCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("model: missing options for SubjectCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *SubjectUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
