package config

import (
	"regexp"
)

var (
	_transformations = []struct {
		match     *regexp.Regexp
		replace   []byte
		replaceFn func([]byte) []byte
	}{
		{
			// Used for interpolation
			// Replace all the `"key" = "$${a.b.c}` for `"key" = a.b.c`.
			match:   regexp.MustCompile(`"\$\${([^$}{]+)\.([^$}{]+)\.([^$}{]+)}"`),
			replace: []byte(`$1.$2.$3`),
		},
		{
			// Used for variables
			// Replace all the `"key" = "$${a.b}` for `"key" = a.b`.
			match:   regexp.MustCompile(`"\$\${([^$}{]+)\.([^$}{]+)}"`),
			replace: []byte(`$1.$2`),
		},
		{
			// replace "providers {" to "providers = {".
			match:   regexp.MustCompile(`providers\s*{`),
			replace: []byte("providers = {"),
		},
	}
)

// Format formats the hcl with the transformations.
func Format(hcl []byte) []byte {
	for _, m := range _transformations {
		if m.replace != nil {
			hcl = m.match.ReplaceAll(hcl, m.replace)
		} else if m.replaceFn != nil {
			hcl = m.match.ReplaceAllFunc(hcl, m.replaceFn)
		}
	}

	return hcl
}
