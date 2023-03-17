package applicationrevision

import (
	"context"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/utils/bus"
)

// BusMessage wraps the changed model.ApplicationRevision as a bus.Message.
type BusMessage struct {
	// TransactionalModelClient holds the model.ClientSet of this calling session,
	// it should be a transactional DAO client,
	// please don't keep for long-term using.
	TransactionalModelClient model.ClientSet
	// Refer holds the updating model.ApplicationRevision item of this calling session.
	Refer *model.ApplicationRevision
}

// Notify notifies the changed model.ApplicationRevision.
func Notify(ctx context.Context, mc model.ClientSet, refer *model.ApplicationRevision) error {
	return bus.Publish(ctx, BusMessage{TransactionalModelClient: mc, Refer: refer})
}

// AddSubscriber add the subscriber to handle the changed notification from model.ApplicationRevision.
func AddSubscriber(n string, h func(context.Context, BusMessage) error) error {
	return bus.Subscribe(n, h)
}
