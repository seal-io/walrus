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
	"github.com/seal-io/seal/pkg/dao/oid"

	"github.com/seal-io/seal/pkg/dao/model/role"
	"github.com/seal-io/seal/pkg/dao/model/setting"
	"github.com/seal-io/seal/pkg/dao/model/subject"
	"github.com/seal-io/seal/pkg/dao/model/token"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
)

// Client is the client that holds all ent builders.
type Client struct {
	config
	// Schema is the client for creating, migrating and dropping schema.
	Schema *migrate.Schema
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
		ctx:     ctx,
		config:  cfg,
		Role:    NewRoleClient(cfg),
		Setting: NewSettingClient(cfg),
		Subject: NewSubjectClient(cfg),
		Token:   NewTokenClient(cfg),
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
		ctx:     ctx,
		config:  cfg,
		Role:    NewRoleClient(cfg),
		Setting: NewSettingClient(cfg),
		Subject: NewSubjectClient(cfg),
		Token:   NewTokenClient(cfg),
	}, nil
}

// Debug returns a new debug-client. It's used to get verbose logging on specific operations.
//
//	client.Debug().
//		Role.
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
	c.Role.Use(hooks...)
	c.Setting.Use(hooks...)
	c.Subject.Use(hooks...)
	c.Token.Use(hooks...)
}

// Intercept adds the query interceptors to all the entity clients.
// In order to add interceptors to a specific client, call: `client.Node.Intercept(...)`.
func (c *Client) Intercept(interceptors ...Interceptor) {
	c.Role.Intercept(interceptors...)
	c.Setting.Intercept(interceptors...)
	c.Subject.Intercept(interceptors...)
	c.Token.Intercept(interceptors...)
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
func (c *RoleClient) UpdateOneID(id oid.ID) *RoleUpdateOne {
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
func (c *RoleClient) DeleteOneID(id oid.ID) *RoleDeleteOne {
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
func (c *RoleClient) Get(ctx context.Context, id oid.ID) (*Role, error) {
	return c.Query().Where(role.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *RoleClient) GetX(ctx context.Context, id oid.ID) *Role {
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
func (c *SettingClient) UpdateOneID(id oid.ID) *SettingUpdateOne {
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
func (c *SettingClient) DeleteOneID(id oid.ID) *SettingDeleteOne {
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
func (c *SettingClient) Get(ctx context.Context, id oid.ID) (*Setting, error) {
	return c.Query().Where(setting.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *SettingClient) GetX(ctx context.Context, id oid.ID) *Setting {
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
func (c *SubjectClient) UpdateOneID(id oid.ID) *SubjectUpdateOne {
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
func (c *SubjectClient) DeleteOneID(id oid.ID) *SubjectDeleteOne {
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
func (c *SubjectClient) Get(ctx context.Context, id oid.ID) (*Subject, error) {
	return c.Query().Where(subject.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *SubjectClient) GetX(ctx context.Context, id oid.ID) *Subject {
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
func (c *TokenClient) UpdateOneID(id oid.ID) *TokenUpdateOne {
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
func (c *TokenClient) DeleteOneID(id oid.ID) *TokenDeleteOne {
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
func (c *TokenClient) Get(ctx context.Context, id oid.ID) (*Token, error) {
	return c.Query().Where(token.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *TokenClient) GetX(ctx context.Context, id oid.ID) *Token {
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
