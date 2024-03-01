package status

const (
	ResourceRunStatusPending  ConditionType = "Pending"
	ResourceRunStatusPlanned  ConditionType = "Planned"
	ResourceRunStatusApplied  ConditionType = "Applied"
	ResourceRunStatusCanceled ConditionType = "Canceled"

	ResourceRunSummaryStatusRunning string = "Running"
	ResourceRunSummaryStatusFailed  string = "Failed"
	ResourceRunSummaryStatusSucceed string = "Succeeded"
)

// resourceRunStatusPaths makes the following decision.
//
//	|  Condition Type  |     Condition Status    | Human Readable Status | Human Sensible Status |
//	| ---------------- | ----------------------- | --------------------- | --------------------- |
//	| Pending          | Unknown                 | Pending               | Transitioning         |
//	| Pending          | False                   | Failed                | Error                 |
//	| Plan             | Unknown                 | Planning              | Transitioning         |
//	| Plan             | False                   | Failed                | Error                 |
//	| Plan             | True                    | Planned               |                       |
//	| Applied          | Unknown                 | Running               | Transitioning         |
//	| Applied          | False                   | Failed                | Error                 |
//	| Applied          | True                    | Succeeded             |                       |
//	| Canceled         | True                    | Canceled              | Canceled              |
var resourceRunStatusPaths = NewWalker(
	[][]ConditionType{
		{
			ResourceRunStatusPending,
			ResourceRunStatusPlanned,
			ResourceRunStatusApplied,
			ResourceRunStatusCanceled,
		},
	},
	func(d Decision[ConditionType]) {
		d.Make(ResourceRunStatusCanceled,
			func(st ConditionStatus, reason string) *Summary {
				switch st {
				case ConditionStatusUnknown:
					return &Summary{
						SummaryStatus: "Canceling",
						Transitioning: true,
					}
				case ConditionStatusFalse:
					return &Summary{
						SummaryStatus: "CancelFailed",
						Error:         true,
					}
				}
				return &Summary{SummaryStatus: "Canceled", Inactive: true}
			})
	},
)

func WalkResourceRun(st *Status) *Summary {
	return resourceRunStatusPaths.Walk(st)
}
