// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// Code generated by "walrus". DO NOT EDIT.

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

	"github.com/seal-io/walrus/pkg/dao/model/environment"
	"github.com/seal-io/walrus/pkg/dao/model/internal"
	"github.com/seal-io/walrus/pkg/dao/model/predicate"
	"github.com/seal-io/walrus/pkg/dao/model/project"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/model/resourcecomponent"
	"github.com/seal-io/walrus/pkg/dao/model/resourcerelationship"
	"github.com/seal-io/walrus/pkg/dao/model/resourcerevision"
	"github.com/seal-io/walrus/pkg/dao/model/templateversion"
	"github.com/seal-io/walrus/pkg/dao/types/object"
)

// ResourceQuery is the builder for querying Resource entities.
type ResourceQuery struct {
	config
	ctx              *QueryContext
	order            []resource.OrderOption
	inters           []Interceptor
	predicates       []predicate.Resource
	withProject      *ProjectQuery
	withEnvironment  *EnvironmentQuery
	withTemplate     *TemplateVersionQuery
	withRevisions    *ResourceRevisionQuery
	withComponents   *ResourceComponentQuery
	withDependencies *ResourceRelationshipQuery
	modifiers        []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the ResourceQuery builder.
func (rq *ResourceQuery) Where(ps ...predicate.Resource) *ResourceQuery {
	rq.predicates = append(rq.predicates, ps...)
	return rq
}

// Limit the number of records to be returned by this query.
func (rq *ResourceQuery) Limit(limit int) *ResourceQuery {
	rq.ctx.Limit = &limit
	return rq
}

// Offset to start from.
func (rq *ResourceQuery) Offset(offset int) *ResourceQuery {
	rq.ctx.Offset = &offset
	return rq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (rq *ResourceQuery) Unique(unique bool) *ResourceQuery {
	rq.ctx.Unique = &unique
	return rq
}

// Order specifies how the records should be ordered.
func (rq *ResourceQuery) Order(o ...resource.OrderOption) *ResourceQuery {
	rq.order = append(rq.order, o...)
	return rq
}

// QueryProject chains the current query on the "project" edge.
func (rq *ResourceQuery) QueryProject() *ProjectQuery {
	query := (&ProjectClient{config: rq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := rq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := rq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(resource.Table, resource.FieldID, selector),
			sqlgraph.To(project.Table, project.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, resource.ProjectTable, resource.ProjectColumn),
		)
		schemaConfig := rq.schemaConfig
		step.To.Schema = schemaConfig.Project
		step.Edge.Schema = schemaConfig.Resource
		fromU = sqlgraph.SetNeighbors(rq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryEnvironment chains the current query on the "environment" edge.
func (rq *ResourceQuery) QueryEnvironment() *EnvironmentQuery {
	query := (&EnvironmentClient{config: rq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := rq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := rq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(resource.Table, resource.FieldID, selector),
			sqlgraph.To(environment.Table, environment.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, resource.EnvironmentTable, resource.EnvironmentColumn),
		)
		schemaConfig := rq.schemaConfig
		step.To.Schema = schemaConfig.Environment
		step.Edge.Schema = schemaConfig.Resource
		fromU = sqlgraph.SetNeighbors(rq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryTemplate chains the current query on the "template" edge.
func (rq *ResourceQuery) QueryTemplate() *TemplateVersionQuery {
	query := (&TemplateVersionClient{config: rq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := rq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := rq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(resource.Table, resource.FieldID, selector),
			sqlgraph.To(templateversion.Table, templateversion.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, resource.TemplateTable, resource.TemplateColumn),
		)
		schemaConfig := rq.schemaConfig
		step.To.Schema = schemaConfig.TemplateVersion
		step.Edge.Schema = schemaConfig.Resource
		fromU = sqlgraph.SetNeighbors(rq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryRevisions chains the current query on the "revisions" edge.
func (rq *ResourceQuery) QueryRevisions() *ResourceRevisionQuery {
	query := (&ResourceRevisionClient{config: rq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := rq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := rq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(resource.Table, resource.FieldID, selector),
			sqlgraph.To(resourcerevision.Table, resourcerevision.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, resource.RevisionsTable, resource.RevisionsColumn),
		)
		schemaConfig := rq.schemaConfig
		step.To.Schema = schemaConfig.ResourceRevision
		step.Edge.Schema = schemaConfig.ResourceRevision
		fromU = sqlgraph.SetNeighbors(rq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryComponents chains the current query on the "components" edge.
func (rq *ResourceQuery) QueryComponents() *ResourceComponentQuery {
	query := (&ResourceComponentClient{config: rq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := rq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := rq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(resource.Table, resource.FieldID, selector),
			sqlgraph.To(resourcecomponent.Table, resourcecomponent.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, resource.ComponentsTable, resource.ComponentsColumn),
		)
		schemaConfig := rq.schemaConfig
		step.To.Schema = schemaConfig.ResourceComponent
		step.Edge.Schema = schemaConfig.ResourceComponent
		fromU = sqlgraph.SetNeighbors(rq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryDependencies chains the current query on the "dependencies" edge.
func (rq *ResourceQuery) QueryDependencies() *ResourceRelationshipQuery {
	query := (&ResourceRelationshipClient{config: rq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := rq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := rq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(resource.Table, resource.FieldID, selector),
			sqlgraph.To(resourcerelationship.Table, resourcerelationship.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, resource.DependenciesTable, resource.DependenciesColumn),
		)
		schemaConfig := rq.schemaConfig
		step.To.Schema = schemaConfig.ResourceRelationship
		step.Edge.Schema = schemaConfig.ResourceRelationship
		fromU = sqlgraph.SetNeighbors(rq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Resource entity from the query.
// Returns a *NotFoundError when no Resource was found.
func (rq *ResourceQuery) First(ctx context.Context) (*Resource, error) {
	nodes, err := rq.Limit(1).All(setContextOp(ctx, rq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{resource.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (rq *ResourceQuery) FirstX(ctx context.Context) *Resource {
	node, err := rq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Resource ID from the query.
// Returns a *NotFoundError when no Resource ID was found.
func (rq *ResourceQuery) FirstID(ctx context.Context) (id object.ID, err error) {
	var ids []object.ID
	if ids, err = rq.Limit(1).IDs(setContextOp(ctx, rq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{resource.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (rq *ResourceQuery) FirstIDX(ctx context.Context) object.ID {
	id, err := rq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Resource entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Resource entity is found.
// Returns a *NotFoundError when no Resource entities are found.
func (rq *ResourceQuery) Only(ctx context.Context) (*Resource, error) {
	nodes, err := rq.Limit(2).All(setContextOp(ctx, rq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{resource.Label}
	default:
		return nil, &NotSingularError{resource.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (rq *ResourceQuery) OnlyX(ctx context.Context) *Resource {
	node, err := rq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Resource ID in the query.
// Returns a *NotSingularError when more than one Resource ID is found.
// Returns a *NotFoundError when no entities are found.
func (rq *ResourceQuery) OnlyID(ctx context.Context) (id object.ID, err error) {
	var ids []object.ID
	if ids, err = rq.Limit(2).IDs(setContextOp(ctx, rq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{resource.Label}
	default:
		err = &NotSingularError{resource.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (rq *ResourceQuery) OnlyIDX(ctx context.Context) object.ID {
	id, err := rq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Resources.
func (rq *ResourceQuery) All(ctx context.Context) ([]*Resource, error) {
	ctx = setContextOp(ctx, rq.ctx, "All")
	if err := rq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Resource, *ResourceQuery]()
	return withInterceptors[[]*Resource](ctx, rq, qr, rq.inters)
}

// AllX is like All, but panics if an error occurs.
func (rq *ResourceQuery) AllX(ctx context.Context) []*Resource {
	nodes, err := rq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Resource IDs.
func (rq *ResourceQuery) IDs(ctx context.Context) (ids []object.ID, err error) {
	if rq.ctx.Unique == nil && rq.path != nil {
		rq.Unique(true)
	}
	ctx = setContextOp(ctx, rq.ctx, "IDs")
	if err = rq.Select(resource.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (rq *ResourceQuery) IDsX(ctx context.Context) []object.ID {
	ids, err := rq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (rq *ResourceQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, rq.ctx, "Count")
	if err := rq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, rq, querierCount[*ResourceQuery](), rq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (rq *ResourceQuery) CountX(ctx context.Context) int {
	count, err := rq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (rq *ResourceQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, rq.ctx, "Exist")
	switch _, err := rq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("model: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (rq *ResourceQuery) ExistX(ctx context.Context) bool {
	exist, err := rq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the ResourceQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (rq *ResourceQuery) Clone() *ResourceQuery {
	if rq == nil {
		return nil
	}
	return &ResourceQuery{
		config:           rq.config,
		ctx:              rq.ctx.Clone(),
		order:            append([]resource.OrderOption{}, rq.order...),
		inters:           append([]Interceptor{}, rq.inters...),
		predicates:       append([]predicate.Resource{}, rq.predicates...),
		withProject:      rq.withProject.Clone(),
		withEnvironment:  rq.withEnvironment.Clone(),
		withTemplate:     rq.withTemplate.Clone(),
		withRevisions:    rq.withRevisions.Clone(),
		withComponents:   rq.withComponents.Clone(),
		withDependencies: rq.withDependencies.Clone(),
		// clone intermediate query.
		sql:  rq.sql.Clone(),
		path: rq.path,
	}
}

// WithProject tells the query-builder to eager-load the nodes that are connected to
// the "project" edge. The optional arguments are used to configure the query builder of the edge.
func (rq *ResourceQuery) WithProject(opts ...func(*ProjectQuery)) *ResourceQuery {
	query := (&ProjectClient{config: rq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	rq.withProject = query
	return rq
}

// WithEnvironment tells the query-builder to eager-load the nodes that are connected to
// the "environment" edge. The optional arguments are used to configure the query builder of the edge.
func (rq *ResourceQuery) WithEnvironment(opts ...func(*EnvironmentQuery)) *ResourceQuery {
	query := (&EnvironmentClient{config: rq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	rq.withEnvironment = query
	return rq
}

// WithTemplate tells the query-builder to eager-load the nodes that are connected to
// the "template" edge. The optional arguments are used to configure the query builder of the edge.
func (rq *ResourceQuery) WithTemplate(opts ...func(*TemplateVersionQuery)) *ResourceQuery {
	query := (&TemplateVersionClient{config: rq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	rq.withTemplate = query
	return rq
}

// WithRevisions tells the query-builder to eager-load the nodes that are connected to
// the "revisions" edge. The optional arguments are used to configure the query builder of the edge.
func (rq *ResourceQuery) WithRevisions(opts ...func(*ResourceRevisionQuery)) *ResourceQuery {
	query := (&ResourceRevisionClient{config: rq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	rq.withRevisions = query
	return rq
}

// WithComponents tells the query-builder to eager-load the nodes that are connected to
// the "components" edge. The optional arguments are used to configure the query builder of the edge.
func (rq *ResourceQuery) WithComponents(opts ...func(*ResourceComponentQuery)) *ResourceQuery {
	query := (&ResourceComponentClient{config: rq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	rq.withComponents = query
	return rq
}

// WithDependencies tells the query-builder to eager-load the nodes that are connected to
// the "dependencies" edge. The optional arguments are used to configure the query builder of the edge.
func (rq *ResourceQuery) WithDependencies(opts ...func(*ResourceRelationshipQuery)) *ResourceQuery {
	query := (&ResourceRelationshipClient{config: rq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	rq.withDependencies = query
	return rq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Name string `json:"name,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Resource.Query().
//		GroupBy(resource.FieldName).
//		Aggregate(model.Count()).
//		Scan(ctx, &v)
func (rq *ResourceQuery) GroupBy(field string, fields ...string) *ResourceGroupBy {
	rq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &ResourceGroupBy{build: rq}
	grbuild.flds = &rq.ctx.Fields
	grbuild.label = resource.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Name string `json:"name,omitempty"`
//	}
//
//	client.Resource.Query().
//		Select(resource.FieldName).
//		Scan(ctx, &v)
func (rq *ResourceQuery) Select(fields ...string) *ResourceSelect {
	rq.ctx.Fields = append(rq.ctx.Fields, fields...)
	sbuild := &ResourceSelect{ResourceQuery: rq}
	sbuild.label = resource.Label
	sbuild.flds, sbuild.scan = &rq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a ResourceSelect configured with the given aggregations.
func (rq *ResourceQuery) Aggregate(fns ...AggregateFunc) *ResourceSelect {
	return rq.Select().Aggregate(fns...)
}

func (rq *ResourceQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range rq.inters {
		if inter == nil {
			return fmt.Errorf("model: uninitialized interceptor (forgotten import model/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, rq); err != nil {
				return err
			}
		}
	}
	for _, f := range rq.ctx.Fields {
		if !resource.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("model: invalid field %q for query", f)}
		}
	}
	if rq.path != nil {
		prev, err := rq.path(ctx)
		if err != nil {
			return err
		}
		rq.sql = prev
	}
	return nil
}

func (rq *ResourceQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Resource, error) {
	var (
		nodes       = []*Resource{}
		_spec       = rq.querySpec()
		loadedTypes = [6]bool{
			rq.withProject != nil,
			rq.withEnvironment != nil,
			rq.withTemplate != nil,
			rq.withRevisions != nil,
			rq.withComponents != nil,
			rq.withDependencies != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Resource).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Resource{config: rq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	_spec.Node.Schema = rq.schemaConfig.Resource
	ctx = internal.NewSchemaConfigContext(ctx, rq.schemaConfig)
	if len(rq.modifiers) > 0 {
		_spec.Modifiers = rq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, rq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := rq.withProject; query != nil {
		if err := rq.loadProject(ctx, query, nodes, nil,
			func(n *Resource, e *Project) { n.Edges.Project = e }); err != nil {
			return nil, err
		}
	}
	if query := rq.withEnvironment; query != nil {
		if err := rq.loadEnvironment(ctx, query, nodes, nil,
			func(n *Resource, e *Environment) { n.Edges.Environment = e }); err != nil {
			return nil, err
		}
	}
	if query := rq.withTemplate; query != nil {
		if err := rq.loadTemplate(ctx, query, nodes, nil,
			func(n *Resource, e *TemplateVersion) { n.Edges.Template = e }); err != nil {
			return nil, err
		}
	}
	if query := rq.withRevisions; query != nil {
		if err := rq.loadRevisions(ctx, query, nodes,
			func(n *Resource) { n.Edges.Revisions = []*ResourceRevision{} },
			func(n *Resource, e *ResourceRevision) { n.Edges.Revisions = append(n.Edges.Revisions, e) }); err != nil {
			return nil, err
		}
	}
	if query := rq.withComponents; query != nil {
		if err := rq.loadComponents(ctx, query, nodes,
			func(n *Resource) { n.Edges.Components = []*ResourceComponent{} },
			func(n *Resource, e *ResourceComponent) { n.Edges.Components = append(n.Edges.Components, e) }); err != nil {
			return nil, err
		}
	}
	if query := rq.withDependencies; query != nil {
		if err := rq.loadDependencies(ctx, query, nodes,
			func(n *Resource) { n.Edges.Dependencies = []*ResourceRelationship{} },
			func(n *Resource, e *ResourceRelationship) { n.Edges.Dependencies = append(n.Edges.Dependencies, e) }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (rq *ResourceQuery) loadProject(ctx context.Context, query *ProjectQuery, nodes []*Resource, init func(*Resource), assign func(*Resource, *Project)) error {
	ids := make([]object.ID, 0, len(nodes))
	nodeids := make(map[object.ID][]*Resource)
	for i := range nodes {
		fk := nodes[i].ProjectID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(project.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "project_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (rq *ResourceQuery) loadEnvironment(ctx context.Context, query *EnvironmentQuery, nodes []*Resource, init func(*Resource), assign func(*Resource, *Environment)) error {
	ids := make([]object.ID, 0, len(nodes))
	nodeids := make(map[object.ID][]*Resource)
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
			return fmt.Errorf(`unexpected foreign-key "environment_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (rq *ResourceQuery) loadTemplate(ctx context.Context, query *TemplateVersionQuery, nodes []*Resource, init func(*Resource), assign func(*Resource, *TemplateVersion)) error {
	ids := make([]object.ID, 0, len(nodes))
	nodeids := make(map[object.ID][]*Resource)
	for i := range nodes {
		fk := nodes[i].TemplateID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(templateversion.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "template_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (rq *ResourceQuery) loadRevisions(ctx context.Context, query *ResourceRevisionQuery, nodes []*Resource, init func(*Resource), assign func(*Resource, *ResourceRevision)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[object.ID]*Resource)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(resourcerevision.FieldResourceID)
	}
	query.Where(predicate.ResourceRevision(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(resource.RevisionsColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.ResourceID
		node, ok := nodeids[fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "resource_id" returned %v for node %v`, fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}
func (rq *ResourceQuery) loadComponents(ctx context.Context, query *ResourceComponentQuery, nodes []*Resource, init func(*Resource), assign func(*Resource, *ResourceComponent)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[object.ID]*Resource)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(resourcecomponent.FieldResourceID)
	}
	query.Where(predicate.ResourceComponent(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(resource.ComponentsColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.ResourceID
		node, ok := nodeids[fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "resource_id" returned %v for node %v`, fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}
func (rq *ResourceQuery) loadDependencies(ctx context.Context, query *ResourceRelationshipQuery, nodes []*Resource, init func(*Resource), assign func(*Resource, *ResourceRelationship)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[object.ID]*Resource)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(resourcerelationship.FieldResourceID)
	}
	query.Where(predicate.ResourceRelationship(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(resource.DependenciesColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.ResourceID
		node, ok := nodeids[fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "resource_id" returned %v for node %v`, fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}

func (rq *ResourceQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := rq.querySpec()
	_spec.Node.Schema = rq.schemaConfig.Resource
	ctx = internal.NewSchemaConfigContext(ctx, rq.schemaConfig)
	if len(rq.modifiers) > 0 {
		_spec.Modifiers = rq.modifiers
	}
	_spec.Node.Columns = rq.ctx.Fields
	if len(rq.ctx.Fields) > 0 {
		_spec.Unique = rq.ctx.Unique != nil && *rq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, rq.driver, _spec)
}

func (rq *ResourceQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(resource.Table, resource.Columns, sqlgraph.NewFieldSpec(resource.FieldID, field.TypeString))
	_spec.From = rq.sql
	if unique := rq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if rq.path != nil {
		_spec.Unique = true
	}
	if fields := rq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, resource.FieldID)
		for i := range fields {
			if fields[i] != resource.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
		if rq.withProject != nil {
			_spec.Node.AddColumnOnce(resource.FieldProjectID)
		}
		if rq.withEnvironment != nil {
			_spec.Node.AddColumnOnce(resource.FieldEnvironmentID)
		}
		if rq.withTemplate != nil {
			_spec.Node.AddColumnOnce(resource.FieldTemplateID)
		}
	}
	if ps := rq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := rq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := rq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := rq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (rq *ResourceQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(rq.driver.Dialect())
	t1 := builder.Table(resource.Table)
	columns := rq.ctx.Fields
	if len(columns) == 0 {
		columns = resource.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if rq.sql != nil {
		selector = rq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if rq.ctx.Unique != nil && *rq.ctx.Unique {
		selector.Distinct()
	}
	t1.Schema(rq.schemaConfig.Resource)
	ctx = internal.NewSchemaConfigContext(ctx, rq.schemaConfig)
	selector.WithContext(ctx)
	for _, m := range rq.modifiers {
		m(selector)
	}
	for _, p := range rq.predicates {
		p(selector)
	}
	for _, p := range rq.order {
		p(selector)
	}
	if offset := rq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := rq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ForUpdate locks the selected rows against concurrent updates, and prevent them from being
// updated, deleted or "selected ... for update" by other sessions, until the transaction is
// either committed or rolled-back.
func (rq *ResourceQuery) ForUpdate(opts ...sql.LockOption) *ResourceQuery {
	if rq.driver.Dialect() == dialect.Postgres {
		rq.Unique(false)
	}
	rq.modifiers = append(rq.modifiers, func(s *sql.Selector) {
		s.ForUpdate(opts...)
	})
	return rq
}

// ForShare behaves similarly to ForUpdate, except that it acquires a shared mode lock
// on any rows that are read. Other sessions can read the rows, but cannot modify them
// until your transaction commits.
func (rq *ResourceQuery) ForShare(opts ...sql.LockOption) *ResourceQuery {
	if rq.driver.Dialect() == dialect.Postgres {
		rq.Unique(false)
	}
	rq.modifiers = append(rq.modifiers, func(s *sql.Selector) {
		s.ForShare(opts...)
	})
	return rq
}

// Modify adds a query modifier for attaching custom logic to queries.
func (rq *ResourceQuery) Modify(modifiers ...func(s *sql.Selector)) *ResourceSelect {
	rq.modifiers = append(rq.modifiers, modifiers...)
	return rq.Select()
}

// WhereP appends storage-level predicates to the ResourceQuery builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (rq *ResourceQuery) WhereP(ps ...func(*sql.Selector)) {
	var wps = make([]predicate.Resource, 0, len(ps))
	for i := 0; i < len(ps); i++ {
		wps = append(wps, predicate.Resource(ps[i]))
	}
	rq.predicates = append(rq.predicates, wps...)
}

// ResourceGroupBy is the group-by builder for Resource entities.
type ResourceGroupBy struct {
	selector
	build *ResourceQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (rgb *ResourceGroupBy) Aggregate(fns ...AggregateFunc) *ResourceGroupBy {
	rgb.fns = append(rgb.fns, fns...)
	return rgb
}

// Scan applies the selector query and scans the result into the given value.
func (rgb *ResourceGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, rgb.build.ctx, "GroupBy")
	if err := rgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*ResourceQuery, *ResourceGroupBy](ctx, rgb.build, rgb, rgb.build.inters, v)
}

func (rgb *ResourceGroupBy) sqlScan(ctx context.Context, root *ResourceQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(rgb.fns))
	for _, fn := range rgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*rgb.flds)+len(rgb.fns))
		for _, f := range *rgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*rgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := rgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// ResourceSelect is the builder for selecting fields of Resource entities.
type ResourceSelect struct {
	*ResourceQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (rs *ResourceSelect) Aggregate(fns ...AggregateFunc) *ResourceSelect {
	rs.fns = append(rs.fns, fns...)
	return rs
}

// Scan applies the selector query and scans the result into the given value.
func (rs *ResourceSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, rs.ctx, "Select")
	if err := rs.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*ResourceQuery, *ResourceSelect](ctx, rs.ResourceQuery, rs, rs.inters, v)
}

func (rs *ResourceSelect) sqlScan(ctx context.Context, root *ResourceQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(rs.fns))
	for _, fn := range rs.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*rs.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := rs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (rs *ResourceSelect) Modify(modifiers ...func(s *sql.Selector)) *ResourceSelect {
	rs.modifiers = append(rs.modifiers, modifiers...)
	return rs
}
