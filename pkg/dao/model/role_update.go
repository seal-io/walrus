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
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/dialect/sql/sqljson"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/walrus/pkg/dao/model/internal"
	"github.com/seal-io/walrus/pkg/dao/model/predicate"
	"github.com/seal-io/walrus/pkg/dao/model/role"
	"github.com/seal-io/walrus/pkg/dao/model/subjectrolerelationship"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
)

// RoleUpdate is the builder for updating Role entities.
type RoleUpdate struct {
	config
	hooks     []Hook
	mutation  *RoleMutation
	modifiers []func(*sql.UpdateBuilder)
	object    *Role
}

// Where appends a list predicates to the RoleUpdate builder.
func (ru *RoleUpdate) Where(ps ...predicate.Role) *RoleUpdate {
	ru.mutation.Where(ps...)
	return ru
}

// SetUpdateTime sets the "update_time" field.
func (ru *RoleUpdate) SetUpdateTime(t time.Time) *RoleUpdate {
	ru.mutation.SetUpdateTime(t)
	return ru
}

// SetDescription sets the "description" field.
func (ru *RoleUpdate) SetDescription(s string) *RoleUpdate {
	ru.mutation.SetDescription(s)
	return ru
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (ru *RoleUpdate) SetNillableDescription(s *string) *RoleUpdate {
	if s != nil {
		ru.SetDescription(*s)
	}
	return ru
}

// ClearDescription clears the value of the "description" field.
func (ru *RoleUpdate) ClearDescription() *RoleUpdate {
	ru.mutation.ClearDescription()
	return ru
}

// SetPolicies sets the "policies" field.
func (ru *RoleUpdate) SetPolicies(tp types.RolePolicies) *RoleUpdate {
	ru.mutation.SetPolicies(tp)
	return ru
}

// AppendPolicies appends tp to the "policies" field.
func (ru *RoleUpdate) AppendPolicies(tp types.RolePolicies) *RoleUpdate {
	ru.mutation.AppendPolicies(tp)
	return ru
}

// AddSubjectIDs adds the "subjects" edge to the SubjectRoleRelationship entity by IDs.
func (ru *RoleUpdate) AddSubjectIDs(ids ...object.ID) *RoleUpdate {
	ru.mutation.AddSubjectIDs(ids...)
	return ru
}

// AddSubjects adds the "subjects" edges to the SubjectRoleRelationship entity.
func (ru *RoleUpdate) AddSubjects(s ...*SubjectRoleRelationship) *RoleUpdate {
	ids := make([]object.ID, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return ru.AddSubjectIDs(ids...)
}

// Mutation returns the RoleMutation object of the builder.
func (ru *RoleUpdate) Mutation() *RoleMutation {
	return ru.mutation
}

// ClearSubjects clears all "subjects" edges to the SubjectRoleRelationship entity.
func (ru *RoleUpdate) ClearSubjects() *RoleUpdate {
	ru.mutation.ClearSubjects()
	return ru
}

// RemoveSubjectIDs removes the "subjects" edge to SubjectRoleRelationship entities by IDs.
func (ru *RoleUpdate) RemoveSubjectIDs(ids ...object.ID) *RoleUpdate {
	ru.mutation.RemoveSubjectIDs(ids...)
	return ru
}

// RemoveSubjects removes "subjects" edges to SubjectRoleRelationship entities.
func (ru *RoleUpdate) RemoveSubjects(s ...*SubjectRoleRelationship) *RoleUpdate {
	ids := make([]object.ID, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return ru.RemoveSubjectIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (ru *RoleUpdate) Save(ctx context.Context) (int, error) {
	if err := ru.defaults(); err != nil {
		return 0, err
	}
	return withHooks(ctx, ru.sqlSave, ru.mutation, ru.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ru *RoleUpdate) SaveX(ctx context.Context) int {
	affected, err := ru.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (ru *RoleUpdate) Exec(ctx context.Context) error {
	_, err := ru.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ru *RoleUpdate) ExecX(ctx context.Context) {
	if err := ru.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ru *RoleUpdate) defaults() error {
	if _, ok := ru.mutation.UpdateTime(); !ok {
		if role.UpdateDefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized role.UpdateDefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := role.UpdateDefaultUpdateTime()
		ru.mutation.SetUpdateTime(v)
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
func (ru *RoleUpdate) Set(obj *Role) *RoleUpdate {
	// Without Default.
	if obj.Description != "" {
		ru.SetDescription(obj.Description)
	} else {
		ru.ClearDescription()
	}
	ru.SetPolicies(obj.Policies)

	// With Default.
	if obj.UpdateTime != nil {
		ru.SetUpdateTime(*obj.UpdateTime)
	}

	// Record the given object.
	ru.object = obj

	return ru
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (ru *RoleUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *RoleUpdate {
	ru.modifiers = append(ru.modifiers, modifiers...)
	return ru
}

func (ru *RoleUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(role.Table, role.Columns, sqlgraph.NewFieldSpec(role.FieldID, field.TypeString))
	if ps := ru.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ru.mutation.UpdateTime(); ok {
		_spec.SetField(role.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := ru.mutation.Description(); ok {
		_spec.SetField(role.FieldDescription, field.TypeString, value)
	}
	if ru.mutation.DescriptionCleared() {
		_spec.ClearField(role.FieldDescription, field.TypeString)
	}
	if value, ok := ru.mutation.Policies(); ok {
		_spec.SetField(role.FieldPolicies, field.TypeJSON, value)
	}
	if value, ok := ru.mutation.AppendedPolicies(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, role.FieldPolicies, value)
		})
	}
	if ru.mutation.SubjectsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   role.SubjectsTable,
			Columns: []string{role.SubjectsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(subjectrolerelationship.FieldID, field.TypeString),
			},
		}
		edge.Schema = ru.schemaConfig.SubjectRoleRelationship
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ru.mutation.RemovedSubjectsIDs(); len(nodes) > 0 && !ru.mutation.SubjectsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   role.SubjectsTable,
			Columns: []string{role.SubjectsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(subjectrolerelationship.FieldID, field.TypeString),
			},
		}
		edge.Schema = ru.schemaConfig.SubjectRoleRelationship
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ru.mutation.SubjectsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   role.SubjectsTable,
			Columns: []string{role.SubjectsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(subjectrolerelationship.FieldID, field.TypeString),
			},
		}
		edge.Schema = ru.schemaConfig.SubjectRoleRelationship
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.Node.Schema = ru.schemaConfig.Role
	ctx = internal.NewSchemaConfigContext(ctx, ru.schemaConfig)
	_spec.AddModifiers(ru.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, ru.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{role.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	ru.mutation.done = true
	return n, nil
}

// RoleUpdateOne is the builder for updating a single Role entity.
type RoleUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *RoleMutation
	modifiers []func(*sql.UpdateBuilder)
	object    *Role
}

// SetUpdateTime sets the "update_time" field.
func (ruo *RoleUpdateOne) SetUpdateTime(t time.Time) *RoleUpdateOne {
	ruo.mutation.SetUpdateTime(t)
	return ruo
}

// SetDescription sets the "description" field.
func (ruo *RoleUpdateOne) SetDescription(s string) *RoleUpdateOne {
	ruo.mutation.SetDescription(s)
	return ruo
}

// SetNillableDescription sets the "description" field if the given value is not nil.
func (ruo *RoleUpdateOne) SetNillableDescription(s *string) *RoleUpdateOne {
	if s != nil {
		ruo.SetDescription(*s)
	}
	return ruo
}

// ClearDescription clears the value of the "description" field.
func (ruo *RoleUpdateOne) ClearDescription() *RoleUpdateOne {
	ruo.mutation.ClearDescription()
	return ruo
}

// SetPolicies sets the "policies" field.
func (ruo *RoleUpdateOne) SetPolicies(tp types.RolePolicies) *RoleUpdateOne {
	ruo.mutation.SetPolicies(tp)
	return ruo
}

// AppendPolicies appends tp to the "policies" field.
func (ruo *RoleUpdateOne) AppendPolicies(tp types.RolePolicies) *RoleUpdateOne {
	ruo.mutation.AppendPolicies(tp)
	return ruo
}

// AddSubjectIDs adds the "subjects" edge to the SubjectRoleRelationship entity by IDs.
func (ruo *RoleUpdateOne) AddSubjectIDs(ids ...object.ID) *RoleUpdateOne {
	ruo.mutation.AddSubjectIDs(ids...)
	return ruo
}

// AddSubjects adds the "subjects" edges to the SubjectRoleRelationship entity.
func (ruo *RoleUpdateOne) AddSubjects(s ...*SubjectRoleRelationship) *RoleUpdateOne {
	ids := make([]object.ID, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return ruo.AddSubjectIDs(ids...)
}

// Mutation returns the RoleMutation object of the builder.
func (ruo *RoleUpdateOne) Mutation() *RoleMutation {
	return ruo.mutation
}

// ClearSubjects clears all "subjects" edges to the SubjectRoleRelationship entity.
func (ruo *RoleUpdateOne) ClearSubjects() *RoleUpdateOne {
	ruo.mutation.ClearSubjects()
	return ruo
}

// RemoveSubjectIDs removes the "subjects" edge to SubjectRoleRelationship entities by IDs.
func (ruo *RoleUpdateOne) RemoveSubjectIDs(ids ...object.ID) *RoleUpdateOne {
	ruo.mutation.RemoveSubjectIDs(ids...)
	return ruo
}

// RemoveSubjects removes "subjects" edges to SubjectRoleRelationship entities.
func (ruo *RoleUpdateOne) RemoveSubjects(s ...*SubjectRoleRelationship) *RoleUpdateOne {
	ids := make([]object.ID, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return ruo.RemoveSubjectIDs(ids...)
}

// Where appends a list predicates to the RoleUpdate builder.
func (ruo *RoleUpdateOne) Where(ps ...predicate.Role) *RoleUpdateOne {
	ruo.mutation.Where(ps...)
	return ruo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (ruo *RoleUpdateOne) Select(field string, fields ...string) *RoleUpdateOne {
	ruo.fields = append([]string{field}, fields...)
	return ruo
}

// Save executes the query and returns the updated Role entity.
func (ruo *RoleUpdateOne) Save(ctx context.Context) (*Role, error) {
	if err := ruo.defaults(); err != nil {
		return nil, err
	}
	return withHooks(ctx, ruo.sqlSave, ruo.mutation, ruo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ruo *RoleUpdateOne) SaveX(ctx context.Context) *Role {
	node, err := ruo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (ruo *RoleUpdateOne) Exec(ctx context.Context) error {
	_, err := ruo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ruo *RoleUpdateOne) ExecX(ctx context.Context) {
	if err := ruo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ruo *RoleUpdateOne) defaults() error {
	if _, ok := ruo.mutation.UpdateTime(); !ok {
		if role.UpdateDefaultUpdateTime == nil {
			return fmt.Errorf("model: uninitialized role.UpdateDefaultUpdateTime (forgotten import model/runtime?)")
		}
		v := role.UpdateDefaultUpdateTime()
		ruo.mutation.SetUpdateTime(v)
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
func (ruo *RoleUpdateOne) Set(obj *Role) *RoleUpdateOne {
	h := func(n ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			mt := m.(*RoleMutation)
			db, err := mt.Client().Role.Get(ctx, *mt.id)
			if err != nil {
				return nil, fmt.Errorf("failed getting Role with id: %v", *mt.id)
			}

			// Without Default.
			if obj.Description != "" {
				if db.Description != obj.Description {
					ruo.SetDescription(obj.Description)
				}
			} else {
				ruo.ClearDescription()
			}
			if !reflect.DeepEqual(db.Policies, obj.Policies) {
				ruo.SetPolicies(obj.Policies)
			}

			// With Default.
			if (obj.UpdateTime != nil) && (!reflect.DeepEqual(db.UpdateTime, obj.UpdateTime)) {
				ruo.SetUpdateTime(*obj.UpdateTime)
			}

			// Record the given object.
			ruo.object = obj

			return n.Mutate(ctx, m)
		})
	}

	ruo.hooks = append(ruo.hooks, h)

	return ruo
}

// getClientSet returns the ClientSet for the given builder.
func (ruo *RoleUpdateOne) getClientSet() (mc ClientSet) {
	if _, ok := ruo.config.driver.(*txDriver); ok {
		tx := &Tx{config: ruo.config}
		tx.init()
		mc = tx
	} else {
		cli := &Client{config: ruo.config}
		cli.init()
		mc = cli
	}
	return mc
}

// SaveE calls the given function after updated the Role entity,
// which is always good for cascading update operations.
func (ruo *RoleUpdateOne) SaveE(ctx context.Context, cbs ...func(ctx context.Context, mc ClientSet, updated *Role) error) (*Role, error) {
	obj, err := ruo.Save(ctx)
	if err != nil &&
		(ruo.object == nil || !errors.Is(err, stdsql.ErrNoRows)) {
		return nil, err
	}

	if len(cbs) == 0 {
		return obj, err
	}

	mc := ruo.getClientSet()

	if obj == nil {
		obj = ruo.object
	} else if x := ruo.object; x != nil {
		if _, set := ruo.mutation.Field(role.FieldDescription); set {
			obj.Description = x.Description
		}
		if _, set := ruo.mutation.Field(role.FieldPolicies); set {
			obj.Policies = x.Policies
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
func (ruo *RoleUpdateOne) SaveEX(ctx context.Context, cbs ...func(ctx context.Context, mc ClientSet, updated *Role) error) *Role {
	obj, err := ruo.SaveE(ctx, cbs...)
	if err != nil {
		panic(err)
	}
	return obj
}

// ExecE calls the given function after executed the query,
// which is always good for cascading update operations.
func (ruo *RoleUpdateOne) ExecE(ctx context.Context, cbs ...func(ctx context.Context, mc ClientSet, updated *Role) error) error {
	_, err := ruo.SaveE(ctx, cbs...)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ruo *RoleUpdateOne) ExecEX(ctx context.Context, cbs ...func(ctx context.Context, mc ClientSet, updated *Role) error) {
	if err := ruo.ExecE(ctx, cbs...); err != nil {
		panic(err)
	}
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (ruo *RoleUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *RoleUpdateOne {
	ruo.modifiers = append(ruo.modifiers, modifiers...)
	return ruo
}

func (ruo *RoleUpdateOne) sqlSave(ctx context.Context) (_node *Role, err error) {
	_spec := sqlgraph.NewUpdateSpec(role.Table, role.Columns, sqlgraph.NewFieldSpec(role.FieldID, field.TypeString))
	id, ok := ruo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`model: missing "Role.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := ruo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, role.FieldID)
		for _, f := range fields {
			if !role.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("model: invalid field %q for query", f)}
			}
			if f != role.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := ruo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ruo.mutation.UpdateTime(); ok {
		_spec.SetField(role.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := ruo.mutation.Description(); ok {
		_spec.SetField(role.FieldDescription, field.TypeString, value)
	}
	if ruo.mutation.DescriptionCleared() {
		_spec.ClearField(role.FieldDescription, field.TypeString)
	}
	if value, ok := ruo.mutation.Policies(); ok {
		_spec.SetField(role.FieldPolicies, field.TypeJSON, value)
	}
	if value, ok := ruo.mutation.AppendedPolicies(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, role.FieldPolicies, value)
		})
	}
	if ruo.mutation.SubjectsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   role.SubjectsTable,
			Columns: []string{role.SubjectsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(subjectrolerelationship.FieldID, field.TypeString),
			},
		}
		edge.Schema = ruo.schemaConfig.SubjectRoleRelationship
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ruo.mutation.RemovedSubjectsIDs(); len(nodes) > 0 && !ruo.mutation.SubjectsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   role.SubjectsTable,
			Columns: []string{role.SubjectsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(subjectrolerelationship.FieldID, field.TypeString),
			},
		}
		edge.Schema = ruo.schemaConfig.SubjectRoleRelationship
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ruo.mutation.SubjectsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   role.SubjectsTable,
			Columns: []string{role.SubjectsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(subjectrolerelationship.FieldID, field.TypeString),
			},
		}
		edge.Schema = ruo.schemaConfig.SubjectRoleRelationship
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.Node.Schema = ruo.schemaConfig.Role
	ctx = internal.NewSchemaConfigContext(ctx, ruo.schemaConfig)
	_spec.AddModifiers(ruo.modifiers...)
	_node = &Role{config: ruo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, ruo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{role.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	ruo.mutation.done = true
	return _node, nil
}
