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
	"github.com/seal-io/seal/pkg/dao/model/moduleversion"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/types"
)

// ModuleVersionQuery is the builder for querying ModuleVersion entities.
type ModuleVersionQuery struct {
	config
	ctx        *QueryContext
	order      []OrderFunc
	inters     []Interceptor
	predicates []predicate.ModuleVersion
	withModule *ModuleQuery
	modifiers  []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the ModuleVersionQuery builder.
func (mvq *ModuleVersionQuery) Where(ps ...predicate.ModuleVersion) *ModuleVersionQuery {
	mvq.predicates = append(mvq.predicates, ps...)
	return mvq
}

// Limit the number of records to be returned by this query.
func (mvq *ModuleVersionQuery) Limit(limit int) *ModuleVersionQuery {
	mvq.ctx.Limit = &limit
	return mvq
}

// Offset to start from.
func (mvq *ModuleVersionQuery) Offset(offset int) *ModuleVersionQuery {
	mvq.ctx.Offset = &offset
	return mvq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (mvq *ModuleVersionQuery) Unique(unique bool) *ModuleVersionQuery {
	mvq.ctx.Unique = &unique
	return mvq
}

// Order specifies how the records should be ordered.
func (mvq *ModuleVersionQuery) Order(o ...OrderFunc) *ModuleVersionQuery {
	mvq.order = append(mvq.order, o...)
	return mvq
}

// QueryModule chains the current query on the "module" edge.
func (mvq *ModuleVersionQuery) QueryModule() *ModuleQuery {
	query := (&ModuleClient{config: mvq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := mvq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := mvq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(moduleversion.Table, moduleversion.FieldID, selector),
			sqlgraph.To(module.Table, module.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, moduleversion.ModuleTable, moduleversion.ModuleColumn),
		)
		schemaConfig := mvq.schemaConfig
		step.To.Schema = schemaConfig.Module
		step.Edge.Schema = schemaConfig.ModuleVersion
		fromU = sqlgraph.SetNeighbors(mvq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first ModuleVersion entity from the query.
// Returns a *NotFoundError when no ModuleVersion was found.
func (mvq *ModuleVersionQuery) First(ctx context.Context) (*ModuleVersion, error) {
	nodes, err := mvq.Limit(1).All(setContextOp(ctx, mvq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{moduleversion.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (mvq *ModuleVersionQuery) FirstX(ctx context.Context) *ModuleVersion {
	node, err := mvq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first ModuleVersion ID from the query.
// Returns a *NotFoundError when no ModuleVersion ID was found.
func (mvq *ModuleVersionQuery) FirstID(ctx context.Context) (id types.ID, err error) {
	var ids []types.ID
	if ids, err = mvq.Limit(1).IDs(setContextOp(ctx, mvq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{moduleversion.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (mvq *ModuleVersionQuery) FirstIDX(ctx context.Context) types.ID {
	id, err := mvq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single ModuleVersion entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one ModuleVersion entity is found.
// Returns a *NotFoundError when no ModuleVersion entities are found.
func (mvq *ModuleVersionQuery) Only(ctx context.Context) (*ModuleVersion, error) {
	nodes, err := mvq.Limit(2).All(setContextOp(ctx, mvq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{moduleversion.Label}
	default:
		return nil, &NotSingularError{moduleversion.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (mvq *ModuleVersionQuery) OnlyX(ctx context.Context) *ModuleVersion {
	node, err := mvq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only ModuleVersion ID in the query.
// Returns a *NotSingularError when more than one ModuleVersion ID is found.
// Returns a *NotFoundError when no entities are found.
func (mvq *ModuleVersionQuery) OnlyID(ctx context.Context) (id types.ID, err error) {
	var ids []types.ID
	if ids, err = mvq.Limit(2).IDs(setContextOp(ctx, mvq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{moduleversion.Label}
	default:
		err = &NotSingularError{moduleversion.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (mvq *ModuleVersionQuery) OnlyIDX(ctx context.Context) types.ID {
	id, err := mvq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of ModuleVersions.
func (mvq *ModuleVersionQuery) All(ctx context.Context) ([]*ModuleVersion, error) {
	ctx = setContextOp(ctx, mvq.ctx, "All")
	if err := mvq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*ModuleVersion, *ModuleVersionQuery]()
	return withInterceptors[[]*ModuleVersion](ctx, mvq, qr, mvq.inters)
}

// AllX is like All, but panics if an error occurs.
func (mvq *ModuleVersionQuery) AllX(ctx context.Context) []*ModuleVersion {
	nodes, err := mvq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of ModuleVersion IDs.
func (mvq *ModuleVersionQuery) IDs(ctx context.Context) (ids []types.ID, err error) {
	if mvq.ctx.Unique == nil && mvq.path != nil {
		mvq.Unique(true)
	}
	ctx = setContextOp(ctx, mvq.ctx, "IDs")
	if err = mvq.Select(moduleversion.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (mvq *ModuleVersionQuery) IDsX(ctx context.Context) []types.ID {
	ids, err := mvq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (mvq *ModuleVersionQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, mvq.ctx, "Count")
	if err := mvq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, mvq, querierCount[*ModuleVersionQuery](), mvq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (mvq *ModuleVersionQuery) CountX(ctx context.Context) int {
	count, err := mvq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (mvq *ModuleVersionQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, mvq.ctx, "Exist")
	switch _, err := mvq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("model: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (mvq *ModuleVersionQuery) ExistX(ctx context.Context) bool {
	exist, err := mvq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the ModuleVersionQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (mvq *ModuleVersionQuery) Clone() *ModuleVersionQuery {
	if mvq == nil {
		return nil
	}
	return &ModuleVersionQuery{
		config:     mvq.config,
		ctx:        mvq.ctx.Clone(),
		order:      append([]OrderFunc{}, mvq.order...),
		inters:     append([]Interceptor{}, mvq.inters...),
		predicates: append([]predicate.ModuleVersion{}, mvq.predicates...),
		withModule: mvq.withModule.Clone(),
		// clone intermediate query.
		sql:  mvq.sql.Clone(),
		path: mvq.path,
	}
}

// WithModule tells the query-builder to eager-load the nodes that are connected to
// the "module" edge. The optional arguments are used to configure the query builder of the edge.
func (mvq *ModuleVersionQuery) WithModule(opts ...func(*ModuleQuery)) *ModuleVersionQuery {
	query := (&ModuleClient{config: mvq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	mvq.withModule = query
	return mvq
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
//	client.ModuleVersion.Query().
//		GroupBy(moduleversion.FieldCreateTime).
//		Aggregate(model.Count()).
//		Scan(ctx, &v)
func (mvq *ModuleVersionQuery) GroupBy(field string, fields ...string) *ModuleVersionGroupBy {
	mvq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &ModuleVersionGroupBy{build: mvq}
	grbuild.flds = &mvq.ctx.Fields
	grbuild.label = moduleversion.Label
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
//	client.ModuleVersion.Query().
//		Select(moduleversion.FieldCreateTime).
//		Scan(ctx, &v)
func (mvq *ModuleVersionQuery) Select(fields ...string) *ModuleVersionSelect {
	mvq.ctx.Fields = append(mvq.ctx.Fields, fields...)
	sbuild := &ModuleVersionSelect{ModuleVersionQuery: mvq}
	sbuild.label = moduleversion.Label
	sbuild.flds, sbuild.scan = &mvq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a ModuleVersionSelect configured with the given aggregations.
func (mvq *ModuleVersionQuery) Aggregate(fns ...AggregateFunc) *ModuleVersionSelect {
	return mvq.Select().Aggregate(fns...)
}

func (mvq *ModuleVersionQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range mvq.inters {
		if inter == nil {
			return fmt.Errorf("model: uninitialized interceptor (forgotten import model/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, mvq); err != nil {
				return err
			}
		}
	}
	for _, f := range mvq.ctx.Fields {
		if !moduleversion.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("model: invalid field %q for query", f)}
		}
	}
	if mvq.path != nil {
		prev, err := mvq.path(ctx)
		if err != nil {
			return err
		}
		mvq.sql = prev
	}
	return nil
}

func (mvq *ModuleVersionQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*ModuleVersion, error) {
	var (
		nodes       = []*ModuleVersion{}
		_spec       = mvq.querySpec()
		loadedTypes = [1]bool{
			mvq.withModule != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*ModuleVersion).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &ModuleVersion{config: mvq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	_spec.Node.Schema = mvq.schemaConfig.ModuleVersion
	ctx = internal.NewSchemaConfigContext(ctx, mvq.schemaConfig)
	if len(mvq.modifiers) > 0 {
		_spec.Modifiers = mvq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, mvq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := mvq.withModule; query != nil {
		if err := mvq.loadModule(ctx, query, nodes, nil,
			func(n *ModuleVersion, e *Module) { n.Edges.Module = e }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (mvq *ModuleVersionQuery) loadModule(ctx context.Context, query *ModuleQuery, nodes []*ModuleVersion, init func(*ModuleVersion), assign func(*ModuleVersion, *Module)) error {
	ids := make([]string, 0, len(nodes))
	nodeids := make(map[string][]*ModuleVersion)
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
			return fmt.Errorf(`unexpected foreign-key "moduleID" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}

func (mvq *ModuleVersionQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := mvq.querySpec()
	_spec.Node.Schema = mvq.schemaConfig.ModuleVersion
	ctx = internal.NewSchemaConfigContext(ctx, mvq.schemaConfig)
	if len(mvq.modifiers) > 0 {
		_spec.Modifiers = mvq.modifiers
	}
	_spec.Node.Columns = mvq.ctx.Fields
	if len(mvq.ctx.Fields) > 0 {
		_spec.Unique = mvq.ctx.Unique != nil && *mvq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, mvq.driver, _spec)
}

func (mvq *ModuleVersionQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(moduleversion.Table, moduleversion.Columns, sqlgraph.NewFieldSpec(moduleversion.FieldID, field.TypeString))
	_spec.From = mvq.sql
	if unique := mvq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if mvq.path != nil {
		_spec.Unique = true
	}
	if fields := mvq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, moduleversion.FieldID)
		for i := range fields {
			if fields[i] != moduleversion.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := mvq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := mvq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := mvq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := mvq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (mvq *ModuleVersionQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(mvq.driver.Dialect())
	t1 := builder.Table(moduleversion.Table)
	columns := mvq.ctx.Fields
	if len(columns) == 0 {
		columns = moduleversion.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if mvq.sql != nil {
		selector = mvq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if mvq.ctx.Unique != nil && *mvq.ctx.Unique {
		selector.Distinct()
	}
	t1.Schema(mvq.schemaConfig.ModuleVersion)
	ctx = internal.NewSchemaConfigContext(ctx, mvq.schemaConfig)
	selector.WithContext(ctx)
	for _, m := range mvq.modifiers {
		m(selector)
	}
	for _, p := range mvq.predicates {
		p(selector)
	}
	for _, p := range mvq.order {
		p(selector)
	}
	if offset := mvq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := mvq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ForUpdate locks the selected rows against concurrent updates, and prevent them from being
// updated, deleted or "selected ... for update" by other sessions, until the transaction is
// either committed or rolled-back.
func (mvq *ModuleVersionQuery) ForUpdate(opts ...sql.LockOption) *ModuleVersionQuery {
	if mvq.driver.Dialect() == dialect.Postgres {
		mvq.Unique(false)
	}
	mvq.modifiers = append(mvq.modifiers, func(s *sql.Selector) {
		s.ForUpdate(opts...)
	})
	return mvq
}

// ForShare behaves similarly to ForUpdate, except that it acquires a shared mode lock
// on any rows that are read. Other sessions can read the rows, but cannot modify them
// until your transaction commits.
func (mvq *ModuleVersionQuery) ForShare(opts ...sql.LockOption) *ModuleVersionQuery {
	if mvq.driver.Dialect() == dialect.Postgres {
		mvq.Unique(false)
	}
	mvq.modifiers = append(mvq.modifiers, func(s *sql.Selector) {
		s.ForShare(opts...)
	})
	return mvq
}

// Modify adds a query modifier for attaching custom logic to queries.
func (mvq *ModuleVersionQuery) Modify(modifiers ...func(s *sql.Selector)) *ModuleVersionSelect {
	mvq.modifiers = append(mvq.modifiers, modifiers...)
	return mvq.Select()
}

// ModuleVersionGroupBy is the group-by builder for ModuleVersion entities.
type ModuleVersionGroupBy struct {
	selector
	build *ModuleVersionQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (mvgb *ModuleVersionGroupBy) Aggregate(fns ...AggregateFunc) *ModuleVersionGroupBy {
	mvgb.fns = append(mvgb.fns, fns...)
	return mvgb
}

// Scan applies the selector query and scans the result into the given value.
func (mvgb *ModuleVersionGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, mvgb.build.ctx, "GroupBy")
	if err := mvgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*ModuleVersionQuery, *ModuleVersionGroupBy](ctx, mvgb.build, mvgb, mvgb.build.inters, v)
}

func (mvgb *ModuleVersionGroupBy) sqlScan(ctx context.Context, root *ModuleVersionQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(mvgb.fns))
	for _, fn := range mvgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*mvgb.flds)+len(mvgb.fns))
		for _, f := range *mvgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*mvgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := mvgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// ModuleVersionSelect is the builder for selecting fields of ModuleVersion entities.
type ModuleVersionSelect struct {
	*ModuleVersionQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (mvs *ModuleVersionSelect) Aggregate(fns ...AggregateFunc) *ModuleVersionSelect {
	mvs.fns = append(mvs.fns, fns...)
	return mvs
}

// Scan applies the selector query and scans the result into the given value.
func (mvs *ModuleVersionSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, mvs.ctx, "Select")
	if err := mvs.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*ModuleVersionQuery, *ModuleVersionSelect](ctx, mvs.ModuleVersionQuery, mvs, mvs.inters, v)
}

func (mvs *ModuleVersionSelect) sqlScan(ctx context.Context, root *ModuleVersionQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(mvs.fns))
	for _, fn := range mvs.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*mvs.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := mvs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (mvs *ModuleVersionSelect) Modify(modifiers ...func(s *sql.Selector)) *ModuleVersionSelect {
	mvs.modifiers = append(mvs.modifiers, modifiers...)
	return mvs
}
