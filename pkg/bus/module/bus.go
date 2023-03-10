package module

import (
	"context"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/utils/bus"
)

// BusMessage wraps the changed model.Module as a bus.Message.
type BusMessage struct {
	ModelClient model.ClientSet
	Refer       *model.Module
}

// Notify notifies the changed model.Module.
func Notify(ctx context.Context, mc model.ClientSet, refer *model.Module) error {
	return bus.Publish(ctx, BusMessage{ModelClient: mc, Refer: refer})
}

// AddSubscriber add the subscriber to handle the changed notification from model.Module.
func AddSubscriber(n string, h func(context.Context, BusMessage) error) error {
	return bus.Subscribe(n, h)
}
