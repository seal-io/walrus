package server

import (
	"context"

	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/setting"
	"github.com/seal-io/seal/pkg/settings"
)

func (r *Server) initSettings(ctx context.Context, opts initOptions) error {
	creates, err := dao.SettingCreates(opts.ModelClient, settings.All()...)
	if err != nil {
		return err
	}

	for i := range creates {
		err = creates[i].
			OnConflictColumns(setting.FieldName).
			Update(func(upsert *model.SettingUpsert) {
				upsert.UpdateHidden()
				upsert.UpdateEditable()
				upsert.UpdatePrivate()
			}).
			Exec(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
