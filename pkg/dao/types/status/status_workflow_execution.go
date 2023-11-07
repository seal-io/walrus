package status

const (
	WorkflowExecutionStatusPending ConditionType = "Pending"
	WorkflowExecutionStatusRunning ConditionType = "Running"
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
var workflowExecutionStatusPaths = NewWalker(
	[][]ConditionType{
		{
			WorkflowExecutionStatusPending,
			WorkflowExecutionStatusRunning,
		},
	},
)

func WalkWorkflowExecution(st *Status) *Summary {
	return workflowExecutionStatusPaths.Walk(st)
}

const (
	WorkflowStageExecutionStatusPending ConditionType = "Pending"
	WorkflowStageExecutionStatusRunning ConditionType = "Running"
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
var workflowStageExecutionStatusPaths = NewWalker(
	[][]ConditionType{
		{
			WorkflowStageExecutionStatusPending,
			WorkflowStageExecutionStatusRunning,
		},
	},
)

func WalkWorkflowStageExecution(st *Status) *Summary {
	return workflowStageExecutionStatusPaths.Walk(st)
}

const (
	WorkflowStepExecutionStatusPending ConditionType = "Pending"
	WorkflowStepExecutionStatusRunning ConditionType = "Running"
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
var workflowStepExecutionStatusPaths = NewWalker(
	[][]ConditionType{
		{
			WorkflowStepExecutionStatusPending,
			WorkflowStepExecutionStatusRunning,
		},
	},
)

func WalkWorkflowStepExecution(st *Status) *Summary {
	return workflowStepExecutionStatusPaths.Walk(st)
}
