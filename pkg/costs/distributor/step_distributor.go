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
	"github.com/seal-io/walrus/utils/sqlx"
)

type stepDistributor struct {
	client model.ClientSet
}

func (r *stepDistributor) distribute(
	ctx context.Context,
	startTime,
	endTime time.Time,
	cond types.QueryCondition,
) ([]Resource, int, error) {
	// Item costs.
	itemCosts, count, err := r.itemCost(ctx, startTime, endTime, cond)
	if err != nil {
		return nil, 0, err
	}

	// Apply shared costs.
	if cond.SharedOptions != nil {
		sharedCosts, err := r.sharedCosts(ctx, startTime, endTime, cond.Step, cond.SharedOptions)
		if err != nil {
			return nil, 0, err
		}

		calInfo, err := r.calculateInfo(ctx, startTime, endTime, cond)
		if err != nil {
			return nil, 0, err
		}

		for i, item := range itemCosts {
			if types.IsIdleOrManagementCost(item.ItemName) {
				continue
			}

			var (
				bucket = item.StartTime.Format(time.RFC3339)
				shared = sharedCosts[bucket]
				info   = calInfo[bucket]
			)

			applySharedCost(&itemCosts[i], shared, info)
		}
	}

	if err = applyItemDisplayName(ctx, r.client, itemCosts, cond.GroupBy); err != nil {
		return nil, 0, err
	}

	return itemCosts, count, nil
}

func (r *stepDistributor) itemCost(
	ctx context.Context,
	startTime,
	endTime time.Time,
	cond types.QueryCondition,
) ([]Resource, int, error) {
	// Condition.
	_, offset := startTime.Zone()

	orderBy, err := orderByWithStepSQL(cond.GroupBy, cond.Step, offset)
	if err != nil {
		return nil, 0, err
	}

	groupBy, err := groupByWithZoneOffsetSQL(cond.GroupBy, offset)
	if err != nil {
		return nil, 0, err
	}

	dateTrunc, err := sqlx.DateTruncWithZoneOffsetSQL(
		costreport.FieldStartTime,
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
			sql.GTE(costreport.FieldStartTime, startTime),
			sql.LTE(costreport.FieldEndTime, endTime),
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
		From(sql.Table(costreport.Table)).
		As("subQuery")

	// Queried count.
	count, err := r.client.CostReports().Query().
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
	query := r.client.CostReports().Query().
		Modify(func(s *sql.Selector) {
			s.
				Where(
					sql.And(ps...),
				).
				SelectExpr(
					sql.Raw(fmt.Sprintf(`%s AS "itemName"`, groupBy)),
					sql.Raw(fmt.Sprintf(`%s AS "startTime"`, dateTrunc)),
					sql.Expr(model.As(model.Sum(costreport.FieldTotalCost), "totalCost")(s)),
					sql.Expr(model.As(model.Sum(costreport.FieldCPUCost), "cpuCost")(s)),
					sql.Expr(model.As(model.Sum(costreport.FieldGPUCost), "gpuCost")(s)),
					sql.Expr(model.As(model.Sum(costreport.FieldRAMCost), "ramCost")(s)),
					sql.Expr(model.As(model.Sum(costreport.FieldPVCost), "pvCost")(s)),
					sql.Expr(model.As(model.Sum(costreport.FieldLoadBalancerCost), "loadBalancerCost")(s)),
				).
				GroupBy(
					groupBys...,
				).
				OrderExpr(
					sql.Expr(orderBy),
				)

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
	for i := range items {
		if items[i].ItemName == "" {
			items[i].ItemName = types.UnallocatedItemName
		}
	}

	return items, count, nil
}

func (r *stepDistributor) sharedCosts(
	ctx context.Context,
	startTime,
	endTime time.Time,
	step types.Step,
	opts *types.SharedCostOptions,
) (map[string]*SharedCostConnectors, error) {
	if (opts == nil) || opts.Idle == nil && opts.Management == nil && len(opts.Item) == 0 {
		return nil, nil
	}

	cpc, err := r.costPerConnectorsPerDays(ctx, startTime, endTime, step)
	if err != nil {
		return nil, err
	}

	scc := make(map[string]*SharedCostConnectors)

	for date := range cpc {
		if _, ok := scc[date]; !ok {
			scc[date] = &SharedCostConnectors{}
		}

		// Idle costs.
		scc[date].Idle = opts.Idle

		// Management costs.
		scc[date].Management = opts.Management
	}

	// Shared item costs.
	for _, v := range opts.Item {
		opt := v

		itemCosts, err := r.itemCostPerConnPerDay(ctx, startTime, endTime, step, v.Filters)
		if err != nil {
			return nil, err
		}

		for date, ic := range itemCosts {
			if _, ok := scc[date]; !ok {
				scc[date] = &SharedCostConnectors{}
			}

			scc[date].Items = append(scc[date].Items, ItemSharedCost{
				Option:      &opt,
				SharedCosts: ic,
			})
		}
	}

	return scc, nil
}

func (r *stepDistributor) calculateInfo(
	ctx context.Context,
	startTime,
	endTime time.Time,
	cond types.QueryCondition,
) (map[string]*CalculateInfo, error) {
	// Condition.
	_, offset := startTime.Zone()

	groupBy, err := groupByWithZoneOffsetSQL(cond.GroupBy, offset)
	if err != nil {
		return nil, err
	}

	dateTrunc, err := sqlx.DateTruncWithZoneOffsetSQL(
		costreport.FieldStartTime,
		string(cond.Step),
		offset,
	)
	if err != nil {
		return nil, err
	}

	var (
		groupBys = []string{
			groupBy,
			dateTrunc,
			costreport.FieldConnectorID,
		}
		ps = []*sql.Predicate{
			sql.GTE(costreport.FieldStartTime, startTime),
			sql.LTE(costreport.FieldEndTime, endTime),
		}
	)

	if filterPs := FilterToSQLPredicates(cond.Filters); filterPs != nil {
		ps = append(ps, filterPs)
	}

	// Query.
	query := r.client.CostReports().Query().
		Modify(func(s *sql.Selector) {
			s.
				Where(
					sql.And(ps...),
				).
				SelectExpr(
					sql.Raw(fmt.Sprintf(`%s AS "itemName"`, groupBy)),
					sql.Raw(fmt.Sprintf(`%s AS "startTime"`, dateTrunc)),
					sql.Expr(model.As(model.Sum(costreport.FieldTotalCost), "totalCost")(s)),
				).
				AppendSelect(
					sql.As(costreport.FieldConnectorID, "connectorID"),
				).
				GroupBy(
					groupBys...,
				)
		})

	var items []struct {
		ConnectorID object.ID `json:"connectorID,omitempty"`
		TotalCost   float64   `json:"totalCost,omitempty"`
		ItemName    string    `json:"itemName,omitempty"`
		StartTime   time.Time `json:"startTime,omitempty"`
	}

	if err = query.Scan(ctx, &items); err != nil {
		return nil, fmt.Errorf("error query item cost: %w", err)
	}

	cpc, err := r.costPerConnectorsPerDays(ctx, startTime, endTime, cond.Step)
	if err != nil {
		return nil, err
	}

	info := make(map[string]*CalculateInfo)

	for _, v := range items {
		// Name.
		itemName := v.ItemName
		if itemName == "" {
			itemName = types.UnallocatedItemName
		}

		// Bucket.
		bucket := v.StartTime.Format(time.RFC3339)
		if _, ok := info[bucket]; !ok {
			info[bucket] = &CalculateInfo{
				ItemCountPerConn:       make(map[object.ID]int),
				ItemCoefficientPerConn: make(map[string]map[object.ID]float64),
				ItemConnIDs:            make(map[string][]object.ID),
			}
		}

		// Count.
		if !types.IsIdleOrManagementCost(itemName) {
			info[bucket].ItemCountPerConn[v.ConnectorID] += 1
		}

		// Coefficient.
		if _, ok := info[bucket].ItemCoefficientPerConn[itemName]; !ok {
			info[bucket].ItemCoefficientPerConn[itemName] = make(map[object.ID]float64)
		}

		if cpc[bucket] != nil && cpc[bucket][v.ConnectorID].WorkloadCost != 0 {
			coef := v.TotalCost / cpc[bucket][v.ConnectorID].WorkloadCost
			info[bucket].ItemCoefficientPerConn[itemName][v.ConnectorID] = coef
		}

		// Cost per connector.
		info[bucket].CostPerConn = cpc[bucket]

		// Item connector ids.
		info[bucket].ItemConnIDs[itemName] = append(info[bucket].ItemConnIDs[itemName], v.ConnectorID)
	}

	return info, nil
}

func (r *stepDistributor) costPerConnectorsPerDays(
	ctx context.Context,
	startTime,
	endTime time.Time,
	step types.Step,
) (map[string]map[object.ID]CostPerConnector, error) {
	connIDs, err := connectorIDs(ctx, r.client)
	if err != nil {
		return nil, err
	}

	idleCosts, err := r.idleCostPerConnPerDay(ctx, startTime, endTime, step, connIDs)
	if err != nil {
		return nil, err
	}

	mgntCosts, err := r.mgntCostPerConnPerDay(ctx, startTime, endTime, step, connIDs)
	if err != nil {
		return nil, err
	}

	totalCosts, err := r.totalCostPerConnPerDay(ctx, startTime, endTime, step, connIDs)
	if err != nil {
		return nil, err
	}

	cpc := make(map[string]map[object.ID]CostPerConnector, len(totalCosts))
	for date, cs := range totalCosts {
		if _, ok := cpc[date]; !ok {
			cpc[date] = make(map[object.ID]CostPerConnector)
		}

		for cid, t := range cs {
			i := idleCosts[date][cid]
			m := mgntCosts[date][cid]
			w := t - m - i
			cpc[date][cid] = CostPerConnector{
				ConnectorID:    cid,
				TotalCost:      t,
				IdleCost:       i,
				ManagementCost: m,
				WorkloadCost:   w,
			}
		}
	}

	return cpc, nil
}

func (r *stepDistributor) idleCostPerConnPerDay(
	ctx context.Context,
	startTime,
	endTime time.Time,
	step types.Step,
	connIDs []object.ID,
) (map[string]map[object.ID]float64, error) {
	if len(connIDs) == 0 {
		return nil, nil
	}

	return r.costPerConnPerDayQuery(
		ctx,
		startTime,
		endTime,
		step,
		sql.EQ(costreport.FieldName, types.IdleCostItemName),
		connIDs,
	)
}

func (r *stepDistributor) mgntCostPerConnPerDay(
	ctx context.Context,
	startTime,
	endTime time.Time,
	step types.Step,
	connIDs []object.ID,
) (map[string]map[object.ID]float64, error) {
	if len(connIDs) == 0 {
		return nil, nil
	}

	return r.costPerConnPerDayQuery(
		ctx,
		startTime,
		endTime,
		step,
		sql.EQ(costreport.FieldName, types.ManagementCostItemName),
		connIDs,
	)
}

func (r *stepDistributor) totalCostPerConnPerDay(
	ctx context.Context,
	startTime,
	endTime time.Time,
	step types.Step,
	connIDs []object.ID,
) (map[string]map[object.ID]float64, error) {
	if len(connIDs) == 0 {
		return nil, nil
	}

	return r.costPerConnPerDayQuery(
		ctx,
		startTime,
		endTime,
		step,
		nil,
		connIDs,
	)
}

func (r *stepDistributor) itemCostPerConnPerDay(
	ctx context.Context,
	startTime,
	endTime time.Time,
	step types.Step,
	filters types.CostFilters,
) (map[string]map[object.ID]float64, error) {
	eps := FilterToSQLPredicates(filters)

	return r.costPerConnPerDayQuery(
		ctx,
		startTime,
		endTime,
		step,
		eps,
		nil,
	)
}

func (r *stepDistributor) costPerConnPerDayQuery(
	ctx context.Context,
	startTime,
	endTime time.Time,
	step types.Step,
	eps *sql.Predicate,
	connIDs []object.ID,
) (map[string]map[object.ID]float64, error) {
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

	_, offset := startTime.Zone()

	dateTrunc, err := sqlx.DateTruncWithZoneOffsetSQL(
		costreport.FieldStartTime,
		string(step),
		offset,
	)
	if err != nil {
		return nil, err
	}

	var costs []struct {
		StartTime time.Time `json:"startTime,omitempty"`
		CostPerConnector
	}

	err = r.client.CostReports().Query().
		Modify(func(s *sql.Selector) {
			s.
				Where(
					sql.And(ps...),
				).
				SelectExpr(
					sql.Raw(fmt.Sprintf(`%s AS "startTime"`, dateTrunc)),
					sql.Expr(model.As(model.Sum(costreport.FieldTotalCost), "totalCost")(s)),
				).
				AppendSelect(
					sql.As(costreport.FieldConnectorID, "connectorID"),
				).
				GroupBy(
					costreport.FieldConnectorID,
					dateTrunc,
				)
		}).
		Scan(ctx, &costs)
	if err != nil {
		return nil, fmt.Errorf("error query cost per connector: %w", err)
	}

	bucket := make(map[string]map[object.ID]float64)

	for _, v := range costs {
		key := v.StartTime.Format(time.RFC3339)

		if _, ok := bucket[key]; !ok {
			bucket[key] = make(map[object.ID]float64)
		}
		bucket[key][v.ConnectorID] += v.TotalCost
	}

	return bucket, nil
}
