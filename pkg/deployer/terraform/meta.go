package terraform

import "k8s.io/apimachinery/pkg/util/sets"

// walrus metadata indicate walrus will set value to these attribute while user module include these attribute.
const (
	/* TODO coordinate with templates and update the following keys.*/

	WalrusMetadataProjectName     = "seal_metadata_project_name"
	WalrusMetadataEnvironmentName = "seal_metadata_environment_name"
	WalrusMetadataServiceName     = "seal_metadata_service_name"
	WalrusMetadataProjectID       = "seal_metadata_project_id"
	WalrusMetadataEnvironmentID   = "seal_metadata_environment_id"
	WalrusMetadataServiceID       = "seal_metadata_service_id"
	// WalrusMetadataNamespaceName is the managed namespace name of an environment,
	// valid when Kubernetes connector is used in the environment.
	WalrusMetadataNamespaceName = "seal_metadata_namespace_name"
)

var WalrusMetadataSet = sets.NewString(
	WalrusMetadataProjectName,
	WalrusMetadataEnvironmentName,
	WalrusMetadataServiceName,
	WalrusMetadataProjectID,
	WalrusMetadataEnvironmentID,
	WalrusMetadataServiceID,
	WalrusMetadataNamespaceName,
)
