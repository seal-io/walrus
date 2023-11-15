package workflow

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflow"
	"github.com/argoproj/argo-workflows/v3/workflow/common"
	corev1 "k8s.io/api/core/v1"

	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/utils/strs"
)

// GetWorkflowStepExecutionLogs gets workflow step execution logs.
func (s *ArgoWorkflowClient) GetLogs(ctx context.Context, opts LogsOptions) ([]byte, error) {
	logsClient, err := s.apiClient.NewWorkflowServiceClient().
		WorkflowLogs(s.apiClient.Ctx, &workflow.WorkflowLogRequest{
			Name:      strs.Join("-", opts.WorkflowExecution.Name, opts.WorkflowExecution.ID.String()),
			Namespace: types.WalrusSystemNamespace,
			LogOptions: &corev1.PodLogOptions{
				Container:    common.MainContainerName,
				Previous:     opts.Previous,
				SinceSeconds: opts.SinceSeconds,
				Timestamps:   opts.Timestamps,
				TailLines:    opts.TailLines,
			},
			Selector: fmt.Sprintf("step-execution-id=%s", opts.StepExecution.ID),
		})
	if err != nil {
		return nil, err
	}

	logs := []byte{}

	for {
		event, err := logsClient.Recv()
		if errors.Is(err, io.EOF) {
			return logs, nil
		}

		if err != nil {
			return nil, err
		}

		logs = append(logs, []byte(event.Content+"\n")...)
	}
}

// StreamWorkflowStepExecutionLogs streams workflow step execution logs.
func (s *ArgoWorkflowClient) StreamLogs(ctx context.Context, opts StreamLogsOptions) error {
	logsClient, err := s.apiClient.NewWorkflowServiceClient().
		WorkflowLogs(s.apiClient.Ctx, &workflow.WorkflowLogRequest{
			Name:      strs.Join("-", opts.WorkflowExecution.Name, opts.WorkflowExecution.ID.String()),
			Namespace: types.WalrusSystemNamespace,
			LogOptions: &corev1.PodLogOptions{
				Container:    common.MainContainerName,
				Follow:       true,
				Previous:     opts.Previous,
				SinceSeconds: opts.SinceSeconds,
				Timestamps:   opts.Timestamps,
				TailLines:    opts.TailLines,
			},
			Selector: fmt.Sprintf("step-execution-id=%s", opts.StepExecution.ID),
		})
	if err != nil {
		return err
	}

	for {
		event, err := logsClient.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}

		if err != nil {
			return err
		}

		_, err = opts.Out.Write([]byte(event.Content + "\n"))
		if err != nil {
			return err
		}
	}
}
