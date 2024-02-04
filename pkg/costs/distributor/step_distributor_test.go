package distributor

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/seal-io/walrus/pkg/dao/model/enttest"
	"github.com/seal-io/walrus/pkg/dao/types"

	_ "github.com/lib/pq"
	_ "github.com/seal-io/walrus/pkg/dao/model/runtime"
)

const (
	testDSAEnv = "TEST_DATA_SOURCE_ADDRESS"
)

func TestStepDistribute(t *testing.T) {
	// StepDistribute include postgres sql function date_trunc, only test while setting postgres dns.
	// example: "postgres://walrus:walrus@localhost:5435/walrus?sslmode=disable"
	dsa := os.Getenv(testDSAEnv)
	if dsa == "" {
		t.Skip("environment TEST_DATA_SOURCE_ADDRESS isn't provided")
		return
	}
	ctx := context.Background()
	client := enttest.Open(t, "postgres", dsa)

	defer func() {
		err := client.Close()
		if err != nil {
			t.Logf("error close sql client: %v", err)
		}
	}()

	var (
		startTime = time.Date(2023, 0o2, 27, 0, 0, 0, 0, time.UTC)
		endTime   = startTime.Add(5 * 24 * time.Hour)

		utc8StartTime = time.Date(2023, 0o2, 27, 0, 0, 0, 0, time.FixedZone("Asia/Shanghai", 28800))
		utc8EndTime   = utc8StartTime.Add(1 * 24 * time.Hour)
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
		outputItems        []stepTestOutputItem
	}{
		{
			name:           "single connector, daily cost",
			inputStartTime: startTime,
			inputEndTime:   endTime,
			inputCondition: types.QueryCondition{
				Filters: filterFirstConn,
				Step:    types.StepDay,
				GroupBy: types.GroupByFieldDay,
			},
			outputTotalItemNum: 5,
			outputItems: []stepTestOutputItem{
				newStepTestOutputItem("2023-03-03T00:00:00Z", "2023-03-03T00:00:00Z", 2520, 0),
				newStepTestOutputItem("2023-03-02T00:00:00Z", "2023-03-02T00:00:00Z", 2520, 0),
				newStepTestOutputItem("2023-03-01T00:00:00Z", "2023-03-01T00:00:00Z", 2520, 0),
				newStepTestOutputItem("2023-02-28T00:00:00Z", "2023-02-28T00:00:00Z", 2520, 0),
				newStepTestOutputItem("2023-02-27T00:00:00Z", "2023-02-27T00:00:00Z", 2520, 0),
			},
		},
		{
			name:           "single connector, daily cost with equally share idle and management costs",
			inputStartTime: startTime,
			inputEndTime:   endTime,
			inputCondition: types.QueryCondition{
				Filters: types.CostFilters{
					{
						filterFirstConn[0][0],
						filterExcludeIdleMgnt[0][0],
					},
				},
				Step:    types.StepDay,
				GroupBy: types.GroupByFieldDay,
				SharedOptions: &types.SharedCostOptions{
					Idle: &types.IdleShareOption{
						SharingStrategy: types.SharingStrategyEqually,
					},
					Management: &types.ManagementShareOption{
						SharingStrategy: types.SharingStrategyEqually,
					},
				},
			},
			outputTotalItemNum: 5,
			outputItems: []stepTestOutputItem{
				newStepTestOutputItem("2023-03-03T00:00:00Z", "2023-03-03T00:00:00Z", 2520, 1080),
				newStepTestOutputItem("2023-03-02T00:00:00Z", "2023-03-02T00:00:00Z", 2520, 1080),
				newStepTestOutputItem("2023-03-01T00:00:00Z", "2023-03-01T00:00:00Z", 2520, 1080),
				newStepTestOutputItem("2023-02-28T00:00:00Z", "2023-02-28T00:00:00Z", 2520, 1080),
				newStepTestOutputItem("2023-02-27T00:00:00Z", "2023-02-27T00:00:00Z", 2520, 1080),
			},
		},
		{
			name:           "single connector, daily cost with time zone",
			inputStartTime: utc8StartTime,
			inputEndTime:   utc8EndTime,
			inputCondition: types.QueryCondition{
				Filters: types.CostFilters{
					{
						filterFirstConn[0][0],
						filterExcludeIdleMgnt[0][0],
					},
				},
				Step:          types.StepDay,
				GroupBy:       types.GroupByFieldDay,
				SharedOptions: sharedOptionSplitIdleMgntProportionallly,
			},
			outputTotalItemNum: 1,
			outputItems: []stepTestOutputItem{
				newStepTestOutputItem("2023-02-27T00:00:00Z", "2023-02-27T00:00:00Z", 1680, 720),
			},
		},
		{
			name:           "single connector, monthly cost",
			inputStartTime: startTime,
			inputEndTime:   endTime,
			inputCondition: types.QueryCondition{
				Filters: types.CostFilters{
					{
						filterFirstConn[0][0],
						filterExcludeIdleMgnt[0][0],
					},
				},
				Step:          types.StepMonth,
				GroupBy:       types.GroupByFieldMonth,
				SharedOptions: sharedOptionSplitIdleMgntProportionallly,
			},
			outputTotalItemNum: 2,
			outputItems: []stepTestOutputItem{
				newStepTestOutputItem("2023-03-01T00:00:00Z", "2023-03-01T00:00:00Z", 7560, 3240),
				newStepTestOutputItem("2023-02-01T00:00:00Z", "2023-02-01T00:00:00Z", 5040, 2160),
			},
		},
		{
			name:           "single connector, daily cost group by namespace",
			inputStartTime: startTime,
			inputEndTime:   endTime,
			inputCondition: types.QueryCondition{
				Filters: types.CostFilters{
					{
						filterFirstConn[0][0],
						filterExcludeIdleMgnt[0][0],
					},
				},
				Step:    types.StepDay,
				GroupBy: types.GroupByFieldNamespace,
			},
			outputTotalItemNum: 15,
			outputItems: []stepTestOutputItem{
				newStepTestOutputItem("namespace-t3", "2023-03-03T00:00:00Z", 720, 0),
				newStepTestOutputItem("namespace-t2", "2023-03-03T00:00:00Z", 480, 0),
				newStepTestOutputItem("namespace-t1", "2023-03-03T00:00:00Z", 240, 0),

				newStepTestOutputItem("namespace-t3", "2023-03-02T00:00:00Z", 720, 0),
				newStepTestOutputItem("namespace-t2", "2023-03-02T00:00:00Z", 480, 0),
				newStepTestOutputItem("namespace-t1", "2023-03-02T00:00:00Z", 240, 0),
			},
		},
		{
			name:           "multiple connector, daily cost group by namespace",
			inputStartTime: startTime,
			inputEndTime:   endTime,
			inputCondition: types.QueryCondition{
				Filters: filterExcludeIdleMgnt,
				Step:    types.StepDay,
				GroupBy: types.GroupByFieldNamespace,
			},
			outputTotalItemNum: 20,
			outputItems: []stepTestOutputItem{
				newStepTestOutputItem("namespace-t3", "2023-03-03T00:00:00Z", 1440, 0),
			},
		},
	}

	for _, v := range cases {
		dsb := stepDistributor{client: client}
		items, count, err := dsb.distribute(ctx, v.inputStartTime, v.inputEndTime, v.inputCondition)

		assert.Nil(t, err, "%s: error get distribute resource cost: %w", v.name, err)
		assert.Equal(t, v.outputTotalItemNum, count, "%s: total item count mismatch", v.name)

		for i := range v.outputItems {
			assert.Equal(t, v.outputItems[i].itemName, items[i].ItemName,
				"%s: item %d name mismatch", v.name, i)

			assert.Equal(t, v.outputItems[i].startTime, items[i].StartTime.Format(time.RFC3339),
				"%s: item %d startTime mismatch", v.name, i)

			assert.Equal(t, v.outputItems[i].totalCost, items[i].TotalCost,
				"%s: item %d total cost mismatch", v.name, i)

			assert.Equal(t, v.outputItems[i].sharedCost, items[i].SharedCost,
				"%s: item %d shared cost mismatch", v.name, i)
		}
	}
}

func newStepTestOutputItem(itemName, startTime string, totalCost, sharedCost float64) stepTestOutputItem {
	return stepTestOutputItem{
		testOutputItem: testOutputItem{
			itemName:   itemName,
			totalCost:  totalCost,
			sharedCost: sharedCost,
		},
		startTime: startTime,
	}
}

type stepTestOutputItem struct {
	testOutputItem
	startTime string
}
