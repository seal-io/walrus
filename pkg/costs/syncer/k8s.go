package syncer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/seal-io/seal/pkg/costs/collector"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/allocationcost"
	"github.com/seal-io/seal/pkg/dao/model/clustercost"
	"github.com/seal-io/seal/pkg/platformk8s"
	"github.com/seal-io/seal/utils/log"
	"github.com/seal-io/seal/utils/timex"
)

const (
	maxCollectTimeRange = 24 * time.Hour
	defaultStep         = 1 * time.Hour
)

type K8sCostSyncer struct {
	client model.ClientSet
	logger log.Logger
}

func NewK8sCostSyncer(client model.ClientSet, logger log.Logger) *K8sCostSyncer {
	if logger == nil {
		logger = log.WithName("cost")
	}
	return &K8sCostSyncer{
		client: client,
		logger: logger,
	}
}

func (in *K8sCostSyncer) SetLogger(logger log.Logger) {
	in.logger = logger
}

func (in *K8sCostSyncer) Sync(ctx context.Context, conn *model.Connector, startTime, endTime *time.Time) error {
	err := in.syncCost(ctx, conn, startTime, endTime)
	return err
}

func (in *K8sCostSyncer) syncCost(ctx context.Context, conn *model.Connector, startTime, endTime *time.Time) error {
	in.logger.Debugf("collect cost for connector: %s", conn.Name)
	apiConfig, _, err := platformk8s.LoadApiConfig(*conn)
	if err != nil {
		return err
	}

	// NB(thxCode): disable timeout as we don't know the maximum time-cost of once full-series costs synchronization,
	// and rely on the session context timeout control,
	// which means we don't close the underlay kubernetes client operation until the `ctx` is cancel or timeout.
	restCfg, err := platformk8s.GetConfig(*conn, platformk8s.WithoutTimeout())
	if err != nil {
		return err
	}

	clusterName := apiConfig.CurrentContext
	collect, err := collector.NewCollector(restCfg, conn, clusterName)
	if err != nil {
		return err
	}

	startTime, endTime, err = in.timeRange(ctx, restCfg, conn, startTime, endTime)
	if err != nil {
		return err
	}
	in.logger.Debugf("connector: %s, current sync costs within %s, %s", conn.Name, startTime, endTime)

	curTimeRange := endTime.Sub(*startTime)
	maxTimeRange := maxCollectTimeRange
	if curTimeRange < maxTimeRange {
		maxTimeRange = curTimeRange
	}

	stepStart := *startTime
	for endTime.After(stepStart) {
		stepEnd := stepStart.Add(maxTimeRange)
		in.logger.Debugf("connector: %s, step sync within %s, %s", conn.Name, stepStart.String(), stepEnd.String())

		cc, ac, err := collect.K8sCosts(&stepStart, &stepEnd, defaultStep)
		if err != nil {
			return err
		}

		if len(cc) == 0 {
			stepStart = stepEnd
			continue
		}

		if err = in.client.WithTx(ctx, func(tx *model.Tx) error {
			if err = in.batchCreateClusterCosts(ctx, cc); err != nil {
				return err
			}
			return in.batchCreateAllocationCosts(ctx, ac)
		}); err != nil {
			return err
		}

		in.logger.Debugf("create %d clusterCosts, %d allocationResourceCosts for connector:%s, within %s, %s",
			len(cc), len(ac), conn.Name, stepStart.String(), stepEnd.String(),
		)
		stepStart = stepEnd
	}
	return nil
}

func (in *K8sCostSyncer) batchCreateClusterCosts(ctx context.Context, costs []*model.ClusterCost) error {
	creates, err := dao.ClusterCostCreates(in.client, costs...)
	if err != nil {
		return err
	}

	err = in.client.ClusterCosts().CreateBulk(creates...).
		OnConflictColumns(
			clustercost.FieldStartTime,
			clustercost.FieldEndTime,
			clustercost.FieldConnectorID,
		).
		DoNothing().
		Exec(ctx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("error batch create cluster costs: %w", err)
	}
	return nil
}

func (in *K8sCostSyncer) batchCreateAllocationCosts(ctx context.Context, costs []*model.AllocationCost) error {
	creates, err := dao.AllocationCostCreates(in.client, costs...)
	if err != nil {
		return err
	}

	err = in.client.AllocationCosts().CreateBulk(creates...).
		OnConflictColumns(
			allocationcost.FieldStartTime,
			allocationcost.FieldEndTime,
			allocationcost.FieldConnectorID,
			allocationcost.FieldFingerprint,
		).
		DoNothing().
		Exec(ctx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("error batch create allocation costs: %w", err)
	}
	return nil
}

func (in *K8sCostSyncer) timeRange(
	ctx context.Context,
	restCfg *rest.Config,
	conn *model.Connector,
	startTime,
	endTime *time.Time,
) (*time.Time, *time.Time, error) {
	// Time range existed.
	if startTime != nil && endTime != nil {
		return startTime, endTime, nil
	}

	// Time range from cluster.
	clientSet, err := kubernetes.NewForConfig(restCfg)
	if err != nil {
		return nil, nil, err
	}

	nodes, err := clientSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, nil, err
	}

	clusterEarliestTime := time.Now()
	for _, v := range nodes.Items {
		if v.CreationTimestamp.Time.Before(clusterEarliestTime) {
			clusterEarliestTime = v.CreationTimestamp.Time
		}
	}

	s := timex.StartTimeOfHour(clusterEarliestTime, time.UTC)
	e := timex.StartTimeOfHour(time.Now(), time.UTC)
	startTime = &s
	endTime = &e

	existed, err := in.client.ClusterCosts().Query().
		Where(clustercost.ConnectorID(conn.ID)).
		Order(model.Desc(clustercost.FieldEndTime)).
		First(ctx)
	if err != nil {
		if model.IsNotFound(err) {
			return startTime, endTime, nil
		}
		return nil, nil, err
	}

	return &existed.EndTime, endTime, nil
}
