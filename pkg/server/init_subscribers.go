package server

import (
	"context"

	"github.com/seal-io/seal/pkg/connectors"
)

func (r *Server) initSubscribers(ctx context.Context, opts initOptions) error {
	err := connectors.AddSubscriber("connector-cost-tools-subscriber", connectors.EnsureCostTools)
	if err != nil {
		return err
	}

	err = connectors.AddSubscriber("connector-cost-custom-pricing-subscriber", connectors.SyncCostCustomPricing)
	if err != nil {
		return err
	}
	return nil
}
