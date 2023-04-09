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
	"github.com/seal-io/seal/pkg/dao/model/applicationinstance"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/model/applicationrevision"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/types/oid"
)

// ApplicationInstanceQuery is the builder for querying ApplicationInstance entities.
type ApplicationInstanceQuery struct {
	config
	ctx             *QueryContext
	order           []OrderFunc
	inters          []Interceptor
	predicates      []predicate.ApplicationInstance
	withApplication *ApplicationQuery
	withEnvironment *EnvironmentQuery
	withRevisions   *ApplicationRevisionQuery
	withResources   *ApplicationResourceQuery
	modifiers       []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the ApplicationInstanceQuery builder.
func (aiq *ApplicationInstanceQuery) Where(ps ...predicate.ApplicationInstance) *ApplicationInstanceQuery {
	aiq.predicates = append(aiq.predicates, ps...)
	return aiq
}

// Limit the number of records to be returned by this query.
func (aiq *ApplicationInstanceQuery) Limit(limit int) *ApplicationInstanceQuery {
	aiq.ctx.Limit = &limit
	return aiq
}

// Offset to start from.
func (aiq *ApplicationInstanceQuery) Offset(offset int) *ApplicationInstanceQuery {
	aiq.ctx.Offset = &offset
	return aiq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (aiq *ApplicationInstanceQuery) Unique(unique bool) *ApplicationInstanceQuery {
	aiq.ctx.Unique = &unique
	return aiq
}

// Order specifies how the records should be ordered.
func (aiq *ApplicationInstanceQuery) Order(o ...OrderFunc) *ApplicationInstanceQuery {
	aiq.order = append(aiq.order, o...)
	return aiq
}

// QueryApplication chains the current query on the "application" edge.
func (aiq *ApplicationInstanceQuery) QueryApplication() *ApplicationQuery {
	query := (&ApplicationClient{config: aiq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := aiq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := aiq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(applicationinstance.Table, applicationinstance.FieldID, selector),
			sqlgraph.To(application.Table, application.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, applicationinstance.ApplicationTable, applicationinstance.ApplicationColumn),
		)
		schemaConfig := aiq.schemaConfig
		step.To.Schema = schemaConfig.Application
		step.Edge.Schema = schemaConfig.ApplicationInstance
		fromU = sqlgraph.SetNeighbors(aiq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryEnvironment chains the current query on the "environment" edge.
func (aiq *ApplicationInstanceQuery) QueryEnvironment() *EnvironmentQuery {
	query := (&EnvironmentClient{config: aiq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := aiq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := aiq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(applicationinstance.Table, applicationinstance.FieldID, selector),
			sqlgraph.To(environment.Table, environment.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, applicationinstance.EnvironmentTable, applicationinstance.EnvironmentColumn),
		)
		schemaConfig := aiq.schemaConfig
		step.To.Schema = schemaConfig.Environment
		step.Edge.Schema = schemaConfig.ApplicationInstance
		fromU = sqlgraph.SetNeighbors(aiq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryRevisions chains the current query on the "revisions" edge.
func (aiq *ApplicationInstanceQuery) QueryRevisions() *ApplicationRevisionQuery {
	query := (&ApplicationRevisionClient{config: aiq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := aiq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := aiq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(applicationinstance.Table, applicationinstance.FieldID, selector),
			sqlgraph.To(applicationrevision.Table, applicationrevision.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, applicationinstance.RevisionsTable, applicationinstance.RevisionsColumn),
		)
		schemaConfig := aiq.schemaConfig
		step.To.Schema = schemaConfig.ApplicationRevision
		step.Edge.Schema = schemaConfig.ApplicationRevision
		fromU = sqlgraph.SetNeighbors(aiq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryResources chains the current query on the "resources" edge.
func (aiq *ApplicationInstanceQuery) QueryResources() *ApplicationResourceQuery {
	query := (&ApplicationResourceClient{config: aiq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := aiq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := aiq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(applicationinstance.Table, applicationinstance.FieldID, selector),
			sqlgraph.To(applicationresource.Table, applicationresource.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, applicationinstance.ResourcesTable, applicationinstance.ResourcesColumn),
		)
		schemaConfig := aiq.schemaConfig
		step.To.Schema = schemaConfig.ApplicationResource
		step.Edge.Schema = schemaConfig.ApplicationResource
		fromU = sqlgraph.SetNeighbors(aiq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first ApplicationInstance entity from the query.
// Returns a *NotFoundError when no ApplicationInstance was found.
func (aiq *ApplicationInstanceQuery) First(ctx context.Context) (*ApplicationInstance, error) {
	nodes, err := aiq.Limit(1).All(setContextOp(ctx, aiq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{applicationinstance.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (aiq *ApplicationInstanceQuery) FirstX(ctx context.Context) *ApplicationInstance {
	node, err := aiq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first ApplicationInstance ID from the query.
// Returns a *NotFoundError when no ApplicationInstance ID was found.
func (aiq *ApplicationInstanceQuery) FirstID(ctx context.Context) (id oid.ID, err error) {
	var ids []oid.ID
	if ids, err = aiq.Limit(1).IDs(setContextOp(ctx, aiq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{applicationinstance.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (aiq *ApplicationInstanceQuery) FirstIDX(ctx context.Context) oid.ID {
	id, err := aiq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single ApplicationInstance entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one ApplicationInstance entity is found.
// Returns a *NotFoundError when no ApplicationInstance entities are found.
func (aiq *ApplicationInstanceQuery) Only(ctx context.Context) (*ApplicationInstance, error) {
	nodes, err := aiq.Limit(2).All(setContextOp(ctx, aiq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{applicationinstance.Label}
	default:
		return nil, &NotSingularError{applicationinstance.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (aiq *ApplicationInstanceQuery) OnlyX(ctx context.Context) *ApplicationInstance {
	node, err := aiq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only ApplicationInstance ID in the query.
// Returns a *NotSingularError when more than one ApplicationInstance ID is found.
// Returns a *NotFoundError when no entities are found.
func (aiq *ApplicationInstanceQuery) OnlyID(ctx context.Context) (id oid.ID, err error) {
	var ids []oid.ID
	if ids, err = aiq.Limit(2).IDs(setContextOp(ctx, aiq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{applicationinstance.Label}
	default:
		err = &NotSingularError{applicationinstance.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (aiq *ApplicationInstanceQuery) OnlyIDX(ctx context.Context) oid.ID {
	id, err := aiq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of ApplicationInstances.
func (aiq *ApplicationInstanceQuery) All(ctx context.Context) ([]*ApplicationInstance, error) {
	ctx = setContextOp(ctx, aiq.ctx, "All")
	if err := aiq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*ApplicationInstance, *ApplicationInstanceQuery]()
	return withInterceptors[[]*ApplicationInstance](ctx, aiq, qr, aiq.inters)
}

// AllX is like All, but panics if an error occurs.
func (aiq *ApplicationInstanceQuery) AllX(ctx context.Context) []*ApplicationInstance {
	nodes, err := aiq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of ApplicationInstance IDs.
func (aiq *ApplicationInstanceQuery) IDs(ctx context.Context) (ids []oid.ID, err error) {
	if aiq.ctx.Unique == nil && aiq.path != nil {
		aiq.Unique(true)
	}
	ctx = setContextOp(ctx, aiq.ctx, "IDs")
	if err = aiq.Select(applicationinstance.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (aiq *ApplicationInstanceQuery) IDsX(ctx context.Context) []oid.ID {
	ids, err := aiq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (aiq *ApplicationInstanceQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, aiq.ctx, "Count")
	if err := aiq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, aiq, querierCount[*ApplicationInstanceQuery](), aiq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (aiq *ApplicationInstanceQuery) CountX(ctx context.Context) int {
	count, err := aiq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (aiq *ApplicationInstanceQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, aiq.ctx, "Exist")
	switch _, err := aiq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("model: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (aiq *ApplicationInstanceQuery) ExistX(ctx context.Context) bool {
	exist, err := aiq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the ApplicationInstanceQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (aiq *ApplicationInstanceQuery) Clone() *ApplicationInstanceQuery {
	if aiq == nil {
		return nil
	}
	return &ApplicationInstanceQuery{
		config:          aiq.config,
		ctx:             aiq.ctx.Clone(),
		order:           append([]OrderFunc{}, aiq.order...),
		inters:          append([]Interceptor{}, aiq.inters...),
		predicates:      append([]predicate.ApplicationInstance{}, aiq.predicates...),
		withApplication: aiq.withApplication.Clone(),
		withEnvironment: aiq.withEnvironment.Clone(),
		withRevisions:   aiq.withRevisions.Clone(),
		withResources:   aiq.withResources.Clone(),
		// clone intermediate query.
		sql:  aiq.sql.Clone(),
		path: aiq.path,
	}
}

// WithApplication tells the query-builder to eager-load the nodes that are connected to
// the "application" edge. The optional arguments are used to configure the query builder of the edge.
func (aiq *ApplicationInstanceQuery) WithApplication(opts ...func(*ApplicationQuery)) *ApplicationInstanceQuery {
	query := (&ApplicationClient{config: aiq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	aiq.withApplication = query
	return aiq
}

// WithEnvironment tells the query-builder to eager-load the nodes that are connected to
// the "environment" edge. The optional arguments are used to configure the query builder of the edge.
func (aiq *ApplicationInstanceQuery) WithEnvironment(opts ...func(*EnvironmentQuery)) *ApplicationInstanceQuery {
	query := (&EnvironmentClient{config: aiq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	aiq.withEnvironment = query
	return aiq
}

// WithRevisions tells the query-builder to eager-load the nodes that are connected to
// the "revisions" edge. The optional arguments are used to configure the query builder of the edge.
func (aiq *ApplicationInstanceQuery) WithRevisions(opts ...func(*ApplicationRevisionQuery)) *ApplicationInstanceQuery {
	query := (&ApplicationRevisionClient{config: aiq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	aiq.withRevisions = query
	return aiq
}

// WithResources tells the query-builder to eager-load the nodes that are connected to
// the "resources" edge. The optional arguments are used to configure the query builder of the edge.
func (aiq *ApplicationInstanceQuery) WithResources(opts ...func(*ApplicationResourceQuery)) *ApplicationInstanceQuery {
	query := (&ApplicationResourceClient{config: aiq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	aiq.withResources = query
	return aiq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Status string `json:"status,omitempty" sql:"status"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.ApplicationInstance.Query().
//		GroupBy(applicationinstance.FieldStatus).
//		Aggregate(model.Count()).
//		Scan(ctx, &v)
func (aiq *ApplicationInstanceQuery) GroupBy(field string, fields ...string) *ApplicationInstanceGroupBy {
	aiq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &ApplicationInstanceGroupBy{build: aiq}
	grbuild.flds = &aiq.ctx.Fields
	grbuild.label = applicationinstance.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Status string `json:"status,omitempty" sql:"status"`
//	}
//
//	client.ApplicationInstance.Query().
//		Select(applicationinstance.FieldStatus).
//		Scan(ctx, &v)
func (aiq *ApplicationInstanceQuery) Select(fields ...string) *ApplicationInstanceSelect {
	aiq.ctx.Fields = append(aiq.ctx.Fields, fields...)
	sbuild := &ApplicationInstanceSelect{ApplicationInstanceQuery: aiq}
	sbuild.label = applicationinstance.Label
	sbuild.flds, sbuild.scan = &aiq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a ApplicationInstanceSelect configured with the given aggregations.
func (aiq *ApplicationInstanceQuery) Aggregate(fns ...AggregateFunc) *ApplicationInstanceSelect {
	return aiq.Select().Aggregate(fns...)
}

func (aiq *ApplicationInstanceQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range aiq.inters {
		if inter == nil {
			return fmt.Errorf("model: uninitialized interceptor (forgotten import model/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, aiq); err != nil {
				return err
			}
		}
	}
	for _, f := range aiq.ctx.Fields {
		if !applicationinstance.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("model: invalid field %q for query", f)}
		}
	}
	if aiq.path != nil {
		prev, err := aiq.path(ctx)
		if err != nil {
			return err
		}
		aiq.sql = prev
	}
	return nil
}

func (aiq *ApplicationInstanceQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*ApplicationInstance, error) {
	var (
		nodes       = []*ApplicationInstance{}
		_spec       = aiq.querySpec()
		loadedTypes = [4]bool{
			aiq.withApplication != nil,
			aiq.withEnvironment != nil,
			aiq.withRevisions != nil,
			aiq.withResources != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*ApplicationInstance).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &ApplicationInstance{config: aiq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	_spec.Node.Schema = aiq.schemaConfig.ApplicationInstance
	ctx = internal.NewSchemaConfigContext(ctx, aiq.schemaConfig)
	if len(aiq.modifiers) > 0 {
		_spec.Modifiers = aiq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, aiq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := aiq.withApplication; query != nil {
		if err := aiq.loadApplication(ctx, query, nodes, nil,
			func(n *ApplicationInstance, e *Application) { n.Edges.Application = e }); err != nil {
			return nil, err
		}
	}
	if query := aiq.withEnvironment; query != nil {
		if err := aiq.loadEnvironment(ctx, query, nodes, nil,
			func(n *ApplicationInstance, e *Environment) { n.Edges.Environment = e }); err != nil {
			return nil, err
		}
	}
	if query := aiq.withRevisions; query != nil {
		if err := aiq.loadRevisions(ctx, query, nodes,
			func(n *ApplicationInstance) { n.Edges.Revisions = []*ApplicationRevision{} },
			func(n *ApplicationInstance, e *ApplicationRevision) { n.Edges.Revisions = append(n.Edges.Revisions, e) }); err != nil {
			return nil, err
		}
	}
	if query := aiq.withResources; query != nil {
		if err := aiq.loadResources(ctx, query, nodes,
			func(n *ApplicationInstance) { n.Edges.Resources = []*ApplicationResource{} },
			func(n *ApplicationInstance, e *ApplicationResource) { n.Edges.Resources = append(n.Edges.Resources, e) }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (aiq *ApplicationInstanceQuery) loadApplication(ctx context.Context, query *ApplicationQuery, nodes []*ApplicationInstance, init func(*ApplicationInstance), assign func(*ApplicationInstance, *Application)) error {
	ids := make([]oid.ID, 0, len(nodes))
	nodeids := make(map[oid.ID][]*ApplicationInstance)
	for i := range nodes {
		fk := nodes[i].ApplicationID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(application.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "applicationID" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (aiq *ApplicationInstanceQuery) loadEnvironment(ctx context.Context, query *EnvironmentQuery, nodes []*ApplicationInstance, init func(*ApplicationInstance), assign func(*ApplicationInstance, *Environment)) error {
	ids := make([]oid.ID, 0, len(nodes))
	nodeids := make(map[oid.ID][]*ApplicationInstance)
	for i := range nodes {
		fk := nodes[i].EnvironmentID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(environment.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "environmentID" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (aiq *ApplicationInstanceQuery) loadRevisions(ctx context.Context, query *ApplicationRevisionQuery, nodes []*ApplicationInstance, init func(*ApplicationInstance), assign func(*ApplicationInstance, *ApplicationRevision)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[oid.ID]*ApplicationInstance)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.Where(predicate.ApplicationRevision(func(s *sql.Selector) {
		s.Where(sql.InValues(applicationinstance.RevisionsColumn, fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.InstanceID
		node, ok := nodeids[fk]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "instanceID" returned %v for node %v`, fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}
func (aiq *ApplicationInstanceQuery) loadResources(ctx context.Context, query *ApplicationResourceQuery, nodes []*ApplicationInstance, init func(*ApplicationInstance), assign func(*ApplicationInstance, *ApplicationResource)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[oid.ID]*ApplicationInstance)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.Where(predicate.ApplicationResource(func(s *sql.Selector) {
		s.Where(sql.InValues(applicationinstance.ResourcesColumn, fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.InstanceID
		node, ok := nodeids[fk]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "instanceID" returned %v for node %v`, fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}

func (aiq *ApplicationInstanceQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := aiq.querySpec()
	_spec.Node.Schema = aiq.schemaConfig.ApplicationInstance
	ctx = internal.NewSchemaConfigContext(ctx, aiq.schemaConfig)
	if len(aiq.modifiers) > 0 {
		_spec.Modifiers = aiq.modifiers
	}
	_spec.Node.Columns = aiq.ctx.Fields
	if len(aiq.ctx.Fields) > 0 {
		_spec.Unique = aiq.ctx.Unique != nil && *aiq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, aiq.driver, _spec)
}

func (aiq *ApplicationInstanceQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(applicationinstance.Table, applicationinstance.Columns, sqlgraph.NewFieldSpec(applicationinstance.FieldID, field.TypeString))
	_spec.From = aiq.sql
	if unique := aiq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if aiq.path != nil {
		_spec.Unique = true
	}
	if fields := aiq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, applicationinstance.FieldID)
		for i := range fields {
			if fields[i] != applicationinstance.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := aiq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := aiq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := aiq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := aiq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (aiq *ApplicationInstanceQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(aiq.driver.Dialect())
	t1 := builder.Table(applicationinstance.Table)
	columns := aiq.ctx.Fields
	if len(columns) == 0 {
		columns = applicationinstance.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if aiq.sql != nil {
		selector = aiq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if aiq.ctx.Unique != nil && *aiq.ctx.Unique {
		selector.Distinct()
	}
	t1.Schema(aiq.schemaConfig.ApplicationInstance)
	ctx = internal.NewSchemaConfigContext(ctx, aiq.schemaConfig)
	selector.WithContext(ctx)
	for _, m := range aiq.modifiers {
		m(selector)
	}
	for _, p := range aiq.predicates {
		p(selector)
	}
	for _, p := range aiq.order {
		p(selector)
	}
	if offset := aiq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := aiq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ForUpdate locks the selected rows against concurrent updates, and prevent them from being
// updated, deleted or "selected ... for update" by other sessions, until the transaction is
// either committed or rolled-back.
func (aiq *ApplicationInstanceQuery) ForUpdate(opts ...sql.LockOption) *ApplicationInstanceQuery {
	if aiq.driver.Dialect() == dialect.Postgres {
		aiq.Unique(false)
	}
	aiq.modifiers = append(aiq.modifiers, func(s *sql.Selector) {
		s.ForUpdate(opts...)
	})
	return aiq
}

// ForShare behaves similarly to ForUpdate, except that it acquires a shared mode lock
// on any rows that are read. Other sessions can read the rows, but cannot modify them
// until your transaction commits.
func (aiq *ApplicationInstanceQuery) ForShare(opts ...sql.LockOption) *ApplicationInstanceQuery {
	if aiq.driver.Dialect() == dialect.Postgres {
		aiq.Unique(false)
	}
	aiq.modifiers = append(aiq.modifiers, func(s *sql.Selector) {
		s.ForShare(opts...)
	})
	return aiq
}

// Modify adds a query modifier for attaching custom logic to queries.
func (aiq *ApplicationInstanceQuery) Modify(modifiers ...func(s *sql.Selector)) *ApplicationInstanceSelect {
	aiq.modifiers = append(aiq.modifiers, modifiers...)
	return aiq.Select()
}

// WhereP appends storage-level predicates to the ApplicationInstanceQuery builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (aiq *ApplicationInstanceQuery) WhereP(ps ...func(*sql.Selector)) {
	var wps = make([]predicate.ApplicationInstance, 0, len(ps))
	for i := 0; i < len(ps); i++ {
		wps = append(wps, predicate.ApplicationInstance(ps[i]))
	}
	aiq.predicates = append(aiq.predicates, wps...)
}

// ApplicationInstanceGroupBy is the group-by builder for ApplicationInstance entities.
type ApplicationInstanceGroupBy struct {
	selector
	build *ApplicationInstanceQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (aigb *ApplicationInstanceGroupBy) Aggregate(fns ...AggregateFunc) *ApplicationInstanceGroupBy {
	aigb.fns = append(aigb.fns, fns...)
	return aigb
}

// Scan applies the selector query and scans the result into the given value.
func (aigb *ApplicationInstanceGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, aigb.build.ctx, "GroupBy")
	if err := aigb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*ApplicationInstanceQuery, *ApplicationInstanceGroupBy](ctx, aigb.build, aigb, aigb.build.inters, v)
}

func (aigb *ApplicationInstanceGroupBy) sqlScan(ctx context.Context, root *ApplicationInstanceQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(aigb.fns))
	for _, fn := range aigb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*aigb.flds)+len(aigb.fns))
		for _, f := range *aigb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*aigb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := aigb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// ApplicationInstanceSelect is the builder for selecting fields of ApplicationInstance entities.
type ApplicationInstanceSelect struct {
	*ApplicationInstanceQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (ais *ApplicationInstanceSelect) Aggregate(fns ...AggregateFunc) *ApplicationInstanceSelect {
	ais.fns = append(ais.fns, fns...)
	return ais
}

// Scan applies the selector query and scans the result into the given value.
func (ais *ApplicationInstanceSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ais.ctx, "Select")
	if err := ais.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*ApplicationInstanceQuery, *ApplicationInstanceSelect](ctx, ais.ApplicationInstanceQuery, ais, ais.inters, v)
}

func (ais *ApplicationInstanceSelect) sqlScan(ctx context.Context, root *ApplicationInstanceQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(ais.fns))
	for _, fn := range ais.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*ais.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ais.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (ais *ApplicationInstanceSelect) Modify(modifiers ...func(s *sql.Selector)) *ApplicationInstanceSelect {
	ais.modifiers = append(ais.modifiers, modifiers...)
	return ais
}
