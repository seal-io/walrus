package workflow

import (
	"context"
	"time"

	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"

	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/workflowstageexecution"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/datalisten/modelchange"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/topic"
)

// ResetWorkflowExecutionStatus resets the workflow execution status and all its stage and step execution status.
// Rerun the workflow execution need reset the running information.
func ResetWorkflowExecutionStatus(
	ctx context.Context,
	mc model.ClientSet,
	workflowExecution *model.WorkflowExecution,
) error {
	err := dao.ResetWorkflowExecution(ctx, mc, workflowExecution)
	if err != nil {
		return err
	}

	stageExecutions, err := mc.WorkflowStageExecutions().Query().
		Select(
			workflowstageexecution.FieldID,
			workflowstageexecution.FieldStatus,
			workflowstageexecution.FieldWorkflowExecutionID,
		).
		WithSteps().
		Where(workflowstageexecution.WorkflowExecutionID(workflowExecution.ID)).
		All(ctx)
	if err != nil {
		return err
	}

	for i := range stageExecutions {
		stageExecution := stageExecutions[i]

		err = dao.ResetWorkflowStageExecution(ctx, mc, stageExecution)
		if err != nil {
			return err
		}

		for i := range stageExecution.Edges.Steps {
			stepExecution := stageExecution.Edges.Steps[i]

			err = dao.ResetWorkflowStepExecution(ctx, mc, stepExecution)
			if err != nil {
				return err
			}
		}
	}

	return topic.Publish(ctx, modelchange.Workflow, modelchange.Event{
		Type: modelchange.EventTypeUpdate,
		Data: []modelchange.EventData{{ID: workflowExecution.WorkflowID}},
	})
}

// StatusSyncer sync the status of workflow execution.
type StatusSyncer struct {
	Logger         log.Logger
	ModelClient    model.ClientSet
	WorkflowClient Client
}

func NewStatusSyncer(
	mc model.ClientSet,
	wc Client,
) *StatusSyncer {
	return &StatusSyncer{
		Logger:         log.WithName("workflow").WithName("status-syncer"),
		ModelClient:    mc,
		WorkflowClient: wc,
	}
}

// SyncWorkflowExecutionStatus syncs the status of workflow execution.
func (m *StatusSyncer) SyncWorkflowExecutionStatus(ctx context.Context, wf *wfv1.Workflow) error {
	workflowExecutionID, ok := wf.Labels[workflowExecutionIDLabel]
	if !ok {
		return nil
	}

	we, err := m.ModelClient.WorkflowExecutions().Get(ctx, object.ID(workflowExecutionID))
	if err != nil && !model.IsNotFound(err) {
		return err
	}

	if we == nil {
		return nil
	}

	update := m.ModelClient.WorkflowExecutions().UpdateOne(we)

	switch wf.Status.Phase {
	case wfv1.WorkflowRunning:
		// If the workflow is running, no need to update status.
		if status.WorkflowExecutionStatusRunning.IsUnknown(we) {
			return nil
		}

		status.WorkflowExecutionStatusPending.True(we, "")
		status.WorkflowExecutionStatusRunning.Unknown(we, "")
		update.SetExecuteTime(wf.Status.StartedAt.Time)

		m.Logger.Debugf("workflow %s is running", we.ID.String())

	case wfv1.WorkflowSucceeded:
		status.WorkflowExecutionStatusPending.True(we, "")
		status.WorkflowExecutionStatusRunning.True(we, "")
		update.SetDuration(int(wf.Status.GetDuration().Seconds()))

		m.Logger.Debugf("workflow %s is succeeded", we.ID.String())

	case wfv1.WorkflowFailed, wfv1.WorkflowError:
		switch {
		case status.WorkflowExecutionStatusCanceled.IsUnknown(we):
			status.WorkflowExecutionStatusCanceled.True(we, wf.Status.Message)
		default:
			status.WorkflowExecutionStatusPending.True(we, "")
			status.WorkflowExecutionStatusRunning.False(we, wf.Status.Message)
		}

		update.SetDuration(int(time.Since(we.ExecuteTime).Seconds()))

		m.Logger.Debugf("workflow %s is failed", we.ID.String())
	default:
		m.Logger.Debugf("workflow %s phase is %s, skip it.", we.ID, wf.Status.Phase)
		return nil
	}

	we.Status.SetSummary(status.WalkWorkflowExecution(&we.Status))

	we, err = update.SetStatus(we.Status).
		Save(ctx)
	if err != nil {
		return err
	}

	// Workflow execution update will trigger workflow topic.
	return topic.Publish(ctx, modelchange.WorkflowExecution, modelchange.Event{
		Type: modelchange.EventTypeUpdate,
		Data: []modelchange.EventData{{ID: we.ID}},
	})
}

// SyncStageExecutionStatus syncs workflow stage execution status.
func (m *StatusSyncer) SyncStageExecutionStatus(
	ctx context.Context,
	node wfv1.NodeStatus,
	stageExecutionID object.ID,
	canceled bool,
) error {
	wse, err := m.ModelClient.WorkflowStageExecutions().Get(ctx, stageExecutionID)
	if err != nil {
		return err
	}

	// Finished stage will not be updated.
	if status.WorkflowStageExecutionStatusRunning.IsTrue(wse) ||
		status.WorkflowStageExecutionStatusRunning.IsFalse(wse) {
		return nil
	}

	update := m.ModelClient.WorkflowStageExecutions().UpdateOne(wse)

	switch node.Phase {
	case wfv1.NodeRunning:
		if status.WorkflowStageExecutionStatusRunning.IsUnknown(wse) {
			return nil
		}

		status.WorkflowStageExecutionStatusPending.True(wse, "")
		status.WorkflowStageExecutionStatusRunning.Unknown(wse, "")
		update.SetExecuteTime(node.StartedAt.Time)
		m.Logger.Debugf("stage %s is running", wse.ID.String())
	case wfv1.NodeSucceeded:
		status.WorkflowStageExecutionStatusPending.True(wse, "")
		status.WorkflowStageExecutionStatusRunning.True(wse, "")
		m.Logger.Debugf("stage %s is succeeded", wse.ID.String())

	case wfv1.NodeSkipped, wfv1.NodeOmitted:
		m.Logger.Debugf("stage %s is skipped or omitted", wse.ID.String(), "message", node.Message)
		return nil
	case wfv1.NodeFailed, wfv1.NodeError:
		if canceled {
			status.WorkflowStageExecutionStatusCanceled.Reset(wse, "")
			status.WorkflowStageExecutionStatusCanceled.True(wse, "")
		} else {
			status.WorkflowStageExecutionStatusPending.True(wse, "")
			status.WorkflowStageExecutionStatusRunning.False(wse, node.Message)
			m.Logger.Debugf("stage %s is failed", wse.ID.String())
		}

	default:
		m.Logger.Debugf("stage %s phase is %s, skip it.", wse.ID, node.Phase)
		return nil
	}

	wse.Status.SetSummary(status.WalkWorkflowStageExecution(&wse.Status))

	if node.Completed() {
		finishTime := node.FinishTime().Time

		if finishTime.IsZero() {
			finishTime = time.Now()
		}

		if wse.ExecuteTime.IsZero() {
			if !node.StartTime().Time.IsZero() {
				wse.ExecuteTime = node.StartTime().Time
			}

			if wse.ExecuteTime.IsZero() {
				wse.ExecuteTime = *wse.CreateTime
			}

			update.SetExecuteTime(wse.ExecuteTime)
		}

		update.SetDuration(int(finishTime.Sub(wse.ExecuteTime).Seconds()))
	}

	wse, err = update.SetStatus(wse.Status).
		Save(ctx)
	if err != nil {
		return err
	}

	// Stage execution update will trigger workflow execution topic.
	err = topic.Publish(ctx, modelchange.WorkflowExecution, modelchange.Event{
		Type: modelchange.EventTypeUpdate,
		Data: []modelchange.EventData{{ID: wse.WorkflowExecutionID}},
	})
	if err != nil {
		return err
	}

	// Stage execution update will trigger workflow topic.
	return topic.Publish(ctx, modelchange.Workflow, modelchange.Event{
		Type: modelchange.EventTypeUpdate,
		Data: []modelchange.EventData{{ID: wse.WorkflowID}},
	})
}

// SyncStepExecutionStatus syncs workflow step execution status.
func (m *StatusSyncer) SyncStepExecutionStatus(
	ctx context.Context,
	node wfv1.NodeStatus,
	stepExecutionID object.ID,
	canceled bool,
) error {
	wse, err := m.ModelClient.WorkflowStepExecutions().Get(ctx, stepExecutionID)
	if err != nil {
		return err
	}

	we, err := m.ModelClient.WorkflowExecutions().Get(ctx, wse.WorkflowExecutionID)
	if err != nil {
		return err
	}

	// Finished step will not be updated.
	if status.WorkflowStepExecutionStatusRunning.IsTrue(wse) ||
		status.WorkflowStepExecutionStatusRunning.IsFalse(wse) {
		return nil
	}

	update := m.ModelClient.WorkflowStepExecutions().UpdateOne(wse)

	switch node.Phase {
	case wfv1.NodeRunning:
		if status.WorkflowStepExecutionStatusRunning.IsUnknown(wse) {
			return nil
		}

		status.WorkflowStepExecutionStatusPending.True(wse, "")
		status.WorkflowStepExecutionStatusRunning.Unknown(wse, "")
		update.SetExecuteTime(node.StartedAt.Time)

		m.Logger.Debugf("step %s is running", wse.ID.String())
	case wfv1.NodeSucceeded:
		status.WorkflowStepExecutionStatusPending.True(wse, "")
		status.WorkflowStepExecutionStatusRunning.True(wse, "")
		m.Logger.Debugf("step %s is succeeded", wse.ID.String())

	case wfv1.NodeSkipped, wfv1.NodeOmitted:
		m.Logger.Debugf("step %s is skipped or omitted, message %s", wse.ID.String(), node.Message)

		return nil
	case wfv1.NodeFailed, wfv1.NodeError:
		if canceled {
			status.WorkflowStepExecutionStatusCanceled.Reset(wse, "")
			status.WorkflowStepExecutionStatusCanceled.True(wse, "")
		} else {
			message := node.Message

			if wse.Type == types.WorkflowStepTypeApproval {
				message, err = dao.GetRejectMessage(ctx, m.ModelClient, wse)
				if err != nil {
					return err
				}
			}

			status.WorkflowStepExecutionStatusPending.True(wse, "")
			status.WorkflowStepExecutionStatusRunning.False(wse, message)
			m.Logger.Debugf("step %s is failed", wse.ID.String())
		}

	default:
		m.Logger.Debugf("step  %s phase is %s, status not will not change.", wse.ID, node.Phase)
		return nil
	}

	wse.Status.SetSummary(status.WalkWorkflowStepExecution(&wse.Status))

	// If the step is completed, set logs to record.
	if node.Completed() {
		finishTime := node.FinishTime().Time

		if finishTime.IsZero() {
			finishTime = time.Now()
		}

		if wse.ExecuteTime.IsZero() {
			if !node.StartedAt.Time.IsZero() {
				wse.ExecuteTime = node.StartedAt.Time
			}

			if wse.ExecuteTime.IsZero() {
				wse.ExecuteTime = *wse.CreateTime
			}

			update.SetExecuteTime(wse.ExecuteTime)
		}

		update.SetDuration(int(finishTime.Sub(wse.ExecuteTime).Seconds()))

		logs, err := m.WorkflowClient.GetLogs(ctx, LogsOptions{
			WorkflowExecution: we,
			StepExecution:     wse,
		})
		if err != nil {
			return err
		}

		if len(logs) > 0 {
			update.SetRecord(string(logs))
		}
	}

	err = update.SetStatus(wse.Status).
		Exec(ctx)
	if err != nil {
		return err
	}

	// Step execution update will trigger workflow execution topic.
	return topic.Publish(ctx, modelchange.WorkflowExecution, modelchange.Event{
		Type: modelchange.EventTypeUpdate,
		Data: []modelchange.EventData{{ID: we.ID}},
	})
}

func (m *StatusSyncer) IsCanceled(ctx context.Context, wf *wfv1.Workflow) (bool, error) {
	workflowExecutionID, ok := wf.Labels[workflowExecutionIDLabel]
	if !ok {
		return false, nil
	}

	we, err := m.ModelClient.WorkflowExecutions().Get(ctx, object.ID(workflowExecutionID))
	if err != nil && !model.IsNotFound(err) {
		return false, err
	}

	if we == nil {
		return false, nil
	}

	return status.WorkflowExecutionStatusCanceled.IsTrue(we), nil
}
