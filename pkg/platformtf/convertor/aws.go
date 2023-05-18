package convertor

import (
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/platformtf/block"
)

type AWSConvertor string

func (m AWSConvertor) IsSupported(connector *model.Connector) bool {
	return connector.Type == types.ConnectorTypeAWS
}

func (m AWSConvertor) ToBlocks(connectors model.Connectors, opts Options) (block.Blocks, error) {
	var blocks block.Blocks

	for _, c := range connectors {
		if !m.IsSupported(c) {
			continue
		}

		b, err := toCloudProviderBlock(string(m), c, opts)
		if err != nil {
			return nil, err
		}

		blocks = append(blocks, b)
	}

	return blocks, nil
}
