package connector

import (
	"context"

	"github.com/seal-io/walrus/pkg/connectors"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/connector"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/utils/gopool"
	"github.com/seal-io/walrus/utils/log"
)

type CollectTask struct {
	logger      log.Logger
	modelClient model.ClientSet
}

func NewCollectTask(logger log.Logger, mc model.ClientSet) (in *CollectTask, err error) {
	in = &CollectTask{
		logger:      logger,
		modelClient: mc,
	}

	return
}

func (in *CollectTask) Process(ctx context.Context, args ...any) error {
	conns, err := in.modelClient.Connectors().Query().Where(
		connector.TypeEQ(types.ConnectorTypeK8s)).All(ctx)
	if err != nil {
		return err
	}

	syncer := connectors.NewStatusSyncer(in.modelClient)

	if err != nil {
		return err
	}

	wg := gopool.Group()

	for i := range conns {
		conn := conns[i]
		if !conn.EnableFinOps {
			continue
		}

		wg.Go(func() error {
			return syncer.SyncFinOpsStatus(ctx, conn)
		})
	}

	return wg.Wait()
}
