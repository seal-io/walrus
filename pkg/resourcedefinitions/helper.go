package resourcedefinitions

import (
	"context"
	"fmt"
	"math"
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tidwall/gjson"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/templateversion"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/templates/openapi"
	"github.com/seal-io/walrus/utils/json"
	"github.com/seal-io/walrus/utils/pointer"
)

const (
	variablesSchemaKey = types.VariableSchemaKey
	outputsSchemaKey   = types.OutputSchemaKey

	originalExtensionKey = openapi.ExtOriginalKey
)

// GenerateSchema generates definition schema with inputs/outputs intersection of matching template versions.
func GenerateSchema(ctx context.Context, mc model.ClientSet, df *model.ResourceDefinition) error {
	// Prepare the match rule schemas.
	var scss []openapi3.Schemas
	{
		// Map default variables.
		defaultsMapping := make(map[object.ID][]byte, len(df.Edges.MatchingRules))
		for _, mr := range df.Edges.MatchingRules {
			defaultsMapping[mr.TemplateID] = json.ShouldMarshal(mr.Attributes)
		}

		// Map schemas.
		schemasMapping := make(map[object.ID]openapi3.Schemas, len(df.Edges.MatchingRules))

		tvs, err := mc.TemplateVersions().Query().
			Where(templateversion.IDIn(sets.KeySet(defaultsMapping).UnsortedList()...)).
			All(ctx)
		if err != nil {
			return err
		}

		for i := range tvs {
			switch s, uis := tvs[i].Schema.OpenAPISchema, tvs[i].UiSchema.OpenAPISchema; {
			case s == nil || s.Components == nil || len(s.Components.Schemas) == 0:
				// Skip mutuallyExclusive schema.
				continue
			case uis == nil || uis.Components == nil || len(uis.Components.Schemas) == 0:
				// Use the original schema directly.
				schemasMapping[tvs[i].ID] = tvs[i].Schema.OpenAPISchema.Components.Schemas
			default:
				// Use the custom variables schema ref.
				ss, uiss := s.Components.Schemas, uis.Components.Schemas
				if uiss[variablesSchemaKey] != nil {
					// Override the original schema with custom schema.
					ss[variablesSchemaKey] = uiss[variablesSchemaKey]
				}
				schemasMapping[tvs[i].ID] = ss
			}

			// Merge defaults.
			for k := range schemasMapping {
				defaultVariablesSchemaRef(schemasMapping[k][variablesSchemaKey], defaultsMapping[k])
			}
		}

		// Flatten.
		scss = maps.Values(schemasMapping)
	}

	// Merge schemas.
	scs := alignSchemas(scss)
	if len(scs) == 0 {
		// Return directly if no schemas.
		return nil
	}

	df.Schema = types.Schema{
		OpenAPISchema: &openapi3.T{
			OpenAPI: openapi.OpenAPIVersion,
			Info: &openapi3.Info{
				Title:   fmt.Sprintf("OpenAPI schema for resource definition type %q", df.Type),
				Version: "v0.0.0",
			},
			Components: &openapi3.Components{
				Schemas: scs,
			},
		},
	}

	return nil
}

// defaultVariablesSchemaRef sets the default values for the variables schema ref.
func defaultVariablesSchemaRef(vsr *openapi3.SchemaRef, defs []byte) {
	if vsr == nil || vsr.Value == nil || vsr.Value.Type == "" {
		return
	}

	jp := json.Get(defs, "@this")
	if !jp.Exists() {
		return
	}

	switch {
	default:
		vsr.Value.Default = jp.Value()
	case vsr.Value.Type == openapi3.TypeObject:
		if !jp.IsObject() || jp.Raw == "{}" {
			return
		}

		// Default pure object.
		if len(vsr.Value.Properties) != 0 {
			for vk, vvsr := range vsr.Value.Properties {
				vr := jp.Get(vk)
				if !vr.Exists() {
					continue
				}

				defaultVariablesSchemaRef(vvsr, []byte(vr.Raw))
			}

			return
		}

		// Default the entry of map.
		defaultVariablesSchemaRef(vsr.Value.AdditionalProperties.Schema, []byte(jp.Raw))
	case vsr.Value.Type == openapi3.TypeArray:
		if !jp.IsArray() || jp.Raw == "[]" {
			return
		}

		// Default the element of array with the first available item.
		var ri gjson.Result
		for i, rs := 0, jp.Array(); i < len(rs) && !ri.Exists(); i++ {
			ri = rs[i]
		}

		if !ri.Exists() {
			return
		}

		defaultVariablesSchemaRef(vsr.Value.Items, []byte(ri.Raw))
	}
}

// alignSchemas aligns the schemas,
// which is performs on the first-level properties of the variables schemas and outputs schemas.
func alignSchemas(scs []openapi3.Schemas) openapi3.Schemas {
	var (
		ret = openapi3.Schemas{}
		nb  = map[string]any{}
	)

	for _, k := range []string{variablesSchemaKey, outputsSchemaKey} {
		var lsr *openapi3.SchemaRef

		for _, s := range scs {
			rsr := s[k]

			// If the schema ref is not object,
			// reset and skip.
			if rsr == nil || rsr.Value == nil || len(rsr.Value.Properties) == 0 {
				lsr = nil
				break
			}

			// Remove original extension.
			for _, p := range rsr.Value.Properties {
				if p.Value == nil || len(p.Value.Extensions) == 0 {
					continue
				}

				delete(p.Value.Extensions, originalExtensionKey)
			}

			// If the schema ref is nil,
			// reuse and go ahead.
			if lsr == nil {
				lsr = rsr
				continue
			}

			// Arrange the schema ref.
			alignSchemaRef(nb, k, lsr, rsr)
		}

		if lsr != nil && len(lsr.Value.Properties) != 0 {
			ret[k] = lsr
		}
	}

	return ret
}

var (
	mutuallyExclusive   struct{} // Used to mark mutually exclusive.
	sentenceTerminators = sets.New('.', '!', '。', '！')
)

// alignSchemaRef aligns the schema,
// which always move in the direction of shrinking at first, if not found, then expanding.
func alignSchemaRef(nb map[string]any, key string, lsr, rsr *openapi3.SchemaRef) {
	if nb == nil ||
		lsr == nil || lsr.Value == nil ||
		rsr == nil || rsr.Value == nil {
		return
	}

	lv, rv := lsr.Value, rsr.Value

	switch lv.Type {
	case openapi3.TypeObject:
		if lv.Properties != nil {
			reqs := sets.NewString(lv.Required...)

			for k := range lv.Properties {
				// Delete the key if the type is not the same.
				if !isSchemaRefTypeEqual(lv.Properties[k], rv.Properties[k]) {
					delete(lv.Properties, k)
					reqs.Delete(k)

					continue
				}

				key = key + "." + k
				if nb[key] == nil {
					nb[key] = map[string]any{}
				}

				alignSchemaRef(nb[key].(map[string]any), key, lv.Properties[k], rv.Properties[k])
			}

			// Reset required keys.
			if reqs.Len() == 0 {
				lsr.Value.Required = nil
			} else {
				lsr.Value.Required = reqs.List()
			}
		} else {
			alignSchemaRef(nb, key, lv.AdditionalProperties.Schema, rv.AdditionalProperties.Schema)
		}
	case openapi3.TypeArray:
		alignSchemaRef(nb, key, lv.Items, rv.Items)
	}

	if lv.Format == "" && rv.Format != "" {
		lv.Format = rv.Format
	}

	if !lv.WriteOnly && rv.WriteOnly {
		lv.WriteOnly = rv.WriteOnly
	}

	const (
		mefDescription = "description"
		mefNumber      = "number"
		mefLength      = "length"
		mefItems       = "items"
		mefEnum        = "enum"
		mefDefault     = "default"
	)

	// Find the most complete sentence.
	// If not found, it cleans up the description.
	if nb[mefDescription] == nil {
		switch {
		case lv.Description == "" && rv.Description != "":
			lv.Description = rv.Description
		case lv.Description != "" && rv.Description != "":
			var (
				d      []rune
				period int
			)
			lvd, rvd := []rune(lv.Description), []rune(rv.Description)
			s := int(math.Min(float64(len(lvd)), float64(len(rvd))))

			for i := 0; i < s; i++ {
				if lvd[i] != rvd[i] {
					break
				}

				d = append(d, lvd[i])

				if sentenceTerminators.Has(lvd[i]) {
					period = i
				}
			}

			if period != 0 {
				lv.Description = string(d)
				if period+1 < len(d) {
					lv.Description = string(d[:period+1])
				}
			} else {
				// Clean description.
				lv.Description = ""
				nb[mefDescription] = &mutuallyExclusive
			}
		}
	}

	// Confirm the default value.
	// If not the same, it cleans up the default value.
	if nb[mefDefault] == nil {
		if (lv.Default != nil && rv.Default == nil) ||
			(lv.Default != nil && rv.Default != nil && !reflect.DeepEqual(lv.Default, rv.Default)) {
			lv.Default = nil
			nb[mefDefault] = &mutuallyExclusive
		}
	}

	// Get the optimal range.
	// If found mutually mutuallyExclusive case, it cleans up all enum-related fields.
	if nb[mefEnum] == nil {
		switch {
		case lv.Enum == nil && len(rv.Enum) != 0: // Copy.
			lv.Enum = rv.Enum

			// Lock the first one as default.
			lv.Default = lv.Enum[0]
			nb[mefDefault] = &mutuallyExclusive
		case len(lv.Enum) != 0 && len(rv.Enum) != 0: // Merge.
			lve, rve := map[string]any{}, map[string]any{}

			for i := range lv.Enum {
				if reflect.ValueOf(lv.Enum[i]).IsZero() {
					continue
				}
				lve[string(json.ShouldMarshal(lv.Enum[i]))] = lv.Enum[i]
			}

			for i := range rv.Enum {
				if reflect.ValueOf(rv.Enum[i]).IsZero() {
					continue
				}
				rve[string(json.ShouldMarshal(rv.Enum[i]))] = rv.Enum[i]
			}

			is := sets.StringKeySet(lve).Intersection(sets.StringKeySet(rve))
			if is.Len() == 0 {
				// Clean enum.
				lv.Enum = nil
				nb[mefEnum] = &mutuallyExclusive
			} else {
				es := make([]any, 0, is.Len())
				for _, e := range is.List() {
					es = append(es, lve[e])
				}
				lv.Enum = es

				// Lock default.
				if !is.Has(string(json.ShouldMarshal(lv.Default))) {
					lv.Default = lv.Enum[0] // Pick the first one.
				}
				nb[mefDefault] = &mutuallyExclusive
			}
		}

		// NB(thxCode): Enum is mutually mutually exclusive with other ranges.
		if lv.Enum != nil && slices.Contains(basicTypes, lv.Type) {
			lv.MinLength = 0
			lv.MaxLength = nil
			nb[mefLength] = &mutuallyExclusive

			lv.MinItems = 0
			lv.MaxItems = nil
			nb[mefItems] = &mutuallyExclusive

			lv.Min = nil
			lv.Max = nil
			lv.ExclusiveMin = false
			lv.ExclusiveMax = false
			nb[mefNumber] = &mutuallyExclusive

			return
		}
	}

	// Get the optimal range.
	// If found mutually mutuallyExclusive case, it cleans up all length-related fields.
	if nb[mefLength] == nil {
		if lv.MinLength < rv.MinLength {
			lv.MinLength = rv.MinLength
		}

		if rv.MaxLength != nil && pointer.Deref(lv.MaxLength, 0) > pointer.Deref(rv.MaxLength, 0) {
			lv.MaxLength = rv.MaxLength
		}

		if lv.MaxLength != nil && lv.MinLength > *lv.MaxLength {
			// Clean length.
			lv.MinLength = 0
			lv.MaxLength = nil
			nb[mefLength] = &mutuallyExclusive
		}
	}

	// Get the optimal range.
	// If found mutually mutuallyExclusive case, it cleans up all items-related fields.
	if nb[mefItems] == nil {
		if lv.MinItems < rv.MinItems {
			lv.MinItems = rv.MinItems
		}

		if rv.MaxItems != nil && pointer.Deref(lv.MaxItems, 0) > pointer.Deref(rv.MaxItems, 0) {
			lv.MaxItems = rv.MaxItems
		}

		if lv.MaxItems != nil && lv.MinItems > *lv.MaxItems {
			// Clean items.
			lv.MinItems = 0
			lv.MaxItems = (*uint64)(nil)
			nb[mefItems] = &mutuallyExclusive
		}
	}

	// Get the optimal range.
	// If found mutually mutuallyExclusive case, it cleans up all number-related fields.
	if nb[mefNumber] == nil {
		lvn, rvn := pointer.Deref(lv.Min, 0), pointer.Deref(rv.Min, 0)
		if lv.ExclusiveMin {
			lvn++
		}

		if rv.ExclusiveMin {
			rvn++
		}

		if rv.Min != nil && lvn < rvn {
			lv.Min = rv.Min
			lv.ExclusiveMin = lv.ExclusiveMin || rv.ExclusiveMin
			lvn = rvn

			// Lock the min value as default.
			lv.Default = lvn
			nb[mefDefault] = &mutuallyExclusive
		}

		lvm, rvm := pointer.Deref(lv.Max, 0), pointer.Deref(rv.Max, 0)
		if lv.ExclusiveMax {
			lvm--
		}

		if rv.ExclusiveMax {
			rvm--
		}

		if rv.Max != nil && lvm > rvm {
			lv.Max = rv.Max
			lv.ExclusiveMax = lv.ExclusiveMax || rv.ExclusiveMax
			lvm = rvm
		}

		if lv.Min != nil && lv.Max != nil && lvn > lvm {
			// Clean default.
			lv.Default = nil
			nb[mefDefault] = &mutuallyExclusive

			// Clean number.
			lv.Min = (*float64)(nil)
			lv.Max = (*float64)(nil)
			lv.ExclusiveMin = false
			lv.ExclusiveMax = false
			nb[mefNumber] = &mutuallyExclusive
		}
	}
}

var basicTypes = []string{
	openapi3.TypeString,
	openapi3.TypeInteger,
	openapi3.TypeNumber,
	openapi3.TypeBoolean,
}

// isSchemaRefTypeEqual compares the type of the schema ref.
func isSchemaRefTypeEqual(lsr, rsr *openapi3.SchemaRef) bool {
	if lsr == nil || lsr.Value == nil || lsr.Value.Type == "" ||
		rsr == nil || rsr.Value == nil || rsr.Value.Type == "" {
		return false
	}

	switch lv, rv := lsr.Value, rsr.Value; {
	case lv.Type != rv.Type:
		return false
	case slices.Contains(basicTypes, lv.Type):
		return true
	case lv.Type == openapi3.TypeObject:
		// Compare type object.
		switch {
		case lv.Properties != nil && rv.Properties != nil:
			if len(lv.Properties) != len(rv.Properties) {
				return false
			}

			for k := range lv.Properties {
				if rv.Properties[k] == nil || !isSchemaRefTypeEqual(lv.Properties[k], rv.Properties[k]) {
					return false
				}
			}

			return true
		case lv.AdditionalProperties.Schema != nil && rv.AdditionalProperties.Schema != nil:
			return isSchemaRefTypeEqual(lv.AdditionalProperties.Schema, rv.AdditionalProperties.Schema)
		}

		return lv.Properties == nil && rv.Properties == nil &&
			lv.AdditionalProperties.Schema == nil && rv.AdditionalProperties.Schema == nil
	default:
		// Compare type array.
		return isSchemaRefTypeEqual(lv.Items, rv.Items)
	}
}
