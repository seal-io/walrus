package workflowstepexecution

import (
	"context"
	"time"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/workflowstepexecution"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/datalisten/modelchange"
	pkgworkflow "github.com/seal-io/walrus/pkg/workflow"
	"github.com/seal-io/walrus/utils/gopool"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/topic"
)

func (h Handler) Get(req GetRequest) (GetResponse, error) {
	entity, err := h.modelClient.WorkflowStepExecutions().Query().
		Where(workflowstepexecution.ID(req.ID)).
		Only(req.Context)
	if err != nil {
		return nil, err
	}

	return model.ExposeWorkflowStepExecution(entity), nil
}

func (h Handler) Update(req UpdateRequest) error {
	entity, err := h.modelClient.WorkflowStepExecutions().Query().
		Where(workflowstepexecution.ID(req.ID)).
		Only(req.Context)
	if err != nil {
		return err
	}

	update := h.modelClient.WorkflowStepExecutions().UpdateOne(entity)

	switch req.Status {
	case types.ExecutionStatusSucceeded:
		status.WorkflowStepExecutionStatusRunning.True(entity, "")
	case types.ExecutionStatusFailed, types.ExecutionStatusError:
		status.WorkflowStepExecutionStatusRunning.False(entity, "")
	case types.ExecutionStatusRunning:
		status.WorkflowExecutionStatusPending.True(entity, "")
		status.WorkflowStepExecutionStatusRunning.Unknown(entity, "")

		update.SetExecuteTime(time.Now())
	default:
		log.WithName("workflowstepexecution").Error(nil, "unknown workflow step execution status", "status", req.Status)
		return nil
	}

	entity.Status.SetSummary(status.WalkWorkflowStepExecution(&entity.Status))

	update.SetStatus(entity.Status)

	// If the workflow step execution is not running, set the duration.
	if req.Status != types.ExecutionStatusRunning {
		update.SetDuration(int(time.Since(entity.ExecuteTime).Seconds()))
	}

	entity, err = update.Save(req.Context)
	if err != nil {
		return err
	}

	// Publish workflow execution topic,
	// step execution update will trigger workflow execution update.
	err = topic.Publish(req.Context, modelchange.WorkflowExecution, modelchange.Event{
		Type: modelchange.EventTypeUpdate,
		IDs:  []object.ID{entity.WorkflowExecutionID},
	})
	if err != nil {
		return err
	}

	gopool.Go(func() {
		logger := log.WithName("workflowstepexecution")
		subCtx := context.Background()
		// If the record is empty, get it from workflow step logs from pod.
		rerr := pkgworkflow.ArchiveWorkflowStepExecutionLogs(subCtx, pkgworkflow.StepExecutionLogOptions{
			RestCfg:       h.k8sConfig,
			ModelClient:   h.modelClient,
			StepExecution: entity,
		})
		if rerr != nil {
			logger.Error(rerr, "failed to set workflow step execution logs")
		}
	})

	return nil
}
