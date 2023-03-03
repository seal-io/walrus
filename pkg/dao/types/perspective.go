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
	// properties
	GroupByFieldConnectorID    GroupByField = "connector_id"
	GroupByFieldNamespace      GroupByField = "namespace"
	GroupByFieldClusterName    GroupByField = "cluster_name"
	GroupByFieldNode           GroupByField = "node"
	GroupByFieldController     GroupByField = "controller"
	GroupByFieldControllerKind GroupByField = "controller_kind"
	GroupByFieldPod            GroupByField = "pod"
	GroupByFieldContainer      GroupByField = "container"
	GroupByFieldWorkload       GroupByField = "workload"

	// time bucket
	GroupByFieldDay   GroupByField = "day"
	GroupByFieldWeek  GroupByField = "week"
	GroupByFieldMonth GroupByField = "month"
	GroupByFieldYear  GroupByField = "year"

	// built-in labels
	GroupByFieldProject     GroupByField = GroupByField("label:" + LabelSealProject)     // "label:seal.io/project"
	GroupByFieldEnvironment GroupByField = GroupByField("label:" + LabelSealEnvironment) // "label:seal.io/environment"
	GroupByFieldApplication GroupByField = GroupByField("label:" + LabelSealApplication) // "label:seal.io/app"
)

func (f *GroupByField) IsLabel() bool {
	if f == nil {
		return false
	}
	return strings.HasPrefix(string(*f), LabelPrefix)
}

func (f *GroupByField) OrderBySQL() (string, error) {
	if f == nil {
		return "", fmt.Errorf("invalid order by: blank")
	}
	return f.OrderByWithOffsetSQL(0)
}

func (f *GroupByField) OrderByWithOffsetSQL(offset int) (string, error) {
	if f == nil {
		return "", fmt.Errorf("invalid order by: blank")
	}

	var timeZone = timeZoneInPosix(offset)
	switch *f {
	case GroupByFieldDay, GroupByFieldWeek, GroupByFieldMonth, GroupByFieldYear:
		return fmt.Sprintf(`date_trunc('%s', start_time AT TIME ZONE '%s') DESC`, *f, timeZone), nil
	default:
		return `SUM(total_cost) DESC`, nil
	}
}

func (f *GroupByField) GroupBySQL() (string, error) {
	if f == nil {
		return "", fmt.Errorf("invalid group by: blank")
	}

	return f.GroupByWithZoneOffsetSQL(0)
}

// GroupByWithZoneOffsetSQL generate the group by sql with timezone offset, offset is in seconds east of UTC
func (f *GroupByField) GroupByWithZoneOffsetSQL(offset int) (string, error) {
	if f == nil {
		return "", fmt.Errorf("invalid group by: blank")
	}

	var (
		groupBy  string
		timeZone = timeZoneInPosix(offset)
	)
	switch {
	case f.IsLabel():
		label := strings.TrimPrefix(string(*f), LabelPrefix)
		groupBy = fmt.Sprintf(`(labels ->> '%s')`, label)
	case *f == GroupByFieldDay:
		groupBy = fmt.Sprintf(`date_trunc('day', (start_time AT TIME ZONE '%s'))`, timeZone)
	case *f == GroupByFieldWeek:
		groupBy = fmt.Sprintf(`date_trunc('week', (start_time AT TIME ZONE '%s'))`, timeZone)
	case *f == GroupByFieldMonth:
		groupBy = fmt.Sprintf(`date_trunc('month', (start_time AT TIME ZONE '%s'))`, timeZone)
	case *f == GroupByFieldYear:
		groupBy = fmt.Sprintf(`date_trunc('year', (start_time AT TIME ZONE '%s'))`, timeZone)
	case *f == GroupByFieldWorkload:
		groupBy = fmt.Sprintf(`CASE WHEN namespace = '' THEN '%s' 
 WHEN controller_kind = '' THEN '%s'
 WHEN controller = '' THEN '%s'
 ELSE  concat_ws('/', namespace, controller_kind, controller)
END`, UnallocatedLabel, UnallocatedLabel, UnallocatedLabel)
	default:
		groupBy = strs.Underscore(string(*f))
	}
	return groupBy, nil
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

// DateTruncSQL generate the date trunc sql from step
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

// DateTruncWithZoneOffsetSQL generate the date trunc sql from step and timezone offset, offset is in seconds east of UTC
func (f *Step) DateTruncWithZoneOffsetSQL(offset int) (string, error) {
	if f == nil {
		return "", fmt.Errorf("invalid step: blank")
	}
	var (
		dateTrunc string
		timeZone  = timeZoneInPosix(offset)
	)
	switch *f {
	case StepDay:
		dateTrunc = fmt.Sprintf(`date_trunc('day', (start_time AT TIME ZONE '%s'))`, timeZone)
	case StepWeek:
		dateTrunc = fmt.Sprintf(`date_trunc('week', (start_time AT TIME ZONE '%s'))`, timeZone)
	case StepMonth:
		dateTrunc = fmt.Sprintf(`date_trunc('month', (start_time AT TIME ZONE '%s'))`, timeZone)
	case StepYear:
		dateTrunc = fmt.Sprintf(`date_trunc('year', (start_time AT TIME ZONE '%s'))`, timeZone)
	}
	return dateTrunc, nil
}

// timeZoneInPosix is in posix timezone string format
// time zone Asia/Shanghai in posix is UTC-8
func timeZoneInPosix(offset int) string {
	timeZone := "UTC"
	if offset != 0 {
		utcOffSig := "-"
		utcOffHrs := offset / 60 / 60

		if utcOffHrs < 0 {
			utcOffSig = "+"
			utcOffHrs = 0 - utcOffHrs
		}

		timeZone = fmt.Sprintf("UTC%s%d", utcOffSig, utcOffHrs)
	}
	return timeZone
}
