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
	if status.ResourceRunStatusCanceled.IsTrue(run) {
		logger.Info("run job is canceled", "run:", run.ID)

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

		res.IsModified = false
		status.ResourceStatusDeployed.Reset(res, "")

		if err = resstatus.UpdateStatus(ctx, mc, res); err != nil {
			return err
		}

		return dp.Apply(ctx, mc, run, deptypes.ApplyOptions{})
	case types.RunTaskTypeDestroy:
		res, err := mc.Resources().Get(ctx, run.ResourceID)
		if err != nil {
			return err
		}
		res.IsModified = false
		// Mark resource status as destroying.（stop or delete）.
		switch types.RunType(run.Type) {
		case types.RunTypeStop:
			status.ResourceStatusStopped.Reset(res, "")
		case types.RunTypeDelete:
			status.ResourceStatusDeleted.Reset(res, "")
		}

		if err = resstatus.UpdateStatus(ctx, mc, res); err != nil {
			return err
		}

		return dp.Destroy(ctx, mc, run, deptypes.DestroyOptions{})
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
