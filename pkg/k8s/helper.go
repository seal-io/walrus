package k8s

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

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
	var cli, err = kubernetes.NewForConfig(cfg)
	if err != nil {
		return fmt.Errorf("failed to create client via cfg: %w", err)
	}
	return wait.PollImmediateUntilWithContext(ctx, 2*time.Second,
		func(ctx context.Context) (bool, error) {
			var _, err = cli.Discovery().ServerVersion()
			if err != nil {
				log.Warnf("waiting for apiserver to be ready: %v", err)
			}
			return err == nil, ctx.Err()
		},
	)
}
