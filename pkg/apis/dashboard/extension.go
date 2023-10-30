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
	"github.com/seal-io/walrus/pkg/dao/model/resourcerevision"
	"github.com/seal-io/walrus/pkg/dao/schema/intercept"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/utils/sqlx"
	"github.com/seal-io/walrus/utils/timex"
)

var getResourceRevisionFields = resourcerevision.WithoutFields(
	resourcerevision.FieldRecord,
	resourcerevision.FieldInputPlan,
	resourcerevision.FieldOutput,
	resourcerevision.FieldTemplateName,
	resourcerevision.FieldTemplateVersion,
	resourcerevision.FieldAttributes,
	resourcerevision.FieldVariables,
)

const summaryStatus = "(status ->> 'summaryStatus')"

func (h Handler) CollectionRouteGetLatestResourceRevisions(
	req CollectionRouteGetLatestResourceRevisionsRequest,
) (CollectionRouteGetLatestResourceRevisionsResponse, int, error) {
	ctx := intercept.WithProjectInterceptor(req.Context)

	entities, err := h.modelClient.ResourceRevisions().Query().
		Order(model.Desc(resourcerevision.FieldCreateTime)).
		Select(getResourceRevisionFields...).
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
		Limit(10).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeResourceRevisions(entities), len(entities), nil
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

	// Count services below owned projects.
	serviceNum, err := h.modelClient.Resources().Query().
		Count(ctx)
	if err != nil {
		return nil, err
	}

	// Count service resources below owned projects if needed.
	var serviceResourceNum int
	if req.WithResourceComponent {
		serviceResourceNum, err = h.modelClient.ResourceComponents().Query().
			Where(resourcecomponent.ModeNEQ(types.ResourceComponentModeData)).
			Count(ctx)
		if err != nil {
			return nil, err
		}
	}

	// Count service revisions below owned projects if needed.
	var serviceRevisionNum int
	if req.WithResourceRevision {
		serviceRevisionNum, err = h.modelClient.ResourceRevisions().Query().
			Count(ctx)
		if err != nil {
			return nil, err
		}
	}

	return &CollectionRouteGetBasicInformationResponse{
		Project:           projectNum,
		Environment:       environmentNum,
		Connector:         connectorNum,
		Service:           serviceNum,
		ResourceComponent: serviceResourceNum,
		ResourceRevision:  serviceRevisionNum,
	}, nil
}

func (h Handler) CollectionRouteGetResourceRevisionStatistics(
	req CollectionRouteGetResourceRevisionStatisticsRequest,
) (*CollectionRouteGetResourceRevisionStatisticsResponse, error) {
	ctx := intercept.WithProjectInterceptor(req.Context)

	query := h.modelClient.Projects().Query().
		QueryResourceRevisions()

	statusStats, err := getResourceRevisionStatusStats(ctx,
		query.Clone(),
		req.StartTime, req.EndTime, req.Step)
	if err != nil {
		return nil, err
	}

	statusCount, err := getResourceRevisionStatusCount(ctx,
		query.Clone())
	if err != nil {
		return nil, err
	}

	return &CollectionRouteGetResourceRevisionStatisticsResponse{
		StatusStats: statusStats,
		StatusCount: statusCount,
	}, nil
}

// getResourceRevisionStatusStats collects the status counts of service revisions
// according to the given time range.
func getResourceRevisionStatusStats(
	ctx context.Context,
	query *model.ResourceRevisionQuery,
	startTime, endTime time.Time,
	step string,
) ([]*RevisionStatusStats, error) {
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

	groupBy, err := sqlx.DateTruncWithZoneOffsetSQL(resourcerevision.FieldCreateTime, step, offset)
	if err != nil {
		return nil, err
	}

	err = query.
		Where(
			resourcerevision.CreateTimeGTE(startTime),
			resourcerevision.CreateTimeLTE(endTime)).
		Modify(func(q *sql.Selector) {
			// Count.
			q.
				Select(
					sql.As(sql.Count(summaryStatus), "count"),
					sql.As(groupBy, resourcerevision.FieldCreateTime),
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
	format := "2006-01-02"

	switch step {
	case timex.Month:
		format = "2006-01"
	case timex.Year:
		format = "2006"
	}

	statMap := make(map[string]*RevisionStatusStats, 0)

	for _, t := range timeSeries {
		// Default status bucket.
		timeString := t.Format(format)
		statMap[timeString] = &RevisionStatusStats{}
	}

	for _, c := range counts {
		t := c.CreateTime.In(loc).Format(format)
		if _, ok := statMap[t]; !ok {
			statMap[t] = &RevisionStatusStats{}
		}

		switch c.SummaryStatus {
		case status.ResourceRevisionSummaryStatusFailed:
			statMap[t].Failed = c.Count
		case status.ResourceRevisionSummaryStatusSucceed:
			statMap[t].Succeed = c.Count
		case status.ResourceRevisionSummaryStatusRunning:
			statMap[t].Running = c.Count
		}
	}

	// Construct result through reducing status by time series.
	r := make([]*RevisionStatusStats, 0, len(statMap))

	for k, sm := range statMap {
		r = append(r, &RevisionStatusStats{
			RevisionStatusCount: RevisionStatusCount{
				Failed:  sm.Failed,
				Succeed: sm.Succeed,
				Running: sm.Running,
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

// getResourceRevisionStatusCount returns the status counts by the service revisions.
func getResourceRevisionStatusCount(
	ctx context.Context,
	query *model.ResourceRevisionQuery,
) (*RevisionStatusCount, error) {
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
	var r RevisionStatusCount

	for _, sc := range counts {
		switch sc.SummaryStatus {
		case status.ResourceRevisionSummaryStatusFailed:
			r.Failed = sc.Count
		case status.ResourceRevisionSummaryStatusSucceed:
			r.Succeed = sc.Count
		case status.ResourceRevisionSummaryStatusRunning:
			r.Running = sc.Count
		}
	}

	return &r, nil
}
