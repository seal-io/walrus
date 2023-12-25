package templateversion

import (
	"context"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/utils/bus"
)

// BusMessage wraps the changed model.TemplateVersion as a bus.Message.
type BusMessage struct {
	// Refer holds the updating model.TemplateVersion item of this calling session.
	Refer *model.TemplateVersion
}

// Notify notifies the changed model.TemplateVersion.
func Notify(ctx context.Context, refer *model.TemplateVersion) error {
	return bus.Publish(ctx, BusMessage{Refer: refer})
}

// AddSubscriber add the subscriber to handle the changed notification from model.TemplateVersion.
func AddSubscriber(n string, h func(context.Context, BusMessage) error) error {
	return bus.Subscribe(n, h)
}
