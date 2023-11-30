package deploy

import (
	"fmt"
	"io"
	"os"
	"time"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/discovery"
	memcached "k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/req"
)

const (
	timeout = 5 * time.Minute
)

type Helm struct {
	logger          log.Logger
	actionConfig    *action.Configuration
	namespace       string
	createNamespace bool
}

type Options struct {
	CreateNamespace bool
	Namespace       string
}

func NewHelm(restCfg *rest.Config, opts Options) (*Helm, error) {
	config := action.Configuration{}
	logger := log.WithName("deployer")

	if err := config.Init(restClientGetter(*restCfg), opts.Namespace, "secrets", func(format string, v ...any) {
		logger.WithName("helm").Debugf(format, v...)
	}); err != nil {
		return nil, err
	}

	return &Helm{
		actionConfig:    &config,
		namespace:       opts.Namespace,
		createNamespace: opts.CreateNamespace,
		logger:          logger,
	}, nil
}

func (h *Helm) Install(name, chartPath string, values map[string]any) error {
	ch, err := loader.Load(chartPath)
	if err != nil {
		return fmt.Errorf("error load chart %s from %s: %w", name, chartPath, err)
	}

	install := action.NewInstall(h.actionConfig)
	install.Force = true
	install.Replace = true
	install.Timeout = timeout
	install.Namespace = h.namespace
	install.CreateNamespace = h.createNamespace
	install.IncludeCRDs = true
	install.ReleaseName = name
	install.Wait = true
	install.WaitForJobs = true

	rel, err := install.Run(ch, values)
	if err != nil {
		return fmt.Errorf("error install chart app %s:%s from %s: %w", h.namespace, name, chartPath, err)
	}

	h.logger.Infof("finished chart install %s:%s, status: %s", h.namespace, name, rel.Info.Status.String())

	return nil
}

func (h *Helm) Uninstall(name string) error {
	uninstall := action.NewUninstall(h.actionConfig)
	uninstall.Wait = true
	uninstall.Timeout = timeout

	_, err := uninstall.Run(name)
	if err != nil {
		return fmt.Errorf("error uninstall %s:%s, %w", h.namespace, name, err)
	}

	h.logger.Infof("finished chart uninstall %s:%s", h.namespace, name)

	return nil
}

func (h *Helm) Download(chartURL, dest string) error {
	h.logger.Debugf("downloading %s", chartURL)

	httpClient := req.HTTP().
		Request().
		WithRedirect()

	resp := httpClient.Get(chartURL)
	if resp.Error() != nil {
		return fmt.Errorf("error download chart %s: %w", chartURL, resp.Error())
	}

	if resp.StatusCode() != 200 {
		errMsg, _ := resp.BodyString()
		return fmt.Errorf("error download chart, http status: %d, body: %s", resp.StatusCode(), errMsg)
	}

	outputFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("error create file %s: %w", dest, err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, resp.BodyOnly())
	if err != nil {
		return fmt.Errorf("error write file %s: %w", dest, err)
	}

	return nil
}

func (h *Helm) GetRelease(name string) (*release.Release, error) {
	get := action.NewGet(h.actionConfig)

	rel, err := get.Run(name)
	if err != nil {
		return nil, fmt.Errorf("error get release %s:%s, %w", h.namespace, name, err)
	}

	return rel, nil
}

func isSucceed(res *release.Release) bool {
	if res.Info == nil {
		return false
	}
	status := res.Info.Status

	return status == release.StatusDeployed || status == release.StatusSuperseded
}

func isUnderway(res *release.Release) bool {
	if res.Info == nil {
		return false
	}
	status := res.Info.Status

	return status == release.StatusUninstalling || status == release.StatusPendingInstall ||
		status == release.StatusPendingUpgrade || status == release.StatusPendingRollback
}

func isFailed(res *release.Release) bool {
	if res.Info == nil {
		return false
	}
	status := res.Info.Status

	return status == release.StatusFailed
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

	return memcached.NewMemCacheClient(discovery.NewDiscoveryClientForConfigOrDie(config)), nil
}

func (g restClientGetter) ToRESTMapper() (meta.RESTMapper, error) {
	discoveryClient, err := g.ToDiscoveryClient()
	if err != nil {
		return nil, err
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(discoveryClient)
	expander := restmapper.NewShortcutExpander(mapper, discoveryClient)

	return expander, nil
}

func (g restClientGetter) ToRawKubeConfigLoader() clientcmd.ClientConfig {
	// Build our config and client.
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)
}
