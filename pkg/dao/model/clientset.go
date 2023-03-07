// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package model

import "context"

// ClientSet is an interface that allows getting all clients.
type ClientSet interface {
	// AllocationCosts returns the client for interacting with the AllocationCost builders.
	AllocationCosts() *AllocationCostClient

	// Applications returns the client for interacting with the Application builders.
	Applications() *ApplicationClient

	// ApplicationInstances returns the client for interacting with the ApplicationInstance builders.
	ApplicationInstances() *ApplicationInstanceClient

	// ApplicationModuleRelationships returns the client for interacting with the ApplicationModuleRelationship builders.
	ApplicationModuleRelationships() *ApplicationModuleRelationshipClient

	// ApplicationResources returns the client for interacting with the ApplicationResource builders.
	ApplicationResources() *ApplicationResourceClient

	// ApplicationRevisions returns the client for interacting with the ApplicationRevision builders.
	ApplicationRevisions() *ApplicationRevisionClient

	// ClusterCosts returns the client for interacting with the ClusterCost builders.
	ClusterCosts() *ClusterCostClient

	// Connectors returns the client for interacting with the Connector builders.
	Connectors() *ConnectorClient

	// Environments returns the client for interacting with the Environment builders.
	Environments() *EnvironmentClient

	// EnvironmentConnectorRelationships returns the client for interacting with the EnvironmentConnectorRelationship builders.
	EnvironmentConnectorRelationships() *EnvironmentConnectorRelationshipClient

	// Modules returns the client for interacting with the Module builders.
	Modules() *ModuleClient

	// Perspectives returns the client for interacting with the Perspective builders.
	Perspectives() *PerspectiveClient

	// Projects returns the client for interacting with the Project builders.
	Projects() *ProjectClient

	// Roles returns the client for interacting with the Role builders.
	Roles() *RoleClient

	// Settings returns the client for interacting with the Setting builders.
	Settings() *SettingClient

	// Subjects returns the client for interacting with the Subject builders.
	Subjects() *SubjectClient

	// Tokens returns the client for interacting with the Token builders.
	Tokens() *TokenClient

	// WithTx gives a new transactional client in the callback function,
	// if already in a transaction, this will keep in the same transaction.
	WithTx(context.Context, func(tx *Tx) error) error

	// Dialect returns the dialect name of the driver.
	Dialect() string
}

// AllocationCostClientGetter is an interface that allows getting AllocationCostClient.
type AllocationCostClientGetter interface {
	// AllocationCosts returns the client for interacting with the AllocationCost builders.
	AllocationCosts() *AllocationCostClient
}

// ApplicationClientGetter is an interface that allows getting ApplicationClient.
type ApplicationClientGetter interface {
	// Applications returns the client for interacting with the Application builders.
	Applications() *ApplicationClient
}

// ApplicationInstanceClientGetter is an interface that allows getting ApplicationInstanceClient.
type ApplicationInstanceClientGetter interface {
	// ApplicationInstances returns the client for interacting with the ApplicationInstance builders.
	ApplicationInstances() *ApplicationInstanceClient
}

// ApplicationModuleRelationshipClientGetter is an interface that allows getting ApplicationModuleRelationshipClient.
type ApplicationModuleRelationshipClientGetter interface {
	// ApplicationModuleRelationships returns the client for interacting with the ApplicationModuleRelationship builders.
	ApplicationModuleRelationships() *ApplicationModuleRelationshipClient
}

// ApplicationResourceClientGetter is an interface that allows getting ApplicationResourceClient.
type ApplicationResourceClientGetter interface {
	// ApplicationResources returns the client for interacting with the ApplicationResource builders.
	ApplicationResources() *ApplicationResourceClient
}

// ApplicationRevisionClientGetter is an interface that allows getting ApplicationRevisionClient.
type ApplicationRevisionClientGetter interface {
	// ApplicationRevisions returns the client for interacting with the ApplicationRevision builders.
	ApplicationRevisions() *ApplicationRevisionClient
}

// ClusterCostClientGetter is an interface that allows getting ClusterCostClient.
type ClusterCostClientGetter interface {
	// ClusterCosts returns the client for interacting with the ClusterCost builders.
	ClusterCosts() *ClusterCostClient
}

// ConnectorClientGetter is an interface that allows getting ConnectorClient.
type ConnectorClientGetter interface {
	// Connectors returns the client for interacting with the Connector builders.
	Connectors() *ConnectorClient
}

// EnvironmentClientGetter is an interface that allows getting EnvironmentClient.
type EnvironmentClientGetter interface {
	// Environments returns the client for interacting with the Environment builders.
	Environments() *EnvironmentClient
}

// EnvironmentConnectorRelationshipClientGetter is an interface that allows getting EnvironmentConnectorRelationshipClient.
type EnvironmentConnectorRelationshipClientGetter interface {
	// EnvironmentConnectorRelationships returns the client for interacting with the EnvironmentConnectorRelationship builders.
	EnvironmentConnectorRelationships() *EnvironmentConnectorRelationshipClient
}

// ModuleClientGetter is an interface that allows getting ModuleClient.
type ModuleClientGetter interface {
	// Modules returns the client for interacting with the Module builders.
	Modules() *ModuleClient
}

// PerspectiveClientGetter is an interface that allows getting PerspectiveClient.
type PerspectiveClientGetter interface {
	// Perspectives returns the client for interacting with the Perspective builders.
	Perspectives() *PerspectiveClient
}

// ProjectClientGetter is an interface that allows getting ProjectClient.
type ProjectClientGetter interface {
	// Projects returns the client for interacting with the Project builders.
	Projects() *ProjectClient
}

// RoleClientGetter is an interface that allows getting RoleClient.
type RoleClientGetter interface {
	// Roles returns the client for interacting with the Role builders.
	Roles() *RoleClient
}

// SettingClientGetter is an interface that allows getting SettingClient.
type SettingClientGetter interface {
	// Settings returns the client for interacting with the Setting builders.
	Settings() *SettingClient
}

// SubjectClientGetter is an interface that allows getting SubjectClient.
type SubjectClientGetter interface {
	// Subjects returns the client for interacting with the Subject builders.
	Subjects() *SubjectClient
}

// TokenClientGetter is an interface that allows getting TokenClient.
type TokenClientGetter interface {
	// Tokens returns the client for interacting with the Token builders.
	Tokens() *TokenClient
}
