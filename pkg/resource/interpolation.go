package resource

import "regexp"

var (
	// VariableReg is the regexp to match the variable reference in attributes.
	VariableReg = regexp.MustCompile(`\${var\.([a-zA-Z0-9_-]+)}`)
	// ResourceReg is the regexp to match the resource output reference in attributes.
	ResourceReg = regexp.MustCompile(`\${res\.([^.}]+)\.([^.}]+)}`)
)
