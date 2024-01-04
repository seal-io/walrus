package annotation

import (
	"fmt"

	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema"

	"github.com/seal-io/walrus/utils/json"
)

var (
	_ schema.Annotation = (*Annotation)(nil)
	_ schema.Merger     = (*Annotation)(nil)
)

type (
	Annotation struct {
		// SkipInput indicates how to skip generating into *Input struct.
		SkipInput struct {
			// Query skips generating the field or edge into *QueryInput struct if true.
			Query bool `json:"Query,omitempty"`
			// Create skips generating the field or edge into *CreateInput struct if true.
			Create bool `json:"Create,omitempty"`
			// Update skips generating the field or edge into *UpdateInput struct if true.
			Update bool `json:"Update,omitempty"`
		} `json:"SkipInput,omitempty"`

		// Input indicates how to generate into *Input struct.
		Input struct {
			// Query generates the field into *QueryInput struct if true.
			Query bool `json:"Query,omitempty"`
			// Create generates the edge into *CreateInput struct if true.
			Create bool `json:"Create,omitempty"`
			// Update generates the immutable field or edge into *UpdateInput struct if true.
			Update bool `json:"Update,omitempty"`
		} `json:"Input,omitempty"`

		// ValidateContextFuncs indicates funcs to call before validating *Input struct.
		ValidateContextFuncs []string `json:"ValidateContextFuncs,omitempty"`

		// SkipValidateIfNotPresent skips validating the field if it is not present.
		SkipValidateIfNotPresent bool `json:"SkipValidateIfNotPresent,omitempty"`

		// SkipOutput skips generating the field or edge into *Output struct.
		SkipOutput bool `json:"SkipOutput,omitempty"`

		// SkipClearing skips generating clearer if the mutable field is optional at updating.
		SkipClearing bool `json:"SkipClearing,omitempty"`

		// SkipStoring treats the field as additional field,
		// which is not stored in the database.
		SkipStoring bool `json:"SkipStoring,omitempty"`
	}
)

func (Annotation) Name() string {
	return "EntX"
}

func (a *Annotation) Decode(annotation any) error {
	buf, err := json.Marshal(annotation)
	if err != nil {
		return err
	}

	return json.Unmarshal(buf, a)
}

func (a Annotation) Merge(other schema.Annotation) schema.Annotation {
	var o Annotation
	switch ot := other.(type) {
	case Annotation:
		o = ot
	case *Annotation:
		if ot != nil {
			o = *ot
		}
	default:
		return a
	}

	if o.SkipInput.Query {
		a.SkipInput.Query = true
	}

	if o.SkipInput.Create {
		a.SkipInput.Create = true
	}

	if o.SkipInput.Update {
		a.SkipInput.Update = true
	}

	if o.Input.Query {
		a.Input.Query = true
	}

	if o.Input.Create {
		a.Input.Create = true
	}

	if o.Input.Update {
		a.Input.Update = true
	}

	if o.ValidateContextFuncs != nil {
		a.ValidateContextFuncs = append(a.ValidateContextFuncs, o.ValidateContextFuncs...)
	}

	if o.SkipValidateIfNotPresent {
		a.SkipValidateIfNotPresent = true
	}

	if o.SkipOutput {
		a.SkipOutput = true
	}

	if o.SkipClearing {
		a.SkipClearing = true
	}

	if o.SkipStoring {
		a.SkipStoring = true
	}

	return a
}

// ExtractAnnotation extracts the entx.Annotation or returns its empty value.
func ExtractAnnotation(ants gen.Annotations) (ant Annotation, err error) {
	n := ant.Name()
	if ants != nil && ants[n] != nil {
		cn := "_" + n + "_"

		// Return caching result if found.
		if v, exist := ants[cn]; exist {
			ant = v.(Annotation)
			return
		}

		// Decode.
		err = ant.Decode(ants[n])
		if err == nil {
			// Cache for speeding up.
			ants[cn] = ant
		}
	}

	return
}

// MustExtractAnnotation extracts the entx.Annotation or panics.
func MustExtractAnnotation(ants gen.Annotations) Annotation {
	ant, err := ExtractAnnotation(ants)
	if err != nil {
		panic(fmt.Errorf("failed extracting %s annotation: %w", ant.Name(), err))
	}

	return ant
}
