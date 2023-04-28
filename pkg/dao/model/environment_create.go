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

	"github.com/seal-io/seal/pkg/dao/model/applicationinstance"
	"github.com/seal-io/seal/pkg/dao/model/applicationrevision"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/types/oid"
)

// EnvironmentCreate is the builder for creating a Environment entity.
type EnvironmentCreate struct {
	config
	mutation *EnvironmentMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetName sets the "name" field.
func (ec *EnvironmentCreate) SetName(s string) *EnvironmentCreate {
	ec.mutation.SetName(s)
	return ec
}

// SetDescription sets the "description" field.
func (ec *EnvironmentCreate) SetDescription(s string) *EnvironmentCreate {
	ec.mutation.SetDescription(s)
	return ec
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (ec *EnvironmentCreate) SetNillableDescription(s *string) *EnvironmentCreate {
	if s != nil {
		ec.SetDescription(*s)
	}
	return ec
}

// SetLabels sets the "labels" field.
func (ec *EnvironmentCreate) SetLabels(m map[string]string) *EnvironmentCreate {
	ec.mutation.SetLabels(m)
	return ec
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

// SetID sets the "id" field.
func (ec *EnvironmentCreate) SetID(o oid.ID) *EnvironmentCreate {
	ec.mutation.SetID(o)
	return ec
}

// AddInstanceIDs adds the "instances" edge to the ApplicationInstance entity by IDs.
func (ec *EnvironmentCreate) AddInstanceIDs(ids ...oid.ID) *EnvironmentCreate {
	ec.mutation.AddInstanceIDs(ids...)
	return ec
}

// AddInstances adds the "instances" edges to the ApplicationInstance entity.
func (ec *EnvironmentCreate) AddInstances(a ...*ApplicationInstance) *EnvironmentCreate {
	ids := make([]oid.ID, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return ec.AddInstanceIDs(ids...)
}

// AddRevisionIDs adds the "revisions" edge to the ApplicationRevision entity by IDs.
func (ec *EnvironmentCreate) AddRevisionIDs(ids ...oid.ID) *EnvironmentCreate {
	ec.mutation.AddRevisionIDs(ids...)
	return ec
}

// AddRevisions adds the "revisions" edges to the ApplicationRevision entity.
func (ec *EnvironmentCreate) AddRevisions(a ...*ApplicationRevision) *EnvironmentCreate {
	ids := make([]oid.ID, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return ec.AddRevisionIDs(ids...)
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
	if _, ok := ec.mutation.Labels(); !ok {
		v := environment.DefaultLabels
		ec.mutation.SetLabels(v)
	}
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
	if _, ok := ec.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`model: missing required field "Environment.name"`)}
	}
	if v, ok := ec.mutation.Name(); ok {
		if err := environment.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`model: validator failed for field "Environment.name": %w`, err)}
		}
	}
	if _, ok := ec.mutation.Labels(); !ok {
		return &ValidationError{Name: "labels", err: errors.New(`model: missing required field "Environment.labels"`)}
	}
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
		_spec = sqlgraph.NewCreateSpec(environment.Table, sqlgraph.NewFieldSpec(environment.FieldID, field.TypeString))
	)
	_spec.Schema = ec.schemaConfig.Environment
	_spec.OnConflict = ec.conflict
	if id, ok := ec.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := ec.mutation.Name(); ok {
		_spec.SetField(environment.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := ec.mutation.Description(); ok {
		_spec.SetField(environment.FieldDescription, field.TypeString, value)
		_node.Description = value
	}
	if value, ok := ec.mutation.Labels(); ok {
		_spec.SetField(environment.FieldLabels, field.TypeJSON, value)
		_node.Labels = value
	}
	if value, ok := ec.mutation.CreateTime(); ok {
		_spec.SetField(environment.FieldCreateTime, field.TypeTime, value)
		_node.CreateTime = &value
	}
	if value, ok := ec.mutation.UpdateTime(); ok {
		_spec.SetField(environment.FieldUpdateTime, field.TypeTime, value)
		_node.UpdateTime = &value
	}
	if nodes := ec.mutation.InstancesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   environment.InstancesTable,
			Columns: []string{environment.InstancesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(applicationinstance.FieldID, field.TypeString),
			},
		}
		edge.Schema = ec.schemaConfig.ApplicationInstance
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ec.mutation.RevisionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   environment.RevisionsTable,
			Columns: []string{environment.RevisionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(applicationrevision.FieldID, field.TypeString),
			},
		}
		edge.Schema = ec.schemaConfig.ApplicationRevision
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Environment.Create().
//		SetName(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.EnvironmentUpsert) {
//			SetName(v+v).
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

// SetName sets the "name" field.
func (u *EnvironmentUpsert) SetName(v string) *EnvironmentUpsert {
	u.Set(environment.FieldName, v)
	return u
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *EnvironmentUpsert) UpdateName() *EnvironmentUpsert {
	u.SetExcluded(environment.FieldName)
	return u
}

// SetDescription sets the "description" field.
func (u *EnvironmentUpsert) SetDescription(v string) *EnvironmentUpsert {
	u.Set(environment.FieldDescription, v)
	return u
}

// UpdateDescription sets the "description" field to the value that was provided on create.
func (u *EnvironmentUpsert) UpdateDescription() *EnvironmentUpsert {
	u.SetExcluded(environment.FieldDescription)
	return u
}

// ClearDescription clears the value of the "description" field.
func (u *EnvironmentUpsert) ClearDescription() *EnvironmentUpsert {
	u.SetNull(environment.FieldDescription)
	return u
}

// SetLabels sets the "labels" field.
func (u *EnvironmentUpsert) SetLabels(v map[string]string) *EnvironmentUpsert {
	u.Set(environment.FieldLabels, v)
	return u
}

// UpdateLabels sets the "labels" field to the value that was provided on create.
func (u *EnvironmentUpsert) UpdateLabels() *EnvironmentUpsert {
	u.SetExcluded(environment.FieldLabels)
	return u
}

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

// SetName sets the "name" field.
func (u *EnvironmentUpsertOne) SetName(v string) *EnvironmentUpsertOne {
	return u.Update(func(s *EnvironmentUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *EnvironmentUpsertOne) UpdateName() *EnvironmentUpsertOne {
	return u.Update(func(s *EnvironmentUpsert) {
		s.UpdateName()
	})
}

// SetDescription sets the "description" field.
func (u *EnvironmentUpsertOne) SetDescription(v string) *EnvironmentUpsertOne {
	return u.Update(func(s *EnvironmentUpsert) {
		s.SetDescription(v)
	})
}

// UpdateDescription sets the "description" field to the value that was provided on create.
func (u *EnvironmentUpsertOne) UpdateDescription() *EnvironmentUpsertOne {
	return u.Update(func(s *EnvironmentUpsert) {
		s.UpdateDescription()
	})
}

// ClearDescription clears the value of the "description" field.
func (u *EnvironmentUpsertOne) ClearDescription() *EnvironmentUpsertOne {
	return u.Update(func(s *EnvironmentUpsert) {
		s.ClearDescription()
	})
}

// SetLabels sets the "labels" field.
func (u *EnvironmentUpsertOne) SetLabels(v map[string]string) *EnvironmentUpsertOne {
	return u.Update(func(s *EnvironmentUpsert) {
		s.SetLabels(v)
	})
}

// UpdateLabels sets the "labels" field to the value that was provided on create.
func (u *EnvironmentUpsertOne) UpdateLabels() *EnvironmentUpsertOne {
	return u.Update(func(s *EnvironmentUpsert) {
		s.UpdateLabels()
	})
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
				var err error
				nodes[i], specs[i] = builder.createSpec()
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
//			SetName(v+v).
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

// SetName sets the "name" field.
func (u *EnvironmentUpsertBulk) SetName(v string) *EnvironmentUpsertBulk {
	return u.Update(func(s *EnvironmentUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *EnvironmentUpsertBulk) UpdateName() *EnvironmentUpsertBulk {
	return u.Update(func(s *EnvironmentUpsert) {
		s.UpdateName()
	})
}

// SetDescription sets the "description" field.
func (u *EnvironmentUpsertBulk) SetDescription(v string) *EnvironmentUpsertBulk {
	return u.Update(func(s *EnvironmentUpsert) {
		s.SetDescription(v)
	})
}

// UpdateDescription sets the "description" field to the value that was provided on create.
func (u *EnvironmentUpsertBulk) UpdateDescription() *EnvironmentUpsertBulk {
	return u.Update(func(s *EnvironmentUpsert) {
		s.UpdateDescription()
	})
}

// ClearDescription clears the value of the "description" field.
func (u *EnvironmentUpsertBulk) ClearDescription() *EnvironmentUpsertBulk {
	return u.Update(func(s *EnvironmentUpsert) {
		s.ClearDescription()
	})
}

// SetLabels sets the "labels" field.
func (u *EnvironmentUpsertBulk) SetLabels(v map[string]string) *EnvironmentUpsertBulk {
	return u.Update(func(s *EnvironmentUpsert) {
		s.SetLabels(v)
	})
}

// UpdateLabels sets the "labels" field to the value that was provided on create.
func (u *EnvironmentUpsertBulk) UpdateLabels() *EnvironmentUpsertBulk {
	return u.Update(func(s *EnvironmentUpsert) {
		s.UpdateLabels()
	})
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
