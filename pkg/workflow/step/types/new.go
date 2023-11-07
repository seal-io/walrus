package types

import "github.com/seal-io/walrus/pkg/dao/model"

// CreateOptions holds the options for creating step manager.
type CreateOptions struct {
	Type        Type
	ModelClient model.ClientSet
}

// Creator is a factory func to create step manager.
type Creator func(CreateOptions) (StepManager, error)
