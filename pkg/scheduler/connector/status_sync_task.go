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

type StatusSyncTask struct {
	mu sync.Mutex

	modelClient model.ClientSet
	logger      log.Logger
}

func NewStatusSyncTask(mc model.ClientSet) (*StatusSyncTask, error) {
	in := &StatusSyncTask{}
	in.modelClient = mc
	in.logger = log.WithName("task").WithName(in.Name())
	return in, nil
}

func (in *StatusSyncTask) Name() string {
	return "connector-status-sync"
}

func (in *StatusSyncTask) Process(ctx context.Context, args ...interface{}) error {
	if !in.mu.TryLock() {
		in.logger.Warn("previous processing is not finished")
		return nil
	}
	var startTs = time.Now()
	defer func() {
		in.mu.Unlock()
		in.logger.Debugf("processed in %v", time.Since(startTs))
	}()

	conns, err := in.modelClient.Connectors().Query().Where(connector.TypeEQ(types.ConnectorTypeK8s)).All(ctx)
	if err != nil {
		return err
	}

	var (
		syncer = connectors.NewStatusSyncer(in.modelClient)
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
