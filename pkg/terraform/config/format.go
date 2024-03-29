package config

import (
	"regexp"
)

var _transformations = []struct {
	match     *regexp.Regexp
	replace   []byte
	replaceFn func([]byte) []byte
}{
	{
		// Used for interpolation
		// replace all the `"key" = "$${a.b.c}` for `"key" = a.b.c`.
		match:   regexp.MustCompile(`"\$\${([^$}{]+)\.([^$}{]+)\.([^$}{]+)}"`),
		replace: []byte(`$1.$2.$3`),
	},
	{
		// Used for variables
		// replace all the `"key" = "$${a.b}` for `"key" = a.b`.
		match:   regexp.MustCompile(`"\$\${([^$}{]+)\.([^$}{]+)}"`),
		replace: []byte(`$1.$2`),
	},
	{
		// Replace `$${xxx}` with `xxx`.
		match:   regexp.MustCompile(`\$\${([^}]+)}`),
		replace: []byte(`${$1}`),
	},
	{
		// Replace "{{xxx}}" with xxx. Quote will be removed.
		match:   regexp.MustCompile(`"{{([^}]+)}}"`),
		replace: []byte(`$1`),
	},
	{},
}

// Format formats the hcl with the transformations.
func Format(hcl []byte) []byte {
	for i := 0; i < len(_transformations); i++ {
		m := _transformations[i]
		if m.replace != nil {
			hcl = m.match.ReplaceAll(hcl, m.replace)
		} else if m.replaceFn != nil {
			hcl = m.match.ReplaceAllFunc(hcl, m.replaceFn)
		}
	}

	return hcl
}
