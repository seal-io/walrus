package deployer

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"path"

	"helm.sh/helm/v3/pkg/repo"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/client-go/rest"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/platformk8s"
	"github.com/seal-io/seal/utils/log"
)

const (
	NamePrometheus             = "seal-prometheus"
	NameOpencost               = "seal-opencost"
	NameKubeStateMetrics       = "seal-prometheus-kube-state-metrics"
	NamePrometheusServer       = "seal-prometheus-server"
	ConfigMapNameOpencost      = "seal-opencost"
	pathOpencostRefreshPricing = "/refreshPricing"
)

const (
	defaultPrometheusChartTgz = "prometheus-19.6.1.tgz"
	defaultPrometheusRepo     = "https://prometheus-community.github.io/helm-charts"
)

var pathServiceProxy = fmt.Sprintf("/api/v1/namespaces/%s/services/http:%s:9003/proxy",
	types.SealSystemNamespace, NameOpencost)

type input struct {
	Name               string
	Namespace          string
	ClusterID          string
	PrometheusEndpoint string
}

func DeployCostTools(ctx context.Context, conn *model.Connector, replace bool) error {
	log.WithName("cost").Debugf("deploying cost tools for connector %s", conn.Name)

	apiConfig, kubeconfig, err := platformk8s.LoadApiConfig(*conn)
	if err != nil {
		return err
	}

	d, err := New(kubeconfig)
	if err != nil {
		return fmt.Errorf("error create deployer: %w", err)
	}

	clusterName := apiConfig.CurrentContext

	yaml, err := opencost(clusterName)
	if err != nil {
		return err
	}

	if err = d.EnsureYaml(ctx, yaml); err != nil {
		return err
	}

	app, err := prometheus()
	if err != nil {
		return err
	}

	return d.EnsureChart(app, replace)
}

func CostToolsStatus(ctx context.Context, conn *model.Connector) error {
	restCfg, err := platformk8s.GetConfig(*conn)
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
			namespace: types.SealSystemNamespace,
			name:      NameOpencost,
		},
		{
			namespace: types.SealSystemNamespace,
			name:      NameKubeStateMetrics,
		},
		{
			namespace: types.SealSystemNamespace,
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

func opencost(clusterName string) ([]byte, error) {
	data := input{
		Name:               NameOpencost,
		Namespace:          types.SealSystemNamespace,
		ClusterID:          clusterName,
		PrometheusEndpoint: fmt.Sprintf("http://%s-server.%s.svc:80", NamePrometheus, types.SealSystemNamespace),
	}

	buf := &bytes.Buffer{}
	if err := tmplOpencost.Execute(buf, data); err != nil {
		return nil, fmt.Errorf("error excute template %s: %w", tmplOpencost.Name(), err)
	}

	return buf.Bytes(), nil
}

func opencostScrape() (string, error) {
	data := input{
		Name:      NameOpencost,
		Namespace: types.SealSystemNamespace,
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

func prometheus() (*ChartApp, error) {
	scrape, err := opencostScrape()
	if err != nil {
		return nil, err
	}

	values := map[string]interface{}{
		"prometheus-pushgateway": map[string]interface{}{
			"enabled": false,
		},
		"alertmanager": map[string]interface{}{
			"enabled": false,
		},
		"server": map[string]interface{}{
			"persistentVolume": map[string]interface{}{
				"enabled": false,
			},
		},
		"extraScrapeConfigs": scrape,
	}

	entry := &repo.Entry{
		Name: "prometheus",
		URL:  defaultPrometheusRepo,
	}

	return &ChartApp{
		Name:         NamePrometheus,
		Namespace:    types.SealSystemNamespace,
		ChartTgzName: defaultPrometheusChartTgz,
		Values:       values,
		Entry:        entry,
	}, nil
}
