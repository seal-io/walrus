package settings

import (
	"context"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/utils/bus"
)

// BusMessage wraps the changed model.Setting as a bus.Message.
type BusMessage struct {
	ModelClient model.ClientSet
	Refer       model.Settings
}

// Notify notifies the changed model.Setting.
func Notify(ctx context.Context, mc model.ClientSet, refer model.Settings) error {
	return bus.Publish(ctx, BusMessage{ModelClient: mc, Refer: refer})
}

// AddSubscriber add the subscriber to handle the changed notification from proxy model.Setting.
func AddSubscriber(n string, h func(context.Context, BusMessage) error) error {
	return bus.Subscribe(n, h)
}
