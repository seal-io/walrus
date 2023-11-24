package openapi

import (
	"bytes"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/getkin/kin-openapi/openapi3"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/utils/log"
)

// IntersectSchema generates a openapi3.Schema with properties intersection of s1 and s2.
func IntersectSchema(s1, s2 *openapi3.Schema) *openapi3.Schema {
	if s1 == nil {
		return s2
	}

	if s2 == nil {
		return s1
	}

	if s1.Type != s2.Type {
		return nil
	}

	r1 := sets.New[string](s1.Required...)
	p1 := sets.KeySet[string](s1.Properties)

	r2 := sets.New[string](s2.Required...)
	p2 := sets.KeySet[string](s2.Properties)

	required := r1.Intersection(r2)
	propertyKeys := p1.Intersection(p2)

	s := &openapi3.Schema{
		Type: s1.Type,
	}
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

	var (
		s2o  = GetExtOriginal(s2.Extensions)
		s1o  = GetExtOriginal(s1.Extensions)
		diff = sets.NewString(s1o.VariablesSequence...).
			Difference(sets.NewString(s2o.VariablesSequence...))
		sequence = s2o.VariablesSequence
	)

	for _, v := range s1o.VariablesSequence {
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
	merged.Extensions = NewExtFromMap(merged.Extensions).WithOriginalVariablesSequence(sequence).Export()

	return &merged, nil
}

// SchemaEqual checks if s1 and s2 are equal.
func SchemaEqual(s1, s2 *openapi3.Schema) (bool, error) {
	if s1 == nil && s2 != nil {
		return false, nil
	}

	if s2 == nil && s1 != nil {
		return false, nil
	}

	s1Byte, err := s1.MarshalJSON()
	if err != nil {
		return false, err
	}

	s2Byte, err := s2.MarshalJSON()
	if err != nil {
		return false, err
	}

	return bytes.Equal(s1Byte, s2Byte), nil
}

// MustSchemaEqual checks if s1 and s2 are equal, just log the error.
func MustSchemaEqual(s1, s2 *openapi3.Schema) bool {
	equal, err := SchemaEqual(s1, s2)
	if err != nil {
		log.Warnf("failed to compare schema, %v", err)
	}

	return equal
}
