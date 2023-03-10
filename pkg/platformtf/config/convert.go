package config

import (
	"fmt"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
)

var (
	// connector type to terraform provider type map.
	_connectorToTFProvider = map[string]string{
		types.ConnectorTypeK8s: "kubernetes",
	}
)

// convertConnectorToBlock returns provider block for the given connector.
func convertConnectorToBlock(connector *model.Connector, secretMountPath string, connSeparator string) (*Block, error) {
	switch connector.Type {
	case types.ConnectorTypeK8s:
		return convertK8sConnectorToBlock(connector, secretMountPath, connSeparator)
	default:
		return nil, fmt.Errorf("connector type %s is not supported", connector.Type)
	}
}

// convertK8sConnectorToBlock returns kubernetes provider block.
func convertK8sConnectorToBlock(connector *model.Connector, secretMountPath string, connSeparator string) (*Block, error) {
	var attributes = make(map[string]interface{})

	if _, ok := _connectorToTFProvider[connector.Type]; !ok {
		return nil, fmt.Errorf("connector type %s is not supported", connector.Type)
	}
	if _, ok := connector.ConfigData["kubeconfig"]; !ok {
		return nil, fmt.Errorf("kubeconfig is not set for kubernetes connector")
	}

	// NB(alex) the config path should keep the same with the secret mount path in deployer.
	attributes["config_path"] = secretMountPath + "/" + GetConnectorSecretConfigName(connector.ID.String())
	attributes["alias"] = connSeparator + connector.ID.String()

	return &Block{
		Type:       BlockTypeProvider,
		Attributes: attributes,
		// convert the connector type to provider type.
		Labels: []string{_connectorToTFProvider[connector.Type]},
	}, nil
}

// convertModuleToBlock returns module block for the given module and variables.
func convertModuleToBlock(m *model.Module, attributes map[string]interface{}) *Block {
	var block Block
	// TODO need check the module required providers from schema?
	attributes["source"] = m.Source
	block = Block{
		Type:       BlockTypeModule,
		Labels:     []string{m.ID},
		Attributes: attributes,
	}

	return &block
}

// GetConnectorSecretConfigName returns the secret config name for the given connector.
// used for kubernetes connector. or other connectors which need to store the kubeconfig in secret.
func GetConnectorSecretConfigName(connectorID string) string {
	return fmt.Sprintf("config%s", connectorID)
}
