package dashboard

import (
	"context"
	"sort"
	"time"

	"entgo.io/ent/dialect/sql"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/environment"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/model/resourcecomponent"
	"github.com/seal-io/walrus/pkg/dao/model/resourcerun"
	"github.com/seal-io/walrus/pkg/dao/schema/intercept"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/utils/sqlx"
	"github.com/seal-io/walrus/utils/timex"
)

var getResourceRunFields = resourcerun.WithoutFields(
	resourcerun.FieldRecord,
	resourcerun.FieldInputConfigs,
	resourcerun.FieldTemplateName,
	resourcerun.FieldTemplateVersion,
	resourcerun.FieldAttributes,
	resourcerun.FieldComputedAttributes,
	resourcerun.FieldVariables,
)

const summaryStatus = "(status ->> 'summaryStatus')"

// CollectionRouteGetLatestResourceRuns returns the latest 10 runs of resources.
func (h Handler) CollectionRouteGetLatestResourceRuns(
	req CollectionRouteGetLatestResourceRunsRequest,
) (CollectionRouteGetLatestResourceRunsResponse, int, error) {
	ctx := intercept.WithProjectInterceptor(req.Context)

	query := h.modelClient.ResourceRuns().Query().
		Order(model.Desc(resourcerun.FieldCreateTime)).
		Select(getResourceRunFields...).
		WithProject(func(pq *model.ProjectQuery) {
			pq.Select(
				resource.FieldID,
				resource.FieldName,
			)
		}).
		WithResource(func(sq *model.ResourceQuery) {
			sq.Select(
				resource.FieldID,
				resource.FieldName,
			)
		}).
		WithEnvironment(
			func(eq *model.EnvironmentQuery) {
				eq.Select(
					environment.FieldID,
					environment.FieldName)
			}).
		Limit(10)

	entities, err := query.All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeResourceRuns(entities), len(entities), nil
}

func (h Handler) CollectionRouteGetBasicInformation(
	req CollectionRouteGetBasicInformationRequest,
) (*CollectionRouteGetBasicInformationResponse, error) {
	ctx := intercept.WithProjectInterceptor(req.Context)

	// Count owned projects.
	projectNum, err := h.modelClient.Projects().Query().
		Count(ctx)
	if err != nil {
		return nil, err
	}

	// Count environments below owned projects.
	environmentNum, err := h.modelClient.Environments().Query().
		Count(ctx)
	if err != nil {
		return nil, err
	}

	// Count connectors below owned projects and global.
	connectorNum, err := h.modelClient.Connectors().Query().
		Count(ctx)
	if err != nil {
		return nil, err
	}

	// Count resources below owned projects.
	resourceNum, err := h.modelClient.Resources().Query().
		Count(ctx)
	if err != nil {
		return nil, err
	}

	// Count resource components below owned projects if needed.
	var resourceComponentNum int
	if req.WithResourceComponent {
		resourceComponentNum, err = h.modelClient.ResourceComponents().Query().
			Where(resourcecomponent.ModeNEQ(types.ResourceComponentModeData)).
			Count(ctx)
		if err != nil {
			return nil, err
		}
	}

	// Count resource runs below owned projects if needed.
	var resourceRunNum int
	if req.WithResourceRun {
		resourceRunNum, err = h.modelClient.ResourceRuns().Query().
			Count(ctx)
		if err != nil {
			return nil, err
		}
	}

	return &CollectionRouteGetBasicInformationResponse{
		Project:           projectNum,
		Environment:       environmentNum,
		Connector:         connectorNum,
		Resource:          resourceNum,
		ResourceComponent: resourceComponentNum,
		ResourceRun:       resourceRunNum,
	}, nil
}

// CollectionRouteGetResourceRunStatistics returns statistics of resource runs.
func (h Handler) CollectionRouteGetResourceRunStatistics(
	req CollectionRouteGetResourceRunStatisticsRequest,
) (*CollectionRouteGetResourceRunStatisticsResponse, error) {
	ctx := intercept.WithProjectInterceptor(req.Context)

	query := h.modelClient.Projects().Query().
		QueryResourceRuns()

	statusStats, err := getResourceRunStatusStats(ctx,
		query.Clone(),
		req.StartTime, req.EndTime, req.Step)
	if err != nil {
		return nil, err
	}

	statusCount, err := getResourceRunStatusCount(ctx,
		query.Clone())
	if err != nil {
		return nil, err
	}

	return &CollectionRouteGetResourceRunStatisticsResponse{
		StatusStats: statusStats,
		StatusCount: statusCount,
	}, nil
}

// getResourceRunStatusStats collects the status counts of resource runs
// according to the given time range.
func getResourceRunStatusStats(
	ctx context.Context,
	query *model.ResourceRunQuery,
	startTime, endTime time.Time,
	step string,
) ([]*RunStatusStats, error) {
	loc := startTime.Location()

	// Get time series by time range.
	timeSeries, err := timex.GetTimeSeries(startTime, endTime, step, loc)
	if err != nil {
		return nil, err
	}

	// Count by the time series and status group.
	var counts []struct {
		Count         int       `json:"count"`
		CreateTime    time.Time `json:"create_time"`
		SummaryStatus string    `json:"summary_status"`
	}
	_, offset := startTime.Zone()

	groupBy, err := sqlx.DateTruncWithZoneOffsetSQL(resourcerun.FieldCreateTime, step, offset)
	if err != nil {
		return nil, err
	}

	err = query.
		Where(
			resourcerun.CreateTimeGTE(startTime),
			resourcerun.CreateTimeLTE(endTime)).
		Modify(func(q *sql.Selector) {
			// Count.
			q.
				Select(
					sql.As(sql.Count(summaryStatus), "count"),
					sql.As(groupBy, resourcerun.FieldCreateTime),
					sql.As(summaryStatus, "summary_status")).
				GroupBy(
					groupBy,
					"summary_status")
		}).
		Scan(ctx, &counts)
	if err != nil {
		return nil, err
	}

	// Map status by time series.
	format := time.DateOnly

	switch step {
	case timex.Month:
		format = "2006-01"
	case timex.Year:
		format = "2006"
	}

	statMap := make(map[string]*RunStatusStats, 0)

	for _, t := range timeSeries {
		// Default status bucket.
		timeString := t.Format(format)
		statMap[timeString] = &RunStatusStats{}
	}

	for _, c := range counts {
		t := c.CreateTime.In(loc).Format(format)
		if _, ok := statMap[t]; !ok {
			statMap[t] = &RunStatusStats{}
		}

		switch c.SummaryStatus {
		case status.ResourceRunSummaryStatusFailed:
			statMap[t].Failed += c.Count
		case status.ResourceRunSummaryStatusSucceed:
			statMap[t].Succeeded += c.Count
		case status.ResourceRunSummaryStatusCanceled:
			statMap[t].Canceled += c.Count
		case status.ResourceRunSummaryStatusPlanned:
			statMap[t].Planned += c.Count
		case status.ResourceRunSummaryStatusRunning,
			status.ResourceRunSummaryStatusPlanning,
			status.ResourceRunSummaryStatusPending:
			statMap[t].Running += c.Count
		}
	}

	// Construct result through reducing status by time series.
	r := make([]*RunStatusStats, 0, len(statMap))

	for k, sm := range statMap {
		r = append(r, &RunStatusStats{
			RunStatusCount: RunStatusCount{
				Failed:    sm.Failed,
				Succeeded: sm.Succeeded,
				Running:   sm.Running,
				Canceled:  sm.Canceled,
				Planned:   sm.Planned,
			},
			StartTime: k,
		})
	}

	// Sort by start time.
	sort.Slice(r, func(i, j int) bool {
		return r[i].StartTime < r[j].StartTime
	})

	return r, nil
}

// getResourceRunStatusCount returns the status counts by the resource runs.
func getResourceRunStatusCount(
	ctx context.Context,
	query *model.ResourceRunQuery,
) (*RunStatusCount, error) {
	// Count by the status group.
	var counts []struct {
		SummaryStatus string `json:"summary_status"`
		Count         int    `json:"count"`
	}

	err := query.
		Modify(func(q *sql.Selector) {
			q.
				Select(
					sql.As(sql.Count(summaryStatus), "count"),
					sql.As(summaryStatus, "summary_status")).
				GroupBy("summary_status")
		}).
		Scan(ctx, &counts)
	if err != nil {
		return nil, err
	}

	// Construct result.
	var r RunStatusCount

	for _, sc := range counts {
		switch sc.SummaryStatus {
		case status.ResourceRunSummaryStatusFailed:
			r.Failed += sc.Count
		case status.ResourceRunSummaryStatusSucceed:
			r.Succeeded += sc.Count
		case status.ResourceRunSummaryStatusCanceled:
			r.Canceled += sc.Count
		case status.ResourceRunSummaryStatusPlanned:
			r.Planned += sc.Count
		case status.ResourceRunSummaryStatusRunning,
			status.ResourceRunSummaryStatusPlanning,
			status.ResourceRunSummaryStatusPending:
			r.Running += sc.Count
		}
	}

	return &r, nil
}
