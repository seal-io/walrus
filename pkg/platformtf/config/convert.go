package config

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/seal/pkg/dao/model"
)

const (
	ProviderK8s  = "kubernetes"
	ProviderHelm = "helm"
)

// _providerConvertors mutate the connector to provider block.
var _providerConvertors = make(map[string]Convertor, 0)

type ProviderConvertOptions struct {
	SecretMountPath string
	ConnSeparator   string
}

func init() {
	convertors := []Convertor{
		K8sConvertor{},
		HelmConvertor{},
		// add more convertors
	}

	for _, c := range convertors {
		_providerConvertors[c.ProviderType()] = c
	}
}

func NewConvertor(providerType string) Convertor {
	if _, ok := _providerConvertors[providerType]; !ok {
		return nil
	}

	return _providerConvertors[providerType]
}

// ToProviderBlocks converts the connectors to provider blocks.
func ToProviderBlocks(providers []string, connectors model.Connectors, createOpts ProviderConvertOptions) (Blocks, error) {
	var blocks []*Block
	for _, p := range providers {
		var (
			opts  ConvertOptions
			conns model.Connectors
		)
		switch p {
		case ProviderK8s, ProviderHelm:
			opts = K8sConvertorOptions{
				ConnSeparator: createOpts.ConnSeparator,
				ConfigPath:    createOpts.SecretMountPath,
			}
		default:
			// TODO add more options
		}

		convertor := NewConvertor(p)
		if convertor == nil {
			// it may be a valid use case. For example,
			// a null provider in a module and it doesn't need to match any connector.
			continue
		}

		conns = convertor.GetConnectors(connectors)
		if conns == nil {
			return nil, fmt.Errorf("failed to get connector for provider %s", p)
		}

		convertBlocks, err := convertor.ToBlocks(conns, opts)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, convertBlocks...)
	}

	if !validateRequiredProviders(providers, blocks) {
		return nil, fmt.Errorf("failed to validate providers: %v, current blockProviders: %v", providers, blocks)
	}

	return blocks, nil
}

// validateRequiredProviders validates blocks contains all required providers.
func validateRequiredProviders(providers []string, blocks Blocks) bool {
	var blockProviders []string
	for _, b := range blocks {
		if len(b.Labels) == 0 {
			continue
		}
		providerType := b.Labels[0]
		blockProviders = append(blockProviders, providerType)
	}

	currentProvidersSet := sets.NewString(blockProviders...)
	return currentProvidersSet.HasAll(providers...)
}

// ToModuleBlock returns module block for the given module and variables.
func ToModuleBlock(m *model.Module, attributes map[string]interface{}) *Block {
	var block Block
	attributes["source"] = m.Source
	block = Block{
		Type:       BlockTypeModule,
		Labels:     []string{m.ID},
		Attributes: attributes,
	}

	return &block
}

// GetSecretK8sConfigName returns the secret config name for the given connector.
// used for kubernetes connector. or other connectors which need to store the kubeconfig in secret.
func GetSecretK8sConfigName(connectorID string) string {
	return fmt.Sprintf("config%s", connectorID)
}
