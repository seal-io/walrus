package status

const (
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
//	| Ready            | Unknown                 | Preparing             | Transitioning         |
//	| Ready            | False                   | NotReady               | Error                 |
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
	},
)

// WalkResource walks the given status by resource flow.
func WalkResource(st *Status) *Summary {
	return resourceStatusPaths.Walk(st)
}
