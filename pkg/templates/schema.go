package templates

import (
	"context"
	"fmt"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/predicate"
	"github.com/seal-io/walrus/pkg/dao/model/templateversion"
	"github.com/seal-io/walrus/pkg/templates/loader"
	"github.com/seal-io/walrus/pkg/templates/openapi"
)

// getSchemas returns template version schema and ui schema from a git repository.
func getSchemas(rootDir, templateName string) (*schemaGroup, error) {
	// Load schema.
	originSchema, err := loader.LoadOriginalSchema(rootDir, templateName)
	if err != nil {
		return nil, err
	}

	if err = originSchema.Validate(); err != nil {
		return nil, fmt.Errorf("error validate schema for %s: %w", templateName, err)
	}

	fileSchema, err := loader.LoadFileSchema(rootDir, templateName)
	if err != nil {
		return nil, fmt.Errorf("error load file schema for %s: %w", templateName, err)
	}

	uiSchema := originSchema.Expose(openapi.WalrusContextVariableName)
	if fileSchema != nil && !fileSchema.IsEmpty() {
		uiSchema = fileSchema.Expose()
	}

	satisfy, err := isConstraintSatisfied(fileSchema)
	if err != nil {
		return nil, fmt.Errorf("error check server version constraint for %s: %w", templateName, err)
	}

	if !satisfy {
		return nil, fmt.Errorf("%s does not satisfy server version constraint", templateName)
	}

	return &schemaGroup{
		UISchema: &uiSchema,
		Schema:   originSchema,
	}, nil
}

// applySchemaDefault applies exist user edit schema to ui schema.
func applyUserEditedUISchema(ctx context.Context, tv *model.TemplateVersion, mc model.ClientSet) error {
	ps := []predicate.TemplateVersion{
		templateversion.Name(tv.Name),
		templateversion.Version(tv.Version),
	}
	if tv.ProjectID.Valid() {
		ps = append(ps, templateversion.ProjectID(tv.ProjectID))
	}

	existed, err := mc.TemplateVersions().Query().
		Where(ps...).
		Select(templateversion.FieldUiSchema).
		Only(ctx)
	if err != nil {
		if model.IsNotFound(err) {
			return nil
		}
		return err
	}

	if existed.UiSchema.IsUserEdited() {
		tv.UiSchema = existed.UiSchema
	}
	return nil
}
