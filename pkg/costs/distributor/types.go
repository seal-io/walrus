package distributor

import (
	"fmt"
	"strings"
	"time"

	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/utils/strs"
)

type SharedCost struct {
	StartTime      time.Time        `json:"startTime"`
	TotalCost      float64          `json:"totalCost"`
	IdleCost       float64          `json:"idleCost"`
	ManagementCost float64          `json:"managementCost"`
	AllocationCost float64          `json:"allocationCost"`
	Condition      types.SharedCost `json:"condition"`
}

// dateTruncWithZoneOffsetSQL generate the date trunc sql from step and timezone offset, offset is in seconds east of UTC
func dateTruncWithZoneOffsetSQL(field types.Step, offset int) (string, error) {
	if field == "" {
		return "", fmt.Errorf("invalid step: blank")
	}

	var timeZone = timeZoneInPosix(offset)
	switch field {
	case types.StepDay, types.StepWeek, types.StepMonth, types.StepYear:
		return fmt.Sprintf(`date_trunc('%s', (start_time AT TIME ZONE '%s'))`, field, timeZone), nil
	default:
		return "", fmt.Errorf("invalid step: unsupport %s", field)
	}
}

// orderByWithOffsetSQL generate the order by sql with groupBy field and timezone offset, offset is in seconds east of UTC
func orderByWithOffsetSQL(field types.GroupByField, offset int) (string, error) {
	if field == "" {
		return "", fmt.Errorf("invalid order by: blank")
	}

	var timeZone = timeZoneInPosix(offset)
	switch field {
	case types.GroupByFieldDay, types.GroupByFieldWeek, types.GroupByFieldMonth, types.GroupByFieldYear:
		return fmt.Sprintf(`date_trunc('%s', start_time AT TIME ZONE '%s') DESC`, field, timeZone), nil
	default:
		return `SUM(total_cost) DESC`, nil
	}
}

// groupByWithZoneOffsetSQL generate the group by sql with timezone offset, offset is in seconds east of UTC
func groupByWithZoneOffsetSQL(field types.GroupByField, offset int) (string, error) {
	if field == "" {
		return "", fmt.Errorf("invalid group by: blank")
	}

	var (
		groupBy  string
		timeZone = timeZoneInPosix(offset)
	)
	switch {
	case field.IsLabel():
		label := strings.TrimPrefix(string(field), types.LabelPrefix)
		groupBy = fmt.Sprintf(`(labels ->> '%s')`, label)
	case field == types.GroupByFieldDay:
		groupBy = fmt.Sprintf(`date_trunc('day', (start_time AT TIME ZONE '%s'))`, timeZone)
	case field == types.GroupByFieldWeek:
		groupBy = fmt.Sprintf(`date_trunc('week', (start_time AT TIME ZONE '%s'))`, timeZone)
	case field == types.GroupByFieldMonth:
		groupBy = fmt.Sprintf(`date_trunc('month', (start_time AT TIME ZONE '%s'))`, timeZone)
	case field == types.GroupByFieldYear:
		groupBy = fmt.Sprintf(`date_trunc('year', (start_time AT TIME ZONE '%s'))`, timeZone)
	case field == types.GroupByFieldWorkload:
		groupBy = fmt.Sprintf(`CASE WHEN namespace = '' THEN '%s' 
 WHEN controller_kind = '' THEN '%s'
 WHEN controller = '' THEN '%s'
 ELSE  concat_ws('/', namespace, controller_kind, controller)
END`, types.UnallocatedLabel, types.UnallocatedLabel, types.UnallocatedLabel)
	default:
		groupBy = strs.Underscore(string(field))
	}
	return groupBy, nil
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
