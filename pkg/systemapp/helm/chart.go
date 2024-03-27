package helm

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/seal-io/utils/osx"
	"github.com/seal-io/utils/pools/bytespool"
	helmaction "helm.sh/helm/v3/pkg/action"
	helmchart "helm.sh/helm/v3/pkg/chart"
	helmloader "helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"sigs.k8s.io/yaml"
)

type (
	Chart struct {
		// Name is the name of the chart.
		Name string
		// Version is the version of the chart.
		Version string
		// Release is the name of the release.
		Release string
		// File is the path of the chart tar ball.
		File string
		// FileDownloadURL is the URL to download the chart tar ball.
		// If the File is not existed, the chart will be downloaded from this URL.
		FileDownloadURL string
		// Values is the values to be passed to the chart.
		Values ChartValues
	}

	ChartValues interface {
		// GetValues returns the values for chart installation.
		GetValues(ctx context.Context) (map[string]any, error)
	}
)

// Validate validates the chart.
func (ch Chart) Validate() error {
	if ch.Name == "" {
		return fmt.Errorf("name is required")
	}
	if ch.Release == "" {
		return fmt.Errorf("release name is required")
	}
	if ch.File == "" && ch.FileDownloadURL == "" {
		return fmt.Errorf("file or file download URL is required")
	}
	return nil
}

// Load loads the chart from the file or remote URL.
func (ch Chart) Load(_ context.Context, cfg *helmaction.Configuration) (*helmchart.Chart, error) {
	f := ch.File
	if f != "" && !osx.ExistsFile(f) {
		f = ""
	}

	if f == "" {
		f = filepath.Join(osx.SubTempDir("charts/"+ch.Version), ch.Name)
		if osx.IsEmptyDir(f) {
			p := helmaction.NewPullWithOpts(helmaction.WithConfig(cfg))
			p.Settings = cli.New()
			p.Version = ch.Version
			p.Untar = true
			p.UntarDir = filepath.Dir(f)

			pr, err := p.Run(ch.FileDownloadURL)
			if err != nil {
				return nil, fmt.Errorf("pull chart from %s: %s: %w", ch.FileDownloadURL, pr, err)
			}
		}
	}

	return helmloader.Load(f)
}

// GetValues returns the values for chart installation.
func (ch Chart) GetValues(ctx context.Context) (map[string]any, error) {
	if ch.Values == nil {
		return nil, nil
	}
	return ch.Values.GetValues(ctx)
}

type MapStaticChartValues map[string]any

func (cv MapStaticChartValues) GetValues(ctx context.Context) (map[string]any, error) {
	return cv, nil
}

type YamlTemplateChartValues struct {
	Template string
	Context  map[string]any
}

func (cv YamlTemplateChartValues) GetValues(ctx context.Context) (map[string]any, error) {
	tmpl, err := template.New("values").
		Funcs(templateFuncMap()).
		Parse(cv.Template)
	if err != nil {
		return nil, err
	}

	buf := bytespool.GetBuffer()
	defer bytespool.Put(buf)

	if err = tmpl.Execute(buf, cv.Context); err != nil {
		return nil, err
	}

	vs := map[string]any{}
	err = yaml.Unmarshal(buf.Bytes(), &vs)
	return vs, err
}

func templateFuncMap() template.FuncMap {
	fm := sprig.TxtFuncMap()
	fm["toYaml"] = toYAML

	return fm
}

// toYAML borrows from helm.sh/helm/pkg/engine/engine.go.
func toYAML(v any) string {
	data, err := yaml.Marshal(v)
	if err != nil {
		// Swallow errors inside of a template.
		return ""
	}

	return strings.TrimSuffix(string(data), "\n")
}
