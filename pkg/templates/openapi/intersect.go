package openapi

import (
	"github.com/getkin/kin-openapi/openapi3"
	"k8s.io/apimachinery/pkg/util/sets"
)

// IntersectSchema generates a openapi3.Schema with properties intersection of s1 and s2.
func IntersectSchema(s1, s2 *openapi3.Schema) *openapi3.Schema {
	if s1 == nil {
		return s2
	}

	if s2 == nil {
		return s1
	}

	r1 := sets.New[string](s1.Required...)
	p1 := sets.KeySet[string](s1.Properties)

	r2 := sets.New[string](s2.Required...)
	p2 := sets.KeySet[string](s2.Properties)

	required := r1.Intersection(r2)
	propertyKeys := p1.Intersection(p2)

	s := &openapi3.Schema{}
	s.Required = required.UnsortedList()
	s.Properties = make(openapi3.Schemas)

	for key := range propertyKeys {
		s.Properties[key] = s1.Properties[key]
	}

	return s
}
