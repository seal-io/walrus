package applicationinstance

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/seal-io/seal/pkg/apis/applicationinstance/view"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/applicationinstance"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/platform"
	"github.com/seal-io/seal/pkg/platform/deployer"
	"github.com/seal-io/seal/pkg/platformk8s"
	"github.com/seal-io/seal/pkg/platformk8s/kube"
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
	return "ApplicationInstance"
}

func (h Handler) Validating() any {
	return h.modelClient
}

// Basic APIs

func (h Handler) Delete(ctx *gin.Context, req view.DeleteRequest) error {
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

	// mark status to deleting.
	entity, err = h.modelClient.ApplicationInstances().UpdateOne(entity).
		SetStatus(status.ApplicationInstanceStatusDeleting).
		SetStatusMessage("").
		Save(ctx)
	if err != nil {
		return err
	}

	// destroy instance.
	var destroyOpts = deployer.DestroyOptions{
		SkipTLSVerify: !h.tlsCertified,
	}
	return dp.Destroy(ctx, entity, destroyOpts)
}

// Batch APIs

func (h Handler) CollectionDelete(ctx *gin.Context, req view.CollectionDeleteRequest) error {
	return h.modelClient.WithTx(ctx, func(tx *model.Tx) (err error) {
		for i := range req {
			err = tx.ApplicationInstances().DeleteOne(req[i].Model()).
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

// Extensional APIs

func (h Handler) RouteUpgrade(ctx *gin.Context, req view.RouteUpgradeRequest) error {
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

	// apply instance.
	var applyOpts = deployer.ApplyOptions{
		SkipTLSVerify: !h.tlsCertified,
	}
	return dp.Apply(ctx, entity, applyOpts)
}

func (h Handler) RouteAccessEndpoints(ctx *gin.Context, req view.AccessEndpointRequest) (*view.AccessEndpointResponse, error) {
	serviceType := "kubernetes_service"
	serviceTypeAlias := "kubernetes_service_v1"
	ingressType := "kubernetes_ingress"
	ingressTypeAlias := "kubernetes_ingress_v1"
	res, err := h.modelClient.ApplicationResources().Query().
		Where(
			applicationresource.InstanceID(req.ID),
			applicationresource.TypeIn(serviceType, serviceTypeAlias, ingressType, ingressTypeAlias),
		).
		Select(
			applicationresource.FieldConnectorID,
			applicationresource.FieldType,
			applicationresource.FieldName,
		).
		All(ctx)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}

	// TODO: query from application resource after implemented store endpoints and conditional status in application resource status
	conns, err := h.modelClient.Connectors().Query().All(ctx)
	if err != nil {
		return nil, err
	}
	var connMap = make(map[types.ID]*model.Connector)
	for _, v := range conns {
		connMap[v.ID] = v
	}

	var connResMapping = make(map[*model.Connector]map[string]sets.Set[string])
	for _, v := range res {
		conn := connMap[v.ConnectorID]
		if conn == nil {
			continue
		}

		if _, ok := connResMapping[conn]; !ok {
			connResMapping[conn] = make(map[string]sets.Set[string])
		}
		if _, ok := connResMapping[conn][v.Type]; !ok {
			connResMapping[conn][v.Type] = sets.Set[string]{}
		}
		connResMapping[conn][v.Type] = connResMapping[conn][v.Type].Insert(v.Name)
	}

	var endpoints []kube.ResourceEndpoint
	for conn, ress := range connResMapping {
		k8sClient, err := k8sClientset(conn)
		if err != nil {
			return nil, err
		}

		for resType, names := range ress {
			switch resType {
			case serviceType, serviceTypeAlias:
				eps, err := kube.ServiceEndpointGetter(k8sClient).Endpoints(ctx, names.UnsortedList()...)
				if err != nil {
					return nil, err
				}
				endpoints = append(endpoints, eps...)
			case ingressType, ingressTypeAlias:
				eps, err := kube.IngressEndpointGetter(k8sClient).Endpoints(ctx, names.UnsortedList()...)
				if err != nil {
					return nil, err
				}
				endpoints = append(endpoints, eps...)
			}
		}
	}
	return &view.AccessEndpointResponse{
		Endpoints: endpoints,
	}, nil
}

func k8sClientset(conn *model.Connector) (*kubernetes.Clientset, error) {
	restCfg, err := platformk8s.GetConfig(*conn)
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(restCfg)
	if err != nil {
		return nil, fmt.Errorf("error creating kubernetes core client: %w", err)
	}
	return client, nil
}
