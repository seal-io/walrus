package applicationinstance

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zclconf/go-cty/cty"
	"k8s.io/client-go/rest"

	"github.com/seal-io/seal/pkg/apis/applicationinstance/view"
	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/model/applicationinstance"
	"github.com/seal-io/seal/pkg/dao/model/applicationmodulerelationship"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/model/applicationrevision"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/pkg/dao/types/property"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/platform"
	"github.com/seal-io/seal/pkg/platform/deployer"
	"github.com/seal-io/seal/pkg/platformk8s/intercept"
	"github.com/seal-io/seal/pkg/platformtf"
	"github.com/seal-io/seal/pkg/topic/datamessage"
	"github.com/seal-io/seal/utils/log"
	"github.com/seal-io/seal/utils/strs"
	"github.com/seal-io/seal/utils/topic"
	"github.com/seal-io/seal/utils/validation"
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

type createInstanceOptions struct {
	Clone               bool
	Tags                []string
	ApplicationInstance *model.ApplicationInstance
}

func (h Handler) Kind() string {
	return "ApplicationInstance"
}

func (h Handler) Validating() any {
	return h.modelClient
}

// Basic APIs.

func (h Handler) Create(ctx *gin.Context, req view.CreateRequest) (resp view.CreateResponse, err error) {
	entity := req.Model()

	return h.createInstance(ctx, createInstanceOptions{
		Tags:                req.RemarkTags,
		ApplicationInstance: entity,
	})
}

func (h Handler) Delete(ctx *gin.Context, req view.DeleteRequest) (err error) {
	logger := log.WithName("api").WithName("application-instance")
	entity := req.Model()

	// Get deployer.
	createOpts := deployer.CreateOptions{
		Type:        platformtf.DeployerType,
		ModelClient: h.modelClient,
		KubeConfig:  h.kubeConfig,
	}

	dp, err := platform.GetDeployer(ctx, createOpts)
	if err != nil {
		return err
	}

	if req.Force != nil && !*req.Force {
		// Get application instance with application id.
		entity, err = h.modelClient.ApplicationInstances().Query().
			Select(applicationinstance.FieldID, applicationinstance.FieldApplicationID).
			Where(applicationinstance.IDEQ(entity.ID)).
			Only(ctx)
		if err != nil {
			return err
		}

		// Do not clean deployed native resources.
		err = h.modelClient.ApplicationInstances().DeleteOne(entity).
			Exec(ctx)
		if err != nil {
			return err
		}

		return publishApplicationUpdate(ctx, entity)
	}

	// Mark status to deleting.
	status.ApplicationInstanceStatusDeleted.Reset(entity, "Deleting")

	update, err := dao.ApplicationInstanceUpdate(h.modelClient, entity)
	if err != nil {
		return err
	}

	entity, err = update.Save(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err == nil {
			return
		}
		// Update a failure status.
		status.ApplicationInstanceStatusDeleted.False(entity, err.Error())

		uerr := h.updateInstanceStatus(ctx, entity)
		if uerr != nil {
			logger.Errorf("error updating status of instance %s: %v",
				entity.ID, uerr)
		}
	}()

	if err = publishApplicationUpdate(ctx, entity); err != nil {
		return err
	}

	// Destroy instance.
	destroyOpts := deployer.DestroyOptions{
		SkipTLSVerify: !h.tlsCertified,
	}

	return dp.Destroy(ctx, entity, destroyOpts)
}

func (h Handler) Get(ctx *gin.Context, req view.GetRequest) (view.GetResponse, error) {
	return h.getEntityOutput(ctx, req.ID)
}

func (h Handler) getEntityOutput(ctx context.Context, id oid.ID) (*model.ApplicationInstanceOutput, error) {
	entity, err := h.modelClient.ApplicationInstances().Query().
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

func (h Handler) Stream(ctx runtime.RequestUnidiStream, req view.StreamRequest) error {
	t, err := topic.Subscribe(datamessage.ApplicationInstance)
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

		dm, ok := event.Data.(datamessage.Message[oid.ID])
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

// Batch APIs.

var (
	queryFields = []string{
		applicationinstance.FieldName,
	}
	getFields = applicationinstance.WithoutFields(
		applicationinstance.FieldApplicationID,
		applicationinstance.FieldUpdateTime)
	sortFields = []string{
		applicationinstance.FieldName,
		applicationinstance.FieldCreateTime,
	}
)

func (h Handler) CollectionGet(
	ctx *gin.Context,
	req view.CollectionGetRequest,
) (view.CollectionGetResponse, int, error) {
	query := h.modelClient.ApplicationInstances().Query().
		Where(applicationinstance.ApplicationID(req.ApplicationID))
	if queries, ok := req.Querying(queryFields); ok {
		query.Where(queries)
	}

	// Get count.
	cnt, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Get entities.
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
		// Allow returning without sorting keys.
		Unique(false).
		// Must extract environment.
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

func (h Handler) CollectionStream(ctx runtime.RequestUnidiStream, req view.CollectionStreamRequest) error {
	t, err := topic.Subscribe(datamessage.ApplicationInstance)
	if err != nil {
		return err
	}

	defer func() { t.Unsubscribe() }()

	query := h.modelClient.ApplicationInstances().Query()
	if req.ApplicationID != "" {
		query.Where(applicationinstance.ApplicationID(req.ApplicationID))
	}

	if fields, ok := req.Extracting(getFields, getFields...); ok {
		query.Select(fields...)
	}

	for {
		var event topic.Event

		event, err = t.Receive(ctx)
		if err != nil {
			return err
		}

		dm, ok := event.Data.(datamessage.Message[oid.ID])
		if !ok {
			continue
		}

		var streamData view.StreamResponse

		switch dm.Type {
		case datamessage.EventCreate, datamessage.EventUpdate:
			entities, err := query.Clone().
				// Allow returning without sorting keys.
				Unique(false).
				// Must extract environment.
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

// Extensional APIs.

func (h Handler) RouteUpgrade(ctx *gin.Context, req view.RouteUpgradeRequest) (err error) {
	logger := log.WithName("api").WithName("application-instance")
	entity := req.Model()

	// Get deployer.
	createOpts := deployer.CreateOptions{
		Type:        platformtf.DeployerType,
		ModelClient: h.modelClient,
		KubeConfig:  h.kubeConfig,
	}

	dp, err := platform.GetDeployer(ctx, createOpts)
	if err != nil {
		return err
	}

	// Update instance, mark status from deploying.
	entity.Variables = req.Variables
	status.ApplicationInstanceStatusDeployed.Reset(entity, "Upgrading")

	update, err := dao.ApplicationInstanceUpdate(h.modelClient, entity)
	if err != nil {
		return err
	}

	entity, err = update.Save(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err == nil {
			return
		}
		// Update a failure status.
		status.ApplicationInstanceStatusDeployed.False(entity, err.Error())

		uerr := h.updateInstanceStatus(ctx, entity)
		if uerr != nil {
			logger.Errorf("error updating status of instance %s: %v",
				entity.ID, uerr)
		}
	}()

	if err = publishApplicationUpdate(ctx, entity); err != nil {
		return err
	}

	// Apply instance.
	applyOpts := deployer.ApplyOptions{
		SkipTLSVerify: !h.tlsCertified,
		Tags:          req.RemarkTags,
	}

	return dp.Apply(ctx, entity, applyOpts)
}

func (h Handler) RouteAccessEndpoints(
	ctx *gin.Context,
	req view.AccessEndpointRequest,
) (view.AccessEndpointResponse, error) {
	return h.accessEndpoints(ctx, req.ID)
}

func (h Handler) accessEndpoints(ctx context.Context, instanceID oid.ID) (view.AccessEndpointResponse, error) {
	// Endpoints from output.
	endpoints, err := h.endpointsFromOutput(ctx, instanceID)
	if err != nil {
		return nil, err
	}

	if len(endpoints) != 0 {
		return endpoints, nil
	}

	// Endpoints from resources.
	return h.endpointsFromResources(ctx, instanceID)
}

func (h Handler) endpointsFromOutput(ctx context.Context, instanceID oid.ID) (view.AccessEndpointResponse, error) {
	outputs, err := h.getInstanceOutputs(ctx, instanceID, true)
	if err != nil {
		return nil, err
	}

	var (
		invalidTypeErr = runtime.Error(http.StatusBadRequest,
			"element type of output endpoints should be string")
		endpoints = make([]view.Endpoint, 0, len(outputs))
	)

	for _, v := range outputs {
		if !view.IsEndpointOuput(v.Name) {
			continue
		}

		prop := property.Property{
			Type:  v.Type,
			Value: v.Value,
		}

		switch {
		case v.Type == cty.String:
			ep, _, err := prop.GetString()
			if err != nil {
				return nil, err
			}

			if err := validation.IsValidEndpoint(ep); err != nil {
				return nil, runtime.Error(http.StatusBadRequest, err)
			}

			endpoints = append(endpoints, view.Endpoint{
				ModuleName: v.ModuleName,
				Endpoints:  []string{ep},
				Name:       v.Name,
			})
		case v.Type.IsListType() || v.Type.IsSetType() || v.Type.IsTupleType():
			if v.Type.IsTupleType() {
				// For tuple: each element has its own type.
				for _, tp := range v.Type.TupleElementTypes() {
					if tp != cty.String {
						return nil, invalidTypeErr
					}
				}
			} else if v.Type.ElementType() != cty.String {
				// For list and set: all elements are the same type.
				return nil, invalidTypeErr
			}

			eps, _, err := property.GetSlice[string](prop)
			if err != nil {
				return nil, err
			}

			if err := validation.IsValidEndpoints(eps); err != nil {
				return nil, runtime.Error(http.StatusBadRequest, err)
			}

			endpoints = append(endpoints, view.Endpoint{
				ModuleName: v.ModuleName,
				Endpoints:  eps,
				Name:       v.Name,
			})
		}
	}

	return endpoints, nil
}

func (h Handler) endpointsFromResources(ctx context.Context, instanceID oid.ID) ([]view.Endpoint, error) {
	res, err := h.modelClient.ApplicationResources().Query().
		Where(
			applicationresource.InstanceID(instanceID),
			applicationresource.TypeIn(
				intercept.TFEndpointsTypes...,
			),
		).
		Select(
			applicationresource.FieldConnectorID,
			applicationresource.FieldType,
			applicationresource.FieldName,
			applicationresource.FieldStatus,
			applicationresource.FieldModule,
		).
		All(ctx)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}

	getModuleName := func(origin string) string {
		if arr := strings.Split(origin, "/"); len(arr) >= 1 {
			return arr[0]
		}

		return ""
	}

	var endpoints []view.Endpoint

	for _, v := range res {
		mn := getModuleName(v.Module)

		for _, eps := range v.Status.ResourceEndpoints {
			endpoints = append(endpoints, view.Endpoint{
				ModuleName: mn,
				Name:       strs.Join("/", eps.EndpointType, v.Name),
				Endpoints:  eps.Endpoints,
			})
		}
	}

	return endpoints, nil
}

func (h Handler) RouteOutputs(ctx *gin.Context, req view.OutputRequest) (view.OutputResponse, error) {
	return h.getInstanceOutputs(ctx, req.ID, true)
}

func (h Handler) getInstanceOutputs(
	ctx context.Context,
	instanceID oid.ID,
	onlySuccess bool,
) ([]types.OutputValue, error) {
	ar, err := h.modelClient.ApplicationRevisions().Query().
		Where(applicationrevision.InstanceID(instanceID)).
		Select(
			applicationrevision.FieldOutput,
			applicationrevision.FieldModules,
			applicationrevision.FieldStatus,
		).
		Order(model.Desc(applicationrevision.FieldCreateTime)).
		First(ctx)
	if err != nil {
		return nil, fmt.Errorf("error get latest application revision")
	}

	if onlySuccess && ar.Status != status.ApplicationRevisionStatusSucceeded {
		return nil, nil
	}

	o, err := platformtf.ParseStateOutput(ar)
	if err != nil {
		return nil, fmt.Errorf("error get outputs: %w", err)
	}

	return o, nil
}

func (h Handler) updateInstanceStatus(
	ctx context.Context,
	entity *model.ApplicationInstance,
) error {
	update, err := dao.ApplicationInstanceStatusUpdate(h.modelClient, entity)
	if err != nil {
		return err
	}

	err = update.Exec(ctx)
	if err != nil && !model.IsNotFound(err) {
		return err
	}

	return nil
}

// CreateClone creates a clone instance of the application instance.
func (h Handler) CreateClone(
	ctx *gin.Context,
	req view.CreateCloneRequest,
) (*model.ApplicationInstanceOutput, error) {
	applicationInstance, err := h.modelClient.ApplicationInstances().Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	applicationInstance.Name = req.Name

	if req.EnviornmentID != "" {
		applicationInstance.EnvironmentID = req.EnviornmentID
	}

	return h.createInstance(ctx, createInstanceOptions{
		Clone:               true,
		Tags:                req.RemarkTags,
		ApplicationInstance: applicationInstance,
	})
}

func (h Handler) createInstance(
	ctx context.Context,
	opts createInstanceOptions,
) (*model.ApplicationInstanceOutput, error) {
	logger := log.WithName("api").WithName("application-instance")

	// Get deployer.
	createOpts := deployer.CreateOptions{
		Type:        platformtf.DeployerType,
		ModelClient: h.modelClient,
		KubeConfig:  h.kubeConfig,
	}

	dp, err := platform.GetDeployer(ctx, createOpts)
	if err != nil {
		return nil, err
	}

	// Create instance, mark status to deploying.
	creates, err := dao.ApplicationInstanceCreates(h.modelClient, opts.ApplicationInstance)
	if err != nil {
		return nil, err
	}

	entity, err := creates[0].Save(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err == nil {
			return
		}
		// Update a failure status.
		status.ApplicationInstanceStatusDeployed.False(entity, err.Error())

		uerr := h.updateInstanceStatus(ctx, entity)
		if uerr != nil {
			logger.Errorf("error updating status of instance %s: %v",
				entity.ID, uerr)
		}
	}()

	// ClonedInstanceRevision is the latest application revision
	// of the cloned application instance.
	var clonedInstanceRevision *model.ApplicationRevision

	if opts.Clone {
		if opts.ApplicationInstance.ID == "" {
			return nil, errors.New("application instance id is empty")
		}

		clonedInstanceRevision, err = h.modelClient.ApplicationRevisions().Query().
			Where(applicationrevision.InstanceID(opts.ApplicationInstance.ID)).
			Order(model.Desc(applicationrevision.FieldCreateTime)).
			First(ctx)
		if err != nil && !model.IsNotFound(err) {
			return nil, err
		}
	}

	// Apply instance.
	applyOpts := deployer.ApplyOptions{
		SkipTLSVerify: !h.tlsCertified,
		Tags:          opts.Tags,
		CloneFrom:     clonedInstanceRevision,
	}

	err = dp.Apply(ctx, entity, applyOpts)
	if err != nil {
		// NB(thxCode): a better approach is to use transaction,
		// however, building the application deployment process is a time-consuming task,
		// to prevent long-time transaction, we use a deletion to achieve this.
		// Usually, the probability of this delete operation failing is very low.
		derr := h.modelClient.ApplicationInstances().DeleteOne(entity).
			Exec(ctx)
		if derr != nil {
			logger.Errorf("error deleting: %v", derr)
		}

		return nil, err
	}

	if err := publishApplicationUpdate(ctx, entity); err != nil {
		return nil, err
	}

	return model.ExposeApplicationInstance(entity), nil
}

func publishApplicationUpdate(ctx context.Context, entity *model.ApplicationInstance) error {
	return datamessage.Publish(ctx, string(datamessage.Application), model.OpUpdate, []oid.ID{entity.ApplicationID})
}

func (h Handler) StreamAccessEndpoint(ctx runtime.RequestUnidiStream, req view.GetRequest) error {
	t, err := topic.Subscribe(datamessage.ApplicationRevision)
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

		dm, ok := event.Data.(datamessage.Message[oid.ID])
		if !ok {
			continue
		}

		var streamData view.StreamAccessEndpointResponse

		for _, id := range dm.Data {
			ar, err := h.getRevisionByID(ctx, id)
			if err != nil {
				return err
			}

			if ar.InstanceID != req.ID {
				continue
			}

			eps, err := h.accessEndpoints(ctx, req.ID)
			if err != nil {
				return err
			}

			if len(eps) == 0 {
				continue
			}

			switch dm.Type {
			case datamessage.EventCreate:
				// While create new application revision,
				// the previous endpoints from outputs and resources need to be deleted.
				streamData = view.StreamAccessEndpointResponse{
					Type:       datamessage.EventDelete,
					Collection: eps,
				}
			case datamessage.EventUpdate:
				// While the application revision status is succeeded, the endpoints is updated to the current revision.
				if ar.Status != status.ApplicationRevisionStatusSucceeded {
					continue
				}
				streamData = view.StreamAccessEndpointResponse{
					Type:       datamessage.EventUpdate,
					Collection: eps,
				}
			}

			err = ctx.SendJSON(streamData)
			if err != nil {
				return err
			}
		}
	}
}

func (h Handler) getRevisionByID(ctx context.Context, id oid.ID) (*model.ApplicationRevision, error) {
	return h.modelClient.ApplicationRevisions().Query().
		Where(applicationrevision.ID(id)).
		Only(ctx)
}

func (h Handler) StreamOutput(ctx runtime.RequestUnidiStream, req view.GetRequest) error {
	t, err := topic.Subscribe(datamessage.ApplicationRevision)
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

		dm, ok := event.Data.(datamessage.Message[oid.ID])
		if !ok {
			continue
		}

		var streamData view.StreamOutputResponse

		for _, id := range dm.Data {
			ar, err := h.getRevisionByID(ctx, id)
			if err != nil {
				return err
			}

			if ar.InstanceID != req.ID {
				continue
			}

			op, err := h.getInstanceOutputs(ctx, ar.InstanceID, false)
			if err != nil {
				return err
			}

			if len(op) == 0 {
				continue
			}

			switch dm.Type {
			case datamessage.EventCreate:
				// While create new application revision, the outputs of new revision is the previous outputs.
				streamData = view.StreamOutputResponse{
					Type:       datamessage.EventDelete,
					Collection: op,
				}
			case datamessage.EventUpdate:
				// While the application revision status is succeeded, the outputs is updated to the current revision.
				if ar.Status != status.ApplicationRevisionStatusSucceeded {
					continue
				}
				streamData = view.StreamOutputResponse{
					Type:       datamessage.EventUpdate,
					Collection: op,
				}
			}

			err = ctx.SendJSON(streamData)
			if err != nil {
				return err
			}
		}
	}
}

// GetDiffLatest diff the latest application config with the current application instance.
func (h Handler) GetDiffLatest(ctx *gin.Context, req view.DiffLatestRequest) (*view.DiffLatestResponse, error) {
	latestRevision, err := h.modelClient.ApplicationRevisions().Query().
		Select(
			applicationrevision.FieldModules,
			applicationrevision.FieldVariables,
		).
		Where(applicationrevision.InstanceID(req.ID)).
		Order(model.Desc(applicationrevision.FieldCreateTime)).
		First(ctx)
	if err != nil {
		return nil, err
	}

	app, err := h.modelClient.ApplicationInstances().Query().
		Where(applicationinstance.ID(req.ID)).
		QueryApplication().
		Select(
			application.FieldID,
			application.FieldVariables,
		).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	relationships, err := h.modelClient.ApplicationModuleRelationships().Query().
		Where(applicationmodulerelationship.ApplicationID(app.ID)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	modules := make([]types.ApplicationModule, 0, len(relationships))
	for _, r := range relationships {
		modules = append(modules, types.ApplicationModule{
			ModuleID:   r.ModuleID,
			Version:    r.Version,
			Name:       r.Name,
			Attributes: r.Attributes,
		})
	}

	return &view.DiffLatestResponse{
		Old: &view.Diff{
			Modules:   latestRevision.Modules,
			Variables: latestRevision.Variables,
		},
		New: &view.Diff{
			Modules:   modules,
			Variables: app.Variables,
		},
	}, nil
}
