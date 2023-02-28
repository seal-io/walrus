package platformk8s

import (
	"fmt"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/utils/strs"
)

// GetConfig returns the rest.Config with the given model.
func GetConfig(conn model.Connector) (*rest.Config, error) {
	var apiConfig, _, err = LoadApiConfig(conn)
	if err != nil {
		return nil, err
	}
	return clientcmd.
		NewNonInteractiveClientConfig(*apiConfig, "", &clientcmd.ConfigOverrides{}, nil).
		ClientConfig()
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
		// }
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

func loadRawConfigV1(data map[string]interface{}) (string, error) {
	// {
	//      "kubeconfig": "..."
	// }
	var kubeconfig, exist = data["kubeconfig"]
	if !exist {
		return "", fmt.Errorf(`not found "kubeconfig"`)
	}
	var kubeconfigText, ok = kubeconfig.(string)
	if !ok {
		return "", fmt.Errorf(`no plain text "kubeconfig"`)
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
