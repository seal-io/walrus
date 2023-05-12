package connectors

import (
	"context"
	"fmt"
	"time"

	"github.com/seal-io/seal/pkg/costs/deployer"
	"github.com/seal-io/seal/pkg/costs/syncer"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/clustercost"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/platform"
	"github.com/seal-io/seal/pkg/platform/operator"
)

type StatusSyncer struct {
	client model.ClientSet
}

func NewStatusSyncer(client model.ClientSet) *StatusSyncer {
	return &StatusSyncer{
		client: client,
	}
}

// SyncStatus sync all connector status.
func (in *StatusSyncer) SyncStatus(ctx context.Context, conn *model.Connector) error {
	// Statuses check in sequence.
	statusChecker := []struct {
		status status.ConditionType
		check  func(context.Context, model.Connector) (string, error)
	}{
		{
			status: status.ConnectorStatusConnected,
			check: func(ctx context.Context, m model.Connector) (string, error) {
				return "", in.checkReachable(ctx, m)
			},
		},
		{
			status: status.ConnectorStatusCostToolsDeployed,
			check: func(ctx context.Context, m model.Connector) (string, error) {
				return "", in.checkCostTool(ctx, m)
			},
		},
		{
			status: status.ConnectorStatusCostSynced,
			check:  in.syncFinOpsData,
		},
		{
			status: status.ConnectorStatusReady,
			check: func(ctx context.Context, connector model.Connector) (string, error) {
				return "", nil
			},
		},
	}

	// Set status and message.
	for _, sc := range statusChecker {
		// Init with unknown.
		sc.status.Unknown(conn, "")

		successMsg, err := sc.check(ctx, *conn)
		if err != nil {
			sc.status.Status(conn, status.ConditionStatusFalse)
			sc.status.Message(conn, err.Error())
			break
		}
		sc.status.Status(conn, status.ConditionStatusTrue)
		sc.status.Message(conn, successMsg)
	}

	return UpdateStatus(ctx, in.client, conn)
}

// SyncFinOpsStatus only sync cost data.
func (in *StatusSyncer) SyncFinOpsStatus(ctx context.Context, conn *model.Connector) error {
	// Statuses check in sequence.
	statusChecker := []struct {
		status status.ConditionType
		check  func(context.Context, model.Connector) (string, error)
	}{
		{
			status: status.ConnectorStatusCostSynced,
			check:  in.syncFinOpsData,
		},
	}

	// Set status and message.
	for _, sc := range statusChecker {
		// Init with unknown.
		sc.status.Unknown(conn, "")

		successMsg, err := sc.check(ctx, *conn)
		if err != nil {
			sc.status.False(conn, err.Error())
			break
		}
		sc.status.True(conn, successMsg)
	}

	return UpdateStatus(ctx, in.client, conn)
}

func (in *StatusSyncer) syncFinOpsData(ctx context.Context, conn model.Connector) (string, error) {
	if !conn.EnableFinOps {
		return "", nil
	}

	if !status.ConnectorStatusReady.IsTrue(&conn) {
		// Skip connector isn't ready.
		return "", nil
	}

	k8sSyncer := syncer.NewK8sCostSyncer(in.client, nil)
	err := k8sSyncer.Sync(ctx, &conn, nil, nil)
	if err != nil {
		return "", err
	}

	existed, err := in.client.ClusterCosts().Query().
		Where(clustercost.ConnectorID(conn.ID)).
		Exist(ctx)
	if err != nil {
		return "", err
	}

	now := time.Now().UTC()
	if now.Sub(*conn.CreateTime) < time.Hour && !existed {
		return "It takes about an hour to generate hour-level cost data", nil
	}

	return fmt.Sprintf("Last sync time %s", now.Format(time.RFC3339)), nil
}

func (in *StatusSyncer) checkReachable(ctx context.Context, conn model.Connector) error {
	if conn.Type != types.ConnectorTypeK8s {
		return nil
	}

	op, err := platform.GetOperator(ctx, operator.CreateOptions{
		Connector: conn,
	})
	if err != nil {
		return fmt.Errorf("invalid connector config: %w", err)
	}
	if err = op.IsConnected(ctx); err != nil {
		return fmt.Errorf("unreachable connector: %w", err)
	}
	return nil
}

func (in *StatusSyncer) checkCostTool(ctx context.Context, conn model.Connector) error {
	if conn.Type != types.ConnectorTypeK8s || !conn.EnableFinOps {
		return nil
	}

	err := deployer.CostToolsStatus(ctx, &conn)
	if err != nil {
		return fmt.Errorf("error check cost tools: %w", err)
	}
	return nil
}

// UpdateStatus set summary and update the connector with locked.
func UpdateStatus(ctx context.Context, client model.ClientSet, conn *model.Connector) error {
	conn.Status.SetSummary(status.WalkConnector(&conn.Status))
	if !conn.Status.Changed() {
		return nil
	}
	return client.WithTx(ctx, func(tx *model.Tx) error {
		_, err := client.Connectors().Query().
			Where(connector.ID(conn.ID)).
			ForUpdate().
			Only(ctx)
		if err != nil {
			return err
		}

		update, err := dao.ConnectorUpdate(tx, conn)
		if err != nil {
			return err
		}
		return update.Exec(ctx)
	})
}
