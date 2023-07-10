package convertor

import (
	"errors"
	"fmt"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	opk8s "github.com/seal-io/seal/pkg/operator/k8s"
	"github.com/seal-io/seal/pkg/terraform/block"
	"github.com/seal-io/seal/pkg/terraform/util"
)

// HelmConvertor mutate the types.ConnectorTypeK8s connector to TypeHelm provider block.
type HelmConvertor string

func (m HelmConvertor) IsSupported(connector *model.Connector) bool {
	return connector.Type == types.ConnectorTypeK8s
}

func (m HelmConvertor) ToBlocks(connectors model.Connectors, opts Options) (block.Blocks, error) {
	var blocks block.Blocks

	for _, c := range connectors {
		if !m.IsSupported(c) {
			continue
		}

		b, err := m.toBlock(c, opts)
		if err != nil {
			return nil, err
		}

		blocks = append(blocks, b)
	}

	return blocks, nil
}

func (m HelmConvertor) toBlock(connector *model.Connector, opts Options) (*block.Block, error) {
	k8sOpts, ok := opts.(K8sConvertorOptions)
	if !ok {
		return nil, errors.New("invalid k8s options")
	}

	if connector.Type != types.ConnectorTypeK8s {
		return nil, fmt.Errorf("connector type is not k8s, connector: %s", connector.ID)
	}

	var (
		// NB(alex) the config path should keep the same with the secret mount path in deployer.
		configPath = k8sOpts.ConfigPath + "/" + util.GetK8sSecretName(connector.ID.String())
		alias      = k8sOpts.ConnSeparator + connector.ID.String()
		attributes = map[string]any{
			"alias": alias,
		}
	)

	_, _, err := opk8s.LoadApiConfig(*connector)
	if err != nil {
		return nil, err
	}

	// Helm provider need a kubernetes block.
	// It is not a regular attribute of the helm provider.
	// E.g.
	// Provider "helm" {
	// 	kubernetes {
	// 		config_path = "xxx"
	// 	}
	// }.

	var (
		helmBlock = &block.Block{
			Type:       block.TypeProvider,
			Attributes: attributes,
			// Convert the connector type to provider type.
			Labels: []string{string(m)},
		}
		k8sBlock = &block.Block{
			Type: block.TypeK8s,
			Attributes: map[string]any{
				"config_path": configPath,
			},
		}
	)

	helmBlock.AppendBlock(k8sBlock)

	return helmBlock, nil
}
