package resourcerevision

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"time"

	"k8s.io/apimachinery/pkg/util/sets"
	coreclient "k8s.io/client-go/kubernetes/typed/core/v1"

	revisionbus "github.com/seal-io/walrus/pkg/bus/resourcerevision"
	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/connector"
	"github.com/seal-io/walrus/pkg/dao/model/environment"
	"github.com/seal-io/walrus/pkg/dao/model/project"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/model/resourcecomponent"
	"github.com/seal-io/walrus/pkg/dao/model/resourcerevision"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/deployer/terraform"
	"github.com/seal-io/walrus/pkg/operator"
	optypes "github.com/seal-io/walrus/pkg/operator/types"
	"github.com/seal-io/walrus/pkg/resourcecomponents"
	tfparser "github.com/seal-io/walrus/pkg/terraform/parser"
	"github.com/seal-io/walrus/utils/gopool"
	"github.com/seal-io/walrus/utils/log"
)

// RouteGetTerraformStates get the terraform states of the service revision deployment.
func (h Handler) RouteGetTerraformStates(
	req RouteGetTerraformStatesRequest,
) (RouteGetTerraformStatesResponse, error) {
	entity, err := h.modelClient.ResourceRevisions().Query().
		Where(resourcerevision.ID(req.ID)).
		Select(resourcerevision.FieldOutput).
		Only(req.Context)
	if err != nil {
		return nil, err
	}

	if entity.Output == "" {
		return nil, nil
	}

	return RouteGetTerraformStatesResponse(entity.Output), nil
}

// RouteUpdateTerraformStates update the terraform states of the service revision deployment.
func (h Handler) RouteUpdateTerraformStates(
	req RouteUpdateTerraformStatesRequest,
) (err error) {
	logger := log.WithName("api").WithName("service-revision")

	entity, err := h.modelClient.ResourceRevisions().Query().
		Where(resourcerevision.ID(req.ID)).
		Select(
			resourcerevision.FieldID,
			resourcerevision.FieldProjectID,
			resourcerevision.FieldEnvironmentID,
			resourcerevision.FieldResourceID,
			resourcerevision.FieldStatus,
			resourcerevision.FieldRecord,
			resourcerevision.FieldDeployerType).
		WithProject(func(pq *model.ProjectQuery) {
			pq.Select(
				project.FieldID,
				project.FieldName)
		}).
		WithEnvironment(func(eq *model.EnvironmentQuery) {
			eq.Select(
				environment.FieldID,
				environment.FieldName)
		}).
		WithResource(func(sq *model.ResourceQuery) {
			sq.Select(
				resource.FieldID,
				resource.FieldName)
		}).
		Only(req.Context)
	if err != nil {
		return err
	}
	entity.Output = string(req.RawMessage)

	err = h.modelClient.ResourceRevisions().UpdateOne(entity).
		SetOutput(entity.Output).
		Exec(req.Context)
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

		status.ResourceRevisionStatusReady.False(entity, err.Error())
		entity.Status.SetSummary(status.WalkResourceRevision(&entity.Status))

		uerr := h.modelClient.ResourceRevisions().UpdateOne(entity).
			SetStatus(entity.Status).
			Exec(updateCtx)
		if uerr != nil {
			logger.Errorf("update status failed: %v", err)
			return
		}

		if nerr := revisionbus.Notify(updateCtx, h.modelClient, entity); nerr != nil {
			logger.Errorf("notify failed: %v", nerr)
		}
	}()

	if err = revisionbus.Notify(req.Context, h.modelClient, entity); err != nil {
		return err
	}

	return manageResources(req.Context, h.modelClient, entity)
}

// manageResources parses the resources from the given revision,
// removes the stale resources from the database,
// creates the new resources to the database,
// and then execute reconcileResources for the new created resources.
func manageResources(
	ctx context.Context,
	modelClient model.ClientSet,
	entity *model.ResourceRevision,
) error {
	var p tfparser.Parser

	observedRess, dependencies, err := p.ParseResourceRevision(entity)
	if err != nil {
		return err
	}

	// Get record resources from local.
	recordRess, err := modelClient.ResourceComponents().Query().
		Where(resourcecomponent.ResourceID(entity.ResourceID)).
		All(ctx)
	if err != nil {
		return err
	}

	// Calculate creating list, deleting list and updating list.
	observedRessIndex := dao.ResourceComponentToMap(observedRess)

	deleteRessIDs := make([]object.ID, 0, len(recordRess))

	updatedRess := make([]*model.ResourceComponent, 0, len(recordRess))

	for _, c := range recordRess {
		k := dao.ResourceComponentGetUniqueKey(c)
		if observedRessIndex[k] != nil {
			c.Edges.Instances = observedRessIndex[k].Edges.Instances
			updatedRess = append(updatedRess, c)

			delete(observedRessIndex, k)

			continue
		}

		deleteRessIDs = append(deleteRessIDs, c.ID)
	}

	createRess := make([]*model.ResourceComponent, 0, len(observedRessIndex))

	for k := range observedRessIndex {
		// Resource instances will be created through edges.
		if observedRessIndex[k].Shape != types.ResourceComponentShapeClass {
			continue
		}

		createRess = append(createRess, observedRessIndex[k])
	}

	// Diff by transactional session.
	replacedRess := make([]*model.ResourceComponent, 0)

	err = modelClient.WithTx(ctx, func(tx *model.Tx) error {
		// Update resources with new instances.
		for _, r := range updatedRess {
			rp, err := dao.ResourceComponentInstancesEdgeSaveWithResult(ctx, tx, r)
			if err != nil {
				return err
			}

			replacedRess = append(replacedRess, rp...)
		}

		// Create new resources.
		if len(createRess) != 0 {
			createRess, err = tx.ResourceComponents().CreateBulk().
				Set(createRess...).
				SaveE(ctx, dao.ResourceComponentInstancesEdgeSave)
			if err != nil {
				return err
			}

			// TODO(thxCode): move the following codes into DAO.

			err = dao.ResourceComponentRelationshipUpdateWithDependencies(ctx, tx, dependencies, recordRess, createRess)
			if err != nil {
				return err
			}
		}

		// Delete stale resources.
		if len(deleteRessIDs) != 0 {
			_, err = tx.ResourceComponents().Delete().
				Where(resourcecomponent.IDIn(deleteRessIDs...)).
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

	reconcileRess := createRess
	reconcileRess = append(reconcileRess, updatedRess...)
	reconcileRess = append(reconcileRess, replacedRess...)

	if len(reconcileRess) == 0 {
		return nil
	}

	// Update the resource component status.
	reconcileRessIndex := dao.ResourceComponentToMap(reconcileRess)

	// Group resources by connector ID,
	// and decorate them with the project/environment/service for latter labeling.
	connRess := make(map[object.ID][]*model.ResourceComponent)

	for k := range reconcileRessIndex {
		if reconcileRessIndex[k].Shape != types.ResourceComponentShapeInstance {
			continue
		}

		reconcileRessIndex[k].Edges.Project = entity.Edges.Project
		reconcileRessIndex[k].Edges.Environment = entity.Edges.Environment
		reconcileRessIndex[k].Edges.Resource = entity.Edges.Resource

		connRess[reconcileRessIndex[k].ConnectorID] = append(connRess[reconcileRessIndex[k].ConnectorID],
			reconcileRessIndex[k])
	}

	gopool.Go(func() {
		logger := log.WithName("api").WithName("service-revision")

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		if err = reconcileResources(ctx, modelClient, entity.ResourceID, connRess); err != nil {
			logger.Errorf("sync service status and resource label failed: %v", err)
		}
	})

	return nil
}

// reconcileResources reconciles the resources,
// including states resources/labels resources/composes resources.
func reconcileResources(
	ctx context.Context,
	modelClient model.ClientSet,
	serviceID object.ID,
	connRess map[object.ID][]*model.ResourceComponent,
) error {
	logger := log.WithName("api").WithName("service-revision")

	// Fetch related connectors at once,
	// and then index these connectors by its id.
	cs, err := modelClient.Connectors().Query().
		Select(
			connector.FieldID,
			connector.FieldName,
			connector.FieldLabels,
			connector.FieldType,
			connector.FieldCategory,
			connector.FieldConfigVersion,
			connector.FieldConfigData).
		Where(connector.IDIn(sets.KeySet(connRess).UnsortedList()...)).
		All(ctx)
	if err != nil {
		return fmt.Errorf("cannot list connectors: %w", err)
	}

	ops := make(map[object.ID]optypes.Operator, len(cs))

	for i := range cs {
		op, err := operator.Get(ctx, optypes.CreateOptions{
			Connector:   *cs[i],
			ModelClient: modelClient,
		})
		if err != nil {
			// Warn out without breaking the whole syncing.
			logger.Warnf("cannot get operator of connector %q: %v", cs[i].ID, err)
			continue
		}

		if err = op.IsConnected(ctx); err != nil {
			// Warn out without breaking the whole syncing.
			logger.Warnf("unreachable connector %q", cs[i].ID)
			// Replace disconnected connector with unknown connector.
			op = operator.UnReachable()
		}
		ops[cs[i].ID] = op
	}

	var sr resourcecomponents.StateResult

	for cid, crs := range connRess {
		op, exist := ops[cid]
		if !exist {
			// Ignore if not found operator.
			continue
		}

		// Discover resources.
		ncrs, err := resourcecomponents.Discover(ctx, op, modelClient, crs)
		if err != nil {
			logger.Errorf("error discovering component resources: %v", err)
		}

		// State resources.
		nsr, err := resourcecomponents.State(ctx, op, modelClient, append(crs, ncrs...))
		if err != nil {
			logger.Errorf("error stating resources: %v", err)
			// Mark error as transitioning,
			// which doesn't flip the status.
			nsr.Transitioning = true
		}

		sr.Merge(nsr)

		// Label resources.
		err = resourcecomponents.Label(ctx, op, crs)
		if err != nil {
			logger.Errorf("error labeling resources: %v", err)
		}
	}

	// State resource.
	svc, err := modelClient.Resources().Query().
		Where(resource.ID(serviceID)).
		Select(
			resource.FieldID,
			resource.FieldStatus).
		Only(ctx)
	if err != nil {
		return fmt.Errorf("cannot get service: %w", err)
	}

	if status.ResourceStatusDeleted.Exist(svc) ||
		status.ResourceStatusStopped.Exist(svc) {
		// Skip if the service is on deleting or stopping.
		return nil
	}

	switch {
	case sr.Error:
		status.ResourceStatusReady.False(svc, "")
	case sr.Transitioning:
		status.ResourceStatusReady.Unknown(svc, "")
	default:
		status.ResourceStatusReady.True(svc, "")
	}

	svc.Status.SetSummary(status.WalkResource(&svc.Status))

	return modelClient.Resources().UpdateOne(svc).
		SetStatus(svc.Status).
		Exec(ctx)
}

func (h Handler) RouteLog(req RouteLogRequest) error {
	// NB(thxCode): disable timeout as we don't know the maximum time-cost of once tracing,
	// and rely on the session context timeout control,
	// which means we don't close the underlay kubernetes client operation until the `ctx` is cancel.
	restConfig := *h.kubeConfig // Copy.
	restConfig.Timeout = 0

	cli, err := coreclient.NewForConfig(&restConfig)
	if err != nil {
		return fmt.Errorf("error creating kubernetes client: %w", err)
	}

	var (
		ctx context.Context
		out io.Writer
	)

	if req.Stream != nil {
		// In stream.
		ctx = req.Stream
		out = req.Stream
	} else {
		ctx = req.Context
		out = req.Context.Writer
	}

	return terraform.StreamJobLogs(ctx, terraform.StreamJobLogsOptions{
		Cli:        cli,
		RevisionID: req.ID,
		JobType:    req.JobType,
		Out:        out,
	})
}

// RouteGetDiffLatest get the revision with the service latest revision diff.
func (h Handler) RouteGetDiffLatest(req RouteGetDiffLatestRequest) (*RouteGetDiffLatestResponse, error) {
	compareRevision, err := h.modelClient.ResourceRevisions().Query().
		Select(
			resourcerevision.FieldID,
			resourcerevision.FieldResourceID,
			resourcerevision.FieldTemplateName,
			resourcerevision.FieldTemplateVersion,
			resourcerevision.FieldAttributes,
		).
		Where(resourcerevision.ID(req.ID)).
		Order(model.Desc(resourcerevision.FieldCreateTime)).
		Only(req.Context)
	if err != nil {
		return nil, err
	}

	latestRevision, err := h.modelClient.ResourceRevisions().Query().
		Select(
			resourcerevision.FieldID,
			resourcerevision.FieldTemplateName,
			resourcerevision.FieldTemplateVersion,
			resourcerevision.FieldAttributes,
		).
		Where(resourcerevision.ResourceID(compareRevision.ResourceID)).
		Order(model.Desc(resourcerevision.FieldCreateTime)).
		First(req.Context)
	if err != nil {
		return nil, err
	}

	return &RouteGetDiffLatestResponse{
		Old: RevisionDiff{
			TemplateName:    latestRevision.TemplateName,
			TemplateVersion: latestRevision.TemplateVersion,
			Attributes:      latestRevision.Attributes,
		},
		New: RevisionDiff{
			TemplateName:    compareRevision.TemplateName,
			TemplateVersion: compareRevision.TemplateVersion,
			Attributes:      compareRevision.Attributes,
		},
	}, nil
}

// RouteGetDiffPrevious get the revision with the service previous revision diff.
func (h Handler) RouteGetDiffPrevious(req RouteGetDiffPreviousRequest) (*RouteGetDiffPreviousResponse, error) {
	compareRevision, err := h.modelClient.ResourceRevisions().Query().
		Select(
			resourcerevision.FieldID,
			resourcerevision.FieldTemplateName,
			resourcerevision.FieldTemplateVersion,
			resourcerevision.FieldAttributes,
			resourcerevision.FieldResourceID,
			resourcerevision.FieldCreateTime,
		).
		Where(resourcerevision.ID(req.ID)).
		Order(model.Desc(resourcerevision.FieldCreateTime)).
		Only(req.Context)
	if err != nil {
		return nil, err
	}

	var old RevisionDiff

	previousRevision, err := h.modelClient.ResourceRevisions().Query().
		Select(
			resourcerevision.FieldID,
			resourcerevision.FieldTemplateName,
			resourcerevision.FieldTemplateVersion,
			resourcerevision.FieldAttributes,
		).
		Where(
			resourcerevision.ResourceID(compareRevision.ResourceID),
			resourcerevision.CreateTimeLT(*compareRevision.CreateTime),
		).
		Order(model.Desc(resourcerevision.FieldCreateTime)).
		First(req.Context)
	if err != nil && !model.IsNotFound(err) {
		return nil, err
	}

	if previousRevision != nil {
		old = RevisionDiff{
			TemplateName:    previousRevision.TemplateName,
			TemplateVersion: previousRevision.TemplateVersion,
			Attributes:      previousRevision.Attributes,
		}
	}

	return &RouteGetDiffPreviousResponse{
		Old: old,
		New: RevisionDiff{
			TemplateName:    compareRevision.TemplateName,
			TemplateVersion: compareRevision.TemplateVersion,
			Attributes:      compareRevision.Attributes,
		},
	}, nil
}
