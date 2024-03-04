package generators

import (
	"fmt"
	"io"
	"strings"

	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
	"k8s.io/utils/ptr"
)

func NewCRDGen(sanitizedName, outputPackage string) generator.Generator {
	return &crdGen{
		DefaultGen: generator.DefaultGen{
			OptionalName: sanitizedName,
		},
		outputPackage: outputPackage,
		imports:       generator.NewImportTrackerForPackage(outputPackage),
	}
}

type CRDTypeDefinition = apiext.CustomResourceDefinition

type crdGen struct {
	generator.DefaultGen

	outputPackage string
	imports       namer.ImportTracker

	types    []*types.Type
	typeDefs map[types.Name]*CRDTypeDefinition
}

func (g *crdGen) Filter(c *generator.Context, t *types.Type) bool {
	if t.Kind == types.Struct && t.Name.Package == g.outputPackage {
		if td := reflectType(c.Universe[g.outputPackage], t); td != nil {
			if g.typeDefs == nil {
				g.typeDefs = map[types.Name]*CRDTypeDefinition{}
			}

			g.types = append(g.types, t)
			g.typeDefs[t.Name] = td
		}
		return true
	}
	return false
}

func (g *crdGen) Namers(c *generator.Context) namer.NameSystems {
	return namer.NameSystems{
		// Have the raw namer for this file track what it imports.
		"raw": namer.NewRawNamer(g.outputPackage, g.imports),
		"private": &namer.NameStrategy{
			Join: func(pre string, in []string, post string) string {
				return strings.Join(in, "_")
			},
			PrependPackageNames: 4,
		},
	}
}

func (g *crdGen) isOtherPackage(pkg string) bool {
	switch {
	default:
		return true
	case pkg == g.outputPackage:
	case strings.HasSuffix(pkg, `"`+g.outputPackage+`"`):
	}
	return false
}

func (g *crdGen) Imports(c *generator.Context) []string {
	var pkgs []string
	for _, pkg := range g.imports.ImportLines() {
		if g.isOtherPackage(pkg) {
			pkgs = append(pkgs, pkg)
		}
	}
	return pkgs
}

func (g *crdGen) Init(c *generator.Context, w io.Writer) error {
	args := getCustomResourceDefinitionTypedArgs()
	sw := generator.NewSnippetWriter(w, c, "$", "$")

	// START: GetCustomResourceDefinitions
	sw.Do("func GetCustomResourceDefinitions() map[string]*$.CustomResourceDefinition|raw$ {\n", args)
	sw.Do("return map[string]*$.CustomResourceDefinition|raw$ {\n", args)

	for _, t := range g.types {
		if g.typeDefs[t.Name] == nil {
			continue
		}
		typedArgs := generator.Args{"TYPE": t}
		sw.Do("\"$.TYPE.Name.Name$\": crd_$.TYPE|private$(),\n", typedArgs)
	}

	sw.Do("}\n", nil)
	sw.Do("}\n\n", nil)
	// END: GetCustomResourceDefinitions

	return sw.Error()
}

func (g *crdGen) GenerateType(c *generator.Context, t *types.Type, w io.Writer) error {
	crd := g.typeDefs[t.Name]
	if crd == nil {
		return nil
	}

	args := getCustomResourceDefinitionTypedArgs().With("TYPE", t)
	sw := generator.NewSnippetWriter(w, c, "$", "$")
	sw.Do("func crd_$.TYPE|private$() *$.CustomResourceDefinition|raw$ {\n", args)
	sw.Do("return &$.CustomResourceDefinition|raw$ {\n", args)

	args = args.With("CRD", crd)
	// START: TypeMeta.
	sw.Do("TypeMeta: $.TypeMeta|raw${\n", args)
	sw.Do("APIVersion: \"apiextensions.k8s.io/v1\",\n", nil)
	sw.Do("Kind: \"CustomResourceDefinition\",\n", nil)
	sw.Do("},\n", nil)
	// END: TypeMeta.
	// START: ObjectMeta.
	sw.Do("ObjectMeta: $.ObjectMeta|raw${\n", args)
	sw.Do("Name: \"$.$\",\n", crd.Name)
	sw.Do("},\n", nil)
	// END: ObjectMeta.

	spec := crd.Spec
	args = args.With("SPEC", spec)
	// START: Spec.
	sw.Do("Spec: $.CustomResourceDefinitionSpec|raw${\n", args)
	sw.Do("Group: \"$.SPEC.Group$\",\n", args)
	sw.Do("Names: $.CustomResourceDefinitionNames|raw${\n", args)
	{
		sw.Do("Plural: \"$.SPEC.Names.Plural$\",\n", args)
		if spec.Names.Singular != "" {
			sw.Do("Singular: \"$.SPEC.Names.Singular$\",\n", args)
		}
		if len(spec.Names.ShortNames) > 0 {
			sw.Do("ShortNames: []string{\n\"$.$\",\n},\n", strings.Join(spec.Names.ShortNames, "\",\n \""))
		}
		sw.Do("Kind: \"$.SPEC.Names.Kind$\",\n", args)
		if spec.Names.ListKind != "" {
			sw.Do("ListKind: \"$.SPEC.Names.ListKind$\",\n", args)
		}
		if len(spec.Names.Categories) > 0 {
			sw.Do("Categories: []string{\n\"$.$\",\n},\n", strings.Join(spec.Names.Categories, "\",\n \""))
		}
	}
	sw.Do("},\n", nil)
	sw.Do("Scope: \"$.SPEC.Scope$\",\n", args)

	ver := spec.Versions[0]
	args = args.With("VERSION", ver)
	// START: Spec/Versions/0.
	sw.Do("Versions: []$.CustomResourceDefinitionVersion|raw${\n", args)
	sw.Do("{\n", nil)
	sw.Do("Name: \"$.VERSION.Name$\",\n", args)
	sw.Do("Served: $.VERSION.Served$,\n", args)
	sw.Do("Storage: $.VERSION.Storage$,\n", args)

	// START: Spec/Versions/0/Schema.
	sw.Do("Schema: &$.CustomResourceValidation|raw${\n", args)
	sw.Do("OpenAPIV3Schema: &$.JSONSchemaProps|raw${\n", args)
	writeSchema(sw, ver.Schema.OpenAPIV3Schema)
	sw.Do("},\n", nil)
	sw.Do("},\n", nil)
	// END: Spec/Versions/0/Schema.

	// START: Spec/Versions/0/SubResources.
	if subResources := ver.Subresources; subResources != nil {
		sw.Do("Subresources: &$.CustomResourceSubresources|raw${\n", args)
		if subResources.Status != nil {
			sw.Do("Status: &$.CustomResourceSubresourceStatus|raw${},\n", args)
		}
		if subResources.Scale != nil {
			sw.Do("Scale: &$.CustomResourceSubresourceScale|raw${},\n", args)
		}
		sw.Do("},\n", nil)
	}
	// END: Spec/Versions/0/SubResources.

	// START: Spec/Versions/0/AdditionalPrinterColumns.
	if printerColumns := ver.AdditionalPrinterColumns; len(printerColumns) > 0 {
		sw.Do("AdditionalPrinterColumns: []$.CustomResourceColumnDefinition|raw${\n", args)
		for _, printerColumn := range printerColumns {
			colArgs := args.With("COLUMN", printerColumn)
			sw.Do("{\n", nil)
			sw.Do("Name: \"$.COLUMN.Name$\",\n", colArgs)
			sw.Do("Type: \"$.COLUMN.Type$\",\n", colArgs)
			sw.Do("Format: \"$.COLUMN.Format$\",\n", colArgs)
			sw.Do("Description: $.$,\n", fmt.Sprintf("%#v", printerColumn.Description))
			sw.Do("Priority: \"$.COLUMN.Priority$\",\n", colArgs)
			sw.Do("JSONPath: \"$.COLUMN.JSONPath$\",\n", colArgs)
			sw.Do("},\n", nil)
		}
		sw.Do("},\n", nil)
	}
	// END: Spec/Versions/0/AdditionalPrinterColumns.

	sw.Do("},\n", nil)
	sw.Do("},\n", nil)
	// END: Spec/Versions/0.

	sw.Do("},\n", nil)
	// END: Spec.

	sw.Do("}\n", nil)
	sw.Do("}\n\n", nil)

	return sw.Error()
}

func writeSchema(sw *generator.SnippetWriter, schema *apiext.JSONSchemaProps) {
	if sw == nil || schema == nil {
		return
	}

	args := getJSONSchemaTypedArgs().With("SCHEMA", schema)

	if schema.Ref != nil {
		sw.Do("Ref: $.PtrTo|raw$(\"$.SCHEMA.Ref$\"),\n", args)
	}
	if schema.Description != "" {
		sw.Do("Description: $.$,\n", fmt.Sprintf("%#v", schema.Description))
	}
	if schema.Type != "" {
		sw.Do("Type: \"$.SCHEMA.Type$\",\n", args)
	}
	if schema.Format != "" {
		sw.Do("Format: \"$.SCHEMA.Format$\",\n", args)
	}
	if schema.Title != "" {
		sw.Do("Title: \"$.SCHEMA.Title$\",\n", args)
	}
	if def := schema.Default; def != nil {
		sw.Do("Default: &$.JSON|raw${\n", args)
		sw.Do("Raw: []byte(`$.$`),\n", string(def.Raw))
		sw.Do("},\n", nil)
	}
	if schema.Maximum != nil {
		sw.Do("Maximum: $.PtrTo|raw$[float64]($.SCHEMA.Maximum$),\n", args)
	}
	if schema.ExclusiveMaximum {
		sw.Do("ExclusiveMaximum: true,\n", nil)
	}
	if schema.Minimum != nil {
		sw.Do("Minimum: $.PtrTo|raw$[float64]($.SCHEMA.Minimum$),\n", args)
	}
	if schema.ExclusiveMinimum {
		sw.Do("ExclusiveMinimum: true,\n", nil)
	}
	if schema.MaxLength != nil {
		sw.Do("MaxLength: $.PtrTo|raw$[int64]($.SCHEMA.MaxLength$),\n", args)
	}
	if schema.MinLength != nil {
		sw.Do("MinLength: $.PtrTo|raw$[int64]($.SCHEMA.MinLength$),\n", args)
	}
	if schema.Pattern != "" {
		sw.Do("Pattern: \"$.SCHEMA.Pattern$\",\n", args)
	}
	if schema.MaxItems != nil {
		sw.Do("MaxItems: $.PtrTo|raw$[int64]($.SCHEMA.MaxItems$),\n", args)
	}
	if schema.MinItems != nil {
		sw.Do("MinItems: $.PtrTo|raw$[int64]($.SCHEMA.MinItems$),\n", args)
	}
	if schema.UniqueItems {
		sw.Do("UniqueItems: true,\n", nil)
	}
	if schema.MultipleOf != nil {
		sw.Do("MultipleOf: $.PtrTo|raw$[float64]($.SCHEMA.MultipleOf$),\n", args)
	}
	if len(schema.Enum) != 0 {
		sw.Do("Enum: []$.JSON|raw${\n", args)
		for _, e := range schema.Enum {
			sw.Do("{\n", nil)
			sw.Do("Raw: []byte(`$.$`),\n", string(e.Raw))
			sw.Do("},\n", nil)
		}
		sw.Do("},\n", nil)
	}
	if schema.MaxProperties != nil {
		sw.Do("MaxProperties: $.PtrTo|raw$[int64]($.SCHEMA.MaxProperties$),\n", args)
	}
	if schema.MinProperties != nil {
		sw.Do("MinProperties: $.PtrTo|raw$[int64]($.SCHEMA.MinProperties$),\n", args)
	}
	if len(schema.Required) != 0 {
		sw.Do("Required: []string{\n\"$.$\",\n},\n", strings.Join(schema.Required, "\",\n \""))
	}
	if items := schema.Items; items != nil {
		sw.Do("Items: &$.JSONSchemaPropsOrArray|raw${\n", args)
		if items.Schema != nil {
			sw.Do("Schema: &$.JSONSchemaProps|raw${\n", args)
			writeSchema(sw, items.Schema)
			sw.Do("},\n", nil)
		} else {
			sw.Do("JSONSchemas: []$.JSONSchemaProps|raw${\n", args)
			for i := range items.JSONSchemas {
				sw.Do("{\n", nil)
				writeSchema(sw, ptr.To(items.JSONSchemas[i]))
				sw.Do("},\n", nil)
			}
			sw.Do("},\n", nil)
		}
		sw.Do("},\n", nil)
	}
	if len(schema.AllOf) != 0 {
		sw.Do("AllOf: []$.JSONSchemaProps|raw${\n", args)
		for i := range schema.AllOf {
			sw.Do("{\n", nil)
			writeSchema(sw, ptr.To(schema.AllOf[i]))
			sw.Do("},\n", nil)
		}
		sw.Do("},\n", nil)
	}
	if len(schema.OneOf) != 0 {
		sw.Do("OneOf: []$.JSONSchemaProps|raw${\n", args)
		for i := range schema.OneOf {
			sw.Do("{\n", nil)
			writeSchema(sw, ptr.To(schema.OneOf[i]))
			sw.Do("},\n", nil)
		}
		sw.Do("},\n", nil)
	}
	if len(schema.AnyOf) != 0 {
		sw.Do("AnyOf: []$.JSONSchemaProps|raw${\n", args)
		for i := range schema.AnyOf {
			sw.Do("{\n", nil)
			writeSchema(sw, ptr.To(schema.AnyOf[i]))
			sw.Do("},\n", nil)
		}
		sw.Do("},\n", nil)
	}
	if schema.Not != nil {
		sw.Do("Not: &$.JSONSchemaProps|raw${\n", args)
		writeSchema(sw, schema.Not)
		sw.Do("},\n", nil)
	}
	if len(schema.Properties) != 0 {
		sw.Do("Properties: map[string]$.JSONSchemaProps|raw${\n", args)
		for _, name := range sets.List(sets.KeySet(schema.Properties)) {
			sw.Do("\"$.$\": {\n", name)
			writeSchema(sw, ptr.To(schema.Properties[name]))
			sw.Do("},\n", nil)
		}
		sw.Do("},\n", nil)
	}
	if addProps := schema.AdditionalProperties; addProps != nil {
		sw.Do("AdditionalProperties: &$.JSONSchemaPropsOrBool|raw${\n", args)
		sw.Do("Allows: $.$,\n", addProps.Allows)
		if addProps.Schema != nil {
			sw.Do("Schema: &$.JSONSchemaProps|raw${\n", args)
			writeSchema(sw, addProps.Schema)
			sw.Do("},\n", nil)
		}
		sw.Do("},\n", nil)
	}
	if len(schema.PatternProperties) != 0 {
		sw.Do("PatternProperties: map[string]$.JSONSchemaProps|raw${\n", args)
		for _, name := range sets.List(sets.KeySet(schema.PatternProperties)) {
			sw.Do("\"$.$\": {\n", name)
			writeSchema(sw, ptr.To(schema.PatternProperties[name]))
			sw.Do("},\n", nil)
		}
		sw.Do("},\n", nil)
	}
	if len(schema.Dependencies) != 0 {
		sw.Do("Dependencies: map[string]$.JSONSchemaPropsOrStringArray|raw${\n", args)
		for _, name := range sets.List(sets.KeySet(schema.Dependencies)) {
			dep := schema.Dependencies[name]
			sw.Do("\"$.$\": {\n", name)
			if dep.Schema != nil {
				sw.Do("Schema: &$.JSONSchemaProps|raw${\n", args)
				writeSchema(sw, dep.Schema)
				sw.Do("},\n", nil)
			} else {
				sw.Do("Property: []string{\n\"$.$\",\n},\n", strings.Join(dep.Property, "\",\n \""))
			}
			sw.Do("},\n", nil)
		}
		sw.Do("},\n", nil)
	}
	if addItems := schema.AdditionalItems; addItems != nil {
		sw.Do("AdditionalItems: &$.JSONSchemaPropsOrBool|raw${\n", args)
		sw.Do("Allows: $.$,\n", addItems.Allows)
		if addItems.Schema != nil {
			sw.Do("Schema: &$.JSONSchemaProps|raw${\n", args)
			writeSchema(sw, addItems.Schema)
			sw.Do("},\n", nil)
		}
		sw.Do("},\n", nil)
	}
	if len(schema.Definitions) != 0 {
		sw.Do("Definitions: map[string]$.JSONSchemaProps|raw${\n", args)
		for _, name := range sets.List(sets.KeySet(schema.Definitions)) {
			sw.Do("\"$.$\": {\n", name)
			writeSchema(sw, ptr.To(schema.Definitions[name]))
			sw.Do("},\n", nil)
		}
		sw.Do("},\n", nil)
	}
	if extDocs := schema.ExternalDocs; extDocs != nil {
		sw.Do("ExternalDocs: &$.ExternalDocumentation|raw${\n", args)
		sw.Do("Description: $.$,\n", fmt.Sprintf("%#v", extDocs.Description))
		sw.Do("URL: \"$.SCHEMA.ExternalDocs.URL$\",\n", args)
		sw.Do("},\n", nil)
	}
	if exp := schema.Example; exp != nil {
		sw.Do("Example: &$.JSON|raw${\n", args)
		sw.Do("Raw: []byte(`$.$`),\n", string(exp.Raw))
		sw.Do("},\n", nil)
	}
	if schema.Nullable {
		sw.Do("Nullable: true,\n", nil)
	}
	if schema.XPreserveUnknownFields != nil {
		sw.Do("XPreserveUnknownFields: $.PtrTo|raw$[bool]($.SCHEMA.XPreserveUnknownFields$),\n", args)
	}
	if schema.XEmbeddedResource {
		sw.Do("XEmbeddedResource: true,\n", nil)
	}
	if schema.XIntOrString {
		sw.Do("XIntOrString: true,\n", nil)
	}
	if len(schema.XListMapKeys) != 0 {
		sw.Do("XListMapKeys: []string{\n\"$.$\",\n},\n", strings.Join(schema.XListMapKeys, "\",\n \""))
	}
	if schema.XListType != nil {
		sw.Do("XListType: $.PtrTo|raw$[string](\"$.SCHEMA.XListType$\"),\n", args)
	}
	if schema.XMapType != nil {
		sw.Do("XMapType: $.PtrTo|raw$[string](\"$.SCHEMA.XMapType$\"),\n", args)
	}
	if len(schema.XValidations) != 0 {
		sw.Do("XValidations: []$.ValidationRule|raw${\n", args)
		for _, val := range schema.XValidations {
			valArgs := args.With("VALIDATION", val)
			sw.Do("{\n", nil)
			sw.Do("Rule: $.$,\n", fmt.Sprintf("%#v", val.Rule))
			if val.Message != "" {
				sw.Do("Message: $.$,\n", fmt.Sprintf("%#v", val.Message))
			}
			if val.MessageExpression != "" {
				sw.Do("MessageExpression: $.$,\n", fmt.Sprintf("%#v", val.MessageExpression))
			}
			if val.Reason != nil {
				sw.Do("Reason: $.PtrTo|raw$[$.FieldValueErrorReason|raw$](\"$.VALIDATION.Reason$\"),\n", valArgs)
			}
			if val.FieldPath != "" {
				sw.Do("FieldPath: $.$,\n", fmt.Sprintf("%#v", val.FieldPath))
			}
			if val.OptionalOldSelf != nil {
				sw.Do("OptionalOldSelf: $.PtrTo|raw$[bool]($.VALIDATION.OptionalOldSelf$),\n", valArgs)
			}
			sw.Do("},\n", nil)
		}
		sw.Do("},\n", nil)
	}
}

const (
	pkgUtilPtr              = "k8s.io/utils/ptr"
	pkgAPIMachineryMetadata = "k8s.io/apimachinery/pkg/apis/meta/v1"
	pkgTypeAPIExtensions    = "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func getCustomResourceDefinitionTypedArgs() generator.Args {
	return generator.Args{
		"TypeMeta":                        types.Ref(pkgAPIMachineryMetadata, "TypeMeta"),
		"ObjectMeta":                      types.Ref(pkgAPIMachineryMetadata, "ObjectMeta"),
		"SchemeGroupVersion":              types.Ref(pkgTypeAPIExtensions, "SchemeGroupVersion"),
		"CustomResourceDefinition":        types.Ref(pkgTypeAPIExtensions, "CustomResourceDefinition"),
		"CustomResourceDefinitionSpec":    types.Ref(pkgTypeAPIExtensions, "CustomResourceDefinitionSpec"),
		"CustomResourceDefinitionNames":   types.Ref(pkgTypeAPIExtensions, "CustomResourceDefinitionNames"),
		"CustomResourceDefinitionVersion": types.Ref(pkgTypeAPIExtensions, "CustomResourceDefinitionVersion"),
		"CustomResourceValidation":        types.Ref(pkgTypeAPIExtensions, "CustomResourceValidation"),
		"CustomResourceSubresources":      types.Ref(pkgTypeAPIExtensions, "CustomResourceSubresources"),
		"CustomResourceSubresourceStatus": types.Ref(pkgTypeAPIExtensions, "CustomResourceSubresourceStatus"),
		"CustomResourceSubresourceScale":  types.Ref(pkgTypeAPIExtensions, "CustomResourceSubresourceScale"),
		"CustomResourceColumnDefinition":  types.Ref(pkgTypeAPIExtensions, "CustomResourceColumnDefinition"),
		"JSONSchemaProps":                 types.Ref(pkgTypeAPIExtensions, "JSONSchemaProps"),
	}
}

func getJSONSchemaTypedArgs() generator.Args {
	return generator.Args{
		"PtrTo":                        types.Ref(pkgUtilPtr, "To"),
		"JSONSchemaProps":              types.Ref(pkgTypeAPIExtensions, "JSONSchemaProps"),
		"JSON":                         types.Ref(pkgTypeAPIExtensions, "JSON"),
		"JSONSchemaPropsOrArray":       types.Ref(pkgTypeAPIExtensions, "JSONSchemaPropsOrArray"),
		"JSONSchemaPropsOrBool":        types.Ref(pkgTypeAPIExtensions, "JSONSchemaPropsOrBool"),
		"JSONSchemaDependencies":       types.Ref(pkgTypeAPIExtensions, "JSONSchemaDependencies"),
		"JSONSchemaPropsOrStringArray": types.Ref(pkgTypeAPIExtensions, "JSONSchemaPropsOrStringArray"),
		"JSONSchemaDefinitions":        types.Ref(pkgTypeAPIExtensions, "JSONSchemaDefinitions"),
		"ExternalDocumentation":        types.Ref(pkgTypeAPIExtensions, "ExternalDocumentation"),
		"ValidationRules":              types.Ref(pkgTypeAPIExtensions, "ValidationRules"),
		"ValidationRule":               types.Ref(pkgTypeAPIExtensions, "ValidationRule"),
		"FieldValueErrorReason":        types.Ref(pkgTypeAPIExtensions, "FieldValueErrorReason"),
	}
}
