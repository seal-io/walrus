package connector

import (
	"context"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/utils/bus"
)

// BusMessage wraps the changed model.Connector as a bus.Message.
type BusMessage struct {
	ModelClient model.ClientSet
	Refer       *model.Connector
	Replace     bool
}

// Notify notifies the changed connector.
func Notify(ctx context.Context, mc model.ClientSet, refer *model.Connector, replace bool) error {
	return bus.Publish(ctx, BusMessage{ModelClient: mc, Refer: refer, Replace: replace})
}

// AddSubscriber add the subscriber to handle the changed notification from connector.
func AddSubscriber(n string, h func(context.Context, BusMessage) error) error {
	return bus.Subscribe(n, h)
}
