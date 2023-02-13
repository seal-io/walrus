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

	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/module"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
)

// ModuleUpdate is the builder for updating Module entities.
type ModuleUpdate struct {
	config
	hooks     []Hook
	mutation  *ModuleMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the ModuleUpdate builder.
func (mu *ModuleUpdate) Where(ps ...predicate.Module) *ModuleUpdate {
	mu.mutation.Where(ps...)
	return mu
}

// SetStatus sets the "status" field.
func (mu *ModuleUpdate) SetStatus(s string) *ModuleUpdate {
	mu.mutation.SetStatus(s)
	return mu
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (mu *ModuleUpdate) SetNillableStatus(s *string) *ModuleUpdate {
	if s != nil {
		mu.SetStatus(*s)
	}
	return mu
}

// ClearStatus clears the value of the "status" field.
func (mu *ModuleUpdate) ClearStatus() *ModuleUpdate {
	mu.mutation.ClearStatus()
	return mu
}

// SetStatusMessage sets the "statusMessage" field.
func (mu *ModuleUpdate) SetStatusMessage(s string) *ModuleUpdate {
	mu.mutation.SetStatusMessage(s)
	return mu
}

// SetNillableStatusMessage sets the "statusMessage" field if the given value is not nil.
func (mu *ModuleUpdate) SetNillableStatusMessage(s *string) *ModuleUpdate {
	if s != nil {
		mu.SetStatusMessage(*s)
	}
	return mu
}

// ClearStatusMessage clears the value of the "statusMessage" field.
func (mu *ModuleUpdate) ClearStatusMessage() *ModuleUpdate {
	mu.mutation.ClearStatusMessage()
	return mu
}

// SetUpdateTime sets the "updateTime" field.
func (mu *ModuleUpdate) SetUpdateTime(t time.Time) *ModuleUpdate {
	mu.mutation.SetUpdateTime(t)
	return mu
}

// SetSource sets the "source" field.
func (mu *ModuleUpdate) SetSource(s string) *ModuleUpdate {
	mu.mutation.SetSource(s)
	return mu
}

// SetVersion sets the "version" field.
func (mu *ModuleUpdate) SetVersion(s string) *ModuleUpdate {
	mu.mutation.SetVersion(s)
	return mu
}

// SetInputSchema sets the "inputSchema" field.
func (mu *ModuleUpdate) SetInputSchema(m map[string]interface{}) *ModuleUpdate {
	mu.mutation.SetInputSchema(m)
	return mu
}

// ClearInputSchema clears the value of the "inputSchema" field.
func (mu *ModuleUpdate) ClearInputSchema() *ModuleUpdate {
	mu.mutation.ClearInputSchema()
	return mu
}

// SetOutputSchema sets the "outputSchema" field.
func (mu *ModuleUpdate) SetOutputSchema(m map[string]interface{}) *ModuleUpdate {
	mu.mutation.SetOutputSchema(m)
	return mu
}

// ClearOutputSchema clears the value of the "outputSchema" field.
func (mu *ModuleUpdate) ClearOutputSchema() *ModuleUpdate {
	mu.mutation.ClearOutputSchema()
	return mu
}

// Mutation returns the ModuleMutation object of the builder.
func (mu *ModuleUpdate) Mutation() *ModuleMutation {
	return mu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (mu *ModuleUpdate) Save(ctx context.Context) (int, error) {
	if err := mu.defaults(); err != nil {
		return 0, err
	}
	return withHooks[int, ModuleMutation](ctx, mu.sqlSave, mu.mutation, mu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (mu *ModuleUpdate) SaveX(ctx context.Context) int {
	affected, err := mu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (mu *ModuleUpdate) Exec(ctx context.Context) error {
	_, err := mu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mu *ModuleUpdate) ExecX(ctx context.Context) {
	if err := mu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (mu *ModuleUpdate) defaults() error {
	if _, ok := mu.mutation.UpdateTime(); !ok {
		if module.UpdateDefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized module.UpdateDefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := module.UpdateDefaultUpdateTime()
		mu.mutation.SetUpdateTime(v)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (mu *ModuleUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ModuleUpdate {
	mu.modifiers = append(mu.modifiers, modifiers...)
	return mu
}

func (mu *ModuleUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   module.Table,
			Columns: module.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeOther,
				Column: module.FieldID,
			},
		},
	}
	if ps := mu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := mu.mutation.Status(); ok {
		_spec.SetField(module.FieldStatus, field.TypeString, value)
	}
	if mu.mutation.StatusCleared() {
		_spec.ClearField(module.FieldStatus, field.TypeString)
	}
	if value, ok := mu.mutation.StatusMessage(); ok {
		_spec.SetField(module.FieldStatusMessage, field.TypeString, value)
	}
	if mu.mutation.StatusMessageCleared() {
		_spec.ClearField(module.FieldStatusMessage, field.TypeString)
	}
	if value, ok := mu.mutation.UpdateTime(); ok {
		_spec.SetField(module.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := mu.mutation.Source(); ok {
		_spec.SetField(module.FieldSource, field.TypeString, value)
	}
	if value, ok := mu.mutation.Version(); ok {
		_spec.SetField(module.FieldVersion, field.TypeString, value)
	}
	if value, ok := mu.mutation.InputSchema(); ok {
		_spec.SetField(module.FieldInputSchema, field.TypeJSON, value)
	}
	if mu.mutation.InputSchemaCleared() {
		_spec.ClearField(module.FieldInputSchema, field.TypeJSON)
	}
	if value, ok := mu.mutation.OutputSchema(); ok {
		_spec.SetField(module.FieldOutputSchema, field.TypeJSON, value)
	}
	if mu.mutation.OutputSchemaCleared() {
		_spec.ClearField(module.FieldOutputSchema, field.TypeJSON)
	}
	_spec.Node.Schema = mu.schemaConfig.Module
	ctx = internal.NewSchemaConfigContext(ctx, mu.schemaConfig)
	_spec.AddModifiers(mu.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, mu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{module.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	mu.mutation.done = true
	return n, nil
}

// ModuleUpdateOne is the builder for updating a single Module entity.
type ModuleUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *ModuleMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetStatus sets the "status" field.
func (muo *ModuleUpdateOne) SetStatus(s string) *ModuleUpdateOne {
	muo.mutation.SetStatus(s)
	return muo
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (muo *ModuleUpdateOne) SetNillableStatus(s *string) *ModuleUpdateOne {
	if s != nil {
		muo.SetStatus(*s)
	}
	return muo
}

// ClearStatus clears the value of the "status" field.
func (muo *ModuleUpdateOne) ClearStatus() *ModuleUpdateOne {
	muo.mutation.ClearStatus()
	return muo
}

// SetStatusMessage sets the "statusMessage" field.
func (muo *ModuleUpdateOne) SetStatusMessage(s string) *ModuleUpdateOne {
	muo.mutation.SetStatusMessage(s)
	return muo
}

// SetNillableStatusMessage sets the "statusMessage" field if the given value is not nil.
func (muo *ModuleUpdateOne) SetNillableStatusMessage(s *string) *ModuleUpdateOne {
	if s != nil {
		muo.SetStatusMessage(*s)
	}
	return muo
}

// ClearStatusMessage clears the value of the "statusMessage" field.
func (muo *ModuleUpdateOne) ClearStatusMessage() *ModuleUpdateOne {
	muo.mutation.ClearStatusMessage()
	return muo
}

// SetUpdateTime sets the "updateTime" field.
func (muo *ModuleUpdateOne) SetUpdateTime(t time.Time) *ModuleUpdateOne {
	muo.mutation.SetUpdateTime(t)
	return muo
}

// SetSource sets the "source" field.
func (muo *ModuleUpdateOne) SetSource(s string) *ModuleUpdateOne {
	muo.mutation.SetSource(s)
	return muo
}

// SetVersion sets the "version" field.
func (muo *ModuleUpdateOne) SetVersion(s string) *ModuleUpdateOne {
	muo.mutation.SetVersion(s)
	return muo
}

// SetInputSchema sets the "inputSchema" field.
func (muo *ModuleUpdateOne) SetInputSchema(m map[string]interface{}) *ModuleUpdateOne {
	muo.mutation.SetInputSchema(m)
	return muo
}

// ClearInputSchema clears the value of the "inputSchema" field.
func (muo *ModuleUpdateOne) ClearInputSchema() *ModuleUpdateOne {
	muo.mutation.ClearInputSchema()
	return muo
}

// SetOutputSchema sets the "outputSchema" field.
func (muo *ModuleUpdateOne) SetOutputSchema(m map[string]interface{}) *ModuleUpdateOne {
	muo.mutation.SetOutputSchema(m)
	return muo
}

// ClearOutputSchema clears the value of the "outputSchema" field.
func (muo *ModuleUpdateOne) ClearOutputSchema() *ModuleUpdateOne {
	muo.mutation.ClearOutputSchema()
	return muo
}

// Mutation returns the ModuleMutation object of the builder.
func (muo *ModuleUpdateOne) Mutation() *ModuleMutation {
	return muo.mutation
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (muo *ModuleUpdateOne) Select(field string, fields ...string) *ModuleUpdateOne {
	muo.fields = append([]string{field}, fields...)
	return muo
}

// Save executes the query and returns the updated Module entity.
func (muo *ModuleUpdateOne) Save(ctx context.Context) (*Module, error) {
	if err := muo.defaults(); err != nil {
		return nil, err
	}
	return withHooks[*Module, ModuleMutation](ctx, muo.sqlSave, muo.mutation, muo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (muo *ModuleUpdateOne) SaveX(ctx context.Context) *Module {
	node, err := muo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (muo *ModuleUpdateOne) Exec(ctx context.Context) error {
	_, err := muo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (muo *ModuleUpdateOne) ExecX(ctx context.Context) {
	if err := muo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (muo *ModuleUpdateOne) defaults() error {
	if _, ok := muo.mutation.UpdateTime(); !ok {
		if module.UpdateDefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized module.UpdateDefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := module.UpdateDefaultUpdateTime()
		muo.mutation.SetUpdateTime(v)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (muo *ModuleUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ModuleUpdateOne {
	muo.modifiers = append(muo.modifiers, modifiers...)
	return muo
}

func (muo *ModuleUpdateOne) sqlSave(ctx context.Context) (_node *Module, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   module.Table,
			Columns: module.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeOther,
				Column: module.FieldID,
			},
		},
	}
	id, ok := muo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`model: missing "Module.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := muo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, module.FieldID)
		for _, f := range fields {
			if !module.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("model: invalid field %q for query", f)}
			}
			if f != module.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := muo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := muo.mutation.Status(); ok {
		_spec.SetField(module.FieldStatus, field.TypeString, value)
	}
	if muo.mutation.StatusCleared() {
		_spec.ClearField(module.FieldStatus, field.TypeString)
	}
	if value, ok := muo.mutation.StatusMessage(); ok {
		_spec.SetField(module.FieldStatusMessage, field.TypeString, value)
	}
	if muo.mutation.StatusMessageCleared() {
		_spec.ClearField(module.FieldStatusMessage, field.TypeString)
	}
	if value, ok := muo.mutation.UpdateTime(); ok {
		_spec.SetField(module.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := muo.mutation.Source(); ok {
		_spec.SetField(module.FieldSource, field.TypeString, value)
	}
	if value, ok := muo.mutation.Version(); ok {
		_spec.SetField(module.FieldVersion, field.TypeString, value)
	}
	if value, ok := muo.mutation.InputSchema(); ok {
		_spec.SetField(module.FieldInputSchema, field.TypeJSON, value)
	}
	if muo.mutation.InputSchemaCleared() {
		_spec.ClearField(module.FieldInputSchema, field.TypeJSON)
	}
	if value, ok := muo.mutation.OutputSchema(); ok {
		_spec.SetField(module.FieldOutputSchema, field.TypeJSON, value)
	}
	if muo.mutation.OutputSchemaCleared() {
		_spec.ClearField(module.FieldOutputSchema, field.TypeJSON)
	}
	_spec.Node.Schema = muo.schemaConfig.Module
	ctx = internal.NewSchemaConfigContext(ctx, muo.schemaConfig)
	_spec.AddModifiers(muo.modifiers...)
	_node = &Module{config: muo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, muo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{module.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	muo.mutation.done = true
	return _node, nil
}
