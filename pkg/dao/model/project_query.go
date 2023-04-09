// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package model

import (
	"context"
	"database/sql/driver"
	"fmt"
	"math"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/project"
	"github.com/seal-io/seal/pkg/dao/model/secret"
	"github.com/seal-io/seal/pkg/dao/types/oid"
)

// ProjectQuery is the builder for querying Project entities.
type ProjectQuery struct {
	config
	ctx              *QueryContext
	order            []OrderFunc
	inters           []Interceptor
	predicates       []predicate.Project
	withApplications *ApplicationQuery
	withSecrets      *SecretQuery
	modifiers        []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the ProjectQuery builder.
func (pq *ProjectQuery) Where(ps ...predicate.Project) *ProjectQuery {
	pq.predicates = append(pq.predicates, ps...)
	return pq
}

// Limit the number of records to be returned by this query.
func (pq *ProjectQuery) Limit(limit int) *ProjectQuery {
	pq.ctx.Limit = &limit
	return pq
}

// Offset to start from.
func (pq *ProjectQuery) Offset(offset int) *ProjectQuery {
	pq.ctx.Offset = &offset
	return pq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (pq *ProjectQuery) Unique(unique bool) *ProjectQuery {
	pq.ctx.Unique = &unique
	return pq
}

// Order specifies how the records should be ordered.
func (pq *ProjectQuery) Order(o ...OrderFunc) *ProjectQuery {
	pq.order = append(pq.order, o...)
	return pq
}

// QueryApplications chains the current query on the "applications" edge.
func (pq *ProjectQuery) QueryApplications() *ApplicationQuery {
	query := (&ApplicationClient{config: pq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := pq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := pq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(project.Table, project.FieldID, selector),
			sqlgraph.To(application.Table, application.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, project.ApplicationsTable, project.ApplicationsColumn),
		)
		schemaConfig := pq.schemaConfig
		step.To.Schema = schemaConfig.Application
		step.Edge.Schema = schemaConfig.Application
		fromU = sqlgraph.SetNeighbors(pq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QuerySecrets chains the current query on the "secrets" edge.
func (pq *ProjectQuery) QuerySecrets() *SecretQuery {
	query := (&SecretClient{config: pq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := pq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := pq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(project.Table, project.FieldID, selector),
			sqlgraph.To(secret.Table, secret.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, project.SecretsTable, project.SecretsColumn),
		)
		schemaConfig := pq.schemaConfig
		step.To.Schema = schemaConfig.Secret
		step.Edge.Schema = schemaConfig.Secret
		fromU = sqlgraph.SetNeighbors(pq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Project entity from the query.
// Returns a *NotFoundError when no Project was found.
func (pq *ProjectQuery) First(ctx context.Context) (*Project, error) {
	nodes, err := pq.Limit(1).All(setContextOp(ctx, pq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{project.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (pq *ProjectQuery) FirstX(ctx context.Context) *Project {
	node, err := pq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Project ID from the query.
// Returns a *NotFoundError when no Project ID was found.
func (pq *ProjectQuery) FirstID(ctx context.Context) (id oid.ID, err error) {
	var ids []oid.ID
	if ids, err = pq.Limit(1).IDs(setContextOp(ctx, pq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{project.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (pq *ProjectQuery) FirstIDX(ctx context.Context) oid.ID {
	id, err := pq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Project entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Project entity is found.
// Returns a *NotFoundError when no Project entities are found.
func (pq *ProjectQuery) Only(ctx context.Context) (*Project, error) {
	nodes, err := pq.Limit(2).All(setContextOp(ctx, pq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{project.Label}
	default:
		return nil, &NotSingularError{project.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (pq *ProjectQuery) OnlyX(ctx context.Context) *Project {
	node, err := pq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Project ID in the query.
// Returns a *NotSingularError when more than one Project ID is found.
// Returns a *NotFoundError when no entities are found.
func (pq *ProjectQuery) OnlyID(ctx context.Context) (id oid.ID, err error) {
	var ids []oid.ID
	if ids, err = pq.Limit(2).IDs(setContextOp(ctx, pq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{project.Label}
	default:
		err = &NotSingularError{project.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (pq *ProjectQuery) OnlyIDX(ctx context.Context) oid.ID {
	id, err := pq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Projects.
func (pq *ProjectQuery) All(ctx context.Context) ([]*Project, error) {
	ctx = setContextOp(ctx, pq.ctx, "All")
	if err := pq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Project, *ProjectQuery]()
	return withInterceptors[[]*Project](ctx, pq, qr, pq.inters)
}

// AllX is like All, but panics if an error occurs.
func (pq *ProjectQuery) AllX(ctx context.Context) []*Project {
	nodes, err := pq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Project IDs.
func (pq *ProjectQuery) IDs(ctx context.Context) (ids []oid.ID, err error) {
	if pq.ctx.Unique == nil && pq.path != nil {
		pq.Unique(true)
	}
	ctx = setContextOp(ctx, pq.ctx, "IDs")
	if err = pq.Select(project.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (pq *ProjectQuery) IDsX(ctx context.Context) []oid.ID {
	ids, err := pq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (pq *ProjectQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, pq.ctx, "Count")
	if err := pq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, pq, querierCount[*ProjectQuery](), pq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (pq *ProjectQuery) CountX(ctx context.Context) int {
	count, err := pq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (pq *ProjectQuery) Exist(ctx context.Context) (bool, error) {
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
func (pq *ProjectQuery) ExistX(ctx context.Context) bool {
	exist, err := pq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the ProjectQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (pq *ProjectQuery) Clone() *ProjectQuery {
	if pq == nil {
		return nil
	}
	return &ProjectQuery{
		config:           pq.config,
		ctx:              pq.ctx.Clone(),
		order:            append([]OrderFunc{}, pq.order...),
		inters:           append([]Interceptor{}, pq.inters...),
		predicates:       append([]predicate.Project{}, pq.predicates...),
		withApplications: pq.withApplications.Clone(),
		withSecrets:      pq.withSecrets.Clone(),
		// clone intermediate query.
		sql:  pq.sql.Clone(),
		path: pq.path,
	}
}

// WithApplications tells the query-builder to eager-load the nodes that are connected to
// the "applications" edge. The optional arguments are used to configure the query builder of the edge.
func (pq *ProjectQuery) WithApplications(opts ...func(*ApplicationQuery)) *ProjectQuery {
	query := (&ApplicationClient{config: pq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	pq.withApplications = query
	return pq
}

// WithSecrets tells the query-builder to eager-load the nodes that are connected to
// the "secrets" edge. The optional arguments are used to configure the query builder of the edge.
func (pq *ProjectQuery) WithSecrets(opts ...func(*SecretQuery)) *ProjectQuery {
	query := (&SecretClient{config: pq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	pq.withSecrets = query
	return pq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Name string `json:"name,omitempty" sql:"name"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Project.Query().
//		GroupBy(project.FieldName).
//		Aggregate(model.Count()).
//		Scan(ctx, &v)
func (pq *ProjectQuery) GroupBy(field string, fields ...string) *ProjectGroupBy {
	pq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &ProjectGroupBy{build: pq}
	grbuild.flds = &pq.ctx.Fields
	grbuild.label = project.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Name string `json:"name,omitempty" sql:"name"`
//	}
//
//	client.Project.Query().
//		Select(project.FieldName).
//		Scan(ctx, &v)
func (pq *ProjectQuery) Select(fields ...string) *ProjectSelect {
	pq.ctx.Fields = append(pq.ctx.Fields, fields...)
	sbuild := &ProjectSelect{ProjectQuery: pq}
	sbuild.label = project.Label
	sbuild.flds, sbuild.scan = &pq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a ProjectSelect configured with the given aggregations.
func (pq *ProjectQuery) Aggregate(fns ...AggregateFunc) *ProjectSelect {
	return pq.Select().Aggregate(fns...)
}

func (pq *ProjectQuery) prepareQuery(ctx context.Context) error {
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
		if !project.ValidColumn(f) {
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

func (pq *ProjectQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Project, error) {
	var (
		nodes       = []*Project{}
		_spec       = pq.querySpec()
		loadedTypes = [2]bool{
			pq.withApplications != nil,
			pq.withSecrets != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Project).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Project{config: pq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	_spec.Node.Schema = pq.schemaConfig.Project
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
	if query := pq.withApplications; query != nil {
		if err := pq.loadApplications(ctx, query, nodes,
			func(n *Project) { n.Edges.Applications = []*Application{} },
			func(n *Project, e *Application) { n.Edges.Applications = append(n.Edges.Applications, e) }); err != nil {
			return nil, err
		}
	}
	if query := pq.withSecrets; query != nil {
		if err := pq.loadSecrets(ctx, query, nodes,
			func(n *Project) { n.Edges.Secrets = []*Secret{} },
			func(n *Project, e *Secret) { n.Edges.Secrets = append(n.Edges.Secrets, e) }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (pq *ProjectQuery) loadApplications(ctx context.Context, query *ApplicationQuery, nodes []*Project, init func(*Project), assign func(*Project, *Application)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[oid.ID]*Project)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.Where(predicate.Application(func(s *sql.Selector) {
		s.Where(sql.InValues(project.ApplicationsColumn, fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.ProjectID
		node, ok := nodeids[fk]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "projectID" returned %v for node %v`, fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}
func (pq *ProjectQuery) loadSecrets(ctx context.Context, query *SecretQuery, nodes []*Project, init func(*Project), assign func(*Project, *Secret)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[oid.ID]*Project)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.Where(predicate.Secret(func(s *sql.Selector) {
		s.Where(sql.InValues(project.SecretsColumn, fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.ProjectID
		node, ok := nodeids[fk]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "projectID" returned %v for node %v`, fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}

func (pq *ProjectQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := pq.querySpec()
	_spec.Node.Schema = pq.schemaConfig.Project
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

func (pq *ProjectQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(project.Table, project.Columns, sqlgraph.NewFieldSpec(project.FieldID, field.TypeString))
	_spec.From = pq.sql
	if unique := pq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if pq.path != nil {
		_spec.Unique = true
	}
	if fields := pq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, project.FieldID)
		for i := range fields {
			if fields[i] != project.FieldID {
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

func (pq *ProjectQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(pq.driver.Dialect())
	t1 := builder.Table(project.Table)
	columns := pq.ctx.Fields
	if len(columns) == 0 {
		columns = project.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if pq.sql != nil {
		selector = pq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if pq.ctx.Unique != nil && *pq.ctx.Unique {
		selector.Distinct()
	}
	t1.Schema(pq.schemaConfig.Project)
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
func (pq *ProjectQuery) ForUpdate(opts ...sql.LockOption) *ProjectQuery {
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
func (pq *ProjectQuery) ForShare(opts ...sql.LockOption) *ProjectQuery {
	if pq.driver.Dialect() == dialect.Postgres {
		pq.Unique(false)
	}
	pq.modifiers = append(pq.modifiers, func(s *sql.Selector) {
		s.ForShare(opts...)
	})
	return pq
}

// Modify adds a query modifier for attaching custom logic to queries.
func (pq *ProjectQuery) Modify(modifiers ...func(s *sql.Selector)) *ProjectSelect {
	pq.modifiers = append(pq.modifiers, modifiers...)
	return pq.Select()
}

// WhereP appends storage-level predicates to the ProjectQuery builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (pq *ProjectQuery) WhereP(ps ...func(*sql.Selector)) {
	var wps = make([]predicate.Project, 0, len(ps))
	for i := 0; i < len(ps); i++ {
		wps = append(wps, predicate.Project(ps[i]))
	}
	pq.predicates = append(pq.predicates, wps...)
}

// ProjectGroupBy is the group-by builder for Project entities.
type ProjectGroupBy struct {
	selector
	build *ProjectQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (pgb *ProjectGroupBy) Aggregate(fns ...AggregateFunc) *ProjectGroupBy {
	pgb.fns = append(pgb.fns, fns...)
	return pgb
}

// Scan applies the selector query and scans the result into the given value.
func (pgb *ProjectGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, pgb.build.ctx, "GroupBy")
	if err := pgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*ProjectQuery, *ProjectGroupBy](ctx, pgb.build, pgb, pgb.build.inters, v)
}

func (pgb *ProjectGroupBy) sqlScan(ctx context.Context, root *ProjectQuery, v any) error {
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

// ProjectSelect is the builder for selecting fields of Project entities.
type ProjectSelect struct {
	*ProjectQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (ps *ProjectSelect) Aggregate(fns ...AggregateFunc) *ProjectSelect {
	ps.fns = append(ps.fns, fns...)
	return ps
}

// Scan applies the selector query and scans the result into the given value.
func (ps *ProjectSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ps.ctx, "Select")
	if err := ps.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*ProjectQuery, *ProjectSelect](ctx, ps.ProjectQuery, ps, ps.inters, v)
}

func (ps *ProjectSelect) sqlScan(ctx context.Context, root *ProjectQuery, v any) error {
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
func (ps *ProjectSelect) Modify(modifiers ...func(s *sql.Selector)) *ProjectSelect {
	ps.modifiers = append(ps.modifiers, modifiers...)
	return ps
}
