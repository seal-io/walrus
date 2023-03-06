package application

import (
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/rest"

	"github.com/seal-io/seal/pkg/apis/application/view"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/model/applicationmodulerelationship"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/model/applicationrevision"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/platform"
	"github.com/seal-io/seal/pkg/platform/deployer"
	"github.com/seal-io/seal/pkg/platform/operator"
	"github.com/seal-io/seal/pkg/platformtf"
)

func Handle(mc model.ClientSet, kc *rest.Config) Handler {
	return Handler{
		modelClient: mc,
		kubeConfig:  kc,
	}
}

type Handler struct {
	modelClient model.ClientSet
	kubeConfig  *rest.Config
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
		WithModules(func(rq *model.ApplicationModuleRelationshipQuery) {
			rq.Select(
				applicationmodulerelationship.FieldApplicationID,
				applicationmodulerelationship.FieldModuleID,
				applicationmodulerelationship.FieldName,
				applicationmodulerelationship.FieldVariables)
		}).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return model.ExposeApplication(entity), nil
}

// Batch APIs

// Extensional APIs

var (
	resourceGetFields = applicationresource.WithoutFields(
		applicationresource.FieldApplicationID,
		applicationresource.FieldUpdateTime)
	resourceSortFields = []string{
		applicationresource.FieldCreateTime,
		applicationresource.FieldModule,
		applicationresource.FieldMode,
		applicationresource.FieldType,
		applicationresource.FieldName}
)

func (h Handler) GetResources(ctx *gin.Context, req view.GetResourcesRequest) (view.GetResourcesResponse, int, error) {
	var query = h.modelClient.ApplicationResources().Query().
		Where(applicationresource.ApplicationID(req.ID))

	// get count.
	cnt, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// get entities.
	if limit, offset, ok := req.Paging(); ok {
		query.Limit(limit).Offset(offset)
	}
	if fields, ok := req.Extracting(resourceGetFields, resourceGetFields...); ok {
		query.Select(fields...)
	}
	if orders, ok := req.Sorting(resourceSortFields, model.Desc(applicationresource.FieldCreateTime)); ok {
		query.Order(orders...)
	}
	entities, err := query.
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

	var resp = make(view.GetResourcesResponse, len(entities))
	for i := 0; i < len(entities); i++ {
		resp[i].ApplicationResourceOutput = model.ExposeApplicationResource(entities[i])

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

var (
	revisionGetFields  = applicationrevision.Columns
	revisionSortFields = []string{
		applicationrevision.FieldCreateTime,
	}
)

func (h Handler) GetRevisions(ctx *gin.Context, req view.GetRevisionsRequest) (view.GetRevisionsResponse, int, error) {
	var query = h.modelClient.ApplicationRevisions().Query().
		Where(applicationrevision.ApplicationID(req.ID))

	// get count.
	cnt, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// get entities.
	if limit, offset, ok := req.Paging(); ok {
		query.Limit(limit).Offset(offset)
	}
	if fields, ok := req.Extracting(revisionGetFields, revisionGetFields...); ok {
		query.Select(fields...)
	}
	if orders, ok := req.Sorting(revisionSortFields, model.Desc(applicationrevision.FieldCreateTime)); ok {
		query.Order(orders...)
	}
	entities, err := query.
		Unique(false). // allow returning without sorting keys.
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeApplicationRevisions(entities), cnt, nil
}

func (h Handler) CreateDeployments(ctx *gin.Context, req view.CreateDeploymentRequest) error {
	createOpts := deployer.CreateOptions{
		Type:        platformtf.DeployerType,
		ModelClient: h.modelClient,
		KubeConfig:  h.kubeConfig,
	}
	d, err := platform.GetDeployer(ctx, createOpts)
	if err != nil {
		return err
	}

	app, err := h.modelClient.Applications().Get(ctx, req.ID)
	if err != nil {
		return err
	}

	err = d.Apply(ctx, app, deployer.ApplyOptions{
		// TODO get from settings
		SkipTLSVerify: true,
	})
	if err != nil {
		return err
	}

	return nil
}

func (h Handler) DeleteDeployments(ctx *gin.Context, req view.DeleteDeploymentRequest) error {
	deployOpts := deployer.CreateOptions{
		Type:        platformtf.DeployerType,
		ModelClient: h.modelClient,
		KubeConfig:  h.kubeConfig,
	}
	d, err := platform.GetDeployer(ctx, deployOpts)
	if err != nil {
		return err
	}

	app, err := h.modelClient.Applications().Get(ctx, req.ID)
	if err != nil {
		return err
	}

	err = d.Destroy(ctx, app, deployer.DestroyOptions{
		// TODO get from settings
		SkipTLSVerify: true,
	})
	if err != nil {
		return err
	}

	return nil
}
