package property

import (
	"fmt"
	"reflect"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/typeexpr"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/seal/utils/json"
	"github.com/seal-io/seal/utils/strs"
)

// TODO(thxCode): support tuple schema?

// Uint64Property wraps uint64 value into a property.
func Uint64Property(v uint64) Property {
	return Property{
		Type:  cty.Number,
		Value: json.MustMarshal(v),
	}
}

// Uint32Property wraps uint32 value into a property.
func Uint32Property(v uint32) Property {
	return Property{
		Type:  cty.Number,
		Value: json.MustMarshal(v),
	}
}

// Uint16Property wraps uint16 value into a property.
func Uint16Property(v uint16) Property {
	return Property{
		Type:  cty.Number,
		Value: json.MustMarshal(v),
	}
}

// Uint8Property wraps uint8 value into a property.
func Uint8Property(v uint8) Property {
	return Property{
		Type:  cty.Number,
		Value: json.MustMarshal(v),
	}
}

// UintProperty wraps uint value into a property.
func UintProperty(v uint) Property {
	return Property{
		Type:  cty.Number,
		Value: json.MustMarshal(v),
	}
}

// Int64Property wraps int64 value into a property.
func Int64Property(v int64) Property {
	return Property{
		Type:  cty.Number,
		Value: json.MustMarshal(v),
	}
}

// Int32Property wraps int32 value into a property.
func Int32Property(v int32) Property {
	return Property{
		Type:  cty.Number,
		Value: json.MustMarshal(v),
	}
}

// Int16Property wraps int16 value into a property.
func Int16Property(v int16) Property {
	return Property{
		Type:  cty.Number,
		Value: json.MustMarshal(v),
	}
}

// Int8Property wraps int8 value into a property.
func Int8Property(v int8) Property {
	return Property{
		Type:  cty.Number,
		Value: json.MustMarshal(v),
	}
}

// IntProperty wraps int value into a property.
func IntProperty(v int) Property {
	return Property{
		Type:  cty.Number,
		Value: json.MustMarshal(v),
	}
}

// Float64Property wraps float64 value into a property.
func Float64Property(v float64) Property {
	return Property{
		Type:  cty.Number,
		Value: json.MustMarshal(v),
	}
}

// Float32Property wraps float32 value into a property.
func Float32Property(v float32) Property {
	return Property{
		Type:  cty.Number,
		Value: json.MustMarshal(v),
	}
}

// DurationProperty wraps time.Duration value into a property.
func DurationProperty(v time.Duration) Property {
	return Property{
		Type:  cty.String,
		Value: json.MustMarshal(v),
	}
}

// BoolProperty wraps bool value into a property.
func BoolProperty(v bool) Property {
	return Property{
		Type:  cty.Bool,
		Value: json.MustMarshal(v),
	}
}

// StringProperty wraps string value into a property.
func StringProperty(v string) Property {
	return Property{
		Type:  cty.String,
		Value: json.MustMarshal(v),
	}
}

// SliceProperty wraps slice value into a property.
func SliceProperty[T any](v []T) Property {
	var t T
	var ty, err = gocty.ImpliedType(t)
	if err != nil {
		panic(fmt.Errorf("error getting implied type: %w", err))
	}
	return Property{
		Type:  cty.List(ty),
		Value: json.MustMarshal(v),
	}
}

// SetProperty wraps set value into a property.
func SetProperty[T comparable](v sets.Set[T]) Property {
	var t T
	var ty, err = gocty.ImpliedType(t)
	if err != nil {
		panic(fmt.Errorf("error getting implied type: %w", err))
	}
	return Property{
		Type:  cty.Set(ty),
		Value: json.MustMarshal(v.UnsortedList()),
	}
}

// MapProperty wraps map value into a property.
func MapProperty[T any](v map[string]T) Property {
	var t T
	var ty, err = gocty.ImpliedType(t)
	if err != nil {
		panic(fmt.Errorf("error getting implied type: %w", err))
	}
	return Property{
		Type:  cty.Map(ty),
		Value: json.MustMarshal(v),
	}
}

// ObjectProperty wraps object value into a property.
func ObjectProperty[T any](v T) Property {
	var t T
	var ty, err = gocty.ImpliedType(t)
	if err != nil {
		panic(fmt.Errorf("error getting implied type: %w", err))
	}
	if !ty.IsObjectType() {
		panic(fmt.Errorf("implied type is not object: %s", ty.GoString()))
	}
	return Property{
		Type:  ty,
		Value: json.MustMarshal(v),
	}
}

// AnyProperty wraps any value into a property.
func AnyProperty(v any) Property {
	return Property{
		Type:  cty.DynamicPseudoType,
		Value: json.MustMarshal(v),
	}
}

// Uint64Schema returns uint64 schema.
func Uint64Schema(n string, d *uint64) Schema {
	var s = Schema{
		Type: cty.Number,
		Name: n,
	}
	if d != nil {
		s.Default = json.MustMarshal(d)
	}
	return s
}

// Uint32Schema returns uint32 schema.
func Uint32Schema(n string, d *uint32) Schema {
	var s = Schema{
		Type: cty.Number,
		Name: n,
	}
	if d != nil {
		s.Default = json.MustMarshal(d)
	}
	return s
}

// Uint16Schema returns uint16 schema.
func Uint16Schema(n string, d *uint16) Schema {
	var s = Schema{
		Type: cty.Number,
		Name: n,
	}
	if d != nil {
		s.Default = json.MustMarshal(d)
	}
	return s
}

// Uint8Schema returns uint8 schema.
func Uint8Schema(n string, d *uint8) Schema {
	var s = Schema{
		Type: cty.Number,
		Name: n,
	}
	if d != nil {
		s.Default = json.MustMarshal(d)
	}
	return s
}

// UintSchema returns uint schema.
func UintSchema(n string, d *uint) Schema {
	var s = Schema{
		Type: cty.Number,
		Name: n,
	}
	if d != nil {
		s.Default = json.MustMarshal(d)
	}
	return s
}

// Int64Schema returns int64 schema.
func Int64Schema(n string, d *int64) Schema {
	var s = Schema{
		Type: cty.Number,
		Name: n,
	}
	if d != nil {
		s.Default = json.MustMarshal(d)
	}
	return s
}

// Int32Schema returns int32 schema.
func Int32Schema(n string, d *int32) Schema {
	var s = Schema{
		Type: cty.Number,
		Name: n,
	}
	if d != nil {
		s.Default = json.MustMarshal(d)
	}
	return s
}

// Int16Schema returns int16 schema.
func Int16Schema(n string, d *int16) Schema {
	var s = Schema{
		Type: cty.Number,
		Name: n,
	}
	if d != nil {
		s.Default = json.MustMarshal(d)
	}
	return s
}

// Int8Schema returns int8 schema.
func Int8Schema(n string, d *int8) Schema {
	var s = Schema{
		Type: cty.Number,
		Name: n,
	}
	if d != nil {
		s.Default = json.MustMarshal(d)
	}
	return s
}

// IntSchema returns int schema.
func IntSchema(n string, d *int) Schema {
	var s = Schema{
		Type: cty.Number,
		Name: n,
	}
	if d != nil {
		s.Default = json.MustMarshal(d)
	}
	return s
}

// Float64Schema returns float64 schema.
func Float64Schema(n string, d *float64) Schema {
	var s = Schema{
		Type: cty.Number,
		Name: n,
	}
	if d != nil {
		s.Default = json.MustMarshal(d)
	}
	return s
}

// Float32Schema returns float32 schema.
func Float32Schema(n string, d *float32) Schema {
	var s = Schema{
		Type: cty.Number,
		Name: n,
	}
	if d != nil {
		s.Default = json.MustMarshal(d)
	}
	return s
}

// DurationSchema returns time.Duration schema.
func DurationSchema(n string, d *time.Duration) Schema {
	var s = Schema{
		Type: cty.String,
		Name: n,
	}
	if d != nil {
		s.Default = json.MustMarshal(d)
	}
	return s
}

// BoolSchema returns bool schema.
func BoolSchema(n string, d *bool) Schema {
	var s = Schema{
		Type: cty.Bool,
		Name: n,
	}
	if d != nil {
		s.Default = json.MustMarshal(d)
	}
	return s
}

// StringSchema returns string schema.
func StringSchema(n string, d *string) Schema {
	var s = Schema{
		Type: cty.String,
		Name: n,
	}
	if d != nil {
		s.Default = json.MustMarshal(d)
	}
	return s
}

// SliceSchema returns []T schema.
func SliceSchema[T any](n string, d []T) Schema {
	var t T
	var ty, err = gocty.ImpliedType(t)
	if err != nil {
		panic(fmt.Errorf("error getting implied type: %w", err))
	}
	var s = Schema{
		Type: cty.List(ty),
		Name: n,
	}
	if d != nil {
		s.Default = json.MustMarshal(d)
	}
	return s
}

// SetSchema returns sets.Set[T] schema.
func SetSchema[T comparable](n string, d []T) Schema {
	var t T
	var ty, err = gocty.ImpliedType(t)
	if err != nil {
		panic(fmt.Errorf("error getting implied type: %w", err))
	}
	var s = Schema{
		Type: cty.Set(ty),
		Name: n,
	}
	if d != nil {
		s.Default = json.MustMarshal(d)
	}
	return s
}

// MapSchema returns map[string]T schema.
func MapSchema[T any](n string, d map[string]T) Schema {
	var t T
	var ty, err = gocty.ImpliedType(t)
	if err != nil {
		panic(fmt.Errorf("error getting implied type: %w", err))
	}
	var s = Schema{
		Type: cty.Map(ty),
		Name: n,
	}
	if d != nil {
		s.Default = json.MustMarshal(d)
	}
	return s
}

// ObjectSchema returns T schema.
func ObjectSchema[T any](n string, d T) Schema {
	var t T
	var ty, err = gocty.ImpliedType(t)
	if err != nil {
		panic(fmt.Errorf("error getting implied type: %w", err))
	}
	if !ty.IsObjectType() {
		panic(fmt.Errorf("implied type is not object: %s", ty.GoString()))
	}
	var s = Schema{
		Type: ty,
		Name: n,
	}
	if !reflect.ValueOf(d).IsZero() {
		s.Default = json.MustMarshal(d)
	}
	return s
}

// AnySchema returns any schema.
func AnySchema(n string, d any) Schema {
	var s = Schema{
		Type: cty.DynamicPseudoType,
		Name: n,
	}
	if d != nil {
		s.Default = json.MustMarshal(d)
	}
	return s
}

// GuessSchema guesses the schema with the given type and data,
// returns any schema if blank type and nil data(in fact, terraform validation must not let this pass),
// returns implied type schema if data is not nil,
// returns parsed type schema if type is not blank.
func GuessSchema(n string, t string, d any) (Schema, error) {
	if t == "" {
		if d == nil {
			// return any schema.
			return AnySchema(n, d), nil
		}

		// guess schema from data.
		var ty, err = gocty.ImpliedType(d)
		if err != nil {
			return Schema{}, err
		}
		return Schema{
			Type:    ty,
			Name:    n,
			Default: json.MustMarshal(d),
		}, nil
	}

	// parse type from type.
	var expr, diags = hclsyntax.ParseExpression(strs.ToBytes(&t), "", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return Schema{}, fmt.Errorf("error parsing expression: %w", diags)
	}
	ty, _, diags := typeexpr.TypeConstraintWithDefaults(expr)
	if diags.HasErrors() {
		return Schema{}, fmt.Errorf("error getting type: %w", diags)
	}
	var s json.RawMessage
	if d != nil {
		s = json.MustMarshal(d)
	}
	return Schema{
		Type:    ty,
		Name:    n,
		Default: s,
	}, nil
}
