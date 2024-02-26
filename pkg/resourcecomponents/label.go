package resourcecomponents

import (
	"context"

	"go.uber.org/multierr"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	optypes "github.com/seal-io/walrus/pkg/operator/types"
)

// Label labels the given model.ResourceComponent list with the given operator.Operator.
//
// The given model.ResourceComponent item must be instance shape and managed mode.
//
// The given model.ResourceComponent item must specify the following fields:
// Shape, Mode, ID, DeployerType, Type and Name,
// and the following edges:
// Project, Environment and Resource.
func Label(
	ctx context.Context,
	op optypes.Operator,
	candidates []*model.ResourceComponent,
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
		if candidates[i].Shape != types.ResourceComponentShapeInstance ||
			candidates[i].Mode != types.ResourceComponentModeManaged {
			continue
		}

		edges := candidates[i].Edges
		if edges.Project == nil || edges.Project.Name == "" ||
			edges.Environment == nil || edges.Environment.Name == "" ||
			edges.Resource == nil || edges.Resource.Name == "" {
			continue
		}

		ls := map[string]string{
			types.LabelWalrusProjectName:     edges.Project.Name,
			types.LabelWalrusEnvironmentName: edges.Environment.Name,
			types.LabelWalrusResourceName:    edges.Resource.Name,
		}

		err := op.Label(ctx, candidates[i], ls)
		if multierr.AppendInto(&berr, err) {
			continue
		}
	}

	return berr
}
