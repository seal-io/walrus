package gopool

import (
	"sync"

	"github.com/alitto/pond"

	"github.com/seal-io/seal/utils/log"
	"github.com/seal-io/seal/utils/runtimex"
)

var (
	once sync.Once
	gp   = setPool(100)
)

func setPool(factor int) *pond.WorkerPool {
	b := runtimex.NumCPU()

	return pond.New(b*factor, b*factor*100,
		pond.MinWorkers(b*factor),
		pond.Strategy(pond.Eager()),
		pond.PanicHandler(func(i interface{}) { log.WithName("gopool").Errorf("panic observing: %v", i) }))
}

func ResetPool(factor int) {
	once.Do(func() {
		gp.Stop()
		gp = setPool(factor)
	})
}

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
