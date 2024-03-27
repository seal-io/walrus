package systemapp

import (
	"context"
	"fmt"
	"path/filepath"

	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/pkg/system"
	"github.com/seal-io/walrus/pkg/systemapp/helm"
	"github.com/seal-io/walrus/pkg/systemsetting"
)

func installMinio(ctx context.Context, cli *helm.Client, globalValuesContext map[string]any, withouts sets.Set[string]) error {
	// NB: please update the following files if changed.
	// - hack/mirror/walrus-images.txt.
	// - pack/walrus/image/Dockerfile.
	name := "minio"
	version := "14.1.3"
	if withouts.Has(name) {
		return nil
	}

	namespace := cli.Namespace()
	release := "walrus-oss"
	file := filepath.Join(system.SubLibDir("charts"), fmt.Sprintf("%s.tgz", name))
	download := fmt.Sprintf("https://charts.bitnami.com/bitnami/%[1]s-%[2]s.tgz", name, version)
	valuesTemplate := `
global:
  imageRegistry: "{{ .ImageRegistry }}"

fullnameOverride: "{{ .Release }}"
namespaceOverride: "{{ .Namespace }}"

commonAnnotations: 
  {{.ManagedLabel}}: "true"

mode: "standalone"

auth:
  rootUser: "admin"
  rootPassword: "admin123"

defaultBuckets: "walrus"

provisioning: 
  enabled: false
`
	valuesContext := globalValuesContext
	valuesContext["Release"] = release
	valuesContext["Namespace"] = namespace

	chart := &helm.Chart{
		Name:            name,
		Version:         version,
		Release:         release,
		File:            file,
		FileDownloadURL: download,
		Values: helm.TemplatedChartValues{
			Template: valuesTemplate,
			Context:  valuesContext,
		},
	}
	err := cli.Install(ctx, chart)
	if err != nil {
		return err
	}

	host := fmt.Sprintf("%s.%s", release, namespace)
	if !system.LoopbackKubeInside.Get() {
		svc, err := cli.KubeClientSet().CoreV1().
			Services(namespace).
			Get(ctx, release, meta.GetOptions{ResourceVersion: "0"})
		if err != nil {
			return fmt.Errorf("get service: %w", err)
		}
		host = svc.Spec.ClusterIP
	}

	endpoint := fmt.Sprintf("s3://admin:admin123@%s:9000/walrus?sslmode=disable", host)
	return systemsetting.ObjectStorageServiceUrl.Configure(ctx, endpoint)
}
