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
	"k8s.io/apimachinery/pkg/util/sets"

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
	versionSchema map[*version.Version]types.TemplateVersionSchema,
) error {
	logger := log.WithName("template")

	var ovs []string
	// Old versions.
	err := mc.TemplateVersions().Query().
		Select(templateversion.FieldVersion).
		Where(templateversion.TemplateID(entity.ID)).
		Modify(func(s *sql.Selector) {
			s.Select(templateversion.FieldVersion)
		}).
		Scan(ctx, &ovs)
	if err != nil {
		return err
	}
	// Old versions set.
	oldSet := sets.NewString(ovs...)
	// New versions set.
	newSet := sets.NewString()
	newVersions := make([]*version.Version, 0, len(versions))

	for i := range versions {
		newSet.Insert(versions[i].Original())

		if oldSet.Has(versions[i].Original()) {
			continue
		}

		newVersions = append(newVersions, versions[i])
	}

	// Create template versions.
	templateVersions, err := GetTemplateVersions(entity, newVersions, versionSchema)
	if err != nil {
		return err
	}

	templateVersionCreates := make([]*model.TemplateVersionCreate, 0, len(templateVersions))

	for i := range templateVersions {
		create := mc.TemplateVersions().Create().Set(templateVersions[i])

		templateVersionCreates = append(templateVersionCreates, create)
	}

	logger.Debugf("create %d versions for template: %s", len(templateVersionCreates), entity.Name)

	_, err = mc.TemplateVersions().CreateBulk(templateVersionCreates...).
		Save(ctx)
	if err != nil {
		return err
	}

	// Delete versions that are not in the repository tag version anymore.
	deleteVersions := sets.NewString(ovs...).Difference(newSet).List()

	_, err = mc.TemplateVersions().Delete().
		Where(
			templateversion.TemplateID(entity.ID),
			templateversion.VersionIn(deleteVersions...),
		).
		Exec(ctx)

	return err
}

// SyncTemplateFromRef clones a git repository create template and template version with reference.
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

	iconFile, err := getGitRepoIcon(r)
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

	// Load schema.
	schema, err := loader.LoadSchema(w.Filesystem.Root(), entity.Name, repo.Reference)
	if err != nil {
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

	var conflictOptions []sql.ConflictOption
	if entity.ProjectID == "" {
		conflictOptions = append(
			conflictOptions,
			sql.ConflictWhere(sql.P().
				IsNull(templateversion.FieldProjectID)),
			sql.ConflictColumns(
				templateversion.FieldName,
				templateversion.FieldVersion,
			),
		)
	} else {
		conflictOptions = append(
			conflictOptions,
			sql.ConflictWhere(sql.P().
				NotNull(templateversion.FieldProjectID)),
			sql.ConflictColumns(
				templateversion.FieldName,
				templateversion.FieldVersion,
				templateversion.FieldProjectID,
			),
		)
	}

	// Create or update a template version.
	return mc.TemplateVersions().Create().
		Set(&model.TemplateVersion{
			TemplateID: entity.ID,
			Name:       entity.Name,
			Version:    ref,
			Source:     source + "?ref=" + repo.Reference,
			Schema:     *schema,
			ProjectID:  entity.ProjectID,
		}).
		OnConflict(conflictOptions...).
		UpdateNewValues().
		Exec(ctx)
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
	iconFile, err := getGitRepoIcon(r)
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

	versions, versionSchema, err := getValidVersions(entity, r, versions)
	if err != nil {
		return err
	}

	if len(versions) == 0 {
		logger.Warnf("no versions found for %s", repo.Name)

		// If template exists, update template status.
		t, err := mc.Templates().Query().
			Where(template.Name(entity.Name)).
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
	entity *model.Template,
	newVersions []*version.Version,
	versionSchema map[*version.Version]types.TemplateVersionSchema,
) (model.TemplateVersions, error) {
	var (
		logger = log.WithName("catalog")
		tvs    = make(model.TemplateVersions, 0, len(versionSchema))
	)

	source, err := vcs.GetGitSource(entity.Source)
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

		tvs = append(tvs, &model.TemplateVersion{
			TemplateID: entity.ID,
			Name:       entity.Name,
			Version:    tag,
			Source:     source + "?ref=" + tag,
			Schema:     schema,
			ProjectID:  entity.ProjectID,
		})
	}

	return tvs, nil
}

// getGitRepoIcon retrieves template icon from a git repository and return icon URL.
func getGitRepoIcon(repoLocal *git.Repository) (string, error) {
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
	}

	if repo.Driver == gitlab.Driver {
		return fmt.Sprintf("%s/-/raw/%s/%s", endpoint.String(), ref, file), nil
	}

	return "", nil
}

// getValidVersions get valid terraform module versions.
func getValidVersions(
	entity *model.Template,
	r *git.Repository,
	versions []*version.Version,
) ([]*version.Version, map[*version.Version]types.TemplateVersionSchema, error) {
	logger := log.WithName("template")

	w, err := r.Worktree()
	if err != nil {
		return nil, nil, err
	}

	validVersions := make([]*version.Version, 0, len(versions))
	versionSchema := make(map[*version.Version]types.TemplateVersionSchema)

	for i := range versions {
		v := versions[i]
		tag := v.Original()

		resetRef, err := vcs.GetRepoRef(r, tag)
		if err != nil {
			logger.Warnf("failed to get \"%s:%s\" of catalog %q git reference: %v",
				entity.Name, tag, entity.CatalogID, err)
			continue
		}

		hash := resetRef.Hash()

		// If tag is not a commit hash, get commit hash from tag object target.
		object, err := r.TagObject(hash)
		if err == nil {
			hash = object.Target
		}

		err = w.Reset(&git.ResetOptions{
			Commit: hash,
			Mode:   git.HardReset,
		})
		if err != nil {
			logger.Warnf("failed set \"%s:%s\" of catalog %q: %v", entity.Name, tag, entity.CatalogID, err)
			continue
		}

		logger.Debugf("get \"%s:%s\" of catalog %q schema", entity.Name, tag, entity.CatalogID)
		dir := w.Filesystem.Root()

		schema, err := loader.LoadSchema(dir, entity.Name, tag)
		if err != nil {
			logger.Warnf("failed to load \"%s:%s\" of catalog %q schema: %v", entity.Name, tag, entity.CatalogID, err)
			continue
		}

		if err = schema.Validate(); err != nil {
			logger.Warnf(
				"failed to validate \"%s:%s\" of catalog %q schema: %v",
				entity.Name,
				tag,
				entity.CatalogID,
				err,
			)

			continue
		}

		validVersions = append(validVersions, v)
		versionSchema[v] = *schema
	}

	return validVersions, versionSchema, nil
}
