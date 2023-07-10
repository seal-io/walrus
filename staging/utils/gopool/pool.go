package gopool

import (
	"errors"
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
		pond.PanicHandler(func(i any) { log.WithName("gopool").Errorf("panic observing: %v", i) }))
}

func ResetPool(factor int) {
	once.Do(func() {
		gp.Stop()
		gp = setPool(factor)
	})
}

// Go submits a task as goroutine.
func Go(f func()) {
	if !gp.TrySubmit(f) {
		log.WithName("gopool").Warn("goroutine pool full")
		gp.Submit(f)
	}
}

// TryGo tries to submit a task as goroutine.
func TryGo(f func()) bool {
	r := gp.TrySubmit(f)

	return r
}

// IsHealthy returns true if the pool has plenty workers.
func IsHealthy(atLeast ...int) error {
	watermark := 0
	if len(atLeast) > 0 {
		watermark = atLeast[0]
	}

	if watermark < 0 {
		watermark = 0
	}

	if gp.IdleWorkers() > watermark {
		return nil
	}

	return errors.New("goroutine pool full")
}
