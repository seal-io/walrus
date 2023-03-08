package deployer

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

	// Apply creates/updates the resources of the given application instance,
	// also cleans stale resources of the given application.
	Apply(context.Context, *model.ApplicationInstance, ApplyOptions) error

	// Destroy cleans all resources of the given application instance.
	Destroy(context.Context, *model.ApplicationInstance, DestroyOptions) error
}

// ApplyOptions holds the options of Deployer's Apply action.
type ApplyOptions struct {
	SkipTLSVerify bool
}

// DestroyOptions holds the options of Deployer's Destroy action.
type DestroyOptions struct {
	SkipTLSVerify bool
}
