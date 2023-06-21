package servicerevision

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

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/apis/servicerevision/view"
	revisionbus "github.com/seal-io/seal/pkg/bus/servicerevision"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/service"
	"github.com/seal-io/seal/pkg/dao/model/serviceresource"
	"github.com/seal-io/seal/pkg/dao/model/servicerevision"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/deployer/terraform"
	"github.com/seal-io/seal/pkg/operator"
	optypes "github.com/seal-io/seal/pkg/operator/types"
	"github.com/seal-io/seal/pkg/serviceresources"
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
	return "ServiceRevision"
}

func (h Handler) Validating() any {
	return h.modelClient
}

// Basic APIs.

func (h Handler) Get(ctx *gin.Context, req view.GetRequest) (view.GetResponse, error) {
	entity, err := h.modelClient.ServiceRevisions().Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return model.ExposeServiceRevision(entity), nil
}

func (h Handler) Stream(ctx runtime.RequestUnidiStream, req view.StreamRequest) error {
	t, err := topic.Subscribe(datamessage.ServiceRevision)
	if err != nil {
		return err
	}

	defer func() { t.Unsubscribe() }()

	query := h.modelClient.ServiceRevisions().Query().
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
					Where(servicerevision.ID(id)).
					WithEnvironment(func(q *model.EnvironmentQuery) {
						q.Select(environment.FieldID, environment.FieldName)
					}).
					Only(ctx)
				if err != nil {
					return err
				}
				streamData = view.StreamResponse{
					Type: dm.Type,
					Collection: []*model.ServiceRevisionOutput{
						model.ExposeServiceRevision(entity),
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
			err = tx.ServiceRevisions().DeleteOne(req[i].Model()).
				Exec(ctx)
			if err != nil {
				return err
			}
		}

		return
	})
}

var (
	getFields = servicerevision.WithoutFields(
		servicerevision.FieldStatusMessage,
		servicerevision.FieldInputPlan,
		servicerevision.FieldOutput,
		servicerevision.FieldTemplateID,
		servicerevision.FieldTemplateVersion,
		servicerevision.FieldAttributes,
		servicerevision.FieldSecrets,
	)
	sortFields = []string{
		servicerevision.FieldCreateTime,
	}
)

func (h Handler) CollectionGet(
	ctx *gin.Context,
	req view.CollectionGetRequest,
) (view.CollectionGetResponse, int, error) {
	query := h.modelClient.ServiceRevisions().Query()
	if req.ServiceID != "" {
		query.Where(servicerevision.ServiceID(req.ServiceID))
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

	if orders, ok := req.Sorting(sortFields, model.Desc(servicerevision.FieldCreateTime)); ok {
		query.Order(orders...)
	}

	if req.ServiceID == "" {
		query.WithService(func(aiq *model.ServiceQuery) {
			aiq.Select(
				service.FieldID,
				service.FieldName,
			)
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

	return model.ExposeServiceRevisions(entities), cnt, nil
}

func (h Handler) CollectionStream(
	ctx runtime.RequestUnidiStream,
	req view.CollectionStreamRequest,
) error {
	t, err := topic.Subscribe(datamessage.ServiceRevision)
	if err != nil {
		return err
	}

	defer func() { t.Unsubscribe() }()

	query := h.modelClient.ServiceRevisions().Query()
	if req.ServiceID != "" {
		query.Where(servicerevision.ServiceID(req.ServiceID))
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
				Where(servicerevision.IDIn(dm.Data...)).
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
				Collection: model.ExposeServiceRevisions(revisions),
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

// GetTerraformStates get the terraform states of the service revision deployment.
func (h Handler) GetTerraformStates(
	ctx *gin.Context,
	req view.GetTerraformStatesRequest,
) (view.GetTerraformStatesResponse, error) {
	entity, err := h.modelClient.ServiceRevisions().Query().
		Where(servicerevision.ID(req.ID)).
		Select(servicerevision.FieldOutput).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	if entity.Output == "" {
		return nil, nil
	}

	return view.GetTerraformStatesResponse(entity.Output), nil
}

// UpdateTerraformStates update the terraform states of the service revision deployment.
func (h Handler) UpdateTerraformStates(
	ctx *gin.Context,
	req view.UpdateTerraformStatesRequest,
) (err error) {
	logger := log.WithName("api").WithName("service-revision")

	entity, err := h.modelClient.ServiceRevisions().Get(ctx, req.ID)
	if err != nil {
		return err
	}
	entity.Output = string(req.RawMessage)

	update, err := dao.ServiceRevisionUpdate(h.modelClient, entity)
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

		entity.Status = status.ServiceRevisionStatusFailed
		entity.StatusMessage = err.Error()

		revisionUpdate, updateErr := dao.ServiceRevisionUpdate(h.modelClient, entity)
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
func (h Handler) manageResources(ctx context.Context, entity *model.ServiceRevision) error {
	// TODO(thxCode): generate by entc.
	key := func(r *model.ServiceResource) string {
		// Align to schema definition.
		return strs.Join("-", string(r.ConnectorID), r.Mode, r.Type, r.Name)
	}

	var p tfparser.Parser

	observedRess, err := p.ParseServiceRevision(entity)
	if err != nil {
		return err
	}

	if observedRess == nil {
		return nil
	}

	// Get record resources from local.
	recordRess, err := h.modelClient.ServiceResources().Query().
		Where(serviceresource.ServiceID(entity.ServiceID)).
		All(ctx)
	if err != nil {
		return err
	}

	// Calculate creating list and deleting list.
	observedRessIndex := make(map[string]*model.ServiceResource, len(observedRess))

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

	createRess := make([]*model.ServiceResource, 0, len(observedRessIndex))
	for k := range observedRessIndex {
		createRess = append(createRess, observedRessIndex[k])
	}

	// Diff by transactional session.
	err = h.modelClient.WithTx(ctx, func(tx *model.Tx) error {
		// Create new resources.
		if len(createRess) != 0 {
			creates, err := dao.ServiceResourceCreates(tx, createRess...)
			if err != nil {
				return err
			}

			createRess, err = tx.ServiceResources().CreateBulk(creates...).
				Save(ctx)
			if err != nil {
				return err
			}
		}
		// Delete stale resources.
		if len(deleteRessIDs) != 0 {
			_, err = tx.ServiceResources().Delete().
				Where(serviceresource.IDIn(deleteRessIDs...)).
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
		logger := log.WithName("api").WithName("service-revision")

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

		var sr serviceresources.StateResult

		for cid, crids := range ids {
			entities, err := serviceresources.ListCandidatesByIDs(ctx, h.modelClient, crids)
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

			nsr, err := serviceresources.State(ctx, op, h.modelClient, entities)
			if err != nil {
				logger.Errorf("error stating entities: %v", err)
				// Mark error as transitioning,
				// which doesn't flip the status.
				nsr.Transitioning = true
			}

			sr.Merge(nsr)

			err = serviceresources.Label(ctx, op, entities)
			if err != nil {
				logger.Errorf("error labeling entities: %v", err)
			}
		}

		// State service.
		i, err := h.modelClient.Services().Query().
			Where(service.ID(entity.ServiceID)).
			Select(
				service.FieldID,
				service.FieldStatus).
			Only(ctx)
		if err != nil {
			logger.Errorf("cannot get service: %v", err)
			return
		}

		if status.ServiceStatusDeleted.Exist(i) {
			// Skip if the service is on deleting.
			return
		}

		switch {
		case sr.Error:
			status.ServiceStatusReady.False(i, "")
		case sr.Transitioning:
			status.ServiceStatusReady.Unknown(i, "")
		default:
			status.ServiceStatusReady.True(i, "")
		}

		update, err := dao.ServiceStatusUpdate(h.modelClient, i)
		if err != nil {
			logger.Errorf("cannot update service: %v", err)
		}

		err = update.Exec(ctx)
		if err != nil {
			logger.Errorf("cannot update service: %v", err)
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

// GetDiffLatest get the revision with the service latest revision diff.
func (h Handler) GetDiffLatest(ctx *gin.Context, req view.DiffLatestRequest) (*view.RevisionDiffResponse, error) {
	compareRevision, err := h.modelClient.ServiceRevisions().Query().
		Select(
			servicerevision.FieldID,
			servicerevision.FieldServiceID,
			servicerevision.FieldTemplateID,
			servicerevision.FieldTemplateVersion,
			servicerevision.FieldAttributes,
		).
		Where(servicerevision.ID(req.ID)).
		Order(model.Desc(servicerevision.FieldCreateTime)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	latestRevision, err := h.modelClient.ServiceRevisions().Query().
		Select(
			servicerevision.FieldID,
			servicerevision.FieldTemplateID,
			servicerevision.FieldTemplateVersion,
			servicerevision.FieldAttributes,
		).
		Where(servicerevision.ServiceID(compareRevision.ServiceID)).
		Order(model.Desc(servicerevision.FieldCreateTime)).
		First(ctx)
	if err != nil {
		return nil, err
	}

	return &view.RevisionDiffResponse{
		Old: view.RevisionDiff{
			TemplateID:      latestRevision.TemplateID,
			TemplateVersion: latestRevision.TemplateVersion,
			Attributes:      latestRevision.Attributes,
		},
		New: view.RevisionDiff{
			TemplateID:      compareRevision.TemplateID,
			TemplateVersion: compareRevision.TemplateVersion,
			Attributes:      compareRevision.Attributes,
		},
	}, nil
}

// GetDiffPrevious get the revision with the service previous revision diff.
func (h Handler) GetDiffPrevious(
	ctx *gin.Context,
	req view.RevisionDiffPreviousRequest,
) (*view.RevisionDiffResponse, error) {
	compareRevision, err := h.modelClient.ServiceRevisions().Query().
		Select(
			servicerevision.FieldID,
			servicerevision.FieldTemplateID,
			servicerevision.FieldTemplateVersion,
			servicerevision.FieldAttributes,
			servicerevision.FieldServiceID,
			servicerevision.FieldCreateTime,
		).
		Where(servicerevision.ID(req.ID)).
		Order(model.Desc(servicerevision.FieldCreateTime)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	var old view.RevisionDiff

	previousRevision, err := h.modelClient.ServiceRevisions().Query().
		Select(
			servicerevision.FieldID,
			servicerevision.FieldTemplateID,
			servicerevision.FieldTemplateVersion,
			servicerevision.FieldAttributes,
		).
		Where(
			servicerevision.ServiceID(compareRevision.ServiceID),
			servicerevision.CreateTimeLT(*compareRevision.CreateTime),
		).
		Order(model.Desc(servicerevision.FieldCreateTime)).
		First(ctx)
	if err != nil && !model.IsNotFound(err) {
		return nil, err
	}

	if previousRevision != nil {
		old = view.RevisionDiff{
			TemplateID:      previousRevision.TemplateID,
			TemplateVersion: previousRevision.TemplateVersion,
			Attributes:      previousRevision.Attributes,
		}
	}

	return &view.RevisionDiffResponse{
		Old: old,
		New: view.RevisionDiff{
			TemplateID:      compareRevision.TemplateID,
			TemplateVersion: compareRevision.TemplateVersion,
			Attributes:      compareRevision.Attributes,
		},
	}, nil
}
