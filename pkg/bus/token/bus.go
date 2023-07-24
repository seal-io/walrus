package token

import (
	"context"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/utils/bus"
)

// BusMessage wraps the deleted model.Token as a bus.Message.
type BusMessage struct {
	// Refers holds the deleted model.Token list of this calling session.
	Refers model.Tokens
}

// Notify notifies the deleted model.Token.
func Notify(ctx context.Context, refers model.Tokens) error {
	return bus.Publish(ctx, BusMessage{Refers: refers})
}

// AddSubscriber add the subscriber to handle the deletion notification from model.Token.
func AddSubscriber(n string, h func(context.Context, BusMessage) error) error {
	return bus.Subscribe(n, h)
}
