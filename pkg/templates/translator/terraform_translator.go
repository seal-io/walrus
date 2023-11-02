package translator

import (
	"fmt"
	"sort"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"

	"github.com/seal-io/walrus/pkg/dao/types/property"
	"github.com/seal-io/walrus/pkg/templates/openapi"
	"github.com/seal-io/walrus/utils/json"
	"github.com/seal-io/walrus/utils/log"
)

// TerraformTranslator translates between terraform types and go types with openapi schema.
type TerraformTranslator struct{}

// NewTerraformTranslator creates a new terraform translator.
func NewTerraformTranslator() TerraformTranslator {
	return TerraformTranslator{}
}

// SchemaOfOriginalType generates openAPI schema from terraform type.
func (t TerraformTranslator) SchemaOfOriginalType(
	tp any,
	name string,
	def any,
	description string,
	sensitive bool,
) *openapi3.Schema {
	// Isn't terraform type.
	typ, ok := tp.(cty.Type)
	if !ok {
		return nil
	}

	switch {
	case typ == cty.Bool:
		s := openapi3.NewBoolSchema().
			WithDefault(def)

		s.Title = name
		s.Description = description
		s.WriteOnly = sensitive
		s.Extensions = openapi.NewExt(s.Extensions).
			SetOriginalType(typ).
			Export()

		return s
	case typ == cty.Number:
		s := openapi3.NewFloat64Schema().
			WithDefault(def)

		s.Title = name
		s.Description = description
		s.WriteOnly = sensitive
		s.Extensions = openapi.NewExt(s.Extensions).
			SetOriginalType(typ).
			Export()

		return s
	case typ == cty.String:
		s := openapi3.NewStringSchema().
			WithDefault(def)

		s.Title = name
		s.Description = description
		s.WriteOnly = sensitive
		s.Extensions = openapi.NewExt(s.Extensions).
			SetOriginalType(typ).
			Export()

		return s
	case typ.IsListType() || typ.IsSetType():
		it := t.SchemaOfOriginalType(typ.ElementType(), "", nil, "", sensitive)

		s := openapi3.NewArraySchema().
			WithDefault(def).
			WithItems(it)

		s.Title = name
		s.Description = description
		s.WriteOnly = sensitive
		s.Extensions = openapi.NewExt(s.Extensions).
			SetOriginalType(typ).
			Export()

		return s
	case typ.IsTupleType():
		var (
			ts   = typ.TupleElementTypes()
			refs = make([]*openapi3.SchemaRef, len(ts))
		)

		for i, tt := range ts {
			refs[i] = t.SchemaOfOriginalType(tt, "", nil, "", sensitive).
				NewRef()
		}

		s := openapi3.NewArraySchema().
			WithDefault(def).
			WithLength(int64(len(ts))).
			WithItems(&openapi3.Schema{
				OneOf: refs,
			})

		s.Title = name
		s.Description = description
		s.WriteOnly = sensitive
		s.Extensions = openapi.NewExt(s.Extensions).
			SetOriginalType(typ).
			Export()

		return s
	case typ.IsMapType():
		var (
			s = openapi3.NewObjectSchema().
				WithDefault(def)
			mtp = typ.MapElementType()
		)

		if mtp != nil {
			it := t.SchemaOfOriginalType(*mtp, "", nil, "", sensitive)
			s.WithAdditionalProperties(it)
		}

		s.Title = name
		s.Description = description
		s.WriteOnly = sensitive
		s.Extensions = openapi.NewExt(s.Extensions).
			SetOriginalType(typ).
			Export()

		return s
	case typ.IsObjectType():
		s := openapi3.NewObjectSchema().
			WithDefault(def)

		for n, tt := range typ.AttributeTypes() {
			st := t.SchemaOfOriginalType(tt, n, nil, "", sensitive)

			s.WithProperty(n, st)

			if !typ.AttributeOptional(n) {
				s.Required = append(s.Required, n)
			}
		}

		s.Title = name
		s.Description = description
		s.WriteOnly = sensitive
		s.Extensions = openapi.NewExt(s.Extensions).
			SetOriginalType(typ).
			Export()
		sort.Strings(s.Required)

		return s
	case typ == cty.DynamicPseudoType:
		// Empty Type.
		s := openapi3.NewSchema().
			WithDefault(def)

		s.Title = name
		s.Description = description
		s.WriteOnly = sensitive
		s.Extensions = openapi.NewExt(s.Extensions).
			SetOriginalType(cty.DynamicPseudoType).
			Export()

		return s
	default:
		log.Warnf("unsupported terraform type %s", typ.FriendlyName())
	}

	return nil
}

// ToGoTypeValues converts the values to go types.
func (t TerraformTranslator) ToGoTypeValues(
	values map[string]json.RawMessage,
	schema openapi3.Schema,
) (map[string]any, error) {
	// Language matching.
	if !t.SchemaMatched(schema) {
		return nil, nil
	}

	// Convert.
	r := make(map[string]any)

	for n, v := range values {
		if schema.Properties[n] == nil || schema.Properties[n].Value == nil {
			continue
		}

		var (
			s   = schema.Properties[n].Value
			err error
		)

		switch s.Type {
		case openapi3.TypeString:
			r[n], _, err = property.GetString(v)
		case openapi3.TypeBoolean:
			r[n], _, err = property.GetBool(v)
		case openapi3.TypeInteger:
			r[n], _, err = property.GetInt(v)
		case openapi3.TypeNumber:
			r[n], _, err = property.GetNumber(v)
		case openapi3.TypeArray:
			r[n], _, err = property.GetSlice[any](v)
		case openapi3.TypeObject:
			r[n], _, err = property.GetMap[any](v)
		default:
			r[n], _, err = property.GetAny[any](v)
		}

		if err != nil {
			log.Errorf("error converting value %v to go type: %v", v, err)
		}
	}

	return r, nil
}

// ToOriginalTypeValues Converts arbitrary go types to a cty Value.
func (t TerraformTranslator) ToOriginalTypeValues(values map[string]any) ([]string, map[string]cty.Value, error) {
	b, err := json.Marshal(values)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal value to json: %w", err)
	}

	var sjv ctyjson.SimpleJSONValue
	if err := sjv.UnmarshalJSON(b); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal json to cty value: %w", err)
	}

	var (
		val  = sjv.Value.AsValueMap()
		keys = make([]string, 0)
	)

	// Sorted Keys.
	for k := range val {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys, val, nil
}

// GetOriginalType returns the original type of the schema.
func (t TerraformTranslator) GetOriginalType(schema *openapi3.Schema) cty.Type {
	if schema != nil {
		ta := openapi.GetOriginalType(schema.Extensions)
		if ta != nil {
			if typ, ok := ta.(cty.Type); ok {
				return typ
			}
		}
	}

	return cty.DynamicPseudoType
}

func (t TerraformTranslator) SchemaMatched(schema openapi3.Schema) bool {
	// Language matching, always true since only terraform template now.
	return true
}
