package hermitcrab

import (
	"context"
	"os"
	"path/filepath"

	core "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/k8s/deploy"
	"github.com/seal-io/walrus/pkg/settings"
	"github.com/seal-io/walrus/utils/osx"
)

const (
	chartURL = "https://github.com/seal-io/helm-charts/releases/download/hermitcrab-0.1.1/" +
		"hermitcrab-0.1.1.tgz"
	chartPathEnv = "DEPLOYER_TERRAFORM_HERMITCRAB_CHART_PATH"

	chartValuesTmplText = `
global:
  imageRegistry: "{{ .ImageRegistry }}"

fullnameOverride: "{{ .Name }}"

commonAnnotations: 
  {{.ManagedLabel}}: "true"

{{ if .Env }}
hermitcrab:
  env: {{ toYaml .Env | nindent 2 }}
{{- end }}
`
)

func Install(ctx context.Context, mc model.ClientSet, config *rest.Config) error {
	var (
		imgReg = settings.ImageRegistry.ShouldValue(ctx, mc)
		name   = "walrus-mirror"
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
			"Env":           getEnv(ctx, mc),
		},
		ValuesTemplate: chartValuesTmplText,
	}

	return d.EnsureChart(c, false, false)
}

func getEnv(ctx context.Context, mc model.ClientSet) (env []core.EnvVar) {
	if v := settings.DeployerAllProxy.ShouldValue(ctx, mc); v != "" {
		env = append(env, core.EnvVar{
			Name:  "ALL_PROXY",
			Value: v,
		})
	}

	if v := settings.DeployerHttpProxy.ShouldValue(ctx, mc); v != "" {
		env = append(env, core.EnvVar{
			Name:  "HTTP_PROXY",
			Value: v,
		})
	}

	if v := settings.DeployerHttpsProxy.ShouldValue(ctx, mc); v != "" {
		env = append(env, core.EnvVar{
			Name:  "HTTPS_PROXY",
			Value: v,
		})
	}

	if v := settings.DeployerNoProxy.ShouldValue(ctx, mc); v != "" {
		env = append(env, core.EnvVar{
			Name:  "NO_PROXY",
			Value: v,
		})
	}

	return env
}
