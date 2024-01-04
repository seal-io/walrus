package k8s

import (
	"context"
	"errors"
	"fmt"
	"os"
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
			// NB(shanewxy): Keep the real error from request.
			lastErr = IsConnected(context.TODO(), cli.RESTClient())
			if lastErr != nil {
				log.Warnf("waiting for apiserver to be ready: %v", lastErr)
			}

			return lastErr == nil, ctx.Err()
		},
	)

	if lastErr != nil {
		err = lastErr // Use last error to overwrite context error while existed.
	}

	return err
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
// Warning this helper is only used for the cases that the caller only accept
// a clientcmdapi.Config. Try to use rest.Config directly if possible.
func ToClientCmdApiConfig(restCfg *rest.Config) (*clientcmdapi.Config, error) {
	logger := log.WithName("k8s")

	contextName := "default-context"
	kubeConfig := clientcmdapi.NewConfig()
	kubeConfig.Contexts = map[string]*clientcmdapi.Context{
		contextName: {
			Cluster:  contextName,
			AuthInfo: contextName,
		},
	}
	kubeConfig.Clusters = map[string]*clientcmdapi.Cluster{
		contextName: {
			Server:                   restCfg.Host,
			InsecureSkipTLSVerify:    restCfg.Insecure,
			CertificateAuthorityData: restCfg.CAData,
			CertificateAuthority:     restCfg.CAFile,
		},
	}
	kubeConfig.AuthInfos = map[string]*clientcmdapi.AuthInfo{
		contextName: {
			Token:                 restCfg.BearerToken,
			TokenFile:             restCfg.BearerTokenFile,
			Impersonate:           restCfg.Impersonate.UserName,
			ImpersonateGroups:     restCfg.Impersonate.Groups,
			ImpersonateUserExtra:  restCfg.Impersonate.Extra,
			ClientCertificate:     restCfg.CertFile,
			ClientCertificateData: restCfg.CertData,
			ClientKey:             restCfg.KeyFile,
			ClientKeyData:         restCfg.KeyData,
			Username:              restCfg.Username,
			Password:              restCfg.Password,
			AuthProvider:          restCfg.AuthProvider,
			Exec:                  restCfg.ExecProvider,
		},
	}
	kubeConfig.CurrentContext = contextName

	// Resolve certificate.
	if kubeConfig.Clusters[contextName].CertificateAuthorityData == nil &&
		kubeConfig.Clusters[contextName].CertificateAuthority != "" {
		o, err := os.ReadFile(kubeConfig.Clusters[contextName].CertificateAuthority)
		if err != nil {
			return nil, err
		}

		kubeConfig.Clusters[contextName].CertificateAuthority = ""
		kubeConfig.Clusters[contextName].CertificateAuthorityData = o
	}

	// Fill in data.
	if kubeConfig.AuthInfos[contextName].ClientCertificateData == nil &&
		kubeConfig.AuthInfos[contextName].ClientCertificate != "" {
		o, err := os.ReadFile(kubeConfig.AuthInfos[contextName].ClientCertificate)
		if err != nil {
			logger.Errorf("failed to read client certificate: %v", err)
			return nil, err
		}

		kubeConfig.AuthInfos[contextName].ClientCertificate = ""
		kubeConfig.AuthInfos[contextName].ClientCertificateData = o
	}

	if kubeConfig.AuthInfos[contextName].ClientKeyData == nil && kubeConfig.AuthInfos[contextName].ClientKey != "" {
		o, err := os.ReadFile(kubeConfig.AuthInfos[contextName].ClientKey)
		if err != nil {
			logger.Errorf("failed to read client key: %v", err)
			return nil, err
		}

		kubeConfig.AuthInfos[contextName].ClientKey = ""
		kubeConfig.AuthInfos[contextName].ClientKeyData = o
	}

	return kubeConfig, nil
}
