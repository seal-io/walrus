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
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/deployer"
	"github.com/seal-io/seal/pkg/deployer/terraform"
	deptypes "github.com/seal-io/seal/pkg/deployer/types"
	"github.com/seal-io/seal/pkg/operator"
	optypes "github.com/seal-io/seal/pkg/operator/types"
	tfparser "github.com/seal-io/seal/pkg/terraform/parser"
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

// Basic APIs.

func (h Handler) Get(ctx *gin.Context, req view.GetRequest) (view.GetResponse, error) {
	entity, err := h.modelClient.ApplicationRevisions().Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return model.ExposeApplicationRevision(entity), nil
}

func (h Handler) Stream(ctx runtime.RequestUnidiStream, req view.StreamRequest) error {
	t, err := topic.Subscribe(datamessage.ApplicationRevision)
	if err != nil {
		return err
	}

	defer func() { t.Unsubscribe() }()

	query := h.modelClient.ApplicationRevisions().Query().
		Select(getFields...)

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

// Batch APIs.

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
		applicationrevision.FieldVariables,
		applicationrevision.FieldSecrets,
	)
	sortFields = []string{
		applicationrevision.FieldCreateTime,
	}
)

func (h Handler) CollectionGet(
	ctx *gin.Context,
	req view.CollectionGetRequest,
) (view.CollectionGetResponse, int, error) {
	query := h.modelClient.ApplicationRevisions().Query()
	if req.InstanceID != "" {
		query.Where(applicationrevision.InstanceID(req.InstanceID))
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
		func(eq *model.EnvironmentQuery) {
			eq.Select(
				environment.FieldID,
				environment.FieldName)
		}).
		Unique(false). // Allow returning without sorting keys.
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeApplicationRevisions(entities), cnt, nil
}

func (h Handler) CollectionStream(
	ctx runtime.RequestUnidiStream,
	req view.CollectionStreamRequest,
) error {
	t, err := topic.Subscribe(datamessage.ApplicationRevision)
	if err != nil {
		return err
	}

	defer func() { t.Unsubscribe() }()

	query := h.modelClient.ApplicationRevisions().Query()
	if req.InstanceID != "" {
		query.Where(applicationrevision.InstanceID(req.InstanceID))
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

// Extensional APIs.

// GetTerraformStates get the terraform states of the application revision deployment.
func (h Handler) GetTerraformStates(
	ctx *gin.Context,
	req view.GetTerraformStatesRequest,
) (view.GetTerraformStatesResponse, error) {
	entity, err := h.modelClient.ApplicationRevisions().Query().
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
func (h Handler) UpdateTerraformStates(
	ctx *gin.Context,
	req view.UpdateTerraformStatesRequest,
) (err error) {
	logger := log.WithName("api").WithName("application-revision")

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

		// Timeout context.
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
			logger.Errorf("update status failed: %v", err)
			return
		}

		if nerr := revisionbus.Notify(ctx, h.modelClient, updateRevision); nerr != nil {
			logger.Errorf("notify failed: %v", nerr)
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
	key := func(r *model.ApplicationResource) string {
		// Align to schema definition.
		return strs.Join("-", string(r.ConnectorID), r.Module, r.Mode, r.Type, r.Name)
	}

	var p tfparser.Parser

	observedRess, err := p.ParseAppRevision(entity)
	if err != nil {
		return err
	}

	if observedRess == nil {
		return nil
	}

	// Get record resources from local.
	recordRess, err := h.modelClient.ApplicationResources().Query().
		Where(applicationresource.InstanceID(entity.InstanceID)).
		All(ctx)
	if err != nil {
		return err
	}

	// Calculate creating list and deleting list.
	observedRessIndex := make(map[string]*model.ApplicationResource, len(observedRess))

	for j := range observedRess {
		c := observedRess[j]
		observedRessIndex[key(c)] = c
	}
	deleteRessIDs := make([]oid.ID, 0, len(recordRess))

	for _, c := range recordRess {
		k := key(c)
		if observedRessIndex[k] != nil {
			delete(observedRessIndex, k)
			continue
		}

		deleteRessIDs = append(deleteRessIDs, c.ID)
	}

	createRess := make([]*model.ApplicationResource, 0, len(observedRessIndex))
	for k := range observedRessIndex {
		createRess = append(createRess, observedRessIndex[k])
	}

	// Diff by transactional session.
	err = h.modelClient.WithTx(ctx, func(tx *model.Tx) error {
		// Create new resources.
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
		// Delete stale resources.
		if len(deleteRessIDs) != 0 {
			_, err = tx.ApplicationResources().Delete().
				Where(applicationresource.IDIn(deleteRessIDs...)).
				Exec(ctx)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	if len(createRess) == 0 {
		return nil
	}

	// State/label the new resources async.
	ids := make(map[oid.ID][]oid.ID)
	for i := range createRess {
		// Group resources by connector.
		ids[createRess[i].ConnectorID] = append(ids[createRess[i].ConnectorID],
			createRess[i].ID)
	}

	gopool.Go(func() {
		logger := log.WithName("api").WithName("application-revision")

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		// Fetch related connectors at once,
		// and then index these connectors by its id.
		cs, err := h.modelClient.Connectors().Query().
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

		csidx := make(map[oid.ID]*model.Connector, len(cs))
		for i := range cs {
			csidx[cs[i].ID] = cs[i]
		}

		var sr applicationresources.StateResult

		for cid, crids := range ids {
			entities, err := applicationresources.ListCandidatesByIDs(ctx, h.modelClient, crids)
			if err != nil {
				logger.Errorf("error listing candidates: %v", err)
				continue
			}

			if len(entities) == 0 {
				continue
			}

			c, exist := csidx[cid]
			if !exist {
				continue
			}

			op, err := operator.Get(ctx, optypes.CreateOptions{
				Connector: *c,
			})
			if err != nil {
				logger.Errorf("error getting operator of connector %s: %v",
					c.ID, err)
				continue
			}

			nsr, err := applicationresources.State(ctx, op, h.modelClient, entities)
			if err != nil {
				logger.Errorf("error stating entities: %v", err)
				// Mark error as transitioning,
				// which doesn't flip the status.
				nsr.Transitioning = true
			}

			sr.Merge(nsr)

			err = applicationresources.Label(ctx, op, entities)
			if err != nil {
				logger.Errorf("error labeling entities: %v", err)
			}
		}

		// State application instance.
		i, err := h.modelClient.ApplicationInstances().Query().
			Where(applicationinstance.ID(entity.InstanceID)).
			Select(
				applicationinstance.FieldID,
				applicationinstance.FieldStatus).
			Only(ctx)
		if err != nil {
			logger.Errorf("cannot get application instance: %v", err)
			return
		}

		if status.ApplicationInstanceStatusDeleted.Exist(i) {
			// Skip if the instance is on deleting.
			return
		}

		switch {
		case sr.Error:
			status.ApplicationInstanceStatusReady.False(i, "")
		case sr.Transitioning:
			status.ApplicationInstanceStatusReady.Unknown(i, "")
		default:
			status.ApplicationInstanceStatusReady.True(i, "")
		}

		update, err := dao.ApplicationInstanceStatusUpdate(h.modelClient, i)
		if err != nil {
			logger.Errorf("cannot update application instance: %v", err)
		}

		err = update.Exec(ctx)
		if err != nil {
			logger.Errorf("cannot update application instance: %v", err)
		}
	})

	return nil
}

func (h Handler) StreamLog(ctx runtime.RequestUnidiStream, req view.StreamLogRequest) error {
	// NB(thxCode): disable timeout as we don't know the maximum time-cost of once tracing,
	// and rely on the session context timeout control,
	// which means we don't close the underlay kubernetes client operation until the `ctx` is cancel.
	restConfig := *h.kubeConfig // Copy.
	restConfig.Timeout = 0

	cli, err := coreclient.NewForConfig(&restConfig)
	if err != nil {
		return fmt.Errorf("error creating kubernetes client: %w", err)
	}

	return terraform.StreamJobLogs(ctx, terraform.StreamJobLogsOptions{
		Cli:        cli,
		RevisionID: req.ID,
		JobType:    req.JobType,
		Out:        ctx,
	})
}

// RouteRollbackInstances rollback instance to a specific revision.
func (h Handler) RouteRollbackInstances(ctx *gin.Context, req view.RollbackInstanceRequest) error {
	applicationRevision, err := h.modelClient.ApplicationRevisions().Get(ctx, req.ID)
	if err != nil {
		return err
	}

	applicationInstance, err := h.modelClient.ApplicationInstances().
		Get(ctx, applicationRevision.InstanceID)
	if err != nil {
		return err
	}

	createOpts := deptypes.CreateOptions{
		Type:        terraform.DeployerType,
		ModelClient: h.modelClient,
		KubeConfig:  h.kubeConfig,
	}

	dp, err := deployer.Get(ctx, createOpts)
	if err != nil {
		return err
	}

	rollbackOpts := deptypes.RollbackOptions{
		SkipTLSVerify: !h.tlsCertified,
		CloneFrom:     applicationRevision,
	}

	return dp.Rollback(ctx, applicationInstance, rollbackOpts)
}

// RouteRollbackApplications rollback application to a specific revision.
func (h Handler) RouteRollbackApplications(
	ctx *gin.Context,
	req view.RollbackApplicationRequest,
) error {
	// Get application revision.
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

	amr := make(model.ApplicationModuleRelationships, len(applicationRevision.Modules))
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
	app.Variables = applicationRevision.Variables

	return h.modelClient.WithTx(ctx, func(tx *model.Tx) error {
		updates, err := dao.ApplicationUpdates(tx, app)
		if err != nil {
			return err
		}

		return updates[0].Exec(ctx)
	})
}

// GetDiffLatest get the revision with the application instance latest revision diff.
func (h Handler) GetDiffLatest(ctx *gin.Context, req view.DiffLatestRequest) (*view.RevisionDiffResponse, error) {
	compareRevision, err := h.modelClient.ApplicationRevisions().Query().
		Select(
			applicationrevision.FieldID,
			applicationrevision.FieldInstanceID,
			applicationrevision.FieldModules,
			applicationrevision.FieldInputVariables,
			applicationrevision.FieldVariables,
		).
		Where(applicationrevision.ID(req.ID)).
		Order(model.Desc(applicationrevision.FieldCreateTime)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	app, err := h.modelClient.ApplicationInstances().Query().
		Where(applicationinstance.ID(compareRevision.InstanceID)).
		QueryApplication().
		Select(application.FieldVariables).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	latestRevision, err := h.modelClient.ApplicationRevisions().Query().
		Select(
			applicationrevision.FieldID,
			applicationrevision.FieldModules,
			applicationrevision.FieldInputVariables,
			applicationrevision.FieldVariables,
		).
		Where(applicationrevision.InstanceID(compareRevision.InstanceID)).
		Order(model.Desc(applicationrevision.FieldCreateTime)).
		First(ctx)
	if err != nil {
		return nil, err
	}

	return &view.RevisionDiffResponse{
		Old: view.RevisionDiff{
			InputVariables: latestRevision.InputVariables,
			Variables:      app.Variables,
			Modules:        latestRevision.Modules,
		},
		New: view.RevisionDiff{
			InputVariables: compareRevision.InputVariables,
			Variables:      compareRevision.Variables,
			Modules:        compareRevision.Modules,
		},
	}, nil
}

// GetDiffPrevious get the revision with the application instance previous revision diff.
func (h Handler) GetDiffPrevious(
	ctx *gin.Context,
	req view.RevisionDiffPreviousRequest,
) (*view.RevisionDiffResponse, error) {
	compareRevision, err := h.modelClient.ApplicationRevisions().Query().
		Select(
			applicationrevision.FieldID,
			applicationrevision.FieldModules,
			applicationrevision.FieldInputVariables,
			applicationrevision.FieldVariables,
			applicationrevision.FieldInstanceID,
			applicationrevision.FieldCreateTime,
		).
		Where(applicationrevision.ID(req.ID)).
		Order(model.Desc(applicationrevision.FieldCreateTime)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	var old view.RevisionDiff

	previousRevision, err := h.modelClient.ApplicationRevisions().Query().
		Select(
			applicationrevision.FieldID,
			applicationrevision.FieldModules,
			applicationrevision.FieldInputVariables,
			applicationrevision.FieldVariables,
		).
		Where(
			applicationrevision.InstanceID(compareRevision.InstanceID),
			applicationrevision.CreateTimeLT(*compareRevision.CreateTime),
		).
		Order(model.Desc(applicationrevision.FieldCreateTime)).
		First(ctx)
	if err != nil && !model.IsNotFound(err) {
		return nil, err
	}

	if previousRevision != nil {
		old = view.RevisionDiff{
			InputVariables: previousRevision.InputVariables,
			Variables:      previousRevision.Variables,
			Modules:        previousRevision.Modules,
		}
	}

	return &view.RevisionDiffResponse{
		Old: old,
		New: view.RevisionDiff{
			InputVariables: compareRevision.InputVariables,
			Variables:      compareRevision.Variables,
			Modules:        compareRevision.Modules,
		},
	}, nil
}
