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

	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
)

// ApplicationResourceUpdate is the builder for updating ApplicationResource entities.
type ApplicationResourceUpdate struct {
	config
	hooks     []Hook
	mutation  *ApplicationResourceMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the ApplicationResourceUpdate builder.
func (aru *ApplicationResourceUpdate) Where(ps ...predicate.ApplicationResource) *ApplicationResourceUpdate {
	aru.mutation.Where(ps...)
	return aru
}

// SetStatus sets the "status" field.
func (aru *ApplicationResourceUpdate) SetStatus(s string) *ApplicationResourceUpdate {
	aru.mutation.SetStatus(s)
	return aru
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (aru *ApplicationResourceUpdate) SetNillableStatus(s *string) *ApplicationResourceUpdate {
	if s != nil {
		aru.SetStatus(*s)
	}
	return aru
}

// ClearStatus clears the value of the "status" field.
func (aru *ApplicationResourceUpdate) ClearStatus() *ApplicationResourceUpdate {
	aru.mutation.ClearStatus()
	return aru
}

// SetStatusMessage sets the "statusMessage" field.
func (aru *ApplicationResourceUpdate) SetStatusMessage(s string) *ApplicationResourceUpdate {
	aru.mutation.SetStatusMessage(s)
	return aru
}

// SetNillableStatusMessage sets the "statusMessage" field if the given value is not nil.
func (aru *ApplicationResourceUpdate) SetNillableStatusMessage(s *string) *ApplicationResourceUpdate {
	if s != nil {
		aru.SetStatusMessage(*s)
	}
	return aru
}

// ClearStatusMessage clears the value of the "statusMessage" field.
func (aru *ApplicationResourceUpdate) ClearStatusMessage() *ApplicationResourceUpdate {
	aru.mutation.ClearStatusMessage()
	return aru
}

// SetUpdateTime sets the "updateTime" field.
func (aru *ApplicationResourceUpdate) SetUpdateTime(t time.Time) *ApplicationResourceUpdate {
	aru.mutation.SetUpdateTime(t)
	return aru
}

// SetModule sets the "module" field.
func (aru *ApplicationResourceUpdate) SetModule(s string) *ApplicationResourceUpdate {
	aru.mutation.SetModule(s)
	return aru
}

// SetType sets the "type" field.
func (aru *ApplicationResourceUpdate) SetType(s string) *ApplicationResourceUpdate {
	aru.mutation.SetType(s)
	return aru
}

// Mutation returns the ApplicationResourceMutation object of the builder.
func (aru *ApplicationResourceUpdate) Mutation() *ApplicationResourceMutation {
	return aru.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (aru *ApplicationResourceUpdate) Save(ctx context.Context) (int, error) {
	if err := aru.defaults(); err != nil {
		return 0, err
	}
	return withHooks[int, ApplicationResourceMutation](ctx, aru.sqlSave, aru.mutation, aru.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (aru *ApplicationResourceUpdate) SaveX(ctx context.Context) int {
	affected, err := aru.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (aru *ApplicationResourceUpdate) Exec(ctx context.Context) error {
	_, err := aru.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (aru *ApplicationResourceUpdate) ExecX(ctx context.Context) {
	if err := aru.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (aru *ApplicationResourceUpdate) defaults() error {
	if _, ok := aru.mutation.UpdateTime(); !ok {
		if applicationresource.UpdateDefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized applicationresource.UpdateDefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := applicationresource.UpdateDefaultUpdateTime()
		aru.mutation.SetUpdateTime(v)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (aru *ApplicationResourceUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ApplicationResourceUpdate {
	aru.modifiers = append(aru.modifiers, modifiers...)
	return aru
}

func (aru *ApplicationResourceUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   applicationresource.Table,
			Columns: applicationresource.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeOther,
				Column: applicationresource.FieldID,
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
		_spec.SetField(applicationresource.FieldStatus, field.TypeString, value)
	}
	if aru.mutation.StatusCleared() {
		_spec.ClearField(applicationresource.FieldStatus, field.TypeString)
	}
	if value, ok := aru.mutation.StatusMessage(); ok {
		_spec.SetField(applicationresource.FieldStatusMessage, field.TypeString, value)
	}
	if aru.mutation.StatusMessageCleared() {
		_spec.ClearField(applicationresource.FieldStatusMessage, field.TypeString)
	}
	if value, ok := aru.mutation.UpdateTime(); ok {
		_spec.SetField(applicationresource.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := aru.mutation.Module(); ok {
		_spec.SetField(applicationresource.FieldModule, field.TypeString, value)
	}
	if value, ok := aru.mutation.GetType(); ok {
		_spec.SetField(applicationresource.FieldType, field.TypeString, value)
	}
	_spec.Node.Schema = aru.schemaConfig.ApplicationResource
	ctx = internal.NewSchemaConfigContext(ctx, aru.schemaConfig)
	_spec.AddModifiers(aru.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, aru.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{applicationresource.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	aru.mutation.done = true
	return n, nil
}

// ApplicationResourceUpdateOne is the builder for updating a single ApplicationResource entity.
type ApplicationResourceUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *ApplicationResourceMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetStatus sets the "status" field.
func (aruo *ApplicationResourceUpdateOne) SetStatus(s string) *ApplicationResourceUpdateOne {
	aruo.mutation.SetStatus(s)
	return aruo
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (aruo *ApplicationResourceUpdateOne) SetNillableStatus(s *string) *ApplicationResourceUpdateOne {
	if s != nil {
		aruo.SetStatus(*s)
	}
	return aruo
}

// ClearStatus clears the value of the "status" field.
func (aruo *ApplicationResourceUpdateOne) ClearStatus() *ApplicationResourceUpdateOne {
	aruo.mutation.ClearStatus()
	return aruo
}

// SetStatusMessage sets the "statusMessage" field.
func (aruo *ApplicationResourceUpdateOne) SetStatusMessage(s string) *ApplicationResourceUpdateOne {
	aruo.mutation.SetStatusMessage(s)
	return aruo
}

// SetNillableStatusMessage sets the "statusMessage" field if the given value is not nil.
func (aruo *ApplicationResourceUpdateOne) SetNillableStatusMessage(s *string) *ApplicationResourceUpdateOne {
	if s != nil {
		aruo.SetStatusMessage(*s)
	}
	return aruo
}

// ClearStatusMessage clears the value of the "statusMessage" field.
func (aruo *ApplicationResourceUpdateOne) ClearStatusMessage() *ApplicationResourceUpdateOne {
	aruo.mutation.ClearStatusMessage()
	return aruo
}

// SetUpdateTime sets the "updateTime" field.
func (aruo *ApplicationResourceUpdateOne) SetUpdateTime(t time.Time) *ApplicationResourceUpdateOne {
	aruo.mutation.SetUpdateTime(t)
	return aruo
}

// SetModule sets the "module" field.
func (aruo *ApplicationResourceUpdateOne) SetModule(s string) *ApplicationResourceUpdateOne {
	aruo.mutation.SetModule(s)
	return aruo
}

// SetType sets the "type" field.
func (aruo *ApplicationResourceUpdateOne) SetType(s string) *ApplicationResourceUpdateOne {
	aruo.mutation.SetType(s)
	return aruo
}

// Mutation returns the ApplicationResourceMutation object of the builder.
func (aruo *ApplicationResourceUpdateOne) Mutation() *ApplicationResourceMutation {
	return aruo.mutation
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (aruo *ApplicationResourceUpdateOne) Select(field string, fields ...string) *ApplicationResourceUpdateOne {
	aruo.fields = append([]string{field}, fields...)
	return aruo
}

// Save executes the query and returns the updated ApplicationResource entity.
func (aruo *ApplicationResourceUpdateOne) Save(ctx context.Context) (*ApplicationResource, error) {
	if err := aruo.defaults(); err != nil {
		return nil, err
	}
	return withHooks[*ApplicationResource, ApplicationResourceMutation](ctx, aruo.sqlSave, aruo.mutation, aruo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (aruo *ApplicationResourceUpdateOne) SaveX(ctx context.Context) *ApplicationResource {
	node, err := aruo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (aruo *ApplicationResourceUpdateOne) Exec(ctx context.Context) error {
	_, err := aruo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (aruo *ApplicationResourceUpdateOne) ExecX(ctx context.Context) {
	if err := aruo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (aruo *ApplicationResourceUpdateOne) defaults() error {
	if _, ok := aruo.mutation.UpdateTime(); !ok {
		if applicationresource.UpdateDefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized applicationresource.UpdateDefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := applicationresource.UpdateDefaultUpdateTime()
		aruo.mutation.SetUpdateTime(v)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (aruo *ApplicationResourceUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ApplicationResourceUpdateOne {
	aruo.modifiers = append(aruo.modifiers, modifiers...)
	return aruo
}

func (aruo *ApplicationResourceUpdateOne) sqlSave(ctx context.Context) (_node *ApplicationResource, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   applicationresource.Table,
			Columns: applicationresource.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeOther,
				Column: applicationresource.FieldID,
			},
		},
	}
	id, ok := aruo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`model: missing "ApplicationResource.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := aruo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, applicationresource.FieldID)
		for _, f := range fields {
			if !applicationresource.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("model: invalid field %q for query", f)}
			}
			if f != applicationresource.FieldID {
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
		_spec.SetField(applicationresource.FieldStatus, field.TypeString, value)
	}
	if aruo.mutation.StatusCleared() {
		_spec.ClearField(applicationresource.FieldStatus, field.TypeString)
	}
	if value, ok := aruo.mutation.StatusMessage(); ok {
		_spec.SetField(applicationresource.FieldStatusMessage, field.TypeString, value)
	}
	if aruo.mutation.StatusMessageCleared() {
		_spec.ClearField(applicationresource.FieldStatusMessage, field.TypeString)
	}
	if value, ok := aruo.mutation.UpdateTime(); ok {
		_spec.SetField(applicationresource.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := aruo.mutation.Module(); ok {
		_spec.SetField(applicationresource.FieldModule, field.TypeString, value)
	}
	if value, ok := aruo.mutation.GetType(); ok {
		_spec.SetField(applicationresource.FieldType, field.TypeString, value)
	}
	_spec.Node.Schema = aruo.schemaConfig.ApplicationResource
	ctx = internal.NewSchemaConfigContext(ctx, aruo.schemaConfig)
	_spec.AddModifiers(aruo.modifiers...)
	_node = &ApplicationResource{config: aruo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, aruo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{applicationresource.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	aruo.mutation.done = true
	return _node, nil
}
