package openapi

import (
	"context"
	"fmt"
	"strings"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/getkin/kin-openapi/openapi3"

	"github.com/seal-io/walrus/utils/gopool"
)

type patchResult struct {
	err   error
	patch []byte
	group string
}

// GenSchemaDefaultPatch generates the default patch for the schema.
func GenSchemaDefaultPatch(ctx context.Context, schema *openapi3.Schema) ([]byte, error) {
	if schema == nil || schema.IsEmpty() {
		return nil, nil
	}

	sortedNodes := PreorderTraversal(&DefaultPatchNode{
		Schema: schema,
	})

	// Group nodes by the first part of the id to speed up.
	gns := groupNodes(sortedNodes)
	resultChan := make(chan patchResult, len(gns))

	for gn, v := range gns {
		nodes := v
		name := gn

		gopool.Go(func() {
			switch {
			default:
				resultChan <- patchResult{
					group: name,
				}
			case len(nodes) == 1:
				np := nodes[0].GetPatch()
				if len(np) == 0 {
					resultChan <- patchResult{
						group: name,
					}

					return
				}

				resultChan <- patchResult{
					group: name,
					patch: nodes[0].GetPatch(),
				}
			case len(nodes) > 1:
				var (
					p   = nodes[0].GetPatch()
					err error
				)

				for i := len(nodes) - 1; i >= 0; i-- {
					np := nodes[i].GetPatch()
					if len(np) == 0 {
						continue
					}

					if len(p) == 0 {
						p = np
						continue
					}

					p, err = jsonpatch.MergePatch(p, np)
					if err != nil {
						resultChan <- patchResult{
							group: name,
							err:   fmt.Errorf("merge patch error for %s %s: %w", nodes[i].GetID(), string(p), err),
						}
					}
				}
				resultChan <- patchResult{
					group: name,
					patch: p,
				}
			}
		})
	}

	var (
		dv    []byte
		count int
		err   error
	)

	for {
		select {
		default:
			if count == len(gns) {
				return dv, nil
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		case r := <-resultChan:
			if r.err != nil {
				return nil, r.err
			}

			if len(dv) == 0 {
				dv = r.patch
			} else if len(r.patch) > 0 {
				dv, err = jsonpatch.MergePatch(dv, r.patch)
				if err != nil {
					return nil, fmt.Errorf("merge patch error for %s %s: %w", r.group, string(dv), err)
				}
			}

			count++
		}
	}
}

func groupNodes(nodes []Node) map[string][]Node {
	group := make(map[string][]Node)

	for i := range nodes {
		k := strings.Split(nodes[i].GetID(), ".")
		if len(k) == 0 {
			continue
		}
		g := k[0]
		group[g] = append(group[g], nodes[i])
	}

	return group
}
