package deployer

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/seal-io/seal/utils/log"
	"github.com/seal-io/seal/utils/version"
)

type Deployer struct {
	restCfg *rest.Config
	kubeCfg string
	logger  log.Logger
}

func New(kubeCfg string) (*Deployer, error) {
	clientCfg, err := clientcmd.NewClientConfigFromBytes([]byte(kubeCfg))
	if err != nil {
		return nil, err
	}

	restCfg, err := clientCfg.ClientConfig()
	if err != nil {
		return nil, err
	}
	restCfg.Timeout = 0
	restCfg.QPS = 16
	restCfg.Burst = 64
	restCfg.UserAgent = version.GetUserAgent()

	return &Deployer{
		restCfg: restCfg,
		kubeCfg: kubeCfg,
		logger:  log.WithName("cost").WithName("deployer"),
	}, nil
}

func (d *Deployer) EnsureYaml(ctx context.Context, yamlContent []byte) error {
	d.logger.Debugf("ensuring yaml")

	dynamicClient, err := dynamic.NewForConfig(d.restCfg)
	if err != nil {
		return fmt.Errorf("error create dynamic client: %w", err)
	}

	var (
		objs    []*unstructured.Unstructured
		decoder = scheme.Codecs.UniversalDeserializer()
		yamls   = bytes.Split(yamlContent, []byte("\n---\n"))
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

			// update
			if exist != nil {
				_, err = dynamicClient.Resource(resource).Namespace(ns).Update(ctx, obj, metav1.UpdateOptions{})
				if err != nil {
					return fmt.Errorf("error update namespaced resource %v from yaml: %w", resource, err)
				}
				continue
			}

			// create
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

func (d *Deployer) EnsureChart(app *ChartApp, replace bool) error {
	d.logger.Debugf("ensuring helm chart")

	helm, err := NewHelm(app.Namespace, d.kubeCfg)
	if err != nil {
		return err
	}
	defer helm.Clean()

	res, err := helm.GetRelease(app.Name)
	if err != nil {
		// error isn't found
		if !strings.Contains(err.Error(), "not found") {
			return fmt.Errorf("error get release %s:%s, %w", app.Namespace, app.Name, err)
		}

		// error is not found, continue to install
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
		if err = helm.Uninstall(app.Name); err != nil && !strings.Contains(err.Error(), "not found") {
			return err
		}
	}

	var chartTgzPath = path.Join(helm.repoCache, app.ChartTgzName)
	if _, err = os.Stat(chartTgzPath); err != nil {
		chartTgzPath, err = helm.Download(app.Entry.URL, app.Entry.Name)
		if err != nil {
			return err
		}
	}

	if err = helm.Install(app.Name, chartTgzPath, app.Values); err != nil {
		return err
	}
	return nil
}
