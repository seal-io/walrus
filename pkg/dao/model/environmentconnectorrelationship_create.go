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

	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/environmentconnectorrelationship"
	"github.com/seal-io/seal/pkg/dao/types/oid"
)

// EnvironmentConnectorRelationshipCreate is the builder for creating a EnvironmentConnectorRelationship entity.
type EnvironmentConnectorRelationshipCreate struct {
	config
	mutation *EnvironmentConnectorRelationshipMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreateTime sets the "createTime" field.
func (ecrc *EnvironmentConnectorRelationshipCreate) SetCreateTime(t time.Time) *EnvironmentConnectorRelationshipCreate {
	ecrc.mutation.SetCreateTime(t)
	return ecrc
}

// SetNillableCreateTime sets the "createTime" field if the given value is not nil.
func (ecrc *EnvironmentConnectorRelationshipCreate) SetNillableCreateTime(t *time.Time) *EnvironmentConnectorRelationshipCreate {
	if t != nil {
		ecrc.SetCreateTime(*t)
	}
	return ecrc
}

// SetEnvironmentID sets the "environment_id" field.
func (ecrc *EnvironmentConnectorRelationshipCreate) SetEnvironmentID(o oid.ID) *EnvironmentConnectorRelationshipCreate {
	ecrc.mutation.SetEnvironmentID(o)
	return ecrc
}

// SetConnectorID sets the "connector_id" field.
func (ecrc *EnvironmentConnectorRelationshipCreate) SetConnectorID(o oid.ID) *EnvironmentConnectorRelationshipCreate {
	ecrc.mutation.SetConnectorID(o)
	return ecrc
}

// SetEnvironment sets the "environment" edge to the Environment entity.
func (ecrc *EnvironmentConnectorRelationshipCreate) SetEnvironment(e *Environment) *EnvironmentConnectorRelationshipCreate {
	return ecrc.SetEnvironmentID(e.ID)
}

// SetConnector sets the "connector" edge to the Connector entity.
func (ecrc *EnvironmentConnectorRelationshipCreate) SetConnector(c *Connector) *EnvironmentConnectorRelationshipCreate {
	return ecrc.SetConnectorID(c.ID)
}

// Mutation returns the EnvironmentConnectorRelationshipMutation object of the builder.
func (ecrc *EnvironmentConnectorRelationshipCreate) Mutation() *EnvironmentConnectorRelationshipMutation {
	return ecrc.mutation
}

// Save creates the EnvironmentConnectorRelationship in the database.
func (ecrc *EnvironmentConnectorRelationshipCreate) Save(ctx context.Context) (*EnvironmentConnectorRelationship, error) {
	ecrc.defaults()
	return withHooks[*EnvironmentConnectorRelationship, EnvironmentConnectorRelationshipMutation](ctx, ecrc.sqlSave, ecrc.mutation, ecrc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (ecrc *EnvironmentConnectorRelationshipCreate) SaveX(ctx context.Context) *EnvironmentConnectorRelationship {
	v, err := ecrc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ecrc *EnvironmentConnectorRelationshipCreate) Exec(ctx context.Context) error {
	_, err := ecrc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ecrc *EnvironmentConnectorRelationshipCreate) ExecX(ctx context.Context) {
	if err := ecrc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ecrc *EnvironmentConnectorRelationshipCreate) defaults() {
	if _, ok := ecrc.mutation.CreateTime(); !ok {
		v := environmentconnectorrelationship.DefaultCreateTime()
		ecrc.mutation.SetCreateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ecrc *EnvironmentConnectorRelationshipCreate) check() error {
	if _, ok := ecrc.mutation.CreateTime(); !ok {
		return &ValidationError{Name: "createTime", err: errors.New(`model: missing required field "EnvironmentConnectorRelationship.createTime"`)}
	}
	if _, ok := ecrc.mutation.EnvironmentID(); !ok {
		return &ValidationError{Name: "environment_id", err: errors.New(`model: missing required field "EnvironmentConnectorRelationship.environment_id"`)}
	}
	if v, ok := ecrc.mutation.EnvironmentID(); ok {
		if err := environmentconnectorrelationship.EnvironmentIDValidator(string(v)); err != nil {
			return &ValidationError{Name: "environment_id", err: fmt.Errorf(`model: validator failed for field "EnvironmentConnectorRelationship.environment_id": %w`, err)}
		}
	}
	if _, ok := ecrc.mutation.ConnectorID(); !ok {
		return &ValidationError{Name: "connector_id", err: errors.New(`model: missing required field "EnvironmentConnectorRelationship.connector_id"`)}
	}
	if v, ok := ecrc.mutation.ConnectorID(); ok {
		if err := environmentconnectorrelationship.ConnectorIDValidator(string(v)); err != nil {
			return &ValidationError{Name: "connector_id", err: fmt.Errorf(`model: validator failed for field "EnvironmentConnectorRelationship.connector_id": %w`, err)}
		}
	}
	if _, ok := ecrc.mutation.EnvironmentID(); !ok {
		return &ValidationError{Name: "environment", err: errors.New(`model: missing required edge "EnvironmentConnectorRelationship.environment"`)}
	}
	if _, ok := ecrc.mutation.ConnectorID(); !ok {
		return &ValidationError{Name: "connector", err: errors.New(`model: missing required edge "EnvironmentConnectorRelationship.connector"`)}
	}
	return nil
}

func (ecrc *EnvironmentConnectorRelationshipCreate) sqlSave(ctx context.Context) (*EnvironmentConnectorRelationship, error) {
	if err := ecrc.check(); err != nil {
		return nil, err
	}
	_node, _spec := ecrc.createSpec()
	if err := sqlgraph.CreateNode(ctx, ecrc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	return _node, nil
}

func (ecrc *EnvironmentConnectorRelationshipCreate) createSpec() (*EnvironmentConnectorRelationship, *sqlgraph.CreateSpec) {
	var (
		_node = &EnvironmentConnectorRelationship{config: ecrc.config}
		_spec = sqlgraph.NewCreateSpec(environmentconnectorrelationship.Table, nil)
	)
	_spec.Schema = ecrc.schemaConfig.EnvironmentConnectorRelationship
	_spec.OnConflict = ecrc.conflict
	if value, ok := ecrc.mutation.CreateTime(); ok {
		_spec.SetField(environmentconnectorrelationship.FieldCreateTime, field.TypeTime, value)
		_node.CreateTime = &value
	}
	if nodes := ecrc.mutation.EnvironmentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   environmentconnectorrelationship.EnvironmentTable,
			Columns: []string{environmentconnectorrelationship.EnvironmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: environment.FieldID,
				},
			},
		}
		edge.Schema = ecrc.schemaConfig.EnvironmentConnectorRelationship
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.EnvironmentID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ecrc.mutation.ConnectorIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   environmentconnectorrelationship.ConnectorTable,
			Columns: []string{environmentconnectorrelationship.ConnectorColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: connector.FieldID,
				},
			},
		}
		edge.Schema = ecrc.schemaConfig.EnvironmentConnectorRelationship
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.ConnectorID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.EnvironmentConnectorRelationship.Create().
//		SetCreateTime(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.EnvironmentConnectorRelationshipUpsert) {
//			SetCreateTime(v+v).
//		}).
//		Exec(ctx)
func (ecrc *EnvironmentConnectorRelationshipCreate) OnConflict(opts ...sql.ConflictOption) *EnvironmentConnectorRelationshipUpsertOne {
	ecrc.conflict = opts
	return &EnvironmentConnectorRelationshipUpsertOne{
		create: ecrc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.EnvironmentConnectorRelationship.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (ecrc *EnvironmentConnectorRelationshipCreate) OnConflictColumns(columns ...string) *EnvironmentConnectorRelationshipUpsertOne {
	ecrc.conflict = append(ecrc.conflict, sql.ConflictColumns(columns...))
	return &EnvironmentConnectorRelationshipUpsertOne{
		create: ecrc,
	}
}

type (
	// EnvironmentConnectorRelationshipUpsertOne is the builder for "upsert"-ing
	//  one EnvironmentConnectorRelationship node.
	EnvironmentConnectorRelationshipUpsertOne struct {
		create *EnvironmentConnectorRelationshipCreate
	}

	// EnvironmentConnectorRelationshipUpsert is the "OnConflict" setter.
	EnvironmentConnectorRelationshipUpsert struct {
		*sql.UpdateSet
	}
)

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.EnvironmentConnectorRelationship.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *EnvironmentConnectorRelationshipUpsertOne) UpdateNewValues() *EnvironmentConnectorRelationshipUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.CreateTime(); exists {
			s.SetIgnore(environmentconnectorrelationship.FieldCreateTime)
		}
		if _, exists := u.create.mutation.EnvironmentID(); exists {
			s.SetIgnore(environmentconnectorrelationship.FieldEnvironmentID)
		}
		if _, exists := u.create.mutation.ConnectorID(); exists {
			s.SetIgnore(environmentconnectorrelationship.FieldConnectorID)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.EnvironmentConnectorRelationship.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *EnvironmentConnectorRelationshipUpsertOne) Ignore() *EnvironmentConnectorRelationshipUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *EnvironmentConnectorRelationshipUpsertOne) DoNothing() *EnvironmentConnectorRelationshipUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the EnvironmentConnectorRelationshipCreate.OnConflict
// documentation for more info.
func (u *EnvironmentConnectorRelationshipUpsertOne) Update(set func(*EnvironmentConnectorRelationshipUpsert)) *EnvironmentConnectorRelationshipUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&EnvironmentConnectorRelationshipUpsert{UpdateSet: update})
	}))
	return u
}

// Exec executes the query.
func (u *EnvironmentConnectorRelationshipUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("model: missing options for EnvironmentConnectorRelationshipCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *EnvironmentConnectorRelationshipUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// EnvironmentConnectorRelationshipCreateBulk is the builder for creating many EnvironmentConnectorRelationship entities in bulk.
type EnvironmentConnectorRelationshipCreateBulk struct {
	config
	builders []*EnvironmentConnectorRelationshipCreate
	conflict []sql.ConflictOption
}

// Save creates the EnvironmentConnectorRelationship entities in the database.
func (ecrcb *EnvironmentConnectorRelationshipCreateBulk) Save(ctx context.Context) ([]*EnvironmentConnectorRelationship, error) {
	specs := make([]*sqlgraph.CreateSpec, len(ecrcb.builders))
	nodes := make([]*EnvironmentConnectorRelationship, len(ecrcb.builders))
	mutators := make([]Mutator, len(ecrcb.builders))
	for i := range ecrcb.builders {
		func(i int, root context.Context) {
			builder := ecrcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*EnvironmentConnectorRelationshipMutation)
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
					_, err = mutators[i+1].Mutate(root, ecrcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = ecrcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ecrcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
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
		if _, err := mutators[0].Mutate(ctx, ecrcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ecrcb *EnvironmentConnectorRelationshipCreateBulk) SaveX(ctx context.Context) []*EnvironmentConnectorRelationship {
	v, err := ecrcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ecrcb *EnvironmentConnectorRelationshipCreateBulk) Exec(ctx context.Context) error {
	_, err := ecrcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ecrcb *EnvironmentConnectorRelationshipCreateBulk) ExecX(ctx context.Context) {
	if err := ecrcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.EnvironmentConnectorRelationship.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.EnvironmentConnectorRelationshipUpsert) {
//			SetCreateTime(v+v).
//		}).
//		Exec(ctx)
func (ecrcb *EnvironmentConnectorRelationshipCreateBulk) OnConflict(opts ...sql.ConflictOption) *EnvironmentConnectorRelationshipUpsertBulk {
	ecrcb.conflict = opts
	return &EnvironmentConnectorRelationshipUpsertBulk{
		create: ecrcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.EnvironmentConnectorRelationship.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (ecrcb *EnvironmentConnectorRelationshipCreateBulk) OnConflictColumns(columns ...string) *EnvironmentConnectorRelationshipUpsertBulk {
	ecrcb.conflict = append(ecrcb.conflict, sql.ConflictColumns(columns...))
	return &EnvironmentConnectorRelationshipUpsertBulk{
		create: ecrcb,
	}
}

// EnvironmentConnectorRelationshipUpsertBulk is the builder for "upsert"-ing
// a bulk of EnvironmentConnectorRelationship nodes.
type EnvironmentConnectorRelationshipUpsertBulk struct {
	create *EnvironmentConnectorRelationshipCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.EnvironmentConnectorRelationship.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *EnvironmentConnectorRelationshipUpsertBulk) UpdateNewValues() *EnvironmentConnectorRelationshipUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.CreateTime(); exists {
				s.SetIgnore(environmentconnectorrelationship.FieldCreateTime)
			}
			if _, exists := b.mutation.EnvironmentID(); exists {
				s.SetIgnore(environmentconnectorrelationship.FieldEnvironmentID)
			}
			if _, exists := b.mutation.ConnectorID(); exists {
				s.SetIgnore(environmentconnectorrelationship.FieldConnectorID)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.EnvironmentConnectorRelationship.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *EnvironmentConnectorRelationshipUpsertBulk) Ignore() *EnvironmentConnectorRelationshipUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *EnvironmentConnectorRelationshipUpsertBulk) DoNothing() *EnvironmentConnectorRelationshipUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the EnvironmentConnectorRelationshipCreateBulk.OnConflict
// documentation for more info.
func (u *EnvironmentConnectorRelationshipUpsertBulk) Update(set func(*EnvironmentConnectorRelationshipUpsert)) *EnvironmentConnectorRelationshipUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&EnvironmentConnectorRelationshipUpsert{UpdateSet: update})
	}))
	return u
}

// Exec executes the query.
func (u *EnvironmentConnectorRelationshipUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("model: OnConflict was set for builder %d. Set it on the EnvironmentConnectorRelationshipCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("model: missing options for EnvironmentConnectorRelationshipCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *EnvironmentConnectorRelationshipUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
