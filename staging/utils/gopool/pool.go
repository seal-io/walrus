package gopool

import (
	"context"
	"runtime"

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

func WithContext(ctx context.Context) (Group, context.Context) {
	var g, c = gp.GroupContext(ctx)
	return Group{g: g}, c
}
