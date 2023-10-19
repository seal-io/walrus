package types

import (
	"sort"
	"strings"

	"golang.org/x/exp/slices"
	"k8s.io/apimachinery/pkg/util/sets"
)

const (
	RoleKindSystem  = "system"
	RoleKindProject = "project"
)

var roleKinds = []string{
	RoleKindSystem,
	RoleKindProject,
}

func RoleKinds() []string {
	return slices.Clone(roleKinds)
}

func IsRoleKind(s string) bool {
	return slices.Contains(roleKinds, s)
}

const (
	SystemRoleAnonymity = "system/anonymity"
	SystemRoleUser      = "system/user"
	SystemRoleManager   = "system/manager"
	SystemRoleAdmin     = "system/admin"

	ProjectRoleViewer = "project/viewer"
	ProjectRoleMember = "project/member"
	ProjectRoleOwner  = "project/owner"
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

type RolePolicy struct {
	// Actions specifies the including action list of this policy.
	Actions []string `json:"actions,omitempty"`

	// Resources specifies the including resource list of this policy.
	Resources []string `json:"resources,omitempty"`
	// ResourceRefers specifies the including object.Refer list of this policy,
	// only works if the Resources specified as ["<a kind of resource>"].
	ResourceRefers []string `json:"resourceRefers,omitempty"`

	// Paths specifies the route-registered path list of this policy,
	// i.e. /resources/:id, only works if the Resources has not been specified.
	Paths []string `json:"paths,omitempty"`
}

func (in RolePolicy) Normalize() RolePolicy {
	// Aggregate Actions if Actions has "*".
	if len(in.Actions) > 1 {
		in.Actions = aggregateList(&in.Actions)
	}

	// Aggregate Resources if Resources has "*".
	if len(in.Resources) > 1 {
		in.Resources = aggregateList(&in.Resources)
	}

	if len(in.Resources) != 0 {
		// Aggregate ResourceRefers if ResourceRefers has "*".
		if len(in.ResourceRefers) > 1 {
			in.ResourceRefers = aggregateList(&in.ResourceRefers)
		}
		// Clean up Paths if Resources is not empty.
		in.Paths = nil
	} else {
		// Clean up ResourceRefers if Resources is empty.
		in.ResourceRefers = nil
	}

	return in
}

// aggregateList returns a slice including only one asterisk if it has asterisk item,
// e.g. [x, y, *] => [*].
func aggregateList(l *[]string) []string {
	if len(*l) < 2 {
		return *l
	}

	ss := sets.NewString()

	for _, s := range *l {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}

		ss.Insert(s)
	}

	if ss.Has("*") {
		return []string{"*"}
	}

	return ss.List()
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
	appendAttributes(&sb, "resources", in.Resources)
	appendAttributes(&sb, "resourceRefers", in.ResourceRefers)
	appendAttributes(&sb, "paths", in.Paths)

	return sb.String()
}

func RolePolicyFields(f ...string) []string {
	return f
}
