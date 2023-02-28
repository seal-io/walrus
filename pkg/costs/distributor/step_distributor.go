package distributor

import (
	"context"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"

	"github.com/seal-io/seal/pkg/apis/cost/view"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/allocationcost"
	"github.com/seal-io/seal/pkg/dao/model/clustercost"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/types"
)

type stepDistributor struct {
	client model.ClientSet
}

func (r *stepDistributor) distribute(ctx context.Context, startTime, endTime time.Time, cond types.QueryCondition) ([]view.Resource, int, error) {
	items, count, err := r.AllocationCosts(ctx, startTime, endTime, cond)
	if err != nil {
		return nil, 0, err
	}

	scb, err := r.SharedCosts(ctx, startTime, endTime, cond.SharedCosts, cond.Step)
	if err != nil {
		return nil, 0, err
	}

	wb, err := r.totalAllocationCost(ctx, startTime, endTime, cond.Step)
	if err != nil {
		return nil, 0, err
	}

	for i, item := range items {
		var (
			bucket   = item.StartTime.String()
			shares   = scb[bucket]
			workload = wb[bucket]
		)

		applySharedCost(count, &items[i].Cost, shares, workload.TotalCost)
	}

	return items, count, nil
}

func (r *stepDistributor) AllocationCosts(ctx context.Context, startTime, endTime time.Time, cond types.QueryCondition) ([]view.Resource, int, error) {
	orderBy, err := cond.GroupBy.OrderBySQL()
	if err != nil {
		return nil, 0, err
	}

	groupBy, err := cond.GroupBy.GroupBySQL()
	if err != nil {
		return nil, 0, err
	}

	dateTrunc, err := cond.Step.DateTruncSQL()
	if err != nil {
		return nil, 0, err
	}

	var (
		groupBys = []string{
			groupBy,
			dateTrunc,
		}
		ps = []*sql.Predicate{
			sql.GTE(allocationcost.FieldStartTime, startTime),
			sql.LTE(allocationcost.FieldEndTime, endTime),
		}
	)

	or := filterToPredicates(cond.Filters)
	if len(or) != 0 {
		ps = append(ps, or...)
	}

	// count
	count, err := r.client.AllocationCosts().Query().
		Modify(func(s *sql.Selector) {
			subQuery := sql.Select(groupBys...).
				Where(
					sql.And(ps...),
				).
				From(sql.Table(allocationcost.Table)).As("subQuery").
				GroupBy(groupBys...)

			s.Count().From(subQuery)
		}).
		Int(ctx)
	if err != nil {
		return nil, 0, err
	}

	// query
	query := r.client.AllocationCosts().Query().
		Modify(func(s *sql.Selector) {
			s.Where(
				sql.And(ps...),
			).SelectExpr(
				sql.Raw(fmt.Sprintf(`%s AS "itemName"`, groupBy)),
				sql.Raw(fmt.Sprintf(`%s AS "startTime"`, dateTrunc)),
				sql.Expr(model.As(model.Sum(allocationcost.FieldTotalCost), "totalCost")(s)),
				sql.Expr(model.As(model.Sum(allocationcost.FieldCpuCost), "cpuCost")(s)),
				sql.Expr(model.As(model.Sum(allocationcost.FieldGpuCost), "gpuCost")(s)),
				sql.Expr(model.As(model.Sum(allocationcost.FieldRamCost), "ramCost")(s)),
				sql.Expr(model.As(model.Sum(allocationcost.FieldPvCost), "pvCost")(s)),
			).GroupBy(
				groupBys...,
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

func (r *stepDistributor) SharedCosts(ctx context.Context, startTime, endTime time.Time, conds types.ShareCosts, step types.Step) (map[string][]*SharedCost, error) {
	if len(conds) == 0 {
		return nil, nil
	}

	dateTrunc, err := step.DateTruncSQL()
	if err != nil {
		return nil, err
	}

	var sharedCosts = make(map[string][]*SharedCost)
	for _, v := range conds {
		saCosts, err := r.sharedAllocationCost(ctx, startTime, endTime, v, dateTrunc)
		if err != nil {
			return nil, err
		}

		idleCosts, err := r.sharedIdleCost(ctx, startTime, endTime, v, dateTrunc)
		if err != nil {
			return nil, err
		}

		mgntCost, err := r.sharedManagementCost(ctx, startTime, endTime, v, dateTrunc)
		if err != nil {
			return nil, err
		}

		bucket := sharedCostBuckets(saCosts, idleCosts, mgntCost, v)
		for key, sc := range bucket {
			sharedCosts[key] = append(sharedCosts[key], sc)
		}
	}
	return sharedCosts, nil
}

func (r *stepDistributor) sharedAllocationCost(ctx context.Context, startTime, endTime time.Time, cond types.SharedCost, dateTrunc string) (map[string]*SharedCost, error) {
	if len(cond.Filters) == 0 {
		return nil, nil
	}

	var filters = []*sql.Predicate{
		sql.GTE(allocationcost.FieldStartTime, startTime),
		sql.LTE(allocationcost.FieldEndTime, endTime),
	}

	var or = filterToPredicates(cond.Filters)
	if len(or) != 0 {
		filters = append(filters, or...)
	}

	var costs []SharedCost
	err := r.client.AllocationCosts().Query().
		Modify(func(s *sql.Selector) {
			s.Where(
				sql.And(filters...),
			)
		}).
		GroupBy(dateTrunc).
		Aggregate(
			model.As(
				model.Sum(allocationcost.FieldTotalCost), "allocationCost",
			),
		).
		Scan(ctx, &costs)
	if err != nil {
		return nil, fmt.Errorf("error query shared allocation cost: %w", err)
	}

	bucket := make(map[string]*SharedCost)
	for _, v := range costs {
		key := v.StartTime.String()
		bucket[key].StartTime = v.StartTime
		bucket[key].AllocationCost += v.AllocationCost
	}

	return bucket, nil
}

func (r *stepDistributor) sharedIdleCost(ctx context.Context, startTime, endTime time.Time, cond types.SharedCost, dateTrunc string) (map[string]*SharedCost, error) {
	if len(cond.IdleCostFilters) == 0 {
		return nil, nil
	}

	var timePs = []*sql.Predicate{
		sql.GTE(allocationcost.FieldStartTime, startTime),
		sql.LTE(allocationcost.FieldEndTime, endTime),
	}

	var ps []*sql.Predicate
	ps = append(ps, timePs...)
	for _, v := range cond.IdleCostFilters {
		if v.ConnectorID.IsNaive() {
			ps = append(ps, sql.EQ("connector_id", v.ConnectorID))
		}
	}

	var costs []SharedCost
	err := r.client.ClusterCosts().Query().
		Modify(func(s *sql.Selector) {
			s.Where(
				sql.And(ps...),
			).SelectExpr(
				sql.Raw(fmt.Sprintf(`%s AS "startTime"`, dateTrunc)),
				sql.Expr(model.As(model.Sum(clustercost.FieldIdleCost), "idleCost")(s)),
			).GroupBy(dateTrunc)
		}).
		Scan(ctx, &costs)
	if err != nil {
		return nil, fmt.Errorf("error query cluster cost: %w", err)
	}

	bucket := make(map[string]*SharedCost)
	for _, v := range costs {
		key := v.StartTime.String()
		if _, ok := bucket[key]; !ok {
			bucket[key] = &SharedCost{}
		}
		bucket[key].StartTime = v.StartTime
		bucket[key].IdleCost += v.IdleCost
	}
	return bucket, nil
}

func (r *stepDistributor) sharedManagementCost(ctx context.Context, startTime, endTime time.Time, cond types.SharedCost, dateTrunc string) (map[string]*SharedCost, error) {
	if len(cond.ManagementCostFilters) == 0 {
		return nil, nil
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

	var costs []SharedCost
	err := r.client.ClusterCosts().Query().
		Modify(func(s *sql.Selector) {
			s.Where(
				sql.And(ps...),
			).SelectExpr(
				sql.Raw(fmt.Sprintf(`%s AS "startTime"`, dateTrunc)),
				sql.Expr(model.As(model.Sum(clustercost.FieldManagementCost), "managementCost")(s)),
			).GroupBy(dateTrunc)
		}).
		Scan(ctx, &costs)
	if err != nil {
		return nil, fmt.Errorf("error query management cost: %w", err)
	}

	bucket := make(map[string]*SharedCost)
	for _, v := range costs {
		key := v.StartTime.String()
		if _, ok := bucket[key]; !ok {
			bucket[key] = &SharedCost{}
		}
		bucket[key].StartTime = v.StartTime
		bucket[key].ManagementCost += v.ManagementCost
	}
	return bucket, nil
}

func (r *stepDistributor) totalAllocationCost(ctx context.Context, startTime, endTime time.Time, step types.Step) (map[string]view.Resource, error) {
	var (
		costs  []view.Resource
		timePs = []predicate.AllocationCost{
			allocationcost.StartTimeGTE(startTime),
			allocationcost.EndTimeLTE(endTime),
		}
	)

	dateTrunc, err := step.DateTruncSQL()
	if err != nil {
		return nil, err
	}

	err = r.client.AllocationCosts().Query().
		Where(timePs...).
		Modify(func(s *sql.Selector) {
			s.SelectExpr(
				sql.Raw(fmt.Sprintf(`%s AS "startTime"`, dateTrunc)),
				sql.Expr(model.As(model.Sum(allocationcost.FieldTotalCost), "totalCost")(s)),
			).GroupBy(dateTrunc)
		}).
		Scan(ctx, &costs)

	if err != nil {
		return nil, fmt.Errorf("error query total allocation cost: %w", err)
	}

	bucket := make(map[string]view.Resource)
	for _, v := range costs {
		key := v.StartTime.String()
		bucket[key] = v
	}
	return bucket, nil
}

func applySharedCost(count int, allocationCost *view.Cost, shares []*SharedCost, workloadCost float64) {
	var coefficients float64
	if allocationCost.TotalCost != 0 && workloadCost != 0 {
		coefficients = allocationCost.TotalCost / workloadCost
	}
	for _, v := range shares {
		var shared float64
		switch v.Condition.SharingStrategy {
		case types.SharingStrategyEqually:
			if count != 0 {
				shared = v.TotalCost / float64(count)
			}
		case types.SharingStrategyProportionally:
			shared = v.TotalCost * coefficients
		}
		allocationCost.SharedCost += shared
		allocationCost.TotalCost += shared
	}
}

func sharedCostBuckets(allocation, idle, management map[string]*SharedCost, cond types.SharedCost) map[string]*SharedCost {
	grouped := make(map[string]*SharedCost)
	for key, v := range allocation {
		if _, ok := grouped[key]; !ok {
			grouped[key] = &SharedCost{}
		}
		grouped[key].TotalCost += v.AllocationCost
		grouped[key].AllocationCost += v.TotalCost
		grouped[key].Condition = cond
	}

	for key, v := range idle {
		if _, ok := grouped[key]; !ok {
			grouped[key] = &SharedCost{}
		}
		grouped[key].TotalCost += v.IdleCost
		grouped[key].IdleCost += v.IdleCost
		grouped[key].Condition = cond
	}

	for key, v := range management {
		if _, ok := grouped[key]; !ok {
			grouped[key] = &SharedCost{}
		}
		grouped[key].TotalCost += v.ManagementCost
		grouped[key].ManagementCost += v.ManagementCost
		grouped[key].Condition = cond
	}

	return grouped
}
