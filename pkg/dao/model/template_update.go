// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "ent". DO NOT EDIT.

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
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/template"
	"github.com/seal-io/seal/pkg/dao/model/templateversion"
	"github.com/seal-io/seal/pkg/dao/types/object"
)

// TemplateUpdate is the builder for updating Template entities.
type TemplateUpdate struct {
	config
	hooks     []Hook
	mutation  *TemplateMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the TemplateUpdate builder.
func (tu *TemplateUpdate) Where(ps ...predicate.Template) *TemplateUpdate {
	tu.mutation.Where(ps...)
	return tu
}

// SetUpdateTime sets the "updateTime" field.
func (tu *TemplateUpdate) SetUpdateTime(t time.Time) *TemplateUpdate {
	tu.mutation.SetUpdateTime(t)
	return tu
}

// SetStatus sets the "status" field.
func (tu *TemplateUpdate) SetStatus(s string) *TemplateUpdate {
	tu.mutation.SetStatus(s)
	return tu
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (tu *TemplateUpdate) SetNillableStatus(s *string) *TemplateUpdate {
	if s != nil {
		tu.SetStatus(*s)
	}
	return tu
}

// ClearStatus clears the value of the "status" field.
func (tu *TemplateUpdate) ClearStatus() *TemplateUpdate {
	tu.mutation.ClearStatus()
	return tu
}

// SetStatusMessage sets the "statusMessage" field.
func (tu *TemplateUpdate) SetStatusMessage(s string) *TemplateUpdate {
	tu.mutation.SetStatusMessage(s)
	return tu
}

// SetNillableStatusMessage sets the "statusMessage" field if the given value is not nil.
func (tu *TemplateUpdate) SetNillableStatusMessage(s *string) *TemplateUpdate {
	if s != nil {
		tu.SetStatusMessage(*s)
	}
	return tu
}

// ClearStatusMessage clears the value of the "statusMessage" field.
func (tu *TemplateUpdate) ClearStatusMessage() *TemplateUpdate {
	tu.mutation.ClearStatusMessage()
	return tu
}

// SetDescription sets the "description" field.
func (tu *TemplateUpdate) SetDescription(s string) *TemplateUpdate {
	tu.mutation.SetDescription(s)
	return tu
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (tu *TemplateUpdate) SetNillableDescription(s *string) *TemplateUpdate {
	if s != nil {
		tu.SetDescription(*s)
	}
	return tu
}

// ClearDescription clears the value of the "description" field.
func (tu *TemplateUpdate) ClearDescription() *TemplateUpdate {
	tu.mutation.ClearDescription()
	return tu
}

// SetIcon sets the "icon" field.
func (tu *TemplateUpdate) SetIcon(s string) *TemplateUpdate {
	tu.mutation.SetIcon(s)
	return tu
}

// SetNillableIcon sets the "icon" field if the given value is not nil.
func (tu *TemplateUpdate) SetNillableIcon(s *string) *TemplateUpdate {
	if s != nil {
		tu.SetIcon(*s)
	}
	return tu
}

// ClearIcon clears the value of the "icon" field.
func (tu *TemplateUpdate) ClearIcon() *TemplateUpdate {
	tu.mutation.ClearIcon()
	return tu
}

// SetLabels sets the "labels" field.
func (tu *TemplateUpdate) SetLabels(m map[string]string) *TemplateUpdate {
	tu.mutation.SetLabels(m)
	return tu
}

// SetSource sets the "source" field.
func (tu *TemplateUpdate) SetSource(s string) *TemplateUpdate {
	tu.mutation.SetSource(s)
	return tu
}

// AddVersionIDs adds the "versions" edge to the TemplateVersion entity by IDs.
func (tu *TemplateUpdate) AddVersionIDs(ids ...object.ID) *TemplateUpdate {
	tu.mutation.AddVersionIDs(ids...)
	return tu
}

// AddVersions adds the "versions" edges to the TemplateVersion entity.
func (tu *TemplateUpdate) AddVersions(t ...*TemplateVersion) *TemplateUpdate {
	ids := make([]object.ID, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return tu.AddVersionIDs(ids...)
}

// Mutation returns the TemplateMutation object of the builder.
func (tu *TemplateUpdate) Mutation() *TemplateMutation {
	return tu.mutation
}

// ClearVersions clears all "versions" edges to the TemplateVersion entity.
func (tu *TemplateUpdate) ClearVersions() *TemplateUpdate {
	tu.mutation.ClearVersions()
	return tu
}

// RemoveVersionIDs removes the "versions" edge to TemplateVersion entities by IDs.
func (tu *TemplateUpdate) RemoveVersionIDs(ids ...object.ID) *TemplateUpdate {
	tu.mutation.RemoveVersionIDs(ids...)
	return tu
}

// RemoveVersions removes "versions" edges to TemplateVersion entities.
func (tu *TemplateUpdate) RemoveVersions(t ...*TemplateVersion) *TemplateUpdate {
	ids := make([]object.ID, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return tu.RemoveVersionIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (tu *TemplateUpdate) Save(ctx context.Context) (int, error) {
	tu.defaults()
	return withHooks[int, TemplateMutation](ctx, tu.sqlSave, tu.mutation, tu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (tu *TemplateUpdate) SaveX(ctx context.Context) int {
	affected, err := tu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (tu *TemplateUpdate) Exec(ctx context.Context) error {
	_, err := tu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tu *TemplateUpdate) ExecX(ctx context.Context) {
	if err := tu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (tu *TemplateUpdate) defaults() {
	if _, ok := tu.mutation.UpdateTime(); !ok {
		v := template.UpdateDefaultUpdateTime()
		tu.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (tu *TemplateUpdate) check() error {
	if v, ok := tu.mutation.Source(); ok {
		if err := template.SourceValidator(v); err != nil {
			return &ValidationError{Name: "source", err: fmt.Errorf(`model: validator failed for field "Template.source": %w`, err)}
		}
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (tu *TemplateUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *TemplateUpdate {
	tu.modifiers = append(tu.modifiers, modifiers...)
	return tu
}

func (tu *TemplateUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := tu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(template.Table, template.Columns, sqlgraph.NewFieldSpec(template.FieldID, field.TypeString))
	if ps := tu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := tu.mutation.UpdateTime(); ok {
		_spec.SetField(template.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := tu.mutation.Status(); ok {
		_spec.SetField(template.FieldStatus, field.TypeString, value)
	}
	if tu.mutation.StatusCleared() {
		_spec.ClearField(template.FieldStatus, field.TypeString)
	}
	if value, ok := tu.mutation.StatusMessage(); ok {
		_spec.SetField(template.FieldStatusMessage, field.TypeString, value)
	}
	if tu.mutation.StatusMessageCleared() {
		_spec.ClearField(template.FieldStatusMessage, field.TypeString)
	}
	if value, ok := tu.mutation.Description(); ok {
		_spec.SetField(template.FieldDescription, field.TypeString, value)
	}
	if tu.mutation.DescriptionCleared() {
		_spec.ClearField(template.FieldDescription, field.TypeString)
	}
	if value, ok := tu.mutation.Icon(); ok {
		_spec.SetField(template.FieldIcon, field.TypeString, value)
	}
	if tu.mutation.IconCleared() {
		_spec.ClearField(template.FieldIcon, field.TypeString)
	}
	if value, ok := tu.mutation.Labels(); ok {
		_spec.SetField(template.FieldLabels, field.TypeJSON, value)
	}
	if value, ok := tu.mutation.Source(); ok {
		_spec.SetField(template.FieldSource, field.TypeString, value)
	}
	if tu.mutation.VersionsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   template.VersionsTable,
			Columns: []string{template.VersionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(templateversion.FieldID, field.TypeString),
			},
		}
		edge.Schema = tu.schemaConfig.TemplateVersion
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tu.mutation.RemovedVersionsIDs(); len(nodes) > 0 && !tu.mutation.VersionsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   template.VersionsTable,
			Columns: []string{template.VersionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(templateversion.FieldID, field.TypeString),
			},
		}
		edge.Schema = tu.schemaConfig.TemplateVersion
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tu.mutation.VersionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   template.VersionsTable,
			Columns: []string{template.VersionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(templateversion.FieldID, field.TypeString),
			},
		}
		edge.Schema = tu.schemaConfig.TemplateVersion
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.Node.Schema = tu.schemaConfig.Template
	ctx = internal.NewSchemaConfigContext(ctx, tu.schemaConfig)
	_spec.AddModifiers(tu.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, tu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{template.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	tu.mutation.done = true
	return n, nil
}

// TemplateUpdateOne is the builder for updating a single Template entity.
type TemplateUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *TemplateMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetUpdateTime sets the "updateTime" field.
func (tuo *TemplateUpdateOne) SetUpdateTime(t time.Time) *TemplateUpdateOne {
	tuo.mutation.SetUpdateTime(t)
	return tuo
}

// SetStatus sets the "status" field.
func (tuo *TemplateUpdateOne) SetStatus(s string) *TemplateUpdateOne {
	tuo.mutation.SetStatus(s)
	return tuo
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (tuo *TemplateUpdateOne) SetNillableStatus(s *string) *TemplateUpdateOne {
	if s != nil {
		tuo.SetStatus(*s)
	}
	return tuo
}

// ClearStatus clears the value of the "status" field.
func (tuo *TemplateUpdateOne) ClearStatus() *TemplateUpdateOne {
	tuo.mutation.ClearStatus()
	return tuo
}

// SetStatusMessage sets the "statusMessage" field.
func (tuo *TemplateUpdateOne) SetStatusMessage(s string) *TemplateUpdateOne {
	tuo.mutation.SetStatusMessage(s)
	return tuo
}

// SetNillableStatusMessage sets the "statusMessage" field if the given value is not nil.
func (tuo *TemplateUpdateOne) SetNillableStatusMessage(s *string) *TemplateUpdateOne {
	if s != nil {
		tuo.SetStatusMessage(*s)
	}
	return tuo
}

// ClearStatusMessage clears the value of the "statusMessage" field.
func (tuo *TemplateUpdateOne) ClearStatusMessage() *TemplateUpdateOne {
	tuo.mutation.ClearStatusMessage()
	return tuo
}

// SetDescription sets the "description" field.
func (tuo *TemplateUpdateOne) SetDescription(s string) *TemplateUpdateOne {
	tuo.mutation.SetDescription(s)
	return tuo
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (tuo *TemplateUpdateOne) SetNillableDescription(s *string) *TemplateUpdateOne {
	if s != nil {
		tuo.SetDescription(*s)
	}
	return tuo
}

// ClearDescription clears the value of the "description" field.
func (tuo *TemplateUpdateOne) ClearDescription() *TemplateUpdateOne {
	tuo.mutation.ClearDescription()
	return tuo
}

// SetIcon sets the "icon" field.
func (tuo *TemplateUpdateOne) SetIcon(s string) *TemplateUpdateOne {
	tuo.mutation.SetIcon(s)
	return tuo
}

// SetNillableIcon sets the "icon" field if the given value is not nil.
func (tuo *TemplateUpdateOne) SetNillableIcon(s *string) *TemplateUpdateOne {
	if s != nil {
		tuo.SetIcon(*s)
	}
	return tuo
}

// ClearIcon clears the value of the "icon" field.
func (tuo *TemplateUpdateOne) ClearIcon() *TemplateUpdateOne {
	tuo.mutation.ClearIcon()
	return tuo
}

// SetLabels sets the "labels" field.
func (tuo *TemplateUpdateOne) SetLabels(m map[string]string) *TemplateUpdateOne {
	tuo.mutation.SetLabels(m)
	return tuo
}

// SetSource sets the "source" field.
func (tuo *TemplateUpdateOne) SetSource(s string) *TemplateUpdateOne {
	tuo.mutation.SetSource(s)
	return tuo
}

// AddVersionIDs adds the "versions" edge to the TemplateVersion entity by IDs.
func (tuo *TemplateUpdateOne) AddVersionIDs(ids ...object.ID) *TemplateUpdateOne {
	tuo.mutation.AddVersionIDs(ids...)
	return tuo
}

// AddVersions adds the "versions" edges to the TemplateVersion entity.
func (tuo *TemplateUpdateOne) AddVersions(t ...*TemplateVersion) *TemplateUpdateOne {
	ids := make([]object.ID, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return tuo.AddVersionIDs(ids...)
}

// Mutation returns the TemplateMutation object of the builder.
func (tuo *TemplateUpdateOne) Mutation() *TemplateMutation {
	return tuo.mutation
}

// ClearVersions clears all "versions" edges to the TemplateVersion entity.
func (tuo *TemplateUpdateOne) ClearVersions() *TemplateUpdateOne {
	tuo.mutation.ClearVersions()
	return tuo
}

// RemoveVersionIDs removes the "versions" edge to TemplateVersion entities by IDs.
func (tuo *TemplateUpdateOne) RemoveVersionIDs(ids ...object.ID) *TemplateUpdateOne {
	tuo.mutation.RemoveVersionIDs(ids...)
	return tuo
}

// RemoveVersions removes "versions" edges to TemplateVersion entities.
func (tuo *TemplateUpdateOne) RemoveVersions(t ...*TemplateVersion) *TemplateUpdateOne {
	ids := make([]object.ID, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return tuo.RemoveVersionIDs(ids...)
}

// Where appends a list predicates to the TemplateUpdate builder.
func (tuo *TemplateUpdateOne) Where(ps ...predicate.Template) *TemplateUpdateOne {
	tuo.mutation.Where(ps...)
	return tuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (tuo *TemplateUpdateOne) Select(field string, fields ...string) *TemplateUpdateOne {
	tuo.fields = append([]string{field}, fields...)
	return tuo
}

// Save executes the query and returns the updated Template entity.
func (tuo *TemplateUpdateOne) Save(ctx context.Context) (*Template, error) {
	tuo.defaults()
	return withHooks[*Template, TemplateMutation](ctx, tuo.sqlSave, tuo.mutation, tuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (tuo *TemplateUpdateOne) SaveX(ctx context.Context) *Template {
	node, err := tuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (tuo *TemplateUpdateOne) Exec(ctx context.Context) error {
	_, err := tuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tuo *TemplateUpdateOne) ExecX(ctx context.Context) {
	if err := tuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (tuo *TemplateUpdateOne) defaults() {
	if _, ok := tuo.mutation.UpdateTime(); !ok {
		v := template.UpdateDefaultUpdateTime()
		tuo.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (tuo *TemplateUpdateOne) check() error {
	if v, ok := tuo.mutation.Source(); ok {
		if err := template.SourceValidator(v); err != nil {
			return &ValidationError{Name: "source", err: fmt.Errorf(`model: validator failed for field "Template.source": %w`, err)}
		}
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (tuo *TemplateUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *TemplateUpdateOne {
	tuo.modifiers = append(tuo.modifiers, modifiers...)
	return tuo
}

func (tuo *TemplateUpdateOne) sqlSave(ctx context.Context) (_node *Template, err error) {
	if err := tuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(template.Table, template.Columns, sqlgraph.NewFieldSpec(template.FieldID, field.TypeString))
	id, ok := tuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`model: missing "Template.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := tuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, template.FieldID)
		for _, f := range fields {
			if !template.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("model: invalid field %q for query", f)}
			}
			if f != template.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := tuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := tuo.mutation.UpdateTime(); ok {
		_spec.SetField(template.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := tuo.mutation.Status(); ok {
		_spec.SetField(template.FieldStatus, field.TypeString, value)
	}
	if tuo.mutation.StatusCleared() {
		_spec.ClearField(template.FieldStatus, field.TypeString)
	}
	if value, ok := tuo.mutation.StatusMessage(); ok {
		_spec.SetField(template.FieldStatusMessage, field.TypeString, value)
	}
	if tuo.mutation.StatusMessageCleared() {
		_spec.ClearField(template.FieldStatusMessage, field.TypeString)
	}
	if value, ok := tuo.mutation.Description(); ok {
		_spec.SetField(template.FieldDescription, field.TypeString, value)
	}
	if tuo.mutation.DescriptionCleared() {
		_spec.ClearField(template.FieldDescription, field.TypeString)
	}
	if value, ok := tuo.mutation.Icon(); ok {
		_spec.SetField(template.FieldIcon, field.TypeString, value)
	}
	if tuo.mutation.IconCleared() {
		_spec.ClearField(template.FieldIcon, field.TypeString)
	}
	if value, ok := tuo.mutation.Labels(); ok {
		_spec.SetField(template.FieldLabels, field.TypeJSON, value)
	}
	if value, ok := tuo.mutation.Source(); ok {
		_spec.SetField(template.FieldSource, field.TypeString, value)
	}
	if tuo.mutation.VersionsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   template.VersionsTable,
			Columns: []string{template.VersionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(templateversion.FieldID, field.TypeString),
			},
		}
		edge.Schema = tuo.schemaConfig.TemplateVersion
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tuo.mutation.RemovedVersionsIDs(); len(nodes) > 0 && !tuo.mutation.VersionsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   template.VersionsTable,
			Columns: []string{template.VersionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(templateversion.FieldID, field.TypeString),
			},
		}
		edge.Schema = tuo.schemaConfig.TemplateVersion
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tuo.mutation.VersionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   template.VersionsTable,
			Columns: []string{template.VersionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(templateversion.FieldID, field.TypeString),
			},
		}
		edge.Schema = tuo.schemaConfig.TemplateVersion
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.Node.Schema = tuo.schemaConfig.Template
	ctx = internal.NewSchemaConfigContext(ctx, tuo.schemaConfig)
	_spec.AddModifiers(tuo.modifiers...)
	_node = &Template{config: tuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, tuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{template.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	tuo.mutation.done = true
	return _node, nil
}
