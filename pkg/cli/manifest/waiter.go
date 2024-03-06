package manifest

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/util/sets"
	utilwait "k8s.io/apimachinery/pkg/util/wait"

	"github.com/seal-io/walrus/pkg/cli/api"
	"github.com/seal-io/walrus/pkg/cli/common"
	"github.com/seal-io/walrus/pkg/cli/config"
	"github.com/seal-io/walrus/pkg/cli/formatter"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/utils/gopool"
	"github.com/seal-io/walrus/utils/json"
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

// DefaultPreviewObjectWaiter used for wait for preview generated.
func DefaultPreviewObjectWaiter(sc *config.Config, timeoutSecond int) *PreviewObjectWaiter {
	timeout := time.Duration(timeoutSecond) * time.Second
	return &PreviewObjectWaiter{
		serverContext: sc,
		backoff: &utilwait.Backoff{
			Duration: 200 * time.Millisecond,
			Factor:   2,
			Steps:    8,
			Cap:      timeout,
		},
		timeout: timeout,
	}
}

// PreviewObjectWaiter is a type that represents preview object waiter.
type PreviewObjectWaiter struct {
	serverContext *config.Config
	timeout       time.Duration
	backoff       *utilwait.Backoff
}

// Wait waits for the preview generated and print the error.
func (o *PreviewObjectWaiter) Wait(ctx context.Context, set ObjectSet, _ chan ObjectSet) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, o.timeout)
	defer cancel()

	var (
		resultChan = make(chan result, 1)
		all        = sets.Set[string]{}
		finished   = sets.Set[string]{}
	)
	for _, s := range set.All() {
		obj := s
		all.Insert(obj.Name)

		gopool.Go(func() {
			err := o.waitPreviewObject(obj)
			resultChan <- result{
				err: err,
				obj: obj,
			}
		})
	}

	for {
		select {
		default:
			if finished.Equal(all) {
				return true, nil
			}
		case <-ctx.Done():
			return false, fmt.Errorf("timeout")
		case ch := <-resultChan:
			finished.Insert(ch.obj.Name)
			if ch.err != nil {
				fmt.Fprintf(os.Stderr, "resource %s preview run error, %v\n", ch.obj.ScopedName(ch.obj.Name), ch.err)
			}
		}
	}
}

// waitPreviewObject waits for the object preview finished.
func (o *PreviewObjectWaiter) waitPreviewObject(obj Object) error {
	listOpt := api.OpenAPI.GetOperation(GroupResourceRuns, operationList)
	if listOpt == nil {
		return fmt.Errorf("not found %s operation for %s", operationList, GroupResourceRuns)
	}

	if o.backoff == nil {
		return o.getResourceRuns(obj, listOpt)
	}

	var (
		err      error
		retryErr *common.RetryableError
	)
	err = utilwait.ExponentialBackoff(*o.backoff, func() (bool, error) {
		err := o.getResourceRuns(obj, listOpt)
		if err != nil {
			if errors.As(err, &retryErr) {
				err = nil
			}
			return false, err
		}

		return true, nil
	})

	return err
}

type collectionGetResourceRun struct {
	Items []*model.ResourceRunOutput `json:"items"`
}

func (o *PreviewObjectWaiter) getResourceRuns(obj Object, listOpt *api.Operation) error {
	m := obj.Map()
	m["resource"] = obj.Name
	m["order"] = "-createTime"
	m["perPage"] = 1
	m["page"] = 1
	m["preview"] = true

	req, err := listOpt.Request(nil, []string{obj.Name}, m, nil, o.serverContext.ServerContext)
	if err != nil {
		return err
	}

	resp, err := o.serverContext.DoRequest(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	err = common.CheckResponseStatus(resp)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error read response from list resource run: %w", err)
	}

	var runs collectionGetResourceRun

	err = json.Unmarshal(body, &runs)
	if err != nil {
		return fmt.Errorf("error unmarshal response from list resource run: %w", err)
	}

	if len(runs.Items) == 0 {
		return common.NewRetryableError("resource run list is empty")
	}

	run := runs.Items[0]
	switch run.Status.SummaryStatus {
	default:
		return fmt.Errorf("resource run %s is in unexpected status %s", run.ID, run.Status.SummaryStatus)
	case "Planning":
		return common.NewRetryableError("preview generate is planning")
	case "Canceled":
		fmt.Printf("resource %s preview generate is canceled\n", obj.ScopedName(obj.Name))
		return nil
	case "Planned":
		fmt.Printf("resource %s preview is generated, changes: %d created, %d updated, %d deleted\n",
			obj.ScopedName(obj.Name), run.ComponentChangeSummary.Created, run.ComponentChangeSummary.Updated, run.ComponentChangeSummary.Deleted)

		if len(run.ComponentChanges) == 0 {
			return nil
		}

		b, err := json.Marshal(run.ComponentChanges)
		if err != nil {
			return err
		}

		columns := []string{"type", "change.type"}
		t := formatter.DefaultTableFormatter(columns, "", "")
		t.ColumnDisplayName = func(s string) string {
			mapping := map[string]string{
				"name":        "NAME",
				"type":        "TYPE",
				"change.type": "CHANGE TYPE",
			}
			if mp, ok := mapping[s]; ok {
				return mp
			}
			return strings.ToUpper(s)
		}

		fmt.Println(t.ResourceItems(b))
	}

	return nil
}
