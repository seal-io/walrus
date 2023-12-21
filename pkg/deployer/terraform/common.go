package terraform

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/seal-io/walrus/pkg/auths"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/settings"
	"github.com/seal-io/walrus/pkg/templates/openapi"
	"github.com/seal-io/walrus/pkg/templates/translator"
	"github.com/seal-io/walrus/pkg/terraform/config"
	"github.com/seal-io/walrus/pkg/terraform/parser"
	"github.com/seal-io/walrus/utils/json"
	"github.com/seal-io/walrus/utils/pointer"
	"k8s.io/apimachinery/pkg/util/sets"
)

func getServerURL(ctx context.Context, mc model.ClientSet) (string, error) {
	// Prepare address for terraform backend.
	serverAddress, err := settings.ServeUrl.Value(ctx, mc)
	if err != nil {
		return "", err
	}

	return serverAddress, nil
}

func getToken(ctx context.Context, mc model.ClientSet, subjectID, revisionID object.ID) (string, error) {
	// Prepare API token for terraform backend.
	const _1Day = 60 * 60 * 24

	at, err := auths.CreateAccessToken(ctx,
		mc, subjectID, types.TokenKindDeployment, revisionID.String(), pointer.Int(_1Day))
	if err != nil {
		return "", err
	}

	return at.AccessToken, nil
}

func getVarConfigOptions(
	variables model.Variables,
	resourceOutputs map[string]parser.OutputState,
) config.CreateOptions {
	varsConfigOpts := config.CreateOptions{
		Attributes: map[string]any{},
	}

	for _, v := range variables {
		varsConfigOpts.Attributes[_variablePrefix+v.Name] = v.Value
	}

	// Setup resource outputs.
	for n, v := range resourceOutputs {
		varsConfigOpts.Attributes[_resourcePrefix+n] = v.Value
	}

	return varsConfigOpts
}

func getModuleConfig(
	revision *model.ResourceRevision,
	template *model.TemplateVersion,
	opts createK8sSecretsOptions,
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

			if n == WalrusContextVariableName {
				mc.Attributes[n] = opts.Context
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

func updateOutputWithVariables(variables model.Variables, moduleConfig *config.ModuleConfig) (map[string]bool, error) {
	var (
		variableOpts         = make(map[string]bool)
		encryptVariableNames = sets.NewString()
	)

	for _, s := range variables {
		variableOpts[s.Name] = s.Sensitive

		if s.Sensitive {
			encryptVariableNames.Insert(_variablePrefix + s.Name)
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
