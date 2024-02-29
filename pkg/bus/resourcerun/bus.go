package resourcerun

import (
	"context"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/utils/bus"
)

// BusMessage wraps the changed model.ResourceRun as a bus.Message.
type BusMessage struct {
	// TransactionalModelClient holds the model.ClientSet of this calling session,
	// it should be a transactional DAO client,
	// please don't keep for long-term using.
	TransactionalModelClient model.ClientSet
	// Refer holds the updating model.ResourceRun item of this calling session.
	Refer *model.ResourceRun
}

// Notify notifies the changed model.ResourceRun.
func Notify(ctx context.Context, mc model.ClientSet, refer *model.ResourceRun) error {
	return bus.Publish(ctx, BusMessage{TransactionalModelClient: mc, Refer: refer})
}

// AddSubscriber add the subscriber to handle the changed notification from model.ResourceRun.
func AddSubscriber(n string, h func(context.Context, BusMessage) error) error {
	return bus.Subscribe(n, h)
}
