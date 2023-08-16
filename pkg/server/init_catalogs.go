package server

import (
	"context"

	buscatalog "github.com/seal-io/walrus/pkg/bus/catalog"
	pkgcatalog "github.com/seal-io/walrus/pkg/catalog"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/catalog"
	"github.com/seal-io/walrus/pkg/dao/types/status"
)

// createBuiltinCatalogs creates the built-in Catalog resources.
func (r *Server) createBuiltinCatalogs(ctx context.Context, opts initOptions) error {
	builtin := pkgcatalog.BuiltinCatalog()

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
	builtin.Status.SetSummary(status.WalkService(&builtin.Status))

	c, err = opts.ModelClient.Catalogs().Create().
		Set(builtin).
		Save(ctx)
	if err != nil {
		return err
	}

	return buscatalog.Notify(ctx, opts.ModelClient, c)
}
