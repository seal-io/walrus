package modelchange

import (
	"context"
	"runtime"
	"sync"
	"time"

	"go.uber.org/multierr"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/utils/gopool"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/topic"
)

type (
	topicalEvent struct {
		topic.Topic
		Event
	}

	eventBuffer struct {
		sync.RWMutex

		logger log.Logger
		ticker *time.Ticker
		buffer []topicalEvent
		end    int
	}
)

func newEventBuffer(ctx context.Context, logger log.Logger) *eventBuffer {
	b := &eventBuffer{
		logger: logger,
		ticker: time.NewTicker(2 * time.Second),
		buffer: make([]topicalEvent, 32),
	}

	// Flush buffer periodically in background.
	gopool.Go(func() {
		for range b.ticker.C {
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			if err := b.flush(ctx, true); err != nil {
				b.logger.Warnf("failed to flush: %v", err)
			}

			cancel()
		}
	})
	runtime.SetFinalizer(b, (*eventBuffer).Close)

	return b
}

// Write writes an event to the buffer.
func (b *eventBuffer) Write(ctx context.Context, topic topic.Topic, event Event) {
	// Merge events if the buffer is full.
	if b.isFull() && !b.merge() {
		// Flush if the buffer is still full.
		b.logger.Debug("buffer is full, flushing")

		if err := b.flush(ctx, false); err != nil {
			b.logger.Warnf("failed to flush: %v", err)
		}
	}

	// Write event to buffer.
	b.Lock()
	defer b.Unlock()
	b.buffer[b.end] = topicalEvent{Topic: topic, Event: event}
	b.end++
}

// isFull returns true if the buffer is full.
func (b *eventBuffer) isFull() bool {
	b.RLock()
	defer b.RUnlock()

	return b.end == len(b.buffer)
}

// popAll pops all events from the buffer.
func (b *eventBuffer) popAll() []topicalEvent {
	b.Lock()
	defer b.Unlock()

	buff := make([]topicalEvent, b.end)
	if b.end != 0 {
		copy(buff, b.buffer[:b.end])
		b.end = 0
	}

	return buff
}

// merge merges the buffer in place,
// returns true if merged.
func (b *eventBuffer) merge() bool {
	b.Lock()
	defer b.Unlock()

	end := len(mergeEvents(b.buffer[:b.end]))
	if b.end != end {
		b.end = end
		return true
	}

	return false
}

func (b *eventBuffer) Close() error {
	b.ticker.Stop()
	return nil
}

// flush flushes the buffer to the topic.
func (b *eventBuffer) flush(ctx context.Context, merge bool) error {
	events := b.popAll()
	if len(events) == 0 {
		return nil
	}

	if merge {
		events = mergeEvents(events)
	}

	var err error

	for i := range events {
		b.logger.V(5).InfoS("sending",
			"topic", events[i].Topic, "event", events[i].Event)
		err = multierr.Append(err,
			topic.Publish(ctx, events[i].Topic, events[i].Event))
	}

	return err
}

// mergeEvents merges the given event slice by type in place.
func mergeEvents(in []topicalEvent) []topicalEvent {
	if len(in) < 2 {
		return in
	}

	// Map events by topic.
	tes := map[topic.Topic][]Event{}
	for i := range in {
		tes[in[i].Topic] = append(tes[in[i].Topic], in[i].Event)
	}

	// Merge events by type.
	for et, es := range tes {
		data := [_EventTypeLength]sets.Set[EventData]{}
		for i := range es {
			if s := data[es[i].Type]; s != nil {
				s.Insert(es[i].Data...)
				continue
			}
			data[es[i].Type] = sets.New(es[i].Data...)
		}

		var i int

		// Ignore unknown events and empty sets,
		// and replace in place.
		for t := range data {
			if EventType(t) == _EventTypeUnknown ||
				data[t] == nil || data[t].Len() == 0 {
				continue
			}

			tes[et][i] = Event{Type: EventType(t), Data: data[t].UnsortedList()}
			i++
		}

		tes[et] = tes[et][:i]
	}

	// Reduce events by topic.
	var i int

	for et, es := range tes {
		// Replace in place.
		for j := range es {
			in[i] = topicalEvent{Topic: et, Event: es[j]}
			i++
		}
	}

	return in[:i]
}
