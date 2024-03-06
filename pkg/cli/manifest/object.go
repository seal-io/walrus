package manifest

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/utils/json"
)

const (
	GroupResources    = "resources"
	GroupResourceRuns = "resourceruns"
)

// GroupSequence defines the sequence of group to create.
var GroupSequence = []string{
	GroupResources,
}

// IDName represents an object with an ID and a name.
type IDName struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// Object represents an object with a name, group, and value.
type Object struct {
	ObjectScope
	IDName

	Group  string
	Value  map[string]any
	Status ObjectStatus
}

// Key returns the key of Object.
func (o Object) Key() string {
	switch {
	case o.Project != "" && o.Environment != "":
		return fmt.Sprintf("%s/%s/%s", o.Project, o.Environment, o.Name)
	case o.Project != "":
		return fmt.Sprintf("%s/%s", o.Project, o.Name)
	default:
		return o.Name
	}
}

// MarshalJSON returns the JSON encoding of Object.
func (o Object) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.Value)
}

// ObjectList is the list of object.
type ObjectList []Object

func (l ObjectList) Names() sets.Set[string] {
	s := sets.Set[string]{}
	for _, v := range l {
		s.Insert(v.Name)
	}

	return s
}

// IDNameMap returns a map of ID to name.
func (l ObjectList) IDNameMap() map[string]string {
	m := make(map[string]string)
	for _, v := range l {
		m[v.ID] = v.Name
	}

	return m
}

// ObjectScope represents the scope of an object within a project and environment.
type ObjectScope struct {
	Project     string
	Environment string
}

// Map returns a map of ObjectScope.
func (o *ObjectScope) Map() map[string]any {
	return map[string]any{
		"project":     &o.Project,
		"environment": &o.Environment,
	}
}

// ScopedName generate object scoped name.
func (o *ObjectScope) ScopedName(name string) string {
	switch {
	case o.Project != "" && o.Environment != "":
		return fmt.Sprintf("%s/%s/%s", o.Project, o.Environment, name)
	case o.Project != "":
		return fmt.Sprintf("%s/%s", o.Project, name)
	default:
		return name
	}
}

// ObjectByScope represents a map of ObjectScope to a slice of Object.
type ObjectByScope map[ObjectScope]ObjectList

// All returns all objects in the ObjectByScope.
func (s *ObjectByScope) All() ObjectList {
	if s == nil || len(*s) == 0 {
		return nil
	}

	var all []Object
	for _, v := range *s {
		all = append(all, v...)
	}

	return all
}

// NewObjectSet returns a new ObjectSet.
func NewObjectSet(objs ...Object) *ObjectSet {
	set := &ObjectSet{}
	set.Add(objs...)

	return set
}

// ObjectByGK represents a map of group to a map of ObjectScope to a slice of Object.
type ObjectByGK map[string]map[ObjectScope]ObjectList

// ObjectSet represents a set of objects.
type ObjectSet struct {
	objectByGK ObjectByGK
	objectList ObjectList
}

// Add adds objects to the ObjectSet.
func (s *ObjectSet) Add(objs ...Object) *ObjectSet {
	for _, o := range objs {
		if s.objectByGK == nil {
			s.objectByGK = make(ObjectByGK)
		}

		if _, ok := s.objectByGK[o.Group]; !ok {
			s.objectByGK[o.Group] = make(map[ObjectScope]ObjectList)
		}

		s.objectByGK[o.Group][o.ObjectScope] = append(s.objectByGK[o.Group][o.ObjectScope], o)
	}

	s.objectList = append(s.objectList, objs...)

	return s
}

// Remove removes objects from the ObjectSet.
func (s *ObjectSet) Remove(objs ...Object) *ObjectSet {
	if s == nil || s.Len() == 0 {
		return s
	}

	if s.objectByGK == nil {
		return s
	}

	if len(objs) == 0 {
		return s
	}

	ks := sets.Set[string]{}

	for _, v := range objs {
		ks.Insert(v.Key())
	}

	var (
		newList ObjectList
		obg     = make(ObjectByGK)
	)

	for i, v := range s.objectList {
		if !ks.Has(s.objectList[i].Key()) {
			newList = append(newList, s.objectList[i])

			if _, ok := obg[v.Group]; !ok {
				obg[v.Group] = make(map[ObjectScope]ObjectList)
			}

			obg[v.Group][v.ObjectScope] = append(obg[v.Group][v.ObjectScope], s.objectList[i])
		}
	}

	s.objectList = newList
	s.objectByGK = obg

	return s
}

// All returns all objects in the ObjectSet.
func (s *ObjectSet) All() []Object {
	return s.objectList
}

// Len returns the number of objects in the ObjectSet.
func (s *ObjectSet) Len() int {
	return len(s.objectList)
}

// ByGroup returns a map of ObjectScope to a slice of Object for the given group.
func (s *ObjectSet) ByGroup(group string) ObjectByScope {
	return s.objectByGK[group]
}

// ObjectByGroup returns a map of ObjectScope to a slice of Object for the given group.
func (s *ObjectSet) ObjectByGroup() ObjectByGK {
	return s.objectByGK
}
