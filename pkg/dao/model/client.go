// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package model

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/seal-io/seal/pkg/dao/model/migrate"
	"github.com/seal-io/seal/pkg/dao/types"

	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/model/applicationmodulerelationship"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/model/applicationrevision"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/environmentconnectorrelationship"
	"github.com/seal-io/seal/pkg/dao/model/module"
	"github.com/seal-io/seal/pkg/dao/model/project"
	"github.com/seal-io/seal/pkg/dao/model/role"
	"github.com/seal-io/seal/pkg/dao/model/setting"
	"github.com/seal-io/seal/pkg/dao/model/subject"
	"github.com/seal-io/seal/pkg/dao/model/token"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

// Client is the client that holds all ent builders.
type Client struct {
	config
	// Schema is the client for creating, migrating and dropping schema.
	Schema *migrate.Schema
	// Application is the client for interacting with the Application builders.
	Application *ApplicationClient
	// ApplicationModuleRelationship is the client for interacting with the ApplicationModuleRelationship builders.
	ApplicationModuleRelationship *ApplicationModuleRelationshipClient
	// ApplicationResource is the client for interacting with the ApplicationResource builders.
	ApplicationResource *ApplicationResourceClient
	// ApplicationRevision is the client for interacting with the ApplicationRevision builders.
	ApplicationRevision *ApplicationRevisionClient
	// Connector is the client for interacting with the Connector builders.
	Connector *ConnectorClient
	// Environment is the client for interacting with the Environment builders.
	Environment *EnvironmentClient
	// EnvironmentConnectorRelationship is the client for interacting with the EnvironmentConnectorRelationship builders.
	EnvironmentConnectorRelationship *EnvironmentConnectorRelationshipClient
	// Module is the client for interacting with the Module builders.
	Module *ModuleClient
	// Project is the client for interacting with the Project builders.
	Project *ProjectClient
	// Role is the client for interacting with the Role builders.
	Role *RoleClient
	// Setting is the client for interacting with the Setting builders.
	Setting *SettingClient
	// Subject is the client for interacting with the Subject builders.
	Subject *SubjectClient
	// Token is the client for interacting with the Token builders.
	Token *TokenClient
}

// NewClient creates a new client configured with the given options.
func NewClient(opts ...Option) *Client {
	cfg := config{log: log.Println, hooks: &hooks{}, inters: &inters{}}
	cfg.options(opts...)
	client := &Client{config: cfg}
	client.init()
	return client
}

func (c *Client) init() {
	c.Schema = migrate.NewSchema(c.driver)
	c.Application = NewApplicationClient(c.config)
	c.ApplicationModuleRelationship = NewApplicationModuleRelationshipClient(c.config)
	c.ApplicationResource = NewApplicationResourceClient(c.config)
	c.ApplicationRevision = NewApplicationRevisionClient(c.config)
	c.Connector = NewConnectorClient(c.config)
	c.Environment = NewEnvironmentClient(c.config)
	c.EnvironmentConnectorRelationship = NewEnvironmentConnectorRelationshipClient(c.config)
	c.Module = NewModuleClient(c.config)
	c.Project = NewProjectClient(c.config)
	c.Role = NewRoleClient(c.config)
	c.Setting = NewSettingClient(c.config)
	c.Subject = NewSubjectClient(c.config)
	c.Token = NewTokenClient(c.config)
}

// Open opens a database/sql.DB specified by the driver name and
// the data source name, and returns a new client attached to it.
// Optional parameters can be added for configuring the client.
func Open(driverName, dataSourceName string, options ...Option) (*Client, error) {
	switch driverName {
	case dialect.MySQL, dialect.Postgres, dialect.SQLite:
		drv, err := sql.Open(driverName, dataSourceName)
		if err != nil {
			return nil, err
		}
		return NewClient(append(options, Driver(drv))...), nil
	default:
		return nil, fmt.Errorf("unsupported driver: %q", driverName)
	}
}

// Tx returns a new transactional client. The provided context
// is used until the transaction is committed or rolled back.
func (c *Client) Tx(ctx context.Context) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, errors.New("model: cannot start a transaction within a transaction")
	}
	tx, err := newTx(ctx, c.driver)
	if err != nil {
		return nil, fmt.Errorf("model: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = tx
	return &Tx{
		ctx:                              ctx,
		config:                           cfg,
		Application:                      NewApplicationClient(cfg),
		ApplicationModuleRelationship:    NewApplicationModuleRelationshipClient(cfg),
		ApplicationResource:              NewApplicationResourceClient(cfg),
		ApplicationRevision:              NewApplicationRevisionClient(cfg),
		Connector:                        NewConnectorClient(cfg),
		Environment:                      NewEnvironmentClient(cfg),
		EnvironmentConnectorRelationship: NewEnvironmentConnectorRelationshipClient(cfg),
		Module:                           NewModuleClient(cfg),
		Project:                          NewProjectClient(cfg),
		Role:                             NewRoleClient(cfg),
		Setting:                          NewSettingClient(cfg),
		Subject:                          NewSubjectClient(cfg),
		Token:                            NewTokenClient(cfg),
	}, nil
}

// BeginTx returns a transactional client with specified options.
func (c *Client) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, errors.New("ent: cannot start a transaction within a transaction")
	}
	tx, err := c.driver.(interface {
		BeginTx(context.Context, *sql.TxOptions) (dialect.Tx, error)
	}).BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = &txDriver{tx: tx, drv: c.driver}
	return &Tx{
		ctx:                              ctx,
		config:                           cfg,
		Application:                      NewApplicationClient(cfg),
		ApplicationModuleRelationship:    NewApplicationModuleRelationshipClient(cfg),
		ApplicationResource:              NewApplicationResourceClient(cfg),
		ApplicationRevision:              NewApplicationRevisionClient(cfg),
		Connector:                        NewConnectorClient(cfg),
		Environment:                      NewEnvironmentClient(cfg),
		EnvironmentConnectorRelationship: NewEnvironmentConnectorRelationshipClient(cfg),
		Module:                           NewModuleClient(cfg),
		Project:                          NewProjectClient(cfg),
		Role:                             NewRoleClient(cfg),
		Setting:                          NewSettingClient(cfg),
		Subject:                          NewSubjectClient(cfg),
		Token:                            NewTokenClient(cfg),
	}, nil
}

// Debug returns a new debug-client. It's used to get verbose logging on specific operations.
//
//	client.Debug().
//		Application.
//		Query().
//		Count(ctx)
func (c *Client) Debug() *Client {
	if c.debug {
		return c
	}
	cfg := c.config
	cfg.driver = dialect.Debug(c.driver, c.log)
	client := &Client{config: cfg}
	client.init()
	return client
}

// Close closes the database connection and prevents new queries from starting.
func (c *Client) Close() error {
	return c.driver.Close()
}

// Use adds the mutation hooks to all the entity clients.
// In order to add hooks to a specific client, call: `client.Node.Use(...)`.
func (c *Client) Use(hooks ...Hook) {
	c.Application.Use(hooks...)
	c.ApplicationModuleRelationship.Use(hooks...)
	c.ApplicationResource.Use(hooks...)
	c.ApplicationRevision.Use(hooks...)
	c.Connector.Use(hooks...)
	c.Environment.Use(hooks...)
	c.EnvironmentConnectorRelationship.Use(hooks...)
	c.Module.Use(hooks...)
	c.Project.Use(hooks...)
	c.Role.Use(hooks...)
	c.Setting.Use(hooks...)
	c.Subject.Use(hooks...)
	c.Token.Use(hooks...)
}

// Intercept adds the query interceptors to all the entity clients.
// In order to add interceptors to a specific client, call: `client.Node.Intercept(...)`.
func (c *Client) Intercept(interceptors ...Interceptor) {
	c.Application.Intercept(interceptors...)
	c.ApplicationModuleRelationship.Intercept(interceptors...)
	c.ApplicationResource.Intercept(interceptors...)
	c.ApplicationRevision.Intercept(interceptors...)
	c.Connector.Intercept(interceptors...)
	c.Environment.Intercept(interceptors...)
	c.EnvironmentConnectorRelationship.Intercept(interceptors...)
	c.Module.Intercept(interceptors...)
	c.Project.Intercept(interceptors...)
	c.Role.Intercept(interceptors...)
	c.Setting.Intercept(interceptors...)
	c.Subject.Intercept(interceptors...)
	c.Token.Intercept(interceptors...)
}

// Applications implements the ClientSet.
func (c *Client) Applications() *ApplicationClient {
	return c.Application
}

// ApplicationModuleRelationships implements the ClientSet.
func (c *Client) ApplicationModuleRelationships() *ApplicationModuleRelationshipClient {
	return c.ApplicationModuleRelationship
}

// ApplicationResources implements the ClientSet.
func (c *Client) ApplicationResources() *ApplicationResourceClient {
	return c.ApplicationResource
}

// ApplicationRevisions implements the ClientSet.
func (c *Client) ApplicationRevisions() *ApplicationRevisionClient {
	return c.ApplicationRevision
}

// Connectors implements the ClientSet.
func (c *Client) Connectors() *ConnectorClient {
	return c.Connector
}

// Environments implements the ClientSet.
func (c *Client) Environments() *EnvironmentClient {
	return c.Environment
}

// EnvironmentConnectorRelationships implements the ClientSet.
func (c *Client) EnvironmentConnectorRelationships() *EnvironmentConnectorRelationshipClient {
	return c.EnvironmentConnectorRelationship
}

// Modules implements the ClientSet.
func (c *Client) Modules() *ModuleClient {
	return c.Module
}

// Projects implements the ClientSet.
func (c *Client) Projects() *ProjectClient {
	return c.Project
}

// Roles implements the ClientSet.
func (c *Client) Roles() *RoleClient {
	return c.Role
}

// Settings implements the ClientSet.
func (c *Client) Settings() *SettingClient {
	return c.Setting
}

// Subjects implements the ClientSet.
func (c *Client) Subjects() *SubjectClient {
	return c.Subject
}

// Tokens implements the ClientSet.
func (c *Client) Tokens() *TokenClient {
	return c.Token
}

// Dialect returns the dialect name of the driver.
func (c *Client) Dialect() string {
	return c.driver.Dialect()
}

// WithTx gives a new transactional client in the callback function,
// if already in a transaction, this will keep in the same transaction.
func (c *Client) WithTx(ctx context.Context, fn func(tx *Tx) error) (err error) {
	var tx *Tx
	tx, err = c.Tx(ctx)
	if err != nil {
		return
	}
	defer func() {
		if v := recover(); v != nil {
			switch vt := v.(type) {
			case error:
				err = fmt.Errorf("panic as %w", vt)
			default:
				err = fmt.Errorf("panic as %v", v)
			}
			if txErr := tx.Rollback(); txErr != nil {
				err = fmt.Errorf("try to rollback as received %w, but failed: %v", err, txErr)
			}
		}
	}()
	if err = fn(tx); err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			err = fmt.Errorf("try to rollback as received %w, but failed: %v", err, txErr)
		}
		return
	}
	if txErr := tx.Commit(); txErr != nil {
		err = fmt.Errorf("try to commit, but failed: %v", txErr)
	}
	return
}

// Mutate implements the ent.Mutator interface.
func (c *Client) Mutate(ctx context.Context, m Mutation) (Value, error) {
	switch m := m.(type) {
	case *ApplicationMutation:
		return c.Application.mutate(ctx, m)
	case *ApplicationModuleRelationshipMutation:
		return c.ApplicationModuleRelationship.mutate(ctx, m)
	case *ApplicationResourceMutation:
		return c.ApplicationResource.mutate(ctx, m)
	case *ApplicationRevisionMutation:
		return c.ApplicationRevision.mutate(ctx, m)
	case *ConnectorMutation:
		return c.Connector.mutate(ctx, m)
	case *EnvironmentMutation:
		return c.Environment.mutate(ctx, m)
	case *EnvironmentConnectorRelationshipMutation:
		return c.EnvironmentConnectorRelationship.mutate(ctx, m)
	case *ModuleMutation:
		return c.Module.mutate(ctx, m)
	case *ProjectMutation:
		return c.Project.mutate(ctx, m)
	case *RoleMutation:
		return c.Role.mutate(ctx, m)
	case *SettingMutation:
		return c.Setting.mutate(ctx, m)
	case *SubjectMutation:
		return c.Subject.mutate(ctx, m)
	case *TokenMutation:
		return c.Token.mutate(ctx, m)
	default:
		return nil, fmt.Errorf("model: unknown mutation type %T", m)
	}
}

// ApplicationClient is a client for the Application schema.
type ApplicationClient struct {
	config
}

// NewApplicationClient returns a client for the Application from the given config.
func NewApplicationClient(c config) *ApplicationClient {
	return &ApplicationClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `application.Hooks(f(g(h())))`.
func (c *ApplicationClient) Use(hooks ...Hook) {
	c.hooks.Application = append(c.hooks.Application, hooks...)
}

// Use adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `application.Intercept(f(g(h())))`.
func (c *ApplicationClient) Intercept(interceptors ...Interceptor) {
	c.inters.Application = append(c.inters.Application, interceptors...)
}

// Create returns a builder for creating a Application entity.
func (c *ApplicationClient) Create() *ApplicationCreate {
	mutation := newApplicationMutation(c.config, OpCreate)
	return &ApplicationCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Application entities.
func (c *ApplicationClient) CreateBulk(builders ...*ApplicationCreate) *ApplicationCreateBulk {
	return &ApplicationCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Application.
func (c *ApplicationClient) Update() *ApplicationUpdate {
	mutation := newApplicationMutation(c.config, OpUpdate)
	return &ApplicationUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *ApplicationClient) UpdateOne(a *Application) *ApplicationUpdateOne {
	mutation := newApplicationMutation(c.config, OpUpdateOne, withApplication(a))
	return &ApplicationUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *ApplicationClient) UpdateOneID(id types.ID) *ApplicationUpdateOne {
	mutation := newApplicationMutation(c.config, OpUpdateOne, withApplicationID(id))
	return &ApplicationUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Application.
func (c *ApplicationClient) Delete() *ApplicationDelete {
	mutation := newApplicationMutation(c.config, OpDelete)
	return &ApplicationDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *ApplicationClient) DeleteOne(a *Application) *ApplicationDeleteOne {
	return c.DeleteOneID(a.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *ApplicationClient) DeleteOneID(id types.ID) *ApplicationDeleteOne {
	builder := c.Delete().Where(application.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &ApplicationDeleteOne{builder}
}

// Query returns a query builder for Application.
func (c *ApplicationClient) Query() *ApplicationQuery {
	return &ApplicationQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeApplication},
		inters: c.Interceptors(),
	}
}

// Get returns a Application entity by its id.
func (c *ApplicationClient) Get(ctx context.Context, id types.ID) (*Application, error) {
	return c.Query().Where(application.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *ApplicationClient) GetX(ctx context.Context, id types.ID) *Application {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryProject queries the project edge of a Application.
func (c *ApplicationClient) QueryProject(a *Application) *ProjectQuery {
	query := (&ProjectClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := a.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(application.Table, application.FieldID, id),
			sqlgraph.To(project.Table, project.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, application.ProjectTable, application.ProjectColumn),
		)
		schemaConfig := a.schemaConfig
		step.To.Schema = schemaConfig.Project
		step.Edge.Schema = schemaConfig.Application
		fromV = sqlgraph.Neighbors(a.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryEnvironment queries the environment edge of a Application.
func (c *ApplicationClient) QueryEnvironment(a *Application) *EnvironmentQuery {
	query := (&EnvironmentClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := a.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(application.Table, application.FieldID, id),
			sqlgraph.To(environment.Table, environment.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, application.EnvironmentTable, application.EnvironmentColumn),
		)
		schemaConfig := a.schemaConfig
		step.To.Schema = schemaConfig.Environment
		step.Edge.Schema = schemaConfig.Application
		fromV = sqlgraph.Neighbors(a.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryResources queries the resources edge of a Application.
func (c *ApplicationClient) QueryResources(a *Application) *ApplicationResourceQuery {
	query := (&ApplicationResourceClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := a.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(application.Table, application.FieldID, id),
			sqlgraph.To(applicationresource.Table, applicationresource.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, application.ResourcesTable, application.ResourcesColumn),
		)
		schemaConfig := a.schemaConfig
		step.To.Schema = schemaConfig.ApplicationResource
		step.Edge.Schema = schemaConfig.ApplicationResource
		fromV = sqlgraph.Neighbors(a.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryRevisions queries the revisions edge of a Application.
func (c *ApplicationClient) QueryRevisions(a *Application) *ApplicationRevisionQuery {
	query := (&ApplicationRevisionClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := a.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(application.Table, application.FieldID, id),
			sqlgraph.To(applicationrevision.Table, applicationrevision.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, application.RevisionsTable, application.RevisionsColumn),
		)
		schemaConfig := a.schemaConfig
		step.To.Schema = schemaConfig.ApplicationRevision
		step.Edge.Schema = schemaConfig.ApplicationRevision
		fromV = sqlgraph.Neighbors(a.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryModules queries the modules edge of a Application.
func (c *ApplicationClient) QueryModules(a *Application) *ModuleQuery {
	query := (&ModuleClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := a.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(application.Table, application.FieldID, id),
			sqlgraph.To(module.Table, module.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, application.ModulesTable, application.ModulesPrimaryKey...),
		)
		schemaConfig := a.schemaConfig
		step.To.Schema = schemaConfig.Module
		step.Edge.Schema = schemaConfig.ApplicationModuleRelationship
		fromV = sqlgraph.Neighbors(a.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryApplicationModuleRelationships queries the applicationModuleRelationships edge of a Application.
func (c *ApplicationClient) QueryApplicationModuleRelationships(a *Application) *ApplicationModuleRelationshipQuery {
	query := (&ApplicationModuleRelationshipClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := a.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(application.Table, application.FieldID, id),
			sqlgraph.To(applicationmodulerelationship.Table, applicationmodulerelationship.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, application.ApplicationModuleRelationshipsTable, application.ApplicationModuleRelationshipsColumn),
		)
		schemaConfig := a.schemaConfig
		step.To.Schema = schemaConfig.ApplicationModuleRelationship
		step.Edge.Schema = schemaConfig.ApplicationModuleRelationship
		fromV = sqlgraph.Neighbors(a.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *ApplicationClient) Hooks() []Hook {
	hooks := c.hooks.Application
	return append(hooks[:len(hooks):len(hooks)], application.Hooks[:]...)
}

// Interceptors returns the client interceptors.
func (c *ApplicationClient) Interceptors() []Interceptor {
	return c.inters.Application
}

func (c *ApplicationClient) mutate(ctx context.Context, m *ApplicationMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&ApplicationCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&ApplicationUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&ApplicationUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&ApplicationDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("model: unknown Application mutation op: %q", m.Op())
	}
}

// ApplicationModuleRelationshipClient is a client for the ApplicationModuleRelationship schema.
type ApplicationModuleRelationshipClient struct {
	config
}

// NewApplicationModuleRelationshipClient returns a client for the ApplicationModuleRelationship from the given config.
func NewApplicationModuleRelationshipClient(c config) *ApplicationModuleRelationshipClient {
	return &ApplicationModuleRelationshipClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `applicationmodulerelationship.Hooks(f(g(h())))`.
func (c *ApplicationModuleRelationshipClient) Use(hooks ...Hook) {
	c.hooks.ApplicationModuleRelationship = append(c.hooks.ApplicationModuleRelationship, hooks...)
}

// Use adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `applicationmodulerelationship.Intercept(f(g(h())))`.
func (c *ApplicationModuleRelationshipClient) Intercept(interceptors ...Interceptor) {
	c.inters.ApplicationModuleRelationship = append(c.inters.ApplicationModuleRelationship, interceptors...)
}

// Create returns a builder for creating a ApplicationModuleRelationship entity.
func (c *ApplicationModuleRelationshipClient) Create() *ApplicationModuleRelationshipCreate {
	mutation := newApplicationModuleRelationshipMutation(c.config, OpCreate)
	return &ApplicationModuleRelationshipCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of ApplicationModuleRelationship entities.
func (c *ApplicationModuleRelationshipClient) CreateBulk(builders ...*ApplicationModuleRelationshipCreate) *ApplicationModuleRelationshipCreateBulk {
	return &ApplicationModuleRelationshipCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for ApplicationModuleRelationship.
func (c *ApplicationModuleRelationshipClient) Update() *ApplicationModuleRelationshipUpdate {
	mutation := newApplicationModuleRelationshipMutation(c.config, OpUpdate)
	return &ApplicationModuleRelationshipUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *ApplicationModuleRelationshipClient) UpdateOne(amr *ApplicationModuleRelationship) *ApplicationModuleRelationshipUpdateOne {
	mutation := newApplicationModuleRelationshipMutation(c.config, OpUpdateOne, withApplicationModuleRelationship(amr))
	return &ApplicationModuleRelationshipUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *ApplicationModuleRelationshipClient) UpdateOneID(id int) *ApplicationModuleRelationshipUpdateOne {
	mutation := newApplicationModuleRelationshipMutation(c.config, OpUpdateOne, withApplicationModuleRelationshipID(id))
	return &ApplicationModuleRelationshipUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for ApplicationModuleRelationship.
func (c *ApplicationModuleRelationshipClient) Delete() *ApplicationModuleRelationshipDelete {
	mutation := newApplicationModuleRelationshipMutation(c.config, OpDelete)
	return &ApplicationModuleRelationshipDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *ApplicationModuleRelationshipClient) DeleteOne(amr *ApplicationModuleRelationship) *ApplicationModuleRelationshipDeleteOne {
	return c.DeleteOneID(amr.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *ApplicationModuleRelationshipClient) DeleteOneID(id int) *ApplicationModuleRelationshipDeleteOne {
	builder := c.Delete().Where(applicationmodulerelationship.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &ApplicationModuleRelationshipDeleteOne{builder}
}

// Query returns a query builder for ApplicationModuleRelationship.
func (c *ApplicationModuleRelationshipClient) Query() *ApplicationModuleRelationshipQuery {
	return &ApplicationModuleRelationshipQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeApplicationModuleRelationship},
		inters: c.Interceptors(),
	}
}

// Get returns a ApplicationModuleRelationship entity by its id.
func (c *ApplicationModuleRelationshipClient) Get(ctx context.Context, id int) (*ApplicationModuleRelationship, error) {
	return c.Query().Where(applicationmodulerelationship.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *ApplicationModuleRelationshipClient) GetX(ctx context.Context, id int) *ApplicationModuleRelationship {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryApplication queries the application edge of a ApplicationModuleRelationship.
func (c *ApplicationModuleRelationshipClient) QueryApplication(amr *ApplicationModuleRelationship) *ApplicationQuery {
	query := (&ApplicationClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := amr.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(applicationmodulerelationship.Table, applicationmodulerelationship.FieldID, id),
			sqlgraph.To(application.Table, application.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, applicationmodulerelationship.ApplicationTable, applicationmodulerelationship.ApplicationColumn),
		)
		schemaConfig := amr.schemaConfig
		step.To.Schema = schemaConfig.Application
		step.Edge.Schema = schemaConfig.ApplicationModuleRelationship
		fromV = sqlgraph.Neighbors(amr.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryModule queries the module edge of a ApplicationModuleRelationship.
func (c *ApplicationModuleRelationshipClient) QueryModule(amr *ApplicationModuleRelationship) *ModuleQuery {
	query := (&ModuleClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := amr.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(applicationmodulerelationship.Table, applicationmodulerelationship.FieldID, id),
			sqlgraph.To(module.Table, module.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, applicationmodulerelationship.ModuleTable, applicationmodulerelationship.ModuleColumn),
		)
		schemaConfig := amr.schemaConfig
		step.To.Schema = schemaConfig.Module
		step.Edge.Schema = schemaConfig.ApplicationModuleRelationship
		fromV = sqlgraph.Neighbors(amr.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *ApplicationModuleRelationshipClient) Hooks() []Hook {
	return c.hooks.ApplicationModuleRelationship
}

// Interceptors returns the client interceptors.
func (c *ApplicationModuleRelationshipClient) Interceptors() []Interceptor {
	return c.inters.ApplicationModuleRelationship
}

func (c *ApplicationModuleRelationshipClient) mutate(ctx context.Context, m *ApplicationModuleRelationshipMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&ApplicationModuleRelationshipCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&ApplicationModuleRelationshipUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&ApplicationModuleRelationshipUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&ApplicationModuleRelationshipDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("model: unknown ApplicationModuleRelationship mutation op: %q", m.Op())
	}
}

// ApplicationResourceClient is a client for the ApplicationResource schema.
type ApplicationResourceClient struct {
	config
}

// NewApplicationResourceClient returns a client for the ApplicationResource from the given config.
func NewApplicationResourceClient(c config) *ApplicationResourceClient {
	return &ApplicationResourceClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `applicationresource.Hooks(f(g(h())))`.
func (c *ApplicationResourceClient) Use(hooks ...Hook) {
	c.hooks.ApplicationResource = append(c.hooks.ApplicationResource, hooks...)
}

// Use adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `applicationresource.Intercept(f(g(h())))`.
func (c *ApplicationResourceClient) Intercept(interceptors ...Interceptor) {
	c.inters.ApplicationResource = append(c.inters.ApplicationResource, interceptors...)
}

// Create returns a builder for creating a ApplicationResource entity.
func (c *ApplicationResourceClient) Create() *ApplicationResourceCreate {
	mutation := newApplicationResourceMutation(c.config, OpCreate)
	return &ApplicationResourceCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of ApplicationResource entities.
func (c *ApplicationResourceClient) CreateBulk(builders ...*ApplicationResourceCreate) *ApplicationResourceCreateBulk {
	return &ApplicationResourceCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for ApplicationResource.
func (c *ApplicationResourceClient) Update() *ApplicationResourceUpdate {
	mutation := newApplicationResourceMutation(c.config, OpUpdate)
	return &ApplicationResourceUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *ApplicationResourceClient) UpdateOne(ar *ApplicationResource) *ApplicationResourceUpdateOne {
	mutation := newApplicationResourceMutation(c.config, OpUpdateOne, withApplicationResource(ar))
	return &ApplicationResourceUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *ApplicationResourceClient) UpdateOneID(id types.ID) *ApplicationResourceUpdateOne {
	mutation := newApplicationResourceMutation(c.config, OpUpdateOne, withApplicationResourceID(id))
	return &ApplicationResourceUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for ApplicationResource.
func (c *ApplicationResourceClient) Delete() *ApplicationResourceDelete {
	mutation := newApplicationResourceMutation(c.config, OpDelete)
	return &ApplicationResourceDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *ApplicationResourceClient) DeleteOne(ar *ApplicationResource) *ApplicationResourceDeleteOne {
	return c.DeleteOneID(ar.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *ApplicationResourceClient) DeleteOneID(id types.ID) *ApplicationResourceDeleteOne {
	builder := c.Delete().Where(applicationresource.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &ApplicationResourceDeleteOne{builder}
}

// Query returns a query builder for ApplicationResource.
func (c *ApplicationResourceClient) Query() *ApplicationResourceQuery {
	return &ApplicationResourceQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeApplicationResource},
		inters: c.Interceptors(),
	}
}

// Get returns a ApplicationResource entity by its id.
func (c *ApplicationResourceClient) Get(ctx context.Context, id types.ID) (*ApplicationResource, error) {
	return c.Query().Where(applicationresource.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *ApplicationResourceClient) GetX(ctx context.Context, id types.ID) *ApplicationResource {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryApplication queries the application edge of a ApplicationResource.
func (c *ApplicationResourceClient) QueryApplication(ar *ApplicationResource) *ApplicationQuery {
	query := (&ApplicationClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := ar.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(applicationresource.Table, applicationresource.FieldID, id),
			sqlgraph.To(application.Table, application.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, applicationresource.ApplicationTable, applicationresource.ApplicationColumn),
		)
		schemaConfig := ar.schemaConfig
		step.To.Schema = schemaConfig.Application
		step.Edge.Schema = schemaConfig.ApplicationResource
		fromV = sqlgraph.Neighbors(ar.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryConnector queries the connector edge of a ApplicationResource.
func (c *ApplicationResourceClient) QueryConnector(ar *ApplicationResource) *ConnectorQuery {
	query := (&ConnectorClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := ar.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(applicationresource.Table, applicationresource.FieldID, id),
			sqlgraph.To(connector.Table, connector.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, applicationresource.ConnectorTable, applicationresource.ConnectorColumn),
		)
		schemaConfig := ar.schemaConfig
		step.To.Schema = schemaConfig.Connector
		step.Edge.Schema = schemaConfig.ApplicationResource
		fromV = sqlgraph.Neighbors(ar.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *ApplicationResourceClient) Hooks() []Hook {
	hooks := c.hooks.ApplicationResource
	return append(hooks[:len(hooks):len(hooks)], applicationresource.Hooks[:]...)
}

// Interceptors returns the client interceptors.
func (c *ApplicationResourceClient) Interceptors() []Interceptor {
	return c.inters.ApplicationResource
}

func (c *ApplicationResourceClient) mutate(ctx context.Context, m *ApplicationResourceMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&ApplicationResourceCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&ApplicationResourceUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&ApplicationResourceUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&ApplicationResourceDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("model: unknown ApplicationResource mutation op: %q", m.Op())
	}
}

// ApplicationRevisionClient is a client for the ApplicationRevision schema.
type ApplicationRevisionClient struct {
	config
}

// NewApplicationRevisionClient returns a client for the ApplicationRevision from the given config.
func NewApplicationRevisionClient(c config) *ApplicationRevisionClient {
	return &ApplicationRevisionClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `applicationrevision.Hooks(f(g(h())))`.
func (c *ApplicationRevisionClient) Use(hooks ...Hook) {
	c.hooks.ApplicationRevision = append(c.hooks.ApplicationRevision, hooks...)
}

// Use adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `applicationrevision.Intercept(f(g(h())))`.
func (c *ApplicationRevisionClient) Intercept(interceptors ...Interceptor) {
	c.inters.ApplicationRevision = append(c.inters.ApplicationRevision, interceptors...)
}

// Create returns a builder for creating a ApplicationRevision entity.
func (c *ApplicationRevisionClient) Create() *ApplicationRevisionCreate {
	mutation := newApplicationRevisionMutation(c.config, OpCreate)
	return &ApplicationRevisionCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of ApplicationRevision entities.
func (c *ApplicationRevisionClient) CreateBulk(builders ...*ApplicationRevisionCreate) *ApplicationRevisionCreateBulk {
	return &ApplicationRevisionCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for ApplicationRevision.
func (c *ApplicationRevisionClient) Update() *ApplicationRevisionUpdate {
	mutation := newApplicationRevisionMutation(c.config, OpUpdate)
	return &ApplicationRevisionUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *ApplicationRevisionClient) UpdateOne(ar *ApplicationRevision) *ApplicationRevisionUpdateOne {
	mutation := newApplicationRevisionMutation(c.config, OpUpdateOne, withApplicationRevision(ar))
	return &ApplicationRevisionUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *ApplicationRevisionClient) UpdateOneID(id types.ID) *ApplicationRevisionUpdateOne {
	mutation := newApplicationRevisionMutation(c.config, OpUpdateOne, withApplicationRevisionID(id))
	return &ApplicationRevisionUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for ApplicationRevision.
func (c *ApplicationRevisionClient) Delete() *ApplicationRevisionDelete {
	mutation := newApplicationRevisionMutation(c.config, OpDelete)
	return &ApplicationRevisionDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *ApplicationRevisionClient) DeleteOne(ar *ApplicationRevision) *ApplicationRevisionDeleteOne {
	return c.DeleteOneID(ar.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *ApplicationRevisionClient) DeleteOneID(id types.ID) *ApplicationRevisionDeleteOne {
	builder := c.Delete().Where(applicationrevision.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &ApplicationRevisionDeleteOne{builder}
}

// Query returns a query builder for ApplicationRevision.
func (c *ApplicationRevisionClient) Query() *ApplicationRevisionQuery {
	return &ApplicationRevisionQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeApplicationRevision},
		inters: c.Interceptors(),
	}
}

// Get returns a ApplicationRevision entity by its id.
func (c *ApplicationRevisionClient) Get(ctx context.Context, id types.ID) (*ApplicationRevision, error) {
	return c.Query().Where(applicationrevision.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *ApplicationRevisionClient) GetX(ctx context.Context, id types.ID) *ApplicationRevision {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryApplication queries the application edge of a ApplicationRevision.
func (c *ApplicationRevisionClient) QueryApplication(ar *ApplicationRevision) *ApplicationQuery {
	query := (&ApplicationClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := ar.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(applicationrevision.Table, applicationrevision.FieldID, id),
			sqlgraph.To(application.Table, application.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, applicationrevision.ApplicationTable, applicationrevision.ApplicationColumn),
		)
		schemaConfig := ar.schemaConfig
		step.To.Schema = schemaConfig.Application
		step.Edge.Schema = schemaConfig.ApplicationRevision
		fromV = sqlgraph.Neighbors(ar.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryEnvironment queries the environment edge of a ApplicationRevision.
func (c *ApplicationRevisionClient) QueryEnvironment(ar *ApplicationRevision) *EnvironmentQuery {
	query := (&EnvironmentClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := ar.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(applicationrevision.Table, applicationrevision.FieldID, id),
			sqlgraph.To(environment.Table, environment.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, applicationrevision.EnvironmentTable, applicationrevision.EnvironmentColumn),
		)
		schemaConfig := ar.schemaConfig
		step.To.Schema = schemaConfig.Environment
		step.Edge.Schema = schemaConfig.ApplicationRevision
		fromV = sqlgraph.Neighbors(ar.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *ApplicationRevisionClient) Hooks() []Hook {
	hooks := c.hooks.ApplicationRevision
	return append(hooks[:len(hooks):len(hooks)], applicationrevision.Hooks[:]...)
}

// Interceptors returns the client interceptors.
func (c *ApplicationRevisionClient) Interceptors() []Interceptor {
	return c.inters.ApplicationRevision
}

func (c *ApplicationRevisionClient) mutate(ctx context.Context, m *ApplicationRevisionMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&ApplicationRevisionCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&ApplicationRevisionUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&ApplicationRevisionUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&ApplicationRevisionDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("model: unknown ApplicationRevision mutation op: %q", m.Op())
	}
}

// ConnectorClient is a client for the Connector schema.
type ConnectorClient struct {
	config
}

// NewConnectorClient returns a client for the Connector from the given config.
func NewConnectorClient(c config) *ConnectorClient {
	return &ConnectorClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `connector.Hooks(f(g(h())))`.
func (c *ConnectorClient) Use(hooks ...Hook) {
	c.hooks.Connector = append(c.hooks.Connector, hooks...)
}

// Use adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `connector.Intercept(f(g(h())))`.
func (c *ConnectorClient) Intercept(interceptors ...Interceptor) {
	c.inters.Connector = append(c.inters.Connector, interceptors...)
}

// Create returns a builder for creating a Connector entity.
func (c *ConnectorClient) Create() *ConnectorCreate {
	mutation := newConnectorMutation(c.config, OpCreate)
	return &ConnectorCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Connector entities.
func (c *ConnectorClient) CreateBulk(builders ...*ConnectorCreate) *ConnectorCreateBulk {
	return &ConnectorCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Connector.
func (c *ConnectorClient) Update() *ConnectorUpdate {
	mutation := newConnectorMutation(c.config, OpUpdate)
	return &ConnectorUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *ConnectorClient) UpdateOne(co *Connector) *ConnectorUpdateOne {
	mutation := newConnectorMutation(c.config, OpUpdateOne, withConnector(co))
	return &ConnectorUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *ConnectorClient) UpdateOneID(id types.ID) *ConnectorUpdateOne {
	mutation := newConnectorMutation(c.config, OpUpdateOne, withConnectorID(id))
	return &ConnectorUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Connector.
func (c *ConnectorClient) Delete() *ConnectorDelete {
	mutation := newConnectorMutation(c.config, OpDelete)
	return &ConnectorDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *ConnectorClient) DeleteOne(co *Connector) *ConnectorDeleteOne {
	return c.DeleteOneID(co.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *ConnectorClient) DeleteOneID(id types.ID) *ConnectorDeleteOne {
	builder := c.Delete().Where(connector.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &ConnectorDeleteOne{builder}
}

// Query returns a query builder for Connector.
func (c *ConnectorClient) Query() *ConnectorQuery {
	return &ConnectorQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeConnector},
		inters: c.Interceptors(),
	}
}

// Get returns a Connector entity by its id.
func (c *ConnectorClient) Get(ctx context.Context, id types.ID) (*Connector, error) {
	return c.Query().Where(connector.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *ConnectorClient) GetX(ctx context.Context, id types.ID) *Connector {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryEnvironments queries the environments edge of a Connector.
func (c *ConnectorClient) QueryEnvironments(co *Connector) *EnvironmentQuery {
	query := (&EnvironmentClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := co.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(connector.Table, connector.FieldID, id),
			sqlgraph.To(environment.Table, environment.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, connector.EnvironmentsTable, connector.EnvironmentsPrimaryKey...),
		)
		schemaConfig := co.schemaConfig
		step.To.Schema = schemaConfig.Environment
		step.Edge.Schema = schemaConfig.EnvironmentConnectorRelationship
		fromV = sqlgraph.Neighbors(co.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryResources queries the resources edge of a Connector.
func (c *ConnectorClient) QueryResources(co *Connector) *ApplicationResourceQuery {
	query := (&ApplicationResourceClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := co.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(connector.Table, connector.FieldID, id),
			sqlgraph.To(applicationresource.Table, applicationresource.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, connector.ResourcesTable, connector.ResourcesColumn),
		)
		schemaConfig := co.schemaConfig
		step.To.Schema = schemaConfig.ApplicationResource
		step.Edge.Schema = schemaConfig.ApplicationResource
		fromV = sqlgraph.Neighbors(co.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryEnvironmentConnectorRelationships queries the environmentConnectorRelationships edge of a Connector.
func (c *ConnectorClient) QueryEnvironmentConnectorRelationships(co *Connector) *EnvironmentConnectorRelationshipQuery {
	query := (&EnvironmentConnectorRelationshipClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := co.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(connector.Table, connector.FieldID, id),
			sqlgraph.To(environmentconnectorrelationship.Table, environmentconnectorrelationship.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, connector.EnvironmentConnectorRelationshipsTable, connector.EnvironmentConnectorRelationshipsColumn),
		)
		schemaConfig := co.schemaConfig
		step.To.Schema = schemaConfig.EnvironmentConnectorRelationship
		step.Edge.Schema = schemaConfig.EnvironmentConnectorRelationship
		fromV = sqlgraph.Neighbors(co.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *ConnectorClient) Hooks() []Hook {
	hooks := c.hooks.Connector
	return append(hooks[:len(hooks):len(hooks)], connector.Hooks[:]...)
}

// Interceptors returns the client interceptors.
func (c *ConnectorClient) Interceptors() []Interceptor {
	return c.inters.Connector
}

func (c *ConnectorClient) mutate(ctx context.Context, m *ConnectorMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&ConnectorCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&ConnectorUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&ConnectorUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&ConnectorDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("model: unknown Connector mutation op: %q", m.Op())
	}
}

// EnvironmentClient is a client for the Environment schema.
type EnvironmentClient struct {
	config
}

// NewEnvironmentClient returns a client for the Environment from the given config.
func NewEnvironmentClient(c config) *EnvironmentClient {
	return &EnvironmentClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `environment.Hooks(f(g(h())))`.
func (c *EnvironmentClient) Use(hooks ...Hook) {
	c.hooks.Environment = append(c.hooks.Environment, hooks...)
}

// Use adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `environment.Intercept(f(g(h())))`.
func (c *EnvironmentClient) Intercept(interceptors ...Interceptor) {
	c.inters.Environment = append(c.inters.Environment, interceptors...)
}

// Create returns a builder for creating a Environment entity.
func (c *EnvironmentClient) Create() *EnvironmentCreate {
	mutation := newEnvironmentMutation(c.config, OpCreate)
	return &EnvironmentCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Environment entities.
func (c *EnvironmentClient) CreateBulk(builders ...*EnvironmentCreate) *EnvironmentCreateBulk {
	return &EnvironmentCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Environment.
func (c *EnvironmentClient) Update() *EnvironmentUpdate {
	mutation := newEnvironmentMutation(c.config, OpUpdate)
	return &EnvironmentUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *EnvironmentClient) UpdateOne(e *Environment) *EnvironmentUpdateOne {
	mutation := newEnvironmentMutation(c.config, OpUpdateOne, withEnvironment(e))
	return &EnvironmentUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *EnvironmentClient) UpdateOneID(id types.ID) *EnvironmentUpdateOne {
	mutation := newEnvironmentMutation(c.config, OpUpdateOne, withEnvironmentID(id))
	return &EnvironmentUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Environment.
func (c *EnvironmentClient) Delete() *EnvironmentDelete {
	mutation := newEnvironmentMutation(c.config, OpDelete)
	return &EnvironmentDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *EnvironmentClient) DeleteOne(e *Environment) *EnvironmentDeleteOne {
	return c.DeleteOneID(e.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *EnvironmentClient) DeleteOneID(id types.ID) *EnvironmentDeleteOne {
	builder := c.Delete().Where(environment.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &EnvironmentDeleteOne{builder}
}

// Query returns a query builder for Environment.
func (c *EnvironmentClient) Query() *EnvironmentQuery {
	return &EnvironmentQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeEnvironment},
		inters: c.Interceptors(),
	}
}

// Get returns a Environment entity by its id.
func (c *EnvironmentClient) Get(ctx context.Context, id types.ID) (*Environment, error) {
	return c.Query().Where(environment.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *EnvironmentClient) GetX(ctx context.Context, id types.ID) *Environment {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryConnectors queries the connectors edge of a Environment.
func (c *EnvironmentClient) QueryConnectors(e *Environment) *ConnectorQuery {
	query := (&ConnectorClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := e.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(environment.Table, environment.FieldID, id),
			sqlgraph.To(connector.Table, connector.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, environment.ConnectorsTable, environment.ConnectorsPrimaryKey...),
		)
		schemaConfig := e.schemaConfig
		step.To.Schema = schemaConfig.Connector
		step.Edge.Schema = schemaConfig.EnvironmentConnectorRelationship
		fromV = sqlgraph.Neighbors(e.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryApplications queries the applications edge of a Environment.
func (c *EnvironmentClient) QueryApplications(e *Environment) *ApplicationQuery {
	query := (&ApplicationClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := e.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(environment.Table, environment.FieldID, id),
			sqlgraph.To(application.Table, application.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, environment.ApplicationsTable, environment.ApplicationsColumn),
		)
		schemaConfig := e.schemaConfig
		step.To.Schema = schemaConfig.Application
		step.Edge.Schema = schemaConfig.Application
		fromV = sqlgraph.Neighbors(e.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryRevisions queries the revisions edge of a Environment.
func (c *EnvironmentClient) QueryRevisions(e *Environment) *ApplicationRevisionQuery {
	query := (&ApplicationRevisionClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := e.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(environment.Table, environment.FieldID, id),
			sqlgraph.To(applicationrevision.Table, applicationrevision.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, environment.RevisionsTable, environment.RevisionsColumn),
		)
		schemaConfig := e.schemaConfig
		step.To.Schema = schemaConfig.ApplicationRevision
		step.Edge.Schema = schemaConfig.ApplicationRevision
		fromV = sqlgraph.Neighbors(e.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryEnvironmentConnectorRelationships queries the environmentConnectorRelationships edge of a Environment.
func (c *EnvironmentClient) QueryEnvironmentConnectorRelationships(e *Environment) *EnvironmentConnectorRelationshipQuery {
	query := (&EnvironmentConnectorRelationshipClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := e.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(environment.Table, environment.FieldID, id),
			sqlgraph.To(environmentconnectorrelationship.Table, environmentconnectorrelationship.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, environment.EnvironmentConnectorRelationshipsTable, environment.EnvironmentConnectorRelationshipsColumn),
		)
		schemaConfig := e.schemaConfig
		step.To.Schema = schemaConfig.EnvironmentConnectorRelationship
		step.Edge.Schema = schemaConfig.EnvironmentConnectorRelationship
		fromV = sqlgraph.Neighbors(e.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *EnvironmentClient) Hooks() []Hook {
	hooks := c.hooks.Environment
	return append(hooks[:len(hooks):len(hooks)], environment.Hooks[:]...)
}

// Interceptors returns the client interceptors.
func (c *EnvironmentClient) Interceptors() []Interceptor {
	return c.inters.Environment
}

func (c *EnvironmentClient) mutate(ctx context.Context, m *EnvironmentMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&EnvironmentCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&EnvironmentUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&EnvironmentUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&EnvironmentDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("model: unknown Environment mutation op: %q", m.Op())
	}
}

// EnvironmentConnectorRelationshipClient is a client for the EnvironmentConnectorRelationship schema.
type EnvironmentConnectorRelationshipClient struct {
	config
}

// NewEnvironmentConnectorRelationshipClient returns a client for the EnvironmentConnectorRelationship from the given config.
func NewEnvironmentConnectorRelationshipClient(c config) *EnvironmentConnectorRelationshipClient {
	return &EnvironmentConnectorRelationshipClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `environmentconnectorrelationship.Hooks(f(g(h())))`.
func (c *EnvironmentConnectorRelationshipClient) Use(hooks ...Hook) {
	c.hooks.EnvironmentConnectorRelationship = append(c.hooks.EnvironmentConnectorRelationship, hooks...)
}

// Use adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `environmentconnectorrelationship.Intercept(f(g(h())))`.
func (c *EnvironmentConnectorRelationshipClient) Intercept(interceptors ...Interceptor) {
	c.inters.EnvironmentConnectorRelationship = append(c.inters.EnvironmentConnectorRelationship, interceptors...)
}

// Create returns a builder for creating a EnvironmentConnectorRelationship entity.
func (c *EnvironmentConnectorRelationshipClient) Create() *EnvironmentConnectorRelationshipCreate {
	mutation := newEnvironmentConnectorRelationshipMutation(c.config, OpCreate)
	return &EnvironmentConnectorRelationshipCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of EnvironmentConnectorRelationship entities.
func (c *EnvironmentConnectorRelationshipClient) CreateBulk(builders ...*EnvironmentConnectorRelationshipCreate) *EnvironmentConnectorRelationshipCreateBulk {
	return &EnvironmentConnectorRelationshipCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for EnvironmentConnectorRelationship.
func (c *EnvironmentConnectorRelationshipClient) Update() *EnvironmentConnectorRelationshipUpdate {
	mutation := newEnvironmentConnectorRelationshipMutation(c.config, OpUpdate)
	return &EnvironmentConnectorRelationshipUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *EnvironmentConnectorRelationshipClient) UpdateOne(ecr *EnvironmentConnectorRelationship) *EnvironmentConnectorRelationshipUpdateOne {
	mutation := newEnvironmentConnectorRelationshipMutation(c.config, OpUpdateOne, withEnvironmentConnectorRelationship(ecr))
	return &EnvironmentConnectorRelationshipUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *EnvironmentConnectorRelationshipClient) UpdateOneID(id int) *EnvironmentConnectorRelationshipUpdateOne {
	mutation := newEnvironmentConnectorRelationshipMutation(c.config, OpUpdateOne, withEnvironmentConnectorRelationshipID(id))
	return &EnvironmentConnectorRelationshipUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for EnvironmentConnectorRelationship.
func (c *EnvironmentConnectorRelationshipClient) Delete() *EnvironmentConnectorRelationshipDelete {
	mutation := newEnvironmentConnectorRelationshipMutation(c.config, OpDelete)
	return &EnvironmentConnectorRelationshipDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *EnvironmentConnectorRelationshipClient) DeleteOne(ecr *EnvironmentConnectorRelationship) *EnvironmentConnectorRelationshipDeleteOne {
	return c.DeleteOneID(ecr.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *EnvironmentConnectorRelationshipClient) DeleteOneID(id int) *EnvironmentConnectorRelationshipDeleteOne {
	builder := c.Delete().Where(environmentconnectorrelationship.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &EnvironmentConnectorRelationshipDeleteOne{builder}
}

// Query returns a query builder for EnvironmentConnectorRelationship.
func (c *EnvironmentConnectorRelationshipClient) Query() *EnvironmentConnectorRelationshipQuery {
	return &EnvironmentConnectorRelationshipQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeEnvironmentConnectorRelationship},
		inters: c.Interceptors(),
	}
}

// Get returns a EnvironmentConnectorRelationship entity by its id.
func (c *EnvironmentConnectorRelationshipClient) Get(ctx context.Context, id int) (*EnvironmentConnectorRelationship, error) {
	return c.Query().Where(environmentconnectorrelationship.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *EnvironmentConnectorRelationshipClient) GetX(ctx context.Context, id int) *EnvironmentConnectorRelationship {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryEnvironment queries the environment edge of a EnvironmentConnectorRelationship.
func (c *EnvironmentConnectorRelationshipClient) QueryEnvironment(ecr *EnvironmentConnectorRelationship) *EnvironmentQuery {
	query := (&EnvironmentClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := ecr.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(environmentconnectorrelationship.Table, environmentconnectorrelationship.FieldID, id),
			sqlgraph.To(environment.Table, environment.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, environmentconnectorrelationship.EnvironmentTable, environmentconnectorrelationship.EnvironmentColumn),
		)
		schemaConfig := ecr.schemaConfig
		step.To.Schema = schemaConfig.Environment
		step.Edge.Schema = schemaConfig.EnvironmentConnectorRelationship
		fromV = sqlgraph.Neighbors(ecr.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryConnector queries the connector edge of a EnvironmentConnectorRelationship.
func (c *EnvironmentConnectorRelationshipClient) QueryConnector(ecr *EnvironmentConnectorRelationship) *ConnectorQuery {
	query := (&ConnectorClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := ecr.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(environmentconnectorrelationship.Table, environmentconnectorrelationship.FieldID, id),
			sqlgraph.To(connector.Table, connector.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, environmentconnectorrelationship.ConnectorTable, environmentconnectorrelationship.ConnectorColumn),
		)
		schemaConfig := ecr.schemaConfig
		step.To.Schema = schemaConfig.Connector
		step.Edge.Schema = schemaConfig.EnvironmentConnectorRelationship
		fromV = sqlgraph.Neighbors(ecr.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *EnvironmentConnectorRelationshipClient) Hooks() []Hook {
	return c.hooks.EnvironmentConnectorRelationship
}

// Interceptors returns the client interceptors.
func (c *EnvironmentConnectorRelationshipClient) Interceptors() []Interceptor {
	return c.inters.EnvironmentConnectorRelationship
}

func (c *EnvironmentConnectorRelationshipClient) mutate(ctx context.Context, m *EnvironmentConnectorRelationshipMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&EnvironmentConnectorRelationshipCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&EnvironmentConnectorRelationshipUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&EnvironmentConnectorRelationshipUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&EnvironmentConnectorRelationshipDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("model: unknown EnvironmentConnectorRelationship mutation op: %q", m.Op())
	}
}

// ModuleClient is a client for the Module schema.
type ModuleClient struct {
	config
}

// NewModuleClient returns a client for the Module from the given config.
func NewModuleClient(c config) *ModuleClient {
	return &ModuleClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `module.Hooks(f(g(h())))`.
func (c *ModuleClient) Use(hooks ...Hook) {
	c.hooks.Module = append(c.hooks.Module, hooks...)
}

// Use adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `module.Intercept(f(g(h())))`.
func (c *ModuleClient) Intercept(interceptors ...Interceptor) {
	c.inters.Module = append(c.inters.Module, interceptors...)
}

// Create returns a builder for creating a Module entity.
func (c *ModuleClient) Create() *ModuleCreate {
	mutation := newModuleMutation(c.config, OpCreate)
	return &ModuleCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Module entities.
func (c *ModuleClient) CreateBulk(builders ...*ModuleCreate) *ModuleCreateBulk {
	return &ModuleCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Module.
func (c *ModuleClient) Update() *ModuleUpdate {
	mutation := newModuleMutation(c.config, OpUpdate)
	return &ModuleUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *ModuleClient) UpdateOne(m *Module) *ModuleUpdateOne {
	mutation := newModuleMutation(c.config, OpUpdateOne, withModule(m))
	return &ModuleUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *ModuleClient) UpdateOneID(id string) *ModuleUpdateOne {
	mutation := newModuleMutation(c.config, OpUpdateOne, withModuleID(id))
	return &ModuleUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Module.
func (c *ModuleClient) Delete() *ModuleDelete {
	mutation := newModuleMutation(c.config, OpDelete)
	return &ModuleDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *ModuleClient) DeleteOne(m *Module) *ModuleDeleteOne {
	return c.DeleteOneID(m.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *ModuleClient) DeleteOneID(id string) *ModuleDeleteOne {
	builder := c.Delete().Where(module.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &ModuleDeleteOne{builder}
}

// Query returns a query builder for Module.
func (c *ModuleClient) Query() *ModuleQuery {
	return &ModuleQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeModule},
		inters: c.Interceptors(),
	}
}

// Get returns a Module entity by its id.
func (c *ModuleClient) Get(ctx context.Context, id string) (*Module, error) {
	return c.Query().Where(module.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *ModuleClient) GetX(ctx context.Context, id string) *Module {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryApplication queries the application edge of a Module.
func (c *ModuleClient) QueryApplication(m *Module) *ApplicationQuery {
	query := (&ApplicationClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := m.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(module.Table, module.FieldID, id),
			sqlgraph.To(application.Table, application.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, module.ApplicationTable, module.ApplicationPrimaryKey...),
		)
		schemaConfig := m.schemaConfig
		step.To.Schema = schemaConfig.Application
		step.Edge.Schema = schemaConfig.ApplicationModuleRelationship
		fromV = sqlgraph.Neighbors(m.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryApplicationModuleRelationships queries the applicationModuleRelationships edge of a Module.
func (c *ModuleClient) QueryApplicationModuleRelationships(m *Module) *ApplicationModuleRelationshipQuery {
	query := (&ApplicationModuleRelationshipClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := m.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(module.Table, module.FieldID, id),
			sqlgraph.To(applicationmodulerelationship.Table, applicationmodulerelationship.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, module.ApplicationModuleRelationshipsTable, module.ApplicationModuleRelationshipsColumn),
		)
		schemaConfig := m.schemaConfig
		step.To.Schema = schemaConfig.ApplicationModuleRelationship
		step.Edge.Schema = schemaConfig.ApplicationModuleRelationship
		fromV = sqlgraph.Neighbors(m.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *ModuleClient) Hooks() []Hook {
	return c.hooks.Module
}

// Interceptors returns the client interceptors.
func (c *ModuleClient) Interceptors() []Interceptor {
	return c.inters.Module
}

func (c *ModuleClient) mutate(ctx context.Context, m *ModuleMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&ModuleCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&ModuleUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&ModuleUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&ModuleDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("model: unknown Module mutation op: %q", m.Op())
	}
}

// ProjectClient is a client for the Project schema.
type ProjectClient struct {
	config
}

// NewProjectClient returns a client for the Project from the given config.
func NewProjectClient(c config) *ProjectClient {
	return &ProjectClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `project.Hooks(f(g(h())))`.
func (c *ProjectClient) Use(hooks ...Hook) {
	c.hooks.Project = append(c.hooks.Project, hooks...)
}

// Use adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `project.Intercept(f(g(h())))`.
func (c *ProjectClient) Intercept(interceptors ...Interceptor) {
	c.inters.Project = append(c.inters.Project, interceptors...)
}

// Create returns a builder for creating a Project entity.
func (c *ProjectClient) Create() *ProjectCreate {
	mutation := newProjectMutation(c.config, OpCreate)
	return &ProjectCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Project entities.
func (c *ProjectClient) CreateBulk(builders ...*ProjectCreate) *ProjectCreateBulk {
	return &ProjectCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Project.
func (c *ProjectClient) Update() *ProjectUpdate {
	mutation := newProjectMutation(c.config, OpUpdate)
	return &ProjectUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *ProjectClient) UpdateOne(pr *Project) *ProjectUpdateOne {
	mutation := newProjectMutation(c.config, OpUpdateOne, withProject(pr))
	return &ProjectUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *ProjectClient) UpdateOneID(id types.ID) *ProjectUpdateOne {
	mutation := newProjectMutation(c.config, OpUpdateOne, withProjectID(id))
	return &ProjectUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Project.
func (c *ProjectClient) Delete() *ProjectDelete {
	mutation := newProjectMutation(c.config, OpDelete)
	return &ProjectDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *ProjectClient) DeleteOne(pr *Project) *ProjectDeleteOne {
	return c.DeleteOneID(pr.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *ProjectClient) DeleteOneID(id types.ID) *ProjectDeleteOne {
	builder := c.Delete().Where(project.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &ProjectDeleteOne{builder}
}

// Query returns a query builder for Project.
func (c *ProjectClient) Query() *ProjectQuery {
	return &ProjectQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeProject},
		inters: c.Interceptors(),
	}
}

// Get returns a Project entity by its id.
func (c *ProjectClient) Get(ctx context.Context, id types.ID) (*Project, error) {
	return c.Query().Where(project.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *ProjectClient) GetX(ctx context.Context, id types.ID) *Project {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryApplications queries the applications edge of a Project.
func (c *ProjectClient) QueryApplications(pr *Project) *ApplicationQuery {
	query := (&ApplicationClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := pr.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(project.Table, project.FieldID, id),
			sqlgraph.To(application.Table, application.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, project.ApplicationsTable, project.ApplicationsColumn),
		)
		schemaConfig := pr.schemaConfig
		step.To.Schema = schemaConfig.Application
		step.Edge.Schema = schemaConfig.Application
		fromV = sqlgraph.Neighbors(pr.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *ProjectClient) Hooks() []Hook {
	hooks := c.hooks.Project
	return append(hooks[:len(hooks):len(hooks)], project.Hooks[:]...)
}

// Interceptors returns the client interceptors.
func (c *ProjectClient) Interceptors() []Interceptor {
	return c.inters.Project
}

func (c *ProjectClient) mutate(ctx context.Context, m *ProjectMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&ProjectCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&ProjectUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&ProjectUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&ProjectDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("model: unknown Project mutation op: %q", m.Op())
	}
}

// RoleClient is a client for the Role schema.
type RoleClient struct {
	config
}

// NewRoleClient returns a client for the Role from the given config.
func NewRoleClient(c config) *RoleClient {
	return &RoleClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `role.Hooks(f(g(h())))`.
func (c *RoleClient) Use(hooks ...Hook) {
	c.hooks.Role = append(c.hooks.Role, hooks...)
}

// Use adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `role.Intercept(f(g(h())))`.
func (c *RoleClient) Intercept(interceptors ...Interceptor) {
	c.inters.Role = append(c.inters.Role, interceptors...)
}

// Create returns a builder for creating a Role entity.
func (c *RoleClient) Create() *RoleCreate {
	mutation := newRoleMutation(c.config, OpCreate)
	return &RoleCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Role entities.
func (c *RoleClient) CreateBulk(builders ...*RoleCreate) *RoleCreateBulk {
	return &RoleCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Role.
func (c *RoleClient) Update() *RoleUpdate {
	mutation := newRoleMutation(c.config, OpUpdate)
	return &RoleUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *RoleClient) UpdateOne(r *Role) *RoleUpdateOne {
	mutation := newRoleMutation(c.config, OpUpdateOne, withRole(r))
	return &RoleUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *RoleClient) UpdateOneID(id types.ID) *RoleUpdateOne {
	mutation := newRoleMutation(c.config, OpUpdateOne, withRoleID(id))
	return &RoleUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Role.
func (c *RoleClient) Delete() *RoleDelete {
	mutation := newRoleMutation(c.config, OpDelete)
	return &RoleDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *RoleClient) DeleteOne(r *Role) *RoleDeleteOne {
	return c.DeleteOneID(r.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *RoleClient) DeleteOneID(id types.ID) *RoleDeleteOne {
	builder := c.Delete().Where(role.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &RoleDeleteOne{builder}
}

// Query returns a query builder for Role.
func (c *RoleClient) Query() *RoleQuery {
	return &RoleQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeRole},
		inters: c.Interceptors(),
	}
}

// Get returns a Role entity by its id.
func (c *RoleClient) Get(ctx context.Context, id types.ID) (*Role, error) {
	return c.Query().Where(role.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *RoleClient) GetX(ctx context.Context, id types.ID) *Role {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// Hooks returns the client hooks.
func (c *RoleClient) Hooks() []Hook {
	hooks := c.hooks.Role
	return append(hooks[:len(hooks):len(hooks)], role.Hooks[:]...)
}

// Interceptors returns the client interceptors.
func (c *RoleClient) Interceptors() []Interceptor {
	return c.inters.Role
}

func (c *RoleClient) mutate(ctx context.Context, m *RoleMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&RoleCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&RoleUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&RoleUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&RoleDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("model: unknown Role mutation op: %q", m.Op())
	}
}

// SettingClient is a client for the Setting schema.
type SettingClient struct {
	config
}

// NewSettingClient returns a client for the Setting from the given config.
func NewSettingClient(c config) *SettingClient {
	return &SettingClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `setting.Hooks(f(g(h())))`.
func (c *SettingClient) Use(hooks ...Hook) {
	c.hooks.Setting = append(c.hooks.Setting, hooks...)
}

// Use adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `setting.Intercept(f(g(h())))`.
func (c *SettingClient) Intercept(interceptors ...Interceptor) {
	c.inters.Setting = append(c.inters.Setting, interceptors...)
}

// Create returns a builder for creating a Setting entity.
func (c *SettingClient) Create() *SettingCreate {
	mutation := newSettingMutation(c.config, OpCreate)
	return &SettingCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Setting entities.
func (c *SettingClient) CreateBulk(builders ...*SettingCreate) *SettingCreateBulk {
	return &SettingCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Setting.
func (c *SettingClient) Update() *SettingUpdate {
	mutation := newSettingMutation(c.config, OpUpdate)
	return &SettingUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *SettingClient) UpdateOne(s *Setting) *SettingUpdateOne {
	mutation := newSettingMutation(c.config, OpUpdateOne, withSetting(s))
	return &SettingUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *SettingClient) UpdateOneID(id types.ID) *SettingUpdateOne {
	mutation := newSettingMutation(c.config, OpUpdateOne, withSettingID(id))
	return &SettingUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Setting.
func (c *SettingClient) Delete() *SettingDelete {
	mutation := newSettingMutation(c.config, OpDelete)
	return &SettingDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *SettingClient) DeleteOne(s *Setting) *SettingDeleteOne {
	return c.DeleteOneID(s.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *SettingClient) DeleteOneID(id types.ID) *SettingDeleteOne {
	builder := c.Delete().Where(setting.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &SettingDeleteOne{builder}
}

// Query returns a query builder for Setting.
func (c *SettingClient) Query() *SettingQuery {
	return &SettingQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeSetting},
		inters: c.Interceptors(),
	}
}

// Get returns a Setting entity by its id.
func (c *SettingClient) Get(ctx context.Context, id types.ID) (*Setting, error) {
	return c.Query().Where(setting.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *SettingClient) GetX(ctx context.Context, id types.ID) *Setting {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// Hooks returns the client hooks.
func (c *SettingClient) Hooks() []Hook {
	hooks := c.hooks.Setting
	return append(hooks[:len(hooks):len(hooks)], setting.Hooks[:]...)
}

// Interceptors returns the client interceptors.
func (c *SettingClient) Interceptors() []Interceptor {
	return c.inters.Setting
}

func (c *SettingClient) mutate(ctx context.Context, m *SettingMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&SettingCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&SettingUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&SettingUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&SettingDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("model: unknown Setting mutation op: %q", m.Op())
	}
}

// SubjectClient is a client for the Subject schema.
type SubjectClient struct {
	config
}

// NewSubjectClient returns a client for the Subject from the given config.
func NewSubjectClient(c config) *SubjectClient {
	return &SubjectClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `subject.Hooks(f(g(h())))`.
func (c *SubjectClient) Use(hooks ...Hook) {
	c.hooks.Subject = append(c.hooks.Subject, hooks...)
}

// Use adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `subject.Intercept(f(g(h())))`.
func (c *SubjectClient) Intercept(interceptors ...Interceptor) {
	c.inters.Subject = append(c.inters.Subject, interceptors...)
}

// Create returns a builder for creating a Subject entity.
func (c *SubjectClient) Create() *SubjectCreate {
	mutation := newSubjectMutation(c.config, OpCreate)
	return &SubjectCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Subject entities.
func (c *SubjectClient) CreateBulk(builders ...*SubjectCreate) *SubjectCreateBulk {
	return &SubjectCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Subject.
func (c *SubjectClient) Update() *SubjectUpdate {
	mutation := newSubjectMutation(c.config, OpUpdate)
	return &SubjectUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *SubjectClient) UpdateOne(s *Subject) *SubjectUpdateOne {
	mutation := newSubjectMutation(c.config, OpUpdateOne, withSubject(s))
	return &SubjectUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *SubjectClient) UpdateOneID(id types.ID) *SubjectUpdateOne {
	mutation := newSubjectMutation(c.config, OpUpdateOne, withSubjectID(id))
	return &SubjectUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Subject.
func (c *SubjectClient) Delete() *SubjectDelete {
	mutation := newSubjectMutation(c.config, OpDelete)
	return &SubjectDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *SubjectClient) DeleteOne(s *Subject) *SubjectDeleteOne {
	return c.DeleteOneID(s.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *SubjectClient) DeleteOneID(id types.ID) *SubjectDeleteOne {
	builder := c.Delete().Where(subject.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &SubjectDeleteOne{builder}
}

// Query returns a query builder for Subject.
func (c *SubjectClient) Query() *SubjectQuery {
	return &SubjectQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeSubject},
		inters: c.Interceptors(),
	}
}

// Get returns a Subject entity by its id.
func (c *SubjectClient) Get(ctx context.Context, id types.ID) (*Subject, error) {
	return c.Query().Where(subject.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *SubjectClient) GetX(ctx context.Context, id types.ID) *Subject {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// Hooks returns the client hooks.
func (c *SubjectClient) Hooks() []Hook {
	hooks := c.hooks.Subject
	return append(hooks[:len(hooks):len(hooks)], subject.Hooks[:]...)
}

// Interceptors returns the client interceptors.
func (c *SubjectClient) Interceptors() []Interceptor {
	return c.inters.Subject
}

func (c *SubjectClient) mutate(ctx context.Context, m *SubjectMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&SubjectCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&SubjectUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&SubjectUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&SubjectDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("model: unknown Subject mutation op: %q", m.Op())
	}
}

// TokenClient is a client for the Token schema.
type TokenClient struct {
	config
}

// NewTokenClient returns a client for the Token from the given config.
func NewTokenClient(c config) *TokenClient {
	return &TokenClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `token.Hooks(f(g(h())))`.
func (c *TokenClient) Use(hooks ...Hook) {
	c.hooks.Token = append(c.hooks.Token, hooks...)
}

// Use adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `token.Intercept(f(g(h())))`.
func (c *TokenClient) Intercept(interceptors ...Interceptor) {
	c.inters.Token = append(c.inters.Token, interceptors...)
}

// Create returns a builder for creating a Token entity.
func (c *TokenClient) Create() *TokenCreate {
	mutation := newTokenMutation(c.config, OpCreate)
	return &TokenCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Token entities.
func (c *TokenClient) CreateBulk(builders ...*TokenCreate) *TokenCreateBulk {
	return &TokenCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Token.
func (c *TokenClient) Update() *TokenUpdate {
	mutation := newTokenMutation(c.config, OpUpdate)
	return &TokenUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *TokenClient) UpdateOne(t *Token) *TokenUpdateOne {
	mutation := newTokenMutation(c.config, OpUpdateOne, withToken(t))
	return &TokenUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *TokenClient) UpdateOneID(id types.ID) *TokenUpdateOne {
	mutation := newTokenMutation(c.config, OpUpdateOne, withTokenID(id))
	return &TokenUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Token.
func (c *TokenClient) Delete() *TokenDelete {
	mutation := newTokenMutation(c.config, OpDelete)
	return &TokenDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *TokenClient) DeleteOne(t *Token) *TokenDeleteOne {
	return c.DeleteOneID(t.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *TokenClient) DeleteOneID(id types.ID) *TokenDeleteOne {
	builder := c.Delete().Where(token.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &TokenDeleteOne{builder}
}

// Query returns a query builder for Token.
func (c *TokenClient) Query() *TokenQuery {
	return &TokenQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeToken},
		inters: c.Interceptors(),
	}
}

// Get returns a Token entity by its id.
func (c *TokenClient) Get(ctx context.Context, id types.ID) (*Token, error) {
	return c.Query().Where(token.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *TokenClient) GetX(ctx context.Context, id types.ID) *Token {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// Hooks returns the client hooks.
func (c *TokenClient) Hooks() []Hook {
	hooks := c.hooks.Token
	return append(hooks[:len(hooks):len(hooks)], token.Hooks[:]...)
}

// Interceptors returns the client interceptors.
func (c *TokenClient) Interceptors() []Interceptor {
	return c.inters.Token
}

func (c *TokenClient) mutate(ctx context.Context, m *TokenMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&TokenCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&TokenUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&TokenUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&TokenDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("model: unknown Token mutation op: %q", m.Op())
	}
}
