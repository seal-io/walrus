package types

import (
	"strings"
)

// QueryCondition indicated the filters, groupBys, step, and shared costs query params.
type (
	QueryCondition struct {
		Filters       CostFilters        `json:"filters,omitempty"`
		GroupBy       GroupByField       `json:"groupBy,omitempty"`
		Step          Step               `json:"step,omitempty"`
		Paging        QueryPagination    `json:"paging,omitempty"`
		Query         string             `json:"query"`
		SharedOptions *SharedCostOptions `json:"sharedOptions,omitempty"`
	}
)

type (
	// CostFilters indicate the filters for cost query.
	CostFilters [][]FilterRule

	// FilterRule indicate the filter rule for cost query.
	FilterRule struct {
		FieldName  FilterField `json:"fieldName,omitempty"`
		Operator   Operator    `json:"operator,omitempty"`
		Values     []string    `json:"values,omitempty"`
		IncludeAll bool        `json:"includeAll,omitempty"`
	}

	// SharedCostOptions indicate the shared cost options for shared cost query.
	SharedCostOptions struct {
		Item       ItemSharedOptions      `json:"item,omitempty"`
		Idle       *IdleShareOption       `json:"idle,omitempty"`
		Management *ManagementShareOption `json:"management,omitempty"`
	}

	// ItemSharedOptions indicate the shared cost options for custom item cost query.
	ItemSharedOptions []ItemSharedOption

	// ItemSharedOption indicate the shared cost option for custom item cost query.
	ItemSharedOption struct {
		Filters         CostFilters     `json:"filters,omitempty"`
		SharingStrategy SharingStrategy `json:"sharingStrategy,omitempty"`
	}

	// IdleShareOption indicate the shared cost option for idle cost.
	IdleShareOption struct {
		SharingStrategy SharingStrategy `json:"sharingStrategy,omitempty"`
	}

	// ManagementShareOption indicate the shared cost option for management cost.
	ManagementShareOption struct {
		SharingStrategy SharingStrategy `json:"sharingStrategy,omitempty"`
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
	FilterFieldName           FilterField = "name"

	// FilterFieldProject is "label:seal.io/project-name".
	FilterFieldProject = FilterField("label:" + LabelWalrusProjectName)

	// FilterFieldEnvironmentPath is "label:seal.io/environment-name".
	FilterFieldEnvironmentPath = FilterField("label:" + LabelWalrusEnvironmentPath)

	// FilterFieldServicePath is "label:seal.io/service-name".
	FilterFieldServicePath = FilterField("label:" + LabelWalrusServicePath)
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
	GroupByFieldProject     = GroupByField("label:" + LabelWalrusProjectName) // "label:walrus.seal.io/project-name".
	GroupByFieldEnvironment = GroupByField(
		"label:" + LabelWalrusEnvironmentName,
	) // "label:walrus.seal.io/environment-name".
	GroupByFieldService = GroupByField(
		"label:" + LabelWalrusServiceName,
	) // "label:walrus.seal.io/service-name".
	GroupByFieldEnvironmentPath = GroupByField(
		"label:" + LabelWalrusEnvironmentPath,
	) // "label:walrus.seal.io/environment-path".
	GroupByFieldServicePath = GroupByField("label:" + LabelWalrusServicePath) // "label:walrus.seal.io/service-path".
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
