package convertor

import (
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/property"
	"github.com/seal-io/walrus/pkg/terraform/block"
	"github.com/seal-io/walrus/utils/log"
)

type AlibabaConvertor string

func (m AlibabaConvertor) IsSupported(connector *model.Connector) bool {
	return connector.Type == types.ConnectorTypeAlibaba
}

func (m AlibabaConvertor) ToBlocks(connectors model.Connectors, opts Options) (block.Blocks, error) {
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

func toCloudProviderBlock(label string, conn *model.Connector) *block.Block {
	var (
		attributes = map[string]any{}
		err        error
	)

	for k, v := range conn.ConfigData {
		attributes[k], _, err = property.GetString(v.Value)
		if err != nil {
			log.Warn("error get config data in connector %s:%s, %w", conn.ID, k, err)
		}
	}

	return &block.Block{
		Type:       block.TypeProvider,
		Attributes: attributes,
		Labels:     []string{label},
	}
}
