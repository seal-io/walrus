package json

import (
	"reflect"

	jsonpatch "github.com/evanphx/json-patch"
)

func PatchObject(obj, patchObject any) (any, error) {
	oc := reflect.New(reflect.TypeOf(obj)).Interface()

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

	err = Unmarshal(patched, oc)
	if err != nil {
		return nil, err
	}

	return oc, nil
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

// CreateMergePatch creates a merge patch from the original and modifies.
// The merge patch returned follows the specification defined at
// http://tools.ietf.org/html/draft-ietf-appsawg-json-merge-patch-07.
func CreateMergePatch(original, modified any) (any, error) {
	mp := reflect.New(reflect.TypeOf(original)).Interface()

	originalJSON, err := Marshal(original)
	if err != nil {
		return nil, err
	}

	modifiedJSON, err := Marshal(modified)
	if err != nil {
		return nil, err
	}

	mergePatch, err := jsonpatch.CreateMergePatch(originalJSON, modifiedJSON)
	if err != nil {
		return nil, err
	}

	err = Unmarshal(mergePatch, mp)
	if err != nil {
		return nil, err
	}

	return mp, nil
}
