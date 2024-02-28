package openapi

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spyzhov/ajson"
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
// The input root schema type should be object type.
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
			patch, err := applyPatches(nodes)
			resultChan <- patchResult{
				group: name,
				patch: patch,
				err:   err,
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

// applyPatches apply nodes' default to the patch.
func applyPatches(nodes []Node) (patch []byte, err error) {
	if len(nodes) == 0 || nodes[0].GetDefault() == nil {
		return
	}

	// 1. Root node.
	root := nodes[0]
	patch, err = sjson.SetBytes(patch, root.GetID(), root.GetDefault())
	if err != nil {
		return patch, fmt.Errorf("set patch node error for %s %s: %w", root.GetID(), string(patch), err)
	}

	if len(nodes) == 1 {
		return patch, nil
	}

	// 2. Child nodes.
	for _, node := range nodes[1:] {
		if node.GetDefault() == nil {
			continue
		}

		switch {
		case strings.Contains(node.GetID(), "*"):
			// 2.1 Apply patch with wildcard.
			patch, err = applyPatchWithWildcard(node, patch)
			if err != nil {
				return patch, fmt.Errorf("error apply patch with wildcard: %w", err)
			}

		default:
			// 2.2 Apply patch with nodeID.
			patch, err = applyPatch(node, patch)
			if err != nil {
				return patch, err
			}
		}
	}

	return patch, nil
}

// applyPatchWithWildcard apply the node's default to the patch, the id of the node could contain *.
func applyPatchWithWildcard(node Node, patch []byte) ([]byte, error) {
	var (
		nodeID      = node.GetID()
		parentID    = node.GetParentID()
		nodeDefault = node.GetDefault()
		nodeSchema  = node.GetSchema()
	)

	var (
		pathes []string
		err    error
	)
	switch {
	case strings.HasSuffix(nodeID, "*"):
		// Value in that nodeID need update.
		pathes, err = getTidwallFormatPathes(nodeID, "", patch)
		if err != nil {
			return nil, err
		}
	default:
		// Value in that nodeID need create/update.
		pathes, err = getTidwallFormatPathes(parentID, getLastName(nodeID), patch)
		if err != nil {
			return nil, err
		}
	}

	for _, id := range pathes {
		ajNode := &DefaultPatchNode{
			Schema: &nodeSchema,
			id:     id,
			def:    nodeDefault,
		}

		patch, err = applyPatch(ajNode, patch)
		if err != nil {
			return patch, err
		}
	}

	return patch, nil
}

// applyPatch apply the node's default to the patch, the id of the node must not contain *.
func applyPatch(node Node, patch []byte) ([]byte, error) {
	var (
		patchJson  = gjson.ParseBytes(patch)
		parentJson = patchJson.Get(node.GetParentID())
	)

	if !parentJson.Exists() {
		// Skip nodes without default or whose ancestors haven't been initialized.
		return patch, nil
	}

	var (
		nodeID    = node.GetID()
		nodeJson  = patchJson.Get(nodeID)
		nodePatch = node.GetDefault()
		err       error
	)

	if nodeJson.Exists() {
		if !nodeJson.IsObject() {
			return patch, nil
		}

		nodePatch, err = json.PatchObject(node.GetDefault(), nodeJson.Value())
		if err != nil {
			return patch, fmt.Errorf("error generate patch for object %s: %w", nodeID, err)
		}
	}

	// Set patch with nodeID.
	patch, err = sjson.SetBytes(patch, nodeID, nodePatch)
	if err != nil {
		return patch, fmt.Errorf("error set patch for %s: %w", nodeID, err)
	}

	return patch, nil
}

// jsonPathToTidwallFormat convert the jsonpath the patch that tidwall/sjson and tidwall/gjson used.
func jsonPathToTidwallFormat(jsonPath string) string {
	re := regexp.MustCompile(`\]\[|\[|\]`)

	jsonPath = strings.ReplaceAll(jsonPath, "'", "")
	jsonPath = re.ReplaceAllString(jsonPath, ".")
	jsonPath = strings.TrimSuffix(strings.TrimPrefix(jsonPath, "$."), ".")
	return jsonPath
}

func getTidwallFormatPathes(nodeID, suffix string, patch []byte) ([]string, error) {
	jsonPath := fmt.Sprintf("$.%s", nodeID)
	ajNodes, err := ajson.JSONPath(patch, jsonPath)
	if err != nil {
		return nil, fmt.Errorf("error get nodes from path %s: %w", jsonPath, err)
	}

	pathes := make([]string, len(ajNodes))
	for i := range ajNodes {
		if suffix != "" {
			pathes[i] = jsonPathToTidwallFormat(ajNodes[i].Path()) + "." + suffix
			continue
		}
		pathes[i] = jsonPathToTidwallFormat(ajNodes[i].Path())
	}
	return pathes, nil
}

// GenSchemaDefaultWithAttribute compute default values with attributes,
// exist default values arranged in ascending order of merging priority.
func GenSchemaDefaultWithAttribute(
	ctx context.Context, schema *openapi3.Schema, attrs property.Values, defaultValuesByte ...[]byte,
) ([]byte, error) {
	var (
		defaultWithAttrsByte []byte
		err                  error
	)

	if schema != nil {
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
		defaultWithAttrsByte, err = GenSchemaDefaultPatch(ctx, copySchema)
		if err != nil {
			return nil, err
		}
	}

	// No need to merge with others.
	if len(defaultValuesByte) == 0 {
		return defaultWithAttrsByte, nil
	}

	// Merge default values in sequence.
	var (
		merged            []byte
		patchesInSequence = append(defaultValuesByte, defaultWithAttrsByte)
	)

	for i := range patchesInSequence {
		if len(patchesInSequence[i]) == 0 {
			continue
		}

		if len(merged) == 0 {
			merged = patchesInSequence[i]
			continue
		}

		merged, err = jsonpatch.MergePatch(merged, patchesInSequence[i])
		if err != nil {
			return nil, err
		}
	}

	return merged, nil
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
