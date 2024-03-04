package systemmeta

import (
	"fmt"
	"strings"

	kmeta "k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
)

const (
	// ResourceTypeLabel is the label key to indicate the system type of delegated resource.
	ResourceTypeLabel = "resource.walrus.seal.io/type"

	// ResourceNoteAnnoPrefix is the annotation prefix to indicate the note of resource.
	ResourceNoteAnnoPrefix = "note.walrus.seal.io/"
)

// NoteResource armors the resource.
func NoteResource(obj MetaObject, resType string, notes map[string]string) {
	if obj == nil {
		panic("object is nil")
	}

	ls := obj.GetLabels()
	if ls == nil {
		ls = make(map[string]string)
	}
	ls[ResourceTypeLabel] = resType
	obj.SetLabels(ls)

	if len(notes) > 0 {
		as := obj.GetAnnotations()
		if as == nil {
			as = make(map[string]string)
		}

		for k, v := range notes {
			as[ResourceNoteAnnoPrefix+k] = v
		}

		obj.SetAnnotations(as)
	}
}

// DescribeResource explains the details of resource.
func DescribeResource(obj MetaObject) (resType string, notes map[string]string) {
	if obj == nil {
		panic("object is nil")
	}

	ls := obj.GetLabels()
	if ls != nil {
		resType = ls[ResourceTypeLabel]
	}

	notes = make(map[string]string)
	as := obj.GetAnnotations()
	for k, v := range as {
		if strings.HasPrefix(k, ResourceNoteAnnoPrefix) {
			notes[k[len(ResourceNoteAnnoPrefix):]] = v
		}
	}
	return
}

// DescribeResourceType explains the resource type of resource.
func DescribeResourceType(obj MetaObject) string {
	if obj == nil {
		panic("object is nil")
	}

	ls := obj.GetLabels()
	if ls != nil {
		return ls[ResourceTypeLabel]
	}

	return ""
}

// DescribeResourceNote explains the note of a resource with the given key.
func DescribeResourceNote(obj MetaObject, noteKey string) string {
	if obj == nil {
		panic("object is nil")
	}

	as := obj.GetAnnotations()
	for k, v := range as {
		if !strings.HasPrefix(k, ResourceNoteAnnoPrefix) {
			continue
		}
		if k[len(ResourceNoteAnnoPrefix):] == noteKey {
			return v
		}
	}
	return ""
}

// DescribeResourceNotes explains the notes of a resource with the given keys.
func DescribeResourceNotes(obj MetaObject, noteKeys []string) map[string]string {
	if obj == nil {
		panic("object is nil")
	}

	nks := sets.New(noteKeys...)
	r := make(map[string]string, len(noteKeys))

	as := obj.GetAnnotations()
	for k, v := range as {
		if !strings.HasPrefix(k, ResourceNoteAnnoPrefix) {
			continue
		}
		if nks.Has(k[len(ResourceNoteAnnoPrefix):]) {
			r[k[len(ResourceNoteAnnoPrefix):]] = v
		}
	}

	return r
}

// UnnoteResource is similar to DescribeResource,
// but it removes the resource labels and annotations.
func UnnoteResource(obj MetaObject) (resType string, notes map[string]string) {
	if obj == nil {
		panic("object is nil")
	}

	ls := obj.GetLabels()
	if ls != nil {
		resType = ls[ResourceTypeLabel]
		delete(ls, ResourceTypeLabel)
		obj.SetLabels(ls)
	}

	notes = make(map[string]string)
	as := obj.GetAnnotations()
	if as != nil {
		for k := range as {
			if !strings.HasPrefix(k, ResourceNoteAnnoPrefix) {
				continue
			}
			notes[k[len(ResourceNoteAnnoPrefix):]] = as[k]
			delete(as, k)
		}
		obj.SetAnnotations(as)
	}

	return
}

// PopResourceNote is similar to DescribeResourceNote,
// but it removes the resource annotation with the given key.
func PopResourceNote(obj MetaObject, noteKey string) string {
	if obj == nil {
		panic("object is nil")
	}

	as := obj.GetAnnotations()
	for k := range as {
		if k == ResourceNoteAnnoPrefix+noteKey {
			v := as[k]
			delete(as, k)
			obj.SetAnnotations(as)
			return v
		}
	}
	return ""
}

// PopResourceNotes is similar to DescribeResourceNotes,
// but it removes the resource annotations with the given keys.
func PopResourceNotes(obj MetaObject, noteKeys []string) map[string]string {
	if obj == nil {
		panic("object is nil")
	}

	nks := sets.New(noteKeys...)
	r := make(map[string]string, len(noteKeys))

	as := obj.GetAnnotations()
	for k := range as {
		if !strings.HasPrefix(k, ResourceNoteAnnoPrefix) {
			continue
		}
		if nks.Has(k[len(ResourceNoteAnnoPrefix):]) {
			r[k[len(ResourceNoteAnnoPrefix):]] = as[k]
			delete(as, k)
		}
	}

	obj.SetAnnotations(as)
	return r
}

// GetResourcesLabelSelectorOfType returns a label selector for the resources of the specified type.
func GetResourcesLabelSelectorOfType(resType string) labels.Selector {
	return labels.SelectorFromSet(labels.Set{
		ResourceTypeLabel: resType,
	})
}

// FilterResourceListByNotes returns a list of resources that matched by notes.
func FilterResourceListByNotes[T MetaObjectList](objList T, noteKey, noteValue string, noteKeyAndValues ...string) T {
	items, err := kmeta.ExtractList(objList)
	if err != nil {
		panic(fmt.Errorf("failed to extract list: %w", err))
	}

	newItems := make([]runtime.Object, 0, len(items))
	for i := range items {
		m := MatchResourceByNotes(items[i].(MetaObject), noteKey, noteValue, noteKeyAndValues...)
		if !m {
			continue
		}
		newItems = append(newItems, items[i])
	}

	err = kmeta.SetList(objList, newItems)
	if err != nil {
		panic(fmt.Errorf("failed to set list: %w", err))
	}

	return objList
}

// MatchResourceByNotes returns true if the given resource matched by notes.
func MatchResourceByNotes(obj MetaObject, noteKey, noteValue string, noteKeyAndValues ...string) bool {
	if obj == nil {
		panic("object is nil")
	}

	if noteValue != DescribeResourceNote(obj, noteKey) {
		return false
	}
	for i := 0; i < len(noteKeyAndValues); i += 2 {
		k := noteKeyAndValues[i]
		v := ""
		if i+1 < len(noteKeyAndValues) {
			v = noteKeyAndValues[i+1]
		}
		if v != DescribeResourceNote(obj, k) {
			return false
		}
	}
	return true
}

// FilterResourceList returns a list of resources that matched by resource type and notes.
func FilterResourceList[T MetaObjectList](objList T, resType string, noteKeyAndValues ...string) T {
	items, err := kmeta.ExtractList(objList)
	if err != nil {
		panic(fmt.Errorf("failed to extract list: %w", err))
	}

	newItems := make([]runtime.Object, 0, len(items))
	for i := range items {
		m := MatchResource(items[i].(MetaObject), resType, noteKeyAndValues...)
		if !m {
			continue
		}
		newItems = append(newItems, items[i])
	}

	err = kmeta.SetList(objList, newItems)
	if err != nil {
		panic(fmt.Errorf("failed to set list: %w", err))
	}

	return objList
}

// MatchResource returns true if the given resource matched by resource type and notes.
func MatchResource(obj MetaObject, resType string, noteKeyAndValues ...string) bool {
	if obj == nil {
		panic("object is nil")
	}

	if DescribeResourceType(obj) != resType {
		return false
	}
	for i := 0; i < len(noteKeyAndValues); i += 2 {
		k := noteKeyAndValues[i]
		v := ""
		if i+1 < len(noteKeyAndValues) {
			v = noteKeyAndValues[i+1]
		}
		if v != DescribeResourceNote(obj, k) {
			return false
		}
	}
	return true
}
