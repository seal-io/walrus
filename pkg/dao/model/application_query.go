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

	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/oid"
)

// ApplicationQuery is the builder for querying Application entities.
type ApplicationQuery struct {
	config
	ctx        *QueryContext
	order      []OrderFunc
	inters     []Interceptor
	predicates []predicate.Application
	modifiers  []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the ApplicationQuery builder.
func (aq *ApplicationQuery) Where(ps ...predicate.Application) *ApplicationQuery {
	aq.predicates = append(aq.predicates, ps...)
	return aq
}

// Limit the number of records to be returned by this query.
func (aq *ApplicationQuery) Limit(limit int) *ApplicationQuery {
	aq.ctx.Limit = &limit
	return aq
}

// Offset to start from.
func (aq *ApplicationQuery) Offset(offset int) *ApplicationQuery {
	aq.ctx.Offset = &offset
	return aq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (aq *ApplicationQuery) Unique(unique bool) *ApplicationQuery {
	aq.ctx.Unique = &unique
	return aq
}

// Order specifies how the records should be ordered.
func (aq *ApplicationQuery) Order(o ...OrderFunc) *ApplicationQuery {
	aq.order = append(aq.order, o...)
	return aq
}

// First returns the first Application entity from the query.
// Returns a *NotFoundError when no Application was found.
func (aq *ApplicationQuery) First(ctx context.Context) (*Application, error) {
	nodes, err := aq.Limit(1).All(setContextOp(ctx, aq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{application.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (aq *ApplicationQuery) FirstX(ctx context.Context) *Application {
	node, err := aq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Application ID from the query.
// Returns a *NotFoundError when no Application ID was found.
func (aq *ApplicationQuery) FirstID(ctx context.Context) (id oid.ID, err error) {
	var ids []oid.ID
	if ids, err = aq.Limit(1).IDs(setContextOp(ctx, aq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{application.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (aq *ApplicationQuery) FirstIDX(ctx context.Context) oid.ID {
	id, err := aq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Application entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Application entity is found.
// Returns a *NotFoundError when no Application entities are found.
func (aq *ApplicationQuery) Only(ctx context.Context) (*Application, error) {
	nodes, err := aq.Limit(2).All(setContextOp(ctx, aq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{application.Label}
	default:
		return nil, &NotSingularError{application.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (aq *ApplicationQuery) OnlyX(ctx context.Context) *Application {
	node, err := aq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Application ID in the query.
// Returns a *NotSingularError when more than one Application ID is found.
// Returns a *NotFoundError when no entities are found.
func (aq *ApplicationQuery) OnlyID(ctx context.Context) (id oid.ID, err error) {
	var ids []oid.ID
	if ids, err = aq.Limit(2).IDs(setContextOp(ctx, aq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{application.Label}
	default:
		err = &NotSingularError{application.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (aq *ApplicationQuery) OnlyIDX(ctx context.Context) oid.ID {
	id, err := aq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Applications.
func (aq *ApplicationQuery) All(ctx context.Context) ([]*Application, error) {
	ctx = setContextOp(ctx, aq.ctx, "All")
	if err := aq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Application, *ApplicationQuery]()
	return withInterceptors[[]*Application](ctx, aq, qr, aq.inters)
}

// AllX is like All, but panics if an error occurs.
func (aq *ApplicationQuery) AllX(ctx context.Context) []*Application {
	nodes, err := aq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Application IDs.
func (aq *ApplicationQuery) IDs(ctx context.Context) ([]oid.ID, error) {
	var ids []oid.ID
	ctx = setContextOp(ctx, aq.ctx, "IDs")
	if err := aq.Select(application.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (aq *ApplicationQuery) IDsX(ctx context.Context) []oid.ID {
	ids, err := aq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (aq *ApplicationQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, aq.ctx, "Count")
	if err := aq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, aq, querierCount[*ApplicationQuery](), aq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (aq *ApplicationQuery) CountX(ctx context.Context) int {
	count, err := aq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (aq *ApplicationQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, aq.ctx, "Exist")
	switch _, err := aq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("model: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (aq *ApplicationQuery) ExistX(ctx context.Context) bool {
	exist, err := aq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the ApplicationQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (aq *ApplicationQuery) Clone() *ApplicationQuery {
	if aq == nil {
		return nil
	}
	return &ApplicationQuery{
		config:     aq.config,
		ctx:        aq.ctx.Clone(),
		order:      append([]OrderFunc{}, aq.order...),
		inters:     append([]Interceptor{}, aq.inters...),
		predicates: append([]predicate.Application{}, aq.predicates...),
		// clone intermediate query.
		sql:  aq.sql.Clone(),
		path: aq.path,
	}
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		CreateTime time.Time `json:"createTime,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Application.Query().
//		GroupBy(application.FieldCreateTime).
//		Aggregate(model.Count()).
//		Scan(ctx, &v)
func (aq *ApplicationQuery) GroupBy(field string, fields ...string) *ApplicationGroupBy {
	aq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &ApplicationGroupBy{build: aq}
	grbuild.flds = &aq.ctx.Fields
	grbuild.label = application.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		CreateTime time.Time `json:"createTime,omitempty"`
//	}
//
//	client.Application.Query().
//		Select(application.FieldCreateTime).
//		Scan(ctx, &v)
func (aq *ApplicationQuery) Select(fields ...string) *ApplicationSelect {
	aq.ctx.Fields = append(aq.ctx.Fields, fields...)
	sbuild := &ApplicationSelect{ApplicationQuery: aq}
	sbuild.label = application.Label
	sbuild.flds, sbuild.scan = &aq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a ApplicationSelect configured with the given aggregations.
func (aq *ApplicationQuery) Aggregate(fns ...AggregateFunc) *ApplicationSelect {
	return aq.Select().Aggregate(fns...)
}

func (aq *ApplicationQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range aq.inters {
		if inter == nil {
			return fmt.Errorf("model: uninitialized interceptor (forgotten import model/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, aq); err != nil {
				return err
			}
		}
	}
	for _, f := range aq.ctx.Fields {
		if !application.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("model: invalid field %q for query", f)}
		}
	}
	if aq.path != nil {
		prev, err := aq.path(ctx)
		if err != nil {
			return err
		}
		aq.sql = prev
	}
	return nil
}

func (aq *ApplicationQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Application, error) {
	var (
		nodes = []*Application{}
		_spec = aq.querySpec()
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Application).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Application{config: aq.config}
		nodes = append(nodes, node)
		return node.assignValues(columns, values)
	}
	_spec.Node.Schema = aq.schemaConfig.Application
	ctx = internal.NewSchemaConfigContext(ctx, aq.schemaConfig)
	if len(aq.modifiers) > 0 {
		_spec.Modifiers = aq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, aq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (aq *ApplicationQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := aq.querySpec()
	_spec.Node.Schema = aq.schemaConfig.Application
	ctx = internal.NewSchemaConfigContext(ctx, aq.schemaConfig)
	if len(aq.modifiers) > 0 {
		_spec.Modifiers = aq.modifiers
	}
	_spec.Node.Columns = aq.ctx.Fields
	if len(aq.ctx.Fields) > 0 {
		_spec.Unique = aq.ctx.Unique != nil && *aq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, aq.driver, _spec)
}

func (aq *ApplicationQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   application.Table,
			Columns: application.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeOther,
				Column: application.FieldID,
			},
		},
		From:   aq.sql,
		Unique: true,
	}
	if unique := aq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := aq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, application.FieldID)
		for i := range fields {
			if fields[i] != application.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := aq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := aq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := aq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := aq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (aq *ApplicationQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(aq.driver.Dialect())
	t1 := builder.Table(application.Table)
	columns := aq.ctx.Fields
	if len(columns) == 0 {
		columns = application.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if aq.sql != nil {
		selector = aq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if aq.ctx.Unique != nil && *aq.ctx.Unique {
		selector.Distinct()
	}
	t1.Schema(aq.schemaConfig.Application)
	ctx = internal.NewSchemaConfigContext(ctx, aq.schemaConfig)
	selector.WithContext(ctx)
	for _, m := range aq.modifiers {
		m(selector)
	}
	for _, p := range aq.predicates {
		p(selector)
	}
	for _, p := range aq.order {
		p(selector)
	}
	if offset := aq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := aq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ForUpdate locks the selected rows against concurrent updates, and prevent them from being
// updated, deleted or "selected ... for update" by other sessions, until the transaction is
// either committed or rolled-back.
func (aq *ApplicationQuery) ForUpdate(opts ...sql.LockOption) *ApplicationQuery {
	if aq.driver.Dialect() == dialect.Postgres {
		aq.Unique(false)
	}
	aq.modifiers = append(aq.modifiers, func(s *sql.Selector) {
		s.ForUpdate(opts...)
	})
	return aq
}

// ForShare behaves similarly to ForUpdate, except that it acquires a shared mode lock
// on any rows that are read. Other sessions can read the rows, but cannot modify them
// until your transaction commits.
func (aq *ApplicationQuery) ForShare(opts ...sql.LockOption) *ApplicationQuery {
	if aq.driver.Dialect() == dialect.Postgres {
		aq.Unique(false)
	}
	aq.modifiers = append(aq.modifiers, func(s *sql.Selector) {
		s.ForShare(opts...)
	})
	return aq
}

// Modify adds a query modifier for attaching custom logic to queries.
func (aq *ApplicationQuery) Modify(modifiers ...func(s *sql.Selector)) *ApplicationSelect {
	aq.modifiers = append(aq.modifiers, modifiers...)
	return aq.Select()
}

// ApplicationGroupBy is the group-by builder for Application entities.
type ApplicationGroupBy struct {
	selector
	build *ApplicationQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (agb *ApplicationGroupBy) Aggregate(fns ...AggregateFunc) *ApplicationGroupBy {
	agb.fns = append(agb.fns, fns...)
	return agb
}

// Scan applies the selector query and scans the result into the given value.
func (agb *ApplicationGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, agb.build.ctx, "GroupBy")
	if err := agb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*ApplicationQuery, *ApplicationGroupBy](ctx, agb.build, agb, agb.build.inters, v)
}

func (agb *ApplicationGroupBy) sqlScan(ctx context.Context, root *ApplicationQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(agb.fns))
	for _, fn := range agb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*agb.flds)+len(agb.fns))
		for _, f := range *agb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*agb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := agb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// ApplicationSelect is the builder for selecting fields of Application entities.
type ApplicationSelect struct {
	*ApplicationQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (as *ApplicationSelect) Aggregate(fns ...AggregateFunc) *ApplicationSelect {
	as.fns = append(as.fns, fns...)
	return as
}

// Scan applies the selector query and scans the result into the given value.
func (as *ApplicationSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, as.ctx, "Select")
	if err := as.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*ApplicationQuery, *ApplicationSelect](ctx, as.ApplicationQuery, as, as.inters, v)
}

func (as *ApplicationSelect) sqlScan(ctx context.Context, root *ApplicationQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(as.fns))
	for _, fn := range as.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*as.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := as.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (as *ApplicationSelect) Modify(modifiers ...func(s *sql.Selector)) *ApplicationSelect {
	as.modifiers = append(as.modifiers, modifiers...)
	return as
}
