package application

import (
	"entgo.io/ent/dialect/sql"
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/rest"

	"github.com/seal-io/seal/pkg/apis/application/view"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/model/applicationinstance"
	"github.com/seal-io/seal/pkg/dao/model/applicationmodulerelationship"
	"github.com/seal-io/seal/pkg/dao/model/environment"
)

func Handle(mc model.ClientSet, kc *rest.Config, tc bool) Handler {
	return Handler{
		modelClient:  mc,
		kubeConfig:   kc,
		tlsCertified: tc,
	}
}

type Handler struct {
	modelClient  model.ClientSet
	kubeConfig   *rest.Config
	tlsCertified bool
}

func (h Handler) Kind() string {
	return "Application"
}

func (h Handler) Validating() any {
	return h.modelClient
}

// Basic APIs

func (h Handler) Create(ctx *gin.Context, req view.CreateRequest) (view.CreateResponse, error) {
	var entity = req.Model()
	var err = h.modelClient.WithTx(ctx, func(tx *model.Tx) error {
		var creates, err = dao.ApplicationCreates(tx, entity)
		if err != nil {
			return err
		}
		entity, err = creates[0].Save(ctx)
		return err
	})
	if err != nil {
		return nil, err
	}

	return model.ExposeApplication(entity), nil
}

func (h Handler) Delete(ctx *gin.Context, req view.DeleteRequest) error {
	return h.modelClient.Applications().DeleteOne(req.Model()).Exec(ctx)
}

func (h Handler) Update(ctx *gin.Context, req view.UpdateRequest) error {
	var entity = req.Model()
	return h.modelClient.WithTx(ctx, func(tx *model.Tx) error {
		var updates, err = dao.ApplicationUpdates(tx, entity)
		if err != nil {
			return err
		}
		return updates[0].Exec(ctx)
	})
}

func (h Handler) Get(ctx *gin.Context, req view.GetRequest) (view.GetResponse, error) {
	var entity, err = h.modelClient.Applications().Query().
		Where(application.ID(req.ID)).
		// must extract modules.
		WithModules(func(rq *model.ApplicationModuleRelationshipQuery) {
			rq.Order(model.Asc(applicationmodulerelationship.FieldCreateTime)).
				Select(
					applicationmodulerelationship.FieldApplicationID,
					applicationmodulerelationship.FieldName,
					applicationmodulerelationship.FieldModuleID,
					applicationmodulerelationship.FieldVersion,
					applicationmodulerelationship.FieldAttributes).
				// allow returning without sorting keys.
				Unique(false)
		}).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return model.ExposeApplication(entity), nil
}

// Batch APIs

func (h Handler) CollectionDelete(ctx *gin.Context, req view.CollectionDeleteRequest) error {
	return h.modelClient.WithTx(ctx, func(tx *model.Tx) (err error) {
		for i := range req {
			err = tx.Applications().DeleteOne(req[i].Model()).
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
		application.FieldName,
	}
	getFields = application.WithoutFields(
		application.FieldProjectID,
		application.FieldUpdateTime)
	sortFields = []string{
		application.FieldName,
		application.FieldCreateTime}
)

func (h Handler) CollectionGet(ctx *gin.Context, req view.CollectionGetRequest) (view.CollectionGetResponse, int, error) {
	var query = h.modelClient.Applications().Query().
		Where(application.ProjectIDIn(req.ProjectIDs...))
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
	if orders, ok := req.Sorting(sortFields, model.Desc(application.FieldCreateTime)); ok {
		query.Order(orders...)
	}

	// get application with instances and environments
	entities, err := query.
		// allow returning without sorting keys.
		Unique(false).
		// must extract application instances.
		WithInstances(func(iq *model.ApplicationInstanceQuery) {
			iq.Select(
				applicationinstance.FieldApplicationID,
				applicationinstance.FieldID,
				applicationinstance.FieldName,
				applicationinstance.FieldStatus).
				Where(func(s *sql.Selector) {
					// sq generate instance with row number.
					var sq = s.Clone().
						AppendSelectExprAs(
							sql.RowNumber().
								PartitionBy(applicationinstance.FieldApplicationID).
								OrderBy(sql.Desc(applicationinstance.FieldCreateTime)),
							"row_number",
						).
						Where(s.P()).
						From(s.Table()).
						As(applicationinstance.Table)

					// query latest 5 instances.
					s.Where(sql.LTE(s.C("row_number"), 5)).
						From(sq)
				}).
				Select(
					applicationinstance.FieldEnvironmentID, // must extract environment.
				).
				WithEnvironment(func(eq *model.EnvironmentQuery) {
					eq.Select(environment.FieldName)
				})
		}).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeApplications(entities), cnt, nil
}

// Extensional APIs
