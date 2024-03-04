package generators

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/seal-io/utils/stringx"
	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/gengo/types"
	"k8s.io/klog/v2"
	apireg "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"

	"github.com/seal-io/code-generator/utils"
)

// reflectPackage reflects the given package into a APIServiceDefinition,
// according to the given package.
func reflectPackage(p *types.Package) *APIServiceDefinition {
	if p == nil {
		return nil
	}

	logger := klog.Background().
		WithName("$").
		WithValues("gen", "apireg-gen", "package", p.Path)

	// Collect package markers.
	//
	// +groupName=
	// +versionName=
	// +k8s:apireg-gen:service:insecureSkipTLSVerify=,groupPriorityMinimum=,versionPriority=
	pm := map[string][]string{}
	collectMarkers(p.Comments, pm)

	if len(pm) == 0 || len(pm["group"]) == 0 || len(pm["version"]) == 0 {
		return nil
	}

	var (
		group   = pm["group"][len(pm["group"])-1]
		version = pm["version"][len(pm["version"])-1]
	)

	apisvc := apireg.APIService{
		TypeMeta: meta.TypeMeta{
			APIVersion: apireg.SchemeGroupVersion.String(),
			Kind:       "APIService",
		},
		ObjectMeta: meta.ObjectMeta{
			Name: fmt.Sprintf("%s.%s", version, group),
		},
		Spec: apireg.APIServiceSpec{
			Group:                 group,
			Version:               version,
			InsecureSkipTLSVerify: false,
			GroupPriorityMinimum:  100,
			VersionPriority:       100,
		},
	}

	for _, svc := range pm["service"] {
		logger := logger.WithValues("markers", "service")

		for mk, mv := range utils.ParseMarker(svc) {
			switch mk {
			case "insecureSkipTLSVerify":
				if err := json.Unmarshal([]byte(mv), &apisvc.Spec.InsecureSkipTLSVerify); err != nil {
					logger.Error(nil, "unmarshal insecureSkipTLSVerify", "value", mv)
				}
			case "groupPriorityMinimum":
				var v int32
				err := json.Unmarshal([]byte(mv), &v)
				switch {
				case err != nil:
					logger.Error(nil, "unmarshal groupPriorityMinimum", "value", mv)
				case v <= 0 || v >= 2000:
					logger.Error(nil, "invalid groupPriorityMinimum range, must between 0 and 2000", "value", mv)
				default:
					apisvc.Spec.GroupPriorityMinimum = v
				}
			case "versionPriority":
				var v int32
				err := json.Unmarshal([]byte(mv), &v)
				switch {
				case err != nil:
					logger.Error(nil, "unmarshal versionPriority", "value", mv)
				case v < 0:
					logger.Error(nil, "invalid versionPriority, must not be negative", "value", mv)
				default:
					apisvc.Spec.VersionPriority = v
				}
			}
		}
	}

	return &apisvc
}

var knownScopes = sets.New(string(apiext.NamespaceScoped), string(apiext.ClusterScoped))

// reflectType reflects the given type into a APIResourceDefinition,
// according to the given type.
func reflectType(t *types.Type) *APIResourceDefinition {
	if t == nil {
		return nil
	}

	logger := klog.Background().
		WithName("$").
		WithValues("gen", "apireg-gen", "type", t.String())

	var (
		kind     = t.Name.Name
		singular = strings.ToLower(kind)
		plural   = strings.ToLower(stringx.Pluralize(kind))
	)

	apires := APIResourceDefinition{
		Kind:     kind,
		Singular: singular,
	}

	// Collect type markers.
	//
	// +k8s:apireg-gen:resource:scope=,categories=,shortName=,plural=,subResources=
	tm := map[string][]string{}
	collectMarkers(t.SecondClosestCommentLines, tm)
	collectMarkers(t.CommentLines, tm)

	if len(tm) == 0 || len(tm["resource"]) == 0 {
		return nil
	}

	for _, res := range tm["resource"] {
		logger := logger.WithValues("markers", "resource")

		for mk, mv := range utils.ParseMarker(res) {
			switch mk {
			case "scope":
				var v string
				err := json.Unmarshal([]byte(mv), &v)
				switch {
				case err != nil:
					logger.Error(err, "unmarshal scope", "value", mv)
				case !knownScopes.Has(v):
					logger.Error(nil, "invalid scope, select from known scopes", "value", mv)
				default:
					apires.Scope = v
				}
			case "categories":
				if err := json.Unmarshal([]byte(mv), &apires.Categories); err != nil {
					logger.Error(nil, "unmarshal categories", "value", mv)
				}
			case "shortName":
				if err := json.Unmarshal([]byte(mv), &apires.ShortNames); err != nil {
					logger.Error(nil, "unmarshal shortName", "value", mv)
				}
			case "plural":
				if err := json.Unmarshal([]byte(mv), &apires.Plural); err != nil {
					logger.Error(nil, "unmarshal plural", "value", mv)
				}
			case "subResources":
				if err := json.Unmarshal([]byte(mv), &apires.SubResources); err != nil {
					logger.Error(nil, "unmarshal subResources", "value", mv)
				}
			}
		}
	}
	if apires.Scope == "" {
		apires.Scope = string(apiext.NamespaceScoped)
	}
	if apires.Plural == "" {
		apires.Plural = plural
	}

	return &apires
}

const (
	groupMarker     = "+groupName="
	versionMarker   = "+versionName="
	apiregGenMarker = "+k8s:apireg-gen:"
)

// collectMarkers collects markers from the given comments into a map.
func collectMarkers(comments []string, into map[string][]string) {
	for _, line := range comments {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		switch {
		default:
			if !strings.HasPrefix(line, "+") {
				into["comment"] = append(into["comment"], line)
			}
		case strings.HasPrefix(line, groupMarker):
			if v := line[len(groupMarker):]; v != "" {
				into["group"] = append(into["group"], v)
			}
		case strings.HasPrefix(line, versionMarker):
			if v := line[len(versionMarker):]; v != "" {
				into["version"] = append(into["version"], v)
			}
		case strings.HasPrefix(line, apiregGenMarker):
			kv := strings.SplitN(line[len(apiregGenMarker):], ":", 2)
			if len(kv) == 2 {
				into[kv[0]] = append(into[kv[0]], kv[1])
			} else if len(kv) == 1 {
				into[kv[0]] = append(into[kv[0]], "")
			}
		}
	}
}
