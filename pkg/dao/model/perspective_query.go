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

	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/perspective"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/types/oid"
)

// PerspectiveQuery is the builder for querying Perspective entities.
type PerspectiveQuery struct {
	config
	ctx        *QueryContext
	order      []OrderFunc
	inters     []Interceptor
	predicates []predicate.Perspective
	modifiers  []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the PerspectiveQuery builder.
func (pq *PerspectiveQuery) Where(ps ...predicate.Perspective) *PerspectiveQuery {
	pq.predicates = append(pq.predicates, ps...)
	return pq
}

// Limit the number of records to be returned by this query.
func (pq *PerspectiveQuery) Limit(limit int) *PerspectiveQuery {
	pq.ctx.Limit = &limit
	return pq
}

// Offset to start from.
func (pq *PerspectiveQuery) Offset(offset int) *PerspectiveQuery {
	pq.ctx.Offset = &offset
	return pq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (pq *PerspectiveQuery) Unique(unique bool) *PerspectiveQuery {
	pq.ctx.Unique = &unique
	return pq
}

// Order specifies how the records should be ordered.
func (pq *PerspectiveQuery) Order(o ...OrderFunc) *PerspectiveQuery {
	pq.order = append(pq.order, o...)
	return pq
}

// First returns the first Perspective entity from the query.
// Returns a *NotFoundError when no Perspective was found.
func (pq *PerspectiveQuery) First(ctx context.Context) (*Perspective, error) {
	nodes, err := pq.Limit(1).All(setContextOp(ctx, pq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{perspective.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (pq *PerspectiveQuery) FirstX(ctx context.Context) *Perspective {
	node, err := pq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Perspective ID from the query.
// Returns a *NotFoundError when no Perspective ID was found.
func (pq *PerspectiveQuery) FirstID(ctx context.Context) (id oid.ID, err error) {
	var ids []oid.ID
	if ids, err = pq.Limit(1).IDs(setContextOp(ctx, pq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{perspective.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (pq *PerspectiveQuery) FirstIDX(ctx context.Context) oid.ID {
	id, err := pq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Perspective entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Perspective entity is found.
// Returns a *NotFoundError when no Perspective entities are found.
func (pq *PerspectiveQuery) Only(ctx context.Context) (*Perspective, error) {
	nodes, err := pq.Limit(2).All(setContextOp(ctx, pq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{perspective.Label}
	default:
		return nil, &NotSingularError{perspective.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (pq *PerspectiveQuery) OnlyX(ctx context.Context) *Perspective {
	node, err := pq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Perspective ID in the query.
// Returns a *NotSingularError when more than one Perspective ID is found.
// Returns a *NotFoundError when no entities are found.
func (pq *PerspectiveQuery) OnlyID(ctx context.Context) (id oid.ID, err error) {
	var ids []oid.ID
	if ids, err = pq.Limit(2).IDs(setContextOp(ctx, pq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{perspective.Label}
	default:
		err = &NotSingularError{perspective.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (pq *PerspectiveQuery) OnlyIDX(ctx context.Context) oid.ID {
	id, err := pq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Perspectives.
func (pq *PerspectiveQuery) All(ctx context.Context) ([]*Perspective, error) {
	ctx = setContextOp(ctx, pq.ctx, "All")
	if err := pq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Perspective, *PerspectiveQuery]()
	return withInterceptors[[]*Perspective](ctx, pq, qr, pq.inters)
}

// AllX is like All, but panics if an error occurs.
func (pq *PerspectiveQuery) AllX(ctx context.Context) []*Perspective {
	nodes, err := pq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Perspective IDs.
func (pq *PerspectiveQuery) IDs(ctx context.Context) (ids []oid.ID, err error) {
	if pq.ctx.Unique == nil && pq.path != nil {
		pq.Unique(true)
	}
	ctx = setContextOp(ctx, pq.ctx, "IDs")
	if err = pq.Select(perspective.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (pq *PerspectiveQuery) IDsX(ctx context.Context) []oid.ID {
	ids, err := pq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (pq *PerspectiveQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, pq.ctx, "Count")
	if err := pq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, pq, querierCount[*PerspectiveQuery](), pq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (pq *PerspectiveQuery) CountX(ctx context.Context) int {
	count, err := pq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (pq *PerspectiveQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, pq.ctx, "Exist")
	switch _, err := pq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("model: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (pq *PerspectiveQuery) ExistX(ctx context.Context) bool {
	exist, err := pq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the PerspectiveQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (pq *PerspectiveQuery) Clone() *PerspectiveQuery {
	if pq == nil {
		return nil
	}
	return &PerspectiveQuery{
		config:     pq.config,
		ctx:        pq.ctx.Clone(),
		order:      append([]OrderFunc{}, pq.order...),
		inters:     append([]Interceptor{}, pq.inters...),
		predicates: append([]predicate.Perspective{}, pq.predicates...),
		// clone intermediate query.
		sql:  pq.sql.Clone(),
		path: pq.path,
	}
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
//	client.Perspective.Query().
//		GroupBy(perspective.FieldCreateTime).
//		Aggregate(model.Count()).
//		Scan(ctx, &v)
func (pq *PerspectiveQuery) GroupBy(field string, fields ...string) *PerspectiveGroupBy {
	pq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &PerspectiveGroupBy{build: pq}
	grbuild.flds = &pq.ctx.Fields
	grbuild.label = perspective.Label
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
//	client.Perspective.Query().
//		Select(perspective.FieldCreateTime).
//		Scan(ctx, &v)
func (pq *PerspectiveQuery) Select(fields ...string) *PerspectiveSelect {
	pq.ctx.Fields = append(pq.ctx.Fields, fields...)
	sbuild := &PerspectiveSelect{PerspectiveQuery: pq}
	sbuild.label = perspective.Label
	sbuild.flds, sbuild.scan = &pq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a PerspectiveSelect configured with the given aggregations.
func (pq *PerspectiveQuery) Aggregate(fns ...AggregateFunc) *PerspectiveSelect {
	return pq.Select().Aggregate(fns...)
}

func (pq *PerspectiveQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range pq.inters {
		if inter == nil {
			return fmt.Errorf("model: uninitialized interceptor (forgotten import model/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, pq); err != nil {
				return err
			}
		}
	}
	for _, f := range pq.ctx.Fields {
		if !perspective.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("model: invalid field %q for query", f)}
		}
	}
	if pq.path != nil {
		prev, err := pq.path(ctx)
		if err != nil {
			return err
		}
		pq.sql = prev
	}
	return nil
}

func (pq *PerspectiveQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Perspective, error) {
	var (
		nodes = []*Perspective{}
		_spec = pq.querySpec()
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Perspective).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Perspective{config: pq.config}
		nodes = append(nodes, node)
		return node.assignValues(columns, values)
	}
	_spec.Node.Schema = pq.schemaConfig.Perspective
	ctx = internal.NewSchemaConfigContext(ctx, pq.schemaConfig)
	if len(pq.modifiers) > 0 {
		_spec.Modifiers = pq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, pq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (pq *PerspectiveQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := pq.querySpec()
	_spec.Node.Schema = pq.schemaConfig.Perspective
	ctx = internal.NewSchemaConfigContext(ctx, pq.schemaConfig)
	if len(pq.modifiers) > 0 {
		_spec.Modifiers = pq.modifiers
	}
	_spec.Node.Columns = pq.ctx.Fields
	if len(pq.ctx.Fields) > 0 {
		_spec.Unique = pq.ctx.Unique != nil && *pq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, pq.driver, _spec)
}

func (pq *PerspectiveQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(perspective.Table, perspective.Columns, sqlgraph.NewFieldSpec(perspective.FieldID, field.TypeString))
	_spec.From = pq.sql
	if unique := pq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if pq.path != nil {
		_spec.Unique = true
	}
	if fields := pq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, perspective.FieldID)
		for i := range fields {
			if fields[i] != perspective.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := pq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := pq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := pq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := pq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (pq *PerspectiveQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(pq.driver.Dialect())
	t1 := builder.Table(perspective.Table)
	columns := pq.ctx.Fields
	if len(columns) == 0 {
		columns = perspective.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if pq.sql != nil {
		selector = pq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if pq.ctx.Unique != nil && *pq.ctx.Unique {
		selector.Distinct()
	}
	t1.Schema(pq.schemaConfig.Perspective)
	ctx = internal.NewSchemaConfigContext(ctx, pq.schemaConfig)
	selector.WithContext(ctx)
	for _, m := range pq.modifiers {
		m(selector)
	}
	for _, p := range pq.predicates {
		p(selector)
	}
	for _, p := range pq.order {
		p(selector)
	}
	if offset := pq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := pq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ForUpdate locks the selected rows against concurrent updates, and prevent them from being
// updated, deleted or "selected ... for update" by other sessions, until the transaction is
// either committed or rolled-back.
func (pq *PerspectiveQuery) ForUpdate(opts ...sql.LockOption) *PerspectiveQuery {
	if pq.driver.Dialect() == dialect.Postgres {
		pq.Unique(false)
	}
	pq.modifiers = append(pq.modifiers, func(s *sql.Selector) {
		s.ForUpdate(opts...)
	})
	return pq
}

// ForShare behaves similarly to ForUpdate, except that it acquires a shared mode lock
// on any rows that are read. Other sessions can read the rows, but cannot modify them
// until your transaction commits.
func (pq *PerspectiveQuery) ForShare(opts ...sql.LockOption) *PerspectiveQuery {
	if pq.driver.Dialect() == dialect.Postgres {
		pq.Unique(false)
	}
	pq.modifiers = append(pq.modifiers, func(s *sql.Selector) {
		s.ForShare(opts...)
	})
	return pq
}

// Modify adds a query modifier for attaching custom logic to queries.
func (pq *PerspectiveQuery) Modify(modifiers ...func(s *sql.Selector)) *PerspectiveSelect {
	pq.modifiers = append(pq.modifiers, modifiers...)
	return pq.Select()
}

// WhereP appends storage-level predicates to the PerspectiveQuery builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (pq *PerspectiveQuery) WhereP(ps ...func(*sql.Selector)) {
	var wps = make([]predicate.Perspective, 0, len(ps))
	for i := 0; i < len(ps); i++ {
		wps = append(wps, predicate.Perspective(ps[i]))
	}
	pq.predicates = append(pq.predicates, wps...)
}

// PerspectiveGroupBy is the group-by builder for Perspective entities.
type PerspectiveGroupBy struct {
	selector
	build *PerspectiveQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (pgb *PerspectiveGroupBy) Aggregate(fns ...AggregateFunc) *PerspectiveGroupBy {
	pgb.fns = append(pgb.fns, fns...)
	return pgb
}

// Scan applies the selector query and scans the result into the given value.
func (pgb *PerspectiveGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, pgb.build.ctx, "GroupBy")
	if err := pgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*PerspectiveQuery, *PerspectiveGroupBy](ctx, pgb.build, pgb, pgb.build.inters, v)
}

func (pgb *PerspectiveGroupBy) sqlScan(ctx context.Context, root *PerspectiveQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(pgb.fns))
	for _, fn := range pgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*pgb.flds)+len(pgb.fns))
		for _, f := range *pgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*pgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := pgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// PerspectiveSelect is the builder for selecting fields of Perspective entities.
type PerspectiveSelect struct {
	*PerspectiveQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (ps *PerspectiveSelect) Aggregate(fns ...AggregateFunc) *PerspectiveSelect {
	ps.fns = append(ps.fns, fns...)
	return ps
}

// Scan applies the selector query and scans the result into the given value.
func (ps *PerspectiveSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ps.ctx, "Select")
	if err := ps.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*PerspectiveQuery, *PerspectiveSelect](ctx, ps.PerspectiveQuery, ps, ps.inters, v)
}

func (ps *PerspectiveSelect) sqlScan(ctx context.Context, root *PerspectiveQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(ps.fns))
	for _, fn := range ps.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*ps.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ps.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (ps *PerspectiveSelect) Modify(modifiers ...func(s *sql.Selector)) *PerspectiveSelect {
	ps.modifiers = append(ps.modifiers, modifiers...)
	return ps
}
