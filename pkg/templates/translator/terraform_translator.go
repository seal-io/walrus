package translator

import (
	"fmt"
	"sort"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hashicorp/hcl/v2/ext/typeexpr"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"

	"github.com/seal-io/walrus/pkg/dao/types/property"
	"github.com/seal-io/walrus/pkg/templates/openapi"
	"github.com/seal-io/walrus/utils/json"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/strs"
)

// TerraformTranslator translates between terraform types and go types with openapi schema.
type TerraformTranslator struct{}

// NewTerraformTranslator creates a new terraform translator.
func NewTerraformTranslator() TerraformTranslator {
	return TerraformTranslator{}
}

// SchemaOfOriginalType generates openAPI schema from terraform type.
func (t TerraformTranslator) SchemaOfOriginalType(tp any, opts Options) *openapi3.Schema {
	// Isn't terraform type.
	typ, ok := tp.(cty.Type)
	if !ok {
		return nil
	}

	var (
		title       = strs.Title(opts.Name)
		def         = opts.Default
		description = opts.Description
		sensitive   = opts.Sensitive
		order       = opts.Order
		nullable    = opts.Nullable
		s           *openapi3.Schema
	)

	switch {
	case typ == cty.Bool:
		// Schema.
		s = openapi3.NewBoolSchema()

		// Default.
		setDefault(s, opts.Default)

		// Extensions.
		s.Extensions = openapi.NewExt().
			WithOriginalType(typ).
			WithUIOrder(order).
			Export()

	case typ == cty.Number:
		// Schema.
		s = openapi3.NewFloat64Schema()

		// Default.
		setDefault(s, def)

		// Extensions.
		s.Extensions = openapi.NewExt().
			WithOriginalType(typ).
			WithUIOrder(order).
			Export()

	case typ == cty.String:
		// Schema.
		s = openapi3.NewStringSchema()

		// Default.
		setDefault(s, def)

		if sensitive {
			s.Format = "password"
		}

		// Extensions.
		s.Extensions = openapi.NewExt().
			WithOriginalType(typ).
			WithUIOrder(order).
			Export()

	case typ.IsListType() || typ.IsSetType():
		// Schema.
		s = openapi3.NewArraySchema()

		// Default.
		var etpDef any
		switch def.(type) {
		case *typeexpr.Defaults:
			etpDef = getPropDefault(def, "")
		default:
			setDefault(s, def)
		}

		// Property.
		it := t.SchemaOfOriginalType(typ.ElementType(), Options{
			Default:     etpDef,
			Sensitive:   sensitive,
			Order:       -1,
			TypeExpress: getListItemExpression(opts.TypeExpress),
		})
		s.WithItems(it)

		// Extensions.
		s.Extensions = openapi.NewExt().
			WithOriginalType(typ).
			WithUIOrder(order).
			WithUIColSpan(12).
			Export()

	case typ.IsTupleType():
		// Schema.
		s = openapi3.NewArraySchema()

		// Property.
		var (
			ts   = typ.TupleElementTypes()
			refs = make([]*openapi3.Schema, len(ts))
			te   []hclsyntax.Expression
		)

		for i, tt := range ts {
			o := Options{
				Sensitive: sensitive,
				Order:     -1,
			}
			if len(te) > i {
				o.TypeExpress = te[i]
			}

			refs[i] = t.SchemaOfOriginalType(tt, o)
		}

		switch {
		case len(refs) == 1:
			s.WithItems(refs[0]).
				WithLength(1)
		case len(refs) > 1:
			var diffSchema bool

			for i := 0; i < len(refs); i++ {
				for j := i + 1; j < len(refs); j++ {
					if refs[i] != refs[j] {
						diffSchema = true
						break
					}
				}
			}

			if diffSchema {
				s.WithItems(openapi3.NewObjectSchema()).
					WithLength(int64(len(ts)))
			} else {
				s.WithItems(refs[0]).
					WithLength(int64(len(ts)))
			}
		default:
			s.WithItems(openapi3.NewObjectSchema())
		}

		// Extensions.
		s.Extensions = openapi.NewExt().
			WithOriginalType(typ).
			WithUIOrder(order).
			WithUIColSpan(12).
			Export()

	case typ.IsMapType():
		// Schema.
		s = openapi3.NewObjectSchema()

		// Default.
		var (
			mtp    = typ.MapElementType()
			mtpDef any
		)

		switch def.(type) {
		case *typeexpr.Defaults:
			mtpDef = getPropDefault(def, "")
		default:
			setDefault(s, def)
		}

		// Property.
		if mtp != nil {
			it := t.SchemaOfOriginalType(*mtp, Options{
				Default:     mtpDef,
				Sensitive:   sensitive,
				Order:       -1,
				TypeExpress: getMapItemExpression(opts.TypeExpress),
			})
			s.WithAdditionalProperties(it)
		}

		// Extensions.
		s.Extensions = openapi.NewExt().
			WithOriginalType(typ).
			WithUIOrder(order).
			WithUIColSpan(12).
			Export()

	case typ.IsObjectType():
		// Schema.
		s = openapi3.NewObjectSchema()

		// Default.
		var (
			defaultValues = make(map[string]cty.Value)
			childDefaults = make(map[string]*typeexpr.Defaults)
		)

		switch dv := def.(type) {
		case *typeexpr.Defaults:
			if dv != nil {
				if dv.DefaultValues != nil && len(dv.DefaultValues) > 0 {
					defaultValues = dv.DefaultValues
				}

				if dv.Children != nil && len(dv.Children) > 0 {
					childDefaults = dv.Children
				}
			}
		default:
			setDefault(s, def)
		}

		// Property Order.
		var (
			propOrder = make(map[string]int)
			propExpr  = make(map[string]any)
		)

		if opts.TypeExpress != nil {
			propOrder, propExpr = getObjectPropExpression(opts.TypeExpress)
		}

		// Property.
		for n, tt := range typ.AttributeTypes() {
			var (
				propDef      any
				propNullable bool
			)

			if defaultValues[n].IsKnown() {
				propDef = defaultValues[n]
			}

			if childDefaults[n] != nil {
				propDef = childDefaults[n]
			}

			if !typ.AttributeOptional(n) {
				s.Required = append(s.Required, n)
			} else {
				propNullable = true
			}

			st := t.SchemaOfOriginalType(tt, Options{
				Name:        n,
				Default:     propDef,
				Sensitive:   sensitive,
				Order:       propOrder[n],
				TypeExpress: propExpr[n],
				Nullable:    propNullable,
			})

			s.WithProperty(n, st)
		}

		// Extensions.
		s.Extensions = openapi.NewExt().
			WithOriginalType(typ).
			WithUIOrder(order).
			WithUIColSpan(12).
			Export()
		sort.Strings(s.Required)

	case typ == cty.DynamicPseudoType:
		// Empty Type.
		s = openapi3.NewObjectSchema()

		// Default.
		setDefault(s, def)

		// Extensions.
		s.Extensions = openapi.NewExt().
			WithOriginalType(cty.DynamicPseudoType).
			WithUIOrder(order).
			WithUIColSpan(12).
			Export()
	}

	if s == nil {
		log.Warnf("unsupported terraform type %s", typ.FriendlyName())
		return nil
	}

	s.Title = title
	s.Description = description
	s.WriteOnly = sensitive
	s.Nullable = nullable

	return s
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
		ta := openapi.GetExtOriginal(schema.Extensions).Type
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

func mustHclToGoValue(in any) any {
	if in == nil {
		return nil
	}

	val, ok := in.(cty.Value)
	if !ok {
		return nil
	}

	if !val.IsWhollyKnown() {
		return nil
	}

	valJSON, err := ctyjson.Marshal(val, val.Type())
	if err != nil {
		log.Warnf("failed to serialize value as JSON: %s", err)
		return nil
	}

	var goValue any

	err = json.Unmarshal(valJSON, &goValue)
	if err != nil {
		log.Warnf("failed re-parse value from JSON: %s", err)
		return nil
	}

	return goValue
}

func setDefault(s *openapi3.Schema, def any) {
	if def == nil {
		return
	}

	switch def.(type) {
	case cty.Value:
		s.WithDefault(mustHclToGoValue(def))
	default:
		s.WithDefault(def)
	}
}

func getPropDefault(def any, key string) any {
	if def == nil {
		return nil
	}

	pdef, ok := def.(*typeexpr.Defaults)
	if !ok || pdef == nil {
		return nil
	}

	if pdef.DefaultValues != nil && len(pdef.DefaultValues) > 0 && pdef.DefaultValues[key].IsKnown() {
		return pdef.DefaultValues[key]
	}

	if pdef.Children != nil && len(pdef.Children) > 0 && pdef.Children[key] != nil {
		return pdef.Children[key]
	}

	return nil
}

func getObjectPropExpression(expr any) (map[string]int, map[string]any) {
	var (
		propOrder = make(map[string]int)
		propExpr  = make(map[string]any)
	)

	fe, ok := expr.(*hclsyntax.FunctionCallExpr)
	if !ok || len(fe.Args) == 0 {
		return propOrder, propExpr
	}

	if fe.Name == "object" {
		switch arg := fe.Args[0].(type) {
		case *hclsyntax.ObjectConsExpr:
			for i, v := range arg.Items {
				obk, ok := v.KeyExpr.(*hclsyntax.ObjectConsKeyExpr)
				if !ok {
					continue
				}

				ste, ok := obk.Wrapped.(*hclsyntax.ScopeTraversalExpr)
				if ok {
					name := ste.Traversal.RootName()
					propOrder[name] = i + 1
					propExpr[name] = v.ValueExpr
				}
			}
		case *hclsyntax.FunctionCallExpr:
			if arg.Name == "optional" && len(arg.Args) != 0 {
				return getObjectPropExpression(arg.Args[0])
			}
		}

		return propOrder, propExpr
	}

	return getObjectPropExpression(fe.Args[0])
}

func getListItemExpression(expr any) hclsyntax.Expression {
	fe, ok := expr.(*hclsyntax.FunctionCallExpr)
	if !ok || len(fe.Args) == 0 {
		return nil
	}

	if fe.Name == "list" {
		return fe.Args[0]
	}

	return getListItemExpression(fe.Args[0])
}

func getMapItemExpression(expr any) hclsyntax.Expression {
	fe, ok := expr.(*hclsyntax.FunctionCallExpr)
	if !ok || len(fe.Args) == 0 {
		return nil
	}

	if fe.Name == "map" {
		return fe.Args[0]
	}

	return getMapItemExpression(fe.Args[0])
}
