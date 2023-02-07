package group

import (
	"net/http"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/group/view"
	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/subject"
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
	return "Group"
}

func (h Handler) Validating() any {
	return h.modelClient
}

// Basic APIs

func (h Handler) Create(ctx *gin.Context, req view.CreateRequest) error {
	var input = &model.Subject{
		Kind:        "group",
		Group:       req.Group,
		Name:        req.Name,
		Description: req.Description,
		Paths:       req.Paths,
		Builtin:     false,
	}

	var creates, err = dao.SubjectCreates(h.modelClient, input)
	if err != nil {
		return err
	}
	_, err = creates[0].Save(ctx)
	return err
}

func (h Handler) Delete(ctx *gin.Context, req view.DeleteRequest) error {
	// get sub subjects.
	var subres, err = h.modelClient.Subjects().Query().
		Where(func(s *sql.Selector) {
			s.Where(sql.And(
				sqljson.ValueContains(subject.FieldPaths, req.Name),
				sql.NEQ(subject.FieldName, req.Name),
			))
		}).
		Select(subject.FieldID, subject.FieldKind, subject.FieldName, subject.FieldMountTo).
		All(ctx)
	if err != nil {
		return runtime.ErrorfP(http.StatusInternalServerError, "failed to get subresource of group: %w", err)
	}

	var inputs = [][]predicate.Subject{
		{
			subject.ID(req.ID),
			subject.Kind("group"),
		},
	}
	for i := 0; i < len(subres); i++ {
		switch subres[i].Kind {
		case "group":
			inputs = append(inputs, []predicate.Subject{
				subject.ID(subres[i].ID),
				subject.Kind("group"),
			})
		case "user":
			if *subres[i].MountTo {
				inputs = append(inputs, []predicate.Subject{
					subject.ID(subres[i].ID),
				})
				continue
			}
			inputs = append(inputs, []predicate.Subject{
				subject.Kind("user"),
				subject.Name(req.Name),
			})
		}
	}
	return h.modelClient.WithTx(ctx, func(tx *model.Tx) error {
		for i := range inputs {
			_, err = tx.Subjects().Delete().
				Where(inputs[i]...).
				Exec(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	})
	// TODO clean cache
}

func (h Handler) Update(ctx *gin.Context, req view.UpdateRequest) error {
	var input = &model.Subject{
		Kind:        "group",
		ID:          req.ID,
		Description: req.Description,
	}

	var updates, err = dao.SubjectUpdates(h.modelClient, input)
	if err != nil {
		return err
	}
	return updates[0].Exec(ctx)
}

// Batch APIs

var (
	getFields  = subject.WithoutFields(subject.FieldCreateTime, subject.FieldLoginTo)
	sortFields = []string{subject.FieldCreateTime, subject.FieldUpdateTime}
)

func (h Handler) CollectionGet(ctx *gin.Context, req view.CollectionGetRequest) (view.CollectionGetResponse, int, error) {
	var input = []predicate.Subject{
		subject.Kind("group"),
	}
	if req.Group != "" {
		input = append(input, subject.Group(req.Group))
	}

	var query = h.modelClient.Subjects().Query().
		Where(input...)

	// get count.
	cnt, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// get entities.
	if limit, offset, ok := req.Paging(); ok {
		query.Limit(limit).Offset(offset)
	}
	if fields, ok := req.Extracting(getFields, getFields...); ok {
		query.Select(fields...)
	}
	if orders, ok := req.Sorting(sortFields); ok {
		query.Order(orders...)
	}
	entities, err := query.
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return entities, cnt, nil
}

// Extensional APIs
