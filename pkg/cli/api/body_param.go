package api

import (
	"github.com/spf13/pflag"

	"github.com/seal-io/seal/utils/strs"
)

// BodyParams represent request body and params type.
type BodyParams struct {
	Type   string       `json:"type,omitempty"`
	Params []*BodyParam `json:"params,omitempty"`
}

// BodyParam represents each field in body.
type BodyParam struct {
	Name        string      `json:"name,omitempty"`
	Type        string      `json:"type,omitempty"`
	Description string      `json:"description,omitempty"`
	Default     interface{} `json:"default,omitempty"`
}

// AddFlag adds a new option flag to a command's flag set for this body param.
func (b BodyParam) AddFlag(flags *pflag.FlagSet) interface{} {
	name := b.OptionName()

	existed := flags.Lookup(name)
	if existed != nil {
		return nil
	}

	return AddFlag(name, b.Type, b.Description, b.Default, flags)
}

// OptionName returns the commandline option name for this parameter.
func (b BodyParam) OptionName() string {
	name := b.Name
	return strs.Dasherize(name)
}
