package openapi

import (
	"context"
	"fmt"
	"strings"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	"github.com/seal-io/walrus/pkg/dao/types/property"
	"github.com/seal-io/walrus/utils/gopool"
	"github.com/seal-io/walrus/utils/json"
)

type patchResult struct {
	err   error
	patch []byte
	group string
}

// GenSchemaDefaultPatch generates the default patch for the variable schema.
// The input root schema type should be object.
func GenSchemaDefaultPatch(ctx context.Context, schema *openapi3.Schema) ([]byte, error) {
	if schema == nil || schema.IsEmpty() {
		return nil, nil
	}

	// 1. Sort the nodes under the schema in a pre-order traversal.
	sortedNodes := PreorderTraversal(&DefaultPatchNode{
		Schema: schema,
	})

	// 2. Group nodes by the first part of the id to speed up,
	// after grouped, each group include the nodes from root to sub nodes in sequence.
	gns := groupNodes(sortedNodes)
	resultChan := make(chan patchResult, len(gns))

	// 3. Use goroutine to speed up apply patches from root to sub nodes.
	for gn, v := range gns {
		nodes := v
		name := gn

		gopool.Go(func() {
			switch {
			default:
				// 3.1 No node found.
				resultChan <- patchResult{
					group: name,
				}
			case len(nodes) == 1:
				// 3.1 Root node without sub nodes will set default and return.
				def := nodes[0].GetDefault()
				if def == nil {
					resultChan <- patchResult{
						group: name,
					}

					return
				}

				// 3.2 Generate patch for the node.
				p, err := sjson.SetBytes([]byte{}, nodes[0].GetID(), def)
				if err != nil {
					resultChan <- patchResult{
						group: name,
						err:   fmt.Errorf("set patch node error for %s %s: %w", nodes[0].GetID(), string(p), err),
					}

					return
				}

				resultChan <- patchResult{
					group: name,
					patch: p,
				}
			case len(nodes) > 1:
				// 3.1 Root node with sub nodes will apply patch in priority root > sub nodes.
				if nodes[0].GetDefault() == nil {
					resultChan <- patchResult{
						group: name,
						patch: nil,
					}

					return
				}

				// 3.2 Generate root patch.
				p, err := sjson.SetBytes([]byte{}, nodes[0].GetID(), nodes[0].GetDefault())
				if err != nil {
					resultChan <- patchResult{
						group: name,
						err:   fmt.Errorf("set patch node error for %s %s: %w", nodes[0].GetID(), string(p), err),
					}

					return
				}

				// 3.3 Generate patch from sub nodes and merge the root into the sub nodes.
				for i := 0; i < len(nodes); i++ {
					np := nodes[i]
					if np.GetDefault() == nil {
						continue
					}

					// Generate sub node patch.
					npb, err := sjson.SetBytes([]byte{}, nodes[i].GetID(), nodes[i].GetDefault())
					if err != nil {
						resultChan <- patchResult{
							group: name,
							err:   fmt.Errorf("set patch node error for %s %s: %w", nodes[0].GetID(), string(p), err),
						}

						return
					}

					var (
						pid = np.GetParentID()
						pr  = gjson.ParseBytes(p)
					)

					// 3.4 Only merge sub node's patch while ancestor already initialize it parent.
					if pe := pr.Get(pid); pe.Exists() {
						p, err = jsonpatch.MergePatch(npb, p)
						if err != nil {
							resultChan <- patchResult{
								group: name,
								err:   fmt.Errorf("merge patch error for %s %s: %w", nodes[i].GetID(), string(p), err),
							}

							return
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

	// 4. Merge all variables into dv.
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

// GenSchemaDefaultWithAttribute compute default values with attributes and exist default values in sequence.
func GenSchemaDefaultWithAttribute(
	ctx context.Context, schema *openapi3.Schema, attrs property.Values, defaultValuesByte ...[]byte,
) ([]byte, error) {
	copySchema := openapi3.NewObjectSchema()

	for n := range schema.Properties {
		if v := attrs[n]; v != nil &&
			schema.Properties[n] != nil &&
			schema.Properties[n].Value != nil {
			copyProp := *schema.Properties[n].Value
			copyProp.Default = v

			copySchema.Properties[n] = &openapi3.SchemaRef{
				Value: &copyProp,
			}
		}
	}

	// Generate default with attributes.
	dv, err := GenSchemaDefaultPatch(ctx, copySchema)
	if err != nil {
		return nil, err
	}

	// Merge the default from attributes and exist default.
	genAttrs := make(property.Values)

	if dv != nil {
		err = json.Unmarshal(dv, &genAttrs)
		if err != nil {
			return nil, err
		}
	}

	for i := range defaultValuesByte {
		var defaultValues property.Values

		err = json.Unmarshal(defaultValuesByte[i], &defaultValues)
		if err != nil {
			return nil, err
		}

		for k := range defaultValues {
			if _, ok := genAttrs[k]; !ok {
				genAttrs[k] = defaultValues[k]
			}
		}
	}

	return json.Marshal(genAttrs)
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
