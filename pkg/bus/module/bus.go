package module

import (
	"context"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/utils/bus"
)

// BusMessage wraps the changed model.Module as a bus.Message.
type BusMessage struct {
	// Refer holds the updating model.Module item of this calling session.
	Refer *model.Module
}

// Notify notifies the changed model.Module.
func Notify(ctx context.Context, refer *model.Module) error {
	return bus.Publish(ctx, BusMessage{Refer: refer})
}

// AddSubscriber add the subscriber to handle the changed notification from model.Module.
func AddSubscriber(n string, h func(context.Context, BusMessage) error) error {
	return bus.Subscribe(n, h)
}
