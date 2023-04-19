// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package connector

import (
	"time"

	"entgo.io/ent"

	"github.com/seal-io/seal/pkg/dao/types/crypto"
)

const (
	// Label holds the string label denoting the connector type in the database.
	Label = "connector"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldDescription holds the string denoting the description field in the database.
	FieldDescription = "description"
	// FieldLabels holds the string denoting the labels field in the database.
	FieldLabels = "labels"
	// FieldCreateTime holds the string denoting the createtime field in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the updatetime field in the database.
	FieldUpdateTime = "update_time"
	// FieldStatus holds the string denoting the status field in the database.
	FieldStatus = "status"
	// FieldType holds the string denoting the type field in the database.
	FieldType = "type"
	// FieldConfigVersion holds the string denoting the configversion field in the database.
	FieldConfigVersion = "config_version"
	// FieldConfigData holds the string denoting the configdata field in the database.
	FieldConfigData = "config_data"
	// FieldEnableFinOps holds the string denoting the enablefinops field in the database.
	FieldEnableFinOps = "enable_fin_ops"
	// FieldFinOpsCustomPricing holds the string denoting the finopscustompricing field in the database.
	FieldFinOpsCustomPricing = "fin_ops_custom_pricing"
	// FieldCategory holds the string denoting the category field in the database.
	FieldCategory = "category"
	// EdgeEnvironments holds the string denoting the environments edge name in mutations.
	EdgeEnvironments = "environments"
	// EdgeResources holds the string denoting the resources edge name in mutations.
	EdgeResources = "resources"
	// EdgeClusterCosts holds the string denoting the clustercosts edge name in mutations.
	EdgeClusterCosts = "clusterCosts"
	// EdgeAllocationCosts holds the string denoting the allocationcosts edge name in mutations.
	EdgeAllocationCosts = "allocationCosts"
	// Table holds the table name of the connector in the database.
	Table = "connectors"
	// EnvironmentsTable is the table that holds the environments relation/edge.
	EnvironmentsTable = "environment_connector_relationships"
	// EnvironmentsInverseTable is the table name for the EnvironmentConnectorRelationship entity.
	// It exists in this package in order to avoid circular dependency with the "environmentconnectorrelationship" package.
	EnvironmentsInverseTable = "environment_connector_relationships"
	// EnvironmentsColumn is the table column denoting the environments relation/edge.
	EnvironmentsColumn = "connector_id"
	// ResourcesTable is the table that holds the resources relation/edge.
	ResourcesTable = "application_resources"
	// ResourcesInverseTable is the table name for the ApplicationResource entity.
	// It exists in this package in order to avoid circular dependency with the "applicationresource" package.
	ResourcesInverseTable = "application_resources"
	// ResourcesColumn is the table column denoting the resources relation/edge.
	ResourcesColumn = "connector_id"
	// ClusterCostsTable is the table that holds the clusterCosts relation/edge.
	ClusterCostsTable = "cluster_costs"
	// ClusterCostsInverseTable is the table name for the ClusterCost entity.
	// It exists in this package in order to avoid circular dependency with the "clustercost" package.
	ClusterCostsInverseTable = "cluster_costs"
	// ClusterCostsColumn is the table column denoting the clusterCosts relation/edge.
	ClusterCostsColumn = "connector_id"
	// AllocationCostsTable is the table that holds the allocationCosts relation/edge.
	AllocationCostsTable = "allocation_costs"
	// AllocationCostsInverseTable is the table name for the AllocationCost entity.
	// It exists in this package in order to avoid circular dependency with the "allocationcost" package.
	AllocationCostsInverseTable = "allocation_costs"
	// AllocationCostsColumn is the table column denoting the allocationCosts relation/edge.
	AllocationCostsColumn = "connector_id"
)

// Columns holds all SQL columns for connector fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldDescription,
	FieldLabels,
	FieldCreateTime,
	FieldUpdateTime,
	FieldStatus,
	FieldType,
	FieldConfigVersion,
	FieldConfigData,
	FieldEnableFinOps,
	FieldFinOpsCustomPricing,
	FieldCategory,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

// Note that the variables below are initialized by the runtime
// package on the initialization of the application. Therefore,
// it should be imported in the main as follows:
//
//	import _ "github.com/seal-io/seal/pkg/dao/model/runtime"
var (
	Hooks [1]ent.Hook
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
	// DefaultLabels holds the default value on creation for the "labels" field.
	DefaultLabels map[string]string
	// DefaultCreateTime holds the default value on creation for the "createTime" field.
	DefaultCreateTime func() time.Time
	// DefaultUpdateTime holds the default value on creation for the "updateTime" field.
	DefaultUpdateTime func() time.Time
	// UpdateDefaultUpdateTime holds the default value on update for the "updateTime" field.
	UpdateDefaultUpdateTime func() time.Time
	// TypeValidator is a validator for the "type" field. It is called by the builders before save.
	TypeValidator func(string) error
	// ConfigVersionValidator is a validator for the "configVersion" field. It is called by the builders before save.
	ConfigVersionValidator func(string) error
	// DefaultConfigData holds the default value on creation for the "configData" field.
	DefaultConfigData crypto.Properties
	// CategoryValidator is a validator for the "category" field. It is called by the builders before save.
	CategoryValidator func(string) error
)

// WithoutFields returns the fields ignored the given list.
func WithoutFields(ignores ...string) []string {
	if len(ignores) == 0 {
		return Columns
	}

	var s = make(map[string]bool, len(ignores))
	for i := range ignores {
		s[ignores[i]] = true
	}

	var r = make([]string, 0, len(Columns)-len(s))
	for i := range Columns {
		if s[Columns[i]] {
			continue
		}
		r = append(r, Columns[i])
	}
	return r
}
