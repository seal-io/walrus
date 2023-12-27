package deploy

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/yaml"

	"github.com/seal-io/walrus/utils/bytespool"
	"github.com/seal-io/walrus/utils/files"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/version"
)

type Deployer struct {
	logger log.Logger

	restCfg *rest.Config
}

type ChartApp struct {
	Name      string
	Namespace string
	ChartPath string
	ChartURL  string
	Values    map[string]any
	// The following fields used for generate Values.
	ValuesContext  map[string]any
	ValuesTemplate string
}

func New(restCfg *rest.Config) (*Deployer, error) {
	copyCfg := rest.CopyConfig(restCfg)

	copyCfg.Timeout = 0
	copyCfg.QPS = 16
	copyCfg.Burst = 64
	copyCfg.UserAgent = version.GetUserAgent()

	return &Deployer{
		restCfg: copyCfg,
		logger:  log.WithName("k8s").WithName("deployer"),
	}, nil
}

func NewWithKubeconfig(kubeconfig string) (*Deployer, error) {
	restCfg, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeconfig))
	if err != nil {
		return nil, err
	}

	return New(restCfg)
}

func (d *Deployer) EnsureYaml(ctx context.Context, yamlContent []byte) error {
	d.logger.Debugf("ensuring yaml")

	dynamicClient, err := dynamic.NewForConfig(d.restCfg)
	if err != nil {
		return fmt.Errorf("error create dynamic client: %w", err)
	}

	var (
		decoder = scheme.Codecs.UniversalDeserializer()
		yamls   = bytes.Split(yamlContent, []byte("\n---\n"))
		objs    = make([]*unstructured.Unstructured, 0, len(yamls))
	)

	for _, v := range yamls {
		obj := &unstructured.Unstructured{}
		if _, _, err = decoder.Decode(v, nil, obj); err != nil {
			return err
		}

		objs = append(objs, obj)
	}

	for _, obj := range objs {
		var (
			name        = obj.GetName()
			ns          = obj.GetNamespace()
			gvk         = obj.GetObjectKind().GroupVersionKind()
			resource, _ = meta.UnsafeGuessKindToResource(gvk)
		)

		switch {
		case ns != "":
			exist, err := dynamicClient.Resource(resource).Namespace(ns).Get(ctx, name, metav1.GetOptions{})
			if err != nil && !apierrors.IsNotFound(err) {
				return fmt.Errorf("error get namespaced resource %v from yaml: %w", resource, err)
			}

			// Update.
			if exist != nil {
				_, err = dynamicClient.Resource(resource).Namespace(ns).Update(ctx, obj, metav1.UpdateOptions{})
				if err != nil {
					return fmt.Errorf("error update namespaced resource %v from yaml: %w", resource, err)
				}

				continue
			}

			// Create.
			_, err = dynamicClient.Resource(resource).Namespace(ns).Create(ctx, obj, metav1.CreateOptions{})
			if err != nil && !apierrors.IsAlreadyExists(err) {
				return fmt.Errorf("error create namespaced resource %v from yaml: %w", resource, err)
			}
		default:
			exist, err := dynamicClient.Resource(resource).Get(ctx, name, metav1.GetOptions{})
			if err != nil && !apierrors.IsNotFound(err) {
				return fmt.Errorf("error get resource %v from yaml: %w", resource, err)
			}

			if exist != nil {
				_, err = dynamicClient.Resource(resource).Update(ctx, obj, metav1.UpdateOptions{})
				if err != nil {
					return fmt.Errorf("error update resource %v from yaml: %w", resource, err)
				}

				continue
			}

			_, err = dynamicClient.Resource(resource).Create(ctx, obj, metav1.CreateOptions{})
			if err != nil && !apierrors.IsAlreadyExists(err) {
				return fmt.Errorf("error create resource %v from yaml: %w", resource, err)
			}
		}
	}

	return nil
}

func (d *Deployer) EnsureChart(app *ChartApp, createNamespace, replace bool) error {
	d.logger.Debugf("ensuring helm chart")

	helm, err := NewHelm(d.restCfg, Options{
		CreateNamespace: createNamespace,
		Namespace:       app.Namespace,
	})
	if err != nil {
		return err
	}

	res, err := helm.GetRelease(app.Name)
	if err != nil {
		// Error isn't found.
		if !strings.Contains(err.Error(), "not found") {
			return fmt.Errorf("error get release %s:%s, %w", app.Namespace, app.Name, err)
		}
		// Error is not found, continue to install.
	} else {
		switch {
		case isSucceed(res) && !replace:
			return nil
		case isUnderway(res):
			return fmt.Errorf("helm chart %s:%s is under status %s", app.Namespace, app.Name, res.Info.Status)
		case isFailed(res) && !replace:
			return fmt.Errorf("helm chart %s:%s deploy failed: %s", app.Namespace, app.Name, res.Info.Description)
		}
	}

	if replace {
		if err = helm.Uninstall(app.Name); err != nil &&
			!strings.Contains(err.Error(), "not found") {
			return err
		}
	}

	// Download chart if isn't existed.
	if !files.Exists(app.ChartPath) {
		err = helm.Download(app.ChartURL, app.ChartPath)
		if err != nil {
			return err
		}
	}

	// Render values if needed.
	if app.ValuesContext != nil && app.ValuesTemplate != "" {
		tmpl, err := template.New(app.Name + "-values").
			Funcs(funcMap()).
			Parse(app.ValuesTemplate)
		if err != nil {
			return err
		}

		buf := bytespool.GetBuffer()
		defer bytespool.Put(buf)

		if err = tmpl.Execute(buf, app.ValuesContext); err != nil {
			return err
		}

		values := map[string]any{}
		if err = yaml.Unmarshal(buf.Bytes(), &values); err != nil {
			return err
		}

		app.Values = values
	}

	if err = helm.Install(app.Name, app.ChartPath, app.Values); err != nil {
		return err
	}

	return nil
}

func funcMap() template.FuncMap {
	fm := sprig.HermeticTxtFuncMap()
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
