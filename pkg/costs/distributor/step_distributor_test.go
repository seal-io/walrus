package distributor

import (
	"context"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"

	"github.com/seal-io/seal/pkg/dao/model/enttest"
	_ "github.com/seal-io/seal/pkg/dao/model/runtime"
	"github.com/seal-io/seal/pkg/dao/types"
)

const (
	testDSAEnv = "TEST_DATA_SOURCE_ADDRESS"
)

func TestStepDistribute(t *testing.T) {
	// StepDistribute include postgres sql function date_trunc, only test while setting postgres dns.
	// example: "postgres://seal:seal@localhost:5435/seal?sslmode=disable"
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
	_, err := testData(ctx, client, startTime, endTime)
	assert.Nil(t, err, "error create cost test data: %w", err)

	cases := []struct {
		name                        string
		inputStartTime              time.Time
		inputEndTime                time.Time
		inputCondition              types.QueryCondition
		outputTotalItemNum          int
		outputQueriedItemNum        int
		outputQueriedItemCost       float64
		outputQueriedItemSharedCost float64
	}{
		{
			name:           "daily cost",
			inputStartTime: startTime,
			inputEndTime:   endTime,
			inputCondition: types.QueryCondition{
				Step:    types.StepDay,
				GroupBy: types.GroupByFieldDay,
				SharedCosts: types.ShareCosts{
					{
						IdleCostFilters: types.IdleCostFilters{
							{
								IncludeAll: true,
							},
						},
						ManagementCostFilters: types.ManagementCostFilters{
							{
								IncludeAll: true,
							},
						},
						SharingStrategy: types.SharingStrategyProportionally,
					},
				},
			},
			outputTotalItemNum:          5,
			outputQueriedItemNum:        5,
			outputQueriedItemCost:       2400,
			outputQueriedItemSharedCost: 1680,
		},
		{
			name:           "daily cost with time zone",
			inputStartTime: utc8StartTime,
			inputEndTime:   utc8EndTime,
			inputCondition: types.QueryCondition{
				Step:    types.StepDay,
				GroupBy: types.GroupByFieldDay,
				SharedCosts: types.ShareCosts{
					{
						IdleCostFilters: types.IdleCostFilters{
							{
								IncludeAll: true,
							},
						},
						ManagementCostFilters: types.ManagementCostFilters{
							{
								IncludeAll: true,
							},
						},
						SharingStrategy: types.SharingStrategyProportionally,
					},
				},
			},
			outputTotalItemNum:          1,
			outputQueriedItemNum:        1,
			outputQueriedItemCost:       1600,
			outputQueriedItemSharedCost: 1120,
		},
		{
			name:           "monthly cost",
			inputStartTime: startTime,
			inputEndTime:   endTime,
			inputCondition: types.QueryCondition{
				Step:    types.StepMonth,
				GroupBy: types.GroupByFieldMonth,
				SharedCosts: types.ShareCosts{
					{
						IdleCostFilters: types.IdleCostFilters{
							{
								IncludeAll: true,
							},
						},
						ManagementCostFilters: types.ManagementCostFilters{
							{
								IncludeAll: true,
							},
						},
						SharingStrategy: types.SharingStrategyProportionally,
					},
				},
			},
			outputTotalItemNum:          2,
			outputQueriedItemNum:        2,
			outputQueriedItemCost:       7200,
			outputQueriedItemSharedCost: 5040,
		},
		{
			name:           "monthly cost group by namespace",
			inputStartTime: startTime,
			inputEndTime:   endTime,
			inputCondition: types.QueryCondition{
				Step:    types.StepMonth,
				GroupBy: types.GroupByFieldNamespace,
				SharedCosts: types.ShareCosts{
					{
						IdleCostFilters: types.IdleCostFilters{
							{
								IncludeAll: true,
							},
						},
						ManagementCostFilters: types.ManagementCostFilters{
							{
								IncludeAll: true,
							},
						},
						SharingStrategy: types.SharingStrategyProportionally,
					},
				},
			},
			outputTotalItemNum:          6,
			outputQueriedItemNum:        6,
			outputQueriedItemCost:       2400,
			outputQueriedItemSharedCost: 1680,
		},
	}

	for _, v := range cases {
		dsb := stepDistributor{client: client}
		items, count, err := dsb.distribute(ctx, v.inputStartTime, v.inputEndTime, v.inputCondition)
		assert.Nil(t, err, "%s: error get distribute resource costs: %w", v.name, err)
		assert.Equal(t, v.outputTotalItemNum, count, "%s: total item count mismatch", v.name)
		assert.Len(t, items, v.outputQueriedItemNum, "%s: queried item length mismatch", v.name)
		assert.Equal(t, v.outputQueriedItemCost, items[0].Cost.TotalCost,
			"%s: first item total cost mismatch", v.name)
		assert.Equal(t, v.outputQueriedItemSharedCost, items[0].Cost.SharedCost,
			"%s: first item shared cost mismatch", v.name)
	}
}
