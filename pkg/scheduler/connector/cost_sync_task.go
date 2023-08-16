package connector

import (
	"context"
	"sync"
	"time"

	"github.com/seal-io/walrus/pkg/connectors"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/connector"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/utils/gopool"
	"github.com/seal-io/walrus/utils/log"
)

type CollectTask struct {
	mu sync.Mutex

	modelClient model.ClientSet
	logger      log.Logger
}

func NewCollectTask(mc model.ClientSet) (*CollectTask, error) {
	in := &CollectTask{}
	in.modelClient = mc
	in.logger = log.WithName("task").WithName(in.Name())

	return in, nil
}

func (in *CollectTask) Name() string {
	return "connector-cost-collect"
}

func (in *CollectTask) Process(ctx context.Context, args ...any) error {
	if !in.mu.TryLock() {
		in.logger.Warn("previous processing is not finished")
		return nil
	}
	startTs := time.Now()

	defer func() {
		in.mu.Unlock()
		in.logger.Debugf("processed in %v", time.Since(startTs))
	}()

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
