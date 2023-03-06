package deployer

import (
	"context"

	"k8s.io/client-go/rest"

	"github.com/seal-io/seal/pkg/dao/model"
)

// CreateOptions holds the options for creating Deployer.
type CreateOptions struct {
	// Type indicates the type for creating.
	Type        Type
	ModelClient model.ClientSet
	KubeConfig  *rest.Config
}

// Creator is a factory func to create Deployer.
type Creator func(context.Context, CreateOptions) (Deployer, error)
