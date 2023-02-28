package main

import (
	"fmt"

	"entgo.io/ent/entc/gen"
	"k8s.io/apimachinery/pkg/util/sets"
)

func getInputFields(n *gen.Type, a string) []*gen.Field {
	var fs []*gen.Field

	// append for query action.
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

	var ignoreSet = sets.New[string]("createTime", "updateTime")
	for _, fk := range n.ForeignKeys {
		if fk == nil || !fk.UserDefined {
			continue
		}
		// ignore defined foreign key.
		ignoreSet.Insert(fk.Field.Name)
	}

	// append for update action.
	if a == "update" {
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

	// append continually.
	for _, f := range n.Fields {
		if f == nil || ignoreSet.Has(f.Name) {
			continue
		}
		switch a {
		default:
			continue
		case "create":
			f.StructTag = getStructTag(f, false)
		case "update":
			if f.Immutable {
				continue
			}
			f.StructTag = getStructTag(f, true)
		}
		fs = append(fs, f)
	}

	// distinct.
	var nfs []*gen.Field
	var fdSet = sets.New[string]()
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
	var ignoreSet = sets.New[string]()
	for _, fk := range n.ForeignKeys {
		if fk == nil || fk.UserDefined {
			continue
		}
		// ignore undefined foreign key.
		ignoreSet.Insert(fk.Edge.Name)
		ignoreSet.Insert(fk.Edge.Ref.Name)
	}

	var es []*gen.Edge

	// append.
	for _, e := range n.Edges {
		if e == nil || ignoreSet.Has(e.Name) {
			// NB(thxCode): cannot process edges that defining without `.Field()`.
			continue
		}
		switch {
		default:
			continue
		case e.O2O() && e.IsInverse():
			// e.g.       from 1-1 to
			//      [entity a] 1-1 [entity b],
			// generate [entity a] into [entity b].
		case e.M2M() && e.IsInverse():
			// e.g.       from *-* to
			//      [entity a] *-* [entity b],
			// generate [entity a] into [entity b].
		case e.M2O(): // inversion.
			// e.g.      from  *-* to          to 1-* from
			//      [entity a] *-1 [relationship] 1-* [entity b],
			// generate [entity a], [entity b] into [relationship]
			// e.g.      from  1-* to
			//      [entity a] 1-* [entity b]
			// generate [entity a] into [entity b]
		case e.O2M() && (e.Through != nil && e == e.Through.EdgeSchema.To):
			// e.g.                          from 1-* to
			//      [entity a] *-1 [relationship] 1-* [entity b],
			// generate [relationship] into [entity b].
		}
		switch a {
		default:
			continue
		case "create":
			e.StructTag = getStructTag(e, false)
		case "update":
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
	var ignoreSet = sets.New[string]()
	for _, fk := range n.ForeignKeys {
		if fk == nil || !fk.UserDefined {
			continue
		}
		// ignore defined foreign key.
		ignoreSet.Insert(fk.Field.Name)
	}
	for _, f := range n.Fields {
		if f == nil || !f.Sensitive() {
			continue
		}
		// ignore sensitive field.
		ignoreSet.Insert(f.Name)
	}

	// append.
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

	// distinct.
	var nfs []*gen.Field
	var fdSet = sets.New[string]()
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
	var ignoreSet = sets.New[string]()
	for _, fk := range n.ForeignKeys {
		if fk == nil || fk.UserDefined {
			continue
		}
		// ignore undefined foreign key.
		ignoreSet.Insert(fk.Edge.Name)
		ignoreSet.Insert(fk.Edge.Ref.Name)
	}

	// append.
	var es []*gen.Edge
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
	var camel = gen.Funcs["camel"].(func(string) string)
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
