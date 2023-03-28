package environment

import (
	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/environment/view"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/environmentconnectorrelationship"
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
	return "Environment"
}

func (h Handler) Validating() any {
	return h.modelClient
}

// Basic APIs

func (h Handler) Create(ctx *gin.Context, req view.CreateRequest) (view.CreateResponse, error) {
	var entity = req.Model()
	var err = h.modelClient.WithTx(ctx, func(tx *model.Tx) error {
		var creates, err = dao.EnvironmentCreates(tx, entity)
		if err != nil {
			return err
		}
		entity, err = creates[0].Save(ctx)
		return err
	})
	if err != nil {
		return nil, err
	}

	return model.ExposeEnvironment(entity), nil
}

func (h Handler) Delete(ctx *gin.Context, req view.DeleteRequest) error {
	return h.modelClient.Environments().DeleteOne(req.Model()).Exec(ctx)
}

func (h Handler) Update(ctx *gin.Context, req view.UpdateRequest) error {
	var entity = req.Model()
	return h.modelClient.WithTx(ctx, func(tx *model.Tx) error {
		var updates, err = dao.EnvironmentUpdates(tx, entity)
		if err != nil {
			return err
		}
		return updates[0].Exec(ctx)
	})
}

func (h Handler) Get(ctx *gin.Context, req view.GetRequest) (view.GetResponse, error) {
	var entity, err = h.modelClient.Environments().Query().
		Where(environment.ID(req.ID)).
		WithConnectors(func(rq *model.EnvironmentConnectorRelationshipQuery) {
			rq.Order(model.Desc(environmentconnectorrelationship.FieldCreateTime)).
				Select(environmentconnectorrelationship.FieldEnvironmentID).
				// allow returning without sorting keys.
				Unique(false).
				Select(environmentconnectorrelationship.FieldConnectorID).
				WithConnector(
					func(cq *model.ConnectorQuery) {
						cq.Select(
							connector.FieldID,
							connector.FieldType,
							connector.FieldName)
					})
		}).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return model.ExposeEnvironment(entity), nil
}

// Batch APIs

func (h Handler) CollectionDelete(ctx *gin.Context, req view.CollectionDeleteRequest) error {
	return h.modelClient.WithTx(ctx, func(tx *model.Tx) (err error) {
		for i := range req {
			err = tx.Environments().DeleteOne(req[i].Model()).
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
		environment.FieldName,
	}
	getFields = environment.WithoutFields(
		environment.FieldUpdateTime)
	sortFields = []string{
		environment.FieldName,
		environment.FieldCreateTime}
)

func (h Handler) CollectionGet(ctx *gin.Context, req view.CollectionGetRequest) (view.CollectionGetResponse, int, error) {
	var query = h.modelClient.Environments().Query()
	if queries, ok := req.Querying(queryFields); ok {
		query.Where(queries)
	}

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
	if orders, ok := req.Sorting(sortFields, model.Desc(environment.FieldCreateTime)); ok {
		query.Order(orders...)
	}
	entities, err := query.
		// allow returning without sorting keys.
		Unique(false).
		// must extract connectors.
		Select(environment.FieldID).
		WithConnectors(func(rq *model.EnvironmentConnectorRelationshipQuery) {
			// includes connectors.
			rq.Order(model.Desc(environmentconnectorrelationship.FieldCreateTime)).
				WithConnector(func(cq *model.ConnectorQuery) {
					cq.Select(
						connector.FieldID,
						connector.FieldType,
						connector.FieldName)
				})
		}).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeEnvironments(entities), cnt, nil
}

// Extensional APIs
