package workflow

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/argoproj/argo-workflows/v3/pkg/apiclient"
	"github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflow"
	"github.com/argoproj/argo-workflows/v3/workflow/common"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/workflowexecution"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/utils/strs"
)

// LogOptions contains options for workflow logs.
type LogOptions struct {
	Workflow   string
	PodName    string
	Grep       string
	Selector   string
	ApiClient  apiclient.Client
	LogOptions *corev1.PodLogOptions
	Out        io.Writer
}

// StreamWorkflowLogs streams workflow logs.
// With selector step-execution-id=stepExecutionID it can filter logs by step name.
func StreamWorkflowLogs(
	ctx context.Context,
	opts LogOptions,
) error {
	serviceClient := opts.ApiClient.NewWorkflowServiceClient()

	stream, err := serviceClient.WorkflowLogs(ctx, &workflow.WorkflowLogRequest{
		Name:       opts.Workflow,
		Namespace:  types.WalrusSystemNamespace,
		PodName:    opts.PodName,
		LogOptions: opts.LogOptions,
		Selector:   opts.Selector,
		Grep:       opts.Grep,
	})
	if err != nil {
		return err
	}

	for {
		event, err := stream.Recv()
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

type StepExecutionLogOptions struct {
	ModelClient   model.ClientSet
	RestCfg       *rest.Config
	StepExecution *model.WorkflowStepExecution
}

// GetWorkflowStepExecutionLogs gets workflow step execution logs.
func GetWorkflowStepExecutionLogs(ctx context.Context, opts StepExecutionLogOptions) ([]byte, error) {
	logsClient, err := getLogsClient(ctx, opts)
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

type StreamWorkflowStepExecutionLogsOptions struct {
	StepExecutionLogOptions

	Out io.Writer
}

func getLogsClient(
	ctx context.Context,
	opts StepExecutionLogOptions,
) (workflow.WorkflowService_WorkflowLogsClient, error) {
	apiClient, err := NewArgoAPIClient(opts.RestCfg)
	if err != nil {
		return nil, err
	}

	workflowExecution, err := opts.ModelClient.WorkflowExecutions().Query().
		Where(workflowexecution.ID(opts.StepExecution.WorkflowExecutionID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return apiClient.NewWorkflowServiceClient().WorkflowLogs(apiClient.Ctx, &workflow.WorkflowLogRequest{
		Name:      strs.Join("-", workflowExecution.Name, workflowExecution.ID.String()),
		Namespace: types.WalrusSystemNamespace,
		LogOptions: &corev1.PodLogOptions{
			Container: common.MainContainerName,
			Follow:    true,
		},
		Selector: fmt.Sprintf("step-execution-id=%s", opts.StepExecution.ID),
	})
}

// StreamWorkflowStepExecutionLogs streams workflow step execution logs.
func StreamWorkflowStepExecutionLogs(ctx context.Context, opts StreamWorkflowStepExecutionLogsOptions) error {
	logsClient, err := getLogsClient(ctx, opts.StepExecutionLogOptions)
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

// ArchiveWorkflowStepExecutionLogs archives workflow step execution logs.
func ArchiveWorkflowStepExecutionLogs(ctx context.Context, opts StepExecutionLogOptions) error {
	logs, err := GetWorkflowStepExecutionLogs(ctx, StepExecutionLogOptions{
		RestCfg:       opts.RestCfg,
		ModelClient:   opts.ModelClient,
		StepExecution: opts.StepExecution,
	})
	if err != nil {
		return err
	}

	return opts.ModelClient.WorkflowStepExecutions().UpdateOne(opts.StepExecution).
		SetRecord(string(logs)).
		Exec(ctx)
}
