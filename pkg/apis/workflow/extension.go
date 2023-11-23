package workflow

import (
	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/workflow"
	"github.com/seal-io/walrus/pkg/dao/model/workflowexecution"
	"github.com/seal-io/walrus/pkg/dao/model/workflowstage"
	"github.com/seal-io/walrus/pkg/dao/model/workflowstageexecution"
	"github.com/seal-io/walrus/pkg/dao/model/workflowstep"
	"github.com/seal-io/walrus/pkg/dao/model/workflowstepexecution"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	pkgworkflow "github.com/seal-io/walrus/pkg/workflow"
	"github.com/seal-io/walrus/utils/log"
)

func (h Handler) RouteGetLatestExecutionRequest(req RouteGetLatestExecutionRequest) (
	RouteGetLatestExecutionResponse,
	error,
) {
	wf, err := h.modelClient.WorkflowExecutions().Query().
		Where(workflowexecution.WorkflowID(req.ID)).
		Order(model.Desc(workflowexecution.FieldCreateTime)).
		WithStages(func(wsq *model.WorkflowStageExecutionQuery) {
			wsq.WithSteps(func(weeq *model.WorkflowStepExecutionQuery) {
				weeq.Order(model.Asc(workflowstepexecution.FieldOrder))
			}).
				Order(model.Asc(workflowstageexecution.FieldOrder))
		}).
		First(req.Context)
	if err != nil && !model.IsNotFound(err) {
		return nil, err
	}

	if wf == nil {
		return nil, nil
	}

	return model.ExposeWorkflowExecution(wf), nil
}

func (h Handler) RouteRun(req RouteRunRequest) (RouteRunResponse, error) {
	var (
		logger = log.WithName("workflow")

		err error
		wf  *model.Workflow
		wfe *model.WorkflowExecution
	)

	err = h.modelClient.WithTx(req.Context, func(tx *model.Tx) error {
		wf, err = tx.Workflows().Query().
			Where(workflow.ID(req.ID)).
			ForUpdate().
			WithStages(func(wsq *model.WorkflowStageQuery) {
				wsq.WithSteps(func(wsq *model.WorkflowStepQuery) {
					wsq.Order(model.Asc(workflowstep.FieldOrder))
				}).Order(model.Asc(workflowstage.FieldOrder))
			}).
			Only(req.Context)
		if err != nil {
			return err
		}

		wfe, err = dao.CreateWorkflowExecution(req.Context, tx, dao.CreateWorkflowExecutionOptions{
			Workflow:    wf,
			Description: req.Description,
			Variables:   req.Variables,
		})
		if err != nil {
			return err
		}
		// Update workflow version.
		return tx.Workflows().UpdateOne(wf).
			AddVersion(1).
			Exec(req.Context)
	})
	if err != nil {
		return nil, err
	}

	err = pkgworkflow.Run(req.Context, h.modelClient, h.workflowClient, wfe)
	if err != nil {
		status.WorkflowExecutionStatusPending.False(wfe, err.Error())
		wfe.Status.SetSummary(status.WalkWorkflowExecution(&wfe.Status))

		rerr := h.modelClient.WorkflowExecutions().UpdateOne(wfe).
			SetStatus(wfe.Status).
			Exec(req.Context)
		if rerr != nil {
			logger.Errorf("failed to update workflow execution status: %v", rerr)
		}

		return nil, err
	}

	return model.ExposeWorkflowExecution(wfe), nil
}
