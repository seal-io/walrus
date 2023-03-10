package config

import (
	"fmt"
	"io"
	"testing"
)

func TestConfig_Print(t *testing.T) {
	blocks := Blocks{
		{
			Type: "resource",
			Labels: []string{
				"aws_instance",
				"test",
			},
			Attributes: map[string]interface{}{
				"ami": "ami-0c55b159cbfafe1f0",
			},
		},
		{
			Type: "provider",
			Labels: []string{
				"aws",
			},
			Attributes: map[string]interface{}{
				"region": "us-east-1",
			},
		},
		{
			Type: "module",
			Labels: []string{
				"mysql",
			},
			Attributes: map[string]interface{}{
				"providers": map[string]interface{}{
					"kubernetes": "$${kubernetes.us-east-1}",
				},
				"source": "github.com/terraform-aws-modules/terraform-aws-rds",
			},
		},
		{
			Type: "data",
			Labels: []string{
				"aws_ami",
				"test",
			},
			Attributes: map[string]interface{}{
				"owners": []string{"amazon"},
			},
		},
	}

	tfBlocks := Blocks{
		{
			Type: "backend",
			Labels: []string{
				"s3",
			},
			Attributes: map[string]interface{}{
				"bucket": "terraform-state",
			},
		},
	}

	testCases := []struct {
		config   *Config
		expected []string
	}{
		{
			config: &Config{
				format:   "hcl",
				Blocks:   blocks,
				TFBlocks: tfBlocks,
			},
			expected: []string{
				`terraform {

  backend "s3" {
    bucket = "terraform-state"
  }
  
}

resource "aws_instance" "test" {
  ami = "ami-0c55b159cbfafe1f0"
}

provider "aws" {
  region = "us-east-1"
}

module "mysql" {
  providers = {
    kubernetes = kubernetes.us-east-1
  }

  source = "github.com/terraform-aws-modules/terraform-aws-rds"
}

data "aws_ami" "test" {
  owners = ["amazon"]
}

`,
				`terraform {
  backend "s3" {
    bucket = "terraform-state"
  }
}

resource "aws_instance" "test" {
  ami = "ami-0c55b159cbfafe1f0"
}

provider "aws" {
  region = "us-east-1"
}

module "mysql" {
  source = "github.com/terraform-aws-modules/terraform-aws-rds"

  providers = {
    kubernetes = kubernetes.us-east-1
  }
}

data "aws_ami" "test" {
  owners = ["amazon"]
}

`,
			},
		},
	}

	for _, tc := range testCases {
		cReader, err := tc.config.Reader()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		actual, err := io.ReadAll(cReader)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		found := false
		for _, expected := range tc.expected {
			if string(actual) == expected {
				found = true
				break
			}
		}

		if !found {
			fmt.Println(string(actual), string(tc.expected[0]))
			t.Errorf("expected %#v, got %#v", tc.expected, string(actual))
		}
	}
}
