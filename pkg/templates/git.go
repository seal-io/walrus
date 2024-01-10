package templates

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"entgo.io/ent/dialect/sql"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/hashicorp/go-version"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/template"
	"github.com/seal-io/walrus/pkg/dao/model/templateversion"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/settings"
	"github.com/seal-io/walrus/pkg/templates/loader"
	"github.com/seal-io/walrus/pkg/vcs"
	"github.com/seal-io/walrus/pkg/vcs/driver/gitlab"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/strs"
)

// CreateTemplateVersionsFromRepo creates template versions and
// return founded version count from a git repository worktree.
func CreateTemplateVersionsFromRepo(
	ctx context.Context,
	mc model.ClientSet,
	entity *model.Template,
	versions []*version.Version,
	versionSchema map[*version.Version]*schemaGroup,
) error {
	logger := log.WithName("template")

	// Create template versions.
	templateVersions, err := GetTemplateVersions(ctx, entity, versions, versionSchema)
	if err != nil {
		return err
	}

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

// syncTemplateFromRef clones a git repository create template and template version with reference.
// The reference can be a commit hash or a tag. The specified reference will not check semver format.
func syncTemplateFromRef(
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

	iconFile, err := getGitRepoIcon(r, repo.SubPath)
	if err != nil {
		logger.Errorf("failed to get icon url: %v", err)
		return err
	}

	icon, err := GetRepoFileRaw(repo, iconFile)
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

	// Load schema.
	originSchema, uiSchema, err := getSchemas(rootDir, entity.Name)
	if err != nil {
		logger.Errorf("%s:%s get schemas failed", entity.Name, repo.Reference)
		return err
	}

	// Create template.
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

	source, err := vcs.GetGitSource(repo.Link)
	if err != nil {
		return err
	}

	// Try to get version from reference.
	ref := repo.Reference
	if v, err := version.NewVersion(repo.Reference); err == nil {
		ref = v.Original()
	}

	conflictOptions := getTemplateVersionUpsertConflictOptions(entity)

	if repo.SubPath != "" {
		source += "//" + repo.SubPath
	}

	// Generate template version.
	tv := &model.TemplateVersion{
		TemplateID: entity.ID,
		Name:       entity.Name,
		Version:    ref,
		Source:     source + "?ref=" + repo.Reference,
		Schema:     *originSchema,
		UiSchema:   *uiSchema,
		ProjectID:  entity.ProjectID,
	}

	err = SetTemplateSchemaDefault(ctx, tv)
	if err != nil {
		return err
	}

	// Create or update a template version.
	return mc.TemplateVersions().Create().
		Set(tv).
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

	// Get icon image name.
	iconFile, err := getGitRepoIcon(r, repo.SubPath)
	if err != nil {
		logger.Errorf("failed to get icon url: %v", err)
		return err
	}

	icon, err := GetRepoFileRaw(repo, iconFile)
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

		t, err := query.
			Only(ctx)
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

	entity.Icon = icon

	// Create template.
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

	return CreateTemplateVersionsFromRepo(ctx, mc, entity, versions, versionSchema)
}

// GetTemplateVersions retrieves template versions from a git repository.
// It will save images to the database if they are found in the repository.
func GetTemplateVersions(
	ctx context.Context,
	entity *model.Template,
	newVersions []*version.Version,
	versionSchema map[*version.Version]*schemaGroup,
) (model.TemplateVersions, error) {
	var (
		logger = log.WithName("catalog")
		tvs    = make(model.TemplateVersions, 0, len(versionSchema))
	)

	source, err := vcs.GetGitSource(entity.Source)
	if err != nil {
		return nil, err
	}

	repo, err := vcs.ParseURLToRepo(entity.Source)
	if err != nil {
		return nil, err
	}

	for i := range newVersions {
		v := newVersions[i]
		tag := v.Original()

		schema, ok := versionSchema[v]
		if !ok {
			logger.Warnf("version schema not found, version: %s", tag)
			continue
		}

		if repo.SubPath != "" {
			source += "//" + repo.SubPath
		}

		// Generate template version.
		tv := &model.TemplateVersion{
			TemplateID: entity.ID,
			Name:       entity.Name,
			Version:    tag,
			Source:     source + "?ref=" + tag,
			Schema:     *schema.Schema,
			UiSchema:   *schema.UISchema,
			ProjectID:  entity.ProjectID,
		}

		err = SetTemplateSchemaDefault(ctx, tv)
		if err != nil {
			return nil, err
		}

		tvs = append(tvs, tv)
	}

	return tvs, nil
}

// getGitRepoIcon retrieves template icon from a git repository and return icon URL.
func getGitRepoIcon(repoLocal *git.Repository, subPath string) (string, error) {
	var (
		err error
		// Valid icon files.
		icons = []string{
			"icon.png",
			"icon.jpg",
			"icon.jpeg",
			"icon.svg",
		}
	)

	w, err := repoLocal.Worktree()
	if err != nil {
		return "", err
	}

	// Get icon URL.
	for _, icon := range icons {
		if subPath != "" {
			icon = filepath.Join(subPath, icon)
		}
		// If icon exists, get icon rawURL.
		if _, err := w.Filesystem.Stat(icon); err == nil {
			return icon, nil
		}
	}

	return "", nil
}

// GetRepoFileRaw returns raw URL of a file in a git repository.
func GetRepoFileRaw(repo *vcs.Repository, file string) (string, error) {
	if file == "" {
		return "", nil
	}

	endpoint, err := transport.NewEndpoint(repo.Link)
	if err != nil {
		return "", err
	}

	var (
		githubRawHost = "raw.githubusercontent.com"
		gitlabRawHost = "gitlab.com"
		giteeRawHost  = "gitee.com"
		ref           = "HEAD"
	)

	if repo.Reference != "" {
		ref = repo.Reference
	}

	switch endpoint.Host {
	case "github.com":
		return fmt.Sprintf("https://%s/%s/%s/%s/%s", githubRawHost, repo.Namespace, repo.Name, ref, file), nil
	case "gitlab.com":
		return fmt.Sprintf("https://%s/%s/%s/-/raw/%s/%s", gitlabRawHost, repo.Namespace, repo.Name, ref, file), nil
	case "gitee.com":
		return fmt.Sprintf("https://%s/%s/%s/raw/%s/%s", giteeRawHost, repo.Namespace, repo.Name, ref, file), nil
	}

	if repo.Driver == gitlab.Driver {
		return fmt.Sprintf("%s/-/raw/%s/%s", endpoint.String(), ref, file), nil
	}

	return "", nil
}

// syncTemplateVersion syncs template version schema and ui schema from remote.
func syncTemplateVersion(ctx context.Context, mc model.ClientSet, tv *model.TemplateVersion) error {
	repo, err := vcs.ParseURLToRepo(tv.Source)
	if err != nil {
		return err
	}

	tempDir := filepath.Join(os.TempDir(), "seal-template-version-"+strs.String(10))
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

	rootDir := w.Filesystem.Root()

	if repo.SubPath != "" {
		rootDir = filepath.Join(rootDir, repo.SubPath)
	}

	// Load schema.
	originSchema, uiSchema, err := getSchemas(rootDir, tv.Name)
	if err != nil {
		return err
	}

	// Update template version.
	return mc.TemplateVersions().UpdateOne(tv).
		SetSchema(*originSchema).
		SetUiSchema(*uiSchema).
		Exec(ctx)
}

// getSchemas returns template version schema and ui schema from a git repository.
func getSchemas(rootDir, templateName string) (*types.TemplateVersionSchema, *types.UISchema, error) {
	// Load schema.
	originSchema, err := loader.LoadOriginalSchema(rootDir, templateName)
	if err != nil {
		return nil, nil, err
	}

	fileSchema, err := loader.LoadFileSchema(rootDir, templateName)
	if err != nil {
		return nil, nil, err
	}

	uiSchema := originSchema.Expose()
	if fileSchema != nil {
		uiSchema = fileSchema.Expose()
	}

	satisfy, err := isConstraintSatisfied(fileSchema)
	if err != nil {
		return nil, nil, err
	}

	if !satisfy {
		return nil, nil, fmt.Errorf("%s does not satisfy server version constraint", templateName)
	}

	return originSchema, &uiSchema, nil
}
