package server

import (
	"context"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/setting"
	"github.com/seal-io/walrus/pkg/settings"
)

// setupSettings creates the global settings into the database,
// it must be the first step to be executed at initialization.
func (r *Server) setupSettings(ctx context.Context, opts initOptions) error {
	return opts.ModelClient.Settings().CreateBulk().
		Set(settings.All()...).
		OnConflictColumns(setting.FieldName).
		Update(func(upsert *model.SettingUpsert) {
			upsert.UpdateHidden()
			upsert.UpdateEditable()
			upsert.UpdateSensitive()
			upsert.UpdatePrivate()
		}).
		Exec(ctx)
}
