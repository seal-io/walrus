package manifest

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/seal-io/walrus/pkg/cli/config"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/utils/strs"
)

// Waiter is an interface that defines methods for waiting for the status of objects.
type Waiter interface {
	Wait(ctx context.Context, set ObjectSet, finished chan ObjectSet) (bool, error)
}

// StatusWaiter creates a Waiter wait for status ready or error.
func StatusWaiter(sc *config.Config, timeoutSecond int) Waiter {
	return &ObjectWaiter{
		serverContext: sc,
		timeout:       time.Duration(timeoutSecond) * time.Second,
		cond: map[string]WaitFor{
			GroupResources: WaitForStatusReadyOrError,
		},
	}
}

// DeleteWaiter creates a Waiter wait for delete.
func DeleteWaiter(sc *config.Config, timeoutSecond int) *ObjectWaiter {
	return &ObjectWaiter{
		serverContext: sc,
		timeout:       time.Duration(timeoutSecond) * time.Second,
		cond: map[string]WaitFor{
			GroupResources: WaitForDelete,
		},
	}
}

// ObjectWaiter is a type that represents an object waiter.
type ObjectWaiter struct {
	serverContext *config.Config
	timeout       time.Duration
	cond          map[string]WaitFor
}

// Wait waits for the object to finished.
func (o *ObjectWaiter) Wait(ctx context.Context, objs ObjectSet, finished chan ObjectSet) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, o.timeout)
	defer cancel()

	return WatchObjects(ctx, o.serverContext, objs, o.cond, finished)
}

var waitCache sync.Map

// WaitFor is a function that defines the condition for waiting for the object.
type WaitFor = func(typ string, e EventItem, group, name string) bool

// WaitForStatusReadyOrError waits for the object to be ready or error.
func WaitForStatusReadyOrError(_ string, e EventItem, group, name string) bool {
	var (
		key  = fmt.Sprintf("%s/%s", group, name)
		msg  = fmt.Sprintf("%s %s is in status: %s", strs.Singularize(group), name, e.Status.SummaryStatus)
		meet = e.Status.SummaryStatus == string(status.ResourceStatusReady) || e.Status.Error
	)

	val, ok := waitCache.Load(key)
	if !ok || val.(string) != msg {
		fmt.Println(msg)
	}

	waitCache.Store(key, msg)

	return meet
}

// WaitForDelete waits for the object to be deleted.
func WaitForDelete(typ string, e EventItem, group, name string) bool {
	var (
		key  = fmt.Sprintf("%s/%s", group, name)
		msg  string
		meet = typ == "delete"
	)

	if meet {
		msg = fmt.Sprintf("%s %s is deleted", strs.Singularize(group), name)
	} else {
		msg = fmt.Sprintf("%s %s is in status: %s", strs.Singularize(group), name, e.Status.SummaryStatus)
	}

	val, ok := waitCache.Load(key)
	if !ok || val.(string) != msg {
		fmt.Println(msg)
	}

	waitCache.Store(key, msg)

	return meet
}
