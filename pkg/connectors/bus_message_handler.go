package connectors

import (
	"context"
	"time"

	"github.com/seal-io/seal/pkg/bus/connector"
	"github.com/seal-io/seal/pkg/costs/deployer"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/utils/gopool"
	"github.com/seal-io/seal/utils/log"
)

const (
	timeout = 3 * time.Minute
)

// SyncStatus receive the bus message and install the cost tools base on message.
func SyncStatus(ctx context.Context, message connector.BusMessage) error {
	logger := log.WithName("cost")

	client := message.ModelClient
	conn, err := client.Connectors().Get(ctx, message.Refer.ID)
	if err != nil {
		return err
	}

	gopool.Go(func() {
		subCtx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// ensure cost tools
		if conn.EnableFinOps {
			logger.Debugf("ensuring cost tools for connector %s", conn.Name)

			// set transition status
			status.ConnectorStatusToolsDeployed.Unknown(conn, "Deploying cost tools")
			if err = UpdateStatus(ctx, client, conn); err != nil {
				logger.Errorf("error update connector %s status: %v", conn.Name, err)
				return
			}

			// deploy
			err = deployer.DeployCostTools(subCtx, conn, message.Replace)
			if err != nil {
				// log instead of return error, then continue to sync the final status to connector
				logger.Errorf("error ensure cost tools for connector %s: %v", conn.Name, err)
			}
		}

		// check and generate final status
		syncer := NewStatusSyncer(client)
		err = syncer.SyncStatus(subCtx, conn)
		if err != nil {
			logger.Errorf("error sync status for connector %s: %v", conn.Name, err)
		}

		// sync cost data and generate status
		err = syncer.SyncFinOpsStatus(subCtx, conn)
		if err != nil {
			logger.Errorf("error sync finOps status for connector %s: %v", conn.Name, err)
		}
	})
	return nil
}

// SyncCostCustomPricing receive bus message and update custom pricing base on the message.
func SyncCostCustomPricing(ctx context.Context, message connector.BusMessage) error {
	conn := message.Refer
	if !conn.EnableFinOps {
		return nil
	}

	return deployer.UpdateCustomPricing(ctx, conn)
}
