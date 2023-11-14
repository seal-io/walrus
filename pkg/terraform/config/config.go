package config

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	"github.com/hashicorp/hcl/v2/hclwrite"

	"github.com/seal-io/walrus/pkg/templates/translator"
	"github.com/seal-io/walrus/pkg/terraform/block"
	"github.com/seal-io/walrus/pkg/terraform/convertor"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/strs"
)

// Config handles the configuration of application to terraform config.
type Config struct {
	// File is the hclwrite.File of the Config.
	file *hclwrite.File

	// Attributes is the attributes of the Config.
	// E.g.
	// Attr1 = "xxx"
	// attr2 = 1
	// attr3 = true.
	Attributes map[string]any

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
	// Terraform block.
	var (
		err        error
		attributes map[string]any
	)

	if opts.Attributes != nil {
		attributes = opts.Attributes
	} else {
		attributes = make(map[string]any)
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

	// Init the config.
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

	translator := translator.NewTerraformTranslator()

	attrKeys, attrMap, err := translator.ToOriginalTypeValues(c.Attributes)
	if err != nil {
		return err
	}

	for _, attr := range attrKeys {
		c.file.Body().SetAttributeValue(attr, attrMap[attr])
	}

	return nil
}

// WriteTo writes the config to the writer.
func (c *Config) WriteTo(w io.Writer) (int64, error) {
	// Format the file.
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
		outputBlocks   block.Blocks
	)
	// Load terraform block.
	if opts.TerraformOptions != nil {
		tfBlocks = block.Blocks{loadTerraformBlock(opts.TerraformOptions)}
	}
	// Other blocks like provider, module, etc.
	// load provider blocks.
	if opts.ProviderOptions != nil {
		providerBlocks, err = loadProviderBlocks(opts.ProviderOptions)
		if err != nil {
			return nil, err
		}
	}
	// Load module blocks.
	if opts.ModuleOptions != nil {
		moduleBlocks = loadModuleBlocks(opts.ModuleOptions.ModuleConfigs, providerBlocks)
	}
	// Load variable blocks.
	if opts.VariableOptions != nil {
		variableBlocks = loadVariableBlocks(opts.VariableOptions)
	}

	if len(opts.OutputOptions) != 0 {
		outputBlocks = loadOutputBlocks(opts.OutputOptions)
	}

	blocks = make(block.Blocks, 0, block.CountLen(tfBlocks, providerBlocks, moduleBlocks, variableBlocks, outputBlocks))
	blocks = block.AppendBlocks(blocks, tfBlocks, providerBlocks, moduleBlocks, variableBlocks, outputBlocks)

	return blocks, nil
}

// loadTerraformBlock loads the terraform block.
func loadTerraformBlock(opts *TerraformOptions) *block.Block {
	var (
		logger         = log.WithName("deployer").WithName("tf")
		terraformBlock = &block.Block{
			Type: block.TypeTerraform,
		}
	)

	if opts.ProviderRequirements != nil {
		requiredProviders := &block.Block{
			Type:       block.TypeRequiredProviders,
			Attributes: map[string]any{},
		}
		for provider, requirement := range opts.ProviderRequirements {
			if _, ok := requiredProviders.Attributes[provider]; ok {
				logger.Warnf("provider already exists, skip", "provider", provider)
				continue
			}
			pr := make(map[string]any)

			if requirement != nil {
				if requirement != nil && len(requirement.VersionConstraints) != 0 {
					pr["version"] = strs.Join(",", requirement.VersionConstraints...)
				}

				if requirement != nil && requirement.Source != "" {
					pr["source"] = requirement.Source
				}
			}
			requiredProviders.Attributes[provider] = pr
		}

		terraformBlock.AppendBlock(requiredProviders)
	}
	backendBlock := &block.Block{
		Type:   block.TypeBackend,
		Labels: []string{"http"},
		Attributes: map[string]any{
			"address": opts.Address,
			// Since the seal server using bearer token and
			// terraform backend only support basic auth.
			// We use the token as the password, and let the username be default.
			"username":               _defaultUsername,
			"password":               opts.Token,
			"skip_cert_verification": opts.SkipTLSVerify,
			// Use PUT method to update the state.
			"update_method":  _updateMethod,
			"retry_max":      10,
			"retry_wait_max": 5,
		},
	}

	terraformBlock.AppendBlock(backendBlock)

	return terraformBlock
}

// loadProviderBlocks returns config providers to get terraform provider config block.
func loadProviderBlocks(opts *ProviderOptions) (block.Blocks, error) {
	return convertor.ToProvidersBlocks(opts.RequiredProviderNames, opts.Connectors, convertor.ConvertOptions{
		SecretMountPath: opts.SecretMonthPath,
		ConnSeparator:   opts.ConnectorSeparator,
		Providers:       opts.RequiredProviderNames,
	})
}

// loadModuleBlocks returns config modules to get terraform module config block.
func loadModuleBlocks(moduleConfigs []*ModuleConfig, providers block.Blocks) block.Blocks {
	var (
		logger       = log.WithName("deployer").WithName("tf").WithName("config")
		blocks       block.Blocks
		providersMap = make(map[string]any)
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
		// Template "{{xxx}}" will be replaced by xxx, the quote will be removed.
		providersMap[name] = fmt.Sprintf("{{%s.%s}}", name, alias)
	}

	for _, mc := range moduleConfigs {
		mb, err := ToModuleBlock(mc)
		if err != nil {
			logger.Warnf("get module mb failed, %w", mc)
			continue
		}
		// Inject providers alias to the module.
		if len(mc.SchemaData.RequiredProviders) != 0 {
			moduleProviders := map[string]any{}

			for _, p := range mc.SchemaData.RequiredProviders {
				if _, ok := providersMap[p.Name]; !ok {
					logger.Warnf("provider not found, skip provider: %s", p.Name)
					continue
				}
				moduleProviders[p.Name] = providersMap[p.Name]
			}
			mb.Attributes["providers"] = moduleProviders
		}

		blocks = append(blocks, mb)
	}

	return blocks
}

// loadVariableBlocks returns config variables to get terraform variable config block.
func loadVariableBlocks(opts *VariableOptions) block.Blocks {
	blocks := make(block.Blocks, 0, len(opts.Variables)+len(opts.DependencyOutputs))

	// Secret variables.
	for name, sensitive := range opts.Variables {
		blocks = append(blocks, &block.Block{
			Type:   block.TypeVariable,
			Labels: []string{opts.VariablePrefix + name},
			Attributes: map[string]any{
				"type":      "{{string}}",
				"sensitive": sensitive,
			},
		})
	}

	// Dependency variables.
	for k, o := range opts.DependencyOutputs {
		blocks = append(blocks, &block.Block{
			Type:   block.TypeVariable,
			Labels: []string{opts.ResourcePrefix + k},
			Attributes: map[string]any{
				"type":      `{{string}}`,
				"sensitive": o.Sensitive,
			},
		})
	}

	return blocks
}

// loadOutputBlocks returns terraform outputs config block.
func loadOutputBlocks(opts OutputOptions) block.Blocks {
	blockConfig := func(output Output) (string, string) {
		label := fmt.Sprintf("%s_%s", output.ResourceName, output.Name)
		value := fmt.Sprintf(`{{module.%s.%s}}`, output.ResourceName, output.Name)

		return label, value
	}

	// Template output.
	blocks := make(block.Blocks, 0, len(opts))

	for _, o := range opts {
		label, value := blockConfig(o)

		blocks = append(blocks, &block.Block{
			Type:   block.TypeOutput,
			Labels: []string{label},
			Attributes: map[string]any{
				"value":     value,
				"sensitive": o.Sensitive,
			},
		})
	}

	return blocks
}

// ToModuleBlock returns module block for the given module and variables.
func ToModuleBlock(mc *ModuleConfig) (*block.Block, error) {
	var b block.Block

	if mc.Attributes == nil {
		mc.Attributes = make(map[string]any, 0)
	}

	mc.Attributes["source"] = mc.Source
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
