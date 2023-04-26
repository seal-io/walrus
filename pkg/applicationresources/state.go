package applicationresources

import (
	"context"

	"go.uber.org/multierr"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/platform"
	"github.com/seal-io/seal/pkg/platform/operator"
)

// State gets status of the given model.ApplicationResource list according to its connector.
func State(ctx context.Context, modelClient model.ClientSet, candidates []*model.ApplicationResource) (berr error) {
	for i := range candidates {
		// get operator.
		var op, err = platform.GetOperator(ctx, operator.CreateOptions{
			Connector: *candidates[i].Edges.Connector,
		})
		if multierr.AppendInto(&berr, err) {
			continue
		}

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

// ListStateCandidatesByPage gets the candidates for State by pagination params.
func ListStateCandidatesByPage(ctx context.Context, modelClient model.ClientSet, offset, limit int) ([]*model.ApplicationResource, error) {
	return queryStateCandidates(modelClient).
		Offset(offset).
		Limit(limit).
		All(ctx)
}

func queryStateCandidates(modelClient model.ClientSet) *model.ApplicationResourceQuery {
	return modelClient.ApplicationResources().Query().
		Order(model.Desc(applicationresource.FieldCreateTime)).
		Unique(false).
		Select(
			applicationresource.FieldID,
			applicationresource.FieldStatus,
			applicationresource.FieldInstanceID,
			applicationresource.FieldConnectorID,
			applicationresource.FieldType,
			applicationresource.FieldName,
			applicationresource.FieldDeployerType).
		WithConnector(func(cq *model.ConnectorQuery) {
			cq.Select(
				connector.FieldName,
				connector.FieldType,
				connector.FieldCategory,
				connector.FieldConfigVersion,
				connector.FieldConfigData)
		})
}
