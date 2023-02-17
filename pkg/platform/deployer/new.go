package deployer

import "context"

// CreateOptions holds the options for creating Deployer.
type CreateOptions struct {
	// Type indicates the type for creating.
	Type Type
}

// Creator is a factory func to create Deployer.
type Creator func(context.Context, CreateOptions) (Deployer, error)
