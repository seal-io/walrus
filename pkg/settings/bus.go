package settings

import (
	"context"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/utils/bus"
)

// BusMessage wraps the changed model.Setting as a bus.Message.
type BusMessage struct {
	// ModelClient holds the model.ClientSet of this calling session,
	// it should be a transactional DAO client,
	// please don't keep for long-term using.
	ModelClient model.ClientSet
	// Refer holds the updating model.Setting list of this calling session.
	Refer model.Settings
}

// Notify notifies the changed model.Setting.
func Notify(ctx context.Context, mc model.ClientSet, refer model.Settings) error {
	return bus.Publish(ctx, BusMessage{ModelClient: mc, Refer: refer})
}

// AddSubscriber add the subscriber to handle the changed notification from proxy model.Setting.
func AddSubscriber(n string, h func(context.Context, BusMessage) error) error {
	return bus.Subscribe(n, h)
}
