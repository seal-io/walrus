package server

import (
	"context"

	"entgo.io/ent/dialect/sql"

	modbus "github.com/seal-io/seal/pkg/bus/module"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/module"
)

func (r *Server) initModules(ctx context.Context, opts initOptions) error {

	var builtin = []*model.Module{
		{
			ID:          "webservice",
			Description: "A long-running, scalable, containerized service that have a stable network endpoint to receive external network traffic.",
			Source:      "github.com/gitlawr/modules/webservice",
		},
		{
			ID:          "webservice-from-source",
			Description: "Build and run a containerized service from source code.",
			Source:      "github.com/gitlawr/modules/webservice-from-source",
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

		err = modbus.Notify(ctx, opts.ModelClient, builtin[i])
		if err != nil {
			return err
		}
	}
	return nil
}
