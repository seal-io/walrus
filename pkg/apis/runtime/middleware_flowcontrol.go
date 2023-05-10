package runtime

import (
	"context"
	"errors"
	"net/http"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RequestCounting limits the request count in the given maximum,
// returns 429 if the new request has waited with the given duration.
func RequestCounting(max int, wait time.Duration) Handle {
	if max <= 0 {
		return func(c *gin.Context) {
			c.AbortWithStatus(http.StatusTooManyRequests)
		}
	}

	limiter := make(chan struct{}, max)
	var token struct{}

	if wait <= 0 {
		return func(c *gin.Context) {
			select {
			default:
				if c.Err() == nil {
					c.AbortWithStatus(http.StatusTooManyRequests)
					return
				}
				c.Abort()
			case limiter <- token:
				defer func() { <-limiter }()
				c.Next()
			}
		}
	}

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, wait)
		defer cancel()
		select {
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				c.AbortWithStatus(http.StatusTooManyRequests)
				return
			}
			c.Abort()
		case limiter <- token:
			defer func() { <-limiter }()
			c.Next()
		}
	}
}

// RequestThrottling controls the request count per second and allows bursting,
// returns 429 if the new request is not allowed.
func RequestThrottling(qps int, burst int) Handle {
	if qps <= 0 || burst <= 0 {
		return func(c *gin.Context) {
			c.AbortWithStatus(http.StatusTooManyRequests)
		}
	}

	limiter := rate.NewLimiter(rate.Limit(qps), burst)
	return func(c *gin.Context) {
		if !limiter.Allow() {
			if c.Err() == nil {
				c.AbortWithStatus(http.StatusTooManyRequests)
				return
			}
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequestShaping arranges all requests to be received on the given qps,
// returns 429 if the new request can be allowed within the given latency,
// if the given latency is not positive, RequestShaping will never return 429.
func RequestShaping(qps int, slack int, latency time.Duration) Handle {
	if qps <= 0 {
		return func(c *gin.Context) {
			c.AbortWithStatus(http.StatusTooManyRequests)
		}
	}

	type state struct {
		arrival time.Time
		sleep   time.Duration
	}
	window := time.Second / time.Duration(qps)
	maxSleep := -1 * time.Duration(slack) * window
	statePointer := func() unsafe.Pointer {
		var s state
		return unsafe.Pointer(&s)
	}()
	return func(c *gin.Context) {
		for {
			select {
			case <-c.Done():
				c.Abort()
				return
			default:
			}

			prevStatePointer := atomic.LoadPointer(&statePointer)
			prevState := (*state)(prevStatePointer)
			currState := state{
				arrival: time.Now(),
				sleep:   prevState.sleep,
			}

			// For first request.
			if prevState.arrival.IsZero() {
				taken := atomic.CompareAndSwapPointer(&statePointer,
					prevStatePointer, unsafe.Pointer(&currState))
				if !taken {
					continue
				}
				// Allow it immediately.
				c.Next()
				return
			}

			// For subsequent requests.
			currState.sleep += window - currState.arrival.Sub(prevState.arrival)
			if currState.sleep < maxSleep {
				currState.sleep = maxSleep
			}
			var wait time.Duration
			if currState.sleep > 0 {
				currState.arrival = currState.arrival.Add(currState.sleep)
				wait, currState.sleep = currState.sleep, 0
			}
			if latency > 0 && wait > latency {
				c.AbortWithStatus(http.StatusTooManyRequests)
				return
			}
			taken := atomic.CompareAndSwapPointer(&statePointer, prevStatePointer, unsafe.Pointer(&currState))
			if !taken {
				continue
			}
			// Allow it after waiting.
			t := time.NewTimer(wait)
			select {
			case <-t.C:
				t.Stop()
				c.Next()
			case <-c.Done():
				t.Stop()
				c.Abort()
			}
			return
		}
	}
}
