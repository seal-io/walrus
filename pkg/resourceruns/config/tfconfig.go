package config

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"k8s.io/apimachinery/pkg/util/sets"

	runbus "github.com/seal-io/walrus/pkg/bus/resourcerun"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/predicate"
	"github.com/seal-io/walrus/pkg/dao/model/templateversion"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/crypto"
	opk8s "github.com/seal-io/walrus/pkg/operator/k8s"
	"github.com/seal-io/walrus/pkg/servervars"
	"github.com/seal-io/walrus/pkg/templates/openapi"
	"github.com/seal-io/walrus/pkg/templates/translator"
	"github.com/seal-io/walrus/pkg/terraform/config"
	"github.com/seal-io/walrus/pkg/terraform/parser"
	"github.com/seal-io/walrus/pkg/terraform/util"
	"github.com/seal-io/walrus/utils/log"
)

const (
	// _backendAPI the API path to terraform deploy backend.
	// Terraform will get and update deployment states from this API.
	_backendAPI = "/v1/projects/%s/environments/%s/resources/%s/runs/%s/terraform-states"
)

// TerraformConfigurator constructs the terraform config files for the run.
type TerraformConfigurator struct {
	logger log.Logger
}

func NewTerraformConfigurator() Configurator {
	return &TerraformConfigurator{
		logger: log.WithName("resource-run").WithName("tf"),
	}
}

func (c *TerraformConfigurator) LoadMain(
	ctx context.Context,
	mc model.ClientSet,
	opts *Options,
) (types.ResourceRunConfigData, error) {
	planConfig, err := c.LoadAll(ctx, mc, opts)
	if err != nil {
		return nil, err
	}

	return planConfig[config.FileMain], nil
}

func (c *TerraformConfigurator) LoadAll(
	ctx context.Context,
	mc model.ClientSet,
	opts *Options,
) (map[string]types.ResourceRunConfigData, error) {
	// Prepare terraform tfConfig.
	//  get module configs from resource run.
	moduleConfig, providerRequirements, err := c.getModuleConfig(ctx, mc, opts)
	if err != nil {
		return nil, err
	}
	// Merge current and previous required providers.
	providerRequirements = append(providerRequirements,
		opts.ResourceRun.PreviousRequiredProviders...)

	requiredProviders := make(map[string]*tfconfig.ProviderRequirement)
	for _, p := range providerRequirements {
		if _, ok := requiredProviders[p.Name]; !ok {
			requiredProviders[p.Name] = p.ProviderRequirement
		} else {
			c.logger.Warnf("duplicate provider requirement: %s", p.Name)
		}
	}

	runOpts := RunOpts{
		ResourceRun:   opts.ResourceRun,
		ResourceName:  opts.Context.Resource.Name,
		ProjectID:     opts.Context.Project.ID,
		EnvironmentID: opts.Context.Environment.ID,
	}
	// Parse module attributes.
	attrs, variables, dependencyOutputs, err := ParseModuleAttributes(
		ctx,
		mc,
		moduleConfig.Attributes,
		false,
		runOpts,
	)
	if err != nil {
		return nil, err
	}

	moduleConfig.Attributes = attrs

	// Update output sensitive with variables.
	wrapVariables, err := updateOutputWithVariables(variables, moduleConfig)
	if err != nil {
		return nil, err
	}

	// Prepare terraform config files to be mounted to secret.
	requiredProviderNames := sets.NewString()
	for _, p := range providerRequirements {
		requiredProviderNames = requiredProviderNames.Insert(p.Name)
	}

	address := fmt.Sprintf("%s%s", opts.SeverULR,
		fmt.Sprintf(_backendAPI,
			opts.Context.Project.ID,
			opts.Context.Environment.ID,
			opts.Context.Resource.ID,
			opts.ResourceRun.ID))

	tfCreateOpts := map[string]config.CreateOptions{
		config.FileMain: {
			TerraformOptions: &config.TerraformOptions{
				Token:                opts.Token,
				Address:              address,
				SkipTLSVerify:        !servervars.TlsCertified.Get(),
				ProviderRequirements: requiredProviders,
			},
			ProviderOptions: &config.ProviderOptions{
				RequiredProviderNames: requiredProviderNames.List(),
				Connectors:            opts.Connectors,
				SecretMonthPath:       opts.SecretMountPath,
				ConnectorSeparator:    parser.ConnectorSeparator,
			},
			ModuleOptions: &config.ModuleOptions{
				ModuleConfigs: []*config.ModuleConfig{moduleConfig},
			},
			VariableOptions: &config.VariableOptions{
				VariablePrefix:    _variablePrefix,
				ResourcePrefix:    _resourcePrefix,
				Variables:         wrapVariables,
				DependencyOutputs: dependencyOutputs,
			},
			OutputOptions: moduleConfig.Outputs,
		},
		config.FileVars: getVarConfigOptions(variables, dependencyOutputs),
	}

	inputConfigs := make(map[string]types.ResourceRunConfigData, len(tfCreateOpts))
	for k, v := range tfCreateOpts {
		inputConfigs[k], err = config.CreateConfigToBytes(v)
		if err != nil {
			return nil, err
		}
	}

	// Save input plan to resource run.
	opts.ResourceRun.InputConfigs = inputConfigs
	// If resource run does not inherit variables from cloned run,
	// then save the parsed variables to resource run.
	if len(opts.ResourceRun.Variables) == 0 {
		variableMap := make(crypto.Map[string, string], len(variables))
		for _, s := range variables {
			variableMap[s.Name] = string(s.Value)
		}
		opts.ResourceRun.Variables = variableMap
	}

	run, err := mc.ResourceRuns().UpdateOne(opts.ResourceRun).
		Set(opts.ResourceRun).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	if err = runbus.Notify(ctx, mc, run); err != nil {
		return nil, err
	}

	return inputConfigs, nil
}

// getModuleConfig returns module configs and required connectors to
// get terraform module config block from resource run.
func (c *TerraformConfigurator) getModuleConfig(
	ctx context.Context,
	mc model.ClientSet,
	opts *Options,
) (*config.ModuleConfig, []types.ProviderRequirement, error) {
	var (
		requiredProviders = make([]types.ProviderRequirement, 0)
		predicates        = make([]predicate.TemplateVersion, 0)
	)

	predicates = append(predicates, templateversion.And(
		templateversion.Version(opts.ResourceRun.TemplateVersion),
		templateversion.TemplateID(opts.ResourceRun.TemplateID),
	))

	templateVersion, err := mc.TemplateVersions().
		Query().
		Select(
			templateversion.FieldID,
			templateversion.FieldTemplateID,
			templateversion.FieldName,
			templateversion.FieldVersion,
			templateversion.FieldSource,
			templateversion.FieldSchema,
			templateversion.FieldUISchema,
		).
		Where(templateversion.Or(predicates...)).
		Only(ctx)
	if err != nil {
		return nil, nil, err
	}

	if len(templateVersion.Schema.RequiredProviders) != 0 {
		requiredProviders = append(requiredProviders, templateVersion.Schema.RequiredProviders...)
	}

	moduleConfig, err := getModuleConfig(opts.ResourceRun, templateVersion, opts)
	if err != nil {
		return nil, nil, err
	}

	return moduleConfig, requiredProviders, err
}

func getModuleConfig(
	run *model.ResourceRun,
	template *model.TemplateVersion,
	opts *Options,
) (*config.ModuleConfig, error) {
	mc := &config.ModuleConfig{
		Name:   opts.Context.Resource.Name,
		Source: template.Source,
	}

	if template.Schema.IsEmpty() {
		return mc, nil
	}

	mc.SchemaData = template.Schema.TemplateVersionSchemaData

	if template.Schema.OpenAPISchema == nil ||
		template.Schema.OpenAPISchema.Components == nil ||
		template.Schema.OpenAPISchema.Components.Schemas == nil {
		return mc, nil
	}

	// Variables.
	var (
		variableSchema     = template.Schema.VariableSchema()
		outputSchema       = template.Schema.OutputSchema()
		sensitiveVariables = sets.Set[string]{}
	)

	if variableSchema != nil {
		attrs, err := translator.ToGoTypeValues(run.ComputedAttributes, *variableSchema)
		if err != nil {
			return nil, err
		}

		mc.Attributes = attrs

		for n, v := range variableSchema.Properties {
			// Add sensitive from schema variable.
			if v.Value.WriteOnly {
				sensitiveVariables.Insert(fmt.Sprintf(`var\.%s`, n))
			}

			if n == types.WalrusContextVariableName {
				mc.Attributes[n] = opts.Context
			}
		}
	}

	// Outputs.
	if outputSchema != nil {
		sps := outputSchema.Properties
		mc.Outputs = make([]config.Output, 0, len(sps))

		sensitiveVariableRegex, err := matchAnyRegex(sensitiveVariables.UnsortedList())
		if err != nil {
			return nil, err
		}

		for k, v := range sps {
			origin := openapi.GetExtOriginal(v.Value.Extensions)
			co := config.Output{
				Sensitive:    v.Value.WriteOnly,
				Name:         k,
				ResourceName: opts.Context.Resource.Name,
				Value:        origin.ValueExpression,
			}

			if !v.Value.WriteOnly {
				// Update sensitive while output is from sensitive data, like secret.
				if sensitiveVariables.Len() != 0 && sensitiveVariableRegex.Match(origin.ValueExpression) {
					co.Sensitive = true
				}
			}

			mc.Outputs = append(mc.Outputs, co)
		}
	}

	return mc, nil
}

func (c *TerraformConfigurator) LoadProviders(
	connectors model.Connectors,
) (map[string]types.ResourceRunConfigData, error) {
	providerConfigs := make(map[string]types.ResourceRunConfigData, len(connectors))

	for _, c := range connectors {
		if c.Type != types.ConnectorTypeKubernetes {
			continue
		}

		_, s, err := opk8s.LoadApiConfig(*c)
		if err != nil {
			return nil, err
		}

		// NB(alex) the secret file name must be config + connector id to
		// match with terraform provider in config convert.
		secretFileName := util.GetK8sSecretName(c.ID.String())
		providerConfigs[secretFileName] = []byte(s)
	}

	return providerConfigs, nil
}
