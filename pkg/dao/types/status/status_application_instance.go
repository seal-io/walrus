package status

const (
	ApplicationInstanceStatusDeployed ConditionType = "Deployed"
	ApplicationInstanceStatusDeleted  ConditionType = "Deleted"
	ApplicationInstanceStatusReady    ConditionType = "Ready"
)

// applicationInstanceStatusPaths makes the following decision.
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
var applicationInstanceStatusPaths = NewWalker(
	[][]ConditionType{
		{
			ApplicationInstanceStatusDeployed,
			ApplicationInstanceStatusReady,
		},
		{
			ApplicationInstanceStatusDeleted,
		},
	},
	func(d Decision[ConditionType]) {
		d.Make(ApplicationInstanceStatusDeleted,
			func(st ConditionStatus, reason string) (display string, isError bool, isTransitioning bool) {
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

// WalkApplicationInstance walks the given status by application instance flow.
func WalkApplicationInstance(st *Status) *Summary {
	return applicationInstanceStatusPaths.Walk(st)
}
