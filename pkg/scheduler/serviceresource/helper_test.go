package serviceresource

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getBatches(t *testing.T) {
	type input struct {
		total, burst, minSize int
	}

	type output struct {
		size, count int
	}

	testCases := []struct {
		name     string
		given    input
		expected output
	}{
		{
			name: "total is negative",
			given: input{
				total:   -1,
				burst:   64,
				minSize: 10,
			},
			expected: output{
				size:  10,
				count: 0,
			},
		},
		{
			name: "burst is negative",
			given: input{
				total:   5000,
				burst:   -1,
				minSize: 10,
			},
			expected: output{
				size:  5000,
				count: 1,
			},
		},
		{
			name: "minSize is negative",
			given: input{
				total:   5000,
				burst:   64,
				minSize: -1,
			},
			expected: output{
				size:  79,
				count: 64,
			},
		},
		{
			name: "total/burst is greater than minSize",
			given: input{
				total:   5000,
				burst:   64,
				minSize: 10,
			},
			expected: output{
				size:  79,
				count: 64,
			},
		},
		{
			name: "total/burst is greater than minSize, and total is not a multiple of burst",
			given: input{
				total:   5001,
				burst:   64,
				minSize: 10,
			},
			expected: output{
				size:  79,
				count: 64,
			},
		},
		{
			name: "total/burst is less than minSize",
			given: input{
				total:   5000,
				burst:   64,
				minSize: 100,
			},
			expected: output{
				size:  100,
				count: 50,
			},
		},
		{
			name: "total/burst is less than minSize, and total is not a multiple of burst",
			given: input{
				total:   5001,
				burst:   64,
				minSize: 100,
			},
			expected: output{
				size:  100,
				count: 51,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var actual output
			actual.size, actual.count = getBatches(tc.given.total, tc.given.burst, tc.given.minSize)
			assert.Equal(t, tc.expected, actual)
			assert.GreaterOrEqual(t, actual.size*actual.count, tc.given.total)

			if tc.given.burst > 0 {
				assert.LessOrEqual(t, actual.count, tc.given.burst)
			}
		})
	}
}
