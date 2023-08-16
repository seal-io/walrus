package database

import (
	"context"
	"time"

	"github.com/lib/pq"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/utils/json"
	"github.com/seal-io/walrus/utils/log"
)

// globalListener is the default global database listener.
var globalListener = NewListener()

// ListenChannel listen to the channel and add event handler to global database listener.
func ListenChannel(name string, handler ListenHandler) {
	globalListener.Listen(name, handler)
}

// StartListener begin to listen and receive event for global database listener.
func StartListener(ctx context.Context, dsa string, client model.ClientSet) error {
	return globalListener.Start(ctx, dsa, client)
}

// Record is the database change event record.
type Record struct {
	TableName string `json:"tableName"`
	Operation string `json:"operation"`
	RecordID  uint64 `json:"recordID"`
}

// ListenHandler is the handler to handle the database table change event.
type ListenHandler func(ctx context.Context, client model.ClientSet, rd Record) error

// Listener listen to the database table changes.
type Listener struct {
	minReconnectInterval time.Duration
	maxReconnectInterval time.Duration

	handlers map[string]ListenHandler
}

// NewListener create database listener with options.
func NewListener(opts ...Option) *Listener {
	l := &Listener{
		maxReconnectInterval: 90 * time.Second,
		minReconnectInterval: 15 * time.Second,
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

// WithReconnectInterval can be used to set the min and max reconnect interval.
func WithReconnectInterval(min, max time.Duration) Option {
	return func(l *Listener) {
		l.minReconnectInterval = min
		l.maxReconnectInterval = max
	}
}

// Option is a function that configures a listener.
type Option func(*Listener)

// Listen the channel and add event handler to listener.
func (l *Listener) Listen(channelName string, handler ListenHandler) {
	if l.handlers == nil {
		l.handlers = make(map[string]ListenHandler)
	}

	l.handlers[channelName] = handler
}

// Start begin to listen and receive event for database listener.
func (l *Listener) Start(ctx context.Context, dsa string, client model.ClientSet) error {
	logger := log.WithName("database-listener")

	// Report error function.
	reportError := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			logger.Errorf("error listen table event, type %s: %v", ev, err)
		}
	}

	// Create listener.
	dl := pq.NewListener(dsa, l.minReconnectInterval, l.maxReconnectInterval, reportError)
	defer func() {
		_ = dl.Close()
	}()

	// Listen channels.
	for channelName := range l.handlers {
		err := dl.Listen(channelName)
		if err != nil {
			return err
		}
	}

	// Receive table event.
	for {
		select {
		case msg := <-dl.Notify:
			if msg.Extra == "" {
				continue
			}

			logger.V(5).Infof("received table event, channel %s: %s", msg.Channel, msg.Extra)

			handler, ok := l.handlers[msg.Channel]
			if !ok {
				logger.V(5).Infof("received table event for unwatched channel %s, skipped", msg.Channel)
				continue
			}

			var rd Record

			err := json.Unmarshal([]byte(msg.Extra), &rd)
			if err != nil {
				logger.Warnf("error unmarshal receiving table event, channel %s: %v", msg.Channel, err)
				continue
			}

			err = handler(ctx, client, rd)
			if err != nil {
				logger.Warnf("error notify table event %s %s %s: %v",
					rd.TableName, rd.Operation, rd.RecordID, err)
			}
		case <-ctx.Done():
			return nil
		}
	}
}
