package json

import (
	jsonpatch "github.com/evanphx/json-patch"
)

func PatchObject(obj, patchObject any) (any, error) {
	ob, err := Marshal(obj)
	if err != nil {
		return nil, err
	}

	pb, err := Marshal(patchObject)
	if err != nil {
		return nil, err
	}

	patched, err := jsonpatch.MergePatch(ob, pb)
	if err != nil {
		return nil, err
	}

	return obj, Unmarshal(patched, obj)
}
