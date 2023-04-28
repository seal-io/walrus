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
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/subject"
	"github.com/seal-io/seal/pkg/dao/types/oid"
)

// SubjectQuery is the builder for querying Subject entities.
type SubjectQuery struct {
	config
	ctx        *QueryContext
	order      []subject.OrderOption
	inters     []Interceptor
	predicates []predicate.Subject
	modifiers  []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the SubjectQuery builder.
func (sq *SubjectQuery) Where(ps ...predicate.Subject) *SubjectQuery {
	sq.predicates = append(sq.predicates, ps...)
	return sq
}

// Limit the number of records to be returned by this query.
func (sq *SubjectQuery) Limit(limit int) *SubjectQuery {
	sq.ctx.Limit = &limit
	return sq
}

// Offset to start from.
func (sq *SubjectQuery) Offset(offset int) *SubjectQuery {
	sq.ctx.Offset = &offset
	return sq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (sq *SubjectQuery) Unique(unique bool) *SubjectQuery {
	sq.ctx.Unique = &unique
	return sq
}

// Order specifies how the records should be ordered.
func (sq *SubjectQuery) Order(o ...subject.OrderOption) *SubjectQuery {
	sq.order = append(sq.order, o...)
	return sq
}

// First returns the first Subject entity from the query.
// Returns a *NotFoundError when no Subject was found.
func (sq *SubjectQuery) First(ctx context.Context) (*Subject, error) {
	nodes, err := sq.Limit(1).All(setContextOp(ctx, sq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{subject.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (sq *SubjectQuery) FirstX(ctx context.Context) *Subject {
	node, err := sq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Subject ID from the query.
// Returns a *NotFoundError when no Subject ID was found.
func (sq *SubjectQuery) FirstID(ctx context.Context) (id oid.ID, err error) {
	var ids []oid.ID
	if ids, err = sq.Limit(1).IDs(setContextOp(ctx, sq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{subject.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (sq *SubjectQuery) FirstIDX(ctx context.Context) oid.ID {
	id, err := sq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Subject entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Subject entity is found.
// Returns a *NotFoundError when no Subject entities are found.
func (sq *SubjectQuery) Only(ctx context.Context) (*Subject, error) {
	nodes, err := sq.Limit(2).All(setContextOp(ctx, sq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{subject.Label}
	default:
		return nil, &NotSingularError{subject.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (sq *SubjectQuery) OnlyX(ctx context.Context) *Subject {
	node, err := sq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Subject ID in the query.
// Returns a *NotSingularError when more than one Subject ID is found.
// Returns a *NotFoundError when no entities are found.
func (sq *SubjectQuery) OnlyID(ctx context.Context) (id oid.ID, err error) {
	var ids []oid.ID
	if ids, err = sq.Limit(2).IDs(setContextOp(ctx, sq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{subject.Label}
	default:
		err = &NotSingularError{subject.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (sq *SubjectQuery) OnlyIDX(ctx context.Context) oid.ID {
	id, err := sq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Subjects.
func (sq *SubjectQuery) All(ctx context.Context) ([]*Subject, error) {
	ctx = setContextOp(ctx, sq.ctx, "All")
	if err := sq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Subject, *SubjectQuery]()
	return withInterceptors[[]*Subject](ctx, sq, qr, sq.inters)
}

// AllX is like All, but panics if an error occurs.
func (sq *SubjectQuery) AllX(ctx context.Context) []*Subject {
	nodes, err := sq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Subject IDs.
func (sq *SubjectQuery) IDs(ctx context.Context) (ids []oid.ID, err error) {
	if sq.ctx.Unique == nil && sq.path != nil {
		sq.Unique(true)
	}
	ctx = setContextOp(ctx, sq.ctx, "IDs")
	if err = sq.Select(subject.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (sq *SubjectQuery) IDsX(ctx context.Context) []oid.ID {
	ids, err := sq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (sq *SubjectQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, sq.ctx, "Count")
	if err := sq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, sq, querierCount[*SubjectQuery](), sq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (sq *SubjectQuery) CountX(ctx context.Context) int {
	count, err := sq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (sq *SubjectQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, sq.ctx, "Exist")
	switch _, err := sq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("model: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (sq *SubjectQuery) ExistX(ctx context.Context) bool {
	exist, err := sq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the SubjectQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (sq *SubjectQuery) Clone() *SubjectQuery {
	if sq == nil {
		return nil
	}
	return &SubjectQuery{
		config:     sq.config,
		ctx:        sq.ctx.Clone(),
		order:      append([]subject.OrderOption{}, sq.order...),
		inters:     append([]Interceptor{}, sq.inters...),
		predicates: append([]predicate.Subject{}, sq.predicates...),
		// clone intermediate query.
		sql:  sq.sql.Clone(),
		path: sq.path,
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
//	client.Subject.Query().
//		GroupBy(subject.FieldCreateTime).
//		Aggregate(model.Count()).
//		Scan(ctx, &v)
func (sq *SubjectQuery) GroupBy(field string, fields ...string) *SubjectGroupBy {
	sq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &SubjectGroupBy{build: sq}
	grbuild.flds = &sq.ctx.Fields
	grbuild.label = subject.Label
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
//	client.Subject.Query().
//		Select(subject.FieldCreateTime).
//		Scan(ctx, &v)
func (sq *SubjectQuery) Select(fields ...string) *SubjectSelect {
	sq.ctx.Fields = append(sq.ctx.Fields, fields...)
	sbuild := &SubjectSelect{SubjectQuery: sq}
	sbuild.label = subject.Label
	sbuild.flds, sbuild.scan = &sq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a SubjectSelect configured with the given aggregations.
func (sq *SubjectQuery) Aggregate(fns ...AggregateFunc) *SubjectSelect {
	return sq.Select().Aggregate(fns...)
}

func (sq *SubjectQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range sq.inters {
		if inter == nil {
			return fmt.Errorf("model: uninitialized interceptor (forgotten import model/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, sq); err != nil {
				return err
			}
		}
	}
	for _, f := range sq.ctx.Fields {
		if !subject.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("model: invalid field %q for query", f)}
		}
	}
	if sq.path != nil {
		prev, err := sq.path(ctx)
		if err != nil {
			return err
		}
		sq.sql = prev
	}
	return nil
}

func (sq *SubjectQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Subject, error) {
	var (
		nodes = []*Subject{}
		_spec = sq.querySpec()
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Subject).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Subject{config: sq.config}
		nodes = append(nodes, node)
		return node.assignValues(columns, values)
	}
	_spec.Node.Schema = sq.schemaConfig.Subject
	ctx = internal.NewSchemaConfigContext(ctx, sq.schemaConfig)
	if len(sq.modifiers) > 0 {
		_spec.Modifiers = sq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, sq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (sq *SubjectQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := sq.querySpec()
	_spec.Node.Schema = sq.schemaConfig.Subject
	ctx = internal.NewSchemaConfigContext(ctx, sq.schemaConfig)
	if len(sq.modifiers) > 0 {
		_spec.Modifiers = sq.modifiers
	}
	_spec.Node.Columns = sq.ctx.Fields
	if len(sq.ctx.Fields) > 0 {
		_spec.Unique = sq.ctx.Unique != nil && *sq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, sq.driver, _spec)
}

func (sq *SubjectQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(subject.Table, subject.Columns, sqlgraph.NewFieldSpec(subject.FieldID, field.TypeString))
	_spec.From = sq.sql
	if unique := sq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if sq.path != nil {
		_spec.Unique = true
	}
	if fields := sq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, subject.FieldID)
		for i := range fields {
			if fields[i] != subject.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := sq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := sq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := sq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := sq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (sq *SubjectQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(sq.driver.Dialect())
	t1 := builder.Table(subject.Table)
	columns := sq.ctx.Fields
	if len(columns) == 0 {
		columns = subject.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if sq.sql != nil {
		selector = sq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if sq.ctx.Unique != nil && *sq.ctx.Unique {
		selector.Distinct()
	}
	t1.Schema(sq.schemaConfig.Subject)
	ctx = internal.NewSchemaConfigContext(ctx, sq.schemaConfig)
	selector.WithContext(ctx)
	for _, m := range sq.modifiers {
		m(selector)
	}
	for _, p := range sq.predicates {
		p(selector)
	}
	for _, p := range sq.order {
		p(selector)
	}
	if offset := sq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := sq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ForUpdate locks the selected rows against concurrent updates, and prevent them from being
// updated, deleted or "selected ... for update" by other sessions, until the transaction is
// either committed or rolled-back.
func (sq *SubjectQuery) ForUpdate(opts ...sql.LockOption) *SubjectQuery {
	if sq.driver.Dialect() == dialect.Postgres {
		sq.Unique(false)
	}
	sq.modifiers = append(sq.modifiers, func(s *sql.Selector) {
		s.ForUpdate(opts...)
	})
	return sq
}

// ForShare behaves similarly to ForUpdate, except that it acquires a shared mode lock
// on any rows that are read. Other sessions can read the rows, but cannot modify them
// until your transaction commits.
func (sq *SubjectQuery) ForShare(opts ...sql.LockOption) *SubjectQuery {
	if sq.driver.Dialect() == dialect.Postgres {
		sq.Unique(false)
	}
	sq.modifiers = append(sq.modifiers, func(s *sql.Selector) {
		s.ForShare(opts...)
	})
	return sq
}

// Modify adds a query modifier for attaching custom logic to queries.
func (sq *SubjectQuery) Modify(modifiers ...func(s *sql.Selector)) *SubjectSelect {
	sq.modifiers = append(sq.modifiers, modifiers...)
	return sq.Select()
}

// WhereP appends storage-level predicates to the SubjectQuery builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (sq *SubjectQuery) WhereP(ps ...func(*sql.Selector)) {
	var wps = make([]predicate.Subject, 0, len(ps))
	for i := 0; i < len(ps); i++ {
		wps = append(wps, predicate.Subject(ps[i]))
	}
	sq.predicates = append(sq.predicates, wps...)
}

// SubjectGroupBy is the group-by builder for Subject entities.
type SubjectGroupBy struct {
	selector
	build *SubjectQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (sgb *SubjectGroupBy) Aggregate(fns ...AggregateFunc) *SubjectGroupBy {
	sgb.fns = append(sgb.fns, fns...)
	return sgb
}

// Scan applies the selector query and scans the result into the given value.
func (sgb *SubjectGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, sgb.build.ctx, "GroupBy")
	if err := sgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*SubjectQuery, *SubjectGroupBy](ctx, sgb.build, sgb, sgb.build.inters, v)
}

func (sgb *SubjectGroupBy) sqlScan(ctx context.Context, root *SubjectQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(sgb.fns))
	for _, fn := range sgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*sgb.flds)+len(sgb.fns))
		for _, f := range *sgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*sgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := sgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// SubjectSelect is the builder for selecting fields of Subject entities.
type SubjectSelect struct {
	*SubjectQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (ss *SubjectSelect) Aggregate(fns ...AggregateFunc) *SubjectSelect {
	ss.fns = append(ss.fns, fns...)
	return ss
}

// Scan applies the selector query and scans the result into the given value.
func (ss *SubjectSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ss.ctx, "Select")
	if err := ss.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*SubjectQuery, *SubjectSelect](ctx, ss.SubjectQuery, ss, ss.inters, v)
}

func (ss *SubjectSelect) sqlScan(ctx context.Context, root *SubjectQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(ss.fns))
	for _, fn := range ss.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*ss.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ss.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (ss *SubjectSelect) Modify(modifiers ...func(s *sql.Selector)) *SubjectSelect {
	ss.modifiers = append(ss.modifiers, modifiers...)
	return ss
}
