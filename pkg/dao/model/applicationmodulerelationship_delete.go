// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package model

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/model/applicationmodulerelationship"
	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
)

// ApplicationModuleRelationshipDelete is the builder for deleting a ApplicationModuleRelationship entity.
type ApplicationModuleRelationshipDelete struct {
	config
	hooks    []Hook
	mutation *ApplicationModuleRelationshipMutation
}

// Where appends a list predicates to the ApplicationModuleRelationshipDelete builder.
func (amrd *ApplicationModuleRelationshipDelete) Where(ps ...predicate.ApplicationModuleRelationship) *ApplicationModuleRelationshipDelete {
	amrd.mutation.Where(ps...)
	return amrd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (amrd *ApplicationModuleRelationshipDelete) Exec(ctx context.Context) (int, error) {
	return withHooks[int, ApplicationModuleRelationshipMutation](ctx, amrd.sqlExec, amrd.mutation, amrd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (amrd *ApplicationModuleRelationshipDelete) ExecX(ctx context.Context) int {
	n, err := amrd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (amrd *ApplicationModuleRelationshipDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: applicationmodulerelationship.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: applicationmodulerelationship.FieldID,
			},
		},
	}
	_spec.Node.Schema = amrd.schemaConfig.ApplicationModuleRelationship
	ctx = internal.NewSchemaConfigContext(ctx, amrd.schemaConfig)
	if ps := amrd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, amrd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	amrd.mutation.done = true
	return affected, err
}

// ApplicationModuleRelationshipDeleteOne is the builder for deleting a single ApplicationModuleRelationship entity.
type ApplicationModuleRelationshipDeleteOne struct {
	amrd *ApplicationModuleRelationshipDelete
}

// Where appends a list predicates to the ApplicationModuleRelationshipDelete builder.
func (amrdo *ApplicationModuleRelationshipDeleteOne) Where(ps ...predicate.ApplicationModuleRelationship) *ApplicationModuleRelationshipDeleteOne {
	amrdo.amrd.mutation.Where(ps...)
	return amrdo
}

// Exec executes the deletion query.
func (amrdo *ApplicationModuleRelationshipDeleteOne) Exec(ctx context.Context) error {
	n, err := amrdo.amrd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{applicationmodulerelationship.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (amrdo *ApplicationModuleRelationshipDeleteOne) ExecX(ctx context.Context) {
	if err := amrdo.Exec(ctx); err != nil {
		panic(err)
	}
}
