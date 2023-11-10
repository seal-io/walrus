package workflowstepexecution

import (
	"context"
	"io"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/workflowexecution"
	"github.com/seal-io/walrus/pkg/dao/model/workflowstepexecution"
	"github.com/seal-io/walrus/pkg/workflow"
)

func (h Handler) RouteLog(req RouteLogRequest) error {
	var (
		ctx context.Context
		out io.Writer
	)

	if req.Stream != nil {
		ctx = req.Stream
		out = req.Stream
	} else {
		ctx = req.Context
		out = req.Context.Writer
	}

	wse, err := h.modelClient.WorkflowStepExecutions().Query().
		Where(workflowstepexecution.ID(req.ID)).
		Only(ctx)
	if err != nil {
		return err
	}

	return workflow.StreamWorkflowStepExecutionLogs(ctx, workflow.StreamWorkflowStepExecutionLogsOptions{
		StepExecutionLogOptions: workflow.StepExecutionLogOptions{
			RestCfg:       h.k8sConfig,
			ModelClient:   h.modelClient,
			StepExecution: wse,
		},
		Out: out,
	})
}

func (h Handler) RouteApprove(req RouteApproveRequest) error {
	stepExecution, err := h.modelClient.WorkflowStepExecutions().Query().
		Where(workflowstepexecution.ID(req.ID)).
		Only(req.Context)
	if err != nil {
		return err
	}

	workflowExecution, err := h.modelClient.WorkflowExecutions().Query().
		Where(workflowexecution.ID(stepExecution.WorkflowExecutionID)).
		Only(req.Context)
	if err != nil {
		return err
	}

	return h.modelClient.WithTx(req.Context, func(tx *model.Tx) error {
		client, err := workflow.NewArgoWorkflowClient(h.modelClient, h.k8sConfig)
		if err != nil {
			return err
		}

		return client.Resume(req.Context, workflow.ResumeOptions{
			Approve:               req.Approve,
			WorkflowExecution:     workflowExecution,
			WorkflowStepExecution: stepExecution,
		})
	})
}
