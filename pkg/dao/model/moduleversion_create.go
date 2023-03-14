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

	"github.com/seal-io/seal/pkg/dao/model/module"
	"github.com/seal-io/seal/pkg/dao/model/moduleversion"
	"github.com/seal-io/seal/pkg/dao/types"
)

// ModuleVersionCreate is the builder for creating a ModuleVersion entity.
type ModuleVersionCreate struct {
	config
	mutation *ModuleVersionMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreateTime sets the "createTime" field.
func (mvc *ModuleVersionCreate) SetCreateTime(t time.Time) *ModuleVersionCreate {
	mvc.mutation.SetCreateTime(t)
	return mvc
}

// SetNillableCreateTime sets the "createTime" field if the given value is not nil.
func (mvc *ModuleVersionCreate) SetNillableCreateTime(t *time.Time) *ModuleVersionCreate {
	if t != nil {
		mvc.SetCreateTime(*t)
	}
	return mvc
}

// SetUpdateTime sets the "updateTime" field.
func (mvc *ModuleVersionCreate) SetUpdateTime(t time.Time) *ModuleVersionCreate {
	mvc.mutation.SetUpdateTime(t)
	return mvc
}

// SetNillableUpdateTime sets the "updateTime" field if the given value is not nil.
func (mvc *ModuleVersionCreate) SetNillableUpdateTime(t *time.Time) *ModuleVersionCreate {
	if t != nil {
		mvc.SetUpdateTime(*t)
	}
	return mvc
}

// SetModuleID sets the "moduleID" field.
func (mvc *ModuleVersionCreate) SetModuleID(s string) *ModuleVersionCreate {
	mvc.mutation.SetModuleID(s)
	return mvc
}

// SetVersion sets the "version" field.
func (mvc *ModuleVersionCreate) SetVersion(s string) *ModuleVersionCreate {
	mvc.mutation.SetVersion(s)
	return mvc
}

// SetSource sets the "source" field.
func (mvc *ModuleVersionCreate) SetSource(s string) *ModuleVersionCreate {
	mvc.mutation.SetSource(s)
	return mvc
}

// SetSchema sets the "schema" field.
func (mvc *ModuleVersionCreate) SetSchema(ts *types.ModuleSchema) *ModuleVersionCreate {
	mvc.mutation.SetSchema(ts)
	return mvc
}

// SetID sets the "id" field.
func (mvc *ModuleVersionCreate) SetID(t types.ID) *ModuleVersionCreate {
	mvc.mutation.SetID(t)
	return mvc
}

// SetModule sets the "module" edge to the Module entity.
func (mvc *ModuleVersionCreate) SetModule(m *Module) *ModuleVersionCreate {
	return mvc.SetModuleID(m.ID)
}

// Mutation returns the ModuleVersionMutation object of the builder.
func (mvc *ModuleVersionCreate) Mutation() *ModuleVersionMutation {
	return mvc.mutation
}

// Save creates the ModuleVersion in the database.
func (mvc *ModuleVersionCreate) Save(ctx context.Context) (*ModuleVersion, error) {
	if err := mvc.defaults(); err != nil {
		return nil, err
	}
	return withHooks[*ModuleVersion, ModuleVersionMutation](ctx, mvc.sqlSave, mvc.mutation, mvc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (mvc *ModuleVersionCreate) SaveX(ctx context.Context) *ModuleVersion {
	v, err := mvc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (mvc *ModuleVersionCreate) Exec(ctx context.Context) error {
	_, err := mvc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mvc *ModuleVersionCreate) ExecX(ctx context.Context) {
	if err := mvc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (mvc *ModuleVersionCreate) defaults() error {
	if _, ok := mvc.mutation.CreateTime(); !ok {
		if moduleversion.DefaultCreateTime == nil {
			return fmt.Errorf("model: uninitialized moduleversion.DefaultCreateTime (forgotten import model/runtime?)")
		}
		v := moduleversion.DefaultCreateTime()
		mvc.mutation.SetCreateTime(v)
	}
	if _, ok := mvc.mutation.UpdateTime(); !ok {
		if moduleversion.DefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized moduleversion.DefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := moduleversion.DefaultUpdateTime()
		mvc.mutation.SetUpdateTime(v)
	}
	if _, ok := mvc.mutation.Schema(); !ok {
		v := moduleversion.DefaultSchema
		mvc.mutation.SetSchema(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (mvc *ModuleVersionCreate) check() error {
	if _, ok := mvc.mutation.CreateTime(); !ok {
		return &ValidationError{Name: "createTime", err: errors.New(`model: missing required field "ModuleVersion.createTime"`)}
	}
	if _, ok := mvc.mutation.UpdateTime(); !ok {
		return &ValidationError{Name: "updateTime", err: errors.New(`model: missing required field "ModuleVersion.updateTime"`)}
	}
	if _, ok := mvc.mutation.ModuleID(); !ok {
		return &ValidationError{Name: "moduleID", err: errors.New(`model: missing required field "ModuleVersion.moduleID"`)}
	}
	if v, ok := mvc.mutation.ModuleID(); ok {
		if err := moduleversion.ModuleIDValidator(v); err != nil {
			return &ValidationError{Name: "moduleID", err: fmt.Errorf(`model: validator failed for field "ModuleVersion.moduleID": %w`, err)}
		}
	}
	if _, ok := mvc.mutation.Version(); !ok {
		return &ValidationError{Name: "version", err: errors.New(`model: missing required field "ModuleVersion.version"`)}
	}
	if _, ok := mvc.mutation.Source(); !ok {
		return &ValidationError{Name: "source", err: errors.New(`model: missing required field "ModuleVersion.source"`)}
	}
	if _, ok := mvc.mutation.Schema(); !ok {
		return &ValidationError{Name: "schema", err: errors.New(`model: missing required field "ModuleVersion.schema"`)}
	}
	if _, ok := mvc.mutation.ModuleID(); !ok {
		return &ValidationError{Name: "module", err: errors.New(`model: missing required edge "ModuleVersion.module"`)}
	}
	return nil
}

func (mvc *ModuleVersionCreate) sqlSave(ctx context.Context) (*ModuleVersion, error) {
	if err := mvc.check(); err != nil {
		return nil, err
	}
	_node, _spec := mvc.createSpec()
	if err := sqlgraph.CreateNode(ctx, mvc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*types.ID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	mvc.mutation.id = &_node.ID
	mvc.mutation.done = true
	return _node, nil
}

func (mvc *ModuleVersionCreate) createSpec() (*ModuleVersion, *sqlgraph.CreateSpec) {
	var (
		_node = &ModuleVersion{config: mvc.config}
		_spec = sqlgraph.NewCreateSpec(moduleversion.Table, sqlgraph.NewFieldSpec(moduleversion.FieldID, field.TypeString))
	)
	_spec.Schema = mvc.schemaConfig.ModuleVersion
	_spec.OnConflict = mvc.conflict
	if id, ok := mvc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := mvc.mutation.CreateTime(); ok {
		_spec.SetField(moduleversion.FieldCreateTime, field.TypeTime, value)
		_node.CreateTime = &value
	}
	if value, ok := mvc.mutation.UpdateTime(); ok {
		_spec.SetField(moduleversion.FieldUpdateTime, field.TypeTime, value)
		_node.UpdateTime = &value
	}
	if value, ok := mvc.mutation.Version(); ok {
		_spec.SetField(moduleversion.FieldVersion, field.TypeString, value)
		_node.Version = value
	}
	if value, ok := mvc.mutation.Source(); ok {
		_spec.SetField(moduleversion.FieldSource, field.TypeString, value)
		_node.Source = value
	}
	if value, ok := mvc.mutation.Schema(); ok {
		_spec.SetField(moduleversion.FieldSchema, field.TypeJSON, value)
		_node.Schema = value
	}
	if nodes := mvc.mutation.ModuleIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   moduleversion.ModuleTable,
			Columns: []string{moduleversion.ModuleColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: module.FieldID,
				},
			},
		}
		edge.Schema = mvc.schemaConfig.ModuleVersion
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.ModuleID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.ModuleVersion.Create().
//		SetCreateTime(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.ModuleVersionUpsert) {
//			SetCreateTime(v+v).
//		}).
//		Exec(ctx)
func (mvc *ModuleVersionCreate) OnConflict(opts ...sql.ConflictOption) *ModuleVersionUpsertOne {
	mvc.conflict = opts
	return &ModuleVersionUpsertOne{
		create: mvc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.ModuleVersion.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (mvc *ModuleVersionCreate) OnConflictColumns(columns ...string) *ModuleVersionUpsertOne {
	mvc.conflict = append(mvc.conflict, sql.ConflictColumns(columns...))
	return &ModuleVersionUpsertOne{
		create: mvc,
	}
}

type (
	// ModuleVersionUpsertOne is the builder for "upsert"-ing
	//  one ModuleVersion node.
	ModuleVersionUpsertOne struct {
		create *ModuleVersionCreate
	}

	// ModuleVersionUpsert is the "OnConflict" setter.
	ModuleVersionUpsert struct {
		*sql.UpdateSet
	}
)

// SetUpdateTime sets the "updateTime" field.
func (u *ModuleVersionUpsert) SetUpdateTime(v time.Time) *ModuleVersionUpsert {
	u.Set(moduleversion.FieldUpdateTime, v)
	return u
}

// UpdateUpdateTime sets the "updateTime" field to the value that was provided on create.
func (u *ModuleVersionUpsert) UpdateUpdateTime() *ModuleVersionUpsert {
	u.SetExcluded(moduleversion.FieldUpdateTime)
	return u
}

// SetVersion sets the "version" field.
func (u *ModuleVersionUpsert) SetVersion(v string) *ModuleVersionUpsert {
	u.Set(moduleversion.FieldVersion, v)
	return u
}

// UpdateVersion sets the "version" field to the value that was provided on create.
func (u *ModuleVersionUpsert) UpdateVersion() *ModuleVersionUpsert {
	u.SetExcluded(moduleversion.FieldVersion)
	return u
}

// SetSource sets the "source" field.
func (u *ModuleVersionUpsert) SetSource(v string) *ModuleVersionUpsert {
	u.Set(moduleversion.FieldSource, v)
	return u
}

// UpdateSource sets the "source" field to the value that was provided on create.
func (u *ModuleVersionUpsert) UpdateSource() *ModuleVersionUpsert {
	u.SetExcluded(moduleversion.FieldSource)
	return u
}

// SetSchema sets the "schema" field.
func (u *ModuleVersionUpsert) SetSchema(v *types.ModuleSchema) *ModuleVersionUpsert {
	u.Set(moduleversion.FieldSchema, v)
	return u
}

// UpdateSchema sets the "schema" field to the value that was provided on create.
func (u *ModuleVersionUpsert) UpdateSchema() *ModuleVersionUpsert {
	u.SetExcluded(moduleversion.FieldSchema)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.ModuleVersion.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(moduleversion.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *ModuleVersionUpsertOne) UpdateNewValues() *ModuleVersionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(moduleversion.FieldID)
		}
		if _, exists := u.create.mutation.CreateTime(); exists {
			s.SetIgnore(moduleversion.FieldCreateTime)
		}
		if _, exists := u.create.mutation.ModuleID(); exists {
			s.SetIgnore(moduleversion.FieldModuleID)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.ModuleVersion.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *ModuleVersionUpsertOne) Ignore() *ModuleVersionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *ModuleVersionUpsertOne) DoNothing() *ModuleVersionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the ModuleVersionCreate.OnConflict
// documentation for more info.
func (u *ModuleVersionUpsertOne) Update(set func(*ModuleVersionUpsert)) *ModuleVersionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&ModuleVersionUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdateTime sets the "updateTime" field.
func (u *ModuleVersionUpsertOne) SetUpdateTime(v time.Time) *ModuleVersionUpsertOne {
	return u.Update(func(s *ModuleVersionUpsert) {
		s.SetUpdateTime(v)
	})
}

// UpdateUpdateTime sets the "updateTime" field to the value that was provided on create.
func (u *ModuleVersionUpsertOne) UpdateUpdateTime() *ModuleVersionUpsertOne {
	return u.Update(func(s *ModuleVersionUpsert) {
		s.UpdateUpdateTime()
	})
}

// SetVersion sets the "version" field.
func (u *ModuleVersionUpsertOne) SetVersion(v string) *ModuleVersionUpsertOne {
	return u.Update(func(s *ModuleVersionUpsert) {
		s.SetVersion(v)
	})
}

// UpdateVersion sets the "version" field to the value that was provided on create.
func (u *ModuleVersionUpsertOne) UpdateVersion() *ModuleVersionUpsertOne {
	return u.Update(func(s *ModuleVersionUpsert) {
		s.UpdateVersion()
	})
}

// SetSource sets the "source" field.
func (u *ModuleVersionUpsertOne) SetSource(v string) *ModuleVersionUpsertOne {
	return u.Update(func(s *ModuleVersionUpsert) {
		s.SetSource(v)
	})
}

// UpdateSource sets the "source" field to the value that was provided on create.
func (u *ModuleVersionUpsertOne) UpdateSource() *ModuleVersionUpsertOne {
	return u.Update(func(s *ModuleVersionUpsert) {
		s.UpdateSource()
	})
}

// SetSchema sets the "schema" field.
func (u *ModuleVersionUpsertOne) SetSchema(v *types.ModuleSchema) *ModuleVersionUpsertOne {
	return u.Update(func(s *ModuleVersionUpsert) {
		s.SetSchema(v)
	})
}

// UpdateSchema sets the "schema" field to the value that was provided on create.
func (u *ModuleVersionUpsertOne) UpdateSchema() *ModuleVersionUpsertOne {
	return u.Update(func(s *ModuleVersionUpsert) {
		s.UpdateSchema()
	})
}

// Exec executes the query.
func (u *ModuleVersionUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("model: missing options for ModuleVersionCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *ModuleVersionUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *ModuleVersionUpsertOne) ID(ctx context.Context) (id types.ID, err error) {
	if u.create.driver.Dialect() == dialect.MySQL {
		// In case of "ON CONFLICT", there is no way to get back non-numeric ID
		// fields from the database since MySQL does not support the RETURNING clause.
		return id, errors.New("model: ModuleVersionUpsertOne.ID is not supported by MySQL driver. Use ModuleVersionUpsertOne.Exec instead")
	}
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *ModuleVersionUpsertOne) IDX(ctx context.Context) types.ID {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// ModuleVersionCreateBulk is the builder for creating many ModuleVersion entities in bulk.
type ModuleVersionCreateBulk struct {
	config
	builders []*ModuleVersionCreate
	conflict []sql.ConflictOption
}

// Save creates the ModuleVersion entities in the database.
func (mvcb *ModuleVersionCreateBulk) Save(ctx context.Context) ([]*ModuleVersion, error) {
	specs := make([]*sqlgraph.CreateSpec, len(mvcb.builders))
	nodes := make([]*ModuleVersion, len(mvcb.builders))
	mutators := make([]Mutator, len(mvcb.builders))
	for i := range mvcb.builders {
		func(i int, root context.Context) {
			builder := mvcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ModuleVersionMutation)
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
					_, err = mutators[i+1].Mutate(root, mvcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = mvcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, mvcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, mvcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (mvcb *ModuleVersionCreateBulk) SaveX(ctx context.Context) []*ModuleVersion {
	v, err := mvcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (mvcb *ModuleVersionCreateBulk) Exec(ctx context.Context) error {
	_, err := mvcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mvcb *ModuleVersionCreateBulk) ExecX(ctx context.Context) {
	if err := mvcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.ModuleVersion.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.ModuleVersionUpsert) {
//			SetCreateTime(v+v).
//		}).
//		Exec(ctx)
func (mvcb *ModuleVersionCreateBulk) OnConflict(opts ...sql.ConflictOption) *ModuleVersionUpsertBulk {
	mvcb.conflict = opts
	return &ModuleVersionUpsertBulk{
		create: mvcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.ModuleVersion.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (mvcb *ModuleVersionCreateBulk) OnConflictColumns(columns ...string) *ModuleVersionUpsertBulk {
	mvcb.conflict = append(mvcb.conflict, sql.ConflictColumns(columns...))
	return &ModuleVersionUpsertBulk{
		create: mvcb,
	}
}

// ModuleVersionUpsertBulk is the builder for "upsert"-ing
// a bulk of ModuleVersion nodes.
type ModuleVersionUpsertBulk struct {
	create *ModuleVersionCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.ModuleVersion.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(moduleversion.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *ModuleVersionUpsertBulk) UpdateNewValues() *ModuleVersionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(moduleversion.FieldID)
			}
			if _, exists := b.mutation.CreateTime(); exists {
				s.SetIgnore(moduleversion.FieldCreateTime)
			}
			if _, exists := b.mutation.ModuleID(); exists {
				s.SetIgnore(moduleversion.FieldModuleID)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.ModuleVersion.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *ModuleVersionUpsertBulk) Ignore() *ModuleVersionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *ModuleVersionUpsertBulk) DoNothing() *ModuleVersionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the ModuleVersionCreateBulk.OnConflict
// documentation for more info.
func (u *ModuleVersionUpsertBulk) Update(set func(*ModuleVersionUpsert)) *ModuleVersionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&ModuleVersionUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdateTime sets the "updateTime" field.
func (u *ModuleVersionUpsertBulk) SetUpdateTime(v time.Time) *ModuleVersionUpsertBulk {
	return u.Update(func(s *ModuleVersionUpsert) {
		s.SetUpdateTime(v)
	})
}

// UpdateUpdateTime sets the "updateTime" field to the value that was provided on create.
func (u *ModuleVersionUpsertBulk) UpdateUpdateTime() *ModuleVersionUpsertBulk {
	return u.Update(func(s *ModuleVersionUpsert) {
		s.UpdateUpdateTime()
	})
}

// SetVersion sets the "version" field.
func (u *ModuleVersionUpsertBulk) SetVersion(v string) *ModuleVersionUpsertBulk {
	return u.Update(func(s *ModuleVersionUpsert) {
		s.SetVersion(v)
	})
}

// UpdateVersion sets the "version" field to the value that was provided on create.
func (u *ModuleVersionUpsertBulk) UpdateVersion() *ModuleVersionUpsertBulk {
	return u.Update(func(s *ModuleVersionUpsert) {
		s.UpdateVersion()
	})
}

// SetSource sets the "source" field.
func (u *ModuleVersionUpsertBulk) SetSource(v string) *ModuleVersionUpsertBulk {
	return u.Update(func(s *ModuleVersionUpsert) {
		s.SetSource(v)
	})
}

// UpdateSource sets the "source" field to the value that was provided on create.
func (u *ModuleVersionUpsertBulk) UpdateSource() *ModuleVersionUpsertBulk {
	return u.Update(func(s *ModuleVersionUpsert) {
		s.UpdateSource()
	})
}

// SetSchema sets the "schema" field.
func (u *ModuleVersionUpsertBulk) SetSchema(v *types.ModuleSchema) *ModuleVersionUpsertBulk {
	return u.Update(func(s *ModuleVersionUpsert) {
		s.SetSchema(v)
	})
}

// UpdateSchema sets the "schema" field to the value that was provided on create.
func (u *ModuleVersionUpsertBulk) UpdateSchema() *ModuleVersionUpsertBulk {
	return u.Update(func(s *ModuleVersionUpsert) {
		s.UpdateSchema()
	})
}

// Exec executes the query.
func (u *ModuleVersionUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("model: OnConflict was set for builder %d. Set it on the ModuleVersionCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("model: missing options for ModuleVersionCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *ModuleVersionUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
