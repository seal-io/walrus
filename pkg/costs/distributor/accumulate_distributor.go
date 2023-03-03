package distributor

import (
	"context"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"

	"github.com/seal-io/seal/pkg/apis/cost/view"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/allocationcost"
	"github.com/seal-io/seal/pkg/dao/model/clustercost"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/utils/strs"
)

type accumulateDistributor struct {
	client model.ClientSet
}

func (r *accumulateDistributor) distribute(ctx context.Context, startTime, endTime time.Time, cond types.QueryCondition) ([]view.Resource, int, error) {
	allocationCosts, count, err := r.allocationResourceCosts(ctx, startTime, endTime, cond)
	if err != nil {
		return nil, 0, err
	}

	sharedCosts, err := r.sharedCosts(ctx, startTime, endTime, cond.SharedCosts)
	if err != nil {
		return nil, 0, err
	}

	workloadCost, err := r.totalAllocationCost(ctx, startTime, endTime)
	if err != nil {
		return nil, 0, err
	}

	for i := range allocationCosts {
		item := &allocationCosts[i]
		if item.ItemName == "" {
			item.ItemName = types.UnallocatedLabel
		}
		applySharedCost(count, &item.Cost, sharedCosts, workloadCost)
	}

	return allocationCosts, count, nil
}

func (r *accumulateDistributor) allocationResourceCosts(ctx context.Context, startTime, endTime time.Time, cond types.QueryCondition) ([]view.Resource, int, error) {
	_, offset := startTime.Zone()
	orderBy, err := orderByWithOffsetSQL(cond.GroupBy, offset)
	if err != nil {
		return nil, 0, err
	}

	groupBy, err := groupByWithZoneOffsetSQL(cond.GroupBy, offset)
	if err != nil {
		return nil, 0, err
	}

	var ps = []*sql.Predicate{
		sql.GTE(allocationcost.FieldStartTime, startTime),
		sql.LTE(allocationcost.FieldEndTime, endTime),
	}

	or := filterToPredicates(cond.Filters)
	if len(or) != 0 {
		ps = append(ps, or...)
	}

	count, err := r.client.AllocationCosts().Query().
		Modify(func(s *sql.Selector) {
			subQuery := sql.Select(groupBy).
				Where(
					sql.And(ps...),
				).
				From(sql.Table(allocationcost.Table)).As("subQuery").
				GroupBy(groupBy)

			s.Count().From(subQuery)
		}).
		Int(ctx)
	if err != nil {
		return nil, 0, err
	}

	query := r.client.AllocationCosts().Query().
		Modify(func(s *sql.Selector) {
			s.Where(
				sql.And(ps...),
			).SelectExpr(
				sql.Raw(fmt.Sprintf(`%s AS "itemName"`, groupBy)),
				sql.Expr(model.As(model.Sum(allocationcost.FieldTotalCost), "totalCost")(s)),
				sql.Expr(model.As(model.Sum(allocationcost.FieldCpuCost), "cpuCost")(s)),
				sql.Expr(model.As(model.Sum(allocationcost.FieldGpuCost), "gpuCost")(s)),
				sql.Expr(model.As(model.Sum(allocationcost.FieldRamCost), "ramCost")(s)),
				sql.Expr(model.As(model.Sum(allocationcost.FieldPvCost), "pvCost")(s)),
			).GroupBy(
				groupBy,
			).OrderExpr(
				sql.Expr(orderBy),
			)
		})

	var (
		page    = cond.Paging.Page
		perPage = cond.Paging.PerPage
	)
	if page != 0 && perPage != 0 {
		query = query.Modify(func(s *sql.Selector) {
			s.Limit(perPage).Offset((page - 1) * perPage)
		})
	}

	var items []view.Resource
	if err = query.Scan(ctx, &items); err != nil {
		return nil, 0, fmt.Errorf("error query allocation cost: %w", err)
	}
	return items, count, nil
}

func (r *accumulateDistributor) sharedCosts(ctx context.Context, startTime, endTime time.Time, conds types.ShareCosts) ([]*SharedCost, error) {
	if len(conds) == 0 {
		return nil, nil
	}

	var sharedCosts []*SharedCost
	for _, v := range conds {
		saCost, err := r.sharedAllocationCost(ctx, startTime, endTime, v)
		if err != nil {
			return nil, err
		}

		idleCost, err := r.sharedIdleCost(ctx, startTime, endTime, v)
		if err != nil {
			return nil, err
		}

		mgntCost, err := r.sharedManagementCost(ctx, startTime, endTime, v)
		if err != nil {
			return nil, err
		}

		sharedCosts = append(sharedCosts, &SharedCost{
			TotalCost:      saCost + idleCost + mgntCost,
			AllocationCost: saCost,
			IdleCost:       idleCost,
			ManagementCost: mgntCost,
			Condition:      v,
		})
	}
	return sharedCosts, nil
}

func (r *accumulateDistributor) sharedAllocationCost(ctx context.Context, startTime, endTime time.Time, cond types.SharedCost) (float64, error) {
	if len(cond.Filters) == 0 {
		return 0, nil
	}

	var filters = []*sql.Predicate{
		sql.GTE(allocationcost.FieldStartTime, startTime),
		sql.LTE(allocationcost.FieldEndTime, endTime),
	}

	var or = filterToPredicates(cond.Filters)
	if len(or) != 0 {
		filters = append(filters, or...)
	}

	cost, err := r.client.AllocationCosts().Query().
		Modify(func(s *sql.Selector) {
			s.Where(
				sql.And(filters...),
			).SelectExpr(
				sql.ExprFunc(func(b *sql.Builder) {
					b.WriteString(fmt.Sprintf(`COALESCE(SUM(%s),0)`, allocationcost.FieldTotalCost))
				}),
			)
		}).
		Float64(ctx)
	if err != nil {
		return 0, fmt.Errorf("error query shared allocation cost: %w", err)
	}

	return cost, nil
}

func (r *accumulateDistributor) sharedIdleCost(ctx context.Context, startTime, endTime time.Time, cond types.SharedCost) (float64, error) {
	if len(cond.IdleCostFilters) == 0 {
		return 0, nil
	}

	var ps = []*sql.Predicate{
		sql.GTE(allocationcost.FieldStartTime, startTime),
		sql.LTE(allocationcost.FieldEndTime, endTime),
	}

	for _, v := range cond.IdleCostFilters {
		if v.ConnectorID.IsNaive() {
			ps = append(ps, sql.EQ(clustercost.FieldConnectorID, v.ConnectorID))
		}
	}

	managementCost, err := r.client.ClusterCosts().Query().
		Modify(func(s *sql.Selector) {
			s.Where(
				sql.And(ps...),
			).SelectExpr(
				sql.ExprFunc(func(b *sql.Builder) {
					b.WriteString(fmt.Sprintf(`COALESCE(SUM(%s),0)`, clustercost.FieldIdleCost))
				}),
			)
		}).
		Float64(ctx)
	if err != nil {
		return 0, fmt.Errorf("error query idle cost: %w", err)
	}

	return managementCost, nil

}

func (r *accumulateDistributor) sharedManagementCost(ctx context.Context, startTime, endTime time.Time, cond types.SharedCost) (float64, error) {
	if len(cond.ManagementCostFilters) == 0 {
		return 0, nil
	}

	var ps = []*sql.Predicate{
		sql.GTE(allocationcost.FieldStartTime, startTime),
		sql.LTE(allocationcost.FieldEndTime, endTime),
	}

	for _, v := range cond.ManagementCostFilters {
		if v.ConnectorID.IsNaive() {
			ps = append(ps, sql.EQ(clustercost.FieldConnectorID, v.ConnectorID))
		}
	}

	managementCost, err := r.client.ClusterCosts().Query().
		Modify(func(s *sql.Selector) {
			s.Where(
				sql.And(ps...),
			).SelectExpr(
				sql.ExprFunc(func(b *sql.Builder) {
					b.WriteString(fmt.Sprintf(`COALESCE(SUM(%s),0)`, clustercost.FieldManagementCost))
				}),
			)
		}).
		Float64(ctx)
	if err != nil {
		return 0, fmt.Errorf("error query management cost: %w", err)
	}

	return managementCost, nil

}

func (r *accumulateDistributor) totalAllocationCost(ctx context.Context, startTime, endTime time.Time) (float64, error) {
	var ps = []*sql.Predicate{
		sql.GTE(allocationcost.FieldStartTime, startTime),
		sql.LTE(allocationcost.FieldEndTime, endTime),
	}
	acCost, err := r.client.AllocationCosts().Query().
		Modify(func(s *sql.Selector) {
			s.Where(
				sql.And(ps...),
			).SelectExpr(
				sql.ExprFunc(func(b *sql.Builder) {
					b.WriteString(fmt.Sprintf(`COALESCE(SUM(%s),0)`, allocationcost.FieldTotalCost))
				}),
			)
		}).
		Float64(ctx)

	if err != nil {
		return 0, fmt.Errorf("error query total allocation cost: %w", err)
	}

	return acCost, nil
}

func filterToPredicates(filters types.AllocationCostFilters) []*sql.Predicate {
	var or []*sql.Predicate
	for _, cond := range filters {
		var and []*sql.Predicate
		for _, andCond := range cond {
			if ps := ruleToPredicates(andCond); ps != nil {
				and = append(and, ps)
			}
		}

		if len(and) != 0 {
			or = append(or, sql.And(and...))
		}
	}
	return or
}

func ruleToPredicates(cond types.FilterRule) *sql.Predicate {
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
