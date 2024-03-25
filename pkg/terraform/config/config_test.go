package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
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
				Attributes: map[string]any{
					"var1":    "ami-0c55b159cbfafe1f0",
					"secret1": "password",
				},
			},
			expected: []byte(`secret1 = "password"
var1    = "ami-0c55b159cbfafe1f0"
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

func TestTypeExprTokens(t *testing.T) {
	cases := []struct {
		name     string
		ctyType  cty.Type
		expected string
		wantErr  bool
	}{
		{
			name:     "test get cty type with string",
			ctyType:  cty.String,
			expected: "string",
			wantErr:  false,
		},
		{
			name:     "test get cty type with number",
			ctyType:  cty.Number,
			expected: "number",
			wantErr:  false,
		},
		{
			name:     "test get cty type with bool",
			ctyType:  cty.Bool,
			expected: "bool",
			wantErr:  false,
		},
		{
			name:     "test get cty type with list",
			ctyType:  cty.List(cty.String),
			expected: "list(string)",
			wantErr:  false,
		},
		{
			name:     "test get cty type with map",
			ctyType:  cty.Map(cty.String),
			expected: "map(string)",
			wantErr:  false,
		},
		{
			name:     "test get cty type with object",
			ctyType:  cty.Object(map[string]cty.Type{"test": cty.String}),
			expected: "{\n  test = string\n}",
			wantErr:  false,
		},
		{
			name:     "test get cty type with tuple",
			ctyType:  cty.Tuple([]cty.Type{cty.String}),
			expected: "[string]",
			wantErr:  false,
		},
		{
			name:     "test get cty type with set",
			ctyType:  cty.Set(cty.String),
			expected: "set(string)",
			wantErr:  false,
		},
		{
			name:     "test get cty type with dynamic",
			ctyType:  cty.DynamicPseudoType,
			expected: "any",
			wantErr:  false,
		},
		{
			name:    "test get cty type with nil",
			ctyType: cty.NilType,
			wantErr: true,
		},
		{
			name:    "test error with nil set",
			ctyType: cty.Set(cty.NilType),
			wantErr: true,
		},
		{
			name:    "test error with nil object",
			ctyType: cty.Object(map[string]cty.Type{"test": cty.NilType}),
			wantErr: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := typeExprTokens(tt.ctyType)
			if err != nil && !tt.wantErr {
				t.Errorf("unexpected error: %s", err)
			}

			if !assert.Equal(t, tt.expected, string(got.Bytes())) {
				t.Errorf("expected %s, got %s", tt.expected, string(got.Bytes()))
			}
		})
	}
}
