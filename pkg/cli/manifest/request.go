package manifest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/seal-io/walrus/pkg/cli/api"
	"github.com/seal-io/walrus/pkg/cli/config"
	"github.com/seal-io/walrus/pkg/cli/formatter"
	"github.com/seal-io/walrus/utils/gopool"
)

const (
	operationBatchCreate = "batch-create"
	operationBatchDelete = "batch-delete"
	operationGet         = "get"
	operationPatch       = "patch"
)

func PatchObjects(sc *config.Config, group string, objs []Object) error {
	if len(objs) == 0 {
		return nil
	}

	patchOpt := api.OpenAPI.GetOperation(group, operationPatch)
	if patchOpt == nil {
		return fmt.Errorf("not found patch operation for %s", group)
	}

	wg := gopool.GroupWithContextIn(context.Background())

	for _, obj := range objs {
		o := obj

		wg.Go(func(ctx context.Context) error {
			csp := sc.ServerContext
			csp.Project = o.Context.Project
			csp.Environment = o.Context.Environment

			req, err := patchOpt.Request(nil, []string{o.Name}, o.Context.FlagsData(), o.Value, csp)
			if err != nil {
				return err
			}

			resp, err := sc.DoRequest(req)
			if err != nil {
				return err
			}

			_, err = formatter.Format(sc.Format, resp)

			return err
		})
	}

	return wg.Wait()
}

func BatchCreateObjects(sc *config.Config, group string, objs map[string][]Object) error {
	if len(objs) == 0 {
		return nil
	}

	createOpt := api.OpenAPI.GetOperation(group, operationBatchCreate)
	if createOpt == nil {
		return fmt.Errorf("not found batch create operation for %s", group)
	}

	wg := gopool.GroupWithContextIn(context.Background())

	for _, obj := range objs {
		o := obj

		wg.Go(func(ctx context.Context) error {
			return batchCreateObjects(sc, o, createOpt)
		})
	}

	return wg.Wait()
}

type collectionCreateInputs struct {
	Items []any `json:"items"`
}

func newCollectionCreateInputs(objs []Object) collectionCreateInputs {
	items := make([]any, 0, len(objs))
	for _, o := range objs {
		items = append(items, o.Value)
	}

	return collectionCreateInputs{Items: items}
}

func batchCreateObjects(sc *config.Config, objs []Object, createOpt *api.Operation) error {
	if len(objs) == 0 {
		return nil
	}

	csp := sc.ServerContext
	csp.Project = objs[0].Context.Project
	csp.Environment = objs[0].Context.Environment

	body := newCollectionCreateInputs(objs)

	req, err := createOpt.Request(nil, nil, objs[0].Context.FlagsData(), body, csp)
	if err != nil {
		return err
	}

	resp, err := sc.DoRequest(req)
	if err != nil {
		return err
	}

	_, err = formatter.Format(sc.Format, resp)

	return err
}

func GetObjects(sc *config.Config, group string, objs []Object) (
	[]Object,
	map[string][]Object,
	error,
) {
	if len(objs) == 0 {
		return nil, nil, nil
	}

	getOpt := api.OpenAPI.GetOperation(group, operationGet)
	if getOpt == nil {
		return nil, nil, fmt.Errorf("not found get operation for %s", group)
	}

	var (
		toCreate = make(map[string][]Object)
		toPatch  = make([]Object, 0)
	)

	type result struct {
		err       error
		operation string
		scope     string
		obj       Object
	}

	results := make(chan result, 1)

	for _, obj := range objs {
		o := obj

		gopool.Go(func() {
			csp := sc.ServerContext
			csp.Project = o.Context.Project
			csp.Environment = o.Context.Environment

			req, err := getOpt.Request(nil, []string{o.Name}, o.Context.FlagsData(), o.Value, csp)
			if err != nil {
				results <- result{err: err, obj: o}
				return
			}

			resp, err := sc.DoRequest(req)
			if err != nil {
				results <- result{err: err, obj: o}
				return
			}

			switch resp.StatusCode {
			case http.StatusNotFound:
				results <- result{operation: operationBatchCreate, scope: o.Scope, obj: o}
			case http.StatusOK:
				results <- result{operation: operationPatch, scope: o.Scope, obj: o}
			default:
				results <- result{
					err: fmt.Errorf("unexpected status code %d from get %s %s", resp.StatusCode, group, o.Name),
					obj: o,
				}
			}
		})
	}

	for i := 0; i < len(objs); i++ {
		r := <-results
		if r.err != nil {
			return nil, nil, r.err
		}

		switch r.operation {
		case operationPatch:
			toPatch = append(toPatch, r.obj)
		case operationBatchCreate:
			toCreate[r.scope] = append(toCreate[r.scope], r.obj)
		default:
			return nil, nil, fmt.Errorf("unknown operation %s", r.operation)
		}
	}

	return toPatch, toCreate, nil
}

type collectionDeleteInputs struct {
	Items []*DeleteInput `json:"items"`
}

type DeleteInput struct {
	Name string `json:"name"`
}

func DeleteObjects(sc *config.Config, group string, objs []Object) error {
	batchDeleteOpt := api.OpenAPI.GetOperation(group, operationBatchDelete)
	if batchDeleteOpt == nil {
		return fmt.Errorf("not found batch delete operation for %s", group)
	}

	scopedObjs := make(map[config.ScopeContext][]*DeleteInput)
	for _, o := range objs {
		scopedObjs[o.Context] = append(scopedObjs[o.Context], &DeleteInput{Name: o.Name})
	}

	for octx, o := range scopedObjs {
		csp := sc.ServerContext
		csp.Project = octx.Project
		csp.Environment = octx.Environment

		body := collectionDeleteInputs{Items: o}

		req, err := batchDeleteOpt.Request(nil, nil, octx.FlagsData(), body, csp)
		if err != nil {
			return err
		}

		resp, err := sc.DoRequest(req)
		if err != nil {
			return err
		}

		_, err = formatter.Format(sc.Format, resp)
		if err != nil {
			return err
		}
	}

	return nil
}
