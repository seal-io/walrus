package connector

import (
	"context"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/utils/bus"
)

// BusMessage wraps the changed model.Connector as a bus.Message.
type BusMessage struct {
	// Refer holds the updating model.Connector item of this calling session.
	Refer *model.Connector
	// ReinstallTools indicates to reinstall the cost tools.
	ReinstallTools bool
}

// Notify notifies the changed model.Connector.
func Notify(ctx context.Context, refer *model.Connector, reinstallTools bool) error {
	return bus.Publish(ctx, BusMessage{Refer: refer, ReinstallTools: reinstallTools})
}

// AddSubscriber add the subscriber to handle the changed notification from model.Connector.
func AddSubscriber(n string, h func(context.Context, BusMessage) error) error {
	return bus.Subscribe(n, h)
}
