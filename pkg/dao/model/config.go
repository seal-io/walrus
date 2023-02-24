// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package model

import (
	"context"
	stdsql "database/sql"
	"fmt"

	"entgo.io/ent"
	"entgo.io/ent/dialect"

	"github.com/seal-io/seal/pkg/dao/model/internal"
)

// Option function to configure the client.
type Option func(*config)

// Config is the configuration for the client and its builder.
type config struct {
	// driver used for executing database requests.
	driver dialect.Driver
	// debug enable a debug logging.
	debug bool
	// log used for logging on debug mode.
	log func(...any)
	// hooks to execute on mutations.
	hooks *hooks
	// interceptors to execute on queries.
	inters *inters
	// schemaConfig contains alternative names for all tables.
	schemaConfig SchemaConfig
}

// hooks and interceptors per client, for fast access.
type (
	hooks struct {
		AllocationCost                   []ent.Hook
		Application                      []ent.Hook
		ApplicationModuleRelationship    []ent.Hook
		ApplicationResource              []ent.Hook
		ApplicationRevision              []ent.Hook
		ClusterCost                      []ent.Hook
		Connector                        []ent.Hook
		Environment                      []ent.Hook
		EnvironmentConnectorRelationship []ent.Hook
		Module                           []ent.Hook
		Perspective                      []ent.Hook
		Project                          []ent.Hook
		Role                             []ent.Hook
		Setting                          []ent.Hook
		Subject                          []ent.Hook
		Token                            []ent.Hook
	}
	inters struct {
		AllocationCost                   []ent.Interceptor
		Application                      []ent.Interceptor
		ApplicationModuleRelationship    []ent.Interceptor
		ApplicationResource              []ent.Interceptor
		ApplicationRevision              []ent.Interceptor
		ClusterCost                      []ent.Interceptor
		Connector                        []ent.Interceptor
		Environment                      []ent.Interceptor
		EnvironmentConnectorRelationship []ent.Interceptor
		Module                           []ent.Interceptor
		Perspective                      []ent.Interceptor
		Project                          []ent.Interceptor
		Role                             []ent.Interceptor
		Setting                          []ent.Interceptor
		Subject                          []ent.Interceptor
		Token                            []ent.Interceptor
	}
)

// Options applies the options on the config object.
func (c *config) options(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
	if c.debug {
		c.driver = dialect.Debug(c.driver, c.log)
	}
}

// Debug enables debug logging on the ent.Driver.
func Debug() Option {
	return func(c *config) {
		c.debug = true
	}
}

// Log sets the logging function for debug mode.
func Log(fn func(...any)) Option {
	return func(c *config) {
		c.log = fn
	}
}

// Driver configures the client driver.
func Driver(driver dialect.Driver) Option {
	return func(c *config) {
		c.driver = driver
	}
}

// SchemaConfig represents alternative schema names for all tables
// that can be passed at runtime.
type SchemaConfig = internal.SchemaConfig

// AlternateSchemas allows alternate schema names to be
// passed into ent operations.
func AlternateSchema(schemaConfig SchemaConfig) Option {
	return func(c *config) {
		c.schemaConfig = schemaConfig
	}
}

// ExecContext allows calling the underlying ExecContext method of the driver if it is supported by it.
// See, database/sql#DB.ExecContext for more information.
func (c *config) ExecContext(ctx context.Context, query string, args ...any) (stdsql.Result, error) {
	ex, ok := c.driver.(interface {
		ExecContext(context.Context, string, ...any) (stdsql.Result, error)
	})
	if !ok {
		return nil, fmt.Errorf("Driver.ExecContext is not supported")
	}
	return ex.ExecContext(ctx, query, args...)
}

// QueryContext allows calling the underlying QueryContext method of the driver if it is supported by it.
// See, database/sql#DB.QueryContext for more information.
func (c *config) QueryContext(ctx context.Context, query string, args ...any) (*stdsql.Rows, error) {
	q, ok := c.driver.(interface {
		QueryContext(context.Context, string, ...any) (*stdsql.Rows, error)
	})
	if !ok {
		return nil, fmt.Errorf("Driver.QueryContext is not supported")
	}
	return q.QueryContext(ctx, query, args...)
}
