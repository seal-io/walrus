// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package internal

import "context"

// SchemaConfig represents alternative schema names for all tables
// that can be passed at runtime.
type SchemaConfig struct {
	AllocationCost                   string // AllocationCost table.
	Application                      string // Application table.
	ApplicationModuleRelationship    string // ApplicationModuleRelationship table.
	ApplicationResource              string // ApplicationResource table.
	ApplicationRevision              string // ApplicationRevision table.
	ClusterCost                      string // ClusterCost table.
	Connector                        string // Connector table.
	Environment                      string // Environment table.
	EnvironmentConnectorRelationship string // EnvironmentConnectorRelationship table.
	Module                           string // Module table.
	Perspective                      string // Perspective table.
	Project                          string // Project table.
	Role                             string // Role table.
	Setting                          string // Setting table.
	Subject                          string // Subject table.
	Token                            string // Token table.
}

type schemaCtxKey struct{}

// SchemaConfigFromContext returns a SchemaConfig stored inside a context, or empty if there isn't one.
func SchemaConfigFromContext(ctx context.Context) SchemaConfig {
	config, _ := ctx.Value(schemaCtxKey{}).(SchemaConfig)
	return config
}

// NewSchemaConfigContext returns a new context with the given SchemaConfig attached.
func NewSchemaConfigContext(parent context.Context, config SchemaConfig) context.Context {
	return context.WithValue(parent, schemaCtxKey{}, config)
}
