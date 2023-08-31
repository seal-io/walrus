package deployer

import (
	"os"
	"path"

	"github.com/seal-io/walrus/utils/osx"
)

// container image.
const (
	imageRepositoryServer       = "sealio/mirrored-prometheus"
	imageRepositoryNodeExporter = "sealio/mirrored-node-exporter"
	imageRepositoryKubeState    = "sealio/mirrored-kube-state-metrics"
	imageRepositoryReload       = "sealio/mirrored-prometheus-config-reloader"
)

var imageOpencost = osx.Getenv("IMAGE_OPENCOST", "sealio/mirrored-kubecost-cost-model:v1.105.2")

// prometheus chart.
const (
	//nolint: lll
	defaultPrometheusChartURL = "https://github.com/prometheus-community/helm-charts/releases/download/prometheus-24.0.0/prometheus-24.0.0.tgz"
	defaultPrometheusChart    = "prometheus.tgz"
)

var defaultPrometheusChartPath = func() string {
	// ChartDir from environment variable.
	if cdir := os.Getenv("CHARTS_DIR"); cdir != "" {
		return path.Join(cdir, defaultPrometheusChart)
	}

	// ChartDir from current working directory.
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	cdir := path.Join(pwd, "./.cache", "charts")
	if err = os.MkdirAll(cdir, 0o777); err != nil {
		panic(err)
	}
	return path.Join(cdir, defaultPrometheusChart)
}()
