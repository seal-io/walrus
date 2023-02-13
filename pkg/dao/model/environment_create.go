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

	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/oid"
)

// EnvironmentCreate is the builder for creating a Environment entity.
type EnvironmentCreate struct {
	config
	mutation *EnvironmentMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreateTime sets the "createTime" field.
func (ec *EnvironmentCreate) SetCreateTime(t time.Time) *EnvironmentCreate {
	ec.mutation.SetCreateTime(t)
	return ec
}

// SetNillableCreateTime sets the "createTime" field if the given value is not nil.
func (ec *EnvironmentCreate) SetNillableCreateTime(t *time.Time) *EnvironmentCreate {
	if t != nil {
		ec.SetCreateTime(*t)
	}
	return ec
}

// SetUpdateTime sets the "updateTime" field.
func (ec *EnvironmentCreate) SetUpdateTime(t time.Time) *EnvironmentCreate {
	ec.mutation.SetUpdateTime(t)
	return ec
}

// SetNillableUpdateTime sets the "updateTime" field if the given value is not nil.
func (ec *EnvironmentCreate) SetNillableUpdateTime(t *time.Time) *EnvironmentCreate {
	if t != nil {
		ec.SetUpdateTime(*t)
	}
	return ec
}

// SetConnectorIDs sets the "connectorIDs" field.
func (ec *EnvironmentCreate) SetConnectorIDs(o []oid.ID) *EnvironmentCreate {
	ec.mutation.SetConnectorIDs(o)
	return ec
}

// SetVariables sets the "variables" field.
func (ec *EnvironmentCreate) SetVariables(m map[string]interface{}) *EnvironmentCreate {
	ec.mutation.SetVariables(m)
	return ec
}

// SetID sets the "id" field.
func (ec *EnvironmentCreate) SetID(o oid.ID) *EnvironmentCreate {
	ec.mutation.SetID(o)
	return ec
}

// Mutation returns the EnvironmentMutation object of the builder.
func (ec *EnvironmentCreate) Mutation() *EnvironmentMutation {
	return ec.mutation
}

// Save creates the Environment in the database.
func (ec *EnvironmentCreate) Save(ctx context.Context) (*Environment, error) {
	if err := ec.defaults(); err != nil {
		return nil, err
	}
	return withHooks[*Environment, EnvironmentMutation](ctx, ec.sqlSave, ec.mutation, ec.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (ec *EnvironmentCreate) SaveX(ctx context.Context) *Environment {
	v, err := ec.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ec *EnvironmentCreate) Exec(ctx context.Context) error {
	_, err := ec.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ec *EnvironmentCreate) ExecX(ctx context.Context) {
	if err := ec.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ec *EnvironmentCreate) defaults() error {
	if _, ok := ec.mutation.CreateTime(); !ok {
		if environment.DefaultCreateTime == nil {
			return fmt.Errorf("model: uninitialized environment.DefaultCreateTime (forgotten import model/runtime?)")
		}
		v := environment.DefaultCreateTime()
		ec.mutation.SetCreateTime(v)
	}
	if _, ok := ec.mutation.UpdateTime(); !ok {
		if environment.DefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized environment.DefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := environment.DefaultUpdateTime()
		ec.mutation.SetUpdateTime(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (ec *EnvironmentCreate) check() error {
	if _, ok := ec.mutation.CreateTime(); !ok {
		return &ValidationError{Name: "createTime", err: errors.New(`model: missing required field "Environment.createTime"`)}
	}
	if _, ok := ec.mutation.UpdateTime(); !ok {
		return &ValidationError{Name: "updateTime", err: errors.New(`model: missing required field "Environment.updateTime"`)}
	}
	return nil
}

func (ec *EnvironmentCreate) sqlSave(ctx context.Context) (*Environment, error) {
	if err := ec.check(); err != nil {
		return nil, err
	}
	_node, _spec := ec.createSpec()
	if err := sqlgraph.CreateNode(ctx, ec.driver, _spec); err != nil {
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
	ec.mutation.id = &_node.ID
	ec.mutation.done = true
	return _node, nil
}

func (ec *EnvironmentCreate) createSpec() (*Environment, *sqlgraph.CreateSpec) {
	var (
		_node = &Environment{config: ec.config}
		_spec = &sqlgraph.CreateSpec{
			Table: environment.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeOther,
				Column: environment.FieldID,
			},
		}
	)
	_spec.Schema = ec.schemaConfig.Environment
	_spec.OnConflict = ec.conflict
	if id, ok := ec.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := ec.mutation.CreateTime(); ok {
		_spec.SetField(environment.FieldCreateTime, field.TypeTime, value)
		_node.CreateTime = &value
	}
	if value, ok := ec.mutation.UpdateTime(); ok {
		_spec.SetField(environment.FieldUpdateTime, field.TypeTime, value)
		_node.UpdateTime = &value
	}
	if value, ok := ec.mutation.ConnectorIDs(); ok {
		_spec.SetField(environment.FieldConnectorIDs, field.TypeJSON, value)
		_node.ConnectorIDs = value
	}
	if value, ok := ec.mutation.Variables(); ok {
		_spec.SetField(environment.FieldVariables, field.TypeJSON, value)
		_node.Variables = value
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Environment.Create().
//		SetCreateTime(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.EnvironmentUpsert) {
//			SetCreateTime(v+v).
//		}).
//		Exec(ctx)
func (ec *EnvironmentCreate) OnConflict(opts ...sql.ConflictOption) *EnvironmentUpsertOne {
	ec.conflict = opts
	return &EnvironmentUpsertOne{
		create: ec,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Environment.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (ec *EnvironmentCreate) OnConflictColumns(columns ...string) *EnvironmentUpsertOne {
	ec.conflict = append(ec.conflict, sql.ConflictColumns(columns...))
	return &EnvironmentUpsertOne{
		create: ec,
	}
}

type (
	// EnvironmentUpsertOne is the builder for "upsert"-ing
	//  one Environment node.
	EnvironmentUpsertOne struct {
		create *EnvironmentCreate
	}

	// EnvironmentUpsert is the "OnConflict" setter.
	EnvironmentUpsert struct {
		*sql.UpdateSet
	}
)

// SetUpdateTime sets the "updateTime" field.
func (u *EnvironmentUpsert) SetUpdateTime(v time.Time) *EnvironmentUpsert {
	u.Set(environment.FieldUpdateTime, v)
	return u
}

// UpdateUpdateTime sets the "updateTime" field to the value that was provided on create.
func (u *EnvironmentUpsert) UpdateUpdateTime() *EnvironmentUpsert {
	u.SetExcluded(environment.FieldUpdateTime)
	return u
}

// SetConnectorIDs sets the "connectorIDs" field.
func (u *EnvironmentUpsert) SetConnectorIDs(v []oid.ID) *EnvironmentUpsert {
	u.Set(environment.FieldConnectorIDs, v)
	return u
}

// UpdateConnectorIDs sets the "connectorIDs" field to the value that was provided on create.
func (u *EnvironmentUpsert) UpdateConnectorIDs() *EnvironmentUpsert {
	u.SetExcluded(environment.FieldConnectorIDs)
	return u
}

// ClearConnectorIDs clears the value of the "connectorIDs" field.
func (u *EnvironmentUpsert) ClearConnectorIDs() *EnvironmentUpsert {
	u.SetNull(environment.FieldConnectorIDs)
	return u
}

// SetVariables sets the "variables" field.
func (u *EnvironmentUpsert) SetVariables(v map[string]interface{}) *EnvironmentUpsert {
	u.Set(environment.FieldVariables, v)
	return u
}

// UpdateVariables sets the "variables" field to the value that was provided on create.
func (u *EnvironmentUpsert) UpdateVariables() *EnvironmentUpsert {
	u.SetExcluded(environment.FieldVariables)
	return u
}

// ClearVariables clears the value of the "variables" field.
func (u *EnvironmentUpsert) ClearVariables() *EnvironmentUpsert {
	u.SetNull(environment.FieldVariables)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.Environment.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(environment.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *EnvironmentUpsertOne) UpdateNewValues() *EnvironmentUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(environment.FieldID)
		}
		if _, exists := u.create.mutation.CreateTime(); exists {
			s.SetIgnore(environment.FieldCreateTime)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Environment.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *EnvironmentUpsertOne) Ignore() *EnvironmentUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *EnvironmentUpsertOne) DoNothing() *EnvironmentUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the EnvironmentCreate.OnConflict
// documentation for more info.
func (u *EnvironmentUpsertOne) Update(set func(*EnvironmentUpsert)) *EnvironmentUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&EnvironmentUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdateTime sets the "updateTime" field.
func (u *EnvironmentUpsertOne) SetUpdateTime(v time.Time) *EnvironmentUpsertOne {
	return u.Update(func(s *EnvironmentUpsert) {
		s.SetUpdateTime(v)
	})
}

// UpdateUpdateTime sets the "updateTime" field to the value that was provided on create.
func (u *EnvironmentUpsertOne) UpdateUpdateTime() *EnvironmentUpsertOne {
	return u.Update(func(s *EnvironmentUpsert) {
		s.UpdateUpdateTime()
	})
}

// SetConnectorIDs sets the "connectorIDs" field.
func (u *EnvironmentUpsertOne) SetConnectorIDs(v []oid.ID) *EnvironmentUpsertOne {
	return u.Update(func(s *EnvironmentUpsert) {
		s.SetConnectorIDs(v)
	})
}

// UpdateConnectorIDs sets the "connectorIDs" field to the value that was provided on create.
func (u *EnvironmentUpsertOne) UpdateConnectorIDs() *EnvironmentUpsertOne {
	return u.Update(func(s *EnvironmentUpsert) {
		s.UpdateConnectorIDs()
	})
}

// ClearConnectorIDs clears the value of the "connectorIDs" field.
func (u *EnvironmentUpsertOne) ClearConnectorIDs() *EnvironmentUpsertOne {
	return u.Update(func(s *EnvironmentUpsert) {
		s.ClearConnectorIDs()
	})
}

// SetVariables sets the "variables" field.
func (u *EnvironmentUpsertOne) SetVariables(v map[string]interface{}) *EnvironmentUpsertOne {
	return u.Update(func(s *EnvironmentUpsert) {
		s.SetVariables(v)
	})
}

// UpdateVariables sets the "variables" field to the value that was provided on create.
func (u *EnvironmentUpsertOne) UpdateVariables() *EnvironmentUpsertOne {
	return u.Update(func(s *EnvironmentUpsert) {
		s.UpdateVariables()
	})
}

// ClearVariables clears the value of the "variables" field.
func (u *EnvironmentUpsertOne) ClearVariables() *EnvironmentUpsertOne {
	return u.Update(func(s *EnvironmentUpsert) {
		s.ClearVariables()
	})
}

// Exec executes the query.
func (u *EnvironmentUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("model: missing options for EnvironmentCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *EnvironmentUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *EnvironmentUpsertOne) ID(ctx context.Context) (id oid.ID, err error) {
	if u.create.driver.Dialect() == dialect.MySQL {
		// In case of "ON CONFLICT", there is no way to get back non-numeric ID
		// fields from the database since MySQL does not support the RETURNING clause.
		return id, errors.New("model: EnvironmentUpsertOne.ID is not supported by MySQL driver. Use EnvironmentUpsertOne.Exec instead")
	}
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *EnvironmentUpsertOne) IDX(ctx context.Context) oid.ID {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// EnvironmentCreateBulk is the builder for creating many Environment entities in bulk.
type EnvironmentCreateBulk struct {
	config
	builders []*EnvironmentCreate
	conflict []sql.ConflictOption
}

// Save creates the Environment entities in the database.
func (ecb *EnvironmentCreateBulk) Save(ctx context.Context) ([]*Environment, error) {
	specs := make([]*sqlgraph.CreateSpec, len(ecb.builders))
	nodes := make([]*Environment, len(ecb.builders))
	mutators := make([]Mutator, len(ecb.builders))
	for i := range ecb.builders {
		func(i int, root context.Context) {
			builder := ecb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*EnvironmentMutation)
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
					_, err = mutators[i+1].Mutate(root, ecb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = ecb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ecb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, ecb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ecb *EnvironmentCreateBulk) SaveX(ctx context.Context) []*Environment {
	v, err := ecb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ecb *EnvironmentCreateBulk) Exec(ctx context.Context) error {
	_, err := ecb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ecb *EnvironmentCreateBulk) ExecX(ctx context.Context) {
	if err := ecb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Environment.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.EnvironmentUpsert) {
//			SetCreateTime(v+v).
//		}).
//		Exec(ctx)
func (ecb *EnvironmentCreateBulk) OnConflict(opts ...sql.ConflictOption) *EnvironmentUpsertBulk {
	ecb.conflict = opts
	return &EnvironmentUpsertBulk{
		create: ecb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Environment.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (ecb *EnvironmentCreateBulk) OnConflictColumns(columns ...string) *EnvironmentUpsertBulk {
	ecb.conflict = append(ecb.conflict, sql.ConflictColumns(columns...))
	return &EnvironmentUpsertBulk{
		create: ecb,
	}
}

// EnvironmentUpsertBulk is the builder for "upsert"-ing
// a bulk of Environment nodes.
type EnvironmentUpsertBulk struct {
	create *EnvironmentCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Environment.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(environment.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *EnvironmentUpsertBulk) UpdateNewValues() *EnvironmentUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(environment.FieldID)
			}
			if _, exists := b.mutation.CreateTime(); exists {
				s.SetIgnore(environment.FieldCreateTime)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Environment.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *EnvironmentUpsertBulk) Ignore() *EnvironmentUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *EnvironmentUpsertBulk) DoNothing() *EnvironmentUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the EnvironmentCreateBulk.OnConflict
// documentation for more info.
func (u *EnvironmentUpsertBulk) Update(set func(*EnvironmentUpsert)) *EnvironmentUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&EnvironmentUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdateTime sets the "updateTime" field.
func (u *EnvironmentUpsertBulk) SetUpdateTime(v time.Time) *EnvironmentUpsertBulk {
	return u.Update(func(s *EnvironmentUpsert) {
		s.SetUpdateTime(v)
	})
}

// UpdateUpdateTime sets the "updateTime" field to the value that was provided on create.
func (u *EnvironmentUpsertBulk) UpdateUpdateTime() *EnvironmentUpsertBulk {
	return u.Update(func(s *EnvironmentUpsert) {
		s.UpdateUpdateTime()
	})
}

// SetConnectorIDs sets the "connectorIDs" field.
func (u *EnvironmentUpsertBulk) SetConnectorIDs(v []oid.ID) *EnvironmentUpsertBulk {
	return u.Update(func(s *EnvironmentUpsert) {
		s.SetConnectorIDs(v)
	})
}

// UpdateConnectorIDs sets the "connectorIDs" field to the value that was provided on create.
func (u *EnvironmentUpsertBulk) UpdateConnectorIDs() *EnvironmentUpsertBulk {
	return u.Update(func(s *EnvironmentUpsert) {
		s.UpdateConnectorIDs()
	})
}

// ClearConnectorIDs clears the value of the "connectorIDs" field.
func (u *EnvironmentUpsertBulk) ClearConnectorIDs() *EnvironmentUpsertBulk {
	return u.Update(func(s *EnvironmentUpsert) {
		s.ClearConnectorIDs()
	})
}

// SetVariables sets the "variables" field.
func (u *EnvironmentUpsertBulk) SetVariables(v map[string]interface{}) *EnvironmentUpsertBulk {
	return u.Update(func(s *EnvironmentUpsert) {
		s.SetVariables(v)
	})
}

// UpdateVariables sets the "variables" field to the value that was provided on create.
func (u *EnvironmentUpsertBulk) UpdateVariables() *EnvironmentUpsertBulk {
	return u.Update(func(s *EnvironmentUpsert) {
		s.UpdateVariables()
	})
}

// ClearVariables clears the value of the "variables" field.
func (u *EnvironmentUpsertBulk) ClearVariables() *EnvironmentUpsertBulk {
	return u.Update(func(s *EnvironmentUpsert) {
		s.ClearVariables()
	})
}

// Exec executes the query.
func (u *EnvironmentUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("model: OnConflict was set for builder %d. Set it on the EnvironmentCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("model: missing options for EnvironmentCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *EnvironmentUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
