// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "seal". DO NOT EDIT.

package model

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/serviceresourcerelationship"
)

// ServiceResourceRelationshipDelete is the builder for deleting a ServiceResourceRelationship entity.
type ServiceResourceRelationshipDelete struct {
	config
	hooks    []Hook
	mutation *ServiceResourceRelationshipMutation
}

// Where appends a list predicates to the ServiceResourceRelationshipDelete builder.
func (srrd *ServiceResourceRelationshipDelete) Where(ps ...predicate.ServiceResourceRelationship) *ServiceResourceRelationshipDelete {
	srrd.mutation.Where(ps...)
	return srrd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (srrd *ServiceResourceRelationshipDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, srrd.sqlExec, srrd.mutation, srrd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (srrd *ServiceResourceRelationshipDelete) ExecX(ctx context.Context) int {
	n, err := srrd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (srrd *ServiceResourceRelationshipDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(serviceresourcerelationship.Table, sqlgraph.NewFieldSpec(serviceresourcerelationship.FieldID, field.TypeString))
	_spec.Node.Schema = srrd.schemaConfig.ServiceResourceRelationship
	ctx = internal.NewSchemaConfigContext(ctx, srrd.schemaConfig)
	if ps := srrd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, srrd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	srrd.mutation.done = true
	return affected, err
}

// ServiceResourceRelationshipDeleteOne is the builder for deleting a single ServiceResourceRelationship entity.
type ServiceResourceRelationshipDeleteOne struct {
	srrd *ServiceResourceRelationshipDelete
}

// Where appends a list predicates to the ServiceResourceRelationshipDelete builder.
func (srrdo *ServiceResourceRelationshipDeleteOne) Where(ps ...predicate.ServiceResourceRelationship) *ServiceResourceRelationshipDeleteOne {
	srrdo.srrd.mutation.Where(ps...)
	return srrdo
}

// Exec executes the deletion query.
func (srrdo *ServiceResourceRelationshipDeleteOne) Exec(ctx context.Context) error {
	n, err := srrdo.srrd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{serviceresourcerelationship.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (srrdo *ServiceResourceRelationshipDeleteOne) ExecX(ctx context.Context) {
	if err := srrdo.Exec(ctx); err != nil {
		panic(err)
	}
}
