package gopool

import (
	"context"
	"runtime"
	"sync"

	"github.com/alitto/pond"

	"github.com/seal-io/seal/utils/log"
)

var gp = pond.New(runtime.NumCPU()*10, runtime.NumCPU()*1000,
	pond.MinWorkers(runtime.NumCPU()*10),
	pond.Strategy(pond.Eager()),
	pond.PanicHandler(func(i interface{}) { log.Error(i) }))

func printState() {
	if !log.Enabled(log.DebugLevel) {
		return
	}
	log.WithName("gopool").Debugf("state: tasks %d/%d workers %d/%d",
		gp.WaitingTasks(), gp.SubmittedTasks(),
		gp.IdleWorkers(), gp.RunningWorkers())
}

func Go(f func()) {
	if !gp.TrySubmit(f) {
		log.WithName("gopool").Warn("goroutine pool full")
		gp.Submit(f)
	}
	printState()
}

func TryGo(f func()) bool {
	var r = gp.TrySubmit(f)
	printState()
	return r
}

type Group struct {
	g *pond.TaskGroupWithContext
}

func (g Group) Go(f func() error) {
	g.g.Submit(f)
	printState()
}

func (g Group) Wait() error {
	return g.g.Wait()
}

func GroupWithContext(ctx context.Context) (Group, context.Context) {
	var g, c = gp.GroupContext(ctx)
	return Group{g: g}, c
}

type WrapWaitGroup struct {
	wg sync.WaitGroup

	errOnce sync.Once
	err     error
}

func (g *WrapWaitGroup) Wait() error {
	g.wg.Wait()
	return g.err
}

func (g *WrapWaitGroup) Go(f func() error) {
	g.wg.Add(1)
	Go(func() {
		defer g.wg.Done()

		if err := f(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
			})
		}
	})
}

func WaitGroup() *WrapWaitGroup {
	return &WrapWaitGroup{}
}
