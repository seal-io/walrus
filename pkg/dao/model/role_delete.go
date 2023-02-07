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
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/role"
)

// RoleDelete is the builder for deleting a Role entity.
type RoleDelete struct {
	config
	hooks    []Hook
	mutation *RoleMutation
}

// Where appends a list predicates to the RoleDelete builder.
func (rd *RoleDelete) Where(ps ...predicate.Role) *RoleDelete {
	rd.mutation.Where(ps...)
	return rd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (rd *RoleDelete) Exec(ctx context.Context) (int, error) {
	return withHooks[int, RoleMutation](ctx, rd.sqlExec, rd.mutation, rd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (rd *RoleDelete) ExecX(ctx context.Context) int {
	n, err := rd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (rd *RoleDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: role.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeOther,
				Column: role.FieldID,
			},
		},
	}
	_spec.Node.Schema = rd.schemaConfig.Role
	ctx = internal.NewSchemaConfigContext(ctx, rd.schemaConfig)
	if ps := rd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, rd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	rd.mutation.done = true
	return affected, err
}

// RoleDeleteOne is the builder for deleting a single Role entity.
type RoleDeleteOne struct {
	rd *RoleDelete
}

// Where appends a list predicates to the RoleDelete builder.
func (rdo *RoleDeleteOne) Where(ps ...predicate.Role) *RoleDeleteOne {
	rdo.rd.mutation.Where(ps...)
	return rdo
}

// Exec executes the deletion query.
func (rdo *RoleDeleteOne) Exec(ctx context.Context) error {
	n, err := rdo.rd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{role.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (rdo *RoleDeleteOne) ExecX(ctx context.Context) {
	if err := rdo.Exec(ctx); err != nil {
		panic(err)
	}
}
