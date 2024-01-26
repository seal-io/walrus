package argoworkflows

import (
	"context"
	"os"
	"path/filepath"

	"k8s.io/client-go/rest"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/k8s/deploy"
	"github.com/seal-io/walrus/pkg/settings"
	"github.com/seal-io/walrus/utils/osx"
)

const (
	chartURL = "https://github.com/argoproj/argo-helm/releases/download/argo-workflows-0.36.1/" +
		"argo-workflows-0.36.1.tgz"
	chartPathEnv = "WORKFLOW_ARGO_WORKFLOWS_CHART_PATH"

	chartValuesTmplText = `
fullnameOverride: "{{ .Name }}"

crds:
  annotations:
    "{{ .ManagedLabel }}": "true"

createAggregateRoles: false

singleNamespace: true

images:
  pullPolicy: "IfNotPresent"

server:
  enabled: false

controller:
  enabled: true
  name: "controller"
  image:
    registry: "{{ .ImageRegistry }}"
    repository: "sealio/mirrored-workflow-controller"
    tag: "{{ .ImageTag }}"

executor:
  image:
    registry: "{{ .ImageRegistry }}"
    repository: "sealio/mirrored-argoexec"
    tag: "{{ .ImageTag }}"

workflow:
  rbac:
    create: true
`
)

func Install(ctx context.Context, mc model.ClientSet, config *rest.Config) error {
	var (
		imgReg = settings.ImageRegistry.ShouldValue(ctx, mc)
		name   = "walrus-workflow"
	)

	d, err := deploy.New(config)
	if err != nil {
		return err
	}

	c := &deploy.ChartApp{
		Name:      name,
		Namespace: types.WalrusSystemNamespace,
		ChartPath: osx.Getenv(chartPathEnv, filepath.Join(os.TempDir(), filepath.Base(chartURL))),
		ChartURL:  chartURL,
		ValuesContext: map[string]any{
			"ManagedLabel":  types.LabelWalrusManaged,
			"Name":          name,
			"ImageRegistry": imgReg,
			"ImageTag":      "v3.5.0",
		},
		ValuesTemplate: chartValuesTmplText,
	}

	return d.EnsureChart(c, false, false)
}
