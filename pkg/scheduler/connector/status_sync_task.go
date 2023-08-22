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

func (in *StatusSyncTask) Process(ctx context.Context, args ...any) error {
	cs, err := in.modelClient.Connectors().Query().
		Where(
			connector.CategoryIn(
				types.ConnectorCategoryKubernetes,
				types.ConnectorCategoryCloudProvider)).
		All(ctx)
	if err != nil {
		return err
	}

	if len(cs) == 0 {
		return nil
	}

	var (
		s  = connectors.NewStatusSyncer(in.modelClient)
		wg = gopool.Group()
	)

	for i := range cs {
		c := cs[i]

		wg.Go(func() error {
			err := s.SyncStatus(ctx, c, false)
			if err != nil {
				return fmt.Errorf("error syncing connector %s: %w",
					c.ID, err)
			}

			return nil
		})
	}

	return wg.Wait()
}
