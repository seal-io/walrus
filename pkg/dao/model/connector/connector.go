// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package connector

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"

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

// OrderOption defines the ordering options for the Connector queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByDescription orders the results by the description field.
func ByDescription(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDescription, opts...).ToFunc()
}

// ByCreateTime orders the results by the createTime field.
func ByCreateTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreateTime, opts...).ToFunc()
}

// ByUpdateTime orders the results by the updateTime field.
func ByUpdateTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdateTime, opts...).ToFunc()
}

// ByType orders the results by the type field.
func ByType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldType, opts...).ToFunc()
}

// ByConfigVersion orders the results by the configVersion field.
func ByConfigVersion(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldConfigVersion, opts...).ToFunc()
}

// ByConfigData orders the results by the configData field.
func ByConfigData(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldConfigData, opts...).ToFunc()
}

// ByEnableFinOps orders the results by the enableFinOps field.
func ByEnableFinOps(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldEnableFinOps, opts...).ToFunc()
}

// ByCategory orders the results by the category field.
func ByCategory(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCategory, opts...).ToFunc()
}

// ByEnvironmentsCount orders the results by environments count.
func ByEnvironmentsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newEnvironmentsStep(), opts...)
	}
}

// ByEnvironments orders the results by environments terms.
func ByEnvironments(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newEnvironmentsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByResourcesCount orders the results by resources count.
func ByResourcesCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newResourcesStep(), opts...)
	}
}

// ByResources orders the results by resources terms.
func ByResources(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newResourcesStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByClusterCostsCount orders the results by clusterCosts count.
func ByClusterCostsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newClusterCostsStep(), opts...)
	}
}

// ByClusterCosts orders the results by clusterCosts terms.
func ByClusterCosts(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newClusterCostsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByAllocationCostsCount orders the results by allocationCosts count.
func ByAllocationCostsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newAllocationCostsStep(), opts...)
	}
}

// ByAllocationCosts orders the results by allocationCosts terms.
func ByAllocationCosts(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newAllocationCostsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newEnvironmentsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(EnvironmentsInverseTable, EnvironmentsColumn),
		sqlgraph.Edge(sqlgraph.O2M, true, EnvironmentsTable, EnvironmentsColumn),
	)
}
func newResourcesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ResourcesInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, ResourcesTable, ResourcesColumn),
	)
}
func newClusterCostsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ClusterCostsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, ClusterCostsTable, ClusterCostsColumn),
	)
}
func newAllocationCostsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(AllocationCostsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, AllocationCostsTable, AllocationCostsColumn),
	)
}

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
