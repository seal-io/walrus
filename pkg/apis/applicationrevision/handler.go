package applicationrevision

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/util/sets"
	coreclient "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"

	"github.com/seal-io/seal/pkg/apis/applicationrevision/view"
	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/applicationresources"
	revisionbus "github.com/seal-io/seal/pkg/bus/applicationrevision"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/model/applicationinstance"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/model/applicationrevision"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/project"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/platform"
	"github.com/seal-io/seal/pkg/platform/deployer"
	"github.com/seal-io/seal/pkg/platform/operator"
	"github.com/seal-io/seal/pkg/platformtf"
	"github.com/seal-io/seal/pkg/topic/datamessage"
	"github.com/seal-io/seal/utils/gopool"
	"github.com/seal-io/seal/utils/log"
	"github.com/seal-io/seal/utils/strs"
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
	return "ApplicationRevision"
}

func (h Handler) Validating() any {
	return h.modelClient
}

// Basic APIs

func (h Handler) Get(ctx *gin.Context, req view.GetRequest) (view.GetResponse, error) {
	var entity, err = h.modelClient.ApplicationRevisions().Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return model.ExposeApplicationRevision(entity), nil
}

func (h Handler) Stream(ctx runtime.RequestStream, req view.StreamRequest) error {
	var t, err = topic.Subscribe(datamessage.ApplicationRevision)
	if err != nil {
		return err
	}
	query := h.modelClient.ApplicationRevisions().Query().
		Select(getFields...)

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
				entity, err := query.Clone().
					Where(applicationrevision.ID(id)).
					WithEnvironment(func(q *model.EnvironmentQuery) {
						q.Select(environment.FieldID, environment.FieldName)
					}).
					Only(ctx)
				if err != nil {
					return err
				}
				streamData = view.StreamResponse{
					Type: dm.Type,
					Collection: []*model.ApplicationRevisionOutput{
						model.ExposeApplicationRevision(entity),
					},
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

func (h Handler) CollectionDelete(ctx *gin.Context, req view.CollectionDeleteRequest) error {
	return h.modelClient.WithTx(ctx, func(tx *model.Tx) (err error) {
		for i := range req {
			err = tx.ApplicationRevisions().DeleteOne(req[i].Model()).
				Exec(ctx)
			if err != nil {
				return err
			}
		}
		return
	})
}

var (
	getFields = applicationrevision.WithoutFields(
		applicationrevision.FieldStatusMessage,
		applicationrevision.FieldInputPlan,
		applicationrevision.FieldOutput,
		applicationrevision.FieldInputVariables,
		applicationrevision.FieldModules,
	)
	sortFields = []string{
		applicationrevision.FieldCreateTime}
)

func (h Handler) CollectionGet(ctx *gin.Context, req view.CollectionGetRequest) (view.CollectionGetResponse, int, error) {
	var query = h.modelClient.ApplicationRevisions().Query()
	if req.InstanceID != "" {
		query.Where(applicationrevision.InstanceID(req.InstanceID))
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
	if orders, ok := req.Sorting(sortFields, model.Desc(applicationrevision.FieldCreateTime)); ok {
		query.Order(orders...)
	}

	if req.InstanceID == "" {
		query.WithInstance(func(aiq *model.ApplicationInstanceQuery) {
			aiq.Select(
				applicationinstance.FieldID,
				applicationinstance.FieldName,
				applicationinstance.FieldApplicationID,
			).WithApplication(func(aq *model.ApplicationQuery) {
				aq.Select(
					application.FieldID,
					application.FieldName,
					application.FieldProjectID).
					WithProject(func(pq *model.ProjectQuery) {
						pq.Select(project.FieldID, project.FieldName)
					})
			})
		})
	}
	entities, err := query.WithEnvironment(
		func(eq *model.EnvironmentQuery) { eq.Select(environment.FieldID, environment.FieldName) }).
		Unique(false). // allow returning without sorting keys.
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeApplicationRevisions(entities), cnt, nil
}

func (h Handler) CollectionStream(ctx runtime.RequestStream, req view.CollectionStreamRequest) error {
	var t, err = topic.Subscribe(datamessage.ApplicationRevision)
	if err != nil {
		return err
	}

	query := h.modelClient.ApplicationRevisions().Query()
	if req.InstanceID != "" {
		query.Where(applicationrevision.InstanceID(req.InstanceID))
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
			revisions, err := query.Clone().
				Where(applicationrevision.IDIn(dm.Data...)).
				WithEnvironment(func(eq *model.EnvironmentQuery) {
					eq.Select(environment.FieldID, environment.FieldName)
				}).
				Unique(false).
				All(ctx)

			if err != nil && !model.IsNotFound(err) {
				return err
			}
			streamData = view.StreamResponse{
				Type:       dm.Type,
				Collection: model.ExposeApplicationRevisions(revisions),
			}
		case datamessage.EventDelete:
			streamData = view.StreamResponse{
				Type: dm.Type,
				IDs:  dm.Data,
			}
		}
		err = ctx.SendJSON(streamData)
		if err != nil {
			return err
		}
	}
}

// Extensional APIs

// GetTerraformStates get the terraform states of the application revision deployment.
func (h Handler) GetTerraformStates(ctx *gin.Context, req view.GetTerraformStatesRequest) (view.GetTerraformStatesResponse, error) {
	var entity, err = h.modelClient.ApplicationRevisions().Query().
		Where(applicationrevision.ID(req.ID)).
		Select(applicationrevision.FieldOutput).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	if entity.Output == "" {
		return nil, nil
	}
	return view.GetTerraformStatesResponse(entity.Output), nil
}

// UpdateTerraformStates update the terraform states of the application revision deployment.
func (h Handler) UpdateTerraformStates(ctx *gin.Context, req view.UpdateTerraformStatesRequest) (err error) {
	var logger = log.WithName("platformtf").WithName("state")

	entity, err := h.modelClient.ApplicationRevisions().Get(ctx, req.ID)
	if err != nil {
		return err
	}
	entity.Output = string(req.RawMessage)
	update, err := dao.ApplicationRevisionUpdate(h.modelClient, entity)
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

		// timeout context
		updateCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		entity.Status = status.ApplicationRevisionStatusFailed
		entity.StatusMessage = err.Error()
		revisionUpdate, updateErr := dao.ApplicationRevisionUpdate(h.modelClient, entity)
		if updateErr != nil {
			logger.Error(updateErr)
			return
		}
		updateRevision, updateErr := revisionUpdate.Save(updateCtx)
		if updateErr != nil {
			logger.Errorf("update application revision status failed: %v", err)
			return
		}

		if nerr := revisionbus.Notify(ctx, h.modelClient, updateRevision); nerr != nil {
			logger.Errorf("notify application revision failed: %v", nerr)
		}
	}()

	if err = revisionbus.Notify(ctx, h.modelClient, entity); err != nil {
		return err
	}

	return h.manageResources(ctx, entity)
}

// manageResources manages the resources of the given revision,
// and states/labels the resources within 3 minutes in the background.
func (h Handler) manageResources(ctx context.Context, entity *model.ApplicationRevision) error {
	// TODO(thxCode): generate by entc.
	var key = func(r *model.ApplicationResource) string {
		// align to schema definition.
		return strs.Join("-", string(r.ConnectorID), r.Module, r.Mode, r.Type, r.Name)
	}

	var p platformtf.Parser
	var observedRess, err = p.ParseAppRevision(entity)
	if err != nil {
		return err
	}
	if observedRess == nil {
		return nil
	}

	// get record resources from local.
	recordRess, err := h.modelClient.ApplicationResources().
		Query().
		Where(applicationresource.InstanceID(entity.InstanceID)).
		All(ctx)
	if err != nil {
		return err
	}

	// calculate creating list and deleting list.
	var observedRessIndex = make(map[string]*model.ApplicationResource, len(observedRess))
	for j := range observedRess {
		var c = observedRess[j]
		observedRessIndex[key(c)] = c
	}
	var deleteRessIDs = make([]types.ID, 0, len(recordRess))
	for _, c := range recordRess {
		var k = key(c)
		if observedRessIndex[k] != nil {
			delete(observedRessIndex, k)
			continue
		}
		deleteRessIDs = append(deleteRessIDs, c.ID)
	}
	var createRess = make([]*model.ApplicationResource, 0, len(observedRessIndex))
	for k := range observedRessIndex {
		createRess = append(createRess, observedRessIndex[k])
	}

	// diff by transactional session.
	err = h.modelClient.WithTx(ctx, func(tx *model.Tx) error {
		// create new resources.
		if len(createRess) != 0 {
			creates, err := dao.ApplicationResourceCreates(tx, createRess...)
			if err != nil {
				return err
			}
			createRess, err = tx.ApplicationResources().CreateBulk(creates...).
				Save(ctx)
			if err != nil {
				return err
			}
		}
		// delete stale resources.
		_, err = tx.ApplicationResources().Delete().
			Where(applicationresource.IDIn(deleteRessIDs...)).
			Exec(ctx)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	if len(createRess) == 0 {
		return nil
	}

	// state/label the new resources async.
	var ids = make(map[types.ID][]types.ID)
	for i := range createRess {
		// group resources by connector.
		ids[createRess[i].ConnectorID] = append(ids[createRess[i].ConnectorID],
			createRess[i].ID)
	}
	gopool.Go(func() {
		var logger = log.WithName("application-revision")
		var ctx, cancel = context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		// fetch related connectors at once,
		// and then index these connectors by its id.
		var cs, err = h.modelClient.Connectors().Query().
			Select(
				connector.FieldID,
				connector.FieldName,
				connector.FieldType,
				connector.FieldCategory,
				connector.FieldConfigVersion,
				connector.FieldConfigData).
			Where(connector.IDIn(sets.KeySet(ids).UnsortedList()...)).
			All(ctx)
		if err != nil {
			logger.Errorf("cannot list connectors: %v", err)
			return
		}
		var csidx = make(map[types.ID]*model.Connector, len(cs))
		for i := range cs {
			csidx[cs[i].ID] = cs[i]
		}

		for cid, crids := range ids {
			entities, err := applicationresources.ListCandidatesByIDs(ctx, h.modelClient, crids)
			if err != nil {
				logger.Errorf("error listing candidates: %v", err)
				continue
			}
			if len(entities) == 0 {
				continue
			}

			var c, exist = csidx[cid]
			if !exist {
				continue
			}

			op, err := platform.GetOperator(ctx, operator.CreateOptions{
				Connector: *c,
			})
			if err != nil {
				logger.Errorf("error getting operator of connector %s: %v",
					c.ID, err)
				continue
			}

			err = applicationresources.State(ctx, op, h.modelClient, entities)
			if err != nil {
				logger.Errorf("error stating entities: %v", err)
			}
			err = applicationresources.Label(ctx, op, entities)
			if err != nil {
				logger.Errorf("error labeling entities: %v", err)
			}
		}
	})
	return nil
}

func (h Handler) StreamLog(ctx runtime.RequestStream, req view.StreamLogRequest) error {
	var cli, err = coreclient.NewForConfig(h.kubeConfig)
	if err != nil {
		return fmt.Errorf("error creating kubernetes client: %w", err)
	}

	return platformtf.StreamJobLogs(ctx, platformtf.StreamJobLogsOptions{
		Cli:        cli,
		RevisionID: req.ID,
		JobType:    req.JobType,
		Out:        ctx,
	})
}

// CreateRollbackInstances rollback instance to a specific revision.
func (h Handler) CreateRollbackInstances(ctx *gin.Context, req view.RollbackInstanceRequest) error {
	applicationRevision, err := h.modelClient.ApplicationRevisions().Get(ctx, req.ID)
	if err != nil {
		return err
	}

	applicationInstance, err := h.modelClient.ApplicationInstances().Get(ctx, applicationRevision.InstanceID)
	if err != nil {
		return err
	}

	var createOpts = deployer.CreateOptions{
		Type:        platformtf.DeployerType,
		ModelClient: h.modelClient,
		KubeConfig:  h.kubeConfig,
	}
	dp, err := platform.GetDeployer(ctx, createOpts)
	if err != nil {
		return err
	}

	var rollbackOpts = deployer.RollbackOptions{
		ApplicationRevision: applicationRevision,
		SkipTLSVerify:       !h.tlsCertified,
	}
	return dp.Rollback(ctx, applicationInstance, rollbackOpts)
}

// CreateRollbackApplications rollback application to a specific revision.
func (h Handler) CreateRollbackApplications(ctx *gin.Context, req view.RollbackApplicationRequest) error {
	// get application revision
	applicationRevision, err := h.modelClient.ApplicationRevisions().Query().
		WithInstance(func(q *model.ApplicationInstanceQuery) {
			q.Select(
				applicationinstance.FieldID,
				applicationinstance.FieldApplicationID,
			)
		}).
		Where(applicationrevision.ID(req.ID)).
		Only(ctx)
	if err != nil {
		return err
	}

	var amr = make(model.ApplicationModuleRelationships, len(applicationRevision.Modules))
	for k, m := range applicationRevision.Modules {
		amr[k] = &model.ApplicationModuleRelationship{
			Name:       m.Name,
			ModuleID:   m.ModuleID,
			Version:    m.Version,
			Attributes: m.Attributes,
		}
	}

	app, err := h.modelClient.Applications().Query().
		Where(application.ID(applicationRevision.Edges.Instance.ApplicationID)).
		Only(ctx)
	if err != nil {
		return err
	}
	app.Edges.Modules = amr
	// update application.
	updates, err := dao.ApplicationUpdates(h.modelClient, app)
	if err != nil {
		return err
	}
	return updates[0].Exec(ctx)
}
