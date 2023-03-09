package applicationrevision

import (
	"context"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/applicationinstance"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/utils/bus"
)

type BusMessage struct {
	ModelClient model.ClientSet
	Refer       *model.ApplicationRevision
}

func Notify(ctx context.Context, mc model.ClientSet, refer *model.ApplicationRevision) error {
	return bus.Publish(ctx, BusMessage{ModelClient: mc, Refer: refer})
}

func AddSubscriber(n string, h func(context.Context, BusMessage) error) error {
	return bus.Subscribe(n, h)
}

func OnRevisionUpdate(ctx context.Context, message BusMessage) error {
	return syncApplicationRevisionStatus(ctx, message.ModelClient, message.Refer)
}

func syncApplicationRevisionStatus(ctx context.Context, mc model.ClientSet, revision *model.ApplicationRevision) error {
	// report to application instance.
	appInstance, err := mc.ApplicationInstances().Query().
		Where(applicationinstance.ID(revision.InstanceID)).
		Select(
			applicationinstance.FieldID,
			applicationinstance.FieldStatus).
		Only(ctx)
	if err != nil {
		return err
	}
	switch revision.Status {
	case status.ApplicationRevisionStatusSucceeded:
		if appInstance.Status == status.ApplicationInstanceStatusDeleting {
			// delete application instance.
			err = mc.ApplicationInstances().DeleteOne(appInstance).
				Exec(ctx)
		} else {
			err = mc.ApplicationInstances().UpdateOne(appInstance).
				SetStatus(status.ApplicationInstanceStatusDeployed).
				Exec(ctx)
		}
	case status.ApplicationRevisionStatusFailed:
		if appInstance.Status == status.ApplicationInstanceStatusDeleting {
			appInstance.Status = status.ApplicationInstanceStatusDeleteFailed
		} else {
			appInstance.Status = status.ApplicationInstanceStatusDeployFailed
		}
		err = mc.ApplicationInstances().UpdateOne(appInstance).
			SetStatus(appInstance.Status).
			SetStatusMessage(revision.StatusMessage).
			Exec(ctx)
	}

	return err
}
