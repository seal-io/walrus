package interpolation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterpolation(t *testing.T) {
	t.Setenv("FOO", "foo")
	t.Setenv("FOO2", "foo2")

	testcases := []struct {
		test     string
		expected string
		isErr    bool
	}{
		{test: "${FOO}", expected: "foo"},
		{test: "${FOO-foo}", expected: "foo"},
		{test: "${FOO:-foo_}", expected: "foo"},
		{test: "${BAR:-bar}", expected: "bar"},

		{test: "${FOO2:?error message}", expected: "foo2"},
		{test: "${FOO2?error message}", expected: "foo2"},
		{test: "${BAR:?error message}", isErr: true},
		{test: "${BAR?error message}", isErr: true},

		// Built-in.
		{test: "${var.test}", expected: "${var.test}"},
		{test: "${res.test.output}", expected: "${res.test.output}"},
		{test: "${unsupported.test}", isErr: true},

		// File.
		{test: "${file(testdata/env)}", expected: "foo"},
		{test: `${file("testdata/env")}`, expected: "foo"},
		{test: "${file('testdata/env')}", expected: "foo"},
		{test: "${file('testdata/.env')}", expected: "foo"},
	}

	for _, tc := range testcases {
		yml := map[string]any{
			"case": tc.test,
		}
		result, err := Interpolate(yml, nil, false)
		assert.Equal(t, err != nil, tc.isErr)

		if err == nil {
			assert.Equal(t, tc.expected, result["case"])
		}
	}
}
