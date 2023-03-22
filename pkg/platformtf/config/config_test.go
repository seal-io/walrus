package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateConfigToBytes(t *testing.T) {
	testCases := []struct {
		name     string
		option   CreateOptions
		expected []byte
	}{
		{
			name:     "test create config to bytes with empty option",
			option:   CreateOptions{},
			expected: nil,
		},
		{
			name: "test create config to bytes with attributes",
			option: CreateOptions{
				Attributes: map[string]interface{}{
					"var1":    "ami-0c55b159cbfafe1f0",
					"secret1": "password",
				},
			},
			expected: []byte(`secret1 = "password"
var1    = "ami-0c55b159cbfafe1f0"
`),
		},
		{
			name: "test create config to bytes with attributes and child blocks",
			option: CreateOptions{
				TerraformOptions: &TerraformOptions{
					Token:         "token",
					Address:       "https://localhost:8080",
					SkipTLSVerify: true,
				},
				VariableOptions: &VariableOptions{
					Variables: map[string]interface{}{
						"var1": "value1",
						"var2": "value2",
					},
				},
			},
			expected: []byte(`terraform {
  backend "http" {
    address                = "https://localhost:8080"
    password               = "token"
    skip_cert_verification = true
    update_method          = "PUT"
    username               = "seal"
  }
}

variable "var1" {
  type = string
}

variable "var2" {
  type = string
}

`),
		},
	}

	for _, tt := range testCases {
		got, err := CreateConfigToBytes(tt.option)
		if err != nil {
			assert.Errorf(t, err, "unexpected error: %v", err)
		}
		if !assert.Equal(t, string(tt.expected), string(got)) {
			assert.Errorf(t, err, "name: %s", tt.name)
		}
	}
}
