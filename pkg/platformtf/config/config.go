package config

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	"github.com/hashicorp/hcl/v2/hclwrite"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/utils/log"
)

// ModuleConfig is a struct with model.Module and its variables.
type ModuleConfig struct {
	// Name is the name of the app module relationship.
	Name   string
	Module *model.Module
	// Attributes is the attributes of the module.
	Attributes map[string]interface{}
}

// CreateOptions represents the CreateOptions of the Config.
type CreateOptions struct {
	// SecretMountPath is the mount path of the secret in the terraform config.
	SecretMountPath string
	// ConnectorSeparator is the separator of the terraform provider alias name and id.
	ConnectorSeparator string
	// RequiredProviders is the required providers of the terraform config.
	// e.g. ["kubernetes", "helm"]
	RequiredProviders []string

	// Address is the backend address of the terraform config.
	Address string
	// Token is the backend token to authenticate with the Seal Server of the terraform config.
	Token string
	// SkipTLSVerify is the backend cert verification of the terraform config.
	SkipTLSVerify bool

	Connectors    model.Connectors
	ModuleConfigs []*ModuleConfig
}

// Config handles the configuration of application to terraform config.
type Config struct {
	once *sync.Once
	// file is the hclwrite.File of the Config.
	file *hclwrite.File

	// TerraformBlocks terraform blocks like backend, required_providers, etc.
	/**
	  terraform {
	  	backend "http" {
	  		xxx
	  	}
	  	xxx
	  }
	*/
	TFBlocks Blocks
	// Blocks other blocks like provider, module, etc.
	/**
	  provider "aws" {
	  	region = "us-east-1"
	  }

	  module "aws" {
	  	source = "xxx"
	  	region = "us-east-1"
	  }
	*/
	Blocks Blocks
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
	// load backend block
	backendBlock := loadBackendBlock(opts.Address, opts.Token, opts.SkipTLSVerify)
	tfBlocks := Blocks{
		backendBlock,
	}

	// other blocks like provider, module, etc.
	// load provider blocks
	providerBlocks, err := loadProviderBlocks(opts)
	if err != nil {
		return nil, err
	}
	// load module blocks
	moduleBlocks := loadModuleBlocks(opts.ModuleConfigs, providerBlocks)
	providerBlocks = append(providerBlocks, moduleBlocks...)

	c := &Config{
		file:     hclwrite.NewEmptyFile(),
		TFBlocks: tfBlocks,
		Blocks:   providerBlocks,
		once:     &sync.Once{},
	}

	if err = c.validate(); err != nil {
		return nil, err
	}
	if err = c.initBlocks(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) validate() error {
	for _, block := range c.TFBlocks {
		if block.Type == "" {
			return fmt.Errorf("invalid block type: %s", block.Type)
		}
	}

	for _, block := range c.Blocks {
		if block.Type == "" {
			return fmt.Errorf("invalid block type: %s", block.Type)
		}
	}

	return nil
}

// AddBlocks adds a block to the configuration.
func (c *Config) AddBlocks(blocks Blocks) error {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	c.Blocks = append(c.Blocks, blocks...)
	for _, block := range blocks {
		tBlock, err := block.ToHCLBlock()
		if err != nil {
			return err
		}

		c.file.Body().AppendBlock(tBlock)
	}

	return nil
}

// appendTerraformBlocks prints terraform blocks of the configuration. e.g. backend, required_providers, etc.
func (c *Config) appendTerraformBlocks() error {
	// terraform block
	tfBlock := hclwrite.NewBlock(BlockTypeTerraform, nil)
	tfBody := tfBlock.Body()

	for i := 0; i < len(c.TFBlocks); i++ {
		childBlock, err := c.TFBlocks[i].ToHCLBlock()
		if err != nil {
			return err
		}
		tfBody.AppendBlock(childBlock)
		if i != len(c.TFBlocks)-1 {
			tfBody.AppendNewline()
		}
	}

	c.file.Body().AppendBlock(tfBlock)
	c.file.Body().AppendNewline()
	return nil
}

// appendBlocks prints other blocks of the configuration. e.g. provider, module, etc.
func (c *Config) appendBlocks() error {
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

func (c *Config) initBlocks() (err error) {
	c.once.Do(
		func() {
			if c.file == nil {
				c.file = hclwrite.NewEmptyFile()
			}

			// terraform block
			if err = c.appendTerraformBlocks(); err != nil {
				return
			}

			// other blocks
			if err = c.appendBlocks(); err != nil {
				return
			}
		})
	return
}

// WriteTo writes the config to the writer.
func (c *Config) WriteTo(w io.Writer) (int64, error) {
	if err := c.initBlocks(); err != nil {
		return 0, err
	}

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

func loadBackendBlock(address, token string, skipTLSVerify bool) *Block {
	return &Block{
		Type:   BlockTypeBackend,
		Labels: []string{"http"},
		Attributes: map[string]interface{}{
			"address": address,
			// since the seal server using bearer token and
			// terraform backend only support basic auth.
			// we use the token as the password, and let the username be default.
			"username":               _defaultUsername,
			"password":               token,
			"skip_cert_verification": skipTLSVerify,
			// use PUT method to update the state
			"update_method": _updateMethod,
		},
	}
}

// loadProviderBlocks returns config providers to get terraform provider config block.
func loadProviderBlocks(opts CreateOptions) (Blocks, error) {
	return ToProviderBlocks(opts.RequiredProviders, opts.Connectors, ProviderConvertOptions{
		SecretMountPath: opts.SecretMountPath,
		ConnSeparator:   opts.ConnectorSeparator,
	})
}

// loadModuleBlocks returns config modules to get terraform module config block.
func loadModuleBlocks(moduleConfigs []*ModuleConfig, providers Blocks) Blocks {
	var (
		blocks       Blocks
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

		providersMap[name] = fmt.Sprintf("${%s.%s}", name, alias)
	}
	for _, mc := range moduleConfigs {
		block, err := ToModuleBlock(mc)
		if err != nil {
			log.WithName("platformtf").WithName("config").Warnf("get module block failed, %w", mc)
			continue
		}
		// inject providers alias to the module
		block.Attributes["providers"] = providersMap
		blocks = append(blocks, block)
	}

	return blocks
}
