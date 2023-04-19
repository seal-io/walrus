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
	"github.com/seal-io/seal/pkg/dao/model/applicationinstance"
	"github.com/seal-io/seal/pkg/dao/model/project"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/pkg/dao/types/property"
)

// ApplicationCreate is the builder for creating a Application entity.
type ApplicationCreate struct {
	config
	mutation *ApplicationMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetName sets the "name" field.
func (ac *ApplicationCreate) SetName(s string) *ApplicationCreate {
	ac.mutation.SetName(s)
	return ac
}

// SetDescription sets the "description" field.
func (ac *ApplicationCreate) SetDescription(s string) *ApplicationCreate {
	ac.mutation.SetDescription(s)
	return ac
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (ac *ApplicationCreate) SetNillableDescription(s *string) *ApplicationCreate {
	if s != nil {
		ac.SetDescription(*s)
	}
	return ac
}

// SetLabels sets the "labels" field.
func (ac *ApplicationCreate) SetLabels(m map[string]string) *ApplicationCreate {
	ac.mutation.SetLabels(m)
	return ac
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

// SetVariables sets the "variables" field.
func (ac *ApplicationCreate) SetVariables(pr property.Schemas) *ApplicationCreate {
	ac.mutation.SetVariables(pr)
	return ac
}

// SetID sets the "id" field.
func (ac *ApplicationCreate) SetID(o oid.ID) *ApplicationCreate {
	ac.mutation.SetID(o)
	return ac
}

// SetProject sets the "project" edge to the Project entity.
func (ac *ApplicationCreate) SetProject(p *Project) *ApplicationCreate {
	return ac.SetProjectID(p.ID)
}

// AddInstanceIDs adds the "instances" edge to the ApplicationInstance entity by IDs.
func (ac *ApplicationCreate) AddInstanceIDs(ids ...oid.ID) *ApplicationCreate {
	ac.mutation.AddInstanceIDs(ids...)
	return ac
}

// AddInstances adds the "instances" edges to the ApplicationInstance entity.
func (ac *ApplicationCreate) AddInstances(a ...*ApplicationInstance) *ApplicationCreate {
	ids := make([]oid.ID, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return ac.AddInstanceIDs(ids...)
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
	if _, ok := ac.mutation.Labels(); !ok {
		v := application.DefaultLabels
		ac.mutation.SetLabels(v)
	}
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
	if _, ok := ac.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`model: missing required field "Application.name"`)}
	}
	if v, ok := ac.mutation.Name(); ok {
		if err := application.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`model: validator failed for field "Application.name": %w`, err)}
		}
	}
	if _, ok := ac.mutation.Labels(); !ok {
		return &ValidationError{Name: "labels", err: errors.New(`model: missing required field "Application.labels"`)}
	}
	if _, ok := ac.mutation.CreateTime(); !ok {
		return &ValidationError{Name: "createTime", err: errors.New(`model: missing required field "Application.createTime"`)}
	}
	if _, ok := ac.mutation.UpdateTime(); !ok {
		return &ValidationError{Name: "updateTime", err: errors.New(`model: missing required field "Application.updateTime"`)}
	}
	if _, ok := ac.mutation.ProjectID(); !ok {
		return &ValidationError{Name: "projectID", err: errors.New(`model: missing required field "Application.projectID"`)}
	}
	if v, ok := ac.mutation.ProjectID(); ok {
		if err := application.ProjectIDValidator(string(v)); err != nil {
			return &ValidationError{Name: "projectID", err: fmt.Errorf(`model: validator failed for field "Application.projectID": %w`, err)}
		}
	}
	if _, ok := ac.mutation.ProjectID(); !ok {
		return &ValidationError{Name: "project", err: errors.New(`model: missing required edge "Application.project"`)}
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
		_spec = sqlgraph.NewCreateSpec(application.Table, sqlgraph.NewFieldSpec(application.FieldID, field.TypeString))
	)
	_spec.Schema = ac.schemaConfig.Application
	_spec.OnConflict = ac.conflict
	if id, ok := ac.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := ac.mutation.Name(); ok {
		_spec.SetField(application.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := ac.mutation.Description(); ok {
		_spec.SetField(application.FieldDescription, field.TypeString, value)
		_node.Description = value
	}
	if value, ok := ac.mutation.Labels(); ok {
		_spec.SetField(application.FieldLabels, field.TypeJSON, value)
		_node.Labels = value
	}
	if value, ok := ac.mutation.CreateTime(); ok {
		_spec.SetField(application.FieldCreateTime, field.TypeTime, value)
		_node.CreateTime = &value
	}
	if value, ok := ac.mutation.UpdateTime(); ok {
		_spec.SetField(application.FieldUpdateTime, field.TypeTime, value)
		_node.UpdateTime = &value
	}
	if value, ok := ac.mutation.Variables(); ok {
		_spec.SetField(application.FieldVariables, field.TypeOther, value)
		_node.Variables = value
	}
	if nodes := ac.mutation.ProjectIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   application.ProjectTable,
			Columns: []string{application.ProjectColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: project.FieldID,
				},
			},
		}
		edge.Schema = ac.schemaConfig.Application
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.ProjectID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ac.mutation.InstancesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   application.InstancesTable,
			Columns: []string{application.InstancesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: applicationinstance.FieldID,
				},
			},
		}
		edge.Schema = ac.schemaConfig.ApplicationInstance
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
//	client.Application.Create().
//		SetName(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.ApplicationUpsert) {
//			SetName(v+v).
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

// SetName sets the "name" field.
func (u *ApplicationUpsert) SetName(v string) *ApplicationUpsert {
	u.Set(application.FieldName, v)
	return u
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *ApplicationUpsert) UpdateName() *ApplicationUpsert {
	u.SetExcluded(application.FieldName)
	return u
}

// SetDescription sets the "description" field.
func (u *ApplicationUpsert) SetDescription(v string) *ApplicationUpsert {
	u.Set(application.FieldDescription, v)
	return u
}

// UpdateDescription sets the "description" field to the value that was provided on create.
func (u *ApplicationUpsert) UpdateDescription() *ApplicationUpsert {
	u.SetExcluded(application.FieldDescription)
	return u
}

// ClearDescription clears the value of the "description" field.
func (u *ApplicationUpsert) ClearDescription() *ApplicationUpsert {
	u.SetNull(application.FieldDescription)
	return u
}

// SetLabels sets the "labels" field.
func (u *ApplicationUpsert) SetLabels(v map[string]string) *ApplicationUpsert {
	u.Set(application.FieldLabels, v)
	return u
}

// UpdateLabels sets the "labels" field to the value that was provided on create.
func (u *ApplicationUpsert) UpdateLabels() *ApplicationUpsert {
	u.SetExcluded(application.FieldLabels)
	return u
}

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

// SetVariables sets the "variables" field.
func (u *ApplicationUpsert) SetVariables(v property.Schemas) *ApplicationUpsert {
	u.Set(application.FieldVariables, v)
	return u
}

// UpdateVariables sets the "variables" field to the value that was provided on create.
func (u *ApplicationUpsert) UpdateVariables() *ApplicationUpsert {
	u.SetExcluded(application.FieldVariables)
	return u
}

// ClearVariables clears the value of the "variables" field.
func (u *ApplicationUpsert) ClearVariables() *ApplicationUpsert {
	u.SetNull(application.FieldVariables)
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

// SetName sets the "name" field.
func (u *ApplicationUpsertOne) SetName(v string) *ApplicationUpsertOne {
	return u.Update(func(s *ApplicationUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *ApplicationUpsertOne) UpdateName() *ApplicationUpsertOne {
	return u.Update(func(s *ApplicationUpsert) {
		s.UpdateName()
	})
}

// SetDescription sets the "description" field.
func (u *ApplicationUpsertOne) SetDescription(v string) *ApplicationUpsertOne {
	return u.Update(func(s *ApplicationUpsert) {
		s.SetDescription(v)
	})
}

// UpdateDescription sets the "description" field to the value that was provided on create.
func (u *ApplicationUpsertOne) UpdateDescription() *ApplicationUpsertOne {
	return u.Update(func(s *ApplicationUpsert) {
		s.UpdateDescription()
	})
}

// ClearDescription clears the value of the "description" field.
func (u *ApplicationUpsertOne) ClearDescription() *ApplicationUpsertOne {
	return u.Update(func(s *ApplicationUpsert) {
		s.ClearDescription()
	})
}

// SetLabels sets the "labels" field.
func (u *ApplicationUpsertOne) SetLabels(v map[string]string) *ApplicationUpsertOne {
	return u.Update(func(s *ApplicationUpsert) {
		s.SetLabels(v)
	})
}

// UpdateLabels sets the "labels" field to the value that was provided on create.
func (u *ApplicationUpsertOne) UpdateLabels() *ApplicationUpsertOne {
	return u.Update(func(s *ApplicationUpsert) {
		s.UpdateLabels()
	})
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

// SetVariables sets the "variables" field.
func (u *ApplicationUpsertOne) SetVariables(v property.Schemas) *ApplicationUpsertOne {
	return u.Update(func(s *ApplicationUpsert) {
		s.SetVariables(v)
	})
}

// UpdateVariables sets the "variables" field to the value that was provided on create.
func (u *ApplicationUpsertOne) UpdateVariables() *ApplicationUpsertOne {
	return u.Update(func(s *ApplicationUpsert) {
		s.UpdateVariables()
	})
}

// ClearVariables clears the value of the "variables" field.
func (u *ApplicationUpsertOne) ClearVariables() *ApplicationUpsertOne {
	return u.Update(func(s *ApplicationUpsert) {
		s.ClearVariables()
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
//			SetName(v+v).
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

// SetName sets the "name" field.
func (u *ApplicationUpsertBulk) SetName(v string) *ApplicationUpsertBulk {
	return u.Update(func(s *ApplicationUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *ApplicationUpsertBulk) UpdateName() *ApplicationUpsertBulk {
	return u.Update(func(s *ApplicationUpsert) {
		s.UpdateName()
	})
}

// SetDescription sets the "description" field.
func (u *ApplicationUpsertBulk) SetDescription(v string) *ApplicationUpsertBulk {
	return u.Update(func(s *ApplicationUpsert) {
		s.SetDescription(v)
	})
}

// UpdateDescription sets the "description" field to the value that was provided on create.
func (u *ApplicationUpsertBulk) UpdateDescription() *ApplicationUpsertBulk {
	return u.Update(func(s *ApplicationUpsert) {
		s.UpdateDescription()
	})
}

// ClearDescription clears the value of the "description" field.
func (u *ApplicationUpsertBulk) ClearDescription() *ApplicationUpsertBulk {
	return u.Update(func(s *ApplicationUpsert) {
		s.ClearDescription()
	})
}

// SetLabels sets the "labels" field.
func (u *ApplicationUpsertBulk) SetLabels(v map[string]string) *ApplicationUpsertBulk {
	return u.Update(func(s *ApplicationUpsert) {
		s.SetLabels(v)
	})
}

// UpdateLabels sets the "labels" field to the value that was provided on create.
func (u *ApplicationUpsertBulk) UpdateLabels() *ApplicationUpsertBulk {
	return u.Update(func(s *ApplicationUpsert) {
		s.UpdateLabels()
	})
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

// SetVariables sets the "variables" field.
func (u *ApplicationUpsertBulk) SetVariables(v property.Schemas) *ApplicationUpsertBulk {
	return u.Update(func(s *ApplicationUpsert) {
		s.SetVariables(v)
	})
}

// UpdateVariables sets the "variables" field to the value that was provided on create.
func (u *ApplicationUpsertBulk) UpdateVariables() *ApplicationUpsertBulk {
	return u.Update(func(s *ApplicationUpsert) {
		s.UpdateVariables()
	})
}

// ClearVariables clears the value of the "variables" field.
func (u *ApplicationUpsertBulk) ClearVariables() *ApplicationUpsertBulk {
	return u.Update(func(s *ApplicationUpsert) {
		s.ClearVariables()
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
