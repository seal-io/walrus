package manager

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlcache "sigs.k8s.io/controller-runtime/pkg/cache"
	ctrlcli "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	ctrlmetricsrv "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	"github.com/seal-io/walrus/pkg/clients/clientset"
	"github.com/seal-io/walrus/pkg/clients/clientset/scheme"
	"github.com/seal-io/walrus/pkg/manager/webhookserver"
	"github.com/seal-io/walrus/pkg/system"
	"github.com/seal-io/walrus/pkg/systemkuberes"
)

type Config struct {
	InformerCacheResyncPeriod time.Duration
	KubeConfigPath            string
	KubeClientConfig          rest.Config
	KubeHTTPClient            *http.Client
	KubeClient                clientset.Interface
	KubeLeaderElection        bool
	KubeLeaderLease           time.Duration
	KubeLeaderRenewTimeout    time.Duration
	ServeListenerCertDir      string
	ServeListener             net.Listener
}

func (c *Config) Apply(_ context.Context) (*Manager, error) {
	ctrlMgrOpts := ctrl.Options{
		// General.
		GracefulShutdownTimeout: ptr.To(30 * time.Second),
		Scheme:                  scheme.Scheme,
		Logger:                  klog.Background().WithName("ctrl"),

		// Client.
		Client: ctrlcli.Options{
			HTTPClient: c.KubeHTTPClient,
		},
		NewClient: func(config *rest.Config, options ctrlcli.Options) (ctrlcli.Client, error) {
			return ctrlcli.NewWithWatch(config, options)
		},

		// Cache.
		Cache: ctrlcache.Options{
			HTTPClient: c.KubeHTTPClient,
			SyncPeriod: ptr.To(c.InformerCacheResyncPeriod),
		},

		// Leader election.
		LeaderElectionReleaseOnCancel: true,
		LeaderElectionNamespace:       systemkuberes.SystemNamespaceName,
		LeaderElectionID:              "walrus-leader",
		LeaderElection:                c.KubeLeaderElection,
		LeaseDuration:                 ptr.To(c.KubeLeaderLease),
		RenewDeadline:                 ptr.To(c.KubeLeaderRenewTimeout),
		RetryPeriod:                   ptr.To(2 * time.Second),

		// Disable default webhook server.
		WebhookServer: webhookserver.Dummy(),
		// Disable default metrics service.
		Metrics: ctrlmetricsrv.Options{BindAddress: "0"},
		// Disable default healthcheck service.
		HealthProbeBindAddress: "0",
		// Disable default profiling service.
		PprofBindAddress: "0",
	}

	// Enable webhook serving,
	// includes configurations installation.
	if c.ServeListener != nil {
		ctrlMgrOpts.WebhookServer = webhookserver.Enhance(c.ServeListener, c.ServeListenerCertDir, c.KubeClient)
	}

	ctrlManager, err := ctrl.NewManager(rest.CopyConfig(&c.KubeClientConfig), ctrlMgrOpts)
	if err != nil {
		return nil, fmt.Errorf("create controller manager: %w", err)
	}

	system.ConfigureLoopbackCtrlRuntime(ctrlManager.GetClient(), ctrlManager.GetCache())

	return &Manager{
		CtrlManager: &_CtrlManager{
			Manager:       ctrlManager,
			IndexedFields: sets.Set[string]{},
		},
	}, nil
}

type (
	// _CtrlManager is a wrapper around ctrl.Manager.
	_CtrlManager struct {
		ctrl.Manager
		IndexedFields sets.Set[string]
	}

	// _CtrlClientFieldIndexer is a wrapper around ctrl.FieldIndexer.
	_CtrlClientFieldIndexer struct {
		ctrl.Manager
		IndexedFields sets.Set[string]
	}
)

func (m *_CtrlManager) GetFieldIndexer() ctrlcli.FieldIndexer {
	return &_CtrlClientFieldIndexer{
		Manager:       m.Manager,
		IndexedFields: m.IndexedFields,
	}
}

func (i *_CtrlClientFieldIndexer) IndexField(ctx context.Context, obj ctrlcli.Object, field string, extractValue ctrlcli.IndexerFunc) error {
	gvk, err := apiutil.GVKForObject(obj, i.Manager.GetScheme())
	if err != nil {
		return err
	}
	key := gvk.String() + "/" + field
	if i.IndexedFields.Has(key) {
		// If the field is already indexed, skip.
		return nil
	}
	i.IndexedFields.Insert(key)
	return i.Manager.GetFieldIndexer().IndexField(ctx, obj, field, extractValue)
}
