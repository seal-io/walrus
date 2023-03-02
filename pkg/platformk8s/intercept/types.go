package intercept

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Convert holds the functions to transfer the given string to a schema descriptor.
type Convert interface {
	// GetGVK returns the GroupVersionKind info with the given alias,
	// and returns false if failed to convert.
	GetGVK(alias string) (gvk schema.GroupVersionKind, ok bool)

	// GetGVR returns the GroupVersionResource info with the given alias,
	// and returns false if failed to convert.
	GetGVR(alias string) (gvr schema.GroupVersionResource, ok bool)
}

// Enforce holds the functions to judge the given schema descriptor,
// whether to be interested in.
type Enforce interface {
	// AllowGVK returns true if the given GroupVersionKind is valid.
	AllowGVK(gvk schema.GroupVersionKind) (valid bool)

	// AllowGVR returns true if the given GroupVersionResource is valid.
	AllowGVR(gvr schema.GroupVersionResource) (valid bool)
}
