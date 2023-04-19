package applicationinstance

import (
	"context"

	"github.com/gin-gonic/gin"
	"k8s.io/client-go/rest"

	"github.com/seal-io/seal/pkg/apis/applicationinstance/view"
	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/applicationinstance"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/platform"
	"github.com/seal-io/seal/pkg/platform/deployer"
	"github.com/seal-io/seal/pkg/platformk8s/intercept"
	"github.com/seal-io/seal/pkg/platformtf"
	"github.com/seal-io/seal/pkg/topic/datamessage"
	"github.com/seal-io/seal/utils/log"
	"github.com/seal-io/seal/utils/topic"
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
	return "ApplicationInstance"
}

func (h Handler) Validating() any {
	return h.modelClient
}

// Basic APIs

func (h Handler) Create(ctx *gin.Context, req view.CreateRequest) (resp view.CreateResponse, err error) {
	var entity = req.Model()

	// get deployer.
	var createOpts = deployer.CreateOptions{
		Type:        platformtf.DeployerType,
		ModelClient: h.modelClient,
		KubeConfig:  h.kubeConfig,
	}
	dp, err := platform.GetDeployer(ctx, createOpts)
	if err != nil {
		return nil, err
	}

	// create instance, mark status to deploying.
	entity.Status = status.ApplicationInstanceStatusDeploying
	creates, err := dao.ApplicationInstanceCreates(h.modelClient, entity)
	if err != nil {
		return nil, err
	}
	entity, err = creates[0].Save(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err == nil {
			return
		}
		_ = h.updateInstanceStatus(ctx, entity, status.ApplicationInstanceStatusDeployFailed, err.Error())
	}()

	// apply instance.
	var applyOpts = deployer.ApplyOptions{
		SkipTLSVerify: !h.tlsCertified,
	}
	err = dp.Apply(ctx, entity, applyOpts)
	if err != nil {
		// NB(thxCode): a better approach is to use transaction,
		// however, building the application deployment process is a time-consuming task,
		// to prevent long-time transaction, we use a deletion to achieve this.
		// usually, the probability of this delete operation failing is very low.
		var derr = h.modelClient.ApplicationInstances().DeleteOne(entity).
			Exec(ctx)
		if derr != nil {
			log.WithName("application-instances").
				Errorf("error deleting: %v", derr)
		}
		return nil, err
	}

	if err := publishApplicationUpdate(ctx, entity); err != nil {
		return nil, err
	}

	return model.ExposeApplicationInstance(entity), nil
}

func (h Handler) Delete(ctx *gin.Context, req view.DeleteRequest) (err error) {
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

	if req.Force != nil && !*req.Force {
		// do not clean deployed native resources.
		err = h.modelClient.ApplicationInstances().DeleteOne(entity).
			Exec(ctx)
		if err != nil {
			return err
		}

		return publishApplicationUpdate(ctx, entity)
	}

	// mark status to deleting.
	entity, err = h.modelClient.ApplicationInstances().UpdateOne(entity).
		SetStatus(status.ApplicationInstanceStatusDeleting).
		SetStatusMessage("").
		Save(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err == nil {
			return
		}
		_ = h.updateInstanceStatus(ctx, entity, status.ApplicationInstanceStatusDeleteFailed, err.Error())
	}()

	if err := publishApplicationUpdate(ctx, entity); err != nil {
		return err
	}

	// destroy instance.
	var destroyOpts = deployer.DestroyOptions{
		SkipTLSVerify: !h.tlsCertified,
	}
	return dp.Destroy(ctx, entity, destroyOpts)
}

func (h Handler) Get(ctx *gin.Context, req view.GetRequest) (view.GetResponse, error) {
	return h.getEntityOutput(ctx, req.ID)
}

func (h Handler) getEntityOutput(ctx context.Context, id types.ID) (*model.ApplicationInstanceOutput, error) {
	var entity, err = h.modelClient.ApplicationInstances().Query().
		Where(applicationinstance.ID(id)).
		WithEnvironment(func(eq *model.EnvironmentQuery) {
			eq.Select(environment.FieldName)
		}).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return model.ExposeApplicationInstance(entity), nil
}

func (h Handler) Stream(ctx runtime.RequestStream, req view.StreamRequest) error {
	var t, err = topic.Subscribe(datamessage.ApplicationInstance)
	if err != nil {
		return err
	}

	defer func() { t.Unsubscribe() }()
	for {
		var event topic.Event
		event, err = t.Receive(ctx)
		if err != nil {
			return err
		}
		dm, ok := event.Data.(datamessage.Message)
		if !ok {
			continue
		}

		var streamData view.StreamResponse
		for _, id := range dm.Data {
			if id != req.ID {
				continue
			}

			switch dm.Type {
			case datamessage.EventCreate, datamessage.EventUpdate:
				entityOutput, err := h.getEntityOutput(ctx, id)
				if err != nil {
					return err
				}
				streamData = view.StreamResponse{
					Type:       dm.Type,
					Collection: []*model.ApplicationInstanceOutput{entityOutput},
				}
			case datamessage.EventDelete:
				streamData = view.StreamResponse{
					Type: dm.Type,
					IDs:  dm.Data,
				}
			}
		}
		err = ctx.SendJSON(streamData)
		if err != nil {
			return err
		}
	}
}

// Batch APIs

var (
	queryFields = []string{
		applicationinstance.FieldName,
	}
	getFields = applicationinstance.WithoutFields(
		applicationinstance.FieldApplicationID,
		applicationinstance.FieldUpdateTime)
	sortFields = []string{
		applicationinstance.FieldName,
		applicationinstance.FieldCreateTime}
)

func (h Handler) CollectionGet(ctx *gin.Context, req view.CollectionGetRequest) (view.CollectionGetResponse, int, error) {
	var query = h.modelClient.ApplicationInstances().Query().
		Where(applicationinstance.ApplicationID(req.ApplicationID))
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
	if orders, ok := req.Sorting(sortFields, model.Desc(applicationinstance.FieldCreateTime)); ok {
		query.Order(orders...)
	}
	entities, err := query.
		// allow returning without sorting keys.
		Unique(false).
		// must extract environment.
		Select(applicationinstance.FieldEnvironmentID).
		WithEnvironment(func(eq *model.EnvironmentQuery) {
			eq.Select(environment.FieldName)
		}).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeApplicationInstances(entities), cnt, nil
}

func (h Handler) CollectionStream(ctx runtime.RequestStream, req view.CollectionStreamRequest) error {
	var t, err = topic.Subscribe(datamessage.ApplicationInstance)
	if err != nil {
		return err
	}

	var query = h.modelClient.ApplicationInstances().Query()
	if req.ApplicationID != "" {
		query.Where(applicationinstance.ApplicationID(req.ApplicationID))
	}
	if fields, ok := req.Extracting(getFields, getFields...); ok {
		query.Select(fields...)
	}

	defer func() { t.Unsubscribe() }()
	for {
		var event topic.Event
		event, err = t.Receive(ctx)
		if err != nil {
			return err
		}
		dm, ok := event.Data.(datamessage.Message)
		if !ok {
			continue
		}

		var streamData view.StreamResponse
		switch dm.Type {
		case datamessage.EventCreate, datamessage.EventUpdate:
			entities, err := query.Clone().
				// allow returning without sorting keys.
				Unique(false).
				// must extract environment.
				Select(applicationinstance.FieldEnvironmentID).
				Where(applicationinstance.IDIn(dm.Data...)).
				WithEnvironment(func(eq *model.EnvironmentQuery) {
					eq.Select(environment.FieldName)
				}).
				All(ctx)

			if err != nil {
				return err
			}
			streamData = view.StreamResponse{
				Type:       dm.Type,
				Collection: model.ExposeApplicationInstances(entities),
			}
		case datamessage.EventDelete:
			streamData = view.StreamResponse{
				Type: dm.Type,
				IDs:  dm.Data,
			}
		}
		if len(streamData.IDs) == 0 && len(streamData.Collection) == 0 {
			continue
		}
		err = ctx.SendJSON(streamData)
		if err != nil {
			return err
		}
	}
}

// Extensional APIs

func (h Handler) RouteUpgrade(ctx *gin.Context, req view.RouteUpgradeRequest) (err error) {
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

	// update instance, mark status to deploying.
	entity, err = h.modelClient.ApplicationInstances().UpdateOne(entity).
		SetVariables(entity.Variables).
		SetStatus(status.ApplicationInstanceStatusDeploying).
		SetStatusMessage("").
		Save(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err == nil {
			return
		}
		_ = h.updateInstanceStatus(ctx, entity, status.ApplicationInstanceStatusDeployFailed, err.Error())
	}()

	if err := publishApplicationUpdate(ctx, entity); err != nil {
		return err
	}

	// apply instance.
	var applyOpts = deployer.ApplyOptions{
		SkipTLSVerify: !h.tlsCertified,
	}
	return dp.Apply(ctx, entity, applyOpts)
}

func (h Handler) RouteAccessEndpoints(ctx *gin.Context, req view.AccessEndpointRequest) (*view.AccessEndpointResponse, error) {
	res, err := h.modelClient.ApplicationResources().Query().
		Where(
			applicationresource.InstanceID(req.ID),
			applicationresource.TypeIn(
				intercept.TFEndpointsTypes...,
			),
		).
		Select(
			applicationresource.FieldConnectorID,
			applicationresource.FieldType,
			applicationresource.FieldName,
			applicationresource.FieldStatus,
		).
		All(ctx)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}

	var endpoints []view.ResourceEndpoint
	for _, v := range res {
		for _, eps := range v.Status.ResourceEndpoints {
			endpoints = append(endpoints, view.ResourceEndpoint{
				ResourceID:      v.Name,
				ResourceKind:    v.Type,
				ResourceSubKind: eps.EndpointType,
				Endpoints:       eps.Endpoints,
			})
		}
	}
	return &view.AccessEndpointResponse{
		Endpoints: endpoints,
	}, nil
}

func (h Handler) updateInstanceStatus(ctx context.Context, entity *model.ApplicationInstance, s, m string) error {
	var logger = log.WithName("application-instance")

	entity.Status = s
	entity.StatusMessage = m
	update, err := dao.ApplicationInstanceUpdate(h.modelClient, entity)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = update.Exec(ctx)
	if err != nil && !model.IsNotFound(err) {
		logger.Errorf("failed to update status of instance %s: %v", entity.ID, err)
		return err
	}

	return nil
}

func publishApplicationUpdate(ctx context.Context, entity *model.ApplicationInstance) error {
	return datamessage.Publish(ctx, string(datamessage.Application), model.OpUpdate, []types.ID{entity.ApplicationID})
}
