package server

import (
	"context"
	"database/sql"
	"errors"

	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/project"
)

func (r *Server) initProjects(ctx context.Context, opts initOptions) error {
	var builtin = []*model.Project{
		// default project
		{
			Name:        "Default",
			Description: "Default project",
		},
	}

	var creates, err = dao.ProjectCreates(opts.ModelClient, builtin...)
	if err != nil {
		return err
	}
	for i := range creates {
		// do nothing if the project has been created.
		err = creates[i].
			OnConflictColumns(project.FieldName).
			DoNothing().
			Exec(ctx)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// no rows error is reasonable for nothing updating.
				continue
			}
			return err
		}
	}
	return nil
}
