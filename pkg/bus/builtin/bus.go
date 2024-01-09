package builtin

import (
	"context"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/utils/bus"
)

type BusMessage struct {
	TransactionalModelClient model.ClientSet
	Refer                    *model.Catalog
}

func Notify(ctx context.Context, mc model.ClientSet, refer *model.Catalog) error {
	return bus.Publish(ctx, BusMessage{TransactionalModelClient: mc, Refer: refer})
}

func AddSubscriber(n string, h func(context.Context, BusMessage) error) error {
	return bus.Subscribe(n, h)
}
