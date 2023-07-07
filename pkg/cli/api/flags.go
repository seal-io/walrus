package api

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/pflag"

	"github.com/seal-io/seal/utils/json"
	"github.com/seal-io/seal/utils/strs"
)

// ObjectFlag creates a custom flag for map[string]interface{}.
type ObjectFlag map[string]interface{}

// String returns a string represent of this flag.
func (i *ObjectFlag) String() string {
	if i != nil {
		b, err := json.Marshal(i)
		if err != nil {
			return ""
		}

		return fmt.Sprintf("%s...", strs.FirstContent(string(b), 20))
	}

	return ""
}

// Set a new value on the flag.
func (i *ObjectFlag) Set(value string) error {
	var val map[string]interface{}

	err := json.Unmarshal([]byte(value), &val)
	if err != nil {
		return err
	}
	*i = val

	return nil
}

// Type returns the type of this custom flag, which will be displayed in `--help` output.
func (i *ObjectFlag) Type() string {
	return "json"
}

// ArrayObjectFlag creates a custom flag for []interface.
type ArrayObjectFlag []interface{}

// String returns a string represent of this flag.
func (i *ArrayObjectFlag) String() string {
	if i != nil {
		b, err := json.Marshal(i)
		if err != nil {
			return ""
		}

		return fmt.Sprintf("%s...", strs.FirstContent(string(b), 20))
	}

	return ""
}

// Set a new value on the flag.
func (i *ArrayObjectFlag) Set(value string) error {
	var val []interface{}

	err := json.Unmarshal([]byte(value), &val)
	if err != nil {
		return err
	}
	*i = val

	return nil
}

// Type returns the type of this custom flag, which will be displayed in `--help` output.
func (i *ArrayObjectFlag) Type() string {
	return "jsonArray"
}

// ObjectIDFlag creates a custom flag for map[string]string{"id": "xxx"}.
type ObjectIDFlag map[string]string

// String returns a string represent of this flag.
func (i ObjectIDFlag) String() string {
	if i == nil {
		return ""
	}

	return i["id"]
}

// Set a new value on the flag.
func (i ObjectIDFlag) Set(value string) error {
	i["id"] = value
	return nil
}

// Type returns the type of this custom flag, which will be displayed in `--help` output.
func (i ObjectIDFlag) Type() string {
	return openapi3.TypeString
}

// AddFlag create flag with name, type, description, default value, and add it to flagSet.
func AddFlag(name, schemaType, description string, value interface{}, flags *pflag.FlagSet) interface{} {
	existed := flags.Lookup(name)
	if existed != nil {
		return nil
	}

	switch schemaType {
	case openapi3.TypeBoolean:
		var def bool
		if value != nil {
			def = value.(bool)
		}

		return flags.Bool(name, def, description)
	case openapi3.TypeInteger:
		def := 0
		if value != nil {
			def = value.(int)
		}

		return flags.Int(name, def, description)
	case openapi3.TypeNumber:
		def := 0.0
		if value != nil {
			def = value.(float64)
		}

		return flags.Float64(name, def, description)
	case openapi3.TypeString:
		def := ""
		if value != nil {
			def = value.(string)
		}

		return flags.String(name, def, description)
	case "map[string]string":
		var def map[string]string
		if value != nil {
			def = value.(map[string]string)
		}

		return flags.StringToString(name, def, description)
	case "map[string]int":
		var def map[string]int
		if value != nil {
			def, _ = value.(map[string]int)
		}

		return flags.StringToInt(name, def, description)
	case "map[string]int64", "map[string]int32":
		var def map[string]int64
		if value != nil {
			def, _ = value.(map[string]int64)
		}

		return flags.StringToInt64(name, def, description)
	case "array[boolean]":
		var def []bool
		if value != nil {
			def = value.([]bool)
		}

		return flags.BoolSlice(name, def, description)
	case "array[integer]":
		var def []int
		if value != nil {
			def = value.([]int)
		}

		return flags.IntSlice(name, def, description)
	case "array[number]":
		var def []float64
		if value != nil {
			def = value.([]float64)
		}

		return flags.Float64Slice(name, def, description)
	case "array[string]":
		var def []string
		if value != nil {
			def = value.([]string)
		}

		return flags.StringSlice(name, def, description)
	case "array[object]":
		ao := &ArrayObjectFlag{}
		flags.Var(ao, name, description)

		return ao
	case "objectID":
		name := fmt.Sprintf("%s-id", name)

		existed := flags.Lookup(name)
		if existed != nil {
			return nil
		}
		oid := ObjectIDFlag{}
		flags.Var(oid, name, description)

		return oid
	default:
		obj := &ObjectFlag{}
		flags.Var(obj, name, description)

		return obj
	}
}
