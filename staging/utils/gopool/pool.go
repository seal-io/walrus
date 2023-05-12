package gopool

import (
	"runtime"

	"github.com/alitto/pond"

	"github.com/seal-io/seal/utils/log"
)

var gp = pond.New(runtime.NumCPU()*10, runtime.NumCPU()*1000,
	pond.MinWorkers(runtime.NumCPU()*10),
	pond.Strategy(pond.Eager()),
	pond.PanicHandler(func(i interface{}) { log.WithName("gopool").Errorf("panic observing: %v", i) }))

func printState() {
	if !log.Enabled(log.DebugLevel) {
		return
	}

	log.WithName("gopool").Debugf("state: tasks %d/%d workers %d/%d",
		gp.WaitingTasks(), gp.SubmittedTasks(),
		gp.IdleWorkers(), gp.RunningWorkers())
}

// Go submits a task as goroutine.
func Go(f func()) {
	if !gp.TrySubmit(f) {
		log.WithName("gopool").Warn("goroutine pool full")
		gp.Submit(f)
	}

	printState()
}

// TryGo tries to submit a task as goroutine.
func TryGo(f func()) bool {
	r := gp.TrySubmit(f)

	printState()

	return r
}
