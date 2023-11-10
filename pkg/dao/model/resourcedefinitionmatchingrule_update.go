// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "walrus". DO NOT EDIT.

package model

import (
	"context"
	stdsql "database/sql"
	"errors"
	"fmt"
	"reflect"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/walrus/pkg/dao/model/internal"
	"github.com/seal-io/walrus/pkg/dao/model/predicate"
	"github.com/seal-io/walrus/pkg/dao/model/resourcedefinitionmatchingrule"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/property"
)

// ResourceDefinitionMatchingRuleUpdate is the builder for updating ResourceDefinitionMatchingRule entities.
type ResourceDefinitionMatchingRuleUpdate struct {
	config
	hooks     []Hook
	mutation  *ResourceDefinitionMatchingRuleMutation
	modifiers []func(*sql.UpdateBuilder)
	object    *ResourceDefinitionMatchingRule
}

// Where appends a list predicates to the ResourceDefinitionMatchingRuleUpdate builder.
func (rdmru *ResourceDefinitionMatchingRuleUpdate) Where(ps ...predicate.ResourceDefinitionMatchingRule) *ResourceDefinitionMatchingRuleUpdate {
	rdmru.mutation.Where(ps...)
	return rdmru
}

// SetSelector sets the "selector" field.
func (rdmru *ResourceDefinitionMatchingRuleUpdate) SetSelector(t types.Selector) *ResourceDefinitionMatchingRuleUpdate {
	rdmru.mutation.SetSelector(t)
	return rdmru
}

// SetAttributes sets the "attributes" field.
func (rdmru *ResourceDefinitionMatchingRuleUpdate) SetAttributes(pr property.Values) *ResourceDefinitionMatchingRuleUpdate {
	rdmru.mutation.SetAttributes(pr)
	return rdmru
}

// ClearAttributes clears the value of the "attributes" field.
func (rdmru *ResourceDefinitionMatchingRuleUpdate) ClearAttributes() *ResourceDefinitionMatchingRuleUpdate {
	rdmru.mutation.ClearAttributes()
	return rdmru
}

// Mutation returns the ResourceDefinitionMatchingRuleMutation object of the builder.
func (rdmru *ResourceDefinitionMatchingRuleUpdate) Mutation() *ResourceDefinitionMatchingRuleMutation {
	return rdmru.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (rdmru *ResourceDefinitionMatchingRuleUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, rdmru.sqlSave, rdmru.mutation, rdmru.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (rdmru *ResourceDefinitionMatchingRuleUpdate) SaveX(ctx context.Context) int {
	affected, err := rdmru.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (rdmru *ResourceDefinitionMatchingRuleUpdate) Exec(ctx context.Context) error {
	_, err := rdmru.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (rdmru *ResourceDefinitionMatchingRuleUpdate) ExecX(ctx context.Context) {
	if err := rdmru.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (rdmru *ResourceDefinitionMatchingRuleUpdate) check() error {
	if _, ok := rdmru.mutation.ResourceDefinitionID(); rdmru.mutation.ResourceDefinitionCleared() && !ok {
		return errors.New(`model: clearing a required unique edge "ResourceDefinitionMatchingRule.resource_definition"`)
	}
	if _, ok := rdmru.mutation.TemplateID(); rdmru.mutation.TemplateCleared() && !ok {
		return errors.New(`model: clearing a required unique edge "ResourceDefinitionMatchingRule.template"`)
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
func (rdmru *ResourceDefinitionMatchingRuleUpdate) Set(obj *ResourceDefinitionMatchingRule) *ResourceDefinitionMatchingRuleUpdate {
	// Without Default.
	rdmru.SetSelector(obj.Selector)
	if !reflect.ValueOf(obj.Attributes).IsZero() {
		rdmru.SetAttributes(obj.Attributes)
	}

	// With Default.

	// Record the given object.
	rdmru.object = obj

	return rdmru
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (rdmru *ResourceDefinitionMatchingRuleUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ResourceDefinitionMatchingRuleUpdate {
	rdmru.modifiers = append(rdmru.modifiers, modifiers...)
	return rdmru
}

func (rdmru *ResourceDefinitionMatchingRuleUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := rdmru.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(resourcedefinitionmatchingrule.Table, resourcedefinitionmatchingrule.Columns, sqlgraph.NewFieldSpec(resourcedefinitionmatchingrule.FieldID, field.TypeString))
	if ps := rdmru.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := rdmru.mutation.Selector(); ok {
		_spec.SetField(resourcedefinitionmatchingrule.FieldSelector, field.TypeJSON, value)
	}
	if value, ok := rdmru.mutation.Attributes(); ok {
		_spec.SetField(resourcedefinitionmatchingrule.FieldAttributes, field.TypeOther, value)
	}
	if rdmru.mutation.AttributesCleared() {
		_spec.ClearField(resourcedefinitionmatchingrule.FieldAttributes, field.TypeOther)
	}
	_spec.Node.Schema = rdmru.schemaConfig.ResourceDefinitionMatchingRule
	ctx = internal.NewSchemaConfigContext(ctx, rdmru.schemaConfig)
	_spec.AddModifiers(rdmru.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, rdmru.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{resourcedefinitionmatchingrule.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	rdmru.mutation.done = true
	return n, nil
}

// ResourceDefinitionMatchingRuleUpdateOne is the builder for updating a single ResourceDefinitionMatchingRule entity.
type ResourceDefinitionMatchingRuleUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *ResourceDefinitionMatchingRuleMutation
	modifiers []func(*sql.UpdateBuilder)
	object    *ResourceDefinitionMatchingRule
}

// SetSelector sets the "selector" field.
func (rdmruo *ResourceDefinitionMatchingRuleUpdateOne) SetSelector(t types.Selector) *ResourceDefinitionMatchingRuleUpdateOne {
	rdmruo.mutation.SetSelector(t)
	return rdmruo
}

// SetAttributes sets the "attributes" field.
func (rdmruo *ResourceDefinitionMatchingRuleUpdateOne) SetAttributes(pr property.Values) *ResourceDefinitionMatchingRuleUpdateOne {
	rdmruo.mutation.SetAttributes(pr)
	return rdmruo
}

// ClearAttributes clears the value of the "attributes" field.
func (rdmruo *ResourceDefinitionMatchingRuleUpdateOne) ClearAttributes() *ResourceDefinitionMatchingRuleUpdateOne {
	rdmruo.mutation.ClearAttributes()
	return rdmruo
}

// Mutation returns the ResourceDefinitionMatchingRuleMutation object of the builder.
func (rdmruo *ResourceDefinitionMatchingRuleUpdateOne) Mutation() *ResourceDefinitionMatchingRuleMutation {
	return rdmruo.mutation
}

// Where appends a list predicates to the ResourceDefinitionMatchingRuleUpdate builder.
func (rdmruo *ResourceDefinitionMatchingRuleUpdateOne) Where(ps ...predicate.ResourceDefinitionMatchingRule) *ResourceDefinitionMatchingRuleUpdateOne {
	rdmruo.mutation.Where(ps...)
	return rdmruo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (rdmruo *ResourceDefinitionMatchingRuleUpdateOne) Select(field string, fields ...string) *ResourceDefinitionMatchingRuleUpdateOne {
	rdmruo.fields = append([]string{field}, fields...)
	return rdmruo
}

// Save executes the query and returns the updated ResourceDefinitionMatchingRule entity.
func (rdmruo *ResourceDefinitionMatchingRuleUpdateOne) Save(ctx context.Context) (*ResourceDefinitionMatchingRule, error) {
	return withHooks(ctx, rdmruo.sqlSave, rdmruo.mutation, rdmruo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (rdmruo *ResourceDefinitionMatchingRuleUpdateOne) SaveX(ctx context.Context) *ResourceDefinitionMatchingRule {
	node, err := rdmruo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (rdmruo *ResourceDefinitionMatchingRuleUpdateOne) Exec(ctx context.Context) error {
	_, err := rdmruo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (rdmruo *ResourceDefinitionMatchingRuleUpdateOne) ExecX(ctx context.Context) {
	if err := rdmruo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (rdmruo *ResourceDefinitionMatchingRuleUpdateOne) check() error {
	if _, ok := rdmruo.mutation.ResourceDefinitionID(); rdmruo.mutation.ResourceDefinitionCleared() && !ok {
		return errors.New(`model: clearing a required unique edge "ResourceDefinitionMatchingRule.resource_definition"`)
	}
	if _, ok := rdmruo.mutation.TemplateID(); rdmruo.mutation.TemplateCleared() && !ok {
		return errors.New(`model: clearing a required unique edge "ResourceDefinitionMatchingRule.template"`)
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
func (rdmruo *ResourceDefinitionMatchingRuleUpdateOne) Set(obj *ResourceDefinitionMatchingRule) *ResourceDefinitionMatchingRuleUpdateOne {
	h := func(n ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			mt := m.(*ResourceDefinitionMatchingRuleMutation)
			db, err := mt.Client().ResourceDefinitionMatchingRule.Get(ctx, *mt.id)
			if err != nil {
				return nil, fmt.Errorf("failed getting ResourceDefinitionMatchingRule with id: %v", *mt.id)
			}

			// Without Default.
			if !reflect.DeepEqual(db.Selector, obj.Selector) {
				rdmruo.SetSelector(obj.Selector)
			}
			if !reflect.ValueOf(obj.Attributes).IsZero() {
				if !reflect.DeepEqual(db.Attributes, obj.Attributes) {
					rdmruo.SetAttributes(obj.Attributes)
				}
			}

			// With Default.

			// Record the given object.
			rdmruo.object = obj

			return n.Mutate(ctx, m)
		})
	}

	rdmruo.hooks = append(rdmruo.hooks, h)

	return rdmruo
}

// getClientSet returns the ClientSet for the given builder.
func (rdmruo *ResourceDefinitionMatchingRuleUpdateOne) getClientSet() (mc ClientSet) {
	if _, ok := rdmruo.config.driver.(*txDriver); ok {
		tx := &Tx{config: rdmruo.config}
		tx.init()
		mc = tx
	} else {
		cli := &Client{config: rdmruo.config}
		cli.init()
		mc = cli
	}
	return mc
}

// SaveE calls the given function after updated the ResourceDefinitionMatchingRule entity,
// which is always good for cascading update operations.
func (rdmruo *ResourceDefinitionMatchingRuleUpdateOne) SaveE(ctx context.Context, cbs ...func(ctx context.Context, mc ClientSet, updated *ResourceDefinitionMatchingRule) error) (*ResourceDefinitionMatchingRule, error) {
	obj, err := rdmruo.Save(ctx)
	if err != nil &&
		(rdmruo.object == nil || !errors.Is(err, stdsql.ErrNoRows)) {
		return nil, err
	}

	if len(cbs) == 0 {
		return obj, err
	}

	mc := rdmruo.getClientSet()

	if obj == nil {
		obj = rdmruo.object
	} else if x := rdmruo.object; x != nil {
		if _, set := rdmruo.mutation.Field(resourcedefinitionmatchingrule.FieldSelector); set {
			obj.Selector = x.Selector
		}
		if _, set := rdmruo.mutation.Field(resourcedefinitionmatchingrule.FieldAttributes); set {
			obj.Attributes = x.Attributes
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
func (rdmruo *ResourceDefinitionMatchingRuleUpdateOne) SaveEX(ctx context.Context, cbs ...func(ctx context.Context, mc ClientSet, updated *ResourceDefinitionMatchingRule) error) *ResourceDefinitionMatchingRule {
	obj, err := rdmruo.SaveE(ctx, cbs...)
	if err != nil {
		panic(err)
	}
	return obj
}

// ExecE calls the given function after executed the query,
// which is always good for cascading update operations.
func (rdmruo *ResourceDefinitionMatchingRuleUpdateOne) ExecE(ctx context.Context, cbs ...func(ctx context.Context, mc ClientSet, updated *ResourceDefinitionMatchingRule) error) error {
	_, err := rdmruo.SaveE(ctx, cbs...)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (rdmruo *ResourceDefinitionMatchingRuleUpdateOne) ExecEX(ctx context.Context, cbs ...func(ctx context.Context, mc ClientSet, updated *ResourceDefinitionMatchingRule) error) {
	if err := rdmruo.ExecE(ctx, cbs...); err != nil {
		panic(err)
	}
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (rdmruo *ResourceDefinitionMatchingRuleUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *ResourceDefinitionMatchingRuleUpdateOne {
	rdmruo.modifiers = append(rdmruo.modifiers, modifiers...)
	return rdmruo
}

func (rdmruo *ResourceDefinitionMatchingRuleUpdateOne) sqlSave(ctx context.Context) (_node *ResourceDefinitionMatchingRule, err error) {
	if err := rdmruo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(resourcedefinitionmatchingrule.Table, resourcedefinitionmatchingrule.Columns, sqlgraph.NewFieldSpec(resourcedefinitionmatchingrule.FieldID, field.TypeString))
	id, ok := rdmruo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`model: missing "ResourceDefinitionMatchingRule.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := rdmruo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, resourcedefinitionmatchingrule.FieldID)
		for _, f := range fields {
			if !resourcedefinitionmatchingrule.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("model: invalid field %q for query", f)}
			}
			if f != resourcedefinitionmatchingrule.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := rdmruo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := rdmruo.mutation.Selector(); ok {
		_spec.SetField(resourcedefinitionmatchingrule.FieldSelector, field.TypeJSON, value)
	}
	if value, ok := rdmruo.mutation.Attributes(); ok {
		_spec.SetField(resourcedefinitionmatchingrule.FieldAttributes, field.TypeOther, value)
	}
	if rdmruo.mutation.AttributesCleared() {
		_spec.ClearField(resourcedefinitionmatchingrule.FieldAttributes, field.TypeOther)
	}
	_spec.Node.Schema = rdmruo.schemaConfig.ResourceDefinitionMatchingRule
	ctx = internal.NewSchemaConfigContext(ctx, rdmruo.schemaConfig)
	_spec.AddModifiers(rdmruo.modifiers...)
	_node = &ResourceDefinitionMatchingRule{config: rdmruo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, rdmruo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{resourcedefinitionmatchingrule.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	rdmruo.mutation.done = true
	return _node, nil
}