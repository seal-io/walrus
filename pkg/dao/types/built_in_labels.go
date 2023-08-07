package types

const (
	// LabelPrefix is used for generate label's field names.
	LabelPrefix = "label:"
)

// built-in labels.
const (
	LabelSealProjectName     string = "seal.io/project-name"
	LabelSealEnvironmentName string = "seal.io/environment-name"
	LabelSealServiceName     string = "seal.io/service-name"

	// LabelSealEnvironmentPath indicate environment with project name, format: projectName/environmentName.
	LabelSealEnvironmentPath string = "seal.io/environment-path"
	// LabelSealServicePath indicate service with project name and environment name,
	// format: projectName/environmentName/serviceName
	LabelSealServicePath string = "seal.io/service-path"

	// LabelSealManaged indicates whether the resource is managed by Seal.
	LabelSealManaged string = "seal.io/managed"
)
