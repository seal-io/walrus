package schema

import (
	"sort"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Subject struct {
	schema
}

func (Subject) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("kind", "group", "name").
			Unique(),
	}
}

func (Subject) Fields() []ent.Field {
	return []ent.Field{
		field.String("kind").
			Comment("The kind of the subject.").
			Immutable().
			Default("user"),
		field.String("group").
			Comment("The group of the subject.").
			Default("default"),
		field.String("name").
			Comment("The name of the subject.").
			Immutable().
			NotEmpty(),
		field.String("description").
			Comment("The detail of the subject.").
			Default(""),
		field.Bool("mountTo").
			Comment("Indicate whether the user mount to the group.").
			Nillable().
			Default(false),
		field.Bool("loginTo").
			Comment("Indicate whether the user login to the group.").
			Nillable().
			Default(true),
		field.JSON("roles", SubjectRoles{}).
			Comment("The role list of the subject.").
			Default(DefaultSubjectRoles()),
		field.Strings("paths").
			Comment("The path of the subject from the root group to itself.").
			Default([]string{}),
		field.Bool("builtin").
			Comment("Indicate whether the subject is builtin.").
			Default(false),
	}
}

func DefaultSubjectRoles() SubjectRoles {
	return make(SubjectRoles, 0)
}

type SubjectRoles []SubjectRole

func (in SubjectRoles) Len() int {
	return len(in)
}

func (in SubjectRoles) Less(i, j int) bool {
	if in[i].Domain == in[j].Domain {
		if in[i].Name == in[j].Name {
			return true
		}

		if in[i].Name == "admin" {
			return true
		} else if in[j].Name == "admin" {
			return false
		}

		if in[i].Name == "edit" {
			return true
		} else if in[j].Name == "edit" {
			return false
		}

		if in[i].Name == "view" {
			return true
		} else if in[j].Name == "view" {
			return false
		}

		return in[i].Name < in[j].Name
	}

	if in[i].Domain == "system" {
		return true
	} else if in[j].Domain == "system" {
		return false
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
