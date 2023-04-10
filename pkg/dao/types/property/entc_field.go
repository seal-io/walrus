package property

import (
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// otherBuilder is the builder for other fields.
type otherBuilder struct {
	desc *field.Descriptor
}

// Default sets the default value of the field.
func (b *otherBuilder) Default(v any) *otherBuilder {
	b.desc.Default = v
	return b
}

// Nillable indicates that this field is a nillable.
// Unlike "Optional" only fields, "Nillable" fields are pointers in the generated struct.
func (b *otherBuilder) Nillable() *otherBuilder {
	b.desc.Nillable = true
	return b
}

// Optional indicates that this field is optional on create.
// Unlike edges, fields are required by default.
func (b *otherBuilder) Optional() *otherBuilder {
	b.desc.Optional = true
	return b
}

// Sensitive fields not printable and not serializable.
func (b *otherBuilder) Sensitive() *otherBuilder {
	b.desc.Sensitive = true
	return b
}

// Unique makes the field unique within all vertices of this type.
// Only supported in PostgreSQL.
func (b *otherBuilder) Unique() *otherBuilder {
	b.desc.Unique = true
	return b
}

// Immutable indicates that this field cannot be updated.
func (b *otherBuilder) Immutable() *otherBuilder {
	b.desc.Immutable = true
	return b
}

// Comment sets the comment of the field.
func (b *otherBuilder) Comment(c string) *otherBuilder {
	b.desc.Comment = c
	return b
}

// StructTag sets the struct tag of the field.
func (b *otherBuilder) StructTag(s string) *otherBuilder {
	b.desc.Tag = s
	return b
}

// StorageKey sets the storage key of the field.
// In SQL dialects is the column name and Gremlin is the property.
func (b *otherBuilder) StorageKey(key string) *otherBuilder {
	b.desc.StorageKey = key
	return b
}

// Annotations adds a list of annotations to the field object to be used by
// codegen extensions.
func (b *otherBuilder) Annotations(annotations ...schema.Annotation) *otherBuilder {
	b.desc.Annotations = append(b.desc.Annotations, annotations...)
	return b
}

// Descriptor implements the ent.Field interface by returning its descriptor.
func (b *otherBuilder) Descriptor() *field.Descriptor {
	return b.desc
}
