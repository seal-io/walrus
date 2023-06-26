package io

import "entgo.io/ent/schema"

// DisableInputWhenCreating doesn't generate the field or edge into *CreateInput entity.
func DisableInputWhenCreating() schema.Annotation {
	return &annotation{
		CreateInputDisabled: true,
	}
}

// DisableInputWhenUpdating doesn't generate the field or edge into *UpdateInput entity.
func DisableInputWhenUpdating() schema.Annotation {
	return &annotation{
		UpdateInputDisabled: true,
	}
}

// DisableInput doesn't generate the field or edge into *Input entity.
func DisableInput() schema.Annotation {
	return &annotation{
		CreateInputDisabled: true,
		UpdateInputDisabled: true,
	}
}

// DisableOutput doesn't generate the field or edge into *Output entity.
func DisableOutput() schema.Annotation {
	return &annotation{
		OutputDisabled: true,
	}
}

// Disable doesn't generate the field or edge into *Input/*Output entity.
func Disable() schema.Annotation {
	return &annotation{
		CreateInputDisabled: true,
		UpdateInputDisabled: true,
		OutputDisabled:      true,
	}
}

const annotationName = "EntIO"

type annotation struct {
	// CreateInputDisabled doesn't generate the field or edge into *CreateInput entity if true.
	CreateInputDisabled bool

	// UpdateInputDisabled doesn't generate the field or edge into *UpdateInput entity if true.
	UpdateInputDisabled bool

	// OutputDisabled doesn't generate the field or edge into *Output entity if true.
	OutputDisabled bool
}

func (annotation) Name() string {
	return annotationName
}
