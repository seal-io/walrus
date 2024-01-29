package config

import (
	"github.com/hashicorp/terraform-config-inspect/tfconfig"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
)

const (
	FileMain = "main.tf"
	FileVars = "terraform.tfvars"
)

// ModuleConfig is a struct with model.Template and its variables.
type ModuleConfig struct {
	// Name is the module name.
	Name string
	// Source is the module source.
	Source string
	// SchemaData is the data module schema.
	SchemaData types.TemplateVersionSchemaData
	// Attributes is the attributes of the module.
	Attributes map[string]any
	// Outputs is the module outputs.
	Outputs []Output
}

// CreateOptions represents the CreateOptions of the Config.
type CreateOptions struct {
	Attributes       map[string]any
	TerraformOptions *TerraformOptions
	ProviderOptions  *ProviderOptions
	ModuleOptions    *ModuleOptions
	VariableOptions  *VariableOptions
	OutputOptions    OutputOptions
}

type (
	// TerraformOptions is the options to create terraform block.
	TerraformOptions struct {
		// Token is the backend token to authenticate with the Seal Server of the terraform config.
		Token string
		// Address is the backend address of the terraform config.
		Address string
		// SkipTLSVerify is the backend cert verification of the terraform config.
		SkipTLSVerify bool

		// ProviderRequirements is the required providers of the terraform config.
		ProviderRequirements map[string]*tfconfig.ProviderRequirement
	}

	// ProviderOptions is the options to create provider blocks.
	ProviderOptions struct {
		// SecretMountPath is the mount path of the secret in the terraform config.
		SecretMonthPath string
		// ConnectorSeparator is the separator of the terraform provider alias name and id.
		ConnectorSeparator string
		// RequiredProviderNames is the required providers of the terraform config.
		// E.g. ["kubernetes", "helm"].
		RequiredProviderNames []string
		Connectors            model.Connectors
	}

	// ModuleOptions is the options to create module blocks.
	ModuleOptions struct {
		// ModuleConfigs is the module configs of the deployment.
		ModuleConfigs []*ModuleConfig
	}

	// VariableOptions is the options to create variables blocks.
	VariableOptions struct {
		// VariablePrefix is the prefix of the variable name.
		VariablePrefix string
		// ResourcePrefix is the prefix of the Walrus resource variable name.
		ResourcePrefix string
		// Variables is map with name in key and sensitive flag in value.
		Variables map[string]bool
		// DependencyOutputs is the map of the variable name and value.
		DependencyOutputs map[string]types.OutputValue
	}

	// OutputOptions is the options to create outputs blocks.
	OutputOptions []Output
	// Output indicate the output name and module.
	Output struct {
		ResourceName string
		Name         string
		Sensitive    bool
		Value        []byte
	}
)
