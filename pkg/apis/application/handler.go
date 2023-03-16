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
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/platform"
	"github.com/seal-io/seal/pkg/platform/deployer"
	"github.com/seal-io/seal/pkg/platformtf"
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
		Where(application.ProjectID(req.ProjectID))
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
	entities, err := query.
		// allow returning without sorting keys.
		Unique(false).
		// must extract application instances.
		WithInstances(func(iq *model.ApplicationInstanceQuery) {
			iq.Order(model.Asc(applicationinstance.FieldCreateTime)).
				// earlier 5 instances.
				Limit(5).
				Select(
					applicationinstance.FieldApplicationID,
					applicationinstance.FieldID,
					applicationinstance.FieldName,
					applicationinstance.FieldStatus).
				// allow returning without sorting keys.
				Unique(false).
				// must extract environment.
				Select(applicationinstance.FieldEnvironmentID).
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

func (h Handler) RouteDeploy(ctx *gin.Context, req view.RouteDeployRequest) error {
	var entity = req.Model()

	// get deployer.
	var createOpts = deployer.CreateOptions{
		Type:        platformtf.DeployerType,
		ModelClient: h.modelClient,
		KubeConfig:  h.kubeConfig,
	}
	dp, err := platform.GetDeployer(ctx, createOpts)
	if err != nil {
		return err
	}

	// create/update instance, mark status to deploying.
	entity.ID, err = h.modelClient.ApplicationInstances().Create().
		SetApplicationID(entity.ApplicationID).
		SetEnvironmentID(entity.EnvironmentID).
		SetName(entity.Name).
		SetVariables(entity.Variables).
		SetStatus(status.ApplicationInstanceStatusDeploying).
		SetStatusMessage("").
		OnConflict(
			sql.ConflictColumns(
				applicationinstance.FieldApplicationID,
				applicationinstance.FieldEnvironmentID,
				applicationinstance.FieldName),
		).
		Update(func(upsert *model.ApplicationInstanceUpsert) {
			if entity.Variables != nil {
				upsert.UpdateVariables()
			}
			upsert.UpdateStatus()
			upsert.UpdateStatusMessage()
			upsert.UpdateUpdateTime()
		}).
		ID(ctx)
	if err != nil {
		return err
	}

	// apply instance.
	var applyOpts = deployer.ApplyOptions{
		SkipTLSVerify: !h.tlsCertified,
	}
	return dp.Apply(ctx, entity, applyOpts)
}
