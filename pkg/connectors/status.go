package connectors

import (
	"context"
	"errors"
	"fmt"

	"github.com/seal-io/seal/pkg/costs/deployer"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/platform"
	"github.com/seal-io/seal/pkg/platform/operator"
)

type StatusChecker struct {
	client model.ClientSet
}

func NewStatusChecker(client model.ClientSet) *StatusChecker {
	return &StatusChecker{
		client: client,
	}
}

func (in *StatusChecker) CheckStatus(ctx context.Context, conn *model.Connector) (string, error) {
	// statuses check in sequence
	statusChecker := []struct {
		status string
		check  func(context.Context, model.Connector) error
	}{
		{
			status: status.ConnectorStatusError,
			check:  in.checkReachable,
		},
		{
			status: status.ConnectorStatusDeploying,
			check:  in.checkCostTool,
		},
	}

	for _, sc := range statusChecker {
		if err := sc.check(ctx, *conn); err != nil {
			return sc.status, err
		}
	}

	// return ready while no error existed
	return status.ConnectorStatusReady, nil
}

func (in *StatusChecker) checkReachable(ctx context.Context, conn model.Connector) error {
	if conn.Type != types.ConnectorTypeK8s {
		return nil
	}

	op, err := platform.GetOperator(ctx, operator.CreateOptions{
		Connector: conn,
	})
	if err != nil {
		return fmt.Errorf("invalid connector config: %w", err)
	}

	connected, err := op.IsConnected(ctx)
	if err != nil {
		return fmt.Errorf("invalid connector: %w", err)
	}
	if !connected {
		return errors.New("invalid connector: unreachable")
	}
	return nil
}

func (in *StatusChecker) checkCostTool(ctx context.Context, conn model.Connector) error {
	if conn.Type != types.ConnectorTypeK8s {
		return nil
	}

	err := deployer.CostToolsStatus(ctx, &conn)
	if err != nil {
		return fmt.Errorf("error check cost tools: %w", err)
	}
	return nil
}
