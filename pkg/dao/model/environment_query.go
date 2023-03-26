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

	"github.com/seal-io/seal/pkg/dao/model/applicationinstance"
	"github.com/seal-io/seal/pkg/dao/model/applicationrevision"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/environmentconnectorrelationship"
	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/types/oid"
)

// EnvironmentQuery is the builder for querying Environment entities.
type EnvironmentQuery struct {
	config
	ctx            *QueryContext
	order          []OrderFunc
	inters         []Interceptor
	predicates     []predicate.Environment
	withConnectors *EnvironmentConnectorRelationshipQuery
	withInstances  *ApplicationInstanceQuery
	withRevisions  *ApplicationRevisionQuery
	modifiers      []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the EnvironmentQuery builder.
func (eq *EnvironmentQuery) Where(ps ...predicate.Environment) *EnvironmentQuery {
	eq.predicates = append(eq.predicates, ps...)
	return eq
}

// Limit the number of records to be returned by this query.
func (eq *EnvironmentQuery) Limit(limit int) *EnvironmentQuery {
	eq.ctx.Limit = &limit
	return eq
}

// Offset to start from.
func (eq *EnvironmentQuery) Offset(offset int) *EnvironmentQuery {
	eq.ctx.Offset = &offset
	return eq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (eq *EnvironmentQuery) Unique(unique bool) *EnvironmentQuery {
	eq.ctx.Unique = &unique
	return eq
}

// Order specifies how the records should be ordered.
func (eq *EnvironmentQuery) Order(o ...OrderFunc) *EnvironmentQuery {
	eq.order = append(eq.order, o...)
	return eq
}

// QueryConnectors chains the current query on the "connectors" edge.
func (eq *EnvironmentQuery) QueryConnectors() *EnvironmentConnectorRelationshipQuery {
	query := (&EnvironmentConnectorRelationshipClient{config: eq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := eq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := eq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(environment.Table, environment.FieldID, selector),
			sqlgraph.To(environmentconnectorrelationship.Table, environmentconnectorrelationship.EnvironmentColumn),
			sqlgraph.Edge(sqlgraph.O2M, true, environment.ConnectorsTable, environment.ConnectorsColumn),
		)
		schemaConfig := eq.schemaConfig
		step.To.Schema = schemaConfig.EnvironmentConnectorRelationship
		step.Edge.Schema = schemaConfig.EnvironmentConnectorRelationship
		fromU = sqlgraph.SetNeighbors(eq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryInstances chains the current query on the "instances" edge.
func (eq *EnvironmentQuery) QueryInstances() *ApplicationInstanceQuery {
	query := (&ApplicationInstanceClient{config: eq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := eq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := eq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(environment.Table, environment.FieldID, selector),
			sqlgraph.To(applicationinstance.Table, applicationinstance.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, environment.InstancesTable, environment.InstancesColumn),
		)
		schemaConfig := eq.schemaConfig
		step.To.Schema = schemaConfig.ApplicationInstance
		step.Edge.Schema = schemaConfig.ApplicationInstance
		fromU = sqlgraph.SetNeighbors(eq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryRevisions chains the current query on the "revisions" edge.
func (eq *EnvironmentQuery) QueryRevisions() *ApplicationRevisionQuery {
	query := (&ApplicationRevisionClient{config: eq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := eq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := eq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(environment.Table, environment.FieldID, selector),
			sqlgraph.To(applicationrevision.Table, applicationrevision.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, environment.RevisionsTable, environment.RevisionsColumn),
		)
		schemaConfig := eq.schemaConfig
		step.To.Schema = schemaConfig.ApplicationRevision
		step.Edge.Schema = schemaConfig.ApplicationRevision
		fromU = sqlgraph.SetNeighbors(eq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Environment entity from the query.
// Returns a *NotFoundError when no Environment was found.
func (eq *EnvironmentQuery) First(ctx context.Context) (*Environment, error) {
	nodes, err := eq.Limit(1).All(setContextOp(ctx, eq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{environment.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (eq *EnvironmentQuery) FirstX(ctx context.Context) *Environment {
	node, err := eq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Environment ID from the query.
// Returns a *NotFoundError when no Environment ID was found.
func (eq *EnvironmentQuery) FirstID(ctx context.Context) (id oid.ID, err error) {
	var ids []oid.ID
	if ids, err = eq.Limit(1).IDs(setContextOp(ctx, eq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{environment.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (eq *EnvironmentQuery) FirstIDX(ctx context.Context) oid.ID {
	id, err := eq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Environment entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Environment entity is found.
// Returns a *NotFoundError when no Environment entities are found.
func (eq *EnvironmentQuery) Only(ctx context.Context) (*Environment, error) {
	nodes, err := eq.Limit(2).All(setContextOp(ctx, eq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{environment.Label}
	default:
		return nil, &NotSingularError{environment.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (eq *EnvironmentQuery) OnlyX(ctx context.Context) *Environment {
	node, err := eq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Environment ID in the query.
// Returns a *NotSingularError when more than one Environment ID is found.
// Returns a *NotFoundError when no entities are found.
func (eq *EnvironmentQuery) OnlyID(ctx context.Context) (id oid.ID, err error) {
	var ids []oid.ID
	if ids, err = eq.Limit(2).IDs(setContextOp(ctx, eq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{environment.Label}
	default:
		err = &NotSingularError{environment.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (eq *EnvironmentQuery) OnlyIDX(ctx context.Context) oid.ID {
	id, err := eq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Environments.
func (eq *EnvironmentQuery) All(ctx context.Context) ([]*Environment, error) {
	ctx = setContextOp(ctx, eq.ctx, "All")
	if err := eq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Environment, *EnvironmentQuery]()
	return withInterceptors[[]*Environment](ctx, eq, qr, eq.inters)
}

// AllX is like All, but panics if an error occurs.
func (eq *EnvironmentQuery) AllX(ctx context.Context) []*Environment {
	nodes, err := eq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Environment IDs.
func (eq *EnvironmentQuery) IDs(ctx context.Context) (ids []oid.ID, err error) {
	if eq.ctx.Unique == nil && eq.path != nil {
		eq.Unique(true)
	}
	ctx = setContextOp(ctx, eq.ctx, "IDs")
	if err = eq.Select(environment.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (eq *EnvironmentQuery) IDsX(ctx context.Context) []oid.ID {
	ids, err := eq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (eq *EnvironmentQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, eq.ctx, "Count")
	if err := eq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, eq, querierCount[*EnvironmentQuery](), eq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (eq *EnvironmentQuery) CountX(ctx context.Context) int {
	count, err := eq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (eq *EnvironmentQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, eq.ctx, "Exist")
	switch _, err := eq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("model: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (eq *EnvironmentQuery) ExistX(ctx context.Context) bool {
	exist, err := eq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the EnvironmentQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (eq *EnvironmentQuery) Clone() *EnvironmentQuery {
	if eq == nil {
		return nil
	}
	return &EnvironmentQuery{
		config:         eq.config,
		ctx:            eq.ctx.Clone(),
		order:          append([]OrderFunc{}, eq.order...),
		inters:         append([]Interceptor{}, eq.inters...),
		predicates:     append([]predicate.Environment{}, eq.predicates...),
		withConnectors: eq.withConnectors.Clone(),
		withInstances:  eq.withInstances.Clone(),
		withRevisions:  eq.withRevisions.Clone(),
		// clone intermediate query.
		sql:  eq.sql.Clone(),
		path: eq.path,
	}
}

// WithConnectors tells the query-builder to eager-load the nodes that are connected to
// the "connectors" edge. The optional arguments are used to configure the query builder of the edge.
func (eq *EnvironmentQuery) WithConnectors(opts ...func(*EnvironmentConnectorRelationshipQuery)) *EnvironmentQuery {
	query := (&EnvironmentConnectorRelationshipClient{config: eq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	eq.withConnectors = query
	return eq
}

// WithInstances tells the query-builder to eager-load the nodes that are connected to
// the "instances" edge. The optional arguments are used to configure the query builder of the edge.
func (eq *EnvironmentQuery) WithInstances(opts ...func(*ApplicationInstanceQuery)) *EnvironmentQuery {
	query := (&ApplicationInstanceClient{config: eq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	eq.withInstances = query
	return eq
}

// WithRevisions tells the query-builder to eager-load the nodes that are connected to
// the "revisions" edge. The optional arguments are used to configure the query builder of the edge.
func (eq *EnvironmentQuery) WithRevisions(opts ...func(*ApplicationRevisionQuery)) *EnvironmentQuery {
	query := (&ApplicationRevisionClient{config: eq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	eq.withRevisions = query
	return eq
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
//	client.Environment.Query().
//		GroupBy(environment.FieldName).
//		Aggregate(model.Count()).
//		Scan(ctx, &v)
func (eq *EnvironmentQuery) GroupBy(field string, fields ...string) *EnvironmentGroupBy {
	eq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &EnvironmentGroupBy{build: eq}
	grbuild.flds = &eq.ctx.Fields
	grbuild.label = environment.Label
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
//	client.Environment.Query().
//		Select(environment.FieldName).
//		Scan(ctx, &v)
func (eq *EnvironmentQuery) Select(fields ...string) *EnvironmentSelect {
	eq.ctx.Fields = append(eq.ctx.Fields, fields...)
	sbuild := &EnvironmentSelect{EnvironmentQuery: eq}
	sbuild.label = environment.Label
	sbuild.flds, sbuild.scan = &eq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a EnvironmentSelect configured with the given aggregations.
func (eq *EnvironmentQuery) Aggregate(fns ...AggregateFunc) *EnvironmentSelect {
	return eq.Select().Aggregate(fns...)
}

func (eq *EnvironmentQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range eq.inters {
		if inter == nil {
			return fmt.Errorf("model: uninitialized interceptor (forgotten import model/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, eq); err != nil {
				return err
			}
		}
	}
	for _, f := range eq.ctx.Fields {
		if !environment.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("model: invalid field %q for query", f)}
		}
	}
	if eq.path != nil {
		prev, err := eq.path(ctx)
		if err != nil {
			return err
		}
		eq.sql = prev
	}
	return nil
}

func (eq *EnvironmentQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Environment, error) {
	var (
		nodes       = []*Environment{}
		_spec       = eq.querySpec()
		loadedTypes = [3]bool{
			eq.withConnectors != nil,
			eq.withInstances != nil,
			eq.withRevisions != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Environment).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Environment{config: eq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	_spec.Node.Schema = eq.schemaConfig.Environment
	ctx = internal.NewSchemaConfigContext(ctx, eq.schemaConfig)
	if len(eq.modifiers) > 0 {
		_spec.Modifiers = eq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, eq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := eq.withConnectors; query != nil {
		if err := eq.loadConnectors(ctx, query, nodes,
			func(n *Environment) { n.Edges.Connectors = []*EnvironmentConnectorRelationship{} },
			func(n *Environment, e *EnvironmentConnectorRelationship) {
				n.Edges.Connectors = append(n.Edges.Connectors, e)
			}); err != nil {
			return nil, err
		}
	}
	if query := eq.withInstances; query != nil {
		if err := eq.loadInstances(ctx, query, nodes,
			func(n *Environment) { n.Edges.Instances = []*ApplicationInstance{} },
			func(n *Environment, e *ApplicationInstance) { n.Edges.Instances = append(n.Edges.Instances, e) }); err != nil {
			return nil, err
		}
	}
	if query := eq.withRevisions; query != nil {
		if err := eq.loadRevisions(ctx, query, nodes,
			func(n *Environment) { n.Edges.Revisions = []*ApplicationRevision{} },
			func(n *Environment, e *ApplicationRevision) { n.Edges.Revisions = append(n.Edges.Revisions, e) }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (eq *EnvironmentQuery) loadConnectors(ctx context.Context, query *EnvironmentConnectorRelationshipQuery, nodes []*Environment, init func(*Environment), assign func(*Environment, *EnvironmentConnectorRelationship)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[oid.ID]*Environment)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.Where(predicate.EnvironmentConnectorRelationship(func(s *sql.Selector) {
		s.Where(sql.InValues(environment.ConnectorsColumn, fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.EnvironmentID
		node, ok := nodeids[fk]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "environment_id" returned %v for node %v`, fk, n)
		}
		assign(node, n)
	}
	return nil
}
func (eq *EnvironmentQuery) loadInstances(ctx context.Context, query *ApplicationInstanceQuery, nodes []*Environment, init func(*Environment), assign func(*Environment, *ApplicationInstance)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[oid.ID]*Environment)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.Where(predicate.ApplicationInstance(func(s *sql.Selector) {
		s.Where(sql.InValues(environment.InstancesColumn, fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.EnvironmentID
		node, ok := nodeids[fk]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "environmentID" returned %v for node %v`, fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}
func (eq *EnvironmentQuery) loadRevisions(ctx context.Context, query *ApplicationRevisionQuery, nodes []*Environment, init func(*Environment), assign func(*Environment, *ApplicationRevision)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[oid.ID]*Environment)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.Where(predicate.ApplicationRevision(func(s *sql.Selector) {
		s.Where(sql.InValues(environment.RevisionsColumn, fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.EnvironmentID
		node, ok := nodeids[fk]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "environmentID" returned %v for node %v`, fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}

func (eq *EnvironmentQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := eq.querySpec()
	_spec.Node.Schema = eq.schemaConfig.Environment
	ctx = internal.NewSchemaConfigContext(ctx, eq.schemaConfig)
	if len(eq.modifiers) > 0 {
		_spec.Modifiers = eq.modifiers
	}
	_spec.Node.Columns = eq.ctx.Fields
	if len(eq.ctx.Fields) > 0 {
		_spec.Unique = eq.ctx.Unique != nil && *eq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, eq.driver, _spec)
}

func (eq *EnvironmentQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(environment.Table, environment.Columns, sqlgraph.NewFieldSpec(environment.FieldID, field.TypeString))
	_spec.From = eq.sql
	if unique := eq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if eq.path != nil {
		_spec.Unique = true
	}
	if fields := eq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, environment.FieldID)
		for i := range fields {
			if fields[i] != environment.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := eq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := eq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := eq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := eq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (eq *EnvironmentQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(eq.driver.Dialect())
	t1 := builder.Table(environment.Table)
	columns := eq.ctx.Fields
	if len(columns) == 0 {
		columns = environment.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if eq.sql != nil {
		selector = eq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if eq.ctx.Unique != nil && *eq.ctx.Unique {
		selector.Distinct()
	}
	t1.Schema(eq.schemaConfig.Environment)
	ctx = internal.NewSchemaConfigContext(ctx, eq.schemaConfig)
	selector.WithContext(ctx)
	for _, m := range eq.modifiers {
		m(selector)
	}
	for _, p := range eq.predicates {
		p(selector)
	}
	for _, p := range eq.order {
		p(selector)
	}
	if offset := eq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := eq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ForUpdate locks the selected rows against concurrent updates, and prevent them from being
// updated, deleted or "selected ... for update" by other sessions, until the transaction is
// either committed or rolled-back.
func (eq *EnvironmentQuery) ForUpdate(opts ...sql.LockOption) *EnvironmentQuery {
	if eq.driver.Dialect() == dialect.Postgres {
		eq.Unique(false)
	}
	eq.modifiers = append(eq.modifiers, func(s *sql.Selector) {
		s.ForUpdate(opts...)
	})
	return eq
}

// ForShare behaves similarly to ForUpdate, except that it acquires a shared mode lock
// on any rows that are read. Other sessions can read the rows, but cannot modify them
// until your transaction commits.
func (eq *EnvironmentQuery) ForShare(opts ...sql.LockOption) *EnvironmentQuery {
	if eq.driver.Dialect() == dialect.Postgres {
		eq.Unique(false)
	}
	eq.modifiers = append(eq.modifiers, func(s *sql.Selector) {
		s.ForShare(opts...)
	})
	return eq
}

// Modify adds a query modifier for attaching custom logic to queries.
func (eq *EnvironmentQuery) Modify(modifiers ...func(s *sql.Selector)) *EnvironmentSelect {
	eq.modifiers = append(eq.modifiers, modifiers...)
	return eq.Select()
}

// EnvironmentGroupBy is the group-by builder for Environment entities.
type EnvironmentGroupBy struct {
	selector
	build *EnvironmentQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (egb *EnvironmentGroupBy) Aggregate(fns ...AggregateFunc) *EnvironmentGroupBy {
	egb.fns = append(egb.fns, fns...)
	return egb
}

// Scan applies the selector query and scans the result into the given value.
func (egb *EnvironmentGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, egb.build.ctx, "GroupBy")
	if err := egb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*EnvironmentQuery, *EnvironmentGroupBy](ctx, egb.build, egb, egb.build.inters, v)
}

func (egb *EnvironmentGroupBy) sqlScan(ctx context.Context, root *EnvironmentQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(egb.fns))
	for _, fn := range egb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*egb.flds)+len(egb.fns))
		for _, f := range *egb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*egb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := egb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// EnvironmentSelect is the builder for selecting fields of Environment entities.
type EnvironmentSelect struct {
	*EnvironmentQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (es *EnvironmentSelect) Aggregate(fns ...AggregateFunc) *EnvironmentSelect {
	es.fns = append(es.fns, fns...)
	return es
}

// Scan applies the selector query and scans the result into the given value.
func (es *EnvironmentSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, es.ctx, "Select")
	if err := es.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*EnvironmentQuery, *EnvironmentSelect](ctx, es.EnvironmentQuery, es, es.inters, v)
}

func (es *EnvironmentSelect) sqlScan(ctx context.Context, root *EnvironmentQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(es.fns))
	for _, fn := range es.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*es.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := es.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (es *EnvironmentSelect) Modify(modifiers ...func(s *sql.Selector)) *EnvironmentSelect {
	es.modifiers = append(es.modifiers, modifiers...)
	return es
}
