// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package intercept

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
)

// The Query interface represents an operation that queries a graph.
// By using this interface, users can write generic code that manipulates
// query builders of different types.
type Query interface {
	// Type returns the string representation of the query type.
	Type() string
	// Limit the number of records to be returned by this query.
	Limit(int)
	// Offset to start from.
	Offset(int)
	// Unique configures the query builder to filter duplicate records.
	Unique(bool)
	// Order specifies how the records should be ordered.
	Order(...model.OrderFunc)
	// WhereP appends storage-level predicates to the query builder. Using this method, users
	// can use type-assertion to append predicates that do not depend on any generated package.
	WhereP(...func(*sql.Selector))
}

// The Func type is an adapter that allows ordinary functions to be used as interceptors.
// Unlike traversal functions, interceptors are skipped during graph traversals. Note that the
// implementation of Func is different from the one defined in entgo.io/ent.InterceptFunc.
type Func func(context.Context, Query) error

// Intercept calls f(ctx, q) and then applied the next Querier.
func (f Func) Intercept(next model.Querier) model.Querier {
	return model.QuerierFunc(func(ctx context.Context, q model.Query) (model.Value, error) {
		query, err := NewQuery(q)
		if err != nil {
			return nil, err
		}
		if err := f(ctx, query); err != nil {
			return nil, err
		}
		return next.Query(ctx, q)
	})
}

// The TraverseFunc type is an adapter to allow the use of ordinary function as Traverser.
// If f is a function with the appropriate signature, TraverseFunc(f) is a Traverser that calls f.
type TraverseFunc func(context.Context, Query) error

// Intercept is a dummy implementation of Intercept that returns the next Querier in the pipeline.
func (f TraverseFunc) Intercept(next model.Querier) model.Querier {
	return next
}

// Traverse calls f(ctx, q).
func (f TraverseFunc) Traverse(ctx context.Context, q model.Query) error {
	query, err := NewQuery(q)
	if err != nil {
		return err
	}
	return f(ctx, query)
}

// The RoleFunc type is an adapter to allow the use of ordinary function as a Querier.
type RoleFunc func(context.Context, *model.RoleQuery) (model.Value, error)

// Query calls f(ctx, q).
func (f RoleFunc) Query(ctx context.Context, q model.Query) (model.Value, error) {
	if q, ok := q.(*model.RoleQuery); ok {
		return f(ctx, q)
	}
	return nil, fmt.Errorf("unexpected query type %T. expect *model.RoleQuery", q)
}

// The TraverseRole type is an adapter to allow the use of ordinary function as Traverser.
type TraverseRole func(context.Context, *model.RoleQuery) error

// Intercept is a dummy implementation of Intercept that returns the next Querier in the pipeline.
func (f TraverseRole) Intercept(next model.Querier) model.Querier {
	return next
}

// Traverse calls f(ctx, q).
func (f TraverseRole) Traverse(ctx context.Context, q model.Query) error {
	if q, ok := q.(*model.RoleQuery); ok {
		return f(ctx, q)
	}
	return fmt.Errorf("unexpected query type %T. expect *model.RoleQuery", q)
}

// The SettingFunc type is an adapter to allow the use of ordinary function as a Querier.
type SettingFunc func(context.Context, *model.SettingQuery) (model.Value, error)

// Query calls f(ctx, q).
func (f SettingFunc) Query(ctx context.Context, q model.Query) (model.Value, error) {
	if q, ok := q.(*model.SettingQuery); ok {
		return f(ctx, q)
	}
	return nil, fmt.Errorf("unexpected query type %T. expect *model.SettingQuery", q)
}

// The TraverseSetting type is an adapter to allow the use of ordinary function as Traverser.
type TraverseSetting func(context.Context, *model.SettingQuery) error

// Intercept is a dummy implementation of Intercept that returns the next Querier in the pipeline.
func (f TraverseSetting) Intercept(next model.Querier) model.Querier {
	return next
}

// Traverse calls f(ctx, q).
func (f TraverseSetting) Traverse(ctx context.Context, q model.Query) error {
	if q, ok := q.(*model.SettingQuery); ok {
		return f(ctx, q)
	}
	return fmt.Errorf("unexpected query type %T. expect *model.SettingQuery", q)
}

// The SubjectFunc type is an adapter to allow the use of ordinary function as a Querier.
type SubjectFunc func(context.Context, *model.SubjectQuery) (model.Value, error)

// Query calls f(ctx, q).
func (f SubjectFunc) Query(ctx context.Context, q model.Query) (model.Value, error) {
	if q, ok := q.(*model.SubjectQuery); ok {
		return f(ctx, q)
	}
	return nil, fmt.Errorf("unexpected query type %T. expect *model.SubjectQuery", q)
}

// The TraverseSubject type is an adapter to allow the use of ordinary function as Traverser.
type TraverseSubject func(context.Context, *model.SubjectQuery) error

// Intercept is a dummy implementation of Intercept that returns the next Querier in the pipeline.
func (f TraverseSubject) Intercept(next model.Querier) model.Querier {
	return next
}

// Traverse calls f(ctx, q).
func (f TraverseSubject) Traverse(ctx context.Context, q model.Query) error {
	if q, ok := q.(*model.SubjectQuery); ok {
		return f(ctx, q)
	}
	return fmt.Errorf("unexpected query type %T. expect *model.SubjectQuery", q)
}

// The TokenFunc type is an adapter to allow the use of ordinary function as a Querier.
type TokenFunc func(context.Context, *model.TokenQuery) (model.Value, error)

// Query calls f(ctx, q).
func (f TokenFunc) Query(ctx context.Context, q model.Query) (model.Value, error) {
	if q, ok := q.(*model.TokenQuery); ok {
		return f(ctx, q)
	}
	return nil, fmt.Errorf("unexpected query type %T. expect *model.TokenQuery", q)
}

// The TraverseToken type is an adapter to allow the use of ordinary function as Traverser.
type TraverseToken func(context.Context, *model.TokenQuery) error

// Intercept is a dummy implementation of Intercept that returns the next Querier in the pipeline.
func (f TraverseToken) Intercept(next model.Querier) model.Querier {
	return next
}

// Traverse calls f(ctx, q).
func (f TraverseToken) Traverse(ctx context.Context, q model.Query) error {
	if q, ok := q.(*model.TokenQuery); ok {
		return f(ctx, q)
	}
	return fmt.Errorf("unexpected query type %T. expect *model.TokenQuery", q)
}

// NewQuery returns the generic Query interface for the given typed query.
func NewQuery(q model.Query) (Query, error) {
	switch q := q.(type) {
	case *model.RoleQuery:
		return &query[*model.RoleQuery, predicate.Role]{typ: model.TypeRole, tq: q}, nil
	case *model.SettingQuery:
		return &query[*model.SettingQuery, predicate.Setting]{typ: model.TypeSetting, tq: q}, nil
	case *model.SubjectQuery:
		return &query[*model.SubjectQuery, predicate.Subject]{typ: model.TypeSubject, tq: q}, nil
	case *model.TokenQuery:
		return &query[*model.TokenQuery, predicate.Token]{typ: model.TypeToken, tq: q}, nil
	default:
		return nil, fmt.Errorf("unknown query type %T", q)
	}
}

type query[T any, P ~func(*sql.Selector)] struct {
	typ string
	tq  interface {
		Limit(int) T
		Offset(int) T
		Unique(bool) T
		Order(...model.OrderFunc) T
		Where(...P) T
	}
}

func (q query[T, P]) Type() string {
	return q.typ
}

func (q query[T, P]) Limit(limit int) {
	q.tq.Limit(limit)
}

func (q query[T, P]) Offset(offset int) {
	q.tq.Offset(offset)
}

func (q query[T, P]) Unique(unique bool) {
	q.tq.Unique(unique)
}

func (q query[T, P]) Order(orders ...model.OrderFunc) {
	q.tq.Order(orders...)
}

func (q query[T, P]) WhereP(ps ...func(*sql.Selector)) {
	p := make([]P, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	q.tq.Where(p...)
}
