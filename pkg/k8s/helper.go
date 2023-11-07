package k8s

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	coreclient "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"github.com/seal-io/walrus/utils/json"
	"github.com/seal-io/walrus/utils/log"
)

func GetConfig(kubeconfigPath string) (*rest.Config, error) {
	// Use the specified config.
	if kubeconfigPath != "" {
		return LoadConfig(kubeconfigPath)
	}

	// Try the in-cluster config.
	if c, err := rest.InClusterConfig(); err == nil {
		return c, nil
	}

	// Try the recommended config.
	loader := clientcmd.NewDefaultClientConfigLoadingRules()
	loader.Precedence = append(loader.Precedence,
		filepath.Join(home, clientcmd.RecommendedHomeDir, clientcmd.RecommendedFileName))

	return loadConfig(loader)
}

func LoadConfig(kubeconfigPath string) (*rest.Config, error) {
	if kubeconfigPath == "" {
		return nil, errors.New("blank kubeconfig path")
	}

	loader := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath}

	return loadConfig(loader)
}

func loadConfig(loader clientcmd.ClientConfigLoader) (*rest.Config, error) {
	overrides := &clientcmd.ConfigOverrides{}

	return clientcmd.
		NewNonInteractiveDeferredLoadingClientConfig(loader, overrides).
		ClientConfig()
}

func Wait(ctx context.Context, cfg *rest.Config) error {
	cli, err := coreclient.NewForConfig(cfg)
	if err != nil {
		return fmt.Errorf("failed to create client via cfg: %w", err)
	}

	var lastErr error

	err = wait.PollUntilContextCancel(ctx, 1*time.Second, true,
		func(ctx context.Context) (bool, error) {
			lastErr = IsConnected(ctx, cli.RESTClient())
			if lastErr != nil {
				log.Warnf("waiting for apiserver to be ready: %v", lastErr)
			}

			return lastErr == nil, ctx.Err()
		},
	)
	if err != nil {
		if lastErr != nil {
			err = lastErr // Use last error to overwrite context error while existed.
		}

		return err
	}

	return nil
}

func IsConnected(ctx context.Context, r rest.Interface) error {
	body, err := r.Get().
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

// ToClientCmdApiConfig using a rest.Config to generate a clientcmdapi.Config.
func ToClientCmdApiConfig(restConfig *rest.Config) clientcmdapi.Config {
	clusters := make(map[string]*clientcmdapi.Cluster)
	clusters["default-cluster"] = &clientcmdapi.Cluster{
		Server:                   restConfig.Host,
		CertificateAuthorityData: restConfig.CAData,
	}
	contexts := make(map[string]*clientcmdapi.Context)
	contexts["default-context"] = &clientcmdapi.Context{
		Cluster:  "default-cluster",
		AuthInfo: "default-user",
	}
	authinfos := make(map[string]*clientcmdapi.AuthInfo)
	authinfos["default-user"] = &clientcmdapi.AuthInfo{
		ClientCertificateData: restConfig.CertData,
		ClientKeyData:         restConfig.KeyData,
	}
	clientConfig := clientcmdapi.Config{
		Kind:           "Config",
		APIVersion:     "v1",
		Clusters:       clusters,
		Contexts:       contexts,
		CurrentContext: "default-context",
		AuthInfos:      authinfos,
	}

	return clientConfig
}
