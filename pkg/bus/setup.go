package bus

import (
	"context"

	"github.com/seal-io/seal/pkg/bus/applicationrevision"
	"github.com/seal-io/seal/pkg/bus/module"
	"github.com/seal-io/seal/pkg/bus/setting"
	"github.com/seal-io/seal/pkg/cron"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/modules"
	"github.com/seal-io/seal/pkg/platformtf"
)

type SetupOptions struct {
	ModelClient model.ClientSet
}

func Setup(ctx context.Context, opts SetupOptions) (err error) {
	// Application revision.
	err = applicationrevision.AddSubscriber("terraform-sync-application-revision-status",
		platformtf.SyncApplicationRevisionStatus)
	if err != nil {
		return
	}

	// Module.
	err = module.AddSubscriber("module-sync-schema", modules.SchemaSync(opts.ModelClient).Do)
	if err != nil {
		return
	}

	// Setting.
	err = setting.AddSubscriber("cron-sync", cron.Sync)
	if err != nil {
		return
	}

	return
}
