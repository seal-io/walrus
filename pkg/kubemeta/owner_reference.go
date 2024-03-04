package kubemeta

import (
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/utils/ptr"
	ctrlcli "sigs.k8s.io/controller-runtime/pkg/client"
)

// MetaObject is the interface for the object with metadata.
type MetaObject = ctrlcli.Object

// IsControlledBy checks if the object has a controller owner reference set to the given owner.
func IsControlledBy(obj, owner MetaObject) bool {
	ref := GetControllerOfNoCopy(obj)
	if ref == nil {
		return false
	}
	return ref.UID == owner.GetUID()
}

// IsControlledByGVK checks if the object has a controller owner reference set to the given owner GVK.
func IsControlledByGVK(obj MetaObject, ownerGVK schema.GroupVersionKind) bool {
	ref := GetOwnerRefOfNoCopy(obj, ownerGVK)
	if ref == nil {
		return false
	}
	return ptr.Deref(ref.Controller, false)
}

// GetControllerOfNoCopy returns a pointer to the controller owner reference if obj has a controller.
func GetControllerOfNoCopy(obj MetaObject) *meta.OwnerReference {
	refs := obj.GetOwnerReferences()
	for i := range refs {
		if refs[i].Controller != nil && *refs[i].Controller {
			return &refs[i]
		}
	}
	return nil
}

// GetControllerOf returns a pointer to a copy of the controller owner reference if obj has a controller.
func GetControllerOf(obj MetaObject) *meta.OwnerReference {
	ref := GetControllerOfNoCopy(obj)
	if ref == nil {
		return nil
	}
	cp := *ref
	return &cp
}

// NewControllerRef creates an OwnerReference pointing to the given owner.
func NewControllerRef(owner MetaObject, gvk schema.GroupVersionKind) *meta.OwnerReference {
	return &meta.OwnerReference{
		APIVersion:         gvk.GroupVersion().String(),
		Kind:               gvk.Kind,
		Name:               owner.GetName(),
		UID:                owner.GetUID(),
		BlockOwnerDeletion: ptr.To(true),
		Controller:         ptr.To(true),
	}
}

// ControlOn sets the owner reference of obj to owner with the given GVK.
//
// ControlOn will append a new owner reference if obj does not have the given kind of owner,
// otherwise it will replace the existing owner reference.
func ControlOn(obj, owner MetaObject, gvk schema.GroupVersionKind) {
	var (
		gv  = gvk.GroupVersion().String()
		ors = obj.GetOwnerReferences()
		ref = NewControllerRef(owner, gvk)
	)

	idx := -1
	for i := range ors {
		if ors[i].APIVersion == gv && ors[i].Kind == gvk.Kind {
			idx = i
			break
		}
	}
	if idx == -1 {
		// Append.
		ors = append(ors, *ref)
	} else {
		// Replace.
		ors[idx] = *ref
	}
	obj.SetOwnerReferences(ors)
}

// ControlOff removes the owner reference of obj to owner with the given GVK.
//
// ControlOff will remove the owner reference if obj has the given kind of owner.
func ControlOff(obj, owner MetaObject, gvk schema.GroupVersionKind) {
	var (
		gv  = gvk.GroupVersion().String()
		ors = obj.GetOwnerReferences()
	)

	for i := range ors {
		if ors[i].APIVersion == gv && ors[i].Kind == gvk.Kind && ors[i].UID == owner.GetUID() {
			ors = append(ors[:i], ors[i+1:]...)
			break
		}
	}
	obj.SetOwnerReferences(ors)
}

// TryControlOn is similar to ControlOn,
// but it will not change the owner reference if obj has the given kind of owner.
func TryControlOn(obj, owner MetaObject, gvk schema.GroupVersionKind) (controlled bool) {
	var (
		gv  = gvk.GroupVersion().String()
		ors = obj.GetOwnerReferences()
		ref = NewControllerRef(owner, gvk)
	)

	for i := range ors {
		if ors[i].APIVersion == gv && ors[i].Kind == gvk.Kind {
			return false
		}
	}

	// Append.
	ors = append(ors, *ref)
	obj.SetOwnerReferences(ors)
	return true
}

// GetOwnerRefOfNoCopy returns a pointer to the owner reference if obj has the given kind of owner.
func GetOwnerRefOfNoCopy(obj MetaObject, gvk schema.GroupVersionKind) *meta.OwnerReference {
	var (
		gv  = gvk.GroupVersion().String()
		ors = obj.GetOwnerReferences()
	)
	for i := range ors {
		if ors[i].APIVersion == gv && ors[i].Kind == gvk.Kind {
			return &ors[i]
		}
	}
	return nil
}

// GetOwnerRefOf returns a pointer to a copy of the owner reference if obj has the given kind of owner.
func GetOwnerRefOf(obj MetaObject, gvk schema.GroupVersionKind) *meta.OwnerReference {
	ownerRef := GetOwnerRefOfNoCopy(obj, gvk)
	if ownerRef == nil {
		return nil
	}
	cp := *ownerRef
	return &cp
}

// GetOwnerRefsOfNoCopy returns a pointer to the owner references if obj has the given kind of owner.
func GetOwnerRefsOfNoCopy(obj MetaObject, gvk schema.GroupVersionKind) []*meta.OwnerReference {
	var (
		gv  = gvk.GroupVersion().String()
		ors = obj.GetOwnerReferences()
	)

	refs := make([]*meta.OwnerReference, 0, 4) // Usually there is only one owner reference.
	for i := range ors {
		if ors[i].APIVersion == gv && ors[i].Kind == gvk.Kind {
			refs = append(refs, &ors[i])
		}
	}
	return refs
}

// GetOwnerRefsOf returns a pointer to the owner references if obj has the given kind of owner.
func GetOwnerRefsOf(obj MetaObject, gvk schema.GroupVersionKind) []*meta.OwnerReference {
	ors := GetOwnerRefsOfNoCopy(obj, gvk)

	refs := make([]*meta.OwnerReference, 0, len(ors))
	for i := range ors {
		cp := *ors[i]
		refs = append(refs, &cp)
	}
	return refs
}
