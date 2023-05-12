package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Owner struct {
	schema
}

func (Owner) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("ownerGroup", "ownerName"),
	}
}

func (Owner) Fields() []ent.Field {
	// Keep the json tag in camel case.
	return []ent.Field{
		field.String("ownerGroup").
			Comment("Describe group of the owner.").
			Default("default"),
		field.String("ownerName").
			Comment("Describe name of the owner.").
			Default("admin"),
	}
}

type OwnerOrg struct {
	schema
}

func (OwnerOrg) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("ownerGroup"),
	}
}

func (OwnerOrg) Fields() []ent.Field {
	// Keep the json tag in camel case.
	return []ent.Field{
		field.String("ownerGroup").
			Comment("Describe group of the owner.").
			Default("default"),
	}
}
