package status

const (
	ResourceRevisionStatusReady ConditionType = "Ready"

	ResourceRevisionSummaryStatusRunning string = "Running"
	ResourceRevisionSummaryStatusFailed  string = "Failed"
	ResourceRevisionSummaryStatusSucceed string = "Succeeded"
)

// resourceRevisionStatusPaths makes the following decision.
//
// |  Condition Type  |     Condition Status    | Human Readable Status | Human Sensible Status |
// | ---------------- | ----------------------- | --------------------- | --------------------- |
// | Ready            | Unknown                 | Running               | Transitioning         |
// | Ready            | False                   | Failed                | Error                 |
// | Ready            | True                    | Succeeded               |                       |.
var resourceRevisionStatusPaths = NewWalker(
	[][]ConditionType{
		{
			ResourceRevisionStatusReady,
		},
	},
	func(d Decision[ConditionType]) {
		d.Make(ResourceRevisionStatusReady,
			func(st ConditionStatus, reason string) *Summary {
				switch st {
				case ConditionStatusUnknown:
					return &Summary{
						SummaryStatus: ResourceRevisionSummaryStatusRunning,
						Transitioning: true,
					}
				case ConditionStatusFalse:
					return &Summary{
						SummaryStatus: ResourceRevisionSummaryStatusFailed,
						Error:         true,
					}
				}
				return &Summary{SummaryStatus: ResourceRevisionSummaryStatusSucceed}
			})
	},
)

func WalkResourceRevision(st *Status) *Summary {
	return resourceRevisionStatusPaths.Walk(st)
}
