package entx

import (
	"github.com/seal-io/seal/pkg/dao/entx/annotation"
)

type (
	// Annotation alias the Annotation from entx/annotation package.
	Annotation = annotation.Annotation

	// InputOption defines the stereotype for generate *Input struct.
	InputOption struct {
		SkipInput func(*Annotation)
		Input     func(*Annotation)
	}
)

// WithCreate is an option works with SkipInput and Input.
//
// By default, all fields will generate into its *CreateInput struct,
// to skip one field, for example:
//
//	func (Todo) Fields() []ent.Field {
//	  return []ent.Field{
//	     field.String("health").
//	         Annotations(
//	            entx.SkipInput(entx.WithCreate())),
//	    }
//	}
//
// By default, all reversible edges will generate into its *CreateInput struct as well,
// to skip one reversible edge, for example:
//
//	func (Todo) Edges() []ent.Edge {
//	  return []ent.Edge{
//	     edge.From("groups", Group.Type).
//	         Ref("users").
//	         Through("groupUserRelationships", GroupUserRelationship.Type),
//	         Annotations(
//	            entx.SkipInput(entx.WithCreate())),
//	    }
//	}
//
// By default, all irreversible edges will not generate into its *CreateInput struct as well,
// to generate one irreversible edge, for example:
//
//	func (Todo) Edges() []ent.Edge {
//	  return []ent.Edge{
//	     edge.To("spouse", User.Type).
//	        Unique().
//	        Field("spouseID").
//	        Annotations(
//	           entx.SkipInput(entx.WithCreate())),
//	    }
//	}
func WithCreate() InputOption {
	return InputOption{
		SkipInput: func(a *Annotation) {
			a.SkipInput.Create = true
		},
	}
}

// WithQuery is an option works with SkipInput and Input.
//
// By default, all non-indexing fields don't be generated into *QueryInput struct,
// to generate one non-indexing field, for example:
//
//	func (Todo) Fields() []ent.Field {
//	  return []ent.Field{
//	     field.String("health").
//	         Annotations(
//	            entx.Input(entx.WithQuery())),
//	    }
//	}
//
// By default, all fields of the longest index will generate into its *QueryInput struct,
// to skip one index, for example:
//
//	func (Todo) Indexes() []ent.Index {
//	  return []ent.Index{
//	     index.Fields("name").
//	         Annotations(
//	            entx.SkipInput(entx.WithQuery())),
//	    }
//	}
//
// By default, all reversible edges will generate into its *QueryInput struct as well,
// to skip one reversible edge, for example:
//
//	func (Todo) Edges() []ent.Edge {
//	  return []ent.Edge{
//	     edge.From("owner", User.Type).
//	        Required().
//	        Ref("pets").
//	        Field("ownerID").
//	        Annotations(
//	           entx.SkipInput(entx.WithQuery())),
//	    }
//	}
func WithQuery() InputOption {
	return InputOption{
		SkipInput: func(a *Annotation) {
			a.SkipInput.Query = true
		},
		Input: func(a *Annotation) {
			a.Input.Query = true
		},
	}
}

// WithUpdate is an option works with SkipInput and Input.
//
// By default, all immutable fields don't be generated into *UpdateInput struct,
// to generate one immutable field, for example:
//
//	func (Todo) Fields() []ent.Field {
//	  return []ent.Field{
//	     field.String("health").
//	         Immutable().
//	         Annotations(
//	            entx.Input(entx.WithUpdate())),
//	    }
//	}
//
// By default, all mutable fields will generate into its *UpdateInput struct,
// to skip one field, for example:
//
//	func (Todo) Fields() []ent.Field {
//	  return []ent.Field{
//	     field.String("health").
//	         Annotations(
//	            entx.SkipInput(entx.WithUpdate())),
//	    }
//	}
//
// By default, all mutable reversible edges will generate into its *UpdateInput struct as well,
// to skip one mutable reversible edge, for example:
//
//	func (Todo) Edges() []ent.Edge {
//	  return []ent.Edge{
//	     edge.From("owner", User.Type).
//	        Required().
//	        Ref("pets").
//	        Field("ownerID").
//	        Annotations(
//	           entx.SkipInput(entx.WithUpdate())),
//	    }
//	}
//
// By default, all mutable irreversible edges will not generate into its *UpdateInput struct as well,
// to generate one mutable irreversible edge, for example:
//
//	func (Todo) Edges() []ent.Edge {
//	  return []ent.Edge{
//	     edge.To("spouse", User.Type).
//	        Unique().
//	        Field("spouseID").
//	        Annotations(
//	           entx.SkipInput(entx.WithUpdate())),
//	    }
//	}
func WithUpdate() InputOption {
	return InputOption{
		SkipInput: func(a *Annotation) {
			a.SkipInput.Update = true
		},
		Input: func(a *Annotation) {
			a.Input.Update = true
		},
	}
}

// SkipInput skips generating the field or reversible edge into *Input struct if no options are provided,
// otherwise, it applies the given options.
func SkipInput(opts ...InputOption) (a Annotation) {
	if len(opts) == 0 {
		a.SkipInput.Query = true
		a.SkipInput.Create = true
		a.SkipInput.Update = true
	}

	for _, opt := range opts {
		opt.SkipInput(&a)
	}

	return
}

// Input generates the immutable field or edge into *UpdateInput/*QueryInput struct if no options are provided,
// otherwise, it applies the given options.
func Input(opts ...InputOption) (a Annotation) {
	if len(opts) == 0 {
		a.Input.Query = true
		a.Input.Update = true
	}

	for _, opt := range opts {
		opt.Input(&a)
	}

	return
}

// SkipOutput skips generating the field or edge into *Output struct.
func SkipOutput() (a Annotation) {
	return Annotation{
		SkipOutput: true,
	}
}

// SkipIO skips generating the field or edge into *Input and *Output struct.
func SkipIO() (a Annotation) {
	a.SkipInput.Query = true
	a.SkipInput.Create = true
	a.SkipInput.Update = true
	a.SkipOutput = true

	return
}

// SkipClearingOptionalField skips generating cleaning if the mutable field is optional at updating.
//
// If an optional mutable field annotates with SkipInput(WithUpdate()),
// it skips generating clearer as well.
//
// By default, all optional fields will be cleared if the given update value is zero,
// to keep the previous value, for example:
//
//	func (Todo) Fields() []ent.Field {
//	  return []ent.Field{
//	     field.String("health").
//	         Optional().
//	         Annotations(
//	            entx.SkipClearingOptionalField()),
//	    }
//	}
//
// To skip generating clearer for all optional mutable fields, you can declare to the struct annotation.
func SkipClearingOptionalField() Annotation {
	return Annotation{
		SkipClearing: true,
	}
}

// SkipStoringField treats the field as additional field,
// which is not stored in the database.
//
// By default, all fields will be stored in the database,
// to keep the field in the struct only, for example:
//
//	func (Todo) Fields() []ent.Field {
//	  return []ent.Field{
//	     field.String("health").
//	         Annotations(
//	            entx.SkipStoringField()),
//	    }
//	}
func SkipStoringField() Annotation {
	return Annotation{
		SkipStoring: true,
	}
}
