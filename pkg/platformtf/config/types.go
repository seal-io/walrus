package config

import (
	"github.com/hashicorp/terraform-config-inspect/tfconfig"

	"github.com/seal-io/seal/pkg/dao/model"
)

const (
	FileMain = "main.tf"
	FileVars = "terraform.tfvars"
)

// ModuleConfig is a struct with model.Module and its variables.
type ModuleConfig struct {
	// Name is the name of the app module relationship.
	Name          string
	ModuleVersion *model.ModuleVersion
	// Attributes is the attributes of the module.
	Attributes map[string]interface{}
}

// CreateOptions represents the CreateOptions of the Config.
type CreateOptions struct {
	Attributes       map[string]interface{}
	TerraformOptions *TerraformOptions
	ProviderOptions  *ProviderOptions
	ModuleOptions    *ModuleOptions
	VariableOptions  *VariableOptions
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
		// e.g. ["kubernetes", "helm"]
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
		// Prefix is the prefix of the variable name.
		Prefix string
		// SecretNames is the list of the secret name, type is always string.
		SecretNames []string
		// VariableNameAndTypes is the map of the variable name and type.
		VariableNameAndTypes map[string]string
	}
)
