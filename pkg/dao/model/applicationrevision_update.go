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
	"entgo.io/ent/dialect/sql/sqljson"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/model/applicationrevision"
	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/types"
)

// ApplicationRevisionUpdate is the builder for updating ApplicationRevision entities.
type ApplicationRevisionUpdate struct {
	config
	hooks     []Hook
	mutation  *ApplicationRevisionMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the ApplicationRevisionUpdate builder.
func (aru *ApplicationRevisionUpdate) Where(ps ...predicate.ApplicationRevision) *ApplicationRevisionUpdate {
	aru.mutation.Where(ps...)
	return aru
}

// SetStatus sets the "status" field.
func (aru *ApplicationRevisionUpdate) SetStatus(s string) *ApplicationRevisionUpdate {
	aru.mutation.SetStatus(s)
	return aru
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (aru *ApplicationRevisionUpdate) SetNillableStatus(s *string) *ApplicationRevisionUpdate {
	if s != nil {
		aru.SetStatus(*s)
	}
	return aru
}

// ClearStatus clears the value of the "status" field.
func (aru *ApplicationRevisionUpdate) ClearStatus() *ApplicationRevisionUpdate {
	aru.mutation.ClearStatus()
	return aru
}

// SetStatusMessage sets the "statusMessage" field.
func (aru *ApplicationRevisionUpdate) SetStatusMessage(s string) *ApplicationRevisionUpdate {
	aru.mutation.SetStatusMessage(s)
	return aru
}

// SetNillableStatusMessage sets the "statusMessage" field if the given value is not nil.
func (aru *ApplicationRevisionUpdate) SetNillableStatusMessage(s *string) *ApplicationRevisionUpdate {
	if s != nil {
		aru.SetStatusMessage(*s)
	}
	return aru
}

// ClearStatusMessage clears the value of the "statusMessage" field.
func (aru *ApplicationRevisionUpdate) ClearStatusMessage() *ApplicationRevisionUpdate {
	aru.mutation.ClearStatusMessage()
	return aru
}

// SetUpdateTime sets the "updateTime" field.
func (aru *ApplicationRevisionUpdate) SetUpdateTime(t time.Time) *ApplicationRevisionUpdate {
	aru.mutation.SetUpdateTime(t)
	return aru
}

// SetModules sets the "modules" field.
func (aru *ApplicationRevisionUpdate) SetModules(tm []types.ApplicationModule) *ApplicationRevisionUpdate {
	aru.mutation.SetModules(tm)
	return aru
}

// AppendModules appends tm to the "modules" field.
func (aru *ApplicationRevisionUpdate) AppendModules(tm []types.ApplicationModule) *ApplicationRevisionUpdate {
	aru.mutation.AppendModules(tm)
	return aru
}

// SetInputVariables sets the "inputVariables" field.
func (aru *ApplicationRevisionUpdate) SetInputVariables(m map[string]interface{}) *ApplicationRevisionUpdate {
	aru.mutation.SetInputVariables(m)
	return aru
}

// SetInputPlan sets the "inputPlan" field.
func (aru *ApplicationRevisionUpdate) SetInputPlan(s string) *ApplicationRevisionUpdate {
	aru.mutation.SetInputPlan(s)
	return aru
}

// SetOutput sets the "output" field.
func (aru *ApplicationRevisionUpdate) SetOutput(s string) *ApplicationRevisionUpdate {
	aru.mutation.SetOutput(s)
	return aru
}

// Mutation returns the ApplicationRevisionMutation object of the builder.
func (aru *ApplicationRevisionUpdate) Mutation() *ApplicationRevisionMutation {
	return aru.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (aru *ApplicationRevisionUpdate) Save(ctx context.Context) (int, error) {
	if err := aru.defaults(); err != nil {
		return 0, err
	}
	return withHooks[int, ApplicationRevisionMutation](ctx, aru.sqlSave, aru.mutation, aru.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (aru *ApplicationRevisionUpdate) SaveX(ctx context.Context) int {
	affected, err := aru.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (aru *ApplicationRevisionUpdate) Exec(ctx context.Context) error {
	_, err := aru.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (aru *ApplicationRevisionUpdate) ExecX(ctx context.Context) {
	if err := aru.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (aru *ApplicationRevisionUpdate) defaults() error {
	if _, ok := aru.mutation.UpdateTime(); !ok {
		if applicationrevision.UpdateDefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized applicationrevision.UpdateDefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := applicationrevision.UpdateDefaultUpdateTime()
		aru.mutation.SetUpdateTime(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (aru *ApplicationRevisionUpdate) check() error {
	if _, ok := aru.mutation.ApplicationID(); aru.mutation.ApplicationCleared() && !ok {
		return errors.New(`model: clearing a required unique edge "ApplicationRevision.application"`)
	}
	if _, ok := aru.mutation.EnvironmentID(); aru.mutation.EnvironmentCleared() && !ok {
		return errors.New(`model: clearing a required unique edge "ApplicationRevision.environment"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (aru *ApplicationRevisionUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ApplicationRevisionUpdate {
	aru.modifiers = append(aru.modifiers, modifiers...)
	return aru
}

func (aru *ApplicationRevisionUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := aru.check(); err != nil {
		return n, err
	}
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   applicationrevision.Table,
			Columns: applicationrevision.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: applicationrevision.FieldID,
			},
		},
	}
	if ps := aru.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := aru.mutation.Status(); ok {
		_spec.SetField(applicationrevision.FieldStatus, field.TypeString, value)
	}
	if aru.mutation.StatusCleared() {
		_spec.ClearField(applicationrevision.FieldStatus, field.TypeString)
	}
	if value, ok := aru.mutation.StatusMessage(); ok {
		_spec.SetField(applicationrevision.FieldStatusMessage, field.TypeString, value)
	}
	if aru.mutation.StatusMessageCleared() {
		_spec.ClearField(applicationrevision.FieldStatusMessage, field.TypeString)
	}
	if value, ok := aru.mutation.UpdateTime(); ok {
		_spec.SetField(applicationrevision.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := aru.mutation.Modules(); ok {
		_spec.SetField(applicationrevision.FieldModules, field.TypeJSON, value)
	}
	if value, ok := aru.mutation.AppendedModules(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, applicationrevision.FieldModules, value)
		})
	}
	if value, ok := aru.mutation.InputVariables(); ok {
		_spec.SetField(applicationrevision.FieldInputVariables, field.TypeJSON, value)
	}
	if value, ok := aru.mutation.InputPlan(); ok {
		_spec.SetField(applicationrevision.FieldInputPlan, field.TypeString, value)
	}
	if value, ok := aru.mutation.Output(); ok {
		_spec.SetField(applicationrevision.FieldOutput, field.TypeString, value)
	}
	_spec.Node.Schema = aru.schemaConfig.ApplicationRevision
	ctx = internal.NewSchemaConfigContext(ctx, aru.schemaConfig)
	_spec.AddModifiers(aru.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, aru.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{applicationrevision.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	aru.mutation.done = true
	return n, nil
}

// ApplicationRevisionUpdateOne is the builder for updating a single ApplicationRevision entity.
type ApplicationRevisionUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *ApplicationRevisionMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetStatus sets the "status" field.
func (aruo *ApplicationRevisionUpdateOne) SetStatus(s string) *ApplicationRevisionUpdateOne {
	aruo.mutation.SetStatus(s)
	return aruo
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (aruo *ApplicationRevisionUpdateOne) SetNillableStatus(s *string) *ApplicationRevisionUpdateOne {
	if s != nil {
		aruo.SetStatus(*s)
	}
	return aruo
}

// ClearStatus clears the value of the "status" field.
func (aruo *ApplicationRevisionUpdateOne) ClearStatus() *ApplicationRevisionUpdateOne {
	aruo.mutation.ClearStatus()
	return aruo
}

// SetStatusMessage sets the "statusMessage" field.
func (aruo *ApplicationRevisionUpdateOne) SetStatusMessage(s string) *ApplicationRevisionUpdateOne {
	aruo.mutation.SetStatusMessage(s)
	return aruo
}

// SetNillableStatusMessage sets the "statusMessage" field if the given value is not nil.
func (aruo *ApplicationRevisionUpdateOne) SetNillableStatusMessage(s *string) *ApplicationRevisionUpdateOne {
	if s != nil {
		aruo.SetStatusMessage(*s)
	}
	return aruo
}

// ClearStatusMessage clears the value of the "statusMessage" field.
func (aruo *ApplicationRevisionUpdateOne) ClearStatusMessage() *ApplicationRevisionUpdateOne {
	aruo.mutation.ClearStatusMessage()
	return aruo
}

// SetUpdateTime sets the "updateTime" field.
func (aruo *ApplicationRevisionUpdateOne) SetUpdateTime(t time.Time) *ApplicationRevisionUpdateOne {
	aruo.mutation.SetUpdateTime(t)
	return aruo
}

// SetModules sets the "modules" field.
func (aruo *ApplicationRevisionUpdateOne) SetModules(tm []types.ApplicationModule) *ApplicationRevisionUpdateOne {
	aruo.mutation.SetModules(tm)
	return aruo
}

// AppendModules appends tm to the "modules" field.
func (aruo *ApplicationRevisionUpdateOne) AppendModules(tm []types.ApplicationModule) *ApplicationRevisionUpdateOne {
	aruo.mutation.AppendModules(tm)
	return aruo
}

// SetInputVariables sets the "inputVariables" field.
func (aruo *ApplicationRevisionUpdateOne) SetInputVariables(m map[string]interface{}) *ApplicationRevisionUpdateOne {
	aruo.mutation.SetInputVariables(m)
	return aruo
}

// SetInputPlan sets the "inputPlan" field.
func (aruo *ApplicationRevisionUpdateOne) SetInputPlan(s string) *ApplicationRevisionUpdateOne {
	aruo.mutation.SetInputPlan(s)
	return aruo
}

// SetOutput sets the "output" field.
func (aruo *ApplicationRevisionUpdateOne) SetOutput(s string) *ApplicationRevisionUpdateOne {
	aruo.mutation.SetOutput(s)
	return aruo
}

// Mutation returns the ApplicationRevisionMutation object of the builder.
func (aruo *ApplicationRevisionUpdateOne) Mutation() *ApplicationRevisionMutation {
	return aruo.mutation
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (aruo *ApplicationRevisionUpdateOne) Select(field string, fields ...string) *ApplicationRevisionUpdateOne {
	aruo.fields = append([]string{field}, fields...)
	return aruo
}

// Save executes the query and returns the updated ApplicationRevision entity.
func (aruo *ApplicationRevisionUpdateOne) Save(ctx context.Context) (*ApplicationRevision, error) {
	if err := aruo.defaults(); err != nil {
		return nil, err
	}
	return withHooks[*ApplicationRevision, ApplicationRevisionMutation](ctx, aruo.sqlSave, aruo.mutation, aruo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (aruo *ApplicationRevisionUpdateOne) SaveX(ctx context.Context) *ApplicationRevision {
	node, err := aruo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (aruo *ApplicationRevisionUpdateOne) Exec(ctx context.Context) error {
	_, err := aruo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (aruo *ApplicationRevisionUpdateOne) ExecX(ctx context.Context) {
	if err := aruo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (aruo *ApplicationRevisionUpdateOne) defaults() error {
	if _, ok := aruo.mutation.UpdateTime(); !ok {
		if applicationrevision.UpdateDefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized applicationrevision.UpdateDefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := applicationrevision.UpdateDefaultUpdateTime()
		aruo.mutation.SetUpdateTime(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (aruo *ApplicationRevisionUpdateOne) check() error {
	if _, ok := aruo.mutation.ApplicationID(); aruo.mutation.ApplicationCleared() && !ok {
		return errors.New(`model: clearing a required unique edge "ApplicationRevision.application"`)
	}
	if _, ok := aruo.mutation.EnvironmentID(); aruo.mutation.EnvironmentCleared() && !ok {
		return errors.New(`model: clearing a required unique edge "ApplicationRevision.environment"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (aruo *ApplicationRevisionUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ApplicationRevisionUpdateOne {
	aruo.modifiers = append(aruo.modifiers, modifiers...)
	return aruo
}

func (aruo *ApplicationRevisionUpdateOne) sqlSave(ctx context.Context) (_node *ApplicationRevision, err error) {
	if err := aruo.check(); err != nil {
		return _node, err
	}
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   applicationrevision.Table,
			Columns: applicationrevision.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: applicationrevision.FieldID,
			},
		},
	}
	id, ok := aruo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`model: missing "ApplicationRevision.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := aruo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, applicationrevision.FieldID)
		for _, f := range fields {
			if !applicationrevision.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("model: invalid field %q for query", f)}
			}
			if f != applicationrevision.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := aruo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := aruo.mutation.Status(); ok {
		_spec.SetField(applicationrevision.FieldStatus, field.TypeString, value)
	}
	if aruo.mutation.StatusCleared() {
		_spec.ClearField(applicationrevision.FieldStatus, field.TypeString)
	}
	if value, ok := aruo.mutation.StatusMessage(); ok {
		_spec.SetField(applicationrevision.FieldStatusMessage, field.TypeString, value)
	}
	if aruo.mutation.StatusMessageCleared() {
		_spec.ClearField(applicationrevision.FieldStatusMessage, field.TypeString)
	}
	if value, ok := aruo.mutation.UpdateTime(); ok {
		_spec.SetField(applicationrevision.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := aruo.mutation.Modules(); ok {
		_spec.SetField(applicationrevision.FieldModules, field.TypeJSON, value)
	}
	if value, ok := aruo.mutation.AppendedModules(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, applicationrevision.FieldModules, value)
		})
	}
	if value, ok := aruo.mutation.InputVariables(); ok {
		_spec.SetField(applicationrevision.FieldInputVariables, field.TypeJSON, value)
	}
	if value, ok := aruo.mutation.InputPlan(); ok {
		_spec.SetField(applicationrevision.FieldInputPlan, field.TypeString, value)
	}
	if value, ok := aruo.mutation.Output(); ok {
		_spec.SetField(applicationrevision.FieldOutput, field.TypeString, value)
	}
	_spec.Node.Schema = aruo.schemaConfig.ApplicationRevision
	ctx = internal.NewSchemaConfigContext(ctx, aruo.schemaConfig)
	_spec.AddModifiers(aruo.modifiers...)
	_node = &ApplicationRevision{config: aruo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, aruo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{applicationrevision.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	aruo.mutation.done = true
	return _node, nil
}
