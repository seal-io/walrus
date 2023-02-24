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

	"github.com/seal-io/seal/pkg/dao/model/allocationcost"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/types"
)

// AllocationCostQuery is the builder for querying AllocationCost entities.
type AllocationCostQuery struct {
	config
	ctx           *QueryContext
	order         []OrderFunc
	inters        []Interceptor
	predicates    []predicate.AllocationCost
	withConnector *ConnectorQuery
	modifiers     []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the AllocationCostQuery builder.
func (acq *AllocationCostQuery) Where(ps ...predicate.AllocationCost) *AllocationCostQuery {
	acq.predicates = append(acq.predicates, ps...)
	return acq
}

// Limit the number of records to be returned by this query.
func (acq *AllocationCostQuery) Limit(limit int) *AllocationCostQuery {
	acq.ctx.Limit = &limit
	return acq
}

// Offset to start from.
func (acq *AllocationCostQuery) Offset(offset int) *AllocationCostQuery {
	acq.ctx.Offset = &offset
	return acq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (acq *AllocationCostQuery) Unique(unique bool) *AllocationCostQuery {
	acq.ctx.Unique = &unique
	return acq
}

// Order specifies how the records should be ordered.
func (acq *AllocationCostQuery) Order(o ...OrderFunc) *AllocationCostQuery {
	acq.order = append(acq.order, o...)
	return acq
}

// QueryConnector chains the current query on the "connector" edge.
func (acq *AllocationCostQuery) QueryConnector() *ConnectorQuery {
	query := (&ConnectorClient{config: acq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := acq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := acq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(allocationcost.Table, allocationcost.FieldID, selector),
			sqlgraph.To(connector.Table, connector.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, allocationcost.ConnectorTable, allocationcost.ConnectorColumn),
		)
		schemaConfig := acq.schemaConfig
		step.To.Schema = schemaConfig.Connector
		step.Edge.Schema = schemaConfig.AllocationCost
		fromU = sqlgraph.SetNeighbors(acq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first AllocationCost entity from the query.
// Returns a *NotFoundError when no AllocationCost was found.
func (acq *AllocationCostQuery) First(ctx context.Context) (*AllocationCost, error) {
	nodes, err := acq.Limit(1).All(setContextOp(ctx, acq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{allocationcost.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (acq *AllocationCostQuery) FirstX(ctx context.Context) *AllocationCost {
	node, err := acq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first AllocationCost ID from the query.
// Returns a *NotFoundError when no AllocationCost ID was found.
func (acq *AllocationCostQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = acq.Limit(1).IDs(setContextOp(ctx, acq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{allocationcost.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (acq *AllocationCostQuery) FirstIDX(ctx context.Context) int {
	id, err := acq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single AllocationCost entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one AllocationCost entity is found.
// Returns a *NotFoundError when no AllocationCost entities are found.
func (acq *AllocationCostQuery) Only(ctx context.Context) (*AllocationCost, error) {
	nodes, err := acq.Limit(2).All(setContextOp(ctx, acq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{allocationcost.Label}
	default:
		return nil, &NotSingularError{allocationcost.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (acq *AllocationCostQuery) OnlyX(ctx context.Context) *AllocationCost {
	node, err := acq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only AllocationCost ID in the query.
// Returns a *NotSingularError when more than one AllocationCost ID is found.
// Returns a *NotFoundError when no entities are found.
func (acq *AllocationCostQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = acq.Limit(2).IDs(setContextOp(ctx, acq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{allocationcost.Label}
	default:
		err = &NotSingularError{allocationcost.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (acq *AllocationCostQuery) OnlyIDX(ctx context.Context) int {
	id, err := acq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of AllocationCosts.
func (acq *AllocationCostQuery) All(ctx context.Context) ([]*AllocationCost, error) {
	ctx = setContextOp(ctx, acq.ctx, "All")
	if err := acq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*AllocationCost, *AllocationCostQuery]()
	return withInterceptors[[]*AllocationCost](ctx, acq, qr, acq.inters)
}

// AllX is like All, but panics if an error occurs.
func (acq *AllocationCostQuery) AllX(ctx context.Context) []*AllocationCost {
	nodes, err := acq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of AllocationCost IDs.
func (acq *AllocationCostQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	ctx = setContextOp(ctx, acq.ctx, "IDs")
	if err := acq.Select(allocationcost.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (acq *AllocationCostQuery) IDsX(ctx context.Context) []int {
	ids, err := acq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (acq *AllocationCostQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, acq.ctx, "Count")
	if err := acq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, acq, querierCount[*AllocationCostQuery](), acq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (acq *AllocationCostQuery) CountX(ctx context.Context) int {
	count, err := acq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (acq *AllocationCostQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, acq.ctx, "Exist")
	switch _, err := acq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("model: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (acq *AllocationCostQuery) ExistX(ctx context.Context) bool {
	exist, err := acq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the AllocationCostQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (acq *AllocationCostQuery) Clone() *AllocationCostQuery {
	if acq == nil {
		return nil
	}
	return &AllocationCostQuery{
		config:        acq.config,
		ctx:           acq.ctx.Clone(),
		order:         append([]OrderFunc{}, acq.order...),
		inters:        append([]Interceptor{}, acq.inters...),
		predicates:    append([]predicate.AllocationCost{}, acq.predicates...),
		withConnector: acq.withConnector.Clone(),
		// clone intermediate query.
		sql:  acq.sql.Clone(),
		path: acq.path,
	}
}

// WithConnector tells the query-builder to eager-load the nodes that are connected to
// the "connector" edge. The optional arguments are used to configure the query builder of the edge.
func (acq *AllocationCostQuery) WithConnector(opts ...func(*ConnectorQuery)) *AllocationCostQuery {
	query := (&ConnectorClient{config: acq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	acq.withConnector = query
	return acq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		StartTime time.Time `json:"startTime"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.AllocationCost.Query().
//		GroupBy(allocationcost.FieldStartTime).
//		Aggregate(model.Count()).
//		Scan(ctx, &v)
func (acq *AllocationCostQuery) GroupBy(field string, fields ...string) *AllocationCostGroupBy {
	acq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &AllocationCostGroupBy{build: acq}
	grbuild.flds = &acq.ctx.Fields
	grbuild.label = allocationcost.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		StartTime time.Time `json:"startTime"`
//	}
//
//	client.AllocationCost.Query().
//		Select(allocationcost.FieldStartTime).
//		Scan(ctx, &v)
func (acq *AllocationCostQuery) Select(fields ...string) *AllocationCostSelect {
	acq.ctx.Fields = append(acq.ctx.Fields, fields...)
	sbuild := &AllocationCostSelect{AllocationCostQuery: acq}
	sbuild.label = allocationcost.Label
	sbuild.flds, sbuild.scan = &acq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a AllocationCostSelect configured with the given aggregations.
func (acq *AllocationCostQuery) Aggregate(fns ...AggregateFunc) *AllocationCostSelect {
	return acq.Select().Aggregate(fns...)
}

func (acq *AllocationCostQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range acq.inters {
		if inter == nil {
			return fmt.Errorf("model: uninitialized interceptor (forgotten import model/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, acq); err != nil {
				return err
			}
		}
	}
	for _, f := range acq.ctx.Fields {
		if !allocationcost.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("model: invalid field %q for query", f)}
		}
	}
	if acq.path != nil {
		prev, err := acq.path(ctx)
		if err != nil {
			return err
		}
		acq.sql = prev
	}
	return nil
}

func (acq *AllocationCostQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*AllocationCost, error) {
	var (
		nodes       = []*AllocationCost{}
		_spec       = acq.querySpec()
		loadedTypes = [1]bool{
			acq.withConnector != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*AllocationCost).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &AllocationCost{config: acq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	_spec.Node.Schema = acq.schemaConfig.AllocationCost
	ctx = internal.NewSchemaConfigContext(ctx, acq.schemaConfig)
	if len(acq.modifiers) > 0 {
		_spec.Modifiers = acq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, acq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := acq.withConnector; query != nil {
		if err := acq.loadConnector(ctx, query, nodes, nil,
			func(n *AllocationCost, e *Connector) { n.Edges.Connector = e }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (acq *AllocationCostQuery) loadConnector(ctx context.Context, query *ConnectorQuery, nodes []*AllocationCost, init func(*AllocationCost), assign func(*AllocationCost, *Connector)) error {
	ids := make([]types.ID, 0, len(nodes))
	nodeids := make(map[types.ID][]*AllocationCost)
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

func (acq *AllocationCostQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := acq.querySpec()
	_spec.Node.Schema = acq.schemaConfig.AllocationCost
	ctx = internal.NewSchemaConfigContext(ctx, acq.schemaConfig)
	if len(acq.modifiers) > 0 {
		_spec.Modifiers = acq.modifiers
	}
	_spec.Node.Columns = acq.ctx.Fields
	if len(acq.ctx.Fields) > 0 {
		_spec.Unique = acq.ctx.Unique != nil && *acq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, acq.driver, _spec)
}

func (acq *AllocationCostQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   allocationcost.Table,
			Columns: allocationcost.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: allocationcost.FieldID,
			},
		},
		From:   acq.sql,
		Unique: true,
	}
	if unique := acq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := acq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, allocationcost.FieldID)
		for i := range fields {
			if fields[i] != allocationcost.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := acq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := acq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := acq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := acq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (acq *AllocationCostQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(acq.driver.Dialect())
	t1 := builder.Table(allocationcost.Table)
	columns := acq.ctx.Fields
	if len(columns) == 0 {
		columns = allocationcost.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if acq.sql != nil {
		selector = acq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if acq.ctx.Unique != nil && *acq.ctx.Unique {
		selector.Distinct()
	}
	t1.Schema(acq.schemaConfig.AllocationCost)
	ctx = internal.NewSchemaConfigContext(ctx, acq.schemaConfig)
	selector.WithContext(ctx)
	for _, m := range acq.modifiers {
		m(selector)
	}
	for _, p := range acq.predicates {
		p(selector)
	}
	for _, p := range acq.order {
		p(selector)
	}
	if offset := acq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := acq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ForUpdate locks the selected rows against concurrent updates, and prevent them from being
// updated, deleted or "selected ... for update" by other sessions, until the transaction is
// either committed or rolled-back.
func (acq *AllocationCostQuery) ForUpdate(opts ...sql.LockOption) *AllocationCostQuery {
	if acq.driver.Dialect() == dialect.Postgres {
		acq.Unique(false)
	}
	acq.modifiers = append(acq.modifiers, func(s *sql.Selector) {
		s.ForUpdate(opts...)
	})
	return acq
}

// ForShare behaves similarly to ForUpdate, except that it acquires a shared mode lock
// on any rows that are read. Other sessions can read the rows, but cannot modify them
// until your transaction commits.
func (acq *AllocationCostQuery) ForShare(opts ...sql.LockOption) *AllocationCostQuery {
	if acq.driver.Dialect() == dialect.Postgres {
		acq.Unique(false)
	}
	acq.modifiers = append(acq.modifiers, func(s *sql.Selector) {
		s.ForShare(opts...)
	})
	return acq
}

// Modify adds a query modifier for attaching custom logic to queries.
func (acq *AllocationCostQuery) Modify(modifiers ...func(s *sql.Selector)) *AllocationCostSelect {
	acq.modifiers = append(acq.modifiers, modifiers...)
	return acq.Select()
}

// AllocationCostGroupBy is the group-by builder for AllocationCost entities.
type AllocationCostGroupBy struct {
	selector
	build *AllocationCostQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (acgb *AllocationCostGroupBy) Aggregate(fns ...AggregateFunc) *AllocationCostGroupBy {
	acgb.fns = append(acgb.fns, fns...)
	return acgb
}

// Scan applies the selector query and scans the result into the given value.
func (acgb *AllocationCostGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, acgb.build.ctx, "GroupBy")
	if err := acgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*AllocationCostQuery, *AllocationCostGroupBy](ctx, acgb.build, acgb, acgb.build.inters, v)
}

func (acgb *AllocationCostGroupBy) sqlScan(ctx context.Context, root *AllocationCostQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(acgb.fns))
	for _, fn := range acgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*acgb.flds)+len(acgb.fns))
		for _, f := range *acgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*acgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := acgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// AllocationCostSelect is the builder for selecting fields of AllocationCost entities.
type AllocationCostSelect struct {
	*AllocationCostQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (acs *AllocationCostSelect) Aggregate(fns ...AggregateFunc) *AllocationCostSelect {
	acs.fns = append(acs.fns, fns...)
	return acs
}

// Scan applies the selector query and scans the result into the given value.
func (acs *AllocationCostSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, acs.ctx, "Select")
	if err := acs.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*AllocationCostQuery, *AllocationCostSelect](ctx, acs.AllocationCostQuery, acs, acs.inters, v)
}

func (acs *AllocationCostSelect) sqlScan(ctx context.Context, root *AllocationCostQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(acs.fns))
	for _, fn := range acs.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*acs.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := acs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (acs *AllocationCostSelect) Modify(modifiers ...func(s *sql.Selector)) *AllocationCostSelect {
	acs.modifiers = append(acs.modifiers, modifiers...)
	return acs
}
