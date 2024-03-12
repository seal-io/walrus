package job

import (
	"context"
	"fmt"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	deptypes "github.com/seal-io/walrus/pkg/deployer/types"
	runstatus "github.com/seal-io/walrus/pkg/resourceruns/status"
	resstatus "github.com/seal-io/walrus/pkg/resources/status"
	"github.com/seal-io/walrus/utils/log"
)

// PerformRunJob performs the run job by the given run.
// Depending on the run type and status, the deployer will perform different actions.
func PerformRunJob(ctx context.Context, mc model.ClientSet, dp deptypes.Deployer, run *model.ResourceRun) (err error) {
	logger := log.WithName("resource-run")
	if runstatus.IsStatusCanceled(run) {
		logger.Info("run job is canceled", "run:", run.ID)

		return nil
	}

	ready, err := runstatus.CheckStatus(ctx, mc, run)
	if err != nil {
		return err
	}

	// If the run is not ready, it will not be performed.
	// Scheduler will check and perform the run job later.
	if !ready {
		return nil
	}

	runJobType, err := GetRunJobType(run)
	if err != nil {
		return err
	}

	defer func() {
		_, rerr := runstatus.UpdateStatusWithErr(ctx, mc, run, err)
		if rerr != nil {
			logger.Error(rerr, "failed to update run status", "run:", run.ID)
		}
	}()

	switch runJobType {
	case types.RunTaskTypePlan:
		return dp.Plan(ctx, mc, run, deptypes.PlanOptions{})
	case types.RunTaskTypeApply:
		// Mark resource status as deploying.
		res, err := mc.Resources().Get(ctx, run.ResourceID)
		if err != nil {
			return err
		}

		status.ResourceStatusDeployed.Reset(res, "")

		if err = resstatus.UpdateStatus(ctx, mc, res); err != nil {
			return err
		}

		err = dp.Apply(ctx, mc, run, deptypes.ApplyOptions{})
		if err != nil {
			status.ResourceStatusDeployed.False(res, err.Error())
			if rerr := resstatus.UpdateStatus(ctx, mc, res); rerr != nil {
				logger.Errorf("failed to update resource status, resource: %s, error: %s", res.ID, rerr)
			}

			return err
		}

		return nil
	case types.RunTaskTypeDestroy:
		res, err := mc.Resources().Get(ctx, run.ResourceID)
		if err != nil {
			return err
		}

		// Mark resource status as destroying.（stop or delete）.
		switch types.RunType(run.Type) {
		case types.RunTypeStop:
			status.ResourceStatusStopped.Reset(res, "")
		case types.RunTypeDelete:
			status.ResourceStatusDeleted.Reset(res, "")
		default:
			return fmt.Errorf("unsupported run type %s for destroy", run.Type)
		}

		if err = resstatus.UpdateStatus(ctx, mc, res); err != nil {
			return err
		}

		err = dp.Destroy(ctx, mc, run, deptypes.DestroyOptions{})
		if err != nil {
			switch types.RunType(run.Type) {
			case types.RunTypeStop:
				status.ResourceStatusStopped.False(res, err.Error())
			case types.RunTypeDelete:
				status.ResourceStatusDeleted.False(res, err.Error())
			}

			if rerr := resstatus.UpdateStatus(ctx, mc, res); rerr != nil {
				logger.Errorf("failed to update resource status, resource: %s, error: %s", res.ID, rerr)
			}

			return err
		}

		return nil
	}

	return fmt.Errorf("unknown run job type %s", runJobType)
}

// GetRunJobType gets the run job type for deployer with its type and status.
// It makes the following decision.
//
//	| Run type         | Run status       | Job type         |
//	| ---------------- | ---------------- | ---------------- |
//	| create           | pending          | plan             |
//	| create           | planed           | apply            |
//	| upgrade          | pending          | plan             |
//	| upgrade          | planed           | apply            |
//	| delete           | pending          | plan             |
//	| delete           | planed           | destroy          |
//	| start            | pending          | plan             |
//	| start            | planed           | apply            |
//	| stop             | pending          | plan             |
//	| stop             | planed           | destroy          |
//	| rollback         | pending          | plan             |
//	| rollback         | planed           | apply            |
func GetRunJobType(run *model.ResourceRun) (types.RunJobType, error) {
	if runstatus.IsStatusPending(run) {
		return types.RunTaskTypePlan, nil
	}

	if runstatus.IsStatusPlanned(run) {
		switch types.RunType(run.Type) {
		case types.RunTypeCreate, types.RunTypeUpdate, types.RunTypeStart, types.RunTypeRollback:
			return types.RunTaskTypeApply, nil
		case types.RunTypeDelete, types.RunTypeStop:
			return types.RunTaskTypeDestroy, nil
		}
	}

	return "", fmt.Errorf("unknown run type %s and status %s", run.Type, run.Status.SummaryStatus)
}
