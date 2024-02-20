package k8sctrls

import (
	"context"

	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/seal-io/walrus/pkg/dao/model"
	runjob "github.com/seal-io/walrus/pkg/resourceruns/job"
	pkgworkflow "github.com/seal-io/walrus/pkg/workflow"
)

func init() {
	// Register none native kubernetes resource scheme below.
	_ = wfv1.AddToScheme(scheme.Scheme)
}

// SetupOptions holds the options for creating the controller.
type SetupOptions struct {
	ReconcileHelper

	ModelClient *model.Client
}

type ReconcileHelper interface {
	// GetLogger returns the logr.Logger.
	GetLogger() logr.Logger
	// GetConfig returns the rest.Config.
	GetConfig() *rest.Config
	// GetScheme returns the runtime.Schema.
	GetScheme() *runtime.Scheme
	// GetClient returns a client.Client configured with the rest.Config.
	// This client may not be a fully "direct" client -- it may read from a cache, for
	// instance.
	GetClient() client.Client
	// GetFieldIndexer returns a client.FieldIndexer configured with the client.
	GetFieldIndexer() client.FieldIndexer
	// GetCache returns the cache.Cache.
	GetCache() cache.Cache
	// GetEventRecorderFor returns a new record.EventRecorder for the provided name.
	GetEventRecorderFor(name string) record.EventRecorder
	// GetRESTMapper returns a meta.RESTMapper.
	GetRESTMapper() meta.RESTMapper
	// GetAPIReader returns a client.Reader that will be configured to use the API server.
	// This should be used sparingly and only when the client does not fit your
	// use case.
	GetAPIReader() client.Reader
}

type Reconciler interface {
	Setup(mgr ctrl.Manager) error
}

func (m *Manager) Setup(ctx context.Context, opts SetupOptions) ([]Reconciler, error) {
	workflowClient, err := pkgworkflow.NewArgoWorkflowClient(opts.ModelClient, opts.GetConfig())
	if err != nil {
		return nil, err
	}

	// Setup reconciler below.
	return []Reconciler{
		runjob.Reconciler{
			Logger:      opts.GetLogger().WithName("deployer").WithName("tf"),
			Kubeconfig:  opts.GetConfig(),
			KubeClient:  opts.GetClient(),
			ModelClient: opts.ModelClient,
		},
		pkgworkflow.WorkflowReconciler{
			Logger:     opts.GetLogger().WithName("workflow"),
			KubeClient: opts.GetClient(),
			StatusSyncer: pkgworkflow.NewStatusSyncer(
				opts.ModelClient,
				workflowClient),
		},
	}, nil
}
