package helm

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/seal-io/utils/funcx"
	helmaction "helm.sh/helm/v3/pkg/action"
	helmregistry "helm.sh/helm/v3/pkg/registry"
	helmrelease "helm.sh/helm/v3/pkg/release"
	helmdriver "helm.sh/helm/v3/pkg/storage/driver"
	core "k8s.io/api/core/v1"
	kmeta "k8s.io/apimachinery/pkg/api/meta"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

type (
	// Client is a Helm client.
	Client struct {
		// config is the Helm action configuration.
		config helmaction.Configuration
		// namespace is the name of Kubernetes Namespace, which to place the release info.
		namespace string
		// Timeout is the timeout of the action.
		timeout time.Duration
	}

	// ClientOption is a function to set the configuration of the Client.
	ClientOption func(*Client)
)

// WithNamespace returns a ClientOption to set the name of Kubernetes Namespace,
// which to place the release info.
func WithNamespace(namespace string) func(*Client) {
	return func(c *Client) {
		c.namespace = namespace
	}
}

// WithTimeout returns a ClientOption to set the timeout of the action.
func WithTimeout(timeout time.Duration) func(*Client) {
	return func(c *Client) {
		c.timeout = timeout
	}
}

// NewClient creates a new Helm configuration from a Kubernetes rest.Config.
//
// The given namespace is used to place the release info.
func NewClient(restCfg *rest.Config, opts ...ClientOption) (*Client, error) {
	logger := klog.Background().WithName("helm")

	c := &Client{
		namespace: core.NamespaceDefault,
		timeout:   5 * time.Minute,
	}
	for i := range opts {
		opts[i](c)
	}

	// Initialize.
	//
	// NB(thxCode): borrowed from helm.sh/helm/cmd/helm/helm.go
	cg := restClientGetter(*restCfg)
	cd := "secret"
	cdl := func(format string, v ...any) {
		logger.Infof(format, v...)
	}
	err := c.config.Init(cg, c.namespace, cd, cdl)
	if err != nil {
		return nil, err
	}

	// Refill the registry client.
	//
	// NB(thxCode): borrowed from helm.sh/helm/cmd/helm/root.go
	ropts := []helmregistry.ClientOption{
		helmregistry.ClientOptDebug(true),
		helmregistry.ClientOptEnableCache(true),
		helmregistry.ClientOptWriter(logger),
	}
	c.config.RegistryClient, err = helmregistry.NewClient(ropts...)
	if err != nil {
		return nil, fmt.Errorf("create registry client: %w", err)
	}

	return c, nil
}

// Namespace returns the name of Kubernetes Namespace, which to place the release info.
func (c Client) Namespace() string {
	return c.namespace
}

// KubeClientSet returns the Kubernetes clientset.
func (c Client) KubeClientSet() kubernetes.Interface {
	return funcx.MustNoError(c.config.KubernetesClientSet())
}

// Install installs the given chart.
//
// If the release has been found, it will be done.
func (c Client) Install(ctx context.Context, chart *Chart) error {
	next := func(r *helmrelease.Release) NextStepType {
		switch r.Info.Status {
		case helmrelease.StatusDeployed, helmrelease.StatusSuperseded:
			return NextStepDone
		case helmrelease.StatusPendingInstall, helmrelease.StatusPendingUpgrade, helmrelease.StatusPendingRollback:
			if time.Since(r.Info.LastDeployed.Time) > c.timeout {
				return _NextStepInstall
			}
			return NextStepRequeue
		default:
			return NextStepReinstall
		}
	}
	return c.InstallWith(ctx, chart, next)
}

// NextStepType is the type of the next step.
type NextStepType uint8

const (
	NextStepDone NextStepType = iota
	NextStepRequeue
	NextStepUpgrade
	NextStepReinstall
	_NextStepInstall
)

// NextStepConditionFunc is a function to determine what to do next by the given release.
type NextStepConditionFunc func(release *helmrelease.Release) (next NextStepType)

// InstallWith installs the given chart with the given condition function.
//
// The condition function is used to determine what to do next by the given release,
// which only calls when the release has been found.
func (c Client) InstallWith(ctx context.Context, chart *Chart, next NextStepConditionFunc) error {
	// Validate.
	if err := chart.Validate(); err != nil {
		return fmt.Errorf("validate chart: %w", err)
	}
	if next == nil {
		return errors.New("next is required")
	}

	logger := klog.Background().WithName(chart.Name).WithValues("release", chart.Release)

	// Get release.
	g := helmaction.NewGet(&c.config)
	r, err := g.Run(chart.Release)
	if err != nil && !errors.Is(err, helmdriver.ErrReleaseNotFound) {
		return fmt.Errorf("helm get: release %s: %w", chart.Release, err)
	}

	// Next.
	for {
		n := _NextStepInstall
		if r != nil {
			n = next(r)
		}

		switch n {
		case NextStepRequeue:
			// Requeue.
			logger.Info("requeueing")
			time.Sleep(10 * time.Second)
			r, err = g.Run(chart.Release)
			if err != nil && !errors.Is(err, helmdriver.ErrReleaseNotFound) {
				return fmt.Errorf("helm get: release %s: %w", chart.Release, err)
			}
		case NextStepUpgrade:
			// Upgrade.
			logger.Info("upgrading")
			u := helmaction.NewUpgrade(&c.config)
			u.Timeout = c.timeout
			u.Atomic = true
			u.Recreate = true
			u.Force = true
			ch, err := chart.Load(ctx, &c.config)
			if err != nil {
				return fmt.Errorf("helm upgrade: load chart: %w", err)
			}
			vs, err := chart.Values.GetValues(ctx)
			if err != nil {
				return fmt.Errorf("helm upgrade: get values: %w", err)
			}
			r, err = u.RunWithContext(ctx, chart.Release, ch, vs)
			if err != nil {
				return fmt.Errorf("helm upgrade: release %s: %w", chart.Release, err)
			}
			logger.Infof("upgraded: %s", r.Info.Status.String())
		case NextStepReinstall:
			// Uninstall.
			logger.Info("uninstalling")
			ui := helmaction.NewUninstall(&c.config)
			ui.Timeout = c.timeout
			ui.IgnoreNotFound = true
			ui.KeepHistory = false
			ui.Wait = true
			ui.DeletionPropagation = string(meta.DeletePropagationForeground)
			r, err := ui.Run(chart.Release)
			if err != nil && errors.Is(err, helmdriver.ErrReleaseNotFound) {
				return fmt.Errorf("helm uninstall: release %s: %w", chart.Release, err)
			}
			logger.Infof("uninstalled: %s", r.Info)
			fallthrough
		case _NextStepInstall:
			// Install.
			logger.Info("installing")
			i := helmaction.NewInstall(&c.config)
			i.Timeout = c.timeout
			i.ReleaseName = chart.Release
			i.Namespace = c.namespace
			i.Atomic = true
			i.IncludeCRDs = true
			ch, err := chart.Load(ctx, &c.config)
			if err != nil {
				return fmt.Errorf("helm install: load chart: %w", err)
			}
			vs, err := chart.Values.GetValues(ctx)
			if err != nil {
				return fmt.Errorf("helm install: get values: %w", err)
			}
			r, err = i.RunWithContext(ctx, ch, vs)
			if err != nil {
				return fmt.Errorf("helm install: release %s: %w", chart.Release, err)
			}
			logger.Infof("installed: %s", r.Info.Status.String())
		default:
			return nil
		}
	}
}

// restClientGetter is a RESTClientGetter interface implementation for the
// Helm Go packages.
type restClientGetter rest.Config

func (g restClientGetter) ToRESTConfig() (*rest.Config, error) {
	r := rest.Config(g)
	return &r, nil
}

func (g restClientGetter) ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error) {
	config, err := g.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	return memory.NewMemCacheClient(discovery.NewDiscoveryClientForConfigOrDie(config)), nil
}

func (g restClientGetter) ToRESTMapper() (kmeta.RESTMapper, error) {
	discoveryClient, err := g.ToDiscoveryClient()
	if err != nil {
		return nil, err
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(discoveryClient)
	expander := restmapper.NewShortcutExpander(mapper, discoveryClient, nil)

	return expander, nil
}

func (g restClientGetter) ToRawKubeConfigLoader() clientcmd.ClientConfig {
	// Build our config and client.
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)
}
