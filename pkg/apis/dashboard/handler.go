package dashboard

import (
	"sort"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/dashboard/view"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/service"
	"github.com/seal-io/seal/pkg/dao/model/servicerevision"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/utils/sqlx"
	"github.com/seal-io/seal/utils/timex"
)

func Handle(mc model.ClientSet) Handler {
	return Handler{
		modelClient: mc,
	}
}

type Handler struct {
	modelClient model.ClientSet
}

func (h Handler) Kind() string {
	return "Dashboard"
}

// Basic APIs.

// Batch APIs.

var getServiceRevisionFields = servicerevision.WithoutFields(
	servicerevision.FieldStatusMessage,
	servicerevision.FieldInputPlan,
	servicerevision.FieldOutput,
	servicerevision.FieldTemplateID,
	servicerevision.FieldTemplateVersion,
	servicerevision.FieldAttributes,
	servicerevision.FieldSecrets,
)

func (h Handler) CollectionGetLatestServiceRevisions(
	ctx *gin.Context,
	_ view.CollectionGetLatestServiceRevisionsRequest,
) (view.CollectionGetLatestServiceRevisionsResponse, int, error) {
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

// Extensional APIs.

func (h Handler) CollectionRouteBasicInformation(
	ctx *gin.Context,
	_ view.BasicInfoRequest,
) (*view.BasicInfoResponse, error) {
	serviceRevisionNum, err := h.modelClient.ServiceRevisions().Query().Count(ctx)
	if err != nil {
		return nil, err
	}

	templateNum, err := h.modelClient.Templates().Query().Count(ctx)
	if err != nil {
		return nil, err
	}

	serviceNum, err := h.modelClient.Services().Query().Count(ctx)
	if err != nil {
		return nil, err
	}

	serviceResourceNum, err := h.modelClient.ServiceResources().Query().Count(ctx)
	if err != nil {
		return nil, err
	}

	environmentNum, err := h.modelClient.Environments().Query().Count(ctx)
	if err != nil {
		return nil, err
	}

	connectorNum, err := h.modelClient.Connectors().Query().Count(ctx)
	if err != nil {
		return nil, err
	}

	return &view.BasicInfoResponse{
		Template:    templateNum,
		Service:     serviceNum,
		Resource:    serviceResourceNum,
		Revision:    serviceRevisionNum,
		Environment: environmentNum,
		Connector:   connectorNum,
	}, nil
}

func (h Handler) CollectionRouteServiceRevisionStatistics(
	ctx *gin.Context,
	req view.ServiceRevisionStatisticsRequest,
) (*view.ServiceRevisionStatisticsResponse, error) {
	var (
		// StatMap map of statistics.
		statMap = make(map[string]*view.RevisionStatusStats, 0)
		// Counts count of each status.
		counts []struct {
			Count      int       `json:"count"`
			CreateTime time.Time `json:"create_time"`
			Status     string    `json:"status"`
		}
		ps = []predicate.ServiceRevision{
			servicerevision.CreateTimeGTE(req.StartTime),
			servicerevision.CreateTimeLTE(req.EndTime),
		}
	)

	// Format.
	var format string

	switch req.Step {
	case timex.Month:
		format = "2006-01"
	case timex.Year:
		format = "2006"
	default:
		format = "2006-01-02"
	}

	// Days.
	_, offset := req.StartTime.Zone()
	loc := req.StartTime.Location()

	timeSeries, err := timex.GetTimeSeries(req.StartTime, req.EndTime, req.Step, loc)
	if err != nil {
		return nil, err
	}

	for _, t := range timeSeries {
		timeString := t.Format(format)
		statMap[timeString] = &view.RevisionStatusStats{}
	}

	groupBy, err := sqlx.DateTruncWithZoneOffsetSQL(servicerevision.FieldCreateTime, req.Step, offset)
	if err != nil {
		return nil, err
	}

	// Group by.
	err = h.modelClient.ServiceRevisions().Query().
		Where(ps...).
		Modify(func(q *sql.Selector) {
			// Count.
			q.Select(
				sql.As(sql.Count(servicerevision.FieldStatus), "count"),
				sql.As(groupBy, servicerevision.FieldCreateTime),
				servicerevision.FieldStatus,
			).
				GroupBy(groupBy).
				GroupBy(servicerevision.FieldStatus)
		}).
		Scan(ctx, &counts)
	if err != nil {
		return nil, err
	}

	for _, c := range counts {
		t := c.CreateTime.In(loc).Format(format)
		if _, ok := statMap[t]; !ok {
			statMap[t] = &view.RevisionStatusStats{}
		}

		switch c.Status {
		case status.ServiceRevisionStatusFailed:
			statMap[t].Failed = c.Count
		case status.ServiceRevisionStatusSucceeded:
			statMap[t].Succeed = c.Count
		case status.ServiceRevisionStatusRunning:
			statMap[t].Running = c.Count
		}
	}

	// StatusStatistics statistics of revision status.
	statusStatistics := make([]*view.RevisionStatusStats, 0, len(statMap))
	for k, sm := range statMap {
		statusStatistics = append(statusStatistics, &view.RevisionStatusStats{
			RevisionStatusCount: view.RevisionStatusCount{
				Failed:  sm.Failed,
				Succeed: sm.Succeed,
				Running: sm.Running,
			},
			StartTime: k,
		})
	}

	sort.Slice(statusStatistics, func(i, j int) bool {
		// Sort by start time.
		return statusStatistics[i].StartTime < statusStatistics[j].StartTime
	})

	statusCount, err := h.getServiceRevisionStatusCount(ctx)
	if err != nil {
		return nil, err
	}

	return &view.ServiceRevisionStatisticsResponse{
		StatusStats: statusStatistics,
		StatusCount: statusCount,
	}, nil
}

// getServiceRevisionStatusCount returns the count of each status of service revisions.
func (h Handler) getServiceRevisionStatusCount(ctx *gin.Context) (*view.RevisionStatusCount, error) {
	var (
		currentStatusCount = &view.RevisionStatusCount{}
		statusCount        []struct {
			Status string `json:"status"`
			Count  int    `json:"count"`
		}
	)

	err := h.modelClient.ServiceRevisions().Query().
		GroupBy(servicerevision.FieldStatus).
		Aggregate(
			model.Count(),
		).
		Scan(ctx, &statusCount)
	if err != nil {
		return nil, err
	}

	for _, s := range statusCount {
		switch s.Status {
		case status.ServiceRevisionStatusFailed:
			currentStatusCount.Failed = s.Count
		case status.ServiceRevisionStatusSucceeded:
			currentStatusCount.Succeed = s.Count
		case status.ServiceRevisionStatusRunning:
			currentStatusCount.Running = s.Count
		}
	}

	return currentStatusCount, nil
}
