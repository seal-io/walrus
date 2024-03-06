package resource

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
	"go.uber.org/multierr"
	"k8s.io/client-go/rest"

	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/model/resourcerelationship"
	"github.com/seal-io/walrus/pkg/dao/model/resourcerun"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/deployer"
	deptypes "github.com/seal-io/walrus/pkg/deployer/types"
	runstatus "github.com/seal-io/walrus/pkg/resourceruns/status"
	pkgresource "github.com/seal-io/walrus/pkg/resources"
	resstatus "github.com/seal-io/walrus/pkg/resources/status"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/strs"
)

// RelationshipCheckTask checks resources pending on relationships and
// proceeds applying/destroying resources when the check pass.
type RelationshipCheckTask struct {
	logger      log.Logger
	modelClient model.ClientSet
	deployer    deptypes.Deployer
}

func NewResourceRelationshipCheckTask(
	logger log.Logger,
	mc model.ClientSet,
	kc *rest.Config,
) (in *RelationshipCheckTask, err error) {
	// Create deployer.
	opts := deptypes.CreateOptions{
		Type:       types.DeployerTypeTF,
		KubeConfig: kc,
	}

	dp, err := deployer.Get(context.Background(), opts)
	if err != nil {
		return nil, err
	}

	in = &RelationshipCheckTask{
		logger:      logger,
		modelClient: mc,
		deployer:    dp,
	}

	return
}

func (in *RelationshipCheckTask) Process(ctx context.Context, args ...any) error {
	checkers := []func(context.Context) error{
		in.applyResources,
		in.destroyResources,
		in.stopResources,
	}

	// Merge the errors to return them all at once,
	// instead of returning the first error.
	var berr error

	for i := range checkers {
		berr = multierr.Append(berr, checkers[i](ctx))

		// Give up the loop if the context is canceled.
		if multierr.AppendInto(&berr, ctx.Err()) {
			break
		}
	}

	return berr
}

// applyResources applies all resources that are in the progressing state.
func (in *RelationshipCheckTask) applyResources(ctx context.Context) error {
	resources, err := in.getPendingRunResources(
		ctx,
		types.RunTypeCreate.String(),
		types.RunTypeUpdate.String(),
		types.RunTypeRollback.String(),
		types.RunTypeStart.String(),
	)
	if err != nil {
		return err
	}

	for _, res := range resources {
		ok, err := in.checkDependencies(ctx, res)
		if err != nil {
			return err
		}

		if !ok {
			continue
		}

		err = pkgresource.PerformResource(ctx, in.modelClient, in.deployer, res.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (in *RelationshipCheckTask) destroyResources(ctx context.Context) error {
	resources, err := in.getPendingRunResources(ctx, types.RunTypeDelete.String())
	if err != nil && !model.IsNotFound(err) {
		return err
	}

	for _, res := range resources {
		ok, err := in.checkDependants(ctx, res)
		if err != nil {
			return err
		}

		if !ok {
			continue
		}

		err = pkgresource.PerformResource(ctx, in.modelClient, in.deployer, res.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

// stopResources stops all resources that are in the progressing and stopping state.
func (in *RelationshipCheckTask) stopResources(ctx context.Context) error {
	resources, err := in.getPendingRunResources(ctx, types.RunTypeStop.String())
	if err != nil && !model.IsNotFound(err) {
		return err
	}

	for _, res := range resources {
		if !status.ResourceStatusStopped.IsUnknown(res) {
			continue
		}

		ok, err := in.checkDependants(ctx, res)
		if err != nil {
			return err
		}

		if !ok {
			continue
		}

		err = pkgresource.PerformResource(ctx, in.modelClient, in.deployer, res.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (in *RelationshipCheckTask) checkDependencies(ctx context.Context, res *model.Resource) (bool, error) {
	dependencies, err := in.modelClient.ResourceRelationships().Query().
		Where(
			resourcerelationship.ResourceIDEQ(res.ID),
			resourcerelationship.DependencyIDNEQ(res.ID),
			resourcerelationship.Type(types.ResourceRelationshipTypeImplicit),
		).
		All(ctx)
	if err != nil && !model.IsNotFound(err) {
		return false, err
	}

	if len(dependencies) == 0 {
		return true, nil
	}

	resourceIDs := make([]object.ID, 0, len(dependencies))

	for _, d := range dependencies {
		if existCycle := dao.ResourceRelationshipCheckCycle(d); existCycle {
			pathIDs := make([]string, 0, len(d.Path))
			for _, id := range d.Path {
				pathIDs = append(pathIDs, id.String())
			}

			pathStr := strs.Join(" -> ", pathIDs...)

			return false, fmt.Errorf("dependency cycle detected, resource id: %s, path: %s", d.ResourceID, pathStr)
		}

		resourceIDs = append(resourceIDs, d.DependencyID)
	}

	dependencyResources, err := in.modelClient.Resources().Query().
		Where(resource.IDIn(resourceIDs...)).
		All(ctx)
	if err != nil {
		return false, err
	}

	for _, depRes := range dependencyResources {
		if resstatus.IsStatusReady(depRes) {
			continue
		}

		err = in.setResourceStatusFalse(ctx, res, depRes)
		if err != nil {
			return false, err
		}

		return false, nil
	}

	return true, nil
}

// checkDependants checks if the resource has dependants and if the dependants are ready.
func (in *RelationshipCheckTask) checkDependants(ctx context.Context, res *model.Resource) (bool, error) {
	dependants, err := dao.GetResourceDependantResource(ctx, in.modelClient, res.ID)
	if err != nil {
		return false, err
	}

	for i := range dependants {
		ok, err := in.checkDependantResourceStatus(ctx, res, dependants[i])
		if err != nil {
			return false, err
		}

		if !ok {
			return false, nil
		}
	}

	return true, nil
}

// setResourceStatusFalse sets a resource status to false if parent dependencies statuses are false or deleted.
func (in *RelationshipCheckTask) setResourceStatusFalse(
	ctx context.Context,
	res, parentResource *model.Resource,
) error {
	if resstatus.IsStatusError(parentResource) {
		errMsg := fmt.Sprintf("Dependency resource %q has encountered an error, please check it",
			parentResource.Name)

		return setResourceStatusFalse(ctx, in.modelClient, res, errMsg)
	}

	if resstatus.IsStatusDeleted(parentResource) {
		errMsg := fmt.Sprintf("Dependency resource %q is in delete status, please check it",
			parentResource.Name)

		return setResourceStatusFalse(ctx, in.modelClient, res, errMsg)
	}

	if resstatus.IsStatusStopped(parentResource) {
		errMsg := fmt.Sprintf("Dependency resource %q is in stop status, please check it",
			parentResource.Name)

		return setResourceStatusFalse(ctx, in.modelClient, res, errMsg)
	}

	return nil
}

// checkDependantResourceStatus sets a resource status to false if dependant resource status is false or deployed.
func (in *RelationshipCheckTask) checkDependantResourceStatus(
	ctx context.Context,
	res, dependantResource *model.Resource,
) (bool, error) {
	if resstatus.IsStatusError(dependantResource) {
		errMsg := fmt.Sprintf("Dependant resource %q has encountered an error, please check it",
			dependantResource.Name)

		err := setResourceStatusFalse(ctx, in.modelClient, res, errMsg)
		if err != nil {
			return false, err
		}

		return false, nil
	}

	// If the dependant resource is deployed or to be deployed, the resource cannot be deleted or stopped.
	if resstatus.IsStatusDeployed(dependantResource) {
		errMsg := fmt.Sprintf("Dependant resource %q is in deploy status, please check it",
			dependantResource.Name)

		err := setResourceStatusFalse(ctx, in.modelClient, res, errMsg)
		if err != nil {
			return false, err
		}

		return false, nil
	}

	return true, nil
}

// setResourceStatusFalse sets resource and resource run status to false.
func setResourceStatusFalse(ctx context.Context, mc model.ClientSet, res *model.Resource, errMsg string) error {
	runs, err := dao.GetResourcesLatestRuns(ctx, mc, res.ID)
	if err != nil {
		return err
	}

	if len(runs) == 0 {
		return fmt.Errorf("no runs found for resource %s", res.ID)
	}

	run := runs[0]

	switch types.RunType(run.Type) {
	case types.RunTypeCreate, types.RunTypeUpdate, types.RunTypeRollback, types.RunTypeStart:
		status.ResourceStatusDeployed.False(res, errMsg)
	case types.RunTypeDelete:
		status.ResourceStatusDeleted.False(res, errMsg)
	case types.RunTypeStop:
		status.ResourceStatusStopped.False(res, errMsg)
	default:
		return fmt.Errorf("unsupported action type: %s", res.Type)
	}

	err = resstatus.UpdateStatus(ctx, mc, res)
	if err != nil {
		return err
	}

	runstatus.SetStatusFalse(run, errMsg)
	_, err = runstatus.UpdateStatus(ctx, mc, run)

	return err
}

// getPendingRuns Retrieve resources from all runs that are pending.
func (in *RelationshipCheckTask) getPendingRunResources(ctx context.Context, runTypes ...string) ([]*model.Resource, error) {
	return in.modelClient.ResourceRuns().Query().
		Where(
			resourcerun.TypeIn(runTypes...),
			func(s *sql.Selector) {
				s.Where(sqljson.ValueEQ(
					resourcerun.FieldStatus,
					status.ResourceRunStatusPending.String(),
					sqljson.Path("summaryStatus"),
				))
			},
			func(s *sql.Selector) {
				s.Where(sqljson.ValueNEQ(
					resourcerun.FieldStatus,
					"",
					sqljson.Path("summaryStatusMessage"),
				))
			},
		).
		QueryResource().
		All(ctx)
}
