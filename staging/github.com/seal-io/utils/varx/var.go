package varx

import (
	"sync"
)

// Once holds the value of the variable and ensures that the value is configured only once.
type Once[T any] struct {
	o sync.Once

	v T
}

// NewOnce returns the Once with the given value.
func NewOnce[T any](value T) *Once[T] {
	return &Once[T]{
		v: value,
	}
}

// Configure sets the value, multiple calls will be ignored.
func (i *Once[T]) Configure(value T) {
	i.o.Do(func() {
		i.v = value
	})
}

// Get returns the configured value,
// or returns initialized value with NewOnce if the variable is not configured yet.
func (i *Once[T]) Get() T {
	return i.v
}

// Notify holds the value of variable and ensures that the value is configured before get.
type Notify[T any] interface {
	// Configure sets the value.
	Configure(value T)
	// IsConfigured returns true if the value is configured.
	IsConfigured() bool
	// Get blocks until the value is configured at the first time,
	// and returns the configured result.
	Get() T
}

type notify[T any] struct {
	m *sync.RWMutex
	c *sync.Cond
	s bool

	v T
}

// NewNotify creates a new Notify.
func NewNotify[T any]() Notify[T] {
	var m sync.RWMutex

	return &notify[T]{
		m: &m,
		c: sync.NewCond(&m),
	}
}

func (i *notify[T]) Configure(value T) {
	i.m.Lock()
	i.v = value
	i.s = true
	i.m.Unlock()

	i.c.Broadcast()
}

func (i *notify[T]) IsConfigured() bool {
	i.m.RLock()
	defer i.m.RUnlock()

	return i.s
}

func (i *notify[T]) Get() T {
	i.m.RLock()
	if i.s {
		i.m.RUnlock()
		return i.v
	}
	i.m.RUnlock()

	i.m.Lock()
	defer i.m.Unlock()
	for !i.s {
		i.c.Wait()
	}

	return i.v
}
