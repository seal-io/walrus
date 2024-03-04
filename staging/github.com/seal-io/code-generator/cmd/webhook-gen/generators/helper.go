package generators

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/seal-io/utils/stringx"
	admreg "k8s.io/api/admissionregistration/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/gengo/types"
	"k8s.io/klog/v2"

	"github.com/seal-io/code-generator/utils"
)

var (
	knownOperations = sets.New(
		admreg.OperationAll,
		admreg.Create,
		admreg.Delete,
		admreg.Update,
		admreg.Connect)
	knownScopes = sets.New(
		admreg.NamespacedScope,
		admreg.ClusterScope,
		admreg.AllScopes)
	knownFailurePolicies = sets.New(
		admreg.Fail,
		admreg.Ignore)
	knownMatchPolicies = sets.New(
		admreg.Exact,
		admreg.Equivalent)
	knownSideEffects = sets.New(
		admreg.SideEffectClassUnknown,
		admreg.SideEffectClassSome,
		admreg.SideEffectClassNone,
		admreg.SideEffectClassNoneOnDryRun)
	knownReinvocationPolicies = sets.New(
		admreg.NeverReinvocationPolicy,
		admreg.IfNeededReinvocationPolicy)
)

// reflectType reflects the type into a WebhookTypeDefinition,
// according to the given package.
func reflectType(t *types.Type) *WebhookTypeDefinition {
	if t == nil {
		return nil
	}

	logger := klog.Background().
		WithName("$").
		WithValues("gen", "webhook-gen", "type", t.String())

	// Collect type markers.
	//
	// +k8s:webhook-gen:validating:group=,version=,resource=,scope=,failurePolicy=,sideEffects=,matchPolicy=,timeoutSeconds=
	// +k8s:webhook-gen:validating:operations=
	// +k8s:webhook-gen:validating:namespaceSelector=
	// +k8s:webhook-gen:validating:objectSelector=
	// +k8s:webhook-gen:validating:matchConditions=
	// +k8s:webhook-gen:mutating:group=,version=,resource=,scope=,failurePolicy=,sideEffects=,matchPolicy=,reinvocationPolicy=,timeoutSeconds=
	// +k8s:webhook-gen:mutating:operations=
	// +k8s:webhook-gen:mutating:namespaceSelector=
	// +k8s:webhook-gen:mutating:objectSelector=
	// +k8s:webhook-gen:mutating:matchConditions=
	tm := map[string][]string{}
	collectMarkers(t.SecondClosestCommentLines, tm)
	collectMarkers(t.CommentLines, tm)

	if len(tm) == 0 || (len(tm["validating"]) == 0 && len(tm["mutating"]) == 0) {
		return nil
	}

	var td WebhookTypeDefinition

	if len(tm["validating"]) != 0 {
		logger := logger.WithValues("markers", "validating")

		wh := &admreg.ValidatingWebhook{
			AdmissionReviewVersions: []string{"v1"},
		}

		for _, webhook := range tm["validating"] {
			for mk, mv := range utils.ParseMarker(webhook) {
				switch mk {
				case "group":
					var v string
					if err := json.Unmarshal([]byte(mv), &v); err != nil {
						logger.Error(err, "unmarshal group", "value", mv)
					} else {
						if len(wh.Rules) == 0 {
							wh.Rules = append(wh.Rules, admreg.RuleWithOperations{})
						}
						wh.Rules[0].APIGroups = []string{strings.ToLower(v)}
					}
				case "version":
					var v string
					if err := json.Unmarshal([]byte(mv), &v); err != nil {
						logger.Error(err, "unmarshal version", "value", mv)
					} else {
						if len(wh.Rules) == 0 {
							wh.Rules = append(wh.Rules, admreg.RuleWithOperations{})
						}
						wh.Rules[0].APIVersions = []string{strings.ToLower(v)}
					}
				case "resource":
					var v string
					if err := json.Unmarshal([]byte(mv), &v); err != nil {
						logger.Error(err, "unmarshal resource", "value", mv)
					} else {
						if len(wh.Rules) == 0 {
							wh.Rules = append(wh.Rules, admreg.RuleWithOperations{})
						}
						wh.Rules[0].Resources = []string{strings.ToLower(v)}
					}
				case "scope":
					var v admreg.ScopeType
					err := json.Unmarshal([]byte(mv), &v)
					switch {
					case err != nil:
						logger.Error(err, "unmarshal scope", "value", mv)
					case !knownScopes.Has(v):
						logger.Error(nil, "invalid scope, select from known scopes", "value", mv)
					default:
						if len(wh.Rules) == 0 {
							wh.Rules = append(wh.Rules, admreg.RuleWithOperations{})
						}
						wh.Rules[0].Scope = &v
					}
				case "operations":
					var v []admreg.OperationType
					err := json.Unmarshal([]byte(mv), &v)
					switch {
					case err != nil:
						logger.Error(err, "unmarshal operations", "value", mv)
					case !knownOperations.HasAll(v...):
						logger.Error(nil, "invalid operations, select from known operations", "value", mv)
					default:
						if len(wh.Rules) == 0 {
							wh.Rules = append(wh.Rules, admreg.RuleWithOperations{})
						}
						wh.Rules[0].Operations = v
					}
				case "failurePolicy":
					var v admreg.FailurePolicyType
					err := json.Unmarshal([]byte(mv), &v)
					switch {
					case err != nil:
						logger.Error(err, "unmarshal failure policy", "value", mv)
					case !knownFailurePolicies.Has(v):
						logger.Error(nil, "invalid failure policy, select from known failure policies", "value", mv)
					default:
						wh.FailurePolicy = &v
					}
				case "sideEffects":
					var v admreg.SideEffectClass
					err := json.Unmarshal([]byte(mv), &v)
					switch {
					case err != nil:
						logger.Error(err, "unmarshal side effects", "value", mv)
					case !knownSideEffects.Has(v):
						logger.Error(nil, "invalid side effects, select from known side effect classes", "value", mv)
					default:
						wh.SideEffects = &v
					}
				case "matchPolicy":
					var v admreg.MatchPolicyType
					err := json.Unmarshal([]byte(mv), &v)
					switch {
					case err != nil:
						logger.Error(err, "unmarshal match policy", "value", mv)
					case !knownMatchPolicies.Has(v):
						logger.Error(nil, "invalid match policy, select from known match policies", "value", mv)
					default:
						wh.MatchPolicy = &v
					}
				case "timeoutSeconds":
					var v int32
					if err := json.Unmarshal([]byte(mv), &v); err != nil {
						logger.Error(err, "unmarshal timeout seconds", "value", mv)
					} else {
						wh.TimeoutSeconds = &v
					}
				case "namespaceSelector":
					var selector meta.LabelSelector
					if err := json.Unmarshal([]byte(mv), &selector); err != nil {
						logger.Error(err, "unmarshal namespace selector", "value", mv)
					} else {
						wh.NamespaceSelector = &selector
					}
				case "objectSelector":
					var selector meta.LabelSelector
					if err := json.Unmarshal([]byte(mv), &selector); err != nil {
						logger.Error(err, "unmarshal object selector", "value", mv)
					} else {
						wh.ObjectSelector = &selector
					}
				case "matchConditions":
					var conditions []admreg.MatchCondition
					err := json.Unmarshal([]byte(mv), &conditions)
					switch {
					case err != nil:
						logger.Error(err, "unmarshal match conditions", "value", mv)
					case len(conditions) > 64:
						logger.Error(nil, "too many match conditions, maximum 64", "value", mv)
					case len(conditions) > 0:
						wh.MatchConditions = conditions
					}
				}
			}
		}

		if len(wh.Rules) == 0 || len(wh.Rules[0].Resources) == 0 {
			logger.Error(nil, "invalid rules specified, must include one resource")
		} else {
			var pre, grp, ver, klow string
			pre = "validate"
			if len(wh.Rules[0].APIGroups) != 0 {
				grp = wh.Rules[0].APIGroups[0]
			}
			if len(wh.Rules[0].APIVersions) != 0 {
				ver = wh.Rules[0].APIVersions[0]
			}
			klow = strings.ToLower(stringx.Singularize(wh.Rules[0].Resources[0]))
			wh.Name = fmt.Sprintf("%s.%s.%s.%s", pre, grp, ver, klow)
			td.Validating = wh
		}
	}

	if len(tm["mutating"]) != 0 {
		logger := logger.WithValues("markers", "mutating")

		wh := &admreg.MutatingWebhook{
			AdmissionReviewVersions: []string{"v1"},
		}

		for _, webhook := range tm["mutating"] {
			for mk, mv := range utils.ParseMarker(webhook) {
				switch mk {
				case "group":
					var v string
					if err := json.Unmarshal([]byte(mv), &v); err != nil {
						logger.Error(err, "unmarshal group", "value", mv)
					} else {
						if len(wh.Rules) == 0 {
							wh.Rules = append(wh.Rules, admreg.RuleWithOperations{})
						}
						wh.Rules[0].APIGroups = []string{strings.ToLower(v)}
					}
				case "version":
					var v string
					if err := json.Unmarshal([]byte(mv), &v); err != nil {
						logger.Error(err, "unmarshal version", "value", mv)
					} else {
						if len(wh.Rules) == 0 {
							wh.Rules = append(wh.Rules, admreg.RuleWithOperations{})
						}
						wh.Rules[0].APIVersions = []string{strings.ToLower(v)}
					}
				case "resource":
					var v string
					if err := json.Unmarshal([]byte(mv), &v); err != nil {
						logger.Error(err, "unmarshal resource", "value", mv)
					} else {
						if len(wh.Rules) == 0 {
							wh.Rules = append(wh.Rules, admreg.RuleWithOperations{})
						}
						wh.Rules[0].Resources = []string{strings.ToLower(v)}
					}
				case "scope":
					var v admreg.ScopeType
					err := json.Unmarshal([]byte(mv), &v)
					switch {
					case err != nil:
						logger.Error(err, "unmarshal scope", "value", mv)
					case !knownScopes.Has(v):
						logger.Error(nil, "invalid scope, select from known scopes", "value", mv)
					default:
						if len(wh.Rules) == 0 {
							wh.Rules = append(wh.Rules, admreg.RuleWithOperations{})
						}
						wh.Rules[0].Scope = &v
					}
				case "operations":
					var v []admreg.OperationType
					err := json.Unmarshal([]byte(mv), &v)
					switch {
					case err != nil:
						logger.Error(err, "unmarshal operations", "value", mv)
					case !knownOperations.HasAll(v...):
						logger.Error(nil, "invalid operations, select from known operations", "value", mv)
					default:
						if len(wh.Rules) == 0 {
							wh.Rules = append(wh.Rules, admreg.RuleWithOperations{})
						}
						wh.Rules[0].Operations = v
					}
				case "failurePolicy":
					var v admreg.FailurePolicyType
					err := json.Unmarshal([]byte(mv), &v)
					switch {
					case err != nil:
						logger.Error(err, "unmarshal failure policy", "value", mv)
					case !knownFailurePolicies.Has(v):
						logger.Error(nil, "invalid failure policy, select from known failure policies", "value", mv)
					default:
						wh.FailurePolicy = &v
					}
				case "sideEffects":
					var v admreg.SideEffectClass
					err := json.Unmarshal([]byte(mv), &v)
					switch {
					case err != nil:
						logger.Error(err, "unmarshal side effects", "value", mv)
					case !knownSideEffects.Has(v):
						logger.Error(nil, "invalid side effects, select from known side effect classes", "value", mv)
					default:
						wh.SideEffects = &v
					}
				case "matchPolicy":
					var v admreg.MatchPolicyType
					err := json.Unmarshal([]byte(mv), &v)
					switch {
					case err != nil:
						logger.Error(err, "unmarshal match policy", "value", mv)
					case !knownMatchPolicies.Has(v):
						logger.Error(nil, "invalid match policy, select from known match policies", "value", mv)
					default:
						wh.MatchPolicy = &v
					}
				case "timeoutSeconds":
					var v int32
					if err := json.Unmarshal([]byte(mv), &v); err != nil {
						logger.Error(err, "unmarshal timeout seconds", "value", mv)
					} else {
						wh.TimeoutSeconds = &v
					}
				case "namespaceSelector":
					var selector meta.LabelSelector
					if err := json.Unmarshal([]byte(mv), &selector); err != nil {
						logger.Error(err, "unmarshal namespace selector", "value", mv)
					} else {
						wh.NamespaceSelector = &selector
					}
				case "objectSelector":
					var selector meta.LabelSelector
					if err := json.Unmarshal([]byte(mv), &selector); err != nil {
						logger.Error(err, "unmarshal object selector", "value", mv)
					} else {
						wh.ObjectSelector = &selector
					}
				case "reinvocationPolicy":
					var v admreg.ReinvocationPolicyType
					err := json.Unmarshal([]byte(mv), &v)
					switch {
					case err != nil:
						logger.Error(err, "unmarshal reinvocation policy", "value", mv)
					case !knownReinvocationPolicies.Has(v):
						logger.Error(nil, "invalid reinvocation policy, select from known reinvocation policies", "value", mv)
					default:
						wh.ReinvocationPolicy = &v
					}
				case "matchConditions":
					var conditions []admreg.MatchCondition
					err := json.Unmarshal([]byte(mv), &conditions)
					switch {
					case err != nil:
						logger.Error(err, "unmarshal match conditions", "value", mv)
					case len(conditions) > 64:
						logger.Error(nil, "too many match conditions, maximum 64", "value", mv)
					case len(conditions) > 0:
						wh.MatchConditions = conditions
					}
				}
			}
		}

		if len(wh.Rules) == 0 || len(wh.Rules[0].Resources) == 0 {
			logger.Error(nil, "invalid rules specified, must include one resource")
		} else {
			var pre, grp, ver, klow string
			pre = "mutate"
			if len(wh.Rules[0].APIGroups) != 0 {
				grp = wh.Rules[0].APIGroups[0]
			}
			if len(wh.Rules[0].APIVersions) != 0 {
				ver = wh.Rules[0].APIVersions[0]
			}
			klow = strings.ToLower(stringx.Singularize(wh.Rules[0].Resources[0]))
			wh.Name = fmt.Sprintf("%s.%s.%s.%s", pre, grp, ver, klow)
			td.Mutating = wh
		}
	}

	if td.Validating == nil && td.Mutating == nil {
		return nil
	}

	return &td
}

const (
	webhookGenMarker = "+k8s:webhook-gen:"
)

// collectMarkers collects markers from the given comment lines.
func collectMarkers(comments []string, into map[string][]string) {
	for _, line := range comments {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		if strings.HasPrefix(line, webhookGenMarker) {
			kv := strings.SplitN(line[len(webhookGenMarker):], ":", 2)
			if len(kv) == 2 {
				into[kv[0]] = append(into[kv[0]], kv[1])
			} else if len(kv) == 1 {
				into[kv[0]] = append(into[kv[0]], "")
			}
		}
	}
}
