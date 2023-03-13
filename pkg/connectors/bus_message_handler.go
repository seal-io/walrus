package connectors

import (
	"context"
	"fmt"
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

// EnsureCostTools receive the bus message and install the cost tools base on message.
func EnsureCostTools(ctx context.Context, message connector.BusMessage) error {
	logger := log.WithName("cost")

	conn, err := message.ModelClient.Connectors().Get(ctx, message.Refer.ID)
	if err != nil {
		return err
	}

	gopool.Go(func() {
		subCtx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		var (
			st  = conn.Status
			msg = conn.StatusMessage
			err error
		)

		// ensure cost tools
		if conn.EnableFinOps {
			logger.Debugf("ensuring cost tools for connector %s", conn.Name)

			err = deployer.DeployCostTools(subCtx, conn, message.Replace)
			if err != nil {
				st = status.ConnectorStatusDeploying
				msg = fmt.Sprintf("error ensure cost tools: %v", err)
				logger.Errorf("error ensure cost tools for connector %s: %v", conn.Name, err)
			}
		}

		// sync connector status
		if err == nil {
			// check and generate connector final status
			checker := NewStatusChecker(message.ModelClient)
			st, err = checker.CheckStatus(ctx, conn)
			if err != nil {
				msg = err.Error()
				logger.Errorf("error check status for connector %s: %v", conn.Name, err)
			}
		}

		updateErr := message.ModelClient.Connectors().
			UpdateOneID(conn.ID).
			SetStatus(st).
			SetStatusMessage(msg).
			Exec(ctx)
		if updateErr != nil {
			logger.Errorf("failed to update connector %s: %v", conn.Name, updateErr)
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
