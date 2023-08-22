package connector

import (
	"context"
	"fmt"

	"github.com/seal-io/walrus/pkg/connectors"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/connector"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/utils/gopool"
	"github.com/seal-io/walrus/utils/log"
)

type StatusSyncTask struct {
	logger      log.Logger
	modelClient model.ClientSet
}

func NewStatusSyncTask(logger log.Logger, mc model.ClientSet) (in *StatusSyncTask, err error) {
	in = &StatusSyncTask{
		logger:      logger,
		modelClient: mc,
	}

	return
}

func (in *StatusSyncTask) Name() string {
	return "connector-status-sync"
}

func (in *StatusSyncTask) Process(ctx context.Context, args ...any) error {
	conns, err := in.modelClient.Connectors().Query().Where(
		connector.CategoryIn(
			types.ConnectorCategoryKubernetes,
			types.ConnectorCategoryCloudProvider,
		),
	).All(ctx)
	if err != nil {
		return err
	}

	var (
		syncer = connectors.NewStatusSyncer(in.modelClient)
		wg     = gopool.Group()
	)

	for i := range conns {
		conn := conns[i]
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
