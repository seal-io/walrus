package workflow

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/argoproj/argo-workflows/v3/pkg/apiclient"
	"github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflow"
	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow"
	"google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/seal-io/walrus/pkg/auths"
	"github.com/seal-io/walrus/pkg/auths/session"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/k8s"
	optypes "github.com/seal-io/walrus/pkg/operator/types"
	"github.com/seal-io/walrus/pkg/workflow/step"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/pointer"
	"github.com/seal-io/walrus/utils/strs"
)

// Client is the interface that defines the operations of workflow engine.
type Client interface {
	// Submit submits a workflow to the workflow engine.
	Submit(context.Context, SubmitOptions) error
	// Resume resumes a workflow step execution of a workflow execution..
	Resume(context.Context, ResumeOptions) error
	// Resubmit resubmits a workflow to the workflow engine.
	Resubmit(context.Context, ResubmitOptions) error
	// Delete deletes a workflow from the workflow engine.
	Delete(context.Context, DeleteOptions) error
	// GetLogs gets logs of a workflow step execution.
	GetLogs(context.Context, LogsOptions) ([]byte, error)
	// StreamLogs streams logs of a workflow execution.
	StreamLogs(context.Context, StreamLogsOptions) error
	// Terminate terminates a workflow execution.
	Terminate(context.Context, TerminateOptions) error
}

type (
	SubmitOptions struct {
		WorkflowExecution *model.WorkflowExecution
		SubjectID         object.ID
	}
	GetOptions struct {
		WorkflowExecution *model.WorkflowExecution
	}

	DeleteOptions struct {
		WorkflowExecution *model.WorkflowExecution
	}

	// SubmitOptions is the options for submitting a workflow.
	// WorkflowExecution's Edge WorkflowStageExecutions and their Edge WorkflowStepExecutions must be set.
	ResumeOptions struct {
		// Approve or deny of the workflow approval step execution.
		Approve bool

		// WorkflowExecution is the workflow execution to be resumed.
		WorkflowExecution *model.WorkflowExecution
		// WorkflowStepExecution is the workflow step execution to be resumed.
		WorkflowStepExecution *model.WorkflowStepExecution
	}

	ResubmitOptions struct {
		WorkflowExecution *model.WorkflowExecution
	}

	// SubmitParamsOpts is the options for submitting a workflow with parameters.
	SubmitParamsOpts struct {
		WorkflowExecution *model.WorkflowExecution
		Params            map[string]string
	}

	// TerminateOptions is the options for terminating a workflow execution.
	TerminateOptions struct {
		WorkflowExecution *model.WorkflowExecution
	}

	// LogsOptions is the options for getting logs of a workflow step execution.
	LogsOptions struct {
		optypes.LogOptions

		WorkflowExecution *model.WorkflowExecution
		StepExecution     *model.WorkflowStepExecution
	}

	// StreamLogsOptions is the options for streaming logs of a workflow execution.
	StreamLogsOptions struct {
		LogsOptions

		Out io.Writer
	}
)

type ArgoWorkflowClient struct {
	Logger log.Logger
	mc     model.ClientSet
	kc     *rest.Config
	tm     *TemplateManager
	// Argo workflow clientset.
	apiClient *ArgoAPIClient
}

func NewArgoWorkflowClient(mc model.ClientSet, restCfg *rest.Config) Client {
	return &ArgoWorkflowClient{
		Logger:    log.WithName("workflow-service"),
		mc:        mc,
		kc:        restCfg,
		tm:        NewTemplateManager(mc),
		apiClient: NewArgoAPIClient(restCfg),
	}
}

func (s *ArgoWorkflowClient) Submit(ctx context.Context, opts SubmitOptions) error {
	token, err := s.createToken(ctx, opts.WorkflowExecution)
	if err != nil {
		return err
	}

	wf, err := s.tm.ToArgoWorkflow(ctx, opts.WorkflowExecution, token)
	if err != nil {
		return err
	}

	secret, err := s.setK8sSecret(ctx, k8sSecretOptions{
		Action: "create",
		Secret: &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("workflow-execution-%s", opts.WorkflowExecution.ID.String()),
				Namespace: types.WalrusSystemNamespace,
			},
			Data: map[string][]byte{
				"token": []byte(token),
			},
		},
	})
	if err != nil {
		return err
	}

	awf, err := s.apiClient.NewWorkflowServiceClient().CreateWorkflow(s.apiClient.Ctx, &workflow.WorkflowCreateRequest{
		Namespace: types.WalrusSystemNamespace,
		Workflow:  wf,
	})
	if err != nil {
		return err
	}

	// Set ownerReferences to secret with the workflow.
	secret.ObjectMeta.OwnerReferences = []metav1.OwnerReference{
		{
			APIVersion: wfv1.APIVersion,
			Kind:       wfv1.WorkflowKind,
			Name:       awf.Name,
			UID:        awf.UID,
		},
	}

	_, err = s.setK8sSecret(ctx, k8sSecretOptions{
		Action: "update",
		Secret: secret,
	})

	return err
}

func (s *ArgoWorkflowClient) Resume(ctx context.Context, opts ResumeOptions) error {
	subject := session.MustGetSubject(ctx)

	if opts.WorkflowStepExecution.Type != types.WorkflowStepTypeApproval {
		return fmt.Errorf("cannot resume workflow execution type %s", opts.WorkflowStepExecution.Type)
	}

	approvalSpec, err := types.NewWorkflowStepApprovalSpec(opts.WorkflowStepExecution.Attributes)
	if err != nil {
		return err
	}

	err = approvalSpec.SetUserApproval(subject.ID, opts.Approve)
	if err != nil {
		return err
	}

	// Update workflow step execution approval attributes.
	err = s.mc.WorkflowStepExecutions().UpdateOne(opts.WorkflowStepExecution).
		SetAttributes(approvalSpec.ToAttributes()).
		Exec(ctx)
	if err != nil {
		return err
	}

	if approvalSpec.IsRejected() {
		// Stop workflow step execution if rejected.
		_, err = s.apiClient.NewWorkflowServiceClient().StopWorkflow(s.apiClient.Ctx, &workflow.WorkflowStopRequest{
			Name:              getWorkflowName(opts.WorkflowExecution),
			Namespace:         types.WalrusSystemNamespace,
			NodeFieldSelector: fmt.Sprintf("templateName=%s", step.StepTemplateName(opts.WorkflowStepExecution)),
		})
	} else {
		// If not approved, do nothing.
		if !approvalSpec.IsApproved() {
			return nil
		}

		// Update secret token.
		err = s.updateWorkflowExecutionToken(ctx, opts.WorkflowExecution)
		if err != nil {
			return err
		}

		// Resume workflow step execution.
		_, err = s.apiClient.NewWorkflowServiceClient().ResumeWorkflow(
			s.apiClient.Ctx,
			&workflow.WorkflowResumeRequest{
				Name:      getWorkflowName(opts.WorkflowExecution),
				Namespace: types.WalrusSystemNamespace,
				NodeFieldSelector: fmt.Sprintf(
					"templateName=%s",
					step.StepTemplateName(opts.WorkflowStepExecution),
				),
			})
	}

	return err
}

func (s *ArgoWorkflowClient) Resubmit(ctx context.Context, opts ResubmitOptions) error {
	if err := s.Delete(ctx, DeleteOptions(opts)); err != nil {
		return err
	}

	subject := session.MustGetSubject(ctx)

	if err := s.Submit(ctx, SubmitOptions{
		WorkflowExecution: opts.WorkflowExecution,
		SubjectID:         subject.ID,
	}); err != nil {
		return err
	}

	return s.mc.WithTx(ctx, func(tx *model.Tx) error {
		return ResetWorkflowExecutionStatus(ctx, tx, opts.WorkflowExecution)
	})
}

func isNotFoundErr(err error) bool {
	if st, ok := grpcstatus.FromError(err); ok {
		return st.Code() == codes.NotFound
	}

	return false
}

func (s *ArgoWorkflowClient) Delete(ctx context.Context, opts DeleteOptions) error {
	_, err := s.apiClient.NewWorkflowServiceClient().DeleteWorkflow(s.apiClient.Ctx, &workflow.WorkflowDeleteRequest{
		Name:      getWorkflowName(opts.WorkflowExecution),
		Namespace: types.WalrusSystemNamespace,
	})
	if err != nil && !isNotFoundErr(err) {
		return err
	}

	_, err = s.setK8sSecret(ctx, k8sSecretOptions{
		Action: "delete",
		Secret: &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("workflow-execution-%s", opts.WorkflowExecution.ID.String()),
				Namespace: types.WalrusSystemNamespace,
			},
		},
	})

	return err
}

// Terminate terminates a workflow execution.
// It will stop all nodes of the workflow execution.
func (s *ArgoWorkflowClient) Terminate(ctx context.Context, opts TerminateOptions) error {
	_, err := s.apiClient.NewWorkflowServiceClient().TerminateWorkflow(
		s.apiClient.Ctx,
		&workflow.WorkflowTerminateRequest{
			Name:      getWorkflowName(opts.WorkflowExecution),
			Namespace: types.WalrusSystemNamespace,
		})
	if err != nil && !kerrors.IsNotFound(err) {
		return err
	}

	return nil
}

// createToken creates a token for a workflow execution.
func (s *ArgoWorkflowClient) createToken(
	ctx context.Context,
	workflowExecution *model.WorkflowExecution,
) (string, error) {
	subject := session.MustGetSubject(ctx)

	const _1Day = 60 * 60 * 24

	at, err := auths.CreateAccessToken(ctx,
		s.mc,
		subject.ID,
		types.TokenKindDeployment,
		fmt.Sprintf("%s-%d", workflowExecution.ID.String(), time.Now().Unix()),
		pointer.Int(_1Day),
	)
	if err != nil {
		return "", err
	}

	return at.AccessToken, nil
}

// updateWorkflowExecutionToken updates the token of a workflow execution.
func (s *ArgoWorkflowClient) updateWorkflowExecutionToken(
	ctx context.Context,
	workflowExecution *model.WorkflowExecution,
) error {
	token, err := s.createToken(ctx, workflowExecution)
	if err != nil {
		return err
	}

	clientSet, err := kubernetes.NewForConfig(s.kc)
	if err != nil {
		return err
	}

	secret, err := clientSet.CoreV1().Secrets(types.WalrusSystemNamespace).
		Get(ctx, fmt.Sprintf("workflow-execution-%s", workflowExecution.ID.String()), metav1.GetOptions{})
	if err != nil {
		return err
	}

	secret.Data["token"] = []byte(token)

	_, err = s.setK8sSecret(ctx, k8sSecretOptions{
		Action: "update",
		Secret: secret,
	})

	return err
}

type k8sSecretOptions struct {
	Action string
	Secret *corev1.Secret
}

// setK8sSecret sets a k8s secret.
func (s *ArgoWorkflowClient) setK8sSecret(
	ctx context.Context,
	opts k8sSecretOptions,
) (secret *corev1.Secret, err error) {
	clientSet, err := kubernetes.NewForConfig(s.kc)
	if err != nil {
		return
	}

	switch opts.Action {
	case "create":
		secret, err = clientSet.CoreV1().Secrets(types.WalrusSystemNamespace).
			Create(ctx, opts.Secret, metav1.CreateOptions{})
		if err != nil && !kerrors.IsAlreadyExists(err) {
			return
		}
	case "update":
		secret, err = clientSet.CoreV1().Secrets(types.WalrusSystemNamespace).
			Update(ctx, opts.Secret, metav1.UpdateOptions{})
		if err != nil {
			return
		}
	case "delete":
		err = clientSet.CoreV1().Secrets(types.WalrusSystemNamespace).
			Delete(ctx, opts.Secret.Name, metav1.DeleteOptions{})
		if err != nil && kerrors.IsNotFound(err) {
			return nil, nil
		}
	default:
		err = fmt.Errorf("invalid action: %s", opts.Action)
	}

	return
}

// ArgoAPIClient is a wrapper of argo workflow client.
// It interacts with argo workflow server.
type ArgoAPIClient struct {
	apiclient.Client

	Ctx context.Context
}

func NewArgoAPIClient(restCfg *rest.Config) *ArgoAPIClient {
	apiConfig := k8s.ToClientCmdApiConfig(restCfg)
	clientConfig := clientcmd.NewDefaultClientConfig(apiConfig, nil)

	ctx, apiClient, _ := apiclient.NewClientFromOpts(apiclient.Opts{
		ClientConfigSupplier: func() clientcmd.ClientConfig {
			return clientConfig
		},
	})

	return &ArgoAPIClient{
		Client: apiClient,
		Ctx:    ctx,
	}
}

// getWorkflowName returns the target workflow name of a workflow execution.
func getWorkflowName(workflowExecution *model.WorkflowExecution) string {
	return strs.Join("-", workflowExecution.Name, workflowExecution.ID.String())
}
