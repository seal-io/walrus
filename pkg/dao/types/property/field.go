package property

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/hcl/v2/ext/typeexpr"
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/seal/utils/json"
)

type (
	// Type indicates the type of property,
	// supported types described as below.
	//
	//	|          Definition           |           JSON Output         |
	//	| ----------------------------- | ----------------------------- |
	//	| any                           | "dynamic"                     |
	//	|                               |                               |
	//	| string                        | "string"                      |
	//	|                               |                               |
	//	| number                        | "number"                      |
	//	|                               |                               |
	//	| bool                          | "bool"                        |
	//	|                               |                               |
	//	| map(string)                   | [                             |
	//	|                               |   "map",                      |
	//	|                               |   "string"                    |
	//	|                               | ]                             |
	//	|                               |                               |
	//	| list(string)                  | [                             |
	//	|                               |   "list",                     |
	//	|                               |   "string"                    |
	//	|                               | ]                             |
	//	|                               |                               |
	//	| list(object({                 | [                             |
	//	|   a = string                  |   "list",                     |
	//	|   b = number                  |   [                           |
	//	|   c = bool                    |     "object",                 |
	//	| }))                           |     {                         |
	//	|                               |       "a":"string",           |
	//	|                               |       "b":"number",           |
	//	|                               |       "c":"bool"              |
	//	|                               |     }                         |
	//	|                               |   ]                           |
	//	|                               | ]                             |
	//	|                               |                               |
	//	| object({                      | [                             |
	//	|   a = string                  |   "object",                   |
	//	|   b = number                  |   {                           |
	//	|   c = bool                    |     "a":"string",             |
	//	| })                            |     "b":"number",             |
	//	|                               |     "c":"bool"                |
	//	|                               |   }                           |
	//	|                               | ]                             |
	//	|                               |                               |
	//	| object({                      | [                             |
	//	|   a = string                  |  "object",                    |
	//	|   b = list(object({           |  {                            |
	//	|     c = bool                  |    "a":"string",              |
	//	|   }))                         |    "b":[                      |
	//	| })                            |         "list",               |
	//	|                               |         [                     |
	//	|                               |          "object",            |
	//	|                               |          {                    |
	//	|                               |            "c":"bool"         |
	//	|                               |          }                    |
	//	|                               |         ]                     |
	//	|                               |        ]                      |
	//	|                               |  }                            |
	//	|                               | ]                             |
	//	|                               |                               |
	//	| tuple([string, bool, number]) | [                             |
	//	|                               |   "tuple",                    |
	//	|                               |   [                           |
	//	|                               |     "string",                 |
	//	|                               |     "bool",                   |
	//	|                               |     "number"                  |
	//	|                               |   ]                           |
	//	|                               | ]                             |
	Type = cty.Type

	// Value indicates the value of property.
	Value = json.RawMessage
)

type (
	// IProperty holds the functions of the property.
	IProperty interface {
		// GetType returns the Type of the property.
		GetType() Type
		// GetValue returns the Value of the property.
		GetValue() Value
		// GetNumber returns the underlay value as a number,
		// if not found, returns false.
		GetNumber() (any, bool, error)
		// GetUint64 returns the underlay value as uint64,
		// if not found, returns false.
		GetUint64() (uint64, bool, error)
		// GetUint32 returns the underlay value as uint32,
		// if not found, returns false.
		GetUint32() (uint32, bool, error)
		// GetUint16 returns the underlay value as uint16,
		// if not found, returns false.
		GetUint16() (uint16, bool, error)
		// GetUint8 returns the underlay value as uint8,
		// if not found, returns false.
		GetUint8() (uint8, bool, error)
		// GetUint returns the underlay value as uint,
		// if not found, returns false.
		GetUint() (uint, bool, error)
		// GetInt64 returns the underlay value as int64,
		// if not found, returns false.
		GetInt64() (int64, bool, error)
		// GetInt32 returns the underlay value as int32,
		// if not found, returns false.
		GetInt32() (int32, bool, error)
		// GetInt16 returns the underlay value as int16,
		// if not found, returns false.
		GetInt16() (int16, bool, error)
		// GetInt8 returns the underlay value as int8,
		// if not found, returns false.
		GetInt8() (int8, bool, error)
		// GetInt returns the underlay value as int,
		// if not found, returns false.
		GetInt() (int, bool, error)
		// GetFloat64 returns the underlay value as float64,
		// if not found, returns false.
		GetFloat64() (float64, bool, error)
		// GetFloat32 returns the underlay value as float32,
		// if not found, returns false.
		GetFloat32() (float32, bool, error)
		// GetDuration returns the underlay value as time.Duration,
		// if not found, returns false.
		GetDuration() (time.Duration, bool, error)
		// GetBool returns the underlay value as bool,
		// if not found, returns false.
		GetBool() (bool, bool, error)
		// GetString returns the underlay value as string,
		// if not found, returns false.
		GetString() (string, bool, error)
		// Cty returns the cty.Type and cty.Value of this value.
		Cty() (cty.Type, cty.Value, error)
	}

	// Property holds the type and underlay value of the property.
	Property struct {
		// Type specifies the type of this property.
		Type Type `json:"type"`
		// Value specifies the value of this property.
		Value Value `json:"value,omitempty"`
	}
)

func (i Property) GetType() Type {
	return i.Type
}

func (i Property) GetValue() Value {
	return i.Value
}

func (i Property) GetNumber() (any, bool, error) {
	var iv, ok, _ = i.GetInt64()
	if ok {
		return iv, ok, nil
	}
	uiv, ok, _ := i.GetUint64()
	if ok {
		return uiv, ok, nil
	}
	fv, ok, err := i.GetFloat64()
	if ok {
		return fv, ok, nil
	}
	return 0, false, err
}

func (i Property) GetUint64() (uint64, bool, error) {
	if i.Type == cty.Number && i.Value != nil {
		var v uint64
		var err = json.Unmarshal(i.Value, &v)
		return v, err == nil, err
	}
	return 0, false, nil
}

func (i Property) GetUint32() (uint32, bool, error) {
	var v, ok, err = i.GetUint64()
	if err == nil && ok {
		return uint32(v), true, nil
	}
	return 0, ok, err
}

func (i Property) GetUint16() (uint16, bool, error) {
	var v, ok, err = i.GetUint64()
	if err == nil && ok {
		return uint16(v), true, nil
	}
	return 0, ok, err
}

func (i Property) GetUint8() (uint8, bool, error) {
	var v, ok, err = i.GetUint64()
	if err == nil && ok {
		return uint8(v), true, nil
	}
	return 0, ok, err
}

func (i Property) GetUint() (uint, bool, error) {
	var v, ok, err = i.GetUint64()
	if err == nil && ok {
		return uint(v), true, nil
	}
	return 0, ok, err
}

func (i Property) GetInt64() (int64, bool, error) {
	if i.Type == cty.Number && i.Value != nil {
		var v int64
		var err = json.Unmarshal(i.Value, &v)
		return v, err == nil, err
	}
	return 0, false, nil
}

func (i Property) GetInt32() (int32, bool, error) {
	var v, ok, err = i.GetInt64()
	if err == nil && ok {
		return int32(v), true, nil
	}
	return 0, ok, err
}

func (i Property) GetInt16() (int16, bool, error) {
	var v, ok, err = i.GetInt64()
	if err == nil && ok {
		return int16(v), true, nil
	}
	return 0, ok, err
}

func (i Property) GetInt8() (int8, bool, error) {
	var v, ok, err = i.GetInt64()
	if err == nil && ok {
		return int8(v), true, nil
	}
	return 0, ok, err
}

func (i Property) GetInt() (int, bool, error) {
	var v, ok, err = i.GetInt64()
	if err == nil && ok {
		return int(v), true, nil
	}
	return 0, ok, err
}

func (i Property) GetFloat64() (float64, bool, error) {
	if i.Type == cty.Number && i.Value != nil {
		var v float64
		var err = json.Unmarshal(i.Value, &v)
		return v, err == nil, err
	}
	return 0, false, nil
}

func (i Property) GetFloat32() (float32, bool, error) {
	var v, ok, err = i.GetFloat64()
	if err == nil && ok {
		return float32(v), true, nil
	}
	return 0, ok, err
}

func (i Property) GetDuration() (time.Duration, bool, error) {
	if i.Type == cty.String && i.Value != nil {
		var v time.Duration
		var err = json.Unmarshal(i.Value, &v)
		return v, err == nil, err
	}
	return 0, false, nil
}

func (i Property) GetBool() (bool, bool, error) {
	if i.Type == cty.Bool && i.Value != nil {
		var v bool
		var err = json.Unmarshal(i.Value, &v)
		return v, err == nil, err
	}
	return false, false, nil
}

func (i Property) GetString() (string, bool, error) {
	if i.Type == cty.String && i.Value != nil {
		var v string
		var err = json.Unmarshal(i.Value, &v)
		return v, err == nil, err
	}
	return "", false, nil
}

// GetSlice returns the underlay value as a slice with the given generic type,
// if not found or parse error, returns false.
func GetSlice[T any](i IProperty) ([]T, bool, error) {
	if i != nil {
		var typ, val = i.GetType(), i.GetValue()
		if (typ.IsListType() || typ.IsTupleType() || typ.IsSetType()) && val != nil {
			var v []T
			var err = json.Unmarshal(val, &v)
			return v, err == nil, err
		}
	}
	return nil, false, nil
}

// GetSet returns the underlay value as a set with the given generic type,
// if not found or parse error, returns false.
func GetSet[T comparable](i IProperty) (sets.Set[T], bool, error) {
	if i != nil {
		var typ, val = i.GetType(), i.GetValue()
		if (typ.IsListType() || typ.IsSetType()) && val != nil {
			var v []T
			var err = json.Unmarshal(val, &v)
			if err != nil {
				return nil, false, err
			}
			return sets.New[T](v...), true, nil
		}
	}
	return nil, false, nil
}

// GetMap returns the underlay value as a string map with the given generic type,
// if not found or parse error, returns false.
func GetMap[T any](i IProperty) (map[string]T, bool, error) {
	if i != nil {
		var typ, val = i.GetType(), i.GetValue()
		if (typ.IsMapType() || typ.IsObjectType()) && val != nil {
			var v map[string]T
			var err = json.Unmarshal(val, &v)
			return v, err == nil, err
		}
	}
	return nil, false, nil
}

// GetObject returns the underlay value as a T object with the given generic type,
// if not found or parse error, returns false.
func GetObject[T any](i IProperty) (T, bool, error) {
	var v T
	if i != nil {
		var typ, val = i.GetType(), i.GetValue()
		if (typ.IsMapType() || typ.IsObjectType()) && val != nil {
			var err = json.Unmarshal(val, &v)
			return v, err == nil, err
		}
	}
	return v, false, nil
}

// GetAny returns the underlay value as the given generic type,
// if not found or parse error, returns false.
func GetAny[T any](i IProperty) (T, bool, error) {
	var v T
	if i != nil {
		var val = i.GetValue()
		if val != nil {
			var err = json.Unmarshal(val, &v)
			return v, err == nil, err
		}
	}
	return v, false, nil
}

func (i Property) Cty() (cty.Type, cty.Value, error) {
	if i.Value == nil {
		return i.Type, cty.NilVal, nil
	}

	if i.Type == cty.NilType || i.Type == cty.DynamicPseudoType {
		// guess
		var v ctyjson.SimpleJSONValue
		if err := json.Unmarshal(i.Value, &v); err != nil {
			return i.Type, cty.NilVal, err
		}
		return v.Type(), v.Value, nil
	}

	var v, err = ctyjson.Unmarshal(i.Value, i.Type)
	return i.Type, v, err
}

// Properties holds the Property collection in map,
// the key of map is the name of Property,
// stores into json.
type Properties map[string]Property

// Value implements driver.Valuer.
func (i Properties) Value() (driver.Value, error) {
	return json.Marshal(i)
}

// Scan implements sql.Scanner.
func (i *Properties) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		return nil
	case []byte:
		return json.Unmarshal(v, i)
	}
	return errors.New("not a valid properties")
}

// Cty returns the type and value of this property.
func (i Properties) Cty() (cty.Type, cty.Value, error) {
	var (
		ot = make(map[string]cty.Type, len(i))
		ov = make(map[string]cty.Value, len(i))
	)
	for x := range i {
		var t, v, err = i[x].Cty()
		if err != nil {
			return cty.NilType, cty.NilVal, err
		}
		ot[x] = t
		ov[x] = v
	}
	return cty.Object(ot), cty.ObjectVal(ov), nil
}

// Values returns a map stores the underlay value.
func (i Properties) Values() Values {
	var m = make(Values, len(i))
	for x := range i {
		m[x] = i[x].GetValue()
	}
	return m
}

// TypedValues returns a map stores the typed value.
func (i Properties) TypedValues() (m map[string]any, err error) {
	m = make(map[string]any, len(i))
	for x := range i {
		var typ = i[x].GetType()
		switch {
		case typ == cty.Number:
			m[x], _, err = i[x].GetNumber()
		case typ == cty.Bool:
			m[x], _, err = i[x].GetBool()
		case typ == cty.String:
			m[x], _, err = i[x].GetString()
		case typ.IsListType() || typ.IsTupleType() || typ.IsSetType():
			m[x], _, err = GetSlice[any](i[x])
		case typ.IsMapType() || typ.IsObjectType():
			m[x], _, err = GetMap[any](i[x])
		default:
			m[x], _, err = GetAny[any](i[x])
		}
		if err != nil {
			return
		}
	}
	return
}

// Values holds the Value collection in map,
// the key of map is the name of Property,
// stores into json.
type Values map[string]Value

// Value implements driver.Valuer.
func (i Values) Value() (driver.Value, error) {
	return json.Marshal(i)
}

// Scan implements sql.Scanner.
func (i *Values) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		return nil
	case []byte:
		return json.Unmarshal(v, i)
	}
	return errors.New("not a valid property values")
}

// TypesWith aligns the property name and type with the given property.Schemas.
func (i Values) TypesWith(schemas Schemas) map[string]Type {
	var m = make(map[string]Type, len(i))
	for k := range i {
		m[k] = cty.DynamicPseudoType
	}
	for idx := range schemas {
		if _, exist := m[schemas[idx].Name]; !exist {
			continue
		}
		m[schemas[idx].Name] = schemas[idx].Type
	}
	return m
}

// StringTypesWith is similar with TypesWith,
// but returns the property type in string.
func (i Values) StringTypesWith(schemas Schemas) map[string]string {
	var m = make(map[string]string, len(i))
	var t = i.TypesWith(schemas)
	for k := range t {
		m[k] = typeexpr.TypeString(t[k])
	}
	return m
}

// ValidateWith validates the property value with the given schemas.
func (i Values) ValidateWith(schemas Schemas) error {
	// convert values to properties with the schemas.
	var props = Properties{}
	for _, s := range schemas {
		var v, exist = i[s.Name]
		if s.Required {
			// validate required value.
			if !exist || len(v) == 0 {
				return fmt.Errorf("not found required value %s", s.Name)
			}
		}
		if !exist {
			// ignore unspecified value.
			continue
		}
		// construct property.
		props[s.Name] = Property{
			Type:  s.Type,
			Value: v,
		}
	}

	// validate undefined value.
	if len(i) != len(props) {
		var l, r = sets.KeySet[string, Value](i), sets.KeySet[string, Property](props)
		return fmt.Errorf("found undefiend values %v", l.Difference(r).UnsortedList())
	}

	// validate converting.
	var _, _, err = props.Cty()
	if err != nil {
		return fmt.Errorf("unexpected value parsed: %w", err)
	}
	return nil
}

// Schema holds the schema of the property.
type Schema struct {
	// Name specifies the name of this property.
	Name string `json:"name"`
	// Description specifies the description of this property.
	Description string `json:"description,omitempty"`
	// Type specifies the type of this property.
	Type Type `json:"type"`
	// Default specifies the default value of this property.
	Default Value `json:"default,omitempty"`
	// Required indicates this property is required or not.
	Required bool `json:"required,omitempty"`
	// Sensitive indicates this property is sensitive or not.
	Sensitive bool `json:"sensitive,omitempty"`
	// Label specifies the UI label of this property.
	Label string `json:"label,omitempty"`
	// Group specifies the UI group of this property,
	// combines multiple levels with a slash.
	Group string `json:"group,omitempty"`
	// Options specifies available options of this property when the type is string.
	Options []Value `json:"options,omitempty"`
	// ShowIf specifies to show this property if the condition is true,
	// e.g. ShowIf: foo=bar.
	ShowIf string `json:"showIf,omitempty"`
	// Hidden specifies the field should be hidden or not,
	// default is visible.
	Hidden bool `json:"hidden,omitempty"`
}

// WithDescription indicates the description of schema.
func (i Schema) WithDescription(d string) Schema {
	i.Description = d
	return i
}

// WithRequired indicates the schema is required.
func (i Schema) WithRequired() Schema {
	i.Required = true
	return i
}

// WithSensitive indicates the schema is sensitive.
func (i Schema) WithSensitive() Schema {
	i.Sensitive = true
	return i
}

// WithLabel indicates the label of schema.
func (i Schema) WithLabel(l string) Schema {
	i.Label = l
	return i
}

// WithGroup indicates the group of schema.
func (i Schema) WithGroup(g string) Schema {
	i.Group = g
	return i
}

// WithOptions indicates the options of schema.
func (i Schema) WithOptions(o ...Value) Schema {
	i.Options = o
	return i
}

// WithShowIf indicates the condition of schema.
func (i Schema) WithShowIf(s string) Schema {
	i.ShowIf = s
	return i
}

// WithHidden indicates the schema is hidden.
func (i Schema) WithHidden() Schema {
	i.Hidden = true
	return i
}

// Schemas holds the Schema collection in slice,
// stores into json.
type Schemas []Schema

// Value implements driver.Valuer.
func (i Schemas) Value() (driver.Value, error) {
	return json.Marshal(i)
}

// Scan implements sql.Scanner.
func (i *Schemas) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		return nil
	case []byte:
		return json.Unmarshal(v, i)
	}
	return errors.New("not a valid named properties")
}
