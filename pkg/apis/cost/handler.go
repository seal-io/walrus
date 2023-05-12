package cost

import (
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/cost/view"
	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/costs/distributor"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/allocationcost"
	"github.com/seal-io/seal/pkg/dao/model/clustercost"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/utils/sqlx"
)

func Handle(mc model.ClientSet) Handler {
	return Handler{
		modelClient: mc,
		distributor: distributor.New(mc),
	}
}

type Handler struct {
	modelClient model.ClientSet
	distributor *distributor.Distributor
}

func (h Handler) Kind() string {
	return "Cost"
}

func (h Handler) Validating() any {
	return h.modelClient
}

// Basic APIs.

// Batch APIs.

// Extensional APIs.

func (h Handler) CollectionRouteAllocationCost(ctx *gin.Context, req view.AllocationCostRequest) (*runtime.ResponseCollection, error) {
	items, count, err := h.distributor.Distribute(ctx, req.StartTime, req.EndTime, req.QueryCondition)
	if err != nil {
		return nil, runtime.Errorw(err, "error query allocation cost")
	}

	resp := runtime.GetResponseCollection(ctx, items, count)
	return &resp, nil
}

func (h Handler) CollectionRouteSummaryCost(ctx *gin.Context, req view.SummaryCostRequest) (*view.SummaryCostResponse, error) {
	// Total.
	var clusterCostPs = []*sql.Predicate{
		sql.GTE(allocationcost.FieldStartTime, req.StartTime),
		sql.LTE(allocationcost.FieldEndTime, req.EndTime),
	}

	totalCost, err := h.modelClient.ClusterCosts().Query().
		Modify(func(s *sql.Selector) {
			s.Where(
				sql.And(clusterCostPs...),
			).SelectExpr(
				sql.ExprFunc(func(b *sql.Builder) {
					b.WriteString(fmt.Sprintf(`COALESCE(SUM(%s),0)`, clustercost.FieldTotalCost))
				}),
			)
		}).
		Float64(ctx)
	if err != nil {
		return nil, fmt.Errorf("error summary total cost: %w", err)
	}

	// Cluster.
	clusters, err := h.modelClient.ClusterCosts().Query().
		Modify(func(s *sql.Selector) {
			s.Where(
				sql.And(clusterCostPs...),
			).Select(
				clustercost.FieldConnectorID,
			).Distinct()
		}).Strings(ctx)
	if err != nil {
		return nil, fmt.Errorf("error summary each cluster's cost: %w", err)
	}

	// Project.
	var (
		projectCostPs = []*sql.Predicate{
			sql.GTE(allocationcost.FieldStartTime, req.StartTime),
			sql.LTE(allocationcost.FieldEndTime, req.EndTime),
			sqljson.ValueIsNotNull(allocationcost.FieldLabels),
		}
		projects []struct {
			Value string `json:"value"`
		}
	)
	err = h.modelClient.AllocationCosts().Query().
		Modify(func(s *sql.Selector) {
			s.Where(
				sql.And(projectCostPs...),
			).SelectExpr(
				sql.Expr(fmt.Sprintf(`DISTINCT(labels ->> '%s') AS value`, types.LabelSealProject)),
			)
		}).Scan(ctx, &projects)
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
	days, err := h.clusterCostExistedDays(ctx, clusterCostPs, offset)
	if err != nil {
		return nil, err
	}

	// Average.
	averageDailyCost := averageDaily(days, totalCost)

	// Collected time range.
	var timeRange *view.CollectedTimeRange
	if totalCost != 0 {
		timeRange, err = h.clusterCostCollectedDataRange(ctx, req.StartTime.Location())
		if err != nil {
			return nil, err
		}
	}

	return &view.SummaryCostResponse{
		TotalCost:          totalCost,
		AverageDailyCost:   averageDailyCost,
		ClusterCount:       len(clusters),
		ProjectCount:       projectCount,
		CollectedTimeRange: timeRange,
	}, nil
}

func (h Handler) CollectionRouteSummaryClusterCost(ctx *gin.Context, req view.SummaryClusterCostRequest) (*view.SummaryClusterCostResponse, error) {
	var ps = []*sql.Predicate{
		sql.GTE(allocationcost.FieldStartTime, req.StartTime),
		sql.LTE(allocationcost.FieldEndTime, req.EndTime),
		sql.EQ(allocationcost.FieldConnectorID, req.ConnectorID),
	}

	var s []view.SummaryClusterCostResponse
	err := h.modelClient.ClusterCosts().Query().
		Modify(func(s *sql.Selector) {
			s.Where(
				sql.And(ps...),
			).SelectExpr(
				sql.Expr(model.As(model.Sum(clustercost.FieldTotalCost), "totalCost")(s)),
				sql.Expr(model.As(model.Sum(clustercost.FieldManagementCost), "managementCost")(s)),
				sql.Expr(model.As(model.Sum(clustercost.FieldIdleCost), "idleCost")(s)),
				sql.Expr(model.As(model.Sum(clustercost.FieldAllocationCost), "allocationCost")(s)),
			)
		}).
		Scan(ctx, &s)
	if err != nil {
		return nil, fmt.Errorf("error summary cluster cost: %w", err)
	}

	if len(s) == 0 {
		return nil, nil
	}

	// Days.
	_, offset := req.StartTime.Zone()
	days, err := h.clusterCostExistedDays(ctx, ps, offset)
	if err != nil {
		return nil, err
	}

	// Collected time range.
	var timeRange *view.CollectedTimeRange
	if s[0].TotalCost != 0 {
		timeRange, err = h.clusterCostCollectedDataRange(ctx,
			req.StartTime.Location(),
			sql.EQ(clustercost.FieldConnectorID, req.ConnectorID))
		if err != nil {
			return nil, err
		}
	}

	summary := s[0]
	summary.CollectedTimeRange = timeRange
	summary.AverageDailyCost = averageDaily(days, summary.TotalCost)
	return &summary, nil
}

func (h Handler) CollectionRouteSummaryProjectCost(ctx *gin.Context, req view.SummaryProjectCostRequest) (*view.SummaryCostCommonResponse, error) {
	var ps = []*sql.Predicate{
		sql.GTE(allocationcost.FieldStartTime, req.StartTime),
		sql.LTE(allocationcost.FieldEndTime, req.EndTime),
		sqljson.ValueEQ(allocationcost.FieldLabels, req.Project, sqljson.Path(types.LabelSealProject)),
	}

	var s []view.SummaryCostCommonResponse
	err := h.modelClient.AllocationCosts().Query().
		Modify(func(s *sql.Selector) {
			s.Where(
				sql.And(ps...),
			).SelectExpr(
				sql.Expr(model.As(model.Sum(clustercost.FieldTotalCost), "totalCost")(s)),
			)
		}).
		Scan(ctx, &s)
	if err != nil {
		return nil, fmt.Errorf("error summary project cost: %w", err)
	}

	if len(s) == 0 {
		return nil, nil
	}

	// Days.
	_, offset := req.StartTime.Zone()
	days, err := h.allocationCostExistedDays(ctx, ps, offset)
	if err != nil {
		return nil, err
	}

	// Collected time range.
	var timeRange *view.CollectedTimeRange
	if s[0].TotalCost != 0 {
		timeRange, err = h.allocationCostCollectedDataRange(ctx,
			req.StartTime.Location(),
			sqljson.ValueEQ(allocationcost.FieldLabels, req.Project, sqljson.Path(types.LabelSealProject)))
		if err != nil {
			return nil, err
		}
	}

	summary := s[0]
	summary.CollectedTimeRange = timeRange
	summary.AverageDailyCost = averageDaily(days, s[0].TotalCost)
	return &summary, nil
}

func (h Handler) CollectionRouteSummaryQueriedCost(ctx *gin.Context, req view.SummaryQueriedCostRequest) (*view.SummaryQueriedCostResponse, error) {
	cond := types.QueryCondition{
		Filters:     req.Filters,
		GroupBy:     types.GroupByFieldConnectorID,
		SharedCosts: req.SharedCosts,
	}

	items, count, err := h.distributor.Distribute(ctx, req.StartTime, req.EndTime, cond)
	if err != nil {
		return nil, runtime.Errorw(err, "error query allocation cost")
	}

	// Days.
	var ps = []*sql.Predicate{
		sql.GTE(allocationcost.FieldStartTime, req.StartTime),
		sql.LTE(allocationcost.FieldEndTime, req.EndTime),
	}

	if filterPs := distributor.FilterToSQLPredicates(cond.Filters); filterPs != nil {
		ps = append(ps, distributor.FilterToSQLPredicates(cond.Filters))
	}

	_, offset := req.StartTime.Zone()
	days, err := h.allocationCostExistedDays(ctx, ps, offset)
	if err != nil {
		return nil, err
	}

	// Summary.
	summary := &view.SummaryQueriedCostResponse{}
	for _, v := range items {
		summary.TotalCost += v.TotalCost
		summary.SharedCost += v.SharedCost
		summary.CpuCost += v.CpuCost
		summary.GpuCost += v.GpuCost
		summary.RamCost += v.RamCost
		summary.PvCost += v.PvCost
	}

	// Collected time range.
	var timeRange *view.CollectedTimeRange
	if summary.TotalCost != 0 {
		timeRange, err = h.allocationCostCollectedDataRange(ctx,
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

func (h Handler) clusterCostExistedDays(ctx *gin.Context, ps []*sql.Predicate, offset int) (int, error) {
	groupBy, err := sqlx.DateTruncWithZoneOffsetSQL(clustercost.FieldStartTime, string(types.StepDay), offset)
	if err != nil {
		return 0, err
	}

	days, err := h.modelClient.ClusterCosts().Query().
		Modify(func(s *sql.Selector) {
			subQuery := sql.Select(groupBy).
				Where(
					sql.And(ps...),
				).
				From(sql.Table(clustercost.Table)).As("subQuery").
				GroupBy(groupBy)

			s.Count().From(subQuery)
		}).
		Int(ctx)
	if err != nil {
		return 0, fmt.Errorf("error get cluster cost time range: %w", err)
	}

	return days, nil
}

func (h Handler) allocationCostExistedDays(ctx *gin.Context, ps []*sql.Predicate, offset int) (int, error) {
	groupBy, err := sqlx.DateTruncWithZoneOffsetSQL(allocationcost.FieldStartTime, string(types.StepDay), offset)
	if err != nil {
		return 0, err
	}

	days, err := h.modelClient.AllocationCosts().Query().
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
		return 0, fmt.Errorf("error get allocation cost time range: %w", err)
	}

	return days, nil
}

func (h Handler) clusterCostCollectedDataRange(ctx *gin.Context, loc *time.Location, ps ...*sql.Predicate) (*view.CollectedTimeRange, error) {
	modifier := func(s *sql.Selector) {
		s.Select(
			clustercost.FieldStartTime,
			clustercost.FieldEndTime,
		)

		if len(ps) != 0 {
			s.Where(
				sql.And(ps...),
			)
		}
	}

	// First.
	first, err := h.modelClient.ClusterCosts().Query().
		Modify(modifier).
		Order(model.Asc(clustercost.FieldStartTime)).
		First(ctx)
	if err != nil {
		return nil, fmt.Errorf("error get first collected cost data: %w", err)
	}

	// Last.
	last, err := h.modelClient.ClusterCosts().Query().
		Modify(modifier).
		Order(model.Desc(clustercost.FieldStartTime)).
		First(ctx)
	if err != nil {
		return nil, fmt.Errorf("error get last collected cost data: %w", err)
	}

	return &view.CollectedTimeRange{
		FirstTime: first.StartTime.In(loc),
		LastTime:  last.EndTime.In(loc),
	}, nil
}

func (h Handler) allocationCostCollectedDataRange(ctx *gin.Context, loc *time.Location, ps ...*sql.Predicate) (*view.CollectedTimeRange, error) {
	modifier := func(s *sql.Selector) {
		s.Select(
			clustercost.FieldStartTime,
			clustercost.FieldEndTime,
		)

		if len(ps) != 0 {
			s.Where(
				sql.And(ps...),
			)
		}
	}

	// First.
	first, err := h.modelClient.AllocationCosts().Query().
		Modify(modifier).
		Order(model.Asc(clustercost.FieldStartTime)).
		First(ctx)
	if err != nil {
		return nil, fmt.Errorf("error get first collected cost data: %w", err)
	}

	// Last.
	last, err := h.modelClient.AllocationCosts().Query().
		Modify(modifier).
		Order(model.Desc(clustercost.FieldStartTime)).
		First(ctx)
	if err != nil {
		return nil, fmt.Errorf("error get last collected cost data: %w", err)
	}

	return &view.CollectedTimeRange{
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
