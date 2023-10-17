package property

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func TestProperty_Cty(t *testing.T) {
	type output struct {
		t   cty.Type
		v   cty.Value
		err error
	}

	testCases := []struct {
		property *Property
		expected output
	}{
		{
			property: &Property{
				Type:  cty.NilType,
				Value: nil,
			},
			expected: output{
				t: cty.NilType,
				v: cty.NilVal,
			},
		},
		{
			property: &Property{
				Type:  cty.NilType,
				Value: []byte(`"test"`),
			},
			expected: output{
				t: cty.String,
				v: cty.StringVal("test"),
			},
		},
		{
			property: &Property{
				Type:  cty.DynamicPseudoType,
				Value: []byte(`true`),
			},
			expected: output{
				t: cty.Bool,
				v: cty.True,
			},
		},
		{
			property: &Property{
				Type:  cty.Bool,
				Value: []byte(`true`),
			},
			expected: output{
				t: cty.Bool,
				v: cty.True,
			},
		},
		{
			property: &Property{
				Type:  cty.List(cty.Bool),
				Value: []byte(`["true", "false"]`),
			},
			expected: output{
				t: cty.List(cty.Bool),
				v: cty.ListVal([]cty.Value{cty.True, cty.False}),
			},
		},
		{
			property: &Property{
				Type:  cty.List(cty.Bool),
				Value: []byte(`[true, false]`),
			},
			expected: output{
				t: cty.List(cty.Bool),
				v: cty.ListVal([]cty.Value{cty.True, cty.False}),
			},
		},
		{
			property: &Property{
				Type:  cty.List(cty.String),
				Value: []byte(`["true", "false"]`),
			},
			expected: output{
				t: cty.List(cty.String),
				v: cty.ListVal([]cty.Value{cty.StringVal("true"), cty.StringVal("false")}),
			},
		},
		{
			property: &Property{
				Type:  cty.List(cty.Map(cty.String)),
				Value: []byte(`[{"hello":"world"}, {"test":"test"}]`),
			},
			expected: output{
				t: cty.List(cty.Map(cty.String)),
				v: cty.ListVal([]cty.Value{
					cty.MapVal(map[string]cty.Value{"hello": cty.StringVal("world")}),
					cty.MapVal(map[string]cty.Value{"test": cty.StringVal("test")}),
				}),
			},
		},
		{
			property: &Property{
				Type:  cty.Object(map[string]cty.Type{"greet": cty.DynamicPseudoType, "true": cty.Bool}),
				Value: []byte(`{"greet":"hello","true":true}`),
			},
			expected: output{
				t: cty.Object(map[string]cty.Type{"greet": cty.DynamicPseudoType, "true": cty.Bool}),
				v: cty.ObjectVal(map[string]cty.Value{"true": cty.True, "greet": cty.StringVal("hello")}),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.property.Type.GoString(), func(t *testing.T) {
			var actual output
			actual.t, actual.v, actual.err = tc.property.Cty()
			assert.Equal(t, tc.expected, actual)
		})
	}
}
