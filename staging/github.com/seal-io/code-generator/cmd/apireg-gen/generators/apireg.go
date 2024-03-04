package generators

import (
	"io"
	"slices"
	"strings"

	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
	apireg "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
)

func NewAPIRegGen(sanitizedName, outputPackage string) generator.Generator {
	return &apiRegGen{
		DefaultGen: generator.DefaultGen{
			OptionalName: sanitizedName,
		},
		outputPackage: outputPackage,
		imports:       generator.NewImportTrackerForPackage(outputPackage),
	}
}

type (
	APIServiceDefinition  = apireg.APIService
	APIResourceDefinition struct {
		Scope        string
		Categories   []string
		ShortNames   []string
		Kind         string
		Singular     string
		Plural       string
		SubResources []string
	}
)

type apiRegGen struct {
	generator.DefaultGen

	outputPackage string
	imports       namer.ImportTracker

	pkgDef   *APIServiceDefinition
	types    []*types.Type
	typeDefs map[types.Name]*APIResourceDefinition
}

func (g *apiRegGen) Filter(c *generator.Context, t *types.Type) bool {
	if t.Kind == types.Struct && t.Name.Package == g.outputPackage {
		g.pkgDef = reflectPackage(c.Universe[g.outputPackage])
		if g.pkgDef == nil {
			return false
		}

		if td := reflectType(t); td != nil {
			if g.typeDefs == nil {
				g.typeDefs = map[types.Name]*APIResourceDefinition{}
			}

			g.types = append(g.types, t)
			g.typeDefs[t.Name] = td
		}
		return true
	}
	return false
}

func (g *apiRegGen) Namers(c *generator.Context) namer.NameSystems {
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

func (g *apiRegGen) isOtherPackage(pkg string) bool {
	switch {
	default:
		return true
	case pkg == g.outputPackage:
	case strings.HasSuffix(pkg, `"`+g.outputPackage+`"`):
	}
	return false
}

func (g *apiRegGen) Imports(c *generator.Context) []string {
	var pkgs []string
	for _, pkg := range g.imports.ImportLines() {
		if g.isOtherPackage(pkg) {
			pkgs = append(pkgs, pkg)
		}
	}
	return pkgs
}

func (g *apiRegGen) Init(c *generator.Context, w io.Writer) error {
	args := getAPIServiceTypedArgs()
	sw := generator.NewSnippetWriter(w, c, "$", "$")

	// START: Interfaces
	sw.Do("type (\n", nil)
	{
		sw.Do("WithStatusSubResource interface {\n", nil)
		sw.Do("$.ObjectMetaAccessor|raw$\n", args)
		sw.Do("$.Object|raw$\n", args)
		sw.Do("CopyStatusTo($.Object|raw$)\n", args)
		sw.Do("}\n\n", nil)
	}
	{
		sw.Do("WithScaleSubResource interface {\n", nil)
		sw.Do("$.ObjectMetaAccessor|raw$\n", args)
		sw.Do("$.Object|raw$\n", args)
		sw.Do("GetScale() *$.Scale|raw$ // TODO: Main struct needs to implement this. \n", args)
		sw.Do("SetScale(*$.Scale|raw$)  // TODO: Main struct needs to implement this. \n", args)
		sw.Do("}\n\n", nil)
	}
	sw.Do(")\n\n", nil)
	// END: Interfaces

	// START: GetAPIService
	sw.Do("func GetAPIService(svc $.ServiceReference|raw$, ca []byte) *$.APIService|raw$ {\n", args)
	sw.Do("return &$.APIService|raw${\n", args)
	args = args.With("APISERVICE", g.pkgDef)
	// START: TypeMeta.
	sw.Do("TypeMeta: $.TypeMeta|raw${\n", args)
	sw.Do("APIVersion: \"$.APISERVICE.TypeMeta.APIVersion$\",\n", args)
	sw.Do("Kind: \"$.APISERVICE.TypeMeta.Kind$\",\n", args)
	sw.Do("},\n", nil)
	// END: TypeMeta.
	// START: ObjectMeta.
	sw.Do("ObjectMeta: $.ObjectMeta|raw${\n", args)
	sw.Do("Name: \"$.APISERVICE.ObjectMeta.Name$\",\n", args)
	sw.Do("},\n", nil)
	// END: ObjectMeta.
	spec := g.pkgDef.Spec
	args = args.With("SPEC", spec)
	// START: Spec.
	sw.Do("Spec: $.APIServiceSpec|raw${\n", args)
	sw.Do("Service: svc.DeepCopy(),\n", args)
	sw.Do("Group: \"$.SPEC.Group$\",\n", args)
	sw.Do("Version: \"$.SPEC.Version$\",\n", args)
	sw.Do("InsecureSkipTLSVerify: $.SPEC.InsecureSkipTLSVerify$,\n", args)
	sw.Do("CABundle: $.BytesClone$(ca),\n", args)
	sw.Do("GroupPriorityMinimum: $.SPEC.GroupPriorityMinimum$,\n", args)
	sw.Do("VersionPriority: $.SPEC.VersionPriority$,\n", args)
	sw.Do("},\n", nil)
	// END: Spec.
	sw.Do("}\n", nil)
	sw.Do("}\n\n", nil)
	// END: GetAPIService

	return sw.Error()
}

func (g *apiRegGen) GenerateType(c *generator.Context, t *types.Type, w io.Writer) error {
	res := g.typeDefs[t.Name]
	if res == nil {
		return nil
	}

	args := getAPIResourceTypedArgs().With("TYPE", t)
	sw := generator.NewSnippetWriter(w, c, "$", "$")

	// START: rest.Scoper
	sw.Do("var _ $.Scoper|raw$ = (*$.TYPE|raw$)(nil)\n", args)
	sw.Do("func (*$.TYPE|raw$) NamespaceScoped() bool{\n", args)
	sw.Do("return $.$\n", res.Scope == string(apiext.NamespaceScoped))
	sw.Do("}\n\n", nil)
	// END: rest.Scoper

	// START: rest.KindProvider
	sw.Do("var _ $.KindProvider|raw$ = (*$.TYPE|raw$)(nil)\n", args)
	sw.Do("func (*$.TYPE|raw$) Kind() string{\n", args)
	sw.Do("return \"$.$\"", res.Kind)
	sw.Do("}\n\n", nil)
	// END: rest.KindProvider

	// START: rest.SingularNameProvider
	sw.Do("var _ $.SingularNameProvider|raw$ = (*$.TYPE|raw$)(nil)\n", args)
	sw.Do("func (*$.TYPE|raw$) GetSingularName() string{\n", args)
	sw.Do("return \"$.$\"\n", res.Singular)
	sw.Do("}\n\n", nil)
	// END: rest.SingularNameProvider

	// START: rest.ShortNamesProvider
	sw.Do("var _ $.ShortNamesProvider|raw$ = (*$.TYPE|raw$)(nil)\n", args)
	sw.Do("func (*$.TYPE|raw$) ShortNames() []string{\n", args)
	if len(res.ShortNames) != 0 {
		sw.Do("return []string{\n\"$.$\",\n}\n", strings.Join(res.ShortNames, "\",\n \""))
	} else {
		sw.Do("return []string{}\n", nil)
	}
	sw.Do("}\n\n", nil)
	// END: rest.ShortNamesProvider

	// START: rest.CategoriesProvider
	sw.Do("var _ $.CategoriesProvider|raw$ = (*$.TYPE|raw$)(nil)\n", args)
	sw.Do("func (*$.TYPE|raw$) Categories() []string{\n", args)
	if len(res.Categories) != 0 {
		sw.Do("return []string{\n\"$.$\",\n}\n", strings.Join(res.Categories, "\",\n \""))
	} else {
		sw.Do("return []string{}\n", nil)
	}
	sw.Do("}\n\n", nil)
	// END: rest.CategoriesProvider

	// START: subresources
	for _, v := range slices.Compact(res.SubResources) {
		switch v {
		case "status":
			sw.Do("var _ WithStatusSubResource = (*$.TYPE|raw$)(nil)\n", args)
			if t.Methods == nil || t.Methods["CopyStatusTo"] == nil {
				sw.Do("func (in *$.TYPE|raw$) CopyStatusTo(out $.Object|raw$) {\n", args)
				sw.Do("out.(*$.TYPE|raw$).Status=in.Status\n", args)
				sw.Do("}\n\n", nil)
			}
		case "scale":
			sw.Do("var _ WithScaleSubResource = (*$.TYPE|raw$)(nil)\n", args)
			if t.Methods == nil || t.Methods["GetScale"] == nil {
				sw.Do("func (in *$.TYPE|raw$) GetScale() *$.Scale|raw$ {\n", args)
				sw.Do("// TODO: Move me to definition file\n", nil)
				sw.Do("return &$.Scale|raw${}\n", args)
				sw.Do("}\n\n", nil)
			}
			if t.Methods == nil || t.Methods["SetScale"] == nil {
				sw.Do("func (in *$.TYPE|raw$) SetScale(s *$.Scale|raw$) {\n", args)
				sw.Do("// TODO: Move me to definition file \n", nil)
				sw.Do("}\n\n", nil)
			}
		}
	}
	// END: subresources

	return sw.Error()
}

const (
	pkgUtilPtr                   = "k8s.io/utils/ptr"
	pkgUtilBytes                 = "bytes"
	pkgAPIMachineryMetadata      = "k8s.io/apimachinery/pkg/apis/meta/v1"
	pkgAPIMachineryRuntimeSchema = "k8s.io/apimachinery/pkg/runtime/schema"
	pkgAPIMachineryRuntime       = "k8s.io/apimachinery/pkg/runtime"
	pkgTypeAPIRegistration       = "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
	pkgTypeAutoscaling           = "k8s.io/api/autoscaling/v1"
	pkgAPIServerRest             = "k8s.io/apiserver/pkg/registry/rest"
)

func getAPIServiceTypedArgs() generator.Args {
	return generator.Args{
		"BytesClone":         types.Ref(pkgUtilBytes, "Clone"),
		"TypeMeta":           types.Ref(pkgAPIMachineryMetadata, "TypeMeta"),
		"ObjectMeta":         types.Ref(pkgAPIMachineryMetadata, "ObjectMeta"),
		"ObjectMetaAccessor": types.Ref(pkgAPIMachineryMetadata, "ObjectMetaAccessor"),
		"Object":             types.Ref(pkgAPIMachineryRuntime, "Object"),
		"APIService":         types.Ref(pkgTypeAPIRegistration, "APIService"),
		"APIServiceSpec":     types.Ref(pkgTypeAPIRegistration, "APIServiceSpec"),
		"ServiceReference":   types.Ref(pkgTypeAPIRegistration, "ServiceReference"),
		"Scale":              types.Ref(pkgTypeAutoscaling, "Scale"),
	}
}

func getAPIResourceTypedArgs() generator.Args {
	return generator.Args{
		"PtrTo":                types.Ref(pkgUtilPtr, "To"),
		"TypeMeta":             types.Ref(pkgAPIMachineryMetadata, "TypeMeta"),
		"ObjectMeta":           types.Ref(pkgAPIMachineryMetadata, "ObjectMeta"),
		"ObjectMetaAccessor":   types.Ref(pkgAPIMachineryMetadata, "ObjectMetaAccessor"),
		"GroupVersionResource": types.Ref(pkgAPIMachineryRuntimeSchema, "GroupVersionResource"),
		"Object":               types.Ref(pkgAPIMachineryRuntime, "Object"),
		"Scale":                types.Ref(pkgTypeAutoscaling, "Scale"),
		"Scoper":               types.Ref(pkgAPIServerRest, "Scoper"),
		"KindProvider":         types.Ref(pkgAPIServerRest, "KindProvider"),
		"SingularNameProvider": types.Ref(pkgAPIServerRest, "SingularNameProvider"),
		"ShortNamesProvider":   types.Ref(pkgAPIServerRest, "ShortNamesProvider"),
		"CategoriesProvider":   types.Ref(pkgAPIServerRest, "CategoriesProvider"),
	}
}
