package intercept

var TFAllTypes = []string{
	// Core.
	"kubernetes_namespace_v1", "kubernetes_namespace",
	"kubernetes_service_v1", "kubernetes_service",
	"kubernetes_service_account_v1", "kubernetes_service_account",
	"kubernetes_default_service_account_v1", "kubernetes_default_service_account",
	"kubernetes_config_map_v1", "kubernetes_config_map", "kubernetes_config_map_v1_data",
	"kubernetes_secret_v1", "kubernetes_secret",
	"kubernetes_pod_v1", "kubernetes_pod",
	"kubernetes_endpoints_v1", "kubernetes_endpoints",
	"kubernetes_limit_range_v1", "kubernetes_limit_range",
	"kubernetes_persistent_volume_v1", "kubernetes_persistent_volume",
	"kubernetes_persistent_volume_claim_v1", "kubernetes_persistent_volume_claim",
	"kubernetes_replication_controller_v1", "kubernetes_replication_controller",
	"kubernetes_resource_quota_v1", "kubernetes_resource_quota",

	// Api registration.
	"kubernetes_api_service_v1", "kubernetes_api_service",

	// Apps.
	"kubernetes_deployment_v1", "kubernetes_deployment",
	"kubernetes_daemon_set_v1", "kubernetes_daemonset", "kubernetes_daemon_set",
	"kubernetes_stateful_set_v1", "kubernetes_stateful_set",

	// Batch.
	"kubernetes_cron_job_v1",
	"kubernetes_cron_job",
	"kubernetes_job_v1", "kubernetes_job",

	// Autoscaling.
	"kubernetes_horizontal_pod_autoscaler_v2",
	"kubernetes_horizontal_pod_autoscaler_v2beta2",
	"kubernetes_horizontal_pod_autoscaler_v1", "kubernetes_horizontal_pod_autoscaler",

	// Certificates.
	"kubernetes_certificate_signing_request_v1",
	"kubernetes_certificate_signing_request",

	// Rbac.
	"kubernetes_role_v1", "kubernetes_role",
	"kubernetes_role_binding_v1", "kubernetes_role_binding",
	"kubernetes_cluster_role_v1", "kubernetes_cluster_role",
	"kubernetes_cluster_role_binding_v1", "kubernetes_cluster_role_binding",

	// Networking.
	"kubernetes_ingress_v1",
	"kubernetes_ingress",
	"kubernetes_ingress_class_v1", "kubernetes_ingress_class",
	"kubernetes_network_policy_v1", "kubernetes_network_policy",

	// Policy.
	"kubernetes_pod_disruption_budget_v1",
	"kubernetes_pod_disruption_budget",
	"kubernetes_pod_security_policy_v1beta1", "kubernetes_pod_security_policy",

	// Scheduling.
	"kubernetes_priority_class_v1", "kubernetes_priority_class",

	// Admission control.
	"kubernetes_validating_webhook_configuration_v1",
	"kubernetes_validating_webhook_configuration",
	"kubernetes_mutating_webhook_configuration_v1",
	"kubernetes_mutating_webhook_configuration",

	// Storage.
	"kubernetes_storage_class_v1", "kubernetes_storage_class",
	"kubernetes_csi_driver_v1", "kubernetes_csi_driver",
}

var TFLabeledTypes = []string{
	// Core.
	"kubernetes_pod_v1", "kubernetes_pod",
	"kubernetes_replication_controller_v1", "kubernetes_replication_controller",
	"kubernetes_persistent_volume_v1", "kubernetes_persistent_volume",
	"kubernetes_persistent_volume_claim_v1", "kubernetes_persistent_volume_claim",
	"kubernetes_service", "kubernetes_service_v1",

	// Apps.
	"kubernetes_deployment_v1", "kubernetes_deployment",
	"kubernetes_daemon_set_v1", "kubernetes_daemonset", "kubernetes_daemon_set",
	"kubernetes_stateful_set_v1", "kubernetes_stateful_set",

	// Batch.
	"kubernetes_cron_job_v1",
	"kubernetes_cron_job",
	"kubernetes_job_v1", "kubernetes_job",

	// Networking.
	"kubernetes_ingress", "kubernetes_ingress_v1",
}

var TFEndpointsTypes = []string{
	// Core.
	"kubernetes_service",
	"kubernetes_service_v1",

	// Networking.
	"kubernetes_ingress",
	"kubernetes_ingress_v1",

	// Kubectl_manifest resources.
	"kubectl_manifest",

	// Helm resources.
	"helm_release",
}
