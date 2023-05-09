package secret

import (
	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/secret/view"
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

// Basic APIs

var (
	queryFields = []string{
		secret.FieldName,
	}
)

func (h Handler) Create(ctx *gin.Context, req view.CreateRequest) (view.CreateResponse, error) {
	var entity = req.Model()
	var err = h.modelClient.WithTx(ctx, func(tx *model.Tx) error {
		var creates, err = dao.SecretCreates(tx, entity)
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
	var entity = req.Model()
	return h.modelClient.WithTx(ctx, func(tx *model.Tx) error {
		var updates, err = dao.SecretUpdates(tx, entity)
		if err != nil {
			return err
		}
		return updates[0].Exec(ctx)
	})
}

// Batch APIs

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

var (
	getFields = secret.WithoutFields(
		secret.FieldValue)
)

func (h Handler) CollectionGet(ctx *gin.Context, req view.CollectionGetRequest) (view.CollectionGetResponse, int, error) {
	var query = h.modelClient.Secrets().Query()
	if len(req.ProjectIDs) != 0 {
		// project scope
		query.Where(secret.ProjectIDIn(req.ProjectIDs...))
	} else {
		// global scope
		query.Where(secret.ProjectIDIsNil())
	}
	if queries, ok := req.Querying(queryFields); ok {
		query.Where(queries)
	}

	// get count.
	cnt, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// get entities.
	if limit, offset, ok := req.Paging(); !ok {
		query.Limit(limit).Offset(offset)
	}
	entities, err := query.
		Order(model.Desc(secret.FieldCreateTime)).
		Select(getFields...).
		// allow returning without sorting keys.
		Unique(false).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeSecrets(entities), cnt, nil
}

// Extensional APIs
