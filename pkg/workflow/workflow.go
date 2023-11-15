package workflow

import (
	"context"

	"github.com/seal-io/walrus/pkg/auths/session"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/utils/log"
)

// Run runs the workflow execution to the argo workflow server.
func Run(
	ctx context.Context,
	mc model.ClientSet,
	client Client,
	wfe *model.WorkflowExecution,
) error {
	s := session.MustGetSubject(ctx)

	return client.Submit(ctx, SubmitOptions{
		WorkflowExecution: wfe,
		SubjectID:         s.ID,
	})
}

// Rerun reruns the workflow execution to the argo workflow server.
func Rerun(
	ctx context.Context,
	mc model.ClientSet,
	client Client,
	wfe *model.WorkflowExecution,
) error {
	logger := log.WithName("workflow")

	err := client.Resubmit(ctx, ResubmitOptions{
		WorkflowExecution: wfe,
	})
	if err != nil {
		{
			// If the workflow execution is not found, reset the status to pending.
			ctx := context.Background()

			status.WorkflowExecutionStatusPending.False(wfe, err.Error())

			err := mc.WorkflowExecutions().UpdateOne(wfe).
				SetStatus(wfe.Status).
				Exec(ctx)
			if err != nil {
				logger.Errorf("failed to update workflow execution status: %v", err)
			}
		}
	}

	return err
}
