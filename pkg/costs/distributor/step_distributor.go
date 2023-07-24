package distributor

import (
	"context"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/seal/pkg/apis/cost/view"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/allocationcost"
	"github.com/seal-io/seal/pkg/dao/model/clustercost"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/utils/sqlx"
)

type stepDistributor struct {
	client model.ClientSet
}

func (r *stepDistributor) distribute(
	ctx context.Context,
	startTime,
	endTime time.Time,
	cond types.QueryCondition,
) ([]view.Resource, int, error) {
	allocationCosts, queriedCount, err := r.AllocationCosts(ctx, startTime, endTime, cond)
	if err != nil {
		return nil, 0, err
	}

	sharedCosts, err := r.SharedCosts(ctx, startTime, endTime, cond.SharedCosts, cond.Step)
	if err != nil {
		return nil, 0, err
	}

	totalAllocationCosts := r.totalAllocationCosts(allocationCosts)

	if err != nil {
		return nil, 0, err
	}

	itemNameSet := sets.Set[string]{}

	for i, item := range allocationCosts {
		if item.ItemName == "" {
			allocationCosts[i].ItemName = types.UnallocatedLabel
		}

		var (
			bucket              = item.StartTime.Format(time.RFC3339)
			shares              = sharedCosts[bucket]
			totalAllocationCost = totalAllocationCosts[bucket]
		)

		itemNameSet.Insert(allocationCosts[i].ItemName)
		applySharedCost(itemNameSet.Len(), &allocationCosts[i].Cost, shares, totalAllocationCost)
	}

	if err = applyItemDisplayName(ctx, r.client, allocationCosts, cond.GroupBy); err != nil {
		return nil, 0, err
	}

	return allocationCosts, queriedCount, nil
}

func (r *stepDistributor) AllocationCosts(
	ctx context.Context,
	startTime,
	endTime time.Time,
	cond types.QueryCondition,
) ([]view.Resource, int, error) {
	// Condition.
	_, offset := startTime.Zone()

	orderBy, err := orderByWithOffsetSQL(cond.GroupBy, offset)
	if err != nil {
		return nil, 0, err
	}

	groupBy, err := groupByWithZoneOffsetSQL(cond.GroupBy, offset)
	if err != nil {
		return nil, 0, err
	}

	dateTrunc, err := sqlx.DateTruncWithZoneOffsetSQL(
		allocationcost.FieldStartTime,
		string(cond.Step),
		offset,
	)
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

	if filterPs := FilterToSQLPredicates(cond.Filters); filterPs != nil {
		ps = append(ps, filterPs)
	}

	var havingPs *sql.Predicate
	if cond.Query != "" {
		havingPs, err = havingSQL(ctx, r.client, cond.GroupBy, groupBy, cond.Query)
		if err != nil {
			return nil, 0, err
		}
	}

	countSubQuery := sql.Select(groupBys...).
		Where(sql.And(ps...)).
		GroupBy(groupBys...).
		From(sql.Table(allocationcost.Table)).
		As("subQuery")

	// Queried count.
	queriedCount, err := r.client.AllocationCosts().Query().
		Modify(func(s *sql.Selector) {
			if havingPs != nil {
				countSubQuery.Having(havingPs)
			}

			s.Count().From(countSubQuery)
		}).
		Int(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Query.
	query := r.client.AllocationCosts().Query().
		Modify(func(s *sql.Selector) {
			s.Where(
				sql.And(ps...),
			).SelectExpr(
				sql.Raw(fmt.Sprintf(`%s AS "itemName"`, groupBy)),
				sql.Raw(fmt.Sprintf(`%s AS "startTime"`, dateTrunc)),
				sql.Expr(model.As(model.Sum(allocationcost.FieldTotalCost), "totalCost")(s)),
				sql.Expr(model.As(model.Sum(allocationcost.FieldCPUCost), "cpuCost")(s)),
				sql.Expr(model.As(model.Sum(allocationcost.FieldGpuCost), "gpuCost")(s)),
				sql.Expr(model.As(model.Sum(allocationcost.FieldRAMCost), "ramCost")(s)),
				sql.Expr(model.As(model.Sum(allocationcost.FieldPvCost), "pvCost")(s)),
				sql.Expr(model.As(model.Sum(allocationcost.FieldLoadBalancerCost), "loadBalancerCost")(s)),
			).
				GroupBy(
					groupBys...,
				).OrderExpr(
				sql.Expr(orderBy),
			)

			if havingPs != nil {
				s.Having(havingPs)
			}
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

	return items, queriedCount, nil
}

func (r *stepDistributor) SharedCosts(
	ctx context.Context,
	startTime,
	endTime time.Time,
	conds types.ShareCosts,
	step types.Step,
) (map[string][]*SharedCost, error) {
	if len(conds) == 0 {
		return nil, nil
	}

	_, offset := startTime.Zone()

	dateTrunc, err := sqlx.DateTruncWithZoneOffsetSQL(
		allocationcost.FieldStartTime,
		string(step),
		offset,
	)
	if err != nil {
		return nil, err
	}

	sharedCosts := make(map[string][]*SharedCost)

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

func (r *stepDistributor) sharedAllocationCost(
	ctx context.Context,
	startTime,
	endTime time.Time,
	cond types.SharedCost,
	dateTrunc string,
) (map[string]*SharedCost, error) {
	if len(cond.Filters) == 0 {
		return nil, nil
	}

	ps := []*sql.Predicate{
		sql.GTE(allocationcost.FieldStartTime, startTime),
		sql.LTE(allocationcost.FieldEndTime, endTime),
	}

	if filterPs := FilterToSQLPredicates(cond.Filters); filterPs != nil {
		ps = append(ps, filterPs)
	}

	var costs []SharedCost

	err := r.client.AllocationCosts().Query().
		Modify(func(s *sql.Selector) {
			s.Where(
				sql.And(ps...),
			).SelectExpr(
				sql.Raw(fmt.Sprintf(`%s AS "startTime"`, dateTrunc)),
				sql.Expr(model.As(model.Sum(allocationcost.FieldTotalCost), "allocationCost")(s)),
			).GroupBy(dateTrunc)
		}).
		Scan(ctx, &costs)
	if err != nil {
		return nil, fmt.Errorf("error query shared allocation cost: %w", err)
	}

	bucket := make(map[string]*SharedCost)

	for _, v := range costs {
		key := v.StartTime.Format(time.RFC3339)
		if _, ok := bucket[key]; !ok {
			bucket[key] = &SharedCost{}
		}
		bucket[key].StartTime = v.StartTime
		bucket[key].AllocationCost += v.AllocationCost
	}

	return bucket, nil
}

func (r *stepDistributor) sharedIdleCost(
	ctx context.Context,
	startTime,
	endTime time.Time,
	cond types.SharedCost,
	dateTrunc string,
) (map[string]*SharedCost, error) {
	if len(cond.IdleCostFilters) == 0 {
		return nil, nil
	}

	timePs := []*sql.Predicate{
		sql.GTE(allocationcost.FieldStartTime, startTime),
		sql.LTE(allocationcost.FieldEndTime, endTime),
	}

	var ps []*sql.Predicate
	ps = append(ps, timePs...)

	for _, v := range cond.IdleCostFilters {
		if v.ConnectorID.Valid() {
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
		key := v.StartTime.Format(time.RFC3339)
		if _, ok := bucket[key]; !ok {
			bucket[key] = &SharedCost{}
		}
		bucket[key].StartTime = v.StartTime
		bucket[key].IdleCost += v.IdleCost
	}

	return bucket, nil
}

func (r *stepDistributor) sharedManagementCost(
	ctx context.Context,
	startTime,
	endTime time.Time,
	cond types.SharedCost,
	dateTrunc string,
) (map[string]*SharedCost, error) {
	if len(cond.ManagementCostFilters) == 0 {
		return nil, nil
	}

	ps := []*sql.Predicate{
		sql.GTE(allocationcost.FieldStartTime, startTime),
		sql.LTE(allocationcost.FieldEndTime, endTime),
	}

	for _, v := range cond.ManagementCostFilters {
		if v.ConnectorID.Valid() {
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
		key := v.StartTime.Format(time.RFC3339)
		if _, ok := bucket[key]; !ok {
			bucket[key] = &SharedCost{}
		}
		bucket[key].StartTime = v.StartTime
		bucket[key].ManagementCost += v.ManagementCost
	}

	return bucket, nil
}

func (r *stepDistributor) totalAllocationCosts(costs []view.Resource) map[string]float64 {
	bucket := make(map[string]float64)

	for _, v := range costs {
		key := v.StartTime.Format(time.RFC3339)
		bucket[key] += v.TotalCost
	}

	return bucket
}

func applyItemDisplayName(
	ctx context.Context,
	client model.ClientSet,
	items []view.Resource,
	groupBy types.GroupByField,
) error {
	if groupBy != types.GroupByFieldConnectorID {
		return nil
	}

	// Group by connector id.
	conns, err := client.Connectors().Query().
		Where(
			connector.TypeEQ(types.ConnectorTypeK8s),
		).
		Select(
			connector.FieldID,
			connector.FieldName,
		).
		All(ctx)
	if err != nil {
		return err
	}

	for i, v := range items {
		for _, conn := range conns {
			if v.ItemName == conn.ID.String() {
				items[i].ItemName = conn.Name
				break
			}
		}
	}

	return nil
}

func applySharedCost(
	count int,
	allocationCost *view.Cost,
	shares []*SharedCost,
	totalAllocationCost float64,
) {
	var coefficients float64
	if allocationCost.TotalCost != 0 && totalAllocationCost != 0 {
		coefficients = allocationCost.TotalCost / totalAllocationCost
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

func sharedCostBuckets(
	allocation, idle, management map[string]*SharedCost,
	cond types.SharedCost,
) map[string]*SharedCost {
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
