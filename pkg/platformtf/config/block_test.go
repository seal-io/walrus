package config

import (
	"testing"
)

func TestBlockEncodeToBytes(t *testing.T) {
	testCases := []struct {
		block    *Block
		expected string
	}{
		{
			block: &Block{
				Type:   "resource",
				Labels: []string{"aws_instance", "test"},
				Attributes: map[string]interface{}{
					"ami":           "ami-0c55b159cbfafe1f0",
					"instance_type": "t2.micro",
				},
			},
			expected: `resource "aws_instance" "test" {
  ami           = "ami-0c55b159cbfafe1f0"
  instance_type = "t2.micro"
}
`,
		},
		{
			block: &Block{
				Type:   "provider",
				Labels: []string{"aws"},
				Attributes: map[string]interface{}{
					"region": "us-east-1",
				},
			},
			expected: `provider "aws" {
  region = "us-east-1"
}
`,
		},
		{
			block: &Block{
				Type:   "module",
				Labels: []string{"mysql"},
				Attributes: map[string]interface{}{
					"source":   "github.com/terraform-aws-modules/terraform-aws-rds",
					"username": "test",
					"password": "test",
				},
			},
			expected: `module "mysql" {
  password = "test"
  source   = "github.com/terraform-aws-modules/terraform-aws-rds"
  username = "test"
}
`,
		},
		{
			block: &Block{
				Type:   "data",
				Labels: []string{"aws_ami", "test"},
				Attributes: map[string]interface{}{
					"owners": []string{"amazon"},
					"filter": map[string]interface{}{
						"name":   "name",
						"values": []string{"amzn-ami-hvm-2018.03.0.20180409-x86_64-gp2"},
					},
				},
			},
			expected: `data "aws_ami" "test" {
  filter = {
    name   = "name"
    values = ["amzn-ami-hvm-2018.03.0.20180409-x86_64-gp2"]
  }
  owners = ["amazon"]
}
`,
		},

		{
			block: &Block{
				Type:   "module",
				Labels: []string{"mysql"},
				Attributes: map[string]interface{}{
					"source": "github.com/terraform-aws-modules/terraform-aws-rds",
					"env": map[string]interface{}{
						"test": map[string]interface{}{
							"identifier": "test",
						},
					},
				},
			},
			expected: `module "mysql" {
  env = {
    test = {
      identifier = "test"
    }
  }
  source = "github.com/terraform-aws-modules/terraform-aws-rds"
}
`,
		},
		{
			block: &Block{
				Type:   "provider",
				Labels: []string{"helm"},
				childBlocks: Blocks{
					&Block{
						Type: "kubernetes",
						Attributes: map[string]interface{}{
							"config_path": "~/.kube/config",
						},
					},
				},
			},
			expected: `provider "helm" {
  kubernetes {
    config_path = "~/.kube/config"
  }
}
`,
		},
	}

	for _, tc := range testCases {
		actual, err := tc.block.EncodeToBytes()
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if string(actual) != tc.expected {
			t.Errorf("expected %#v, got %#v", tc.expected, string(actual))
		}
	}
}
