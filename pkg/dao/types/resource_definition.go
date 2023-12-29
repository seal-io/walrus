package types

type Selector struct {
	ProjectNames      []string          `json:"projectNames,omitempty"`
	ProjectLabels     map[string]string `json:"projectLabels,omitempty"`
	EnvironmentNames  []string          `json:"environmentNames,omitempty"`
	EnvironmentTypes  []string          `json:"environmentTypes,omitempty"`
	EnvironmentLabels map[string]string `json:"environmentLabels,omitempty"`
	ResourceLabels    map[string]string `json:"resourceLabels,omitempty"`
}

type MatchResourceMetadata struct {
	ProjectName       string
	EnvironmentName   string
	EnvironmentType   string
	ProjectLabels     map[string]string
	EnvironmentLabels map[string]string
	ResourceLabels    map[string]string
}
