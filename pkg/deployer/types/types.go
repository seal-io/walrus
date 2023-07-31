package types

import (
	"context"

	"github.com/seal-io/seal/pkg/dao/model"
)

// Type indicates the type of Deployer,
// e.g. Terraform, KubeVela, etc.
type Type = string

// Deployer holds the actions that a deployer must satisfy.
type Deployer interface {
	// Type returns Type.
	Type() Type

	// Apply creates/updates the resources of the given service,
	// also cleans stale resources.
	Apply(context.Context, *model.Service, ApplyOptions) error

	// Destroy cleans all resources of the given service.
	Destroy(context.Context, *model.Service, DestroyOptions) error

	// Refresh sync the resources from remote system of the given service.
	Refresh(context.Context, *model.Service, RefreshOptions) error

	// Detect will detect resource changes from remote system of given service.
	Detect(context.Context, *model.Service, DetectOptions) error
}

// ApplyOptions holds the options of Deployer's Apply action.
type ApplyOptions struct {
	SkipTLSVerify bool
	// Tags is the service revision tags.
	Tags []string
}

// DestroyOptions holds the options of Deployer's Destroy action.
type DestroyOptions struct {
	SkipTLSVerify bool
}

// RefreshOptions holds the options of Deployer's Refresh action.
type RefreshOptions struct {
	SkipTLSVerify bool
}

// DetectOptions holds the options of Deployer's detect action.
type DetectOptions struct {
	SkipTLSVerify bool
}
