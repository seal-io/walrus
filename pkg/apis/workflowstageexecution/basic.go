package workflowstageexecution

import (
	"time"

	"github.com/seal-io/walrus/pkg/dao/model/workflowstageexecution"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/datalisten/modelchange"
	"github.com/seal-io/walrus/utils/topic"
)

func (h Handler) Update(req UpdateRequest) error {
	entity, err := h.modelClient.WorkflowStageExecutions().Query().
		Where(workflowstageexecution.ID(req.ID)).
		Only(req.Context)
	if err != nil {
		return err
	}

	update := h.modelClient.WorkflowStageExecutions().UpdateOne(entity)

	switch req.Status {
	case types.ExecutionStatusSucceeded:
		status.WorkflowStageExecutionStatusRunning.True(entity, "")
	case types.ExecutionStatusFailed, types.ExecutionStatusError:
		status.WorkflowStageExecutionStatusRunning.False(entity, "")
	case types.ExecutionStatusRunning:
		status.WorkflowExecutionStatusPending.True(entity, "")
		status.WorkflowStageExecutionStatusRunning.Unknown(entity, "")
	default:
		return nil
	}

	entity.Status.SetSummary(status.WalkWorkflowStageExecution(&entity.Status))
	update.SetStatus(entity.Status)

	// If the workflow stage execution is not running, set the duration.
	if req.Status != types.ExecutionStatusRunning {
		update.SetDuration(int(time.Since(entity.ExecuteTime).Seconds()))
	}

	entity, err = update.Save(req.Context)
	if err != nil {
		return err
	}

	// Stage execution update will trigger workflow execution topic.
	err = topic.Publish(req.Context, modelchange.WorkflowExecution, modelchange.Event{
		Type: modelchange.EventTypeUpdate,
		IDs:  []object.ID{entity.WorkflowExecutionID},
	})
	if err != nil {
		return err
	}

	// Stage execution update will trigger workflow update.
	return topic.Publish(req.Context, modelchange.Workflow, modelchange.Event{
		Type: modelchange.EventTypeUpdate,
		IDs:  []object.ID{entity.WorkflowID},
	})
}
