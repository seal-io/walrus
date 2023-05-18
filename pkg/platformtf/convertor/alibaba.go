package convertor

import (
	"errors"

	"github.com/seal-io/seal/pkg/cloudprovider"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/platformtf/block"
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

func toCloudProviderBlock(label string, conn *model.Connector, opts interface{}) (*block.Block, error) {
	convertOpts, ok := opts.(CloudProviderConvertorOptions)
	if !ok {
		return nil, errors.New("invalid options type")
	}

	cred, err := cloudprovider.CredentialFromConnector(conn)
	if err != nil {
		return nil, err
	}

	var (
		alias      = convertOpts.ConnSeparator + conn.ID.String()
		attributes = map[string]interface{}{
			"access_key": cred.AccessKey,
			"secret_key": cred.SecretKey,
			"region":     cred.Region,
			"alias":      alias,
		}
	)

	return &block.Block{
		Type:       block.TypeProvider,
		Attributes: attributes,
		Labels:     []string{label},
	}, nil
}
