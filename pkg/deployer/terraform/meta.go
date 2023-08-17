package terraform

import "k8s.io/apimachinery/pkg/util/sets"

// walrus metadata indicate walrus will set value to these attribute while user module include these attribute.
const (
	WalrusMetadataProjectName     = "walrus_metadata_project_name"
	WalrusMetadataEnvironmentName = "walrus_metadata_environment_name"
	WalrusMetadataServiceName     = "walrus_metadata_service_name"
	WalrusMetadataProjectID       = "walrus_metadata_project_id"
	WalrusMetadataEnvironmentID   = "walrus_metadata_environment_id"
	WalrusMetadataServiceID       = "walrus_metadata_service_id"
	// WalrusMetadataNamespaceName is the managed namespace name of an environment,
	// valid when Kubernetes connector is used in the environment.
	WalrusMetadataNamespaceName = "walrus_metadata_namespace_name"
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
