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

// The AllocationCostFunc type is an adapter to allow the use of ordinary function as a Querier.
type AllocationCostFunc func(context.Context, *model.AllocationCostQuery) (model.Value, error)

// Query calls f(ctx, q).
func (f AllocationCostFunc) Query(ctx context.Context, q model.Query) (model.Value, error) {
	if q, ok := q.(*model.AllocationCostQuery); ok {
		return f(ctx, q)
	}
	return nil, fmt.Errorf("unexpected query type %T. expect *model.AllocationCostQuery", q)
}

// The TraverseAllocationCost type is an adapter to allow the use of ordinary function as Traverser.
type TraverseAllocationCost func(context.Context, *model.AllocationCostQuery) error

// Intercept is a dummy implementation of Intercept that returns the next Querier in the pipeline.
func (f TraverseAllocationCost) Intercept(next model.Querier) model.Querier {
	return next
}

// Traverse calls f(ctx, q).
func (f TraverseAllocationCost) Traverse(ctx context.Context, q model.Query) error {
	if q, ok := q.(*model.AllocationCostQuery); ok {
		return f(ctx, q)
	}
	return fmt.Errorf("unexpected query type %T. expect *model.AllocationCostQuery", q)
}

// The ApplicationFunc type is an adapter to allow the use of ordinary function as a Querier.
type ApplicationFunc func(context.Context, *model.ApplicationQuery) (model.Value, error)

// Query calls f(ctx, q).
func (f ApplicationFunc) Query(ctx context.Context, q model.Query) (model.Value, error) {
	if q, ok := q.(*model.ApplicationQuery); ok {
		return f(ctx, q)
	}
	return nil, fmt.Errorf("unexpected query type %T. expect *model.ApplicationQuery", q)
}

// The TraverseApplication type is an adapter to allow the use of ordinary function as Traverser.
type TraverseApplication func(context.Context, *model.ApplicationQuery) error

// Intercept is a dummy implementation of Intercept that returns the next Querier in the pipeline.
func (f TraverseApplication) Intercept(next model.Querier) model.Querier {
	return next
}

// Traverse calls f(ctx, q).
func (f TraverseApplication) Traverse(ctx context.Context, q model.Query) error {
	if q, ok := q.(*model.ApplicationQuery); ok {
		return f(ctx, q)
	}
	return fmt.Errorf("unexpected query type %T. expect *model.ApplicationQuery", q)
}

// The ApplicationModuleRelationshipFunc type is an adapter to allow the use of ordinary function as a Querier.
type ApplicationModuleRelationshipFunc func(context.Context, *model.ApplicationModuleRelationshipQuery) (model.Value, error)

// Query calls f(ctx, q).
func (f ApplicationModuleRelationshipFunc) Query(ctx context.Context, q model.Query) (model.Value, error) {
	if q, ok := q.(*model.ApplicationModuleRelationshipQuery); ok {
		return f(ctx, q)
	}
	return nil, fmt.Errorf("unexpected query type %T. expect *model.ApplicationModuleRelationshipQuery", q)
}

// The TraverseApplicationModuleRelationship type is an adapter to allow the use of ordinary function as Traverser.
type TraverseApplicationModuleRelationship func(context.Context, *model.ApplicationModuleRelationshipQuery) error

// Intercept is a dummy implementation of Intercept that returns the next Querier in the pipeline.
func (f TraverseApplicationModuleRelationship) Intercept(next model.Querier) model.Querier {
	return next
}

// Traverse calls f(ctx, q).
func (f TraverseApplicationModuleRelationship) Traverse(ctx context.Context, q model.Query) error {
	if q, ok := q.(*model.ApplicationModuleRelationshipQuery); ok {
		return f(ctx, q)
	}
	return fmt.Errorf("unexpected query type %T. expect *model.ApplicationModuleRelationshipQuery", q)
}

// The ApplicationResourceFunc type is an adapter to allow the use of ordinary function as a Querier.
type ApplicationResourceFunc func(context.Context, *model.ApplicationResourceQuery) (model.Value, error)

// Query calls f(ctx, q).
func (f ApplicationResourceFunc) Query(ctx context.Context, q model.Query) (model.Value, error) {
	if q, ok := q.(*model.ApplicationResourceQuery); ok {
		return f(ctx, q)
	}
	return nil, fmt.Errorf("unexpected query type %T. expect *model.ApplicationResourceQuery", q)
}

// The TraverseApplicationResource type is an adapter to allow the use of ordinary function as Traverser.
type TraverseApplicationResource func(context.Context, *model.ApplicationResourceQuery) error

// Intercept is a dummy implementation of Intercept that returns the next Querier in the pipeline.
func (f TraverseApplicationResource) Intercept(next model.Querier) model.Querier {
	return next
}

// Traverse calls f(ctx, q).
func (f TraverseApplicationResource) Traverse(ctx context.Context, q model.Query) error {
	if q, ok := q.(*model.ApplicationResourceQuery); ok {
		return f(ctx, q)
	}
	return fmt.Errorf("unexpected query type %T. expect *model.ApplicationResourceQuery", q)
}

// The ApplicationRevisionFunc type is an adapter to allow the use of ordinary function as a Querier.
type ApplicationRevisionFunc func(context.Context, *model.ApplicationRevisionQuery) (model.Value, error)

// Query calls f(ctx, q).
func (f ApplicationRevisionFunc) Query(ctx context.Context, q model.Query) (model.Value, error) {
	if q, ok := q.(*model.ApplicationRevisionQuery); ok {
		return f(ctx, q)
	}
	return nil, fmt.Errorf("unexpected query type %T. expect *model.ApplicationRevisionQuery", q)
}

// The TraverseApplicationRevision type is an adapter to allow the use of ordinary function as Traverser.
type TraverseApplicationRevision func(context.Context, *model.ApplicationRevisionQuery) error

// Intercept is a dummy implementation of Intercept that returns the next Querier in the pipeline.
func (f TraverseApplicationRevision) Intercept(next model.Querier) model.Querier {
	return next
}

// Traverse calls f(ctx, q).
func (f TraverseApplicationRevision) Traverse(ctx context.Context, q model.Query) error {
	if q, ok := q.(*model.ApplicationRevisionQuery); ok {
		return f(ctx, q)
	}
	return fmt.Errorf("unexpected query type %T. expect *model.ApplicationRevisionQuery", q)
}

// The ClusterCostFunc type is an adapter to allow the use of ordinary function as a Querier.
type ClusterCostFunc func(context.Context, *model.ClusterCostQuery) (model.Value, error)

// Query calls f(ctx, q).
func (f ClusterCostFunc) Query(ctx context.Context, q model.Query) (model.Value, error) {
	if q, ok := q.(*model.ClusterCostQuery); ok {
		return f(ctx, q)
	}
	return nil, fmt.Errorf("unexpected query type %T. expect *model.ClusterCostQuery", q)
}

// The TraverseClusterCost type is an adapter to allow the use of ordinary function as Traverser.
type TraverseClusterCost func(context.Context, *model.ClusterCostQuery) error

// Intercept is a dummy implementation of Intercept that returns the next Querier in the pipeline.
func (f TraverseClusterCost) Intercept(next model.Querier) model.Querier {
	return next
}

// Traverse calls f(ctx, q).
func (f TraverseClusterCost) Traverse(ctx context.Context, q model.Query) error {
	if q, ok := q.(*model.ClusterCostQuery); ok {
		return f(ctx, q)
	}
	return fmt.Errorf("unexpected query type %T. expect *model.ClusterCostQuery", q)
}

// The ConnectorFunc type is an adapter to allow the use of ordinary function as a Querier.
type ConnectorFunc func(context.Context, *model.ConnectorQuery) (model.Value, error)

// Query calls f(ctx, q).
func (f ConnectorFunc) Query(ctx context.Context, q model.Query) (model.Value, error) {
	if q, ok := q.(*model.ConnectorQuery); ok {
		return f(ctx, q)
	}
	return nil, fmt.Errorf("unexpected query type %T. expect *model.ConnectorQuery", q)
}

// The TraverseConnector type is an adapter to allow the use of ordinary function as Traverser.
type TraverseConnector func(context.Context, *model.ConnectorQuery) error

// Intercept is a dummy implementation of Intercept that returns the next Querier in the pipeline.
func (f TraverseConnector) Intercept(next model.Querier) model.Querier {
	return next
}

// Traverse calls f(ctx, q).
func (f TraverseConnector) Traverse(ctx context.Context, q model.Query) error {
	if q, ok := q.(*model.ConnectorQuery); ok {
		return f(ctx, q)
	}
	return fmt.Errorf("unexpected query type %T. expect *model.ConnectorQuery", q)
}

// The EnvironmentFunc type is an adapter to allow the use of ordinary function as a Querier.
type EnvironmentFunc func(context.Context, *model.EnvironmentQuery) (model.Value, error)

// Query calls f(ctx, q).
func (f EnvironmentFunc) Query(ctx context.Context, q model.Query) (model.Value, error) {
	if q, ok := q.(*model.EnvironmentQuery); ok {
		return f(ctx, q)
	}
	return nil, fmt.Errorf("unexpected query type %T. expect *model.EnvironmentQuery", q)
}

// The TraverseEnvironment type is an adapter to allow the use of ordinary function as Traverser.
type TraverseEnvironment func(context.Context, *model.EnvironmentQuery) error

// Intercept is a dummy implementation of Intercept that returns the next Querier in the pipeline.
func (f TraverseEnvironment) Intercept(next model.Querier) model.Querier {
	return next
}

// Traverse calls f(ctx, q).
func (f TraverseEnvironment) Traverse(ctx context.Context, q model.Query) error {
	if q, ok := q.(*model.EnvironmentQuery); ok {
		return f(ctx, q)
	}
	return fmt.Errorf("unexpected query type %T. expect *model.EnvironmentQuery", q)
}

// The EnvironmentConnectorRelationshipFunc type is an adapter to allow the use of ordinary function as a Querier.
type EnvironmentConnectorRelationshipFunc func(context.Context, *model.EnvironmentConnectorRelationshipQuery) (model.Value, error)

// Query calls f(ctx, q).
func (f EnvironmentConnectorRelationshipFunc) Query(ctx context.Context, q model.Query) (model.Value, error) {
	if q, ok := q.(*model.EnvironmentConnectorRelationshipQuery); ok {
		return f(ctx, q)
	}
	return nil, fmt.Errorf("unexpected query type %T. expect *model.EnvironmentConnectorRelationshipQuery", q)
}

// The TraverseEnvironmentConnectorRelationship type is an adapter to allow the use of ordinary function as Traverser.
type TraverseEnvironmentConnectorRelationship func(context.Context, *model.EnvironmentConnectorRelationshipQuery) error

// Intercept is a dummy implementation of Intercept that returns the next Querier in the pipeline.
func (f TraverseEnvironmentConnectorRelationship) Intercept(next model.Querier) model.Querier {
	return next
}

// Traverse calls f(ctx, q).
func (f TraverseEnvironmentConnectorRelationship) Traverse(ctx context.Context, q model.Query) error {
	if q, ok := q.(*model.EnvironmentConnectorRelationshipQuery); ok {
		return f(ctx, q)
	}
	return fmt.Errorf("unexpected query type %T. expect *model.EnvironmentConnectorRelationshipQuery", q)
}

// The ModuleFunc type is an adapter to allow the use of ordinary function as a Querier.
type ModuleFunc func(context.Context, *model.ModuleQuery) (model.Value, error)

// Query calls f(ctx, q).
func (f ModuleFunc) Query(ctx context.Context, q model.Query) (model.Value, error) {
	if q, ok := q.(*model.ModuleQuery); ok {
		return f(ctx, q)
	}
	return nil, fmt.Errorf("unexpected query type %T. expect *model.ModuleQuery", q)
}

// The TraverseModule type is an adapter to allow the use of ordinary function as Traverser.
type TraverseModule func(context.Context, *model.ModuleQuery) error

// Intercept is a dummy implementation of Intercept that returns the next Querier in the pipeline.
func (f TraverseModule) Intercept(next model.Querier) model.Querier {
	return next
}

// Traverse calls f(ctx, q).
func (f TraverseModule) Traverse(ctx context.Context, q model.Query) error {
	if q, ok := q.(*model.ModuleQuery); ok {
		return f(ctx, q)
	}
	return fmt.Errorf("unexpected query type %T. expect *model.ModuleQuery", q)
}

// The PerspectiveFunc type is an adapter to allow the use of ordinary function as a Querier.
type PerspectiveFunc func(context.Context, *model.PerspectiveQuery) (model.Value, error)

// Query calls f(ctx, q).
func (f PerspectiveFunc) Query(ctx context.Context, q model.Query) (model.Value, error) {
	if q, ok := q.(*model.PerspectiveQuery); ok {
		return f(ctx, q)
	}
	return nil, fmt.Errorf("unexpected query type %T. expect *model.PerspectiveQuery", q)
}

// The TraversePerspective type is an adapter to allow the use of ordinary function as Traverser.
type TraversePerspective func(context.Context, *model.PerspectiveQuery) error

// Intercept is a dummy implementation of Intercept that returns the next Querier in the pipeline.
func (f TraversePerspective) Intercept(next model.Querier) model.Querier {
	return next
}

// Traverse calls f(ctx, q).
func (f TraversePerspective) Traverse(ctx context.Context, q model.Query) error {
	if q, ok := q.(*model.PerspectiveQuery); ok {
		return f(ctx, q)
	}
	return fmt.Errorf("unexpected query type %T. expect *model.PerspectiveQuery", q)
}

// The ProjectFunc type is an adapter to allow the use of ordinary function as a Querier.
type ProjectFunc func(context.Context, *model.ProjectQuery) (model.Value, error)

// Query calls f(ctx, q).
func (f ProjectFunc) Query(ctx context.Context, q model.Query) (model.Value, error) {
	if q, ok := q.(*model.ProjectQuery); ok {
		return f(ctx, q)
	}
	return nil, fmt.Errorf("unexpected query type %T. expect *model.ProjectQuery", q)
}

// The TraverseProject type is an adapter to allow the use of ordinary function as Traverser.
type TraverseProject func(context.Context, *model.ProjectQuery) error

// Intercept is a dummy implementation of Intercept that returns the next Querier in the pipeline.
func (f TraverseProject) Intercept(next model.Querier) model.Querier {
	return next
}

// Traverse calls f(ctx, q).
func (f TraverseProject) Traverse(ctx context.Context, q model.Query) error {
	if q, ok := q.(*model.ProjectQuery); ok {
		return f(ctx, q)
	}
	return fmt.Errorf("unexpected query type %T. expect *model.ProjectQuery", q)
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
	case *model.AllocationCostQuery:
		return &query[*model.AllocationCostQuery, predicate.AllocationCost]{typ: model.TypeAllocationCost, tq: q}, nil
	case *model.ApplicationQuery:
		return &query[*model.ApplicationQuery, predicate.Application]{typ: model.TypeApplication, tq: q}, nil
	case *model.ApplicationModuleRelationshipQuery:
		return &query[*model.ApplicationModuleRelationshipQuery, predicate.ApplicationModuleRelationship]{typ: model.TypeApplicationModuleRelationship, tq: q}, nil
	case *model.ApplicationResourceQuery:
		return &query[*model.ApplicationResourceQuery, predicate.ApplicationResource]{typ: model.TypeApplicationResource, tq: q}, nil
	case *model.ApplicationRevisionQuery:
		return &query[*model.ApplicationRevisionQuery, predicate.ApplicationRevision]{typ: model.TypeApplicationRevision, tq: q}, nil
	case *model.ClusterCostQuery:
		return &query[*model.ClusterCostQuery, predicate.ClusterCost]{typ: model.TypeClusterCost, tq: q}, nil
	case *model.ConnectorQuery:
		return &query[*model.ConnectorQuery, predicate.Connector]{typ: model.TypeConnector, tq: q}, nil
	case *model.EnvironmentQuery:
		return &query[*model.EnvironmentQuery, predicate.Environment]{typ: model.TypeEnvironment, tq: q}, nil
	case *model.EnvironmentConnectorRelationshipQuery:
		return &query[*model.EnvironmentConnectorRelationshipQuery, predicate.EnvironmentConnectorRelationship]{typ: model.TypeEnvironmentConnectorRelationship, tq: q}, nil
	case *model.ModuleQuery:
		return &query[*model.ModuleQuery, predicate.Module]{typ: model.TypeModule, tq: q}, nil
	case *model.PerspectiveQuery:
		return &query[*model.PerspectiveQuery, predicate.Perspective]{typ: model.TypePerspective, tq: q}, nil
	case *model.ProjectQuery:
		return &query[*model.ProjectQuery, predicate.Project]{typ: model.TypeProject, tq: q}, nil
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
