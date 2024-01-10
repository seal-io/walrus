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

	err = Unmarshal(patched, obj)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func ApplyPatches(doc []byte, patches ...[]byte) ([]byte, error) {
	var err error

	for i := range patches {
		if len(patches[i]) == 0 {
			continue
		}

		if len(doc) == 0 {
			doc = patches[i]
			continue
		}

		doc, err = jsonpatch.MergePatch(doc, patches[i])
		if err != nil {
			return nil, err
		}
	}

	return doc, nil
}
