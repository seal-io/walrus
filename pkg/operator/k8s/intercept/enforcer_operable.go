package intercept

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
)

func init() {
	// Emit, transfer and record.
	//
	// Only consider operable types.
	//
	for _, gvk := range []schema.GroupVersionKind{
		// Select pod directly.
		corev1.SchemeGroupVersion.WithKind("Pod"),
	} {
		opEnforcer.gvks.Insert(gvk)
		gvr, _ := meta.UnsafeGuessKindToResource(gvk)
		opEnforcer.gvrs.Insert(gvr)
	}
}

// Operable returns Enforcer to detect if the given Kubernetes GVK/GVR is operable enforcer.
func Operable() Enforcer {
	// Singleton pattern.
	return opEnforcer
}

// operableEnforcer implements Enforcer.
type operableEnforcer struct {
	gvks sets.Set[schema.GroupVersionKind]
	gvrs sets.Set[schema.GroupVersionResource]
}

func (e operableEnforcer) AllowGVK(gvk schema.GroupVersionKind) bool {
	return e.gvks.Has(gvk)
}

func (e operableEnforcer) AllowGVR(gvr schema.GroupVersionResource) bool {
	return e.gvrs.Has(gvr)
}

var opEnforcer = operableEnforcer{
	gvks: sets.Set[schema.GroupVersionKind]{},
	gvrs: sets.Set[schema.GroupVersionResource]{},
}
