package server

import (
	"context"
	"database/sql"
	"errors"

	tmplbus "github.com/seal-io/seal/pkg/bus/template"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model/template"
	"github.com/seal-io/seal/pkg/templates"
)

func (r *Server) initTemplates(ctx context.Context, opts initOptions) error {
	builtin, err := templates.BuiltinTemplates(ctx, opts.ModelClient)
	if err != nil {
		return err
	}

	creates, err := dao.TemplateCreates(opts.ModelClient, builtin...)
	if err != nil {
		return err
	}

	for i := range creates {
		// Do nothing if the template has been created.
		err = creates[i].
			OnConflictColumns(template.FieldID).
			DoNothing().
			Exec(ctx)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// No rows error is reasonable for nothing updating.
				continue
			}

			return err
		}

		err = tmplbus.Notify(ctx, builtin[i])
		if err != nil {
			return err
		}
	}

	return nil
}
