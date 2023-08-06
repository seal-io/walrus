package distributor

import (
	"context"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/object"
)

func connectorIDs(ctx context.Context, client model.ClientSet) ([]object.ID, error) {
	return client.Connectors().Query().
		Where(connector.Type(types.ConnectorTypeK8s)).
		IDs(ctx)
}

func applyItemDisplayName(
	ctx context.Context,
	client model.ClientSet,
	items []Resource,
	groupBy types.GroupByField,
) error {
	if groupBy != types.GroupByFieldConnectorID {
		return nil
	}

	// Group by connector id.
	conns, err := client.Connectors().Query().
		Where(
			connector.TypeEQ(types.ConnectorTypeK8s),
		).
		Select(
			connector.FieldID,
			connector.FieldName,
		).
		All(ctx)
	if err != nil {
		return err
	}

	for i, v := range items {
		for _, conn := range conns {
			if v.ItemName == conn.ID.String() {
				items[i].ItemName = conn.Name
				break
			}
		}
	}

	return nil
}
