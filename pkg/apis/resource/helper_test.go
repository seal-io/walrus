package resource

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/seal-io/walrus/pkg/dao/types/property"
	"github.com/seal-io/walrus/utils/json"
)

func TestInjectAttributes(t *testing.T) {
	cases := []struct {
		name         string
		attrByte     []byte
		varValues    map[string]json.RawMessage
		outputValues map[string]property.Value
		expected     property.Values
		err          bool
	}{
		{
			name:     "no interpolation",
			attrByte: []byte(`{"a": "b"}`),
			expected: property.Values{
				"a": json.RawMessage(`"b"`),
			},
		},
		{
			name:     "interpolation with variable",
			attrByte: []byte(`{"a": "${var.b}"}`),
			varValues: map[string]json.RawMessage{
				"${var.b}": json.RawMessage(`c`),
			},
			expected: property.Values{
				"a": json.RawMessage(`"c"`),
			},
		},
		{
			name:     "interpolation with output",
			attrByte: []byte(`{"a": "${res.b.c}"}`),
			outputValues: map[string]property.Value{
				"${res.b.c}": property.Value(`d`),
			},
			expected: property.Values{
				"a": json.RawMessage(`"d"`),
			},
		},
		{
			name:     "interpolation with variable have newline",
			attrByte: []byte(`{"a": "${var.b}"}`),
			varValues: map[string]json.RawMessage{
				"${var.b}": json.RawMessage(`-----BEGIN RSA PRIVATE KEY-----
xxx
-----END RSA PRIVATE KEY-----`),
			},
			expected: property.Values{
				"a": json.RawMessage(`"-----BEGIN RSA PRIVATE KEY-----\nxxx\n-----END RSA PRIVATE KEY-----"`),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual, err := replaceAttributes(c.attrByte, c.varValues, c.outputValues)
			if (err != nil) != c.err {
				t.Errorf("expected error: %v, got: %v", c.err, err)
			}
			if err == nil {
				if !assert.Equal(t, c.expected, actual) {
					t.Errorf("expected: %v, got: %v", c.expected, actual)
				}
			}
		})
	}
}
