package job

import (
	"context"
	"fmt"

	"k8s.io/client-go/rest"

	runbus "github.com/seal-io/walrus/pkg/bus/resourcerun"
	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/deployer"
	deptypes "github.com/seal-io/walrus/pkg/deployer/types"
	runstatus "github.com/seal-io/walrus/pkg/resourceruns/status"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/strs"
)

func Syncer(kc *rest.Config) syncer {
	return syncer{
		logger: log.WithName("resource-run").WithName("syncer"),
		kc:     kc,
	}
}

type syncer struct {
	logger log.Logger
	kc     *rest.Config
}

// Do handler update of the resource run.
func (s syncer) Do(ctx context.Context, bm runbus.BusMessage) (err error) {
	var (
		mc  = bm.TransactionalModelClient
		run = bm.Refer
	)

	// Report to resource.
	entity, err := mc.Resources().Query().
		Where(resource.ID(run.ResourceID)).
		Only(ctx)
	if err != nil {
		return err
	}

	dp, err := deployer.Get(ctx, deptypes.CreateOptions{
		Type:       types.DeployerTypeTF,
		KubeConfig: s.kc,
	})
	if err != nil {
		return err
	}

	switch {
	case runstatus.IsStatusPlanned(run):
		if !run.Preview {
			err = PerformRunJob(ctx, mc, dp, run)
			if err != nil {
				return err
			}

			return nil
		}

		s.logger.Debugf("resource %s %q run %s %q, wait to be applied", entity.Name, run.Type, run.ID, run.Status.SummaryStatus)

		return nil
	case runstatus.IsStatusSucceeded(run):
		s.logger.Debugf("resource %s %q run %s status %q", entity.Name, run.Type, run.ID, run.Status.SummaryStatus)

		switch {
		case status.ResourceStatusDeleted.IsUnknown(entity):
			err = mc.Resources().DeleteOne(entity).
				Exec(ctx)
			if err == nil {
				s.logger.Debugf("resource %s is deleted", entity.Name, run.ID, run.Status.SummaryStatus)
				return nil
			}

			msg := err.Error()
			// Check dependants.
			dependants, rerr := dao.GetResourceDependantNames(ctx, mc, entity)
			if rerr != nil {
				s.logger.Errorf("failed to get dependants of resource %s: %v", entity.Name, rerr)
			}

			if len(dependants) > 0 {
				msg = fmt.Sprintf("resource to be deleted is the dependency of: %s", strs.Join(", ", dependants...))
			}

			// Mark resource delete failed.
			status.ResourceStatusDeleted.False(entity, msg)

		case status.ResourceStatusStopped.IsUnknown(entity):
			// Stopping -> Stopped.
			status.ResourceStatusStopped.True(entity, "")
		default:
			// Deployed.
			status.ResourceStatusDeployed.True(entity, "")
			status.ResourceStatusReady.Unknown(entity, "")
		}
	case runstatus.IsStatusFailed(run):
		// If job fail, and preview is true, then we should not update the resource status.
		if runstatus.IsPreviewPlanFailed(run) {
			s.logger.Debugf("resource %s %q run %s plan Failed", entity.Name, run.Type, run.ID)
			return nil
		}

		s.logger.Debugf("resource %s %q run %s status %q", entity.Name, run.Type, run.ID, run.Status.SummaryStatus)

		switch run.Type {
		case types.RunTypeCreate.String(), types.RunTypeUpdate.String(), types.RunTypeRollback.String(), types.RunTypeStart.String():
			status.ResourceStatusDeployed.False(entity, "")
		case types.RunTypeDelete.String():
			status.ResourceStatusDeleted.False(entity, "")
		case types.RunTypeStop.String():
			status.ResourceStatusStopped.False(entity, "")
		default:
			s.logger.Debugf("run %s unsupported action type %q", run.ID, run.Type)
			return nil
		}

		entity.Status.SummaryStatusMessage = run.Status.SummaryStatusMessage
	default:
		s.logger.Debugf("skip resource run %s status %q", run.ID, run.Status.SummaryStatus)
		return nil
	}

	entity.Status.SetSummary(status.WalkResource(&entity.Status))

	s.logger.Debugf("set resource %s status to %q", entity.Name, entity.Status.SummaryStatus)

	return mc.Resources().UpdateOne(entity).
		SetStatus(entity.Status).
		Exec(ctx)
}
