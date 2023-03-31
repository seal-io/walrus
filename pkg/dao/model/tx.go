// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package model

import (
	"context"
	stdsql "database/sql"
	"fmt"
	"sync"

	"entgo.io/ent/dialect"
)

// Tx is a transactional client that is created by calling Client.Tx().
type Tx struct {
	config
	// AllocationCost is the client for interacting with the AllocationCost builders.
	AllocationCost *AllocationCostClient
	// Application is the client for interacting with the Application builders.
	Application *ApplicationClient
	// ApplicationInstance is the client for interacting with the ApplicationInstance builders.
	ApplicationInstance *ApplicationInstanceClient
	// ApplicationModuleRelationship is the client for interacting with the ApplicationModuleRelationship builders.
	ApplicationModuleRelationship *ApplicationModuleRelationshipClient
	// ApplicationResource is the client for interacting with the ApplicationResource builders.
	ApplicationResource *ApplicationResourceClient
	// ApplicationRevision is the client for interacting with the ApplicationRevision builders.
	ApplicationRevision *ApplicationRevisionClient
	// ClusterCost is the client for interacting with the ClusterCost builders.
	ClusterCost *ClusterCostClient
	// Connector is the client for interacting with the Connector builders.
	Connector *ConnectorClient
	// Environment is the client for interacting with the Environment builders.
	Environment *EnvironmentClient
	// EnvironmentConnectorRelationship is the client for interacting with the EnvironmentConnectorRelationship builders.
	EnvironmentConnectorRelationship *EnvironmentConnectorRelationshipClient
	// Module is the client for interacting with the Module builders.
	Module *ModuleClient
	// ModuleVersion is the client for interacting with the ModuleVersion builders.
	ModuleVersion *ModuleVersionClient
	// Perspective is the client for interacting with the Perspective builders.
	Perspective *PerspectiveClient
	// Project is the client for interacting with the Project builders.
	Project *ProjectClient
	// Role is the client for interacting with the Role builders.
	Role *RoleClient
	// Secret is the client for interacting with the Secret builders.
	Secret *SecretClient
	// Setting is the client for interacting with the Setting builders.
	Setting *SettingClient
	// Subject is the client for interacting with the Subject builders.
	Subject *SubjectClient
	// Token is the client for interacting with the Token builders.
	Token *TokenClient

	// lazily loaded.
	client     *Client
	clientOnce sync.Once
	// ctx lives for the life of the transaction. It is
	// the same context used by the underlying connection.
	ctx context.Context
}

type (
	// Committer is the interface that wraps the Commit method.
	Committer interface {
		Commit(context.Context, *Tx) error
	}

	// The CommitFunc type is an adapter to allow the use of ordinary
	// function as a Committer. If f is a function with the appropriate
	// signature, CommitFunc(f) is a Committer that calls f.
	CommitFunc func(context.Context, *Tx) error

	// CommitHook defines the "commit middleware". A function that gets a Committer
	// and returns a Committer. For example:
	//
	//	hook := func(next ent.Committer) ent.Committer {
	//		return ent.CommitFunc(func(ctx context.Context, tx *ent.Tx) error {
	//			// Do some stuff before.
	//			if err := next.Commit(ctx, tx); err != nil {
	//				return err
	//			}
	//			// Do some stuff after.
	//			return nil
	//		})
	//	}
	//
	CommitHook func(Committer) Committer
)

// Commit calls f(ctx, m).
func (f CommitFunc) Commit(ctx context.Context, tx *Tx) error {
	return f(ctx, tx)
}

// Commit commits the transaction.
func (tx *Tx) Commit() error {
	txDriver := tx.config.driver.(*txDriver)
	var fn Committer = CommitFunc(func(context.Context, *Tx) error {
		return txDriver.tx.Commit()
	})
	txDriver.mu.Lock()
	hooks := append([]CommitHook(nil), txDriver.onCommit...)
	txDriver.mu.Unlock()
	for i := len(hooks) - 1; i >= 0; i-- {
		fn = hooks[i](fn)
	}
	return fn.Commit(tx.ctx, tx)
}

// OnCommit adds a hook to call on commit.
func (tx *Tx) OnCommit(f CommitHook) {
	txDriver := tx.config.driver.(*txDriver)
	txDriver.mu.Lock()
	txDriver.onCommit = append(txDriver.onCommit, f)
	txDriver.mu.Unlock()
}

type (
	// Rollbacker is the interface that wraps the Rollback method.
	Rollbacker interface {
		Rollback(context.Context, *Tx) error
	}

	// The RollbackFunc type is an adapter to allow the use of ordinary
	// function as a Rollbacker. If f is a function with the appropriate
	// signature, RollbackFunc(f) is a Rollbacker that calls f.
	RollbackFunc func(context.Context, *Tx) error

	// RollbackHook defines the "rollback middleware". A function that gets a Rollbacker
	// and returns a Rollbacker. For example:
	//
	//	hook := func(next ent.Rollbacker) ent.Rollbacker {
	//		return ent.RollbackFunc(func(ctx context.Context, tx *ent.Tx) error {
	//			// Do some stuff before.
	//			if err := next.Rollback(ctx, tx); err != nil {
	//				return err
	//			}
	//			// Do some stuff after.
	//			return nil
	//		})
	//	}
	//
	RollbackHook func(Rollbacker) Rollbacker
)

// Rollback calls f(ctx, m).
func (f RollbackFunc) Rollback(ctx context.Context, tx *Tx) error {
	return f(ctx, tx)
}

// Rollback rollbacks the transaction.
func (tx *Tx) Rollback() error {
	txDriver := tx.config.driver.(*txDriver)
	var fn Rollbacker = RollbackFunc(func(context.Context, *Tx) error {
		return txDriver.tx.Rollback()
	})
	txDriver.mu.Lock()
	hooks := append([]RollbackHook(nil), txDriver.onRollback...)
	txDriver.mu.Unlock()
	for i := len(hooks) - 1; i >= 0; i-- {
		fn = hooks[i](fn)
	}
	return fn.Rollback(tx.ctx, tx)
}

// OnRollback adds a hook to call on rollback.
func (tx *Tx) OnRollback(f RollbackHook) {
	txDriver := tx.config.driver.(*txDriver)
	txDriver.mu.Lock()
	txDriver.onRollback = append(txDriver.onRollback, f)
	txDriver.mu.Unlock()
}

// Client returns a Client that binds to current transaction.
func (tx *Tx) Client() *Client {
	tx.clientOnce.Do(func() {
		tx.client = &Client{config: tx.config}
		tx.client.init()
	})
	return tx.client
}

func (tx *Tx) init() {
	tx.AllocationCost = NewAllocationCostClient(tx.config)
	tx.Application = NewApplicationClient(tx.config)
	tx.ApplicationInstance = NewApplicationInstanceClient(tx.config)
	tx.ApplicationModuleRelationship = NewApplicationModuleRelationshipClient(tx.config)
	tx.ApplicationResource = NewApplicationResourceClient(tx.config)
	tx.ApplicationRevision = NewApplicationRevisionClient(tx.config)
	tx.ClusterCost = NewClusterCostClient(tx.config)
	tx.Connector = NewConnectorClient(tx.config)
	tx.Environment = NewEnvironmentClient(tx.config)
	tx.EnvironmentConnectorRelationship = NewEnvironmentConnectorRelationshipClient(tx.config)
	tx.Module = NewModuleClient(tx.config)
	tx.ModuleVersion = NewModuleVersionClient(tx.config)
	tx.Perspective = NewPerspectiveClient(tx.config)
	tx.Project = NewProjectClient(tx.config)
	tx.Role = NewRoleClient(tx.config)
	tx.Secret = NewSecretClient(tx.config)
	tx.Setting = NewSettingClient(tx.config)
	tx.Subject = NewSubjectClient(tx.config)
	tx.Token = NewTokenClient(tx.config)
}

// txDriver wraps the given dialect.Tx with a nop dialect.Driver implementation.
// The idea is to support transactions without adding any extra code to the builders.
// When a builder calls to driver.Tx(), it gets the same dialect.Tx instance.
// Commit and Rollback are nop for the internal builders and the user must call one
// of them in order to commit or rollback the transaction.
//
// If a closed transaction is embedded in one of the generated entities, and the entity
// applies a query, for example: AllocationCost.QueryXXX(), the query will be executed
// through the driver which created this transaction.
//
// Note that txDriver is not goroutine safe.
type txDriver struct {
	// the driver we started the transaction from.
	drv dialect.Driver
	// tx is the underlying transaction.
	tx dialect.Tx
	// completion hooks.
	mu         sync.Mutex
	onCommit   []CommitHook
	onRollback []RollbackHook
}

// newTx creates a new transactional driver.
func newTx(ctx context.Context, drv dialect.Driver) (*txDriver, error) {
	tx, err := drv.Tx(ctx)
	if err != nil {
		return nil, err
	}
	return &txDriver{tx: tx, drv: drv}, nil
}

// Tx returns the transaction wrapper (txDriver) to avoid Commit or Rollback calls
// from the internal builders. Should be called only by the internal builders.
func (tx *txDriver) Tx(context.Context) (dialect.Tx, error) { return tx, nil }

// Dialect returns the dialect of the driver we started the transaction from.
func (tx *txDriver) Dialect() string { return tx.drv.Dialect() }

// Close is a nop close.
func (*txDriver) Close() error { return nil }

// Commit is a nop commit for the internal builders.
// User must call `Tx.Commit` in order to commit the transaction.
func (*txDriver) Commit() error { return nil }

// Rollback is a nop rollback for the internal builders.
// User must call `Tx.Rollback` in order to rollback the transaction.
func (*txDriver) Rollback() error { return nil }

// Exec calls tx.Exec.
func (tx *txDriver) Exec(ctx context.Context, query string, args, v any) error {
	return tx.tx.Exec(ctx, query, args, v)
}

// Query calls tx.Query.
func (tx *txDriver) Query(ctx context.Context, query string, args, v any) error {
	return tx.tx.Query(ctx, query, args, v)
}

var _ dialect.Driver = (*txDriver)(nil)

// AllocationCosts implements the ClientSet.
func (tx *Tx) AllocationCosts() *AllocationCostClient {
	return tx.AllocationCost
}

// Applications implements the ClientSet.
func (tx *Tx) Applications() *ApplicationClient {
	return tx.Application
}

// ApplicationInstances implements the ClientSet.
func (tx *Tx) ApplicationInstances() *ApplicationInstanceClient {
	return tx.ApplicationInstance
}

// ApplicationModuleRelationships implements the ClientSet.
func (tx *Tx) ApplicationModuleRelationships() *ApplicationModuleRelationshipClient {
	return tx.ApplicationModuleRelationship
}

// ApplicationResources implements the ClientSet.
func (tx *Tx) ApplicationResources() *ApplicationResourceClient {
	return tx.ApplicationResource
}

// ApplicationRevisions implements the ClientSet.
func (tx *Tx) ApplicationRevisions() *ApplicationRevisionClient {
	return tx.ApplicationRevision
}

// ClusterCosts implements the ClientSet.
func (tx *Tx) ClusterCosts() *ClusterCostClient {
	return tx.ClusterCost
}

// Connectors implements the ClientSet.
func (tx *Tx) Connectors() *ConnectorClient {
	return tx.Connector
}

// Environments implements the ClientSet.
func (tx *Tx) Environments() *EnvironmentClient {
	return tx.Environment
}

// EnvironmentConnectorRelationships implements the ClientSet.
func (tx *Tx) EnvironmentConnectorRelationships() *EnvironmentConnectorRelationshipClient {
	return tx.EnvironmentConnectorRelationship
}

// Modules implements the ClientSet.
func (tx *Tx) Modules() *ModuleClient {
	return tx.Module
}

// ModuleVersions implements the ClientSet.
func (tx *Tx) ModuleVersions() *ModuleVersionClient {
	return tx.ModuleVersion
}

// Perspectives implements the ClientSet.
func (tx *Tx) Perspectives() *PerspectiveClient {
	return tx.Perspective
}

// Projects implements the ClientSet.
func (tx *Tx) Projects() *ProjectClient {
	return tx.Project
}

// Roles implements the ClientSet.
func (tx *Tx) Roles() *RoleClient {
	return tx.Role
}

// Secrets implements the ClientSet.
func (tx *Tx) Secrets() *SecretClient {
	return tx.Secret
}

// Settings implements the ClientSet.
func (tx *Tx) Settings() *SettingClient {
	return tx.Setting
}

// Subjects implements the ClientSet.
func (tx *Tx) Subjects() *SubjectClient {
	return tx.Subject
}

// Tokens implements the ClientSet.
func (tx *Tx) Tokens() *TokenClient {
	return tx.Token
}

// Debug returns the debug value of the driver.
func (tx *Tx) Debug() *Client {
	return tx.client.Debug()
}

// Dialect returns the dialect name of the driver.
func (tx *Tx) Dialect() string {
	return tx.driver.Dialect()
}

// ExecContext allows calling the underlying ExecContext method of the transaction if it is supported by it.
// See, database/sql#Tx.ExecContext for more information.
func (tx *txDriver) ExecContext(ctx context.Context, query string, args ...any) (stdsql.Result, error) {
	ex, ok := tx.tx.(interface {
		ExecContext(context.Context, string, ...any) (stdsql.Result, error)
	})
	if !ok {
		return nil, fmt.Errorf("Tx.ExecContext is not supported")
	}
	return ex.ExecContext(ctx, query, args...)
}

// QueryContext allows calling the underlying QueryContext method of the transaction if it is supported by it.
// See, database/sql#Tx.QueryContext for more information.
func (tx *txDriver) QueryContext(ctx context.Context, query string, args ...any) (*stdsql.Rows, error) {
	q, ok := tx.tx.(interface {
		QueryContext(context.Context, string, ...any) (*stdsql.Rows, error)
	})
	if !ok {
		return nil, fmt.Errorf("Tx.QueryContext is not supported")
	}
	return q.QueryContext(ctx, query, args...)
}

// Use adds the mutation hooks to all the entity clients.
func (tx *Tx) Use(hooks ...Hook) {
	tx.AllocationCost.Use(hooks...)
	tx.Application.Use(hooks...)
	tx.ApplicationInstance.Use(hooks...)
	tx.ApplicationModuleRelationship.Use(hooks...)
	tx.ApplicationResource.Use(hooks...)
	tx.ApplicationRevision.Use(hooks...)
	tx.ClusterCost.Use(hooks...)
	tx.Connector.Use(hooks...)
	tx.Environment.Use(hooks...)
	tx.EnvironmentConnectorRelationship.Use(hooks...)
	tx.Module.Use(hooks...)
	tx.ModuleVersion.Use(hooks...)
	tx.Perspective.Use(hooks...)
	tx.Project.Use(hooks...)
	tx.Role.Use(hooks...)
	tx.Secret.Use(hooks...)
	tx.Setting.Use(hooks...)
	tx.Subject.Use(hooks...)
	tx.Token.Use(hooks...)
}

// WithTx gives a new transactional client in the callback function,
// if already in a transaction, this will keep in the same transaction.
func (tx *Tx) WithTx(ctx context.Context, callback func(tx *Tx) error) error {
	return callback(tx)
}
