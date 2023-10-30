package resourcerevision

import (
	"context"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/utils/bus"
)

// BusMessage wraps the changed model.ResourceRevision as a bus.Message.
type BusMessage struct {
	// TransactionalModelClient holds the model.ClientSet of this calling session,
	// it should be a transactional DAO client,
	// please don't keep for long-term using.
	TransactionalModelClient model.ClientSet
	// Refer holds the updating model.ResourceRevision item of this calling session.
	Refer *model.ResourceRevision
}

// Notify notifies the changed model.ResourceRevision.
func Notify(ctx context.Context, mc model.ClientSet, refer *model.ResourceRevision) error {
	return bus.Publish(ctx, BusMessage{TransactionalModelClient: mc, Refer: refer})
}

// AddSubscriber add the subscriber to handle the changed notification from model.ResourceRevision.
func AddSubscriber(n string, h func(context.Context, BusMessage) error) error {
	return bus.Subscribe(n, h)
}
