package entx

import (
	"bytes"
	"sort"
	"strings"
	"text/template"

	"entgo.io/ent/entc/gen"
	"golang.org/x/exp/slices"

	"github.com/seal-io/seal/pkg/dao/entx/annotation"
	"github.com/seal-io/seal/pkg/dao/entx/extension/view"
)

func loadHooks() []gen.Hook {
	return []gen.Hook{
		preHook(generateExtensionView),
		preHook(generateBuilderFields),
		preHook(generateNotStoredFields),
		preHook(clipM2MEdgeIndexes),
		preHook(flattenM2MEdges),
	}
}

func preHook(fn gen.GenerateFunc) gen.Hook {
	return func(n gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(g *gen.Graph) error {
			if err := fn(g); err != nil {
				return err
			}

			return n.Generate(g)
		})
	}
}

// generateExtensionView handles the generated *Input or *Output structs.
func generateExtensionView(g *gen.Graph) error {
	return view.Generate(g)
}

// generateBuilderFields handles the fields that generate into the builder.
func generateBuilderFields(g *gen.Graph) error {
	/* For example:
	{{ define "dialect/sql/create/fields/additional/x_pet" }}
	    {{- if eq $.Name "Pet" }}
	    object *Pet
	    fromUpsert bool
	    {{- end }}
	{{- end }}

	{{ define "dialect/sql/create_bulk/fields/additional/x_pet" }}
	    {{- if eq $.Name "Pet" }}
	    objects []*Pet
	    fromUpsert bool
	    {{- end }}
	{{- end }}

	{{ define "dialect/sql/update/fields/additional/x_pet" }}
	    {{- if eq $.Name "Pet" }}
	    object *Pet
	    {{- end }}
	{{- end }}
	*/
	const genericTemplate = `
<<$title := lower $.Name>>

{{ define "dialect/sql/create/fields/additional/x_<< $title >>" }}
    {{- if eq $.Name "<< $.Name >>" }}
    object *<< $.Name >>
    fromUpsert bool
    {{- end }}
{{- end }}

{{ define "dialect/sql/create_bulk/fields/additional/x_<< $title >>" }}
    {{- if eq $.Name "<< $.Name >>" }}
    objects []*<< $.Name >>
    fromUpsert bool
    {{- end }}
{{- end }}

{{ define "dialect/sql/update/fields/additional/x_<< $title >>" }}
    {{- if eq $.Name "<< $.Name >>" }}
    object *<< $.Name >>
    {{- end }}
{{- end }}

`

	generic := template.Must(template.New("generic").
		Delims("<<", ">>").
		Funcs(map[string]any{
			"lower": strings.ToLower,
		}).
		Parse(genericTemplate))

	for i := range g.Nodes {
		var b bytes.Buffer

		err := generic.Execute(&b, g.Nodes[i])
		if err != nil {
			return err
		}

		t, err := gen.NewTemplate("external").
			Parse(b.String())
		if err != nil {
			return err
		}

		g.Templates = append(g.Templates, t)
	}

	return nil
}

// generateNotStoredFields handles the fields that are not stored in the database.
func generateNotStoredFields(g *gen.Graph) error {
	nodeNotStoredFields := make(map[*gen.Type][]*gen.Field, len(g.Nodes))

	for i := range g.Nodes {
		storedFs := make([]*gen.Field, 0, len(g.Nodes[i].Fields))

		for j := range g.Nodes[i].Fields {
			if !annotation.MustExtractAnnotation(g.Nodes[i].Fields[j].Annotations).SkipStoring {
				storedFs = append(storedFs, g.Nodes[i].Fields[j])
				continue
			}

			nodeNotStoredFields[g.Nodes[i]] = append(nodeNotStoredFields[g.Nodes[i]],
				g.Nodes[i].Fields[j])
		}
		g.Nodes[i].Fields = storedFs
	}

	if len(nodeNotStoredFields) == 0 {
		return nil
	}

	/* For example:
	{{ define "model/fields/additional" }}
	    {{- if eq $.Name "Pet" }}
	    // Test holds the value of the "test" field.
	    Test string `json:"test,omitempty"`
	    {{- end }}
	{{- end }}
	*/

	const genericTemplate = `
{{ define "model/fields/additional" }}
	<< range $node, $fields := . >>
    {{- if eq $.Name "<< $node.Name >>" }}
        << range $field := $fields >>
			<<- $commend := $field.Comment >>
			<<- $name := $field.StructField >>
			// << if $commend >><< $commend >><< else >><< $name >> holds the value of the "<< $name >>" field.<< end >>
			// << $name >> does not store in the database.
			<< $name >> << if $field.NillableValue >>*<< end >><< $field.Type >> ` + "`<< $field.StructTag >>`" + `
		<< end >>
    {{- end }}
	<< end >>
{{- end }}

`

	generic := template.Must(template.New("generic").
		Delims("<<", ">>").
		Parse(genericTemplate))

	var b bytes.Buffer

	err := generic.Execute(&b, nodeNotStoredFields)
	if err != nil {
		return err
	}

	t, err := gen.NewTemplate("external").
		Parse(b.String())
	if err != nil {
		return err
	}

	g.Templates = append(g.Templates, t)

	return nil
}

func clipM2MEdgeIndexes(g *gen.Graph) error {
	for _, n := range g.Nodes {
		for _, e := range n.Edges {
			if !e.M2M() {
				continue
			}

			if e.Through == nil {
				continue
			}

			t := e.Through

			// Get the M2M edge unique indexes to discard.
			var discardIndexes []int

			for i, idx := range t.Indexes {
				if !idx.Unique || len(idx.Columns) < 2 {
					continue
				}

				if !slices.Contains(idx.Columns, e.Rel.Columns[0]) ||
					!slices.Contains(idx.Columns, e.Rel.Columns[1]) ||
					idx.Annotations != nil {
					continue
				}

				discardIndexes = append(discardIndexes, i)
			}

			if len(discardIndexes) <= 1 {
				continue
			}

			// Discard the M2M edge unique indexes.
			sort.Slice(discardIndexes, func(i, j int) bool {
				return len(t.Indexes[i].Columns) > len(t.Indexes[j].Columns)
			})

			switch x := discardIndexes[0]; x {
			case 0:
				t.Indexes = t.Indexes[1:]
			case len(t.Indexes) - 1:
				t.Indexes = t.Indexes[:x]
			default:
				t.Indexes = append(t.Indexes[:x], t.Indexes[x+1:]...)
			}
		}
	}

	return nil
}

// flattenM2MEdges handles the M2M edges that are created by M2M through edge.
func flattenM2MEdges(g *gen.Graph) error {
	for _, n := range g.Nodes {
		// Index O2M edges.
		o2mEdgesIndex := make(map[string]*gen.Edge, 0)

		for i := range n.Edges {
			if !n.Edges[i].O2M() {
				continue
			}
			o2mEdgesIndex[n.Edges[i].Name] = n.Edges[i]
		}

		// Index O2M edges created by M2M through edge.
		o2mThroughEdgesIndex := make(map[string]*gen.Edge, 0)

		for i := range n.Edges {
			if !n.Edges[i].M2M() || n.Edges[i].Through == nil {
				continue
			}

			for o2mEdgeName, o2mEdge := range o2mEdgesIndex {
				if o2mEdge.Rel.Table != n.Edges[i].Rel.Table {
					continue
				}
				o2mThroughEdgesIndex[o2mEdgeName] = n.Edges[i]
			}
		}

		if len(o2mThroughEdgesIndex) == 0 {
			continue
		}

		// Clear O2M edges created by M2M through edge.
		edges := make([]*gen.Edge, 0, len(n.Edges))

		for i := range n.Edges {
			if _, exist := o2mThroughEdgesIndex[n.Edges[i].Name]; exist {
				continue
			}

			edges = append(edges, n.Edges[i])
		}
		n.Edges = edges

		// Modify M2M through edge to point to the O2M edge.
		for o2mEdgeName, throughEdge := range o2mThroughEdgesIndex {
			o2mEdge := o2mEdgesIndex[o2mEdgeName]
			o2mEdge.Name = throughEdge.Name
			o2mEdge.Owner = throughEdge.Owner
			o2mEdge.Optional = throughEdge.Optional
			o2mEdge.StructTag = throughEdge.StructTag
			o2mEdge.Annotations = throughEdge.Annotations
			*throughEdge = *o2mEdge
		}
	}

	return nil
}
