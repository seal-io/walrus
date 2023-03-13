package connector

import (
	"context"
	"fmt"

	"github.com/seal-io/seal/pkg/connectors"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/utils/gopool"
	"github.com/seal-io/seal/utils/log"
)

type StatusCheckTask struct {
	client model.ClientSet
	logger log.Logger
}

func NewStatusCheckTask(client model.ClientSet) (*StatusCheckTask, error) {
	return &StatusCheckTask{
		client: client,
		logger: log.WithName("schedule-task").WithName("status-check"),
	}, nil
}

func (in *StatusCheckTask) Process(ctx context.Context, args ...interface{}) error {
	conns, err := in.client.Connectors().Query().Where(connector.TypeEQ(types.ConnectorTypeK8s)).All(ctx)
	if err != nil {
		return err
	}

	var (
		checker = connectors.NewStatusChecker(in.client)
		wg      = gopool.Group()
	)
	for i := range conns {
		var conn = conns[i]
		if !conn.EnableFinOps {
			continue
		}

		in.logger.Debugf("check cost tools status for connector: %s", conn.Name)
		wg.Go(func() error {
			var errMsg string
			connStatus, err := checker.CheckStatus(ctx, conn)
			if err != nil {
				errMsg = err.Error()
				in.logger.Errorf("error check connector %s status %s: %v", conn.Name, connStatus, err)
			}

			updateErr := in.client.Connectors().
				UpdateOneID(conn.ID).
				SetStatus(connStatus).
				SetStatusMessage(errMsg).
				Exec(ctx)
			if updateErr != nil {
				return fmt.Errorf("error update connector %s: %w", conn.Name, updateErr)
			}
			return nil
		})
	}
	return wg.Wait()
}
