package crypto

import (
	"errors"
	"regexp"

	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// stringBuilder is the builder for string fields.
type stringBuilder struct {
	desc *field.Descriptor
}

// Default sets the default value of the field.
func (b *stringBuilder) Default(v String) *stringBuilder {
	b.desc.Default = v
	return b
}

// Nillable indicates that this field is a nillable.
// Unlike "Optional" only fields, "Nillable" fields are pointers in the generated field.
func (b *stringBuilder) Nillable() *stringBuilder {
	b.desc.Nillable = true
	return b
}

// Optional indicates that this field is optional on create.
// Unlike edges, fields are required by default.
func (b *stringBuilder) Optional() *stringBuilder {
	b.desc.Optional = true
	return b
}

// Sensitive fields not printable and not serializable.
func (b *stringBuilder) Sensitive() *stringBuilder {
	b.desc.Sensitive = true
	return b
}

// Unique makes the field unique within all vertices of this type.
// Only supported in PostgreSQL.
func (b *stringBuilder) Unique() *stringBuilder {
	b.desc.Unique = true
	return b
}

// Immutable indicates that this field cannot be updated.
func (b *stringBuilder) Immutable() *stringBuilder {
	b.desc.Immutable = true
	return b
}

// Comment sets the comment of the field.
func (b *stringBuilder) Comment(c string) *stringBuilder {
	b.desc.Comment = c
	return b
}

// StructTag sets the struct tag of the field.
func (b *stringBuilder) StructTag(s string) *stringBuilder {
	b.desc.Tag = s
	return b
}

// StorageKey sets the storage key of the field.
// In SQL dialects is the column name and Gremlin is the property.
func (b *stringBuilder) StorageKey(key string) *stringBuilder {
	b.desc.StorageKey = key
	return b
}

// Annotations adds a list of annotations to the field object to be used by
// codegen extensions.
func (b *stringBuilder) Annotations(annotations ...schema.Annotation) *stringBuilder {
	b.desc.Annotations = append(b.desc.Annotations, annotations...)
	return b
}

// Descriptor implements the ent.Field interface by returning its descriptor.
func (b *stringBuilder) Descriptor() *field.Descriptor {
	return b.desc
}

// MaxLen sets the max-length of the bytes type in the database.
// In MySQL, this affects the BLOB type (tiny 2^8-1, regular 2^16-1, medium 2^24-1, long 2^32-1).
// In SQLite, it does not have any effect on the type size, which is default to 1B bytes.
func (b *stringBuilder) MaxLen(i int) *stringBuilder {
	b.desc.Size = i
	b.desc.Validators = append(b.desc.Validators, func(v string) error {
		if len(v) > i {
			return errors.New("value is greater than the required length")
		}
		return nil
	})
	return b
}

// MinLen adds a length validator for this field.
// Operation fails if the length of the string is less than the given value.
func (b *stringBuilder) MinLen(i int) *stringBuilder {
	b.desc.Validators = append(b.desc.Validators, func(v string) error {
		if len(v) < i {
			return errors.New("value is less than the required length")
		}
		return nil
	})
	return b
}

// NotEmpty adds a length validator for this field.
// Operation fails if the length of the string is zero.
func (b *stringBuilder) NotEmpty() *stringBuilder {
	return b.MinLen(1)
}

// Validate adds a validator for this field.
// Operation fails if the validation fails.
func (b *stringBuilder) Validate(fn func(string) error) *stringBuilder {
	b.desc.Validators = append(b.desc.Validators, fn)
	return b
}

// Match adds a regex matcher for this field.
// Operation fails if the regex fails.
func (b *stringBuilder) Match(re *regexp.Regexp) *stringBuilder {
	b.desc.Validators = append(b.desc.Validators, func(v string) error {
		if !re.MatchString(v) {
			return errors.New("value does not match validation")
		}
		return nil
	})
	return b
}

// bytesBuilder is the builder for bytes fields.
type bytesBuilder struct {
	desc *field.Descriptor
}

// Default sets the default value of the field.
func (b *bytesBuilder) Default(v Bytes) *bytesBuilder {
	b.desc.Default = v
	return b
}

// Nillable indicates that this field is a nillable.
// Unlike "Optional" only fields, "Nillable" fields are pointers in the generated struct.
func (b *bytesBuilder) Nillable() *bytesBuilder {
	b.desc.Nillable = true
	return b
}

// Optional indicates that this field is optional on create.
// Unlike edges, fields are required by default.
func (b *bytesBuilder) Optional() *bytesBuilder {
	b.desc.Optional = true
	return b
}

// Sensitive fields not printable and not serializable.
func (b *bytesBuilder) Sensitive() *bytesBuilder {
	b.desc.Sensitive = true
	return b
}

// Unique makes the field unique within all vertices of this type.
// Only supported in PostgreSQL.
func (b *bytesBuilder) Unique() *bytesBuilder {
	b.desc.Unique = true
	return b
}

// Immutable indicates that this field cannot be updated.
func (b *bytesBuilder) Immutable() *bytesBuilder {
	b.desc.Immutable = true
	return b
}

// Comment sets the comment of the field.
func (b *bytesBuilder) Comment(c string) *bytesBuilder {
	b.desc.Comment = c
	return b
}

// StructTag sets the struct tag of the field.
func (b *bytesBuilder) StructTag(s string) *bytesBuilder {
	b.desc.Tag = s
	return b
}

// StorageKey sets the storage key of the field.
// In SQL dialects is the column name and Gremlin is the property.
func (b *bytesBuilder) StorageKey(key string) *bytesBuilder {
	b.desc.StorageKey = key
	return b
}

// Annotations adds a list of annotations to the field object to be used by
// codegen extensions.
func (b *bytesBuilder) Annotations(annotations ...schema.Annotation) *bytesBuilder {
	b.desc.Annotations = append(b.desc.Annotations, annotations...)
	return b
}

// Descriptor implements the ent.Field interface by returning its descriptor.
func (b *bytesBuilder) Descriptor() *field.Descriptor {
	return b.desc
}

// MaxLen sets the max-length of the bytes type in the database.
// In MySQL, this affects the BLOB type (tiny 2^8-1, regular 2^16-1, medium 2^24-1, long 2^32-1).
// In SQLite, it does not have any effect on the type size, which is default to 1B bytes.
func (b *bytesBuilder) MaxLen(i int) *bytesBuilder {
	b.desc.Size = i
	b.desc.Validators = append(b.desc.Validators, func(buf []byte) error {
		if len(buf) > i {
			return errors.New("value is greater than the required length")
		}
		return nil
	})
	return b
}

// MinLen adds a length validator for this field.
// Operation fails if the length of the buffer is less than the given value.
func (b *bytesBuilder) MinLen(i int) *bytesBuilder {
	b.desc.Validators = append(b.desc.Validators, func(b []byte) error {
		if len(b) < i {
			return errors.New("value is less than the required length")
		}
		return nil
	})
	return b
}

// NotEmpty adds a length validator for this field.
// Operation fails if the length of the buffer is zero.
func (b *bytesBuilder) NotEmpty() *bytesBuilder {
	return b.MinLen(1)
}

// Validate adds a validator for this field.
// Operation fails if the validation fails.
func (b *bytesBuilder) Validate(fn func([]byte) error) *bytesBuilder {
	b.desc.Validators = append(b.desc.Validators, fn)
	return b
}

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
