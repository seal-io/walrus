package workflowexecution

import (
	"fmt"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/workflowexecution"
	"github.com/seal-io/walrus/pkg/dao/model/workflowstageexecution"
	"github.com/seal-io/walrus/pkg/dao/model/workflowstepexecution"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	pkgworkflow "github.com/seal-io/walrus/pkg/workflow"
)

func (h Handler) RouteRerunRequest(req RouteRerunRequest) error {
	entity, err := h.modelClient.WorkflowExecutions().Query().
		Where(workflowexecution.ID(req.ID)).
		WithStages(func(wsgq *model.WorkflowStageExecutionQuery) {
			wsgq.WithSteps(func(wseq *model.WorkflowStepExecutionQuery) {
				wseq.
					Select(workflowstepexecution.WithoutFields(workflowstepexecution.FieldRecord)...).
					Order(model.Asc(workflowstepexecution.FieldOrder))
			}).
				Order(model.Asc(workflowstageexecution.FieldOrder))
		}).
		Only(req.Context)
	if err != nil {
		return err
	}

	if status.WorkflowExecutionStatusPending.IsUnknown(entity) ||
		status.WorkflowExecutionStatusRunning.IsUnknown(entity) {
		return fmt.Errorf("workflow execution is pending or running")
	}

	return h.modelClient.WithTx(req.Context, func(tx *model.Tx) error {
		return pkgworkflow.Rerun(req.Context, h.modelClient, h.workflowClient, entity)
	})
}

// RouteStopRequest terminates the workflow execution.
func (h Handler) RouteStopRequest(req RouteStopRequest) error {
	entity, err := h.modelClient.WorkflowExecutions().Query().
		Where(workflowexecution.ID(req.ID)).
		Only(req.Context)
	if err != nil {
		return err
	}

	return h.workflowClient.Terminate(req.Context, pkgworkflow.TerminateOptions{
		WorkflowExecution: entity,
	})
}
