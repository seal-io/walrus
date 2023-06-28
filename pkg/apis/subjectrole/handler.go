package subjectrole

import (
	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/subjectrole/view"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/subject"
	"github.com/seal-io/seal/pkg/dao/model/subjectrolerelationship"
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
	return "SubjectRole"
}

// Basic APIs.

func (h Handler) Create(ctx *gin.Context, req view.CreateRequest) (view.CreateResponse, error) {
	entity := req.Model()

	err := h.modelClient.WithTx(ctx, func(tx *model.Tx) error {
		creates, err := dao.SubjectRoleRelationshipCreates(tx, entity)
		if err != nil {
			return err
		}
		entity, err = creates[0].Save(ctx)

		return err
	})
	if err != nil {
		return nil, err
	}

	return model.ExposeSubjectRoleRelationship(entity), nil
}

func (h Handler) Delete(ctx *gin.Context, req view.DeleteRequest) error {
	return h.modelClient.SubjectRoleRelationships().DeleteOne(req.Model()).Exec(ctx)
}

// Batch APIs.

func (h Handler) CollectionDelete(ctx *gin.Context, req view.CollectionDeleteRequest) error {
	return h.modelClient.WithTx(ctx, func(tx *model.Tx) (err error) {
		for i := range req {
			err = tx.SubjectRoleRelationships().DeleteOne(req[i].Model()).
				Exec(ctx)
			if err != nil {
				return err
			}
		}

		return
	})
}

func (h Handler) CollectionGet(
	ctx *gin.Context,
	req view.CollectionGetRequest,
) (view.CollectionGetResponse, int, error) {
	query := h.modelClient.SubjectRoleRelationships().Query()
	if req.ProjectID != "" {
		// Project scope.
		query.Where(subjectrolerelationship.ProjectIDIn(req.ProjectID))
	} else {
		// Global scope.
		query.Where(subjectrolerelationship.ProjectIDIsNil())
	}

	// Get count.
	cnt, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Get entities.
	if limit, offset, ok := req.Paging(); !ok {
		query.Limit(limit).Offset(offset)
	}

	entities, err := query.
		Order(model.Desc(subjectrolerelationship.FieldCreateTime)).
		WithSubject(func(sq *model.SubjectQuery) {
			sq.Select(
				subject.FieldID,
				subject.FieldKind,
				subject.FieldDomain,
				subject.FieldName)
		}).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeSubjectRoleRelationships(entities), cnt, nil
}

// Extensional APIs.
