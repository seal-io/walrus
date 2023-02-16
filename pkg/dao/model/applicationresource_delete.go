// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package model

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
)

// ApplicationResourceDelete is the builder for deleting a ApplicationResource entity.
type ApplicationResourceDelete struct {
	config
	hooks    []Hook
	mutation *ApplicationResourceMutation
}

// Where appends a list predicates to the ApplicationResourceDelete builder.
func (ard *ApplicationResourceDelete) Where(ps ...predicate.ApplicationResource) *ApplicationResourceDelete {
	ard.mutation.Where(ps...)
	return ard
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (ard *ApplicationResourceDelete) Exec(ctx context.Context) (int, error) {
	return withHooks[int, ApplicationResourceMutation](ctx, ard.sqlExec, ard.mutation, ard.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (ard *ApplicationResourceDelete) ExecX(ctx context.Context) int {
	n, err := ard.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (ard *ApplicationResourceDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: applicationresource.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: applicationresource.FieldID,
			},
		},
	}
	_spec.Node.Schema = ard.schemaConfig.ApplicationResource
	ctx = internal.NewSchemaConfigContext(ctx, ard.schemaConfig)
	if ps := ard.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, ard.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	ard.mutation.done = true
	return affected, err
}

// ApplicationResourceDeleteOne is the builder for deleting a single ApplicationResource entity.
type ApplicationResourceDeleteOne struct {
	ard *ApplicationResourceDelete
}

// Where appends a list predicates to the ApplicationResourceDelete builder.
func (ardo *ApplicationResourceDeleteOne) Where(ps ...predicate.ApplicationResource) *ApplicationResourceDeleteOne {
	ardo.ard.mutation.Where(ps...)
	return ardo
}

// Exec executes the deletion query.
func (ardo *ApplicationResourceDeleteOne) Exec(ctx context.Context) error {
	n, err := ardo.ard.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{applicationresource.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (ardo *ApplicationResourceDeleteOne) ExecX(ctx context.Context) {
	if err := ardo.Exec(ctx); err != nil {
		panic(err)
	}
}
