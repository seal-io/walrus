package contextx

import (
	"context"
	"time"

	"github.com/seal-io/utils/pools/gopool"
)

func Background(stop <-chan struct{}) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	gopool.Go(func() {
		select {
		case <-stop:
		case <-ctx.Done():
		}
		cancel()
	})
	return ctx
}

func TODO(stop <-chan struct{}) context.Context {
	ctx, cancel := context.WithCancel(context.TODO())
	gopool.Go(func() {
		select {
		case <-stop:
		case <-ctx.Done():
		}
		cancel()
	})
	return ctx
}

func WithCancel(stop <-chan struct{}) (context.Context, context.CancelFunc) {
	return context.WithCancel(Background(stop))
}

func WithCancelCause(stop <-chan struct{}) (context.Context, context.CancelCauseFunc) {
	return context.WithCancelCause(Background(stop))
}

func WithDeadline(stop <-chan struct{}, d time.Time) (context.Context, context.CancelFunc) {
	return context.WithDeadline(Background(stop), d)
}

func WithDeadlineCause(stop <-chan struct{}, d time.Time, cause error) (context.Context, context.CancelFunc) {
	return context.WithDeadlineCause(Background(stop), d, cause)
}

func WithTimeout(stop <-chan struct{}, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(Background(stop), timeout)
}

func WithTimeoutCause(stop <-chan struct{}, timeout time.Duration, cause error) (context.Context, context.CancelFunc) {
	return context.WithTimeoutCause(Background(stop), timeout, cause)
}

func WithValue(stop <-chan struct{}, key, val any) context.Context {
	return context.WithValue(Background(stop), key, val)
}
