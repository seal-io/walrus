package service

import (
	"context"
	"fmt"
	"net/http"

	"entgo.io/ent/dialect/sql"
	"github.com/gin-gonic/gin"
	"github.com/zclconf/go-cty/cty"
	"k8s.io/client-go/rest"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/apis/service/view"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/project"
	"github.com/seal-io/seal/pkg/dao/model/service"
	"github.com/seal-io/seal/pkg/dao/model/servicerelationship"
	"github.com/seal-io/seal/pkg/dao/model/serviceresource"
	"github.com/seal-io/seal/pkg/dao/model/servicerevision"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/pkg/dao/types/property"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/deployer"
	deployertf "github.com/seal-io/seal/pkg/deployer/terraform"
	deptypes "github.com/seal-io/seal/pkg/deployer/types"
	"github.com/seal-io/seal/pkg/operator"
	"github.com/seal-io/seal/pkg/operator/k8s/intercept"
	optypes "github.com/seal-io/seal/pkg/operator/types"
	pkgservice "github.com/seal-io/seal/pkg/service"
	tfparser "github.com/seal-io/seal/pkg/terraform/parser"
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

type createServiceOptions struct {
	Tags    []string
	Service *model.Service
}

func (h Handler) Kind() string {
	return "Service"
}

func (h Handler) Validating() any {
	return h.modelClient
}

// Basic APIs.

func (h Handler) Create(ctx *gin.Context, req view.CreateRequest) (resp view.CreateResponse, err error) {
	entity := req.Model()

	dp, err := h.getDeployer(ctx)
	if err != nil {
		return nil, err
	}

	createOpts := pkgservice.Options{
		TlsCertified: h.tlsCertified,
		Tags:         req.RemarkTags,
	}

	return pkgservice.Create(ctx,
		h.modelClient,
		dp,
		entity,
		createOpts,
	)
}

func (h Handler) Delete(ctx *gin.Context, req view.DeleteRequest) (err error) {
	entity := req.Model()
	if req.Force != nil && !*req.Force {
		// Do not clean deployed native resources.
		return h.modelClient.Services().DeleteOne(entity).
			Exec(ctx)
	}

	// Mark status to deleting.
	status.ServiceStatusDeleted.Reset(entity, "Deleting")

	update, err := dao.ServiceUpdate(h.modelClient, entity)
	if err != nil {
		return err
	}

	entity, err = update.Save(ctx)
	if err != nil {
		return err
	}

	dp, err := h.getDeployer(ctx)
	if err != nil {
		return err
	}

	destroyOpts := pkgservice.Options{
		TlsCertified: h.tlsCertified,
	}

	return pkgservice.Destroy(ctx, h.modelClient, dp, entity, destroyOpts)
}

func (h Handler) Get(ctx *gin.Context, req view.GetRequest) (view.GetResponse, error) {
	return h.getEntityOutput(ctx, req.ID)
}

func (h Handler) getEntityOutput(ctx context.Context, id oid.ID) (*model.ServiceOutput, error) {
	entity, err := h.modelClient.Services().Query().
		Where(service.ID(id)).
		WithProject(func(pq *model.ProjectQuery) {
			pq.Select(project.FieldName)
		}).
		WithEnvironment(func(eq *model.EnvironmentQuery) {
			eq.Select(environment.FieldName)
		}).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return model.ExposeService(entity), nil
}

func (h Handler) Stream(ctx runtime.RequestUnidiStream, req view.StreamRequest) error {
	t, err := topic.Subscribe(datamessage.Service)
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
					Collection: []*model.ServiceOutput{entityOutput},
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
		service.FieldName,
	}
	getFields = service.WithoutFields(
		service.FieldEnvironmentID,
		service.FieldUpdateTime)
	sortFields = []string{
		service.FieldName,
		service.FieldCreateTime,
	}
)

func (h Handler) CollectionGet(
	ctx *gin.Context,
	req view.CollectionGetRequest,
) (view.CollectionGetResponse, int, error) {
	query := h.modelClient.Services().Query().
		Where(service.EnvironmentID(req.EnvironmentID))
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

	if orders, ok := req.Sorting(sortFields, model.Desc(service.FieldCreateTime)); ok {
		query.Order(orders...)
	}

	entities, err := query.
		// Allow returning without sorting keys.
		Unique(false).
		// Must extract environment.
		Select(service.FieldEnvironmentID).
		WithEnvironment(func(eq *model.EnvironmentQuery) {
			eq.Select(environment.FieldName)
		}).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeServices(entities), cnt, nil
}

func (h Handler) CollectionCreate(ctx *gin.Context, req view.CollectionCreateRequest) error {
	return h.modelClient.WithTx(ctx, func(tx *model.Tx) error {
		for _, envID := range req.EnvironmentIDs {
			var services model.Services

			for _, s := range req.Services {
				svc := s.Model()
				svc.EnvironmentID = envID
				services = append(services, svc)
			}

			_, err := pkgservice.CreateScheduledServices(ctx, tx, services)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (h Handler) CollectionStream(ctx runtime.RequestUnidiStream, req view.CollectionStreamRequest) error {
	t, err := topic.Subscribe(datamessage.Service)
	if err != nil {
		return err
	}

	defer func() { t.Unsubscribe() }()

	query := h.modelClient.Services().Query()
	if req.EnvironmentID != "" {
		query.Where(service.EnvironmentID(req.EnvironmentID))
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
				Select(service.FieldEnvironmentID).
				Where(service.IDIn(dm.Data...)).
				WithEnvironment(func(eq *model.EnvironmentQuery) {
					eq.Select(environment.FieldName)
				}).
				All(ctx)
			if err != nil {
				return err
			}
			streamData = view.StreamResponse{
				Type:       dm.Type,
				Collection: model.ExposeServices(entities),
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
	entity := req.Model()
	// Update service, mark status from deploying.
	status.ServiceStatusDeployed.Reset(entity, "Upgrading")

	err = h.modelClient.WithTx(ctx, func(tx *model.Tx) error {
		update, err := dao.ServiceUpdate(tx, entity)
		if err != nil {
			return err
		}

		entity, err = update.Save(ctx)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	dp, err := h.getDeployer(ctx)
	if err != nil {
		return err
	}

	// Apply service.
	applyOpts := pkgservice.Options{
		TlsCertified: h.tlsCertified,
		Tags:         req.RemarkTags,
	}

	return pkgservice.Apply(ctx, h.modelClient, dp, entity, applyOpts)
}

// RouteRollback rolls back a service to a specific revision.
func (h Handler) RouteRollback(ctx *gin.Context, req view.RouteRollbackRequest) error {
	service, err := h.modelClient.Services().Get(ctx, req.ID)
	if err != nil {
		return err
	}

	serviceRevision, err := h.modelClient.ServiceRevisions().Get(ctx, req.RevisionID)
	if err != nil {
		return err
	}

	service.Template.Version = serviceRevision.TemplateVersion
	service.Attributes = serviceRevision.Attributes
	status.ServiceStatusDeployed.Reset(service, "Rolling back")

	err = h.modelClient.WithTx(ctx, func(tx *model.Tx) error {
		update, err := dao.ServiceUpdate(tx, service)
		if err != nil {
			return err
		}

		service, err = update.Save(ctx)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	dp, err := h.getDeployer(ctx)
	if err != nil {
		return err
	}

	applyOpts := pkgservice.Options{
		TlsCertified: h.tlsCertified,
	}

	return pkgservice.Apply(ctx, h.modelClient, dp, service, applyOpts)
}

func (h Handler) RouteAccessEndpoints(
	ctx *gin.Context,
	req view.RouteAccessEndpointRequest,
) (view.RouteAccessEndpointResponse, error) {
	return h.accessEndpoints(ctx, req.ID)
}

func (h Handler) accessEndpoints(ctx context.Context, serviceID oid.ID) (view.RouteAccessEndpointResponse, error) {
	// Endpoints from output.
	endpoints, err := h.endpointsFromOutput(ctx, serviceID)
	if err != nil {
		return nil, err
	}

	if len(endpoints) != 0 {
		return endpoints, nil
	}

	// Endpoints from resources.
	return h.endpointsFromResources(ctx, serviceID)
}

func (h Handler) endpointsFromOutput(ctx context.Context, serviceID oid.ID) (view.RouteAccessEndpointResponse, error) {
	outputs, err := h.getServiceOutputs(ctx, serviceID, true)
	if err != nil {
		return nil, err
	}

	var (
		invalidTypeErr = runtime.Error(http.StatusBadRequest,
			"element type of output endpoints should be string")
		endpoints = make([]view.AccessEndpoint, 0, len(outputs))
	)

	for _, v := range outputs {
		if !view.IsEndpointOutput(v.Name) {
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

			endpoints = append(endpoints, view.AccessEndpoint{
				Endpoints: []string{ep},
				Name:      v.Name,
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

			endpoints = append(endpoints, view.AccessEndpoint{
				Endpoints: eps,
				Name:      v.Name,
			})
		}
	}

	return endpoints, nil
}

func (h Handler) endpointsFromResources(ctx context.Context, serviceID oid.ID) ([]view.AccessEndpoint, error) {
	res, err := h.modelClient.ServiceResources().Query().
		Where(
			serviceresource.ServiceID(serviceID),
			serviceresource.TypeIn(
				intercept.TFEndpointsTypes...,
			),
		).
		Select(
			serviceresource.FieldConnectorID,
			serviceresource.FieldType,
			serviceresource.FieldName,
			serviceresource.FieldStatus,
		).
		All(ctx)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}

	var endpoints []view.AccessEndpoint

	for _, v := range res {
		for _, eps := range v.Status.ResourceEndpoints {
			endpoints = append(endpoints, view.AccessEndpoint{
				Name:      strs.Join("/", eps.EndpointType, v.Name),
				Endpoints: eps.Endpoints,
			})
		}
	}

	return endpoints, nil
}

func (h Handler) RouteOutputs(ctx *gin.Context, req view.RouteOutputRequest) (view.RouteOutputResponse, error) {
	return h.getServiceOutputs(ctx, req.ID, true)
}

func (h Handler) getServiceOutputs(
	ctx context.Context,
	serviceID oid.ID,
	onlySuccess bool,
) ([]types.OutputValue, error) {
	sr, err := h.modelClient.ServiceRevisions().Query().
		Where(servicerevision.ServiceID(serviceID)).
		Select(
			servicerevision.FieldOutput,
			servicerevision.FieldTemplateID,
			servicerevision.FieldTemplateVersion,
			servicerevision.FieldAttributes,
			servicerevision.FieldStatus,
		).
		WithService(func(sq *model.ServiceQuery) {
			sq.Select(service.FieldName)
		}).
		Order(model.Desc(servicerevision.FieldCreateTime)).
		First(ctx)
	if err != nil && !model.IsNotFound(err) {
		return nil, fmt.Errorf("error getting the latest service revision")
	}

	if sr == nil {
		return nil, nil
	}

	if onlySuccess && sr.Status != status.ServiceRevisionStatusSucceeded {
		return nil, nil
	}

	o, err := tfparser.ParseStateOutput(sr)
	if err != nil {
		return nil, fmt.Errorf("error get outputs: %w", err)
	}

	return o, nil
}

func (h Handler) updateServiceStatus(
	ctx context.Context,
	entity *model.Service,
) error {
	update, err := dao.ServiceStatusUpdate(h.modelClient, entity)
	if err != nil {
		return err
	}

	err = update.Exec(ctx)
	if err != nil && !model.IsNotFound(err) {
		return err
	}

	return nil
}

// CreateClone creates a clone of the service.
func (h Handler) CreateClone(
	ctx *gin.Context,
	req view.CreateCloneRequest,
) (*model.ServiceOutput, error) {
	service, err := h.modelClient.Services().Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	service.Name = req.Name

	if req.EnvironmentID != "" {
		service.EnvironmentID = req.EnvironmentID
	}

	return h.createService(ctx, createServiceOptions{
		Tags:    req.RemarkTags,
		Service: service,
	})
}

func (h Handler) createService(
	ctx context.Context,
	opts createServiceOptions,
) (*model.ServiceOutput, error) {
	logger := log.WithName("api").WithName("service")

	// Get deployer.
	createOpts := deptypes.CreateOptions{
		Type:        deployertf.DeployerType,
		ModelClient: h.modelClient,
		KubeConfig:  h.kubeConfig,
	}

	dp, err := deployer.Get(ctx, createOpts)
	if err != nil {
		return nil, err
	}

	// Create service, mark status to deploying.
	creates, err := dao.ServiceCreates(h.modelClient, opts.Service)
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
		status.ServiceStatusDeployed.False(entity, err.Error())

		uerr := h.updateServiceStatus(ctx, entity)
		if uerr != nil {
			logger.Errorf("error updating status of service %s: %v",
				entity.ID, uerr)
		}
	}()

	// Apply service.
	applyOpts := deptypes.ApplyOptions{
		SkipTLSVerify: !h.tlsCertified,
		Tags:          opts.Tags,
	}

	err = dp.Apply(ctx, entity, applyOpts)
	if err != nil {
		// NB(thxCode): a better approach is to use transaction,
		// however, building the service deployment process is a time-consuming task,
		// to prevent long-time transaction, we use a deletion to achieve this.
		// Usually, the probability of this delete operation failing is very low.
		derr := h.modelClient.Services().DeleteOne(entity).
			Exec(ctx)
		if derr != nil {
			logger.Errorf("error deleting: %v", derr)
		}

		return nil, err
	}

	return model.ExposeService(entity), nil
}

func (h Handler) StreamAccessEndpoint(ctx runtime.RequestUnidiStream, req view.GetRequest) error {
	t, err := topic.Subscribe(datamessage.ServiceRevision)
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

			if ar.ServiceID != req.ID {
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
				// While create new service revision,
				// the previous endpoints from outputs and resources need to be deleted.
				streamData = view.StreamAccessEndpointResponse{
					Type:       datamessage.EventDelete,
					Collection: eps,
				}
			case datamessage.EventUpdate:
				// While the service revision status is succeeded, the endpoints is updated to the current revision.
				if ar.Status != status.ServiceRevisionStatusSucceeded {
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

func (h Handler) getRevisionByID(ctx context.Context, id oid.ID) (*model.ServiceRevision, error) {
	return h.modelClient.ServiceRevisions().Query().
		Where(servicerevision.ID(id)).
		Only(ctx)
}

func (h Handler) StreamOutput(ctx runtime.RequestUnidiStream, req view.GetRequest) error {
	t, err := topic.Subscribe(datamessage.ServiceRevision)
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

			if ar.ServiceID != req.ID {
				continue
			}

			op, err := h.getServiceOutputs(ctx, ar.ServiceID, false)
			if err != nil {
				return err
			}

			if len(op) == 0 {
				continue
			}

			switch dm.Type {
			case datamessage.EventCreate:
				// While create new service revision, the outputs of new revision is the previous outputs.
				streamData = view.StreamOutputResponse{
					Type:       datamessage.EventDelete,
					Collection: op,
				}
			case datamessage.EventUpdate:
				// While the service revision status is succeeded, the outputs is updated to the current revision.
				if ar.Status != status.ServiceRevisionStatusSucceeded {
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

func (h Handler) getDeployer(ctx context.Context) (deptypes.Deployer, error) {
	return deployer.Get(ctx, deptypes.CreateOptions{
		Type:        deployertf.DeployerType,
		ModelClient: h.modelClient,
		KubeConfig:  h.kubeConfig,
	})
}

func (h Handler) CollectionGetGraph(
	ctx *gin.Context,
	req view.CollectionGetGraphRequest,
) (*view.CollectionGetGraphResponse, error) {
	logger := log.WithName("api").WithName("service")

	// Fetch entities.
	entities, err := h.modelClient.Services().Query().
		Order(model.Asc(service.FieldCreateTime)).
		Where(service.EnvironmentID(req.EnvironmentID)).
		Select(getFields...).
		// Must extract dependency.
		WithDependencies(func(dq *model.ServiceRelationshipQuery) {
			dq.Select(servicerelationship.FieldDependencyID).
				Where(func(s *sql.Selector) {
					s.Where(sql.ColumnsNEQ(servicerelationship.FieldServiceID, servicerelationship.FieldDependencyID))
				})
		}).
		// Must extract resource.
		WithResources(func(rq *model.ServiceResourceQuery) {
			rq.Order(model.Desc(serviceresource.FieldCreateTime)).
				Where(serviceresource.ModeNEQ(types.ServiceResourceModeDiscovered)).
				Select(
					serviceresource.FieldServiceID,
					serviceresource.FieldDeployerType,
					serviceresource.FieldType,
					serviceresource.FieldID,
					serviceresource.FieldName,
					serviceresource.FieldCreateTime,
					serviceresource.FieldUpdateTime,
					serviceresource.FieldStatus).
				// Must extract connector.
				Select(serviceresource.FieldConnectorID).
				WithConnector(func(cq *model.ConnectorQuery) {
					cq.Select(
						connector.FieldName,
						connector.FieldType,
						connector.FieldCategory,
						connector.FieldConfigVersion,
						connector.FieldConfigData)
				}).
				// Must extract components(resources).
				WithComponents(func(rq *model.ServiceResourceQuery) {
					rq.Order(model.Desc(serviceresource.FieldCreateTime)).
						Where(serviceresource.Mode(types.ServiceResourceModeDiscovered)).
						Select(
							serviceresource.FieldServiceID,
							serviceresource.FieldDeployerType,
							serviceresource.FieldType,
							serviceresource.FieldID,
							serviceresource.FieldName,
							serviceresource.FieldCreateTime,
							serviceresource.FieldUpdateTime,
							serviceresource.FieldStatus)
				})
		}).
		All(ctx)
	if err != nil {
		return nil, err
	}

	// Index entities by connector.
	m := make(map[*model.Connector][][2]int)

	for si := range entities {
		for i := range entities[si].Edges.Resources {
			// NB(thxCode): we can safety index the connector with its pointer here,
			// as the ent can keep the connector pointer is the same
			// between those resources related by the same connector.
			m[entities[si].Edges.Resources[i].Edges.Connector] = append(
				m[entities[si].Edges.Resources[i].Edges.Connector], [2]int{si, i})
		}
	}

	// Fetch keys for each resource without error returning.
	for c, idxs := range m {
		// Get operator by connector.
		op, err := operator.Get(ctx, optypes.CreateOptions{Connector: *c})
		if err != nil {
			logger.Warnf("cannot get operator of connector: %v", err)
			continue
		}

		if err = op.IsConnected(ctx); err != nil {
			logger.Warnf("unreachable connector: %v", err)
			continue
		}

		// Fetch keys for the resources that related to same connector.
		for _, idx := range idxs {
			si, i := idx[0], idx[1]
			entities[si].Edges.Resources[i].Keys, err = op.GetKeys(ctx, entities[si].Edges.Resources[i])

			if err != nil {
				logger.Errorf("error getting keys for %q: %v",
					entities[si].Edges.Resources[i].ID,
					err)

				continue
			}

			// Fetch keys of components as well.
			for j := range entities[si].Edges.Resources[i].Edges.Components {
				entities[si].Edges.Resources[i].Edges.Components[j].Keys, err = op.GetKeys(
					ctx,
					entities[si].Edges.Resources[i].Edges.Components[j],
				)
				if err != nil {
					logger.Errorf("error getting keys for %q: %v",
						entities[si].Edges.Resources[i].Edges.Components[j].ID,
						err)
				}
			}
		}
	}

	// Calculate capacity for allocation.
	var verticesCap, edgesCap int
	{
		// Count the number of Service.
		verticesCap = len(entities)
		for i := 0; i < len(entities); i++ {
			// Count the vertex size of ServiceResource,
			// and the edge size from Service to ServiceResource.
			verticesCap += len(entities[i].Edges.Resources)
			edgesCap += len(entities[i].Edges.Dependencies)

			for j := 0; j < len(entities[i].Edges.Resources); j++ {
				// Count the vertex size of sub ServiceResource,
				// and the edge size from ServiceResource to sub ServiceResource.
				verticesCap += len(entities[i].Edges.Resources[j].Edges.Components)
				edgesCap += len(entities[i].Edges.Resources[j].Edges.Components)
			}
		}
	}

	// Construct response.
	var (
		vertices = make([]view.GraphVertex, 0, verticesCap)
		edges    = make([]view.GraphEdge, 0, edgesCap)
	)

	for i := 0; i < len(entities); i++ {
		// Append Service to vertices.
		vertices = append(vertices, view.GraphVertex{
			GraphVertexID: view.GraphVertexID{
				Kind: "Service",
				ID:   entities[i].ID,
			},
			Name:        entities[i].Name,
			Description: entities[i].Description,
			Labels:      entities[i].Labels,
			CreateTime:  entities[i].CreateTime,
			UpdateTime:  entities[i].UpdateTime,
			Status:      entities[i].Status.Summary,
		})

		// Append the edge of Service to Service.
		for j := 0; j < len(entities[i].Edges.Dependencies); j++ {
			edges = append(edges, view.GraphEdge{
				Type: "Dependency",
				Start: view.GraphVertexID{
					Kind: "Service",
					ID:   entities[i].ID,
				},
				End: view.GraphVertexID{
					Kind: "Service",
					ID:   entities[i].Edges.Dependencies[j].DependencyID,
				},
			})
		}

		for j := 0; j < len(entities[i].Edges.Resources); j++ {
			// Append ServiceResource to vertices.
			vertices = append(vertices, view.GraphVertex{
				GraphVertexID: view.GraphVertexID{
					Kind: "ServiceResource",
					ID:   entities[i].Edges.Resources[j].ID,
				},
				Name:       entities[i].Edges.Resources[j].Name,
				CreateTime: entities[i].Edges.Resources[j].CreateTime,
				UpdateTime: entities[i].Edges.Resources[j].UpdateTime,
				Status:     entities[i].Edges.Resources[j].Status.Summary,
				Extensions: map[string]any{
					"type": entities[i].Edges.Resources[j].Type,
					"keys": entities[i].Edges.Resources[j].Keys,
				},
			})

			// Append the edge of Service to ServiceResource.
			edges = append(edges, view.GraphEdge{
				Type: "Composition",
				Start: view.GraphVertexID{
					Kind: "Service",
					ID:   entities[i].ID,
				},
				End: view.GraphVertexID{
					Kind: "ServiceResource",
					ID:   entities[i].Edges.Resources[j].ID,
				},
			})

			for k := 0; k < len(entities[i].Edges.Resources[j].Edges.Components); k++ {
				// Append sub ServiceResource to vertices.
				vertices = append(vertices, view.GraphVertex{
					GraphVertexID: view.GraphVertexID{
						Kind: "ServiceResource",
						ID:   entities[i].Edges.Resources[j].Edges.Components[k].ID,
					},
					Name:       entities[i].Edges.Resources[j].Edges.Components[k].Name,
					CreateTime: entities[i].Edges.Resources[j].Edges.Components[k].CreateTime,
					UpdateTime: entities[i].Edges.Resources[j].Edges.Components[k].UpdateTime,
					Status:     entities[i].Edges.Resources[j].Edges.Components[k].Status.Summary,
					Extensions: map[string]any{
						"type": entities[i].Edges.Resources[j].Edges.Components[k].Type,
						"keys": entities[i].Edges.Resources[j].Edges.Components[k].Keys,
					},
				})

				// Append the edge of ServiceResource to sub ServiceResource.
				edges = append(edges, view.GraphEdge{
					Type: "Composition",
					Start: view.GraphVertexID{
						Kind: "ServiceResource",
						ID:   entities[i].Edges.Resources[j].ID,
					},
					End: view.GraphVertexID{
						Kind: "ServiceResource",
						ID:   entities[i].Edges.Resources[j].Edges.Components[k].ID,
					},
				})
			}
		}
	}

	return &view.CollectionGetGraphResponse{
		Vertices: vertices,
		Edges:    edges,
	}, nil
}
