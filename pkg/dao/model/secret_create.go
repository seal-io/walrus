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

	"github.com/seal-io/seal/pkg/dao/model/project"
	"github.com/seal-io/seal/pkg/dao/model/secret"
	"github.com/seal-io/seal/pkg/dao/types/crypto"
	"github.com/seal-io/seal/pkg/dao/types/oid"
)

// SecretCreate is the builder for creating a Secret entity.
type SecretCreate struct {
	config
	mutation *SecretMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreateTime sets the "createTime" field.
func (sc *SecretCreate) SetCreateTime(t time.Time) *SecretCreate {
	sc.mutation.SetCreateTime(t)
	return sc
}

// SetNillableCreateTime sets the "createTime" field if the given value is not nil.
func (sc *SecretCreate) SetNillableCreateTime(t *time.Time) *SecretCreate {
	if t != nil {
		sc.SetCreateTime(*t)
	}
	return sc
}

// SetUpdateTime sets the "updateTime" field.
func (sc *SecretCreate) SetUpdateTime(t time.Time) *SecretCreate {
	sc.mutation.SetUpdateTime(t)
	return sc
}

// SetNillableUpdateTime sets the "updateTime" field if the given value is not nil.
func (sc *SecretCreate) SetNillableUpdateTime(t *time.Time) *SecretCreate {
	if t != nil {
		sc.SetUpdateTime(*t)
	}
	return sc
}

// SetProjectID sets the "projectID" field.
func (sc *SecretCreate) SetProjectID(o oid.ID) *SecretCreate {
	sc.mutation.SetProjectID(o)
	return sc
}

// SetNillableProjectID sets the "projectID" field if the given value is not nil.
func (sc *SecretCreate) SetNillableProjectID(o *oid.ID) *SecretCreate {
	if o != nil {
		sc.SetProjectID(*o)
	}
	return sc
}

// SetName sets the "name" field.
func (sc *SecretCreate) SetName(s string) *SecretCreate {
	sc.mutation.SetName(s)
	return sc
}

// SetValue sets the "value" field.
func (sc *SecretCreate) SetValue(c crypto.String) *SecretCreate {
	sc.mutation.SetValue(c)
	return sc
}

// SetID sets the "id" field.
func (sc *SecretCreate) SetID(o oid.ID) *SecretCreate {
	sc.mutation.SetID(o)
	return sc
}

// SetProject sets the "project" edge to the Project entity.
func (sc *SecretCreate) SetProject(p *Project) *SecretCreate {
	return sc.SetProjectID(p.ID)
}

// Mutation returns the SecretMutation object of the builder.
func (sc *SecretCreate) Mutation() *SecretMutation {
	return sc.mutation
}

// Save creates the Secret in the database.
func (sc *SecretCreate) Save(ctx context.Context) (*Secret, error) {
	if err := sc.defaults(); err != nil {
		return nil, err
	}
	return withHooks[*Secret, SecretMutation](ctx, sc.sqlSave, sc.mutation, sc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (sc *SecretCreate) SaveX(ctx context.Context) *Secret {
	v, err := sc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (sc *SecretCreate) Exec(ctx context.Context) error {
	_, err := sc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (sc *SecretCreate) ExecX(ctx context.Context) {
	if err := sc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (sc *SecretCreate) defaults() error {
	if _, ok := sc.mutation.CreateTime(); !ok {
		if secret.DefaultCreateTime == nil {
			return fmt.Errorf("model: uninitialized secret.DefaultCreateTime (forgotten import model/runtime?)")
		}
		v := secret.DefaultCreateTime()
		sc.mutation.SetCreateTime(v)
	}
	if _, ok := sc.mutation.UpdateTime(); !ok {
		if secret.DefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized secret.DefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := secret.DefaultUpdateTime()
		sc.mutation.SetUpdateTime(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (sc *SecretCreate) check() error {
	if _, ok := sc.mutation.CreateTime(); !ok {
		return &ValidationError{Name: "createTime", err: errors.New(`model: missing required field "Secret.createTime"`)}
	}
	if _, ok := sc.mutation.UpdateTime(); !ok {
		return &ValidationError{Name: "updateTime", err: errors.New(`model: missing required field "Secret.updateTime"`)}
	}
	if _, ok := sc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`model: missing required field "Secret.name"`)}
	}
	if v, ok := sc.mutation.Name(); ok {
		if err := secret.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`model: validator failed for field "Secret.name": %w`, err)}
		}
	}
	if _, ok := sc.mutation.Value(); !ok {
		return &ValidationError{Name: "value", err: errors.New(`model: missing required field "Secret.value"`)}
	}
	if v, ok := sc.mutation.Value(); ok {
		if err := secret.ValueValidator(string(v)); err != nil {
			return &ValidationError{Name: "value", err: fmt.Errorf(`model: validator failed for field "Secret.value": %w`, err)}
		}
	}
	return nil
}

func (sc *SecretCreate) sqlSave(ctx context.Context) (*Secret, error) {
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

func (sc *SecretCreate) createSpec() (*Secret, *sqlgraph.CreateSpec) {
	var (
		_node = &Secret{config: sc.config}
		_spec = sqlgraph.NewCreateSpec(secret.Table, sqlgraph.NewFieldSpec(secret.FieldID, field.TypeString))
	)
	_spec.Schema = sc.schemaConfig.Secret
	_spec.OnConflict = sc.conflict
	if id, ok := sc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := sc.mutation.CreateTime(); ok {
		_spec.SetField(secret.FieldCreateTime, field.TypeTime, value)
		_node.CreateTime = &value
	}
	if value, ok := sc.mutation.UpdateTime(); ok {
		_spec.SetField(secret.FieldUpdateTime, field.TypeTime, value)
		_node.UpdateTime = &value
	}
	if value, ok := sc.mutation.Name(); ok {
		_spec.SetField(secret.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := sc.mutation.Value(); ok {
		_spec.SetField(secret.FieldValue, field.TypeString, value)
		_node.Value = value
	}
	if nodes := sc.mutation.ProjectIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   secret.ProjectTable,
			Columns: []string{secret.ProjectColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: project.FieldID,
				},
			},
		}
		edge.Schema = sc.schemaConfig.Secret
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.ProjectID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Secret.Create().
//		SetCreateTime(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.SecretUpsert) {
//			SetCreateTime(v+v).
//		}).
//		Exec(ctx)
func (sc *SecretCreate) OnConflict(opts ...sql.ConflictOption) *SecretUpsertOne {
	sc.conflict = opts
	return &SecretUpsertOne{
		create: sc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Secret.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (sc *SecretCreate) OnConflictColumns(columns ...string) *SecretUpsertOne {
	sc.conflict = append(sc.conflict, sql.ConflictColumns(columns...))
	return &SecretUpsertOne{
		create: sc,
	}
}

type (
	// SecretUpsertOne is the builder for "upsert"-ing
	//  one Secret node.
	SecretUpsertOne struct {
		create *SecretCreate
	}

	// SecretUpsert is the "OnConflict" setter.
	SecretUpsert struct {
		*sql.UpdateSet
	}
)

// SetUpdateTime sets the "updateTime" field.
func (u *SecretUpsert) SetUpdateTime(v time.Time) *SecretUpsert {
	u.Set(secret.FieldUpdateTime, v)
	return u
}

// UpdateUpdateTime sets the "updateTime" field to the value that was provided on create.
func (u *SecretUpsert) UpdateUpdateTime() *SecretUpsert {
	u.SetExcluded(secret.FieldUpdateTime)
	return u
}

// SetValue sets the "value" field.
func (u *SecretUpsert) SetValue(v crypto.String) *SecretUpsert {
	u.Set(secret.FieldValue, v)
	return u
}

// UpdateValue sets the "value" field to the value that was provided on create.
func (u *SecretUpsert) UpdateValue() *SecretUpsert {
	u.SetExcluded(secret.FieldValue)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.Secret.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(secret.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *SecretUpsertOne) UpdateNewValues() *SecretUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(secret.FieldID)
		}
		if _, exists := u.create.mutation.CreateTime(); exists {
			s.SetIgnore(secret.FieldCreateTime)
		}
		if _, exists := u.create.mutation.ProjectID(); exists {
			s.SetIgnore(secret.FieldProjectID)
		}
		if _, exists := u.create.mutation.Name(); exists {
			s.SetIgnore(secret.FieldName)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Secret.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *SecretUpsertOne) Ignore() *SecretUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *SecretUpsertOne) DoNothing() *SecretUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the SecretCreate.OnConflict
// documentation for more info.
func (u *SecretUpsertOne) Update(set func(*SecretUpsert)) *SecretUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&SecretUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdateTime sets the "updateTime" field.
func (u *SecretUpsertOne) SetUpdateTime(v time.Time) *SecretUpsertOne {
	return u.Update(func(s *SecretUpsert) {
		s.SetUpdateTime(v)
	})
}

// UpdateUpdateTime sets the "updateTime" field to the value that was provided on create.
func (u *SecretUpsertOne) UpdateUpdateTime() *SecretUpsertOne {
	return u.Update(func(s *SecretUpsert) {
		s.UpdateUpdateTime()
	})
}

// SetValue sets the "value" field.
func (u *SecretUpsertOne) SetValue(v crypto.String) *SecretUpsertOne {
	return u.Update(func(s *SecretUpsert) {
		s.SetValue(v)
	})
}

// UpdateValue sets the "value" field to the value that was provided on create.
func (u *SecretUpsertOne) UpdateValue() *SecretUpsertOne {
	return u.Update(func(s *SecretUpsert) {
		s.UpdateValue()
	})
}

// Exec executes the query.
func (u *SecretUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("model: missing options for SecretCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *SecretUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *SecretUpsertOne) ID(ctx context.Context) (id oid.ID, err error) {
	if u.create.driver.Dialect() == dialect.MySQL {
		// In case of "ON CONFLICT", there is no way to get back non-numeric ID
		// fields from the database since MySQL does not support the RETURNING clause.
		return id, errors.New("model: SecretUpsertOne.ID is not supported by MySQL driver. Use SecretUpsertOne.Exec instead")
	}
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *SecretUpsertOne) IDX(ctx context.Context) oid.ID {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// SecretCreateBulk is the builder for creating many Secret entities in bulk.
type SecretCreateBulk struct {
	config
	builders []*SecretCreate
	conflict []sql.ConflictOption
}

// Save creates the Secret entities in the database.
func (scb *SecretCreateBulk) Save(ctx context.Context) ([]*Secret, error) {
	specs := make([]*sqlgraph.CreateSpec, len(scb.builders))
	nodes := make([]*Secret, len(scb.builders))
	mutators := make([]Mutator, len(scb.builders))
	for i := range scb.builders {
		func(i int, root context.Context) {
			builder := scb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*SecretMutation)
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
func (scb *SecretCreateBulk) SaveX(ctx context.Context) []*Secret {
	v, err := scb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (scb *SecretCreateBulk) Exec(ctx context.Context) error {
	_, err := scb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (scb *SecretCreateBulk) ExecX(ctx context.Context) {
	if err := scb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Secret.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.SecretUpsert) {
//			SetCreateTime(v+v).
//		}).
//		Exec(ctx)
func (scb *SecretCreateBulk) OnConflict(opts ...sql.ConflictOption) *SecretUpsertBulk {
	scb.conflict = opts
	return &SecretUpsertBulk{
		create: scb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Secret.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (scb *SecretCreateBulk) OnConflictColumns(columns ...string) *SecretUpsertBulk {
	scb.conflict = append(scb.conflict, sql.ConflictColumns(columns...))
	return &SecretUpsertBulk{
		create: scb,
	}
}

// SecretUpsertBulk is the builder for "upsert"-ing
// a bulk of Secret nodes.
type SecretUpsertBulk struct {
	create *SecretCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Secret.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(secret.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *SecretUpsertBulk) UpdateNewValues() *SecretUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(secret.FieldID)
			}
			if _, exists := b.mutation.CreateTime(); exists {
				s.SetIgnore(secret.FieldCreateTime)
			}
			if _, exists := b.mutation.ProjectID(); exists {
				s.SetIgnore(secret.FieldProjectID)
			}
			if _, exists := b.mutation.Name(); exists {
				s.SetIgnore(secret.FieldName)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Secret.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *SecretUpsertBulk) Ignore() *SecretUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *SecretUpsertBulk) DoNothing() *SecretUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the SecretCreateBulk.OnConflict
// documentation for more info.
func (u *SecretUpsertBulk) Update(set func(*SecretUpsert)) *SecretUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&SecretUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdateTime sets the "updateTime" field.
func (u *SecretUpsertBulk) SetUpdateTime(v time.Time) *SecretUpsertBulk {
	return u.Update(func(s *SecretUpsert) {
		s.SetUpdateTime(v)
	})
}

// UpdateUpdateTime sets the "updateTime" field to the value that was provided on create.
func (u *SecretUpsertBulk) UpdateUpdateTime() *SecretUpsertBulk {
	return u.Update(func(s *SecretUpsert) {
		s.UpdateUpdateTime()
	})
}

// SetValue sets the "value" field.
func (u *SecretUpsertBulk) SetValue(v crypto.String) *SecretUpsertBulk {
	return u.Update(func(s *SecretUpsert) {
		s.SetValue(v)
	})
}

// UpdateValue sets the "value" field to the value that was provided on create.
func (u *SecretUpsertBulk) UpdateValue() *SecretUpsertBulk {
	return u.Update(func(s *SecretUpsert) {
		s.UpdateValue()
	})
}

// Exec executes the query.
func (u *SecretUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("model: OnConflict was set for builder %d. Set it on the SecretCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("model: missing options for SecretCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *SecretUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
