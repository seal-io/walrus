package intercept

import (
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Terraform returns Convert to convert Terraform provider resource type to raw Kubernetes GVK/GVR.
func Terraform() Convert {
	// singleton pattern.
	return tc
}

type terraformConvert struct {
	gvkm map[string]schema.GroupVersionKind
	gvrm map[string]schema.GroupVersionResource
}

func (c terraformConvert) GetGVK(alias string) (gvk schema.GroupVersionKind, ok bool) {
	gvk, ok = c.gvkm[alias]
	return
}

func (c terraformConvert) GetGVR(alias string) (gvr schema.GroupVersionResource, ok bool) {
	gvr, ok = c.gvrm[alias]
	return
}

var tc = terraformConvert{
	gvkm: map[string]schema.GroupVersionKind{},
	gvrm: map[string]schema.GroupVersionResource{},
}

func init() {
	// emit, transfer and record.
	//
	// only consider operable types,
	// from https://registry.terraform.io/providers/hashicorp/kubernetes/2.18.1.
	//
	for _, alias := range []string{
		"kubernetes_daemon_set_v1", "kubernetes_daemonset", "kubernetes_daemon_set",
		"kubernetes_deployment_v1", "kubernetes_deployment",
		"kubernetes_stateful_set_v1", "kubernetes_stateful_set",
		"kubernetes_stateful_set_v1", "kubernetes_stateful_set",
		"kubernetes_cron_job_v1",
		"kubernetes_cron_job",
		"kubernetes_job_v1", "kubernetes_job",
		"kubernetes_replication_controller_v1", "kubernetes_replication_controller",
		"kubernetes_pod_v1", "kubernetes_pod",
	} {
		var gvk = func() schema.GroupVersionKind {
			switch alias {
			case "kubernetes_daemon_set_v1", "kubernetes_daemonset", "kubernetes_daemon_set":
				return appsv1.SchemeGroupVersion.WithKind("DaemonSet")
			case "kubernetes_deployment_v1", "kubernetes_deployment":
				return appsv1.SchemeGroupVersion.WithKind("Deployment")
			case "kubernetes_stateful_set_v1", "kubernetes_stateful_set":
				return appsv1.SchemeGroupVersion.WithKind("StatefulSet")
			case "kubernetes_cron_job_v1":
				return batchv1.SchemeGroupVersion.WithKind("CronJob")
			case "kubernetes_cron_job":
				return batchv1beta1.SchemeGroupVersion.WithKind("CronJob")
			case "kubernetes_job_v1", "kubernetes_job":
				return batchv1.SchemeGroupVersion.WithKind("Job")
			case "kubernetes_replication_controller_v1", "kubernetes_replication_controller":
				return corev1.SchemeGroupVersion.WithKind("ReplicationController")
			case "kubernetes_pod_v1", "kubernetes_pod":
				return corev1.SchemeGroupVersion.WithKind("Pod")
			}
			panic("it will never happen")
		}()
		tc.gvkm[alias] = gvk
		tc.gvrm[alias], _ = meta.UnsafeGuessKindToResource(gvk)
	}
}
