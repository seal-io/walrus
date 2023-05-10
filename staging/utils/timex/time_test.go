package timex

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetTimeUnitSeries(t *testing.T) {
	testCases := []struct {
		name      string
		startTime time.Time
		endTime   time.Time
		unit      string
		expected  []time.Time
	}{
		{
			name:      "test day",
			startTime: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			endTime:   time.Date(2022, 1, 3, 0, 0, 0, 0, time.UTC),
			unit:      "day",
			expected: []time.Time{
				time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:      "test week",
			startTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			endTime:   time.Date(2023, 1, 12, 0, 0, 0, 0, time.UTC),
			unit:      "week",
			expected: []time.Time{
				time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 1, 8, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:      "test month",
			startTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			endTime:   time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
			unit:      "month",
			expected: []time.Time{
				time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:      "test quarter",
			startTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			endTime:   time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC),
			unit:      "quarter",
			expected: []time.Time{
				time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name:      "test year",
			startTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			endTime:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			unit:      "year",
			expected: []time.Time{
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tc := range testCases {
		got, err := GetTimeSeries(tc.startTime, tc.endTime, tc.unit, time.UTC)
		if err != nil {
			assert.Errorf(t, err, "GetTimeSeries(%v, %v, %v, %v) error", tc.startTime, tc.endTime, tc.unit, time.UTC)
		}
		if !assert.Equal(t, tc.expected, got) {
			t.Errorf("GetTimeSeries(%v, %v, %v, %v) = %v, want %v", tc.startTime,
				tc.endTime, tc.unit, time.UTC, got, tc.expected)
		}
	}
}
