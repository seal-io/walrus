package convertor

import (
	"errors"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/property"
	"github.com/seal-io/walrus/pkg/terraform/block"
	"github.com/seal-io/walrus/utils/log"
)

type DockerConvertor string

func (m DockerConvertor) IsSupported(connector *model.Connector) bool {
	return connector.Type == types.ConnectorTypeDocker
}

func (m DockerConvertor) ToBlocks(connectors model.Connectors, opts Options) (block.Blocks, error) {
	var blocks block.Blocks

	for _, c := range connectors {
		if !m.IsSupported(c) {
			continue
		}

		b, err := toProviderBlock(string(m), c, opts)
		if err != nil {
			return nil, err
		}

		blocks = append(blocks, b)
	}

	return blocks, nil
}

func toProviderBlock(label string, conn *model.Connector, opts any) (*block.Block, error) {
	convertOpts, ok := opts.(ConvertOptions)
	if !ok {
		return nil, errors.New("invalid options type")
	}

	var (
		alias      = convertOpts.ConnSeparator + conn.ID.String()
		attributes = map[string]any{
			"alias": alias,
		}
		err error
	)

	for k, v := range conn.ConfigData {
		attributes[k], _, err = property.GetString(v.Value)
		if err != nil {
			log.Warnf("error get config data in connector %s:%s, %w", conn.ID, k, err)
		}
	}

	return &block.Block{
		Type:       block.TypeProvider,
		Attributes: attributes,
		Labels:     []string{label},
	}, nil
}
