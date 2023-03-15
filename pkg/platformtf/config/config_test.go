package config

import (
	"io"
	"sync"
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
				"env": map[string]interface{}{
					"test": "test",
				},
				"null": nil,
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
					"kubernetes": "${kubernetes.us-east-1}",
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
		expected string
	}{
		{
			config: &Config{
				Blocks:   blocks,
				TFBlocks: tfBlocks,
				once:     &sync.Once{},
			},
			expected: `terraform {
  backend "s3" {
    bucket = "terraform-state"
  }
}

resource "aws_instance" "test" {
  ami = "ami-0c55b159cbfafe1f0"
  env = {
    test = "test"
  }
  null = null
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

		if string(actual) != tc.expected {
			t.Errorf("expected %#v, got %#v", tc.expected, string(actual))
		}
	}
}
