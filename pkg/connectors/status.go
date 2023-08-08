package connectors

import (
	"context"
	"fmt"
	"time"

	"github.com/seal-io/seal/pkg/costs/deployer"
	"github.com/seal-io/seal/pkg/costs/syncer"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/model/costreport"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/operator"
	optypes "github.com/seal-io/seal/pkg/operator/types"
	"github.com/seal-io/seal/utils/slice"
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
		check  func(context.Context, model.Connector) (bool, string, error)
	}{
		{
			status: status.ConnectorStatusConnected,
			check: func(ctx context.Context, m model.Connector) (bool, string, error) {
				related, err := in.checkReachable(ctx, m)
				return related, "", err
			},
		},
		{
			status: status.ConnectorStatusCostToolsDeployed,
			check: func(ctx context.Context, m model.Connector) (bool, string, error) {
				related, err := in.checkCostTool(ctx, m)
				return related, "", err
			},
		},
		{
			status: status.ConnectorStatusCostSynced,
			check:  in.syncFinOpsData,
		},
		{
			status: status.ConnectorStatusReady,
			check: func(ctx context.Context, connector model.Connector) (bool, string, error) {
				return true, "", nil
			},
		},
	}

	// Set status and message.
	for _, sc := range statusChecker {
		related, successMsg, err := sc.check(ctx, *conn)
		if !related {
			continue
		}
		// Init with unknown.
		sc.status.Unknown(conn, "")

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
		check  func(context.Context, model.Connector) (bool, string, error)
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

		related, successMsg, err := sc.check(ctx, *conn)
		if !related {
			continue
		}

		if err != nil {
			sc.status.False(conn, err.Error())
			break
		}

		sc.status.True(conn, successMsg)
	}

	return UpdateStatus(ctx, in.client, conn)
}

func (in *StatusSyncer) syncFinOpsData(ctx context.Context, conn model.Connector) (bool, string, error) {
	if conn.Type != types.ConnectorTypeK8s {
		return false, "", nil
	}

	if !conn.EnableFinOps {
		return true, "", nil
	}

	if !status.ConnectorStatusReady.IsTrue(&conn) {
		// Skip connector isn't ready.
		return true, "", nil
	}

	k8sSyncer := syncer.NewK8sCostSyncer(in.client, nil)

	err := k8sSyncer.Sync(ctx, &conn, nil, nil)
	if err != nil {
		return true, "", err
	}

	existed, err := in.client.CostReports().Query().
		Where(costreport.ConnectorID(conn.ID)).
		Exist(ctx)
	if err != nil {
		return true, "", err
	}

	now := time.Now().UTC()
	if now.Sub(*conn.CreateTime) < time.Hour && !existed {
		return true, "It takes about an hour to generate hour-level cost data", nil
	}

	return true, fmt.Sprintf("Last sync time %s", now.Format(time.RFC3339)), nil
}

func (in *StatusSyncer) checkReachable(ctx context.Context, conn model.Connector) (bool, error) {
	if !slice.ContainsAny(
		[]string{
			types.ConnectorCategoryCloudProvider,
			types.ConnectorCategoryKubernetes,
		},
		conn.Category,
	) {
		return false, nil
	}

	op, err := operator.Get(ctx, optypes.CreateOptions{
		Connector: conn,
	})
	if err != nil {
		return true, fmt.Errorf("invalid connector config: %w", err)
	}

	if err = op.IsConnected(ctx); err != nil {
		return true, fmt.Errorf("unreachable connector: %w", err)
	}

	return true, nil
}

func (in *StatusSyncer) checkCostTool(ctx context.Context, conn model.Connector) (bool, error) {
	if conn.Type != types.ConnectorTypeK8s {
		return false, nil
	}

	if !conn.EnableFinOps {
		return true, nil
	}

	err := deployer.CostToolsStatus(ctx, &conn)
	if err != nil {
		return true, fmt.Errorf("error check cost tools: %w", err)
	}

	return true, nil
}

// UpdateStatus set summary and update the connector with locked.
func UpdateStatus(ctx context.Context, client model.ClientSet, conn *model.Connector) error {
	conn.Status.SetSummary(status.WalkConnector(&conn.Status))

	if !conn.Status.Changed() {
		return nil
	}

	return client.WithTx(ctx, func(tx *model.Tx) error {
		_, err := tx.Connectors().Query().
			Where(connector.ID(conn.ID)).
			Select(
				connector.FieldID,
				connector.FieldStatus).
			ForUpdate().
			Only(ctx)
		if err != nil {
			return err
		}

		return tx.Connectors().UpdateOne(conn).
			SetStatus(conn.Status).
			Exec(ctx)
	})
}
