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
