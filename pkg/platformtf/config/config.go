package config

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	"github.com/hashicorp/hcl/v2/hclwrite"

	"github.com/seal-io/seal/pkg/platformtf/block"
	"github.com/seal-io/seal/pkg/platformtf/convertor"
	"github.com/seal-io/seal/utils/log"
)

// Config handles the configuration of application to terraform config.
type Config struct {
	// file is the hclwrite.File of the Config.
	file *hclwrite.File

	// Attributes is the attributes of the Config.
	// e.g.
	// attr1 = "xxx"
	// attr2 = 1
	// attr3 = true
	Attributes map[string]interface{}

	// Blocks blocks like terraform, provider, module, etc.
	/**
	  terraform {
	  	backend "http" {
	  		xxx
	  	}
	  	xxx
	  }
	  provider "aws" {
	  	region = "us-east-1"
	  }

	  module "aws" {
	  	source = "xxx"
	  	region = "us-east-1"
	  }
	*/
	Blocks block.Blocks
}

const (
	// _defaultUsername is the default username in the backend.
	_defaultUsername = "seal"

	// _updateMethod is the method to update state in the backend.
	_updateMethod = "PUT"
)

// NewConfig returns a new Config.
func NewConfig(opts CreateOptions) (*Config, error) {
	// terraform block
	var (
		err        error
		attributes map[string]interface{}
	)

	if opts.Attributes != nil {
		attributes = opts.Attributes
	} else {
		attributes = make(map[string]interface{})
	}

	blocks, err := loadBlocks(opts)
	if err != nil {
		return nil, err
	}

	c := &Config{
		file:       hclwrite.NewEmptyFile(),
		Attributes: attributes,
		Blocks:     blocks,
	}

	if err = c.validate(); err != nil {
		return nil, err
	}

	// init the config.
	if err = c.initAttributes(); err != nil {
		return nil, err
	}
	if err = c.initBlocks(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) validate() error {
	for _, b := range c.Blocks {
		if b.Type == "" {
			return fmt.Errorf("invalid b type: %s", b.Type)
		}
	}

	return nil
}

// AddBlocks adds a block to the configuration.
func (c *Config) AddBlocks(blocks block.Blocks) error {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	c.Blocks = append(c.Blocks, blocks...)
	for _, b := range blocks {
		tBlock, err := b.ToHCLBlock()
		if err != nil {
			return err
		}

		c.file.Body().AppendBlock(tBlock)
	}

	return nil
}

// initBlocks initializes the Blocks of the configuration.
func (c *Config) initBlocks() error {
	for i := 0; i < len(c.Blocks); i++ {
		childBlock, err := c.Blocks[i].ToHCLBlock()
		if err != nil {
			return err
		}
		c.file.Body().AppendBlock(childBlock)
		c.file.Body().AppendNewline()
	}

	return nil
}

// initAttributes initializes the attributes of the configuration.
func (c *Config) initAttributes() error {
	if len(c.Attributes) == 0 {
		return nil
	}

	attributes, err := block.ConvertToCtyWithJson(c.Attributes)
	if err != nil {
		return err
	}
	attrKeys := block.SortValueKeys(attributes)
	if len(attrKeys) == 0 {
		return nil
	}
	attrMap := attributes.AsValueMap()
	for _, attr := range attrKeys {
		c.file.Body().SetAttributeValue(attr, attrMap[attr])
	}

	return nil
}

// WriteTo writes the config to the writer.
func (c *Config) WriteTo(w io.Writer) (int64, error) {

	// format the file
	formatted := hclwrite.Format(Format(c.file.Bytes()))

	return io.Copy(w, bytes.NewReader(formatted))
}

// Reader returns a reader of the config.
func (c *Config) Reader() (io.Reader, error) {
	var buf bytes.Buffer
	if _, err := c.WriteTo(&buf); err != nil {
		return nil, err
	}

	return bytes.NewReader(buf.Bytes()), nil
}

// Bytes returns the bytes of the config.
func (c *Config) Bytes() ([]byte, error) {
	return hclwrite.Format(Format(c.file.Bytes())), nil
}

// loadBlocks loads the blocks of the configuration.
func loadBlocks(opts CreateOptions) (blocks block.Blocks, err error) {
	var (
		tfBlocks       block.Blocks
		providerBlocks block.Blocks
		moduleBlocks   block.Blocks
		variableBlocks block.Blocks
	)
	// load terraform block
	if opts.TerraformOptions != nil {
		tfBlocks = block.Blocks{loadTerraformBlock(opts.TerraformOptions)}
	}
	// other blocks like provider, module, etc.
	// load provider blocks
	if opts.ProviderOptions != nil {
		providerBlocks, err = loadProviderBlocks(opts.ProviderOptions)
		if err != nil {
			return nil, err
		}
	}
	// load module blocks
	if opts.ModuleOptions != nil {
		moduleBlocks = loadModuleBlocks(opts.ModuleOptions.ModuleConfigs, providerBlocks)
	}
	// load variable blocks
	if opts.VariableOptions != nil {
		variableBlocks = loadVariableBlocks(opts.VariableOptions)
	}

	blocks = make(block.Blocks, 0, block.CountLen(tfBlocks, providerBlocks, moduleBlocks, variableBlocks))
	blocks = block.AppendBlocks(blocks, tfBlocks, providerBlocks, moduleBlocks, variableBlocks)

	return
}

// loadTerraformBlock loads the terraform block.
func loadTerraformBlock(opts *TerraformOptions) *block.Block {
	var (
		terraformBlock = &block.Block{
			Type: block.TypeTerraform,
		}
		backendBlock = &block.Block{
			Type:   block.TypeBackend,
			Labels: []string{"http"},
			Attributes: map[string]interface{}{
				"address": opts.Address,
				// since the seal server using bearer token and
				// terraform backend only support basic auth.
				// we use the token as the password, and let the username be default.
				"username":               _defaultUsername,
				"password":               opts.Token,
				"skip_cert_verification": opts.SkipTLSVerify,
				// use PUT method to update the state
				"update_method": _updateMethod,
			},
		}
	)
	terraformBlock.AppendBlock(backendBlock)

	return terraformBlock
}

// loadProviderBlocks returns config providers to get terraform provider config block.
func loadProviderBlocks(opts *ProviderOptions) (block.Blocks, error) {
	return convertor.ToProvidersBlocks(opts.RequiredProviders, opts.Connectors, convertor.ConvertOptions{
		SecretMountPath: opts.SecretMonthPath,
		ConnSeparator:   opts.ConnectorSeparator,
	})
}

// loadModuleBlocks returns config modules to get terraform module config block.
func loadModuleBlocks(moduleConfigs []*ModuleConfig, providers block.Blocks) block.Blocks {
	var (
		logger       = log.WithName("platformtf").WithName("config")
		blocks       block.Blocks
		providersMap = make(map[string]interface{})
	)

	for _, p := range providers {
		alias, ok := p.Attributes["alias"].(string)
		if !ok {
			continue
		}

		if len(p.Labels) == 0 {
			continue
		}
		name := p.Labels[0]
		// template "{{xxx}}" will be replaced by xxx, the quote will be removed.
		providersMap[name] = fmt.Sprintf("{{%s.%s}}", name, alias)
	}
	for _, mc := range moduleConfigs {
		mb, err := ToModuleBlock(mc)
		if err != nil {
			logger.Warnf("get module mb failed, %w", mc)
			continue
		}
		// inject providers alias to the module
		mb.Attributes["providers"] = providersMap
		blocks = append(blocks, mb)
	}

	return blocks
}

// loadVariableBlocks returns config variables to get terraform variable config block.
func loadVariableBlocks(opts *VariableOptions) block.Blocks {
	var (
		// TODO: support other types for secrets and variables
		variableType = "{{string}}"
		blocks       block.Blocks
	)

	// secret variables.
	for _, s := range opts.Secrets {
		blocks = append(blocks, &block.Block{
			Type:   block.TypeVariable,
			Labels: []string{s.Name},
			Attributes: map[string]interface{}{
				"type":      variableType,
				"sensitive": true,
			},
		})
	}

	// application variables.
	for k, v := range opts.Variables {
		if _, ok := v.(string); !ok {
			log.WithName("platformtf").WithName("config").Warnf("application variable %s is not string type, skip", k)
			continue
		}

		blocks = append(blocks, &block.Block{
			Type:   block.TypeVariable,
			Labels: []string{k},
			Attributes: map[string]interface{}{
				"type": variableType,
			},
		})
	}

	return blocks
}

// ToModuleBlock returns module block for the given module and variables.
func ToModuleBlock(mc *ModuleConfig) (*block.Block, error) {
	var b block.Block
	if mc == nil || mc.ModuleVersion == nil {
		return nil, fmt.Errorf("invalid module config: blank")
	}

	if mc.Attributes == nil {
		mc.Attributes = make(map[string]interface{}, 0)
	}

	mc.Attributes["source"] = mc.ModuleVersion.Source
	b = block.Block{
		Type:       block.TypeModule,
		Labels:     []string{mc.Name},
		Attributes: mc.Attributes,
	}

	return &b, nil
}

func CreateConfigToBytes(opts CreateOptions) ([]byte, error) {
	conf, err := NewConfig(opts)
	if err != nil {
		return nil, err
	}

	return conf.Bytes()
}
