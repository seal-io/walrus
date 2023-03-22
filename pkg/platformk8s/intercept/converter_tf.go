package intercept

import (
	"fmt"

	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	certificatesv1 "k8s.io/api/certificates/v1"
	certificatesv1beta1 "k8s.io/api/certificates/v1beta1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1 "k8s.io/api/policy/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	schedulingv1 "k8s.io/api/scheduling/v1"
	storagev1 "k8s.io/api/storage/v1"
	storagev1beta1 "k8s.io/api/storage/v1beta1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
)

// Terraform returns Converter to convert Terraform provider resource type to raw Kubernetes GVK/GVR.
func Terraform() Converter {
	// singleton pattern.
	return tfConvert
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

var tfConvert = terraformConvert{
	gvkm: map[string]schema.GroupVersionKind{},
	gvrm: map[string]schema.GroupVersionResource{},
}

func init() {
	// emit, transfer and record.
	//
	// ref to https://registry.terraform.io/providers/hashicorp/kubernetes/2.18.1.
	//
	for _, alias := range TFAllTypes {
		var gvk = func() schema.GroupVersionKind {
			switch alias {
			case "kubernetes_namespace_v1", "kubernetes_namespace":
				return corev1.SchemeGroupVersion.WithKind("Namespace")
			case "kubernetes_service_v1", "kubernetes_service":
				return corev1.SchemeGroupVersion.WithKind("Service")
			case "kubernetes_service_account_v1", "kubernetes_service_account", "kubernetes_default_service_account_v1", "kubernetes_default_service_account":
				return corev1.SchemeGroupVersion.WithKind("ServiceAccount")
			case "kubernetes_config_map_v1", "kubernetes_config_map", "kubernetes_config_map_v1_data":
				return corev1.SchemeGroupVersion.WithKind("ConfigMap")
			case "kubernetes_secret_v1", "kubernetes_secret":
				return corev1.SchemeGroupVersion.WithKind("Secret")
			case "kubernetes_pod_v1", "kubernetes_pod":
				return corev1.SchemeGroupVersion.WithKind("Pod")
			case "kubernetes_endpoints_v1", "kubernetes_endpoints":
				return corev1.SchemeGroupVersion.WithKind("Endpoints")
			case "kubernetes_limit_range_v1", "kubernetes_limit_range":
				return corev1.SchemeGroupVersion.WithKind("LimitRange")
			case "kubernetes_persistent_volume_v1", "kubernetes_persistent_volume":
				return corev1.SchemeGroupVersion.WithKind("PersistentVolume")
			case "kubernetes_persistent_volume_claim_v1", "kubernetes_persistent_volume_claim":
				return corev1.SchemeGroupVersion.WithKind("PersistentVolumeClaim")
			case "kubernetes_replication_controller_v1", "kubernetes_replication_controller":
				return corev1.SchemeGroupVersion.WithKind("ReplicationController")
			case "kubernetes_resource_quota_v1", "kubernetes_resource_quota":
				return corev1.SchemeGroupVersion.WithKind("ResourceQuota")

			case "kubernetes_api_service_v1", "kubernetes_api_service":
				return apiregistrationv1.SchemeGroupVersion.WithKind("APIService")

			case "kubernetes_deployment_v1", "kubernetes_deployment":
				return appsv1.SchemeGroupVersion.WithKind("Deployment")
			case "kubernetes_daemon_set_v1", "kubernetes_daemonset", "kubernetes_daemon_set":
				return appsv1.SchemeGroupVersion.WithKind("DaemonSet")
			case "kubernetes_stateful_set_v1", "kubernetes_stateful_set":
				return appsv1.SchemeGroupVersion.WithKind("StatefulSet")

			case "kubernetes_cron_job_v1":
				return batchv1.SchemeGroupVersion.WithKind("CronJob")
			case "kubernetes_cron_job":
				return batchv1beta1.SchemeGroupVersion.WithKind("CronJob")
			case "kubernetes_job_v1", "kubernetes_job":
				return batchv1.SchemeGroupVersion.WithKind("Job")

			case "kubernetes_horizontal_pod_autoscaler_v2":
				return autoscalingv2.SchemeGroupVersion.WithKind("HorizontalPodAutoscaler")
			case "kubernetes_horizontal_pod_autoscaler_v2beta2":
				return autoscalingv2beta2.SchemeGroupVersion.WithKind("HorizontalPodAutoscaler")
			case "kubernetes_horizontal_pod_autoscaler_v1", "kubernetes_horizontal_pod_autoscaler":
				return autoscalingv1.SchemeGroupVersion.WithKind("HorizontalPodAutoscaler")

			case "kubernetes_certificate_signing_request_v1":
				return certificatesv1.SchemeGroupVersion.WithKind("CertificateSigningRequest")
			case "kubernetes_certificate_signing_request":
				return certificatesv1beta1.SchemeGroupVersion.WithKind("CertificateSigningRequest")

			case "kubernetes_role_v1", "kubernetes_role":
				return rbacv1.SchemeGroupVersion.WithKind("Role")
			case "kubernetes_role_binding_v1", "kubernetes_role_binding":
				return rbacv1.SchemeGroupVersion.WithKind("RoleBinding")
			case "kubernetes_cluster_role_v1", "kubernetes_cluster_role":
				return rbacv1.SchemeGroupVersion.WithKind("ClusterRole")
			case "kubernetes_cluster_role_binding_v1", "kubernetes_cluster_role_binding":
				return rbacv1.SchemeGroupVersion.WithKind("ClusterRoleBinding")

			case "kubernetes_ingress_v1":
				return networkingv1.SchemeGroupVersion.WithKind("Ingress")
			case "kubernetes_ingress":
				return extensionsv1beta1.SchemeGroupVersion.WithKind("Ingress")
			case "kubernetes_ingress_class_v1", "kubernetes_ingress_class":
				return networkingv1.SchemeGroupVersion.WithKind("IngressClass")
			case "kubernetes_network_policy_v1", "kubernetes_network_policy":
				return networkingv1.SchemeGroupVersion.WithKind("NetworkPolicy")

			case "kubernetes_pod_disruption_budget_v1":
				return policyv1.SchemeGroupVersion.WithKind("PodDisruptionBudget")
			case "kubernetes_pod_disruption_budget":
				return policyv1beta1.SchemeGroupVersion.WithKind("PodDisruptionBudget")
			case "kubernetes_pod_security_policy_v1beta1", "kubernetes_pod_security_policy":
				return policyv1beta1.SchemeGroupVersion.WithKind("PodSecurityPolicy")

			case "kubernetes_priority_class_v1", "kubernetes_priority_class":
				return schedulingv1.SchemeGroupVersion.WithKind("PriorityClass")

			case "kubernetes_validating_webhook_configuration_v1":
				return admissionregistrationv1.SchemeGroupVersion.WithKind("ValidatingWebhookConfiguration")
			case "kubernetes_validating_webhook_configuration":
				return admissionregistrationv1beta1.SchemeGroupVersion.WithKind("ValidatingWebhookConfiguration")
			case "kubernetes_mutating_webhook_configuration_v1":
				return admissionregistrationv1.SchemeGroupVersion.WithKind("MutatingWebhookConfiguration")
			case "kubernetes_mutating_webhook_configuration":
				return admissionregistrationv1beta1.SchemeGroupVersion.WithKind("MutatingWebhookConfiguration")

			case "kubernetes_storage_class_v1", "kubernetes_storage_class":
				return storagev1.SchemeGroupVersion.WithKind("StorageClass")
			case "kubernetes_csi_driver_v1":
				return storagev1.SchemeGroupVersion.WithKind("CSIDriver")
			case "kubernetes_csi_driver":
				return storagev1beta1.SchemeGroupVersion.WithKind("CSIDriver")
			}

			panic(fmt.Sprintf("needs transferring new alias %s to Kubernetes GVK", alias))
		}()
		tfConvert.gvkm[alias] = gvk
		tfConvert.gvrm[alias], _ = meta.UnsafeGuessKindToResource(gvk)
	}
}
