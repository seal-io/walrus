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

	"github.com/seal-io/seal/pkg/dao/model/clustercost"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/types/oid"
)

// ClusterCostQuery is the builder for querying ClusterCost entities.
type ClusterCostQuery struct {
	config
	ctx           *QueryContext
	order         []OrderFunc
	inters        []Interceptor
	predicates    []predicate.ClusterCost
	withConnector *ConnectorQuery
	modifiers     []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the ClusterCostQuery builder.
func (ccq *ClusterCostQuery) Where(ps ...predicate.ClusterCost) *ClusterCostQuery {
	ccq.predicates = append(ccq.predicates, ps...)
	return ccq
}

// Limit the number of records to be returned by this query.
func (ccq *ClusterCostQuery) Limit(limit int) *ClusterCostQuery {
	ccq.ctx.Limit = &limit
	return ccq
}

// Offset to start from.
func (ccq *ClusterCostQuery) Offset(offset int) *ClusterCostQuery {
	ccq.ctx.Offset = &offset
	return ccq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (ccq *ClusterCostQuery) Unique(unique bool) *ClusterCostQuery {
	ccq.ctx.Unique = &unique
	return ccq
}

// Order specifies how the records should be ordered.
func (ccq *ClusterCostQuery) Order(o ...OrderFunc) *ClusterCostQuery {
	ccq.order = append(ccq.order, o...)
	return ccq
}

// QueryConnector chains the current query on the "connector" edge.
func (ccq *ClusterCostQuery) QueryConnector() *ConnectorQuery {
	query := (&ConnectorClient{config: ccq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := ccq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := ccq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(clustercost.Table, clustercost.FieldID, selector),
			sqlgraph.To(connector.Table, connector.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, clustercost.ConnectorTable, clustercost.ConnectorColumn),
		)
		schemaConfig := ccq.schemaConfig
		step.To.Schema = schemaConfig.Connector
		step.Edge.Schema = schemaConfig.ClusterCost
		fromU = sqlgraph.SetNeighbors(ccq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first ClusterCost entity from the query.
// Returns a *NotFoundError when no ClusterCost was found.
func (ccq *ClusterCostQuery) First(ctx context.Context) (*ClusterCost, error) {
	nodes, err := ccq.Limit(1).All(setContextOp(ctx, ccq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{clustercost.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (ccq *ClusterCostQuery) FirstX(ctx context.Context) *ClusterCost {
	node, err := ccq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first ClusterCost ID from the query.
// Returns a *NotFoundError when no ClusterCost ID was found.
func (ccq *ClusterCostQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = ccq.Limit(1).IDs(setContextOp(ctx, ccq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{clustercost.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (ccq *ClusterCostQuery) FirstIDX(ctx context.Context) int {
	id, err := ccq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single ClusterCost entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one ClusterCost entity is found.
// Returns a *NotFoundError when no ClusterCost entities are found.
func (ccq *ClusterCostQuery) Only(ctx context.Context) (*ClusterCost, error) {
	nodes, err := ccq.Limit(2).All(setContextOp(ctx, ccq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{clustercost.Label}
	default:
		return nil, &NotSingularError{clustercost.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (ccq *ClusterCostQuery) OnlyX(ctx context.Context) *ClusterCost {
	node, err := ccq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only ClusterCost ID in the query.
// Returns a *NotSingularError when more than one ClusterCost ID is found.
// Returns a *NotFoundError when no entities are found.
func (ccq *ClusterCostQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = ccq.Limit(2).IDs(setContextOp(ctx, ccq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{clustercost.Label}
	default:
		err = &NotSingularError{clustercost.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (ccq *ClusterCostQuery) OnlyIDX(ctx context.Context) int {
	id, err := ccq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of ClusterCosts.
func (ccq *ClusterCostQuery) All(ctx context.Context) ([]*ClusterCost, error) {
	ctx = setContextOp(ctx, ccq.ctx, "All")
	if err := ccq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*ClusterCost, *ClusterCostQuery]()
	return withInterceptors[[]*ClusterCost](ctx, ccq, qr, ccq.inters)
}

// AllX is like All, but panics if an error occurs.
func (ccq *ClusterCostQuery) AllX(ctx context.Context) []*ClusterCost {
	nodes, err := ccq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of ClusterCost IDs.
func (ccq *ClusterCostQuery) IDs(ctx context.Context) (ids []int, err error) {
	if ccq.ctx.Unique == nil && ccq.path != nil {
		ccq.Unique(true)
	}
	ctx = setContextOp(ctx, ccq.ctx, "IDs")
	if err = ccq.Select(clustercost.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (ccq *ClusterCostQuery) IDsX(ctx context.Context) []int {
	ids, err := ccq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (ccq *ClusterCostQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, ccq.ctx, "Count")
	if err := ccq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, ccq, querierCount[*ClusterCostQuery](), ccq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (ccq *ClusterCostQuery) CountX(ctx context.Context) int {
	count, err := ccq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (ccq *ClusterCostQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, ccq.ctx, "Exist")
	switch _, err := ccq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("model: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (ccq *ClusterCostQuery) ExistX(ctx context.Context) bool {
	exist, err := ccq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the ClusterCostQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (ccq *ClusterCostQuery) Clone() *ClusterCostQuery {
	if ccq == nil {
		return nil
	}
	return &ClusterCostQuery{
		config:        ccq.config,
		ctx:           ccq.ctx.Clone(),
		order:         append([]OrderFunc{}, ccq.order...),
		inters:        append([]Interceptor{}, ccq.inters...),
		predicates:    append([]predicate.ClusterCost{}, ccq.predicates...),
		withConnector: ccq.withConnector.Clone(),
		// clone intermediate query.
		sql:  ccq.sql.Clone(),
		path: ccq.path,
	}
}

// WithConnector tells the query-builder to eager-load the nodes that are connected to
// the "connector" edge. The optional arguments are used to configure the query builder of the edge.
func (ccq *ClusterCostQuery) WithConnector(opts ...func(*ConnectorQuery)) *ClusterCostQuery {
	query := (&ConnectorClient{config: ccq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	ccq.withConnector = query
	return ccq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		StartTime time.Time `json:"startTime,omitempty" sql:"startTime"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.ClusterCost.Query().
//		GroupBy(clustercost.FieldStartTime).
//		Aggregate(model.Count()).
//		Scan(ctx, &v)
func (ccq *ClusterCostQuery) GroupBy(field string, fields ...string) *ClusterCostGroupBy {
	ccq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &ClusterCostGroupBy{build: ccq}
	grbuild.flds = &ccq.ctx.Fields
	grbuild.label = clustercost.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		StartTime time.Time `json:"startTime,omitempty" sql:"startTime"`
//	}
//
//	client.ClusterCost.Query().
//		Select(clustercost.FieldStartTime).
//		Scan(ctx, &v)
func (ccq *ClusterCostQuery) Select(fields ...string) *ClusterCostSelect {
	ccq.ctx.Fields = append(ccq.ctx.Fields, fields...)
	sbuild := &ClusterCostSelect{ClusterCostQuery: ccq}
	sbuild.label = clustercost.Label
	sbuild.flds, sbuild.scan = &ccq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a ClusterCostSelect configured with the given aggregations.
func (ccq *ClusterCostQuery) Aggregate(fns ...AggregateFunc) *ClusterCostSelect {
	return ccq.Select().Aggregate(fns...)
}

func (ccq *ClusterCostQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range ccq.inters {
		if inter == nil {
			return fmt.Errorf("model: uninitialized interceptor (forgotten import model/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, ccq); err != nil {
				return err
			}
		}
	}
	for _, f := range ccq.ctx.Fields {
		if !clustercost.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("model: invalid field %q for query", f)}
		}
	}
	if ccq.path != nil {
		prev, err := ccq.path(ctx)
		if err != nil {
			return err
		}
		ccq.sql = prev
	}
	return nil
}

func (ccq *ClusterCostQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*ClusterCost, error) {
	var (
		nodes       = []*ClusterCost{}
		_spec       = ccq.querySpec()
		loadedTypes = [1]bool{
			ccq.withConnector != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*ClusterCost).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &ClusterCost{config: ccq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	_spec.Node.Schema = ccq.schemaConfig.ClusterCost
	ctx = internal.NewSchemaConfigContext(ctx, ccq.schemaConfig)
	if len(ccq.modifiers) > 0 {
		_spec.Modifiers = ccq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, ccq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := ccq.withConnector; query != nil {
		if err := ccq.loadConnector(ctx, query, nodes, nil,
			func(n *ClusterCost, e *Connector) { n.Edges.Connector = e }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (ccq *ClusterCostQuery) loadConnector(ctx context.Context, query *ConnectorQuery, nodes []*ClusterCost, init func(*ClusterCost), assign func(*ClusterCost, *Connector)) error {
	ids := make([]oid.ID, 0, len(nodes))
	nodeids := make(map[oid.ID][]*ClusterCost)
	for i := range nodes {
		fk := nodes[i].ConnectorID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(connector.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "connectorID" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}

func (ccq *ClusterCostQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := ccq.querySpec()
	_spec.Node.Schema = ccq.schemaConfig.ClusterCost
	ctx = internal.NewSchemaConfigContext(ctx, ccq.schemaConfig)
	if len(ccq.modifiers) > 0 {
		_spec.Modifiers = ccq.modifiers
	}
	_spec.Node.Columns = ccq.ctx.Fields
	if len(ccq.ctx.Fields) > 0 {
		_spec.Unique = ccq.ctx.Unique != nil && *ccq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, ccq.driver, _spec)
}

func (ccq *ClusterCostQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(clustercost.Table, clustercost.Columns, sqlgraph.NewFieldSpec(clustercost.FieldID, field.TypeInt))
	_spec.From = ccq.sql
	if unique := ccq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if ccq.path != nil {
		_spec.Unique = true
	}
	if fields := ccq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, clustercost.FieldID)
		for i := range fields {
			if fields[i] != clustercost.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := ccq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := ccq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := ccq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := ccq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (ccq *ClusterCostQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(ccq.driver.Dialect())
	t1 := builder.Table(clustercost.Table)
	columns := ccq.ctx.Fields
	if len(columns) == 0 {
		columns = clustercost.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if ccq.sql != nil {
		selector = ccq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if ccq.ctx.Unique != nil && *ccq.ctx.Unique {
		selector.Distinct()
	}
	t1.Schema(ccq.schemaConfig.ClusterCost)
	ctx = internal.NewSchemaConfigContext(ctx, ccq.schemaConfig)
	selector.WithContext(ctx)
	for _, m := range ccq.modifiers {
		m(selector)
	}
	for _, p := range ccq.predicates {
		p(selector)
	}
	for _, p := range ccq.order {
		p(selector)
	}
	if offset := ccq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := ccq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ForUpdate locks the selected rows against concurrent updates, and prevent them from being
// updated, deleted or "selected ... for update" by other sessions, until the transaction is
// either committed or rolled-back.
func (ccq *ClusterCostQuery) ForUpdate(opts ...sql.LockOption) *ClusterCostQuery {
	if ccq.driver.Dialect() == dialect.Postgres {
		ccq.Unique(false)
	}
	ccq.modifiers = append(ccq.modifiers, func(s *sql.Selector) {
		s.ForUpdate(opts...)
	})
	return ccq
}

// ForShare behaves similarly to ForUpdate, except that it acquires a shared mode lock
// on any rows that are read. Other sessions can read the rows, but cannot modify them
// until your transaction commits.
func (ccq *ClusterCostQuery) ForShare(opts ...sql.LockOption) *ClusterCostQuery {
	if ccq.driver.Dialect() == dialect.Postgres {
		ccq.Unique(false)
	}
	ccq.modifiers = append(ccq.modifiers, func(s *sql.Selector) {
		s.ForShare(opts...)
	})
	return ccq
}

// Modify adds a query modifier for attaching custom logic to queries.
func (ccq *ClusterCostQuery) Modify(modifiers ...func(s *sql.Selector)) *ClusterCostSelect {
	ccq.modifiers = append(ccq.modifiers, modifiers...)
	return ccq.Select()
}

// ClusterCostGroupBy is the group-by builder for ClusterCost entities.
type ClusterCostGroupBy struct {
	selector
	build *ClusterCostQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (ccgb *ClusterCostGroupBy) Aggregate(fns ...AggregateFunc) *ClusterCostGroupBy {
	ccgb.fns = append(ccgb.fns, fns...)
	return ccgb
}

// Scan applies the selector query and scans the result into the given value.
func (ccgb *ClusterCostGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ccgb.build.ctx, "GroupBy")
	if err := ccgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*ClusterCostQuery, *ClusterCostGroupBy](ctx, ccgb.build, ccgb, ccgb.build.inters, v)
}

func (ccgb *ClusterCostGroupBy) sqlScan(ctx context.Context, root *ClusterCostQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(ccgb.fns))
	for _, fn := range ccgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*ccgb.flds)+len(ccgb.fns))
		for _, f := range *ccgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*ccgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ccgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// ClusterCostSelect is the builder for selecting fields of ClusterCost entities.
type ClusterCostSelect struct {
	*ClusterCostQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (ccs *ClusterCostSelect) Aggregate(fns ...AggregateFunc) *ClusterCostSelect {
	ccs.fns = append(ccs.fns, fns...)
	return ccs
}

// Scan applies the selector query and scans the result into the given value.
func (ccs *ClusterCostSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ccs.ctx, "Select")
	if err := ccs.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*ClusterCostQuery, *ClusterCostSelect](ctx, ccs.ClusterCostQuery, ccs, ccs.inters, v)
}

func (ccs *ClusterCostSelect) sqlScan(ctx context.Context, root *ClusterCostQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(ccs.fns))
	for _, fn := range ccs.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*ccs.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ccs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (ccs *ClusterCostSelect) Modify(modifiers ...func(s *sql.Selector)) *ClusterCostSelect {
	ccs.modifiers = append(ccs.modifiers, modifiers...)
	return ccs
}
