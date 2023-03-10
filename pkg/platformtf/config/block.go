package config

import (
	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/seal-io/seal/utils/maps"
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
	}

	Blocks []*Block
)

// ToNestedMap returns a nested map with the block type as the root key
// and labels are nested in the map with the attributes.
func (b *Block) ToNestedMap() map[string]interface{} {
	m := b.Attributes

	for i := len(b.Labels) - 1; i >= 0; i-- {
		m = map[string]interface{}{
			b.Labels[i]: m,
		}
	}

	m = map[string]interface{}{
		b.Type: m,
	}

	return m
}

// Print returns the block as a config string.
// mapObjects is a map of objects that have been printed already.
func (b *Block) Print(format string, mapObjects map[string]struct{}) ([]byte, error) {
	nestedMap := maps.RemoveNullsCopy(b.ToNestedMap())
	outputBytes, err := terraformutils.Print(nestedMap, mapObjects, format)
	if err != nil {
		return nil, err
	}

	return outputBytes, nil
}
