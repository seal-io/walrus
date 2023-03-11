package config

import (
	"bytes"
	"fmt"
	"io"
	"text/template"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/utils/bytespool"
	"github.com/seal-io/seal/utils/strs"
)

// ModuleConfig is a struct with model.Module and its variables.
type ModuleConfig struct {
	Module *model.Module
	// Attributes is the attributes of the module.
	Attributes map[string]interface{}
}

// CreateOptions represents the CreateOptions of the Config.
type CreateOptions struct {
	// Format is the default print format of the Config, support hcl, json.
	Format string
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
	format       string
	outputBuffer bytes.Buffer
	// mapObjects is a map of objects that have been printed already.
	mapObjects map[string]struct{}

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
		format:     opts.Format,
		mapObjects: make(map[string]struct{}),
		TFBlocks:   tfBlocks,
		Blocks:     providerBlocks,
	}

	if err = c.validate(); err != nil {
		return nil, err
	}
	if err = c.initOutput(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) validate() error {
	if c.format != BlockFormatHCL && c.format != BlockFormatJSON {
		return fmt.Errorf("invalid format: %s", c.format)
	}

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

// AddBlock adds a block to the configuration.
func (c *Config) AddBlock(blocks Blocks) error {
	var blocksOutput []byte

	c.Blocks = append(c.Blocks, blocks...)
	for _, block := range blocks {
		output, err := block.Print(c.format, c.mapObjects)
		if err != nil {
			return err
		}

		blocksOutput = append(blocksOutput, output...)
	}
	_, err := c.outputBuffer.Write(blocksOutput)
	if err != nil {
		return err
	}

	return nil
}

// printTerraformBlocks prints terraform blocks of the configuration. e.g. backend, required_providers, etc.
func (c *Config) printTerraformBlocks() ([]byte, error) {
	terraformTpl := `terraform {
{{ range $block := .Blocks }}
{{ $block }}
{{ end -}}
}

`
	data := struct {
		Blocks []string
	}{}
	tpl, err := template.New("terraform").Parse(terraformTpl)
	if err != nil {
		return nil, err
	}

	for _, block := range c.TFBlocks {
		blockOutput, err := block.Print(c.format, c.mapObjects)
		if err != nil {
			return nil, err
		}

		data.Blocks = append(data.Blocks, strs.Indent(2, string(blockOutput)))
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// printBlocks prints other blocks of the configuration. e.g. provider, module, etc.
func (c *Config) printBlocks() ([]byte, error) {
	var b = bytespool.GetBuffer()
	for _, block := range c.Blocks {
		blockOutput, err := block.Print(c.format, c.mapObjects)
		if err != nil {
			return nil, err
		}
		blockOutput = append(blockOutput, []byte("\n")...)
		_, err = b.Write(blockOutput)
		if err != nil {
			return nil, err
		}
	}

	return b.Bytes(), nil
}

func (c *Config) initOutput() error {
	tfOutput, err := c.printTerraformBlocks()
	if err != nil {
		return err
	}

	blocksOutput, err := c.printBlocks()
	if err != nil {
		return err
	}

	tfOutput = append(tfOutput, blocksOutput...)
	tfOutput = Format(tfOutput)

	_, err = c.outputBuffer.Write(tfOutput)
	if err != nil {
		return err
	}

	return nil
}

// WriteTo writes the config to the writer.
func (c *Config) WriteTo(w io.Writer) (int64, error) {
	return c.outputBuffer.WriteTo(w)
}

// Reader returns a reader of the config.
func (c *Config) Reader() (io.Reader, error) {
	if c.outputBuffer.Len() == 0 {
		err := c.initOutput()
		if err != nil {
			return nil, err
		}
	}

	return bytes.NewReader(c.outputBuffer.Bytes()), nil
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
func loadModuleBlocks(moduleBlocks []*ModuleConfig, providers Blocks) Blocks {
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

		providersMap[name] = fmt.Sprintf("$${%s.%s}", name, alias)
	}
	for _, m := range moduleBlocks {
		block := ToModuleBlock(m.Module, m.Attributes)
		// inject providers alias to the module
		block.Attributes["providers"] = providersMap
		blocks = append(blocks, block)
	}

	return blocks
}
