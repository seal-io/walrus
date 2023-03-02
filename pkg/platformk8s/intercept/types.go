package intercept

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type Convert interface {
	// GetGVK returns the GroupVersionKind info with the given alias.
	GetGVK(string) (schema.GroupVersionKind, bool)

	// GetGVR returns the GroupVersionResource info with the given alias.
	GetGVR(string) (schema.GroupVersionResource, bool)
}

type Enforce interface {
	// AllowGVK returns true if the given GroupVersionKind is valid.
	AllowGVK(schema.GroupVersionKind) bool

	// AllowGVR returns true if the given GroupVersionResource is valid.
	AllowGVR(schema.GroupVersionResource) bool
}
