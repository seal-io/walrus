package distributor

import (
	"context"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/costreport"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
)

type accumulateDistributor struct {
	client model.ClientSet
}

func (r *accumulateDistributor) distribute(
	ctx context.Context,
	startTime,
	endTime time.Time,
	cond types.QueryCondition,
) ([]Resource, int, error) {
	// Item costs.
	itemCosts, count, err := r.itemCosts(ctx, startTime, endTime, cond)
	if err != nil {
		return nil, 0, err
	}

	// Apply shared costs.
	if cond.SharedOptions != nil {
		sharedCosts, err := r.sharedCosts(ctx, startTime, endTime, cond.SharedOptions)
		if err != nil {
			return nil, 0, err
		}

		info, err := r.calculateInfo(ctx, startTime, endTime, cond)
		if err != nil {
			return nil, 0, err
		}

		for i, item := range itemCosts {
			if types.IsIdleOrManagementCost(item.ItemName) {
				continue
			}

			applySharedCost(&itemCosts[i], sharedCosts, info)
		}
	}

	// Apply item display name.
	if err = applyItemDisplayName(ctx, r.client, itemCosts, cond.GroupBy); err != nil {
		return nil, 0, err
	}

	return itemCosts, count, nil
}

func applySharedCost(
	itemCost *Resource,
	sharedCosts *SharedCostConnectors,
	calInfo *CalculateInfo,
) {
	if sharedCosts == nil {
		return
	}

	// Kubernetes cluster resource cost, will not be shared across the clusters.
	applySharedIdleCost(itemCost, sharedCosts.Idle, calInfo)

	applySharedManagementCost(itemCost, sharedCosts.Management, calInfo)

	applySharedItemCost(itemCost, sharedCosts.Items, calInfo)
}

func applySharedIdleCost(
	itemCost *Resource,
	opts *types.IdleShareOption,
	calInfo *CalculateInfo,
) {
	if opts == nil || calInfo == nil {
		return
	}

	for _, connID := range calInfo.ItemConnIDs[itemCost.ItemName] {
		var (
			count      = calInfo.ItemCountPerConn[connID]
			sharedCost = calInfo.CostPerConn[connID].IdleCost

			coef    = 1.0
			perItem = calInfo.ItemCoefficientPerConn[itemCost.ItemName]
		)

		if perItem != nil && perItem[connID] != 0 {
			coef = perItem[connID]
		}

		applySharedStrategy(itemCost, opts.SharingStrategy, sharedCost, count, coef)
	}
}

func applySharedManagementCost(
	itemCost *Resource,
	opts *types.ManagementShareOption,
	calInfo *CalculateInfo,
) {
	if opts == nil || calInfo == nil {
		return
	}

	for _, connID := range calInfo.ItemConnIDs[itemCost.ItemName] {
		var (
			count      = calInfo.ItemCountPerConn[connID]
			sharedCost = calInfo.CostPerConn[connID].ManagementCost

			coef    = 1.0
			perItem = calInfo.ItemCoefficientPerConn[itemCost.ItemName]
		)

		if perItem != nil && perItem[connID] != 0 {
			coef = perItem[connID]
		}

		applySharedStrategy(itemCost, opts.SharingStrategy, sharedCost, count, coef)
	}
}

func applySharedItemCost(
	itemCost *Resource,
	sharedCosts []ItemSharedCost,
	calInfo *CalculateInfo,
) {
	if len(sharedCosts) == 0 || calInfo == nil {
		return
	}

	for _, v := range sharedCosts {
		if v.Option == nil {
			continue
		}

		for _, connID := range calInfo.ItemConnIDs[itemCost.ItemName] {
			if _, ok := v.SharedCosts[connID]; !ok {
				continue
			}

			var (
				count      = calInfo.ItemCountPerConn[connID]
				sharedCost = v.SharedCosts[connID]

				coef    = 1.0
				perItem = calInfo.ItemCoefficientPerConn[itemCost.ItemName]
			)

			if perItem != nil && perItem[connID] != 0 {
				coef = perItem[connID]
			}

			applySharedStrategy(itemCost, v.Option.SharingStrategy, sharedCost, count, coef)
		}
	}
}

func applySharedStrategy(
	itemCost *Resource,
	strategy types.SharingStrategy,
	sharedCost float64,
	count int,
	coef float64,
) {
	var shared float64

	switch strategy {
	case types.SharingStrategyEqually:
		if count != 0 {
			shared = sharedCost / float64(count)
		}
	case types.SharingStrategyProportionally:
		shared = sharedCost * coef
	}

	itemCost.SharedCost += shared
	itemCost.TotalCost += shared
}

func (r *accumulateDistributor) itemCosts(
	ctx context.Context,
	startTime,
	endTime time.Time,
	cond types.QueryCondition,
) ([]Resource, int, error) {
	// Condition.
	_, offset := startTime.Zone()

	orderBy, err := orderBySQL(cond.GroupBy, offset)
	if err != nil {
		return nil, 0, err
	}

	groupBy, err := groupByWithZoneOffsetSQL(cond.GroupBy, offset)
	if err != nil {
		return nil, 0, err
	}

	ps := []*sql.Predicate{
		sql.GTE(costreport.FieldStartTime, startTime),
		sql.LTE(costreport.FieldEndTime, endTime),
	}

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

	countSubQuery := sql.Select(groupBy).
		Where(sql.And(ps...)).
		GroupBy(groupBy).
		From(sql.Table(costreport.Table)).
		As("subQuery")

	// Queried count.
	queriedCount, err := r.client.CostReports().Query().
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

	// Queried items.
	query := r.client.CostReports().Query().
		Modify(func(s *sql.Selector) {
			s.
				Where(sql.And(ps...)).
				SelectExpr(
					sql.Raw(fmt.Sprintf(`%s AS "itemName"`, groupBy)),
					sql.Expr(model.As(model.Sum(costreport.FieldTotalCost), "totalCost")(s)),
					sql.Expr(model.As(model.Sum(costreport.FieldCPUCost), "cpuCost")(s)),
					sql.Expr(model.As(model.Sum(costreport.FieldGPUCost), "gpuCost")(s)),
					sql.Expr(model.As(model.Sum(costreport.FieldRAMCost), "ramCost")(s)),
					sql.Expr(model.As(model.Sum(costreport.FieldPVCost), "pvCost")(s)),
					sql.Expr(model.As(model.Sum(costreport.FieldLoadBalancerCost), "loadBalancerCost")(s)),
				).
				GroupBy(groupBy).
				OrderExpr(sql.Expr(orderBy))

			if havingPs != nil {
				s.Having(havingPs)
			}
		})

	// Paging.
	var (
		page    = cond.Paging.Page
		perPage = cond.Paging.PerPage
	)

	if page != 0 && perPage != 0 {
		query = query.Modify(func(s *sql.Selector) {
			s.Limit(perPage).Offset((page - 1) * perPage)
		})
	}

	var items []Resource
	if err = query.Scan(ctx, &items); err != nil {
		return nil, 0, fmt.Errorf("error query item cost: %w", err)
	}

	// Rename unallocated item name.
	for i, item := range items {
		if item.ItemName == "" {
			items[i].ItemName = types.UnallocatedItemName
			break
		}
	}

	return items, queriedCount, nil
}

func (r *accumulateDistributor) sharedCosts(
	ctx context.Context,
	startTime,
	endTime time.Time,
	opts *types.SharedCostOptions,
) (*SharedCostConnectors, error) {
	if (opts == nil) || opts.Idle == nil && opts.Management == nil && len(opts.Item) == 0 {
		return nil, nil
	}

	scc := &SharedCostConnectors{}

	// Idle costs.
	scc.Idle = opts.Idle

	// Management costs.
	scc.Management = opts.Management

	// Shared item costs.
	if len(opts.Item) != 0 {
		iscs := make([]ItemSharedCost, 0, len(opts.Item))

		for _, v := range opts.Item {
			opt := v

			itemCosts, err := r.customItemCostPerConn(ctx, startTime, endTime, v.Filters)
			if err != nil {
				return nil, err
			}

			iscs = append(iscs, ItemSharedCost{
				Option:      &opt,
				SharedCosts: itemCosts,
			})
		}

		scc.Items = iscs
	}

	return scc, nil
}

func (r *accumulateDistributor) calculateInfo(
	ctx context.Context,
	startTime,
	endTime time.Time,
	cond types.QueryCondition,
) (*CalculateInfo, error) {
	_, offset := startTime.Zone()

	groupBy, err := groupByWithZoneOffsetSQL(cond.GroupBy, offset)
	if err != nil {
		return nil, err
	}

	ps := []*sql.Predicate{
		sql.GTE(costreport.FieldStartTime, startTime),
		sql.LTE(costreport.FieldEndTime, endTime),
	}

	if filterPs := FilterToSQLPredicates(cond.Filters); filterPs != nil {
		ps = append(ps, filterPs)
	}

	// Queried items.
	query := r.client.CostReports().Query().
		Modify(func(s *sql.Selector) {
			s.
				Where(sql.And(ps...)).
				SelectExpr(
					sql.Raw(fmt.Sprintf(`%s AS "itemName"`, groupBy)),
					sql.Expr(model.As(model.Sum(costreport.FieldTotalCost), "totalCost")(s)),
				).
				AppendSelect(
					sql.As(costreport.FieldConnectorID, "connectorID"),
				).
				GroupBy(
					groupBy,
					costreport.FieldConnectorID,
				)
		})

	var items []struct {
		ConnectorID object.ID `json:"connectorID,omitempty"`
		TotalCost   float64   `json:"totalCost,omitempty"`
		ItemName    string    `json:"itemName,omitempty"`
	}

	if err = query.Scan(ctx, &items); err != nil {
		return nil, fmt.Errorf("error query item cost: %w", err)
	}

	// Cost per connector.
	cpc, err := r.costPerConnectors(ctx, startTime, endTime)
	if err != nil {
		return nil, err
	}

	var (
		coepc = make(map[string]map[object.ID]float64)
		ctpc  = make(map[object.ID]int)
		icids = make(map[string][]object.ID)
	)

	for _, v := range items {
		// Name.
		itemName := v.ItemName
		if itemName == "" {
			itemName = types.UnallocatedItemName
		}

		// Count.
		if !types.IsIdleOrManagementCost(itemName) {
			ctpc[v.ConnectorID] += 1
		}

		// Coefficient.
		if _, ok := coepc[itemName]; !ok {
			coepc[itemName] = make(map[object.ID]float64)
		}
		coepc[itemName][v.ConnectorID] = v.TotalCost / cpc[v.ConnectorID].WorkloadCost

		// Item connector ids.
		icids[itemName] = append(icids[itemName], v.ConnectorID)
	}

	return &CalculateInfo{
		ItemCoefficientPerConn: coepc,
		ItemCountPerConn:       ctpc,
		CostPerConn:            cpc,
		ItemConnIDs:            icids,
	}, nil
}

func (r *accumulateDistributor) costPerConnectors(
	ctx context.Context,
	startTime, endTime time.Time,
) (map[object.ID]CostPerConnector, error) {
	connIDs, err := connectorIDs(ctx, r.client)
	if err != nil {
		return nil, err
	}

	idleCosts, err := r.idleCostPerConn(ctx, startTime, endTime, connIDs)
	if err != nil {
		return nil, err
	}

	mgntCosts, err := r.mgntCostPerConn(ctx, startTime, endTime, connIDs)
	if err != nil {
		return nil, err
	}

	totalCosts, err := r.totalCostPerConn(ctx, startTime, endTime, connIDs)
	if err != nil {
		return nil, err
	}

	cpc := make(map[object.ID]CostPerConnector, len(totalCosts))

	for c, v := range totalCosts {
		i := idleCosts[c]
		m := mgntCosts[c]
		w := v - m - i
		cpc[c] = CostPerConnector{
			ConnectorID:    c,
			TotalCost:      v,
			IdleCost:       i,
			ManagementCost: m,
			WorkloadCost:   w,
		}
	}

	return cpc, nil
}

func (r *accumulateDistributor) totalCostPerConn(
	ctx context.Context,
	startTime,
	endTime time.Time,
	connIDs []object.ID,
) (map[object.ID]float64, error) {
	return r.costPerConnQuery(
		ctx,
		startTime,
		endTime,
		nil,
		connIDs,
	)
}

func (r *accumulateDistributor) idleCostPerConn(
	ctx context.Context,
	startTime,
	endTime time.Time,
	connIDs []object.ID,
) (map[object.ID]float64, error) {
	eps := sql.EQ(costreport.FieldName, types.IdleCostItemName)

	return r.costPerConnQuery(
		ctx,
		startTime,
		endTime,
		eps,
		connIDs,
	)
}

func (r *accumulateDistributor) mgntCostPerConn(
	ctx context.Context,
	startTime,
	endTime time.Time,
	connIDs []object.ID,
) (map[object.ID]float64, error) {
	eps := sql.EQ(costreport.FieldName, types.ManagementCostItemName)

	return r.costPerConnQuery(
		ctx,
		startTime,
		endTime,
		eps,
		connIDs,
	)
}

func (r *accumulateDistributor) customItemCostPerConn(
	ctx context.Context,
	startTime,
	endTime time.Time,
	filters types.CostFilters,
) (map[object.ID]float64, error) {
	eps := FilterToSQLPredicates(filters)

	return r.costPerConnQuery(
		ctx,
		startTime,
		endTime,
		eps,
		nil,
	)
}

func (r *accumulateDistributor) costPerConnQuery(
	ctx context.Context,
	startTime,
	endTime time.Time,
	eps *sql.Predicate,
	connIDs []object.ID,
) (map[object.ID]float64, error) {
	ps := []*sql.Predicate{
		sql.GTE(costreport.FieldStartTime, startTime),
		sql.LTE(costreport.FieldEndTime, endTime),
	}

	if len(connIDs) != 0 {
		ids := make([]any, len(connIDs))
		for i := range connIDs {
			ids[i] = connIDs[i]
		}

		ps = append(ps, sql.In(costreport.FieldConnectorID, ids...))
	}

	if eps != nil {
		ps = append(ps, eps)
	}

	var costs []CostPerConnector

	err := r.client.CostReports().Query().
		Modify(func(s *sql.Selector) {
			s.
				Where(sql.And(ps...)).
				SelectExpr(
					sql.Expr(
						model.As(
							model.Sum(costreport.FieldTotalCost), "totalCost",
						)(s),
					),
				).
				AppendSelect(
					sql.As(costreport.FieldConnectorID, "connectorID"),
				).
				GroupBy(costreport.FieldConnectorID)
		}).
		Scan(ctx, &costs)
	if err != nil {
		return nil, fmt.Errorf("error query cost per connector: %w", err)
	}

	costPerConn := make(map[object.ID]float64)
	for _, v := range costs {
		costPerConn[v.ConnectorID] = v.TotalCost
	}

	return costPerConn, nil
}
