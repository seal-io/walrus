package config

const (
	// BlockTypeTerraform represents the terraform block.
	BlockTypeTerraform = "terraform"
	// BlockTypeBackend represents the backend block inside terraform block.
	BlockTypeBackend = "backend"
	// BlockTypeRequiredProviders represents the required_providers block inside terraform block.
	BlockTypeRequiredProviders = "required_providers"

	// BlockTypeProvider represents the provider block.
	BlockTypeProvider = "provider"
	// BlockTypeModule represents the module block.
	BlockTypeModule = "module"
	// BlockTypeVariable represents the variable block.
	BlockTypeVariable = "variable"
	// BlockTypeOutput represents the output block.
	BlockTypeOutput = "output"
	// BlockTypeResource represents the resource block.
	BlockTypeResource = "resource"

	// BlockFormatHCL represents the hcl format.
	BlockFormatHCL = "hcl"
	// BlockFormatJSON represents the json format.
	BlockFormatJSON = "json"
)
