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
	"github.com/seal-io/seal/pkg/dao/model/secret"
)

// SecretDelete is the builder for deleting a Secret entity.
type SecretDelete struct {
	config
	hooks    []Hook
	mutation *SecretMutation
}

// Where appends a list predicates to the SecretDelete builder.
func (sd *SecretDelete) Where(ps ...predicate.Secret) *SecretDelete {
	sd.mutation.Where(ps...)
	return sd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (sd *SecretDelete) Exec(ctx context.Context) (int, error) {
	return withHooks[int, SecretMutation](ctx, sd.sqlExec, sd.mutation, sd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (sd *SecretDelete) ExecX(ctx context.Context) int {
	n, err := sd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (sd *SecretDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(secret.Table, sqlgraph.NewFieldSpec(secret.FieldID, field.TypeString))
	_spec.Node.Schema = sd.schemaConfig.Secret
	ctx = internal.NewSchemaConfigContext(ctx, sd.schemaConfig)
	if ps := sd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, sd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	sd.mutation.done = true
	return affected, err
}

// SecretDeleteOne is the builder for deleting a single Secret entity.
type SecretDeleteOne struct {
	sd *SecretDelete
}

// Where appends a list predicates to the SecretDelete builder.
func (sdo *SecretDeleteOne) Where(ps ...predicate.Secret) *SecretDeleteOne {
	sdo.sd.mutation.Where(ps...)
	return sdo
}

// Exec executes the deletion query.
func (sdo *SecretDeleteOne) Exec(ctx context.Context) error {
	n, err := sdo.sd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{secret.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (sdo *SecretDeleteOne) ExecX(ctx context.Context) {
	if err := sdo.Exec(ctx); err != nil {
		panic(err)
	}
}
