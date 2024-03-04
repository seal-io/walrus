package walruscore

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
	ctrlreconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"

	walruscore "github.com/seal-io/walrus/pkg/apis/walruscore/v1"
	"github.com/seal-io/walrus/pkg/controller"
)

// ResourceReconciler reconciles a v1.Resource object.
type ResourceReconciler struct{}

var _ ctrlreconcile.Reconciler = (*ResourceReconciler)(nil)

func (r *ResourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = ctrllog.FromContext(ctx)

	// TODO: your logic here

	return ctrl.Result{}, nil
}

func (r *ResourceReconciler) SetupController(_ context.Context, opts controller.SetupOptions) error {
	return ctrl.NewControllerManagedBy(opts.Manager).
		For(&walruscore.Resource{}).
		Complete(r)
}
