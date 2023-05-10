package platformk8s

import (
	"fmt"
	"time"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types/crypto"
	"github.com/seal-io/seal/utils/strs"
	"github.com/seal-io/seal/utils/version"
)

// ConfigOption holds the modification to modify the return rest.Config.
type ConfigOption func(*rest.Config)

// WithoutTimeout disables the timeout.
func WithoutTimeout() ConfigOption {
	return func(c *rest.Config) {
		c.Timeout = 0
	}
}

// WithTimeout sets the timeout.
func WithTimeout(t time.Duration) ConfigOption {
	return func(c *rest.Config) {
		c.Timeout = t
	}
}

// WithRateLimit sets rate limitation.
func WithRateLimit(qps float32, burst int) ConfigOption {
	return func(c *rest.Config) {
		c.QPS = qps
		c.Burst = burst
	}
}

// GetConfig returns the rest.Config with the given model,
// by default, the rest.Config configures with 15s timeout/16 qps/64 burst,
// please modify the default configuration with ConfigOption as need.
func GetConfig(conn model.Connector, opts ...ConfigOption) (restConfig *rest.Config, err error) {
	apiConfig, _, err := LoadApiConfig(conn)
	if err != nil {
		return nil, err
	}

	restConfig, err = clientcmd.
		NewNonInteractiveClientConfig(*apiConfig, "", &clientcmd.ConfigOverrides{}, nil).
		ClientConfig()
	if err != nil {
		err = fmt.Errorf("cannot construct rest config from api config: %w", err)
		return
	}
	restConfig.Timeout = 15 * time.Second
	restConfig.QPS = 16
	restConfig.Burst = 64
	restConfig.UserAgent = version.GetUserAgent()
	for i := range opts {
		if opts[i] == nil {
			continue
		}
		opts[i](restConfig)
	}

	return
}

// LoadApiConfig returns the client api.Config with the given model.
func LoadApiConfig(conn model.Connector) (apiConfig *api.Config, raw string, err error) {
	switch conn.ConfigVersion {
	default:
		return nil, "", fmt.Errorf("unknown config version: %v", conn.ConfigVersion)
	case "v1":
		// {
		//      "configVersion": "v1",
		//      "configData": {
		//          "kubeconfig": "..."
		//      }
		// }.
		raw, err = loadRawConfigV1(conn.ConfigData)
		if err != nil {
			return nil, "", fmt.Errorf("error load config from connector %s: %w", conn.Name, err)
		}

		apiConfig, err = loadApiConfigV1(raw)
		if err != nil {
			return nil, "", fmt.Errorf("error load version %s config: %w", conn.ConfigVersion, err)
		}
	}
	return
}

func loadRawConfigV1(data crypto.Properties) (string, error) {
	// {
	//      "kubeconfig": "..."
	// }.
	kubeconfigText, exist, err := data["kubeconfig"].GetString()
	if err != nil {
		return "", fmt.Errorf(`failed to get "kubeconfig": %w`, err)
	}
	if !exist {
		return "", fmt.Errorf(`not found "kubeconfig"`)
	}
	return kubeconfigText, nil
}

func loadApiConfigV1(kubeconfigText string) (*api.Config, error) {
	config, err := clientcmd.Load(strs.ToBytes(&kubeconfigText))
	if err != nil {
		return nil, fmt.Errorf("error load api config: %w", err)
	}

	err = clientcmd.Validate(*config)
	if err != nil {
		return nil, fmt.Errorf("error validate api config: %w", err)
	}
	return config, nil
}
