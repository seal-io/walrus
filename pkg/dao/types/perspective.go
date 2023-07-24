package types

import (
	"strings"

	"github.com/seal-io/seal/pkg/dao/types/object"
)

// QueryCondition indicated the filters, groupBys, step, and shared costs query params.
type (
	QueryCondition struct {
		Filters     AllocationCostFilters `json:"filters,omitempty"`
		SharedCosts ShareCosts            `json:"shareCosts,omitempty"`
		GroupBy     GroupByField          `json:"groupBy,omitempty"`
		Step        Step                  `json:"step,omitempty"`
		Paging      QueryPagination       `json:"paging,omitempty"`
		Query       string                `json:"query"`
	}
)

// Filters: allocation, idle and management filters.
type (
	AllocationCostFilters [][]FilterRule
	FilterRule            struct {
		FieldName  FilterField `json:"fieldName,omitempty"`
		Operator   Operator    `json:"operator,omitempty"`
		Values     []string    `json:"values,omitempty"`
		IncludeAll bool        `json:"includeAll,omitempty"`
	}

	ShareCosts []SharedCost
	SharedCost struct {
		Filters               AllocationCostFilters `json:"filters,omitempty"`
		IdleCostFilters       IdleCostFilters       `json:"idleCostFilters,omitempty"`
		ManagementCostFilters ManagementCostFilters `json:"managementCostFilters,omitempty"`
		SharingStrategy       SharingStrategy       `json:"sharingStrategy,omitempty"`
	}

	IdleCostFilters []IdleCostFilter
	IdleCostFilter  struct {
		ConnectorID object.ID `json:"connectorID,omitempty"`
		IncludeAll  bool      `json:"includeAll,omitempty"`
	}

	ManagementCostFilters []ManagementCostFilter
	ManagementCostFilter  struct {
		ConnectorID object.ID `json:"connectorID,omitempty"`
		IncludeAll  bool      `json:"includeAll,omitempty"`
	}
)

// Operator for filter rule.
type Operator string

const (
	OperatorIn    Operator = "in"
	OperatorNotIn Operator = "notin"
)

// QueryPagination indicate the pagination query config.
type QueryPagination struct {
	Page    int `json:"page,omitempty"`
	PerPage int `json:"perPage,omitempty"`
}

// FilterField indicate type for filter field.
type FilterField string

// built-in filter fieldName.
const (
	FilterFieldConnectorID    FilterField = "connector_id"
	FilterFieldNamespace      FilterField = "namespace"
	FilterFieldClusterName    FilterField = "cluster_name"
	FilterFieldNode           FilterField = "node"
	FilterFieldController     FilterField = "controller"
	FilterFieldControllerKind FilterField = "controller_kind"
	FilterFieldPod            FilterField = "pod"
	FilterFieldContainer      FilterField = "container"

	// FilterFieldProject is "label:seal.io/project-name".
	FilterFieldProject = FilterField("label:" + LabelSealProjectName)

	// FilterFieldEnvironmentPath is "label:seal.io/environment-name".
	FilterFieldEnvironmentPath = FilterField("label:" + LabelSealEnvironmentPath)

	// FilterFieldServicePath is "label:seal.io/service-name".
	FilterFieldServicePath = FilterField("label:" + LabelSealServicePath)
)

func (f *FilterField) IsLabel() bool {
	if f == nil {
		return false
	}

	return strings.HasPrefix(string(*f), LabelPrefix)
}

// GroupByField indicate type for groupBy field.
type GroupByField string

// built-in groupBy fieldName.
const (
	// Properties.
	GroupByFieldConnectorID    GroupByField = "connector_id"
	GroupByFieldNamespace      GroupByField = "namespace"
	GroupByFieldClusterName    GroupByField = "cluster_name"
	GroupByFieldNode           GroupByField = "node"
	GroupByFieldController     GroupByField = "controller"
	GroupByFieldControllerKind GroupByField = "controller_kind"
	GroupByFieldPod            GroupByField = "pod"
	GroupByFieldContainer      GroupByField = "container"
	GroupByFieldWorkload       GroupByField = "workload"

	// Time bucket.
	GroupByFieldDay   GroupByField = "day"
	GroupByFieldWeek  GroupByField = "week"
	GroupByFieldMonth GroupByField = "month"

	// Built-in labels.
	GroupByFieldProject         = GroupByField("label:" + LabelSealProjectName)     // "label:seal.io/project-name".
	GroupByFieldEnvironment     = GroupByField("label:" + LabelSealEnvironmentName) // "label:seal.io/environment-name".
	GroupByFieldService         = GroupByField("label:" + LabelSealServiceName)     // "label:seal.io/service-name".
	GroupByFieldEnvironmentPath = GroupByField("label:" + LabelSealEnvironmentPath) // "label:seal.io/environment-path".
	GroupByFieldServicePath     = GroupByField("label:" + LabelSealServicePath)     // "label:seal.io/service-path".
)

func (f *GroupByField) IsLabel() bool {
	if f == nil {
		return false
	}

	return strings.HasPrefix(string(*f), LabelPrefix)
}

// SharingStrategy indicate the share cost strategy.
type SharingStrategy string

const (
	SharingStrategyEqually        SharingStrategy = "equally"
	SharingStrategyProportionally SharingStrategy = "proportionally"
)

// Step indicate the time step to aggregate cost.
type Step string

const (
	StepDay   Step = "day"
	StepWeek  Step = "week"
	StepMonth Step = "month"
)
