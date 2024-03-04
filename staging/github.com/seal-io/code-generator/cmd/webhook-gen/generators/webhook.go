package generators

import (
	"fmt"
	"io"
	"strings"

	"github.com/seal-io/utils/stringx"
	admreg "k8s.io/api/admissionregistration/v1"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
)

func NewWebhookGen(sanitizedName, outputPackage string) generator.Generator {
	return &webhookGen{
		DefaultGen: generator.DefaultGen{
			OptionalName: sanitizedName,
		},
		outputPackage: outputPackage,
		imports:       generator.NewImportTrackerForPackage(outputPackage),
	}
}

type WebhookTypeDefinition struct {
	Validating *admreg.ValidatingWebhook `json:"validating,omitempty"`
	Mutating   *admreg.MutatingWebhook   `json:"mutating,omitempty"`
}

type webhookGen struct {
	generator.DefaultGen

	outputPackage string
	imports       namer.ImportTracker

	types    []*types.Type
	typeDefs map[types.Name]*WebhookTypeDefinition
}

func (g *webhookGen) Filter(c *generator.Context, t *types.Type) bool {
	if t.Kind == types.Struct && t.Name.Package == g.outputPackage {
		if td := reflectType(t); td != nil {
			if g.typeDefs == nil {
				g.typeDefs = map[types.Name]*WebhookTypeDefinition{}
			}

			g.types = append(g.types, t)
			g.typeDefs[t.Name] = td
		}
		return true
	}
	return false
}

func (g *webhookGen) Namers(c *generator.Context) namer.NameSystems {
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

func (g *webhookGen) isOtherPackage(pkg string) bool {
	switch {
	default:
		return true
	case pkg == g.outputPackage:
	case strings.HasSuffix(pkg, `"`+g.outputPackage+`"`):
	}
	return false
}

func (g *webhookGen) Imports(c *generator.Context) []string {
	var pkgs []string
	for _, pkg := range g.imports.ImportLines() {
		if g.isOtherPackage(pkg) {
			pkgs = append(pkgs, pkg)
		}
	}
	return pkgs
}

func (g *webhookGen) Init(c *generator.Context, w io.Writer) error {
	var validateTypes []*types.Type
	var mutateTypes []*types.Type
	for _, t := range g.types {
		td := g.typeDefs[t.Name]
		if td == nil {
			continue
		}

		if td.Validating != nil {
			validateTypes = append(validateTypes, t)
		}

		if td.Mutating != nil {
			mutateTypes = append(mutateTypes, t)
		}
	}

	args := getAdmissionRegistrationTypedArgs()
	sw := generator.NewSnippetWriter(w, c, "$", "$")

	sw.Do("func GetWebhookConfigurations(np string, c $.WebhookClientConfig|raw$) "+
		"(*$.ValidatingWebhookConfiguration|raw$, *$.MutatingWebhookConfiguration|raw$) {\n", args)
	sw.Do("vwc := GetValidatingWebhookConfiguration(np+\"-validation\", c)\n", nil)
	sw.Do("mwc := GetMutatingWebhookConfiguration(np+\"-mutation\", c)\n", nil)
	sw.Do("return vwc, mwc\n", nil)
	sw.Do("}\n\n", nil)

	// START: ValidatingWebhookConfiguration.
	sw.Do("func GetValidatingWebhookConfiguration(n string, c $.WebhookClientConfig|raw$) *$.ValidatingWebhookConfiguration|raw$ {\n", args)

	if len(validateTypes) == 0 {
		sw.Do("return nil\n", nil)
	} else {
		sw.Do("return &$.ValidatingWebhookConfiguration|raw${\n", args)
		sw.Do("TypeMeta: $.TypeMeta|raw${\n", args)
		sw.Do("APIVersion: \"admissionregistration.k8s.io/v1\",\n", nil)
		sw.Do("Kind: \"ValidatingWebhookConfiguration\",\n", nil)
		sw.Do("},\n", nil)
		sw.Do("ObjectMeta: $.ObjectMeta|raw${\n", args)
		sw.Do("Name: n,\n", nil)
		sw.Do("},\n", nil)
		sw.Do("Webhooks: []$.ValidatingWebhook|raw${\n", args)
		for _, t := range validateTypes {
			typedArgs := generator.Args{"TYPE": t}
			sw.Do("vwh_$.TYPE|private$(c),\n", typedArgs)
		}
		sw.Do("},\n", nil)
		sw.Do("}\n", nil)
	}

	sw.Do("}\n\n", nil)
	// END: ValidatingWebhookConfiguration.

	// START: MutatingWebhookConfiguration.
	sw.Do("func GetMutatingWebhookConfiguration(n string, c $.WebhookClientConfig|raw$) *$.MutatingWebhookConfiguration|raw$ {\n", args)

	if len(mutateTypes) == 0 {
		sw.Do("return nil\n", nil)
	} else {
		sw.Do("return &$.MutatingWebhookConfiguration|raw${\n", args)
		sw.Do("TypeMeta: $.TypeMeta|raw${\n", args)
		sw.Do("APIVersion: \"admissionregistration.k8s.io/v1\",\n", nil)
		sw.Do("Kind: \"MutatingWebhookConfiguration\",\n", nil)
		sw.Do("},\n", nil)
		sw.Do("ObjectMeta: $.ObjectMeta|raw${\n", args)
		sw.Do("Name: n,\n", nil)
		sw.Do("},\n", nil)
		sw.Do("Webhooks: []$.MutatingWebhook|raw${\n", args)
		for _, t := range mutateTypes {
			typedArgs := generator.Args{"TYPE": t}
			sw.Do("mwh_$.TYPE|private$(c),\n", typedArgs)
		}
		sw.Do("},\n", nil)
		sw.Do("}\n", nil)
	}

	sw.Do("}\n\n", nil)
	// END: MutatingWebhookConfiguration.

	return sw.Error()
}

func (g *webhookGen) GenerateType(c *generator.Context, t *types.Type, w io.Writer) error {
	td := g.typeDefs[t.Name]
	if td == nil {
		return nil
	}

	if td.Validating != nil {
		err := g.generateValidating(c, t, w, td.Validating)
		if err != nil {
			return err
		}
	}

	if td.Mutating != nil {
		err := g.generateMutating(c, t, w, td.Mutating)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *webhookGen) generateValidating(c *generator.Context, t *types.Type, w io.Writer, spec *admreg.ValidatingWebhook) error {
	path := fmt.Sprintf("/%s", strings.ReplaceAll(spec.Name, ".", "-"))

	args := getAdmissionRegistrationTypedArgs().With("TYPE", t)
	sw := generator.NewSnippetWriter(w, c, "$", "$")

	// ValidatePath
	sw.Do("func (*$.TYPE|raw$) ValidatePath() string {\n", args)
	sw.Do("return \"$.$\"\n", path)
	sw.Do("}\n\n", nil)

	// vwh_***
	sw.Do("func vwh_$.TYPE|private$(c $.WebhookClientConfig|raw$) $.ValidatingWebhook|raw$ {\n", args)

	sw.Do("path := \"$.$\"\n\n", path)
	sw.Do("cc := c.DeepCopy()\n", nil)
	sw.Do("if cc.Service != nil {\n", nil)
	sw.Do("cc.Service.Path = &path\n", nil)
	sw.Do("} else if c.URL != nil {\n", nil)
	sw.Do("cc.URL = ptr.To(*c.URL + path)\n", nil)
	sw.Do("}\n\n", nil)
	sw.Do("return $.ValidatingWebhook|raw${\n", args)

	args = args.With("SPEC", spec)
	// START: Webhook.
	sw.Do("Name: \"$.$\",\n", spec.Name)
	sw.Do("ClientConfig: *cc,\n", nil)

	rule := spec.Rules[0]
	args = args.With("RULE", rule)
	// START: Webhook/Rules/0.
	sw.Do("Rules: []$.RuleWithOperations|raw${\n", args)
	sw.Do("{\n", nil)
	sw.Do("Rule: $.Rule|raw${\n", args)
	if len(rule.APIGroups) != 0 {
		sw.Do("APIGroups: []string{\n\"$.$\",\n},\n", strings.Join(rule.APIGroups, "\",\n \""))
	}
	if len(rule.APIVersions) != 0 {
		sw.Do("APIVersions: []string{\n\"$.$\",\n},\n", strings.Join(rule.APIVersions, "\",\n \""))
	}
	sw.Do("Resources: []string{\n\"$.$\",\n},\n", strings.Join(rule.Resources, "\",\n \""))
	if rule.Scope != nil {
		sw.Do("Scope: ptr.To[$.ScopeType|raw$](\"$.RULE.Scope$\"),\n", args)
	}
	sw.Do("},\n", nil)
	if len(rule.Operations) != 0 {
		sw.Do("Operations: []$.OperationType|raw${\n\"$.OPERS$\",\n},\n",
			args.With("OPERS", stringx.Join("\",\n \"", rule.Operations...)))
	}
	sw.Do("},\n", nil)
	sw.Do("},\n", nil)
	// END: Webhook/Rules/0.

	if spec.FailurePolicy != nil {
		sw.Do("FailurePolicy: ptr.To[$.FailurePolicyType|raw$](\"$.SPEC.FailurePolicy$\"),\n", args)
	}
	if spec.MatchPolicy != nil {
		sw.Do("MatchPolicy: ptr.To[$.MatchPolicyType|raw$](\"$.SPEC.MatchPolicy$\"),\n", args)
	}
	if selector := spec.NamespaceSelector; selector != nil {
		sw.Do("NamespaceSelector: ptr.To($.LabelSelector|raw${\n", args)
		if len(selector.MatchLabels) != 0 {
			sw.Do("MatchLabels: map[string]string{\n", nil)
			for k, v := range selector.MatchLabels {
				sw.Do("\"$.KEY$\": \"$.VAL$\",\n", generator.Args{"KEY": k, "VAL": v})
			}
			sw.Do("},\n", nil)
		}
		if len(selector.MatchExpressions) != 0 {
			sw.Do("MatchExpressions: []$.LabelSelectorRequirement|raw${\n", args)
			for _, expr := range selector.MatchExpressions {
				sw.Do("{\n", nil)
				sw.Do("Key: \"$.KEY$\",\n", generator.Args{"KEY": expr.Key})
				sw.Do("Operator: $.LabelSelectorOperator|raw$(\"$.OPER$\"),\n", args.With("OPER", expr.Operator))
				if expr.Values != nil {
					sw.Do("Values: []string{\n\"$.$\",\n},\n", strings.Join(expr.Values, "\",\n \""))
				}
				sw.Do("},\n", nil)
			}
			sw.Do("},\n", nil)
		}
		sw.Do("}),\n", nil)
	}
	if selector := spec.ObjectSelector; selector != nil {
		sw.Do("ObjectSelector: ptr.To($.LabelSelector|raw${\n", args)
		if len(selector.MatchLabels) != 0 {
			sw.Do("MatchLabels: map[string]string{\n", nil)
			for k, v := range selector.MatchLabels {
				sw.Do("\"$.KEY$\": \"$.VAL$\",\n", generator.Args{"KEY": k, "VAL": v})
			}
			sw.Do("},\n", nil)
		}
		if len(selector.MatchExpressions) != 0 {
			sw.Do("MatchExpressions: []$.LabelSelectorRequirement|raw${\n", args)
			for _, expr := range selector.MatchExpressions {
				sw.Do("{\n", nil)
				sw.Do("Key: \"$.KEY$\",\n", generator.Args{"KEY": expr.Key})
				sw.Do("Operator: $.LabelSelectorOperator|raw$(\"$.OPER$\"),\n", args.With("OPER", expr.Operator))
				if expr.Values != nil {
					sw.Do("Values: []string{\n\"$.$\",\n},\n", strings.Join(expr.Values, "\",\n \""))
				}
				sw.Do("},\n", nil)
			}
			sw.Do("},\n", nil)
		}
		sw.Do("}),\n", nil)
	}
	if spec.SideEffects != nil {
		sw.Do("SideEffects: ptr.To[$.SideEffectClass|raw$](\"$.SPEC.SideEffects$\"),\n", args)
	}
	if spec.TimeoutSeconds != nil {
		sw.Do("TimeoutSeconds: ptr.To[int32]($.SPEC.TimeoutSeconds$),\n", args)
	}
	if len(spec.AdmissionReviewVersions) != 0 {
		sw.Do("AdmissionReviewVersions: []string{\n\"$.$\",\n},\n", strings.Join(spec.AdmissionReviewVersions, "\",\n \""))
	}
	if len(spec.MatchConditions) != 0 {
		sw.Do("MatchConditions: []$.MatchCondition|raw${\n", args)
		for _, cond := range spec.MatchConditions {
			sw.Do("{\n", nil)
			sw.Do("Name: \"$.$\",\n", cond.Name)
			sw.Do("Expression: $.$,\n", fmt.Sprintf("%#v", cond.Expression))
			sw.Do("},\n", nil)
		}
		sw.Do("},\n", nil)
	}

	sw.Do("}\n", nil)
	// END: Webhook.

	sw.Do("}\n\n", nil)

	return sw.Error()
}

func (g *webhookGen) generateMutating(c *generator.Context, t *types.Type, w io.Writer, spec *admreg.MutatingWebhook) error {
	path := fmt.Sprintf("/%s", strings.ReplaceAll(spec.Name, ".", "-"))

	args := getAdmissionRegistrationTypedArgs().With("TYPE", t)
	sw := generator.NewSnippetWriter(w, c, "$", "$")

	// DefaultPath
	sw.Do("func (*$.TYPE|raw$) DefaultPath() string {\n", args)
	sw.Do("return \"$.$\"\n", path)
	sw.Do("}\n\n", nil)

	// mwh_***
	sw.Do("func mwh_$.TYPE|private$(c $.WebhookClientConfig|raw$) $.MutatingWebhook|raw$ {\n", args)

	sw.Do("path := \"$.$\"\n\n", path)
	sw.Do("cc := c.DeepCopy()\n", nil)
	sw.Do("if cc.Service != nil {\n", nil)
	sw.Do("cc.Service.Path = &path\n", nil)
	sw.Do("} else if c.URL != nil {\n", nil)
	sw.Do("cc.URL = ptr.To(*c.URL + path)\n", nil)
	sw.Do("}\n\n", nil)
	sw.Do("return $.MutatingWebhook|raw${\n", args)

	args = args.With("SPEC", spec)
	// START: Webhook.
	sw.Do("Name: \"$.$\",\n", spec.Name)
	sw.Do("ClientConfig: *cc,\n", nil)

	rule := spec.Rules[0]
	args = args.With("RULE", rule)
	// START: Webhook/Rules/0.
	sw.Do("Rules: []$.RuleWithOperations|raw${\n", args)
	sw.Do("{\n", nil)
	sw.Do("Rule: $.Rule|raw${\n", args)
	if len(rule.APIGroups) != 0 {
		sw.Do("APIGroups: []string{\n\"$.$\",\n},\n", strings.Join(rule.APIGroups, "\",\n \""))
	}
	if len(rule.APIVersions) != 0 {
		sw.Do("APIVersions: []string{\n\"$.$\",\n},\n", strings.Join(rule.APIVersions, "\",\n \""))
	}
	sw.Do("Resources: []string{\n\"$.$\",\n},\n", strings.Join(rule.Resources, "\",\n \""))
	if rule.Scope != nil {
		sw.Do("Scope: ptr.To[$.ScopeType|raw$](\"$.RULE.Scope$\"),\n", args)
	}
	sw.Do("},\n", nil)
	if len(rule.Operations) != 0 {
		sw.Do("Operations: []$.OperationType|raw${\n\"$.OPERS$\",\n},\n",
			args.With("OPERS", stringx.Join("\",\n \"", rule.Operations...)))
	}
	sw.Do("},\n", nil)
	sw.Do("},\n", nil)
	// END: Webhook/Rules/0.

	if spec.FailurePolicy != nil {
		sw.Do("FailurePolicy: ptr.To[$.FailurePolicyType|raw$](\"$.SPEC.FailurePolicy$\"),\n", args)
	}
	if spec.MatchPolicy != nil {
		sw.Do("MatchPolicy: ptr.To[$.MatchPolicyType|raw$](\"$.SPEC.MatchPolicy$\"),\n", args)
	}
	if selector := spec.NamespaceSelector; selector != nil {
		sw.Do("NamespaceSelector: ptr.To($.LabelSelector|raw${\n", args)
		if len(selector.MatchLabels) != 0 {
			sw.Do("MatchLabels: map[string]string{\n", nil)
			for k, v := range selector.MatchLabels {
				sw.Do("\"$.KEY$\": \"$.VAL$\",\n", generator.Args{"KEY": k, "VAL": v})
			}
			sw.Do("},\n", nil)
		}
		if len(selector.MatchExpressions) != 0 {
			sw.Do("MatchExpressions: []$.LabelSelectorRequirement|raw${\n", args)
			for _, expr := range selector.MatchExpressions {
				sw.Do("{\n", nil)
				sw.Do("Key: \"$.KEY$\",\n", generator.Args{"KEY": expr.Key})
				sw.Do("Operator: $.LabelSelectorOperator|raw$(\"$.OPER$\"),\n", args.With("OPER", expr.Operator))
				if expr.Values != nil {
					sw.Do("Values: []string{\n\"$.$\",\n},\n", strings.Join(expr.Values, "\",\n \""))
				}
				sw.Do("},\n", nil)
			}
			sw.Do("},\n", nil)
		}
		sw.Do("}),\n", nil)
	}
	if selector := spec.ObjectSelector; selector != nil {
		sw.Do("ObjectSelector: ptr.To($.LabelSelector|raw${\n", args)
		if len(selector.MatchLabels) != 0 {
			sw.Do("MatchLabels: map[string]string{\n", nil)
			for k, v := range selector.MatchLabels {
				sw.Do("\"$.KEY$\": \"$.VAL$\",\n", generator.Args{"KEY": k, "VAL": v})
			}
			sw.Do("},\n", nil)
		}
		if len(selector.MatchExpressions) != 0 {
			sw.Do("MatchExpressions: []$.LabelSelectorRequirement|raw${\n", nil)
			for _, expr := range selector.MatchExpressions {
				sw.Do("{\n", nil)
				sw.Do("Key: \"$.KEY$\",\n", generator.Args{"KEY": expr.Key})
				sw.Do("Operator: $.LabelSelectorOperator|raw$(\"$.OPER$\"),\n", args.With("OPER", expr.Operator))
				if expr.Values != nil {
					sw.Do("Values: []string{\n\"$.$\",\n},\n", strings.Join(expr.Values, "\",\n \""))
				}
				sw.Do("},\n", nil)
			}
			sw.Do("},\n", nil)
		}
		sw.Do("}),\n", nil)
	}
	if spec.SideEffects != nil {
		sw.Do("SideEffects: ptr.To[$.SideEffectClass|raw$](\"$.SPEC.SideEffects$\"),\n", args)
	}
	if spec.TimeoutSeconds != nil {
		sw.Do("TimeoutSeconds: ptr.To[int32]($.SPEC.TimeoutSeconds$),\n", args)
	}
	if len(spec.AdmissionReviewVersions) != 0 {
		sw.Do("AdmissionReviewVersions: []string{\n\"$.$\",\n},\n", strings.Join(spec.AdmissionReviewVersions, "\",\n \""))
	}
	if spec.ReinvocationPolicy != nil {
		sw.Do("ReinvocationPolicy: ptr.To[$.ReinvocationPolicyType|raw$](\"$.SPEC.ReinvocationPolicy$\"),\n", args)
	}
	if len(spec.MatchConditions) != 0 {
		sw.Do("MatchConditions: []$.MatchCondition|raw${\n", args)
		for _, cond := range spec.MatchConditions {
			sw.Do("{\n", nil)
			sw.Do("Name: \"$.$\",\n", cond.Name)
			sw.Do("Expression: $.$,\n", fmt.Sprintf("%#v", cond.Expression))
			sw.Do("},\n", nil)
		}
		sw.Do("},\n", nil)
	}

	sw.Do("}\n", nil)
	// END: Webhook.

	sw.Do("}\n\n", nil)

	return sw.Error()
}

const (
	pkgUtilPtr                   = "k8s.io/utils/ptr"
	pkgAPIMachineryMetadata      = "k8s.io/apimachinery/pkg/apis/meta/v1"
	pkgTypeAdmissionRegistration = "k8s.io/api/admissionregistration/v1"
)

func getAdmissionRegistrationTypedArgs() generator.Args {
	return generator.Args{
		"PtrTo":                          types.Ref(pkgUtilPtr, "To"),
		"TypeMeta":                       types.Ref(pkgAPIMachineryMetadata, "TypeMeta"),
		"ObjectMeta":                     types.Ref(pkgAPIMachineryMetadata, "ObjectMeta"),
		"LabelSelector":                  types.Ref(pkgAPIMachineryMetadata, "LabelSelector"),
		"LabelSelectorRequirement":       types.Ref(pkgAPIMachineryMetadata, "LabelSelectorRequirement"),
		"LabelSelectorOperator":          types.Ref(pkgAPIMachineryMetadata, "LabelSelectorOperator"),
		"ValidatingWebhookConfiguration": types.Ref(pkgTypeAdmissionRegistration, "ValidatingWebhookConfiguration"),
		"ValidatingWebhook":              types.Ref(pkgTypeAdmissionRegistration, "ValidatingWebhook"),
		"MutatingWebhookConfiguration":   types.Ref(pkgTypeAdmissionRegistration, "MutatingWebhookConfiguration"),
		"MutatingWebhook":                types.Ref(pkgTypeAdmissionRegistration, "MutatingWebhook"),
		"RuleWithOperations":             types.Ref(pkgTypeAdmissionRegistration, "RuleWithOperations"),
		"WebhookClientConfig":            types.Ref(pkgTypeAdmissionRegistration, "WebhookClientConfig"),
		"Rule":                           types.Ref(pkgTypeAdmissionRegistration, "Rule"),
		"ScopeType":                      types.Ref(pkgTypeAdmissionRegistration, "ScopeType"),
		"OperationType":                  types.Ref(pkgTypeAdmissionRegistration, "OperationType"),
		"FailurePolicyType":              types.Ref(pkgTypeAdmissionRegistration, "FailurePolicyType"),
		"MatchPolicyType":                types.Ref(pkgTypeAdmissionRegistration, "MatchPolicyType"),
		"SideEffectClass":                types.Ref(pkgTypeAdmissionRegistration, "SideEffectClass"),
		"ReinvocationPolicyType":         types.Ref(pkgTypeAdmissionRegistration, "ReinvocationPolicyType"),
		"MatchCondition":                 types.Ref(pkgTypeAdmissionRegistration, "MatchCondition"),
	}
}
