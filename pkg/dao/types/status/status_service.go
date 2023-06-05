package status

const (
	ServiceStatusDeployed ConditionType = "Deployed"
	ServiceStatusDeleted  ConditionType = "Deleted"
	ServiceStatusReady    ConditionType = "Ready"
)

// serviceStatusPaths makes the following decision.
//
//	|  Condition Type  |     Condition Status    | Human Readable Status | Human Sensible Status |
//	| ---------------- | ----------------------- | --------------------- | --------------------- |
//	| Deployed         | Unknown                 | Deploying             | Transitioning         |
//	| Deployed         | False                   | DeployFailed          | Error                 |
//	| Deployed         | True                    | Deployed              |                       |
//	| Ready            | Unknown                 | Preparing             | Transitioning         |
//	| Ready            | False                   | Unready               | Error                 |
//	| Ready            | True                    | Ready                 |                       |
//	| Deleted          | Unknown                 | Deleting              | Transitioning         |
//	| Deleted          | False                   | DeleteFailed          | Error                 |
//	| Deleted          | True                    | Deleted               |                       |
var serviceStatusPaths = NewWalker(
	[][]ConditionType{
		{
			ServiceStatusDeployed,
			ServiceStatusReady,
		},
		{
			ServiceStatusDeleted,
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
				return "Deleted", false, false
			})
	},
)

// WalkService walks the given status by service flow.
func WalkService(st *Status) *Summary {
	return serviceStatusPaths.Walk(st)
}
