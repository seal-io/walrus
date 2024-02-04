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
	"github.com/seal-io/walrus/pkg/dao/model/workflowstageexecution"
)

// WorkflowStageExecutionDelete is the builder for deleting a WorkflowStageExecution entity.
type WorkflowStageExecutionDelete struct {
	config
	hooks    []Hook
	mutation *WorkflowStageExecutionMutation
}

// Where appends a list predicates to the WorkflowStageExecutionDelete builder.
func (wsed *WorkflowStageExecutionDelete) Where(ps ...predicate.WorkflowStageExecution) *WorkflowStageExecutionDelete {
	wsed.mutation.Where(ps...)
	return wsed
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (wsed *WorkflowStageExecutionDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, wsed.sqlExec, wsed.mutation, wsed.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (wsed *WorkflowStageExecutionDelete) ExecX(ctx context.Context) int {
	n, err := wsed.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (wsed *WorkflowStageExecutionDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(workflowstageexecution.Table, sqlgraph.NewFieldSpec(workflowstageexecution.FieldID, field.TypeString))
	_spec.Node.Schema = wsed.schemaConfig.WorkflowStageExecution
	ctx = internal.NewSchemaConfigContext(ctx, wsed.schemaConfig)
	if ps := wsed.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, wsed.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	wsed.mutation.done = true
	return affected, err
}

// WorkflowStageExecutionDeleteOne is the builder for deleting a single WorkflowStageExecution entity.
type WorkflowStageExecutionDeleteOne struct {
	wsed *WorkflowStageExecutionDelete
}

// Where appends a list predicates to the WorkflowStageExecutionDelete builder.
func (wsedo *WorkflowStageExecutionDeleteOne) Where(ps ...predicate.WorkflowStageExecution) *WorkflowStageExecutionDeleteOne {
	wsedo.wsed.mutation.Where(ps...)
	return wsedo
}

// Exec executes the deletion query.
func (wsedo *WorkflowStageExecutionDeleteOne) Exec(ctx context.Context) error {
	n, err := wsedo.wsed.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{workflowstageexecution.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (wsedo *WorkflowStageExecutionDeleteOne) ExecX(ctx context.Context) {
	if err := wsedo.Exec(ctx); err != nil {
		panic(err)
	}
}
