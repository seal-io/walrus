package kubeconfig

import (
	"errors"
	"fmt"
	"os"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd/api"
)

// ConvertRestConfigToApiConfig converts a rest config to an api config,
// all data locate with a file path will be load into the api config.
func ConvertRestConfigToApiConfig(cfg *rest.Config) (api.Config, error) {
	if cfg == nil {
		return api.Config{}, errors.New("nil rest config")
	}

	const context = "default"

	cc := *api.NewConfig()

	// Convert context info.
	{
		info := api.NewContext()

		info.Cluster = context
		info.AuthInfo = context

		// TODO: Namespace, Extensions.

		cc.Contexts[context] = info
	}

	cc.CurrentContext = context

	// Convert cluster info.
	{
		info := api.NewCluster()

		info.Server = cfg.Host
		info.TLSServerName = cfg.TLSClientConfig.ServerName
		info.InsecureSkipTLSVerify = cfg.Insecure
		info.DisableCompression = cfg.DisableCompression

		info.CertificateAuthorityData = cfg.CAData
		if cfg.CAFile != "" {
			// Load the CA data from the file.
			d, err := os.ReadFile(cfg.CAFile)
			if err != nil {
				return api.Config{}, fmt.Errorf("read CA file: %w", err)
			}
			info.CertificateAuthorityData = d
		}

		if cfg.Proxy != nil {
			// Get the proxy URL with a nil request.
			u, err := cfg.Proxy(nil)
			if err == nil {
				info.ProxyURL = u.String()
			}
		}

		// TODO: Extensions.

		cc.Clusters[context] = info
	}

	// Convert auth info.
	{
		info := api.NewAuthInfo()

		info.ClientCertificateData = cfg.CertData
		if cfg.CertFile != "" {
			// Load the certificate data from the file.
			d, err := os.ReadFile(cfg.CertFile)
			if err != nil {
				return api.Config{}, fmt.Errorf("read certificate file: %w", err)
			}
			info.ClientCertificateData = d
		}

		info.ClientKeyData = cfg.KeyData
		if cfg.KeyFile != "" {
			// Load the key data from the file.
			d, err := os.ReadFile(cfg.KeyFile)
			if err != nil {
				return api.Config{}, fmt.Errorf("read key file: %w", err)
			}
			info.ClientKeyData = d
		}

		info.Token = cfg.BearerToken
		if cfg.BearerTokenFile != "" {
			// Load the token from the file.
			d, err := os.ReadFile(cfg.BearerTokenFile)
			if err != nil {
				return api.Config{}, fmt.Errorf("read bearer token file: %w", err)
			}
			info.Token = string(d)
		}

		info.Impersonate = cfg.Impersonate.UserName
		info.ImpersonateGroups = cfg.Impersonate.Groups
		info.ImpersonateUserExtra = cfg.Impersonate.Extra

		info.Username = cfg.Username
		info.Password = cfg.Password

		info.AuthProvider = cfg.AuthProvider
		info.Exec = cfg.ExecProvider

		// TODO: Extensions.

		cc.AuthInfos[context] = info
	}

	return cc, nil
}
