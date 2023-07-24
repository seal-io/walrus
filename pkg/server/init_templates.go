package server

import (
	"context"

	tmplbus "github.com/seal-io/seal/pkg/bus/template"
	"github.com/seal-io/seal/pkg/dao/model/template"
	"github.com/seal-io/seal/pkg/templates"
)

func (r *Server) initTemplates(ctx context.Context, opts initOptions) error {
	builtin, err := templates.BuiltinTemplates(ctx, opts.ModelClient)
	if err != nil {
		return err
	}

	err = opts.ModelClient.Templates().CreateBulk().
		Set(builtin...).
		OnConflictColumns(template.FieldID).
		DoNothing().
		Exec(ctx)
	if err != nil {
		return err
	}

	for i := range builtin {
		err = tmplbus.Notify(ctx, builtin[i])
		if err != nil {
			return err
		}
	}

	return nil
}
