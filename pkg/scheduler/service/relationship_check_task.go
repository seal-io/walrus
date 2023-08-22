package service

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
	"go.uber.org/multierr"
	"k8s.io/client-go/rest"

	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/service"
	"github.com/seal-io/walrus/pkg/dao/model/servicerelationship"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/deployer"
	deployertf "github.com/seal-io/walrus/pkg/deployer/terraform"
	deptypes "github.com/seal-io/walrus/pkg/deployer/types"
	pkgservice "github.com/seal-io/walrus/pkg/service"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/strs"
)

const (
	summaryStatusProgressing = "Progressing"
	summaryStatusDeleting    = "Deleting"
)

// RelationshipCheckTask checks services pending on relationships and
// proceeds applying/destroying services when the check pass.
type RelationshipCheckTask struct {
	logger        log.Logger
	modelClient   model.ClientSet
	skipTLSVerify bool
	deployer      deptypes.Deployer
}

func NewServiceRelationshipCheckTask(
	logger log.Logger,
	mc model.ClientSet,
	kc *rest.Config,
	skipTLSVerify bool,
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
		logger:        logger,
		modelClient:   mc,
		skipTLSVerify: skipTLSVerify,
		deployer:      dp,
	}

	return
}

func (in *RelationshipCheckTask) Process(ctx context.Context, args ...any) error {
	checkers := []func(context.Context) error{
		in.applyServices,
		in.destroyServices,
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

// applyServices applies all services that are in the progressing state.
func (in *RelationshipCheckTask) applyServices(ctx context.Context) error {
	services, err := in.modelClient.Services().Query().
		Where(
			func(s *sql.Selector) {
				s.Where(sqljson.ValueEQ(
					service.FieldStatus,
					summaryStatusProgressing,
					sqljson.Path("summaryStatus"),
				))
			},
			func(s *sql.Selector) {
				s.Where(sqljson.ValueEQ(
					service.FieldStatus,
					true,
					sqljson.Path("transitioning"),
				))
			},
		).
		All(ctx)
	if err != nil && !model.IsNotFound(err) {
		return err
	}

	for _, svc := range services {
		ok, err := in.checkDependencies(ctx, svc)
		if err != nil {
			return err
		}

		if !ok {
			continue
		}

		// Deploy.
		err = in.deployService(ctx, svc)
		if err != nil {
			return err
		}
	}

	return nil
}

func (in *RelationshipCheckTask) destroyServices(ctx context.Context) error {
	services, err := in.modelClient.Services().Query().
		Where(
			func(s *sql.Selector) {
				s.Where(sqljson.ValueEQ(
					service.FieldStatus,
					summaryStatusDeleting,
					sqljson.Path("summaryStatus"),
				))
			},
		).
		All(ctx)
	if err != nil && !model.IsNotFound(err) {
		return err
	}

	opts := pkgservice.Options{
		TlsCertified: in.skipTLSVerify,
	}

	for _, svc := range services {
		if status.ServiceStatusProgressing.IsTrue(svc) {
			// Dependencies resolved and destruction in progress.
			continue
		}

		err = pkgservice.Destroy(ctx, in.modelClient, in.deployer, svc, opts)
		if err != nil {
			return err
		}
	}

	return nil
}

func (in *RelationshipCheckTask) checkDependencies(ctx context.Context, svc *model.Service) (bool, error) {
	dependencies, err := in.modelClient.ServiceRelationships().Query().
		Where(
			servicerelationship.ServiceIDEQ(svc.ID),
			servicerelationship.DependencyIDNEQ(svc.ID),
			servicerelationship.Type(types.ServiceRelationshipTypeImplicit),
		).
		All(ctx)
	if err != nil && !model.IsNotFound(err) {
		return false, err
	}

	if len(dependencies) == 0 {
		return true, nil
	}

	serviceIDs := make([]object.ID, 0, len(dependencies))

	for _, d := range dependencies {
		if existCycle := dao.ServiceRelationshipCheckCycle(d); existCycle {
			pathIDs := make([]string, 0, len(d.Path))
			for _, id := range d.Path {
				pathIDs = append(pathIDs, id.String())
			}

			pathStr := strs.Join(" -> ", pathIDs...)

			return false, fmt.Errorf("dependency cycle detected, service id: %s, path: %s", d.ServiceID, pathStr)
		}

		serviceIDs = append(serviceIDs, d.DependencyID)
	}

	dependencyServices, err := in.modelClient.Services().Query().
		Where(service.IDIn(serviceIDs...)).
		All(ctx)
	if err != nil {
		return false, err
	}

	for _, depSvc := range dependencyServices {
		if pkgservice.IsStatusReady(depSvc) {
			continue
		}

		err = in.setServiceStatusFalse(ctx, svc, depSvc)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func (in *RelationshipCheckTask) deployService(ctx context.Context, entity *model.Service) error {
	// Reset status.
	status.ServiceStatusDeployed.Reset(entity, "")

	err := pkgservice.UpdateStatus(ctx, in.modelClient, entity)
	if err != nil {
		return err
	}

	opts := pkgservice.Options{
		TlsCertified: in.skipTLSVerify,
	}

	return pkgservice.Apply(ctx, in.modelClient, in.deployer, entity, opts)
}

// setServiceStatusFalse sets a service status to false if parent dependencies statuses are false or deleted.
func (in *RelationshipCheckTask) setServiceStatusFalse(
	ctx context.Context,
	svc, parentService *model.Service,
) error {
	if pkgservice.IsStatusFalse(parentService) {
		status.ServiceStatusProgressing.False(
			svc,
			fmt.Sprintf("Parent service status is false, service name: %s", parentService.Name),
		)

		err := pkgservice.UpdateStatus(ctx, in.modelClient, svc)
		if err != nil {
			return err
		}
	}

	if pkgservice.IsStatusDeleted(parentService) {
		status.ServiceStatusProgressing.False(svc,
			fmt.Sprintf("Parent service status is deleted, service name: %s", parentService.Name),
		)

		err := pkgservice.UpdateStatus(ctx, in.modelClient, svc)
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}
