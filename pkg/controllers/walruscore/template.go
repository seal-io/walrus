package walruscore

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
	ctrlreconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"

	walruscore "github.com/seal-io/walrus/pkg/apis/walruscore/v1"
	"github.com/seal-io/walrus/pkg/controller"
)

// TemplateReconciler reconciles a v1.Template object.
type TemplateReconciler struct{}

var _ ctrlreconcile.Reconciler = (*TemplateReconciler)(nil)

func (r *TemplateReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = ctrllog.FromContext(ctx)

	// TODO: your logic here

	return ctrl.Result{}, nil
}

func (r *TemplateReconciler) SetupController(_ context.Context, opts controller.SetupOptions) error {
	return ctrl.NewControllerManagedBy(opts.Manager).
		For(&walruscore.Template{}).
		Complete(r)
}
