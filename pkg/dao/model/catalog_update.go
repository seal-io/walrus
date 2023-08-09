// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "seal". DO NOT EDIT.

package model

import (
	"context"
	stdsql "database/sql"
	"errors"
	"fmt"
	"reflect"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/model/catalog"
	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/template"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/object"
	"github.com/seal-io/seal/pkg/dao/types/status"
)

// CatalogUpdate is the builder for updating Catalog entities.
type CatalogUpdate struct {
	config
	hooks     []Hook
	mutation  *CatalogMutation
	modifiers []func(*sql.UpdateBuilder)
	object    *Catalog
}

// Where appends a list predicates to the CatalogUpdate builder.
func (cu *CatalogUpdate) Where(ps ...predicate.Catalog) *CatalogUpdate {
	cu.mutation.Where(ps...)
	return cu
}

// SetName sets the "name" field.
func (cu *CatalogUpdate) SetName(s string) *CatalogUpdate {
	cu.mutation.SetName(s)
	return cu
}

// SetDescription sets the "description" field.
func (cu *CatalogUpdate) SetDescription(s string) *CatalogUpdate {
	cu.mutation.SetDescription(s)
	return cu
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (cu *CatalogUpdate) SetNillableDescription(s *string) *CatalogUpdate {
	if s != nil {
		cu.SetDescription(*s)
	}
	return cu
}

// ClearDescription clears the value of the "description" field.
func (cu *CatalogUpdate) ClearDescription() *CatalogUpdate {
	cu.mutation.ClearDescription()
	return cu
}

// SetLabels sets the "labels" field.
func (cu *CatalogUpdate) SetLabels(m map[string]string) *CatalogUpdate {
	cu.mutation.SetLabels(m)
	return cu
}

// ClearLabels clears the value of the "labels" field.
func (cu *CatalogUpdate) ClearLabels() *CatalogUpdate {
	cu.mutation.ClearLabels()
	return cu
}

// SetAnnotations sets the "annotations" field.
func (cu *CatalogUpdate) SetAnnotations(m map[string]string) *CatalogUpdate {
	cu.mutation.SetAnnotations(m)
	return cu
}

// ClearAnnotations clears the value of the "annotations" field.
func (cu *CatalogUpdate) ClearAnnotations() *CatalogUpdate {
	cu.mutation.ClearAnnotations()
	return cu
}

// SetUpdateTime sets the "update_time" field.
func (cu *CatalogUpdate) SetUpdateTime(t time.Time) *CatalogUpdate {
	cu.mutation.SetUpdateTime(t)
	return cu
}

// SetStatus sets the "status" field.
func (cu *CatalogUpdate) SetStatus(s status.Status) *CatalogUpdate {
	cu.mutation.SetStatus(s)
	return cu
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (cu *CatalogUpdate) SetNillableStatus(s *status.Status) *CatalogUpdate {
	if s != nil {
		cu.SetStatus(*s)
	}
	return cu
}

// ClearStatus clears the value of the "status" field.
func (cu *CatalogUpdate) ClearStatus() *CatalogUpdate {
	cu.mutation.ClearStatus()
	return cu
}

// SetSource sets the "source" field.
func (cu *CatalogUpdate) SetSource(s string) *CatalogUpdate {
	cu.mutation.SetSource(s)
	return cu
}

// SetSync sets the "sync" field.
func (cu *CatalogUpdate) SetSync(ts *types.CatalogSync) *CatalogUpdate {
	cu.mutation.SetSync(ts)
	return cu
}

// ClearSync clears the value of the "sync" field.
func (cu *CatalogUpdate) ClearSync() *CatalogUpdate {
	cu.mutation.ClearSync()
	return cu
}

// AddTemplateIDs adds the "templates" edge to the Template entity by IDs.
func (cu *CatalogUpdate) AddTemplateIDs(ids ...object.ID) *CatalogUpdate {
	cu.mutation.AddTemplateIDs(ids...)
	return cu
}

// AddTemplates adds the "templates" edges to the Template entity.
func (cu *CatalogUpdate) AddTemplates(t ...*Template) *CatalogUpdate {
	ids := make([]object.ID, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return cu.AddTemplateIDs(ids...)
}

// Mutation returns the CatalogMutation object of the builder.
func (cu *CatalogUpdate) Mutation() *CatalogMutation {
	return cu.mutation
}

// ClearTemplates clears all "templates" edges to the Template entity.
func (cu *CatalogUpdate) ClearTemplates() *CatalogUpdate {
	cu.mutation.ClearTemplates()
	return cu
}

// RemoveTemplateIDs removes the "templates" edge to Template entities by IDs.
func (cu *CatalogUpdate) RemoveTemplateIDs(ids ...object.ID) *CatalogUpdate {
	cu.mutation.RemoveTemplateIDs(ids...)
	return cu
}

// RemoveTemplates removes "templates" edges to Template entities.
func (cu *CatalogUpdate) RemoveTemplates(t ...*Template) *CatalogUpdate {
	ids := make([]object.ID, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return cu.RemoveTemplateIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (cu *CatalogUpdate) Save(ctx context.Context) (int, error) {
	if err := cu.defaults(); err != nil {
		return 0, err
	}
	return withHooks(ctx, cu.sqlSave, cu.mutation, cu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (cu *CatalogUpdate) SaveX(ctx context.Context) int {
	affected, err := cu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (cu *CatalogUpdate) Exec(ctx context.Context) error {
	_, err := cu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cu *CatalogUpdate) ExecX(ctx context.Context) {
	if err := cu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (cu *CatalogUpdate) defaults() error {
	if _, ok := cu.mutation.UpdateTime(); !ok {
		if catalog.UpdateDefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized catalog.UpdateDefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := catalog.UpdateDefaultUpdateTime()
		cu.mutation.SetUpdateTime(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (cu *CatalogUpdate) check() error {
	if v, ok := cu.mutation.Name(); ok {
		if err := catalog.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`model: validator failed for field "Catalog.name": %w`, err)}
		}
	}
	if v, ok := cu.mutation.Source(); ok {
		if err := catalog.SourceValidator(v); err != nil {
			return &ValidationError{Name: "source", err: fmt.Errorf(`model: validator failed for field "Catalog.source": %w`, err)}
		}
	}
	return nil
}

// Set is different from other Set* methods,
// it sets the value by judging the definition of each field within the entire object.
//
// For default fields, Set calls if the value is not zero.
//
// For no default but required fields, Set calls directly.
//
// For no default but optional fields, Set calls if the value is not zero,
// or clears if the value is zero.
//
// For example:
//
//	## Without Default
//
//	### Required
//
//	db.SetX(obj.X)
//
//	### Optional or Default
//
//	if _is_zero_value_(obj.X) {
//	   db.SetX(obj.X)
//	} else {
//	   db.ClearX()
//	}
//
//	## With Default
//
//	if _is_zero_value_(obj.X) {
//	   db.SetX(obj.X)
//	}
func (cu *CatalogUpdate) Set(obj *Catalog) *CatalogUpdate {
	// Without Default.
	cu.SetName(obj.Name)
	if obj.Description != "" {
		cu.SetDescription(obj.Description)
	} else {
		cu.ClearDescription()
	}
	if !reflect.ValueOf(obj.Labels).IsZero() {
		cu.SetLabels(obj.Labels)
	} else {
		cu.ClearLabels()
	}
	if !reflect.ValueOf(obj.Annotations).IsZero() {
		cu.SetAnnotations(obj.Annotations)
	}
	if !reflect.ValueOf(obj.Status).IsZero() {
		cu.SetStatus(obj.Status)
	}
	cu.SetSource(obj.Source)
	if !reflect.ValueOf(obj.Sync).IsZero() {
		cu.SetSync(obj.Sync)
	}

	// With Default.
	if obj.UpdateTime != nil {
		cu.SetUpdateTime(*obj.UpdateTime)
	}

	// Record the given object.
	cu.object = obj

	return cu
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (cu *CatalogUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *CatalogUpdate {
	cu.modifiers = append(cu.modifiers, modifiers...)
	return cu
}

func (cu *CatalogUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := cu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(catalog.Table, catalog.Columns, sqlgraph.NewFieldSpec(catalog.FieldID, field.TypeString))
	if ps := cu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := cu.mutation.Name(); ok {
		_spec.SetField(catalog.FieldName, field.TypeString, value)
	}
	if value, ok := cu.mutation.Description(); ok {
		_spec.SetField(catalog.FieldDescription, field.TypeString, value)
	}
	if cu.mutation.DescriptionCleared() {
		_spec.ClearField(catalog.FieldDescription, field.TypeString)
	}
	if value, ok := cu.mutation.Labels(); ok {
		_spec.SetField(catalog.FieldLabels, field.TypeJSON, value)
	}
	if cu.mutation.LabelsCleared() {
		_spec.ClearField(catalog.FieldLabels, field.TypeJSON)
	}
	if value, ok := cu.mutation.Annotations(); ok {
		_spec.SetField(catalog.FieldAnnotations, field.TypeJSON, value)
	}
	if cu.mutation.AnnotationsCleared() {
		_spec.ClearField(catalog.FieldAnnotations, field.TypeJSON)
	}
	if value, ok := cu.mutation.UpdateTime(); ok {
		_spec.SetField(catalog.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := cu.mutation.Status(); ok {
		_spec.SetField(catalog.FieldStatus, field.TypeJSON, value)
	}
	if cu.mutation.StatusCleared() {
		_spec.ClearField(catalog.FieldStatus, field.TypeJSON)
	}
	if value, ok := cu.mutation.Source(); ok {
		_spec.SetField(catalog.FieldSource, field.TypeString, value)
	}
	if value, ok := cu.mutation.Sync(); ok {
		_spec.SetField(catalog.FieldSync, field.TypeJSON, value)
	}
	if cu.mutation.SyncCleared() {
		_spec.ClearField(catalog.FieldSync, field.TypeJSON)
	}
	if cu.mutation.TemplatesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   catalog.TemplatesTable,
			Columns: []string{catalog.TemplatesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(template.FieldID, field.TypeString),
			},
		}
		edge.Schema = cu.schemaConfig.Template
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cu.mutation.RemovedTemplatesIDs(); len(nodes) > 0 && !cu.mutation.TemplatesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   catalog.TemplatesTable,
			Columns: []string{catalog.TemplatesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(template.FieldID, field.TypeString),
			},
		}
		edge.Schema = cu.schemaConfig.Template
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cu.mutation.TemplatesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   catalog.TemplatesTable,
			Columns: []string{catalog.TemplatesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(template.FieldID, field.TypeString),
			},
		}
		edge.Schema = cu.schemaConfig.Template
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.Node.Schema = cu.schemaConfig.Catalog
	ctx = internal.NewSchemaConfigContext(ctx, cu.schemaConfig)
	_spec.AddModifiers(cu.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, cu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{catalog.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	cu.mutation.done = true
	return n, nil
}

// CatalogUpdateOne is the builder for updating a single Catalog entity.
type CatalogUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *CatalogMutation
	modifiers []func(*sql.UpdateBuilder)
	object    *Catalog
}

// SetName sets the "name" field.
func (cuo *CatalogUpdateOne) SetName(s string) *CatalogUpdateOne {
	cuo.mutation.SetName(s)
	return cuo
}

// SetDescription sets the "description" field.
func (cuo *CatalogUpdateOne) SetDescription(s string) *CatalogUpdateOne {
	cuo.mutation.SetDescription(s)
	return cuo
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (cuo *CatalogUpdateOne) SetNillableDescription(s *string) *CatalogUpdateOne {
	if s != nil {
		cuo.SetDescription(*s)
	}
	return cuo
}

// ClearDescription clears the value of the "description" field.
func (cuo *CatalogUpdateOne) ClearDescription() *CatalogUpdateOne {
	cuo.mutation.ClearDescription()
	return cuo
}

// SetLabels sets the "labels" field.
func (cuo *CatalogUpdateOne) SetLabels(m map[string]string) *CatalogUpdateOne {
	cuo.mutation.SetLabels(m)
	return cuo
}

// ClearLabels clears the value of the "labels" field.
func (cuo *CatalogUpdateOne) ClearLabels() *CatalogUpdateOne {
	cuo.mutation.ClearLabels()
	return cuo
}

// SetAnnotations sets the "annotations" field.
func (cuo *CatalogUpdateOne) SetAnnotations(m map[string]string) *CatalogUpdateOne {
	cuo.mutation.SetAnnotations(m)
	return cuo
}

// ClearAnnotations clears the value of the "annotations" field.
func (cuo *CatalogUpdateOne) ClearAnnotations() *CatalogUpdateOne {
	cuo.mutation.ClearAnnotations()
	return cuo
}

// SetUpdateTime sets the "update_time" field.
func (cuo *CatalogUpdateOne) SetUpdateTime(t time.Time) *CatalogUpdateOne {
	cuo.mutation.SetUpdateTime(t)
	return cuo
}

// SetStatus sets the "status" field.
func (cuo *CatalogUpdateOne) SetStatus(s status.Status) *CatalogUpdateOne {
	cuo.mutation.SetStatus(s)
	return cuo
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (cuo *CatalogUpdateOne) SetNillableStatus(s *status.Status) *CatalogUpdateOne {
	if s != nil {
		cuo.SetStatus(*s)
	}
	return cuo
}

// ClearStatus clears the value of the "status" field.
func (cuo *CatalogUpdateOne) ClearStatus() *CatalogUpdateOne {
	cuo.mutation.ClearStatus()
	return cuo
}

// SetSource sets the "source" field.
func (cuo *CatalogUpdateOne) SetSource(s string) *CatalogUpdateOne {
	cuo.mutation.SetSource(s)
	return cuo
}

// SetSync sets the "sync" field.
func (cuo *CatalogUpdateOne) SetSync(ts *types.CatalogSync) *CatalogUpdateOne {
	cuo.mutation.SetSync(ts)
	return cuo
}

// ClearSync clears the value of the "sync" field.
func (cuo *CatalogUpdateOne) ClearSync() *CatalogUpdateOne {
	cuo.mutation.ClearSync()
	return cuo
}

// AddTemplateIDs adds the "templates" edge to the Template entity by IDs.
func (cuo *CatalogUpdateOne) AddTemplateIDs(ids ...object.ID) *CatalogUpdateOne {
	cuo.mutation.AddTemplateIDs(ids...)
	return cuo
}

// AddTemplates adds the "templates" edges to the Template entity.
func (cuo *CatalogUpdateOne) AddTemplates(t ...*Template) *CatalogUpdateOne {
	ids := make([]object.ID, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return cuo.AddTemplateIDs(ids...)
}

// Mutation returns the CatalogMutation object of the builder.
func (cuo *CatalogUpdateOne) Mutation() *CatalogMutation {
	return cuo.mutation
}

// ClearTemplates clears all "templates" edges to the Template entity.
func (cuo *CatalogUpdateOne) ClearTemplates() *CatalogUpdateOne {
	cuo.mutation.ClearTemplates()
	return cuo
}

// RemoveTemplateIDs removes the "templates" edge to Template entities by IDs.
func (cuo *CatalogUpdateOne) RemoveTemplateIDs(ids ...object.ID) *CatalogUpdateOne {
	cuo.mutation.RemoveTemplateIDs(ids...)
	return cuo
}

// RemoveTemplates removes "templates" edges to Template entities.
func (cuo *CatalogUpdateOne) RemoveTemplates(t ...*Template) *CatalogUpdateOne {
	ids := make([]object.ID, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return cuo.RemoveTemplateIDs(ids...)
}

// Where appends a list predicates to the CatalogUpdate builder.
func (cuo *CatalogUpdateOne) Where(ps ...predicate.Catalog) *CatalogUpdateOne {
	cuo.mutation.Where(ps...)
	return cuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (cuo *CatalogUpdateOne) Select(field string, fields ...string) *CatalogUpdateOne {
	cuo.fields = append([]string{field}, fields...)
	return cuo
}

// Save executes the query and returns the updated Catalog entity.
func (cuo *CatalogUpdateOne) Save(ctx context.Context) (*Catalog, error) {
	if err := cuo.defaults(); err != nil {
		return nil, err
	}
	return withHooks(ctx, cuo.sqlSave, cuo.mutation, cuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (cuo *CatalogUpdateOne) SaveX(ctx context.Context) *Catalog {
	node, err := cuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (cuo *CatalogUpdateOne) Exec(ctx context.Context) error {
	_, err := cuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cuo *CatalogUpdateOne) ExecX(ctx context.Context) {
	if err := cuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (cuo *CatalogUpdateOne) defaults() error {
	if _, ok := cuo.mutation.UpdateTime(); !ok {
		if catalog.UpdateDefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized catalog.UpdateDefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := catalog.UpdateDefaultUpdateTime()
		cuo.mutation.SetUpdateTime(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (cuo *CatalogUpdateOne) check() error {
	if v, ok := cuo.mutation.Name(); ok {
		if err := catalog.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`model: validator failed for field "Catalog.name": %w`, err)}
		}
	}
	if v, ok := cuo.mutation.Source(); ok {
		if err := catalog.SourceValidator(v); err != nil {
			return &ValidationError{Name: "source", err: fmt.Errorf(`model: validator failed for field "Catalog.source": %w`, err)}
		}
	}
	return nil
}

// Set is different from other Set* methods,
// it sets the value by judging the definition of each field within the entire object.
//
// For default fields, Set calls if the value changes from the original.
//
// For no default but required fields, Set calls if the value changes from the original.
//
// For no default but optional fields, Set calls if the value changes from the original,
// or clears if changes to zero.
//
// For example:
//
//	## Without Default
//
//	### Required
//
//	db.SetX(obj.X)
//
//	### Optional or Default
//
//	if _is_zero_value_(obj.X) {
//	   if _is_not_equal_(db.X, obj.X) {
//	      db.SetX(obj.X)
//	   }
//	} else {
//	   db.ClearX()
//	}
//
//	## With Default
//
//	if _is_zero_value_(obj.X) && _is_not_equal_(db.X, obj.X) {
//	   db.SetX(obj.X)
//	}
func (cuo *CatalogUpdateOne) Set(obj *Catalog) *CatalogUpdateOne {
	h := func(n ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			mt := m.(*CatalogMutation)
			db, err := mt.Client().Catalog.Get(ctx, *mt.id)
			if err != nil {
				return nil, fmt.Errorf("failed getting Catalog with id: %v", *mt.id)
			}

			// Without Default.
			if db.Name != obj.Name {
				cuo.SetName(obj.Name)
			}
			if obj.Description != "" {
				if db.Description != obj.Description {
					cuo.SetDescription(obj.Description)
				}
			} else {
				cuo.ClearDescription()
			}
			if !reflect.ValueOf(obj.Labels).IsZero() {
				if !reflect.DeepEqual(db.Labels, obj.Labels) {
					cuo.SetLabels(obj.Labels)
				}
			} else {
				cuo.ClearLabels()
			}
			if !reflect.ValueOf(obj.Annotations).IsZero() {
				if !reflect.DeepEqual(db.Annotations, obj.Annotations) {
					cuo.SetAnnotations(obj.Annotations)
				}
			}
			if !reflect.ValueOf(obj.Status).IsZero() {
				if !db.Status.Equal(obj.Status) {
					cuo.SetStatus(obj.Status)
				}
			}
			if db.Source != obj.Source {
				cuo.SetSource(obj.Source)
			}
			if !reflect.ValueOf(obj.Sync).IsZero() {
				if !reflect.DeepEqual(db.Sync, obj.Sync) {
					cuo.SetSync(obj.Sync)
				}
			}

			// With Default.
			if (obj.UpdateTime != nil) && (!reflect.DeepEqual(db.UpdateTime, obj.UpdateTime)) {
				cuo.SetUpdateTime(*obj.UpdateTime)
			}

			// Record the given object.
			cuo.object = obj

			return n.Mutate(ctx, m)
		})
	}

	cuo.hooks = append(cuo.hooks, h)

	return cuo
}

// getClientSet returns the ClientSet for the given builder.
func (cuo *CatalogUpdateOne) getClientSet() (mc ClientSet) {
	if _, ok := cuo.config.driver.(*txDriver); ok {
		tx := &Tx{config: cuo.config}
		tx.init()
		mc = tx
	} else {
		cli := &Client{config: cuo.config}
		cli.init()
		mc = cli
	}
	return mc
}

// SaveE calls the given function after updated the Catalog entity,
// which is always good for cascading update operations.
func (cuo *CatalogUpdateOne) SaveE(ctx context.Context, cbs ...func(ctx context.Context, mc ClientSet, updated *Catalog) error) (*Catalog, error) {
	obj, err := cuo.Save(ctx)
	if err != nil &&
		(cuo.object == nil || !errors.Is(err, stdsql.ErrNoRows)) {
		return nil, err
	}

	if len(cbs) == 0 {
		return obj, err
	}

	mc := cuo.getClientSet()

	if obj == nil {
		obj = cuo.object
	} else if x := cuo.object; x != nil {
		if _, set := cuo.mutation.Field(catalog.FieldName); set {
			obj.Name = x.Name
		}
		if _, set := cuo.mutation.Field(catalog.FieldDescription); set {
			obj.Description = x.Description
		}
		if _, set := cuo.mutation.Field(catalog.FieldLabels); set {
			obj.Labels = x.Labels
		}
		if _, set := cuo.mutation.Field(catalog.FieldAnnotations); set {
			obj.Annotations = x.Annotations
		}
		if _, set := cuo.mutation.Field(catalog.FieldStatus); set {
			obj.Status = x.Status
		}
		if _, set := cuo.mutation.Field(catalog.FieldSource); set {
			obj.Source = x.Source
		}
		if _, set := cuo.mutation.Field(catalog.FieldSync); set {
			obj.Sync = x.Sync
		}
		obj.Edges = x.Edges
	}

	for i := range cbs {
		if err = cbs[i](ctx, mc, obj); err != nil {
			return nil, err
		}
	}

	return obj, nil
}

// SaveEX is like SaveE, but panics if an error occurs.
func (cuo *CatalogUpdateOne) SaveEX(ctx context.Context, cbs ...func(ctx context.Context, mc ClientSet, updated *Catalog) error) *Catalog {
	obj, err := cuo.SaveE(ctx, cbs...)
	if err != nil {
		panic(err)
	}
	return obj
}

// ExecE calls the given function after executed the query,
// which is always good for cascading update operations.
func (cuo *CatalogUpdateOne) ExecE(ctx context.Context, cbs ...func(ctx context.Context, mc ClientSet, updated *Catalog) error) error {
	_, err := cuo.SaveE(ctx, cbs...)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cuo *CatalogUpdateOne) ExecEX(ctx context.Context, cbs ...func(ctx context.Context, mc ClientSet, updated *Catalog) error) {
	if err := cuo.ExecE(ctx, cbs...); err != nil {
		panic(err)
	}
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (cuo *CatalogUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *CatalogUpdateOne {
	cuo.modifiers = append(cuo.modifiers, modifiers...)
	return cuo
}

func (cuo *CatalogUpdateOne) sqlSave(ctx context.Context) (_node *Catalog, err error) {
	if err := cuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(catalog.Table, catalog.Columns, sqlgraph.NewFieldSpec(catalog.FieldID, field.TypeString))
	id, ok := cuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`model: missing "Catalog.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := cuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, catalog.FieldID)
		for _, f := range fields {
			if !catalog.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("model: invalid field %q for query", f)}
			}
			if f != catalog.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := cuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := cuo.mutation.Name(); ok {
		_spec.SetField(catalog.FieldName, field.TypeString, value)
	}
	if value, ok := cuo.mutation.Description(); ok {
		_spec.SetField(catalog.FieldDescription, field.TypeString, value)
	}
	if cuo.mutation.DescriptionCleared() {
		_spec.ClearField(catalog.FieldDescription, field.TypeString)
	}
	if value, ok := cuo.mutation.Labels(); ok {
		_spec.SetField(catalog.FieldLabels, field.TypeJSON, value)
	}
	if cuo.mutation.LabelsCleared() {
		_spec.ClearField(catalog.FieldLabels, field.TypeJSON)
	}
	if value, ok := cuo.mutation.Annotations(); ok {
		_spec.SetField(catalog.FieldAnnotations, field.TypeJSON, value)
	}
	if cuo.mutation.AnnotationsCleared() {
		_spec.ClearField(catalog.FieldAnnotations, field.TypeJSON)
	}
	if value, ok := cuo.mutation.UpdateTime(); ok {
		_spec.SetField(catalog.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := cuo.mutation.Status(); ok {
		_spec.SetField(catalog.FieldStatus, field.TypeJSON, value)
	}
	if cuo.mutation.StatusCleared() {
		_spec.ClearField(catalog.FieldStatus, field.TypeJSON)
	}
	if value, ok := cuo.mutation.Source(); ok {
		_spec.SetField(catalog.FieldSource, field.TypeString, value)
	}
	if value, ok := cuo.mutation.Sync(); ok {
		_spec.SetField(catalog.FieldSync, field.TypeJSON, value)
	}
	if cuo.mutation.SyncCleared() {
		_spec.ClearField(catalog.FieldSync, field.TypeJSON)
	}
	if cuo.mutation.TemplatesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   catalog.TemplatesTable,
			Columns: []string{catalog.TemplatesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(template.FieldID, field.TypeString),
			},
		}
		edge.Schema = cuo.schemaConfig.Template
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cuo.mutation.RemovedTemplatesIDs(); len(nodes) > 0 && !cuo.mutation.TemplatesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   catalog.TemplatesTable,
			Columns: []string{catalog.TemplatesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(template.FieldID, field.TypeString),
			},
		}
		edge.Schema = cuo.schemaConfig.Template
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cuo.mutation.TemplatesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   catalog.TemplatesTable,
			Columns: []string{catalog.TemplatesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(template.FieldID, field.TypeString),
			},
		}
		edge.Schema = cuo.schemaConfig.Template
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.Node.Schema = cuo.schemaConfig.Catalog
	ctx = internal.NewSchemaConfigContext(ctx, cuo.schemaConfig)
	_spec.AddModifiers(cuo.modifiers...)
	_node = &Catalog{config: cuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, cuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{catalog.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	cuo.mutation.done = true
	return _node, nil
}
