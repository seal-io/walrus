package cost

import (
	"context"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/costs/distributor"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/costreport"
	"github.com/seal-io/walrus/pkg/dao/model/predicate"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/utils/sqlx"
)

func (h Handler) CollectionRouteGetCostReports(
	req CollectionRouteGetCostReportsRequest,
) (CollectionRouteGetCostReportsResponse, error) {
	items, count, err := h.distributor.Distribute(req.Context, req.StartTime, req.EndTime, req.QueryCondition)
	if err != nil {
		return nil, fmt.Errorf("error query allocation cost: %w", err)
	}

	return runtime.FullPageResponse(items, count), nil
}

func (h Handler) CollectionRouteGetSummaryCosts(
	req CollectionRouteGetSummaryCostsRequest,
) (*CollectionRouteGetSummaryCostsResponse, error) {
	// Total.
	clusterCostPs := []*sql.Predicate{
		sql.GTE(costreport.FieldStartTime, req.StartTime),
		sql.LTE(costreport.FieldEndTime, req.EndTime),
	}

	totalCost, err := h.modelClient.CostReports().Query().
		Modify(func(s *sql.Selector) {
			s.Where(
				sql.And(clusterCostPs...),
			).SelectExpr(
				sql.ExprFunc(func(b *sql.Builder) {
					b.WriteString(fmt.Sprintf(`COALESCE(SUM(%s),0)`, costreport.FieldTotalCost))
				}),
			)
		}).
		Float64(req.Context)
	if err != nil {
		return nil, fmt.Errorf("error summary total cost: %w", err)
	}

	// Cluster.
	clusters, err := h.modelClient.CostReports().Query().
		Modify(func(s *sql.Selector) {
			s.Where(
				sql.And(clusterCostPs...),
			).Select(
				costreport.FieldConnectorID,
			).Distinct()
		}).Strings(req.Context)
	if err != nil {
		return nil, fmt.Errorf("error summary each cluster's cost: %w", err)
	}

	// Project.
	var (
		projectCostPs = []*sql.Predicate{
			sql.GTE(costreport.FieldStartTime, req.StartTime),
			sql.LTE(costreport.FieldEndTime, req.EndTime),
			sqljson.ValueIsNotNull(costreport.FieldLabels),
		}
		projects []struct {
			Value string `json:"value"`
		}
	)

	err = h.modelClient.CostReports().Query().
		Modify(func(s *sql.Selector) {
			s.Where(
				sql.And(projectCostPs...),
			).SelectExpr(
				sql.Expr(fmt.Sprintf(`DISTINCT(labels ->> '%s') AS value`, types.LabelWalrusProjectName)),
			)
		}).Scan(req.Context, &projects)
	if err != nil {
		return nil, fmt.Errorf("error summary each project's cost: %w", err)
	}

	var projectCount int

	for _, v := range projects {
		if v.Value != "" {
			projectCount += 1
		}
	}

	// Days.
	_, offset := req.StartTime.Zone()

	days, err := h.costReportExistedDays(req.Context, clusterCostPs, offset)
	if err != nil {
		return nil, err
	}

	// Average.
	averageDailyCost := averageDaily(days, totalCost)

	// Collected time range.
	var timeRange *CollectedTimeRange
	if totalCost != 0 {
		timeRange, err = h.costCollectedDataRange(req.Context, req.StartTime.Location())
		if err != nil {
			return nil, err
		}
	}

	return &CollectionRouteGetSummaryCostsResponse{
		TotalCost:          totalCost,
		AverageDailyCost:   averageDailyCost,
		ClusterCount:       len(clusters),
		ProjectCount:       projectCount,
		CollectedTimeRange: timeRange,
	}, nil
}

func (h Handler) CollectionRouteGetSummaryClusterCosts(
	req CollectionRouteGetSummaryClusterCostsRequest,
) (*CollectionRouteGetSummaryClusterCostsResponse, error) {
	ps := []predicate.CostReport{
		costreport.StartTimeGTE(req.StartTime),
		costreport.EndTimeLTE(req.EndTime),
		costreport.ConnectorID(req.ConnectorID),
	}

	// Total cost, use modify instead of aggregate to handle when result is null, can't aggregate in null.
	totalCost, err := h.modelClient.CostReports().Query().
		Where(ps...).
		Modify(func(s *sql.Selector) {
			s.SelectExpr(
				sql.ExprFunc(func(b *sql.Builder) {
					b.WriteString(fmt.Sprintf(`COALESCE(SUM(%s),0)`, costreport.FieldTotalCost))
				}),
			)
		}).
		Float64(req.Context)
	if err != nil {
		return nil, fmt.Errorf("error summary cluster cost: %w", err)
	}

	// Management cost.
	managementCost, err := h.modelClient.CostReports().Query().
		Where(
			costreport.Name(types.ManagementCostItemName),
			costreport.And(ps...),
		).
		Modify(func(s *sql.Selector) {
			s.SelectExpr(
				sql.ExprFunc(func(b *sql.Builder) {
					b.WriteString(fmt.Sprintf(`COALESCE(SUM(%s),0)`, costreport.FieldTotalCost))
				}),
			)
		}).
		Float64(req.Context)
	if err != nil {
		return nil, fmt.Errorf("error summary cluster management cost: %w", err)
	}

	// Idle cost.
	idleCost, err := h.modelClient.CostReports().Query().
		Where(
			costreport.Name(types.IdleCostItemName),
			costreport.And(ps...),
		).
		Modify(func(s *sql.Selector) {
			s.SelectExpr(
				sql.ExprFunc(func(b *sql.Builder) {
					b.WriteString(fmt.Sprintf(`COALESCE(SUM(%s),0)`, costreport.FieldTotalCost))
				}),
			)
		}).
		Float64(req.Context)
	if err != nil {
		return nil, fmt.Errorf("error summary cluster idle cost: %w", err)
	}

	// Days.
	_, offset := req.StartTime.Zone()
	sps := []*sql.Predicate{
		sql.GTE(costreport.FieldStartTime, req.StartTime),
		sql.LTE(costreport.FieldEndTime, req.EndTime),
		sql.EQ(costreport.FieldConnectorID, req.ConnectorID),
	}

	days, err := h.costReportExistedDays(req.Context, sps, offset)
	if err != nil {
		return nil, err
	}

	// Collected time range.
	var timeRange *CollectedTimeRange
	if totalCost != 0 {
		timeRange, err = h.costCollectedDataRange(
			req.Context,
			req.StartTime.Location(),
			sql.EQ(costreport.FieldConnectorID, req.ConnectorID))
		if err != nil {
			return nil, err
		}
	}

	return &CollectionRouteGetSummaryClusterCostsResponse{
		TotalCost:          totalCost,
		ManagementCost:     managementCost,
		IdleCost:           idleCost,
		ItemCost:           totalCost - managementCost - idleCost,
		CollectedTimeRange: timeRange,
		AverageDailyCost:   averageDaily(days, totalCost),
	}, nil
}

func (h Handler) CollectionRouteGetSummaryProjectCosts(
	req CollectionRouteGetSummaryProjectCostsRequest,
) (*CollectionRouteGetSummaryProjectCostsResponse, error) {
	ps := []*sql.Predicate{
		sql.GTE(costreport.FieldStartTime, req.StartTime),
		sql.LTE(costreport.FieldEndTime, req.EndTime),
		sqljson.ValueEQ(costreport.FieldLabels, req.Project, sqljson.Path(types.LabelWalrusProjectName)),
	}

	totalCost, err := h.modelClient.CostReports().Query().
		Modify(func(s *sql.Selector) {
			s.Where(
				sql.And(ps...),
			).SelectExpr(
				sql.ExprFunc(func(b *sql.Builder) {
					b.WriteString(fmt.Sprintf(`COALESCE(SUM(%s),0)`, costreport.FieldTotalCost))
				}),
			)
		}).
		Float64(req.Context)
	if err != nil {
		return nil, fmt.Errorf("error summary project cost: %w", err)
	}

	if totalCost == 0 {
		return nil, nil
	}

	// Days.
	_, offset := req.StartTime.Zone()

	days, err := h.costReportExistedDays(req.Context, ps, offset)
	if err != nil {
		return nil, err
	}

	// Collected time range.
	var timeRange *CollectedTimeRange
	if totalCost != 0 {
		timeRange, err = h.itemCostCollectedDataRange(
			req.Context,
			req.StartTime.Location(),
			sqljson.ValueEQ(costreport.FieldLabels, req.Project, sqljson.Path(types.LabelWalrusProjectName)))
		if err != nil {
			return nil, err
		}
	}

	return &CollectionRouteGetSummaryProjectCostsResponse{
		TotalCost:          totalCost,
		CollectedTimeRange: timeRange,
		AverageDailyCost:   averageDaily(days, totalCost),
	}, nil
}

func (h Handler) CollectionRouteGetSummaryQueriedCosts(
	req CollectionRouteGetSummaryQueriedCostsRequest,
) (*CollectionRouteGetSummaryQueriedCostsResponse, error) {
	cond := types.QueryCondition{
		Filters:       req.Filters,
		GroupBy:       types.GroupByFieldConnectorID,
		SharedOptions: req.SharedOptions,
	}

	items, count, err := h.distributor.Distribute(req.Context, req.StartTime, req.EndTime, cond)
	if err != nil {
		return nil, fmt.Errorf("error query item cost: %w", err)
	}

	// Days.
	ps := []*sql.Predicate{
		sql.GTE(costreport.FieldStartTime, req.StartTime),
		sql.LTE(costreport.FieldEndTime, req.EndTime),
	}

	if filterPs := distributor.FilterToSQLPredicates(cond.Filters); filterPs != nil {
		ps = append(ps, distributor.FilterToSQLPredicates(cond.Filters))
	}

	_, offset := req.StartTime.Zone()

	days, err := h.costReportExistedDays(req.Context, ps, offset)
	if err != nil {
		return nil, err
	}

	// Summary.
	summary := &CollectionRouteGetSummaryQueriedCostsResponse{}
	for _, v := range items {
		summary.TotalCost += v.TotalCost
		summary.SharedCost += v.SharedCost
		summary.CPUCost += v.CPUCost
		summary.GPUCost += v.GPUCost
		summary.RAMCost += v.RAMCost
		summary.PVCost += v.PVCost
	}

	// Collected time range.
	var timeRange *CollectedTimeRange
	if summary.TotalCost != 0 {
		timeRange, err = h.itemCostCollectedDataRange(req.Context,
			req.StartTime.Location(),
			distributor.FilterToSQLPredicates(cond.Filters))
		if err != nil {
			return nil, err
		}
	}

	summary.CollectedTimeRange = timeRange
	summary.AverageDailyCost = averageDaily(days, summary.TotalCost)
	summary.ConnectorCount = count

	return summary, nil
}

func (h Handler) costReportExistedDays(ctx context.Context, ps []*sql.Predicate, offset int) (int, error) {
	groupBy, err := sqlx.DateTruncWithZoneOffsetSQL(costreport.FieldStartTime, string(types.StepDay), offset)
	if err != nil {
		return 0, err
	}

	days, err := h.modelClient.CostReports().Query().
		Modify(func(s *sql.Selector) {
			subQuery := sql.Select(groupBy).
				Where(
					sql.And(ps...),
				).
				From(sql.Table(costreport.Table)).As("subQuery").
				GroupBy(groupBy)

			s.Count().From(subQuery)
		}).
		Int(ctx)
	if err != nil {
		return 0, fmt.Errorf("error get cost report time range: %w", err)
	}

	return days, nil
}

func (h Handler) costCollectedDataRange(
	ctx context.Context,
	loc *time.Location,
	ps ...*sql.Predicate,
) (*CollectedTimeRange, error) {
	modifier := func(s *sql.Selector) {
		s.Select(
			costreport.FieldStartTime,
			costreport.FieldEndTime,
		)

		if len(ps) != 0 {
			s.Where(
				sql.And(ps...),
			)
		}
	}

	// First.
	first, err := h.modelClient.CostReports().Query().
		Modify(modifier).
		Order(model.Asc(costreport.FieldStartTime)).
		First(ctx)
	if err != nil {
		return nil, fmt.Errorf("error get first collected cost data: %w", err)
	}

	// Last.
	last, err := h.modelClient.CostReports().Query().
		Modify(modifier).
		Order(model.Desc(costreport.FieldStartTime)).
		First(ctx)
	if err != nil {
		return nil, fmt.Errorf("error get last collected cost data: %w", err)
	}

	return &CollectedTimeRange{
		FirstTime: first.StartTime.In(loc),
		LastTime:  last.EndTime.In(loc),
	}, nil
}

func (h Handler) itemCostCollectedDataRange(
	ctx context.Context,
	loc *time.Location,
	ps ...*sql.Predicate,
) (*CollectedTimeRange, error) {
	modifier := func(s *sql.Selector) {
		s.Select(
			costreport.FieldStartTime,
			costreport.FieldEndTime,
		)

		if len(ps) != 0 {
			s.Where(
				sql.And(ps...),
			)
		}
	}

	// First.
	first, err := h.modelClient.CostReports().Query().
		Modify(modifier).
		Order(model.Asc(costreport.FieldStartTime)).
		First(ctx)
	if err != nil {
		return nil, fmt.Errorf("error get first collected cost data: %w", err)
	}

	// Last.
	last, err := h.modelClient.CostReports().Query().
		Modify(modifier).
		Order(model.Desc(costreport.FieldStartTime)).
		First(ctx)
	if err != nil {
		return nil, fmt.Errorf("error get last collected cost data: %w", err)
	}

	return &CollectedTimeRange{
		FirstTime: first.StartTime.In(loc),
		LastTime:  last.EndTime.In(loc),
	}, nil
}

func averageDaily(days int, total float64) float64 {
	if total == 0 || days == 0 {
		return 0
	}

	// Average.
	return total / float64(days)
}
