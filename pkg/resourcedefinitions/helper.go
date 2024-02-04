package resourcedefinitions

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"reflect"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"golang.org/x/exp/slices"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/templateversion"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/templates/openapi"
	"github.com/seal-io/walrus/utils/json"
	"github.com/seal-io/walrus/utils/pointer"
	"github.com/seal-io/walrus/utils/strs"
)

const (
	variablesSchemaKey = types.VariableSchemaKey
	outputsSchemaKey   = types.OutputSchemaKey

	originalExtensionKey = openapi.ExtOriginalKey

	immutableExtensionKey = "immutable"
)

// GenSchema generates the schema for the resource definition.
//
// This function performs the following logic:
//  1. Quickly perform schema aligning to select common parts between multiple rules' template.
//  2. Refill the same default value to the variable schema if allowed.
//     The default value comes from the SchemaDefaultValue of the matching rule.
func GenSchema(ctx context.Context, mc model.ClientSet, df *model.ResourceDefinition) error {
	// Fetch schemas.
	scss := make([]openapi3.Schemas, 0, len(df.Edges.MatchingRules))
	{
		tvIDs := sets.New[object.ID]()
		for _, mr := range df.Edges.MatchingRules {
			tvIDs.Insert(mr.TemplateID)
		}

		tvs, err := mc.TemplateVersions().Query().
			Where(templateversion.IDIn(tvIDs.UnsortedList()...)).
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
				scss = append(scss, s.Components.Schemas)
			default:
				// Use the custom variables schema ref.
				ss, uiss := s.Components.Schemas, uis.Components.Schemas
				if uiss[variablesSchemaKey] != nil {
					// Override the original schema with custom schema.
					ss[variablesSchemaKey] = uiss[variablesSchemaKey]
				}

				scss = append(scss, ss)
			}
		}
	}

	// Align schemas.
	var (
		nb = map[string]any{
			variablesSchemaKey: map[string]any{},
			outputsSchemaKey:   map[string]any{},
		}
		scs = alignSchemas(nb, scss)
	)
	// Return directly if no schemas.
	if len(scs) == 0 {
		return nil
	}

	// Refill default value to variable schema.
	if sr, ok := scs[variablesSchemaKey]; ok {
		defs := make([][]byte, 0, len(df.Edges.MatchingRules))
		for _, mr := range df.Edges.MatchingRules {
			defs = append(defs, mr.SchemaDefaultValue)
		}

		refillVariableSchemaRef(
			nb[variablesSchemaKey].(map[string]any), "", sr, defs,
			[]string{""}, []*openapi3.Schema{sr.Value})
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

// alignSchemas aligns the schemas,
// which is performs on the first-level properties of the variables schemas and outputs schemas.
func alignSchemas(nb map[string]any, scs []openapi3.Schemas) openapi3.Schemas {
	ret := openapi3.Schemas{}

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

const (
	mefDescription = "description"
	mefNumber      = "number"
	mefLength      = "length"
	mefItems       = "items"
	mefEnum        = "enum"
	mefDefault     = "default"
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

				key = getKey(key, k)
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
		key = getKey(key, "@this.0")
		if nb[key] == nil {
			nb[key] = map[string]any{}
		}

		alignSchemaRef(nb, key, lv.Items, rv.Items)
	}

	if lv.Format == "" && rv.Format != "" {
		lv.Format = rv.Format
	}

	if !lv.WriteOnly && rv.WriteOnly {
		lv.WriteOnly = rv.WriteOnly
	}

	if !lv.ReadOnly && rv.ReadOnly {
		lv.ReadOnly = rv.ReadOnly
	}

	if len(lv.Extensions) != 0 && len(rv.Extensions) != 0 &&
		!isImmutable(lv.Extensions) && isImmutable(rv.Extensions) {
		if extensions, ok := lv.Extensions[openapi.ExtUIKey].(map[string]any); ok {
			extensions[immutableExtensionKey] = true
		}
	}

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

	// Get the optimal range.
	// If found mutually exclusive case, it cleans up all enum-related fields.
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

				// Lock default.
				lv.Default = nil
				nb[mefDefault] = &mutuallyExclusive
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
	// If found mutually exclusive case, it cleans up all length-related fields.
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
	// If found mutually exclusive case, it cleans up all items-related fields.
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
	// If found mutually exclusive case, it cleans up all number-related fields.
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

// isImmutable checks whether the schema is immutable in the UI extension.
func isImmutable(extensions map[string]any) bool {
	ui, ok := extensions[openapi.ExtUIKey]

	if !ok {
		return false
	}

	e, ok := ui.(map[string]any)

	if !ok {
		return false
	}

	return e[immutableExtensionKey] == true
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

// refillVariableSchemaRef refills the default value to the given schema ref.
//
// It compares the default values have the same key,
// only refills if all data are the same. If not the same, it will erase the default value and turn a required property to optional.
//
// The given notebook(`nb`) shares by the alignSchemas func,
// refill if a property doesn't mutually exclusive in the default value.
//
// The given property name(`pNames`) and schema(`pSchemas`) are used to turn the required property to optional.
func refillVariableSchemaRef(
	nb map[string]any, key string, sr *openapi3.SchemaRef, defs [][]byte,
	pNames []string, pSchemas []*openapi3.Schema,
) {
	if sr == nil || sr.Value == nil {
		return
	}

	v := sr.Value

	switch v.Type {
	case openapi3.TypeObject:
		switch {
		case v.Properties != nil:
			if isDefaultable(nb) {
				def, ok := alignDefaults(key, defs)
				if !ok {
					// Return directly if not found.
					return
				}

				// Refill the default value to the root,
				// but ignore this refilling when it is an array item.
				if !strings.HasSuffix(key, ".0") {
					v.Default = def
				}
			}

			for k := range v.Properties {
				nkey := getKey(key, k)
				if nb[nkey] == nil {
					nb[nkey] = map[string]any{}
				}

				if nb := nb[nkey].(map[string]any); isDefaultable(nb) {
					// Record all node from root to current point into the slice.
					ns := make([]string, 0, len(pNames)+1)
					ns = append(ns, pNames...)
					ns = append(ns, k)

					ss := make([]*openapi3.Schema, 0, len(pSchemas)+1)
					ss = append(ss, pSchemas...)
					ss = append(ss, v)

					refillVariableSchemaRef(
						nb, nkey, v.Properties[k], defs,
						ns, ss)
				}
			}
		case isDefaultable(nb) &&
			v.AdditionalProperties.Schema != nil && v.AdditionalProperties.Schema.Value != nil:
			def, ok := alignDefaults(key, defs)
			if !ok {
				// Return directly if not found.
				return
			}

			// Refill the default value to the additional properties,
			// but ignore this refilling when it is an array item.
			if !strings.HasSuffix(key, ".0") {
				v.Default = def
			}
		}

		// Turn property to optional if found different defaults.
		pName := pNames[len(pNames)-1]
		pSchema := pSchemas[len(pSchemas)-1]
		dropRequired(pSchema, v, pName)

		// Effect to the parent node also if the node has no required property.
		if len(pSchemas) > 1 {
			ppName := pNames[len(pNames)-2]
			ppSchema := pSchemas[len(pSchemas)-2]
			dropRequired(ppSchema, pSchema, ppName)
		}

		return
	case openapi3.TypeArray:
		nkey := getKey(key, "@this")

		def, ok := alignDefaults(nkey, defs)
		if !ok {
			// Return directly if not found.
			return
		}
		v.Default = def

		nkey = getKey(nkey, "0")
		if nb[nkey] == nil {
			nb[nkey] = map[string]any{}
		}

		if nb := nb[nkey].(map[string]any); isDefaultable(nb) {
			refillVariableSchemaRef(
				nb, nkey, v.Items, defs,
				pNames, pSchemas)
		}

		return
	}

	if isDefaultable(nb) {
		def, ok := alignDefaults(key, defs)
		if !ok {
			// Return directly if not found.
			return
		}

		v.Default = def

		// Turn property to optional if found different defaults.
		pName := pNames[len(pNames)-1]
		pSchema := pSchemas[len(pSchemas)-1]
		dropRequired(pSchema, v, pName)

		// Effect to the parent node also if the node has no required property.
		if len(pSchemas) > 1 {
			ppName := pNames[len(pNames)-2]
			ppSchema := pSchemas[len(pSchemas)-2]
			dropRequired(ppSchema, pSchema, ppName)
		}
	}
}

// dropRequired drops the given property name from the parent schema required list,
// if the property schema's required list is empty.
func dropRequired(parentSchema, propSchema *openapi3.Schema, propName string) {
	// Return directly, if required list is not empty.
	if len(propSchema.Required) != 0 {
		return
	}

	// Otherwise, drop the property from the required list.
	switch parentSchema.Type {
	default:
		return
	case openapi3.TypeObject:
	case openapi3.TypeArray:
		parentSchema = parentSchema.Items.Value
	}

	var i int
	for ; i < len(parentSchema.Required); i++ {
		if parentSchema.Required[i] == propName {
			break
		}
	}

	switch i {
	case len(parentSchema.Required):
		return
	case 0:
		parentSchema.Required = parentSchema.Required[1:]
	case len(parentSchema.Required) - 1:
		parentSchema.Required = parentSchema.Required[:i]
	default:
		parentSchema.Required = append(parentSchema.Required[:i], parentSchema.Required[i+1:]...)
	}
}

// isDefaultable checks whether the schema ref can set default.
func isDefaultable(nb map[string]any) bool {
	return nb == nil || nb[mefDefault] == nil
}

// getKey returns the key with prefix,
// if the prefix is empty, it returns the key directly.
func getKey(prefix, key string) string {
	if prefix == "" {
		return key
	}

	return prefix + "." + key
}

// alignDefaults compares the given key's value between all items of the given default value slice,
// if the value be same between all items, returns (default value, true),
// otherwise, returns (nil, true).
//
// This function skips the value if the key's prefix is not found in all items.
// For example, if any value of `a.b` from the given default value slice is not found,
// any key starts with `a.b` will be skipped, and returns (nil, false).
func alignDefaults(key string, defs [][]byte) (any, bool) {
	prefix := key
	if i := strings.LastIndex(key, "."); i > 0 {
		prefix = key[:i]
	}

	if prefix != "" {
		var found bool

		for _, d := range defs {
			jq := json.Get(d, key)
			if jq.Exists() {
				found = true
				break
			}
		}

		// Return directly if key is not found.
		if !found {
			return nil, false
		}
	}

	var def json.RawMessage

	for _, d := range defs {
		jq := json.Get(d, key)
		if !jq.Exists() {
			break
		}

		if def == nil {
			def = strs.ToBytes(&jq.Raw)
			continue
		}

		if !bytes.Equal(def, strs.ToBytes(&jq.Raw)) {
			def = nil
			break
		}
	}

	if def != nil {
		return def, true
	}

	// Avoid nil wrapped by type, e.g. json.RawMessage(nil).
	return nil, true
}
