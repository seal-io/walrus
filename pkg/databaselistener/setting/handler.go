package setting

import (
	"context"
	"fmt"

	settingbus "github.com/seal-io/walrus/pkg/bus/setting"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/database"
)

const ChannelName = "setting_update"

func Handler(ctx context.Context, client model.ClientSet, rd database.Record) error {
	oid := object.NewID(rd.RecordID)

	setting, err := client.Settings().Get(ctx, oid)
	if err != nil {
		return fmt.Errorf("error get setting %s: %w", oid, err)
	}

	err = settingbus.Notify(ctx, client, model.Settings{
		setting,
	})

	return err
}
