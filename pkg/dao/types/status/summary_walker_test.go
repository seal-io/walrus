package status

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleWalkerExample(t *testing.T) {
	// 1. Define resource with status.
	type ExampleResource struct {
		Status Status
	}

	// 2. Define the status of ExampleResource.
	const (
		StatusCreating     = "Creating"
		StatusWaiting      = "Waiting"
		StatusAvailable    = "Available"
		StatusUnAvailable  = "UnAvailable"
		StatusCreateFailed = "CreateFailed"
	)

	// 2.1  clarify the status and its sensible status.
	//  | Human Readable Status | Human Sensible Status |
	//  | --------------------- | --------------------- |
	//  | StatusCreating        | Transitioning         |
	//  | StatusWaiting         | Transitioning         |
	//  | StatusAvailable       |                       |
	//  | StatusUnAvailable     | Error                 |
	//  | StatusCreateFailed    | Error                 |

	// 3. Create walker, indicate the normal status, error status and transitioning status.
	f := NewSummaryWalker(
		[]string{
			StatusAvailable,
		}, []string{
			StatusUnAvailable,
			StatusCreateFailed,
		}, []string{
			StatusCreating,
			StatusWaiting,
		})

	var p printer

	// 4. Create an instance of ExampleResource.
	var r ExampleResource

	// 4.  At beginning, the status is empty(we haven't configured any conditions or summary result),
	//      so we should get a empty summary.
	p.Dump("Empty Status [N]", f.Walk(&r.Status))

	// 4.2  Set the SummaryStatus status to Creating, which means progressing,
	//      we should get a transitioning summary.
	r.Status.SummaryStatus = StatusCreating
	p.Dump("Creating [T]", f.Walk(&r.Status))

	t.Log(p.String())
}

func TestSimpleWalker(t *testing.T) {
	const (
		StatusCreating     = "Creating"
		StatusWaiting      = "Waiting"
		StatusAvailable    = "Available"
		StatusUnAvailable  = "UnAvailable"
		StatusCreateFailed = "CreateFailed"
	)

	f := NewSummaryWalker(
		[]string{
			StatusAvailable,
		}, []string{
			StatusUnAvailable,
			StatusCreateFailed,
		}, []string{
			StatusCreating,
			StatusWaiting,
		},
		func(current Summary) (change bool, update Summary) {
			if current.SummaryStatus == StatusCreateFailed {
				return true,
					Summary{
						SummaryStatus:        "Error",
						SummaryStatusMessage: "failed to create resource",
						Error:                true,
						Transitioning:        false,
					}
			}

			return false, current
		},
	)

	testCases := []struct {
		name     string
		input    *Status
		expected *Summary
	}{
		{
			name:  "empty",
			input: &Status{},
			expected: &Summary{
				SummaryStatus: "",
			},
		},
		{
			name: "creating",
			input: &Status{
				Summary: Summary{
					SummaryStatus: StatusCreating,
				},
			},
			expected: &Summary{
				SummaryStatus: StatusCreating,
				Transitioning: true,
			},
		},
		{
			name: "unavailable",
			input: &Status{
				Summary: Summary{
					SummaryStatus: StatusUnAvailable,
				},
			},
			expected: &Summary{
				SummaryStatus: StatusUnAvailable,
				Error:         true,
			},
		},
		{
			name: "create failed",
			input: &Status{
				Summary: Summary{
					SummaryStatus: StatusCreateFailed,
				},
			},
			expected: &Summary{
				SummaryStatus:        "Error",
				SummaryStatusMessage: "failed to create resource",
				Error:                true,
				Transitioning:        false,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := f.Walk(tc.input)
			assert.Equal(t, tc.expected, actual, "case %q", tc.name)
		})
	}
}
