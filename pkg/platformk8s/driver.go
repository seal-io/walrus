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
	var apiConfig, err = LoadApiConfig(conn)
	if err != nil {
		return nil, err
	}
	return clientcmd.
		NewNonInteractiveClientConfig(*apiConfig, "", &clientcmd.ConfigOverrides{}, nil).
		ClientConfig()
}

// LoadApiConfig returns the client api.Config with the given model.
func LoadApiConfig(conn model.Connector) (apiConfig *api.Config, err error) {
	switch conn.ConfigVersion {
	default:
		return nil, fmt.Errorf("unknown config version: %v", conn.ConfigVersion)
	case "v1":
		// {
		//      "configVersion": "v1",
		//      "configData": {
		//          "kubeconfig": "..."
		//      }
		// }
		apiConfig, err = loadApiConfigV1(conn.ConfigData)
	}
	if err != nil {
		return nil, fmt.Errorf("error load version %s config: %w", conn.ConfigVersion, err)
	}
	return
}

func loadApiConfigV1(data map[string]interface{}) (*api.Config, error) {
	// {
	//      "kubeconfig": "..."
	// }
	var kubeconfig, exist = data["kubeconfig"]
	if !exist {
		return nil, fmt.Errorf(`not found "kubeconfig"`)
	}
	var kubeconfigText, ok = kubeconfig.(string)
	if !ok {
		return nil, fmt.Errorf(`no plain text "kubeconfig"`)
	}
	return clientcmd.Load(strs.ToBytes(&kubeconfigText))
}
