package convertor

import (
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/platformtf/block"
)

type (
	Options = interface{}
	// Convertor converts the connector to provider block.
	// E.g. ConnectorType(kubernetes) connector to ProviderType(kubernetes) provider block.
	// ConnectorType(kubernetes) connector to ProviderType(helm) provider block.
	Convertor interface {
		// IsSupported checks if the connector is supported by the convertor.
		IsSupported(*model.Connector) bool
		// ToBlocks converts the connectors to provider blocks.
		ToBlocks(model.Connectors, Options) (block.Blocks, error)
	}
)
