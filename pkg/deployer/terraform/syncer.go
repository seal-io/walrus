package terraform

import (
	"context"

	revisionbus "github.com/seal-io/walrus/pkg/bus/resourcerevision"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/types/status"
)

// SyncResourceRevisionStatus updates the status of the service according to its recent finished service revision.
func SyncResourceRevisionStatus(ctx context.Context, bm revisionbus.BusMessage) (err error) {
	var (
		mc       = bm.TransactionalModelClient
		revision = bm.Refer
	)

	// Report to resource.
	entity, err := mc.Resources().Query().
		Where(resource.ID(revision.ResourceID)).
		Select(
			resource.FieldID,
			resource.FieldStatus).
		Only(ctx)
	if err != nil {
		return err
	}

	if status.ResourceRevisionStatusReady.IsTrue(revision) {
		if status.ResourceStatusDeleted.IsUnknown(entity) {
			return mc.Resources().DeleteOne(entity).
				Exec(ctx)
		}

		status.ResourceStatusDeployed.True(entity, "")
		status.ResourceStatusReady.Unknown(entity, "")
	} else if status.ResourceRevisionStatusReady.IsFalse(revision) {
		if status.ResourceStatusDeleted.IsUnknown(entity) {
			status.ResourceStatusDeleted.False(entity, "")
		} else {
			status.ResourceStatusDeployed.False(entity, "")
		}

		entity.Status.SummaryStatusMessage = revision.Status.SummaryStatusMessage
	}

	entity.Status.SetSummary(status.WalkResource(&entity.Status))

	return mc.Resources().UpdateOne(entity).
		SetStatus(entity.Status).
		Exec(ctx)
}
