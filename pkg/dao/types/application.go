package types

// ApplicationModule is a snapshot of model.ApplicationModuleRelationship to avoid cycle importing.
type ApplicationModule struct {
	// ID of module that configure to the application.
	ModuleID string `json:"moduleID"`
	// Name of the module customized to the application.
	Name string `json:"name"`
	// attributes to configure the module.
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

// ApplicationVariable holds the definition of a variable of application.
// TODO(thxCode): provide a general variable definition to migrate this and ModuleVariable.
type ApplicationVariable struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Description string      `json:"description,omitempty"`
	Default     interface{} `json:"default,omitempty"`
	Required    bool        `json:"required,omitempty"`
}
