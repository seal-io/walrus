// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package model

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/moduleversion"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
)

// ModuleVersionDelete is the builder for deleting a ModuleVersion entity.
type ModuleVersionDelete struct {
	config
	hooks    []Hook
	mutation *ModuleVersionMutation
}

// Where appends a list predicates to the ModuleVersionDelete builder.
func (mvd *ModuleVersionDelete) Where(ps ...predicate.ModuleVersion) *ModuleVersionDelete {
	mvd.mutation.Where(ps...)
	return mvd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (mvd *ModuleVersionDelete) Exec(ctx context.Context) (int, error) {
	return withHooks[int, ModuleVersionMutation](ctx, mvd.sqlExec, mvd.mutation, mvd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (mvd *ModuleVersionDelete) ExecX(ctx context.Context) int {
	n, err := mvd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (mvd *ModuleVersionDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(moduleversion.Table, sqlgraph.NewFieldSpec(moduleversion.FieldID, field.TypeString))
	_spec.Node.Schema = mvd.schemaConfig.ModuleVersion
	ctx = internal.NewSchemaConfigContext(ctx, mvd.schemaConfig)
	if ps := mvd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, mvd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	mvd.mutation.done = true
	return affected, err
}

// ModuleVersionDeleteOne is the builder for deleting a single ModuleVersion entity.
type ModuleVersionDeleteOne struct {
	mvd *ModuleVersionDelete
}

// Where appends a list predicates to the ModuleVersionDelete builder.
func (mvdo *ModuleVersionDeleteOne) Where(ps ...predicate.ModuleVersion) *ModuleVersionDeleteOne {
	mvdo.mvd.mutation.Where(ps...)
	return mvdo
}

// Exec executes the deletion query.
func (mvdo *ModuleVersionDeleteOne) Exec(ctx context.Context) error {
	n, err := mvdo.mvd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{moduleversion.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (mvdo *ModuleVersionDeleteOne) ExecX(ctx context.Context) {
	if err := mvdo.Exec(ctx); err != nil {
		panic(err)
	}
}
