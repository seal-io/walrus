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

	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/model/applicationmodulerelationship"
	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/module"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/types/oid"
)

// ApplicationModuleRelationshipQuery is the builder for querying ApplicationModuleRelationship entities.
type ApplicationModuleRelationshipQuery struct {
	config
	ctx             *QueryContext
	order           []OrderFunc
	inters          []Interceptor
	predicates      []predicate.ApplicationModuleRelationship
	withApplication *ApplicationQuery
	withModule      *ModuleQuery
	modifiers       []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the ApplicationModuleRelationshipQuery builder.
func (amrq *ApplicationModuleRelationshipQuery) Where(ps ...predicate.ApplicationModuleRelationship) *ApplicationModuleRelationshipQuery {
	amrq.predicates = append(amrq.predicates, ps...)
	return amrq
}

// Limit the number of records to be returned by this query.
func (amrq *ApplicationModuleRelationshipQuery) Limit(limit int) *ApplicationModuleRelationshipQuery {
	amrq.ctx.Limit = &limit
	return amrq
}

// Offset to start from.
func (amrq *ApplicationModuleRelationshipQuery) Offset(offset int) *ApplicationModuleRelationshipQuery {
	amrq.ctx.Offset = &offset
	return amrq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (amrq *ApplicationModuleRelationshipQuery) Unique(unique bool) *ApplicationModuleRelationshipQuery {
	amrq.ctx.Unique = &unique
	return amrq
}

// Order specifies how the records should be ordered.
func (amrq *ApplicationModuleRelationshipQuery) Order(o ...OrderFunc) *ApplicationModuleRelationshipQuery {
	amrq.order = append(amrq.order, o...)
	return amrq
}

// QueryApplication chains the current query on the "application" edge.
func (amrq *ApplicationModuleRelationshipQuery) QueryApplication() *ApplicationQuery {
	query := (&ApplicationClient{config: amrq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := amrq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := amrq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(applicationmodulerelationship.Table, applicationmodulerelationship.ApplicationColumn, selector),
			sqlgraph.To(application.Table, application.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, applicationmodulerelationship.ApplicationTable, applicationmodulerelationship.ApplicationColumn),
		)
		schemaConfig := amrq.schemaConfig
		step.To.Schema = schemaConfig.Application
		step.Edge.Schema = schemaConfig.ApplicationModuleRelationship
		fromU = sqlgraph.SetNeighbors(amrq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryModule chains the current query on the "module" edge.
func (amrq *ApplicationModuleRelationshipQuery) QueryModule() *ModuleQuery {
	query := (&ModuleClient{config: amrq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := amrq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := amrq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(applicationmodulerelationship.Table, applicationmodulerelationship.ModuleColumn, selector),
			sqlgraph.To(module.Table, module.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, applicationmodulerelationship.ModuleTable, applicationmodulerelationship.ModuleColumn),
		)
		schemaConfig := amrq.schemaConfig
		step.To.Schema = schemaConfig.Module
		step.Edge.Schema = schemaConfig.ApplicationModuleRelationship
		fromU = sqlgraph.SetNeighbors(amrq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first ApplicationModuleRelationship entity from the query.
// Returns a *NotFoundError when no ApplicationModuleRelationship was found.
func (amrq *ApplicationModuleRelationshipQuery) First(ctx context.Context) (*ApplicationModuleRelationship, error) {
	nodes, err := amrq.Limit(1).All(setContextOp(ctx, amrq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{applicationmodulerelationship.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (amrq *ApplicationModuleRelationshipQuery) FirstX(ctx context.Context) *ApplicationModuleRelationship {
	node, err := amrq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// Only returns a single ApplicationModuleRelationship entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one ApplicationModuleRelationship entity is found.
// Returns a *NotFoundError when no ApplicationModuleRelationship entities are found.
func (amrq *ApplicationModuleRelationshipQuery) Only(ctx context.Context) (*ApplicationModuleRelationship, error) {
	nodes, err := amrq.Limit(2).All(setContextOp(ctx, amrq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{applicationmodulerelationship.Label}
	default:
		return nil, &NotSingularError{applicationmodulerelationship.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (amrq *ApplicationModuleRelationshipQuery) OnlyX(ctx context.Context) *ApplicationModuleRelationship {
	node, err := amrq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// All executes the query and returns a list of ApplicationModuleRelationships.
func (amrq *ApplicationModuleRelationshipQuery) All(ctx context.Context) ([]*ApplicationModuleRelationship, error) {
	ctx = setContextOp(ctx, amrq.ctx, "All")
	if err := amrq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*ApplicationModuleRelationship, *ApplicationModuleRelationshipQuery]()
	return withInterceptors[[]*ApplicationModuleRelationship](ctx, amrq, qr, amrq.inters)
}

// AllX is like All, but panics if an error occurs.
func (amrq *ApplicationModuleRelationshipQuery) AllX(ctx context.Context) []*ApplicationModuleRelationship {
	nodes, err := amrq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// Count returns the count of the given query.
func (amrq *ApplicationModuleRelationshipQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, amrq.ctx, "Count")
	if err := amrq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, amrq, querierCount[*ApplicationModuleRelationshipQuery](), amrq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (amrq *ApplicationModuleRelationshipQuery) CountX(ctx context.Context) int {
	count, err := amrq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (amrq *ApplicationModuleRelationshipQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, amrq.ctx, "Exist")
	switch _, err := amrq.First(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("model: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (amrq *ApplicationModuleRelationshipQuery) ExistX(ctx context.Context) bool {
	exist, err := amrq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the ApplicationModuleRelationshipQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (amrq *ApplicationModuleRelationshipQuery) Clone() *ApplicationModuleRelationshipQuery {
	if amrq == nil {
		return nil
	}
	return &ApplicationModuleRelationshipQuery{
		config:          amrq.config,
		ctx:             amrq.ctx.Clone(),
		order:           append([]OrderFunc{}, amrq.order...),
		inters:          append([]Interceptor{}, amrq.inters...),
		predicates:      append([]predicate.ApplicationModuleRelationship{}, amrq.predicates...),
		withApplication: amrq.withApplication.Clone(),
		withModule:      amrq.withModule.Clone(),
		// clone intermediate query.
		sql:  amrq.sql.Clone(),
		path: amrq.path,
	}
}

// WithApplication tells the query-builder to eager-load the nodes that are connected to
// the "application" edge. The optional arguments are used to configure the query builder of the edge.
func (amrq *ApplicationModuleRelationshipQuery) WithApplication(opts ...func(*ApplicationQuery)) *ApplicationModuleRelationshipQuery {
	query := (&ApplicationClient{config: amrq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	amrq.withApplication = query
	return amrq
}

// WithModule tells the query-builder to eager-load the nodes that are connected to
// the "module" edge. The optional arguments are used to configure the query builder of the edge.
func (amrq *ApplicationModuleRelationshipQuery) WithModule(opts ...func(*ModuleQuery)) *ApplicationModuleRelationshipQuery {
	query := (&ModuleClient{config: amrq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	amrq.withModule = query
	return amrq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		CreateTime time.Time `json:"createTime,omitempty" sql:"createTime"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.ApplicationModuleRelationship.Query().
//		GroupBy(applicationmodulerelationship.FieldCreateTime).
//		Aggregate(model.Count()).
//		Scan(ctx, &v)
func (amrq *ApplicationModuleRelationshipQuery) GroupBy(field string, fields ...string) *ApplicationModuleRelationshipGroupBy {
	amrq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &ApplicationModuleRelationshipGroupBy{build: amrq}
	grbuild.flds = &amrq.ctx.Fields
	grbuild.label = applicationmodulerelationship.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		CreateTime time.Time `json:"createTime,omitempty" sql:"createTime"`
//	}
//
//	client.ApplicationModuleRelationship.Query().
//		Select(applicationmodulerelationship.FieldCreateTime).
//		Scan(ctx, &v)
func (amrq *ApplicationModuleRelationshipQuery) Select(fields ...string) *ApplicationModuleRelationshipSelect {
	amrq.ctx.Fields = append(amrq.ctx.Fields, fields...)
	sbuild := &ApplicationModuleRelationshipSelect{ApplicationModuleRelationshipQuery: amrq}
	sbuild.label = applicationmodulerelationship.Label
	sbuild.flds, sbuild.scan = &amrq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a ApplicationModuleRelationshipSelect configured with the given aggregations.
func (amrq *ApplicationModuleRelationshipQuery) Aggregate(fns ...AggregateFunc) *ApplicationModuleRelationshipSelect {
	return amrq.Select().Aggregate(fns...)
}

func (amrq *ApplicationModuleRelationshipQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range amrq.inters {
		if inter == nil {
			return fmt.Errorf("model: uninitialized interceptor (forgotten import model/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, amrq); err != nil {
				return err
			}
		}
	}
	for _, f := range amrq.ctx.Fields {
		if !applicationmodulerelationship.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("model: invalid field %q for query", f)}
		}
	}
	if amrq.path != nil {
		prev, err := amrq.path(ctx)
		if err != nil {
			return err
		}
		amrq.sql = prev
	}
	return nil
}

func (amrq *ApplicationModuleRelationshipQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*ApplicationModuleRelationship, error) {
	var (
		nodes       = []*ApplicationModuleRelationship{}
		_spec       = amrq.querySpec()
		loadedTypes = [2]bool{
			amrq.withApplication != nil,
			amrq.withModule != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*ApplicationModuleRelationship).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &ApplicationModuleRelationship{config: amrq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	_spec.Node.Schema = amrq.schemaConfig.ApplicationModuleRelationship
	ctx = internal.NewSchemaConfigContext(ctx, amrq.schemaConfig)
	if len(amrq.modifiers) > 0 {
		_spec.Modifiers = amrq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, amrq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := amrq.withApplication; query != nil {
		if err := amrq.loadApplication(ctx, query, nodes, nil,
			func(n *ApplicationModuleRelationship, e *Application) { n.Edges.Application = e }); err != nil {
			return nil, err
		}
	}
	if query := amrq.withModule; query != nil {
		if err := amrq.loadModule(ctx, query, nodes, nil,
			func(n *ApplicationModuleRelationship, e *Module) { n.Edges.Module = e }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (amrq *ApplicationModuleRelationshipQuery) loadApplication(ctx context.Context, query *ApplicationQuery, nodes []*ApplicationModuleRelationship, init func(*ApplicationModuleRelationship), assign func(*ApplicationModuleRelationship, *Application)) error {
	ids := make([]oid.ID, 0, len(nodes))
	nodeids := make(map[oid.ID][]*ApplicationModuleRelationship)
	for i := range nodes {
		fk := nodes[i].ApplicationID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(application.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "application_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (amrq *ApplicationModuleRelationshipQuery) loadModule(ctx context.Context, query *ModuleQuery, nodes []*ApplicationModuleRelationship, init func(*ApplicationModuleRelationship), assign func(*ApplicationModuleRelationship, *Module)) error {
	ids := make([]string, 0, len(nodes))
	nodeids := make(map[string][]*ApplicationModuleRelationship)
	for i := range nodes {
		fk := nodes[i].ModuleID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(module.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "module_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}

func (amrq *ApplicationModuleRelationshipQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := amrq.querySpec()
	_spec.Node.Schema = amrq.schemaConfig.ApplicationModuleRelationship
	ctx = internal.NewSchemaConfigContext(ctx, amrq.schemaConfig)
	if len(amrq.modifiers) > 0 {
		_spec.Modifiers = amrq.modifiers
	}
	_spec.Unique = false
	_spec.Node.Columns = nil
	return sqlgraph.CountNodes(ctx, amrq.driver, _spec)
}

func (amrq *ApplicationModuleRelationshipQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(applicationmodulerelationship.Table, applicationmodulerelationship.Columns, nil)
	_spec.From = amrq.sql
	if unique := amrq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if amrq.path != nil {
		_spec.Unique = true
	}
	if fields := amrq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		for i := range fields {
			_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
		}
	}
	if ps := amrq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := amrq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := amrq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := amrq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (amrq *ApplicationModuleRelationshipQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(amrq.driver.Dialect())
	t1 := builder.Table(applicationmodulerelationship.Table)
	columns := amrq.ctx.Fields
	if len(columns) == 0 {
		columns = applicationmodulerelationship.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if amrq.sql != nil {
		selector = amrq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if amrq.ctx.Unique != nil && *amrq.ctx.Unique {
		selector.Distinct()
	}
	t1.Schema(amrq.schemaConfig.ApplicationModuleRelationship)
	ctx = internal.NewSchemaConfigContext(ctx, amrq.schemaConfig)
	selector.WithContext(ctx)
	for _, m := range amrq.modifiers {
		m(selector)
	}
	for _, p := range amrq.predicates {
		p(selector)
	}
	for _, p := range amrq.order {
		p(selector)
	}
	if offset := amrq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := amrq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ForUpdate locks the selected rows against concurrent updates, and prevent them from being
// updated, deleted or "selected ... for update" by other sessions, until the transaction is
// either committed or rolled-back.
func (amrq *ApplicationModuleRelationshipQuery) ForUpdate(opts ...sql.LockOption) *ApplicationModuleRelationshipQuery {
	if amrq.driver.Dialect() == dialect.Postgres {
		amrq.Unique(false)
	}
	amrq.modifiers = append(amrq.modifiers, func(s *sql.Selector) {
		s.ForUpdate(opts...)
	})
	return amrq
}

// ForShare behaves similarly to ForUpdate, except that it acquires a shared mode lock
// on any rows that are read. Other sessions can read the rows, but cannot modify them
// until your transaction commits.
func (amrq *ApplicationModuleRelationshipQuery) ForShare(opts ...sql.LockOption) *ApplicationModuleRelationshipQuery {
	if amrq.driver.Dialect() == dialect.Postgres {
		amrq.Unique(false)
	}
	amrq.modifiers = append(amrq.modifiers, func(s *sql.Selector) {
		s.ForShare(opts...)
	})
	return amrq
}

// Modify adds a query modifier for attaching custom logic to queries.
func (amrq *ApplicationModuleRelationshipQuery) Modify(modifiers ...func(s *sql.Selector)) *ApplicationModuleRelationshipSelect {
	amrq.modifiers = append(amrq.modifiers, modifiers...)
	return amrq.Select()
}

// ApplicationModuleRelationshipGroupBy is the group-by builder for ApplicationModuleRelationship entities.
type ApplicationModuleRelationshipGroupBy struct {
	selector
	build *ApplicationModuleRelationshipQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (amrgb *ApplicationModuleRelationshipGroupBy) Aggregate(fns ...AggregateFunc) *ApplicationModuleRelationshipGroupBy {
	amrgb.fns = append(amrgb.fns, fns...)
	return amrgb
}

// Scan applies the selector query and scans the result into the given value.
func (amrgb *ApplicationModuleRelationshipGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, amrgb.build.ctx, "GroupBy")
	if err := amrgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*ApplicationModuleRelationshipQuery, *ApplicationModuleRelationshipGroupBy](ctx, amrgb.build, amrgb, amrgb.build.inters, v)
}

func (amrgb *ApplicationModuleRelationshipGroupBy) sqlScan(ctx context.Context, root *ApplicationModuleRelationshipQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(amrgb.fns))
	for _, fn := range amrgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*amrgb.flds)+len(amrgb.fns))
		for _, f := range *amrgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*amrgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := amrgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// ApplicationModuleRelationshipSelect is the builder for selecting fields of ApplicationModuleRelationship entities.
type ApplicationModuleRelationshipSelect struct {
	*ApplicationModuleRelationshipQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (amrs *ApplicationModuleRelationshipSelect) Aggregate(fns ...AggregateFunc) *ApplicationModuleRelationshipSelect {
	amrs.fns = append(amrs.fns, fns...)
	return amrs
}

// Scan applies the selector query and scans the result into the given value.
func (amrs *ApplicationModuleRelationshipSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, amrs.ctx, "Select")
	if err := amrs.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*ApplicationModuleRelationshipQuery, *ApplicationModuleRelationshipSelect](ctx, amrs.ApplicationModuleRelationshipQuery, amrs, amrs.inters, v)
}

func (amrs *ApplicationModuleRelationshipSelect) sqlScan(ctx context.Context, root *ApplicationModuleRelationshipQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(amrs.fns))
	for _, fn := range amrs.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*amrs.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := amrs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (amrs *ApplicationModuleRelationshipSelect) Modify(modifiers ...func(s *sql.Selector)) *ApplicationModuleRelationshipSelect {
	amrs.modifiers = append(amrs.modifiers, modifiers...)
	return amrs
}
