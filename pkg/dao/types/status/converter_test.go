package status

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConverterExample(t *testing.T) {
	// 1. Define the status of resource.
	const (
		StatusCreating     = "Creating"
		StatusAvailable    = "Available"
		StatusInActive     = "InActive"
		StatusUnAvailable  = "UnAvailable"
		StatusCreateFailed = "CreateFailed"
	)

	// 2.  Clarify the status and its sensible status,
	// 		only need to define normal and error status.
	//  | Human Readable Status | Human Sensible Status |
	//  | --------------------- | --------------------- |
	//  | StatusAvailable       |                       |
	//  | StatusInActive        | Inactive              |
	//  | StatusUnAvailable     | Error                 |
	//  | StatusCreateFailed    | Error                 |

	// 3. Create converter, indicate the normal status, error status.
	f := NewConverter(
		[]string{
			StatusAvailable,
		},
		[]string{
			StatusInActive,
		},
		[]string{
			StatusUnAvailable,
			StatusCreateFailed,
		})

	var p printer

	// 4.  Input the status Creating, which means progressing,
	//      we should get a transitioning summary.
	st := f.Convert(StatusCreating, "creating")
	p.Dump("Creating [T]", &st.Summary)

	t.Log(p.String())
}

func TestConverter(t *testing.T) {
	const (
		StatusCreating     = "Creating"
		StatusAvailable    = "Available"
		StatusInActive     = "InActive"
		StatusUnAvailable  = "UnAvailable"
		StatusCreateFailed = "CreateFailed"
	)

	f := NewConverter(
		[]string{
			StatusAvailable,
		},
		[]string{
			StatusInActive,
		},
		[]string{
			StatusUnAvailable,
			StatusCreateFailed,
		},
	)

	testCases := []struct {
		name     string
		input    []string
		expected *Status
	}{
		{
			name:  "empty",
			input: []string{"", ""},
			expected: &Status{
				Summary: Summary{
					SummaryStatus: "",
				},
			},
		},
		{
			name:  "creating",
			input: []string{StatusCreating, ""},
			expected: &Status{
				Summary: Summary{
					SummaryStatus: StatusCreating,
					Transitioning: true,
				},
			},
		},
		{
			name:  "inactive",
			input: []string{StatusInActive, ""},
			expected: &Status{
				Summary: Summary{
					SummaryStatus: StatusInActive,
					Inactive:      true,
				},
			},
		},
		{
			name:  "unavailable",
			input: []string{StatusUnAvailable, ""},
			expected: &Status{
				Summary: Summary{
					SummaryStatus: StatusUnAvailable,
					Error:         true,
				},
			},
		},
		{
			name:  "create failed",
			input: []string{StatusCreateFailed, "failed to create resource"},
			expected: &Status{
				Summary: Summary{
					SummaryStatus:        "CreateFailed",
					SummaryStatusMessage: "failed to create resource",
					Error:                true,
					Transitioning:        false,
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := f.Convert(tc.input[0], tc.input[1])
			assert.Equal(t, tc.expected, actual, "case %q", tc.name)
		})
	}
}
