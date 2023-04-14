package convertor

import (
	"errors"

	"github.com/seal-io/seal/pkg/connectors/config"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/platformtf/block"
)

// DefaultConvertor is the convertor for custom category connector.
type DefaultConvertor string

func (m DefaultConvertor) IsSupported(connector *model.Connector) bool {
	return connector.Category == types.ConnectorCategoryCustom &&
		connector.Type == string(m)
}

func (m DefaultConvertor) ToBlocks(connectors model.Connectors, opts Options) (block.Blocks, error) {
	toBlockOpts, ok := opts.(ConvertOptions)
	if !ok {
		return nil, errors.New("invalid convert options")
	}

	var blocks block.Blocks
	for _, conn := range connectors {
		if !m.IsSupported(conn) {
			continue
		}
		b, err := m.toBlock(conn, toBlockOpts)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, b)
	}
	return blocks, nil
}

func (m DefaultConvertor) toBlock(connector *model.Connector, opts ConvertOptions) (*block.Block, error) {
	customConfig, err := config.LoadCustomConfig(connector)
	if err != nil {
		return nil, err
	}

	attributes := customConfig.Attributes
	attributes["alias"] = opts.ConnSeparator + connector.ID.String()
	return &block.Block{
		Type:       block.TypeProvider,
		Labels:     []string{string(m)},
		Attributes: attributes,
	}, nil
}
