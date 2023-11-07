package workflow

import (
	"context"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/workflowstageexecution"
	"github.com/seal-io/walrus/pkg/dao/model/workflowstepexecution"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/datalisten/modelchange"
	"github.com/seal-io/walrus/utils/topic"
)

// ResetWorkflowExecutionStatus resets the workflow execution status to pending.
func ResetWorkflowExecutionStatus(
	ctx context.Context,
	mc model.ClientSet,
	workflowExecution *model.WorkflowExecution,
) error {
	status.WorkflowExecutionStatusPending.Reset(workflowExecution, "")
	workflowExecution.Status.SetSummary(status.WalkWorkflowExecution(&workflowExecution.Status))

	err := mc.WorkflowExecutions().UpdateOne(workflowExecution).
		SetStatus(workflowExecution.Status).
		AddTimes(1).
		Exec(ctx)
	if err != nil {
		return err
	}

	stageExecutions, err := mc.WorkflowStageExecutions().Query().
		Select(
			workflowstageexecution.FieldID,
			workflowstageexecution.FieldStatus,
			workflowstageexecution.FieldWorkflowExecutionID,
		).
		WithSteps(func(wseq *model.WorkflowStepExecutionQuery) {
			wseq.Select(
				workflowstepexecution.FieldID,
				workflowstepexecution.FieldStatus,
				workflowstepexecution.FieldWorkflowStageExecutionID,
			)
		}).
		Where(workflowstageexecution.WorkflowExecutionID(workflowExecution.ID)).
		All(ctx)
	if err != nil {
		return err
	}

	for i := range stageExecutions {
		stageExecution := stageExecutions[i]
		status.WorkflowStageExecutionStatusPending.Reset(stageExecution, "")
		stageExecution.Status.SetSummary(status.WalkWorkflowStageExecution(&stageExecution.Status))

		err = mc.WorkflowStageExecutions().UpdateOne(stageExecution).
			SetStatus(stageExecution.Status).
			Exec(ctx)
		if err != nil {
			return err
		}

		for i := range stageExecution.Edges.Steps {
			stepExecution := stageExecution.Edges.Steps[i]
			status.WorkflowStepExecutionStatusPending.Reset(stepExecution, "")
			stepExecution.Status.SetSummary(status.WalkWorkflowStepExecution(&stepExecution.Status))

			err = mc.WorkflowStepExecutions().UpdateOne(stepExecution).
				SetStatus(stepExecution.Status).
				Exec(ctx)
			if err != nil {
				return err
			}
		}
	}

	return topic.Publish(ctx, modelchange.WorkflowExecution, modelchange.Event{
		Type: modelchange.EventTypeUpdate,
		IDs:  []object.ID{workflowExecution.ID},
	})
}
