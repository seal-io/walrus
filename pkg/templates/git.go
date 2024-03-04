package templates

import (
	"context"
	"os"
	"path/filepath"

	"entgo.io/ent/dialect/sql"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/hashicorp/go-version"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/template"
	"github.com/seal-io/walrus/pkg/dao/model/templateversion"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/settings"
	"github.com/seal-io/walrus/pkg/vcs"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/strs"
)

// SyncTemplateFromGitRepo clones a git repository, retrieves all tags, creates a template and template versions.
// If the template already exists, it will update the template and template versions.
// Only semver tags will be used to create template versions.
func SyncTemplateFromGitRepo(
	ctx context.Context,
	mc model.ClientSet,
	entity *model.Template,
	repo *vcs.Repository,
) (err error) {
	logger := log.WithName("template")

	tempDir := filepath.Join(os.TempDir(), "seal-template-"+strs.String(10))
	defer os.RemoveAll(tempDir)

	u, err := transport.NewEndpoint(repo.Link)
	if err != nil {
		return err
	}

	// Clone git repository.
	r, err := vcs.CloneGitRepo(ctx, u.String(), tempDir, settings.SkipRemoteTLSVerify.ShouldValueBool(ctx, mc))
	if err != nil {
		return err
	}

	versions, err := vcs.GetGitRepoVersions(r)
	if err != nil {
		return err
	}

	versions, versionSchema, err := getValidVersions(entity, r, versions, repo.SubPath)
	if err != nil {
		return err
	}

	if len(versions) == 0 {
		logger.Warnf("no versions found for %s", repo.Name)

		// If template exists, update template status.
		query := mc.Templates().Query().
			Where(template.Name(entity.Name))

		if entity.ProjectID.Valid() {
			query.Where(template.ProjectID(entity.ProjectID))
		} else {
			query.Where(template.ProjectIDIsNil())
		}

		t, err := query.Only(ctx)
		if err != nil {
			if !model.IsNotFound(err) {
				return err
			}

			return nil
		}

		// When template source contains no valid versions, set template status to initialized.
		status.TemplateStatusInitialized.True(t, "")
		t.Status.SetSummary(status.WalkTemplate(&t.Status))

		return mc.Templates().UpdateOne(t).
			SetStatus(t.Status).
			Exec(ctx)
	}

	// Create template.
	iconURL, err := gitRepoIconURL(r, *repo)
	if err != nil {
		return err
	}

	entity.Icon = iconURL
	entity, err = CreateTemplate(ctx, mc, entity)
	if err != nil {
		return err
	}

	defer func() {
		rerr := updateTemplateStatus(ctx, mc, entity, err)
		if rerr != nil {
			logger.Errorf("failed to update template %s status: %v", entity.Name, rerr)
		}
	}()

	// Create template versions.
	templateVersions, err := genTemplateVersionsFromGitRepo(ctx, entity, versions, versionSchema, repo, mc)
	if err != nil {
		return err
	}
	return createTemplateVersions(ctx, mc, entity, templateVersions)
}

// SyncTemplateFromGitRef clones a git repository create template and template version with reference.
// The reference can be a commit hash or a tag. The specified reference will not check semver format.
func SyncTemplateFromGitRef(
	ctx context.Context,
	mc model.ClientSet,
	entity *model.Template,
	repo *vcs.Repository,
) (err error) {
	logger := log.WithName("template")

	tempDir := filepath.Join(os.TempDir(), "seal-template-"+strs.String(10))
	defer os.RemoveAll(tempDir)

	// Clone git repository.
	r, err := vcs.CloneGitRepo(ctx, repo.Link, tempDir, settings.SkipRemoteTLSVerify.ShouldValueBool(ctx, mc))
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	err = vcs.HardResetGitRepo(r, repo.Reference)
	if err != nil {
		return err
	}

	rootDir := w.Filesystem.Root()
	if repo.SubPath != "" {
		rootDir = filepath.Join(rootDir, repo.SubPath)
	}

	// Create template.
	icon, err := gitRepoIconURL(r, *repo)
	if err != nil {
		return err
	}

	entity.Icon = icon
	entity, err = CreateTemplate(ctx, mc, entity)
	if err != nil {
		return err
	}

	defer func() {
		rerr := updateTemplateStatus(ctx, mc, entity, err)
		if rerr != nil {
			logger.Errorf("failed to update template status: %v", rerr)
		}
	}()

	// Create template version.

	// Try to get version from reference.
	ref := repo.Reference
	if v, err := version.NewVersion(repo.Reference); err == nil {
		ref = v.Original()
	}

	source, err := generateSource(repo, repo.Reference)
	if err != nil {
		return err
	}

	// Load schema.
	schema, err := getSchemas(rootDir, entity.Name)
	if err != nil {
		logger.Errorf("%s:%s get schemas failed", entity.Name, repo.Reference)
		return err
	}

	// Generate template version.
	tv, err := genTemplateVersion(ctx, source, ref, entity, schema, mc)
	if err != nil {
		return err
	}

	return createTemplateVersions(ctx, mc, entity, []*model.TemplateVersion{tv})
}

// SyncTemplateVersion syncs template version schema and ui schema from remote.
func SyncTemplateVersion(ctx context.Context, mc model.ClientSet, tv *model.TemplateVersion) error {
	tempDir := filepath.Join(os.TempDir(), "seal-template-version-"+strs.String(10))
	defer os.RemoveAll(tempDir)

	// Clone git repository.
	repo, err := vcs.ParseURLToRepo(tv.Source)
	if err != nil {
		return err
	}

	r, err := vcs.CloneGitRepo(ctx, repo.Link, tempDir, settings.SkipRemoteTLSVerify.ShouldValueBool(ctx, mc))
	if err != nil {
		return err
	}

	if err = vcs.HardResetGitRepo(r, repo.Reference); err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	rootDir := w.Filesystem.Root()

	if repo.SubPath != "" {
		rootDir = filepath.Join(rootDir, repo.SubPath)
	}

	// Regenerate template version.
	schema, err := getSchemas(rootDir, tv.Name)
	if err != nil {
		return err
	}

	tp := &model.Template{
		ID:        tv.TemplateID,
		Name:      tv.Name,
		ProjectID: tv.ProjectID,
	}
	update, err := genTemplateVersion(ctx, tv.Source, tv.Version, tp, schema, mc)
	if err != nil {
		return err
	}

	return mc.TemplateVersions().UpdateOneID(tv.ID).
		SetSchema(update.Schema).
		SetUiSchema(update.UiSchema).
		SetSchemaDefaultValue(update.SchemaDefaultValue).
		Exec(ctx)
}

// createTemplateVersions creates template versions.
func createTemplateVersions(
	ctx context.Context,
	mc model.ClientSet,
	entity *model.Template,
	templateVersions []*model.TemplateVersion,
) error {
	logger := log.WithName("template")

	templateVersionCreates := make([]*model.TemplateVersionCreate, 0, len(templateVersions))

	for i := range templateVersions {
		create := mc.TemplateVersions().Create().Set(templateVersions[i])
		templateVersionCreates = append(templateVersionCreates, create)
	}

	logger.Debugf("create %d versions for template: %s", len(templateVersionCreates), entity.Name)

	conflictOptions := getTemplateVersionUpsertConflictOptions(entity)

	return mc.TemplateVersions().CreateBulk(templateVersionCreates...).
		OnConflict(conflictOptions...).
		UpdateNewValues().
		Exec(ctx)
}

// getTemplateVersionUpsertConflictOptions returns template version conflict options with template project id.
func getTemplateVersionUpsertConflictOptions(entity *model.Template) []sql.ConflictOption {
	if entity.ProjectID == "" {
		return []sql.ConflictOption{
			sql.ConflictWhere(sql.P().
				IsNull(templateversion.FieldProjectID)),
			sql.ConflictColumns(
				templateversion.FieldName,
				templateversion.FieldVersion,
			),
		}
	}

	return []sql.ConflictOption{
		sql.ConflictWhere(sql.P().
			NotNull(templateversion.FieldProjectID)),
		sql.ConflictColumns(
			templateversion.FieldName,
			templateversion.FieldVersion,
			templateversion.FieldProjectID,
		),
	}
}

// updateTemplateStatus updates template status.
func updateTemplateStatus(ctx context.Context, mc model.ClientSet, entity *model.Template, err error) error {
	if err != nil {
		status.TemplateStatusInitialized.False(entity, err.Error())
	} else {
		status.TemplateStatusReady.Reset(entity, "")
		status.TemplateStatusReady.True(entity, "")
	}

	// Update template status.
	entity.Status.SetSummary(status.WalkTemplate(&entity.Status))

	return mc.Templates().UpdateOne(entity).
		SetStatus(entity.Status).
		Exec(ctx)
}

// genTemplateVersionsFromGitRepo retrieves template versions from a git repository.
func genTemplateVersionsFromGitRepo(
	ctx context.Context,
	entity *model.Template,
	newVersions []*version.Version,
	versionSchema map[*version.Version]*loader.SchemaGroup,
	repo *vcs.Repository,
	mc model.ClientSet,
) (model.TemplateVersions, error) {
	var (
		logger = log.WithName("catalog")
		tvs    = make(model.TemplateVersions, 0, len(versionSchema))
	)

	for i := range newVersions {
		var (
			v   = newVersions[i]
			tag = v.Original()
		)

		s, err := generateSource(repo, tag)
		if err != nil {
			logger.Warnf("failed to generate source for %s:%s: %v", entity.Name, tag, err)
			continue
		}

		schema, ok := versionSchema[v]
		if !ok {
			logger.Warnf("version schema not found, version: %s", tag)
			continue
		}

		// Generate template version.
		tv, err := genTemplateVersion(ctx, s, tag, entity, schema, mc)
		if err != nil {
			return nil, err
		}

		tvs = append(tvs, tv)
	}

	return tvs, nil
}

// genTemplateVersion generates a template version.
func genTemplateVersion(
	ctx context.Context,
	source, ref string,
	entity *model.Template,
	schema *loader.SchemaGroup, mc model.ClientSet,
) (*model.TemplateVersion, error) {
	// Generate template version.
	tv := &model.TemplateVersion{
		TemplateID:       entity.ID,
		Name:             entity.Name,
		Version:          ref,
		Source:           source,
		Schema:           *schema.Schema,
		OriginalUISchema: *schema.UISchema,
		UiSchema:         *schema.UISchema,
		ProjectID:        entity.ProjectID,
	}

	// Set Schema Default.
	err := SetTemplateSchemaDefault(ctx, tv)
	if err != nil {
		return nil, err
	}

	// Set User Edited UI Schema.
	err = applyUserEditedUISchema(ctx, tv, mc)
	if err != nil {
		return nil, err
	}
	return tv, nil
}

// generateSource generates the source URL for template version.
func generateSource(repo *vcs.Repository, ref string) (string, error) {
	source, err := vcs.GetGitSource(repo.Link)
	if err != nil {
		return "", err
	}

	if repo != nil {
		if v, err := version.NewVersion(ref); err == nil {
			ref = v.Original()
		}

		if repo.SubPath != "" {
			source += "//" + repo.SubPath
		}

		source += "?ref=" + ref
	}

	return source, nil
}
