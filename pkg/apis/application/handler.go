package application

import (
	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/application/view"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/model/applicationmodulerelationship"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/model/applicationrevision"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/platform"
	"github.com/seal-io/seal/pkg/platform/operator"
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
	return "Application"
}

func (h Handler) Validating() any {
	return h.modelClient
}

// Basic APIs

func (h Handler) Create(ctx *gin.Context, req view.CreateRequest) (resp view.CreateResponse, err error) {
	var input = req.Model()

	err = h.modelClient.WithTx(ctx, func(tx *model.Tx) error {
		var creates, err = dao.ApplicationCreates(tx, input)
		if err != nil {
			return err
		}

		resp.Application, err = creates[0].Save(ctx)
		if err != nil {
			return err
		}
		return nil
	})
	return
}

func (h Handler) Delete(ctx *gin.Context, req view.DeleteRequest) error {
	return h.modelClient.Applications().DeleteOneID(req.ID).Exec(ctx)
}

func (h Handler) Update(ctx *gin.Context, req view.UpdateRequest) error {
	var input = req.Model()

	return h.modelClient.WithTx(ctx, func(tx *model.Tx) error {
		var updates, err = dao.ApplicationUpdates(tx, input)
		if err != nil {
			return err
		}

		return updates[0].Exec(ctx)
	})
}

func (h Handler) Get(ctx *gin.Context, req view.GetRequest) (resp view.GetResponse, err error) {
	resp.Application, err = h.modelClient.Applications().Query().
		Where(application.ID(req.ID)).
		WithApplicationModuleRelationships(func(rq *model.ApplicationModuleRelationshipQuery) {
			rq.Select(
				applicationmodulerelationship.FieldModuleID,
				applicationmodulerelationship.FieldName,
				applicationmodulerelationship.FieldVariables)
		}).
		Only(ctx)
	return
}

// Batch APIs

// Extensional APIs

func (h Handler) GetResources(ctx *gin.Context, req view.GetResourcesRequest) (view.GetResourcesResponse, int, error) {
	var query = h.modelClient.ApplicationResources().Query().
		Where(applicationresource.ApplicationID(req.ID))

	// get count.
	cnt, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// get entities.
	var sortFields = []string{
		applicationresource.FieldCreateTime,
		applicationresource.FieldModule,
		applicationresource.FieldMode,
		applicationresource.FieldType,
		applicationresource.FieldName,
	}
	if limit, offset, ok := req.Paging(); ok {
		query.Limit(limit).Offset(offset)
	}
	if orders, ok := req.Sorting(sortFields, model.Desc(applicationresource.FieldCreateTime)); ok {
		query.Order(orders...)
	}
	entities, err := query.
		Select(applicationresource.WithoutFields(
				applicationresource.FieldApplicationID,
				applicationresource.FieldUpdateTime)...).
		Unique(false). // allow returning without sorting keys.
		WithConnector(func(cq *model.ConnectorQuery) {
			cq.Select(
				connector.FieldName,
				connector.FieldType,
				connector.FieldConfigVersion,
				connector.FieldConfigData)
		}).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	// construct response.
	var resp = make(view.GetResourcesResponse, len(entities))
	for i := 0; i < len(entities); i++ {
		resp[i].ApplicationResource = entities[i]

		if !req.WithoutKeys {
			// fetch operable keys.
			var op operator.Operator
			op, err = platform.GetOperator(ctx,
				operator.CreateOptions{Connector: *entities[i].Edges.Connector})
			if err != nil {
				return nil, 0, err
			}
			resp[i].OperatorKeys, err = op.GetKeys(ctx, *entities[i])
			if err != nil {
				return nil, 0, err
			}
		}
	}
	return resp, cnt, nil
}

func (h Handler) GetRevisions(ctx *gin.Context, req view.GetRevisionsRequest) (view.GetRevisionsResponse, int, error) {
	var query = h.modelClient.ApplicationRevisions().Query().
		Where(applicationrevision.ApplicationID(req.ID))

	// get count.
	cnt, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// get entities.
	var sortFields = []string{applicationrevision.FieldCreateTime}
	if limit, offset, ok := req.Paging(); ok {
		query.Limit(limit).Offset(offset)
	}
	if orders, ok := req.Sorting(sortFields, model.Desc(applicationrevision.FieldCreateTime)); ok {
		query.Order(orders...)
	}
	entities, err := query.
		Select(
				applicationrevision.FieldID,
				applicationrevision.FieldCreateTime,
				applicationrevision.FieldStatus,
				applicationrevision.FieldStatusMessage).
		Unique(false). // allow returning without sorting keys.
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return entities, cnt, nil
}
