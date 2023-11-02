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
	cs, err := in.modelClient.Connectors().Query().
		Where(
			connector.TypeEQ(types.ConnectorTypeKubernetes),
			connector.EnableFinOps(true)).
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
			err := s.SyncFinOpsStatus(ctx, c)
			if err != nil {
				return fmt.Errorf("error syncing cost of connector %s: %w",
					c.ID, err)
			}

			return nil
		})
	}

	return wg.Wait()
}
