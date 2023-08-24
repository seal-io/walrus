package serviceresources

import (
	"context"

	"go.uber.org/multierr"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	optypes "github.com/seal-io/walrus/pkg/operator/types"
)

// Label labels the given model.ServiceResource list with the given operator.Operator.
//
// The given model.ServiceResource item must be instance shape and managed mode.
//
// The given model.ServiceResource item must specify the following fields:
// Shape, Mode, ID, DeployerType, Type and Name,
// and the following edges:
// Project, Environment and Service.
func Label(
	ctx context.Context,
	op optypes.Operator,
	candidates []*model.ServiceResource,
) error {
	if op == nil {
		return nil
	}

	// Merge the errors to return them all at once,
	// instead of returning the first error.
	var berr error

	for i := range candidates {
		// Give up the loop if the context is canceled.
		if multierr.AppendInto(&berr, ctx.Err()) {
			break
		}

		// Skip the resource if it is not instance shape or not managed mode.
		if candidates[i].Shape != types.ServiceResourceShapeInstance ||
			candidates[i].Mode != types.ServiceResourceModeManaged {
			continue
		}

		edges := candidates[i].Edges
		if edges.Project == nil || edges.Project.Name == "" ||
			edges.Environment == nil || edges.Environment.Name == "" ||
			edges.Service == nil || edges.Service.Name == "" {
			continue
		}

		ls := map[string]string{
			types.LabelWalrusProjectName:     edges.Project.Name,
			types.LabelWalrusEnvironmentName: edges.Environment.Name,
			types.LabelWalrusServiceName:     edges.Service.Name,
		}

		err := op.Label(ctx, candidates[i], ls)
		if multierr.AppendInto(&berr, err) {
			continue
		}
	}

	return berr
}
