package config

import (
	"fmt"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/platformk8s"
)

type (
	ConvertOptions = interface{}
	// Convertor converts the connector to provider block.
	// e.g. ConnectorType(kubernetes) connector to ProviderType(kubernetes) provider block.
	// ConnectorType(kubernetes) connector to ProviderType(helm) provider block.
	Convertor interface {
		// ProviderType returns the provider type.
		ProviderType() string
		// ConnectorType returns the connector type.
		ConnectorType() string
		// GetConnectors returns the model.Connectors of the provider.
		GetConnectors(model.Connectors) model.Connectors
		// ToBlocks converts the connectors to provider blocks.
		ToBlocks(model.Connectors, ConvertOptions) (Blocks, error)
	}
)

type K8sConvertorOptions struct {
	ConfigPath    string
	ConnSeparator string
}

// K8sConvertor mutate the types.ConnectorTypeK8s connector to ProviderK8s provider block.
type K8sConvertor struct{}

func (m K8sConvertor) ProviderType() string {
	return ProviderK8s
}

func (m K8sConvertor) ConnectorType() string {
	return types.ConnectorTypeK8s
}

func (m K8sConvertor) GetConnectors(connectors model.Connectors) model.Connectors {
	return getConnectorsWithType(m.ConnectorType(), connectors)
}

func (m K8sConvertor) ToBlocks(connectors model.Connectors, opts ConvertOptions) (Blocks, error) {
	var blocks Blocks
	for _, c := range connectors {
		b, err := m.ToBlock(c, opts)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, b)
	}
	return blocks, nil
}

func (m K8sConvertor) ToBlock(connector *model.Connector, opts ConvertOptions) (*Block, error) {
	k8sOpts, ok := opts.(K8sConvertorOptions)
	if !ok {
		return nil, fmt.Errorf("invalid k8s options")
	}

	if connector.Type != types.ConnectorTypeK8s {
		return nil, fmt.Errorf("connector type is not k8s")
	}

	var (
		// NB(alex) the config path should keep the same with the secret mount path in deployer.
		configPath = k8sOpts.ConfigPath + "/" + GetSecretK8sConfigName(connector.ID.String())
		alias      = k8sOpts.ConnSeparator + connector.ID.String()
		attributes = map[string]interface{}{
			"config_path": configPath,
			"alias":       alias,
		}
	)

	_, _, err := platformk8s.LoadApiConfig(*connector)
	if err != nil {
		return nil, err
	}

	return &Block{
		Type:       BlockTypeProvider,
		Attributes: attributes,
		// convert the connector type to provider type.
		Labels: []string{ProviderK8s},
	}, nil
}

// HelmConvertor mutate the types.ConnectorTypeK8s connector to ProviderHelm provider block.
type HelmConvertor struct{}

func (m HelmConvertor) ProviderType() string {
	return ProviderHelm
}

func (m HelmConvertor) ConnectorType() string {
	return types.ConnectorTypeK8s
}

func (m HelmConvertor) GetConnectors(conns model.Connectors) model.Connectors {
	return getConnectorsWithType(m.ConnectorType(), conns)
}

func (m HelmConvertor) ToBlock(connector *model.Connector, opts ConvertOptions) (*Block, error) {
	k8sOpts, ok := opts.(K8sConvertorOptions)
	if !ok {
		return nil, fmt.Errorf("invalid k8s options")
	}

	if connector.Type != types.ConnectorTypeK8s {
		return nil, fmt.Errorf("connector type is not k8s")
	}
	var (
		// NB(alex) the config path should keep the same with the secret mount path in deployer.
		configPath = k8sOpts.ConfigPath + "/" + GetSecretK8sConfigName(connector.ID.String())
		alias      = k8sOpts.ConnSeparator + connector.ID.String()
		attributes = map[string]interface{}{
			"kubernetes": []map[string]interface{}{
				{
					"config_path": configPath,
				},
			},
			"alias": alias,
		}
	)

	_, _, err := platformk8s.LoadApiConfig(*connector)
	if err != nil {
		return nil, err
	}

	return &Block{
		Type:       BlockTypeProvider,
		Attributes: attributes,
		// convert the connector type to provider type.
		Labels: []string{ProviderHelm},
	}, nil
}

func (m HelmConvertor) ToBlocks(connectors model.Connectors, opts ConvertOptions) (Blocks, error) {
	var blocks Blocks
	for _, c := range connectors {
		b, err := m.ToBlock(c, opts)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, b)
	}
	return blocks, nil
}

// getConnectorsWithType returns the connectors for the given connector type.
func getConnectorsWithType(connectorType string, connectors model.Connectors) model.Connectors {
	var matchedConnectors model.Connectors
	for _, c := range connectors {
		if c.Type == connectorType {
			matchedConnectors = append(matchedConnectors, c)
		}
	}

	return matchedConnectors
}
