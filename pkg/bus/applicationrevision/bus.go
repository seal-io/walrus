package applicationrevision

import (
	"context"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/utils/bus"
)

type BusMessage struct {
	ModelClient model.ClientSet
	Refer       *model.ApplicationRevision
}

func Notify(ctx context.Context, mc model.ClientSet, refer *model.ApplicationRevision) error {
	return bus.Publish(ctx, BusMessage{ModelClient: mc, Refer: refer})
}

func AddSubscriber(n string, h func(context.Context, BusMessage) error) error {
	return bus.Subscribe(n, h)
}
