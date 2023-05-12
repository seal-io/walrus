package types

import (
	"net/http"
	"sort"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"
)

func DefaultRolePolicies() RolePolicies {
	return make(RolePolicies, 0)
}

type RolePolicies []RolePolicy

func (in RolePolicies) Len() int {
	return len(in)
}

func (in RolePolicies) Less(i, j int) bool {
	return in[i].String() < in[j].String()
}

func (in RolePolicies) Swap(i, j int) {
	in[i], in[j] = in[j], in[i]
}

func (in RolePolicies) Deduplicate() RolePolicies {
	if len(in) == 0 {
		return in
	}

	type setEntry struct {
		v RolePolicy
		i int
	}
	set := map[string]setEntry{}

	for i := range in {
		if in[i].IsZero() {
			continue
		}

		k := in[i].String()
		if _, existed := set[k]; !existed {
			set[k] = setEntry{v: in[i], i: len(set)}
		}
	}

	out := make(RolePolicies, len(set))
	for _, o := range set {
		out[o.i] = o.v
	}

	return out
}

func (in RolePolicies) Sort() RolePolicies {
	sort.Sort(in)
	return in
}

func (in RolePolicies) Clone() RolePolicies {
	out := make(RolePolicies, 0, len(in))
	for i := 0; i < len(in); i++ {
		out = append(out, in[i])
	}

	return out
}

func (in RolePolicies) Delete(rs ...RolePolicy) RolePolicies {
	type setEntry struct{}

	set := map[string]setEntry{}
	for i := 0; i < len(rs); i++ {
		set[rs[i].String()] = setEntry{}
	}
	out := RolePolicies{}

	for i := 0; i < len(in); i++ {
		if _, existed := set[in[i].String()]; existed {
			continue
		}

		out = append(out, in[i])
	}

	return out
}

func (in RolePolicies) Normalize() RolePolicies {
	for i := 0; i < len(in); i++ {
		in[i] = in[i].Normalize()
	}

	return in
}

type RolePolicyResourceScope = string

const (
	RolePolicyResourceScopePrivate     RolePolicyResourceScope = "P"
	RolePolicyResourceScopeSubordinate RolePolicyResourceScope = "S"
	RolePolicyResourceScopeInherit     RolePolicyResourceScope = "I"
	RolePolicyResourceScopeGlobal      RolePolicyResourceScope = "G"
)

type RolePolicy struct {
	// Actions specifies the including action list of this policy.
	Actions []string `json:"actions,omitempty"`
	// ActionExcludes specifies the excluding action list of this policy,
	// only works if the Actions specified as ["*"].
	ActionExcludes []string `json:"actionExcludes,omitempty"`

	// Resources specifies the including action list of this policy.
	Resources []string `json:"resources,omitempty"`
	// ResourceExcludes specifies the excluding action list of this policy,
	// only works if the Resources specified as ["*"].
	ResourceExcludes []string `json:"resourceExcludes,omitempty"`
	// ObjectIDs specifies the including ID list of this policy,
	// only works if the Resources specified as ["<a kind of resource>"].
	ObjectIDs []string `json:"objectIDs,omitempty"`
	// ObjectIDExcludes specifies the excluding action list of this policy,
	// only works if the Resources specified as ["<a kind of resource>"]
	// and ObjectIDs specified as ["*"].
	ObjectIDExcludes []string `json:"objectIDExcludes,omitempty"`
	// Scope specifies the scope of the resource is granted for this policy,
	// supports configuring with "P"(private), "S"(subordinate), "I"(inherit) and "G"(global),
	// and must configure as "G" if specify ObjectIDs and ObjectIDExcludes.
	Scope RolePolicyResourceScope `json:"scope,omitempty"`

	// Paths specifies the route-registered path list of this policy,
	// i.e. /resources/:id, only works if the Resources has not been specified.
	Paths []string `json:"paths,omitempty"`
}

func (in RolePolicy) Normalize() RolePolicy {
	if len(in.Actions) > 1 {
		// Aggregate Actions if Actions has "*".
		in.Actions = aggregateList(&in.Actions)
	}

	if len(in.Actions) != 1 || in.Actions[0] != "*" {
		// Clean up ActionExcludes if Action is not ["*"].
		in.ActionExcludes = nil
	}

	if len(in.Resources) > 1 {
		// Aggregate Resources if Resources has "*".
		in.Resources = aggregateList(&in.Resources)
	}

	if len(in.Resources) != 1 || in.Resources[0] != "*" {
		// Clean up ResourceExcludes if Resources is not ["*"].
		in.ResourceExcludes = nil
	}

	if len(in.Resources) != 1 || in.Resources[0] == "*" {
		// Clean up ObjectIDs if Resources is not ["<a kind of resource>"].
		in.ObjectIDs = nil
	}

	if len(in.ObjectIDs) > 1 {
		// Aggregate ObjectIDs if ObjectIDs has "*".
		in.ObjectIDs = aggregateList(&in.ObjectIDs)
	}

	if len(in.ObjectIDs) != 1 || in.ObjectIDs[0] != "*" {
		// Clean up ObjectIDExcludes if ObjectIDs is not ["*"].
		in.ObjectIDExcludes = nil
	}

	if len(in.ObjectIDs) != 0 {
		// Correct Scope to "G" if ObjectIDs is not empty.
		in.Scope = RolePolicyResourceScopeGlobal
	}

	if len(in.Resources) != 0 {
		// Clean up Paths if Resources is not empty.
		in.Paths = nil
	}

	if len(in.Paths) != 0 {
		// Clean up Scope if Paths is not empty.
		in.Scope = ""
	}

	return in
}

// aggregateList returns a slice including only one asterisk if it has asterisk item,
// e.g. [x, y, *] => [*].
func aggregateList(l *[]string) []string {
	if len(*l) < 2 {
		return *l
	}

	s := sets.NewString(*l...)
	if s.Has("*") {
		return []string{"*"}
	}

	return s.List()
}

func (in RolePolicy) IsZero() bool {
	if len(in.Actions) == 0 {
		return true
	}

	if len(in.Resources) == 0 && len(in.Paths) == 0 {
		return true
	}

	return false
}

func (in RolePolicy) String() string {
	appendAttributes := func(sb *strings.Builder, prefix string, ss []string) {
		s := len(ss)
		if s == 0 {
			return
		}

		if s > 1 {
			sort.Strings(ss)
		}

		sb.WriteString(prefix)
		sb.WriteString(": ")

		for i := range ss {
			sb.WriteString(ss[i])

			if i < s-1 {
				sb.WriteString(", ")
			}
		}

		sb.WriteString(";")
	}

	var sb strings.Builder

	appendAttributes(&sb, "actions", in.Actions)
	appendAttributes(&sb, "actionExcludes", in.ActionExcludes)
	appendAttributes(&sb, "resources", in.Resources)
	appendAttributes(&sb, "resourceExcludes", in.ResourceExcludes)
	appendAttributes(&sb, "objectIDs", in.ObjectIDs)
	appendAttributes(&sb, "objectIDExcludes", in.ObjectIDExcludes)
	appendAttributes(&sb, "scope", []string{in.Scope})
	appendAttributes(&sb, "paths", in.Paths)

	return sb.String()
}

func RolePolicyFields(f ...string) []string {
	return f
}

func RolePolicyResourceAdminFor(resources ...string) RolePolicy {
	return RolePolicy{
		Actions:   RolePolicyFields("*"),
		Resources: resources,
		Scope:     RolePolicyResourceScopeGlobal,
	}
}

func RolePolicyResourceEditFor(resources ...string) RolePolicy {
	return RolePolicy{
		Actions:   RolePolicyFields("*"),
		Resources: resources,
		Scope:     RolePolicyResourceScopeSubordinate,
	}
}

func RolePolicyResourceViewFor(resources ...string) RolePolicy {
	return RolePolicy{
		Actions:   RolePolicyFields(http.MethodGet),
		Resources: resources,
		Scope:     RolePolicyResourceScopeSubordinate,
	}
}
