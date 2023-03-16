package scheduler

import (
	"context"
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
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/platformk8s"
	"github.com/seal-io/seal/utils/gopool"
	"github.com/seal-io/seal/utils/log"
	"github.com/seal-io/seal/utils/timex"
)

const (
	maxCollectTimeRange = 24 * time.Hour
	defaultStep         = 1 * time.Hour
)

type CostSyncTask struct {
	client model.ClientSet
	logger log.Logger
}

func NewCostSyncTask(client model.ClientSet) (*CostSyncTask, error) {
	return &CostSyncTask{
		client: client,
		logger: log.WithName("cost").WithName("sync-cost"),
	}, nil
}

func (in *CostSyncTask) Process(ctx context.Context, args ...interface{}) error {
	conns, err := in.client.Connectors().Query().Where(connector.TypeEQ(types.ConnectorTypeK8s)).All(ctx)
	if err != nil {
		return err
	}

	wg := gopool.Group()
	for i := range conns {
		var conn = conns[i]
		if !conn.EnableFinOps {
			continue
		}

		if conn.FinOpsStatus != status.Ready {
			in.logger.Debugf("connector %s status is:%s, skip collect cost data", conn.Name, conn.FinOpsStatus)
			continue
		}

		wg.Go(func() error {
			return in.SyncK8sCost(ctx, conn, args)
		})
	}

	return wg.Wait()
}

func (in *CostSyncTask) SyncK8sCost(ctx context.Context, conn *model.Connector, args []interface{}) error {
	in.logger.Debugf("collect cost for connector: %s", conn.Name)
	apiConfig, _, err := platformk8s.LoadApiConfig(*conn)
	if err != nil {
		return err
	}

	restCfg, err := platformk8s.GetConfig(*conn)
	if err != nil {
		return err
	}

	clusterName := apiConfig.CurrentContext
	collect, err := collector.NewCollector(restCfg, conn, clusterName)
	if err != nil {
		return err
	}

	startTime, endTime, err := in.timeRange(ctx, restCfg, conn, args)
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

func (in *CostSyncTask) batchCreateClusterCosts(ctx context.Context, costs []*model.ClusterCost) error {
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
	if err != nil {
		return fmt.Errorf("error batch create cluster costs: %w", err)
	}
	return nil
}

func (in *CostSyncTask) batchCreateAllocationCosts(ctx context.Context, costs []*model.AllocationCost) error {
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
	if err != nil {
		return fmt.Errorf("error batch create allocation costs: %w", err)
	}
	return nil
}

func (in *CostSyncTask) timeRange(ctx context.Context, restCfg *rest.Config, conn *model.Connector, args []interface{}) (*time.Time, *time.Time, error) {
	// time range from args.
	var startTime, endTime *time.Time
	for i, v := range args {
		switch i {
		case 1:
			s, ok := v.(*time.Time)
			if ok {
				startTime = s
			}
		case 2:
			s, ok := v.(*time.Time)
			if ok {
				endTime = s
			}
		}
	}

	if startTime != nil && endTime != nil {
		return startTime, endTime, nil
	}

	// time range from cluster.
	clientSet, err := kubernetes.NewForConfig(restCfg)
	if err != nil {
		return nil, nil, err
	}

	nodes, err := clientSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, nil, err
	}

	var clusterEarliestTime = time.Now()
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
