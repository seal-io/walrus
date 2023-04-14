package cost

import (
	"context"

	"github.com/seal-io/seal/pkg/connectors"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/utils/gopool"
	"github.com/seal-io/seal/utils/log"
)

type CollectTask struct {
	client model.ClientSet
	logger log.Logger
}

func NewCollectTask(client model.ClientSet) (*CollectTask, error) {
	return &CollectTask{
		client: client,
		logger: log.WithName("task").WithName("cost-collect"),
	}, nil
}

func (in *CollectTask) Process(ctx context.Context, args ...interface{}) error {
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
