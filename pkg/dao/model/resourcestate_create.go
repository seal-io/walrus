// SPDX-FileCopyrightText: 2024 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "walrus". DO NOT EDIT.

package model

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/model/resourcestate"
	"github.com/seal-io/walrus/pkg/dao/types/object"
)

// ResourceStateCreate is the builder for creating a ResourceState entity.
type ResourceStateCreate struct {
	config
	mutation   *ResourceStateMutation
	hooks      []Hook
	conflict   []sql.ConflictOption
	object     *ResourceState
	fromUpsert bool
}

// SetData sets the "data" field.
func (rsc *ResourceStateCreate) SetData(s string) *ResourceStateCreate {
	rsc.mutation.SetData(s)
	return rsc
}

// SetNillableData sets the "data" field if the given value is not nil.
func (rsc *ResourceStateCreate) SetNillableData(s *string) *ResourceStateCreate {
	if s != nil {
		rsc.SetData(*s)
	}
	return rsc
}

// SetResourceID sets the "resource_id" field.
func (rsc *ResourceStateCreate) SetResourceID(o object.ID) *ResourceStateCreate {
	rsc.mutation.SetResourceID(o)
	return rsc
}

// SetID sets the "id" field.
func (rsc *ResourceStateCreate) SetID(o object.ID) *ResourceStateCreate {
	rsc.mutation.SetID(o)
	return rsc
}

// SetResource sets the "resource" edge to the Resource entity.
func (rsc *ResourceStateCreate) SetResource(r *Resource) *ResourceStateCreate {
	return rsc.SetResourceID(r.ID)
}

// Mutation returns the ResourceStateMutation object of the builder.
func (rsc *ResourceStateCreate) Mutation() *ResourceStateMutation {
	return rsc.mutation
}

// Save creates the ResourceState in the database.
func (rsc *ResourceStateCreate) Save(ctx context.Context) (*ResourceState, error) {
	if err := rsc.defaults(); err != nil {
		return nil, err
	}
	return withHooks(ctx, rsc.sqlSave, rsc.mutation, rsc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (rsc *ResourceStateCreate) SaveX(ctx context.Context) *ResourceState {
	v, err := rsc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (rsc *ResourceStateCreate) Exec(ctx context.Context) error {
	_, err := rsc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (rsc *ResourceStateCreate) ExecX(ctx context.Context) {
	if err := rsc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (rsc *ResourceStateCreate) defaults() error {
	if _, ok := rsc.mutation.Data(); !ok {
		v := resourcestate.DefaultData
		rsc.mutation.SetData(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (rsc *ResourceStateCreate) check() error {
	if _, ok := rsc.mutation.Data(); !ok {
		return &ValidationError{Name: "data", err: errors.New(`model: missing required field "ResourceState.data"`)}
	}
	if _, ok := rsc.mutation.ResourceID(); !ok {
		return &ValidationError{Name: "resource_id", err: errors.New(`model: missing required field "ResourceState.resource_id"`)}
	}
	if v, ok := rsc.mutation.ResourceID(); ok {
		if err := resourcestate.ResourceIDValidator(string(v)); err != nil {
			return &ValidationError{Name: "resource_id", err: fmt.Errorf(`model: validator failed for field "ResourceState.resource_id": %w`, err)}
		}
	}
	if _, ok := rsc.mutation.ResourceID(); !ok {
		return &ValidationError{Name: "resource", err: errors.New(`model: missing required edge "ResourceState.resource"`)}
	}
	return nil
}

func (rsc *ResourceStateCreate) sqlSave(ctx context.Context) (*ResourceState, error) {
	if err := rsc.check(); err != nil {
		return nil, err
	}
	_node, _spec := rsc.createSpec()
	if err := sqlgraph.CreateNode(ctx, rsc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*object.ID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	rsc.mutation.id = &_node.ID
	rsc.mutation.done = true
	return _node, nil
}

func (rsc *ResourceStateCreate) createSpec() (*ResourceState, *sqlgraph.CreateSpec) {
	var (
		_node = &ResourceState{config: rsc.config}
		_spec = sqlgraph.NewCreateSpec(resourcestate.Table, sqlgraph.NewFieldSpec(resourcestate.FieldID, field.TypeString))
	)
	_spec.Schema = rsc.schemaConfig.ResourceState
	_spec.OnConflict = rsc.conflict
	if id, ok := rsc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := rsc.mutation.Data(); ok {
		_spec.SetField(resourcestate.FieldData, field.TypeString, value)
		_node.Data = value
	}
	if nodes := rsc.mutation.ResourceIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   resourcestate.ResourceTable,
			Columns: []string{resourcestate.ResourceColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(resource.FieldID, field.TypeString),
			},
		}
		edge.Schema = rsc.schemaConfig.ResourceState
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.ResourceID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// Set is different from other Set* methods,
// it sets the value by judging the definition of each field within the entire object.
//
// For required fields, Set calls directly.
//
// For optional fields, Set calls if the value is not zero.
//
// For example:
//
//	## Required
//
//	db.SetX(obj.X)
//
//	## Optional or Default
//
//	if _is_zero_value_(obj.X) {
//	   db.SetX(obj.X)
//	}
func (rsc *ResourceStateCreate) Set(obj *ResourceState) *ResourceStateCreate {
	// Required.
	rsc.SetData(obj.Data)
	rsc.SetResourceID(obj.ResourceID)

	// Optional.

	// Record the given object.
	rsc.object = obj

	return rsc
}

// getClientSet returns the ClientSet for the given builder.
func (rsc *ResourceStateCreate) getClientSet() (mc ClientSet) {
	if _, ok := rsc.config.driver.(*txDriver); ok {
		tx := &Tx{config: rsc.config}
		tx.init()
		mc = tx
	} else {
		cli := &Client{config: rsc.config}
		cli.init()
		mc = cli
	}
	return mc
}

// SaveE calls the given function after created the ResourceState entity,
// which is always good for cascading create operations.
func (rsc *ResourceStateCreate) SaveE(ctx context.Context, cbs ...func(ctx context.Context, mc ClientSet, created *ResourceState) error) (*ResourceState, error) {
	obj, err := rsc.Save(ctx)
	if err != nil {
		return nil, err
	}

	if len(cbs) == 0 {
		return obj, err
	}

	mc := rsc.getClientSet()

	if x := rsc.object; x != nil {
		if _, set := rsc.mutation.Field(resourcestate.FieldResourceID); set {
			obj.ResourceID = x.ResourceID
		}
		obj.Edges = x.Edges
	}

	for i := range cbs {
		if err = cbs[i](ctx, mc, obj); err != nil {
			return nil, err
		}
	}

	return obj, nil
}

// SaveEX is like SaveE, but panics if an error occurs.
func (rsc *ResourceStateCreate) SaveEX(ctx context.Context, cbs ...func(ctx context.Context, mc ClientSet, created *ResourceState) error) *ResourceState {
	obj, err := rsc.SaveE(ctx, cbs...)
	if err != nil {
		panic(err)
	}
	return obj
}

// ExecE calls the given function after executed the query,
// which is always good for cascading create operations.
func (rsc *ResourceStateCreate) ExecE(ctx context.Context, cbs ...func(ctx context.Context, mc ClientSet, created *ResourceState) error) error {
	_, err := rsc.SaveE(ctx, cbs...)
	return err
}

// ExecEX is like ExecE, but panics if an error occurs.
func (rsc *ResourceStateCreate) ExecEX(ctx context.Context, cbs ...func(ctx context.Context, mc ClientSet, created *ResourceState) error) {
	if err := rsc.ExecE(ctx, cbs...); err != nil {
		panic(err)
	}
}

// Set leverages the ResourceStateCreate Set method,
// it sets the value by judging the definition of each field within the entire item of the given list.
//
// For required fields, Set calls directly.
//
// For optional fields, Set calls if the value is not zero.
//
// For example:
//
//	## Required
//
//	db.SetX(obj.X)
//
//	## Optional or Default
//
//	if _is_zero_value_(obj.X) {
//	   db.SetX(obj.X)
//	}
func (rscb *ResourceStateCreateBulk) Set(objs ...*ResourceState) *ResourceStateCreateBulk {
	if len(objs) != 0 {
		client := NewResourceStateClient(rscb.config)

		rscb.builders = make([]*ResourceStateCreate, len(objs))
		for i := range objs {
			rscb.builders[i] = client.Create().Set(objs[i])
		}

		// Record the given objects.
		rscb.objects = objs
	}

	return rscb
}

// getClientSet returns the ClientSet for the given builder.
func (rscb *ResourceStateCreateBulk) getClientSet() (mc ClientSet) {
	if _, ok := rscb.config.driver.(*txDriver); ok {
		tx := &Tx{config: rscb.config}
		tx.init()
		mc = tx
	} else {
		cli := &Client{config: rscb.config}
		cli.init()
		mc = cli
	}
	return mc
}

// SaveE calls the given function after created the ResourceState entities,
// which is always good for cascading create operations.
func (rscb *ResourceStateCreateBulk) SaveE(ctx context.Context, cbs ...func(ctx context.Context, mc ClientSet, created *ResourceState) error) ([]*ResourceState, error) {
	objs, err := rscb.Save(ctx)
	if err != nil {
		return nil, err
	}

	if len(cbs) == 0 {
		return objs, err
	}

	mc := rscb.getClientSet()

	if x := rscb.objects; x != nil {
		for i := range x {
			if _, set := rscb.builders[i].mutation.Field(resourcestate.FieldResourceID); set {
				objs[i].ResourceID = x[i].ResourceID
			}
			objs[i].Edges = x[i].Edges
		}
	}

	for i := range objs {
		for j := range cbs {
			if err = cbs[j](ctx, mc, objs[i]); err != nil {
				return nil, err
			}
		}
	}

	return objs, nil
}

// SaveEX is like SaveE, but panics if an error occurs.
func (rscb *ResourceStateCreateBulk) SaveEX(ctx context.Context, cbs ...func(ctx context.Context, mc ClientSet, created *ResourceState) error) []*ResourceState {
	objs, err := rscb.SaveE(ctx, cbs...)
	if err != nil {
		panic(err)
	}
	return objs
}

// ExecE calls the given function after executed the query,
// which is always good for cascading create operations.
func (rscb *ResourceStateCreateBulk) ExecE(ctx context.Context, cbs ...func(ctx context.Context, mc ClientSet, created *ResourceState) error) error {
	_, err := rscb.SaveE(ctx, cbs...)
	return err
}

// ExecEX is like ExecE, but panics if an error occurs.
func (rscb *ResourceStateCreateBulk) ExecEX(ctx context.Context, cbs ...func(ctx context.Context, mc ClientSet, created *ResourceState) error) {
	if err := rscb.ExecE(ctx, cbs...); err != nil {
		panic(err)
	}
}

// ExecE calls the given function after executed the query,
// which is always good for cascading create operations.
func (u *ResourceStateUpsertOne) ExecE(ctx context.Context, cbs ...func(ctx context.Context, mc ClientSet, updated *ResourceState) error) error {
	if len(u.create.conflict) == 0 {
		return errors.New("model: missing options for ResourceStateUpsertOne.OnConflict")
	}
	u.create.fromUpsert = true
	return u.create.ExecE(ctx, cbs...)
}

// ExecEX is like ExecE, but panics if an error occurs.
func (u *ResourceStateUpsertOne) ExecEX(ctx context.Context, cbs ...func(ctx context.Context, mc ClientSet, updated *ResourceState) error) {
	if err := u.ExecE(ctx, cbs...); err != nil {
		panic(err)
	}
}

// ExecE calls the given function after executed the query,
// which is always good for cascading create operations.
func (u *ResourceStateUpsertBulk) ExecE(ctx context.Context, cbs ...func(ctx context.Context, mc ClientSet, updated *ResourceState) error) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("model: OnConflict was set for builder %d. Set it on the ResourceStateUpsertBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("model: missing options for ResourceStateUpsertBulk.OnConflict")
	}
	u.create.fromUpsert = true
	return u.create.ExecE(ctx, cbs...)
}

// ExecEX is like ExecE, but panics if an error occurs.
func (u *ResourceStateUpsertBulk) ExecEX(ctx context.Context, cbs ...func(ctx context.Context, mc ClientSet, updated *ResourceState) error) {
	if err := u.ExecE(ctx, cbs...); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.ResourceState.Create().
//		SetData(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.ResourceStateUpsert) {
//			SetData(v+v).
//		}).
//		Exec(ctx)
func (rsc *ResourceStateCreate) OnConflict(opts ...sql.ConflictOption) *ResourceStateUpsertOne {
	rsc.conflict = opts
	return &ResourceStateUpsertOne{
		create: rsc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.ResourceState.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (rsc *ResourceStateCreate) OnConflictColumns(columns ...string) *ResourceStateUpsertOne {
	rsc.conflict = append(rsc.conflict, sql.ConflictColumns(columns...))
	return &ResourceStateUpsertOne{
		create: rsc,
	}
}

type (
	// ResourceStateUpsertOne is the builder for "upsert"-ing
	//  one ResourceState node.
	ResourceStateUpsertOne struct {
		create *ResourceStateCreate
	}

	// ResourceStateUpsert is the "OnConflict" setter.
	ResourceStateUpsert struct {
		*sql.UpdateSet
	}
)

// SetData sets the "data" field.
func (u *ResourceStateUpsert) SetData(v string) *ResourceStateUpsert {
	u.Set(resourcestate.FieldData, v)
	return u
}

// UpdateData sets the "data" field to the value that was provided on create.
func (u *ResourceStateUpsert) UpdateData() *ResourceStateUpsert {
	u.SetExcluded(resourcestate.FieldData)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.ResourceState.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(resourcestate.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *ResourceStateUpsertOne) UpdateNewValues() *ResourceStateUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(resourcestate.FieldID)
		}
		if _, exists := u.create.mutation.ResourceID(); exists {
			s.SetIgnore(resourcestate.FieldResourceID)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.ResourceState.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *ResourceStateUpsertOne) Ignore() *ResourceStateUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *ResourceStateUpsertOne) DoNothing() *ResourceStateUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the ResourceStateCreate.OnConflict
// documentation for more info.
func (u *ResourceStateUpsertOne) Update(set func(*ResourceStateUpsert)) *ResourceStateUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&ResourceStateUpsert{UpdateSet: update})
	}))
	return u
}

// SetData sets the "data" field.
func (u *ResourceStateUpsertOne) SetData(v string) *ResourceStateUpsertOne {
	return u.Update(func(s *ResourceStateUpsert) {
		s.SetData(v)
	})
}

// UpdateData sets the "data" field to the value that was provided on create.
func (u *ResourceStateUpsertOne) UpdateData() *ResourceStateUpsertOne {
	return u.Update(func(s *ResourceStateUpsert) {
		s.UpdateData()
	})
}

// Exec executes the query.
func (u *ResourceStateUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("model: missing options for ResourceStateCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *ResourceStateUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *ResourceStateUpsertOne) ID(ctx context.Context) (id object.ID, err error) {
	if u.create.driver.Dialect() == dialect.MySQL {
		// In case of "ON CONFLICT", there is no way to get back non-numeric ID
		// fields from the database since MySQL does not support the RETURNING clause.
		return id, errors.New("model: ResourceStateUpsertOne.ID is not supported by MySQL driver. Use ResourceStateUpsertOne.Exec instead")
	}
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *ResourceStateUpsertOne) IDX(ctx context.Context) object.ID {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// ResourceStateCreateBulk is the builder for creating many ResourceState entities in bulk.
type ResourceStateCreateBulk struct {
	config
	err        error
	builders   []*ResourceStateCreate
	conflict   []sql.ConflictOption
	objects    []*ResourceState
	fromUpsert bool
}

// Save creates the ResourceState entities in the database.
func (rscb *ResourceStateCreateBulk) Save(ctx context.Context) ([]*ResourceState, error) {
	if rscb.err != nil {
		return nil, rscb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(rscb.builders))
	nodes := make([]*ResourceState, len(rscb.builders))
	mutators := make([]Mutator, len(rscb.builders))
	for i := range rscb.builders {
		func(i int, root context.Context) {
			builder := rscb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ResourceStateMutation)
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
					_, err = mutators[i+1].Mutate(root, rscb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = rscb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, rscb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, rscb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (rscb *ResourceStateCreateBulk) SaveX(ctx context.Context) []*ResourceState {
	v, err := rscb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (rscb *ResourceStateCreateBulk) Exec(ctx context.Context) error {
	_, err := rscb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (rscb *ResourceStateCreateBulk) ExecX(ctx context.Context) {
	if err := rscb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.ResourceState.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.ResourceStateUpsert) {
//			SetData(v+v).
//		}).
//		Exec(ctx)
func (rscb *ResourceStateCreateBulk) OnConflict(opts ...sql.ConflictOption) *ResourceStateUpsertBulk {
	rscb.conflict = opts
	return &ResourceStateUpsertBulk{
		create: rscb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.ResourceState.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (rscb *ResourceStateCreateBulk) OnConflictColumns(columns ...string) *ResourceStateUpsertBulk {
	rscb.conflict = append(rscb.conflict, sql.ConflictColumns(columns...))
	return &ResourceStateUpsertBulk{
		create: rscb,
	}
}

// ResourceStateUpsertBulk is the builder for "upsert"-ing
// a bulk of ResourceState nodes.
type ResourceStateUpsertBulk struct {
	create *ResourceStateCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.ResourceState.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(resourcestate.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *ResourceStateUpsertBulk) UpdateNewValues() *ResourceStateUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(resourcestate.FieldID)
			}
			if _, exists := b.mutation.ResourceID(); exists {
				s.SetIgnore(resourcestate.FieldResourceID)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.ResourceState.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *ResourceStateUpsertBulk) Ignore() *ResourceStateUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *ResourceStateUpsertBulk) DoNothing() *ResourceStateUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the ResourceStateCreateBulk.OnConflict
// documentation for more info.
func (u *ResourceStateUpsertBulk) Update(set func(*ResourceStateUpsert)) *ResourceStateUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&ResourceStateUpsert{UpdateSet: update})
	}))
	return u
}

// SetData sets the "data" field.
func (u *ResourceStateUpsertBulk) SetData(v string) *ResourceStateUpsertBulk {
	return u.Update(func(s *ResourceStateUpsert) {
		s.SetData(v)
	})
}

// UpdateData sets the "data" field to the value that was provided on create.
func (u *ResourceStateUpsertBulk) UpdateData() *ResourceStateUpsertBulk {
	return u.Update(func(s *ResourceStateUpsert) {
		s.UpdateData()
	})
}

// Exec executes the query.
func (u *ResourceStateUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("model: OnConflict was set for builder %d. Set it on the ResourceStateCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("model: missing options for ResourceStateCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *ResourceStateUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
