package resourcerun

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"time"

	"k8s.io/apimachinery/pkg/util/sets"
	coreclient "k8s.io/client-go/kubernetes/typed/core/v1"

	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/connector"
	"github.com/seal-io/walrus/pkg/dao/model/environment"
	"github.com/seal-io/walrus/pkg/dao/model/project"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/model/resourcecomponent"
	"github.com/seal-io/walrus/pkg/dao/model/resourcerun"
	"github.com/seal-io/walrus/pkg/dao/model/resourcestate"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/deployer"
	"github.com/seal-io/walrus/pkg/deployer/terraform"
	deptypes "github.com/seal-io/walrus/pkg/deployer/types"
	"github.com/seal-io/walrus/pkg/operator"
	optypes "github.com/seal-io/walrus/pkg/operator/types"
	"github.com/seal-io/walrus/pkg/resourcecomponents"
	pkgrun "github.com/seal-io/walrus/pkg/resourceruns"
	runstatus "github.com/seal-io/walrus/pkg/resourceruns/status"
	"github.com/seal-io/walrus/pkg/servervars"
	"github.com/seal-io/walrus/pkg/settings"
	tfparser "github.com/seal-io/walrus/pkg/terraform/parser"
	"github.com/seal-io/walrus/utils/errorx"
	"github.com/seal-io/walrus/utils/gopool"
	"github.com/seal-io/walrus/utils/json"
	"github.com/seal-io/walrus/utils/log"
)

// RouteGetTerraformStates get the terraform states of the service run deployment.
func (h Handler) RouteGetTerraformStates(
	req RouteGetTerraformStatesRequest,
) (RouteGetTerraformStatesResponse, error) {
	entity, err := h.modelClient.ResourceRuns().Query().
		Where(resourcerun.ID(req.ID)).
		WithResource(func(rq *model.ResourceQuery) {
			rq.WithState()
		}).
		Only(req.Context)
	if err != nil {
		return nil, err
	}

	if entity.Edges.Resource.Edges.State == nil || entity.Edges.Resource.Edges.State.Data == "" {
		return nil, nil
	}
	stateData := entity.Edges.Resource.Edges.State.Data

	return RouteGetTerraformStatesResponse(stateData), nil
}

// RouteUpdateTerraformStates update the terraform states of the service run deployment.
func (h Handler) RouteUpdateTerraformStates(
	req RouteUpdateTerraformStatesRequest,
) (err error) {
	logger := log.WithName("api").WithName("resource-run")

	entity, err := h.modelClient.ResourceRuns().Query().
		Where(resourcerun.ID(req.ID)).
		Select(
			resourcerun.FieldID,
			resourcerun.FieldProjectID,
			resourcerun.FieldEnvironmentID,
			resourcerun.FieldResourceID,
			resourcerun.FieldStatus,
			resourcerun.FieldRecord,
			resourcerun.FieldDeployerType).
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

	state, err := h.modelClient.ResourceStates().Query().
		Where(resourcestate.ResourceID(entity.ResourceID)).
		Only(req.Context)
	if err != nil {
		return err
	}

	state.Data = string(req.RawMessage)

	err = h.modelClient.ResourceStates().UpdateOne(state).
		SetData(state.Data).
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

		runstatus.SetStatusFalse(entity, err.Error())

		if _, uerr := runstatus.UpdateStatus(updateCtx, h.modelClient, entity); uerr != nil {
			logger.Errorf("update status failed: %v", err)
		}
	}()

	return manageResourceComponentsAndEndpoints(req.Context, h.modelClient, entity, state.Data)
}

// manageResourceComponentsAndEndpoints parses and updates the resource components/endpoints,
// and execute reconcileResourceComponents for the new created resource components.
func manageResourceComponentsAndEndpoints(
	ctx context.Context,
	mc model.ClientSet,
	entity *model.ResourceRun,
	stateData string,
) error {
	var p tfparser.StateParser

	mappedOutputs, err := p.GetOriginalOutputsMap(stateData, entity.Edges.Resource.Name)
	if err != nil {
		return fmt.Errorf("error getting original outputs: %w", err)
	}

	observedComponents, dependencies, err := p.GetComponentsAndExtractDependencies(ctx, mc, entity)
	if err != nil {
		return fmt.Errorf("error getting resources and dependencies: %w", err)
	}

	// Get record components from local.
	recordComponents, err := mc.ResourceComponents().Query().
		Where(resourcecomponent.ResourceID(entity.ResourceID)).
		All(ctx)
	if err != nil {
		return err
	}

	// Calculate creating list, deleting list and updating list.
	var (
		deleteComponentIDs = make([]object.ID, 0, len(recordComponents))
		updatedComponents  = make([]*model.ResourceComponent, 0, len(recordComponents))

		observedComponentsIndex = dao.ResourceComponentToMap(observedComponents)
	)

	for _, c := range recordComponents {
		k := dao.ResourceComponentGetUniqueKey(c)
		if observedComponentsIndex[k] != nil {
			c.Edges.Instances = observedComponentsIndex[k].Edges.Instances
			updatedComponents = append(updatedComponents, c)

			delete(observedComponentsIndex, k)

			continue
		}

		deleteComponentIDs = append(deleteComponentIDs, c.ID)
	}

	createComponents := make([]*model.ResourceComponent, 0, len(observedComponentsIndex))

	for k := range observedComponentsIndex {
		// Component instances will be created through edges.
		if observedComponentsIndex[k].Shape != types.ResourceComponentShapeClass {
			continue
		}

		createComponents = append(createComponents, observedComponentsIndex[k])
	}

	// Diff by transactional session.
	replacedComponents := make([]*model.ResourceComponent, 0)

	// TODO(alex): refactor the following codes, make it more readable.
	err = mc.WithTx(ctx, func(tx *model.Tx) error {
		// Update components with new items.
		for _, r := range updatedComponents {
			rp, err := dao.ResourceComponentInstancesEdgeSaveWithResult(ctx, tx, r)
			if err != nil {
				return err
			}

			replacedComponents = append(replacedComponents, rp...)
		}

		// Some components may be removed when updating,
		// make sure the components still exist.
		recordComponents, err = dao.GetCleanResourceComponents(ctx, tx, recordComponents)
		if err != nil {
			return err
		}

		// Create new components.
		if len(createComponents) != 0 {
			createComponents, err = tx.ResourceComponents().CreateBulk().
				Set(createComponents...).
				SaveE(ctx, dao.ResourceComponentInstancesEdgeSave)
			if err != nil {
				return err
			}

			// TODO(thxCode): move the following codes into DAO.
			err = dao.ResourceComponentRelationshipUpdateWithDependencies(
				ctx,
				tx,
				dependencies,
				recordComponents,
				createComponents,
			)
			if err != nil {
				return err
			}
		}

		// Delete stale components.
		if len(deleteComponentIDs) != 0 {
			_, err = tx.ResourceComponents().Delete().
				Where(resourcecomponent.IDIn(deleteComponentIDs...)).
				Exec(ctx)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				return err
			}
		}

		// Parse endpoints from outputs,
		// the expected output is named by `walrus_endpoints` or `endpoints`,
		// and structures in `map[string]string` format.
		var m map[string]string

		for _, l := range []string{"walrus_endpoints", "endpoints"} {
			if v, ok := mappedOutputs[l]; ok {
				if err := json.Unmarshal(v.Value, &m); err == nil {
					break
				}
			}
		}

		// NB(thxCode): in order to access some endpoints pointed to local IP,
		// we can parse one-by-one to replace the local IP with the server's host,
		// especially be useful for deploying on the embedded Kubernetes cluster.
		if ips := servervars.NonLoopBackIPs.Get(); len(m) != 0 && ips.Len() != 0 {
			su := settings.ServeUrl.ShouldValueURL(ctx, tx)
			if su != nil {
				h := su.Hostname()
				// Replace the endpoint's url host with the server's host,
				// if the endpoint's url host is a local IP.
				for k, v := range m {
					u, _ := url.Parse(v)
					if u == nil || !ips.Has(u.Hostname()) {
						continue
					}
					u.Host = net.JoinHostPort(h, u.Port())
					m[k] = u.String()
				}
			}
		}

		// Update endpoints.
		err = tx.Resources().UpdateOneID(entity.ResourceID).
			SetEndpoints(types.ResourceEndpointsFromMap(m).Sort()).
			Exec(ctx)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	reconcileComponents := createComponents
	reconcileComponents = append(reconcileComponents, updatedComponents...)
	reconcileComponents = append(reconcileComponents, replacedComponents...)

	if len(reconcileComponents) == 0 {
		return nil
	}

	// Update the resource component status.
	reconcileComponentIndex := dao.ResourceComponentToMap(reconcileComponents)

	// Group components by connector ID,
	// and decorate them with the project/environment/component for latter labeling.
	connComponents := make(map[object.ID][]*model.ResourceComponent)

	for k := range reconcileComponentIndex {
		if reconcileComponentIndex[k].Shape != types.ResourceComponentShapeInstance {
			continue
		}

		reconcileComponentIndex[k].Edges.Project = entity.Edges.Project
		reconcileComponentIndex[k].Edges.Environment = entity.Edges.Environment
		reconcileComponentIndex[k].Edges.Resource = entity.Edges.Resource

		connComponents[reconcileComponentIndex[k].ConnectorID] = append(
			connComponents[reconcileComponentIndex[k].ConnectorID],
			reconcileComponentIndex[k],
		)
	}

	gopool.Go(func() {
		logger := log.WithName("api").WithName("resource-run")

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		if err = reconcileResourceComponents(ctx, mc, entity.ResourceID, connComponents); err != nil {
			logger.Errorf("error reconciling resource components: %v", err)
		}
	})

	return nil
}

// reconcileResourceComponents reconciles the resource components,
// including states/labels/composes resource components.
func reconcileResourceComponents(
	ctx context.Context,
	modelClient model.ClientSet,
	resourceID object.ID,
	connComponents map[object.ID][]*model.ResourceComponent,
) error {
	logger := log.WithName("api").WithName("resource-run")

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
		Where(connector.IDIn(sets.KeySet(connComponents).UnsortedList()...)).
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

	for cid, crs := range connComponents {
		op, exist := ops[cid]
		if !exist {
			// Ignore if not found operator.
			continue
		}

		// Discover components.
		ncrs, err := resourcecomponents.Discover(ctx, op, modelClient, crs)
		if err != nil {
			logger.Errorf("error discovering component resources: %v", err)
		}

		// State components.
		nsr, err := resourcecomponents.State(ctx, op, modelClient, append(crs, ncrs...))
		if err != nil {
			logger.Errorf("error stating resources: %v", err)
			// Mark error as transitioning,
			// which doesn't flip the status.
			nsr.Transitioning = true
		}

		sr.Merge(nsr)

		// Label components.
		err = resourcecomponents.Label(ctx, op, crs)
		if err != nil {
			logger.Errorf("error labeling resources: %v", err)
		}
	}

	// State resource.
	res, err := modelClient.Resources().Query().
		Where(resource.ID(resourceID)).
		Select(
			resource.FieldID,
			resource.FieldStatus).
		Only(ctx)
	if err != nil {
		return fmt.Errorf("cannot get resource: %w", err)
	}

	if status.ResourceStatusDeleted.Exist(res) ||
		status.ResourceStatusStopped.Exist(res) {
		// Skip if the service is on deleting or stopping.
		return nil
	}

	switch {
	case sr.Error:
		status.ResourceStatusReady.False(res, "")
	case sr.Transitioning:
		status.ResourceStatusReady.Unknown(res, "")
	default:
		status.ResourceStatusReady.True(res, "")
	}

	res.Status.SetSummary(status.WalkResource(&res.Status))

	return modelClient.Resources().UpdateOne(res).
		SetStatus(res.Status).
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
		Cli:     cli,
		RunID:   req.ID,
		JobType: req.JobType,
		Out:     out,
	})
}

// RouteGetDiffLatest get the run with the service latest run diff.
func (h Handler) RouteGetDiffLatest(req RouteGetDiffLatestRequest) (*RouteGetDiffLatestResponse, error) {
	compareRun, err := h.modelClient.ResourceRuns().Query().
		Select(
			resourcerun.FieldID,
			resourcerun.FieldResourceID,
			resourcerun.FieldTemplateName,
			resourcerun.FieldTemplateVersion,
			resourcerun.FieldAttributes,
			resourcerun.FieldComputedAttributes,
		).
		Where(resourcerun.ID(req.ID)).
		Order(model.Desc(resourcerun.FieldCreateTime)).
		Only(req.Context)
	if err != nil {
		return nil, err
	}

	latestRun, err := h.modelClient.ResourceRuns().Query().
		Select(
			resourcerun.FieldID,
			resourcerun.FieldTemplateName,
			resourcerun.FieldTemplateVersion,
			resourcerun.FieldAttributes,
			resourcerun.FieldComputedAttributes,
		).
		Where(resourcerun.ResourceID(compareRun.ResourceID)).
		Order(model.Desc(resourcerun.FieldCreateTime)).
		First(req.Context)
	if err != nil {
		return nil, err
	}

	return &RouteGetDiffLatestResponse{
		Old: RunDiff{
			TemplateName:       latestRun.TemplateName,
			TemplateVersion:    latestRun.TemplateVersion,
			Attributes:         latestRun.Attributes,
			ComputedAttributes: latestRun.ComputedAttributes,
		},
		New: RunDiff{
			TemplateName:       compareRun.TemplateName,
			TemplateVersion:    compareRun.TemplateVersion,
			Attributes:         compareRun.Attributes,
			ComputedAttributes: compareRun.ComputedAttributes,
		},
	}, nil
}

// RouteGetDiffPrevious get the run with the service previous run diff.
func (h Handler) RouteGetDiffPrevious(req RouteGetDiffPreviousRequest) (*RouteGetDiffPreviousResponse, error) {
	compareRun, err := h.modelClient.ResourceRuns().Query().
		Select(
			resourcerun.FieldID,
			resourcerun.FieldTemplateName,
			resourcerun.FieldTemplateVersion,
			resourcerun.FieldAttributes,
			resourcerun.FieldResourceID,
			resourcerun.FieldCreateTime,
			resourcerun.FieldComputedAttributes,
		).
		Where(resourcerun.ID(req.ID)).
		Order(model.Desc(resourcerun.FieldCreateTime)).
		Only(req.Context)
	if err != nil {
		return nil, err
	}

	var old RunDiff

	previousRun, err := h.modelClient.ResourceRuns().Query().
		Select(
			resourcerun.FieldID,
			resourcerun.FieldTemplateName,
			resourcerun.FieldTemplateVersion,
			resourcerun.FieldAttributes,
			resourcerun.FieldComputedAttributes,
		).
		Where(
			resourcerun.ResourceID(compareRun.ResourceID),
			resourcerun.CreateTimeLT(*compareRun.CreateTime),
		).
		Order(model.Desc(resourcerun.FieldCreateTime)).
		First(req.Context)
	if err != nil && !model.IsNotFound(err) {
		return nil, err
	}

	if previousRun != nil {
		old = RunDiff{
			TemplateName:       previousRun.TemplateName,
			TemplateVersion:    previousRun.TemplateVersion,
			Attributes:         previousRun.Attributes,
			ComputedAttributes: previousRun.ComputedAttributes,
		}
	}

	return &RouteGetDiffPreviousResponse{
		Old: old,
		New: RunDiff{
			TemplateName:       compareRun.TemplateName,
			TemplateVersion:    compareRun.TemplateVersion,
			Attributes:         compareRun.Attributes,
			ComputedAttributes: compareRun.ComputedAttributes,
		},
	}, nil
}

// RouteApply apply the planned run.
func (h Handler) RouteApply(req RouteApplyRequest) error {
	run, err := h.modelClient.ResourceRuns().Get(req.Context, req.ID)
	if err != nil {
		return err
	}

	if !runstatus.IsStatusPlanned(run) {
		return errorx.Errorf("can not approve a non-planned run: %s", run.Status.SummaryStatus)
	}

	if !req.Approve {
		status.ResourceRunStatusCanceled.True(run, "")
		run.Status.SetSummary(status.WalkResourceRun(&run.Status))

		_, err = runstatus.UpdateStatus(req.Context, h.modelClient, run)

		return err
	}

	dp, err := deployer.Get(req.Context, deptypes.CreateOptions{
		Type:       types.DeployerTypeTF,
		KubeConfig: h.kubeConfig,
	})
	if err != nil {
		return err
	}

	return pkgrun.Apply(req.Context, h.modelClient, dp, run)
}

func (h Handler) RouteSetPlan(req RouteSetPlanRequest) error {
	run, err := h.modelClient.ResourceRuns().Get(req.Context, req.ID)
	if err != nil {
		return err
	}

	jsonPlanHeader, err := req.Context.FormFile("jsonplan")
	if err != nil {
		return err
	}

	// Get change file from form change field.
	jsonPlanFile, err := jsonPlanHeader.Open()
	if err != nil {
		return err
	}
	defer jsonPlanFile.Close()

	jsonPlanBytes, err := io.ReadAll(jsonPlanFile)
	if err != nil {
		return err
	}

	var runPlanChanges *types.Plan
	if err = json.Unmarshal(jsonPlanBytes, &runPlanChanges); err != nil {
		return err
	}

	runPlanChanges.ResourceComponentChanges, err = resourcecomponents.FilterResourceComponentChange(
		req.Context,
		h.modelClient,
		run.ResourceID,
		runPlanChanges.ResourceComponentChanges,
	)
	if err != nil {
		return err
	}

	run.ComponentChangeSummary = runPlanChanges.GetResourceChangeSummary()
	run.ComponentChanges = runPlanChanges.ResourceComponentChanges

	err = h.modelClient.ResourceRuns().UpdateOne(run).
		SetComponentChangeSummary(run.ComponentChangeSummary).
		SetComponentChanges(run.ComponentChanges).
		Exec(req.Context)
	if err != nil {
		return err
	}

	planHeader, err := req.Context.FormFile("plan")
	if err != nil {
		return err
	}

	// Get plan file from form plan field.
	planFile, err := planHeader.Open()
	if err != nil {
		return err
	}
	defer planFile.Close()

	// Set plan file to storage.
	planBytes, err := io.ReadAll(planFile)
	if err != nil {
		return err
	}

	return h.storageManager.SetRunPlan(req.Context, run, planBytes)
}

func (h Handler) RouteGetPlan(req RouteGetPlanRequest) error {
	run, err := h.modelClient.ResourceRuns().Get(req.Context, req.ID)
	if err != nil {
		return err
	}

	// TODO Encrypt the plan file with key.
	plan, err := h.storageManager.GetRunPlan(req.Context, run)
	if err != nil {
		return err
	}

	req.Context.Writer.Header().Set("Content-Type", "application/zip")
	req.Context.Writer.Header().Set("Content-Disposition", "attachment; filename=plan.out")

	_, err = req.Context.Writer.Write(plan)

	return err
}
