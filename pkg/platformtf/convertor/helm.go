package convertor

import (
	"fmt"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/platformk8s"
	"github.com/seal-io/seal/pkg/platformtf/block"
	"github.com/seal-io/seal/pkg/platformtf/util"
)

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

func (m HelmConvertor) ToBlock(connector *model.Connector, opts Options) (*block.Block, error) {
	k8sOpts, ok := opts.(K8sConvertorOptions)
	if !ok {
		return nil, fmt.Errorf("invalid k8s options")
	}

	if connector.Type != types.ConnectorTypeK8s {
		return nil, fmt.Errorf("connector type is not k8s")
	}
	var (
		// NB(alex) the config path should keep the same with the secret mount path in deployer.
		configPath = k8sOpts.ConfigPath + "/" + util.GetSecretK8sConfigName(connector.ID.String())
		alias      = k8sOpts.ConnSeparator + connector.ID.String()
		attributes = map[string]interface{}{
			"alias": alias,
		}
	)

	_, _, err := platformk8s.LoadApiConfig(*connector)
	if err != nil {
		return nil, err
	}

	// helm provider need a kubernetes block.
	// it is not a regular attribute of the helm provider.
	// e.g.
	// provider "helm" {
	// 	kubernetes {
	// 		config_path = "xxx"
	// 	}
	// }

	var (
		helmBlock = &block.Block{
			Type:       block.TypeProvider,
			Attributes: attributes,
			// convert the connector type to provider type.
			Labels: []string{ProviderHelm},
		}
		k8sBlock = &block.Block{
			Type: block.TypeK8s,
			Attributes: map[string]interface{}{
				"config_path": configPath,
			},
		}
	)
	helmBlock.AppendBlock(k8sBlock)

	return helmBlock, nil
}

func (m HelmConvertor) ToBlocks(connectors model.Connectors, opts Options) (block.Blocks, error) {
	return connectorsToBlocks(connectors, m.ToBlock, opts)
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
