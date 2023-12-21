package types

import (
	"context"

	"github.com/seal-io/walrus/pkg/dao/model"
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
	Apply(context.Context, *model.Resource, ApplyOptions) error

	// Destroy cleans all resources of the given resource.
	Destroy(context.Context, *model.Resource, DestroyOptions) error

	// Sync syncs the state of the given resource.
	Sync(context.Context, *model.Resource, SyncOptions) error

	// Detect detect drift of the given resource.
	Detect(context.Context, *model.Resource, DetectOptions) error
}

// ApplyOptions holds the options of Deployer's Apply action.
type ApplyOptions struct{}

// DestroyOptions holds the options of Deployer's Destroy action.
type DestroyOptions struct{}

// SyncOptions holds the options of Deployer's Sync action.
type SyncOptions struct{}

// DeployerFactory creates a Deployer.
type DetectOptions struct{}
