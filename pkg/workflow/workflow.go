package workflow

import (
	"context"

	"k8s.io/client-go/rest"

	"github.com/seal-io/walrus/pkg/auths/session"
	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/utils/log"
)

// Run runs the workflow execution to the argo workflow server.
func Run(
	ctx context.Context,
	mc model.ClientSet,
	wf *model.Workflow,
	opts dao.ExecuteOptions,
) (*model.WorkflowExecution, error) {
	logger := log.WithName("workflow")

	client, err := NewArgoWorkflowClient(mc, opts.RestCfg)
	if err != nil {
		return nil, err
	}

	s := session.MustGetSubject(ctx)

	wfe, err := dao.CreateWorkflowExecution(ctx, mc, dao.CreateWorkflowExecutionOptions{
		ExecuteOptions: opts,
		Workflow:       wf,
	})
	if err != nil {
		return nil, err
	}

	err = client.Submit(ctx, SubmitOptions{
		WorkflowExecution: wfe,
		SubjectID:         s.ID,
	})

	if err != nil {
		{
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

	return wfe, err
}

// Rerun reruns the workflow execution to the argo workflow server.
func Rerun(
	ctx context.Context,
	mc model.ClientSet,
	restCfg *rest.Config,
	wfe *model.WorkflowExecution,
) error {
	logger := log.WithName("workflow")

	client, err := NewArgoWorkflowClient(mc, restCfg)
	if err != nil {
		return err
	}

	err = client.Resubmit(ctx, ResubmitOptions{
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
