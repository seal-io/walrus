// SPDX-FileCopyrightText: 2024 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "walrus". DO NOT EDIT.

package model

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/walrus/pkg/dao/model/internal"
	"github.com/seal-io/walrus/pkg/dao/model/predicate"
	"github.com/seal-io/walrus/pkg/dao/model/resourcerun"
)

// ResourceRunDelete is the builder for deleting a ResourceRun entity.
type ResourceRunDelete struct {
	config
	hooks    []Hook
	mutation *ResourceRunMutation
}

// Where appends a list predicates to the ResourceRunDelete builder.
func (rrd *ResourceRunDelete) Where(ps ...predicate.ResourceRun) *ResourceRunDelete {
	rrd.mutation.Where(ps...)
	return rrd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (rrd *ResourceRunDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, rrd.sqlExec, rrd.mutation, rrd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (rrd *ResourceRunDelete) ExecX(ctx context.Context) int {
	n, err := rrd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (rrd *ResourceRunDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(resourcerun.Table, sqlgraph.NewFieldSpec(resourcerun.FieldID, field.TypeString))
	_spec.Node.Schema = rrd.schemaConfig.ResourceRun
	ctx = internal.NewSchemaConfigContext(ctx, rrd.schemaConfig)
	if ps := rrd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, rrd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	rrd.mutation.done = true
	return affected, err
}

// ResourceRunDeleteOne is the builder for deleting a single ResourceRun entity.
type ResourceRunDeleteOne struct {
	rrd *ResourceRunDelete
}

// Where appends a list predicates to the ResourceRunDelete builder.
func (rrdo *ResourceRunDeleteOne) Where(ps ...predicate.ResourceRun) *ResourceRunDeleteOne {
	rrdo.rrd.mutation.Where(ps...)
	return rrdo
}

// Exec executes the deletion query.
func (rrdo *ResourceRunDeleteOne) Exec(ctx context.Context) error {
	n, err := rrdo.rrd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{resourcerun.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (rrdo *ResourceRunDeleteOne) ExecX(ctx context.Context) {
	if err := rrdo.Exec(ctx); err != nil {
		panic(err)
	}
}
