package k8s

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/seal-io/seal/utils/json"
	"github.com/seal-io/seal/utils/log"
)

func GetConfig(kubeconfigPath string) (*rest.Config, error) {
	// use the specified config.
	if kubeconfigPath != "" {
		return LoadConfig(kubeconfigPath)
	}

	// try the in-cluster config.
	if c, err := rest.InClusterConfig(); err == nil {
		return c, nil
	}

	// try the recommended config.
	var loader = clientcmd.NewDefaultClientConfigLoadingRules()
	loader.Precedence = append(loader.Precedence,
		filepath.Join(home, clientcmd.RecommendedHomeDir, clientcmd.RecommendedFileName))
	return loadConfig(loader)
}

func LoadConfig(kubeconfigPath string) (*rest.Config, error) {
	if kubeconfigPath == "" {
		return nil, errors.New("blank kubeconfig path")
	}

	var loader = &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath}
	return loadConfig(loader)
}

func loadConfig(loader clientcmd.ClientConfigLoader) (*rest.Config, error) {
	var overrides = &clientcmd.ConfigOverrides{}
	return clientcmd.
		NewNonInteractiveDeferredLoadingClientConfig(loader, overrides).
		ClientConfig()
}

func Wait(ctx context.Context, cfg *rest.Config) error {
	var cli, err = discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return fmt.Errorf("failed to create client via cfg: %w", err)
	}

	var lastErr error
	err = wait.PollImmediateUntilWithContext(ctx, 1*time.Second,
		func(ctx context.Context) (bool, error) {
			lastErr = IsConnected(ctx, cli.RESTClient())
			if lastErr != nil {
				log.Warnf("waiting for apiserver to be ready: %v", lastErr)
			}
			return lastErr == nil, ctx.Err()
		},
	)
	if err != nil && lastErr != nil {
		err = lastErr // use last error to overwrite context error while existed
	}
	return err
}

func IsConnected(ctx context.Context, r rest.Interface) error {
	var body, err = r.Get().
		AbsPath("/version").
		Do(ctx).
		Raw()
	if err != nil {
		return err
	}
	var info struct {
		Major    string `json:"major"`
		Minor    string `json:"minor"`
		Compiler string `json:"compiler"`
		Platform string `json:"platform"`
	}
	err = json.Unmarshal(body, &info)
	if err != nil {
		return fmt.Errorf("unable to parse the server version: %w", err)
	}
	return nil
}
