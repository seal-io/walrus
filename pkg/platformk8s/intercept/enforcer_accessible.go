package intercept

import (
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
)

// Accessible returns Enforcer to detect if the given Kubernetes GVK/GVR is accessible enforcer.
func Accessible() Enforcer {
	// Singleton pattern.
	return acEnforcer
}

// accessibleEnforcer implements Enforcer.
type accessibleEnforcer struct {
	gvks sets.Set[schema.GroupVersionKind]
	gvrs sets.Set[schema.GroupVersionResource]
}

func (e accessibleEnforcer) AllowGVK(gvk schema.GroupVersionKind) bool {
	return e.gvks.Has(gvk)
}

func (e accessibleEnforcer) AllowGVR(gvr schema.GroupVersionResource) bool {
	return e.gvrs.Has(gvr)
}

var acEnforcer = accessibleEnforcer{
	gvks: sets.Set[schema.GroupVersionKind]{},
	gvrs: sets.Set[schema.GroupVersionResource]{},
}

func init() {
	// Emit, transfer and record.
	//
	// Only consider accessible types.
	//
	for _, gvk := range []schema.GroupVersionKind{
		corev1.SchemeGroupVersion.WithKind("Service"),
		networkingv1.SchemeGroupVersion.WithKind("Ingress"),
		extensionsv1beta1.SchemeGroupVersion.WithKind("Ingress"),
	} {
		acEnforcer.gvks.Insert(gvk)
		gvr, _ := meta.UnsafeGuessKindToResource(gvk)
		acEnforcer.gvrs.Insert(gvr)
	}
}
