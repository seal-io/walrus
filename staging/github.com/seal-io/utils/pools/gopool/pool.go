package gopool

import (
	"errors"
	"fmt"
	"runtime"
	"sync"

	"github.com/alitto/pond"
	"k8s.io/klog/v2"
)

var (
	once sync.Once
	gp   = NewPool(100)
)

func NewPool(factor int) *pond.WorkerPool {
	// NB(thxCode): Go allows us to create goroutines at will, but if we create too many goroutines,
	// it will cause the program to crash due to insufficient memory,
	// so we need to limit the number of goroutines with pooling.
	//
	// The advantage of pooling is that space is exchanged for time and the reuse rate is improved.
	//
	// - MaxWorkers is the total goroutine number should the pool creates,
	//   we take the number of available CPU cores as the basic value at present,
	//   then times the given factor to get the max workers.
	//
	// - MinWorkers is the goroutine number value should the pool creates at begin,
	//   we take the MaxWorkers as the result if it is less than 500 at present,
	//   otherwise, we limit the MinWorkers value to avoid the pool creates too many goroutines at begin.
	//
	// - MaxCapacity is the max value of submitting goroutine number should be accepted at the same time,
	//   we take 80% of the MaxWorkers as the result if it is greater than 100 at present.
	maxWorkers := runtime.GOMAXPROCS(0) * factor

	minWorkers := maxWorkers
	if minWorkers > 500 {
		minWorkers = 500
	}

	maxCapacity := maxWorkers * 8 / 10
	if maxCapacity < 100 {
		maxCapacity = 100
	}

	return pond.New(maxWorkers, maxCapacity,
		pond.MinWorkers(minWorkers),
		pond.Strategy(pond.Eager()),
		pond.PanicHandler(func(i any) {
			klog.Background().WithName("gopool").
				Error(fmt.Errorf("%v", i), "panic observing")
		}))
}

// Configure sets the goroutine pool with a new factor once,
// multiple calls will be ignored.
func Configure(factor int) {
	once.Do(func() {
		gp.Stop()
		gp = NewPool(factor)
	})
}

// Go submits a task as goroutine.
func Go(f func()) {
	if !gp.TrySubmit(f) {
		klog.Background().WithName("gopool").
			V(5).
			Info("goroutine pool full")
		gp.Submit(f)
	}
}

// TryGo tries to submit a task as goroutine.
func TryGo(f func()) bool {
	r := gp.TrySubmit(f)

	return r
}

// IsHealthy returns nil if the pool has plenty workers.
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

// Burst returns the maximum number of goroutines to submit at the same moment.
func Burst() int {
	l, r := gp.IdleWorkers(), gp.MaxCapacity()
	if l > r {
		return r
	}

	return l
}
