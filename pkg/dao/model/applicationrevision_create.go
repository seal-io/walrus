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
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/crypto"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/pkg/dao/types/property"
)

// ApplicationRevisionCreate is the builder for creating a ApplicationRevision entity.
type ApplicationRevisionCreate struct {
	config
	mutation *ApplicationRevisionMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetStatus sets the "status" field.
func (arc *ApplicationRevisionCreate) SetStatus(s string) *ApplicationRevisionCreate {
	arc.mutation.SetStatus(s)
	return arc
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (arc *ApplicationRevisionCreate) SetNillableStatus(s *string) *ApplicationRevisionCreate {
	if s != nil {
		arc.SetStatus(*s)
	}
	return arc
}

// SetStatusMessage sets the "statusMessage" field.
func (arc *ApplicationRevisionCreate) SetStatusMessage(s string) *ApplicationRevisionCreate {
	arc.mutation.SetStatusMessage(s)
	return arc
}

// SetNillableStatusMessage sets the "statusMessage" field if the given value is not nil.
func (arc *ApplicationRevisionCreate) SetNillableStatusMessage(s *string) *ApplicationRevisionCreate {
	if s != nil {
		arc.SetStatusMessage(*s)
	}
	return arc
}

// SetCreateTime sets the "createTime" field.
func (arc *ApplicationRevisionCreate) SetCreateTime(t time.Time) *ApplicationRevisionCreate {
	arc.mutation.SetCreateTime(t)
	return arc
}

// SetNillableCreateTime sets the "createTime" field if the given value is not nil.
func (arc *ApplicationRevisionCreate) SetNillableCreateTime(t *time.Time) *ApplicationRevisionCreate {
	if t != nil {
		arc.SetCreateTime(*t)
	}
	return arc
}

// SetInstanceID sets the "instanceID" field.
func (arc *ApplicationRevisionCreate) SetInstanceID(o oid.ID) *ApplicationRevisionCreate {
	arc.mutation.SetInstanceID(o)
	return arc
}

// SetEnvironmentID sets the "environmentID" field.
func (arc *ApplicationRevisionCreate) SetEnvironmentID(o oid.ID) *ApplicationRevisionCreate {
	arc.mutation.SetEnvironmentID(o)
	return arc
}

// SetModules sets the "modules" field.
func (arc *ApplicationRevisionCreate) SetModules(tm []types.ApplicationModule) *ApplicationRevisionCreate {
	arc.mutation.SetModules(tm)
	return arc
}

// SetSecrets sets the "secrets" field.
func (arc *ApplicationRevisionCreate) SetSecrets(c crypto.Map[string, string]) *ApplicationRevisionCreate {
	arc.mutation.SetSecrets(c)
	return arc
}

// SetVariables sets the "variables" field.
func (arc *ApplicationRevisionCreate) SetVariables(pr property.Schemas) *ApplicationRevisionCreate {
	arc.mutation.SetVariables(pr)
	return arc
}

// SetInputVariables sets the "inputVariables" field.
func (arc *ApplicationRevisionCreate) SetInputVariables(pr property.Values) *ApplicationRevisionCreate {
	arc.mutation.SetInputVariables(pr)
	return arc
}

// SetInputPlan sets the "inputPlan" field.
func (arc *ApplicationRevisionCreate) SetInputPlan(s string) *ApplicationRevisionCreate {
	arc.mutation.SetInputPlan(s)
	return arc
}

// SetOutput sets the "output" field.
func (arc *ApplicationRevisionCreate) SetOutput(s string) *ApplicationRevisionCreate {
	arc.mutation.SetOutput(s)
	return arc
}

// SetDeployerType sets the "deployerType" field.
func (arc *ApplicationRevisionCreate) SetDeployerType(s string) *ApplicationRevisionCreate {
	arc.mutation.SetDeployerType(s)
	return arc
}

// SetNillableDeployerType sets the "deployerType" field if the given value is not nil.
func (arc *ApplicationRevisionCreate) SetNillableDeployerType(s *string) *ApplicationRevisionCreate {
	if s != nil {
		arc.SetDeployerType(*s)
	}
	return arc
}

// SetDuration sets the "duration" field.
func (arc *ApplicationRevisionCreate) SetDuration(i int) *ApplicationRevisionCreate {
	arc.mutation.SetDuration(i)
	return arc
}

// SetNillableDuration sets the "duration" field if the given value is not nil.
func (arc *ApplicationRevisionCreate) SetNillableDuration(i *int) *ApplicationRevisionCreate {
	if i != nil {
		arc.SetDuration(*i)
	}
	return arc
}

// SetPreviousRequiredProviders sets the "previousRequiredProviders" field.
func (arc *ApplicationRevisionCreate) SetPreviousRequiredProviders(tr []types.ProviderRequirement) *ApplicationRevisionCreate {
	arc.mutation.SetPreviousRequiredProviders(tr)
	return arc
}

// SetID sets the "id" field.
func (arc *ApplicationRevisionCreate) SetID(o oid.ID) *ApplicationRevisionCreate {
	arc.mutation.SetID(o)
	return arc
}

// SetInstance sets the "instance" edge to the ApplicationInstance entity.
func (arc *ApplicationRevisionCreate) SetInstance(a *ApplicationInstance) *ApplicationRevisionCreate {
	return arc.SetInstanceID(a.ID)
}

// SetEnvironment sets the "environment" edge to the Environment entity.
func (arc *ApplicationRevisionCreate) SetEnvironment(e *Environment) *ApplicationRevisionCreate {
	return arc.SetEnvironmentID(e.ID)
}

// Mutation returns the ApplicationRevisionMutation object of the builder.
func (arc *ApplicationRevisionCreate) Mutation() *ApplicationRevisionMutation {
	return arc.mutation
}

// Save creates the ApplicationRevision in the database.
func (arc *ApplicationRevisionCreate) Save(ctx context.Context) (*ApplicationRevision, error) {
	if err := arc.defaults(); err != nil {
		return nil, err
	}
	return withHooks[*ApplicationRevision, ApplicationRevisionMutation](ctx, arc.sqlSave, arc.mutation, arc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (arc *ApplicationRevisionCreate) SaveX(ctx context.Context) *ApplicationRevision {
	v, err := arc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (arc *ApplicationRevisionCreate) Exec(ctx context.Context) error {
	_, err := arc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (arc *ApplicationRevisionCreate) ExecX(ctx context.Context) {
	if err := arc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (arc *ApplicationRevisionCreate) defaults() error {
	if _, ok := arc.mutation.CreateTime(); !ok {
		if applicationrevision.DefaultCreateTime == nil {
			return fmt.Errorf("model: uninitialized applicationrevision.DefaultCreateTime (forgotten import model/runtime?)")
		}
		v := applicationrevision.DefaultCreateTime()
		arc.mutation.SetCreateTime(v)
	}
	if _, ok := arc.mutation.Modules(); !ok {
		v := applicationrevision.DefaultModules
		arc.mutation.SetModules(v)
	}
	if _, ok := arc.mutation.Secrets(); !ok {
		v := applicationrevision.DefaultSecrets
		arc.mutation.SetSecrets(v)
	}
	if _, ok := arc.mutation.InputVariables(); !ok {
		v := applicationrevision.DefaultInputVariables
		arc.mutation.SetInputVariables(v)
	}
	if _, ok := arc.mutation.DeployerType(); !ok {
		v := applicationrevision.DefaultDeployerType
		arc.mutation.SetDeployerType(v)
	}
	if _, ok := arc.mutation.Duration(); !ok {
		v := applicationrevision.DefaultDuration
		arc.mutation.SetDuration(v)
	}
	if _, ok := arc.mutation.PreviousRequiredProviders(); !ok {
		v := applicationrevision.DefaultPreviousRequiredProviders
		arc.mutation.SetPreviousRequiredProviders(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (arc *ApplicationRevisionCreate) check() error {
	if _, ok := arc.mutation.CreateTime(); !ok {
		return &ValidationError{Name: "createTime", err: errors.New(`model: missing required field "ApplicationRevision.createTime"`)}
	}
	if _, ok := arc.mutation.InstanceID(); !ok {
		return &ValidationError{Name: "instanceID", err: errors.New(`model: missing required field "ApplicationRevision.instanceID"`)}
	}
	if v, ok := arc.mutation.InstanceID(); ok {
		if err := applicationrevision.InstanceIDValidator(string(v)); err != nil {
			return &ValidationError{Name: "instanceID", err: fmt.Errorf(`model: validator failed for field "ApplicationRevision.instanceID": %w`, err)}
		}
	}
	if _, ok := arc.mutation.EnvironmentID(); !ok {
		return &ValidationError{Name: "environmentID", err: errors.New(`model: missing required field "ApplicationRevision.environmentID"`)}
	}
	if v, ok := arc.mutation.EnvironmentID(); ok {
		if err := applicationrevision.EnvironmentIDValidator(string(v)); err != nil {
			return &ValidationError{Name: "environmentID", err: fmt.Errorf(`model: validator failed for field "ApplicationRevision.environmentID": %w`, err)}
		}
	}
	if _, ok := arc.mutation.Modules(); !ok {
		return &ValidationError{Name: "modules", err: errors.New(`model: missing required field "ApplicationRevision.modules"`)}
	}
	if _, ok := arc.mutation.Secrets(); !ok {
		return &ValidationError{Name: "secrets", err: errors.New(`model: missing required field "ApplicationRevision.secrets"`)}
	}
	if _, ok := arc.mutation.InputVariables(); !ok {
		return &ValidationError{Name: "inputVariables", err: errors.New(`model: missing required field "ApplicationRevision.inputVariables"`)}
	}
	if _, ok := arc.mutation.InputPlan(); !ok {
		return &ValidationError{Name: "inputPlan", err: errors.New(`model: missing required field "ApplicationRevision.inputPlan"`)}
	}
	if _, ok := arc.mutation.Output(); !ok {
		return &ValidationError{Name: "output", err: errors.New(`model: missing required field "ApplicationRevision.output"`)}
	}
	if _, ok := arc.mutation.DeployerType(); !ok {
		return &ValidationError{Name: "deployerType", err: errors.New(`model: missing required field "ApplicationRevision.deployerType"`)}
	}
	if _, ok := arc.mutation.Duration(); !ok {
		return &ValidationError{Name: "duration", err: errors.New(`model: missing required field "ApplicationRevision.duration"`)}
	}
	if _, ok := arc.mutation.PreviousRequiredProviders(); !ok {
		return &ValidationError{Name: "previousRequiredProviders", err: errors.New(`model: missing required field "ApplicationRevision.previousRequiredProviders"`)}
	}
	if _, ok := arc.mutation.InstanceID(); !ok {
		return &ValidationError{Name: "instance", err: errors.New(`model: missing required edge "ApplicationRevision.instance"`)}
	}
	if _, ok := arc.mutation.EnvironmentID(); !ok {
		return &ValidationError{Name: "environment", err: errors.New(`model: missing required edge "ApplicationRevision.environment"`)}
	}
	return nil
}

func (arc *ApplicationRevisionCreate) sqlSave(ctx context.Context) (*ApplicationRevision, error) {
	if err := arc.check(); err != nil {
		return nil, err
	}
	_node, _spec := arc.createSpec()
	if err := sqlgraph.CreateNode(ctx, arc.driver, _spec); err != nil {
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
	arc.mutation.id = &_node.ID
	arc.mutation.done = true
	return _node, nil
}

func (arc *ApplicationRevisionCreate) createSpec() (*ApplicationRevision, *sqlgraph.CreateSpec) {
	var (
		_node = &ApplicationRevision{config: arc.config}
		_spec = sqlgraph.NewCreateSpec(applicationrevision.Table, sqlgraph.NewFieldSpec(applicationrevision.FieldID, field.TypeString))
	)
	_spec.Schema = arc.schemaConfig.ApplicationRevision
	_spec.OnConflict = arc.conflict
	if id, ok := arc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := arc.mutation.Status(); ok {
		_spec.SetField(applicationrevision.FieldStatus, field.TypeString, value)
		_node.Status = value
	}
	if value, ok := arc.mutation.StatusMessage(); ok {
		_spec.SetField(applicationrevision.FieldStatusMessage, field.TypeString, value)
		_node.StatusMessage = value
	}
	if value, ok := arc.mutation.CreateTime(); ok {
		_spec.SetField(applicationrevision.FieldCreateTime, field.TypeTime, value)
		_node.CreateTime = &value
	}
	if value, ok := arc.mutation.Modules(); ok {
		_spec.SetField(applicationrevision.FieldModules, field.TypeJSON, value)
		_node.Modules = value
	}
	if value, ok := arc.mutation.Secrets(); ok {
		_spec.SetField(applicationrevision.FieldSecrets, field.TypeOther, value)
		_node.Secrets = value
	}
	if value, ok := arc.mutation.Variables(); ok {
		_spec.SetField(applicationrevision.FieldVariables, field.TypeOther, value)
		_node.Variables = value
	}
	if value, ok := arc.mutation.InputVariables(); ok {
		_spec.SetField(applicationrevision.FieldInputVariables, field.TypeOther, value)
		_node.InputVariables = value
	}
	if value, ok := arc.mutation.InputPlan(); ok {
		_spec.SetField(applicationrevision.FieldInputPlan, field.TypeString, value)
		_node.InputPlan = value
	}
	if value, ok := arc.mutation.Output(); ok {
		_spec.SetField(applicationrevision.FieldOutput, field.TypeString, value)
		_node.Output = value
	}
	if value, ok := arc.mutation.DeployerType(); ok {
		_spec.SetField(applicationrevision.FieldDeployerType, field.TypeString, value)
		_node.DeployerType = value
	}
	if value, ok := arc.mutation.Duration(); ok {
		_spec.SetField(applicationrevision.FieldDuration, field.TypeInt, value)
		_node.Duration = value
	}
	if value, ok := arc.mutation.PreviousRequiredProviders(); ok {
		_spec.SetField(applicationrevision.FieldPreviousRequiredProviders, field.TypeJSON, value)
		_node.PreviousRequiredProviders = value
	}
	if nodes := arc.mutation.InstanceIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   applicationrevision.InstanceTable,
			Columns: []string{applicationrevision.InstanceColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: applicationinstance.FieldID,
				},
			},
		}
		edge.Schema = arc.schemaConfig.ApplicationRevision
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.InstanceID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := arc.mutation.EnvironmentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   applicationrevision.EnvironmentTable,
			Columns: []string{applicationrevision.EnvironmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: environment.FieldID,
				},
			},
		}
		edge.Schema = arc.schemaConfig.ApplicationRevision
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.EnvironmentID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.ApplicationRevision.Create().
//		SetStatus(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.ApplicationRevisionUpsert) {
//			SetStatus(v+v).
//		}).
//		Exec(ctx)
func (arc *ApplicationRevisionCreate) OnConflict(opts ...sql.ConflictOption) *ApplicationRevisionUpsertOne {
	arc.conflict = opts
	return &ApplicationRevisionUpsertOne{
		create: arc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.ApplicationRevision.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (arc *ApplicationRevisionCreate) OnConflictColumns(columns ...string) *ApplicationRevisionUpsertOne {
	arc.conflict = append(arc.conflict, sql.ConflictColumns(columns...))
	return &ApplicationRevisionUpsertOne{
		create: arc,
	}
}

type (
	// ApplicationRevisionUpsertOne is the builder for "upsert"-ing
	//  one ApplicationRevision node.
	ApplicationRevisionUpsertOne struct {
		create *ApplicationRevisionCreate
	}

	// ApplicationRevisionUpsert is the "OnConflict" setter.
	ApplicationRevisionUpsert struct {
		*sql.UpdateSet
	}
)

// SetStatus sets the "status" field.
func (u *ApplicationRevisionUpsert) SetStatus(v string) *ApplicationRevisionUpsert {
	u.Set(applicationrevision.FieldStatus, v)
	return u
}

// UpdateStatus sets the "status" field to the value that was provided on create.
func (u *ApplicationRevisionUpsert) UpdateStatus() *ApplicationRevisionUpsert {
	u.SetExcluded(applicationrevision.FieldStatus)
	return u
}

// ClearStatus clears the value of the "status" field.
func (u *ApplicationRevisionUpsert) ClearStatus() *ApplicationRevisionUpsert {
	u.SetNull(applicationrevision.FieldStatus)
	return u
}

// SetStatusMessage sets the "statusMessage" field.
func (u *ApplicationRevisionUpsert) SetStatusMessage(v string) *ApplicationRevisionUpsert {
	u.Set(applicationrevision.FieldStatusMessage, v)
	return u
}

// UpdateStatusMessage sets the "statusMessage" field to the value that was provided on create.
func (u *ApplicationRevisionUpsert) UpdateStatusMessage() *ApplicationRevisionUpsert {
	u.SetExcluded(applicationrevision.FieldStatusMessage)
	return u
}

// ClearStatusMessage clears the value of the "statusMessage" field.
func (u *ApplicationRevisionUpsert) ClearStatusMessage() *ApplicationRevisionUpsert {
	u.SetNull(applicationrevision.FieldStatusMessage)
	return u
}

// SetModules sets the "modules" field.
func (u *ApplicationRevisionUpsert) SetModules(v []types.ApplicationModule) *ApplicationRevisionUpsert {
	u.Set(applicationrevision.FieldModules, v)
	return u
}

// UpdateModules sets the "modules" field to the value that was provided on create.
func (u *ApplicationRevisionUpsert) UpdateModules() *ApplicationRevisionUpsert {
	u.SetExcluded(applicationrevision.FieldModules)
	return u
}

// SetSecrets sets the "secrets" field.
func (u *ApplicationRevisionUpsert) SetSecrets(v crypto.Map[string, string]) *ApplicationRevisionUpsert {
	u.Set(applicationrevision.FieldSecrets, v)
	return u
}

// UpdateSecrets sets the "secrets" field to the value that was provided on create.
func (u *ApplicationRevisionUpsert) UpdateSecrets() *ApplicationRevisionUpsert {
	u.SetExcluded(applicationrevision.FieldSecrets)
	return u
}

// SetVariables sets the "variables" field.
func (u *ApplicationRevisionUpsert) SetVariables(v property.Schemas) *ApplicationRevisionUpsert {
	u.Set(applicationrevision.FieldVariables, v)
	return u
}

// UpdateVariables sets the "variables" field to the value that was provided on create.
func (u *ApplicationRevisionUpsert) UpdateVariables() *ApplicationRevisionUpsert {
	u.SetExcluded(applicationrevision.FieldVariables)
	return u
}

// ClearVariables clears the value of the "variables" field.
func (u *ApplicationRevisionUpsert) ClearVariables() *ApplicationRevisionUpsert {
	u.SetNull(applicationrevision.FieldVariables)
	return u
}

// SetInputVariables sets the "inputVariables" field.
func (u *ApplicationRevisionUpsert) SetInputVariables(v property.Values) *ApplicationRevisionUpsert {
	u.Set(applicationrevision.FieldInputVariables, v)
	return u
}

// UpdateInputVariables sets the "inputVariables" field to the value that was provided on create.
func (u *ApplicationRevisionUpsert) UpdateInputVariables() *ApplicationRevisionUpsert {
	u.SetExcluded(applicationrevision.FieldInputVariables)
	return u
}

// SetInputPlan sets the "inputPlan" field.
func (u *ApplicationRevisionUpsert) SetInputPlan(v string) *ApplicationRevisionUpsert {
	u.Set(applicationrevision.FieldInputPlan, v)
	return u
}

// UpdateInputPlan sets the "inputPlan" field to the value that was provided on create.
func (u *ApplicationRevisionUpsert) UpdateInputPlan() *ApplicationRevisionUpsert {
	u.SetExcluded(applicationrevision.FieldInputPlan)
	return u
}

// SetOutput sets the "output" field.
func (u *ApplicationRevisionUpsert) SetOutput(v string) *ApplicationRevisionUpsert {
	u.Set(applicationrevision.FieldOutput, v)
	return u
}

// UpdateOutput sets the "output" field to the value that was provided on create.
func (u *ApplicationRevisionUpsert) UpdateOutput() *ApplicationRevisionUpsert {
	u.SetExcluded(applicationrevision.FieldOutput)
	return u
}

// SetDeployerType sets the "deployerType" field.
func (u *ApplicationRevisionUpsert) SetDeployerType(v string) *ApplicationRevisionUpsert {
	u.Set(applicationrevision.FieldDeployerType, v)
	return u
}

// UpdateDeployerType sets the "deployerType" field to the value that was provided on create.
func (u *ApplicationRevisionUpsert) UpdateDeployerType() *ApplicationRevisionUpsert {
	u.SetExcluded(applicationrevision.FieldDeployerType)
	return u
}

// SetDuration sets the "duration" field.
func (u *ApplicationRevisionUpsert) SetDuration(v int) *ApplicationRevisionUpsert {
	u.Set(applicationrevision.FieldDuration, v)
	return u
}

// UpdateDuration sets the "duration" field to the value that was provided on create.
func (u *ApplicationRevisionUpsert) UpdateDuration() *ApplicationRevisionUpsert {
	u.SetExcluded(applicationrevision.FieldDuration)
	return u
}

// AddDuration adds v to the "duration" field.
func (u *ApplicationRevisionUpsert) AddDuration(v int) *ApplicationRevisionUpsert {
	u.Add(applicationrevision.FieldDuration, v)
	return u
}

// SetPreviousRequiredProviders sets the "previousRequiredProviders" field.
func (u *ApplicationRevisionUpsert) SetPreviousRequiredProviders(v []types.ProviderRequirement) *ApplicationRevisionUpsert {
	u.Set(applicationrevision.FieldPreviousRequiredProviders, v)
	return u
}

// UpdatePreviousRequiredProviders sets the "previousRequiredProviders" field to the value that was provided on create.
func (u *ApplicationRevisionUpsert) UpdatePreviousRequiredProviders() *ApplicationRevisionUpsert {
	u.SetExcluded(applicationrevision.FieldPreviousRequiredProviders)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.ApplicationRevision.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(applicationrevision.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *ApplicationRevisionUpsertOne) UpdateNewValues() *ApplicationRevisionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(applicationrevision.FieldID)
		}
		if _, exists := u.create.mutation.CreateTime(); exists {
			s.SetIgnore(applicationrevision.FieldCreateTime)
		}
		if _, exists := u.create.mutation.InstanceID(); exists {
			s.SetIgnore(applicationrevision.FieldInstanceID)
		}
		if _, exists := u.create.mutation.EnvironmentID(); exists {
			s.SetIgnore(applicationrevision.FieldEnvironmentID)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.ApplicationRevision.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *ApplicationRevisionUpsertOne) Ignore() *ApplicationRevisionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *ApplicationRevisionUpsertOne) DoNothing() *ApplicationRevisionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the ApplicationRevisionCreate.OnConflict
// documentation for more info.
func (u *ApplicationRevisionUpsertOne) Update(set func(*ApplicationRevisionUpsert)) *ApplicationRevisionUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&ApplicationRevisionUpsert{UpdateSet: update})
	}))
	return u
}

// SetStatus sets the "status" field.
func (u *ApplicationRevisionUpsertOne) SetStatus(v string) *ApplicationRevisionUpsertOne {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.SetStatus(v)
	})
}

// UpdateStatus sets the "status" field to the value that was provided on create.
func (u *ApplicationRevisionUpsertOne) UpdateStatus() *ApplicationRevisionUpsertOne {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.UpdateStatus()
	})
}

// ClearStatus clears the value of the "status" field.
func (u *ApplicationRevisionUpsertOne) ClearStatus() *ApplicationRevisionUpsertOne {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.ClearStatus()
	})
}

// SetStatusMessage sets the "statusMessage" field.
func (u *ApplicationRevisionUpsertOne) SetStatusMessage(v string) *ApplicationRevisionUpsertOne {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.SetStatusMessage(v)
	})
}

// UpdateStatusMessage sets the "statusMessage" field to the value that was provided on create.
func (u *ApplicationRevisionUpsertOne) UpdateStatusMessage() *ApplicationRevisionUpsertOne {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.UpdateStatusMessage()
	})
}

// ClearStatusMessage clears the value of the "statusMessage" field.
func (u *ApplicationRevisionUpsertOne) ClearStatusMessage() *ApplicationRevisionUpsertOne {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.ClearStatusMessage()
	})
}

// SetModules sets the "modules" field.
func (u *ApplicationRevisionUpsertOne) SetModules(v []types.ApplicationModule) *ApplicationRevisionUpsertOne {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.SetModules(v)
	})
}

// UpdateModules sets the "modules" field to the value that was provided on create.
func (u *ApplicationRevisionUpsertOne) UpdateModules() *ApplicationRevisionUpsertOne {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.UpdateModules()
	})
}

// SetSecrets sets the "secrets" field.
func (u *ApplicationRevisionUpsertOne) SetSecrets(v crypto.Map[string, string]) *ApplicationRevisionUpsertOne {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.SetSecrets(v)
	})
}

// UpdateSecrets sets the "secrets" field to the value that was provided on create.
func (u *ApplicationRevisionUpsertOne) UpdateSecrets() *ApplicationRevisionUpsertOne {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.UpdateSecrets()
	})
}

// SetVariables sets the "variables" field.
func (u *ApplicationRevisionUpsertOne) SetVariables(v property.Schemas) *ApplicationRevisionUpsertOne {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.SetVariables(v)
	})
}

// UpdateVariables sets the "variables" field to the value that was provided on create.
func (u *ApplicationRevisionUpsertOne) UpdateVariables() *ApplicationRevisionUpsertOne {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.UpdateVariables()
	})
}

// ClearVariables clears the value of the "variables" field.
func (u *ApplicationRevisionUpsertOne) ClearVariables() *ApplicationRevisionUpsertOne {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.ClearVariables()
	})
}

// SetInputVariables sets the "inputVariables" field.
func (u *ApplicationRevisionUpsertOne) SetInputVariables(v property.Values) *ApplicationRevisionUpsertOne {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.SetInputVariables(v)
	})
}

// UpdateInputVariables sets the "inputVariables" field to the value that was provided on create.
func (u *ApplicationRevisionUpsertOne) UpdateInputVariables() *ApplicationRevisionUpsertOne {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.UpdateInputVariables()
	})
}

// SetInputPlan sets the "inputPlan" field.
func (u *ApplicationRevisionUpsertOne) SetInputPlan(v string) *ApplicationRevisionUpsertOne {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.SetInputPlan(v)
	})
}

// UpdateInputPlan sets the "inputPlan" field to the value that was provided on create.
func (u *ApplicationRevisionUpsertOne) UpdateInputPlan() *ApplicationRevisionUpsertOne {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.UpdateInputPlan()
	})
}

// SetOutput sets the "output" field.
func (u *ApplicationRevisionUpsertOne) SetOutput(v string) *ApplicationRevisionUpsertOne {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.SetOutput(v)
	})
}

// UpdateOutput sets the "output" field to the value that was provided on create.
func (u *ApplicationRevisionUpsertOne) UpdateOutput() *ApplicationRevisionUpsertOne {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.UpdateOutput()
	})
}

// SetDeployerType sets the "deployerType" field.
func (u *ApplicationRevisionUpsertOne) SetDeployerType(v string) *ApplicationRevisionUpsertOne {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.SetDeployerType(v)
	})
}

// UpdateDeployerType sets the "deployerType" field to the value that was provided on create.
func (u *ApplicationRevisionUpsertOne) UpdateDeployerType() *ApplicationRevisionUpsertOne {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.UpdateDeployerType()
	})
}

// SetDuration sets the "duration" field.
func (u *ApplicationRevisionUpsertOne) SetDuration(v int) *ApplicationRevisionUpsertOne {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.SetDuration(v)
	})
}

// AddDuration adds v to the "duration" field.
func (u *ApplicationRevisionUpsertOne) AddDuration(v int) *ApplicationRevisionUpsertOne {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.AddDuration(v)
	})
}

// UpdateDuration sets the "duration" field to the value that was provided on create.
func (u *ApplicationRevisionUpsertOne) UpdateDuration() *ApplicationRevisionUpsertOne {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.UpdateDuration()
	})
}

// SetPreviousRequiredProviders sets the "previousRequiredProviders" field.
func (u *ApplicationRevisionUpsertOne) SetPreviousRequiredProviders(v []types.ProviderRequirement) *ApplicationRevisionUpsertOne {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.SetPreviousRequiredProviders(v)
	})
}

// UpdatePreviousRequiredProviders sets the "previousRequiredProviders" field to the value that was provided on create.
func (u *ApplicationRevisionUpsertOne) UpdatePreviousRequiredProviders() *ApplicationRevisionUpsertOne {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.UpdatePreviousRequiredProviders()
	})
}

// Exec executes the query.
func (u *ApplicationRevisionUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("model: missing options for ApplicationRevisionCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *ApplicationRevisionUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *ApplicationRevisionUpsertOne) ID(ctx context.Context) (id oid.ID, err error) {
	if u.create.driver.Dialect() == dialect.MySQL {
		// In case of "ON CONFLICT", there is no way to get back non-numeric ID
		// fields from the database since MySQL does not support the RETURNING clause.
		return id, errors.New("model: ApplicationRevisionUpsertOne.ID is not supported by MySQL driver. Use ApplicationRevisionUpsertOne.Exec instead")
	}
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *ApplicationRevisionUpsertOne) IDX(ctx context.Context) oid.ID {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// ApplicationRevisionCreateBulk is the builder for creating many ApplicationRevision entities in bulk.
type ApplicationRevisionCreateBulk struct {
	config
	builders []*ApplicationRevisionCreate
	conflict []sql.ConflictOption
}

// Save creates the ApplicationRevision entities in the database.
func (arcb *ApplicationRevisionCreateBulk) Save(ctx context.Context) ([]*ApplicationRevision, error) {
	specs := make([]*sqlgraph.CreateSpec, len(arcb.builders))
	nodes := make([]*ApplicationRevision, len(arcb.builders))
	mutators := make([]Mutator, len(arcb.builders))
	for i := range arcb.builders {
		func(i int, root context.Context) {
			builder := arcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ApplicationRevisionMutation)
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
					_, err = mutators[i+1].Mutate(root, arcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = arcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, arcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, arcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (arcb *ApplicationRevisionCreateBulk) SaveX(ctx context.Context) []*ApplicationRevision {
	v, err := arcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (arcb *ApplicationRevisionCreateBulk) Exec(ctx context.Context) error {
	_, err := arcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (arcb *ApplicationRevisionCreateBulk) ExecX(ctx context.Context) {
	if err := arcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.ApplicationRevision.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.ApplicationRevisionUpsert) {
//			SetStatus(v+v).
//		}).
//		Exec(ctx)
func (arcb *ApplicationRevisionCreateBulk) OnConflict(opts ...sql.ConflictOption) *ApplicationRevisionUpsertBulk {
	arcb.conflict = opts
	return &ApplicationRevisionUpsertBulk{
		create: arcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.ApplicationRevision.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (arcb *ApplicationRevisionCreateBulk) OnConflictColumns(columns ...string) *ApplicationRevisionUpsertBulk {
	arcb.conflict = append(arcb.conflict, sql.ConflictColumns(columns...))
	return &ApplicationRevisionUpsertBulk{
		create: arcb,
	}
}

// ApplicationRevisionUpsertBulk is the builder for "upsert"-ing
// a bulk of ApplicationRevision nodes.
type ApplicationRevisionUpsertBulk struct {
	create *ApplicationRevisionCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.ApplicationRevision.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(applicationrevision.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *ApplicationRevisionUpsertBulk) UpdateNewValues() *ApplicationRevisionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(applicationrevision.FieldID)
			}
			if _, exists := b.mutation.CreateTime(); exists {
				s.SetIgnore(applicationrevision.FieldCreateTime)
			}
			if _, exists := b.mutation.InstanceID(); exists {
				s.SetIgnore(applicationrevision.FieldInstanceID)
			}
			if _, exists := b.mutation.EnvironmentID(); exists {
				s.SetIgnore(applicationrevision.FieldEnvironmentID)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.ApplicationRevision.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *ApplicationRevisionUpsertBulk) Ignore() *ApplicationRevisionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *ApplicationRevisionUpsertBulk) DoNothing() *ApplicationRevisionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the ApplicationRevisionCreateBulk.OnConflict
// documentation for more info.
func (u *ApplicationRevisionUpsertBulk) Update(set func(*ApplicationRevisionUpsert)) *ApplicationRevisionUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&ApplicationRevisionUpsert{UpdateSet: update})
	}))
	return u
}

// SetStatus sets the "status" field.
func (u *ApplicationRevisionUpsertBulk) SetStatus(v string) *ApplicationRevisionUpsertBulk {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.SetStatus(v)
	})
}

// UpdateStatus sets the "status" field to the value that was provided on create.
func (u *ApplicationRevisionUpsertBulk) UpdateStatus() *ApplicationRevisionUpsertBulk {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.UpdateStatus()
	})
}

// ClearStatus clears the value of the "status" field.
func (u *ApplicationRevisionUpsertBulk) ClearStatus() *ApplicationRevisionUpsertBulk {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.ClearStatus()
	})
}

// SetStatusMessage sets the "statusMessage" field.
func (u *ApplicationRevisionUpsertBulk) SetStatusMessage(v string) *ApplicationRevisionUpsertBulk {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.SetStatusMessage(v)
	})
}

// UpdateStatusMessage sets the "statusMessage" field to the value that was provided on create.
func (u *ApplicationRevisionUpsertBulk) UpdateStatusMessage() *ApplicationRevisionUpsertBulk {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.UpdateStatusMessage()
	})
}

// ClearStatusMessage clears the value of the "statusMessage" field.
func (u *ApplicationRevisionUpsertBulk) ClearStatusMessage() *ApplicationRevisionUpsertBulk {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.ClearStatusMessage()
	})
}

// SetModules sets the "modules" field.
func (u *ApplicationRevisionUpsertBulk) SetModules(v []types.ApplicationModule) *ApplicationRevisionUpsertBulk {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.SetModules(v)
	})
}

// UpdateModules sets the "modules" field to the value that was provided on create.
func (u *ApplicationRevisionUpsertBulk) UpdateModules() *ApplicationRevisionUpsertBulk {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.UpdateModules()
	})
}

// SetSecrets sets the "secrets" field.
func (u *ApplicationRevisionUpsertBulk) SetSecrets(v crypto.Map[string, string]) *ApplicationRevisionUpsertBulk {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.SetSecrets(v)
	})
}

// UpdateSecrets sets the "secrets" field to the value that was provided on create.
func (u *ApplicationRevisionUpsertBulk) UpdateSecrets() *ApplicationRevisionUpsertBulk {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.UpdateSecrets()
	})
}

// SetVariables sets the "variables" field.
func (u *ApplicationRevisionUpsertBulk) SetVariables(v property.Schemas) *ApplicationRevisionUpsertBulk {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.SetVariables(v)
	})
}

// UpdateVariables sets the "variables" field to the value that was provided on create.
func (u *ApplicationRevisionUpsertBulk) UpdateVariables() *ApplicationRevisionUpsertBulk {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.UpdateVariables()
	})
}

// ClearVariables clears the value of the "variables" field.
func (u *ApplicationRevisionUpsertBulk) ClearVariables() *ApplicationRevisionUpsertBulk {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.ClearVariables()
	})
}

// SetInputVariables sets the "inputVariables" field.
func (u *ApplicationRevisionUpsertBulk) SetInputVariables(v property.Values) *ApplicationRevisionUpsertBulk {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.SetInputVariables(v)
	})
}

// UpdateInputVariables sets the "inputVariables" field to the value that was provided on create.
func (u *ApplicationRevisionUpsertBulk) UpdateInputVariables() *ApplicationRevisionUpsertBulk {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.UpdateInputVariables()
	})
}

// SetInputPlan sets the "inputPlan" field.
func (u *ApplicationRevisionUpsertBulk) SetInputPlan(v string) *ApplicationRevisionUpsertBulk {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.SetInputPlan(v)
	})
}

// UpdateInputPlan sets the "inputPlan" field to the value that was provided on create.
func (u *ApplicationRevisionUpsertBulk) UpdateInputPlan() *ApplicationRevisionUpsertBulk {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.UpdateInputPlan()
	})
}

// SetOutput sets the "output" field.
func (u *ApplicationRevisionUpsertBulk) SetOutput(v string) *ApplicationRevisionUpsertBulk {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.SetOutput(v)
	})
}

// UpdateOutput sets the "output" field to the value that was provided on create.
func (u *ApplicationRevisionUpsertBulk) UpdateOutput() *ApplicationRevisionUpsertBulk {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.UpdateOutput()
	})
}

// SetDeployerType sets the "deployerType" field.
func (u *ApplicationRevisionUpsertBulk) SetDeployerType(v string) *ApplicationRevisionUpsertBulk {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.SetDeployerType(v)
	})
}

// UpdateDeployerType sets the "deployerType" field to the value that was provided on create.
func (u *ApplicationRevisionUpsertBulk) UpdateDeployerType() *ApplicationRevisionUpsertBulk {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.UpdateDeployerType()
	})
}

// SetDuration sets the "duration" field.
func (u *ApplicationRevisionUpsertBulk) SetDuration(v int) *ApplicationRevisionUpsertBulk {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.SetDuration(v)
	})
}

// AddDuration adds v to the "duration" field.
func (u *ApplicationRevisionUpsertBulk) AddDuration(v int) *ApplicationRevisionUpsertBulk {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.AddDuration(v)
	})
}

// UpdateDuration sets the "duration" field to the value that was provided on create.
func (u *ApplicationRevisionUpsertBulk) UpdateDuration() *ApplicationRevisionUpsertBulk {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.UpdateDuration()
	})
}

// SetPreviousRequiredProviders sets the "previousRequiredProviders" field.
func (u *ApplicationRevisionUpsertBulk) SetPreviousRequiredProviders(v []types.ProviderRequirement) *ApplicationRevisionUpsertBulk {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.SetPreviousRequiredProviders(v)
	})
}

// UpdatePreviousRequiredProviders sets the "previousRequiredProviders" field to the value that was provided on create.
func (u *ApplicationRevisionUpsertBulk) UpdatePreviousRequiredProviders() *ApplicationRevisionUpsertBulk {
	return u.Update(func(s *ApplicationRevisionUpsert) {
		s.UpdatePreviousRequiredProviders()
	})
}

// Exec executes the query.
func (u *ApplicationRevisionUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("model: OnConflict was set for builder %d. Set it on the ApplicationRevisionCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("model: missing options for ApplicationRevisionCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *ApplicationRevisionUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
