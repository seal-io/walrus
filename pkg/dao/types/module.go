package types

type ModuleSchema struct {
	Readme                 string
	Variables              []ModuleVariable
	Outputs                []ModuleOutput
	RequiredConnectorTypes []string
}

type ModuleVariable struct {
	Name        string
	Description string

	// Supported types are string, number, bool, list(<TYPE>), map(<TYPE>), set(<TYPE>),
	// object({<ATTR NAME> = <TYPE>, ... }), tuple([<TYPE>, ...]).
	Type string
	// Default value of the variable.
	Default interface{}
	// Whether the variable is required or not.
	Required bool
	// Whether the variable is sensitive or not.
	Sensitive bool
	// UI label of the variable.
	Label string
	// UI group of the variable.
	Group string
	// Specify available options when the variable type is string.
	Options []string
	// Show the variable if the condition is true. For example, ShowIf: foo=bar.
	ShowIf string
}

type ModuleOutput struct {
	Name        string
	Description string
	Sensitive   bool
}
