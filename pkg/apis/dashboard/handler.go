package dashboard

import (
	"context"
	"sort"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/dashboard/view"
	"github.com/seal-io/seal/pkg/auths/session"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/service"
	"github.com/seal-io/seal/pkg/dao/model/serviceresource"
	"github.com/seal-io/seal/pkg/dao/model/servicerevision"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/object"
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
	servicerevision.FieldTemplateName,
	servicerevision.FieldTemplateVersion,
	servicerevision.FieldAttributes,
	servicerevision.FieldVariables,
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
	req view.CollectionRouteBasicInformationRequest,
) (*view.CollectionRouteBasicInformationResponse, error) {
	sj := session.MustGetSubject(ctx)

	sj.IncognitoOn()
	defer sj.IncognitoOff()

	var (
		isAdmin = sj.IsAdmin()
		ids     []object.ID
	)

	if !isAdmin {
		// Get owned project id list.
		ids = make([]object.ID, len(sj.ProjectRoles))
		for i := range sj.ProjectRoles {
			ids[i] = sj.ProjectRoles[i].Project.ID
		}
	}

	// Count owned projects.
	projectNum, err := h.modelClient.Projects().Query().
		Where(predicateIn[predicate.Project](isAdmin, "id", ids)...).
		Count(ctx)
	if err != nil {
		return nil, err
	}

	// Count environments below owned projects.
	environmentNum, err := h.modelClient.Environments().Query().
		Where(predicateIn[predicate.Environment](isAdmin, "project_id", ids)...).
		Count(ctx)
	if err != nil {
		return nil, err
	}

	// Count connectors below owned projects and global.
	connectorNum, err := h.modelClient.Connectors().Query().
		Where(predicateOr(
			connector.ProjectIDIsNil(), // Nil project id means configuring in global.
			predicateIn[predicate.Connector](isAdmin, "project_id", ids)...)...).
		Count(ctx)
	if err != nil {
		return nil, err
	}

	// Count services below owned projects.
	serviceNum, err := h.modelClient.Services().Query().
		Where(predicateIn[predicate.Service](isAdmin, "project_id", ids)...).
		Count(ctx)
	if err != nil {
		return nil, err
	}

	// Count service resources below owned projects if needed.
	var serviceResourceNum int
	if req.WithServiceResource {
		serviceResourceNum, err = h.modelClient.ServiceResources().Query().
			Where(serviceresource.ModeNEQ(types.ServiceResourceModeData)).
			Where(predicateIn[predicate.ServiceResource](isAdmin, "project_id", ids)...).
			Count(ctx)
		if err != nil {
			return nil, err
		}
	}

	// Count service revisions below owned projects if needed.
	var serviceRevisionNum int
	if req.WithServiceRevision {
		serviceRevisionNum, err = h.modelClient.ServiceRevisions().Query().
			Where(predicateIn[predicate.ServiceRevision](isAdmin, "project_id", ids)...).
			Count(ctx)
		if err != nil {
			return nil, err
		}
	}

	return &view.CollectionRouteBasicInformationResponse{
		Project:         projectNum,
		Environment:     environmentNum,
		Connector:       connectorNum,
		Service:         serviceNum,
		ServiceResource: serviceResourceNum,
		ServiceRevision: serviceRevisionNum,
	}, nil
}

func (h Handler) CollectionRouteServiceRevisionStatistics(
	ctx *gin.Context,
	req view.CollectionRouteServiceRevisionStatisticsRequest,
) (*view.CollectionRouteServiceRevisionStatisticsResponse, error) {
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

	return &view.CollectionRouteServiceRevisionStatisticsResponse{
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
) ([]*view.RevisionStatusStats, error) {
	loc := startTime.Location()

	// Get time series by time range.
	timeSeries, err := timex.GetTimeSeries(startTime, endTime, step, loc)
	if err != nil {
		return nil, err
	}

	// Count by the time series and status group.
	var counts []struct {
		Count      int       `json:"count"`
		CreateTime time.Time `json:"create_time"`
		Status     string    `json:"status"`
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
					sql.As(sql.Count(servicerevision.FieldStatus), "count"),
					sql.As(groupBy, servicerevision.FieldCreateTime),
					servicerevision.FieldStatus).
				GroupBy(
					groupBy,
					servicerevision.FieldStatus)
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

	statMap := make(map[string]*view.RevisionStatusStats, 0)

	for _, t := range timeSeries {
		// Default status bucket.
		timeString := t.Format(format)
		statMap[timeString] = &view.RevisionStatusStats{}
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

	// Construct result through reducing status by time series.
	r := make([]*view.RevisionStatusStats, 0, len(statMap))

	for k, sm := range statMap {
		r = append(r, &view.RevisionStatusStats{
			RevisionStatusCount: view.RevisionStatusCount{
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
) (*view.RevisionStatusCount, error) {
	// Count by the status group.
	var counts []struct {
		Status string `json:"status"`
		Count  int    `json:"count"`
	}

	err := query.
		GroupBy(servicerevision.FieldStatus).
		Aggregate(model.Count()).
		Scan(ctx, &counts)
	if err != nil {
		return nil, err
	}

	// Construct result.
	var r view.RevisionStatusCount

	for _, sc := range counts {
		switch sc.Status {
		case status.ServiceRevisionStatusFailed:
			r.Failed = sc.Count
		case status.ServiceRevisionStatusSucceeded:
			r.Succeed = sc.Count
		case status.ServiceRevisionStatusRunning:
			r.Running = sc.Count
		}
	}

	return &r, nil
}
