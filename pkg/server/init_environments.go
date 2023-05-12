package server

import (
	"context"
	"database/sql"
	"errors"

	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/environment"
)

func (r *Server) initEnvironments(ctx context.Context, opts initOptions) error {
	builtin := []*model.Environment{
		// Default environment.
		{
			Name:        "default",
			Description: "Default environment",
		},
	}

	creates, err := dao.EnvironmentCreates(opts.ModelClient, builtin...)
	if err != nil {
		return err
	}

	for i := range creates {
		err = creates[i].
			OnConflictColumns(environment.FieldName).
			DoNothing().
			Exec(ctx)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// No rows error is reasonable for nothing updating.
				continue
			}

			return err
		}
	}

	return nil
}
