package distributor

import (
	"context"
	"fmt"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"

	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/enttest"
	_ "github.com/seal-io/seal/pkg/dao/model/runtime"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/crypto"
	"github.com/seal-io/seal/utils/timex"
)

const (
	testDriverName     = "sqlite3"
	testDataSourceName = "file:ent?mode=memory&cache=shared&_fk=1"
)

const (
	testLabelApp = "seal/app"
	testLabelEnv = "seal/environment"

	// TestLabelApp replace the types.FilterFieldEnvironment = "seal.io/environment",
	// because the sqlite use JSON path expression, can't handle dot in json key.
	testFilterFieldEnv = types.FilterField(types.LabelPrefix + testLabelEnv)
	// TestLabelApp replace the types.GroupByFieldApplication.
	testGroupByFieldApp = types.GroupByField(types.LabelPrefix + testLabelApp)
)

func TestAccumulateDistribute(t *testing.T) {
	ctx := context.Background()

	client := enttest.Open(t, testDriverName, testDataSourceName)
	defer client.Close()

	var (
		startTime = timex.StartTimeOfHour(time.Now(), time.UTC)
		endTime   = startTime.Add(5 * time.Hour)
	)
	conn, err := testData(ctx, client, startTime, endTime)
	assert.Nil(t, err, "error create cost test data: %w", err)

	cases := []struct {
		name                 string
		inputStartTime       time.Time
		inputEndTime         time.Time
		inputCondition       types.QueryCondition
		outputTotalItemNum   int
		outputItemCost       float64
		outputItemSharedCost float64
	}{
		{
			name:           "time range with no data",
			inputStartTime: time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC),
			inputEndTime:   time.Date(2000, 2, 1, 1, 1, 1, 1, time.UTC),
			inputCondition: types.QueryCondition{
				GroupBy: types.GroupByFieldNamespace,
				SharedCosts: types.ShareCosts{
					{
						Filters: types.AllocationCostFilters{
							{
								{
									FieldName: types.FilterFieldNamespace,
									Operator:  types.OperatorIn,
									Values:    []string{"namespace-t1"},
								},
							},
						},
						SharingStrategy: types.SharingStrategyEqually,
					},
				},
			},
			outputTotalItemNum:   0,
			outputItemCost:       0,
			outputItemSharedCost: 0,
		},
		{
			name:           "equally share allocation cost with filter",
			inputStartTime: startTime,
			inputEndTime:   endTime,
			inputCondition: types.QueryCondition{
				GroupBy: types.GroupByFieldNamespace,
				SharedCosts: types.ShareCosts{
					{
						Filters: types.AllocationCostFilters{
							{
								{
									FieldName: types.FilterFieldNamespace,
									Operator:  types.OperatorIn,
									Values:    []string{"namespace-t1"},
								},
							},
						},
						SharingStrategy: types.SharingStrategyEqually,
					},
				},
			},
			outputTotalItemNum:   3,
			outputItemCost:       66.66666666666667,
			outputItemSharedCost: 16.66666666666667,
		},
		{
			name:           "equally share allocation cost with management and idle cost",
			inputStartTime: startTime,
			inputEndTime:   endTime,
			inputCondition: types.QueryCondition{
				GroupBy: types.GroupByFieldNamespace,
				SharedCosts: types.ShareCosts{
					{
						IdleCostFilters: types.IdleCostFilters{
							{
								ConnectorID: conn.ID,
							},
						},
						ManagementCostFilters: types.ManagementCostFilters{
							{
								ConnectorID: conn.ID,
							},
						},
						SharingStrategy: types.SharingStrategyEqually,
					},
				},
			},
			outputTotalItemNum:   3,
			outputItemCost:       166.66666666666669,
			outputItemSharedCost: 116.66666666666669,
		},
		{
			name:           "equally share allocation cost with filter, management and idle cost with filter",
			inputStartTime: startTime,
			inputEndTime:   endTime,
			inputCondition: types.QueryCondition{
				GroupBy: types.GroupByFieldNamespace,
				SharedCosts: types.ShareCosts{
					{
						Filters: types.AllocationCostFilters{
							{
								{
									FieldName: types.FilterFieldNamespace,
									Operator:  types.OperatorIn,
									Values:    []string{"namespace-t1"},
								},
							},
						},
						IdleCostFilters: types.IdleCostFilters{
							{
								ConnectorID: conn.ID,
							},
						},
						ManagementCostFilters: types.ManagementCostFilters{
							{
								ConnectorID: conn.ID,
							},
						},
						SharingStrategy: types.SharingStrategyEqually,
					},
				},
			},
			outputTotalItemNum:   3,
			outputItemCost:       183.33333333333333,
			outputItemSharedCost: 133.33333333333333,
		},
		{
			name:           "proportionally share allocation cost with filter",
			inputStartTime: startTime,
			inputEndTime:   endTime,
			inputCondition: types.QueryCondition{
				GroupBy: types.GroupByFieldNamespace,
				SharedCosts: types.ShareCosts{
					{
						Filters: types.AllocationCostFilters{
							{
								{
									FieldName: types.FilterFieldNamespace,
									Operator:  types.OperatorIn,
									Values:    []string{"namespace-t1"},
								},
							},
						},
						SharingStrategy: types.SharingStrategyProportionally,
					},
				},
			},
			outputTotalItemNum:   3,
			outputItemCost:       66.66666666666666,
			outputItemSharedCost: 16.66666666666667,
		},
	}

	for _, v := range cases {
		dsb := accumulateDistributor{client: client}
		items, count, err := dsb.distribute(ctx, v.inputStartTime, v.inputEndTime, v.inputCondition)
		assert.Equal(t, v.outputTotalItemNum, count, "%s: total item count mismatch", v.name)
		assert.Nil(t, err, "%s: error get distribute resource cost: %w", v.name, err)

		if len(items) != 0 {
			assert.Equal(t, v.outputItemCost, items[0].Cost.TotalCost,
				"%s: first item total cost mismatch", v.name)
		}
	}
}

func TestAllocationResourceCosts(t *testing.T) {
	// Init data.
	ctx := context.Background()

	client := enttest.Open(t, testDriverName, testDataSourceName)
	defer client.Close()

	var (
		startTime = timex.StartTimeOfHour(time.Now(), time.UTC)
		endTime   = startTime.Add(5 * time.Hour)
	)
	_, err := testData(ctx, client, startTime, endTime)
	assert.Nil(t, err, "error create cost test data: %w", err)

	cases := []struct {
		name                 string
		inputStartTime       time.Time
		inputEndTime         time.Time
		inputCondition       types.QueryCondition
		outputTotalItemNum   int
		outputQueriedItemNum int
		outputItemCost       float64
	}{
		{
			name:           "empty filters, group by namespace",
			inputStartTime: startTime,
			inputEndTime:   endTime,
			inputCondition: types.QueryCondition{
				GroupBy: types.GroupByFieldNamespace,
			},
			outputTotalItemNum:   3,
			outputQueriedItemNum: 3,
			outputItemCost:       50,
		},
		{
			name:           "filter by namespace, group by namespace",
			inputStartTime: startTime,
			inputEndTime:   endTime,
			inputCondition: types.QueryCondition{
				Filters: types.AllocationCostFilters{
					{
						{
							FieldName: types.FilterFieldNamespace,
							Operator:  types.OperatorIn,
							Values:    []string{"namespace-t1"},
						},
					},
				},
				GroupBy: types.GroupByFieldNamespace,
			},
			outputTotalItemNum:   1,
			outputQueriedItemNum: 1,
			outputItemCost:       50,
		},
		{
			name:           "filter by namespace, group by label",
			inputStartTime: startTime,
			inputEndTime:   endTime,
			inputCondition: types.QueryCondition{
				Filters: types.AllocationCostFilters{
					{
						{
							FieldName: types.FilterFieldNamespace,
							Operator:  types.OperatorIn,
							Values:    []string{"namespace-t1"},
						},
					},
				},
				GroupBy: types.GroupByFieldEnvironment,
			},
			outputTotalItemNum:   1,
			outputQueriedItemNum: 1,
			outputItemCost:       50,
		},
		{
			name:           "filter by label, group by namespace",
			inputStartTime: startTime,
			inputEndTime:   endTime,
			inputCondition: types.QueryCondition{
				Filters: types.AllocationCostFilters{
					{
						{
							FieldName: testFilterFieldEnv,
							Operator:  types.OperatorIn,
							Values:    []string{"dev"},
						},
					},
				},
				GroupBy: types.GroupByFieldNamespace,
			},
			outputTotalItemNum:   3,
			outputQueriedItemNum: 3,
			outputItemCost:       50,
		},
		{
			name:           "filter by label, group by label",
			inputStartTime: startTime,
			inputEndTime:   endTime,
			inputCondition: types.QueryCondition{
				Filters: types.AllocationCostFilters{
					{
						{
							FieldName: testFilterFieldEnv,
							Operator:  types.OperatorIn,
							Values:    []string{"dev"},
						},
					},
				},
				GroupBy: testGroupByFieldApp,
			},
			outputTotalItemNum:   3,
			outputQueriedItemNum: 3,
			outputItemCost:       50,
		},
		{
			name:           "2 hours time range",
			inputStartTime: startTime,
			inputEndTime:   startTime.Add(2 * time.Hour),
			inputCondition: types.QueryCondition{
				Filters: types.AllocationCostFilters{
					{
						{
							FieldName: testFilterFieldEnv,
							Operator:  types.OperatorIn,
							Values:    []string{"dev"},
						},
					},
				},
				GroupBy: testGroupByFieldApp,
			},
			outputTotalItemNum:   3,
			outputQueriedItemNum: 3,
			outputItemCost:       20,
		},
		{
			name:           "paging",
			inputStartTime: startTime,
			inputEndTime:   endTime,
			inputCondition: types.QueryCondition{
				Filters: types.AllocationCostFilters{
					{
						{
							FieldName: testFilterFieldEnv,
							Operator:  types.OperatorIn,
							Values:    []string{"dev"},
						},
					},
				},
				GroupBy: types.GroupByFieldNamespace,
				Paging: types.QueryPagination{
					Page:    2,
					PerPage: 2,
				},
			},
			outputTotalItemNum:   3,
			outputQueriedItemNum: 1,
			outputItemCost:       50,
		},
		{
			name:           "include query",
			inputStartTime: startTime,
			inputEndTime:   endTime,
			inputCondition: types.QueryCondition{
				Filters: types.AllocationCostFilters{
					{
						{
							FieldName: testFilterFieldEnv,
							Operator:  types.OperatorIn,
							Values:    []string{"dev"},
						},
					},
				},
				GroupBy: types.GroupByFieldNamespace,
				Query:   "namespace-t1",
			},
			outputTotalItemNum:   1,
			outputQueriedItemNum: 1,
			outputItemCost:       50,
		},
	}

	for _, v := range cases {
		dsb := accumulateDistributor{client: client}
		items, _, queried, err := dsb.allocationResourceCosts(ctx, v.inputStartTime, v.inputEndTime, v.inputCondition)
		assert.Nil(t, err, "%s: error get allocation resource cost: %w", v.name, err)
		assert.Equal(t, v.outputTotalItemNum, queried, "%s: total item number mismatch", v.name)
		assert.Len(t, items, v.outputQueriedItemNum, "%s: queried item length mismatch", v.name)
		assert.Equal(t, v.outputItemCost, items[0].Cost.TotalCost,
			"%s: first item total cost mismatch", v.name)
	}
}

func testData(ctx context.Context, client *model.Client, startTime, endTime time.Time) (*model.Connector, error) {
	// Clean.
	if _, err := client.Connector.Delete().Exec(ctx); err != nil {
		return nil, err
	}

	if _, err := client.ClusterCost.Delete().Exec(ctx); err != nil {
		return nil, err
	}

	if _, err := client.AllocationCost.Delete().Exec(ctx); err != nil {
		return nil, err
	}

	// Init.
	conn, err := newTestConn(ctx, client)
	if err != nil {
		return nil, err
	}

	var ac []*model.AllocationCost

	var cc []*model.ClusterCost

	hours := endTime.Sub(startTime).Hours()
	for i := 0; i < int(hours); i++ {
		ac = append(ac, testAc("t1", startTime.Add(time.Duration(i)*time.Hour), conn))
		ac = append(ac, testAc("t2", startTime.Add(time.Duration(i)*time.Hour), conn))
		ac = append(ac, testAc("t3", startTime.Add(time.Duration(i)*time.Hour), conn))
		cc = append(cc, testCc("c1", startTime.Add(time.Duration(i)*time.Hour), conn))
	}

	acs, err := dao.AllocationCostCreates(client, ac...)
	if err != nil {
		return nil, err
	}

	err = client.AllocationCosts().CreateBulk(acs...).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("error batch create allocation costs: %w", err)
	}

	ccs, err := dao.ClusterCostCreates(client, cc...)
	if err != nil {
		return nil, err
	}

	err = client.ClusterCosts().CreateBulk(ccs...).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("error batch create cluster costs: %w", err)
	}

	return conn, nil
}

func testAc(name string, startTime time.Time, conn *model.Connector) *model.AllocationCost {
	return &model.AllocationCost{
		StartTime:      startTime,
		EndTime:        startTime.Add(1 * time.Hour),
		ConnectorID:    conn.ID,
		Name:           name,
		Fingerprint:    "",
		ClusterName:    "cluster-test",
		Namespace:      "namespace-" + name,
		Node:           "node-" + name,
		Controller:     "controller-" + name,
		ControllerKind: "controllerKind-" + name,
		Pod:            "pod-" + name,
		Container:      "container-" + name,
		Labels: map[string]string{
			testLabelApp: name,
			testLabelEnv: "dev",
		},
		TotalCost: 10,
		Currency:  1,
		CpuCost:   1,
		GpuCost:   2,
		RamCost:   3,
		PvCost:    4,

		CpuCoreRequest:      100,
		GpuCount:            200,
		RamByteRequest:      300,
		PvBytes:             400,
		CpuCoreUsageAverage: 100,
		CpuCoreUsageMax:     100,
		RamByteUsageAverage: 300,
		RamByteUsageMax:     300,
	}
}

func testCc(name string, startTime time.Time, conn *model.Connector) *model.ClusterCost {
	return &model.ClusterCost{
		StartTime:      startTime,
		EndTime:        startTime.Add(1 * time.Hour),
		ConnectorID:    conn.ID,
		ClusterName:    name,
		TotalCost:      100,
		Currency:       1,
		ManagementCost: 30,
		IdleCost:       40,
		AllocationCost: 30,
	}
}

func newTestConn(ctx context.Context, client *model.Client) (*model.Connector, error) {
	conn, err := client.Connector.Create().
		SetName(time.Now().String()).
		SetType(types.ConnectorTypeK8s).
		SetCategory(types.ConnectorCategoryKubernetes).
		SetConfigVersion("test").
		SetEnableFinOps(true).
		SetConfigData(crypto.Properties{
			"kubeconfig": crypto.StringProperty(""),
		}).
		Save(ctx)

	return conn, err
}
