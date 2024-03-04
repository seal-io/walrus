package manager

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"path/filepath"
	"time"

	certcache "github.com/seal-io/utils/certs/cache"
	"github.com/seal-io/utils/certs/kubecert"
	"github.com/seal-io/utils/osx"
	"github.com/seal-io/utils/pools/gopool"
	"github.com/seal-io/utils/stringx"
	"github.com/seal-io/utils/version"
	"github.com/spf13/pflag"
	"go.uber.org/automaxprocs/maxprocs"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/certwatcher"

	"github.com/seal-io/walrus/pkg/clients/clientset"
	"github.com/seal-io/walrus/pkg/internalprocesses/kubernetes"
	"github.com/seal-io/walrus/pkg/kubeconfig"
	"github.com/seal-io/walrus/pkg/kuberest"
	"github.com/seal-io/walrus/pkg/system"
	"github.com/seal-io/walrus/pkg/systemkuberes"
)

type Options struct {
	// Establish.
	BindAddress net.IP
	BindPort    int
	CertDir     string

	// Control.
	GopoolWorkerFactor        int
	InformerCacheResyncPeriod time.Duration

	// Connect Kubernetes.
	KubeConnTimeout        time.Duration
	KubeConnQPS            float64
	KubeConnBurst          int
	KubeLeaderElection     bool
	KubeLeaderLease        time.Duration
	KubeLeaderRenewTimeout time.Duration

	// Internal.
	Serve bool
}

func NewOptions() *Options {
	return &Options{
		// Establish.
		BindAddress: net.ParseIP("0.0.0.0"),
		BindPort:    443,

		// Control.
		GopoolWorkerFactor:        100,
		InformerCacheResyncPeriod: 1 * time.Hour,

		// Connect Kubernetes.
		KubeConnTimeout:        5 * time.Minute,
		KubeConnQPS:            200,
		KubeConnBurst:          400,
		KubeLeaderElection:     true,
		KubeLeaderLease:        15 * time.Second,
		KubeLeaderRenewTimeout: 10 * time.Second,

		// Internal.
		Serve: true,
	}
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	// Establish.
	fs.IPVar(&o.BindAddress, "bind-address", o.BindAddress,
		"the IP address(without port) on which to serve.")
	fs.IntVar(&o.BindPort, "secure-port", o.BindPort,
		"the port on which to serve HTTPS.")
	fs.StringVar(&o.CertDir, "cert-dir", o.CertDir,
		"the directory where the TLS certs are located. "+
			"if provided, must place tls.crt and tls.key under --cert-dir. ")

	// Control.
	fs.IntVar(&o.GopoolWorkerFactor, "gopool-worker-factor", o.GopoolWorkerFactor,
		"the number of tasks of the goroutine worker pool, "+
			"it is calculated by the number of CPU cores multiplied by this factor.")
	fs.DurationVar(&o.InformerCacheResyncPeriod, "informer-cache-resync-period", o.InformerCacheResyncPeriod,
		"the period at which the informer's cache is resynced.")

	// Connect Kubernetes.
	fs.DurationVar(&o.KubeConnTimeout, "kube-conn-timeout", o.KubeConnTimeout,
		"the timeout for dialing the loopback Kubernetes cluster.")
	fs.Float64Var(&o.KubeConnQPS, "kube-conn-qps", o.KubeConnQPS,
		"the QPS(maximum average number per second) when dialing the loopback Kubernetes cluster.")
	fs.IntVar(&o.KubeConnBurst, "kube-conn-burst", o.KubeConnBurst,
		"the burst(maximum number at the same moment) when dialing the loopback Kubernetes cluster.")
	fs.BoolVar(&o.KubeLeaderElection, "kube-leader-election", o.KubeLeaderElection,
		"the config to determines whether or not to use leader election, "+
			"leader election is primarily used in multi-instance deployments.")
	fs.DurationVar(&o.KubeLeaderLease, "kube-leader-lease", o.KubeLeaderLease,
		"the duration to keep the leadership. "+
			"if --kube-leader-election=false, this flag will be ignored. "+
			"when the network environment is not ideal or do not want to cause frequent access to the cluster, "+
			"please increase the value appropriately.")
	fs.DurationVar(&o.KubeLeaderRenewTimeout, "kube-leader-renew-timeout", o.KubeLeaderRenewTimeout,
		"the duration to renew the leadership before give up, "+
			"must be less than the duration of --kube-leader-lease."+
			"if --kube-leader-election=false, this flag will be ignored. "+
			"when the network environment is not ideal, please increase the value appropriately.")
}

func (o *Options) Validate(_ context.Context) error {
	// Establish.
	if o.BindPort < 1 || o.BindPort > 65535 {
		return errors.New("--secure-port: out of range")
	}
	if o.CertDir != "" && !osx.ExistsDir(o.CertDir) {
		return errors.New("--cert-dir: no found directory")
	}

	// Control.
	if o.GopoolWorkerFactor < 100 {
		return errors.New("--gopool-worker-factor: less than 100")
	}
	if o.InformerCacheResyncPeriod < 5*time.Minute {
		return errors.New("--informer-cache-resync-period: less than 5 minutes")
	}

	// Connect Kubernetes.
	if o.KubeConnTimeout < 10*time.Second {
		return errors.New("--kube-conn-timeout: less than 10 seconds")
	}
	switch {
	case o.KubeConnQPS < 10:
		return errors.New("--kube-conn-qps: less than 10")
	case o.KubeConnBurst < 10:
		return errors.New("--kube-conn-burst: less than 10")
	case float64(o.KubeConnBurst) <= o.KubeConnQPS:
		return errors.New("--kube-conn-burst: less than --kube-conn-qps")
	}
	if o.KubeLeaderElection {
		switch {
		case o.KubeLeaderLease < 5*time.Second:
			return errors.New("--kube-leader-lease: less than 5 seconds")
		case o.KubeLeaderRenewTimeout < 5*time.Second:
			return errors.New("--kube-leader-renew-timeout: less than 5 seconds")
		case o.KubeLeaderLease <= o.KubeLeaderRenewTimeout:
			return errors.New("--kube-leader-lease: less than --kube-leader-renew-timeout")
		}
	}

	return nil
}

func (o *Options) Complete(ctx context.Context) (*Config, error) {
	// Configure goruntime,
	// which is able to query by `runtime.GOMAXPROCS(0)`.
	_, err := maxprocs.Set(maxprocs.Logger(klog.NewStandardLogger("INFO").Printf))
	if err != nil {
		return nil, fmt.Errorf("set maxprocs: %w", err)
	}

	// Configure goroutine pool.
	gopool.Configure(o.GopoolWorkerFactor)

	// Configure system network.
	err = system.ConfigureNetwork()
	if err != nil {
		return nil, fmt.Errorf("configure system network: %w", err)
	}

	// Get loopback config.
	lpCfgPath, lpCliCfg, lpInside, err := kubeconfig.LoadRestConfigNonInteractive()
	if err != nil {
		var (
			embedded    kubernetes.Embedded
			ctx, cancel = context.WithCancel(ctx)
		)
		gopool.Go(func() {
			defer cancel()
			klog.Info("!!! starting embedded Kubernetes !!!")

			err := embedded.Start(ctx)
			if err != nil {
				klog.Error(err, "start embedded Kubernetes")
			}
		})

		lpCfgPath, lpCliCfg, err = embedded.GetConfig(ctx)
		if err != nil {
			return nil, fmt.Errorf("get embedded Kubernetes config: %w", err)
		}
	}

	// Set the timeout, QPS and burst of the rest config.
	lpCliCfg.Timeout = o.KubeConnTimeout
	lpCliCfg.QPS = float32(o.KubeConnQPS)
	lpCliCfg.Burst = o.KubeConnBurst
	lpCliCfg.UserAgent = version.GetUserAgent()
	lpNearby := lpInside || isLoopbackClusterNearby(lpCliCfg)

	// Get loopback client.
	lpHttpCli, err := rest.HTTPClientFor(lpCliCfg)
	if err != nil {
		return nil, fmt.Errorf("create http client: %w", err)
	}

	lpCli, err := clientset.NewForConfigAndClient(rest.CopyConfig(lpCliCfg), lpHttpCli)
	if err != nil {
		return nil, fmt.Errorf("create authorization client: %w", err)
	} else {
		klog.Info("waiting loopback Kubernetes cluster to be available")
		err = kuberest.WaitUntilAvailable(ctx, lpCli.RESTClient())
		if err != nil {
			return nil, fmt.Errorf("wait loopback Kubernetes cluster ready: %w", err)
		}
	}

	// Get serve listener.
	var srvListener net.Listener
	if o.Serve {
		if !lpNearby {
			return nil, errors.New("loopback Kubernetes cluster is not nearby, must provide a nearby cluster")
		}

		tlsCfg := &tls.Config{
			NextProtos: []string{"h2", "http/1.1"},
			MinVersion: tls.VersionTLS12,
		}
		if o.CertDir == "" {
			klog.Info("no cert dir provided, going to use Kubecert to serve")
			certCache, err := certcache.NewK8sCache(ctx, "walrus", lpCli.CoreV1().Secrets(systemkuberes.SystemNamespaceName))
			if err != nil {
				return nil, fmt.Errorf("create cert cache: %w", err)
			}
			certMgr := &kubecert.DynamicManager{
				CertCli: lpCli.CertificatesV1().CertificateSigningRequests(),
				Cache:   certCache,
			}
			tlsCfg.GetCertificate = certMgr.GetCertificate
		} else {
			klog.Info("cert dir provided, going to use cert dir to serve")
			certPath, keyPath := filepath.Join(o.CertDir, "tls.crt"), filepath.Join(o.CertDir, "tls.key")
			certWatcher, err := certwatcher.New(certPath, keyPath)
			if err != nil {
				return nil, fmt.Errorf("create cert watcher: %w", err)
			}
			gopool.Go(func() {
				err := certWatcher.Start(ctx)
				if err != nil {
					klog.Error("certificate watcher error", err)
				}
			})
			tlsCfg.GetCertificate = certWatcher.GetCertificate
		}

		address := net.JoinHostPort(o.BindAddress.String(), stringx.FromInt(o.BindPort))
		srvListener, err = tls.Listen("tcp", address, tlsCfg)
		if err != nil {
			return nil, fmt.Errorf("create tls listener: %w", err)
		}
	}

	system.ConfigureLoopbackKube(lpInside, lpNearby, lpCfgPath, *lpCliCfg, lpCli)

	return &Config{
		InformerCacheResyncPeriod: o.InformerCacheResyncPeriod,
		KubeConfigPath:            lpCfgPath,
		KubeClientConfig:          *lpCliCfg,
		KubeHTTPClient:            lpHttpCli,
		KubeClient:                lpCli,
		KubeLeaderElection:        o.KubeLeaderElection,
		KubeLeaderLease:           o.KubeLeaderLease,
		KubeLeaderRenewTimeout:    o.KubeLeaderRenewTimeout,
		ServeListenerCertDir:      o.CertDir,
		ServeListener:             srvListener,
	}, nil
}
