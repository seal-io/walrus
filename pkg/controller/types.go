package controller

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
	ctrlreconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type (
	// SetupOptions is the options for setting up a controller.
	SetupOptions struct {
		// Manager is the controller-runtime manager.
		Manager ctrl.Manager
	}

	// Setup is the interface for the controller setup.
	Setup interface {
		ctrlreconcile.Reconciler
		// SetupController sets up the controller.
		//
		// SetupController is called before the Cache is started,
		// you should not do anything that requires the Cache to be started.
		// Instead, you can configure the Cache, like IndexField or something else.
		SetupController(ctx context.Context, opts SetupOptions) error
	}
)
