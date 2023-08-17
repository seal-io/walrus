package server

import (
	"context"

	"entgo.io/ent/dialect/sql/schema"

	"github.com/seal-io/walrus/pkg/dao/migration"
	"github.com/seal-io/walrus/pkg/datalisten/modelchange"
)

// applyModelSchemas creates the model schemas into the database,
// it must be the first step to be executed at any initializations.
func (r *Server) applyModelSchemas(ctx context.Context, opts initOptions) error {
	var createOpts []schema.MigrateOption

	// Apply subscription trigger to the listenable tables.
	createOpts = append(createOpts,
		migration.ApplyModelChangeTrigger(modelchange.TableNames())...)

	return opts.ModelClient.Schema.Create(ctx, createOpts...)
}
