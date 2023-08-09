package server

import (
	"context"

	buscatalog "github.com/seal-io/seal/pkg/bus/catalog"
	pkgcatalog "github.com/seal-io/seal/pkg/catalog"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/catalog"
	"github.com/seal-io/seal/pkg/dao/types/status"
)

func (r *Server) initCatalog(ctx context.Context, opts initOptions) error {
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
