package deployer

import (
	"fmt"
	"os"
	"path"
	"time"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"

	"github.com/seal-io/walrus/utils/log"
)

const (
	timeout = 5 * time.Minute
)

type ChartApp struct {
	Name         string
	Namespace    string
	ChartTgzName string
	Values       map[string]any
	Entry        *repo.Entry
}

type Helm struct {
	settings     *cli.EnvSettings
	actionConfig *action.Configuration
	kubeCfgFile  *os.File
	namespace    string
	repoCache    string
	logger       log.Logger
}

func NewHelm(namespace, kubeconfig string) (*Helm, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	storeBaseDir := path.Join(pwd, "./.cache")
	if err = os.MkdirAll(storeBaseDir, 0o777); err != nil {
		return nil, err
	}

	repoPath := path.Join(storeBaseDir, "repository")
	if err := os.MkdirAll(repoPath, 0o777); err != nil {
		return nil, err
	}

	repoCachePath := path.Join(storeBaseDir, "repository-cache")
	if err := os.MkdirAll(repoCachePath, 0o777); err != nil {
		return nil, err
	}

	kubeconfigFile, err := os.CreateTemp(storeBaseDir, "kubeconfig")
	if err != nil {
		return nil, err
	}

	if _, err = kubeconfigFile.WriteString(kubeconfig); err != nil {
		return nil, err
	}

	settings := cli.New()
	settings.RepositoryConfig = path.Join(repoPath, "repositories.yaml")
	settings.RepositoryCache = repoCachePath
	settings.KubeConfig = kubeconfigFile.Name()

	config := action.Configuration{}
	logger := log.WithName("cost")

	if err = config.Init(settings.RESTClientGetter(), namespace, "secrets", func(format string, v ...any) {
		logger.WithName("helm").Debugf(format, v...)
	}); err != nil {
		return nil, err
	}

	return &Helm{
		settings:     settings,
		actionConfig: &config,
		kubeCfgFile:  kubeconfigFile,
		namespace:    namespace,
		repoCache:    repoCachePath,
		logger:       logger,
	}, nil
}

func (h *Helm) ChartCacheDir() string {
	return h.repoCache
}

func (h *Helm) Download(repoURL, chartName string) (string, error) {
	h.logger.Debugf("downloading %s from %s", chartName, repoURL)
	chartOps := action.ChartPathOptions{
		RepoURL: repoURL,
	}

	outputPath, err := chartOps.LocateChart(chartName, h.settings)
	if err != nil {
		return "", fmt.Errorf("error download chart %s:%s, %w", repoURL, chartName, err)
	}

	return outputPath, nil
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
	install.CreateNamespace = true
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

func (h *Helm) GetRelease(name string) (*release.Release, error) {
	get := action.NewGet(h.actionConfig)

	rel, err := get.Run(name)
	if err != nil {
		return nil, fmt.Errorf("error get release %s:%s, %w", h.namespace, name, err)
	}

	return rel, nil
}

func (h *Helm) Clean() {
	err := os.RemoveAll(h.kubeCfgFile.Name())
	if err != nil {
		h.logger.Warnf("error clean temp kubeconfig %s", h.kubeCfgFile.Name())
	}
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
