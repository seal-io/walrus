package convertor

import (
	"errors"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/terraform/block"
	"github.com/seal-io/seal/utils/log"
)

type CloudProviderConvertorOptions struct {
	ConnSeparator string
}

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

		b, err := toCloudProviderBlock(string(m), c, opts)
		if err != nil {
			return nil, err
		}

		blocks = append(blocks, b)
	}

	return blocks, nil
}

func toCloudProviderBlock(label string, conn *model.Connector, opts any) (*block.Block, error) {
	convertOpts, ok := opts.(CloudProviderConvertorOptions)
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
		attributes[k], _, err = v.GetString()
		if err != nil {
			log.Warn("error get config data in connector %s:%s, %w", conn.ID, k, err)
		}
	}

	return &block.Block{
		Type:       block.TypeProvider,
		Attributes: attributes,
		Labels:     []string{label},
	}, nil
}
