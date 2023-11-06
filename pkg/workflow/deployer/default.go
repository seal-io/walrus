package deployer

import (
	"context"
	"os"
	"path"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/k8s"
	"github.com/seal-io/walrus/pkg/k8s/deploy"
)

const NameWorkflow = "walrus-workflow"

const (
	defaultWorkflowChartURL = "https://github.com/argoproj/argo-helm/releases/download/" +
		"argo-workflows-0.36.1/argo-workflows-0.36.1.tgz"
	defaultWorkflowChart = "argo-workflows.tgz"
)

var defaultWorkflowChartPath = func() string {
	// ChartDir from environment variable.
	if cdir := os.Getenv("CHARTS_DIR"); cdir != "" {
		return path.Join(cdir, defaultWorkflowChart)
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

	return path.Join(cdir, defaultWorkflowChart)
}()

func workflow() *deploy.ChartApp {
	values := map[string]any{
		"controller": map[string]any{
			"name": "controller",
		},
		"fullnameOverride": NameWorkflow,
		"crds": map[string]any{
			"keep": false,
		},
	}

	return &deploy.ChartApp{
		Name:      NameWorkflow,
		Namespace: types.WalrusSystemNamespace,
		ChartPath: defaultWorkflowChartPath,
		ChartURL:  defaultWorkflowChartURL,
		Values:    values,
	}
}

func DeployArgoWorkflow(ctx context.Context, config *rest.Config) error {
	clientConfig := k8s.ToClientCmdApiConfig(config)

	kubeCfg, err := clientcmd.Write(clientConfig)
	if err != nil {
		return err
	}

	d, err := deploy.New(string(kubeCfg))
	if err != nil {
		return err
	}

	return d.EnsureChart(workflow(), false)
}
