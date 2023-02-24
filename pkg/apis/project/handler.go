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

func (h Handler) Create(ctx *gin.Context, req view.CreateRequest) (*view.CreateResponse, error) {
	var input = &model.Project{
		Name:        req.Name,
		Description: req.Description,
		Labels:      req.Labels,
	}

	var creates, err = dao.ProjectCreates(h.modelClient, input)
	if err != nil {
		return nil, err
	}

	return creates[0].Save(ctx)
}

func (h Handler) Delete(ctx *gin.Context, req view.DeleteRequest) error {
	return h.modelClient.Projects().DeleteOneID(req.ID).Exec(ctx)
}

func (h Handler) Update(ctx *gin.Context, req view.UpdateRequest) error {
	var input = &model.Project{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		Labels:      req.Labels,
	}

	var updates, err = dao.ProjectUpdates(h.modelClient, input)
	if err != nil {
		return err
	}

	return updates[0].Exec(ctx)
}

func (h Handler) Get(ctx *gin.Context, req view.GetRequest) (*model.Project, error) {
	return h.modelClient.Projects().Get(ctx, req.ID)
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

	return entities, cnt, nil
}

// Extensional APIs

func (h Handler) GetApplications(ctx *gin.Context, req view.GetApplicationsRequest) (view.GetApplicationsResponse, int, error) {
	var query = h.modelClient.Applications().Query()

	// get count.
	cnt, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// get entities.
	var sortFields = []string{application.FieldCreateTime, application.FieldUpdateTime}
	if limit, offset, ok := req.Paging(); ok {
		query.Limit(limit).Offset(offset)
	}
	if orders, ok := req.Sorting(sortFields, model.Desc(application.FieldCreateTime)); ok {
		query.Order(orders...)
	}
	entities, err := query.
		Where(application.ProjectID(req.ID)).
		Select(application.FieldID, application.FieldName, application.FieldDescription, application.FieldEnvironmentID).
		Unique(false). // allow returning without sorting keys.
		WithEnvironment(func(eq *model.EnvironmentQuery) {
			eq.Select(environment.FieldName)
		}).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	var resp = make(view.GetApplicationsResponse, len(entities))
	for i := 0; i < len(entities); i++ {
		// move `.Edges.Environment.Name` to `.EnvironmentName`.
		resp[i] = view.GetApplicationResponse{
			Application:     entities[i],
			EnvironmentName: entities[i].Edges.Environment.Name,
		}
		entities[i].Edges.Environment = nil // release
	}
	return resp, cnt, nil
}
