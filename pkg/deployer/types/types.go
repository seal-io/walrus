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

	// Apply creates/updates the resources of the given ResourceRun,
	// also cleans stale resources.
	Apply(context.Context, model.ClientSet, *model.ResourceRun, ApplyOptions) error

	// Destroy cleans all resources of the given ResourceRun.
	Destroy(context.Context, model.ClientSet, *model.ResourceRun, DestroyOptions) error

	// Plan plans the resources of the given ResourceRun.
	Plan(context.Context, model.ClientSet, *model.ResourceRun, PlanOptions) error
}

// ApplyOptions holds the options of Deployer's Apply action.
type ApplyOptions struct{}

// DestroyOptions holds the options of Deployer's Destroy action.
type DestroyOptions struct{}

// PlanOptions holds the options of Deployer's Plan action.
type PlanOptions struct{}
