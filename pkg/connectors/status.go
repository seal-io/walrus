package connectors

import (
	"context"
	"errors"
	"fmt"

	"github.com/seal-io/seal/pkg/costs/deployer"
	"github.com/seal-io/seal/pkg/costs/syncer"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
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

// SyncStatus sync all connector status
func (in *StatusSyncer) SyncStatus(ctx context.Context, conn *model.Connector) error {
	// statuses check in sequence
	statusChecker := []struct {
		status status.ConditionType
		check  func(context.Context, model.Connector) error
	}{
		{
			status: status.ConnectorStatusProvisioned,
			check:  in.checkReachable,
		},
		{
			status: status.ConnectorStatusToolsDeployed,
			check:  in.checkCostTool,
		},
		{
			status: status.ConnectorStatusCostSynced,
			check:  in.syncFinOpsData,
		},
		{
			status: status.ConnectorStatusReady,
			check: func(ctx context.Context, connector model.Connector) error {
				return nil
			},
		},
	}

	// init with unknown
	for _, sc := range statusChecker {
		sc.status.Unknown(conn, "")
	}

	// set status and message
	for _, sc := range statusChecker {
		err := sc.check(ctx, *conn)
		if err != nil {
			sc.status.Status(conn, status.ConditionStatusFalse)
			sc.status.Message(conn, err.Error())
			break
		}
		sc.status.Status(conn, status.ConditionStatusTrue)
		sc.status.Message(conn, "")
	}

	return UpdateStatus(ctx, in.client, conn)
}

// SyncFinOpsStatus only sync cost data
func (in *StatusSyncer) SyncFinOpsStatus(ctx context.Context, conn *model.Connector) error {
	// statuses check in sequence
	statusChecker := []struct {
		status status.ConditionType
		check  func(context.Context, model.Connector) error
	}{
		{
			status: status.ConnectorStatusCostSynced,
			check:  in.syncFinOpsData,
		},
	}

	// init with unknown
	for _, sc := range statusChecker {
		sc.status.Unknown(conn, "")
	}

	// set status and message
	for _, sc := range statusChecker {
		err := sc.check(ctx, *conn)
		if err != nil {
			sc.status.False(conn, err.Error())
			break
		}
		sc.status.True(conn, "")
	}

	return UpdateStatus(ctx, in.client, conn)
}

func (in *StatusSyncer) syncFinOpsData(ctx context.Context, conn model.Connector) error {
	if !conn.EnableFinOps {
		return nil
	}

	if !status.ConnectorStatusReady.IsTrue(&conn) {
		// skip connector isn't ready
		return nil
	}

	k8sSyncer := syncer.NewK8sCostSyncer(in.client, nil)
	return k8sSyncer.Sync(ctx, &conn, nil, nil)
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

	connected, err := op.IsConnected(ctx)
	if err != nil {
		return fmt.Errorf("invalid connector: %w", err)
	}
	if !connected {
		return errors.New("invalid connector: unreachable")
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

// UpdateStatus set summary and update the connector with locked
func UpdateStatus(ctx context.Context, client model.ClientSet, conn *model.Connector) error {
	if !conn.Status.ConditionChanged() {
		return nil
	}
	StatusSummarizer.SetSummarize(&conn.Status)
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
