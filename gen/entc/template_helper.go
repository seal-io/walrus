package main

import (
	"fmt"

	"entgo.io/ent/entc/gen"
	"k8s.io/apimachinery/pkg/util/sets"
)

const actionUpdate = "update"

func getInputFields(n *gen.Type, a string) []*gen.Field {
	//nolint:prealloc
	var fs []*gen.Field

	// Append for query action.
	if a == "query" {
		if n.HasOneFieldID() {
			n.ID.StructTag = `uri:"id,omitempty" json:"id,omitempty"`
			fs = append(fs, n.ID)
		} else {
			for _, f := range n.EdgeSchema.ID {
				f.StructTag = getStructTag(f, false)
				fs = append(fs, f)
			}
		}

		return fs
	}

	ignoreSet := sets.New[string]("createTime", "updateTime")

	for _, fk := range n.ForeignKeys {
		if fk == nil || !fk.UserDefined {
			continue
		}
		// Ignore defined foreign key.
		ignoreSet.Insert(fk.Field.Name)
	}

	// Append for update action.
	if a == actionUpdate {
		if n.HasOneFieldID() {
			n.ID.StructTag = `uri:"id" json:"-"`
			fs = append(fs, n.ID)
		}

		for _, f := range n.EdgeSchema.ID {
			if f == nil || ignoreSet.Has(f.Name) {
				continue
			}
			f.StructTag = getStructTag(f, false)
			fs = append(fs, f)
		}
	}

	// Append continually.
	for _, f := range n.Fields {
		if f == nil || ignoreSet.Has(f.Name) {
			continue
		}

		switch a {
		default:
			continue
		case "create":
			f.StructTag = getStructTag(f, false)
		case actionUpdate:
			if f.Immutable {
				continue
			}
			f.StructTag = getStructTag(f, true)
		}

		fs = append(fs, f)
	}

	// Distinct.
	nfs := make([]*gen.Field, 0, len(fs))
	fdSet := sets.New[string]()

	for i := range fs {
		if fdSet.Has(fs[i].Name) {
			continue
		}

		fdSet.Insert(fs[i].Name)
		nfs = append(nfs, fs[i])
	}

	return nfs
}

func getInputEdges(n *gen.Type, a string) []*gen.Edge {
	ignoreSet := sets.New[string]()

	for _, fk := range n.ForeignKeys {
		if fk == nil || fk.UserDefined {
			continue
		}
		// Ignore undefined foreign key.
		ignoreSet.Insert(fk.Edge.Name)
		ignoreSet.Insert(fk.Edge.Ref.Name)
	}

	es := make([]*gen.Edge, 0, len(n.Edges))

	// Append.
	for _, e := range n.Edges {
		if e == nil || ignoreSet.Has(e.Name) {
			// NB(thxCode): cannot process edges that defining without `.Field()`.
			continue
		}

		switch {
		default:
			continue
		case e.O2O() && e.IsInverse():
			// E.g.       From 1-1 to
			//      [entity a] 1-1 [entity b],
			// generate [entity a] into [entity b].
		case e.M2M() && e.IsInverse():
			// E.g.       From *-* to
			//      [entity a] *-* [entity b],
			// generate [entity a] into [entity b].
		case e.M2O(): // Inversion.
			// E.g.      From  *-* to          to 1-* from
			//      [entity a] *-1 [relationship] 1-* [entity b],
			// generate [entity a], [entity b] into [relationship]
			// e.g.      from  1-* to
			//      [entity a] 1-* [entity b]
			// generate [entity a] into [entity b].
		case e.O2M() && (e.Through != nil && e == e.Through.EdgeSchema.To):
			// E.g.                          From 1-* to
			//      [entity a] *-1 [relationship] 1-* [entity b],
			// generate [relationship] into [entity b].
		}

		switch a {
		default:
			continue
		case "create":
			e.StructTag = getStructTag(e, false)
		case actionUpdate:
			if !n.IsEdgeSchema() && e.Immutable {
				continue
			}
			e.StructTag = getStructTag(e, true)
		}

		es = append(es, e)
	}

	return es
}

func getOutputFields(n *gen.Type) []*gen.Field {
	ignoreSet := sets.New[string]()

	for _, fk := range n.ForeignKeys {
		if fk == nil || !fk.UserDefined {
			continue
		}
		// Ignore defined foreign key.
		ignoreSet.Insert(fk.Field.Name)
	}

	for _, f := range n.Fields {
		if f == nil || !f.Sensitive() {
			continue
		}
		// Ignore sensitive field.
		ignoreSet.Insert(f.Name)
	}

	// Append.
	//nolint:prealloc
	var fs []*gen.Field

	if n.HasOneFieldID() {
		n.ID.StructTag = `json:"id,omitempty"`
		fs = append(fs, n.ID)
	} else {
		for _, f := range n.EdgeSchema.ID {
			if f == nil || ignoreSet.Has(f.Name) {
				continue
			}
			f.StructTag = getStructTag(f, true)
			fs = append(fs, f)
		}
	}

	for _, f := range n.Fields {
		if f == nil || ignoreSet.Has(f.Name) {
			continue
		}
		f.StructTag = getStructTag(f, true)
		fs = append(fs, f)
	}

	// Distinct.
	nfs := make([]*gen.Field, 0, len(fs))
	fdSet := sets.New[string]()

	for i := range fs {
		if fdSet.Has(fs[i].Name) {
			continue
		}

		fdSet.Insert(fs[i].Name)
		nfs = append(nfs, fs[i])
	}

	return nfs
}

func getOutputEdges(n *gen.Type) []*gen.Edge {
	ignoreSet := sets.New[string]()

	for _, fk := range n.ForeignKeys {
		if fk == nil || fk.UserDefined {
			continue
		}
		// Ignore undefined foreign key.
		ignoreSet.Insert(fk.Edge.Name)
		ignoreSet.Insert(fk.Edge.Ref.Name)
	}

	// Append.
	es := make([]*gen.Edge, 0, len(n.Edges))

	for _, e := range n.Edges {
		if e == nil || ignoreSet.Has(e.Name) {
			continue
		}

		e.StructTag = getStructTag(e, true)
		es = append(es, e)
	}

	return es
}

func getStructTag(v any, mustOmit bool) string {
	camel := gen.Funcs["camel"].(func(string) string)

	switch f := v.(type) {
	case *gen.Field:
		if mustOmit || f.Nillable || f.Optional || f.Default || f.UpdateDefault || f.Validators == 0 || f.Sensitive() {
			return fmt.Sprintf(`json:"%s,omitempty"`, camel(f.Name))
		}

		return fmt.Sprintf(`json:"%s"`, camel(f.Name))
	case *gen.Edge:
		if mustOmit || f.Optional || !f.Unique {
			return fmt.Sprintf(`json:"%s,omitempty"`, camel(f.Name))
		}

		return fmt.Sprintf(`json:"%s"`, camel(f.Name))
	}

	return `json:"-"`
}
