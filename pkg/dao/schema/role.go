package schema

import (
	"net/http"
	"sort"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"k8s.io/apimachinery/pkg/util/sets"
)

type Role struct {
	schema
}

func (Role) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("domain", "name").
			Unique(),
	}
}

func (Role) Fields() []ent.Field {
	return []ent.Field{
		field.String("domain").
			Comment("The domain of the role.").
			Immutable().
			Default("system"),
		field.String("name").
			Comment("The name of the role.").
			Immutable().
			NotEmpty(),
		field.String("description").
			Comment("The detail of the role.").
			Default(""),
		field.JSON("policies", RolePolicies{}).
			Comment("The policy list of the role.").
			Default(DefaultRolePolicies()),
		field.Bool("builtin").
			Comment("Indicate whether the subject is builtin, decide when creating.").
			Immutable().
			Default(false),
		field.Bool("session").
			Comment("Indicate whether the subject is session level, decide when creating.").
			Immutable().
			Default(false),
	}
}

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
	var set = map[string]setEntry{}
	for i := range in {
		if in[i].IsZero() {
			continue
		}
		var k = in[i].String()
		if _, existed := set[k]; !existed {
			set[k] = setEntry{v: in[i], i: len(set)}
		}
	}
	var out = make(RolePolicies, len(set))
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
	var out = make(RolePolicies, 0, len(in))
	for i := 0; i < len(in); i++ {
		out = append(out, in[i])
	}
	return out
}

func (in RolePolicies) Delete(rs ...RolePolicy) RolePolicies {
	type setEntry struct{}
	var set = map[string]setEntry{}
	for i := 0; i < len(rs); i++ {
		set[rs[i].String()] = setEntry{}
	}
	var out = RolePolicies{}
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

	// Scope specifies the scope of the resource is granted for this policy,
	// supports configuring with "P"(private), "S"(subordinate), "I"(inherit) and "G"(global).
	Scope RolePolicyResourceScope `json:"scope,omitempty"`
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

	// Paths specifies the route-registered path list of this policy,
	// i.e. /resources/:id, only works if the Resources has not been specified.
	Paths []string `json:"paths,omitempty"`
}

func (in RolePolicy) Normalize() RolePolicy {
	if len(in.Actions) > 1 {
		var actions = sets.NewString(in.Actions...)
		if actions.Has("*") {
			in.Actions = []string{"*"}
		} else {
			in.Actions = actions.List()
		}
	}
	if len(in.Actions) != 1 || in.Actions[0] != "*" {
		in.ActionExcludes = nil
	}

	if len(in.Resources) > 1 {
		var resources = sets.NewString(in.Resources...)
		if resources.Has("*") {
			in.Resources = []string{"*"}
		} else {
			in.Resources = resources.List()
		}
	}
	if len(in.Resources) != 1 || in.Resources[0] != "*" {
		in.ResourceExcludes = nil
	}

	if len(in.Resources) != 1 || in.Resources[0] == "*" {
		in.ObjectIDs = nil
	}
	if len(in.ObjectIDs) > 1 {
		var objectIDs = sets.NewString(in.ObjectIDs...)
		if objectIDs.Has("*") {
			in.ObjectIDs = []string{"*"}
		} else {
			in.ObjectIDs = objectIDs.List()
		}
	}
	if len(in.ObjectIDs) != 1 || in.ObjectIDs[0] != "*" {
		in.ObjectIDExcludes = nil
	}

	if len(in.Resources) != 0 {
		in.Paths = nil
	}

	if len(in.Paths) != 0 {
		in.Scope = ""
	}

	return in
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
	var appendAttributes = func(sb *strings.Builder, prefix string, ss []string) {
		var s = len(ss)
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
	appendAttributes(&sb, "scope", []string{in.Scope})
	appendAttributes(&sb, "resources", in.Resources)
	appendAttributes(&sb, "resourceExcludes", in.ResourceExcludes)
	appendAttributes(&sb, "objectIDs", in.ObjectIDs)
	appendAttributes(&sb, "objectIDExcludes", in.ObjectIDExcludes)
	appendAttributes(&sb, "nonResourceURLs", in.Paths)
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
