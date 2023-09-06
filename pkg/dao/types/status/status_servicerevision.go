package status

const (
	ServiceRevisionStatusReady ConditionType = "Ready"

	ServiceRevisionSummaryStatusRunning string = "Running"
	ServiceRevisionSummaryStatusFailed  string = "Failed"
	ServiceRevisionSummaryStatusSucceed string = "Succeed"
)

// serviceRevisionStatusPaths makes the following decision.
//
// |  Condition Type  |     Condition Status    | Human Readable Status | Human Sensible Status |
// | ---------------- | ----------------------- | --------------------- | --------------------- |
// | Ready            | Unknown                 | Running               | Transitioning         |
// | Ready            | False                   | Failed                | Error                 |
// | Ready            | True                    | Succeed               |                       |.
var serviceRevisionStatusPaths = NewWalker(
	[][]ConditionType{
		{
			ServiceRevisionStatusReady,
		},
	},
	func(d Decision[ConditionType]) {
		d.Make(ServiceRevisionStatusReady,
			func(st ConditionStatus, reason string) (display string, isError, isTransitioning bool) {
				switch st {
				case ConditionStatusUnknown:
					return ServiceRevisionSummaryStatusRunning, false, true
				case ConditionStatusFalse:
					return ServiceRevisionSummaryStatusFailed, true, false
				}
				return ServiceRevisionSummaryStatusSucceed, false, false
			})
	},
)

func WalkServiceRevision(st *Status) *Summary {
	return serviceRevisionStatusPaths.Walk(st)
}
