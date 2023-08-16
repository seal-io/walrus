package environment

import (
	"context"

	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/utils/bus"
)

type Event int

const (
	EventCreate = iota
	EventUpdate
	EventDelete
)

// BusMessage wraps the changed model.Environment as a bus.Message.
type BusMessage struct {
	// Event holds event type.
	Event Event
	// TransactionalModelClient holds the model.ClientSet of this calling session,
	// it should be a transactional DAO client,
	// please don't keep for long-term using.
	TransactionalModelClient model.ClientSet
	// Refers holds the updating model.Environment list of this calling session.
	Refers model.Environments
}

// NotifyIDs notifies the changed model.Environment IDs.
func NotifyIDs(ctx context.Context, mc model.ClientSet, event Event, ids ...object.ID) error {
	envs, err := dao.GetEnvironmentsByIDs(ctx, mc, ids...)
	if err != nil {
		return err
	}

	return bus.Publish(ctx, BusMessage{
		Event:                    event,
		TransactionalModelClient: mc,
		Refers:                   envs,
	})
}

// Notify notifies the changed model.Environment.
func Notify(ctx context.Context, mc model.ClientSet, event Event, refers model.Environments) error {
	return bus.Publish(ctx, BusMessage{
		Event:                    event,
		TransactionalModelClient: mc,
		Refers:                   refers,
	})
}

// AddSubscriber add the subscriber to handle the changed notification from model.Environment.
func AddSubscriber(n string, h func(context.Context, BusMessage) error) error {
	return bus.Subscribe(n, h)
}
