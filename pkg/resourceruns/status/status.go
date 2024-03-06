package status

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
	"k8s.io/apimachinery/pkg/util/sets"

	runbus "github.com/seal-io/walrus/pkg/bus/resourcerun"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/model/resourcerelationship"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	resstatus "github.com/seal-io/walrus/pkg/resources/status"
	"github.com/seal-io/walrus/utils/strs"
)

const (
	ActionDelete = "delete"
	ActionStop   = "stop"

	summaryStatusStopping = "Stopping"
)

// IsStatusRunning checks if the resource run is in the running status.
func IsStatusRunning(run *model.ResourceRun) bool {
	return status.ResourceRunStatusPending.IsUnknown(run) ||
		status.ResourceRunStatusPlanned.IsUnknown(run) ||
		status.ResourceRunStatusApplied.IsUnknown(run)
}

func IsStatusPending(run *model.ResourceRun) bool {
	return status.ResourceRunStatusPending.IsUnknown(run)
}

// IsStatusApplying checks if the resource run is in the applying status.
func IsStatusApplying(run *model.ResourceRun) bool {
	return status.ResourceRunStatusApplied.IsUnknown(run)
}

func IsStatusFailed(run *model.ResourceRun) bool {
	return status.ResourceRunStatusPending.IsFalse(run) ||
		status.ResourceRunStatusPlanned.IsFalse(run) ||
		status.ResourceRunStatusApplied.IsFalse(run)
}

func IsStatusPlanned(run *model.ResourceRun) bool {
	return status.ResourceRunStatusPlanned.IsTrue(run) &&
		!status.ResourceRunStatusApplied.Exist(run)
}

// IsStatusPlanCondition checks if the resource run is in the plan condition.
func IsStatusPlanCondition(run *model.ResourceRun) bool {
	return status.ResourceRunStatusPlanned.Exist(run) &&
		!status.ResourceRunStatusApplied.Exist(run)
}

func IsStatusSucceeded(run *model.ResourceRun) bool {
	return status.ResourceRunStatusApplied.IsTrue(run)
}

// SetStatusFalse sets the status of the resource run to false.
func SetStatusFalse(run *model.ResourceRun, errMsg string) {
	switch {
	case status.ResourceRunStatusPlanned.IsUnknown(run):
		errMsg = fmt.Sprintf("plan failed: %s", errMsg)
		status.ResourceRunStatusPlanned.False(run, errMsg)
	case status.ResourceRunStatusApplied.IsUnknown(run):
		errMsg = fmt.Sprintf("apply failed: %s", errMsg)
		status.ResourceRunStatusApplied.False(run, errMsg)
	case status.ResourceRunStatusPending.IsUnknown(run):
		errMsg = fmt.Sprintf("pending failed: %s", errMsg)
		status.ResourceRunStatusPending.False(run, errMsg)
	}

	run.Status.SetSummary(status.WalkResourceRun(&run.Status))
}

// SetStatusTrue sets the status of the resource run to true.
// It marks the status of the resource run as "Succeeded".
func SetStatusTrue(run *model.ResourceRun, msg string) {
	switch {
	case status.ResourceRunStatusPlanned.IsUnknown(run):
		status.ResourceRunStatusPlanned.True(run, msg)
	case status.ResourceRunStatusApplied.IsUnknown(run):
		status.ResourceRunStatusApplied.True(run, msg)
	}

	run.Status.SetSummary(status.WalkResourceRun(&run.Status))
}

// UpdateStatus updates the status of the resource run.
func UpdateStatus(ctx context.Context, mc model.ClientSet, run *model.ResourceRun) (*model.ResourceRun, error) {
	if run == nil {
		return nil, nil
	}

	// Report to resource run.
	run.Status.SetSummary(status.WalkResourceRun(&run.Status))

	run, err := mc.ResourceRuns().UpdateOne(run).
		SetStatus(run.Status).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	err = runbus.Notify(ctx, mc, run)
	if err != nil {
		return nil, err
	}

	return run, nil
}

// UpdateStatusWithErr updates the status of the resource run with the given error.
func UpdateStatusWithErr(ctx context.Context, mc model.ClientSet, run *model.ResourceRun, err error) (*model.ResourceRun, error) {
	if err == nil {
		return run, nil
	}

	SetStatusFalse(run, err.Error())

	return UpdateStatus(ctx, mc, run)
}

// CheckDependencyStatus checks the resource dependency status of the resource run.
// Works for create, update, rollback and start actions.
func CheckDependencyStatus(
	ctx context.Context,
	mc model.ClientSet,
	run *model.ResourceRun,
) (bool, error) {
	// Check dependencies.
	dependencies, err := mc.ResourceRelationships().Query().
		Where(
			resourcerelationship.ResourceID(run.ResourceID),
			resourcerelationship.DependencyIDNEQ(run.ResourceID),
		).
		QueryDependency().
		All(ctx)
	if err != nil {
		return false, err
	}

	if len(dependencies) == 0 {
		return true, nil
	}

	notReadyDependencies := make([]*model.Resource, 0, len(dependencies))

	for _, d := range dependencies {
		if resstatus.IsStatusError(d) {
			msg := fmt.Sprintf("failed as the dependency resource %s has an error status", d.Name)
			SetStatusFalse(run, msg)
			if _, err := UpdateStatus(ctx, mc, run); err != nil {
				return false, err
			}

			return false, nil
		}

		if resstatus.IsInactive(d) {
			msg := fmt.Sprintf("failed as the dependency resource %s is inactive", d.Name)
			SetStatusFalse(run, msg)
			if _, err := UpdateStatus(ctx, mc, run); err != nil {
				return false, err
			}
		}

		if !resstatus.IsStatusReady(d) {
			notReadyDependencies = append(notReadyDependencies, d)
		}
	}

	if len(notReadyDependencies) > 0 {
		// Get the names of the not ready dependencies.
		dependencyNames := sets.NewString()
		for _, d := range notReadyDependencies {
			dependencyNames.Insert(d.Name)
		}

		msg := fmt.Sprintf("Waiting for dependent resources to be ready: %s", strs.Join(", ", dependencyNames.List()...))

		status.ResourceRunStatusPending.Reset(run, msg)

		if _, err := UpdateStatus(ctx, mc, run); err != nil {
			return false, err
		}

		return false, nil
	}

	return true, nil
}

// CheckDependantStatus checks the resource dependant status of the resource run.
// Works for both stop and delete actions.
func CheckDependantStatus(
	ctx context.Context,
	mc model.ClientSet,
	run *model.ResourceRun,
) (bool, error) {
	// Check dependants.
	query := mc.ResourceRelationships().Query().
		Where(
			resourcerelationship.ResourceIDNEQ(run.ResourceID),
			resourcerelationship.DependencyID(run.ResourceID),
		).
		QueryDependency()

	// When stop resource, stopped dependant resource should be exclude.
	if run.Type == types.RunTypeStop.String() {
		query.Where(func(s *sql.Selector) {
			s.Where(sqljson.ValueNEQ(
				resource.FieldStatus,
				summaryStatusStopping,
				sqljson.Path("summaryStatus"),
			))
		})
	}

	// Check dependants.
	dependants, err := query.All(ctx)
	if err != nil {
		return false, err
	}

	if len(dependants) == 0 {
		return true, nil
	}

	dependantsNames := sets.NewString()
	for _, d := range dependants {
		if resstatus.IsStatusError(d) {
			msg := fmt.Sprintf("failed as the dependant resource %s has an error status", d.Name)
			SetStatusFalse(run, msg)

			if _, err = UpdateStatus(ctx, mc, run); err != nil {
				return false, fmt.Errorf("failed to update resource run status: %w", err)
			}

			return false, nil
		}

		dependantsNames.Insert(d.Name)
	}

	var msg string
	switch run.Type {
	case types.RunTypeDelete.String():
		msg = fmt.Sprintf("Waiting for dependants to be deleted: %s", strs.Join(", ", dependantsNames.List()...))

	case types.RunTypeStop.String():
		msg = fmt.Sprintf("Waiting for dependants to be stopped: %s", strs.Join(", ", dependantsNames.List()...))
	default:
		return false, fmt.Errorf("unsupported action type: %s", run.Type)
	}

	status.ResourceRunStatusPending.Reset(run, msg)

	if _, err = UpdateStatus(ctx, mc, run); err != nil {
		return false, fmt.Errorf("failed to update resource run status: %w", err)
	}

	return false, nil
}

// CheckStatus checks the status of the resource run dependant or dependency with the given resource run.
func CheckStatus(ctx context.Context, mc model.ClientSet, run *model.ResourceRun) (bool, error) {
	switch types.RunType(run.Type) {
	case types.RunTypeCreate, types.RunTypeUpdate, types.RunTypeRollback, types.RunTypeStart:
		return CheckDependencyStatus(ctx, mc, run)
	case types.RunTypeDelete, types.RunTypeStop:
		return CheckDependantStatus(ctx, mc, run)
	default:
		return false, fmt.Errorf("unsupported action type: %s", run.Type)
	}
}
