package configmap

import (
	"context"

	"github.com/go-logr/logr"
	core "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/seal-io/seal/pkg/dao/model"
)

type Reconciler struct {
	Logger      logr.Logger
	KubeClient  client.Client
	ModelClient *model.Client
}

func (r Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	return ctrl.Result{}, nil
}

func (r Reconciler) Setup(mgr ctrl.Manager) error {
	return ctrl.
		NewControllerManagedBy(mgr).
		For(&core.ConfigMap{}).
		Complete(r)
}
