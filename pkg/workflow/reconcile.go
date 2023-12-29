package workflow

import (
	"context"

	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/go-logr/logr"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/seal-io/walrus/utils/log"
)

const (
	workflowExecutionIDLabel     = "walrus.seal.io/workflow-execution-id"
	workflowStepExecutionIDLabel = "walrus.seal.io/workflow-step-execution-id"
)

// WorkflowReconciler reconciles a Workflow object.
type WorkflowReconciler struct {
	Logger       logr.Logger
	KubeClient   client.Client
	StatusSyncer *StatusSyncer
}

// Reconcile reconciles the workflow.
func (r WorkflowReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.WithName("workflow").WithName("reconcile")
	wf := &wfv1.Workflow{}

	err := r.KubeClient.Get(ctx, req.NamespacedName, wf)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, err
	}

	logger.Debugf("workflow: %s, phase: %s", wf.Name, wf.Status.Phase)

	if !wf.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, nil
	}

	err = r.StatusSyncer.SyncWorkflowExecutionStatus(ctx, wf)
	if err != nil {
		return ctrl.Result{}, err
	}

	workflowExecutionCanceled, err := r.StatusSyncer.IsCanceled(ctx, wf)
	if err != nil {
		return ctrl.Result{}, err
	}

	for i := range wf.Status.Nodes {
		node := wf.Status.Nodes[i]

		templateType, stage, id, ok := parseTemplateName(node.DisplayName)
		if !ok {
			continue
		}

		logger.Debugf("node: %s, status: %s", node.DisplayName, node.Phase)

		switch {
		case templateType == templateTypeStage && stage == templateStageEnter:
			err = r.StatusSyncer.SyncStageExecutionStatus(ctx, node, id, workflowExecutionCanceled)

		case templateType == templateTypeStep && stage == templateStageMain:
			err = r.StatusSyncer.SyncStepExecutionStatus(ctx, node, id, workflowExecutionCanceled)
		}

		if err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r WorkflowReconciler) Setup(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&wfv1.Workflow{}).
		Complete(r)
}
