package convertor

import (
	"github.com/seal-io/seal/pkg/connectors/config"
	"github.com/seal-io/seal/pkg/connectors/types"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/platformtf/block"
)

// DefaultConvertor is the convertor for custom category connector.
type DefaultConvertor string

func (m DefaultConvertor) IsSupported(connector *model.Connector) bool {
	return types.IsCustom(connector) && connector.Type == string(m)
}

func (m DefaultConvertor) ToBlocks(connectors model.Connectors, _ Options) (block.Blocks, error) {
	var blocks block.Blocks
	for _, conn := range connectors {
		if !m.IsSupported(conn) {
			continue
		}
		b, err := m.toBlock(conn)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, b)
	}
	return blocks, nil
}

func (m DefaultConvertor) toBlock(connector *model.Connector) (*block.Block, error) {
	customConfig, err := config.LoadCustomConfig(connector)
	if err != nil {
		return nil, err
	}

	return &block.Block{
		Type:       block.TypeProvider,
		Labels:     []string{string(m)},
		Attributes: customConfig.Attributes,
	}, nil
}
