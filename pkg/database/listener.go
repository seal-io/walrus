package database

import (
	"context"
	"errors"
	"time"

	"github.com/lib/pq"

	"github.com/seal-io/walrus/utils/log"
)

type (
	// ListenHandler holds the operations for the  database listen handler.
	ListenHandler interface {
		// Channels returns the name list of multiple channels to establish.
		Channels() []string

		// Handle handles the payload according to the given channel.
		Handle(ctx context.Context, channel, payload string)
	}

	// Listener subscribes the database listen channel and handle the event.
	Listener struct {
		logger               log.Logger
		minReconnectInterval time.Duration
		maxReconnectInterval time.Duration
		dataSourceAddress    string
		handlers             map[string][]ListenHandler
	}
)

// NewListener create database listener with options.
func NewListener(dataSourceAddress string, opts ...Option) (*Listener, error) {
	if dataSourceAddress == "" {
		return nil, errors.New("blank data source address")
	}

	l := &Listener{
		logger:               log.WithName("database-listener"),
		maxReconnectInterval: 90 * time.Second,
		minReconnectInterval: 15 * time.Second,
		dataSourceAddress:    dataSourceAddress,
		handlers:             make(map[string][]ListenHandler),
	}
	for _, opt := range opts {
		opt(l)
	}

	return l, nil
}

// Option is a function that configures a listener.
type Option func(*Listener)

// WithReconnectInterval can be used to set the min and max reconnect interval.
func WithReconnectInterval(min, max time.Duration) Option {
	return func(l *Listener) {
		l.minReconnectInterval = min
		l.maxReconnectInterval = max
	}
}

// Register listens the channel returns by the given handler
// and handles by the given handler if event received.
func (l *Listener) Register(handler ListenHandler) error {
	if handler == nil {
		return errors.New("nil handler")
	}

	for _, ch := range handler.Channels() {
		l.handlers[ch] = append(l.handlers[ch], handler)
	}

	return nil
}

// eventCallBack is the callback function for the database listener.
func (l *Listener) eventCallBack(ev pq.ListenerEventType, err error) {
	switch ev {
	case pq.ListenerEventConnected:
		l.logger.Info("connected")
	case pq.ListenerEventDisconnected:
		l.logger.Errorf("disconnected: %v", err)
	case pq.ListenerEventReconnected:
		l.logger.Warnf("reconnected")
	case pq.ListenerEventConnectionAttemptFailed:
		l.logger.Errorf("connection attempt failed: %v", err)
	}
}

// Start setups a database listener to subscribe the established channels,
// and stops when the given context is done.
func (l *Listener) Start(ctx context.Context) error {
	// Create listener.
	dl := pq.NewListener(l.dataSourceAddress,
		l.minReconnectInterval, l.maxReconnectInterval, l.eventCallBack)
	defer func() {
		_ = dl.Close()
	}()

	// Listen channels.
	for ch := range l.handlers {
		err := dl.Listen(ch)
		if err != nil {
			return err
		}
	}

	// Receive events.
	for {
		select {
		case msg := <-dl.Notify:
			l.logger.V(5).InfoS("received event",
				"channel", msg.Channel, "payload", msg.Extra)

			for _, h := range l.handlers[msg.Channel] {
				h.Handle(ctx, msg.Channel, msg.Extra)
			}
		case <-ctx.Done():
			return nil
		}
	}
}
