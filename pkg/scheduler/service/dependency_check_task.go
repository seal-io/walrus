package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
	"k8s.io/client-go/rest"

	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/service"
	"github.com/seal-io/seal/pkg/dao/model/servicedependency"
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

const summaryStatusProgressing = "Progressing"

type DependencyCheckTask struct {
	mu sync.Mutex

	skipTLSVerify bool
	logger        log.Logger
	modelClient   model.ClientSet
	kubeConfig    *rest.Config
	deployer      deptypes.Deployer
}

func NewServiceDependencyCheckTask(
	mc model.ClientSet,
	kc *rest.Config,
	skipTLSVerify bool,
) (*DependencyCheckTask, error) {
	in := &DependencyCheckTask{
		modelClient:   mc,
		kubeConfig:    kc,
		skipTLSVerify: skipTLSVerify,
	}

	in.logger = log.WithName("task").WithName(in.Name())

	return in, nil
}

func (in *DependencyCheckTask) Name() string {
	return "service-dependency-check"
}

func (in *DependencyCheckTask) Process(ctx context.Context, args ...interface{}) error {
	if !in.mu.TryLock() {
		in.logger.Warn("previous processing is not finished")
		return nil
	}

	startTs := time.Now()

	defer func() {
		in.mu.Unlock()
		in.logger.Debugf("processed in %v", time.Since(startTs))
	}()

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

func (in *DependencyCheckTask) checkDependencies(ctx context.Context, svc *model.Service) (bool, error) {
	dependencies, err := in.modelClient.ServiceDependencies().Query().
		Where(
			servicedependency.ServiceIDEQ(svc.ID),
			servicedependency.Type(types.ServiceDependencyTypeImplicit),
		).
		All(ctx)
	if err != nil && !model.IsNotFound(err) {
		return false, err
	}

	if len(dependencies) == 0 {
		return true, nil
	}

	serviceIDs := make([]oid.ID, 0, len(dependencies))

	for _, dep := range dependencies {
		if existCycle := dao.CheckDependencyCycle(dep); existCycle {
			pathIDs := make([]string, 0, len(dep.Path))
			for _, id := range dep.Path {
				pathIDs = append(pathIDs, id.String())
			}

			pathStr := strs.Join(" -> ", pathIDs...)

			return false, fmt.Errorf("dependency cycle detected, service id: %s, path: %s", dep.ServiceID, pathStr)
		}

		serviceIDs = append(serviceIDs, dep.DependentID)
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

func (in *DependencyCheckTask) deployService(ctx context.Context, entity *model.Service) error {
	// Reset status.
	status.ServiceStatusDeployed.Reset(entity, "Deploying service")

	err := pkgservice.UpdateStatus(ctx, in.modelClient, entity)
	if err != nil {
		return err
	}

	dp, err := in.getDeployer(ctx, deptypes.CreateOptions{
		Type:        deployertf.DeployerType,
		ModelClient: in.modelClient,
		KubeConfig:  in.kubeConfig,
	})
	if err != nil {
		return err
	}

	deployOpts := pkgservice.Options{
		TlsCertified: in.skipTLSVerify,
	}

	return pkgservice.Apply(ctx, in.modelClient, dp, entity, deployOpts)
}

// setServiceStatusFalse sets a service status to false if parent dependencies statuses are false or deleted.
func (in *DependencyCheckTask) setServiceStatusFalse(
	ctx context.Context,
	svc, parentService *model.Service,
) (err error) {
	if pkgservice.IsStatusFalse(parentService) {
		status.ServiceStatusProgressing.False(
			svc,
			fmt.Sprintf("Parent service status is false, service id: %s", parentService.ID),
		)
		svc.Status.SetSummary(status.WalkService(&svc.Status))

		err = pkgservice.UpdateStatus(ctx, in.modelClient, svc)
		if err != nil {
			return err
		}
	}

	if pkgservice.IsStatusDeleted(parentService) {
		status.ServiceStatusProgressing.False(svc,
			fmt.Sprintf("Parent service status is deleted, service id: %s", parentService.ID),
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

func (in *DependencyCheckTask) getDeployer(
	ctx context.Context,
	opts deptypes.CreateOptions,
) (deptypes.Deployer, error) {
	if in.deployer != nil {
		return in.deployer, nil
	}

	createOpts := deptypes.CreateOptions{
		Type:        deployertf.DeployerType,
		ModelClient: opts.ModelClient,
		KubeConfig:  opts.KubeConfig,
	}

	dp, err := deployer.Get(ctx, createOpts)
	if err != nil {
		return nil, err
	}

	return dp, nil
}
