package workflow

import (
	"context"
	"fmt"
	"time"

	"github.com/argoproj/argo-workflows/v3/pkg/apiclient"
	"github.com/argoproj/argo-workflows/v3/pkg/apiclient/workflow"
	apisworkflow "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/seal-io/walrus/pkg/auths"
	"github.com/seal-io/walrus/pkg/auths/session"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/workflowexecution"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/k8s"
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
}

type (
	SubmitOptions struct {
		WorkflowExecution *model.WorkflowExecution
		SubjectID         object.ID
	}
	GetOptions struct {
		Workflow *model.WorkflowExecution
	}

	DeleteOptions struct {
		Workflow *model.WorkflowExecution
	}

	// SubmitOptions is the options for submitting a workflow.
	// WorkflowExecution's Edge WorkflowStageExecutions and their Edge WorkflowStepExecutions must be set.
	ResumeOptions struct {
		// Only used when Type is "workflow".
		WorkflowExecution *model.WorkflowExecution
		// Only used when Type is "step".
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
)

type ArgoWorkflowClient struct {
	Logger log.Logger
	mc     model.ClientSet
	kc     *rest.Config
	tm     *TemplateManager
	// Argo workflow clientset.
	apiClient *ArgoAPIClient
}

func NewArgoWorkflowClient(mc model.ClientSet, restCfg *rest.Config) (Client, error) {
	apiClient, err := NewArgoAPIClient(restCfg)
	if err != nil {
		return nil, err
	}

	return &ArgoWorkflowClient{
		Logger:    log.WithName("workflow-service"),
		mc:        mc,
		kc:        restCfg,
		tm:        NewTemplateManager(mc),
		apiClient: apiClient,
	}, nil
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
			APIVersion: apisworkflow.APIVersion,
			Kind:       apisworkflow.WorkflowKind,
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

	if err = approvalSpec.SetApprovedUser(subject.ID); err != nil {
		return err
	}

	// Update workflow step execution approval attributes.
	err = s.mc.WorkflowStepExecutions().UpdateOne(opts.WorkflowStepExecution).
		SetAttributes(approvalSpec.ToAttributes()).
		Exec(ctx)
	if err != nil {
		return err
	}

	// If not approved, do nothing.
	if !approvalSpec.IsApproved() {
		return nil
	}

	workflowExecution, err := s.mc.WorkflowExecutions().Query().
		Where(workflowexecution.ID(opts.WorkflowStepExecution.WorkflowExecutionID)).
		Only(ctx)
	if err != nil {
		return err
	}

	// Update secret token.
	err = s.updateWorkflowExecutionToken(ctx, workflowExecution)
	if err != nil {
		return err
	}

	// Resume workflow step execution.
	_, err = s.apiClient.NewWorkflowServiceClient().ResumeWorkflow(s.apiClient.Ctx, &workflow.WorkflowResumeRequest{
		Name:              getWorkflowName(workflowExecution),
		Namespace:         types.WalrusSystemNamespace,
		NodeFieldSelector: fmt.Sprintf("templateName=suspend-%s", opts.WorkflowStepExecution.ID.String()),
	})

	return err
}

func (s *ArgoWorkflowClient) Resubmit(ctx context.Context, opts ResubmitOptions) error {
	awf, err := s.apiClient.NewWorkflowServiceClient().GetWorkflow(s.apiClient.Ctx, &workflow.WorkflowGetRequest{
		Name:      getWorkflowName(opts.WorkflowExecution),
		Namespace: types.WalrusSystemNamespace,
	})
	if err != nil && kerrors.IsNotFound(err) {
		return err
	}

	if awf != nil {
		_, err = s.apiClient.NewWorkflowServiceClient().
			ResubmitWorkflow(s.apiClient.Ctx, &workflow.WorkflowResubmitRequest{
				Name:      getWorkflowName(opts.WorkflowExecution),
				Namespace: types.WalrusSystemNamespace,
				Memoized:  false,
			})
		if err != nil {
			return err
		}
	} else {
		subject := session.MustGetSubject(ctx)

		err = s.Submit(ctx, SubmitOptions{
			WorkflowExecution: opts.WorkflowExecution,
			SubjectID:         subject.ID,
		})
		if err != nil {
			return err
		}
	}

	return ResetWorkflowExecutionStatus(ctx, s.mc, opts.WorkflowExecution)
}

func (s *ArgoWorkflowClient) Delete(ctx context.Context, opts DeleteOptions) error {
	_, err := s.apiClient.NewWorkflowServiceClient().DeleteWorkflow(s.apiClient.Ctx, &workflow.WorkflowDeleteRequest{
		Name:      opts.Workflow.Name,
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

func NewArgoAPIClient(restCfg *rest.Config) (*ArgoAPIClient, error) {
	apiConfig := k8s.ToClientCmdApiConfig(restCfg)
	clientConfig := clientcmd.NewDefaultClientConfig(apiConfig, nil)

	ctx, apiClient, err := apiclient.NewClientFromOpts(apiclient.Opts{
		ClientConfigSupplier: func() clientcmd.ClientConfig {
			return clientConfig
		},
	})
	if err != nil {
		return nil, err
	}

	return &ArgoAPIClient{
		Client: apiClient,
		Ctx:    ctx,
	}, nil
}

// getWorkflowName returns the target workflow name of a workflow execution.
func getWorkflowName(workflowExecution *model.WorkflowExecution) string {
	return strs.Join("-", workflowExecution.Name, workflowExecution.ID.String())
}
