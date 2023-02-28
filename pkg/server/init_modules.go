package server

import (
	"context"

	"entgo.io/ent/dialect/sql"

	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/module"
	"github.com/seal-io/seal/pkg/modules"
)

func (r *Server) initModules(ctx context.Context, opts initOptions) error {
	if err := modules.AddSubscriber("sync-module-schema-handler", modules.SyncSchema); err != nil {
		return err
	}

	var builtin = []*model.Module{
		{
			ID:          "webservice",
			Description: "A long-running, scalable, containerized service that have a stable network endpoint to receive external network traffic.",
			Source:      "github.com/gitlawr/modules/webservice",
		},
		{
			ID:          "aws-rds",
			Description: "An AWS RDS instance.",
			Source:      "github.com/gitlawr/modules/aws-rds",
		},
		{
			ID:          "mysql",
			Description: "A containerized mysql instance.",
			Source:      "github.com/gitlawr/modules/mysql",
		},
	}

	var creates, err = dao.ModuleCreates(opts.ModelClient, builtin...)
	if err != nil {
		return err
	}
	for i := range creates {
		err = creates[i].
			OnConflict(
				sql.ConflictColumns(
					module.FieldID,
				),
			).
			Update(func(upsert *model.ModuleUpsert) {
				upsert.UpdateDescription()
				upsert.UpdateSource()
				upsert.UpdateVersion()
			}).
			Exec(ctx)
		if err != nil {
			return err
		}

		modules.Notify(ctx, opts.ModelClient, builtin[i])
	}
	return nil
}
