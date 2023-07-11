package api

import (
	"fmt"
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/pflag"

	"github.com/seal-io/seal/utils/strs"
)

// Param represents an API operation input parameter.
type Param struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Description string      `json:"description,omitempty"`
	Style       string      `json:"style,omitempty"`
	Explode     bool        `json:"explode,omitempty"`
	Default     interface{} `json:"default,omitempty"`
}

// AddFlag adds a new option flag to a command's flag set for this parameter.
func (p Param) AddFlag(flags *pflag.FlagSet) interface{} {
	name := p.OptionName()

	existed := flags.Lookup(name)
	if existed != nil {
		return nil
	}

	return AddFlag(name, p.Type, p.Description, p.Default, flags)
}

// OptionName returns the formatted commandline option name for this parameter.
func (p Param) OptionName() string {
	name := p.Name
	return strs.Dasherize(name)
}

// Serialize the parameter based on the type/style/explode.
func (p Param) Serialize(value interface{}) []string {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		value = v.Interface()
	}

	switch p.Type {
	case openapi3.TypeBoolean, openapi3.TypeInteger, openapi3.TypeNumber, openapi3.TypeString:
		switch p.Style {
		case openapi3.SerializationForm:
			return []string{fmt.Sprintf("%s=%v", p.Name, value)}
		case openapi3.SerializationSimple:
			return []string{fmt.Sprintf("%v", value)}
		}

	case ValueTypeArrayBoolean, ValueTypeArrayInt, ValueTypeArrayNumber, ValueTypeArrayString:
		var encoded []string

		switch p.Style {
		case openapi3.SerializationForm:
			for i := 0; i < v.Len(); i++ {
				item := v.Index(i)
				if p.Explode {
					encoded = append(encoded, fmt.Sprintf("%v", item.Interface()))
				} else {
					if len(encoded) == 0 {
						encoded = append(encoded, "")
					}

					encoded[0] += fmt.Sprintf("%v", item.Interface())
					if i < v.Len()-1 {
						encoded[0] += ","
					}
				}
			}
		case openapi3.SerializationSimple:
			encoded = append(encoded, "")

			for i := 0; i < v.Len(); i++ {
				item := v.Index(i)

				encoded[0] += fmt.Sprintf("%v", item.Interface())
				if i < v.Len()-1 {
					encoded[0] += ","
				}
			}
		}

		return encoded
	}

	return nil
}
