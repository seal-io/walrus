package server

import (
	"context"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/setting"
	"github.com/seal-io/seal/pkg/settings"
)

func (r *Server) initSettings(ctx context.Context, opts initOptions) error {
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
