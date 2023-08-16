package catalog

import (
	"context"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/utils/bus"
)

type BusMessage struct {
	// TransactionalModelClient holds the model.ClientSet of this calling session,
	// it should be a transactional DAO client,
	// please don't keep for long-term using.
	TransactionalModelClient model.ClientSet
	// Refer holds the updating model.Catalog item of this calling session.
	Refer *model.Catalog
}

func Notify(ctx context.Context, mc model.ClientSet, refer *model.Catalog) error {
	return bus.Publish(ctx, BusMessage{TransactionalModelClient: mc, Refer: refer})
}

func AddSubscriber(n string, h func(context.Context, BusMessage) error) error {
	return bus.Subscribe(n, h)
}
