package types

const (
	// LabelPrefix is used for generate label's field names.
	LabelPrefix = "label:"
)

// built-in labels.
const (
	LabelWalrusProjectName     string = "walrus.seal.io/project-name"
	LabelWalrusEnvironmentName string = "walrus.seal.io/environment-name"
	LabelWalrusServiceName     string = "walrus.seal.io/service-name"

	// LabelWalrusEnvironmentPath indicate environment with project name, format: projectName/environmentName.
	LabelWalrusEnvironmentPath string = "walrus.seal.io/environment-path"
	// LabelWalrusServicePath indicate service with project name and environment name,
	// format: projectName/environmentName/serviceName
	LabelWalrusServicePath string = "walrus.seal.io/service-path"

	// LabelWalrusManaged indicates whether the resource is managed by Seal.
	LabelWalrusManaged string = "walrus.seal.io/managed"
)
