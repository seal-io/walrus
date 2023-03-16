package connector

import (
	"context"
	"fmt"

	"github.com/seal-io/seal/pkg/connectors"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/utils/gopool"
	"github.com/seal-io/seal/utils/log"
)

type StatusSyncTask struct {
	client model.ClientSet
	logger log.Logger
}

func NewStatusCheckTask(client model.ClientSet) (*StatusSyncTask, error) {
	return &StatusSyncTask{
		client: client,
		logger: log.WithName("schedule-task").WithName("status-sync"),
	}, nil
}

func (in *StatusSyncTask) Process(ctx context.Context, args ...interface{}) error {
	conns, err := in.client.Connectors().Query().Where(connector.TypeEQ(types.ConnectorTypeK8s)).All(ctx)
	if err != nil {
		return err
	}

	var (
		syncer = connectors.NewStatusSyncer(in.client)
		wg     = gopool.Group()
	)
	for i := range conns {
		var conn = conns[i]
		in.logger.Debugf("sync status for connector: %s", conn.Name)
		wg.Go(func() error {
			if err := syncer.SyncStatus(ctx, conn); err != nil {
				return fmt.Errorf("error sync connector %s: %w", conn.Name, err)
			}
			return nil
		})
	}
	return wg.Wait()
}
