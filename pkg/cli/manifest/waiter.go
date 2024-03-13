package manifest

import (
	"context"
	"errors"
	"fmt"
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

// DefaultPreviewWaiter wait for preview planned.
func DefaultPreviewWaiter(sc *config.Config, timeoutSecond int) *PreviewWaiter {
	timeout := time.Duration(timeoutSecond) * time.Second
	return &PreviewWaiter{
		serverContext: sc,
		backoff: &utilwait.Backoff{
			Duration: 200 * time.Millisecond,
			Factor:   2,
			Steps:    10,
			Cap:      timeout,
		},
		timeout: timeout,
	}
}

// PreviewWaiter is a type that represents preview object waiter.
type PreviewWaiter struct {
	serverContext *config.Config
	timeout       time.Duration
	backoff       *utilwait.Backoff
}

// Wait waits for the preview planned and print the error.
func (o *PreviewWaiter) Wait(ctx context.Context, set ObjectSet, _ chan ObjectSet) (bool, error) {
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
			err := o.waitPreview(obj)
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

// waitPreview waits for the object preview finished.
func (o *PreviewWaiter) waitPreview(obj Object) error {
	listOpt := api.OpenAPI.GetOperation(GroupResourceRuns, operationList)
	if listOpt == nil {
		return fmt.Errorf("not found %s operation for %s", operationList, GroupResourceRuns)
	}

	if o.backoff == nil {
		return o.getResourceRun(obj, listOpt)
	}

	var (
		err      error
		retryErr *common.RetryableError
	)
	err = utilwait.ExponentialBackoff(*o.backoff, func() (bool, error) {
		err := o.getResourceRun(obj, listOpt)
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

// getResourceRun get resource run and print changes.
func (o *PreviewWaiter) getResourceRun(obj Object, listOpt *api.Operation) error {
	runs, err := ListResourceRuns(o.serverContext, obj, listOpt, "")
	if err != nil {
		return err
	}

	run := runs.Items[0]
	switch run.Status.SummaryStatus {
	default:
		return fmt.Errorf("resource run %s is in unexpected status %s", run.ID, run.Status.SummaryStatus)
	case status.ResourceRunSummaryStatusPlanning:
		return common.NewRetryableError("preview is planning")
	case status.ResourceRunSummaryStatusPending:
		fmt.Printf("resource %s preview is pending\n", obj.ScopedName(obj.Name))
		return nil
	case status.ResourceRunSummaryStatusCanceled:
		fmt.Printf("resource %s preview is canceled\n", obj.ScopedName(obj.Name))
		return nil
	case status.ResourceRunSummaryStatusPlanned:
		return printResourceRunChanges(obj.ScopedName(obj.Name), "preview plan", run)
	}
}

// DefaultPreviewApplyWaiter wait for preview apply finished.
func DefaultPreviewApplyWaiter(sc *config.Config, timeoutSecond int) *PreviewApplyWaiter {
	timeout := time.Duration(timeoutSecond) * time.Second
	return &PreviewApplyWaiter{
		serverContext: sc,
		backoff: &utilwait.Backoff{
			Duration: 200 * time.Millisecond,
			Factor:   2,
			Steps:    10,
			Cap:      timeout,
		},
		timeout: timeout,
	}
}

// PreviewApplyWaiter is a type that represents preview apply waiter.
type PreviewApplyWaiter struct {
	serverContext *config.Config
	timeout       time.Duration
	backoff       *utilwait.Backoff
}

// Wait waits for the preview applied and print the error.
func (o *PreviewApplyWaiter) Wait(ctx context.Context, _ ObjectSet, setChan chan ObjectSet) (bool, error) {
	for {
		select {
		case <-ctx.Done():
			return false, fmt.Errorf("timeout")
		case resultSet := <-setChan:
			all := resultSet.All()
			if len(all) == 0 {
				return true, nil
			}

			for i := range all {
				err := o.waitPreviewApply(all[i])
				if err != nil {
					return false, err
				}
			}
			return true, nil
		}
	}
}

// waitPreviewApply waits for the preview apply finished.
func (o *PreviewApplyWaiter) waitPreviewApply(obj Object) error {
	getOpt := api.OpenAPI.GetOperation(GroupResourceRuns, operationGet)
	if getOpt == nil {
		return fmt.Errorf("not found %s operation for %s", operationGet, GroupResourceRuns)
	}

	runID := fmt.Sprintf("%v", obj.GetValue(FieldResourceRunID))

	if o.backoff == nil {
		return o.getResourceRun(obj, getOpt, runID)
	}

	var retryErr *common.RetryableError
	err := utilwait.ExponentialBackoff(*o.backoff, func() (bool, error) {
		err := o.getResourceRun(obj, getOpt, runID)
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

func (o *PreviewApplyWaiter) getResourceRun(obj Object, getOpt *api.Operation, runID string) error {
	run, err := GetResourceRun(o.serverContext, obj, getOpt, runID)
	if err != nil {
		return err
	}

	switch run.Status.SummaryStatus {
	default:
		return fmt.Errorf("resource %s run %s is in unexpected status %s: %s",
			obj.ScopedName(obj.Name), run.ID.String(), run.Status.SummaryStatus, run.Status.SummaryStatusMessage)
	case
		status.ResourceRunSummaryStatusPending,
		status.ResourceRunSummaryStatusRunning,
		status.ResourceRunSummaryStatusPlanning,
		status.ResourceRunSummaryStatusPlanned:

		return common.NewRetryableError(fmt.Sprintf("preview apply is %s", run.Status.SummaryStatus))
	case status.ResourceRunSummaryStatusCanceled:
		return fmt.Errorf("resource %s run %s is canceled", obj.ScopedName(obj.Name), run.ID)
	case status.ResourceRunSummaryStatusSucceed:
		return printResourceRunChanges(obj.ScopedName(obj.Name), "preview applied", run)
	}
}

func printResourceRunChanges(resName, action string, run *model.ResourceRunOutput) error {
	actionMsg := fmt.Sprintf("%s succeed", action)
	summary := fmt.Sprintf("resource %s %s, changes: %d created, %d updated, %d deleted.",
		resName, actionMsg, run.ComponentChangeSummary.Created, run.ComponentChangeSummary.Updated, run.ComponentChangeSummary.Deleted)

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

	fmt.Printf("%s\n%s\n\n", summary, t.ResourceItems(b))
	return nil
}
