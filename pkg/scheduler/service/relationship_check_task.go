package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
	"go.uber.org/multierr"
	"k8s.io/client-go/rest"

	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/service"
	"github.com/seal-io/seal/pkg/dao/model/servicerelationship"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/deployer"
	deployertf "github.com/seal-io/seal/pkg/deployer/terraform"
	deptypes "github.com/seal-io/seal/pkg/deployer/types"
	pkgservice "github.com/seal-io/seal/pkg/service"
	"github.com/seal-io/seal/utils/log"
	"github.com/seal-io/seal/utils/strs"
)

const (
	summaryStatusProgressing = "Progressing"
	summaryStatusDeleting    = "Deleting"
)

// RelationshipCheckTask checks services pending on relationships and
// proceeds applying/destroying services when the check pass.
type RelationshipCheckTask struct {
	mu sync.Mutex

	skipTLSVerify bool
	logger        log.Logger
	modelClient   model.ClientSet
	kubeConfig    *rest.Config
	deployer      deptypes.Deployer
}

func NewServiceRelationshipCheckTask(
	mc model.ClientSet,
	kc *rest.Config,
	skipTLSVerify bool,
) (*RelationshipCheckTask, error) {
	in := &RelationshipCheckTask{
		modelClient:   mc,
		kubeConfig:    kc,
		skipTLSVerify: skipTLSVerify,
	}

	in.logger = log.WithName("task").WithName(in.Name())

	return in, nil
}

func (in *RelationshipCheckTask) Name() string {
	return "service-relationship-check"
}

func (in *RelationshipCheckTask) Process(ctx context.Context, args ...interface{}) error {
	if !in.mu.TryLock() {
		in.logger.Warn("previous processing is not finished")
		return nil
	}

	startTs := time.Now()

	defer func() {
		in.mu.Unlock()
		in.logger.Debugf("processed in %v", time.Since(startTs))
	}()

	var (
		checkers = []func(context.Context) error{
			in.applyServices,
			in.destroyServices,
		}
		err error
	)

	for i := range checkers {
		// Don't return directly when error occurs,
		// but records it and continue to handle the next connector,
		// the final error collect all errors,
		// and reports this time task running as failure at observing.
		err = multierr.Append(err, checkers[i](ctx))

		if ctx.Err() != nil {
			// Give up the loop if the context is canceled.
			break
		}
	}

	return err
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

	dp, err := in.getDeployer(ctx)
	if err != nil {
		return err
	}

	destroyOpts := pkgservice.Options{
		TlsCertified: in.skipTLSVerify,
	}

	for _, svc := range services {
		if status.ServiceStatusProgressing.IsTrue(svc) {
			// Dependencies resolved and destruction in progress.
			continue
		}

		if err = pkgservice.Destroy(ctx, in.modelClient, dp, svc, destroyOpts); err != nil {
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

	serviceIDs := make([]oid.ID, 0, len(dependencies))

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
		if !pkgservice.IsStatusReady(depSvc) {
			if err := in.setServiceStatusFalse(ctx, svc, depSvc); err != nil {
				return false, err
			}

			return false, nil
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

	dp, err := in.getDeployer(ctx)
	if err != nil {
		return err
	}

	deployOpts := pkgservice.Options{
		TlsCertified: in.skipTLSVerify,
	}

	return pkgservice.Apply(ctx, in.modelClient, dp, entity, deployOpts)
}

// setServiceStatusFalse sets a service status to false if parent dependencies statuses are false or deleted.
func (in *RelationshipCheckTask) setServiceStatusFalse(
	ctx context.Context,
	svc, parentService *model.Service,
) (err error) {
	if pkgservice.IsStatusFalse(parentService) {
		status.ServiceStatusProgressing.False(
			svc,
			fmt.Sprintf("Parent service status is false, service name: %s", parentService.Name),
		)
		svc.Status.SetSummary(status.WalkService(&svc.Status))

		err = pkgservice.UpdateStatus(ctx, in.modelClient, svc)
		if err != nil {
			return err
		}
	}

	if pkgservice.IsStatusDeleted(parentService) {
		status.ServiceStatusProgressing.False(svc,
			fmt.Sprintf("Parent service status is deleted, service name: %s", parentService.Name),
		)
		svc.Status.SetSummary(status.WalkService(&svc.Status))

		err = pkgservice.UpdateStatus(ctx, in.modelClient, svc)
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}

func (in *RelationshipCheckTask) getDeployer(ctx context.Context) (deptypes.Deployer, error) {
	if in.deployer != nil {
		return in.deployer, nil
	}

	createOpts := deptypes.CreateOptions{
		Type:        deployertf.DeployerType,
		ModelClient: in.modelClient,
		KubeConfig:  in.kubeConfig,
	}

	dp, err := deployer.Get(ctx, createOpts)
	if err != nil {
		return nil, err
	}

	return dp, nil
}
