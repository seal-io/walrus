package alias

import (
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Terraform returns Convert to convert Terraform provider resource type to raw Kubernetes GVK.
func Terraform() Convert {
	return terraformConvert{}
}

// terraformConvert implements Convert,
// only consider operable types,
// from https://registry.terraform.io/providers/hashicorp/kubernetes/2.18.1.
type terraformConvert struct{}

func (terraformConvert) GetGVK(alias string) (gvk schema.GroupVersionKind, ok bool) {
	switch alias {
	case "kubernetes_daemon_set_v1", "kubernetes_daemonset", "kubernetes_daemon_set":
		return appsv1.SchemeGroupVersion.WithKind("DaemonSet"), true
	case "kubernetes_deployment_v1", "kubernetes_deployment":
		return appsv1.SchemeGroupVersion.WithKind("Deployment"), true
	case "kubernetes_stateful_set_v1", "kubernetes_stateful_set":
		return appsv1.SchemeGroupVersion.WithKind("StatefulSet"), true
	case "kubernetes_cron_job_v1":
		return batchv1.SchemeGroupVersion.WithKind("CronJob"), true
	case "kubernetes_cron_job":
		return batchv1beta1.SchemeGroupVersion.WithKind("CronJob"), true
	case "kubernetes_job_v1", "kubernetes_job":
		return batchv1.SchemeGroupVersion.WithKind("Job"), true
	case "kubernetes_replication_controller_v1", "kubernetes_replication_controller":
		return corev1.SchemeGroupVersion.WithKind("ReplicationController"), true
	case "kubernetes_pod_v1", "kubernetes_pod":
		return corev1.SchemeGroupVersion.WithKind("Pod"), true
	}
	return
}

func (terraformConvert) GetGVR(alias string) (gvr schema.GroupVersionResource, ok bool) {
	switch alias {
	case "kubernetes_daemon_set_v1", "kubernetes_daemonset", "kubernetes_daemon_set":
		return appsv1.SchemeGroupVersion.WithResource("daemonsets"), true
	case "kubernetes_deployment_v1", "kubernetes_deployment":
		return appsv1.SchemeGroupVersion.WithResource("deployments"), true
	case "kubernetes_stateful_set_v1", "kubernetes_stateful_set":
		return appsv1.SchemeGroupVersion.WithResource("statefulsets"), true
	case "kubernetes_cron_job_v1":
		return batchv1.SchemeGroupVersion.WithResource("cronjobs"), true
	case "kubernetes_cron_job":
		return batchv1beta1.SchemeGroupVersion.WithResource("cronjobs"), true
	case "kubernetes_job_v1", "kubernetes_job":
		return batchv1.SchemeGroupVersion.WithResource("jobs"), true
	case "kubernetes_replication_controller_v1", "kubernetes_replication_controller":
		return corev1.SchemeGroupVersion.WithResource("replicationcontrollers"), true
	case "kubernetes_pod_v1", "kubernetes_pod":
		return corev1.SchemeGroupVersion.WithResource("pods"), true
	}
	return
}
