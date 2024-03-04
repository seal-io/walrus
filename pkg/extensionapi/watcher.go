package extensionapi

import (
	"runtime"

	"go.uber.org/atomic"
)

type WatchBookmark struct {
	rv *atomic.String
}

// NewWatchBookmark creates a new WatchBookmark.
func NewWatchBookmark() WatchBookmark {
	return WatchBookmark{
		rv: atomic.NewString("0"),
	}
}

// SwapResourceVersion swaps the current resource version with the new one if the new one is greater.
//
// SwapResourceVersion returns false if the new resource version is not swapped.
func (c WatchBookmark) SwapResourceVersion(newRv string) bool {
	for {
		oldRv := c.rv.Load()
		if oldRv == newRv ||
			len(oldRv) > len(newRv) ||
			(len(oldRv) == len(newRv) && oldRv > newRv) {
			return false
		}
		if c.rv.CompareAndSwap(oldRv, newRv) {
			return true
		}
		runtime.Gosched() // Yield to other goroutines.
	}
}

// GetResourceVersion returns the latest resource version.
func (c WatchBookmark) GetResourceVersion() string {
	return c.rv.Load()
}

// DeepCopy returns a deep copy of the WatchBookmark.
func (c WatchBookmark) DeepCopy() WatchBookmark {
	return WatchBookmark{
		rv: atomic.NewString(c.rv.Load()),
	}
}
