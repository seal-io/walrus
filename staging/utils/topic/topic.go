package topic

import (
	"context"
	"errors"
	"runtime"
	"sync"

	"github.com/google/uuid"
)

type Event struct {
	Topic Topic
	Data  Message
}

// Subscriber defines the interface to process the asynchronous logic.
type Subscriber interface {
	// Receive receives an Event from the subscribing Topic.
	Receive(context.Context) (Event, error)
	// Unsubscribe quits from the subscribing Topic.
	Unsubscribe()
}

// Hub defines the interface to decouple the asynchronous logic.
type Hub interface {
	// Subscribe gets a Subscriber from the Hub with the given name Topic.
	Subscribe(Topic) (Subscriber, error)
	// Unsubscribe quits all subscribing Subscriber from the given name Topic.
	Unsubscribe(n Topic) error
	// Publish publishes a Message to all Subscriber who is subscribing the given name Topic.
	Publish(context.Context, Topic, Message) error
}

// Message defines the message to be transferred.
type Message interface{}

// Topic defines the name of the Subscriber.
type Topic string

type hub struct {
	p *hub
	t Topic
	m sync.Map
}

func (h *hub) Subscribe(t Topic) (Subscriber, error) {
	if t == "" {
		// Topic scope.
		var n = uuid.NewString()
		var c = make(chan Event, runtime.NumCPU()*2)
		h.m.Store(n, c)
		return subscriber{p: h, n: n, c: c}, nil
	}
	// Hub scope.
	var v, _ = h.m.LoadOrStore(t, &hub{p: h, t: t})
	var sh = v.(*hub)
	return sh.Subscribe("")
}

func (h *hub) Unsubscribe(t Topic) error {
	if t == "" {
		// Topic scope.
		h.m.Range(func(n, v any) bool {
			var c = v.(chan Event)
			close(c)
			h.m.Delete(n)
			return true
		})
		return nil
	}
	// Hub scope.
	var v, ok = h.m.Load(t)
	if !ok {
		return nil
	}
	var sh = v.(*hub)
	return sh.Unsubscribe("")
}

func (h *hub) Publish(ctx context.Context, n Topic, m Message) error {
	if n == "" {
		// Topic scope.
		h.m.Range(func(n, v any) bool {
			var c = v.(chan Event)
			select {
			case <-ctx.Done():
				return false
			case c <- Event{Topic: h.t, Data: m}:
				return true
			default:
				// If chan is blocking.
				close(c)
				h.m.Delete(n)
				return true
			}
		})
		return nil
	}
	// Hub scope.
	var v, ok = h.m.Load(n)
	if !ok {
		return nil
	}
	var sh = v.(*hub)
	return sh.Publish(ctx, "", m)
}

type subscriber struct {
	p *hub
	n string
	c chan Event
}

func (s subscriber) Receive(ctx context.Context) (Event, error) {
	select {
	case <-ctx.Done():
		return Event{}, ctx.Err()
	case e, ok := <-s.c:
		if !ok {
			return Event{}, errors.New("topic is closed")
		}
		return e, nil
	}
}

func (s subscriber) Unsubscribe() {
	s.p.m.Delete(s.n)
}

var globalHub = New()

// New returns a new Hub.
func New() Hub {
	return &hub{}
}

// Subscribe gets a Subscriber from global Hub with the given topic Topic.
func Subscribe(n Topic) (Subscriber, error) {
	return globalHub.Subscribe(n)
}

// MustSubscribe likes Subscribe, but panic if error found.
func MustSubscribe(n Topic) Subscriber {
	var s, err = Subscribe(n)
	if err != nil {
		panic(err)
	}
	return s
}

// Publish publishes a Message to global Hub with the given topic Topic.
func Publish(ctx context.Context, n Topic, m Message) error {
	return globalHub.Publish(ctx, n, m)
}

// MustPublish likes Publish, but panic if error found.
func MustPublish(ctx context.Context, n Topic, m Message) {
	var err = Publish(ctx, n, m)
	if err != nil {
		panic(err)
	}
}
