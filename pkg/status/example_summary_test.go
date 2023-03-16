package status

import (
	"fmt"

	"github.com/seal-io/seal/pkg/dao/types/status"
)

// 1. define resource with status
type ExampleObject struct {
	Status status.Status
}

// 2. define condition
const (
	// finished condition
	ApplicationInstanceStatusDeployed status.ConditionType = "Deployed"
	// finished condition
	ApplicationInstanceStatusReady status.ConditionType = "Ready"

	// transitioning display
	ApplicationInstanceStatusDeployedTransitioning string = "Deploying"
)

// 3. create summarizer and define default condition
var ExampleSummarize *Summarizer

func Example_summary() {
	// 4. decide condition type and add to summarizer in sequence
	/*

	   For condition type ApplicationInstanceStatusDeployed is ErrorFalseTransitioningUnknown, so should call AddErrorFalseTransitioningUnknown

	   Success:
	   {
	   	"conditionType": "Deployed",
	   	"conditionStatus": "True"
	   }

	   Error:
	   {
	   	"conditionType": "Deployed",
	   	"conditionStatus": "False",
	   	"message": "error happened"
	   }

	   In transitioning:
	   {
	   	"conditionType": "Deployed",
	   	"conditionStatus": "Unknown",
	   	"message": "deploying"
	   }

	*/
	ExampleSummarize = NewSummarizer(ApplicationInstanceStatusReady)
	ExampleSummarize.AddErrorFalseTransitioningUnknown(ApplicationInstanceStatusDeployed, ApplicationInstanceStatusDeployedTransitioning)

	obj := &ExampleObject{}
	// 5. set status to UnKnown while action begin, so object will in transitioning
	ApplicationInstanceStatusDeployed.Unknown(obj, "")

	// 6. do actual logic, set status to False while error occur, set status to True while no error
	do := func() error {
		return fmt.Errorf("example error")
	}
	if err := do(); err != nil {
		ApplicationInstanceStatusDeployed.False(obj, err.Error())
	}

	// 6. summarize will set the summary status to obj
	ExampleSummarize.SetSummarize(&obj.Status)

	// 7. the summary status should include error
	fmt.Println(obj.Status.Status, obj.Status.Error, obj.Status.StatusMessage)
	// Output: Deployed true example error
}
