package bus

import (
	"context"

	"github.com/seal-io/seal/pkg/bus/connector"
	"github.com/seal-io/seal/pkg/bus/module"
	"github.com/seal-io/seal/pkg/dao/model"
)

type SetupOptions struct {
	ModelClient model.ClientSet
}

func Setup(ctx context.Context, opts SetupOptions) error {
	if err := module.AddSubscriber("sync-module-schema-handler", module.SyncSchema); err != nil {
		return err
	}

	if err := connector.AddSubscriber("connector-cost-tools-subscriber", connector.EnsureCostTools); err != nil {
		return err
	}

	if err := connector.AddSubscriber("connector-cost-custom-pricing-subscriber", connector.SyncCostCustomPricing); err != nil {
		return err
	}

	return nil
}
