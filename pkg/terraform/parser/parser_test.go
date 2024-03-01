package parser

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseInstanceProviderConnector(t *testing.T) {
	testCases := []struct {
		input          string
		expectedOutput string
		expectedError  bool
	}{
		{
			input:          "provider.connector--instance",
			expectedOutput: "",
			expectedError:  true,
		},
		{
			input:          "provider.connector",
			expectedOutput: "",
			expectedError:  true,
		},
		{
			input:          "invalid format",
			expectedOutput: "",
			expectedError:  true,
		},
		{
			input:          "provider[\"registry.terraform.io/hashicorp/kubernetes\"].connector--kubernetes",
			expectedOutput: "kubernetes",
			expectedError:  false,
		},
	}

	for _, tc := range testCases {
		actualOutput, actualError := ParseInstanceProviderConnector(tc.input)
		assert.Equal(t, tc.expectedOutput, actualOutput)

		if tc.expectedError {
			assert.Error(t, actualError)
		} else {
			assert.NoError(t, actualError)
		}
	}
}

func TestParseInstanceID(t *testing.T) {
	testCases := []struct {
		rs             resourceState
		input          instanceObjectState
		expectedOutput string
		expectedError  bool
	}{
		{
			rs:             resourceState{Module: "a", Name: "test", Type: "test"},
			input:          instanceObjectState{Attributes: []byte(`{"id":"123"}`)},
			expectedOutput: "123",
			expectedError:  false,
		},
		{
			rs:             resourceState{Name: "test", Type: "test"},
			input:          instanceObjectState{AttributesFlat: map[string]string{"id": "123"}},
			expectedOutput: "123",
			expectedError:  false,
		},
		{
			rs:             resourceState{Module: "a", Name: "foo", Type: "helm_release"},
			input:          instanceObjectState{IndexKey: 0},
			expectedOutput: "a.helm_release.foo[0]",
			expectedError:  false,
		},
		{
			rs: resourceState{Module: "a", Name: "bar", Type: "helm_release"},
			input: instanceObjectState{
				Attributes: []byte(`{"id":"123"}`),
				IndexKey:   "foo",
			},
			expectedOutput: "123",
			expectedError:  false,
		},
		{
			rs:             resourceState{Module: "b", Name: "bar", Type: "helm_release"},
			input:          instanceObjectState{IndexKey: "foo"},
			expectedOutput: "b.helm_release.bar[\"foo\"]",
			expectedError:  false,
		},
		{
			rs:             resourceState{},
			input:          instanceObjectState{},
			expectedOutput: "",
			expectedError:  true,
		},
	}

	for _, tc := range testCases {
		actualOutput, actualError := ParseInstanceID(tc.rs, tc.input)
		assert.Equal(t, tc.expectedOutput, actualOutput)

		if tc.expectedError {
			assert.Error(t, actualError)
		} else {
			assert.NoError(t, actualError)
		}
	}
}

func TestParseProviderString(t *testing.T) {
	testCases := []struct {
		input          string
		expectedOutput *AbsProviderConfig
		expectedError  error
	}{
		{
			input: `provider["registry.terraform.io/hashicorp/kubernetes"].connector--kubernetes`,
			expectedOutput: &AbsProviderConfig{
				Provider: Provider{
					Type:      "kubernetes",
					Namespace: "hashicorp",
					Hostname:  "registry.terraform.io",
				},
				Alias: "connector--kubernetes",
			},
		},
		{
			input: `provider["registry.terraform.io/hashicorp/kubernetes"].connector`,
			expectedOutput: &AbsProviderConfig{
				Provider: Provider{
					Type:      "kubernetes",
					Namespace: "hashicorp",
					Hostname:  "registry.terraform.io",
				},
				Alias: "connector",
			},
		},
		{
			input: `provider["registry.terraform.io/hashicorp/kubernetes"]`,
			expectedOutput: &AbsProviderConfig{
				Provider: Provider{
					Type:      "kubernetes",
					Namespace: "hashicorp",
					Hostname:  "registry.terraform.io",
				},
				Alias: "",
			},
		},
		{
			input: `module.baz.provider["registry.terraform.io/hashicorp/aws"].foo`,
			expectedOutput: &AbsProviderConfig{
				Provider: Provider{
					Type:      "aws",
					Namespace: "hashicorp",
					Hostname:  "registry.terraform.io",
				},
				Alias: "foo",
			},
		},
		{
			input: `module.baz["foo"].provider["registry.terraform.io/hashicorp/aws"]`,
			expectedError: errors.New(
				"invalid provider configuration address \"module.baz[\"foo\"].provider[\"registry.terraform.io/hashicorp/aws\"]",
			),
		},
		{
			input: `module.baz[1].provider["registry.terraform.io/hashicorp/aws"]`,
			expectedError: errors.New(
				"invalid provider configuration address \"module.baz[1].provider[\"registry.terraform.io/hashicorp/aws\"]",
			),
		},
		{
			input: `module.baz[1].module.bar.provider["registry.terraform.io/hashicorp/aws"]`,
			expectedError: errors.New(
				"invalid provider configuration address \"module.baz[1].module.bar.provider[\"registry.terraform.io/hashicorp/aws\"]",
			),
		},
		{
			input: `aws`,
			expectedError: errors.New(
				"provider address must begin with \"provider.\", followed by a provider type name",
			),
		},
		{
			input:         `provider.`,
			expectedError: errors.New("the prefix \"provider.\" must be followed by a provider type name"),
		},
		{
			input:         `aws.foo`,
			expectedError: errors.New("the prefix \"provider.\" must be followed by a provider type name"),
		},
		{
			input:         `provider.aws.foo.bar`,
			expectedError: errors.New("provider type name must be followed by a configuration alias name"),
		},
		{
			input:         `provider["aws"]["foo"]`,
			expectedError: errors.New("provider type name must be followed by a configuration alias name"),
		},
		{
			input:         `provider[0]`,
			expectedError: errors.New("provider type name must be followed by a configuration alias name"),
		},
	}

	for _, tc := range testCases {
		actualOutput, actualError := ParseAbsProviderString(tc.input)
		if tc.expectedError != nil {
			assert.Error(t, tc.expectedError, actualError)
		} else {
			assert.NoError(t, actualError)
			assert.Equal(t, tc.expectedOutput, actualOutput)
		}
	}
}
