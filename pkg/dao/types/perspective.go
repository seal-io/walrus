package types

import (
	"fmt"
	"strings"

	"github.com/seal-io/seal/utils/strs"
)

// QueryCondition indicated the filters, groupBys, step, and shared costs query params.
type (
	QueryCondition struct {
		Filters     AllocationCostFilters `json:"filters,omitempty"`
		SharedCosts ShareCosts            `json:"shareCosts,omitempty"`
		GroupBy     GroupByField          `json:"groupBy,omitempty"`
		Step        Step                  `json:"step,omitempty"`
		Paging      QueryPagination       `json:"paging,omitempty"`
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
		ConnectorID ID   `json:"connectorID,omitempty"`
		IncludeAll  bool `json:"includeAll,omitempty"`
	}

	ManagementCostFilters []ManagementCostFilter
	ManagementCostFilter  struct {
		ConnectorID ID   `json:"connectorID,omitempty"`
		IncludeAll  bool `json:"includeAll,omitempty"`
	}
)

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
	FilterFieldProject        FilterField = FilterField("label:" + LabelSealProject)     // "label:seal.io/project"
	FilterFieldEnvironment    FilterField = FilterField("label:" + LabelSealEnvironment) // "label:seal.io/environment"
	FilterFieldApplication    FilterField = FilterField("label:" + LabelSealApplication) // "label:seal.io/app"
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
	GroupByFieldConnectorID    GroupByField = "connector_id"
	GroupByFieldNamespace      GroupByField = "namespace"
	GroupByFieldClusterName    GroupByField = "cluster_name"
	GroupByFieldNode           GroupByField = "node"
	GroupByFieldController     GroupByField = "controller"
	GroupByFieldControllerKind GroupByField = "controller_kind"
	GroupByFieldPod            GroupByField = "pod"
	GroupByFieldContainer      GroupByField = "container"
	GroupByFieldDay            GroupByField = "day"
	GroupByFieldWeek           GroupByField = "week"
	GroupByFieldMonth          GroupByField = "month"
	GroupByFieldYear           GroupByField = "year"
	GroupByFieldProject        GroupByField = GroupByField("label:" + LabelSealProject)     // "label:seal.io/project"
	GroupByFieldEnvironment    GroupByField = GroupByField("label:" + LabelSealEnvironment) // "label:seal.io/environment"
	GroupByFieldApplication    GroupByField = GroupByField("label:" + LabelSealApplication) // "label:seal.io/app"
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
	StepYear  Step = "year"
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

func (f *GroupByField) OrderBySQL() (string, error) {
	if f == nil {
		return "", fmt.Errorf("invalid order by: blank")
	}
	switch *f {
	case GroupByFieldDay, GroupByFieldWeek, GroupByFieldMonth, GroupByFieldYear:
		return fmt.Sprintf(`date_trunc('%s', start_time) DESC`, *f), nil
	default:
		return `SUM(total_cost) DESC`, nil
	}
}

func (f *GroupByField) GroupBySQL() (string, error) {
	if f == nil {
		return "", fmt.Errorf("invalid group by: blank")
	}

	var groupBy string
	switch {
	case f.IsLabel():
		label := strings.TrimPrefix(string(*f), LabelPrefix)
		groupBy = fmt.Sprintf(`(labels ->> '%s')`, label)
	case *f == GroupByFieldDay:
		groupBy = `date_trunc('day', start_time)`
	case *f == GroupByFieldWeek:
		groupBy = `date_trunc('week', start_time)`
	case *f == GroupByFieldMonth:
		groupBy = `date_trunc('month', start_time)`
	case *f == GroupByFieldYear:
		groupBy = `date_trunc('year', start_time)`
	default:
		groupBy = strs.Underscore(string(*f))
	}
	return groupBy, nil
}

func (f *Step) DateTruncSQL() (string, error) {
	if f == nil {
		return "", fmt.Errorf("blank step")
	}
	var dateTrunc string
	switch *f {
	case StepDay:
		dateTrunc = `date_trunc('day', start_time)`
	case StepWeek:
		dateTrunc = `date_trunc('week', start_time)`
	case StepMonth:
		dateTrunc = `date_trunc('month', start_time)`
	case StepYear:
		dateTrunc = `date_trunc('year', start_time)`
	}
	return dateTrunc, nil
}
