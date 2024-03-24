package convertor

import (
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/terraform/block"
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

		b := toCloudProviderBlock(string(m), c)

		blocks = append(blocks, b)
	}

	return blocks, nil
}
