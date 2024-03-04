package webhook

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlwebhook "sigs.k8s.io/controller-runtime/pkg/webhook"
	ctrladmission "sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type (
	// SetupOptions is the options for setting up a webhook.
	SetupOptions struct {
		// Manager is the controller-runtime manager.
		Manager ctrl.Manager
	}

	// Setup is the interface for the webhook setup.
	Setup interface {
		// SetupWebhook returns the webhook affected object if successful.
		//
		// SetupWebhook is called before the Cache is started,
		// you should not do anything that requires the Cache to be started.
		// Instead, you can configure the Cache, like AddIndex or something else.
		SetupWebhook(context.Context, SetupOptions) (runtime.Object, error)
	}
)

// DefaultCustomValidator implements webhook.CustomValidator,
// which is used to combine and override the required methods.
type DefaultCustomValidator struct{}

var _ ctrlwebhook.CustomValidator = (*DefaultCustomValidator)(nil)

func (DefaultCustomValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (ctrladmission.Warnings, error) {
	return nil, nil
}

func (DefaultCustomValidator) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (ctrladmission.Warnings, error) {
	return nil, nil
}

func (DefaultCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (ctrladmission.Warnings, error) {
	return nil, nil
}
