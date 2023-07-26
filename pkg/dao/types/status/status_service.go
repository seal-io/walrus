package status

const (
	ServiceStatusDeployed    ConditionType = "Deployed"
	ServiceStatusDeleted     ConditionType = "Deleted"
	ServiceStatusReady       ConditionType = "Ready"
	ServiceStatusProgressing ConditionType = "Progressing"
	ServiceStatusSynced      ConditionType = "Synced"
	ServiceStatusDetected    ConditionType = "Detected"
)

// serviceStatusPaths makes the following decision.
//
//	|  Condition Type  |     Condition Status    | Human Readable Status | Human Sensible Status |
//	| ---------------- | ----------------------- | --------------------- | --------------------- |
//	| Progressing      | Unknown                 | Progressing           | Transitioning         |
//	| Progressing      | False                   | Progressing           | Error                 |
//	| Progressing      | True                    | Progressed            |                       |
//	| Deployed         | Unknown                 | Deploying             | Transitioning         |
//	| Deployed         | False                   | DeployFailed          | Error                 |
//	| Deployed         | True                    | Deployed              |                       |
//	| Ready            | Unknown                 | Preparing             | Transitioning         |
//	| Ready            | False                   | NotReady              | Error                 |
//	| Ready            | True                    | Ready                 |                       |
//	| Synced           | Unknown                 | Syncing               | Transitioning         |
//	| Synced		   | False                   | SyncFailed            | Error                 |
//	| Synced		   | True                    | Synced                |                       |
//	| Detected         | Unknown                 | Detecting             | Transitioning         |
//	| Detected         | False                   | DetectFailed          | Error                 |
//	| Detected         | True                    | Detected              |                       |
//	| Deleted          | Unknown                 | Deleting              | Transitioning         |
//	| Deleted          | False                   | DeleteFailed          | Error                 |
//	| Deleted          | True                    | Deleted               |                       |
var serviceStatusPaths = NewWalker(
	[][]ConditionType{
		{
			ServiceStatusDeleted,
			ServiceStatusProgressing,
			ServiceStatusDeployed,
			ServiceStatusReady,
		},
	},
	func(d Decision[ConditionType]) {
		d.Make(ServiceStatusDeleted,
			func(st ConditionStatus, reason string) (display string, isError, isTransitioning bool) {
				switch st {
				case ConditionStatusUnknown:
					return "Deleting", false, true
				case ConditionStatusFalse:
					return "DeleteFailed", true, false
				}
				return "", false, false
			})
	},
)

// WalkService walks the given status by service flow.
func WalkService(st *Status) *Summary {
	return serviceStatusPaths.Walk(st)
}
