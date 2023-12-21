package terraform

import (
	"context"

	revisionbus "github.com/seal-io/walrus/pkg/bus/resourcerevision"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	pkgresource "github.com/seal-io/walrus/pkg/resource"
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
		Only(ctx)
	if err != nil {
		return err
	}

	if status.ResourceRevisionStatusReady.IsTrue(revision) {
		switch {
		case status.ResourceStatusDeleted.IsUnknown(entity):
			return mc.Resources().DeleteOne(entity).
				Exec(ctx)
		case status.ResourceStatusStopped.IsUnknown(entity):
			// Stopping -> Stopped.
			status.ResourceStatusStopped.True(entity, "")
		default:
			switch {
			case status.ResourceStatusDetected.IsUnknown(entity):
				// Detecting -> Detected.
				status.ResourceStatusDetected.True(entity, "")

				err = pkgresource.UpdateServiceDriftResult(ctx, mc, entity, revision)
				if err != nil {
					return err
				}
			case status.ResourceStatusSynced.IsUnknown(entity):
				// Syncing -> Synced.
				status.ResourceStatusSynced.True(entity, "")
			}

			// Deployed.
			status.ResourceStatusDeployed.True(entity, "")
			status.ResourceStatusReady.Unknown(entity, "")
		}
	} else if status.ResourceRevisionStatusReady.IsFalse(revision) {
		switch {
		case status.ResourceStatusDeleted.IsUnknown(entity):
			status.ResourceStatusDeleted.False(entity, "")
		case status.ResourceStatusSynced.IsUnknown(entity):
			status.ResourceStatusSynced.False(entity, "")
		case status.ResourceStatusDetected.IsUnknown(entity):
			status.ResourceStatusDetected.False(entity, "")
		default:
			status.ResourceStatusDeployed.False(entity, "")
		}

		entity.Status.SummaryStatusMessage = revision.Status.SummaryStatusMessage
	}

	entity.Status.SetSummary(status.WalkResource(&entity.Status))

	return mc.Resources().UpdateOne(entity).
		SetStatus(entity.Status).
		Exec(ctx)
}
