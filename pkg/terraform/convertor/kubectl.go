package convertor

import (
	"errors"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	opk8s "github.com/seal-io/seal/pkg/operator/k8s"
	"github.com/seal-io/seal/pkg/terraform/block"
	"github.com/seal-io/seal/pkg/terraform/util"
)

type KubectlConvertor string

func (m KubectlConvertor) IsSupported(connector *model.Connector) bool {
	return connector.Type == types.ConnectorTypeK8s
}

func (m KubectlConvertor) ToBlocks(connectors model.Connectors, opts Options) (block.Blocks, error) {
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

func (m KubectlConvertor) toBlock(connector *model.Connector, opts any) (*block.Block, error) {
	convertOpts, ok := opts.(K8sConvertorOptions)
	if !ok {
		return nil, errors.New("invalid options type")
	}

	var (
		// NB(alex) the config path should keep the same with the secret mount path in deployer.
		configPath = convertOpts.ConfigPath + "/" + util.GetK8sSecretName(connector.ID.String())
		alias      = convertOpts.ConnSeparator + connector.ID.String()
		attributes = map[string]any{
			"config_path": configPath,
			"alias":       alias,
		}
	)

	_, _, err := opk8s.LoadApiConfig(*connector)
	if err != nil {
		return nil, err
	}

	return &block.Block{
		Type:       block.TypeProvider,
		Attributes: attributes,
		// Convert the connector type to provider type.
		Labels: []string{string(m)},
	}, nil
}
