package systemmeta

import (
	"strings"

	"k8s.io/apimachinery/pkg/labels"
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

// LabelSelectorOf returns a label selector for a resources of the specified type.
//
// The returned selector should use to select the resources not include in github.com/seal-io/walrus/pkg/apis.
func LabelSelectorOf(resType string) labels.Selector {
	return labels.SelectorFromSet(labels.Set{
		ResourceTypeLabel: resType,
	})
}
