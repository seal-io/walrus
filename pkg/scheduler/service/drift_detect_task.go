package service

import (
	"context"
	"sync"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
	"k8s.io/client-go/rest"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/service"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/deployer"
	deployertf "github.com/seal-io/seal/pkg/deployer/terraform"
	deptypes "github.com/seal-io/seal/pkg/deployer/types"
	pkgservice "github.com/seal-io/seal/pkg/service"
	"github.com/seal-io/seal/pkg/settings"
	"github.com/seal-io/seal/utils/log"
)

const (
	summaryStatusReady = "Ready"
)

type DriftDetectTask struct {
	mu sync.Mutex

	tlsCertified bool
	logger       log.Logger
	modelClient  model.ClientSet
	kubeConfig   *rest.Config
	deployer     deptypes.Deployer
}

func NewServiceDriftDetectTask(
	mc model.ClientSet,
	kc *rest.Config,
	tlsCertified bool,
) (*DriftDetectTask, error) {
	in := &DriftDetectTask{
		modelClient:  mc,
		kubeConfig:   kc,
		tlsCertified: tlsCertified,
	}

	in.logger = log.WithName("task").WithName(in.Name())

	return in, nil
}

func (in *DriftDetectTask) Name() string {
	return "service-drift-detect"
}

func (in *DriftDetectTask) Process(ctx context.Context, args ...interface{}) error {
	if !in.mu.TryLock() {
		in.logger.Warn("previous processing is not finished")
		return nil
	}

	startTs := time.Now()

	in.logger.Info("start processing")

	defer func() {
		in.mu.Unlock()
		in.logger.Debugf("processed in %v", time.Since(startTs))
	}()

	return in.DriftDetect(ctx)
}

func (in *DriftDetectTask) DriftDetect(ctx context.Context) error {
	if !settings.EnableDriftDetection.ShouldValueBool(ctx, in.modelClient) {
		// Disable drift detection.
		return nil
	}

	entities, err := in.modelClient.Services().Query().
		Where(
			func(s *sql.Selector) {
				s.Where(sqljson.ValueEQ(
					service.FieldStatus,
					summaryStatusReady,
					sqljson.Path("summaryStatus"),
				))
			},
			service.Or(
				service.DriftResultIsNil(),
				func(s *sql.Selector) {
					s.Where(sqljson.ValueLTE(
						service.FieldDriftResult,
						time.Now().Add(-time.Hour),
						sqljson.Path("time"),
					))
				},
			),
		).
		All(ctx)
	if err != nil {
		return err
	}

	dp, err := in.getDeployer(ctx)
	if err != nil {
		return err
	}

	for _, entity := range entities {
		if err = in.updateServiceStatus(ctx, entity, status.ServiceStatusDetected); err != nil {
			return err
		}

		opts := pkgservice.Options{
			TlsCertified: in.tlsCertified,
		}
		if err = pkgservice.Detect(ctx, in.modelClient, dp, entity, opts); err != nil {
			return err
		}
	}

	return nil
}

func (in *DriftDetectTask) updateServiceStatus(
	ctx context.Context,
	entity *model.Service,
	s status.ConditionType,
) error {
	s.Reset(entity, "Detect service drift")

	return pkgservice.UpdateStatus(ctx, in.modelClient, entity)
}

func (in *DriftDetectTask) getDeployer(ctx context.Context) (deptypes.Deployer, error) {
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
