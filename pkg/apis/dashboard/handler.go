package dashboard

import (
	"sort"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/dashboard/view"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/model/applicationinstance"
	"github.com/seal-io/seal/pkg/dao/model/applicationrevision"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/project"
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

var getApplicationRevisionFields = applicationrevision.WithoutFields(
	applicationrevision.FieldStatusMessage,
	applicationrevision.FieldInputPlan,
	applicationrevision.FieldOutput,
	applicationrevision.FieldInputVariables,
	applicationrevision.FieldModules,
	applicationrevision.FieldVariables,
	applicationrevision.FieldSecrets,
)

func (h Handler) CollectionGetLatestApplicationRevisions(
	ctx *gin.Context,
	_ view.CollectionGetLatestApplicationRevisionsRequest,
) (view.CollectionGetLatestApplicationRevisionsResponse, int, error) {
	entities, err := h.modelClient.ApplicationRevisions().Query().
		Order(model.Desc(applicationrevision.FieldCreateTime)).
		Select(getApplicationRevisionFields...).
		WithInstance(func(aiq *model.ApplicationInstanceQuery) {
			aiq.Select(
				applicationinstance.FieldID,
				applicationinstance.FieldName,
				applicationinstance.FieldApplicationID,
			).WithApplication(func(aq *model.ApplicationQuery) {
				aq.Select(
					application.FieldID,
					application.FieldName,
					application.FieldProjectID).
					WithProject(func(pq *model.ProjectQuery) {
						pq.Select(project.FieldID, project.FieldName)
					})
			})
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

	return model.ExposeApplicationRevisions(entities), len(entities), nil
}

// Extensional APIs.

func (h Handler) CollectionRouteBasicInformation(
	ctx *gin.Context,
	_ view.BasicInfoRequest,
) (*view.BasicInfoResponse, error) {
	applicationNum, err := h.modelClient.Applications().Query().Count(ctx)
	if err != nil {
		return nil, err
	}

	applicationRevisionNum, err := h.modelClient.ApplicationRevisions().Query().Count(ctx)
	if err != nil {
		return nil, err
	}

	moduleNum, err := h.modelClient.Modules().Query().Count(ctx)
	if err != nil {
		return nil, err
	}

	instanceNum, err := h.modelClient.ApplicationInstances().Query().Count(ctx)
	if err != nil {
		return nil, err
	}

	applicationResourceNum, err := h.modelClient.ApplicationResources().Query().Count(ctx)
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
		Application: applicationNum,
		Module:      moduleNum,
		Instance:    instanceNum,
		Resource:    applicationResourceNum,
		Revision:    applicationRevisionNum,
		Environment: environmentNum,
		Connector:   connectorNum,
	}, nil
}

func (h Handler) CollectionRouteApplicationRevisionStatistics(
	ctx *gin.Context,
	req view.ApplicationRevisionStatisticsRequest,
) (*view.ApplicationRevisionStatisticsResponse, error) {
	var (
		// StatMap map of statistics.
		statMap = make(map[string]*view.RevisionStatusStats, 0)
		// Counts count of each status.
		counts []struct {
			Count      int       `json:"count"`
			CreateTime time.Time `json:"create_time"`
			Status     string    `json:"status"`
		}
		ps = []predicate.ApplicationRevision{
			applicationrevision.CreateTimeGTE(req.StartTime),
			applicationrevision.CreateTimeLTE(req.EndTime),
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

	groupBy, err := sqlx.DateTruncWithZoneOffsetSQL(applicationrevision.FieldCreateTime, req.Step, offset)
	if err != nil {
		return nil, err
	}

	// Group by.
	err = h.modelClient.ApplicationRevisions().Query().
		Where(ps...).
		Modify(func(q *sql.Selector) {
			// Count.
			q.Select(
				sql.As(sql.Count(applicationrevision.FieldStatus), "count"),
				sql.As(groupBy, applicationrevision.FieldCreateTime),
				applicationrevision.FieldStatus,
			).
				GroupBy(groupBy).
				GroupBy(applicationrevision.FieldStatus)
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
		case status.ApplicationRevisionStatusFailed:
			statMap[t].Failed = c.Count
		case status.ApplicationRevisionStatusSucceeded:
			statMap[t].Succeed = c.Count
		case status.ApplicationRevisionStatusRunning:
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

	statusCount, err := h.getApplicationRevisionStatusCount(ctx)
	if err != nil {
		return nil, err
	}

	return &view.ApplicationRevisionStatisticsResponse{
		StatusStats: statusStatistics,
		StatusCount: statusCount,
	}, nil
}

// getApplicationRevisionStatusCount returns the count of each status of application revisions.
func (h Handler) getApplicationRevisionStatusCount(ctx *gin.Context) (*view.RevisionStatusCount, error) {
	var (
		currentStatusCount = &view.RevisionStatusCount{}
		statusCount        []struct {
			Status string `json:"status"`
			Count  int    `json:"count"`
		}
	)

	err := h.modelClient.ApplicationRevisions().Query().
		GroupBy(applicationrevision.FieldStatus).
		Aggregate(
			model.Count(),
		).
		Scan(ctx, &statusCount)
	if err != nil {
		return nil, err
	}

	for _, s := range statusCount {
		switch s.Status {
		case status.ApplicationRevisionStatusFailed:
			currentStatusCount.Failed = s.Count
		case status.ApplicationRevisionStatusSucceeded:
			currentStatusCount.Succeed = s.Count
		case status.ApplicationRevisionStatusRunning:
			currentStatusCount.Running = s.Count
		}
	}

	return currentStatusCount, nil
}
