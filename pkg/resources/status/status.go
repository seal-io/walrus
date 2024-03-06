package status

import (
	"context"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types/status"
)

// IsStatusReady returns true if the resource is ready.
func IsStatusReady(entity *model.Resource) bool {
	switch entity.Status.SummaryStatus {
	case "Preparing", "NotReady", "Ready":
		return true
	}

	return false
}

// IsStatusError returns true if the resource is in error status.
func IsStatusError(entity *model.Resource) bool {
	switch entity.Status.SummaryStatus {
	case "DeployFailed", "DeleteFailed":
		return true
	case "Progressing":
		return entity.Status.Error
	}

	return false
}

// IsStatusDeleted returns true if the resource is deleted or to be deleted.
func IsStatusDeleted(entity *model.Resource) bool {
	switch entity.Status.SummaryStatus {
	case "Deleted", "Deleting":
		return true
	}

	// If the resource is in progressing status to be deleted.
	if status.ResourceStatusProgressing.IsUnknown(entity) &&
		status.ResourceStatusDeleted.IsUnknown(entity) {
		return true
	}

	return false
}

// IsStatusStopped returns true if the resource is stopped or to be stopped.
func IsStatusStopped(entity *model.Resource) bool {
	switch entity.Status.SummaryStatus {
	case "Stopped", "Stopping":
		return true
	}

	// If the resource is in progressing status to be stopped.
	if status.ResourceStatusProgressing.IsUnknown(entity) &&
		status.ResourceStatusStopped.IsUnknown(entity) {
		return true
	}

	return false
}

// IsStatusDeployed returns true if the resource is deployed or to be deployed.
func IsStatusDeployed(entity *model.Resource) bool {
	switch entity.Status.SummaryStatus {
	case "Deployed", "Deploying":
		return true
	}

	// If the resource is in progressing status to be deployed.
	if status.ResourceStatusProgressing.IsUnknown(entity) &&
		status.ResourceStatusDeployed.IsUnknown(entity) {
		return true
	}

	return false
}

// IsInactive tells whether the given resource is inactive.
func IsInactive(r *model.Resource) bool {
	if r == nil {
		return false
	}

	return r.Status.SummaryStatus == status.ResourceStatusUnDeployed.String() ||
		r.Status.SummaryStatus == status.ResourceStatusStopped.String()
}

// UpdateStatus updates the status of the given resource.
func UpdateStatus(
	ctx context.Context,
	mc model.ClientSet,
	entity *model.Resource,
) error {
	entity.Status.SetSummary(status.WalkResource(&entity.Status))

	err := mc.Resources().UpdateOne(entity).
		SetStatus(entity.Status).
		Exec(ctx)
	if err != nil && !model.IsNotFound(err) {
		return err
	}

	return nil
}
