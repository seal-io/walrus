package status

const (
	ResourceStatusUnDeployed  ConditionType = "Undeployed"
	ResourceStatusStopped     ConditionType = "Stopped"
	ResourceStatusPlanned     ConditionType = "Planned"
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
//	| Planned          | Unknown                 | Planning              | Transitioning         |
//	| Planned          | False                   | Planned               | Error                 |
//	| Planned          | True                    | Planned               |                       |
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
			ResourceStatusPlanned,
			ResourceStatusDeployed,
			ResourceStatusUnDeployed,
			ResourceStatusStopped,
			ResourceStatusReady,
		},
	},
	func(d Decision[ConditionType]) {
		d.Make(ResourceStatusDeleted,
			func(st ConditionStatus, reason string) *Summary {
				switch st {
				case ConditionStatusUnknown:
					return &Summary{
						SummaryStatus: "Deleting",
						Transitioning: true,
					}
				case ConditionStatusFalse:
					return &Summary{
						SummaryStatus: "DeleteFailed",
						Error:         true,
					}
				}
				return &Summary{}
			})
		d.Make(ResourceStatusUnDeployed,
			func(st ConditionStatus, reason string) *Summary {
				switch st {
				case ConditionStatusUnknown:
					return &Summary{
						SummaryStatus: "Transitioning",
						Transitioning: true,
					}
				case ConditionStatusFalse:
					return &Summary{
						SummaryStatus: "Error",
						Error:         true,
					}
				}
				return &Summary{SummaryStatus: "Undeployed", Inactive: true}
			})
		d.Make(ResourceStatusStopped,
			func(st ConditionStatus, reason string) *Summary {
				switch st {
				case ConditionStatusUnknown:
					return &Summary{
						SummaryStatus: "Stopping",
						Transitioning: true,
					}
				case ConditionStatusFalse:
					return &Summary{
						SummaryStatus: "StopFailed",
						Error:         true,
					}
				}
				return &Summary{SummaryStatus: "Stopped", Inactive: true}
			})
	},
)

// WalkResource walks the given status by resource flow.
func WalkResource(st *Status) *Summary {
	return resourceStatusPaths.Walk(st)
}
