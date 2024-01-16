package status

const (
	WorkflowExecutionStatusPending  ConditionType = "Pending"
	WorkflowExecutionStatusRunning  ConditionType = "Running"
	WorkflowExecutionStatusCanceled ConditionType = "Canceled"
)

// workflowExecutionStatusPaths makes the following decision.
//
//	|  Condition Type  |     Condition Status    | Human Readable Status | Human Sensible Status |
//	| ---------------- | ----------------------- | --------------------- | --------------------- |
//	| Pending          | Unknown                 | Pending               | Transitioning         |
//	| Pending          | False                   | Failed                | Error                 |
//	| Running          | Unknown                 | Running               | Transitioning         |
//	| Running          | False                   | Failed                | Error                 |
//	| Running          | True                    | Completed             | Completed             |
//	| Canceled         | Unknown                 | Canceling             | Transitioning         |
//	| Canceled         | False                   | CancelFailed          | Error                 |
//	| Canceled         | True                    | Canceled              | Canceled              |
var workflowExecutionStatusPaths = NewWalker(
	[][]ConditionType{
		{
			WorkflowExecutionStatusPending,
			WorkflowExecutionStatusRunning,
			WorkflowExecutionStatusCanceled,
		},
	},
	func(d Decision[ConditionType]) {
		d.Make(WorkflowExecutionStatusCanceled,
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

func WalkWorkflowExecution(st *Status) *Summary {
	return workflowExecutionStatusPaths.Walk(st)
}

const (
	WorkflowStageExecutionStatusPending  ConditionType = "Pending"
	WorkflowStageExecutionStatusRunning  ConditionType = "Running"
	WorkflowStageExecutionStatusCanceled ConditionType = "Canceled"
)

// workflowStageExecutionStatusPaths makes the following decision.
//
//	|  Condition Type  |     Condition Status    | Human Readable Status | Human Sensible Status |
//	| ---------------- | ----------------------- | --------------------- | --------------------- |
//	| Pending          | Unknown                 | Pending               | Transitioning         |
//	| Pending          | False                   | Failed                | Error                 |
//	| Running          | Unknown                 | Running               | Transitioning         |
//	| Running          | False                   | Failed                / Error                 |
//	| Running          | True                    | Running               | Completed             |
//	| Canceled         | Unknown                 | Canceling             | Transitioning         |
//	| Canceled         | False                   | CancelFailed          | Error                 |
//	| Canceled         | True                    | Canceled              | Canceled              |
var workflowStageExecutionStatusPaths = NewWalker(
	[][]ConditionType{
		{
			WorkflowStageExecutionStatusPending,
			WorkflowStageExecutionStatusRunning,
			WorkflowStageExecutionStatusCanceled,
		},
	},
	func(d Decision[ConditionType]) {
		d.Make(WorkflowStageExecutionStatusCanceled,
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

func WalkWorkflowStageExecution(st *Status) *Summary {
	return workflowStageExecutionStatusPaths.Walk(st)
}

const (
	WorkflowStepExecutionStatusPending  ConditionType = "Pending"
	WorkflowStepExecutionStatusRunning  ConditionType = "Running"
	WorkflowStepExecutionStatusCanceled ConditionType = "Canceled"
)

// workflowStepExecutionStatusPaths makes the following decision.
//
//	|  Condition Type  |     Condition Status    | Human Readable Status | Human Sensible Status |
//	| ---------------- | ----------------------- | --------------------- | --------------------- |
//	| Pending          | Unknown                 | Pending               | Transitioning         |
//	| Pending          | False                   | Failed                | Error                 |
//	| Running          | Unknown                 | Running               | Transitioning         |
//	| Running          | False                   | Failed                | Error                 |
//	| Running          | True                    | Running               | Completed             |
//	| Canceled         | Unknown                 | Canceling             | Transitioning         |
//	| Canceled         | False                   | CancelFailed          | Error                 |
//	| Canceled         | True                    | Canceled              | Canceled              |
var workflowStepExecutionStatusPaths = NewWalker(
	[][]ConditionType{
		{
			WorkflowStepExecutionStatusPending,
			WorkflowStepExecutionStatusRunning,
			WorkflowStepExecutionStatusCanceled,
		},
	},
	func(d Decision[ConditionType]) {
		d.Make(WorkflowStepExecutionStatusCanceled,
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

func WalkWorkflowStepExecution(st *Status) *Summary {
	return workflowStepExecutionStatusPaths.Walk(st)
}
