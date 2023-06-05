package template

import (
	"context"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/utils/bus"
)

// BusMessage wraps the changed model.Template as a bus.Message.
type BusMessage struct {
	// Refer holds the updating model.Template item of this calling session.
	Refer *model.Template
}

// Notify notifies the changed model.Template.
func Notify(ctx context.Context, refer *model.Template) error {
	return bus.Publish(ctx, BusMessage{Refer: refer})
}

// AddSubscriber add the subscriber to handle the changed notification from model.Template.
func AddSubscriber(n string, h func(context.Context, BusMessage) error) error {
	return bus.Subscribe(n, h)
}
