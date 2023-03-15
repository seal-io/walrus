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
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/model/applicationrevision"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/types/oid"
)

// ApplicationInstanceCreate is the builder for creating a ApplicationInstance entity.
type ApplicationInstanceCreate struct {
	config
	mutation *ApplicationInstanceMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetStatus sets the "status" field.
func (aic *ApplicationInstanceCreate) SetStatus(s string) *ApplicationInstanceCreate {
	aic.mutation.SetStatus(s)
	return aic
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (aic *ApplicationInstanceCreate) SetNillableStatus(s *string) *ApplicationInstanceCreate {
	if s != nil {
		aic.SetStatus(*s)
	}
	return aic
}

// SetStatusMessage sets the "statusMessage" field.
func (aic *ApplicationInstanceCreate) SetStatusMessage(s string) *ApplicationInstanceCreate {
	aic.mutation.SetStatusMessage(s)
	return aic
}

// SetNillableStatusMessage sets the "statusMessage" field if the given value is not nil.
func (aic *ApplicationInstanceCreate) SetNillableStatusMessage(s *string) *ApplicationInstanceCreate {
	if s != nil {
		aic.SetStatusMessage(*s)
	}
	return aic
}

// SetCreateTime sets the "createTime" field.
func (aic *ApplicationInstanceCreate) SetCreateTime(t time.Time) *ApplicationInstanceCreate {
	aic.mutation.SetCreateTime(t)
	return aic
}

// SetNillableCreateTime sets the "createTime" field if the given value is not nil.
func (aic *ApplicationInstanceCreate) SetNillableCreateTime(t *time.Time) *ApplicationInstanceCreate {
	if t != nil {
		aic.SetCreateTime(*t)
	}
	return aic
}

// SetUpdateTime sets the "updateTime" field.
func (aic *ApplicationInstanceCreate) SetUpdateTime(t time.Time) *ApplicationInstanceCreate {
	aic.mutation.SetUpdateTime(t)
	return aic
}

// SetNillableUpdateTime sets the "updateTime" field if the given value is not nil.
func (aic *ApplicationInstanceCreate) SetNillableUpdateTime(t *time.Time) *ApplicationInstanceCreate {
	if t != nil {
		aic.SetUpdateTime(*t)
	}
	return aic
}

// SetApplicationID sets the "applicationID" field.
func (aic *ApplicationInstanceCreate) SetApplicationID(o oid.ID) *ApplicationInstanceCreate {
	aic.mutation.SetApplicationID(o)
	return aic
}

// SetEnvironmentID sets the "environmentID" field.
func (aic *ApplicationInstanceCreate) SetEnvironmentID(o oid.ID) *ApplicationInstanceCreate {
	aic.mutation.SetEnvironmentID(o)
	return aic
}

// SetName sets the "name" field.
func (aic *ApplicationInstanceCreate) SetName(s string) *ApplicationInstanceCreate {
	aic.mutation.SetName(s)
	return aic
}

// SetVariables sets the "variables" field.
func (aic *ApplicationInstanceCreate) SetVariables(m map[string]interface{}) *ApplicationInstanceCreate {
	aic.mutation.SetVariables(m)
	return aic
}

// SetID sets the "id" field.
func (aic *ApplicationInstanceCreate) SetID(o oid.ID) *ApplicationInstanceCreate {
	aic.mutation.SetID(o)
	return aic
}

// SetApplication sets the "application" edge to the Application entity.
func (aic *ApplicationInstanceCreate) SetApplication(a *Application) *ApplicationInstanceCreate {
	return aic.SetApplicationID(a.ID)
}

// SetEnvironment sets the "environment" edge to the Environment entity.
func (aic *ApplicationInstanceCreate) SetEnvironment(e *Environment) *ApplicationInstanceCreate {
	return aic.SetEnvironmentID(e.ID)
}

// AddRevisionIDs adds the "revisions" edge to the ApplicationRevision entity by IDs.
func (aic *ApplicationInstanceCreate) AddRevisionIDs(ids ...oid.ID) *ApplicationInstanceCreate {
	aic.mutation.AddRevisionIDs(ids...)
	return aic
}

// AddRevisions adds the "revisions" edges to the ApplicationRevision entity.
func (aic *ApplicationInstanceCreate) AddRevisions(a ...*ApplicationRevision) *ApplicationInstanceCreate {
	ids := make([]oid.ID, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return aic.AddRevisionIDs(ids...)
}

// AddResourceIDs adds the "resources" edge to the ApplicationResource entity by IDs.
func (aic *ApplicationInstanceCreate) AddResourceIDs(ids ...oid.ID) *ApplicationInstanceCreate {
	aic.mutation.AddResourceIDs(ids...)
	return aic
}

// AddResources adds the "resources" edges to the ApplicationResource entity.
func (aic *ApplicationInstanceCreate) AddResources(a ...*ApplicationResource) *ApplicationInstanceCreate {
	ids := make([]oid.ID, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return aic.AddResourceIDs(ids...)
}

// Mutation returns the ApplicationInstanceMutation object of the builder.
func (aic *ApplicationInstanceCreate) Mutation() *ApplicationInstanceMutation {
	return aic.mutation
}

// Save creates the ApplicationInstance in the database.
func (aic *ApplicationInstanceCreate) Save(ctx context.Context) (*ApplicationInstance, error) {
	if err := aic.defaults(); err != nil {
		return nil, err
	}
	return withHooks[*ApplicationInstance, ApplicationInstanceMutation](ctx, aic.sqlSave, aic.mutation, aic.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (aic *ApplicationInstanceCreate) SaveX(ctx context.Context) *ApplicationInstance {
	v, err := aic.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (aic *ApplicationInstanceCreate) Exec(ctx context.Context) error {
	_, err := aic.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (aic *ApplicationInstanceCreate) ExecX(ctx context.Context) {
	if err := aic.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (aic *ApplicationInstanceCreate) defaults() error {
	if _, ok := aic.mutation.CreateTime(); !ok {
		if applicationinstance.DefaultCreateTime == nil {
			return fmt.Errorf("model: uninitialized applicationinstance.DefaultCreateTime (forgotten import model/runtime?)")
		}
		v := applicationinstance.DefaultCreateTime()
		aic.mutation.SetCreateTime(v)
	}
	if _, ok := aic.mutation.UpdateTime(); !ok {
		if applicationinstance.DefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized applicationinstance.DefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := applicationinstance.DefaultUpdateTime()
		aic.mutation.SetUpdateTime(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (aic *ApplicationInstanceCreate) check() error {
	if _, ok := aic.mutation.CreateTime(); !ok {
		return &ValidationError{Name: "createTime", err: errors.New(`model: missing required field "ApplicationInstance.createTime"`)}
	}
	if _, ok := aic.mutation.UpdateTime(); !ok {
		return &ValidationError{Name: "updateTime", err: errors.New(`model: missing required field "ApplicationInstance.updateTime"`)}
	}
	if _, ok := aic.mutation.ApplicationID(); !ok {
		return &ValidationError{Name: "applicationID", err: errors.New(`model: missing required field "ApplicationInstance.applicationID"`)}
	}
	if v, ok := aic.mutation.ApplicationID(); ok {
		if err := applicationinstance.ApplicationIDValidator(string(v)); err != nil {
			return &ValidationError{Name: "applicationID", err: fmt.Errorf(`model: validator failed for field "ApplicationInstance.applicationID": %w`, err)}
		}
	}
	if _, ok := aic.mutation.EnvironmentID(); !ok {
		return &ValidationError{Name: "environmentID", err: errors.New(`model: missing required field "ApplicationInstance.environmentID"`)}
	}
	if v, ok := aic.mutation.EnvironmentID(); ok {
		if err := applicationinstance.EnvironmentIDValidator(string(v)); err != nil {
			return &ValidationError{Name: "environmentID", err: fmt.Errorf(`model: validator failed for field "ApplicationInstance.environmentID": %w`, err)}
		}
	}
	if _, ok := aic.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`model: missing required field "ApplicationInstance.name"`)}
	}
	if v, ok := aic.mutation.Name(); ok {
		if err := applicationinstance.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`model: validator failed for field "ApplicationInstance.name": %w`, err)}
		}
	}
	if _, ok := aic.mutation.ApplicationID(); !ok {
		return &ValidationError{Name: "application", err: errors.New(`model: missing required edge "ApplicationInstance.application"`)}
	}
	if _, ok := aic.mutation.EnvironmentID(); !ok {
		return &ValidationError{Name: "environment", err: errors.New(`model: missing required edge "ApplicationInstance.environment"`)}
	}
	return nil
}

func (aic *ApplicationInstanceCreate) sqlSave(ctx context.Context) (*ApplicationInstance, error) {
	if err := aic.check(); err != nil {
		return nil, err
	}
	_node, _spec := aic.createSpec()
	if err := sqlgraph.CreateNode(ctx, aic.driver, _spec); err != nil {
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
	aic.mutation.id = &_node.ID
	aic.mutation.done = true
	return _node, nil
}

func (aic *ApplicationInstanceCreate) createSpec() (*ApplicationInstance, *sqlgraph.CreateSpec) {
	var (
		_node = &ApplicationInstance{config: aic.config}
		_spec = sqlgraph.NewCreateSpec(applicationinstance.Table, sqlgraph.NewFieldSpec(applicationinstance.FieldID, field.TypeString))
	)
	_spec.Schema = aic.schemaConfig.ApplicationInstance
	_spec.OnConflict = aic.conflict
	if id, ok := aic.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := aic.mutation.Status(); ok {
		_spec.SetField(applicationinstance.FieldStatus, field.TypeString, value)
		_node.Status = value
	}
	if value, ok := aic.mutation.StatusMessage(); ok {
		_spec.SetField(applicationinstance.FieldStatusMessage, field.TypeString, value)
		_node.StatusMessage = value
	}
	if value, ok := aic.mutation.CreateTime(); ok {
		_spec.SetField(applicationinstance.FieldCreateTime, field.TypeTime, value)
		_node.CreateTime = &value
	}
	if value, ok := aic.mutation.UpdateTime(); ok {
		_spec.SetField(applicationinstance.FieldUpdateTime, field.TypeTime, value)
		_node.UpdateTime = &value
	}
	if value, ok := aic.mutation.Name(); ok {
		_spec.SetField(applicationinstance.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := aic.mutation.Variables(); ok {
		_spec.SetField(applicationinstance.FieldVariables, field.TypeJSON, value)
		_node.Variables = value
	}
	if nodes := aic.mutation.ApplicationIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   applicationinstance.ApplicationTable,
			Columns: []string{applicationinstance.ApplicationColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: application.FieldID,
				},
			},
		}
		edge.Schema = aic.schemaConfig.ApplicationInstance
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.ApplicationID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := aic.mutation.EnvironmentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   applicationinstance.EnvironmentTable,
			Columns: []string{applicationinstance.EnvironmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: environment.FieldID,
				},
			},
		}
		edge.Schema = aic.schemaConfig.ApplicationInstance
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.EnvironmentID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := aic.mutation.RevisionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   applicationinstance.RevisionsTable,
			Columns: []string{applicationinstance.RevisionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: applicationrevision.FieldID,
				},
			},
		}
		edge.Schema = aic.schemaConfig.ApplicationRevision
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := aic.mutation.ResourcesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   applicationinstance.ResourcesTable,
			Columns: []string{applicationinstance.ResourcesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: applicationresource.FieldID,
				},
			},
		}
		edge.Schema = aic.schemaConfig.ApplicationResource
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
//	client.ApplicationInstance.Create().
//		SetStatus(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.ApplicationInstanceUpsert) {
//			SetStatus(v+v).
//		}).
//		Exec(ctx)
func (aic *ApplicationInstanceCreate) OnConflict(opts ...sql.ConflictOption) *ApplicationInstanceUpsertOne {
	aic.conflict = opts
	return &ApplicationInstanceUpsertOne{
		create: aic,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.ApplicationInstance.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (aic *ApplicationInstanceCreate) OnConflictColumns(columns ...string) *ApplicationInstanceUpsertOne {
	aic.conflict = append(aic.conflict, sql.ConflictColumns(columns...))
	return &ApplicationInstanceUpsertOne{
		create: aic,
	}
}

type (
	// ApplicationInstanceUpsertOne is the builder for "upsert"-ing
	//  one ApplicationInstance node.
	ApplicationInstanceUpsertOne struct {
		create *ApplicationInstanceCreate
	}

	// ApplicationInstanceUpsert is the "OnConflict" setter.
	ApplicationInstanceUpsert struct {
		*sql.UpdateSet
	}
)

// SetStatus sets the "status" field.
func (u *ApplicationInstanceUpsert) SetStatus(v string) *ApplicationInstanceUpsert {
	u.Set(applicationinstance.FieldStatus, v)
	return u
}

// UpdateStatus sets the "status" field to the value that was provided on create.
func (u *ApplicationInstanceUpsert) UpdateStatus() *ApplicationInstanceUpsert {
	u.SetExcluded(applicationinstance.FieldStatus)
	return u
}

// ClearStatus clears the value of the "status" field.
func (u *ApplicationInstanceUpsert) ClearStatus() *ApplicationInstanceUpsert {
	u.SetNull(applicationinstance.FieldStatus)
	return u
}

// SetStatusMessage sets the "statusMessage" field.
func (u *ApplicationInstanceUpsert) SetStatusMessage(v string) *ApplicationInstanceUpsert {
	u.Set(applicationinstance.FieldStatusMessage, v)
	return u
}

// UpdateStatusMessage sets the "statusMessage" field to the value that was provided on create.
func (u *ApplicationInstanceUpsert) UpdateStatusMessage() *ApplicationInstanceUpsert {
	u.SetExcluded(applicationinstance.FieldStatusMessage)
	return u
}

// ClearStatusMessage clears the value of the "statusMessage" field.
func (u *ApplicationInstanceUpsert) ClearStatusMessage() *ApplicationInstanceUpsert {
	u.SetNull(applicationinstance.FieldStatusMessage)
	return u
}

// SetUpdateTime sets the "updateTime" field.
func (u *ApplicationInstanceUpsert) SetUpdateTime(v time.Time) *ApplicationInstanceUpsert {
	u.Set(applicationinstance.FieldUpdateTime, v)
	return u
}

// UpdateUpdateTime sets the "updateTime" field to the value that was provided on create.
func (u *ApplicationInstanceUpsert) UpdateUpdateTime() *ApplicationInstanceUpsert {
	u.SetExcluded(applicationinstance.FieldUpdateTime)
	return u
}

// SetVariables sets the "variables" field.
func (u *ApplicationInstanceUpsert) SetVariables(v map[string]interface{}) *ApplicationInstanceUpsert {
	u.Set(applicationinstance.FieldVariables, v)
	return u
}

// UpdateVariables sets the "variables" field to the value that was provided on create.
func (u *ApplicationInstanceUpsert) UpdateVariables() *ApplicationInstanceUpsert {
	u.SetExcluded(applicationinstance.FieldVariables)
	return u
}

// ClearVariables clears the value of the "variables" field.
func (u *ApplicationInstanceUpsert) ClearVariables() *ApplicationInstanceUpsert {
	u.SetNull(applicationinstance.FieldVariables)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.ApplicationInstance.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(applicationinstance.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *ApplicationInstanceUpsertOne) UpdateNewValues() *ApplicationInstanceUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(applicationinstance.FieldID)
		}
		if _, exists := u.create.mutation.CreateTime(); exists {
			s.SetIgnore(applicationinstance.FieldCreateTime)
		}
		if _, exists := u.create.mutation.ApplicationID(); exists {
			s.SetIgnore(applicationinstance.FieldApplicationID)
		}
		if _, exists := u.create.mutation.EnvironmentID(); exists {
			s.SetIgnore(applicationinstance.FieldEnvironmentID)
		}
		if _, exists := u.create.mutation.Name(); exists {
			s.SetIgnore(applicationinstance.FieldName)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.ApplicationInstance.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *ApplicationInstanceUpsertOne) Ignore() *ApplicationInstanceUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *ApplicationInstanceUpsertOne) DoNothing() *ApplicationInstanceUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the ApplicationInstanceCreate.OnConflict
// documentation for more info.
func (u *ApplicationInstanceUpsertOne) Update(set func(*ApplicationInstanceUpsert)) *ApplicationInstanceUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&ApplicationInstanceUpsert{UpdateSet: update})
	}))
	return u
}

// SetStatus sets the "status" field.
func (u *ApplicationInstanceUpsertOne) SetStatus(v string) *ApplicationInstanceUpsertOne {
	return u.Update(func(s *ApplicationInstanceUpsert) {
		s.SetStatus(v)
	})
}

// UpdateStatus sets the "status" field to the value that was provided on create.
func (u *ApplicationInstanceUpsertOne) UpdateStatus() *ApplicationInstanceUpsertOne {
	return u.Update(func(s *ApplicationInstanceUpsert) {
		s.UpdateStatus()
	})
}

// ClearStatus clears the value of the "status" field.
func (u *ApplicationInstanceUpsertOne) ClearStatus() *ApplicationInstanceUpsertOne {
	return u.Update(func(s *ApplicationInstanceUpsert) {
		s.ClearStatus()
	})
}

// SetStatusMessage sets the "statusMessage" field.
func (u *ApplicationInstanceUpsertOne) SetStatusMessage(v string) *ApplicationInstanceUpsertOne {
	return u.Update(func(s *ApplicationInstanceUpsert) {
		s.SetStatusMessage(v)
	})
}

// UpdateStatusMessage sets the "statusMessage" field to the value that was provided on create.
func (u *ApplicationInstanceUpsertOne) UpdateStatusMessage() *ApplicationInstanceUpsertOne {
	return u.Update(func(s *ApplicationInstanceUpsert) {
		s.UpdateStatusMessage()
	})
}

// ClearStatusMessage clears the value of the "statusMessage" field.
func (u *ApplicationInstanceUpsertOne) ClearStatusMessage() *ApplicationInstanceUpsertOne {
	return u.Update(func(s *ApplicationInstanceUpsert) {
		s.ClearStatusMessage()
	})
}

// SetUpdateTime sets the "updateTime" field.
func (u *ApplicationInstanceUpsertOne) SetUpdateTime(v time.Time) *ApplicationInstanceUpsertOne {
	return u.Update(func(s *ApplicationInstanceUpsert) {
		s.SetUpdateTime(v)
	})
}

// UpdateUpdateTime sets the "updateTime" field to the value that was provided on create.
func (u *ApplicationInstanceUpsertOne) UpdateUpdateTime() *ApplicationInstanceUpsertOne {
	return u.Update(func(s *ApplicationInstanceUpsert) {
		s.UpdateUpdateTime()
	})
}

// SetVariables sets the "variables" field.
func (u *ApplicationInstanceUpsertOne) SetVariables(v map[string]interface{}) *ApplicationInstanceUpsertOne {
	return u.Update(func(s *ApplicationInstanceUpsert) {
		s.SetVariables(v)
	})
}

// UpdateVariables sets the "variables" field to the value that was provided on create.
func (u *ApplicationInstanceUpsertOne) UpdateVariables() *ApplicationInstanceUpsertOne {
	return u.Update(func(s *ApplicationInstanceUpsert) {
		s.UpdateVariables()
	})
}

// ClearVariables clears the value of the "variables" field.
func (u *ApplicationInstanceUpsertOne) ClearVariables() *ApplicationInstanceUpsertOne {
	return u.Update(func(s *ApplicationInstanceUpsert) {
		s.ClearVariables()
	})
}

// Exec executes the query.
func (u *ApplicationInstanceUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("model: missing options for ApplicationInstanceCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *ApplicationInstanceUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *ApplicationInstanceUpsertOne) ID(ctx context.Context) (id oid.ID, err error) {
	if u.create.driver.Dialect() == dialect.MySQL {
		// In case of "ON CONFLICT", there is no way to get back non-numeric ID
		// fields from the database since MySQL does not support the RETURNING clause.
		return id, errors.New("model: ApplicationInstanceUpsertOne.ID is not supported by MySQL driver. Use ApplicationInstanceUpsertOne.Exec instead")
	}
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *ApplicationInstanceUpsertOne) IDX(ctx context.Context) oid.ID {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// ApplicationInstanceCreateBulk is the builder for creating many ApplicationInstance entities in bulk.
type ApplicationInstanceCreateBulk struct {
	config
	builders []*ApplicationInstanceCreate
	conflict []sql.ConflictOption
}

// Save creates the ApplicationInstance entities in the database.
func (aicb *ApplicationInstanceCreateBulk) Save(ctx context.Context) ([]*ApplicationInstance, error) {
	specs := make([]*sqlgraph.CreateSpec, len(aicb.builders))
	nodes := make([]*ApplicationInstance, len(aicb.builders))
	mutators := make([]Mutator, len(aicb.builders))
	for i := range aicb.builders {
		func(i int, root context.Context) {
			builder := aicb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ApplicationInstanceMutation)
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
					_, err = mutators[i+1].Mutate(root, aicb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = aicb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, aicb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, aicb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (aicb *ApplicationInstanceCreateBulk) SaveX(ctx context.Context) []*ApplicationInstance {
	v, err := aicb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (aicb *ApplicationInstanceCreateBulk) Exec(ctx context.Context) error {
	_, err := aicb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (aicb *ApplicationInstanceCreateBulk) ExecX(ctx context.Context) {
	if err := aicb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.ApplicationInstance.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.ApplicationInstanceUpsert) {
//			SetStatus(v+v).
//		}).
//		Exec(ctx)
func (aicb *ApplicationInstanceCreateBulk) OnConflict(opts ...sql.ConflictOption) *ApplicationInstanceUpsertBulk {
	aicb.conflict = opts
	return &ApplicationInstanceUpsertBulk{
		create: aicb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.ApplicationInstance.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (aicb *ApplicationInstanceCreateBulk) OnConflictColumns(columns ...string) *ApplicationInstanceUpsertBulk {
	aicb.conflict = append(aicb.conflict, sql.ConflictColumns(columns...))
	return &ApplicationInstanceUpsertBulk{
		create: aicb,
	}
}

// ApplicationInstanceUpsertBulk is the builder for "upsert"-ing
// a bulk of ApplicationInstance nodes.
type ApplicationInstanceUpsertBulk struct {
	create *ApplicationInstanceCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.ApplicationInstance.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(applicationinstance.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *ApplicationInstanceUpsertBulk) UpdateNewValues() *ApplicationInstanceUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(applicationinstance.FieldID)
			}
			if _, exists := b.mutation.CreateTime(); exists {
				s.SetIgnore(applicationinstance.FieldCreateTime)
			}
			if _, exists := b.mutation.ApplicationID(); exists {
				s.SetIgnore(applicationinstance.FieldApplicationID)
			}
			if _, exists := b.mutation.EnvironmentID(); exists {
				s.SetIgnore(applicationinstance.FieldEnvironmentID)
			}
			if _, exists := b.mutation.Name(); exists {
				s.SetIgnore(applicationinstance.FieldName)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.ApplicationInstance.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *ApplicationInstanceUpsertBulk) Ignore() *ApplicationInstanceUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *ApplicationInstanceUpsertBulk) DoNothing() *ApplicationInstanceUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the ApplicationInstanceCreateBulk.OnConflict
// documentation for more info.
func (u *ApplicationInstanceUpsertBulk) Update(set func(*ApplicationInstanceUpsert)) *ApplicationInstanceUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&ApplicationInstanceUpsert{UpdateSet: update})
	}))
	return u
}

// SetStatus sets the "status" field.
func (u *ApplicationInstanceUpsertBulk) SetStatus(v string) *ApplicationInstanceUpsertBulk {
	return u.Update(func(s *ApplicationInstanceUpsert) {
		s.SetStatus(v)
	})
}

// UpdateStatus sets the "status" field to the value that was provided on create.
func (u *ApplicationInstanceUpsertBulk) UpdateStatus() *ApplicationInstanceUpsertBulk {
	return u.Update(func(s *ApplicationInstanceUpsert) {
		s.UpdateStatus()
	})
}

// ClearStatus clears the value of the "status" field.
func (u *ApplicationInstanceUpsertBulk) ClearStatus() *ApplicationInstanceUpsertBulk {
	return u.Update(func(s *ApplicationInstanceUpsert) {
		s.ClearStatus()
	})
}

// SetStatusMessage sets the "statusMessage" field.
func (u *ApplicationInstanceUpsertBulk) SetStatusMessage(v string) *ApplicationInstanceUpsertBulk {
	return u.Update(func(s *ApplicationInstanceUpsert) {
		s.SetStatusMessage(v)
	})
}

// UpdateStatusMessage sets the "statusMessage" field to the value that was provided on create.
func (u *ApplicationInstanceUpsertBulk) UpdateStatusMessage() *ApplicationInstanceUpsertBulk {
	return u.Update(func(s *ApplicationInstanceUpsert) {
		s.UpdateStatusMessage()
	})
}

// ClearStatusMessage clears the value of the "statusMessage" field.
func (u *ApplicationInstanceUpsertBulk) ClearStatusMessage() *ApplicationInstanceUpsertBulk {
	return u.Update(func(s *ApplicationInstanceUpsert) {
		s.ClearStatusMessage()
	})
}

// SetUpdateTime sets the "updateTime" field.
func (u *ApplicationInstanceUpsertBulk) SetUpdateTime(v time.Time) *ApplicationInstanceUpsertBulk {
	return u.Update(func(s *ApplicationInstanceUpsert) {
		s.SetUpdateTime(v)
	})
}

// UpdateUpdateTime sets the "updateTime" field to the value that was provided on create.
func (u *ApplicationInstanceUpsertBulk) UpdateUpdateTime() *ApplicationInstanceUpsertBulk {
	return u.Update(func(s *ApplicationInstanceUpsert) {
		s.UpdateUpdateTime()
	})
}

// SetVariables sets the "variables" field.
func (u *ApplicationInstanceUpsertBulk) SetVariables(v map[string]interface{}) *ApplicationInstanceUpsertBulk {
	return u.Update(func(s *ApplicationInstanceUpsert) {
		s.SetVariables(v)
	})
}

// UpdateVariables sets the "variables" field to the value that was provided on create.
func (u *ApplicationInstanceUpsertBulk) UpdateVariables() *ApplicationInstanceUpsertBulk {
	return u.Update(func(s *ApplicationInstanceUpsert) {
		s.UpdateVariables()
	})
}

// ClearVariables clears the value of the "variables" field.
func (u *ApplicationInstanceUpsertBulk) ClearVariables() *ApplicationInstanceUpsertBulk {
	return u.Update(func(s *ApplicationInstanceUpsert) {
		s.ClearVariables()
	})
}

// Exec executes the query.
func (u *ApplicationInstanceUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("model: OnConflict was set for builder %d. Set it on the ApplicationInstanceCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("model: missing options for ApplicationInstanceCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *ApplicationInstanceUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
