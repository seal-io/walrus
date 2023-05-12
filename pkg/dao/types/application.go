package types

import (
	"github.com/seal-io/seal/pkg/dao/types/property"
)

// ApplicationModule is a snapshot of model.ApplicationModuleRelationship to avoid cycle importing.
type ApplicationModule struct {
	// ID of module that configure to the application.
	ModuleID string `json:"moduleID"`
	// Version of the module that configure to the application.
	Version string `json:"version"`
	// Name of the module customized to the application.
	Name string `json:"name"`
	// Attributes to configure the module.
	Attributes property.Values `json:"attributes,omitempty"`
}
