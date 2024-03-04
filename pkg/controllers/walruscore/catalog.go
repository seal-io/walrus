package walruscore

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
	ctrlreconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"

	walruscore "github.com/seal-io/walrus/pkg/apis/walruscore/v1"
	"github.com/seal-io/walrus/pkg/controller"
)

// CatalogReconciler reconciles a v1.Catalog object.
type CatalogReconciler struct{}

var _ ctrlreconcile.Reconciler = (*CatalogReconciler)(nil)

func (r *CatalogReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = ctrllog.FromContext(ctx)

	// TODO: your logic here

	return ctrl.Result{}, nil
}

func (r *CatalogReconciler) SetupController(_ context.Context, opts controller.SetupOptions) error {
	return ctrl.NewControllerManagedBy(opts.Manager).
		For(&walruscore.Catalog{}).
		Complete(r)
}
