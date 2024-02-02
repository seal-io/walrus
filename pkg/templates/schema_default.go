package templates

import (
	"context"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/templates/openapi"
	"github.com/seal-io/walrus/utils/json"
)

// SetResourceDefinitionSchemaDefault generate the schema default for resource definition.
// Required rd.Edges.MatchingRules.
func SetResourceDefinitionSchemaDefault(
	ctx context.Context,
	client model.ClientSet,
	rd *model.ResourceDefinition,
) (err error) {
	for i := range rd.Edges.MatchingRules {
		rule := rd.Edges.MatchingRules[i]

		tv, err := client.TemplateVersions().Get(ctx, rule.TemplateID)
		if err != nil {
			return err
		}

		rd.Edges.MatchingRules[i].SchemaDefaultValue, err = openapi.GenSchemaDefaultWithAttribute(
			ctx,
			tv.UiSchema.VariableSchema(),
			rule.Attributes,
			tv.SchemaDefaultValue,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// SetTemplateSchemaDefault set the schema default for template.
func SetTemplateSchemaDefault(ctx context.Context, tv *model.TemplateVersion) error {
	originSchemaDefault, err := openapi.GenSchemaDefaultPatch(ctx, tv.Schema.VariableSchema())
	if err != nil {
		return err
	}

	uiSchemaDefault, err := openapi.GenSchemaDefaultPatch(ctx, tv.UiSchema.VariableSchema())
	if err != nil {
		return err
	}

	dv, err := json.ApplyPatches(originSchemaDefault, uiSchemaDefault)
	if err != nil {
		return err
	}

	tv.SchemaDefaultValue = dv

	return nil
}
