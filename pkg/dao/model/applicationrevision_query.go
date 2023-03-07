// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package model

import (
	"context"
	"fmt"
	"math"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/model/applicationinstance"
	"github.com/seal-io/seal/pkg/dao/model/applicationrevision"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/types"
)

// ApplicationRevisionQuery is the builder for querying ApplicationRevision entities.
type ApplicationRevisionQuery struct {
	config
	ctx             *QueryContext
	order           []OrderFunc
	inters          []Interceptor
	predicates      []predicate.ApplicationRevision
	withInstance    *ApplicationInstanceQuery
	withEnvironment *EnvironmentQuery
	modifiers       []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the ApplicationRevisionQuery builder.
func (arq *ApplicationRevisionQuery) Where(ps ...predicate.ApplicationRevision) *ApplicationRevisionQuery {
	arq.predicates = append(arq.predicates, ps...)
	return arq
}

// Limit the number of records to be returned by this query.
func (arq *ApplicationRevisionQuery) Limit(limit int) *ApplicationRevisionQuery {
	arq.ctx.Limit = &limit
	return arq
}

// Offset to start from.
func (arq *ApplicationRevisionQuery) Offset(offset int) *ApplicationRevisionQuery {
	arq.ctx.Offset = &offset
	return arq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (arq *ApplicationRevisionQuery) Unique(unique bool) *ApplicationRevisionQuery {
	arq.ctx.Unique = &unique
	return arq
}

// Order specifies how the records should be ordered.
func (arq *ApplicationRevisionQuery) Order(o ...OrderFunc) *ApplicationRevisionQuery {
	arq.order = append(arq.order, o...)
	return arq
}

// QueryInstance chains the current query on the "instance" edge.
func (arq *ApplicationRevisionQuery) QueryInstance() *ApplicationInstanceQuery {
	query := (&ApplicationInstanceClient{config: arq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := arq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := arq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(applicationrevision.Table, applicationrevision.FieldID, selector),
			sqlgraph.To(applicationinstance.Table, applicationinstance.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, applicationrevision.InstanceTable, applicationrevision.InstanceColumn),
		)
		schemaConfig := arq.schemaConfig
		step.To.Schema = schemaConfig.ApplicationInstance
		step.Edge.Schema = schemaConfig.ApplicationRevision
		fromU = sqlgraph.SetNeighbors(arq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryEnvironment chains the current query on the "environment" edge.
func (arq *ApplicationRevisionQuery) QueryEnvironment() *EnvironmentQuery {
	query := (&EnvironmentClient{config: arq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := arq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := arq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(applicationrevision.Table, applicationrevision.FieldID, selector),
			sqlgraph.To(environment.Table, environment.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, applicationrevision.EnvironmentTable, applicationrevision.EnvironmentColumn),
		)
		schemaConfig := arq.schemaConfig
		step.To.Schema = schemaConfig.Environment
		step.Edge.Schema = schemaConfig.ApplicationRevision
		fromU = sqlgraph.SetNeighbors(arq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first ApplicationRevision entity from the query.
// Returns a *NotFoundError when no ApplicationRevision was found.
func (arq *ApplicationRevisionQuery) First(ctx context.Context) (*ApplicationRevision, error) {
	nodes, err := arq.Limit(1).All(setContextOp(ctx, arq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{applicationrevision.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (arq *ApplicationRevisionQuery) FirstX(ctx context.Context) *ApplicationRevision {
	node, err := arq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first ApplicationRevision ID from the query.
// Returns a *NotFoundError when no ApplicationRevision ID was found.
func (arq *ApplicationRevisionQuery) FirstID(ctx context.Context) (id types.ID, err error) {
	var ids []types.ID
	if ids, err = arq.Limit(1).IDs(setContextOp(ctx, arq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{applicationrevision.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (arq *ApplicationRevisionQuery) FirstIDX(ctx context.Context) types.ID {
	id, err := arq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single ApplicationRevision entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one ApplicationRevision entity is found.
// Returns a *NotFoundError when no ApplicationRevision entities are found.
func (arq *ApplicationRevisionQuery) Only(ctx context.Context) (*ApplicationRevision, error) {
	nodes, err := arq.Limit(2).All(setContextOp(ctx, arq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{applicationrevision.Label}
	default:
		return nil, &NotSingularError{applicationrevision.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (arq *ApplicationRevisionQuery) OnlyX(ctx context.Context) *ApplicationRevision {
	node, err := arq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only ApplicationRevision ID in the query.
// Returns a *NotSingularError when more than one ApplicationRevision ID is found.
// Returns a *NotFoundError when no entities are found.
func (arq *ApplicationRevisionQuery) OnlyID(ctx context.Context) (id types.ID, err error) {
	var ids []types.ID
	if ids, err = arq.Limit(2).IDs(setContextOp(ctx, arq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{applicationrevision.Label}
	default:
		err = &NotSingularError{applicationrevision.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (arq *ApplicationRevisionQuery) OnlyIDX(ctx context.Context) types.ID {
	id, err := arq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of ApplicationRevisions.
func (arq *ApplicationRevisionQuery) All(ctx context.Context) ([]*ApplicationRevision, error) {
	ctx = setContextOp(ctx, arq.ctx, "All")
	if err := arq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*ApplicationRevision, *ApplicationRevisionQuery]()
	return withInterceptors[[]*ApplicationRevision](ctx, arq, qr, arq.inters)
}

// AllX is like All, but panics if an error occurs.
func (arq *ApplicationRevisionQuery) AllX(ctx context.Context) []*ApplicationRevision {
	nodes, err := arq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of ApplicationRevision IDs.
func (arq *ApplicationRevisionQuery) IDs(ctx context.Context) (ids []types.ID, err error) {
	if arq.ctx.Unique == nil && arq.path != nil {
		arq.Unique(true)
	}
	ctx = setContextOp(ctx, arq.ctx, "IDs")
	if err = arq.Select(applicationrevision.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (arq *ApplicationRevisionQuery) IDsX(ctx context.Context) []types.ID {
	ids, err := arq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (arq *ApplicationRevisionQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, arq.ctx, "Count")
	if err := arq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, arq, querierCount[*ApplicationRevisionQuery](), arq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (arq *ApplicationRevisionQuery) CountX(ctx context.Context) int {
	count, err := arq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (arq *ApplicationRevisionQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, arq.ctx, "Exist")
	switch _, err := arq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("model: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (arq *ApplicationRevisionQuery) ExistX(ctx context.Context) bool {
	exist, err := arq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the ApplicationRevisionQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (arq *ApplicationRevisionQuery) Clone() *ApplicationRevisionQuery {
	if arq == nil {
		return nil
	}
	return &ApplicationRevisionQuery{
		config:          arq.config,
		ctx:             arq.ctx.Clone(),
		order:           append([]OrderFunc{}, arq.order...),
		inters:          append([]Interceptor{}, arq.inters...),
		predicates:      append([]predicate.ApplicationRevision{}, arq.predicates...),
		withInstance:    arq.withInstance.Clone(),
		withEnvironment: arq.withEnvironment.Clone(),
		// clone intermediate query.
		sql:  arq.sql.Clone(),
		path: arq.path,
	}
}

// WithInstance tells the query-builder to eager-load the nodes that are connected to
// the "instance" edge. The optional arguments are used to configure the query builder of the edge.
func (arq *ApplicationRevisionQuery) WithInstance(opts ...func(*ApplicationInstanceQuery)) *ApplicationRevisionQuery {
	query := (&ApplicationInstanceClient{config: arq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	arq.withInstance = query
	return arq
}

// WithEnvironment tells the query-builder to eager-load the nodes that are connected to
// the "environment" edge. The optional arguments are used to configure the query builder of the edge.
func (arq *ApplicationRevisionQuery) WithEnvironment(opts ...func(*EnvironmentQuery)) *ApplicationRevisionQuery {
	query := (&EnvironmentClient{config: arq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	arq.withEnvironment = query
	return arq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Status string `json:"status,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.ApplicationRevision.Query().
//		GroupBy(applicationrevision.FieldStatus).
//		Aggregate(model.Count()).
//		Scan(ctx, &v)
func (arq *ApplicationRevisionQuery) GroupBy(field string, fields ...string) *ApplicationRevisionGroupBy {
	arq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &ApplicationRevisionGroupBy{build: arq}
	grbuild.flds = &arq.ctx.Fields
	grbuild.label = applicationrevision.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Status string `json:"status,omitempty"`
//	}
//
//	client.ApplicationRevision.Query().
//		Select(applicationrevision.FieldStatus).
//		Scan(ctx, &v)
func (arq *ApplicationRevisionQuery) Select(fields ...string) *ApplicationRevisionSelect {
	arq.ctx.Fields = append(arq.ctx.Fields, fields...)
	sbuild := &ApplicationRevisionSelect{ApplicationRevisionQuery: arq}
	sbuild.label = applicationrevision.Label
	sbuild.flds, sbuild.scan = &arq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a ApplicationRevisionSelect configured with the given aggregations.
func (arq *ApplicationRevisionQuery) Aggregate(fns ...AggregateFunc) *ApplicationRevisionSelect {
	return arq.Select().Aggregate(fns...)
}

func (arq *ApplicationRevisionQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range arq.inters {
		if inter == nil {
			return fmt.Errorf("model: uninitialized interceptor (forgotten import model/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, arq); err != nil {
				return err
			}
		}
	}
	for _, f := range arq.ctx.Fields {
		if !applicationrevision.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("model: invalid field %q for query", f)}
		}
	}
	if arq.path != nil {
		prev, err := arq.path(ctx)
		if err != nil {
			return err
		}
		arq.sql = prev
	}
	return nil
}

func (arq *ApplicationRevisionQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*ApplicationRevision, error) {
	var (
		nodes       = []*ApplicationRevision{}
		_spec       = arq.querySpec()
		loadedTypes = [2]bool{
			arq.withInstance != nil,
			arq.withEnvironment != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*ApplicationRevision).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &ApplicationRevision{config: arq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	_spec.Node.Schema = arq.schemaConfig.ApplicationRevision
	ctx = internal.NewSchemaConfigContext(ctx, arq.schemaConfig)
	if len(arq.modifiers) > 0 {
		_spec.Modifiers = arq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, arq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := arq.withInstance; query != nil {
		if err := arq.loadInstance(ctx, query, nodes, nil,
			func(n *ApplicationRevision, e *ApplicationInstance) { n.Edges.Instance = e }); err != nil {
			return nil, err
		}
	}
	if query := arq.withEnvironment; query != nil {
		if err := arq.loadEnvironment(ctx, query, nodes, nil,
			func(n *ApplicationRevision, e *Environment) { n.Edges.Environment = e }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (arq *ApplicationRevisionQuery) loadInstance(ctx context.Context, query *ApplicationInstanceQuery, nodes []*ApplicationRevision, init func(*ApplicationRevision), assign func(*ApplicationRevision, *ApplicationInstance)) error {
	ids := make([]types.ID, 0, len(nodes))
	nodeids := make(map[types.ID][]*ApplicationRevision)
	for i := range nodes {
		fk := nodes[i].InstanceID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(applicationinstance.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "instanceID" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (arq *ApplicationRevisionQuery) loadEnvironment(ctx context.Context, query *EnvironmentQuery, nodes []*ApplicationRevision, init func(*ApplicationRevision), assign func(*ApplicationRevision, *Environment)) error {
	ids := make([]types.ID, 0, len(nodes))
	nodeids := make(map[types.ID][]*ApplicationRevision)
	for i := range nodes {
		fk := nodes[i].EnvironmentID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(environment.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "environmentID" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}

func (arq *ApplicationRevisionQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := arq.querySpec()
	_spec.Node.Schema = arq.schemaConfig.ApplicationRevision
	ctx = internal.NewSchemaConfigContext(ctx, arq.schemaConfig)
	if len(arq.modifiers) > 0 {
		_spec.Modifiers = arq.modifiers
	}
	_spec.Node.Columns = arq.ctx.Fields
	if len(arq.ctx.Fields) > 0 {
		_spec.Unique = arq.ctx.Unique != nil && *arq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, arq.driver, _spec)
}

func (arq *ApplicationRevisionQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(applicationrevision.Table, applicationrevision.Columns, sqlgraph.NewFieldSpec(applicationrevision.FieldID, field.TypeString))
	_spec.From = arq.sql
	if unique := arq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if arq.path != nil {
		_spec.Unique = true
	}
	if fields := arq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, applicationrevision.FieldID)
		for i := range fields {
			if fields[i] != applicationrevision.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := arq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := arq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := arq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := arq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (arq *ApplicationRevisionQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(arq.driver.Dialect())
	t1 := builder.Table(applicationrevision.Table)
	columns := arq.ctx.Fields
	if len(columns) == 0 {
		columns = applicationrevision.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if arq.sql != nil {
		selector = arq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if arq.ctx.Unique != nil && *arq.ctx.Unique {
		selector.Distinct()
	}
	t1.Schema(arq.schemaConfig.ApplicationRevision)
	ctx = internal.NewSchemaConfigContext(ctx, arq.schemaConfig)
	selector.WithContext(ctx)
	for _, m := range arq.modifiers {
		m(selector)
	}
	for _, p := range arq.predicates {
		p(selector)
	}
	for _, p := range arq.order {
		p(selector)
	}
	if offset := arq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := arq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ForUpdate locks the selected rows against concurrent updates, and prevent them from being
// updated, deleted or "selected ... for update" by other sessions, until the transaction is
// either committed or rolled-back.
func (arq *ApplicationRevisionQuery) ForUpdate(opts ...sql.LockOption) *ApplicationRevisionQuery {
	if arq.driver.Dialect() == dialect.Postgres {
		arq.Unique(false)
	}
	arq.modifiers = append(arq.modifiers, func(s *sql.Selector) {
		s.ForUpdate(opts...)
	})
	return arq
}

// ForShare behaves similarly to ForUpdate, except that it acquires a shared mode lock
// on any rows that are read. Other sessions can read the rows, but cannot modify them
// until your transaction commits.
func (arq *ApplicationRevisionQuery) ForShare(opts ...sql.LockOption) *ApplicationRevisionQuery {
	if arq.driver.Dialect() == dialect.Postgres {
		arq.Unique(false)
	}
	arq.modifiers = append(arq.modifiers, func(s *sql.Selector) {
		s.ForShare(opts...)
	})
	return arq
}

// Modify adds a query modifier for attaching custom logic to queries.
func (arq *ApplicationRevisionQuery) Modify(modifiers ...func(s *sql.Selector)) *ApplicationRevisionSelect {
	arq.modifiers = append(arq.modifiers, modifiers...)
	return arq.Select()
}

// ApplicationRevisionGroupBy is the group-by builder for ApplicationRevision entities.
type ApplicationRevisionGroupBy struct {
	selector
	build *ApplicationRevisionQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (argb *ApplicationRevisionGroupBy) Aggregate(fns ...AggregateFunc) *ApplicationRevisionGroupBy {
	argb.fns = append(argb.fns, fns...)
	return argb
}

// Scan applies the selector query and scans the result into the given value.
func (argb *ApplicationRevisionGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, argb.build.ctx, "GroupBy")
	if err := argb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*ApplicationRevisionQuery, *ApplicationRevisionGroupBy](ctx, argb.build, argb, argb.build.inters, v)
}

func (argb *ApplicationRevisionGroupBy) sqlScan(ctx context.Context, root *ApplicationRevisionQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(argb.fns))
	for _, fn := range argb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*argb.flds)+len(argb.fns))
		for _, f := range *argb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*argb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := argb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// ApplicationRevisionSelect is the builder for selecting fields of ApplicationRevision entities.
type ApplicationRevisionSelect struct {
	*ApplicationRevisionQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (ars *ApplicationRevisionSelect) Aggregate(fns ...AggregateFunc) *ApplicationRevisionSelect {
	ars.fns = append(ars.fns, fns...)
	return ars
}

// Scan applies the selector query and scans the result into the given value.
func (ars *ApplicationRevisionSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ars.ctx, "Select")
	if err := ars.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*ApplicationRevisionQuery, *ApplicationRevisionSelect](ctx, ars.ApplicationRevisionQuery, ars, ars.inters, v)
}

func (ars *ApplicationRevisionSelect) sqlScan(ctx context.Context, root *ApplicationRevisionQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(ars.fns))
	for _, fn := range ars.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*ars.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ars.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (ars *ApplicationRevisionSelect) Modify(modifiers ...func(s *sql.Selector)) *ApplicationRevisionSelect {
	ars.modifiers = append(ars.modifiers, modifiers...)
	return ars
}
