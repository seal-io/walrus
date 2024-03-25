package convertor

import (
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/terraform/block"
)

type AzureRMConvertor string

func (m AzureRMConvertor) IsSupported(connector *model.Connector) bool {
	return connector.Type == types.ConnectorTypeAzure
}

func (m AzureRMConvertor) ToBlocks(connectors model.Connectors, opts Options) (block.Blocks, error) {
	var blocks block.Blocks

	for _, c := range connectors {
		if !m.IsSupported(c) {
			continue
		}

		b := toCloudProviderBlock(string(m), c)

		b.AppendBlock(&block.Block{
			Type: block.TypeFeatures,
		})

		blocks = append(blocks, b)
	}

	return blocks, nil
}
