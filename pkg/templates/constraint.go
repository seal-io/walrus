package templates

import (
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	"github.com/hashicorp/go-version"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/templates/loader"
	"github.com/seal-io/walrus/pkg/templates/openapi"
	"github.com/seal-io/walrus/pkg/vcs"
	"github.com/seal-io/walrus/utils/log"
	utilversion "github.com/seal-io/walrus/utils/version"
)

// getValidVersions get valid terraform module versions.
func getValidVersions(
	entity *model.Template,
	r *git.Repository,
	versions []*version.Version,
	subPath string,
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

		if subPath != "" {
			dir = filepath.Join(dir, subPath)
		}

		schema, err := loader.LoadSchemaPreferFile(dir, entity.Name)
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

		satisfy, err := isConstraintSatisfied(schema)
		if err != nil {
			logger.Warnf(
				"failed to check constraint satisfy \"%s:%s\" of catalog %q: %v",
				entity.Name,
				tag,
				entity.CatalogID,
				err,
			)

			continue
		}

		if !satisfy {
			logger.Infof(
				"\"%s:%s\" of catalog %q does not satisfy server version constraint",
				entity.Name,
				tag,
				entity.CatalogID,
			)

			continue
		}

		validVersions = append(validVersions, v)
		versionSchema[v] = *schema
	}

	return validVersions, versionSchema, nil
}

func isConstraintSatisfied(schema *types.TemplateVersionSchema) (bool, error) {
	v := utilversion.Version

	if utilversion.IsDevVersion() {
		return true, nil
	}

	if schema == nil || schema.OpenAPISchema == nil || schema.OpenAPISchema.Info == nil {
		return false, nil
	}

	semv, err := semver.NewVersion(strings.TrimPrefix(v, "v"))
	if err != nil {
		return false, err
	}

	semc := openapi.GetExtWalrusVersion(schema.OpenAPISchema.Info.Extensions)
	if semc == "" {
		return true, nil
	}

	semtc, err := semver.NewConstraint(semc)
	if err != nil {
		return false, err
	}

	return semtc.Check(semv), nil
}
