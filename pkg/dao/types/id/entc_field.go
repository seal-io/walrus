package id

import (
	"errors"
	"strings"

	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/seal/pkg/dao/types"
)

type fieldBuilder struct {
	desc *field.Descriptor
}

// Unique makes the field unique within all vertices of this type.
func (b *fieldBuilder) Unique() *fieldBuilder {
	b.desc.Unique = true
	return b
}

// NotEmpty adds a blank validator for this field.
// Operation fails if the length of the ID is blank.
func (b *fieldBuilder) NotEmpty() *fieldBuilder {
	b.desc.Validators = append(b.desc.Validators, func(v string) error {
		if strings.TrimSpace(v) == "" {
			return errors.New("value is blank")
		}
		return nil
	})
	return b
}

// Default sets the default value of the field.
func (b *fieldBuilder) Default(v types.ID) *fieldBuilder {
	b.desc.Default = v
	return b
}

// Nillable indicates that this field is a nillable.
// Unlike "Optional" only fields, "Nillable" fields are pointers in the generated field.
func (b *fieldBuilder) Nillable() *fieldBuilder {
	b.desc.Nillable = true
	return b
}

// Optional indicates that this field is optional on create.
// Unlike edges, fields are required by default.
func (b *fieldBuilder) Optional() *fieldBuilder {
	b.desc.Optional = true
	return b
}

// Immutable indicates that this field cannot be updated.
func (b *fieldBuilder) Immutable() *fieldBuilder {
	b.desc.Immutable = true
	return b
}

// Comment sets the comment of the field.
func (b *fieldBuilder) Comment(c string) *fieldBuilder {
	b.desc.Comment = c
	return b
}

// StructTag sets the struct tag of the field.
func (b *fieldBuilder) StructTag(s string) *fieldBuilder {
	b.desc.Tag = s
	return b
}

// StorageKey sets the storage key of the field.
// In SQL dialects is the column name and Gremlin is the property.
func (b *fieldBuilder) StorageKey(key string) *fieldBuilder {
	b.desc.StorageKey = key
	return b
}

// Annotations adds a list of annotations to the field object to be used by
// codegen extensions.
//
//	oid.Field("link").
//		Annotations(
//			entgql.OrderField("LINK"),
//		)
func (b *fieldBuilder) Annotations(annotations ...schema.Annotation) *fieldBuilder {
	b.desc.Annotations = append(b.desc.Annotations, annotations...)
	return b
}

// Descriptor implements the ent.Field interface by returning its descriptor.
func (b *fieldBuilder) Descriptor() *field.Descriptor {
	return b.desc
}
