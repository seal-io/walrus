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

	"github.com/seal-io/seal/pkg/dao/model/token"
	"github.com/seal-io/seal/pkg/dao/types/oid"
)

// TokenCreate is the builder for creating a Token entity.
type TokenCreate struct {
	config
	mutation *TokenMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreateTime sets the "createTime" field.
func (tc *TokenCreate) SetCreateTime(t time.Time) *TokenCreate {
	tc.mutation.SetCreateTime(t)
	return tc
}

// SetNillableCreateTime sets the "createTime" field if the given value is not nil.
func (tc *TokenCreate) SetNillableCreateTime(t *time.Time) *TokenCreate {
	if t != nil {
		tc.SetCreateTime(*t)
	}
	return tc
}

// SetUpdateTime sets the "updateTime" field.
func (tc *TokenCreate) SetUpdateTime(t time.Time) *TokenCreate {
	tc.mutation.SetUpdateTime(t)
	return tc
}

// SetNillableUpdateTime sets the "updateTime" field if the given value is not nil.
func (tc *TokenCreate) SetNillableUpdateTime(t *time.Time) *TokenCreate {
	if t != nil {
		tc.SetUpdateTime(*t)
	}
	return tc
}

// SetCasdoorTokenName sets the "casdoorTokenName" field.
func (tc *TokenCreate) SetCasdoorTokenName(s string) *TokenCreate {
	tc.mutation.SetCasdoorTokenName(s)
	return tc
}

// SetCasdoorTokenOwner sets the "casdoorTokenOwner" field.
func (tc *TokenCreate) SetCasdoorTokenOwner(s string) *TokenCreate {
	tc.mutation.SetCasdoorTokenOwner(s)
	return tc
}

// SetName sets the "name" field.
func (tc *TokenCreate) SetName(s string) *TokenCreate {
	tc.mutation.SetName(s)
	return tc
}

// SetExpiration sets the "expiration" field.
func (tc *TokenCreate) SetExpiration(i int) *TokenCreate {
	tc.mutation.SetExpiration(i)
	return tc
}

// SetNillableExpiration sets the "expiration" field if the given value is not nil.
func (tc *TokenCreate) SetNillableExpiration(i *int) *TokenCreate {
	if i != nil {
		tc.SetExpiration(*i)
	}
	return tc
}

// SetID sets the "id" field.
func (tc *TokenCreate) SetID(o oid.ID) *TokenCreate {
	tc.mutation.SetID(o)
	return tc
}

// Mutation returns the TokenMutation object of the builder.
func (tc *TokenCreate) Mutation() *TokenMutation {
	return tc.mutation
}

// Save creates the Token in the database.
func (tc *TokenCreate) Save(ctx context.Context) (*Token, error) {
	if err := tc.defaults(); err != nil {
		return nil, err
	}
	return withHooks[*Token, TokenMutation](ctx, tc.sqlSave, tc.mutation, tc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (tc *TokenCreate) SaveX(ctx context.Context) *Token {
	v, err := tc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tc *TokenCreate) Exec(ctx context.Context) error {
	_, err := tc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tc *TokenCreate) ExecX(ctx context.Context) {
	if err := tc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (tc *TokenCreate) defaults() error {
	if _, ok := tc.mutation.CreateTime(); !ok {
		if token.DefaultCreateTime == nil {
			return fmt.Errorf("model: uninitialized token.DefaultCreateTime (forgotten import model/runtime?)")
		}
		v := token.DefaultCreateTime()
		tc.mutation.SetCreateTime(v)
	}
	if _, ok := tc.mutation.UpdateTime(); !ok {
		if token.DefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized token.DefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := token.DefaultUpdateTime()
		tc.mutation.SetUpdateTime(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (tc *TokenCreate) check() error {
	if _, ok := tc.mutation.CreateTime(); !ok {
		return &ValidationError{Name: "createTime", err: errors.New(`model: missing required field "Token.createTime"`)}
	}
	if _, ok := tc.mutation.UpdateTime(); !ok {
		return &ValidationError{Name: "updateTime", err: errors.New(`model: missing required field "Token.updateTime"`)}
	}
	if _, ok := tc.mutation.CasdoorTokenName(); !ok {
		return &ValidationError{Name: "casdoorTokenName", err: errors.New(`model: missing required field "Token.casdoorTokenName"`)}
	}
	if v, ok := tc.mutation.CasdoorTokenName(); ok {
		if err := token.CasdoorTokenNameValidator(v); err != nil {
			return &ValidationError{Name: "casdoorTokenName", err: fmt.Errorf(`model: validator failed for field "Token.casdoorTokenName": %w`, err)}
		}
	}
	if _, ok := tc.mutation.CasdoorTokenOwner(); !ok {
		return &ValidationError{Name: "casdoorTokenOwner", err: errors.New(`model: missing required field "Token.casdoorTokenOwner"`)}
	}
	if v, ok := tc.mutation.CasdoorTokenOwner(); ok {
		if err := token.CasdoorTokenOwnerValidator(v); err != nil {
			return &ValidationError{Name: "casdoorTokenOwner", err: fmt.Errorf(`model: validator failed for field "Token.casdoorTokenOwner": %w`, err)}
		}
	}
	if _, ok := tc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`model: missing required field "Token.name"`)}
	}
	if v, ok := tc.mutation.Name(); ok {
		if err := token.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`model: validator failed for field "Token.name": %w`, err)}
		}
	}
	return nil
}

func (tc *TokenCreate) sqlSave(ctx context.Context) (*Token, error) {
	if err := tc.check(); err != nil {
		return nil, err
	}
	_node, _spec := tc.createSpec()
	if err := sqlgraph.CreateNode(ctx, tc.driver, _spec); err != nil {
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
	tc.mutation.id = &_node.ID
	tc.mutation.done = true
	return _node, nil
}

func (tc *TokenCreate) createSpec() (*Token, *sqlgraph.CreateSpec) {
	var (
		_node = &Token{config: tc.config}
		_spec = sqlgraph.NewCreateSpec(token.Table, sqlgraph.NewFieldSpec(token.FieldID, field.TypeString))
	)
	_spec.Schema = tc.schemaConfig.Token
	_spec.OnConflict = tc.conflict
	if id, ok := tc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := tc.mutation.CreateTime(); ok {
		_spec.SetField(token.FieldCreateTime, field.TypeTime, value)
		_node.CreateTime = &value
	}
	if value, ok := tc.mutation.UpdateTime(); ok {
		_spec.SetField(token.FieldUpdateTime, field.TypeTime, value)
		_node.UpdateTime = &value
	}
	if value, ok := tc.mutation.CasdoorTokenName(); ok {
		_spec.SetField(token.FieldCasdoorTokenName, field.TypeString, value)
		_node.CasdoorTokenName = value
	}
	if value, ok := tc.mutation.CasdoorTokenOwner(); ok {
		_spec.SetField(token.FieldCasdoorTokenOwner, field.TypeString, value)
		_node.CasdoorTokenOwner = value
	}
	if value, ok := tc.mutation.Name(); ok {
		_spec.SetField(token.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := tc.mutation.Expiration(); ok {
		_spec.SetField(token.FieldExpiration, field.TypeInt, value)
		_node.Expiration = &value
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Token.Create().
//		SetCreateTime(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.TokenUpsert) {
//			SetCreateTime(v+v).
//		}).
//		Exec(ctx)
func (tc *TokenCreate) OnConflict(opts ...sql.ConflictOption) *TokenUpsertOne {
	tc.conflict = opts
	return &TokenUpsertOne{
		create: tc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Token.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (tc *TokenCreate) OnConflictColumns(columns ...string) *TokenUpsertOne {
	tc.conflict = append(tc.conflict, sql.ConflictColumns(columns...))
	return &TokenUpsertOne{
		create: tc,
	}
}

type (
	// TokenUpsertOne is the builder for "upsert"-ing
	//  one Token node.
	TokenUpsertOne struct {
		create *TokenCreate
	}

	// TokenUpsert is the "OnConflict" setter.
	TokenUpsert struct {
		*sql.UpdateSet
	}
)

// SetUpdateTime sets the "updateTime" field.
func (u *TokenUpsert) SetUpdateTime(v time.Time) *TokenUpsert {
	u.Set(token.FieldUpdateTime, v)
	return u
}

// UpdateUpdateTime sets the "updateTime" field to the value that was provided on create.
func (u *TokenUpsert) UpdateUpdateTime() *TokenUpsert {
	u.SetExcluded(token.FieldUpdateTime)
	return u
}

// SetCasdoorTokenName sets the "casdoorTokenName" field.
func (u *TokenUpsert) SetCasdoorTokenName(v string) *TokenUpsert {
	u.Set(token.FieldCasdoorTokenName, v)
	return u
}

// UpdateCasdoorTokenName sets the "casdoorTokenName" field to the value that was provided on create.
func (u *TokenUpsert) UpdateCasdoorTokenName() *TokenUpsert {
	u.SetExcluded(token.FieldCasdoorTokenName)
	return u
}

// SetCasdoorTokenOwner sets the "casdoorTokenOwner" field.
func (u *TokenUpsert) SetCasdoorTokenOwner(v string) *TokenUpsert {
	u.Set(token.FieldCasdoorTokenOwner, v)
	return u
}

// UpdateCasdoorTokenOwner sets the "casdoorTokenOwner" field to the value that was provided on create.
func (u *TokenUpsert) UpdateCasdoorTokenOwner() *TokenUpsert {
	u.SetExcluded(token.FieldCasdoorTokenOwner)
	return u
}

// SetName sets the "name" field.
func (u *TokenUpsert) SetName(v string) *TokenUpsert {
	u.Set(token.FieldName, v)
	return u
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *TokenUpsert) UpdateName() *TokenUpsert {
	u.SetExcluded(token.FieldName)
	return u
}

// SetExpiration sets the "expiration" field.
func (u *TokenUpsert) SetExpiration(v int) *TokenUpsert {
	u.Set(token.FieldExpiration, v)
	return u
}

// UpdateExpiration sets the "expiration" field to the value that was provided on create.
func (u *TokenUpsert) UpdateExpiration() *TokenUpsert {
	u.SetExcluded(token.FieldExpiration)
	return u
}

// AddExpiration adds v to the "expiration" field.
func (u *TokenUpsert) AddExpiration(v int) *TokenUpsert {
	u.Add(token.FieldExpiration, v)
	return u
}

// ClearExpiration clears the value of the "expiration" field.
func (u *TokenUpsert) ClearExpiration() *TokenUpsert {
	u.SetNull(token.FieldExpiration)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.Token.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(token.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *TokenUpsertOne) UpdateNewValues() *TokenUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(token.FieldID)
		}
		if _, exists := u.create.mutation.CreateTime(); exists {
			s.SetIgnore(token.FieldCreateTime)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Token.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *TokenUpsertOne) Ignore() *TokenUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *TokenUpsertOne) DoNothing() *TokenUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the TokenCreate.OnConflict
// documentation for more info.
func (u *TokenUpsertOne) Update(set func(*TokenUpsert)) *TokenUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&TokenUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdateTime sets the "updateTime" field.
func (u *TokenUpsertOne) SetUpdateTime(v time.Time) *TokenUpsertOne {
	return u.Update(func(s *TokenUpsert) {
		s.SetUpdateTime(v)
	})
}

// UpdateUpdateTime sets the "updateTime" field to the value that was provided on create.
func (u *TokenUpsertOne) UpdateUpdateTime() *TokenUpsertOne {
	return u.Update(func(s *TokenUpsert) {
		s.UpdateUpdateTime()
	})
}

// SetCasdoorTokenName sets the "casdoorTokenName" field.
func (u *TokenUpsertOne) SetCasdoorTokenName(v string) *TokenUpsertOne {
	return u.Update(func(s *TokenUpsert) {
		s.SetCasdoorTokenName(v)
	})
}

// UpdateCasdoorTokenName sets the "casdoorTokenName" field to the value that was provided on create.
func (u *TokenUpsertOne) UpdateCasdoorTokenName() *TokenUpsertOne {
	return u.Update(func(s *TokenUpsert) {
		s.UpdateCasdoorTokenName()
	})
}

// SetCasdoorTokenOwner sets the "casdoorTokenOwner" field.
func (u *TokenUpsertOne) SetCasdoorTokenOwner(v string) *TokenUpsertOne {
	return u.Update(func(s *TokenUpsert) {
		s.SetCasdoorTokenOwner(v)
	})
}

// UpdateCasdoorTokenOwner sets the "casdoorTokenOwner" field to the value that was provided on create.
func (u *TokenUpsertOne) UpdateCasdoorTokenOwner() *TokenUpsertOne {
	return u.Update(func(s *TokenUpsert) {
		s.UpdateCasdoorTokenOwner()
	})
}

// SetName sets the "name" field.
func (u *TokenUpsertOne) SetName(v string) *TokenUpsertOne {
	return u.Update(func(s *TokenUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *TokenUpsertOne) UpdateName() *TokenUpsertOne {
	return u.Update(func(s *TokenUpsert) {
		s.UpdateName()
	})
}

// SetExpiration sets the "expiration" field.
func (u *TokenUpsertOne) SetExpiration(v int) *TokenUpsertOne {
	return u.Update(func(s *TokenUpsert) {
		s.SetExpiration(v)
	})
}

// AddExpiration adds v to the "expiration" field.
func (u *TokenUpsertOne) AddExpiration(v int) *TokenUpsertOne {
	return u.Update(func(s *TokenUpsert) {
		s.AddExpiration(v)
	})
}

// UpdateExpiration sets the "expiration" field to the value that was provided on create.
func (u *TokenUpsertOne) UpdateExpiration() *TokenUpsertOne {
	return u.Update(func(s *TokenUpsert) {
		s.UpdateExpiration()
	})
}

// ClearExpiration clears the value of the "expiration" field.
func (u *TokenUpsertOne) ClearExpiration() *TokenUpsertOne {
	return u.Update(func(s *TokenUpsert) {
		s.ClearExpiration()
	})
}

// Exec executes the query.
func (u *TokenUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("model: missing options for TokenCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *TokenUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *TokenUpsertOne) ID(ctx context.Context) (id oid.ID, err error) {
	if u.create.driver.Dialect() == dialect.MySQL {
		// In case of "ON CONFLICT", there is no way to get back non-numeric ID
		// fields from the database since MySQL does not support the RETURNING clause.
		return id, errors.New("model: TokenUpsertOne.ID is not supported by MySQL driver. Use TokenUpsertOne.Exec instead")
	}
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *TokenUpsertOne) IDX(ctx context.Context) oid.ID {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// TokenCreateBulk is the builder for creating many Token entities in bulk.
type TokenCreateBulk struct {
	config
	builders []*TokenCreate
	conflict []sql.ConflictOption
}

// Save creates the Token entities in the database.
func (tcb *TokenCreateBulk) Save(ctx context.Context) ([]*Token, error) {
	specs := make([]*sqlgraph.CreateSpec, len(tcb.builders))
	nodes := make([]*Token, len(tcb.builders))
	mutators := make([]Mutator, len(tcb.builders))
	for i := range tcb.builders {
		func(i int, root context.Context) {
			builder := tcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*TokenMutation)
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
					_, err = mutators[i+1].Mutate(root, tcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = tcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, tcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, tcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (tcb *TokenCreateBulk) SaveX(ctx context.Context) []*Token {
	v, err := tcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tcb *TokenCreateBulk) Exec(ctx context.Context) error {
	_, err := tcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tcb *TokenCreateBulk) ExecX(ctx context.Context) {
	if err := tcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Token.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.TokenUpsert) {
//			SetCreateTime(v+v).
//		}).
//		Exec(ctx)
func (tcb *TokenCreateBulk) OnConflict(opts ...sql.ConflictOption) *TokenUpsertBulk {
	tcb.conflict = opts
	return &TokenUpsertBulk{
		create: tcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Token.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (tcb *TokenCreateBulk) OnConflictColumns(columns ...string) *TokenUpsertBulk {
	tcb.conflict = append(tcb.conflict, sql.ConflictColumns(columns...))
	return &TokenUpsertBulk{
		create: tcb,
	}
}

// TokenUpsertBulk is the builder for "upsert"-ing
// a bulk of Token nodes.
type TokenUpsertBulk struct {
	create *TokenCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Token.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(token.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *TokenUpsertBulk) UpdateNewValues() *TokenUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(token.FieldID)
			}
			if _, exists := b.mutation.CreateTime(); exists {
				s.SetIgnore(token.FieldCreateTime)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Token.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *TokenUpsertBulk) Ignore() *TokenUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *TokenUpsertBulk) DoNothing() *TokenUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the TokenCreateBulk.OnConflict
// documentation for more info.
func (u *TokenUpsertBulk) Update(set func(*TokenUpsert)) *TokenUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&TokenUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdateTime sets the "updateTime" field.
func (u *TokenUpsertBulk) SetUpdateTime(v time.Time) *TokenUpsertBulk {
	return u.Update(func(s *TokenUpsert) {
		s.SetUpdateTime(v)
	})
}

// UpdateUpdateTime sets the "updateTime" field to the value that was provided on create.
func (u *TokenUpsertBulk) UpdateUpdateTime() *TokenUpsertBulk {
	return u.Update(func(s *TokenUpsert) {
		s.UpdateUpdateTime()
	})
}

// SetCasdoorTokenName sets the "casdoorTokenName" field.
func (u *TokenUpsertBulk) SetCasdoorTokenName(v string) *TokenUpsertBulk {
	return u.Update(func(s *TokenUpsert) {
		s.SetCasdoorTokenName(v)
	})
}

// UpdateCasdoorTokenName sets the "casdoorTokenName" field to the value that was provided on create.
func (u *TokenUpsertBulk) UpdateCasdoorTokenName() *TokenUpsertBulk {
	return u.Update(func(s *TokenUpsert) {
		s.UpdateCasdoorTokenName()
	})
}

// SetCasdoorTokenOwner sets the "casdoorTokenOwner" field.
func (u *TokenUpsertBulk) SetCasdoorTokenOwner(v string) *TokenUpsertBulk {
	return u.Update(func(s *TokenUpsert) {
		s.SetCasdoorTokenOwner(v)
	})
}

// UpdateCasdoorTokenOwner sets the "casdoorTokenOwner" field to the value that was provided on create.
func (u *TokenUpsertBulk) UpdateCasdoorTokenOwner() *TokenUpsertBulk {
	return u.Update(func(s *TokenUpsert) {
		s.UpdateCasdoorTokenOwner()
	})
}

// SetName sets the "name" field.
func (u *TokenUpsertBulk) SetName(v string) *TokenUpsertBulk {
	return u.Update(func(s *TokenUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *TokenUpsertBulk) UpdateName() *TokenUpsertBulk {
	return u.Update(func(s *TokenUpsert) {
		s.UpdateName()
	})
}

// SetExpiration sets the "expiration" field.
func (u *TokenUpsertBulk) SetExpiration(v int) *TokenUpsertBulk {
	return u.Update(func(s *TokenUpsert) {
		s.SetExpiration(v)
	})
}

// AddExpiration adds v to the "expiration" field.
func (u *TokenUpsertBulk) AddExpiration(v int) *TokenUpsertBulk {
	return u.Update(func(s *TokenUpsert) {
		s.AddExpiration(v)
	})
}

// UpdateExpiration sets the "expiration" field to the value that was provided on create.
func (u *TokenUpsertBulk) UpdateExpiration() *TokenUpsertBulk {
	return u.Update(func(s *TokenUpsert) {
		s.UpdateExpiration()
	})
}

// ClearExpiration clears the value of the "expiration" field.
func (u *TokenUpsertBulk) ClearExpiration() *TokenUpsertBulk {
	return u.Update(func(s *TokenUpsert) {
		s.ClearExpiration()
	})
}

// Exec executes the query.
func (u *TokenUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("model: OnConflict was set for builder %d. Set it on the TokenCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("model: missing options for TokenCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *TokenUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
