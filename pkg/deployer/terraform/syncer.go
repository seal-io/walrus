package terraform

import (
	"context"
	"fmt"

	revisionbus "github.com/seal-io/walrus/pkg/bus/resourcerevision"
	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/strs"
)

// SyncResourceRevisionStatus updates the status of the service according to its recent finished service revision.
func SyncResourceRevisionStatus(ctx context.Context, bm revisionbus.BusMessage) (err error) {
	var (
		logger = log.WithName("deployer").WithName("tf")

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
		switch {
		case status.ResourceStatusDeleted.IsUnknown(entity):
			err = mc.Resources().DeleteOne(entity).
				Exec(ctx)
			if err == nil {
				return nil
			}

			msg := err.Error()
			// Check dependants.
			dependants, rerr := dao.GetResourceDependantNames(ctx, mc, entity)
			if rerr != nil {
				logger.Errorf("failed to get dependants of resource %s: %v", entity.Name, rerr)
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
	} else if status.ResourceRevisionStatusReady.IsFalse(revision) {
		switch {
		case status.ResourceStatusDeleted.IsUnknown(entity):
			status.ResourceStatusDeleted.False(entity, "")
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
