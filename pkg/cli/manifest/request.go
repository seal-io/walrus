package manifest

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	jsonpatch "github.com/evanphx/json-patch"
	"go.uber.org/multierr"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/r3labs/sse"

	"github.com/seal-io/walrus/pkg/cli/api"
	"github.com/seal-io/walrus/pkg/cli/common"
	"github.com/seal-io/walrus/pkg/cli/config"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/utils/gopool"
	"github.com/seal-io/walrus/utils/json"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/pointer"
	"github.com/seal-io/walrus/utils/strs"
)

const (
	operationWatch       = "list"
	operationBatchCreate = "collection-create"
	operationBatchDelete = "collection-delete"
	operationGet         = "get"
	operationPatch       = "patch"
)

type ObjectStatus uint

const (
	statusNotFound ObjectStatus = iota
	statusUnchanged
	statusChanged
)

type result struct {
	err  error
	obj  Object
	objs []Object
}

// PatchObjects send patches objects request.
func PatchObjects(sc *config.Config, group string, objs ObjectByScope) (success, failed ObjectSet, err error) {
	if len(objs) == 0 {
		return success, failed, nil
	}

	patchOpt := api.OpenAPI.GetOperation(group, operationPatch)
	if patchOpt == nil {
		return success, failed, fmt.Errorf("not found patch operation for %s", group)
	}

	var (
		results = make(chan result, 1)
		all     = objs.All()
	)

	for _, obj := range all {
		o := obj

		gopool.Go(func() {
			csp := sc.ServerContext
			csp.Project = o.Project
			csp.Environment = o.Environment

			req, err := patchOpt.Request(nil, []string{o.Name}, o.ObjectScope.Map(), o.Value, csp)
			if err != nil {
				results <- result{
					err: err,
					obj: o,
				}

				return
			}

			resp, err := sc.DoRequest(req)
			if err != nil {
				results <- result{
					err: err,
					obj: o,
				}

				return
			}

			results <- result{
				err: common.CheckResponseStatus(resp, ""),
				obj: o,
			}
		})
	}

	for i := 0; i < len(all); i++ {
		r := <-results
		if r.err != nil {
			err = multierr.Append(err, r.err)
			failed.Add(r.obj)

			continue
		}

		success.Add(r.obj)
	}

	return success, failed, err
}

// BatchCreateObjects send batch create objects request.
func BatchCreateObjects(sc *config.Config, group string, objs ObjectByScope) (
	success, failed ObjectSet, err error,
) {
	if len(objs) == 0 {
		return
	}

	createOpt := api.OpenAPI.GetOperation(group, operationBatchCreate)
	if createOpt == nil {
		return success, failed, fmt.Errorf("not found batch create operation for %s", group)
	}

	results := make(chan result, 1)

	for scope, obj := range objs {
		o := obj
		s := scope

		gopool.Go(func() {
			err := batchCreateObjects(sc, createOpt, s, o)
			results <- result{
				err:  err,
				objs: o,
			}
		})
	}

	for i := 0; i < len(objs); i++ {
		r := <-results
		if r.err != nil {
			err = multierr.Append(err, r.err)
			failed.Add(r.objs...)

			continue
		}

		success.Add(r.objs...)
	}

	return success, failed, err
}

type collectionCreateInput struct {
	Items []any `json:"items"`
}

func newCollectionCreateInputs(objs []Object) collectionCreateInput {
	items := make([]any, 0, len(objs))
	for _, o := range objs {
		items = append(items, o.Value)
	}

	return collectionCreateInput{Items: items}
}

func batchCreateObjects(sc *config.Config, createOpt *api.Operation, scope ObjectScope, objs []Object) error {
	if len(objs) == 0 {
		return nil
	}

	csp := sc.ServerContext
	csp.Project = scope.Project
	csp.Environment = scope.Environment

	body := newCollectionCreateInputs(objs)

	req, err := createOpt.Request(nil, nil, scope.Map(), body, csp)
	if err != nil {
		return err
	}

	resp, err := sc.DoRequest(req)
	if err != nil {
		return err
	}

	return common.CheckResponseStatus(resp, "")
}

// GetObjects send get objects request.
func GetObjects(sc *config.Config, group string, objs ObjectByScope, detectChange bool) (
	unchanged ObjectSet,
	notFound ObjectSet,
	changed ObjectSet,
	err error,
) {
	if len(objs) == 0 {
		return unchanged, notFound, changed, nil
	}

	getOpt := api.OpenAPI.GetOperation(group, operationGet)
	if getOpt == nil {
		return unchanged, notFound, changed, fmt.Errorf("not found get operation for %s", group)
	}

	var (
		results = make(chan result, 1)
		all     = objs.All()
	)

	for _, obj := range all {
		o := obj

		gopool.Go(func() {
			csp := sc.ServerContext
			csp.Project = o.Project
			csp.Environment = o.Environment

			req, err := getOpt.Request(nil, []string{o.Name}, o.ObjectScope.Map(), o.Value, csp)
			if err != nil {
				results <- result{err: err, obj: o}
				return
			}

			resp, err := sc.DoRequest(req)
			if err != nil {
				results <- result{err: err, obj: o}
				return
			}
			defer resp.Body.Close()

			switch resp.StatusCode {
			case http.StatusTooManyRequests:
				results <- result{
					err: common.NewRetryableError("too many request"),
					obj: o,
				}
			case http.StatusNotFound:
				o.Status = statusNotFound
				results <- result{
					obj: o,
				}
			case http.StatusOK:
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					results <- result{
						err: fmt.Errorf("failed to read response from get %s %s", group, o.Name),
						obj: o,
					}

					return
				}

				var identify IDName

				err = json.Unmarshal(body, &identify)
				if err != nil {
					results <- result{
						err: fmt.Errorf("failed to unmarshal response from get %s %s", group, o.Name),
						obj: o,
					}

					return
				}

				o.ID = identify.ID
				o.Status = statusUnchanged

				if detectChange {
					pb, err := json.Marshal(o)
					if err != nil {
						results <- result{
							err: fmt.Errorf("failed to marshal object %s %s", group, o.Name),
							obj: o,
						}

						return
					}

					patched, err := jsonpatch.MergePatch(body, pb)
					if err != nil {
						results <- result{
							err: fmt.Errorf("failed to merge patch from get %s %s", group, o.Name),
							obj: o,
						}

						return
					}

					if !jsonpatch.Equal(body, patched) {
						o.Status = statusChanged
					}
				}

				results <- result{
					obj: o,
				}

			default:
				results <- result{
					err: fmt.Errorf("unexpected status code %d from get %s %s", resp.StatusCode, group, o.Name),
					obj: o,
				}
			}
		})
	}

	for i := 0; i < len(all); i++ {
		r := <-results
		if r.err != nil {
			return unchanged, notFound, changed, r.err
		}

		switch r.obj.Status {
		case statusUnchanged:
			unchanged.Add(r.obj)
		case statusNotFound:
			notFound.Add(r.obj)
		case statusChanged:
			changed.Add(r.obj)
		}
	}

	return unchanged, notFound, changed, nil
}

type deleteInput struct {
	Name string `json:"name"`
}

type collectionDeleteInput struct {
	Items []*deleteInput `json:"items"`
}

func newCollectionDeleteInput(objs []Object) collectionDeleteInput {
	input := collectionDeleteInput{}
	for _, v := range objs {
		input.Items = append(input.Items, &deleteInput{
			Name: v.Name,
		})
	}

	return input
}

// DeleteObjects send delete objects request.
func DeleteObjects(sc *config.Config, group string, objs ObjectByScope) (*ObjectSet, *ObjectSet, error) {
	batchDeleteOpt := api.OpenAPI.GetOperation(group, operationBatchDelete)
	if batchDeleteOpt == nil {
		return nil, nil, fmt.Errorf("not found batch delete operation for %s", group)
	}

	results := make(chan result, 1)

	for scope, obj := range objs {
		o := obj
		s := scope

		gopool.Go(func() {
			csp := sc.ServerContext
			csp.Project = s.Project
			csp.Environment = s.Environment
			body := newCollectionDeleteInput(o)

			req, err := batchDeleteOpt.Request(nil, nil, s.Map(), body, csp)
			if err != nil {
				results <- result{
					err:  err,
					objs: o,
				}

				return
			}

			resp, err := sc.DoRequest(req)
			if err != nil {
				results <- result{
					err:  err,
					objs: o,
				}

				return
			}

			results <- result{
				err:  common.CheckResponseStatus(resp, ""),
				objs: o,
			}
		})
	}

	var (
		success = &ObjectSet{}
		failed  = &ObjectSet{}
		err     error
	)

	for i := 0; i < len(objs); i++ {
		r := <-results
		if r.err != nil {
			err = multierr.Append(err, r.err)
			failed.Add(r.objs...)

			continue
		}

		success.Add(r.objs...)
	}

	return success, failed, err
}

// WatchObjects send watch objects request.
func WatchObjects(
	ctx context.Context,
	sc *config.Config,
	set ObjectSet,
	waitConds map[string]WaitFor,
	finishedChan chan ObjectSet,
) (bool, error) {
	wg := gopool.GroupWithContextIn(ctx)

	for group, condFunc := range waitConds {
		watchOpt := api.OpenAPI.GetOperation(group, operationWatch)
		if watchOpt == nil {
			return false, fmt.Errorf("not found list operation for %s", group)
		}

		for s, objs := range set.ByGroup(group) {
			var (
				names     = objs.Names()
				idNameMap = objs.IDNameMap()
			)

			wg.Go(func(ctx context.Context) error {
				var (
					errChan   = make(chan error)
					eventChan = make(chan Event)
					finished  = sets.Set[string]{}
				)

				gopool.Go(func() {
					m := s.Map()
					m["watch"] = pointer.String("true")

					req, err := watchOpt.Request(nil, nil, m, nil, sc.ServerContext)
					if err != nil {
						errChan <- err
					}

					sc.SetHeaders(req)

					err = sc.SetHost(req)
					if err != nil {
						errChan <- err
					}

					httpClient := sc.HttpClient()
					httpClient.Timeout = 0

					resp, err := httpClient.Do(req)
					if err != nil {
						errChan <- err
					}

					reader := sse.NewEventStreamReader(resp.Body)
					startReadLoop(ctx, reader, eventChan, errChan)
				})

				for {
					select {
					default:
					case <-ctx.Done():
						return fmt.Errorf("timeout")
					case err := <-errChan:
						return err
					case resultSet := <-finishedChan:
						finishedResult := resultSet.All()
						for _, fr := range finishedResult {
							if fr.Group != group {
								continue
							}

							if fr.Project != s.Project || fr.Environment != s.Environment {
								continue
							}

							if !names.Has(fr.Name) {
								continue
							}

							finished.Insert(fr.Name)
						}
					case e := <-eventChan:
						for i := range e.Items {
							// Set name while empty.
							ev := e.Items[i]
							if ev.Name == "" {
								ev.Name = idNameMap[ev.ID]
							}

							log.Debugf(
								"received %s event: %s %s %s \n",
								e.Type,
								strs.Singularize(group),
								s.ScopedName(ev.Name),
								ev.Status.SummaryStatus,
							)

							if !names.Has(ev.Name) {
								continue
							}

							if finished.Has(ev.Name) {
								continue
							}

							meet := condFunc(e.Type, ev, group, s.ScopedName(ev.Name))
							if meet {
								finished.Insert(ev.Name)
							}
						}
					}

					if names.Equal(finished) {
						return nil
					}
				}
			})
		}
	}

	err := wg.Wait()
	if err != nil {
		return false, err
	}

	return true, nil
}

func startReadLoop(ctx context.Context, reader *sse.EventStreamReader, eventChan chan Event, errChan chan error) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		b, err := reader.ReadEvent()
		if err != nil {
			if errors.Is(err, io.EOF) {
				errChan <- nil
				return
			}
			errChan <- err

			return
		}

		var event Event

		err = json.Unmarshal(b, &event)
		if err != nil {
			errChan <- err
			return
		}

		eventChan <- event
	}
}

type Event struct {
	Type  string      `json:"type"`
	Items []EventItem `json:"items"`
}

type EventItem struct {
	ID     string        `json:"id"`
	Name   string        `json:"name"`
	Status status.Status `json:"status"`
}
