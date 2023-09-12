package k8sctrls

import (
	"context"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/utils/gopool"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/pointer"
)

// ManagerOptions holds the options for creating a new manager.
type ManagerOptions struct {
	// IsReady observes whether the manager is ready,
	// the caller can leverage this symbol to be aware the manager's progress.
	IsReady *atomic.Bool
	// LeaderElection indicates whether to enable leader election.
	LeaderElection bool
	// LeaderLease indicates the duration of the lease that keeps the leadership.
	LeaderLease time.Duration
	// LeaderRenewTimeout indicates the timeout of renewing the leadership.
	LeaderRenewTimeout time.Duration
}

func NewManager(opts ManagerOptions) (*Manager, error) {
	logger := log.WithName("k8sctrl")

	// Defaults.
	if opts.IsReady == nil {
		opts.IsReady = &atomic.Bool{}
	}

	if opts.LeaderLease == 0 {
		opts.LeaderLease = 15 * time.Second
	}

	if opts.LeaderRenewTimeout == 0 {
		opts.LeaderRenewTimeout = 10 * time.Second
	}

	// Build options for creating controller manager.
	options := ctrl.Options{
		// General.
		Scheme:     scheme.Scheme,
		SyncPeriod: pointer.Duration(time.Hour),
		Logger:     log.AsLogr(logger),
		Namespace:  types.WalrusSystemNamespace,

		// Leader election.
		LeaderElection:                opts.LeaderElection,
		LeaderElectionID:              "walrus-leader-election",
		LeaderElectionNamespace:       types.WalrusSystemNamespace,
		LeaderElectionReleaseOnCancel: true,
		LeaseDuration:                 pointer.Duration(opts.LeaderLease),
		RenewDeadline:                 pointer.Duration(opts.LeaderRenewTimeout),
		RetryPeriod:                   pointer.Duration(2 * time.Second),

		// Disable unexposed services.
		MetricsBindAddress:     "0", // Controller metrics expose by Walrus's server as well.
		HealthProbeBindAddress: "0",
	}

	return &Manager{
		logger:  logger,
		isReady: opts.IsReady,
		options: options,
	}, nil
}

type Manager struct {
	logger  log.Logger
	isReady *atomic.Bool
	options ctrl.Options
}

// StartOptions holds the options for starting the manager.
type StartOptions struct {
	// RestConfig indicates the rest config for connecting Kubernetes.
	RestConfig *rest.Config
	// SetupOptions holds the options for creating the Kubernetes controllers.
	SetupOptions SetupOptions
}

func (m *Manager) Start(ctx context.Context, opts StartOptions) error {
	m.logger.Info("starting")

	if !m.options.LeaderElection {
		return m.doStart(ctx, opts.RestConfig, opts.SetupOptions)
	}

	// After enabled the leader election, the control manager may be forced to quit
	// when it cannot access the Kubernetes cluster for a short time.
	//
	// Quitting is usually not a bad thing, but in order to avoid restarts due to short-time jitter,
	// we add a simple counting retry here for mitigation.
	//
	// The retry mechanism continues trying 10 times, and the retry period is equal to leadership renew timeout.
	// So, a retry window is 10 times of the retry period, which is 50 seconds by default.
	// If the manager has been running for more than one retry window,
	// either restart or not, the retry counter will be reset.
	var (
		retryLimit  = 10
		retryPeriod = *m.options.RetryPeriod
		retryWindow = time.Duration(retryLimit) * retryPeriod

		retries        int
		lastStartEpoch time.Time
	)

	return wait.PollUntilContextCancel(ctx, retryPeriod, true,
		func(ctx context.Context) (done bool, err error) {
			lastStartEpoch = time.Now()

			err = m.doStart(ctx, opts.RestConfig, opts.SetupOptions)
			if err != nil {
				// Reset retries.
				if time.Since(lastStartEpoch) >= retryWindow {
					retries = 0
				}

				// Retry by error message.
				switch errMsg := err.Error(); {
				case strings.Contains(errMsg, "leader election lost"):
					// Restart from leader election lost, restart.
					if retries < retryLimit {
						m.logger.Info("lost leader election, restarting %d times", retryLimit-retries)

						retries++
						err = nil
					}
				case time.Since(lastStartEpoch) <= retryPeriod && strings.Contains(errMsg, "connection refused"):
					// Connection refused in short time, restart.
					if retries < retryLimit {
						m.logger.Infof("cluster connection refused, restarting %d times", retryLimit-retries)

						retries++
						err = nil
					}
				}
			}

			return
		},
	)
}

func (m *Manager) doStart(ctx context.Context, restConfig *rest.Config, setupOpts SetupOptions) error {
	// Notify the manager is not ready.
	m.isReady.Store(false)

	// Create manager.
	mgr, err := ctrl.NewManager(restConfig, m.options)
	if err != nil {
		return fmt.Errorf("error creating kubernetes controller manager: %w", err)
	}

	// Setup manager and get controllers.
	setupOpts.ReconcileHelper = mgr

	controllers, err := m.Setup(ctx, setupOpts)
	if err != nil {
		return err
	}

	// Setup controllers.
	for i := 0; i < len(controllers); i++ {
		if err = controllers[i].Setup(mgr); err != nil {
			return fmt.Errorf("error setting up kubernetes controller: %w", err)
		}
	}

	// Watch for cache sync.
	gopool.Go(func() {
		err := wait.PollUntilContextCancel(ctx, time.Second, true,
			func(ctx context.Context) (done bool, err error) {
				return mgr.GetCache().WaitForCacheSync(ctx), ctx.Err()
			})
		if err != nil {
			return
		}

		// Notify the manager is ready.
		m.isReady.Store(true)
	})

	// Start manager.
	return mgr.Start(ctx)
}
