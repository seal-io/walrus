package k8sctrls

import (
	"context"
	"fmt"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/utils/log"
)

func NewManager(opts ManagerOptions) (*Manager, error) {
	logger := log.WithName("k8sctrl")

	mgr, err := ctrl.NewManager(opts.K8sConfig, ctrl.Options{
		Scheme:    scheme.Scheme,
		Logger:    log.AsLogr(logger),
		Namespace: types.WalrusSystemNamespace,

		// Leader election.
		LeaderElection:          opts.LeaderElection,
		LeaderElectionID:        "walrus-leader-election",
		LeaderElectionNamespace: types.WalrusSystemNamespace,
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

type ManagerOptions struct {
	K8sConfig      *rest.Config
	LeaderElection bool
}

type StartOptions struct {
	SetupOptions
}

func (m *Manager) Start(ctx context.Context, opts StartOptions) error {
	m.logger.Info("starting")

	mgr := m.mgr
	opts.SetupOptions.ReconcileHelper = mgr

	reconcilers, err := m.Setup(ctx, opts.SetupOptions)
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

func (m *Manager) IsReady(ctx context.Context) bool {
	if m.mgr == nil {
		return false
	}

	return m.mgr.GetCache().WaitForCacheSync(ctx)
}
