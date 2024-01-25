package server

import (
	"context"
	"fmt"

	buscatalog "github.com/seal-io/walrus/pkg/bus/catalog"
	pkgcatalog "github.com/seal-io/walrus/pkg/catalog"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/catalog"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/settings"
)

// createBuiltinCatalogs creates the built-in Catalog resources.
func (r *Server) createBuiltinCatalogs(ctx context.Context, opts initOptions) error {
	enableBuiltinCatalog := settings.EnableBuiltinCatalog.ShouldValueBool(ctx, opts.ModelClient)
	if !enableBuiltinCatalog {
		return nil
	}

	builtin := pkgcatalog.BuiltinCatalog()

	switch opts.BuiltinCatalogProvider {
	case "github":
		builtin.Type = types.GitDriverGithub
	case "gitee":
		builtin.Type = types.GitDriverGitee
	default:
		return fmt.Errorf("invalid builtin catalog provider: %s", opts.BuiltinCatalogProvider)
	}

	c, err := opts.ModelClient.Catalogs().Query().
		Where(catalog.Name(builtin.Name)).
		Only(ctx)
	if err != nil && !model.IsNotFound(err) {
		return err
	}

	if c != nil {
		return nil
	}

	status.CatalogStatusInitialized.Unknown(builtin, "Initializing catalog template")
	builtin.Status.SetSummary(status.WalkResource(&builtin.Status))

	c, err = opts.ModelClient.Catalogs().Create().
		Set(builtin).
		Save(ctx)
	if err != nil {
		return err
	}

	return buscatalog.Notify(ctx, opts.ModelClient, c)
}
