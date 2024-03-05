package api

import (
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/pflag"

	"github.com/seal-io/walrus/utils/strs"
)

// BodyParams represent request body and params type.
type BodyParams struct {
	Type   string       `json:"type,omitempty"`
	Params []*BodyParam `json:"params,omitempty"`
}

// BodyParam represents each field in body.
type BodyParam struct {
	Name        string `json:"name,omitempty"`
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
	Default     any    `json:"default,omitempty"`
}

// AddFlag adds a new option flag to a command's flag set for this body param.
func (b BodyParam) AddFlag(flags *pflag.FlagSet) any {
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

// Serialize the parameter based on the type and remove empty value.
func (b BodyParam) Serialize(value any, flagSet *pflag.FlagSet) any {
	if value == nil {
		return nil
	}

	switch b.Type {
	case openapi3.TypeBoolean, openapi3.TypeInteger, openapi3.TypeNumber, openapi3.TypeString:
	case ValueTypeObjectID:
	case ValueTypeArrayBoolean, ValueTypeArrayInt, ValueTypeArrayNumber, ValueTypeArrayString, ValueTypeArrayObject:
		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		if (v.Kind() == reflect.Array || v.Kind() == reflect.Slice) &&
			v.Len() == 0 &&
			!flagSet.Changed(b.OptionName()) {
			return nil
		}
	case ValueTypeMapStringInt64, ValueTypeMapStringInt32, ValueTypeMapStringInt, ValueTypeMapStringString:
		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		if v.Kind() == reflect.Map &&
			v.Len() == 0 &&
			!flagSet.Changed(b.OptionName()) {
			return nil
		}

	default:
		v, ok := value.(*ObjectFlag)
		if ok &&
			(v == nil || reflect.DeepEqual(*v, ObjectFlag{})) &&
			!flagSet.Changed(b.OptionName()) {
			return nil
		}
	}

	return value
}
