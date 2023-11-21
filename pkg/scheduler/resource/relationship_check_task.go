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
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/deployer"
	deployertf "github.com/seal-io/walrus/pkg/deployer/terraform"
	deptypes "github.com/seal-io/walrus/pkg/deployer/types"
	pkgresource "github.com/seal-io/walrus/pkg/resource"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/strs"
)

const (
	summaryStatusProgressing = "Progressing"
	summaryStatusDeleting    = "Deleting"
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
		Type:        deployertf.DeployerType,
		ModelClient: mc,
		KubeConfig:  kc,
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
	resources, err := in.modelClient.Resources().Query().
		Where(
			func(s *sql.Selector) {
				s.Where(sqljson.ValueEQ(
					resource.FieldStatus,
					summaryStatusProgressing,
					sqljson.Path("summaryStatus"),
				))
			},
			func(s *sql.Selector) {
				s.Where(sqljson.ValueEQ(
					resource.FieldStatus,
					true,
					sqljson.Path("transitioning"),
				))
			},
		).
		All(ctx)
	if err != nil && !model.IsNotFound(err) {
		return err
	}

	for _, svc := range resources {
		if status.ResourceStatusStopped.Exist(svc) {
			continue
		}

		ok, err := in.checkDependencies(ctx, svc)
		if err != nil {
			return err
		}

		if !ok {
			continue
		}

		// Deploy.
		err = in.deployResource(ctx, svc)
		if err != nil {
			return err
		}
	}

	return nil
}

func (in *RelationshipCheckTask) destroyResources(ctx context.Context) error {
	resources, err := in.modelClient.Resources().Query().
		Where(
			func(s *sql.Selector) {
				s.Where(sqljson.ValueEQ(
					resource.FieldStatus,
					summaryStatusDeleting,
					sqljson.Path("summaryStatus"),
				))
			},
		).
		All(ctx)
	if err != nil && !model.IsNotFound(err) {
		return err
	}

	for _, res := range resources {
		if status.ResourceStatusProgressing.IsTrue(res) {
			// Dependencies resolved and destruction in progress.
			continue
		}

		err = pkgresource.Destroy(ctx, in.modelClient, res, pkgresource.Options{
			Deployer: in.deployer,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// stopResources stops all resources that are in the progressing and stopping state.
func (in *RelationshipCheckTask) stopResources(ctx context.Context) error {
	resources, err := in.modelClient.Resources().Query().
		Where(
			func(s *sql.Selector) {
				s.Where(sqljson.ValueEQ(
					resource.FieldStatus,
					summaryStatusProgressing,
					sqljson.Path("summaryStatus"),
				))
			},
		).
		All(ctx)
	if err != nil && !model.IsNotFound(err) {
		return err
	}

	for _, res := range resources {
		if !status.ResourceStatusStopped.IsUnknown(res) {
			continue
		}

		err = pkgresource.Stop(ctx, in.modelClient, res, pkgresource.Options{
			Deployer: in.deployer,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (in *RelationshipCheckTask) checkDependencies(ctx context.Context, svc *model.Resource) (bool, error) {
	dependencies, err := in.modelClient.ResourceRelationships().Query().
		Where(
			resourcerelationship.ResourceIDEQ(svc.ID),
			resourcerelationship.DependencyIDNEQ(svc.ID),
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

	for _, depSvc := range dependencyResources {
		if pkgresource.IsStatusReady(depSvc) {
			continue
		}

		err = in.setResourceStatusFalse(ctx, svc, depSvc)
		if err != nil {
			return false, err
		}

		return false, nil
	}

	return true, nil
}

func (in *RelationshipCheckTask) deployResource(ctx context.Context, entity *model.Resource) error {
	// Reset status.
	status.ResourceStatusDeployed.Reset(entity, "")

	err := pkgresource.UpdateStatus(ctx, in.modelClient, entity)
	if err != nil {
		return err
	}

	return pkgresource.Apply(ctx, in.modelClient, entity, pkgresource.Options{
		Deployer: in.deployer,
	})
}

// setResourceStatusFalse sets a resource status to false if parent dependencies statuses are false or deleted.
func (in *RelationshipCheckTask) setResourceStatusFalse(
	ctx context.Context,
	svc, parentResource *model.Resource,
) error {
	if pkgresource.IsStatusFalse(parentResource) {
		status.ResourceStatusProgressing.False(
			svc,
			fmt.Sprintf("Parent resource status is false, name: %s", parentResource.Name),
		)

		err := pkgresource.UpdateStatus(ctx, in.modelClient, svc)
		if err != nil {
			return err
		}
	}

	if pkgresource.IsStatusDeleted(parentResource) {
		status.ResourceStatusProgressing.False(svc,
			fmt.Sprintf("Parent resource status is deleted, name: %s", parentResource.Name),
		)

		err := pkgresource.UpdateStatus(ctx, in.modelClient, svc)
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}
