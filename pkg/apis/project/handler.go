package project

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/project/view"
	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/auths/session"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/project"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/object"
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

// Basic APIs.

func (h Handler) Create(ctx *gin.Context, req view.CreateRequest) (view.CreateResponse, error) {
	entity := req.Model()

	// Add a subject role for the project.
	s := session.MustGetSubject(ctx)
	entity.Edges.SubjectRoles = []*model.SubjectRoleRelationship{
		{
			SubjectID: s.ID,
			RoleID:    types.ProjectRoleOwner,
		},
	}

	creates, err := dao.ProjectCreates(h.modelClient, entity)
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
	s := session.MustGetSubject(ctx)
	if !s.Enforce(req.ID, "", "projects", http.MethodDelete, string(req.ID), "") {
		return runtime.Errorc(http.StatusForbidden)
	}

	return h.modelClient.Projects().DeleteOne(req.Model()).Exec(ctx)
}

func (h Handler) Update(ctx *gin.Context, req view.UpdateRequest) error {
	s := session.MustGetSubject(ctx)
	if !s.Enforce(req.ID, "", "projects", http.MethodPut, string(req.ID), "") {
		return runtime.Errorc(http.StatusForbidden)
	}

	entity := req.Model()

	updates, err := dao.ProjectUpdates(h.modelClient, entity)
	if err != nil {
		return err
	}

	return updates[0].Exec(ctx)
}

func (h Handler) Get(ctx *gin.Context, req view.GetRequest) (view.GetResponse, error) {
	s := session.MustGetSubject(ctx)
	if !s.Enforce(req.ID, "", "projects", http.MethodGet, string(req.ID), "") {
		return nil, runtime.Errorc(http.StatusForbidden)
	}

	entity, err := h.modelClient.Projects().Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return model.ExposeProject(entity), nil
}

// Batch APIs.

func (h Handler) CollectionDelete(ctx *gin.Context, req view.CollectionDeleteRequest) error {
	s := session.MustGetSubject(ctx)
	for i := range req {
		if !s.Enforce(req[i].ID, "", "projects", http.MethodDelete, string(req[i].ID), "") {
			return runtime.Errorc(http.StatusForbidden)
		}
	}

	return h.modelClient.WithTx(ctx, func(tx *model.Tx) (err error) {
		for i := range req {
			err = tx.Projects().DeleteOne(req[i].Model()).
				Exec(ctx)
			if err != nil {
				return err
			}
		}

		return
	})
}

var (
	queryFields = []string{
		project.FieldName,
	}
	getFields = project.WithoutFields(
		project.FieldUpdateTime)
	sortFields = []string{
		project.FieldName,
		project.FieldCreateTime,
	}
)

func (h Handler) CollectionGet(
	ctx *gin.Context,
	req view.CollectionGetRequest,
) (view.CollectionGetResponse, int, error) {
	query := h.modelClient.Projects().Query()

	s := session.MustGetSubject(ctx)
	if !s.IsAdmin() {
		pids := make([]object.ID, len(s.ProjectRoles))
		for i := range s.ProjectRoles {
			pids[i] = s.ProjectRoles[i].Project.ID
		}

		query.Where(project.IDIn(pids...))
	}

	if queries, ok := req.Querying(queryFields); ok {
		query.Where(queries)
	}

	// Get count.
	cnt, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Get entities.
	if limit, offset, ok := req.Paging(); ok {
		query.Limit(limit).Offset(offset)
	}

	if fields, ok := req.Extracting(getFields, getFields...); ok {
		query.Select(fields...)
	}

	if orders, ok := req.Sorting(sortFields, model.Desc(project.FieldCreateTime)); ok {
		query.Order(orders...)
	}

	entities, err := query.
		// Allow returning without sorting keys.
		Unique(false).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeProjects(entities), cnt, nil
}

// Extensional APIs.
