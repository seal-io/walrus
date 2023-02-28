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

	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/model/applicationmodulerelationship"
	"github.com/seal-io/seal/pkg/dao/model/module"
	"github.com/seal-io/seal/pkg/dao/types"
)

// ApplicationModuleRelationshipCreate is the builder for creating a ApplicationModuleRelationship entity.
type ApplicationModuleRelationshipCreate struct {
	config
	mutation *ApplicationModuleRelationshipMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreateTime sets the "createTime" field.
func (amrc *ApplicationModuleRelationshipCreate) SetCreateTime(t time.Time) *ApplicationModuleRelationshipCreate {
	amrc.mutation.SetCreateTime(t)
	return amrc
}

// SetNillableCreateTime sets the "createTime" field if the given value is not nil.
func (amrc *ApplicationModuleRelationshipCreate) SetNillableCreateTime(t *time.Time) *ApplicationModuleRelationshipCreate {
	if t != nil {
		amrc.SetCreateTime(*t)
	}
	return amrc
}

// SetUpdateTime sets the "updateTime" field.
func (amrc *ApplicationModuleRelationshipCreate) SetUpdateTime(t time.Time) *ApplicationModuleRelationshipCreate {
	amrc.mutation.SetUpdateTime(t)
	return amrc
}

// SetNillableUpdateTime sets the "updateTime" field if the given value is not nil.
func (amrc *ApplicationModuleRelationshipCreate) SetNillableUpdateTime(t *time.Time) *ApplicationModuleRelationshipCreate {
	if t != nil {
		amrc.SetUpdateTime(*t)
	}
	return amrc
}

// SetApplicationID sets the "application_id" field.
func (amrc *ApplicationModuleRelationshipCreate) SetApplicationID(t types.ID) *ApplicationModuleRelationshipCreate {
	amrc.mutation.SetApplicationID(t)
	return amrc
}

// SetModuleID sets the "module_id" field.
func (amrc *ApplicationModuleRelationshipCreate) SetModuleID(s string) *ApplicationModuleRelationshipCreate {
	amrc.mutation.SetModuleID(s)
	return amrc
}

// SetName sets the "name" field.
func (amrc *ApplicationModuleRelationshipCreate) SetName(s string) *ApplicationModuleRelationshipCreate {
	amrc.mutation.SetName(s)
	return amrc
}

// SetVariables sets the "variables" field.
func (amrc *ApplicationModuleRelationshipCreate) SetVariables(m map[string]interface{}) *ApplicationModuleRelationshipCreate {
	amrc.mutation.SetVariables(m)
	return amrc
}

// SetApplication sets the "application" edge to the Application entity.
func (amrc *ApplicationModuleRelationshipCreate) SetApplication(a *Application) *ApplicationModuleRelationshipCreate {
	return amrc.SetApplicationID(a.ID)
}

// SetModule sets the "module" edge to the Module entity.
func (amrc *ApplicationModuleRelationshipCreate) SetModule(m *Module) *ApplicationModuleRelationshipCreate {
	return amrc.SetModuleID(m.ID)
}

// Mutation returns the ApplicationModuleRelationshipMutation object of the builder.
func (amrc *ApplicationModuleRelationshipCreate) Mutation() *ApplicationModuleRelationshipMutation {
	return amrc.mutation
}

// Save creates the ApplicationModuleRelationship in the database.
func (amrc *ApplicationModuleRelationshipCreate) Save(ctx context.Context) (*ApplicationModuleRelationship, error) {
	amrc.defaults()
	return withHooks[*ApplicationModuleRelationship, ApplicationModuleRelationshipMutation](ctx, amrc.sqlSave, amrc.mutation, amrc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (amrc *ApplicationModuleRelationshipCreate) SaveX(ctx context.Context) *ApplicationModuleRelationship {
	v, err := amrc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (amrc *ApplicationModuleRelationshipCreate) Exec(ctx context.Context) error {
	_, err := amrc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (amrc *ApplicationModuleRelationshipCreate) ExecX(ctx context.Context) {
	if err := amrc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (amrc *ApplicationModuleRelationshipCreate) defaults() {
	if _, ok := amrc.mutation.CreateTime(); !ok {
		v := applicationmodulerelationship.DefaultCreateTime()
		amrc.mutation.SetCreateTime(v)
	}
	if _, ok := amrc.mutation.UpdateTime(); !ok {
		v := applicationmodulerelationship.DefaultUpdateTime()
		amrc.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (amrc *ApplicationModuleRelationshipCreate) check() error {
	if _, ok := amrc.mutation.CreateTime(); !ok {
		return &ValidationError{Name: "createTime", err: errors.New(`model: missing required field "ApplicationModuleRelationship.createTime"`)}
	}
	if _, ok := amrc.mutation.UpdateTime(); !ok {
		return &ValidationError{Name: "updateTime", err: errors.New(`model: missing required field "ApplicationModuleRelationship.updateTime"`)}
	}
	if _, ok := amrc.mutation.ApplicationID(); !ok {
		return &ValidationError{Name: "application_id", err: errors.New(`model: missing required field "ApplicationModuleRelationship.application_id"`)}
	}
	if v, ok := amrc.mutation.ApplicationID(); ok {
		if err := applicationmodulerelationship.ApplicationIDValidator(string(v)); err != nil {
			return &ValidationError{Name: "application_id", err: fmt.Errorf(`model: validator failed for field "ApplicationModuleRelationship.application_id": %w`, err)}
		}
	}
	if _, ok := amrc.mutation.ModuleID(); !ok {
		return &ValidationError{Name: "module_id", err: errors.New(`model: missing required field "ApplicationModuleRelationship.module_id"`)}
	}
	if v, ok := amrc.mutation.ModuleID(); ok {
		if err := applicationmodulerelationship.ModuleIDValidator(v); err != nil {
			return &ValidationError{Name: "module_id", err: fmt.Errorf(`model: validator failed for field "ApplicationModuleRelationship.module_id": %w`, err)}
		}
	}
	if _, ok := amrc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`model: missing required field "ApplicationModuleRelationship.name"`)}
	}
	if v, ok := amrc.mutation.Name(); ok {
		if err := applicationmodulerelationship.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`model: validator failed for field "ApplicationModuleRelationship.name": %w`, err)}
		}
	}
	if _, ok := amrc.mutation.ApplicationID(); !ok {
		return &ValidationError{Name: "application", err: errors.New(`model: missing required edge "ApplicationModuleRelationship.application"`)}
	}
	if _, ok := amrc.mutation.ModuleID(); !ok {
		return &ValidationError{Name: "module", err: errors.New(`model: missing required edge "ApplicationModuleRelationship.module"`)}
	}
	return nil
}

func (amrc *ApplicationModuleRelationshipCreate) sqlSave(ctx context.Context) (*ApplicationModuleRelationship, error) {
	if err := amrc.check(); err != nil {
		return nil, err
	}
	_node, _spec := amrc.createSpec()
	if err := sqlgraph.CreateNode(ctx, amrc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	return _node, nil
}

func (amrc *ApplicationModuleRelationshipCreate) createSpec() (*ApplicationModuleRelationship, *sqlgraph.CreateSpec) {
	var (
		_node = &ApplicationModuleRelationship{config: amrc.config}
		_spec = sqlgraph.NewCreateSpec(applicationmodulerelationship.Table, nil)
	)
	_spec.Schema = amrc.schemaConfig.ApplicationModuleRelationship
	_spec.OnConflict = amrc.conflict
	if value, ok := amrc.mutation.CreateTime(); ok {
		_spec.SetField(applicationmodulerelationship.FieldCreateTime, field.TypeTime, value)
		_node.CreateTime = &value
	}
	if value, ok := amrc.mutation.UpdateTime(); ok {
		_spec.SetField(applicationmodulerelationship.FieldUpdateTime, field.TypeTime, value)
		_node.UpdateTime = &value
	}
	if value, ok := amrc.mutation.Name(); ok {
		_spec.SetField(applicationmodulerelationship.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := amrc.mutation.Variables(); ok {
		_spec.SetField(applicationmodulerelationship.FieldVariables, field.TypeJSON, value)
		_node.Variables = value
	}
	if nodes := amrc.mutation.ApplicationIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   applicationmodulerelationship.ApplicationTable,
			Columns: []string{applicationmodulerelationship.ApplicationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: application.FieldID,
				},
			},
		}
		edge.Schema = amrc.schemaConfig.ApplicationModuleRelationship
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.ApplicationID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := amrc.mutation.ModuleIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   applicationmodulerelationship.ModuleTable,
			Columns: []string{applicationmodulerelationship.ModuleColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: module.FieldID,
				},
			},
		}
		edge.Schema = amrc.schemaConfig.ApplicationModuleRelationship
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
//	client.ApplicationModuleRelationship.Create().
//		SetCreateTime(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.ApplicationModuleRelationshipUpsert) {
//			SetCreateTime(v+v).
//		}).
//		Exec(ctx)
func (amrc *ApplicationModuleRelationshipCreate) OnConflict(opts ...sql.ConflictOption) *ApplicationModuleRelationshipUpsertOne {
	amrc.conflict = opts
	return &ApplicationModuleRelationshipUpsertOne{
		create: amrc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.ApplicationModuleRelationship.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (amrc *ApplicationModuleRelationshipCreate) OnConflictColumns(columns ...string) *ApplicationModuleRelationshipUpsertOne {
	amrc.conflict = append(amrc.conflict, sql.ConflictColumns(columns...))
	return &ApplicationModuleRelationshipUpsertOne{
		create: amrc,
	}
}

type (
	// ApplicationModuleRelationshipUpsertOne is the builder for "upsert"-ing
	//  one ApplicationModuleRelationship node.
	ApplicationModuleRelationshipUpsertOne struct {
		create *ApplicationModuleRelationshipCreate
	}

	// ApplicationModuleRelationshipUpsert is the "OnConflict" setter.
	ApplicationModuleRelationshipUpsert struct {
		*sql.UpdateSet
	}
)

// SetUpdateTime sets the "updateTime" field.
func (u *ApplicationModuleRelationshipUpsert) SetUpdateTime(v time.Time) *ApplicationModuleRelationshipUpsert {
	u.Set(applicationmodulerelationship.FieldUpdateTime, v)
	return u
}

// UpdateUpdateTime sets the "updateTime" field to the value that was provided on create.
func (u *ApplicationModuleRelationshipUpsert) UpdateUpdateTime() *ApplicationModuleRelationshipUpsert {
	u.SetExcluded(applicationmodulerelationship.FieldUpdateTime)
	return u
}

// SetVariables sets the "variables" field.
func (u *ApplicationModuleRelationshipUpsert) SetVariables(v map[string]interface{}) *ApplicationModuleRelationshipUpsert {
	u.Set(applicationmodulerelationship.FieldVariables, v)
	return u
}

// UpdateVariables sets the "variables" field to the value that was provided on create.
func (u *ApplicationModuleRelationshipUpsert) UpdateVariables() *ApplicationModuleRelationshipUpsert {
	u.SetExcluded(applicationmodulerelationship.FieldVariables)
	return u
}

// ClearVariables clears the value of the "variables" field.
func (u *ApplicationModuleRelationshipUpsert) ClearVariables() *ApplicationModuleRelationshipUpsert {
	u.SetNull(applicationmodulerelationship.FieldVariables)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.ApplicationModuleRelationship.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *ApplicationModuleRelationshipUpsertOne) UpdateNewValues() *ApplicationModuleRelationshipUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.CreateTime(); exists {
			s.SetIgnore(applicationmodulerelationship.FieldCreateTime)
		}
		if _, exists := u.create.mutation.ApplicationID(); exists {
			s.SetIgnore(applicationmodulerelationship.FieldApplicationID)
		}
		if _, exists := u.create.mutation.ModuleID(); exists {
			s.SetIgnore(applicationmodulerelationship.FieldModuleID)
		}
		if _, exists := u.create.mutation.Name(); exists {
			s.SetIgnore(applicationmodulerelationship.FieldName)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.ApplicationModuleRelationship.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *ApplicationModuleRelationshipUpsertOne) Ignore() *ApplicationModuleRelationshipUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *ApplicationModuleRelationshipUpsertOne) DoNothing() *ApplicationModuleRelationshipUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the ApplicationModuleRelationshipCreate.OnConflict
// documentation for more info.
func (u *ApplicationModuleRelationshipUpsertOne) Update(set func(*ApplicationModuleRelationshipUpsert)) *ApplicationModuleRelationshipUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&ApplicationModuleRelationshipUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdateTime sets the "updateTime" field.
func (u *ApplicationModuleRelationshipUpsertOne) SetUpdateTime(v time.Time) *ApplicationModuleRelationshipUpsertOne {
	return u.Update(func(s *ApplicationModuleRelationshipUpsert) {
		s.SetUpdateTime(v)
	})
}

// UpdateUpdateTime sets the "updateTime" field to the value that was provided on create.
func (u *ApplicationModuleRelationshipUpsertOne) UpdateUpdateTime() *ApplicationModuleRelationshipUpsertOne {
	return u.Update(func(s *ApplicationModuleRelationshipUpsert) {
		s.UpdateUpdateTime()
	})
}

// SetVariables sets the "variables" field.
func (u *ApplicationModuleRelationshipUpsertOne) SetVariables(v map[string]interface{}) *ApplicationModuleRelationshipUpsertOne {
	return u.Update(func(s *ApplicationModuleRelationshipUpsert) {
		s.SetVariables(v)
	})
}

// UpdateVariables sets the "variables" field to the value that was provided on create.
func (u *ApplicationModuleRelationshipUpsertOne) UpdateVariables() *ApplicationModuleRelationshipUpsertOne {
	return u.Update(func(s *ApplicationModuleRelationshipUpsert) {
		s.UpdateVariables()
	})
}

// ClearVariables clears the value of the "variables" field.
func (u *ApplicationModuleRelationshipUpsertOne) ClearVariables() *ApplicationModuleRelationshipUpsertOne {
	return u.Update(func(s *ApplicationModuleRelationshipUpsert) {
		s.ClearVariables()
	})
}

// Exec executes the query.
func (u *ApplicationModuleRelationshipUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("model: missing options for ApplicationModuleRelationshipCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *ApplicationModuleRelationshipUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// ApplicationModuleRelationshipCreateBulk is the builder for creating many ApplicationModuleRelationship entities in bulk.
type ApplicationModuleRelationshipCreateBulk struct {
	config
	builders []*ApplicationModuleRelationshipCreate
	conflict []sql.ConflictOption
}

// Save creates the ApplicationModuleRelationship entities in the database.
func (amrcb *ApplicationModuleRelationshipCreateBulk) Save(ctx context.Context) ([]*ApplicationModuleRelationship, error) {
	specs := make([]*sqlgraph.CreateSpec, len(amrcb.builders))
	nodes := make([]*ApplicationModuleRelationship, len(amrcb.builders))
	mutators := make([]Mutator, len(amrcb.builders))
	for i := range amrcb.builders {
		func(i int, root context.Context) {
			builder := amrcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ApplicationModuleRelationshipMutation)
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
					_, err = mutators[i+1].Mutate(root, amrcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = amrcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, amrcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, amrcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (amrcb *ApplicationModuleRelationshipCreateBulk) SaveX(ctx context.Context) []*ApplicationModuleRelationship {
	v, err := amrcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (amrcb *ApplicationModuleRelationshipCreateBulk) Exec(ctx context.Context) error {
	_, err := amrcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (amrcb *ApplicationModuleRelationshipCreateBulk) ExecX(ctx context.Context) {
	if err := amrcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.ApplicationModuleRelationship.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.ApplicationModuleRelationshipUpsert) {
//			SetCreateTime(v+v).
//		}).
//		Exec(ctx)
func (amrcb *ApplicationModuleRelationshipCreateBulk) OnConflict(opts ...sql.ConflictOption) *ApplicationModuleRelationshipUpsertBulk {
	amrcb.conflict = opts
	return &ApplicationModuleRelationshipUpsertBulk{
		create: amrcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.ApplicationModuleRelationship.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (amrcb *ApplicationModuleRelationshipCreateBulk) OnConflictColumns(columns ...string) *ApplicationModuleRelationshipUpsertBulk {
	amrcb.conflict = append(amrcb.conflict, sql.ConflictColumns(columns...))
	return &ApplicationModuleRelationshipUpsertBulk{
		create: amrcb,
	}
}

// ApplicationModuleRelationshipUpsertBulk is the builder for "upsert"-ing
// a bulk of ApplicationModuleRelationship nodes.
type ApplicationModuleRelationshipUpsertBulk struct {
	create *ApplicationModuleRelationshipCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.ApplicationModuleRelationship.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *ApplicationModuleRelationshipUpsertBulk) UpdateNewValues() *ApplicationModuleRelationshipUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.CreateTime(); exists {
				s.SetIgnore(applicationmodulerelationship.FieldCreateTime)
			}
			if _, exists := b.mutation.ApplicationID(); exists {
				s.SetIgnore(applicationmodulerelationship.FieldApplicationID)
			}
			if _, exists := b.mutation.ModuleID(); exists {
				s.SetIgnore(applicationmodulerelationship.FieldModuleID)
			}
			if _, exists := b.mutation.Name(); exists {
				s.SetIgnore(applicationmodulerelationship.FieldName)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.ApplicationModuleRelationship.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *ApplicationModuleRelationshipUpsertBulk) Ignore() *ApplicationModuleRelationshipUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *ApplicationModuleRelationshipUpsertBulk) DoNothing() *ApplicationModuleRelationshipUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the ApplicationModuleRelationshipCreateBulk.OnConflict
// documentation for more info.
func (u *ApplicationModuleRelationshipUpsertBulk) Update(set func(*ApplicationModuleRelationshipUpsert)) *ApplicationModuleRelationshipUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&ApplicationModuleRelationshipUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdateTime sets the "updateTime" field.
func (u *ApplicationModuleRelationshipUpsertBulk) SetUpdateTime(v time.Time) *ApplicationModuleRelationshipUpsertBulk {
	return u.Update(func(s *ApplicationModuleRelationshipUpsert) {
		s.SetUpdateTime(v)
	})
}

// UpdateUpdateTime sets the "updateTime" field to the value that was provided on create.
func (u *ApplicationModuleRelationshipUpsertBulk) UpdateUpdateTime() *ApplicationModuleRelationshipUpsertBulk {
	return u.Update(func(s *ApplicationModuleRelationshipUpsert) {
		s.UpdateUpdateTime()
	})
}

// SetVariables sets the "variables" field.
func (u *ApplicationModuleRelationshipUpsertBulk) SetVariables(v map[string]interface{}) *ApplicationModuleRelationshipUpsertBulk {
	return u.Update(func(s *ApplicationModuleRelationshipUpsert) {
		s.SetVariables(v)
	})
}

// UpdateVariables sets the "variables" field to the value that was provided on create.
func (u *ApplicationModuleRelationshipUpsertBulk) UpdateVariables() *ApplicationModuleRelationshipUpsertBulk {
	return u.Update(func(s *ApplicationModuleRelationshipUpsert) {
		s.UpdateVariables()
	})
}

// ClearVariables clears the value of the "variables" field.
func (u *ApplicationModuleRelationshipUpsertBulk) ClearVariables() *ApplicationModuleRelationshipUpsertBulk {
	return u.Update(func(s *ApplicationModuleRelationshipUpsert) {
		s.ClearVariables()
	})
}

// Exec executes the query.
func (u *ApplicationModuleRelationshipUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("model: OnConflict was set for builder %d. Set it on the ApplicationModuleRelationshipCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("model: missing options for ApplicationModuleRelationshipCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *ApplicationModuleRelationshipUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
