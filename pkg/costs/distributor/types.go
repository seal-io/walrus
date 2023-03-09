package distributor

import (
	"context"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/allocationcost"
	"github.com/seal-io/seal/pkg/dao/model/connector"
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

// DateTruncWithZoneOffsetSQL generate the date trunc sql from step and timezone offset, offset is in seconds east of UTC
func DateTruncWithZoneOffsetSQL(field types.Step, offset int) (string, error) {
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

// havingSQL generate the having sql with group by and query keyword
func havingSQL(
	ctx context.Context,
	client model.ClientSet,
	groupBy types.GroupByField,
	groupBySQL string,
	query string,
) (*sql.Predicate, error) {
	if query == "" {
		return nil, fmt.Errorf("invalid query: blank")
	}
	if groupBy == "" || groupBySQL == "" {
		return nil, fmt.Errorf("invalid group by: blank")
	}

	var having *sql.Predicate
	switch {
	case groupBy == types.GroupByFieldConnectorID:
		connIDs, err := client.Connectors().Query().
			Where(
				connector.NameContainsFold(query),
				connector.TypeEQ(types.ConnectorTypeK8s),
			).IDs(ctx)
		if err != nil {
			return nil, err
		}

		args := make([]any, len(connIDs))
		for i := range connIDs {
			args[i] = connIDs[i]
		}

		having = sql.In(allocationcost.FieldConnectorID, args...)
	default:
		col := sql.Max(fmt.Sprintf(`CAST((%s) AS varchar)`, groupBySQL))
		pattern := fmt.Sprintf("%%%s%%", query)
		having = sql.Like(col, pattern)
	}
	return having, nil
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

// FilterToSQLPredicates create sql predicate from filters
func FilterToSQLPredicates(filters types.AllocationCostFilters) []*sql.Predicate {
	var or []*sql.Predicate
	for _, cond := range filters {
		var and []*sql.Predicate
		for _, andCond := range cond {
			if ps := ruleToSQLPredicates(andCond); ps != nil {
				and = append(and, ps)
			}
		}

		if len(and) != 0 {
			or = append(or, sql.And(and...))
		}
	}
	return or
}

func ruleToSQLPredicates(cond types.FilterRule) *sql.Predicate {
	if cond.IncludeAll {
		return nil
	}

	toArgs := func(values []string) []any {
		var args []any
		for _, v := range cond.Values {
			args = append(args, v)
		}
		return args
	}

	var pred *sql.Predicate
	// label query
	if strings.HasPrefix(string(cond.FieldName), types.LabelPrefix) {
		labelName := strings.TrimPrefix(string(cond.FieldName), types.LabelPrefix)
		switch cond.Operator {
		case types.OperatorIn:
			pred = sqljson.ValueIn(allocationcost.FieldLabels, toArgs(cond.Values), sqljson.Path(labelName))
		case types.OperatorNotIn:
			pred = sqljson.ValueNotIn(allocationcost.FieldLabels, toArgs(cond.Values), sqljson.Path(labelName))
		}
		return pred
	}

	// other field query
	fieldName := strs.Underscore(string(cond.FieldName))
	switch cond.Operator {
	case types.OperatorIn:
		pred = sql.In(fieldName, toArgs(cond.Values)...)
	case types.OperatorNotIn:
		pred = sql.NotIn(fieldName, toArgs(cond.Values))
	}

	return pred
}
