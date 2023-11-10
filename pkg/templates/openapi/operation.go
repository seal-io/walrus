package openapi

import (
	jsonpatch "github.com/evanphx/json-patch"
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

// UnionSchema generates a openapi3.Schema union of s1 and s2, for the same key s2 will overwrite s1's.
func UnionSchema(s1, s2 *openapi3.Schema) (*openapi3.Schema, error) {
	if s1 == nil {
		return s2, nil
	}

	if s2 == nil {
		return s1, nil
	}

	s1Byte, err := s1.MarshalJSON()
	if err != nil {
		return nil, err
	}

	s2Byte, err := s2.MarshalJSON()
	if err != nil {
		return nil, err
	}

	// Set sequence extension while existed.
	sequence := GetOriginalVariablesSequence(s2.Extensions)
	s1Sequence := GetOriginalVariablesSequence(s1.Extensions)
	diff := sets.NewString(s1Sequence...).Difference(sets.NewString(sequence...))

	for _, v := range s1Sequence {
		if diff.Has(v) {
			sequence = append(sequence, v)
		}
	}

	// Merge patch.
	combined, err := jsonpatch.MergeMergePatches(s1Byte, s2Byte)
	if err != nil {
		return nil, err
	}

	merged := openapi3.Schema{}

	err = merged.UnmarshalJSON(combined)
	if err != nil {
		return nil, err
	}
	merged.Extensions = NewExt(merged.Extensions).SetOriginalVariablesSequence(sequence).Export()

	return &merged, nil
}
