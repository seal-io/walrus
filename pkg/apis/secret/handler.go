package secret

import (
	"entgo.io/ent/dialect/sql"
	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/secret/view"
	"github.com/seal-io/seal/pkg/auths/session"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/secret"
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
	return "Secret"
}

func (h Handler) Validating() any {
	return h.modelClient
}

// Basic APIs.

var queryFields = []string{
	secret.FieldName,
}

func (h Handler) Create(ctx *gin.Context, req view.CreateRequest) (view.CreateResponse, error) {
	entity := req.Model()

	err := h.modelClient.WithTx(ctx, func(tx *model.Tx) error {
		creates, err := dao.SecretCreates(tx, entity)
		if err != nil {
			return err
		}
		entity, err = creates[0].Save(ctx)

		return err
	})
	if err != nil {
		return nil, err
	}

	return model.ExposeSecret(entity), nil
}

func (h Handler) Delete(ctx *gin.Context, req view.DeleteRequest) error {
	return h.modelClient.Secrets().DeleteOne(req.Model()).Exec(ctx)
}

func (h Handler) Update(ctx *gin.Context, req view.UpdateRequest) error {
	entity := req.Model()

	return h.modelClient.WithTx(ctx, func(tx *model.Tx) error {
		updates, err := dao.SecretUpdates(tx, entity)
		if err != nil {
			return err
		}

		return updates[0].Exec(ctx)
	})
}

// Batch APIs.

func (h Handler) CollectionDelete(ctx *gin.Context, req view.CollectionDeleteRequest) error {
	return h.modelClient.WithTx(ctx, func(tx *model.Tx) (err error) {
		for i := range req {
			err = tx.Secrets().DeleteOne(req[i].Model()).
				Exec(ctx)
			if err != nil {
				return err
			}
		}

		return
	})
}

var getFields = secret.WithoutFields(
	secret.FieldValue)

func (h Handler) CollectionGet(
	ctx *gin.Context,
	req view.CollectionGetRequest,
) (view.CollectionGetResponse, int, error) {
	s := session.MustGetSubject(ctx)

	s.IncognitoOn()
	defer s.IncognitoOff()

	query := h.modelClient.Secrets().Query()

	if len(req.ProjectIDs) != 0 {
		if req.WithGlobal {
			// With global scope.
			query.Where(secret.Or(
				secret.ProjectIDIsNil(),
				secret.ProjectIDIn(req.ProjectIDs...)))
		} else {
			// Project scope only.
			query.Where(secret.ProjectIDIn(req.ProjectIDs...))
		}
	} else {
		// Global scope.
		query.Where(secret.ProjectIDIsNil())
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
	if limit, offset, ok := req.Paging(); !ok {
		query.Limit(limit).Offset(offset)
	}

	if req.WithGlobal {
		// With global scope, make it unique.
		query.Modify(func(s *sql.Selector) {
			s.Select(secret.FieldName).
				Distinct()
		})
	} else {
		if fields, ok := req.Extracting(getFields, getFields...); ok {
			query.Select(fields...)
		}
		// Allow returning without sorting keys.
		query.Order(model.Desc(secret.FieldCreateTime)).
			Unique(false)
	}

	entities, err := query.All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeSecrets(entities), cnt, nil
}

// Extensional APIs.
