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

func installMinio(ctx context.Context, cli *helm.Client, globalValuesContext map[string]any, disable sets.Set[string]) error {
	// NB: please update the following files if changed.
	// - hack/mirror/walrus-images.txt.
	// - pack/walrus/image/Dockerfile.
	// - github.com/seal-io/helm-charts/charts/walrus.

	name := "minio"
	version := "14.1.3"
	if disable.Has(name) {
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
  {{ .ManagedLabel }}: "true"

mode: "standalone"

auth:
  rootUser: "admin"
  rootPassword: {{ randAlphaNum 10 | quote }}

defaultBuckets: "walrus"

image:
  repository: "bitnami/minio"
  tag: "2024.3.26-debian-12-r0"

provisioning: 
  enabled: false

volumePermissions:
  enabled: true
  image: 
    repository: "bitnami/os-shell"
    tag: "12-debian-12-r17"

persistence:
  enabled: true
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
		Values: helm.YamlTemplateChartValues{
			Template: valuesTemplate,
			Context:  valuesContext,
		},
	}
	values, err := cli.Install(ctx, chart)
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

	user, err := values.PathValue("auth.rootUser")
	if err != nil {
		return fmt.Errorf("get root user: %w", err)
	}
	pass, err := values.PathValue("auth.rootPassword")
	if err != nil {
		return fmt.Errorf("get root password: %w", err)
	}
	endpoint := fmt.Sprintf("s3://%s:%s@%s:9000/walrus?sslmode=disable", user, pass, host)
	return systemsetting.ServeObjectStorageUrl.Configure(ctx, endpoint)
}
