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

	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/oid"
)

// ApplicationResourceQuery is the builder for querying ApplicationResource entities.
type ApplicationResourceQuery struct {
	config
	ctx        *QueryContext
	order      []OrderFunc
	inters     []Interceptor
	predicates []predicate.ApplicationResource
	modifiers  []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the ApplicationResourceQuery builder.
func (arq *ApplicationResourceQuery) Where(ps ...predicate.ApplicationResource) *ApplicationResourceQuery {
	arq.predicates = append(arq.predicates, ps...)
	return arq
}

// Limit the number of records to be returned by this query.
func (arq *ApplicationResourceQuery) Limit(limit int) *ApplicationResourceQuery {
	arq.ctx.Limit = &limit
	return arq
}

// Offset to start from.
func (arq *ApplicationResourceQuery) Offset(offset int) *ApplicationResourceQuery {
	arq.ctx.Offset = &offset
	return arq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (arq *ApplicationResourceQuery) Unique(unique bool) *ApplicationResourceQuery {
	arq.ctx.Unique = &unique
	return arq
}

// Order specifies how the records should be ordered.
func (arq *ApplicationResourceQuery) Order(o ...OrderFunc) *ApplicationResourceQuery {
	arq.order = append(arq.order, o...)
	return arq
}

// First returns the first ApplicationResource entity from the query.
// Returns a *NotFoundError when no ApplicationResource was found.
func (arq *ApplicationResourceQuery) First(ctx context.Context) (*ApplicationResource, error) {
	nodes, err := arq.Limit(1).All(setContextOp(ctx, arq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{applicationresource.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (arq *ApplicationResourceQuery) FirstX(ctx context.Context) *ApplicationResource {
	node, err := arq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first ApplicationResource ID from the query.
// Returns a *NotFoundError when no ApplicationResource ID was found.
func (arq *ApplicationResourceQuery) FirstID(ctx context.Context) (id oid.ID, err error) {
	var ids []oid.ID
	if ids, err = arq.Limit(1).IDs(setContextOp(ctx, arq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{applicationresource.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (arq *ApplicationResourceQuery) FirstIDX(ctx context.Context) oid.ID {
	id, err := arq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single ApplicationResource entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one ApplicationResource entity is found.
// Returns a *NotFoundError when no ApplicationResource entities are found.
func (arq *ApplicationResourceQuery) Only(ctx context.Context) (*ApplicationResource, error) {
	nodes, err := arq.Limit(2).All(setContextOp(ctx, arq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{applicationresource.Label}
	default:
		return nil, &NotSingularError{applicationresource.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (arq *ApplicationResourceQuery) OnlyX(ctx context.Context) *ApplicationResource {
	node, err := arq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only ApplicationResource ID in the query.
// Returns a *NotSingularError when more than one ApplicationResource ID is found.
// Returns a *NotFoundError when no entities are found.
func (arq *ApplicationResourceQuery) OnlyID(ctx context.Context) (id oid.ID, err error) {
	var ids []oid.ID
	if ids, err = arq.Limit(2).IDs(setContextOp(ctx, arq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{applicationresource.Label}
	default:
		err = &NotSingularError{applicationresource.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (arq *ApplicationResourceQuery) OnlyIDX(ctx context.Context) oid.ID {
	id, err := arq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of ApplicationResources.
func (arq *ApplicationResourceQuery) All(ctx context.Context) ([]*ApplicationResource, error) {
	ctx = setContextOp(ctx, arq.ctx, "All")
	if err := arq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*ApplicationResource, *ApplicationResourceQuery]()
	return withInterceptors[[]*ApplicationResource](ctx, arq, qr, arq.inters)
}

// AllX is like All, but panics if an error occurs.
func (arq *ApplicationResourceQuery) AllX(ctx context.Context) []*ApplicationResource {
	nodes, err := arq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of ApplicationResource IDs.
func (arq *ApplicationResourceQuery) IDs(ctx context.Context) ([]oid.ID, error) {
	var ids []oid.ID
	ctx = setContextOp(ctx, arq.ctx, "IDs")
	if err := arq.Select(applicationresource.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (arq *ApplicationResourceQuery) IDsX(ctx context.Context) []oid.ID {
	ids, err := arq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (arq *ApplicationResourceQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, arq.ctx, "Count")
	if err := arq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, arq, querierCount[*ApplicationResourceQuery](), arq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (arq *ApplicationResourceQuery) CountX(ctx context.Context) int {
	count, err := arq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (arq *ApplicationResourceQuery) Exist(ctx context.Context) (bool, error) {
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
func (arq *ApplicationResourceQuery) ExistX(ctx context.Context) bool {
	exist, err := arq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the ApplicationResourceQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (arq *ApplicationResourceQuery) Clone() *ApplicationResourceQuery {
	if arq == nil {
		return nil
	}
	return &ApplicationResourceQuery{
		config:     arq.config,
		ctx:        arq.ctx.Clone(),
		order:      append([]OrderFunc{}, arq.order...),
		inters:     append([]Interceptor{}, arq.inters...),
		predicates: append([]predicate.ApplicationResource{}, arq.predicates...),
		// clone intermediate query.
		sql:  arq.sql.Clone(),
		path: arq.path,
	}
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
//	client.ApplicationResource.Query().
//		GroupBy(applicationresource.FieldStatus).
//		Aggregate(model.Count()).
//		Scan(ctx, &v)
func (arq *ApplicationResourceQuery) GroupBy(field string, fields ...string) *ApplicationResourceGroupBy {
	arq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &ApplicationResourceGroupBy{build: arq}
	grbuild.flds = &arq.ctx.Fields
	grbuild.label = applicationresource.Label
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
//	client.ApplicationResource.Query().
//		Select(applicationresource.FieldStatus).
//		Scan(ctx, &v)
func (arq *ApplicationResourceQuery) Select(fields ...string) *ApplicationResourceSelect {
	arq.ctx.Fields = append(arq.ctx.Fields, fields...)
	sbuild := &ApplicationResourceSelect{ApplicationResourceQuery: arq}
	sbuild.label = applicationresource.Label
	sbuild.flds, sbuild.scan = &arq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a ApplicationResourceSelect configured with the given aggregations.
func (arq *ApplicationResourceQuery) Aggregate(fns ...AggregateFunc) *ApplicationResourceSelect {
	return arq.Select().Aggregate(fns...)
}

func (arq *ApplicationResourceQuery) prepareQuery(ctx context.Context) error {
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
		if !applicationresource.ValidColumn(f) {
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

func (arq *ApplicationResourceQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*ApplicationResource, error) {
	var (
		nodes = []*ApplicationResource{}
		_spec = arq.querySpec()
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*ApplicationResource).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &ApplicationResource{config: arq.config}
		nodes = append(nodes, node)
		return node.assignValues(columns, values)
	}
	_spec.Node.Schema = arq.schemaConfig.ApplicationResource
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
	return nodes, nil
}

func (arq *ApplicationResourceQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := arq.querySpec()
	_spec.Node.Schema = arq.schemaConfig.ApplicationResource
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

func (arq *ApplicationResourceQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   applicationresource.Table,
			Columns: applicationresource.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeOther,
				Column: applicationresource.FieldID,
			},
		},
		From:   arq.sql,
		Unique: true,
	}
	if unique := arq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := arq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, applicationresource.FieldID)
		for i := range fields {
			if fields[i] != applicationresource.FieldID {
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

func (arq *ApplicationResourceQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(arq.driver.Dialect())
	t1 := builder.Table(applicationresource.Table)
	columns := arq.ctx.Fields
	if len(columns) == 0 {
		columns = applicationresource.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if arq.sql != nil {
		selector = arq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if arq.ctx.Unique != nil && *arq.ctx.Unique {
		selector.Distinct()
	}
	t1.Schema(arq.schemaConfig.ApplicationResource)
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
func (arq *ApplicationResourceQuery) ForUpdate(opts ...sql.LockOption) *ApplicationResourceQuery {
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
func (arq *ApplicationResourceQuery) ForShare(opts ...sql.LockOption) *ApplicationResourceQuery {
	if arq.driver.Dialect() == dialect.Postgres {
		arq.Unique(false)
	}
	arq.modifiers = append(arq.modifiers, func(s *sql.Selector) {
		s.ForShare(opts...)
	})
	return arq
}

// Modify adds a query modifier for attaching custom logic to queries.
func (arq *ApplicationResourceQuery) Modify(modifiers ...func(s *sql.Selector)) *ApplicationResourceSelect {
	arq.modifiers = append(arq.modifiers, modifiers...)
	return arq.Select()
}

// ApplicationResourceGroupBy is the group-by builder for ApplicationResource entities.
type ApplicationResourceGroupBy struct {
	selector
	build *ApplicationResourceQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (argb *ApplicationResourceGroupBy) Aggregate(fns ...AggregateFunc) *ApplicationResourceGroupBy {
	argb.fns = append(argb.fns, fns...)
	return argb
}

// Scan applies the selector query and scans the result into the given value.
func (argb *ApplicationResourceGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, argb.build.ctx, "GroupBy")
	if err := argb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*ApplicationResourceQuery, *ApplicationResourceGroupBy](ctx, argb.build, argb, argb.build.inters, v)
}

func (argb *ApplicationResourceGroupBy) sqlScan(ctx context.Context, root *ApplicationResourceQuery, v any) error {
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

// ApplicationResourceSelect is the builder for selecting fields of ApplicationResource entities.
type ApplicationResourceSelect struct {
	*ApplicationResourceQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (ars *ApplicationResourceSelect) Aggregate(fns ...AggregateFunc) *ApplicationResourceSelect {
	ars.fns = append(ars.fns, fns...)
	return ars
}

// Scan applies the selector query and scans the result into the given value.
func (ars *ApplicationResourceSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ars.ctx, "Select")
	if err := ars.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*ApplicationResourceQuery, *ApplicationResourceSelect](ctx, ars.ApplicationResourceQuery, ars, ars.inters, v)
}

func (ars *ApplicationResourceSelect) sqlScan(ctx context.Context, root *ApplicationResourceQuery, v any) error {
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
func (ars *ApplicationResourceSelect) Modify(modifiers ...func(s *sql.Selector)) *ApplicationResourceSelect {
	ars.modifiers = append(ars.modifiers, modifiers...)
	return ars
}
