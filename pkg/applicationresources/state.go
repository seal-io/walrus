package applicationresources

import (
	"context"

	"go.uber.org/multierr"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/platform/operator"
)

// State gets status of the given model.ApplicationResource list with the given operator.Operator.
func State(ctx context.Context, op operator.Operator, modelClient model.ClientSet, candidates []*model.ApplicationResource) (berr error) {
	if op == nil {
		return
	}

	for i := range candidates {
		// get status of the application resource.
		st, err := op.GetStatus(ctx, candidates[i])
		if err != nil {
			berr = multierr.Append(berr, err)
		}
		// get endpoints of the application resource.
		eps, err := op.GetEndpoints(ctx, candidates[i])
		if err != nil {
			berr = multierr.Append(berr, err)
		}

		// new application resource status.
		newStatus := types.ApplicationResourceStatus{
			Status:            *st,
			ResourceEndpoints: eps,
		}
		if candidates[i].Status.Equal(newStatus) {
			// do not update if the status is same as previous.
			continue
		}

		err = modelClient.ApplicationResources().UpdateOne(candidates[i]).
			SetStatus(newStatus).
			Exec(ctx)
		if err != nil {
			if model.IsNotFound(err) {
				// application resource has been deleted by other thread processing.
				continue
			}
			berr = multierr.Append(berr, err)
		}
	}
	return
}
