package workflowstepexecution

import (
	"context"
	"io"
	"log"

	"github.com/seal-io/walrus/pkg/dao/model/workflowexecution"
	"github.com/seal-io/walrus/pkg/dao/model/workflowstepexecution"
	"github.com/seal-io/walrus/pkg/datalisten/modelchange"
	optypes "github.com/seal-io/walrus/pkg/operator/types"
	pkgworkflow "github.com/seal-io/walrus/pkg/workflow"
	"github.com/seal-io/walrus/utils/topic"
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

	wfe, err := h.modelClient.WorkflowExecutions().Query().
		Where(workflowexecution.ID(wse.WorkflowExecutionID)).
		Only(ctx)
	if err != nil {
		return err
	}

	return h.workflowClient.StreamLogs(ctx, pkgworkflow.StreamLogsOptions{
		LogsOptions: pkgworkflow.LogsOptions{
			LogOptions: optypes.LogOptions{
				Previous:     req.Previous,
				SinceSeconds: req.SinceSeconds,
				TailLines:    req.TailLines,
				Timestamps:   req.Timestamps,
			},
			WorkflowExecution: wfe,
			StepExecution:     wse,
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

	err = h.workflowClient.Resume(req.Context, pkgworkflow.ResumeOptions{
		Approve:               req.Approve,
		WorkflowExecution:     workflowExecution,
		WorkflowStepExecution: stepExecution,
	})
	if err != nil {
		// Reset approval attributes.
		rerr := h.modelClient.WorkflowStepExecutions().UpdateOne(stepExecution).
			SetAttributes(stepExecution.Attributes).
			Exec(req.Context)
		if rerr != nil {
			log.Printf("failed to update workflow step execution status: %v", rerr)
		}

		return err
	}

	return topic.Publish(req.Context, modelchange.WorkflowExecution, modelchange.Event{
		Type: modelchange.EventTypeUpdate,
		Data: []modelchange.EventData{{ID: workflowExecution.ID}},
	})
}
