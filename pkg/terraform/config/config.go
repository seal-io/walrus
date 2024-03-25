package config

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/pkg/templates/translator"
	"github.com/seal-io/walrus/pkg/terraform/block"
	"github.com/seal-io/walrus/pkg/terraform/convertor"
	"github.com/seal-io/walrus/utils/json"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/strs"
)

// Config handles the configuration of resource to terraform config.
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

// loadVariableBlocks returns config variables to get terraform variable config block.
func loadVariableBlocks(opts *VariableOptions) block.Blocks {
	var (
		logger               = log.WithName("terraform").WithName("config")
		blocks               = make(block.Blocks, 0, len(opts.Variables)+len(opts.DependencyOutputs))
		encryptVariableNames = sets.NewString()
	)

	// Secret variables.
	for name, sensitive := range opts.Variables {
		if sensitive {
			encryptVariableNames.Insert(name)
		}
	}

	// Dependency variables.
	for k, o := range opts.DependencyOutputs {
		if o.Sensitive {
			encryptVariableNames.Insert(k)
		}
	}

	if encryptVariableNames.Len() == 0 {
		return blocks
	}

	reg, err := matchAnyRegex(encryptVariableNames.UnsortedList())
	if err != nil {
		logger.Errorf("match any regex failed, %s", err.Error())
		return nil
	}

	for k, v := range opts.Attributes {
		b, err := json.Marshal(v)
		if err != nil {
			logger.Errorf("marshal failed, %s", err.Error())
			return nil
		}

		matches := reg.FindAllString(string(b), -1)
		if len(matches) != 0 {
			blocks = append(blocks, &block.Block{
				Type:   block.TypeVariable,
				Labels: []string{k},
				Attributes: map[string]any{
					"sensitive": true,
				},
			})
		}
	}

	return blocks
}

// loadOutputBlocks returns terraform outputs config block.
func loadOutputBlocks(opts OutputOptions) block.Blocks {
	// Template output.
	blocks := make(block.Blocks, 0, len(opts))

	for _, o := range opts {
		blocks = append(blocks, &block.Block{
			Type:   block.TypeOutput,
			Labels: []string{o.Name},
			Attributes: map[string]any{
				"sensitive": o.Sensitive,
			},
		})
	}

	return blocks
}

func CreateConfigToBytes(opts CreateOptions) ([]byte, error) {
	conf, err := NewConfig(opts)
	if err != nil {
		return nil, err
	}

	return conf.Bytes()
}

// typeExprTokens returns the HCL tokens for a type expression.
func typeExprTokens(ty cty.Type) (hclwrite.Tokens, error) {
	switch ty {
	case cty.String:
		return hclwrite.TokensForIdentifier("string"), nil
	case cty.Bool:
		return hclwrite.TokensForIdentifier("bool"), nil
	case cty.Number:
		return hclwrite.TokensForIdentifier("number"), nil
	case cty.DynamicPseudoType:
		return hclwrite.TokensForIdentifier("any"), nil
	}

	if ty.IsCollectionType() {
		etyTokens, err := typeExprTokens(ty.ElementType())
		if err != nil {
			return nil, err
		}

		switch {
		case ty.IsListType():
			return hclwrite.TokensForFunctionCall("list", etyTokens), nil
		case ty.IsSetType():
			return hclwrite.TokensForFunctionCall("set", etyTokens), nil
		case ty.IsMapType():
			return hclwrite.TokensForFunctionCall("map", etyTokens), nil
		default:
			// Should never happen because the above is exhaustive.
			return nil, fmt.Errorf("unsupported collection type: %s", ty.FriendlyName())
		}
	}

	if ty.IsObjectType() {
		atys := ty.AttributeTypes()
		names := make([]string, 0, len(atys))

		for name := range atys {
			names = append(names, name)
		}

		sort.Strings(names)

		items := make([]hclwrite.ObjectAttrTokens, len(names))

		for i, name := range names {
			value, err := typeExprTokens(atys[name])
			if err != nil {
				return nil, err
			}

			items[i] = hclwrite.ObjectAttrTokens{
				Name:  hclwrite.TokensForIdentifier(name),
				Value: value,
			}
		}

		return hclwrite.TokensForObject(items), nil
	}

	if ty.IsTupleType() {
		etys := ty.TupleElementTypes()
		items := make([]hclwrite.Tokens, len(etys))

		for i, ety := range etys {
			value, err := typeExprTokens(ety)
			if err != nil {
				return nil, err
			}

			items[i] = value
		}

		return hclwrite.TokensForTuple(items), nil
	}

	return nil, fmt.Errorf("unsupported type: %s", ty.GoString())
}

func matchAnyRegex(list []string) (*regexp.Regexp, error) {
	var sb strings.Builder

	sb.WriteString("(")

	for i, v := range list {
		sb.WriteString(v)

		if i < len(list)-1 {
			sb.WriteString("|")
		}
	}

	sb.WriteString(")")

	return regexp.Compile(sb.String())
}
