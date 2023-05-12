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
	"github.com/seal-io/seal/pkg/dao/types"
)

type accumulateDistributor struct {
	client model.ClientSet
}

func (r *accumulateDistributor) distribute(
	ctx context.Context,
	startTime,
	endTime time.Time,
	cond types.QueryCondition,
) ([]view.Resource, int, error) {
	allocationCosts, totalCount, queriedCount, err := r.allocationResourceCosts(
		ctx,
		startTime,
		endTime,
		cond,
	)
	if err != nil {
		return nil, 0, err
	}

	sharedCosts, err := r.sharedCosts(ctx, startTime, endTime, cond.SharedCosts)
	if err != nil {
		return nil, 0, err
	}

	totalAllocationCosts := r.totalAllocationCosts(allocationCosts)

	if err != nil {
		return nil, 0, err
	}

	for i := range allocationCosts {
		item := &allocationCosts[i]
		if item.ItemName == "" {
			item.ItemName = types.UnallocatedLabel
		}

		applySharedCost(totalCount, &item.Cost, sharedCosts, totalAllocationCosts)
	}

	if err = applyItemDisplayName(ctx, r.client, allocationCosts, cond.GroupBy); err != nil {
		return nil, 0, err
	}

	return allocationCosts, queriedCount, nil
}

func (r *accumulateDistributor) allocationResourceCosts(
	ctx context.Context,
	startTime,
	endTime time.Time,
	cond types.QueryCondition,
) ([]view.Resource, int, int, error) {
	// Condition.
	_, offset := startTime.Zone()

	orderBy, err := orderByWithOffsetSQL(cond.GroupBy, offset)
	if err != nil {
		return nil, 0, 0, err
	}

	groupBy, err := groupByWithZoneOffsetSQL(cond.GroupBy, offset)
	if err != nil {
		return nil, 0, 0, err
	}

	ps := []*sql.Predicate{
		sql.GTE(allocationcost.FieldStartTime, startTime),
		sql.LTE(allocationcost.FieldEndTime, endTime),
	}

	if filterPs := FilterToSQLPredicates(cond.Filters); filterPs != nil {
		ps = append(ps, filterPs)
	}

	var havingPs *sql.Predicate
	if cond.Query != "" {
		havingPs, err = havingSQL(ctx, r.client, cond.GroupBy, groupBy, cond.Query)
		if err != nil {
			return nil, 0, 0, err
		}
	}

	countSubQuery := sql.Select(groupBy).
		Where(sql.And(ps...)).
		GroupBy(groupBy).
		From(sql.Table(allocationcost.Table)).
		As("subQuery")

	// Total count.
	totalCount, err := r.client.AllocationCosts().Query().
		Modify(func(s *sql.Selector) {
			s.Count().From(countSubQuery)
		}).
		Int(ctx)
	if err != nil {
		return nil, 0, 0, err
	}

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
		return nil, 0, 0, err
	}

	// Queried items.
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
				sql.Expr(model.As(model.Sum(allocationcost.FieldLoadBalancerCost), "loadBalancerCost")(s)),
			).GroupBy(
				groupBy,
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
		return nil, 0, 0, fmt.Errorf("error query allocation cost: %w", err)
	}

	return items, totalCount, queriedCount, nil
}

func (r *accumulateDistributor) sharedCosts(
	ctx context.Context,
	startTime,
	endTime time.Time,
	conds types.ShareCosts,
) ([]*SharedCost, error) {
	if len(conds) == 0 {
		return nil, nil
	}

	sharedCosts := make([]*SharedCost, 0, len(conds))

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

func (r *accumulateDistributor) sharedAllocationCost(
	ctx context.Context,
	startTime,
	endTime time.Time,
	cond types.SharedCost,
) (float64, error) {
	if len(cond.Filters) == 0 {
		return 0, nil
	}

	ps := []*sql.Predicate{
		sql.GTE(allocationcost.FieldStartTime, startTime),
		sql.LTE(allocationcost.FieldEndTime, endTime),
	}

	if filterPs := FilterToSQLPredicates(cond.Filters); filterPs != nil {
		ps = append(ps, filterPs)
	}

	cost, err := r.client.AllocationCosts().Query().
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
		return 0, fmt.Errorf("error query shared allocation cost: %w", err)
	}

	return cost, nil
}

func (r *accumulateDistributor) sharedIdleCost(
	ctx context.Context,
	startTime,
	endTime time.Time,
	cond types.SharedCost,
) (float64, error) {
	if len(cond.IdleCostFilters) == 0 {
		return 0, nil
	}

	ps := []*sql.Predicate{
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

func (r *accumulateDistributor) sharedManagementCost(
	ctx context.Context,
	startTime,
	endTime time.Time,
	cond types.SharedCost,
) (float64, error) {
	if len(cond.ManagementCostFilters) == 0 {
		return 0, nil
	}

	ps := []*sql.Predicate{
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

func (r *accumulateDistributor) totalAllocationCosts(cost []view.Resource) float64 {
	var total float64
	for _, v := range cost {
		total += v.TotalCost
	}

	return total
}
