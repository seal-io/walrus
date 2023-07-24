// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "seal". DO NOT EDIT.

package model

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/subjectrolerelationship"
)

// SubjectRoleRelationshipUpdate is the builder for updating SubjectRoleRelationship entities.
type SubjectRoleRelationshipUpdate struct {
	config
	hooks     []Hook
	mutation  *SubjectRoleRelationshipMutation
	modifiers []func(*sql.UpdateBuilder)
	object    *SubjectRoleRelationship
}

// Where appends a list predicates to the SubjectRoleRelationshipUpdate builder.
func (srru *SubjectRoleRelationshipUpdate) Where(ps ...predicate.SubjectRoleRelationship) *SubjectRoleRelationshipUpdate {
	srru.mutation.Where(ps...)
	return srru
}

// Mutation returns the SubjectRoleRelationshipMutation object of the builder.
func (srru *SubjectRoleRelationshipUpdate) Mutation() *SubjectRoleRelationshipMutation {
	return srru.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (srru *SubjectRoleRelationshipUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, srru.sqlSave, srru.mutation, srru.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (srru *SubjectRoleRelationshipUpdate) SaveX(ctx context.Context) int {
	affected, err := srru.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (srru *SubjectRoleRelationshipUpdate) Exec(ctx context.Context) error {
	_, err := srru.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (srru *SubjectRoleRelationshipUpdate) ExecX(ctx context.Context) {
	if err := srru.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (srru *SubjectRoleRelationshipUpdate) check() error {
	if _, ok := srru.mutation.SubjectID(); srru.mutation.SubjectCleared() && !ok {
		return errors.New(`model: clearing a required unique edge "SubjectRoleRelationship.subject"`)
	}
	if _, ok := srru.mutation.RoleID(); srru.mutation.RoleCleared() && !ok {
		return errors.New(`model: clearing a required unique edge "SubjectRoleRelationship.role"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (srru *SubjectRoleRelationshipUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *SubjectRoleRelationshipUpdate {
	srru.modifiers = append(srru.modifiers, modifiers...)
	return srru
}

func (srru *SubjectRoleRelationshipUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := srru.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(subjectrolerelationship.Table, subjectrolerelationship.Columns, sqlgraph.NewFieldSpec(subjectrolerelationship.FieldID, field.TypeString))
	if ps := srru.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	_spec.Node.Schema = srru.schemaConfig.SubjectRoleRelationship
	ctx = internal.NewSchemaConfigContext(ctx, srru.schemaConfig)
	_spec.AddModifiers(srru.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, srru.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{subjectrolerelationship.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	srru.mutation.done = true
	return n, nil
}

// SubjectRoleRelationshipUpdateOne is the builder for updating a single SubjectRoleRelationship entity.
type SubjectRoleRelationshipUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *SubjectRoleRelationshipMutation
	modifiers []func(*sql.UpdateBuilder)
	object    *SubjectRoleRelationship
}

// Mutation returns the SubjectRoleRelationshipMutation object of the builder.
func (srruo *SubjectRoleRelationshipUpdateOne) Mutation() *SubjectRoleRelationshipMutation {
	return srruo.mutation
}

// Where appends a list predicates to the SubjectRoleRelationshipUpdate builder.
func (srruo *SubjectRoleRelationshipUpdateOne) Where(ps ...predicate.SubjectRoleRelationship) *SubjectRoleRelationshipUpdateOne {
	srruo.mutation.Where(ps...)
	return srruo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (srruo *SubjectRoleRelationshipUpdateOne) Select(field string, fields ...string) *SubjectRoleRelationshipUpdateOne {
	srruo.fields = append([]string{field}, fields...)
	return srruo
}

// Save executes the query and returns the updated SubjectRoleRelationship entity.
func (srruo *SubjectRoleRelationshipUpdateOne) Save(ctx context.Context) (*SubjectRoleRelationship, error) {
	return withHooks(ctx, srruo.sqlSave, srruo.mutation, srruo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (srruo *SubjectRoleRelationshipUpdateOne) SaveX(ctx context.Context) *SubjectRoleRelationship {
	node, err := srruo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (srruo *SubjectRoleRelationshipUpdateOne) Exec(ctx context.Context) error {
	_, err := srruo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (srruo *SubjectRoleRelationshipUpdateOne) ExecX(ctx context.Context) {
	if err := srruo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (srruo *SubjectRoleRelationshipUpdateOne) check() error {
	if _, ok := srruo.mutation.SubjectID(); srruo.mutation.SubjectCleared() && !ok {
		return errors.New(`model: clearing a required unique edge "SubjectRoleRelationship.subject"`)
	}
	if _, ok := srruo.mutation.RoleID(); srruo.mutation.RoleCleared() && !ok {
		return errors.New(`model: clearing a required unique edge "SubjectRoleRelationship.role"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (srruo *SubjectRoleRelationshipUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *SubjectRoleRelationshipUpdateOne {
	srruo.modifiers = append(srruo.modifiers, modifiers...)
	return srruo
}

func (srruo *SubjectRoleRelationshipUpdateOne) sqlSave(ctx context.Context) (_node *SubjectRoleRelationship, err error) {
	if err := srruo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(subjectrolerelationship.Table, subjectrolerelationship.Columns, sqlgraph.NewFieldSpec(subjectrolerelationship.FieldID, field.TypeString))
	id, ok := srruo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`model: missing "SubjectRoleRelationship.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := srruo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, subjectrolerelationship.FieldID)
		for _, f := range fields {
			if !subjectrolerelationship.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("model: invalid field %q for query", f)}
			}
			if f != subjectrolerelationship.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := srruo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	_spec.Node.Schema = srruo.schemaConfig.SubjectRoleRelationship
	ctx = internal.NewSchemaConfigContext(ctx, srruo.schemaConfig)
	_spec.AddModifiers(srruo.modifiers...)
	_node = &SubjectRoleRelationship{config: srruo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, srruo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{subjectrolerelationship.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	srruo.mutation.done = true
	return _node, nil
}
