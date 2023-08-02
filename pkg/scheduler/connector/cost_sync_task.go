package connector

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/seal-io/seal/pkg/connectors"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/utils/gopool"
	"github.com/seal-io/seal/utils/log"
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

func (in *CollectTask) Process(ctx context.Context, args ...interface{}) error {
	if !in.mu.TryLock() {
		in.logger.Warn("previous processing is not finished")
		return nil
	}
	startTs := time.Now()

	defer func() {
		in.mu.Unlock()
		in.logger.Debugf("processed in %v", time.Since(startTs))
	}()

	conns, err := in.modelClient.Connectors().Query().
		Where(
			connector.TypeEQ(types.ConnectorTypeK8s),
			connector.EnableFinOps(true)).
		All(ctx)
	if err != nil {
		return err
	}

	if len(conns) == 0 {
		return nil
	}

	var (
		syncer = connectors.NewStatusSyncer(in.modelClient)
		wg     = gopool.Group()
	)

	for i := range conns {
		c := conns[i]

		in.logger.Debugf("syncing cost of connector %q", c.ID)

		wg.Go(func() error {
			if err := syncer.SyncFinOpsStatus(ctx, c); err != nil {
				return fmt.Errorf("error sync cost of connector %s: %w", c.ID, err)
			}

			return nil
		})
	}

	return wg.Wait()
}
