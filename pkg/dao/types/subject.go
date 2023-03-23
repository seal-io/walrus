package types

import "sort"

func DefaultSubjectRoles() SubjectRoles {
	return make(SubjectRoles, 0)
}

type SubjectRoles []SubjectRole

func (in SubjectRoles) Len() int {
	return len(in)
}

var builtinNames = []string{"admin", "edit", "view"}

func (in SubjectRoles) Less(i, j int) bool {
	// 1. the less, the priority is higher.
	// 2. the roles of "system" domain have higher priority.
	// 3. with the same domain, order by name as the below rule:
	//    "admin" / "edit" / "view" / <others>...
	switch {
	case in[i].Domain == "system":
		return true
	case in[j].Domain == "system":
		return false
	case in[i].Domain == in[j].Domain:
		for k := range builtinNames {
			switch {
			case in[i].Name == builtinNames[k]:
				return true
			case in[j].Name == builtinNames[k]:
				return false
			}
		}
		return in[i].Name < in[j].Name
	}
	return in[i].Domain < in[j].Domain
}

func (in SubjectRoles) Swap(i, j int) {
	in[i], in[j] = in[j], in[i]
}

func (in SubjectRoles) Deduplicate() SubjectRoles {
	if len(in) == 0 {
		return in
	}

	type setEntry struct {
		v SubjectRole
		i int
	}
	var set = map[string]setEntry{}
	for i := 0; i < len(in); i++ {
		if in[i].IsZero() {
			continue
		}
		var k = in[i].String()
		if _, existed := set[k]; !existed {
			set[k] = setEntry{v: in[i], i: len(set)}
		}
	}
	var out = make(SubjectRoles, len(set))
	for _, o := range set {
		out[o.i] = o.v
	}
	return out
}

func (in SubjectRoles) Sort() SubjectRoles {
	sort.Sort(in)
	return in
}

func (in SubjectRoles) Clone() SubjectRoles {
	var out = make(SubjectRoles, 0, len(in))
	for i := 0; i < len(in); i++ {
		out = append(out, in[i])
	}
	return out
}

func (in SubjectRoles) Delete(rs ...SubjectRole) SubjectRoles {
	type setEntry struct{}
	var set = map[string]setEntry{}
	for i := 0; i < len(rs); i++ {
		set[rs[i].String()] = setEntry{}
	}
	var out = SubjectRoles{}
	for i := 0; i < len(in); i++ {
		if _, existed := set[in[i].String()]; existed {
			continue
		}
		out = append(out, in[i])
	}
	return out
}

type SubjectRole struct {
	// Domain specifies the domain of this role.
	Domain string `json:"domain"`
	// Name specifies the name of this role.
	Name string `json:"name"`
}

func (in SubjectRole) IsZero() bool {
	return in.Domain == "" || in.Name == ""
}

func (in SubjectRole) String() string {
	return in.Domain + "/" + in.Name
}
