package entx

import (
	"embed"
	"reflect"

	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"

	"github.com/seal-io/walrus/pkg/dao/entx/annotation"
)

//go:embed template/*
var templateDir embed.FS

func loadTemplate() *gen.Template {
	p := []string{
		"template/*/*.tmpl",
		"template/*/*/*/*.tmpl",
	}

	fn := map[string]any{
		"genUpdater":        genUpdater,
		"genUpdaterClearer": genUpdaterClearer,
		"getIndexFields":    getIndexFields,
		"isPointer":         isPointer,
		"hasEqual":          hasEqual,
		"hasIsZero":         hasIsZero,
	}

	return gen.MustParse(gen.NewTemplate("entio").
		Funcs(fn).
		ParseFS(templateDir, p...))
}

// genUpdater is a template function,
// it returns true if the given value has mutable fields.
func genUpdater(v any) bool {
	switch t := v.(type) {
	case gen.Type:
		for i := range t.Fields {
			if !t.Fields[i].Immutable {
				return true
			}
		}
	case *gen.Type:
		for i := range t.Fields {
			if !t.Fields[i].Immutable {
				return true
			}
		}
	case gen.Field:
		return !t.Immutable
	case *gen.Field:
		return !t.Immutable
	}

	return false
}

// genUpdaterClearer is a template function,
// it returns true if the given value has skipped generating clearer.
func genUpdaterClearer(v any) bool {
	switch t := v.(type) {
	case gen.Field:
		a := annotation.MustExtractAnnotation(t.Annotations)
		return !a.SkipClearing && !a.SkipInput.Update
	case *gen.Field:
		a := annotation.MustExtractAnnotation(t.Annotations)
		return !a.SkipClearing && !a.SkipInput.Update
	case *gen.Type:
		a := annotation.MustExtractAnnotation(t.Annotations)
		return !a.SkipClearing
	case gen.Type:
		a := annotation.MustExtractAnnotation(t.Annotations)
		return !a.SkipClearing
	}

	return false
}

// getIndexFields is a template function,
// it returns the fields of the longest unique index.
func getIndexFields(v any) []*gen.Field {
	var n *gen.Type

	switch t := v.(type) {
	default:
		return nil
	case gen.Type:
		n = &t
	case *gen.Type:
		n = t
	}

	// Append fields of the longest unique index.
	var indexColumns []string

	for _, i := range n.Indexes {
		if !i.Unique {
			continue
		}

		if len(i.Columns) > len(indexColumns) {
			indexColumns = i.Columns
		}
	}

	if len(indexColumns) == 0 {
		return nil
	}

	fm := map[string]*gen.Field{}
	for i := range n.Fields {
		fm[n.Fields[i].Name] = n.Fields[i]
	}

	fs := make([]*gen.Field, 0, len(indexColumns))

	for _, col := range indexColumns {
		f, exist := fm[col]
		if !exist {
			continue
		}

		if !f.Immutable ||
			(!f.IsInt() && !f.IsInt64() && !f.IsBool() && !f.IsString()) {
			return nil
		}

		fs = append(fs, f)
	}

	return fs
}

// isPointer is a template function,
// it returns true if the given value is a pointer.
func isPointer(v any) bool {
	var rt *field.RType

	switch t := v.(type) {
	default:
		return false
	case gen.Field:
		rt = t.Type.RType
	case *gen.Field:
		rt = t.Type.RType
	}

	return rt.IsPtr()
}

// hasIsZero is a template function,
// it returns true if the given value has implemented `IsZero() bool` function.
func hasIsZero(v any) bool {
	var rt *field.RType

	switch t := v.(type) {
	default:
		return false
	case gen.Field:
		rt = t.Type.RType
	case *gen.Field:
		rt = t.Type.RType
	}

	if rt != nil {
		f, exist := rt.Methods["IsZero"]

		return exist &&
			len(f.In) == 0 &&
			len(f.Out) == 1 && f.Out[0].Kind == reflect.Bool
	}

	return false
}

// hasEqual is a template function,
// it returns true if the given value has implemented `Equal(other T) bool` function.
func hasEqual(v any) bool {
	var rt *field.RType

	switch t := v.(type) {
	default:
		return false
	case gen.Field:
		rt = t.Type.RType
	case *gen.Field:
		rt = t.Type.RType
	}

	if rt != nil {
		f, exist := rt.Methods["Equal"]

		return exist &&
			len(f.In) == 1 && f.In[0].String() == rt.String() &&
			len(f.Out) == 1 && f.Out[0].Kind == reflect.Bool
	}

	return false
}
