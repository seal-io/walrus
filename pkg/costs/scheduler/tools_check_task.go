package scheduler

import (
	"context"
	"fmt"

	"github.com/seal-io/seal/pkg/costs/deployer"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/utils/gopool"
	"github.com/seal-io/seal/utils/log"
)

type ToolsCheckTask struct {
	client model.ClientSet
	logger log.Logger
}

func NewToolsCheckTask(client model.ClientSet) (*ToolsCheckTask, error) {
	return &ToolsCheckTask{
		client: client,
		logger: log.WithName("cost").WithName("check-status"),
	}, nil
}

func (in *ToolsCheckTask) Process(ctx context.Context, args ...interface{}) error {
	conns, err := in.client.Connectors().Query().Where(connector.TypeEQ(types.ConnectorTypeK8s)).All(ctx)
	if err != nil {
		return err
	}

	wg := gopool.Group()
	for i := range conns {
		var conn = conns[i]
		if !conn.EnableFinOps {
			continue
		}

		in.logger.Debugf("check cost tools status for connector: %s", conn.Name)
		wg.Go(func() error {
			var (
				st  = status.Ready
				msg string
			)
			err := deployer.CostToolsStatus(ctx, conn)
			if err != nil {
				st = status.Error
				msg = err.Error()
				in.logger.Errorf("error check cost tools for connector %s: %v", conn.Name, err)
			}

			updateErr := in.client.Connectors().
				UpdateOneID(conn.ID).
				SetFinOpsStatus(st).
				SetFinOpsStatusMessage(msg).
				Exec(ctx)
			if updateErr != nil {
				return fmt.Errorf("error update connector %s: %w", conn.Name, updateErr)
			}
			return nil
		})
	}
	return wg.Wait()
}
