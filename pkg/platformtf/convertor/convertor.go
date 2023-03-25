package convertor

import (
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/platformtf/block"
)

type (
	Options = interface{}
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
		ToBlocks(model.Connectors, Options) (block.Blocks, error)
	}
)

func connectorsToBlocks(
	connectors model.Connectors,
	h func(*model.Connector, Options) (*block.Block, error),
	opts Options,
) (block.Blocks, error) {
	var blocks block.Blocks
	for _, c := range connectors {
		b, err := h(c, opts)
		if err != nil {
			return nil, err
		}
		if b == nil {
			continue
		}
		blocks = append(blocks, b)
	}
	return blocks, nil
}
