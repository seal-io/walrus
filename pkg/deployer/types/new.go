package types

import (
	"context"

	"k8s.io/client-go/rest"
)

// CreateOptions holds the options for creating Deployer.
type CreateOptions struct {
	// Type indicates the type for creating.
	Type       Type
	KubeConfig *rest.Config
}

// Creator is a factory func to create Deployer.
type Creator func(context.Context, CreateOptions) (Deployer, error)
