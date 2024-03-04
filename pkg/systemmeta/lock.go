package systemmeta

import (
	"slices"
)

const (
	// LockedResourceFinalizer is the finalizer to indicate the resource is locked by system.
	LockedResourceFinalizer = "walrus.seal.io/controlled"
)

// Lock adds a finalizer to the given resource.
//
// If the resource has been controlled, returns true,
// otherwise, returns false.
func Lock(obj MetaObject) (locked bool) {
	if obj == nil {
		panic("object is nil")
	}

	fs := obj.GetFinalizers()
	if slices.Contains(fs, LockedResourceFinalizer) {
		return true
	}

	fs = append(fs, LockedResourceFinalizer)
	obj.SetFinalizers(fs)
	return false
}

// Unlock removes a finalizer from the given resource.
//
// If the resource is not be controlled, returns true,
// otherwise, returns false.
func Unlock(obj MetaObject) (unlocked bool) {
	if obj == nil {
		panic("object is nil")
	}

	fs := obj.GetFinalizers()
	fs2 := slices.DeleteFunc(fs, func(s string) bool {
		return s == LockedResourceFinalizer
	})
	if len(fs) == len(fs2) {
		return true
	}
	obj.SetFinalizers(fs2)
	return false
}
