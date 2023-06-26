// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "ent". DO NOT EDIT.

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

	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/model/internal"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/service"
	"github.com/seal-io/seal/pkg/dao/model/serviceresource"
	"github.com/seal-io/seal/pkg/dao/types/oid"
)

// ServiceResourceQuery is the builder for querying ServiceResource entities.
type ServiceResourceQuery struct {
	config
	ctx             *QueryContext
	order           []serviceresource.OrderOption
	inters          []Interceptor
	predicates      []predicate.ServiceResource
	withService     *ServiceQuery
	withConnector   *ConnectorQuery
	withComposition *ServiceResourceQuery
	withComponents  *ServiceResourceQuery
	modifiers       []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the ServiceResourceQuery builder.
func (srq *ServiceResourceQuery) Where(ps ...predicate.ServiceResource) *ServiceResourceQuery {
	srq.predicates = append(srq.predicates, ps...)
	return srq
}

// Limit the number of records to be returned by this query.
func (srq *ServiceResourceQuery) Limit(limit int) *ServiceResourceQuery {
	srq.ctx.Limit = &limit
	return srq
}

// Offset to start from.
func (srq *ServiceResourceQuery) Offset(offset int) *ServiceResourceQuery {
	srq.ctx.Offset = &offset
	return srq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (srq *ServiceResourceQuery) Unique(unique bool) *ServiceResourceQuery {
	srq.ctx.Unique = &unique
	return srq
}

// Order specifies how the records should be ordered.
func (srq *ServiceResourceQuery) Order(o ...serviceresource.OrderOption) *ServiceResourceQuery {
	srq.order = append(srq.order, o...)
	return srq
}

// QueryService chains the current query on the "service" edge.
func (srq *ServiceResourceQuery) QueryService() *ServiceQuery {
	query := (&ServiceClient{config: srq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := srq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := srq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(serviceresource.Table, serviceresource.FieldID, selector),
			sqlgraph.To(service.Table, service.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, serviceresource.ServiceTable, serviceresource.ServiceColumn),
		)
		schemaConfig := srq.schemaConfig
		step.To.Schema = schemaConfig.Service
		step.Edge.Schema = schemaConfig.ServiceResource
		fromU = sqlgraph.SetNeighbors(srq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryConnector chains the current query on the "connector" edge.
func (srq *ServiceResourceQuery) QueryConnector() *ConnectorQuery {
	query := (&ConnectorClient{config: srq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := srq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := srq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(serviceresource.Table, serviceresource.FieldID, selector),
			sqlgraph.To(connector.Table, connector.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, serviceresource.ConnectorTable, serviceresource.ConnectorColumn),
		)
		schemaConfig := srq.schemaConfig
		step.To.Schema = schemaConfig.Connector
		step.Edge.Schema = schemaConfig.ServiceResource
		fromU = sqlgraph.SetNeighbors(srq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryComposition chains the current query on the "composition" edge.
func (srq *ServiceResourceQuery) QueryComposition() *ServiceResourceQuery {
	query := (&ServiceResourceClient{config: srq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := srq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := srq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(serviceresource.Table, serviceresource.FieldID, selector),
			sqlgraph.To(serviceresource.Table, serviceresource.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, serviceresource.CompositionTable, serviceresource.CompositionColumn),
		)
		schemaConfig := srq.schemaConfig
		step.To.Schema = schemaConfig.ServiceResource
		step.Edge.Schema = schemaConfig.ServiceResource
		fromU = sqlgraph.SetNeighbors(srq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryComponents chains the current query on the "components" edge.
func (srq *ServiceResourceQuery) QueryComponents() *ServiceResourceQuery {
	query := (&ServiceResourceClient{config: srq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := srq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := srq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(serviceresource.Table, serviceresource.FieldID, selector),
			sqlgraph.To(serviceresource.Table, serviceresource.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, serviceresource.ComponentsTable, serviceresource.ComponentsColumn),
		)
		schemaConfig := srq.schemaConfig
		step.To.Schema = schemaConfig.ServiceResource
		step.Edge.Schema = schemaConfig.ServiceResource
		fromU = sqlgraph.SetNeighbors(srq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first ServiceResource entity from the query.
// Returns a *NotFoundError when no ServiceResource was found.
func (srq *ServiceResourceQuery) First(ctx context.Context) (*ServiceResource, error) {
	nodes, err := srq.Limit(1).All(setContextOp(ctx, srq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{serviceresource.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (srq *ServiceResourceQuery) FirstX(ctx context.Context) *ServiceResource {
	node, err := srq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first ServiceResource ID from the query.
// Returns a *NotFoundError when no ServiceResource ID was found.
func (srq *ServiceResourceQuery) FirstID(ctx context.Context) (id oid.ID, err error) {
	var ids []oid.ID
	if ids, err = srq.Limit(1).IDs(setContextOp(ctx, srq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{serviceresource.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (srq *ServiceResourceQuery) FirstIDX(ctx context.Context) oid.ID {
	id, err := srq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single ServiceResource entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one ServiceResource entity is found.
// Returns a *NotFoundError when no ServiceResource entities are found.
func (srq *ServiceResourceQuery) Only(ctx context.Context) (*ServiceResource, error) {
	nodes, err := srq.Limit(2).All(setContextOp(ctx, srq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{serviceresource.Label}
	default:
		return nil, &NotSingularError{serviceresource.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (srq *ServiceResourceQuery) OnlyX(ctx context.Context) *ServiceResource {
	node, err := srq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only ServiceResource ID in the query.
// Returns a *NotSingularError when more than one ServiceResource ID is found.
// Returns a *NotFoundError when no entities are found.
func (srq *ServiceResourceQuery) OnlyID(ctx context.Context) (id oid.ID, err error) {
	var ids []oid.ID
	if ids, err = srq.Limit(2).IDs(setContextOp(ctx, srq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{serviceresource.Label}
	default:
		err = &NotSingularError{serviceresource.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (srq *ServiceResourceQuery) OnlyIDX(ctx context.Context) oid.ID {
	id, err := srq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of ServiceResources.
func (srq *ServiceResourceQuery) All(ctx context.Context) ([]*ServiceResource, error) {
	ctx = setContextOp(ctx, srq.ctx, "All")
	if err := srq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*ServiceResource, *ServiceResourceQuery]()
	return withInterceptors[[]*ServiceResource](ctx, srq, qr, srq.inters)
}

// AllX is like All, but panics if an error occurs.
func (srq *ServiceResourceQuery) AllX(ctx context.Context) []*ServiceResource {
	nodes, err := srq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of ServiceResource IDs.
func (srq *ServiceResourceQuery) IDs(ctx context.Context) (ids []oid.ID, err error) {
	if srq.ctx.Unique == nil && srq.path != nil {
		srq.Unique(true)
	}
	ctx = setContextOp(ctx, srq.ctx, "IDs")
	if err = srq.Select(serviceresource.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (srq *ServiceResourceQuery) IDsX(ctx context.Context) []oid.ID {
	ids, err := srq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (srq *ServiceResourceQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, srq.ctx, "Count")
	if err := srq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, srq, querierCount[*ServiceResourceQuery](), srq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (srq *ServiceResourceQuery) CountX(ctx context.Context) int {
	count, err := srq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (srq *ServiceResourceQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, srq.ctx, "Exist")
	switch _, err := srq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("model: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (srq *ServiceResourceQuery) ExistX(ctx context.Context) bool {
	exist, err := srq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the ServiceResourceQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (srq *ServiceResourceQuery) Clone() *ServiceResourceQuery {
	if srq == nil {
		return nil
	}
	return &ServiceResourceQuery{
		config:          srq.config,
		ctx:             srq.ctx.Clone(),
		order:           append([]serviceresource.OrderOption{}, srq.order...),
		inters:          append([]Interceptor{}, srq.inters...),
		predicates:      append([]predicate.ServiceResource{}, srq.predicates...),
		withService:     srq.withService.Clone(),
		withConnector:   srq.withConnector.Clone(),
		withComposition: srq.withComposition.Clone(),
		withComponents:  srq.withComponents.Clone(),
		// clone intermediate query.
		sql:  srq.sql.Clone(),
		path: srq.path,
	}
}

// WithService tells the query-builder to eager-load the nodes that are connected to
// the "service" edge. The optional arguments are used to configure the query builder of the edge.
func (srq *ServiceResourceQuery) WithService(opts ...func(*ServiceQuery)) *ServiceResourceQuery {
	query := (&ServiceClient{config: srq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	srq.withService = query
	return srq
}

// WithConnector tells the query-builder to eager-load the nodes that are connected to
// the "connector" edge. The optional arguments are used to configure the query builder of the edge.
func (srq *ServiceResourceQuery) WithConnector(opts ...func(*ConnectorQuery)) *ServiceResourceQuery {
	query := (&ConnectorClient{config: srq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	srq.withConnector = query
	return srq
}

// WithComposition tells the query-builder to eager-load the nodes that are connected to
// the "composition" edge. The optional arguments are used to configure the query builder of the edge.
func (srq *ServiceResourceQuery) WithComposition(opts ...func(*ServiceResourceQuery)) *ServiceResourceQuery {
	query := (&ServiceResourceClient{config: srq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	srq.withComposition = query
	return srq
}

// WithComponents tells the query-builder to eager-load the nodes that are connected to
// the "components" edge. The optional arguments are used to configure the query builder of the edge.
func (srq *ServiceResourceQuery) WithComponents(opts ...func(*ServiceResourceQuery)) *ServiceResourceQuery {
	query := (&ServiceResourceClient{config: srq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	srq.withComponents = query
	return srq
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
//	client.ServiceResource.Query().
//		GroupBy(serviceresource.FieldCreateTime).
//		Aggregate(model.Count()).
//		Scan(ctx, &v)
func (srq *ServiceResourceQuery) GroupBy(field string, fields ...string) *ServiceResourceGroupBy {
	srq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &ServiceResourceGroupBy{build: srq}
	grbuild.flds = &srq.ctx.Fields
	grbuild.label = serviceresource.Label
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
//	client.ServiceResource.Query().
//		Select(serviceresource.FieldCreateTime).
//		Scan(ctx, &v)
func (srq *ServiceResourceQuery) Select(fields ...string) *ServiceResourceSelect {
	srq.ctx.Fields = append(srq.ctx.Fields, fields...)
	sbuild := &ServiceResourceSelect{ServiceResourceQuery: srq}
	sbuild.label = serviceresource.Label
	sbuild.flds, sbuild.scan = &srq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a ServiceResourceSelect configured with the given aggregations.
func (srq *ServiceResourceQuery) Aggregate(fns ...AggregateFunc) *ServiceResourceSelect {
	return srq.Select().Aggregate(fns...)
}

func (srq *ServiceResourceQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range srq.inters {
		if inter == nil {
			return fmt.Errorf("model: uninitialized interceptor (forgotten import model/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, srq); err != nil {
				return err
			}
		}
	}
	for _, f := range srq.ctx.Fields {
		if !serviceresource.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("model: invalid field %q for query", f)}
		}
	}
	if srq.path != nil {
		prev, err := srq.path(ctx)
		if err != nil {
			return err
		}
		srq.sql = prev
	}
	return nil
}

func (srq *ServiceResourceQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*ServiceResource, error) {
	var (
		nodes       = []*ServiceResource{}
		_spec       = srq.querySpec()
		loadedTypes = [4]bool{
			srq.withService != nil,
			srq.withConnector != nil,
			srq.withComposition != nil,
			srq.withComponents != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*ServiceResource).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &ServiceResource{config: srq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	_spec.Node.Schema = srq.schemaConfig.ServiceResource
	ctx = internal.NewSchemaConfigContext(ctx, srq.schemaConfig)
	if len(srq.modifiers) > 0 {
		_spec.Modifiers = srq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, srq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := srq.withService; query != nil {
		if err := srq.loadService(ctx, query, nodes, nil,
			func(n *ServiceResource, e *Service) { n.Edges.Service = e }); err != nil {
			return nil, err
		}
	}
	if query := srq.withConnector; query != nil {
		if err := srq.loadConnector(ctx, query, nodes, nil,
			func(n *ServiceResource, e *Connector) { n.Edges.Connector = e }); err != nil {
			return nil, err
		}
	}
	if query := srq.withComposition; query != nil {
		if err := srq.loadComposition(ctx, query, nodes, nil,
			func(n *ServiceResource, e *ServiceResource) { n.Edges.Composition = e }); err != nil {
			return nil, err
		}
	}
	if query := srq.withComponents; query != nil {
		if err := srq.loadComponents(ctx, query, nodes,
			func(n *ServiceResource) { n.Edges.Components = []*ServiceResource{} },
			func(n *ServiceResource, e *ServiceResource) { n.Edges.Components = append(n.Edges.Components, e) }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (srq *ServiceResourceQuery) loadService(ctx context.Context, query *ServiceQuery, nodes []*ServiceResource, init func(*ServiceResource), assign func(*ServiceResource, *Service)) error {
	ids := make([]oid.ID, 0, len(nodes))
	nodeids := make(map[oid.ID][]*ServiceResource)
	for i := range nodes {
		fk := nodes[i].ServiceID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(service.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "serviceID" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (srq *ServiceResourceQuery) loadConnector(ctx context.Context, query *ConnectorQuery, nodes []*ServiceResource, init func(*ServiceResource), assign func(*ServiceResource, *Connector)) error {
	ids := make([]oid.ID, 0, len(nodes))
	nodeids := make(map[oid.ID][]*ServiceResource)
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
func (srq *ServiceResourceQuery) loadComposition(ctx context.Context, query *ServiceResourceQuery, nodes []*ServiceResource, init func(*ServiceResource), assign func(*ServiceResource, *ServiceResource)) error {
	ids := make([]oid.ID, 0, len(nodes))
	nodeids := make(map[oid.ID][]*ServiceResource)
	for i := range nodes {
		fk := nodes[i].CompositionID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(serviceresource.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "compositionID" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (srq *ServiceResourceQuery) loadComponents(ctx context.Context, query *ServiceResourceQuery, nodes []*ServiceResource, init func(*ServiceResource), assign func(*ServiceResource, *ServiceResource)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[oid.ID]*ServiceResource)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(serviceresource.FieldCompositionID)
	}
	query.Where(predicate.ServiceResource(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(serviceresource.ComponentsColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.CompositionID
		node, ok := nodeids[fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "compositionID" returned %v for node %v`, fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}

func (srq *ServiceResourceQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := srq.querySpec()
	_spec.Node.Schema = srq.schemaConfig.ServiceResource
	ctx = internal.NewSchemaConfigContext(ctx, srq.schemaConfig)
	if len(srq.modifiers) > 0 {
		_spec.Modifiers = srq.modifiers
	}
	_spec.Node.Columns = srq.ctx.Fields
	if len(srq.ctx.Fields) > 0 {
		_spec.Unique = srq.ctx.Unique != nil && *srq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, srq.driver, _spec)
}

func (srq *ServiceResourceQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(serviceresource.Table, serviceresource.Columns, sqlgraph.NewFieldSpec(serviceresource.FieldID, field.TypeString))
	_spec.From = srq.sql
	if unique := srq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if srq.path != nil {
		_spec.Unique = true
	}
	if fields := srq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, serviceresource.FieldID)
		for i := range fields {
			if fields[i] != serviceresource.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
		if srq.withService != nil {
			_spec.Node.AddColumnOnce(serviceresource.FieldServiceID)
		}
		if srq.withConnector != nil {
			_spec.Node.AddColumnOnce(serviceresource.FieldConnectorID)
		}
		if srq.withComposition != nil {
			_spec.Node.AddColumnOnce(serviceresource.FieldCompositionID)
		}
	}
	if ps := srq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := srq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := srq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := srq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (srq *ServiceResourceQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(srq.driver.Dialect())
	t1 := builder.Table(serviceresource.Table)
	columns := srq.ctx.Fields
	if len(columns) == 0 {
		columns = serviceresource.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if srq.sql != nil {
		selector = srq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if srq.ctx.Unique != nil && *srq.ctx.Unique {
		selector.Distinct()
	}
	t1.Schema(srq.schemaConfig.ServiceResource)
	ctx = internal.NewSchemaConfigContext(ctx, srq.schemaConfig)
	selector.WithContext(ctx)
	for _, m := range srq.modifiers {
		m(selector)
	}
	for _, p := range srq.predicates {
		p(selector)
	}
	for _, p := range srq.order {
		p(selector)
	}
	if offset := srq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := srq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ForUpdate locks the selected rows against concurrent updates, and prevent them from being
// updated, deleted or "selected ... for update" by other sessions, until the transaction is
// either committed or rolled-back.
func (srq *ServiceResourceQuery) ForUpdate(opts ...sql.LockOption) *ServiceResourceQuery {
	if srq.driver.Dialect() == dialect.Postgres {
		srq.Unique(false)
	}
	srq.modifiers = append(srq.modifiers, func(s *sql.Selector) {
		s.ForUpdate(opts...)
	})
	return srq
}

// ForShare behaves similarly to ForUpdate, except that it acquires a shared mode lock
// on any rows that are read. Other sessions can read the rows, but cannot modify them
// until your transaction commits.
func (srq *ServiceResourceQuery) ForShare(opts ...sql.LockOption) *ServiceResourceQuery {
	if srq.driver.Dialect() == dialect.Postgres {
		srq.Unique(false)
	}
	srq.modifiers = append(srq.modifiers, func(s *sql.Selector) {
		s.ForShare(opts...)
	})
	return srq
}

// Modify adds a query modifier for attaching custom logic to queries.
func (srq *ServiceResourceQuery) Modify(modifiers ...func(s *sql.Selector)) *ServiceResourceSelect {
	srq.modifiers = append(srq.modifiers, modifiers...)
	return srq.Select()
}

// WhereP appends storage-level predicates to the ServiceResourceQuery builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (srq *ServiceResourceQuery) WhereP(ps ...func(*sql.Selector)) {
	var wps = make([]predicate.ServiceResource, 0, len(ps))
	for i := 0; i < len(ps); i++ {
		wps = append(wps, predicate.ServiceResource(ps[i]))
	}
	srq.predicates = append(srq.predicates, wps...)
}

// ServiceResourceGroupBy is the group-by builder for ServiceResource entities.
type ServiceResourceGroupBy struct {
	selector
	build *ServiceResourceQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (srgb *ServiceResourceGroupBy) Aggregate(fns ...AggregateFunc) *ServiceResourceGroupBy {
	srgb.fns = append(srgb.fns, fns...)
	return srgb
}

// Scan applies the selector query and scans the result into the given value.
func (srgb *ServiceResourceGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, srgb.build.ctx, "GroupBy")
	if err := srgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*ServiceResourceQuery, *ServiceResourceGroupBy](ctx, srgb.build, srgb, srgb.build.inters, v)
}

func (srgb *ServiceResourceGroupBy) sqlScan(ctx context.Context, root *ServiceResourceQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(srgb.fns))
	for _, fn := range srgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*srgb.flds)+len(srgb.fns))
		for _, f := range *srgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*srgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := srgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// ServiceResourceSelect is the builder for selecting fields of ServiceResource entities.
type ServiceResourceSelect struct {
	*ServiceResourceQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (srs *ServiceResourceSelect) Aggregate(fns ...AggregateFunc) *ServiceResourceSelect {
	srs.fns = append(srs.fns, fns...)
	return srs
}

// Scan applies the selector query and scans the result into the given value.
func (srs *ServiceResourceSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, srs.ctx, "Select")
	if err := srs.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*ServiceResourceQuery, *ServiceResourceSelect](ctx, srs.ServiceResourceQuery, srs, srs.inters, v)
}

func (srs *ServiceResourceSelect) sqlScan(ctx context.Context, root *ServiceResourceQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(srs.fns))
	for _, fn := range srs.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*srs.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := srs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (srs *ServiceResourceSelect) Modify(modifiers ...func(s *sql.Selector)) *ServiceResourceSelect {
	srs.modifiers = append(srs.modifiers, modifiers...)
	return srs
}
