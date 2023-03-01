package project

import (
	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/project/view"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/project"
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
	return "Project"
}

func (h Handler) Validating() any {
	return h.modelClient
}

// Basic APIs

func (h Handler) Create(ctx *gin.Context, req view.CreateRequest) (view.CreateResponse, error) {
	var entity = req.Model()

	var creates, err = dao.ProjectCreates(h.modelClient, entity)
	if err != nil {
		return nil, err
	}
	entity, err = creates[0].Save(ctx)
	if err != nil {
		return nil, err
	}

	return model.ExposeProject(entity), nil
}

func (h Handler) Delete(ctx *gin.Context, req view.DeleteRequest) error {
	return h.modelClient.Projects().DeleteOne(req.Model()).Exec(ctx)
}

func (h Handler) Update(ctx *gin.Context, req view.UpdateRequest) error {
	var entity = req.Model()

	var updates, err = dao.ProjectUpdates(h.modelClient, entity)
	if err != nil {
		return err
	}
	return updates[0].Exec(ctx)
}

func (h Handler) Get(ctx *gin.Context, req view.GetRequest) (view.GetResponse, error) {
	var entity, err = h.modelClient.Projects().Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return model.ExposeProject(entity), nil
}

// Batch APIs

func (h Handler) CollectionGet(ctx *gin.Context, req view.CollectionGetRequest) (view.CollectionGetResponse, int, error) {
	var query = h.modelClient.Projects().Query()

	// get count.
	cnt, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// get entities.
	var sortFields = []string{project.FieldCreateTime, project.FieldUpdateTime}
	if limit, offset, ok := req.Paging(); ok {
		query.Limit(limit).Offset(offset)
	}
	if orders, ok := req.Sorting(sortFields, model.Desc(project.FieldCreateTime)); ok {
		query.Order(orders...)
	}
	entities, err := query.
		Select(project.WithoutFields(project.FieldUpdateTime)...).
		Unique(false). // allow returning without sorting keys.
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeProjects(entities), cnt, nil
}

// Extensional APIs

var (
	applicationGetFields = application.WithoutFields(
		application.FieldProjectID,
		application.FieldUpdateTime)
	applicationSortFields = []string{
		application.FieldName,
		application.FieldCreateTime}
)

func (h Handler) GetApplications(ctx *gin.Context, req view.GetApplicationsRequest) (view.GetApplicationsResponse, int, error) {
	var query = h.modelClient.Applications().Query().
		Where(application.ProjectID(req.ID))

	// get count.
	cnt, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// get entities.
	if limit, offset, ok := req.Paging(); ok {
		query.Limit(limit).Offset(offset)
	}
	if fields, ok := req.Extracting(applicationGetFields, applicationGetFields...); ok {
		query.Select(fields...)
	}
	if orders, ok := req.Sorting(applicationSortFields, model.Desc(application.FieldCreateTime)); ok {
		query.Order(orders...)
	}
	entities, err := query.
		Unique(false). // allow returning without sorting keys.
		WithEnvironment(func(eq *model.EnvironmentQuery) {
			eq.Select(environment.FieldName)
		}).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeApplications(entities), cnt, nil
}
