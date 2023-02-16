package types

// ApplicationModule is a snapshot of model.ApplicationModuleRelationship to avoid cycle importing.
type ApplicationModule struct {
	// ID of module that configure to the application.
	ModuleID string `json:"moduleID"`
	// Name of the module customized to the application.
	Name string `json:"name"`
	// Variables to configure the module.
	Variables map[string]interface{} `json:"variables,omitempty"`
}
