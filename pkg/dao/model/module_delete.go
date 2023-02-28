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
	"github.com/seal-io/seal/pkg/dao/model/module"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
)

// ModuleDelete is the builder for deleting a Module entity.
type ModuleDelete struct {
	config
	hooks    []Hook
	mutation *ModuleMutation
}

// Where appends a list predicates to the ModuleDelete builder.
func (md *ModuleDelete) Where(ps ...predicate.Module) *ModuleDelete {
	md.mutation.Where(ps...)
	return md
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (md *ModuleDelete) Exec(ctx context.Context) (int, error) {
	return withHooks[int, ModuleMutation](ctx, md.sqlExec, md.mutation, md.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (md *ModuleDelete) ExecX(ctx context.Context) int {
	n, err := md.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (md *ModuleDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(module.Table, sqlgraph.NewFieldSpec(module.FieldID, field.TypeString))
	_spec.Node.Schema = md.schemaConfig.Module
	ctx = internal.NewSchemaConfigContext(ctx, md.schemaConfig)
	if ps := md.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, md.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	md.mutation.done = true
	return affected, err
}

// ModuleDeleteOne is the builder for deleting a single Module entity.
type ModuleDeleteOne struct {
	md *ModuleDelete
}

// Where appends a list predicates to the ModuleDelete builder.
func (mdo *ModuleDeleteOne) Where(ps ...predicate.Module) *ModuleDeleteOne {
	mdo.md.mutation.Where(ps...)
	return mdo
}

// Exec executes the deletion query.
func (mdo *ModuleDeleteOne) Exec(ctx context.Context) error {
	n, err := mdo.md.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{module.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (mdo *ModuleDeleteOne) ExecX(ctx context.Context) {
	if err := mdo.Exec(ctx); err != nil {
		panic(err)
	}
}
