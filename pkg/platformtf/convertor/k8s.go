package convertor

import (
	"fmt"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/platformk8s"
	"github.com/seal-io/seal/pkg/platformtf/block"
	"github.com/seal-io/seal/pkg/platformtf/util"
)

type K8sConvertorOptions struct {
	ConfigPath    string
	ConnSeparator string
}

// K8sConvertor mutate the types.ConnectorTypeK8s connector to ProviderK8s provider block.
type K8sConvertor struct{}

func (m K8sConvertor) ProviderType() string {
	return ProviderK8s
}

func (m K8sConvertor) ConnectorType() string {
	return types.ConnectorTypeK8s
}

func (m K8sConvertor) GetConnectors(connectors model.Connectors) model.Connectors {
	return getConnectorsWithType(m.ConnectorType(), connectors)
}

func (m K8sConvertor) ToBlocks(connectors model.Connectors, opts Options) (block.Blocks, error) {
	return connectorsToBlocks(connectors, m.ToBlock, opts)
}

func (m K8sConvertor) ToBlock(connector *model.Connector, opts Options) (*block.Block, error) {
	k8sOpts, ok := opts.(K8sConvertorOptions)
	if !ok {
		return nil, fmt.Errorf("invalid k8s options")
	}

	if connector.Type != types.ConnectorTypeK8s {
		return nil, fmt.Errorf("connector type is not k8s")
	}

	var (
		// NB(alex) the config path should keep the same with the secret mount path in deployer.
		configPath = k8sOpts.ConfigPath + "/" + util.GetSecretK8sConfigName(connector.ID.String())
		alias      = k8sOpts.ConnSeparator + connector.ID.String()
		attributes = map[string]interface{}{
			"config_path": configPath,
			"alias":       alias,
		}
	)

	_, _, err := platformk8s.LoadApiConfig(*connector)
	if err != nil {
		return nil, err
	}

	return &block.Block{
		Type:       block.TypeProvider,
		Attributes: attributes,
		// convert the connector type to provider type.
		Labels: []string{ProviderK8s},
	}, nil
}
