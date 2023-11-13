package status

const (
	ResourceStatusUnDeployed  ConditionType = "Undeployed"
	ResourceStatusStopped     ConditionType = "Stopped"
	ResourceStatusDeployed    ConditionType = "Deployed"
	ResourceStatusDeleted     ConditionType = "Deleted"
	ResourceStatusReady       ConditionType = "Ready"
	ResourceStatusProgressing ConditionType = "Progressing"
)

// resourceStatusPaths makes the following decision.
//
//	|  Condition Type  |     Condition Status    | Human Readable Status | Human Sensible Status |
//	| ---------------- | ----------------------- | --------------------- | --------------------- |
//	| Progressing      | Unknown                 | Progressing           | Transitioning         |
//	| Progressing      | False                   | Progressing           | Error                 |
//	| Progressing      | True                    | Progressed            |                       |
//	| Deployed         | Unknown                 | Deploying             | Transitioning         |
//	| Deployed         | False                   | DeployFailed          | Error                 |
//	| Deployed         | True                    | Deployed              |                       |
//	| UnDeployed       | Unknown                 | Transitioning         | Transitioning         |
//	| UnDeployed       | False                   | Error                 | Error                 |
//	| UnDeployed       | True                    | Undeployed            |                       |
//	| Stopped          | Unknown                 | Stopping              | Transitioning         |
//	| Stopped          | False                   | StopFailed            | Error                 |
//	| Stopped          | True                    | Stopped               |                       |
//	| Ready            | Unknown                 | Preparing             | Transitioning         |
//	| Ready            | False                   | NotReady              | Error                 |
//	| Ready            | True                    | Ready                 |                       |
//	| Deleted          | Unknown                 | Deleting              | Transitioning         |
//	| Deleted          | False                   | DeleteFailed          | Error                 |
//	| Deleted          | True                    | Deleted               |                       |
var resourceStatusPaths = NewWalker(
	[][]ConditionType{
		{
			ResourceStatusDeleted,
			ResourceStatusProgressing,
			ResourceStatusDeployed,
			ResourceStatusUnDeployed,
			ResourceStatusStopped,
			ResourceStatusReady,
		},
	},
	func(d Decision[ConditionType]) {
		d.Make(ResourceStatusDeleted,
			func(st ConditionStatus, reason string) (display string, isError, isTransitioning bool) {
				switch st {
				case ConditionStatusUnknown:
					return "Deleting", false, true
				case ConditionStatusFalse:
					return "DeleteFailed", true, false
				}
				return "", false, false
			})
		d.Make(ResourceStatusStopped,
			func(st ConditionStatus, reason string) (display string, isError, isTransitioning bool) {
				// Both Stopping and Stopped are considered as transitioning.
				switch st {
				case ConditionStatusUnknown:
					return "Stopping", false, true
				case ConditionStatusFalse:
					return "StopFailed", true, false
				}
				return "Stopped", false, true
			})
	},
)

// WalkResource walks the given status by resource flow.
func WalkResource(st *Status) *Summary {
	return resourceStatusPaths.Walk(st)
}
