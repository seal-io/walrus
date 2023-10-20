package types

type Selector struct {
	ProjectName       string            `json:"projectName,omitempty"`
	EnvironmentName   string            `json:"environmentName,omitempty"`
	EnvironmentType   string            `json:"environmentType,omitempty"`
	ResourceLabels    map[string]string `json:"resourceLabels,omitempty"`
	EnvironmentLabels map[string]string `json:"environmentLabels,omitempty"`
}
