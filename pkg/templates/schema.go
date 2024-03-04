package templates

import (
	"context"
	"fmt"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/predicate"
	"github.com/seal-io/walrus/pkg/dao/model/templateversion"
	"github.com/seal-io/walrus/pkg/templates/loader"
)

// getSchemas returns template version schema and ui schema from a git repository.
func getSchemas(rootDir, templateName string) (*loader.SchemaGroup, error) {
	// Load schema.
	s, err := loader.LoadSchema(rootDir, templateName)
	if err != nil {
		return nil, err
	}

	satisfy, err := isConstraintSatisfied(s.UISchema)
	if err != nil {
		return nil, fmt.Errorf("error check server version constraint for %s: %w", templateName, err)
	}

	if !satisfy {
		return nil, fmt.Errorf("%s does not satisfy server version constraint", templateName)
	}

	return s, nil
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
