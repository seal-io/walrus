package types

import (
	"github.com/hashicorp/terraform-config-inspect/tfconfig"

	"github.com/seal-io/seal/pkg/dao/types/property"
)

type ProviderRequirement struct {
	*tfconfig.ProviderRequirement

	Name string `json:"name"`
}

type ModuleSchema struct {
	Readme            string                `json:"readme"`
	Variables         property.Schemas      `json:"variables"`
	Outputs           property.Schemas      `json:"outputs"`
	RequiredProviders []ProviderRequirement `json:"requiredProviders"`
}
