package project

import (
	"net/http"

	"github.com/seal-io/walrus/pkg/auths/session"
	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/project"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/utils/errorx"
)

func (h Handler) Create(req CreateRequest) (CreateResponse, error) {
	entity := req.Model()

	// Add a subject role for the project.
	sj := session.MustGetSubject(req.Context)
	entity.Edges.SubjectRoles = []*model.SubjectRoleRelationship{
		{
			SubjectID: sj.ID,
			RoleID:    types.ProjectRoleOwner,
		},
	}

	err := h.modelClient.WithTx(req.Context, func(tx *model.Tx) (err error) {
		entity, err = tx.Projects().Create().
			Set(entity).
			SaveE(req.Context, dao.ProjectSubjectRolesEdgeSave)

		return err
	})
	if err != nil {
		return nil, err
	}

	return model.ExposeProject(entity), nil
}

func (h Handler) Get(req GetRequest) (GetResponse, error) {
	entity, err := h.modelClient.Projects().Get(req.Context, req.ID)
	if err != nil {
		return nil, err
	}

	return model.ExposeProject(entity), nil
}

func (h Handler) Update(req UpdateRequest) error {
	entity := req.Model()

	return h.modelClient.Projects().UpdateOne(entity).
		Set(entity).
		ExecE(req.Context, dao.ProjectSubjectRolesEdgeSave)
}

func (h Handler) Delete(req DeleteRequest) error {
	return h.modelClient.Projects().DeleteOneID(req.ID).
		Exec(req.Context)
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

func (h Handler) CollectionGet(req CollectionGetRequest) (CollectionGetResponse, int, error) {
	query := h.modelClient.Projects().Query()

	// Filter out projects that the subject has permission to access.
	sj := session.MustGetSubject(req.Context)
	if !sj.IsAdmin() {
		pids := make([]object.ID, len(sj.ProjectRoles))
		for i := range sj.ProjectRoles {
			pids[i] = sj.ProjectRoles[i].Project.ID
		}

		query.Where(project.IDIn(pids...))
	}

	if queries, ok := req.Querying(queryFields); ok {
		query.Where(queries)
	}

	// Get count.
	cnt, err := query.Clone().Count(req.Context)
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
		Unique(false).
		All(req.Context)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeProjects(entities), cnt, nil
}

func (h Handler) CollectionDelete(req CollectionDeleteRequest) error {
	ids := req.IDs()

	// Validate whether the subject has permission to delete the projects.
	sj := session.MustGetSubject(req.Context)
	if !sj.IsAdmin() {
		for i := range ids {
			if !sj.Enforce(string(ids[i]), http.MethodDelete, "projects", string(ids[i]), req.Context.FullPath()) {
				return errorx.NewHttpError(http.StatusForbidden, "")
			}
		}
	}

	return h.modelClient.WithTx(req.Context, func(tx *model.Tx) error {
		_, err := tx.Projects().Delete().
			Where(project.IDIn(ids...)).
			Exec(req.Context)

		return err
	})
}
