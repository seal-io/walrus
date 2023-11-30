package deployer

import (
	"context"
	"os"
	"path"

	"k8s.io/client-go/rest"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/k8s/deploy"
	"github.com/seal-io/walrus/pkg/settings"
)

const NameWorkflow = "walrus-workflow"

const (
	imageServer     = "sealio/mirrored-argocli"
	imageController = "sealio/mirrored-workflow-controller"
	imageExecutor   = "sealio/mirrored-argoexec"
	tag             = "v3.5.0"
)

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

func workflow(imageRegistry string) *deploy.ChartApp {
	imageConfig := func(repo string) map[string]any {
		cfg := map[string]any{
			"registry": imageRegistry,
			"tag":      tag,
		}

		if repo != "" {
			cfg["repository"] = repo
		}

		return cfg
	}

	values := map[string]any{
		"images": map[string]any{
			"pullPolicy": "IfNotPresent",
		},
		"server": map[string]any{
			"image":   imageConfig(imageServer),
			"enabled": false,
		},
		"controller": map[string]any{
			"name":  "controller",
			"image": imageConfig(imageController),
			// Note(alex): Disable clusterWorkflowTemplates for now.
			"clusterWorkflowTemplates": map[string]any{
				"enabled": false,
			},
		},
		"executor": map[string]any{
			"image": imageConfig(imageExecutor),
		},
		"fullnameOverride": NameWorkflow,
		"crds": map[string]any{
			"annotations": map[string]any{
				types.LabelWalrusManaged: "true",
			},
		},
		"singleNamespace": true,
		"workflow": map[string]any{
			"rbac": map[string]any{
				"create": false,
			},
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

func DeployArgoWorkflow(ctx context.Context, mc model.ClientSet, config *rest.Config) error {
	imageRegistry := settings.ImageRegistry.ShouldValue(ctx, mc)

	d, err := deploy.New(config)
	if err != nil {
		return err
	}

	return d.EnsureChart(workflow(imageRegistry), false, false)
}
