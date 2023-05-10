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

// Composite returns Enforcer to detect if the given Kubernetes GVK/GVR is composite enforcer.
func Composite() Enforcer {
	// Singleton pattern.
	return compEnforcer
}

// compositeEnforcer implements Enforcer.
type compositeEnforcer struct {
	gvks sets.Set[schema.GroupVersionKind]
	gvrs sets.Set[schema.GroupVersionResource]
}

func (e compositeEnforcer) AllowGVK(gvk schema.GroupVersionKind) bool {
	return e.gvks.Has(gvk)
}

func (e compositeEnforcer) AllowGVR(gvr schema.GroupVersionResource) bool {
	return e.gvrs.Has(gvr)
}

var compEnforcer = compositeEnforcer{
	gvks: sets.Set[schema.GroupVersionKind]{},
	gvrs: sets.Set[schema.GroupVersionResource]{},
}

func init() {
	// Emit, transfer and record.
	//
	// Only consider composite types.
	//
	for _, gvk := range []schema.GroupVersionKind{
		// Select generated job list by the cronjob label selector,
		// then select pod list by the job label selector.
		batchv1.SchemeGroupVersion.WithKind("CronJob"),
		batchv1beta1.SchemeGroupVersion.WithKind("CronJob"),

		// Select related persistent volume by the persistent volume claim.
		corev1.SchemeGroupVersion.WithKind("PersistentVolumeClaim"),

		// Select generated pod list by the following kinds' label selector.
		appsv1.SchemeGroupVersion.WithKind("DaemonSet"),
		appsv1.SchemeGroupVersion.WithKind("Deployment"),
		appsv1.SchemeGroupVersion.WithKind("StatefulSet"),
		appsv1.SchemeGroupVersion.WithKind("ReplicaSet"),
		batchv1.SchemeGroupVersion.WithKind("Job"),
		corev1.SchemeGroupVersion.WithKind("ReplicationController"),
	} {
		compEnforcer.gvks.Insert(gvk)
		gvr, _ := meta.UnsafeGuessKindToResource(gvk)
		compEnforcer.gvrs.Insert(gvr)
	}
}
