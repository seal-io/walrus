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
	"github.com/seal-io/seal/pkg/dao/model/moduleversion"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/types"
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

// SetDescription sets the "description" field.
func (mu *ModuleUpdate) SetDescription(s string) *ModuleUpdate {
	mu.mutation.SetDescription(s)
	return mu
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (mu *ModuleUpdate) SetNillableDescription(s *string) *ModuleUpdate {
	if s != nil {
		mu.SetDescription(*s)
	}
	return mu
}

// ClearDescription clears the value of the "description" field.
func (mu *ModuleUpdate) ClearDescription() *ModuleUpdate {
	mu.mutation.ClearDescription()
	return mu
}

// SetIcon sets the "icon" field.
func (mu *ModuleUpdate) SetIcon(s string) *ModuleUpdate {
	mu.mutation.SetIcon(s)
	return mu
}

// SetNillableIcon sets the "icon" field if the given value is not nil.
func (mu *ModuleUpdate) SetNillableIcon(s *string) *ModuleUpdate {
	if s != nil {
		mu.SetIcon(*s)
	}
	return mu
}

// ClearIcon clears the value of the "icon" field.
func (mu *ModuleUpdate) ClearIcon() *ModuleUpdate {
	mu.mutation.ClearIcon()
	return mu
}

// SetLabels sets the "labels" field.
func (mu *ModuleUpdate) SetLabels(m map[string]string) *ModuleUpdate {
	mu.mutation.SetLabels(m)
	return mu
}

// SetSource sets the "source" field.
func (mu *ModuleUpdate) SetSource(s string) *ModuleUpdate {
	mu.mutation.SetSource(s)
	return mu
}

// AddVersionIDs adds the "versions" edge to the ModuleVersion entity by IDs.
func (mu *ModuleUpdate) AddVersionIDs(ids ...types.ID) *ModuleUpdate {
	mu.mutation.AddVersionIDs(ids...)
	return mu
}

// AddVersions adds the "versions" edges to the ModuleVersion entity.
func (mu *ModuleUpdate) AddVersions(m ...*ModuleVersion) *ModuleUpdate {
	ids := make([]types.ID, len(m))
	for i := range m {
		ids[i] = m[i].ID
	}
	return mu.AddVersionIDs(ids...)
}

// Mutation returns the ModuleMutation object of the builder.
func (mu *ModuleUpdate) Mutation() *ModuleMutation {
	return mu.mutation
}

// ClearVersions clears all "versions" edges to the ModuleVersion entity.
func (mu *ModuleUpdate) ClearVersions() *ModuleUpdate {
	mu.mutation.ClearVersions()
	return mu
}

// RemoveVersionIDs removes the "versions" edge to ModuleVersion entities by IDs.
func (mu *ModuleUpdate) RemoveVersionIDs(ids ...types.ID) *ModuleUpdate {
	mu.mutation.RemoveVersionIDs(ids...)
	return mu
}

// RemoveVersions removes "versions" edges to ModuleVersion entities.
func (mu *ModuleUpdate) RemoveVersions(m ...*ModuleVersion) *ModuleUpdate {
	ids := make([]types.ID, len(m))
	for i := range m {
		ids[i] = m[i].ID
	}
	return mu.RemoveVersionIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (mu *ModuleUpdate) Save(ctx context.Context) (int, error) {
	mu.defaults()
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
func (mu *ModuleUpdate) defaults() {
	if _, ok := mu.mutation.UpdateTime(); !ok {
		v := module.UpdateDefaultUpdateTime()
		mu.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (mu *ModuleUpdate) check() error {
	if v, ok := mu.mutation.Source(); ok {
		if err := module.SourceValidator(v); err != nil {
			return &ValidationError{Name: "source", err: fmt.Errorf(`model: validator failed for field "Module.source": %w`, err)}
		}
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (mu *ModuleUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ModuleUpdate {
	mu.modifiers = append(mu.modifiers, modifiers...)
	return mu
}

func (mu *ModuleUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := mu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(module.Table, module.Columns, sqlgraph.NewFieldSpec(module.FieldID, field.TypeString))
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
	if value, ok := mu.mutation.Description(); ok {
		_spec.SetField(module.FieldDescription, field.TypeString, value)
	}
	if mu.mutation.DescriptionCleared() {
		_spec.ClearField(module.FieldDescription, field.TypeString)
	}
	if value, ok := mu.mutation.Icon(); ok {
		_spec.SetField(module.FieldIcon, field.TypeString, value)
	}
	if mu.mutation.IconCleared() {
		_spec.ClearField(module.FieldIcon, field.TypeString)
	}
	if value, ok := mu.mutation.Labels(); ok {
		_spec.SetField(module.FieldLabels, field.TypeJSON, value)
	}
	if value, ok := mu.mutation.Source(); ok {
		_spec.SetField(module.FieldSource, field.TypeString, value)
	}
	if mu.mutation.VersionsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   module.VersionsTable,
			Columns: []string{module.VersionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: moduleversion.FieldID,
				},
			},
		}
		edge.Schema = mu.schemaConfig.ModuleVersion
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := mu.mutation.RemovedVersionsIDs(); len(nodes) > 0 && !mu.mutation.VersionsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   module.VersionsTable,
			Columns: []string{module.VersionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: moduleversion.FieldID,
				},
			},
		}
		edge.Schema = mu.schemaConfig.ModuleVersion
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := mu.mutation.VersionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   module.VersionsTable,
			Columns: []string{module.VersionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: moduleversion.FieldID,
				},
			},
		}
		edge.Schema = mu.schemaConfig.ModuleVersion
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
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

// SetDescription sets the "description" field.
func (muo *ModuleUpdateOne) SetDescription(s string) *ModuleUpdateOne {
	muo.mutation.SetDescription(s)
	return muo
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (muo *ModuleUpdateOne) SetNillableDescription(s *string) *ModuleUpdateOne {
	if s != nil {
		muo.SetDescription(*s)
	}
	return muo
}

// ClearDescription clears the value of the "description" field.
func (muo *ModuleUpdateOne) ClearDescription() *ModuleUpdateOne {
	muo.mutation.ClearDescription()
	return muo
}

// SetIcon sets the "icon" field.
func (muo *ModuleUpdateOne) SetIcon(s string) *ModuleUpdateOne {
	muo.mutation.SetIcon(s)
	return muo
}

// SetNillableIcon sets the "icon" field if the given value is not nil.
func (muo *ModuleUpdateOne) SetNillableIcon(s *string) *ModuleUpdateOne {
	if s != nil {
		muo.SetIcon(*s)
	}
	return muo
}

// ClearIcon clears the value of the "icon" field.
func (muo *ModuleUpdateOne) ClearIcon() *ModuleUpdateOne {
	muo.mutation.ClearIcon()
	return muo
}

// SetLabels sets the "labels" field.
func (muo *ModuleUpdateOne) SetLabels(m map[string]string) *ModuleUpdateOne {
	muo.mutation.SetLabels(m)
	return muo
}

// SetSource sets the "source" field.
func (muo *ModuleUpdateOne) SetSource(s string) *ModuleUpdateOne {
	muo.mutation.SetSource(s)
	return muo
}

// AddVersionIDs adds the "versions" edge to the ModuleVersion entity by IDs.
func (muo *ModuleUpdateOne) AddVersionIDs(ids ...types.ID) *ModuleUpdateOne {
	muo.mutation.AddVersionIDs(ids...)
	return muo
}

// AddVersions adds the "versions" edges to the ModuleVersion entity.
func (muo *ModuleUpdateOne) AddVersions(m ...*ModuleVersion) *ModuleUpdateOne {
	ids := make([]types.ID, len(m))
	for i := range m {
		ids[i] = m[i].ID
	}
	return muo.AddVersionIDs(ids...)
}

// Mutation returns the ModuleMutation object of the builder.
func (muo *ModuleUpdateOne) Mutation() *ModuleMutation {
	return muo.mutation
}

// ClearVersions clears all "versions" edges to the ModuleVersion entity.
func (muo *ModuleUpdateOne) ClearVersions() *ModuleUpdateOne {
	muo.mutation.ClearVersions()
	return muo
}

// RemoveVersionIDs removes the "versions" edge to ModuleVersion entities by IDs.
func (muo *ModuleUpdateOne) RemoveVersionIDs(ids ...types.ID) *ModuleUpdateOne {
	muo.mutation.RemoveVersionIDs(ids...)
	return muo
}

// RemoveVersions removes "versions" edges to ModuleVersion entities.
func (muo *ModuleUpdateOne) RemoveVersions(m ...*ModuleVersion) *ModuleUpdateOne {
	ids := make([]types.ID, len(m))
	for i := range m {
		ids[i] = m[i].ID
	}
	return muo.RemoveVersionIDs(ids...)
}

// Where appends a list predicates to the ModuleUpdate builder.
func (muo *ModuleUpdateOne) Where(ps ...predicate.Module) *ModuleUpdateOne {
	muo.mutation.Where(ps...)
	return muo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (muo *ModuleUpdateOne) Select(field string, fields ...string) *ModuleUpdateOne {
	muo.fields = append([]string{field}, fields...)
	return muo
}

// Save executes the query and returns the updated Module entity.
func (muo *ModuleUpdateOne) Save(ctx context.Context) (*Module, error) {
	muo.defaults()
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
func (muo *ModuleUpdateOne) defaults() {
	if _, ok := muo.mutation.UpdateTime(); !ok {
		v := module.UpdateDefaultUpdateTime()
		muo.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (muo *ModuleUpdateOne) check() error {
	if v, ok := muo.mutation.Source(); ok {
		if err := module.SourceValidator(v); err != nil {
			return &ValidationError{Name: "source", err: fmt.Errorf(`model: validator failed for field "Module.source": %w`, err)}
		}
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (muo *ModuleUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ModuleUpdateOne {
	muo.modifiers = append(muo.modifiers, modifiers...)
	return muo
}

func (muo *ModuleUpdateOne) sqlSave(ctx context.Context) (_node *Module, err error) {
	if err := muo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(module.Table, module.Columns, sqlgraph.NewFieldSpec(module.FieldID, field.TypeString))
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
	if value, ok := muo.mutation.Description(); ok {
		_spec.SetField(module.FieldDescription, field.TypeString, value)
	}
	if muo.mutation.DescriptionCleared() {
		_spec.ClearField(module.FieldDescription, field.TypeString)
	}
	if value, ok := muo.mutation.Icon(); ok {
		_spec.SetField(module.FieldIcon, field.TypeString, value)
	}
	if muo.mutation.IconCleared() {
		_spec.ClearField(module.FieldIcon, field.TypeString)
	}
	if value, ok := muo.mutation.Labels(); ok {
		_spec.SetField(module.FieldLabels, field.TypeJSON, value)
	}
	if value, ok := muo.mutation.Source(); ok {
		_spec.SetField(module.FieldSource, field.TypeString, value)
	}
	if muo.mutation.VersionsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   module.VersionsTable,
			Columns: []string{module.VersionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: moduleversion.FieldID,
				},
			},
		}
		edge.Schema = muo.schemaConfig.ModuleVersion
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := muo.mutation.RemovedVersionsIDs(); len(nodes) > 0 && !muo.mutation.VersionsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   module.VersionsTable,
			Columns: []string{module.VersionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: moduleversion.FieldID,
				},
			},
		}
		edge.Schema = muo.schemaConfig.ModuleVersion
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := muo.mutation.VersionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   module.VersionsTable,
			Columns: []string{module.VersionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: moduleversion.FieldID,
				},
			},
		}
		edge.Schema = muo.schemaConfig.ModuleVersion
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
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
