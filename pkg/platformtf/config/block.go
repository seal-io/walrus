package config

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

type (
	// Block represents to terraform config block like terraform block, provider block, module block, etc.
	// Block.Print() will generate the config string of the block.
	// e.g. provider block:
	/**
	  block1 = &Block{
	  	Type: "provider",
	      Labels: []string{"aws"},
	  	Attributes: map[string]interface{}{
	  		"region": "us-east-1",
	  	},
	  },
	  block1.Print("hcl") will generate the string:
	  provider "aws" {
	    region = "us-east-1"
	  }


	  block2 = &Block{
	  	Type: "data",
	  	Labels: []string{"lable1", "label2"},
	  	Attributes: map[string]interface{}{
	  		"test": "test"
	  	},
	  }
	  block2.Print("hcl") will generate the string:
	  data "lable1" "label2" {
	  	  test = "test"
	  }

	*/
	Block struct {
		// Type the type of the block, e.g. provider, module, resource, etc.
		Type string
		// Label the label of the block, e.g. aws, aws_instance, etc.
		Labels []string

		// Attributes the Attributes of the block.
		Attributes map[string]interface{}

		hclBlock *hclwrite.Block

		// childBlocks holds information about any child Blocks that the Block may have. This can be empty.
		childBlocks Blocks
	}

	Blocks []*Block
)

// EncodeToBytes returns the block as a config bytes.
func (b *Block) EncodeToBytes() ([]byte, error) {
	hclBlock, err := b.ToHCLBlock()
	if err != nil {
		return nil, err
	}

	f := hclwrite.NewEmptyFile()
	body := f.Body()
	body.AppendBlock(hclBlock)

	return f.Bytes(), nil
}

// removeRing removes the ring in the block tree.
func (b *Block) removeRing() {
	stack := make([]*Block, 0)

	removeRing(b, stack)
}

func (b *Block) AppendBlock(block *Block) {
	if block == nil {
		return
	}

	b.childBlocks = append(b.childBlocks, block)
	// remove ring
	b.removeRing()
	// append block will cause the tree structure change,
	// so we need to reset the hclBlock.
	b.hclBlock = nil
}

// ToHCLBlock returns the block as a hclwrite.Block.
func (b *Block) ToHCLBlock() (*hclwrite.Block, error) {
	if b.Type == "" {
		return nil, fmt.Errorf("block type is empty")
	}

	block := hclwrite.NewBlock(b.Type, b.Labels)
	// append child blocks.
	for i := 0; i < len(b.childBlocks); i++ {
		cb, err := b.childBlocks[i].ToHCLBlock()
		if err != nil {
			return nil, err
		}
		block.Body().AppendBlock(cb)
		if i != len(b.childBlocks)-1 {
			block.Body().AppendNewline()
		}
	}
	attributes, err := convertToCtyWithJson(b.Attributes)
	if err != nil {
		return nil, err
	}
	attrKeys := sortValueKeys(attributes)
	if len(attrKeys) == 0 {
		return block, nil
	}
	attrMap := attributes.AsValueMap()
	for _, attr := range attrKeys {
		block.Body().SetAttributeValue(attr, attrMap[attr])
	}
	b.hclBlock = block

	return b.hclBlock, nil
}

// sortValueKeys will return a sorted list of the keys the val has.
func sortValueKeys(val cty.Value) []string {
	if !val.CanIterateElements() {
		return nil
	}
	keys := make([]string, 0, val.LengthInt())
	for it := val.ElementIterator(); it.Next(); {
		k, _ := it.Element()
		keys = append(keys, k.AsString())
	}
	sort.Strings(keys)

	return keys
}

// convertToCtyWithJson Converts arbitrary go types that are json serializable to a cty Value
// by using json as an intermediary representation.
func convertToCtyWithJson(val interface{}) (cty.Value, error) {
	jsonBytes, err := json.Marshal(val)
	if err != nil {
		return cty.NilVal, fmt.Errorf("failed to marshal value to json: %w", err)
	}
	var ctyJsonVal ctyjson.SimpleJSONValue
	if err := ctyJsonVal.UnmarshalJSON(jsonBytes); err != nil {
		return cty.NilVal, fmt.Errorf("failed to unmarshal json to cty value: %w", err)
	}
	return ctyJsonVal.Value, nil
}

// removeRing will remove the ring in the block tree.
func removeRing(root *Block, stack []*Block) {
	for i, child := range root.childBlocks {
		for _, node := range stack {
			if child == node {
				root.childBlocks = append(root.childBlocks[:i], root.childBlocks[i+1:]...)
				return
			}
		}
		stack = append(stack, child)
		removeRing(child, stack)
		stack = stack[:len(stack)-1]
	}
}

// AppendBlocks will append the blocks to the tagert.
func AppendBlocks(tagert Blocks, appendBlocks ...Blocks) Blocks {
	for _, b := range appendBlocks {
		tagert = append(tagert, b...)
	}

	return tagert
}
