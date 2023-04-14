package cost

import (
	"context"
	"sync"
	"time"

	"github.com/seal-io/seal/pkg/connectors"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/utils/gopool"
	"github.com/seal-io/seal/utils/log"
)

type CostSyncTask struct {
	mu sync.Mutex

	client model.ClientSet
	logger log.Logger
}

func NewCostSyncTask(client model.ClientSet) (*CostSyncTask, error) {
	return &CostSyncTask{
		client: client,
		logger: log.WithName("schedule-task").WithName("cost-sync"),
	}, nil
}

func (in *CostSyncTask) Process(ctx context.Context, args ...interface{}) error {
	if !in.mu.TryLock() {
		in.logger.Warn("previous processing is not finished")
		return nil
	}
	var startTs = time.Now()
	defer func() {
		in.mu.Unlock()
		in.logger.Debugf("processed in %v", time.Since(startTs))
	}()

	conns, err := in.client.Connectors().Query().Where(connector.TypeEQ(types.ConnectorTypeK8s)).All(ctx)
	if err != nil {
		return err
	}

	syncer := connectors.NewStatusSyncer(in.client)
	if err != nil {
		return err
	}

	wg := gopool.Group()
	for i := range conns {
		var conn = conns[i]
		if !conn.EnableFinOps {
			continue
		}

		wg.Go(func() error {
			return syncer.SyncFinOpsStatus(ctx, conn)
		})
	}

	return wg.Wait()
}
