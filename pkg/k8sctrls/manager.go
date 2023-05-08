package k8sctrls

import (
	"context"
	"fmt"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/utils/log"
)

func NewManager(cfg *rest.Config) (*Manager, error) {
	var logger = log.WithName("k8sctrl")
	var mgr, err = ctrl.NewManager(cfg, ctrl.Options{
		Scheme:    scheme.Scheme,
		Logger:    log.AsLogr(logger),
		Namespace: types.SealSystemNamespace,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating kubernetes controller manager: %w", err)
	}
	return &Manager{
		logger: logger,
		mgr:    mgr,
	}, nil
}

type Manager struct {
	logger log.Logger
	mgr    ctrl.Manager
}

type StartOptions struct {
	SetupOptions
}

func (m *Manager) Start(ctx context.Context, opts StartOptions) error {
	m.logger.Info("starting")

	var mgr = m.mgr
	opts.SetupOptions.ReconcileHelper = mgr
	var reconcilers, err = m.Setup(ctx, opts.SetupOptions)
	if err != nil {
		return err
	}
	for i := 0; i < len(reconcilers); i++ {
		if err = reconcilers[i].Setup(mgr); err != nil {
			return fmt.Errorf("error setting up kubernetes controller: %w", err)
		}
	}
	return mgr.Start(ctx)
}
