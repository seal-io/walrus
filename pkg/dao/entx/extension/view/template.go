package view

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"entgo.io/ent/entc/gen"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/pkg/dao/entx/annotation"
	"github.com/seal-io/walrus/utils/strs"
)

var (
	// Templates holds the template information for a file that the graph is generating.
	Templates = []gen.TypeTemplate{
		{
			Name:   "view",
			Format: pkgf("%s_view.go"),
		},
	}
	// GraphTemplates holds the templates applied on the graph.
	GraphTemplates = []gen.GraphTemplate{
		{
			Name:   "base",
			Format: "entview.go",
		},
	}
	// Templates holds the Go templates for the code generation.
	templates *gen.Template
	//go:embed template/*
	templateDir embed.FS
)

func loadTemplate() *gen.Template {
	p := []string{
		"template/*.tmpl",
	}

	fn := map[string]any{
		"xtemplate":         xtemplate,
		"hasTemplate":       hasTemplate,
		"matchTemplate":     matchTemplate,
		"getInput":          getInput,
		"getOutput":         getOutput,
		"getStructTag":      getStructTag,
		"extractAnnotation": extractAnnotation,
	}

	return gen.MustParse(gen.NewTemplate("view").
		Funcs(fn).
		ParseFS(templateDir, p...))
}

func pkgf(s string) func(t *gen.Type) string {
	return func(t *gen.Type) string {
		return fmt.Sprintf(s, t.PackageDir())
	}
}

// xtemplate is a template function,
// it dynamically executes templates by their names.
func xtemplate(name string, v any) (string, error) {
	buf := bytes.NewBuffer(nil)
	if err := templates.ExecuteTemplate(buf, name, v); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// hasTemplate is a template function,
// it checks whether a template exists in the loaded templates.
func hasTemplate(name string) bool {
	for _, t := range templates.Templates() {
		if t.Name() == name {
			return true
		}
	}

	return false
}

// matchTemplate is a template function,
// it returns all template names that match the given patterns.
func matchTemplate(patterns ...string) []string {
	var (
		names   []string
		visited = sets.Set[string]{}
	)

	for _, pattern := range patterns {
		for _, t := range templates.Templates() {
			name := t.Name()
			if visited.Has(name) {
				continue
			}

			if match, _ := filepath.Match(pattern, name); match {
				names = append(names, name)
				visited.Insert(name)
			}
		}
	}

	sort.Strings(names)

	return names
}

const (
	inputTypeQuery  = "Query"
	inputTypeCreate = "Create"
	inputTypeUpdate = "Update"
	inputTypePatch  = "Patch"
)

type InputRef struct {
	PrerequisiteEdges  []*gen.Edge
	PrerequisiteFields []*gen.Field
	IndexFields        []*gen.Field
	ImmutableFields    []*gen.Field
	Fields             []*gen.Field
	AdditionalEdges    []*gen.Edge
}

// FieldsWithoutIndexing returns the Fields that not found in IndexFields.
func (r InputRef) FieldsWithoutIndexing() []*gen.Field {
	indexFieldsSet := sets.Set[*gen.Field]{}
	for i := range r.IndexFields {
		indexFieldsSet.Insert(r.IndexFields[i])
	}

	fs := make([]*gen.Field, 0, len(r.Fields))

	for i := range r.Fields {
		if indexFieldsSet.Has(r.Fields[i]) {
			continue
		}

		fs = append(fs, r.Fields[i])
	}

	return fs
}

// FieldsSkipWrite returns the Fields should not create and update from api.
func (r InputRef) FieldsSkipWrite() []*gen.Field {
	fs := make([]*gen.Field, 0, len(r.Fields))

	for i := range r.Fields {
		if r.Fields[i].StorageKey() == "" {
			continue
		}

		a, err := annotation.ExtractAnnotation(r.Fields[i].Annotations)
		if err != nil {
			continue
		}

		if a.SkipStoring {
			continue
		}

		if a.SkipInput.Create && a.SkipInput.Update {
			fs = append(fs, r.Fields[i])
		}
	}

	return fs
}

// getInput is a template function,
// it returns the Fields and Edges for generating *Input struct.
func getInput(v any, typ string) (r InputRef, err error) {
	var n *gen.Type

	switch t := v.(type) {
	default:
		return r, nil
	case *gen.Type:
		n = t
	case gen.Type:
		n = &t
	}

	if !n.HasOneFieldID() {
		return r, errors.New("composited id is not supported, replace it with indexer")
	}

	// Append fields.
	for _, f := range n.Fields {
		fa := annotation.MustExtractAnnotation(f.Annotations)

		switch {
		case typ == inputTypeCreate && fa.SkipInput.Create:
			continue
		case typ == inputTypeUpdate && (f.Immutable && !fa.Input.Update || fa.SkipInput.Update):
			continue
		default:
		}

		if f.IsEdgeField() {
			e, _ := f.Edge()
			ea := annotation.MustExtractAnnotation(e.Annotations)

			switch {
			case typ == inputTypeCreate && ea.SkipInput.Create:
				continue
			case typ == inputTypeUpdate && (e.Immutable && !ea.Input.Update || ea.SkipInput.Update):
				continue
			case typ == inputTypePatch && (ea.SkipInput.Create || ea.SkipInput.Update):
				continue
			default:
			}

			if !e.IsInverse() || e.Type == e.Ref.Type {
				continue
			}

			switch {
			case typ == inputTypeCreate && ea.Input.Create:
				continue
			case typ == inputTypeUpdate && ea.Input.Update:
				continue
			case typ == inputTypePatch && (ea.Input.Create || ea.Input.Update):
				continue
			}

			r.PrerequisiteFields = append(r.PrerequisiteFields, f)

			continue
		}

		if typ == inputTypeQuery && !fa.Input.Query {
			continue
		}

		if f.Immutable {
			r.ImmutableFields = append(r.ImmutableFields, f)
		}

		r.Fields = append(r.Fields, f)
	}

	sort.SliceStable(r.Fields, func(i, j int) bool {
		if !r.Fields[i].Optional {
			switch {
			case typ == inputTypeCreate && !r.Fields[i].Default:
				return true
			case typ == inputTypeUpdate && !r.Fields[i].UpdateDefault:
			}
		}

		return false
	})

	// Append fields of the longest unique index.
	var indexColumns []string

	for _, i := range n.Indexes {
		if !i.Unique ||
			annotation.MustExtractAnnotation(i.Annotations).SkipInput.Query {
			continue
		}

		if len(i.Columns) > len(indexColumns) {
			indexColumns = i.Columns
		}
	}

	if len(indexColumns) != 0 {
		fm := map[string]*gen.Field{}

		for i := range n.Fields {
			fm[n.Fields[i].Name] = n.Fields[i]
		}

		cfs := make([]*gen.Field, 0, len(indexColumns))

		for _, col := range indexColumns {
			f, exist := fm[col]

			if !exist || f.IsEdgeField() {
				continue
			}

			if !f.Immutable ||
				(!f.IsInt() && !f.IsInt64() && !f.IsBool() && !f.IsString()) {
				cfs = nil
				break
			}

			cfs = append(cfs, f)
		}

		r.IndexFields = append(r.IndexFields, cfs...)
	}

	// Append edges.
	for _, e := range n.Edges {
		ea := annotation.MustExtractAnnotation(e.Annotations)

		if !e.IsInverse() || e.O2M() {
			switch {
			case typ == inputTypeCreate && ea.SkipInput.Create:
				continue
			case typ == inputTypeUpdate && (e.Immutable && !ea.Input.Update || ea.SkipInput.Update):
				continue
			case typ == inputTypePatch && (ea.SkipInput.Create || ea.SkipInput.Update):
				continue
			case e.Type.EdgeSchema.To != nil && e.Type.EdgeSchema.To.Through != nil:
				continue
			}

			r.AdditionalEdges = append(r.AdditionalEdges, e)

			continue
		}

		if e.Type != e.Ref.Type {
			switch {
			case !ea.SkipInput.Query:
				r.PrerequisiteEdges = append(r.PrerequisiteEdges, e)
			case typ == inputTypeCreate && ea.Input.Create:
				r.AdditionalEdges = append(r.AdditionalEdges, e)
			case typ == inputTypeUpdate && ea.Input.Update:
				r.AdditionalEdges = append(r.AdditionalEdges, e)
			case typ == inputTypePatch && (ea.Input.Create || ea.Input.Update):
				r.AdditionalEdges = append(r.AdditionalEdges, e)
			}
		}
	}

	return r, nil
}

type OutputRef struct {
	Fields []*gen.Field
	Edges  []*gen.Edge
}

// getOutput is a template function,
// it returns the Fields and Edges for generating *Output struct.
func getOutput(v any) (r OutputRef, err error) {
	var n *gen.Type

	switch t := v.(type) {
	default:
		return r, nil
	case *gen.Type:
		n = t
	case gen.Type:
		n = &t
	}

	if !n.HasOneFieldID() {
		return r, errors.New("composited id is not supported, replace it with indexer")
	}

	// Append fields.
	for _, f := range n.Fields {
		if f.Sensitive() ||
			annotation.MustExtractAnnotation(f.Annotations).SkipOutput {
			continue
		}

		if f.IsEdgeField() {
			continue
		}

		r.Fields = append(r.Fields, f)
	}

	// Append edges, excluding O2M reversible edges.
	for _, e := range n.Edges {
		if (e.O2M() && e.IsInverse()) ||
			annotation.MustExtractAnnotation(e.Annotations).SkipOutput {
			continue
		}

		r.Edges = append(r.Edges, e)
	}

	return r, nil
}

// getStructTag is a template function,
// it returns the struct tag with the given tag name.
func getStructTag(v any, tag, typ string) string {
	var (
		st        string
		omitempty bool
	)

	switch t := v.(type) {
	default:
		return ""
	case *gen.Field:
		st = t.StructTag
		omitempty = (typ != inputTypeCreate && typ != inputTypeUpdate) ||
			(typ == inputTypeCreate && (t.Default || t.Optional)) ||
			(typ == inputTypeUpdate && (t.UpdateDefault || t.Optional))
	case gen.Field:
		st = t.StructTag
		omitempty = (typ != inputTypeCreate && typ != inputTypeUpdate) ||
			(typ == inputTypeCreate && (t.Default || t.Optional)) ||
			(typ == inputTypeUpdate && (t.UpdateDefault || t.Optional))
	case *gen.Edge:
		st = t.StructTag
		omitempty = (typ != inputTypeCreate && typ != inputTypeUpdate) ||
			(typ == inputTypeCreate && t.Optional) ||
			(typ == inputTypeUpdate && t.Optional)
	case gen.Edge:
		st = t.StructTag
		omitempty = (typ != inputTypeCreate && typ != inputTypeUpdate) ||
			(typ == inputTypeCreate && t.Optional) ||
			(typ == inputTypeUpdate && t.Optional)
	}

	attrs, ok := reflect.StructTag(st).Lookup(tag)
	if !ok || attrs == "" || attrs == "-" {
		return st
	}

	const (
		sep           = ","
		omitemptyAttr = "omitempty"
	)

	name, attrs, ok := strings.Cut(attrs, sep)
	if !ok || name == "" {
		return st
	}
	name = strs.CamelizeDownFirst(name)

	as := strings.Split(attrs, sep)
	for i := range as {
		if as[i] != omitemptyAttr {
			continue
		}

		switch i {
		case 0:
			as = as[1:]
		case len(as) - 1:
			as = as[:i]
		default:
			as = append(as[:i], as[i+1:]...)
		}

		break
	}

	if omitempty {
		as = append(as, omitemptyAttr)
	}

	attrs = strings.Join(append([]string{name}, as...), sep)

	return fmt.Sprintf(`%s:"%s"`, tag, attrs)
}

// extractAnnotation is a template function,
// it returns the annotation.Annotation of the given value.
func extractAnnotation(v any) (annotation.Annotation, error) {
	switch t := v.(type) {
	default:
		return annotation.Annotation{}, nil
	case *gen.Type:
		return annotation.ExtractAnnotation(t.Annotations)
	case gen.Type:
		return annotation.ExtractAnnotation(t.Annotations)
	case *gen.Field:
		return annotation.ExtractAnnotation(t.Annotations)
	case gen.Field:
		return annotation.ExtractAnnotation(t.Annotations)
	case *gen.Edge:
		return annotation.ExtractAnnotation(t.Annotations)
	case gen.Edge:
		return annotation.ExtractAnnotation(t.Annotations)
	}
}
