package gopool

import (
	"context"
	"fmt"
	"sync"

	"github.com/alitto/pond"
	"go.uber.org/multierr"

	"github.com/seal-io/seal/utils/log"
)

// Group returns a waiting group,
// which closes at all tasks finishing and aggregates errors from tasks.
func Group() *waitGroup {
	return &waitGroup{}
}

type waitGroup struct {
	g   sync.WaitGroup
	m   sync.Mutex
	err error
}

// Wait blocks until all tasks completed and aggregates errors from tasks.
func (g *waitGroup) Wait() error {
	g.g.Wait()
	return g.err
}

// Go submits a task as goroutine.
func (g *waitGroup) Go(f func() error) {
	if f == nil {
		return
	}

	var wf = func() (err error) {
		defer func() {
			if v := recover(); v != nil {
				switch vt := v.(type) {
				case error:
					err = fmt.Errorf("panic as %w", vt)
				default:
					err = fmt.Errorf("panic as %v", v)
				}
				log.WithName("gopool").Errorf("panic observing: %v", err)
			}
		}()
		return f()
	}

	g.g.Add(1)
	Go(func() {
		defer g.g.Done()
		var err = wf()
		if err != nil {
			g.m.Lock()
			g.err = multierr.Append(g.err, err)
			g.m.Unlock()
		}
	})
}

// GroupWithContext returns a waiting group and a context derived by the given context.Context.
// Waiting group notifies closing when any task raises error,
// any submitting task should use the returning context to receive quiting.
func GroupWithContext(ctx context.Context) (contextWaitGroup, context.Context) {
	var g, c = gp.GroupContext(ctx)
	return contextWaitGroup{g: g}, c
}

type contextWaitGroup struct {
	g *pond.TaskGroupWithContext
}

// Wait blocks until either all tasks completed or
// one of them returned a non-nil error or the context associated to this group
// was canceled.
func (g contextWaitGroup) Wait() error {
	return g.g.Wait()
}

// Go submits a task as goroutine.
func (g contextWaitGroup) Go(f func() error) {
	if f == nil {
		return
	}

	var wf = func() (err error) {
		defer func() {
			if v := recover(); v != nil {
				switch vt := v.(type) {
				case error:
					err = fmt.Errorf("panic as %w", vt)
				default:
					err = fmt.Errorf("panic as %v", v)
				}
				log.WithName("gopool").Errorf("panic observing: %v", err)
			}
		}()
		return f()
	}

	g.g.Submit(wf)
	printState()
}

// GroupWithContextIn is similar as GroupWithContext but doesn't return a derived context,
// all tasks can receive the derived context at submitting, a kind of more compact usage.
func GroupWithContextIn(ctx context.Context) (g embeddedContextWaitGroup) {
	g.g, g.c = GroupWithContext(ctx)
	return
}

type embeddedContextWaitGroup struct {
	g contextWaitGroup
	c context.Context
}

// Wait blocks until either all tasks completed or
// one of them returned a non-nil error or the context associated to this group
// was canceled.
func (g embeddedContextWaitGroup) Wait() error {
	return g.g.Wait()
}

// Go submits a task as goroutine.
func (g embeddedContextWaitGroup) Go(f func(context.Context) error) {
	if f == nil {
		return
	}

	g.g.Go(func() error {
		return f(g.c)
	})
}
