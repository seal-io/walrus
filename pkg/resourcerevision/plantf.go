package resourcerevision

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"k8s.io/apimachinery/pkg/util/sets"

	apiconfig "github.com/seal-io/walrus/pkg/apis/config"
	"github.com/seal-io/walrus/pkg/auths"
	revisionbus "github.com/seal-io/walrus/pkg/bus/resourcerevision"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/predicate"
	"github.com/seal-io/walrus/pkg/dao/model/templateversion"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/crypto"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	opk8s "github.com/seal-io/walrus/pkg/operator/k8s"
	pkgresource "github.com/seal-io/walrus/pkg/resource"
	"github.com/seal-io/walrus/pkg/settings"
	"github.com/seal-io/walrus/pkg/templates/openapi"
	"github.com/seal-io/walrus/pkg/templates/translator"
	"github.com/seal-io/walrus/pkg/terraform/config"
	"github.com/seal-io/walrus/pkg/terraform/parser"
	"github.com/seal-io/walrus/pkg/terraform/util"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/pointer"
)

// TerraformPlan handler the revision Plan of terraform deployer.
type TerraformPlan struct {
	logger log.Logger
}

// NewTerraformPlan creates a new terraform plan of service revision.
func NewTerraformPlan() *TerraformPlan {
	return &TerraformPlan{
		logger: log.WithName("service.revision-plan"),
	}
}

// LoadPlan will get the revision input plan of terraform.
// It will generate the main.tf for the revision.
func (t TerraformPlan) LoadPlan(ctx context.Context, mc model.ClientSet, opts *PlanOptions) ([]byte, error) {
	configBytes, err := t.LoadConfigs(ctx, mc, opts)
	if err != nil {
		return nil, err
	}

	return configBytes[config.FileMain], nil
}

func (t TerraformPlan) LoadConfigs(
	ctx context.Context,
	mc model.ClientSet,
	opts *PlanOptions,
) (map[string][]byte, error) {
	// Prepare terraform tfConfig.
	//  get module configs from service revision.
	moduleConfig, providerRequirements, err := t.getModuleConfig(ctx, mc, opts)
	if err != nil {
		return nil, err
	}

	// Merge current and previous required providers.
	providerRequirements = append(providerRequirements,
		opts.ResourceRevision.PreviousRequiredProviders...)

	requiredProviders := make(map[string]*tfconfig.ProviderRequirement, 0)
	for _, p := range providerRequirements {
		if _, ok := requiredProviders[p.Name]; !ok {
			requiredProviders[p.Name] = p.ProviderRequirement
		} else {
			t.logger.Warnf("duplicate provider requirement: %s", p.Name)
		}
	}

	resourceOpts := pkgresource.ParseAttributesOptions{
		ResourceRevision: opts.ResourceRevision,
		ProjectID:        opts.Context.Project.ID,
		EnvironmentID:    opts.Context.Environment.ID,
		ResourceName:     opts.Context.Resource.Name,
	}
	// Parse module attributes.
	attrs, variables, dependencyOutputs, err := pkgresource.ParseModuleAttributes(
		ctx,
		mc,
		moduleConfig.Attributes,
		resourceOpts,
	)
	if err != nil {
		return nil, err
	}

	moduleConfig.Attributes = attrs

	// Update output sensitive with variables.
	wrapVariables, err := setOutputSensitiveWithVariables(variables, moduleConfig)
	if err != nil {
		return nil, err
	}

	// Prepare terraform config files to be mounted to secret.
	requiredProviderNames := sets.NewString()
	for _, p := range providerRequirements {
		requiredProviderNames = requiredProviderNames.Insert(p.Name)
	}

	address, token, err := t.getBackendConfig(ctx, mc, opts)
	if err != nil {
		return nil, err
	}

	// Options for create terraform config.
	planConfigOptions := map[string]config.CreateOptions{
		config.FileMain: {
			TerraformOptions: &config.TerraformOptions{
				Token:                token,
				Address:              address,
				SkipTLSVerify:        !apiconfig.TlsCertified.Get(),
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
				VariablePrefix:    pkgresource.VariablePrefix,
				ResourcePrefix:    pkgresource.ResourcePrefix,
				Variables:         wrapVariables,
				DependencyOutputs: dependencyOutputs,
			},
			OutputOptions: moduleConfig.Outputs,
		},
		config.FileVars: getVarConfigOptions(variables, dependencyOutputs),
	}
	planConfigs := make(map[string][]byte, 0)

	for k, v := range planConfigOptions {
		planConfigs[k], err = config.CreateConfigToBytes(v)
		if err != nil {
			return nil, err
		}
	}

	// Save input plan to service revision.
	opts.ResourceRevision.InputPlan = string(planConfigs[config.FileMain])
	// If service revision does not inherit variables from cloned revision,
	// then save the parsed variables to service revision.
	if len(opts.ResourceRevision.Variables) == 0 {
		variableMap := make(crypto.Map[string, string], len(variables))
		for _, s := range variables {
			variableMap[s.Name] = string(s.Value)
		}
		opts.ResourceRevision.Variables = variableMap
	}

	status.ResourceRevisionStatusReady.Reset(opts.ResourceRevision, "")
	opts.ResourceRevision.Status.SetSummary(status.WalkResourceRevision(&opts.ResourceRevision.Status))

	revision, err := mc.ResourceRevisions().UpdateOne(opts.ResourceRevision).
		Set(opts.ResourceRevision).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	if err = revisionbus.Notify(ctx, mc, revision); err != nil {
		return nil, err
	}

	return planConfigs, nil
}

// LoadConnectorConfigs loads the connector for terraform provider.
func (t TerraformPlan) LoadConnectorConfigs(connectors model.Connectors) (map[string][]byte, error) {
	secretData := make(map[string][]byte)

	for _, c := range connectors {
		// Note(alex) only load k8s connector config for deployment.
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
		secretData[secretFileName] = []byte(s)
	}

	return secretData, nil
}

func (t TerraformPlan) getBackendConfig(
	ctx context.Context,
	mc model.ClientSet,
	opts *PlanOptions,
) (address, token string, err error) {
	// Prepare address for terraform backend.
	serverAddress, err := settings.ServeUrl.Value(ctx, mc)
	if err != nil {
		return "", "", err
	}

	if serverAddress == "" {
		return "", "", errors.New("server address is empty")
	}
	address = fmt.Sprintf("%s%s", serverAddress,
		fmt.Sprintf(_backendAPI,
			opts.Context.Project.ID,
			opts.Context.Environment.ID,
			opts.Context.Resource.ID,
			opts.ResourceRevision.ID))

	// Prepare API token for terraform backend.
	const _1Day = 60 * 60 * 24

	at, err := auths.CreateAccessToken(ctx,
		mc, opts.SubjectID, types.TokenKindDeployment, string(opts.ResourceRevision.ID), pointer.Int(_1Day))
	if err != nil {
		return "", "", err
	}

	token = at.AccessToken

	return
}

// getModuleConfig returns module configs and required connectors to
// get terraform module config block from service revision.
func (t TerraformPlan) getModuleConfig(
	ctx context.Context,
	mc model.ClientSet,
	opts *PlanOptions,
) (*config.ModuleConfig, []types.ProviderRequirement, error) {
	var (
		requiredProviders = make([]types.ProviderRequirement, 0)
		predicates        = make([]predicate.TemplateVersion, 0)
	)

	predicates = append(predicates, templateversion.And(
		templateversion.Name(opts.ResourceRevision.TemplateName),
		templateversion.Version(opts.ResourceRevision.TemplateVersion),
		templateversion.TemplateID(opts.ResourceRevision.TemplateID),
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
		).
		Where(templateversion.Or(predicates...)).
		Only(ctx)
	if err != nil {
		return nil, nil, err
	}

	if len(templateVersion.Schema.RequiredProviders) != 0 {
		requiredProviders = append(requiredProviders, templateVersion.Schema.RequiredProviders...)
	}

	moduleConfig, err := getModuleConfig(opts.ResourceRevision, templateVersion, opts)
	if err != nil {
		return nil, nil, err
	}

	return moduleConfig, requiredProviders, err
}

// getModuleConfig get module config of terraform.
func getModuleConfig(
	revision *model.ResourceRevision,
	template *model.TemplateVersion,
	opts *PlanOptions,
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
		variablesSchema    = template.Schema.VariableSchema()
		outputsSchemas     = template.Schema.OutputSchema()
		sensitiveVariables = sets.Set[string]{}
	)

	if variablesSchema != nil {
		attrs, err := translator.ToGoTypeValues(revision.Attributes, *variablesSchema)
		if err != nil {
			return nil, err
		}

		mc.Attributes = attrs

		for n, v := range variablesSchema.Properties {
			// Add sensitive from schema variable.
			if v.Value.WriteOnly {
				sensitiveVariables.Insert(fmt.Sprintf(`var\.%s`, n))
			}

			mc.Attributes = attrs

			for n, v := range variablesSchema.Properties {
				// Add sensitive from schema variable.
				if v.Value.WriteOnly {
					sensitiveVariables.Insert(fmt.Sprintf(`var\.%s`, n))
				}

				if n == pkgresource.WalrusContextVariableName {
					mc.Attributes[n] = opts.Context
				}
			}
		}
	}

	// Outputs.
	if outputsSchemas != nil {
		sps := outputsSchemas.Properties
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

// matchAnyRegex get regex of match any list string.
func matchAnyRegex(list []string) (*regexp.Regexp, error) {
	var sb strings.Builder

	sb.WriteString("(")

	for i, v := range list {
		sb.WriteString(v)

		if i < len(list)-1 {
			sb.WriteString("|")
		}
	}

	sb.WriteString(")")

	return regexp.Compile(sb.String())
}

// setOutputSensitiveWithVariables update output with variables.
// Sensitive output should not show the value.
func setOutputSensitiveWithVariables(
	variables model.Variables,
	moduleConfig *config.ModuleConfig,
) (map[string]bool, error) {
	var (
		variableOpts         = make(map[string]bool)
		encryptVariableNames = sets.NewString()
	)

	for _, s := range variables {
		variableOpts[s.Name] = s.Sensitive

		if s.Sensitive {
			encryptVariableNames.Insert(pkgresource.VariablePrefix + s.Name)
		}
	}

	if encryptVariableNames.Len() == 0 {
		return variableOpts, nil
	}

	reg, err := matchAnyRegex(encryptVariableNames.UnsortedList())
	if err != nil {
		return nil, err
	}

	var shouldEncryptAttr []string

	for k, v := range moduleConfig.Attributes {
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}

		matches := reg.FindAllString(string(b), -1)
		if len(matches) != 0 {
			shouldEncryptAttr = append(shouldEncryptAttr, fmt.Sprintf(`var\.%s`, k))
		}
	}

	// Outputs use encrypted variable should set to sensitive.
	for i, v := range moduleConfig.Outputs {
		if v.Sensitive {
			continue
		}

		reg, err := matchAnyRegex(shouldEncryptAttr)
		if err != nil {
			return nil, err
		}

		if reg.MatchString(string(v.Value)) {
			moduleConfig.Outputs[i].Sensitive = true
		}
	}

	return variableOpts, nil
}

// getVarConfigOptions get terraform tf.vars config.
func getVarConfigOptions(variables model.Variables, serviceOutputs map[string]parser.OutputState) config.CreateOptions {
	varsConfigOpts := config.CreateOptions{
		Attributes: map[string]any{},
	}

	for _, v := range variables {
		varsConfigOpts.Attributes[pkgresource.VariablePrefix+v.Name] = v.Value
	}

	// Setup service outputs.
	for n, v := range serviceOutputs {
		varsConfigOpts.Attributes[pkgresource.ResourcePrefix+n] = v.Value
	}

	return varsConfigOpts
}
