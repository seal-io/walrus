package intercept

import (
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
)

// Operable returns Enforcer to detect if the given Kubernetes GVK/GVR is operable enforcer.
func Operable() Enforcer {
	// singleton pattern.
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

func init() {
	// emit, transfer and record.
	//
	// only consider operableEnforcer types.
	//
	for _, gvk := range []schema.GroupVersionKind{
		appsv1.SchemeGroupVersion.WithKind("DaemonSet"),
		appsv1.SchemeGroupVersion.WithKind("Deployment"),
		appsv1.SchemeGroupVersion.WithKind("StatefulSet"),
		batchv1.SchemeGroupVersion.WithKind("CronJob"),
		batchv1beta1.SchemeGroupVersion.WithKind("CronJob"),
		batchv1.SchemeGroupVersion.WithKind("Job"),
		corev1.SchemeGroupVersion.WithKind("ReplicationController"),
		corev1.SchemeGroupVersion.WithKind("Pod"),
	} {
		opEnforcer.gvks.Insert(gvk)
		var gvr, _ = meta.UnsafeGuessKindToResource(gvk)
		opEnforcer.gvrs.Insert(gvr)
	}
}
