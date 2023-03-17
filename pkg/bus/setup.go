package bus

import (
	"context"

	"github.com/seal-io/seal/pkg/bus/applicationrevision"
	"github.com/seal-io/seal/pkg/bus/connector"
	"github.com/seal-io/seal/pkg/bus/module"
	"github.com/seal-io/seal/pkg/bus/setting"
	"github.com/seal-io/seal/pkg/connectors"
	"github.com/seal-io/seal/pkg/cron"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/modules"
	"github.com/seal-io/seal/pkg/platformtf"
)

type SetupOptions struct {
	ModelClient model.ClientSet
}

func Setup(ctx context.Context, opts SetupOptions) (err error) {
	// application revision
	err = applicationrevision.AddSubscriber("terraform-sync-application-revision-status", platformtf.SyncApplicationRevisionStatus)
	if err != nil {
		return
	}

	// connector
	err = connector.AddSubscriber("connector-sync-status", connectors.StatusSync(opts.ModelClient).Do)
	if err != nil {
		return
	}
	err = connector.AddSubscriber("connector-sync-cost-custom-pricing", connectors.SyncCostCustomPricing)
	if err != nil {
		return
	}

	// module
	err = module.AddSubscriber("module-sync-schema", modules.SchemaSync(opts.ModelClient).Do)
	if err != nil {
		return
	}

	// setting
	err = setting.AddSubscriber("cron-sync", cron.Sync)
	if err != nil {
		return
	}
	return
}
