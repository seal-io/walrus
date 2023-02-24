package types

type ModuleSchema struct {
	Readme                 string
	Variables              []ModuleVariable
	Outputs                []ModuleOutput
	RequiredConnectorTypes []string
}

type ModuleVariable struct {
	Name        string
	Type        string
	Description string

	Default   interface{}
	Required  bool
	Sensitive bool
}

type ModuleOutput struct {
	Name        string
	Description string
	Sensitive   bool
}
