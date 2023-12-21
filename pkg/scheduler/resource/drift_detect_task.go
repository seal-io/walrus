package resource

import (
	"context"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/deployer"
	deployertf "github.com/seal-io/walrus/pkg/deployer/terraform"
	deptypes "github.com/seal-io/walrus/pkg/deployer/types"
	pkgresource "github.com/seal-io/walrus/pkg/resource"
	"github.com/seal-io/walrus/pkg/settings"
	"github.com/seal-io/walrus/utils/log"
)

const (
	summaryStatusReady  = "Ready"
	driftDetectionLimit = 10
)

type DriftDetectionTask struct {
	logger log.Logger

	modelClient model.ClientSet
	clientset   *kubernetes.Clientset
	deployer    deptypes.Deployer
}

func NewResourceDriftDetectTask(
	logger log.Logger,
	mc model.ClientSet,
	kc *rest.Config,
) (*DriftDetectionTask, error) {
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

	cs, err := kubernetes.NewForConfig(kc)
	if err != nil {
		return nil, err
	}

	return &DriftDetectionTask{
		logger:      logger,
		modelClient: mc,
		clientset:   cs,
		deployer:    dp,
	}, nil
}

func (in *DriftDetectionTask) Process(ctx context.Context, args ...any) error {
	if !settings.EnableDriftDetection.ShouldValueBool(ctx, in.modelClient) {
		// Disable drift detection.
		return nil
	}

	entities, err := in.modelClient.Resources().Query().
		Where(
			func(s *sql.Selector) {
				s.Where(sqljson.ValueEQ(
					resource.FieldStatus,
					summaryStatusReady,
					sqljson.Path("summaryStatus"),
				))
			},
			resource.Or(
				resource.DriftDetectionIsNil(),
				func(s *sql.Selector) {
					s.Where(sqljson.ValueLTE(
						resource.FieldDriftDetection,
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

	for i := range entities {
		entity := entities[i]

		status.ResourceStatusDetected.Unknown(entity, "")
		entity.Status.SetSummary(status.WalkResource(&entity.Status))

		err = in.modelClient.Resources().UpdateOne(entity).
			SetStatus(entity.Status).
			Exec(ctx)
		if err != nil {
			return err
		}

		waitErr := wait.PollUntilContextTimeout(ctx, 5*time.Second, 1*time.Hour, true,
			func(ctx context.Context) (bool, error) {
				jobs, err := in.clientset.BatchV1().Jobs(types.WalrusSystemNamespace).List(ctx,
					metav1.ListOptions{LabelSelector: types.LabelWalrusDriftDetection + "=true"})
				if err != nil {
					return false, err
				}

				runningJobs := 0

				for i := range jobs.Items {
					j := jobs.Items[i]
					if j.Status.Active > 0 {
						runningJobs++
					}
				}

				if runningJobs >= driftDetectionLimit {
					return false, nil
				}

				in.logger.Debugf("drift detection jobs %d less than limit %d", runningJobs, driftDetectionLimit)

				return true, nil
			})
		if waitErr != nil {
			in.logger.Warnf("wait for detection job timeout: %v", waitErr)
			return nil
		}

		err = pkgresource.Detect(ctx, in.modelClient, entities[i], pkgresource.Options{
			Deployer: in.deployer,
		})
		if err != nil {
			status.ResourceStatusDetected.False(entity, err.Error())
			entity.Status.SetSummary(status.WalkResource(&entity.Status))

			updateErr := in.modelClient.Resources().UpdateOne(entity).
				SetStatus(entity.Status).
				Exec(ctx)
			if updateErr != nil {
				in.logger.Errorf("update resource %q status failed: %v", entity.ID, updateErr)
			}

			return err
		}

		if err != nil {
			return err
		}
	}

	return nil
}
