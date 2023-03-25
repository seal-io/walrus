package convertor

import (
	"errors"
	"fmt"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/seal/pkg/connectors/config"
	conntypes "github.com/seal-io/seal/pkg/connectors/types"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/platformtf/block"
)

type CustomConvertorOptions struct {
	// Providers is the required providers of the terraform modules.
	Providers []string `json:"providers"`
}

// CustomConvertor is the convertor for custom category connector.
type CustomConvertor struct{}

func (m CustomConvertor) ProviderType() string {
	return ProviderCustom
}

func (m CustomConvertor) ConnectorType() string {
	return types.ConnectorCategoryCustom
}

func (m CustomConvertor) GetConnectors(connectors model.Connectors) model.Connectors {
	var matchedConnectors model.Connectors
	for _, c := range connectors {
		if conntypes.IsCustom(c) {
			matchedConnectors = append(matchedConnectors, c)
		}
	}

	return matchedConnectors
}

func (m CustomConvertor) ToBlock(connector *model.Connector, opts Options) (*block.Block, error) {
	customOpts, ok := opts.(CustomConvertorOptions)
	if !ok {
		return nil, errors.New("invalid custom options")
	}
	if connector.Category != types.ConnectorCategoryCustom {
		return nil, fmt.Errorf("connector category is not custom connector: %s", connector.ID)
	}

	customConfig, err := config.LoadCustomConfig(connector)
	if err != nil {
		return nil, err
	}

	requiredProviders := sets.NewString(customOpts.Providers...)
	// if the required providers not contains the custom provider type, return nil.
	// NB(alex): add terraform providers that not in the required providers
	// will cause terraform init slowly.
	if !requiredProviders.Has(connector.Type) {
		return nil, nil
	}

	return &block.Block{
		Type:       block.TypeProvider,
		Labels:     []string{connector.Type},
		Attributes: customConfig.Attributes,
	}, nil
}

func (m CustomConvertor) ToBlocks(conns model.Connectors, opts Options) (block.Blocks, error) {
	return connectorsToBlocks(conns, m.ToBlock, opts)
}
