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
	"github.com/seal-io/seal/pkg/dao/model/module"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/oid"
)

// ModuleQuery is the builder for querying Module entities.
type ModuleQuery struct {
	config
	ctx        *QueryContext
	order      []OrderFunc
	inters     []Interceptor
	predicates []predicate.Module
	modifiers  []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the ModuleQuery builder.
func (mq *ModuleQuery) Where(ps ...predicate.Module) *ModuleQuery {
	mq.predicates = append(mq.predicates, ps...)
	return mq
}

// Limit the number of records to be returned by this query.
func (mq *ModuleQuery) Limit(limit int) *ModuleQuery {
	mq.ctx.Limit = &limit
	return mq
}

// Offset to start from.
func (mq *ModuleQuery) Offset(offset int) *ModuleQuery {
	mq.ctx.Offset = &offset
	return mq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (mq *ModuleQuery) Unique(unique bool) *ModuleQuery {
	mq.ctx.Unique = &unique
	return mq
}

// Order specifies how the records should be ordered.
func (mq *ModuleQuery) Order(o ...OrderFunc) *ModuleQuery {
	mq.order = append(mq.order, o...)
	return mq
}

// First returns the first Module entity from the query.
// Returns a *NotFoundError when no Module was found.
func (mq *ModuleQuery) First(ctx context.Context) (*Module, error) {
	nodes, err := mq.Limit(1).All(setContextOp(ctx, mq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{module.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (mq *ModuleQuery) FirstX(ctx context.Context) *Module {
	node, err := mq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Module ID from the query.
// Returns a *NotFoundError when no Module ID was found.
func (mq *ModuleQuery) FirstID(ctx context.Context) (id oid.ID, err error) {
	var ids []oid.ID
	if ids, err = mq.Limit(1).IDs(setContextOp(ctx, mq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{module.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (mq *ModuleQuery) FirstIDX(ctx context.Context) oid.ID {
	id, err := mq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Module entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Module entity is found.
// Returns a *NotFoundError when no Module entities are found.
func (mq *ModuleQuery) Only(ctx context.Context) (*Module, error) {
	nodes, err := mq.Limit(2).All(setContextOp(ctx, mq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{module.Label}
	default:
		return nil, &NotSingularError{module.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (mq *ModuleQuery) OnlyX(ctx context.Context) *Module {
	node, err := mq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Module ID in the query.
// Returns a *NotSingularError when more than one Module ID is found.
// Returns a *NotFoundError when no entities are found.
func (mq *ModuleQuery) OnlyID(ctx context.Context) (id oid.ID, err error) {
	var ids []oid.ID
	if ids, err = mq.Limit(2).IDs(setContextOp(ctx, mq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{module.Label}
	default:
		err = &NotSingularError{module.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (mq *ModuleQuery) OnlyIDX(ctx context.Context) oid.ID {
	id, err := mq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Modules.
func (mq *ModuleQuery) All(ctx context.Context) ([]*Module, error) {
	ctx = setContextOp(ctx, mq.ctx, "All")
	if err := mq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Module, *ModuleQuery]()
	return withInterceptors[[]*Module](ctx, mq, qr, mq.inters)
}

// AllX is like All, but panics if an error occurs.
func (mq *ModuleQuery) AllX(ctx context.Context) []*Module {
	nodes, err := mq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Module IDs.
func (mq *ModuleQuery) IDs(ctx context.Context) ([]oid.ID, error) {
	var ids []oid.ID
	ctx = setContextOp(ctx, mq.ctx, "IDs")
	if err := mq.Select(module.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (mq *ModuleQuery) IDsX(ctx context.Context) []oid.ID {
	ids, err := mq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (mq *ModuleQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, mq.ctx, "Count")
	if err := mq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, mq, querierCount[*ModuleQuery](), mq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (mq *ModuleQuery) CountX(ctx context.Context) int {
	count, err := mq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (mq *ModuleQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, mq.ctx, "Exist")
	switch _, err := mq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("model: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (mq *ModuleQuery) ExistX(ctx context.Context) bool {
	exist, err := mq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the ModuleQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (mq *ModuleQuery) Clone() *ModuleQuery {
	if mq == nil {
		return nil
	}
	return &ModuleQuery{
		config:     mq.config,
		ctx:        mq.ctx.Clone(),
		order:      append([]OrderFunc{}, mq.order...),
		inters:     append([]Interceptor{}, mq.inters...),
		predicates: append([]predicate.Module{}, mq.predicates...),
		// clone intermediate query.
		sql:  mq.sql.Clone(),
		path: mq.path,
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
//	client.Module.Query().
//		GroupBy(module.FieldStatus).
//		Aggregate(model.Count()).
//		Scan(ctx, &v)
func (mq *ModuleQuery) GroupBy(field string, fields ...string) *ModuleGroupBy {
	mq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &ModuleGroupBy{build: mq}
	grbuild.flds = &mq.ctx.Fields
	grbuild.label = module.Label
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
//	client.Module.Query().
//		Select(module.FieldStatus).
//		Scan(ctx, &v)
func (mq *ModuleQuery) Select(fields ...string) *ModuleSelect {
	mq.ctx.Fields = append(mq.ctx.Fields, fields...)
	sbuild := &ModuleSelect{ModuleQuery: mq}
	sbuild.label = module.Label
	sbuild.flds, sbuild.scan = &mq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a ModuleSelect configured with the given aggregations.
func (mq *ModuleQuery) Aggregate(fns ...AggregateFunc) *ModuleSelect {
	return mq.Select().Aggregate(fns...)
}

func (mq *ModuleQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range mq.inters {
		if inter == nil {
			return fmt.Errorf("model: uninitialized interceptor (forgotten import model/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, mq); err != nil {
				return err
			}
		}
	}
	for _, f := range mq.ctx.Fields {
		if !module.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("model: invalid field %q for query", f)}
		}
	}
	if mq.path != nil {
		prev, err := mq.path(ctx)
		if err != nil {
			return err
		}
		mq.sql = prev
	}
	return nil
}

func (mq *ModuleQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Module, error) {
	var (
		nodes = []*Module{}
		_spec = mq.querySpec()
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Module).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Module{config: mq.config}
		nodes = append(nodes, node)
		return node.assignValues(columns, values)
	}
	_spec.Node.Schema = mq.schemaConfig.Module
	ctx = internal.NewSchemaConfigContext(ctx, mq.schemaConfig)
	if len(mq.modifiers) > 0 {
		_spec.Modifiers = mq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, mq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (mq *ModuleQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := mq.querySpec()
	_spec.Node.Schema = mq.schemaConfig.Module
	ctx = internal.NewSchemaConfigContext(ctx, mq.schemaConfig)
	if len(mq.modifiers) > 0 {
		_spec.Modifiers = mq.modifiers
	}
	_spec.Node.Columns = mq.ctx.Fields
	if len(mq.ctx.Fields) > 0 {
		_spec.Unique = mq.ctx.Unique != nil && *mq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, mq.driver, _spec)
}

func (mq *ModuleQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   module.Table,
			Columns: module.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeOther,
				Column: module.FieldID,
			},
		},
		From:   mq.sql,
		Unique: true,
	}
	if unique := mq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := mq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, module.FieldID)
		for i := range fields {
			if fields[i] != module.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := mq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := mq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := mq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := mq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (mq *ModuleQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(mq.driver.Dialect())
	t1 := builder.Table(module.Table)
	columns := mq.ctx.Fields
	if len(columns) == 0 {
		columns = module.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if mq.sql != nil {
		selector = mq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if mq.ctx.Unique != nil && *mq.ctx.Unique {
		selector.Distinct()
	}
	t1.Schema(mq.schemaConfig.Module)
	ctx = internal.NewSchemaConfigContext(ctx, mq.schemaConfig)
	selector.WithContext(ctx)
	for _, m := range mq.modifiers {
		m(selector)
	}
	for _, p := range mq.predicates {
		p(selector)
	}
	for _, p := range mq.order {
		p(selector)
	}
	if offset := mq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := mq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ForUpdate locks the selected rows against concurrent updates, and prevent them from being
// updated, deleted or "selected ... for update" by other sessions, until the transaction is
// either committed or rolled-back.
func (mq *ModuleQuery) ForUpdate(opts ...sql.LockOption) *ModuleQuery {
	if mq.driver.Dialect() == dialect.Postgres {
		mq.Unique(false)
	}
	mq.modifiers = append(mq.modifiers, func(s *sql.Selector) {
		s.ForUpdate(opts...)
	})
	return mq
}

// ForShare behaves similarly to ForUpdate, except that it acquires a shared mode lock
// on any rows that are read. Other sessions can read the rows, but cannot modify them
// until your transaction commits.
func (mq *ModuleQuery) ForShare(opts ...sql.LockOption) *ModuleQuery {
	if mq.driver.Dialect() == dialect.Postgres {
		mq.Unique(false)
	}
	mq.modifiers = append(mq.modifiers, func(s *sql.Selector) {
		s.ForShare(opts...)
	})
	return mq
}

// Modify adds a query modifier for attaching custom logic to queries.
func (mq *ModuleQuery) Modify(modifiers ...func(s *sql.Selector)) *ModuleSelect {
	mq.modifiers = append(mq.modifiers, modifiers...)
	return mq.Select()
}

// ModuleGroupBy is the group-by builder for Module entities.
type ModuleGroupBy struct {
	selector
	build *ModuleQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (mgb *ModuleGroupBy) Aggregate(fns ...AggregateFunc) *ModuleGroupBy {
	mgb.fns = append(mgb.fns, fns...)
	return mgb
}

// Scan applies the selector query and scans the result into the given value.
func (mgb *ModuleGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, mgb.build.ctx, "GroupBy")
	if err := mgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*ModuleQuery, *ModuleGroupBy](ctx, mgb.build, mgb, mgb.build.inters, v)
}

func (mgb *ModuleGroupBy) sqlScan(ctx context.Context, root *ModuleQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(mgb.fns))
	for _, fn := range mgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*mgb.flds)+len(mgb.fns))
		for _, f := range *mgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*mgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := mgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// ModuleSelect is the builder for selecting fields of Module entities.
type ModuleSelect struct {
	*ModuleQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (ms *ModuleSelect) Aggregate(fns ...AggregateFunc) *ModuleSelect {
	ms.fns = append(ms.fns, fns...)
	return ms
}

// Scan applies the selector query and scans the result into the given value.
func (ms *ModuleSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ms.ctx, "Select")
	if err := ms.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*ModuleQuery, *ModuleSelect](ctx, ms.ModuleQuery, ms, ms.inters, v)
}

func (ms *ModuleSelect) sqlScan(ctx context.Context, root *ModuleQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(ms.fns))
	for _, fn := range ms.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*ms.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ms.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (ms *ModuleSelect) Modify(modifiers ...func(s *sql.Selector)) *ModuleSelect {
	ms.modifiers = append(ms.modifiers, modifiers...)
	return ms
}
