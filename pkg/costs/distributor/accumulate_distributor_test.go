package distributor

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/enttest"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/crypto"
	"github.com/seal-io/walrus/utils/timex"

	_ "github.com/mattn/go-sqlite3"
	_ "github.com/seal-io/walrus/pkg/dao/model/runtime"
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

var (
	filterExcludeIdleMgnt = types.CostFilters{
		{
			{
				FieldName: types.FilterFieldName,
				Operator:  types.OperatorNotIn,
				Values: []string{
					types.IdleCostItemName,
					types.ManagementCostItemName,
				},
			},
		},
	}

	sharedOptionSplitIdleEqually = &types.IdleShareOption{
		SharingStrategy: types.SharingStrategyEqually,
	}

	sharedOptionSplitMgntEqually = &types.ManagementShareOption{
		SharingStrategy: types.SharingStrategyEqually,
	}

	sharedOptionSplitIdleMgntEqually = &types.SharedCostOptions{
		Idle: &types.IdleShareOption{
			SharingStrategy: types.SharingStrategyEqually,
		},
		Management: &types.ManagementShareOption{
			SharingStrategy: types.SharingStrategyEqually,
		},
	}

	sharedOptionSplitIdleMgntProportionallly = &types.SharedCostOptions{
		Idle: &types.IdleShareOption{
			SharingStrategy: types.SharingStrategyProportionally,
		},
		Management: &types.ManagementShareOption{
			SharingStrategy: types.SharingStrategyProportionally,
		},
	}

	sharedOptionSplitIdleProportionally = &types.IdleShareOption{
		SharingStrategy: types.SharingStrategyProportionally,
	}

	sharedOptionSplitMgntProportionally = &types.ManagementShareOption{
		SharingStrategy: types.SharingStrategyProportionally,
	}
)

type testOutputItem struct {
	itemName   string
	totalCost  float64
	sharedCost float64
}

func newTestOutputItem(itemName string, totalCost, sharedCost float64) testOutputItem {
	return testOutputItem{
		itemName:   itemName,
		totalCost:  totalCost,
		sharedCost: sharedCost,
	}
}

func TestAccumulateDistribute(t *testing.T) {
	ctx := context.Background()

	client := enttest.Open(t, testDriverName, testDataSourceName)
	defer client.Close()

	var (
		startTime = timex.StartTimeOfHour(time.Now(), time.UTC)
		endTime   = startTime.Add(5 * time.Hour)
	)

	conns, err := testData(ctx, client, startTime, endTime)
	assert.Nil(t, err, "error create cost test data: %v", err)

	filterFirstConn := types.CostFilters{
		{
			{
				FieldName: types.FilterFieldConnectorID,
				Operator:  types.OperatorIn,
				Values: []string{
					conns[0].ID.String(),
				},
			},
		},
	}

	cases := []struct {
		name               string
		inputStartTime     time.Time
		inputEndTime       time.Time
		inputCondition     types.QueryCondition
		outputTotalItemNum int
		outputItems        []testOutputItem
	}{
		{
			name:           "single connector, time range with no data",
			inputStartTime: time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC),
			inputEndTime:   time.Date(2000, 2, 1, 1, 1, 1, 1, time.UTC),
			inputCondition: types.QueryCondition{
				Filters: filterFirstConn,
				GroupBy: types.GroupByFieldNamespace,
			},
			outputTotalItemNum: 0,
		},
		{
			name:           "single connector, equally share item cost with idle costs without exclude idle cost",
			inputStartTime: startTime,
			inputEndTime:   endTime,
			inputCondition: types.QueryCondition{
				Filters: filterFirstConn,
				GroupBy: types.GroupByFieldNamespace,
				SharedOptions: &types.SharedCostOptions{
					Idle: sharedOptionSplitIdleEqually,
				},
			},
			outputTotalItemNum: 5,
			outputItems: []testOutputItem{
				newTestOutputItem(types.ManagementCostItemName, 75, 0),
				newTestOutputItem(types.IdleCostItemName, 150, 0),
				newTestOutputItem("namespace-t3", 200, 50),
				newTestOutputItem("namespace-t2", 150, 50),
				newTestOutputItem("namespace-t1", 100, 50),
			},
		},
		{
			name:           "single connector, equally share item cost with management costs without exclude management cost",
			inputStartTime: startTime,
			inputEndTime:   endTime,
			inputCondition: types.QueryCondition{
				Filters: filterFirstConn,
				GroupBy: types.GroupByFieldNamespace,
				SharedOptions: &types.SharedCostOptions{
					Management: sharedOptionSplitMgntEqually,
				},
			},
			outputTotalItemNum: 5,
			outputItems: []testOutputItem{
				newTestOutputItem(types.ManagementCostItemName, 75, 0),
				newTestOutputItem(types.IdleCostItemName, 150, 0),
				newTestOutputItem("namespace-t3", 175, 25),
				newTestOutputItem("namespace-t2", 125, 25),
				newTestOutputItem("namespace-t1", 75, 25),
			},
		},
		{
			name:           "single connector, equally share item cost with idle and management costs without exclude them",
			inputStartTime: startTime,
			inputEndTime:   endTime,
			inputCondition: types.QueryCondition{
				Filters: filterFirstConn,
				GroupBy: types.GroupByFieldNamespace,
				SharedOptions: &types.SharedCostOptions{
					Idle:       sharedOptionSplitIdleEqually,
					Management: sharedOptionSplitMgntEqually,
				},
			},
			outputTotalItemNum: 5,
			outputItems: []testOutputItem{
				newTestOutputItem(types.ManagementCostItemName, 75, 0),
				newTestOutputItem(types.IdleCostItemName, 150, 0),
				newTestOutputItem("namespace-t3", 225, 75),
				newTestOutputItem("namespace-t2", 175, 75),
				newTestOutputItem("namespace-t1", 125, 75),
			},
		},
		{
			name:           "single connector, equally share item cost with idle and management costs with exclude them",
			inputStartTime: startTime,
			inputEndTime:   endTime,
			inputCondition: types.QueryCondition{
				Filters: types.CostFilters{
					{
						filterFirstConn[0][0],
						filterExcludeIdleMgnt[0][0],
					},
				},
				GroupBy:       types.GroupByFieldNamespace,
				SharedOptions: sharedOptionSplitIdleMgntEqually,
			},
			outputTotalItemNum: 3,
			outputItems: []testOutputItem{
				newTestOutputItem("namespace-t3", 225, 75),
				newTestOutputItem("namespace-t2", 175, 75),
				newTestOutputItem("namespace-t1", 125, 75),
			},
		},

		{
			name:           "single connector, equally share item cost with filter",
			inputStartTime: startTime,
			inputEndTime:   endTime,
			inputCondition: types.QueryCondition{
				Filters: types.CostFilters{
					{
						filterFirstConn[0][0],
						filterExcludeIdleMgnt[0][0],
						{
							FieldName: types.FilterFieldNamespace,
							Operator:  types.OperatorNotIn,
							Values:    []string{"namespace-t1"},
						},
					},
				},
				GroupBy: types.GroupByFieldNamespace,
				SharedOptions: &types.SharedCostOptions{
					Item: types.ItemSharedOptions{
						{
							Filters: types.CostFilters{
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
			},
			outputTotalItemNum: 2,
			outputItems: []testOutputItem{
				newTestOutputItem("namespace-t3", 175, 25),
				newTestOutputItem("namespace-t2", 125, 25),
			},
		},
		{
			name:           "single connector, equally share item cost with management, idle, custom resource costs",
			inputStartTime: startTime,
			inputEndTime:   endTime,
			inputCondition: types.QueryCondition{
				Filters: types.CostFilters{
					{
						filterFirstConn[0][0],
						filterExcludeIdleMgnt[0][0],
						{
							FieldName: types.FilterFieldNamespace,
							Operator:  types.OperatorNotIn,
							Values:    []string{"namespace-t1"},
						},
					},
				},
				GroupBy: types.GroupByFieldNamespace,
				SharedOptions: &types.SharedCostOptions{
					Idle:       sharedOptionSplitIdleEqually,
					Management: sharedOptionSplitMgntEqually,
					Item: types.ItemSharedOptions{
						{
							Filters: types.CostFilters{
								{
									{
										FieldName: types.FilterFieldNamespace,
										Operator:  types.OperatorIn,
										Values: []string{
											"namespace-t1",
										},
									},
								},
							},
							SharingStrategy: types.SharingStrategyEqually,
						},
					},
				},
			},
			outputTotalItemNum: 2,
			outputItems: []testOutputItem{
				newTestOutputItem("namespace-t3", 287.5, 137.5),
				newTestOutputItem("namespace-t2", 237.5, 137.5),
			},
		},
		{
			name:           "single connector, proportionally share idle cost with exclude idle cost",
			inputStartTime: startTime,
			inputEndTime:   endTime,
			inputCondition: types.QueryCondition{
				Filters: types.CostFilters{
					{
						filterFirstConn[0][0],
						{
							FieldName: types.FilterFieldName,
							Operator:  types.OperatorNotIn,
							Values: []string{
								types.IdleCostItemName,
							},
						},
					},
				},
				GroupBy: types.GroupByFieldNamespace,
				SharedOptions: &types.SharedCostOptions{
					Idle: &types.IdleShareOption{
						SharingStrategy: types.SharingStrategyProportionally,
					},
				},
			},
			outputTotalItemNum: 4,
			outputItems: []testOutputItem{
				newTestOutputItem(types.ManagementCostItemName, 75, 0),
				newTestOutputItem("namespace-t3", 225, 75),
				newTestOutputItem("namespace-t2", 150, 50),
				newTestOutputItem("namespace-t1", 75, 25),
			},
		},
		{
			name:           "single connector, proportionally share idle cost and management cost with exclude idle and management cost",
			inputStartTime: startTime,
			inputEndTime:   endTime,
			inputCondition: types.QueryCondition{
				Filters: types.CostFilters{
					{
						filterFirstConn[0][0],
						filterExcludeIdleMgnt[0][0],
					},
				},
				GroupBy: types.GroupByFieldNamespace,
				SharedOptions: &types.SharedCostOptions{
					Idle:       sharedOptionSplitIdleProportionally,
					Management: sharedOptionSplitMgntProportionally,
				},
			},
			outputTotalItemNum: 3,
			outputItems: []testOutputItem{
				newTestOutputItem("namespace-t3", 262.5, 112.5),
				newTestOutputItem("namespace-t2", 175, 75),
				newTestOutputItem("namespace-t1", 87.5, 37.5),
			},
		},
		{
			name:           "multiple connectors, group by connector id without shared cost options",
			inputStartTime: startTime,
			inputEndTime:   endTime,
			inputCondition: types.QueryCondition{
				GroupBy: types.GroupByFieldConnectorID,
			},
			outputTotalItemNum: 2,
			outputItems: []testOutputItem{
				newTestOutputItem(conns[1].Name, 725, 0),
				newTestOutputItem(conns[0].Name, 525, 0),
			},
		},
		{
			name:           "multiple connectors, equally share idle cost and management cost with exclude idle and management cost",
			inputStartTime: startTime,
			inputEndTime:   endTime,
			inputCondition: types.QueryCondition{
				Filters: filterExcludeIdleMgnt,
				GroupBy: types.GroupByFieldNamespace,
				SharedOptions: &types.SharedCostOptions{
					Idle:       sharedOptionSplitIdleEqually,
					Management: sharedOptionSplitMgntEqually,
				},
			},
			outputTotalItemNum: 4,
			outputItems: []testOutputItem{
				newTestOutputItem("namespace-t3", 431.25, 131.25),
				newTestOutputItem("namespace-t2", 331.25, 131.25),
				newTestOutputItem("namespace-t4", 256.25, 56.25),
				newTestOutputItem("namespace-t1", 231.25, 131.25),
			},
		},
		{
			name:           "multiple connectors, proportionally share idle cost and management cost with exclude idle and management cost",
			inputStartTime: startTime,
			inputEndTime:   endTime,
			inputCondition: types.QueryCondition{
				Filters:       filterExcludeIdleMgnt,
				GroupBy:       types.GroupByFieldNamespace,
				SharedOptions: sharedOptionSplitIdleMgntProportionallly,
			},
			outputTotalItemNum: 4,
			outputItems: []testOutputItem{
				newTestOutputItem("namespace-t3", 480, 180),
				newTestOutputItem("namespace-t2", 320, 120),
				newTestOutputItem("namespace-t4", 290, 90),
				newTestOutputItem("namespace-t1", 160, 60),
			},
		},
		{
			name:           "multiple connectors, proportionally share item cost with exclude idle and management cost",
			inputStartTime: startTime,
			inputEndTime:   endTime,
			inputCondition: types.QueryCondition{
				Filters: filterExcludeIdleMgnt,
				GroupBy: types.GroupByFieldNamespace,
				SharedOptions: &types.SharedCostOptions{
					Item: types.ItemSharedOptions{
						{
							Filters: types.CostFilters{
								{
									filterFirstConn[0][0],
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
			},
			outputTotalItemNum: 4,
			outputItems: []testOutputItem{
				newTestOutputItem("namespace-t3", 325, 25),
				newTestOutputItem("namespace-t2", 216.66666666666666, 16.666666666666664),
				newTestOutputItem("namespace-t4", 200, 0),
				newTestOutputItem("namespace-t1", 108.33333333333333, 8.333333333333332),
			},
		},
		{
			name:           "multiple connectors, filter by namespace, group by namespace",
			inputStartTime: startTime,
			inputEndTime:   endTime,
			inputCondition: types.QueryCondition{
				Filters: types.CostFilters{
					{
						{
							FieldName: types.FilterFieldNamespace,
							Operator:  types.OperatorIn,
							Values:    []string{"namespace-t1", "namespace-t2"},
						},
					},
				},
				GroupBy: types.GroupByFieldNamespace,
			},
			outputTotalItemNum: 2,
			outputItems: []testOutputItem{
				newTestOutputItem("namespace-t2", 200, 0),
				newTestOutputItem("namespace-t1", 100, 0),
			},
		},
		{
			name:           "multiple connectors, filter by label, group by label",
			inputStartTime: startTime,
			inputEndTime:   endTime,
			inputCondition: types.QueryCondition{
				Filters: types.CostFilters{
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
			outputTotalItemNum: 4,
			outputItems: []testOutputItem{
				newTestOutputItem("t3", 300, 0),
			},
		},
		{
			name:           "multiple connectors, paging",
			inputStartTime: startTime,
			inputEndTime:   endTime,
			inputCondition: types.QueryCondition{
				Filters:       filterExcludeIdleMgnt,
				GroupBy:       types.GroupByFieldNamespace,
				SharedOptions: sharedOptionSplitIdleMgntProportionallly,
				Paging: types.QueryPagination{
					Page:    1,
					PerPage: 3,
				},
			},
			outputTotalItemNum: 4,
			outputItems: []testOutputItem{
				newTestOutputItem("namespace-t3", 480, 180),
				newTestOutputItem("namespace-t2", 320, 120),
				newTestOutputItem("namespace-t4", 290, 90),
			},
		},
		{
			name:           "multiple connectors, include query",
			inputStartTime: startTime,
			inputEndTime:   endTime,
			inputCondition: types.QueryCondition{
				Filters:       filterExcludeIdleMgnt,
				GroupBy:       types.GroupByFieldNamespace,
				SharedOptions: sharedOptionSplitIdleMgntProportionallly,
				Query:         "namespace-t3",
			},
			outputTotalItemNum: 1,
			outputItems: []testOutputItem{
				newTestOutputItem("namespace-t3", 480, 180),
			},
		},
		{
			name:           "multiple connectors, include all",
			inputStartTime: startTime,
			inputEndTime:   endTime,
			inputCondition: types.QueryCondition{
				GroupBy: types.GroupByFieldConnectorID,
			},
			outputTotalItemNum: 2,
			outputItems: []testOutputItem{
				newTestOutputItem("connector2", 725, 0),
				newTestOutputItem("connector1", 525, 0),
			},
		},
	}

	for _, v := range cases {
		dsb := accumulateDistributor{client: client}
		items, count, err := dsb.distribute(ctx, v.inputStartTime, v.inputEndTime, v.inputCondition)

		assert.Nil(t, err, "%s: error get distribute resource cost: %w", v.name, err)
		assert.Equal(t, v.outputTotalItemNum, count, "%s: total item count mismatch", v.name)

		for i := range v.outputItems {
			assert.Equal(t, v.outputItems[i].itemName, items[i].ItemName,
				"%s: item %d name mismatch", v.name, i)

			assert.Equal(t, v.outputItems[i].totalCost, items[i].TotalCost,
				"%s: item %d total cost mismatch", v.name, i)

			assert.Equal(t, v.outputItems[i].sharedCost, items[i].SharedCost,
				"%s: item %d shared cost mismatch", v.name, i)
		}
	}
}

func testData(ctx context.Context, client *model.Client, startTime, endTime time.Time) ([]*model.Connector, error) {
	// Clean.
	if _, err := client.Connector.Delete().Exec(ctx); err != nil {
		return nil, err
	}

	if _, err := client.CostReports().Delete().Exec(ctx); err != nil {
		return nil, err
	}

	// Init.
	var (
		conns = make([]*model.Connector, 2)
		err   error
	)

	conns[0], err = testDataForConnector(ctx, client, startTime, endTime, "1", 3)
	if err != nil {
		return nil, err
	}

	conns[1], err = testDataForConnector(ctx, client, startTime, endTime, "2", 4)
	if err != nil {
		return nil, err
	}

	return conns, nil
}

func testDataForConnector(
	ctx context.Context,
	client *model.Client,
	startTime,
	endTime time.Time,
	nameSuffix string,
	itemCount int,
) (*model.Connector, error) {
	conn, err := newTestConn(ctx, client, nameSuffix)
	if err != nil {
		return nil, err
	}

	var ac []*model.CostReport

	hours := endTime.Sub(startTime).Hours()

	// Management and idle cost.
	for i := 0; i < int(hours); i++ {
		ac = append(ac, testAc(types.ManagementCostItemName, startTime.Add(time.Duration(i)*time.Hour), 15, conn))
		ac = append(ac, testAc(types.IdleCostItemName, startTime.Add(time.Duration(i)*time.Hour), 30, conn))
	}

	// Item cost.
	for i := 0; i < int(hours); i++ {
		for ic := 0; ic < itemCount; ic++ {
			var (
				name      = "t" + strconv.Itoa(ic+1)
				st        = startTime.Add(time.Duration(i) * time.Hour)
				totalCost = 10 * float64(ic+1)
			)

			ac = append(ac, testAc(name, st, totalCost, conn))
		}
	}

	err = client.CostReport.CreateBulk().
		Set(ac...).
		Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("error batch create item costs: %w", err)
	}

	return conn, nil
}

func testAc(name string, startTime time.Time, totalCost float64, conn *model.Connector) *model.CostReport {
	var (
		projName = "project-" + conn.Name
		env      = "dev"
	)

	cr := &model.CostReport{
		StartTime:   startTime,
		EndTime:     startTime.Add(1 * time.Hour),
		ConnectorID: conn.ID,
		Name:        name,
		ClusterName: conn.Name,

		Namespace:      "namespace-" + name,
		Node:           "node-" + name,
		Controller:     "controller-" + name,
		ControllerKind: "controllerKind-" + name,
		Pod:            "pod-" + name,
		Container:      "container-" + name,
		Labels: map[string]string{
			// Labels for sqlite test.
			testLabelApp: name,
			testLabelEnv: env,

			// Original labels.
			types.LabelWalrusServiceName:     name,
			types.LabelWalrusEnvironmentName: env,
			types.LabelWalrusProjectName:     projName,
			types.LabelWalrusServicePath:     fmt.Sprintf("%s/%s/%s", projName, env, name),
			types.LabelWalrusEnvironmentPath: fmt.Sprintf("%s/%s", projName, env),
		},

		Fingerprint: "",
		TotalCost:   totalCost,
		Currency:    1,
		CPUCost:     totalCost * 0.5,
		GPUCost:     totalCost * 0.1,
		RAMCost:     totalCost * 0.3,
		PVCost:      totalCost * 0.1,

		CPUCoreRequest:      100,
		GPUCount:            200,
		RAMByteRequest:      300,
		PVBytes:             400,
		CPUCoreUsageAverage: 100,
		CPUCoreUsageMax:     100,
		RAMByteUsageAverage: 300,
		RAMByteUsageMax:     300,
	}

	if slices.Contains(
		[]string{
			types.IdleCostItemName,
			types.ManagementCostItemName,
		}, name) {
		cr.Namespace = name
		cr.Node = name
		cr.Controller = name
		cr.ControllerKind = name
		cr.Pod = name
		cr.Container = name
		cr.Labels = nil
	}

	return cr
}

func newTestConn(ctx context.Context, client *model.Client, nameSuffix string) (*model.Connector, error) {
	conn, err := client.Connector.Create().
		SetName("connector" + nameSuffix).
		SetApplicableEnvironmentType(types.EnvironmentDevelopment).
		SetType(types.ConnectorTypeKubernetes).
		SetCategory(types.ConnectorCategoryKubernetes).
		SetConfigVersion("test").
		SetEnableFinOps(true).
		SetConfigData(crypto.Properties{
			"kubeconfig": crypto.StringProperty(""),
		}).
		Save(ctx)

	return conn, err
}
