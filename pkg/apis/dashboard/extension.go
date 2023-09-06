package dashboard

import (
	"context"
	"sort"
	"time"

	"entgo.io/ent/dialect/sql"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/environment"
	"github.com/seal-io/walrus/pkg/dao/model/service"
	"github.com/seal-io/walrus/pkg/dao/model/serviceresource"
	"github.com/seal-io/walrus/pkg/dao/model/servicerevision"
	"github.com/seal-io/walrus/pkg/dao/schema/intercept"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/utils/sqlx"
	"github.com/seal-io/walrus/utils/timex"
)

var getServiceRevisionFields = servicerevision.WithoutFields(
	servicerevision.FieldRecord,
	servicerevision.FieldInputPlan,
	servicerevision.FieldOutput,
	servicerevision.FieldTemplateName,
	servicerevision.FieldTemplateVersion,
	servicerevision.FieldAttributes,
	servicerevision.FieldVariables,
)

const summaryStatus = "(status ->> 'summaryStatus')"

func (h Handler) CollectionRouteGetLatestServiceRevisions(
	req CollectionRouteGetLatestServiceRevisionsRequest,
) (CollectionRouteGetLatestServiceRevisionsResponse, int, error) {
	ctx := intercept.WithProjectInterceptor(req.Context)

	entities, err := h.modelClient.ServiceRevisions().Query().
		Order(model.Desc(servicerevision.FieldCreateTime)).
		Select(getServiceRevisionFields...).
		WithProject(func(pq *model.ProjectQuery) {
			pq.Select(
				service.FieldID,
				service.FieldName,
			)
		}).
		WithService(func(sq *model.ServiceQuery) {
			sq.Select(
				service.FieldID,
				service.FieldName,
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

	return model.ExposeServiceRevisions(entities), len(entities), nil
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
	serviceNum, err := h.modelClient.Services().Query().
		Count(ctx)
	if err != nil {
		return nil, err
	}

	// Count service resources below owned projects if needed.
	var serviceResourceNum int
	if req.WithServiceResource {
		serviceResourceNum, err = h.modelClient.ServiceResources().Query().
			Where(serviceresource.ModeNEQ(types.ServiceResourceModeData)).
			Count(ctx)
		if err != nil {
			return nil, err
		}
	}

	// Count service revisions below owned projects if needed.
	var serviceRevisionNum int
	if req.WithServiceRevision {
		serviceRevisionNum, err = h.modelClient.ServiceRevisions().Query().
			Count(ctx)
		if err != nil {
			return nil, err
		}
	}

	return &CollectionRouteGetBasicInformationResponse{
		Project:         projectNum,
		Environment:     environmentNum,
		Connector:       connectorNum,
		Service:         serviceNum,
		ServiceResource: serviceResourceNum,
		ServiceRevision: serviceRevisionNum,
	}, nil
}

func (h Handler) CollectionRouteGetServiceRevisionStatistics(
	req CollectionRouteGetServiceRevisionStatisticsRequest,
) (*CollectionRouteGetServiceRevisionStatisticsResponse, error) {
	ctx := intercept.WithProjectInterceptor(req.Context)

	query := h.modelClient.Projects().Query().
		QueryServiceRevisions()

	statusStats, err := getServiceRevisionStatusStats(ctx,
		query.Clone(),
		req.StartTime, req.EndTime, req.Step)
	if err != nil {
		return nil, err
	}

	statusCount, err := getServiceRevisionStatusCount(ctx,
		query.Clone())
	if err != nil {
		return nil, err
	}

	return &CollectionRouteGetServiceRevisionStatisticsResponse{
		StatusStats: statusStats,
		StatusCount: statusCount,
	}, nil
}

// getServiceRevisionStatusStats collects the status counts of service revisions
// according to the given time range.
func getServiceRevisionStatusStats(
	ctx context.Context,
	query *model.ServiceRevisionQuery,
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

	groupBy, err := sqlx.DateTruncWithZoneOffsetSQL(servicerevision.FieldCreateTime, step, offset)
	if err != nil {
		return nil, err
	}

	err = query.
		Where(
			servicerevision.CreateTimeGTE(startTime),
			servicerevision.CreateTimeLTE(endTime)).
		Modify(func(q *sql.Selector) {
			// Count.
			q.
				Select(
					sql.As(sql.Count(summaryStatus), "count"),
					sql.As(groupBy, servicerevision.FieldCreateTime),
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
		case status.ServiceRevisionSummaryStatusFailed:
			statMap[t].Failed = c.Count
		case status.ServiceRevisionSummaryStatusSucceed:
			statMap[t].Succeed = c.Count
		case status.ServiceRevisionSummaryStatusRunning:
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

// getServiceRevisionStatusCount returns the status counts by the service revisions.
func getServiceRevisionStatusCount(
	ctx context.Context,
	query *model.ServiceRevisionQuery,
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
		case status.ServiceRevisionSummaryStatusFailed:
			r.Failed = sc.Count
		case status.ServiceRevisionSummaryStatusSucceed:
			r.Succeed = sc.Count
		case status.ServiceRevisionSummaryStatusRunning:
			r.Running = sc.Count
		}
	}

	return &r, nil
}
