package status

const (
	ResourceRunStatusReady ConditionType = "Ready"

	ResourceRunSummaryStatusRunning string = "Running"
	ResourceRunSummaryStatusFailed  string = "Failed"
	ResourceRunSummaryStatusSucceed string = "Succeeded"
)

// resourceRunStatusPaths makes the following decision.
//
// |  Condition Type  |     Condition Status    | Human Readable Status | Human Sensible Status |
// | ---------------- | ----------------------- | --------------------- | --------------------- |
// | Ready            | Unknown                 | Running               | Transitioning         |
// | Ready            | False                   | Failed                | Error                 |
// | Ready            | True                    | Succeeded               |                       |.
var resourceRunStatusPaths = NewWalker(
	[][]ConditionType{
		{
			ResourceRunStatusReady,
		},
	},
	func(d Decision[ConditionType]) {
		d.Make(ResourceRunStatusReady,
			func(st ConditionStatus, reason string) *Summary {
				switch st {
				case ConditionStatusUnknown:
					return &Summary{
						SummaryStatus: ResourceRunSummaryStatusRunning,
						Transitioning: true,
					}
				case ConditionStatusFalse:
					return &Summary{
						SummaryStatus: ResourceRunSummaryStatusFailed,
						Error:         true,
					}
				}
				return &Summary{SummaryStatus: ResourceRunSummaryStatusSucceed}
			})
	},
)

func WalkResourceRun(st *Status) *Summary {
	return resourceRunStatusPaths.Walk(st)
}
