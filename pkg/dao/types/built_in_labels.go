package types

const (
	// LabelPrefix is used for generate label's field names.
	LabelPrefix = "label:"
) // Built-in labels.
const (
	LabelWalrusProjectName     string = "walrus.seal.io/project-name"
	LabelWalrusEnvironmentName string = "walrus.seal.io/environment-name"
	LabelWalrusResourceName    string = "walrus.seal.io/resource-name"
	// LabelWalrusEnvironmentPath indicate environment with project name, format: projectName/environmentName.
	LabelWalrusEnvironmentPath string = "walrus.seal.io/environment-path"
	// LabelWalrusResourcePath indicate resource with project name and environment name,
	// format: projectName/environmentName/resourceName
	LabelWalrusResourcePath string = "walrus.seal.io/resource-path"
	// LabelWalrusManaged indicates whether the resource is managed by Seal.
	LabelWalrusManaged string = "walrus.seal.io/managed"
	// LabelWalrusCategory indicates the category of the resource.
	LabelWalrusCategory string = "walrus.seal.io/category"
	// LabelWalrusConnectorType indicates the connector type of the resource.
	LabelWalrusConnectorType string = "walrus.seal.io/connector-type"
	// LabelWalrusResourceType indicates the type of the resource.
	LabelWalrusResourceType string = "walrus.seal.io/resource-type"
	// LabelWalrusResourceDefinition indicates if the template is for resource definition.
	LabelWalrusResourceDefinition string = "walrus.seal.io/resource-definition"
	// LabelResourceStoppable indicates if the resource is stoppable.
	LabelResourceStoppable string = "walrus.seal.io/stoppable"
	// LabelEmbeddedKubernetes indicates whether a connector is the embedded kubernetes.
	LabelEmbeddedKubernetes = "walrus.seal.io/embedded-kubernetes"
	// LabelWalrusResourceRunID is the label for resource run id of the run job.
	LabelWalrusResourceRunID = "walrus.seal.io/resource-run-id"
	// LabelWalrusResourceRunTaskType is the label for resource run task type of the run job.
	LabelWalrusResourceRunTaskType = "walrus.seal.io/resource-run-task-type"
)
