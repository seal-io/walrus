package deployer

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"path"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/client-go/rest"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/k8s/deploy"
	opk8s "github.com/seal-io/walrus/pkg/operator/k8s"
	"github.com/seal-io/walrus/pkg/settings"
	"github.com/seal-io/walrus/utils/log"
)

const (
	NamePrometheus             = "walrus-prometheus"
	NameOpencost               = "walrus-opencost"
	NameKubeStateMetrics       = "walrus-prometheus-kube-state-metrics"
	NamePrometheusServer       = "walrus-prometheus-server"
	ConfigMapNameOpencost      = "walrus-opencost"
	pathOpencostRefreshPricing = "/refreshPricing"
)

var pathServiceProxy = fmt.Sprintf("/api/v1/namespaces/%s/services/http:%s:9003/proxy",
	types.WalrusSystemNamespace, NameOpencost)

func DeployCostTools(ctx context.Context, mc model.ClientSet, conn *model.Connector, replace bool) error {
	log.WithName("cost").Debugf("deploying cost tools for connector %s", conn.Name)

	apiConfig, kubeconfig, err := opk8s.LoadApiConfig(*conn)
	if err != nil {
		return err
	}

	d, err := deploy.New(kubeconfig)
	if err != nil {
		return fmt.Errorf("error create deployer: %w", err)
	}

	var (
		clusterName   = apiConfig.CurrentContext
		imageRegistry = settings.ImageRegistry.ShouldValue(ctx, mc)
	)

	yaml, err := opencost(clusterName, imageRegistry)
	if err != nil {
		return err
	}

	if err = d.EnsureYaml(ctx, yaml); err != nil {
		return err
	}

	app, err := prometheus(imageRegistry)
	if err != nil {
		return err
	}

	return d.EnsureChart(app, replace)
}

func CostToolsStatus(ctx context.Context, conn *model.Connector) error {
	restCfg, err := opk8s.GetConfig(*conn)
	if err != nil {
		return err
	}

	appv1Client, err := appv1.NewForConfig(restCfg)
	if err != nil {
		return fmt.Errorf("error creating kubernetes core client: %w", err)
	}

	isDeploymentReady := func(namespace, name string) error {
		dep, err := appv1Client.Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("error get tool %s:%s, %w", namespace, name, err)
		}

		if dep.Status.ReadyReplicas != *dep.Spec.Replicas {
			return fmt.Errorf("tool %s:%s, expected %d replica, actual ready %d replica, check deployment details",
				namespace, name, *dep.Spec.Replicas, dep.Status.ReadyReplicas)
		}

		return nil
	}

	for _, v := range []struct {
		namespace string
		name      string
	}{
		{
			namespace: types.WalrusSystemNamespace,
			name:      NameOpencost,
		},
		{
			namespace: types.WalrusSystemNamespace,
			name:      NameKubeStateMetrics,
		},
		{
			namespace: types.WalrusSystemNamespace,
			name:      NamePrometheusServer,
		},
	} {
		err := isDeploymentReady(v.namespace, v.name)
		if err != nil {
			return err
		}
	}

	return nil
}

func opencost(clusterName, imageRegistry string) ([]byte, error) {
	image := path.Join(imageRegistry, imageOpencost)

	data := struct {
		Name               string
		Namespace          string
		ClusterID          string
		PrometheusEndpoint string
		Image              string
	}{
		Name:               NameOpencost,
		Namespace:          types.WalrusSystemNamespace,
		ClusterID:          clusterName,
		PrometheusEndpoint: fmt.Sprintf("http://%s-server.%s.svc:80", NamePrometheus, types.WalrusSystemNamespace),
		Image:              image,
	}

	buf := &bytes.Buffer{}
	if err := tmplOpencost.Execute(buf, data); err != nil {
		return nil, fmt.Errorf("error excute template %s: %w", tmplOpencost.Name(), err)
	}

	return buf.Bytes(), nil
}

func opencostScrape() (string, error) {
	data := struct {
		Name      string
		Namespace string
	}{
		Name:      NameOpencost,
		Namespace: types.WalrusSystemNamespace,
	}

	buf := &bytes.Buffer{}
	if err := tmplPrometheusScrapeJob.Execute(buf, data); err != nil {
		return "", fmt.Errorf("error excute template %s: %w", tmplPrometheusScrapeJob.Name(), err)
	}

	return buf.String(), nil
}

func opencostCustomPricingConfigMap(conn *model.Connector) *v1.ConfigMap {
	data := map[string]string{
		"CPU":     conn.FinOpsCustomPricing.CPU,
		"SpotCPU": conn.FinOpsCustomPricing.SpotCPU,
		"RAM":     conn.FinOpsCustomPricing.RAM,
		"SpotRAM": conn.FinOpsCustomPricing.SpotRAM,
		"GPU":     conn.FinOpsCustomPricing.GPU,
		"Storage": conn.FinOpsCustomPricing.Storage,
	}

	return &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: ConfigMapNameOpencost,
		},
		Data: data,
	}
}

func opencostRefreshPricingURL(restCfg *rest.Config) (string, error) {
	u, err := url.Parse(restCfg.Host)
	if err != nil {
		return "", err
	}

	u.Path = path.Join(u.Path, pathServiceProxy, pathOpencostRefreshPricing)

	return u.String(), nil
}

func prometheus(imageRegistry string) (*deploy.ChartApp, error) {
	scrape, err := opencostScrape()
	if err != nil {
		return nil, err
	}

	imageConfig := func(repo string) map[string]any {
		cfg := map[string]any{
			"registry": imageRegistry,
		}

		if repo != "" {
			cfg["repository"] = repo
		}

		return cfg
	}

	values := map[string]any{
		"prometheus-pushgateway": map[string]any{
			"enabled": false,
		},
		"alertmanager": map[string]any{
			"enabled": false,
		},

		"kube-state-metrics": map[string]any{
			"image": imageConfig(imageRepositoryKubeState),
		},
		"prometheus-node-exporter": map[string]any{
			"image": imageConfig(imageRepositoryNodeExporter),
		},
		"extraScrapeConfigs": scrape,

		// Configmap reload and prometheus only support include registry in repository.
		"configmapReload": map[string]any{
			"prometheus": map[string]any{
				"image": map[string]any{
					"repository": path.Join(imageRegistry, imageRepositoryReload),
				},
			},
		},
		"server": map[string]any{
			"persistentVolume": map[string]any{
				"enabled": false,
			},
			"image": map[string]any{
				"repository": path.Join(imageRegistry, imageRepositoryServer),
			},
		},
	}

	return &deploy.ChartApp{
		Name:      NamePrometheus,
		Namespace: types.WalrusSystemNamespace,
		ChartPath: defaultPrometheusChartPath,
		ChartURL:  defaultPrometheusChartURL,
		Values:    values,
	}, nil
}
