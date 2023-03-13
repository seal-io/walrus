package cost

import (
	"context"

	"github.com/seal-io/seal/pkg/costs/syncer"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/utils/gopool"
	"github.com/seal-io/seal/utils/log"
)

type CostSyncTask struct {
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
	conns, err := in.client.Connectors().Query().Where(connector.TypeEQ(types.ConnectorTypeK8s)).All(ctx)
	if err != nil {
		return err
	}

	// init cost syncer
	k8sCostSyncer, err := syncer.NewK8sCostSyncer(in.client)
	if err != nil {
		return err
	}
	k8sCostSyncer.SetLogger(in.logger)

	wg := gopool.Group()
	for i := range conns {
		var conn = conns[i]
		if !conn.EnableFinOps {
			continue
		}

		if conn.Status != status.ConnectorStatusReady {
			in.logger.Debugf("connector %s status is:%s, skip collect cost data", conn.Name, conn.Status)
			continue
		}

		wg.Go(func() error {
			return k8sCostSyncer.Sync(ctx, conn, nil, nil)
		})
	}

	return wg.Wait()
}
