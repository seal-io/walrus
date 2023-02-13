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

	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/oid"
	"github.com/seal-io/seal/pkg/dao/schema"
)

// ApplicationCreate is the builder for creating a Application entity.
type ApplicationCreate struct {
	config
	mutation *ApplicationMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreateTime sets the "createTime" field.
func (ac *ApplicationCreate) SetCreateTime(t time.Time) *ApplicationCreate {
	ac.mutation.SetCreateTime(t)
	return ac
}

// SetNillableCreateTime sets the "createTime" field if the given value is not nil.
func (ac *ApplicationCreate) SetNillableCreateTime(t *time.Time) *ApplicationCreate {
	if t != nil {
		ac.SetCreateTime(*t)
	}
	return ac
}

// SetUpdateTime sets the "updateTime" field.
func (ac *ApplicationCreate) SetUpdateTime(t time.Time) *ApplicationCreate {
	ac.mutation.SetUpdateTime(t)
	return ac
}

// SetNillableUpdateTime sets the "updateTime" field if the given value is not nil.
func (ac *ApplicationCreate) SetNillableUpdateTime(t *time.Time) *ApplicationCreate {
	if t != nil {
		ac.SetUpdateTime(*t)
	}
	return ac
}

// SetProjectID sets the "projectID" field.
func (ac *ApplicationCreate) SetProjectID(o oid.ID) *ApplicationCreate {
	ac.mutation.SetProjectID(o)
	return ac
}

// SetEnvironmentID sets the "environmentID" field.
func (ac *ApplicationCreate) SetEnvironmentID(o oid.ID) *ApplicationCreate {
	ac.mutation.SetEnvironmentID(o)
	return ac
}

// SetModules sets the "modules" field.
func (ac *ApplicationCreate) SetModules(sm []schema.ApplicationModule) *ApplicationCreate {
	ac.mutation.SetModules(sm)
	return ac
}

// SetID sets the "id" field.
func (ac *ApplicationCreate) SetID(o oid.ID) *ApplicationCreate {
	ac.mutation.SetID(o)
	return ac
}

// Mutation returns the ApplicationMutation object of the builder.
func (ac *ApplicationCreate) Mutation() *ApplicationMutation {
	return ac.mutation
}

// Save creates the Application in the database.
func (ac *ApplicationCreate) Save(ctx context.Context) (*Application, error) {
	if err := ac.defaults(); err != nil {
		return nil, err
	}
	return withHooks[*Application, ApplicationMutation](ctx, ac.sqlSave, ac.mutation, ac.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (ac *ApplicationCreate) SaveX(ctx context.Context) *Application {
	v, err := ac.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ac *ApplicationCreate) Exec(ctx context.Context) error {
	_, err := ac.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ac *ApplicationCreate) ExecX(ctx context.Context) {
	if err := ac.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ac *ApplicationCreate) defaults() error {
	if _, ok := ac.mutation.CreateTime(); !ok {
		if application.DefaultCreateTime == nil {
			return fmt.Errorf("model: uninitialized application.DefaultCreateTime (forgotten import model/runtime?)")
		}
		v := application.DefaultCreateTime()
		ac.mutation.SetCreateTime(v)
	}
	if _, ok := ac.mutation.UpdateTime(); !ok {
		if application.DefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized application.DefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := application.DefaultUpdateTime()
		ac.mutation.SetUpdateTime(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (ac *ApplicationCreate) check() error {
	if _, ok := ac.mutation.CreateTime(); !ok {
		return &ValidationError{Name: "createTime", err: errors.New(`model: missing required field "Application.createTime"`)}
	}
	if _, ok := ac.mutation.UpdateTime(); !ok {
		return &ValidationError{Name: "updateTime", err: errors.New(`model: missing required field "Application.updateTime"`)}
	}
	if _, ok := ac.mutation.ProjectID(); !ok {
		return &ValidationError{Name: "projectID", err: errors.New(`model: missing required field "Application.projectID"`)}
	}
	if _, ok := ac.mutation.EnvironmentID(); !ok {
		return &ValidationError{Name: "environmentID", err: errors.New(`model: missing required field "Application.environmentID"`)}
	}
	if _, ok := ac.mutation.Modules(); !ok {
		return &ValidationError{Name: "modules", err: errors.New(`model: missing required field "Application.modules"`)}
	}
	return nil
}

func (ac *ApplicationCreate) sqlSave(ctx context.Context) (*Application, error) {
	if err := ac.check(); err != nil {
		return nil, err
	}
	_node, _spec := ac.createSpec()
	if err := sqlgraph.CreateNode(ctx, ac.driver, _spec); err != nil {
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
	ac.mutation.id = &_node.ID
	ac.mutation.done = true
	return _node, nil
}

func (ac *ApplicationCreate) createSpec() (*Application, *sqlgraph.CreateSpec) {
	var (
		_node = &Application{config: ac.config}
		_spec = &sqlgraph.CreateSpec{
			Table: application.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeOther,
				Column: application.FieldID,
			},
		}
	)
	_spec.Schema = ac.schemaConfig.Application
	_spec.OnConflict = ac.conflict
	if id, ok := ac.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := ac.mutation.CreateTime(); ok {
		_spec.SetField(application.FieldCreateTime, field.TypeTime, value)
		_node.CreateTime = &value
	}
	if value, ok := ac.mutation.UpdateTime(); ok {
		_spec.SetField(application.FieldUpdateTime, field.TypeTime, value)
		_node.UpdateTime = &value
	}
	if value, ok := ac.mutation.ProjectID(); ok {
		_spec.SetField(application.FieldProjectID, field.TypeOther, value)
		_node.ProjectID = value
	}
	if value, ok := ac.mutation.EnvironmentID(); ok {
		_spec.SetField(application.FieldEnvironmentID, field.TypeOther, value)
		_node.EnvironmentID = value
	}
	if value, ok := ac.mutation.Modules(); ok {
		_spec.SetField(application.FieldModules, field.TypeJSON, value)
		_node.Modules = value
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Application.Create().
//		SetCreateTime(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.ApplicationUpsert) {
//			SetCreateTime(v+v).
//		}).
//		Exec(ctx)
func (ac *ApplicationCreate) OnConflict(opts ...sql.ConflictOption) *ApplicationUpsertOne {
	ac.conflict = opts
	return &ApplicationUpsertOne{
		create: ac,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Application.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (ac *ApplicationCreate) OnConflictColumns(columns ...string) *ApplicationUpsertOne {
	ac.conflict = append(ac.conflict, sql.ConflictColumns(columns...))
	return &ApplicationUpsertOne{
		create: ac,
	}
}

type (
	// ApplicationUpsertOne is the builder for "upsert"-ing
	//  one Application node.
	ApplicationUpsertOne struct {
		create *ApplicationCreate
	}

	// ApplicationUpsert is the "OnConflict" setter.
	ApplicationUpsert struct {
		*sql.UpdateSet
	}
)

// SetUpdateTime sets the "updateTime" field.
func (u *ApplicationUpsert) SetUpdateTime(v time.Time) *ApplicationUpsert {
	u.Set(application.FieldUpdateTime, v)
	return u
}

// UpdateUpdateTime sets the "updateTime" field to the value that was provided on create.
func (u *ApplicationUpsert) UpdateUpdateTime() *ApplicationUpsert {
	u.SetExcluded(application.FieldUpdateTime)
	return u
}

// SetModules sets the "modules" field.
func (u *ApplicationUpsert) SetModules(v []schema.ApplicationModule) *ApplicationUpsert {
	u.Set(application.FieldModules, v)
	return u
}

// UpdateModules sets the "modules" field to the value that was provided on create.
func (u *ApplicationUpsert) UpdateModules() *ApplicationUpsert {
	u.SetExcluded(application.FieldModules)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.Application.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(application.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *ApplicationUpsertOne) UpdateNewValues() *ApplicationUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(application.FieldID)
		}
		if _, exists := u.create.mutation.CreateTime(); exists {
			s.SetIgnore(application.FieldCreateTime)
		}
		if _, exists := u.create.mutation.ProjectID(); exists {
			s.SetIgnore(application.FieldProjectID)
		}
		if _, exists := u.create.mutation.EnvironmentID(); exists {
			s.SetIgnore(application.FieldEnvironmentID)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Application.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *ApplicationUpsertOne) Ignore() *ApplicationUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *ApplicationUpsertOne) DoNothing() *ApplicationUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the ApplicationCreate.OnConflict
// documentation for more info.
func (u *ApplicationUpsertOne) Update(set func(*ApplicationUpsert)) *ApplicationUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&ApplicationUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdateTime sets the "updateTime" field.
func (u *ApplicationUpsertOne) SetUpdateTime(v time.Time) *ApplicationUpsertOne {
	return u.Update(func(s *ApplicationUpsert) {
		s.SetUpdateTime(v)
	})
}

// UpdateUpdateTime sets the "updateTime" field to the value that was provided on create.
func (u *ApplicationUpsertOne) UpdateUpdateTime() *ApplicationUpsertOne {
	return u.Update(func(s *ApplicationUpsert) {
		s.UpdateUpdateTime()
	})
}

// SetModules sets the "modules" field.
func (u *ApplicationUpsertOne) SetModules(v []schema.ApplicationModule) *ApplicationUpsertOne {
	return u.Update(func(s *ApplicationUpsert) {
		s.SetModules(v)
	})
}

// UpdateModules sets the "modules" field to the value that was provided on create.
func (u *ApplicationUpsertOne) UpdateModules() *ApplicationUpsertOne {
	return u.Update(func(s *ApplicationUpsert) {
		s.UpdateModules()
	})
}

// Exec executes the query.
func (u *ApplicationUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("model: missing options for ApplicationCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *ApplicationUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *ApplicationUpsertOne) ID(ctx context.Context) (id oid.ID, err error) {
	if u.create.driver.Dialect() == dialect.MySQL {
		// In case of "ON CONFLICT", there is no way to get back non-numeric ID
		// fields from the database since MySQL does not support the RETURNING clause.
		return id, errors.New("model: ApplicationUpsertOne.ID is not supported by MySQL driver. Use ApplicationUpsertOne.Exec instead")
	}
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *ApplicationUpsertOne) IDX(ctx context.Context) oid.ID {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// ApplicationCreateBulk is the builder for creating many Application entities in bulk.
type ApplicationCreateBulk struct {
	config
	builders []*ApplicationCreate
	conflict []sql.ConflictOption
}

// Save creates the Application entities in the database.
func (acb *ApplicationCreateBulk) Save(ctx context.Context) ([]*Application, error) {
	specs := make([]*sqlgraph.CreateSpec, len(acb.builders))
	nodes := make([]*Application, len(acb.builders))
	mutators := make([]Mutator, len(acb.builders))
	for i := range acb.builders {
		func(i int, root context.Context) {
			builder := acb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ApplicationMutation)
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
					_, err = mutators[i+1].Mutate(root, acb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = acb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, acb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, acb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (acb *ApplicationCreateBulk) SaveX(ctx context.Context) []*Application {
	v, err := acb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (acb *ApplicationCreateBulk) Exec(ctx context.Context) error {
	_, err := acb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (acb *ApplicationCreateBulk) ExecX(ctx context.Context) {
	if err := acb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Application.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.ApplicationUpsert) {
//			SetCreateTime(v+v).
//		}).
//		Exec(ctx)
func (acb *ApplicationCreateBulk) OnConflict(opts ...sql.ConflictOption) *ApplicationUpsertBulk {
	acb.conflict = opts
	return &ApplicationUpsertBulk{
		create: acb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Application.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (acb *ApplicationCreateBulk) OnConflictColumns(columns ...string) *ApplicationUpsertBulk {
	acb.conflict = append(acb.conflict, sql.ConflictColumns(columns...))
	return &ApplicationUpsertBulk{
		create: acb,
	}
}

// ApplicationUpsertBulk is the builder for "upsert"-ing
// a bulk of Application nodes.
type ApplicationUpsertBulk struct {
	create *ApplicationCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Application.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(application.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *ApplicationUpsertBulk) UpdateNewValues() *ApplicationUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(application.FieldID)
			}
			if _, exists := b.mutation.CreateTime(); exists {
				s.SetIgnore(application.FieldCreateTime)
			}
			if _, exists := b.mutation.ProjectID(); exists {
				s.SetIgnore(application.FieldProjectID)
			}
			if _, exists := b.mutation.EnvironmentID(); exists {
				s.SetIgnore(application.FieldEnvironmentID)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Application.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *ApplicationUpsertBulk) Ignore() *ApplicationUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *ApplicationUpsertBulk) DoNothing() *ApplicationUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the ApplicationCreateBulk.OnConflict
// documentation for more info.
func (u *ApplicationUpsertBulk) Update(set func(*ApplicationUpsert)) *ApplicationUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&ApplicationUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdateTime sets the "updateTime" field.
func (u *ApplicationUpsertBulk) SetUpdateTime(v time.Time) *ApplicationUpsertBulk {
	return u.Update(func(s *ApplicationUpsert) {
		s.SetUpdateTime(v)
	})
}

// UpdateUpdateTime sets the "updateTime" field to the value that was provided on create.
func (u *ApplicationUpsertBulk) UpdateUpdateTime() *ApplicationUpsertBulk {
	return u.Update(func(s *ApplicationUpsert) {
		s.UpdateUpdateTime()
	})
}

// SetModules sets the "modules" field.
func (u *ApplicationUpsertBulk) SetModules(v []schema.ApplicationModule) *ApplicationUpsertBulk {
	return u.Update(func(s *ApplicationUpsert) {
		s.SetModules(v)
	})
}

// UpdateModules sets the "modules" field to the value that was provided on create.
func (u *ApplicationUpsertBulk) UpdateModules() *ApplicationUpsertBulk {
	return u.Update(func(s *ApplicationUpsert) {
		s.UpdateModules()
	})
}

// Exec executes the query.
func (u *ApplicationUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("model: OnConflict was set for builder %d. Set it on the ApplicationCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("model: missing options for ApplicationCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *ApplicationUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
