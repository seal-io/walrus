package status

const (
	ResourceRevisionStatusReady ConditionType = "Ready"

	ResourceRevisionSummaryStatusRunning string = "Running"
	ResourceRevisionSummaryStatusFailed  string = "Failed"
	ResourceRevisionSummaryStatusSucceed string = "Succeed"
)

// resourceRevisionStatusPaths makes the following decision.
//
// |  Condition Type  |     Condition Status    | Human Readable Status | Human Sensible Status |
// | ---------------- | ----------------------- | --------------------- | --------------------- |
// | Ready            | Unknown                 | Running               | Transitioning         |
// | Ready            | False                   | Failed                | Error                 |
// | Ready            | True                    | Succeed               |                       |.
var resourceRevisionStatusPaths = NewWalker(
	[][]ConditionType{
		{
			ResourceRevisionStatusReady,
		},
	},
	func(d Decision[ConditionType]) {
		d.Make(ResourceRevisionStatusReady,
			func(st ConditionStatus, reason string) (display string, isError, isTransitioning bool) {
				switch st {
				case ConditionStatusUnknown:
					return ResourceRevisionSummaryStatusRunning, false, true
				case ConditionStatusFalse:
					return ResourceRevisionSummaryStatusFailed, true, false
				}
				return ResourceRevisionSummaryStatusSucceed, false, false
			})
	},
)

func WalkResourceRevision(st *Status) *Summary {
	return resourceRevisionStatusPaths.Walk(st)
}
