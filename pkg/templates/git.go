package templates

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"entgo.io/ent/dialect/sql"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/hashicorp/go-version"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/templateversion"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/vcs"
	"github.com/seal-io/seal/pkg/vcs/driver/gitlab"
	"github.com/seal-io/seal/utils/log"
	"github.com/seal-io/seal/utils/strs"
)

// CreateTemplateVersionsFromRepo creates template versions and
// return founded version count from a git repository worktree.
func CreateTemplateVersionsFromRepo(
	ctx context.Context,
	mc model.ClientSet,
	entity *model.Template,
	r *git.Repository,
	versions []*version.Version,
) error {
	logger := log.WithName("template")

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	var ovs []string
	// Old versions.
	err = mc.TemplateVersions().Query().
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
		newSet.Insert(versions[i].String())

		if oldSet.Has(versions[i].String()) {
			continue
		}

		newVersions = append(newVersions, versions[i])
	}

	// Create template versions.
	templateVersions, err := GetTemplateVersions(entity, w, newVersions)
	if err != nil {
		return err
	}

	templateVersionCreates := make([]*model.TemplateVersionCreate, 0, len(templateVersions))

	for i := range templateVersions {
		create := mc.TemplateVersions().Create().Set(templateVersions[i])

		templateVersionCreates = append(templateVersionCreates, create)
	}

	logger.Debugf("create %s versions number: %v", entity.Name, len(templateVersionCreates))

	err = mc.TemplateVersions().CreateBulk(templateVersionCreates...).
		Exec(ctx)
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
	repo *vcs.Repository,
) (err error) {
	logger := log.WithName("template")

	tempDir := filepath.Join(os.TempDir(), "seal-template-"+strs.String(10))
	defer os.RemoveAll(tempDir)

	// Clone git repository.
	r, err := vcs.CloneGitRepo(repo.Link, tempDir)
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

	schema, err := loadTerraformTemplateSchema(w.Filesystem.Root())
	if err != nil {
		return err
	}

	iconFile, err := getGitRepoIcon(r)
	if err != nil {
		logger.Errorf("failed to get icon url: %v", err)
		return err
	}

	icon, err := GetRepoFileRaw(repo, iconFile, nil)
	if err != nil {
		return err
	}

	// Create template.
	entity, err := CreateTemplate(ctx, mc, &model.Template{
		Name:        repo.Name,
		Source:      repo.Link,
		Description: repo.Description,
		Icon:        icon,
	})
	if err != nil {
		return err
	}

	defer func() {
		rerr := updateTemplateStatus(ctx, mc, entity, err)
		if rerr != nil {
			logger.Errorf("failed to update template status: %v", rerr)
		}
	}()

	u, err := transport.NewEndpoint(repo.Link)
	if err != nil {
		return err
	}

	// Try to get version from reference.
	ref := repo.Reference
	if v, err := version.NewVersion(repo.Reference); err == nil {
		ref = v.String()
	}

	return createTemplateVersion(ctx, mc, &model.TemplateVersion{
		TemplateID:   entity.ID,
		TemplateName: entity.Name,
		Version:      ref,
		Source:       u.Host + u.Path + "?ref=" + repo.Reference,
		Schema:       schema,
	})
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

// createTemplateVersion creates a template version.
func createTemplateVersion(ctx context.Context, mc model.ClientSet, tv *model.TemplateVersion) error {
	// Delete old template version.
	_, err := mc.TemplateVersions().Delete().
		Where(templateversion.TemplateID(tv.TemplateID)).
		Exec(ctx)
	if err != nil {
		return err
	}

	// Create new template version.
	return mc.TemplateVersions().Create().
		Set(tv).
		Exec(ctx)
}

// SyncTemplateFromGitRepo clones a git repository, retrieves all tags, creates a template and template versions.
// If the template already exists, it will update the template and template versions.
// Only semver tags will be used to create template versions.
func SyncTemplateFromGitRepo(
	ctx context.Context,
	mc model.ClientSet,
	c *model.Catalog,
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
	r, err := vcs.CloneGitRepo(u.String(), tempDir)
	if err != nil {
		return err
	}

	versions, err := vcs.GetGitRepoVersions(r)
	if err != nil {
		return err
	}

	if len(versions) == 0 {
		logger.Warnf("no versions found for %s", repo.Name)
		return nil
	}

	// Get icon image name.
	iconFile, err := getGitRepoIcon(r)
	if err != nil {
		logger.Errorf("failed to get icon url: %v", err)
		return err
	}

	icon, err := GetRepoFileRaw(repo, iconFile, c)
	if err != nil {
		return err
	}

	entity := &model.Template{
		Name:        repo.Name,
		Source:      repo.Link,
		Description: repo.Description,
		Icon:        icon,
	}
	if c != nil {
		entity.CatalogID = c.ID
	}
	// Create template.
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

	return CreateTemplateVersionsFromRepo(ctx, mc, entity, r, versions)
}

// GetTemplateVersions retrieves template versions from a git repository.
// It will save images to the database if they are found in the repository.
func GetTemplateVersions(
	entity *model.Template,
	w *git.Worktree,
	versions []*version.Version,
) (model.TemplateVersions, error) {
	var (
		logger = log.WithName("catalog")
		tvs    = make(model.TemplateVersions, 0, len(versions))
	)

	u, err := transport.NewEndpoint(entity.Source)
	if err != nil {
		return nil, err
	}

	source := u.Host + u.Path

	for i := range versions {
		v := versions[i]
		tag := v.Original()
		err := w.Reset(&git.ResetOptions{
			Commit: plumbing.NewHash(tag),
			Mode:   git.HardReset,
		})
		if err != nil {
			logger.Warnf("failed to reset to tag %s: %v", tag, err)
			continue
		}

		dir := w.Filesystem.Root()

		schema, err := loadTerraformTemplateSchema(dir)
		if err != nil {
			logger.Warnf("failed to load terraform template schema: %v", err)
			continue
		}

		if schema == nil {
			logger.Warnf("terraform template schema is nil")
			continue
		}

		tvs = append(tvs, &model.TemplateVersion{
			TemplateID:   entity.ID,
			TemplateName: entity.Name,
			Version:      v.String(),
			Source:       source + "?ref=" + tag,
			Schema:       schema,
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
func GetRepoFileRaw(repo *vcs.Repository, file string, c *model.Catalog) (string, error) {
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

	if c != nil && c.Type == gitlab.Driver {
		return fmt.Sprintf("https://%s/%s/%s/-/raw/%s/%s", endpoint.Host, repo.Namespace, repo.Name, ref, file), nil
	}

	return "", nil
}