package dao

import (
	"context"
	stdsql "database/sql"
	"errors"

	"entgo.io/ent/dialect/sql"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/environmentconnectorrelationship"
	"github.com/seal-io/seal/pkg/dao/model/project"
	"github.com/seal-io/seal/pkg/dao/types/object"
	"github.com/seal-io/seal/utils/strs"
)

// EnvironmentConnectorsEdgeSave saves the edge connectors of model.Environment entity.
func EnvironmentConnectorsEdgeSave(ctx context.Context, mc model.ClientSet, entity *model.Environment) error {
	if entity.Edges.Connectors == nil {
		return nil
	}

	// Default new items and create key set for items.
	var (
		newItems       = entity.Edges.Connectors
		newItemsKeySet = sets.New[string]()
	)

	for i := range newItems {
		if newItems[i] == nil {
			return errors.New("invalid input: nil relationship")
		}
		newItems[i].EnvironmentID = entity.ID

		newItemsKeySet.Insert(strs.Join("/", newItems[i].EnvironmentID, newItems[i].ConnectorID))
	}

	// Add/Update new items.
	if len(newItems) != 0 {
		err := mc.EnvironmentConnectorRelationships().CreateBulk().
			Set(newItems...).
			OnConflict(
				sql.ConflictColumns(
					environmentconnectorrelationship.FieldEnvironmentID,
					environmentconnectorrelationship.FieldConnectorID,
				)).
			DoNothing().
			Exec(ctx)
		if err != nil && !errors.Is(err, stdsql.ErrNoRows) {
			return err
		}

		entity.Edges.Connectors = newItems // Feedback.
	}

	// Delete stale items.
	oldItems, err := mc.EnvironmentConnectorRelationships().Query().
		Where(environmentconnectorrelationship.EnvironmentID(entity.ID)).
		All(ctx)
	if err != nil {
		return err
	}

	deletedIDs := make([]object.ID, 0, len(oldItems))

	for i := range oldItems {
		if newItemsKeySet.Has(strs.Join("/", oldItems[i].EnvironmentID, oldItems[i].ConnectorID)) {
			continue
		}

		deletedIDs = append(deletedIDs, oldItems[i].ID)
	}

	if len(deletedIDs) != 0 {
		_, err = mc.EnvironmentConnectorRelationships().Delete().
			Where(environmentconnectorrelationship.IDIn(deletedIDs...)).
			Exec(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetEnvironmentByID gets an environment including project & connectors edges by ID.
func GetEnvironmentByID(ctx context.Context, mc model.ClientSet, id object.ID) (*model.Environment, error) {
	envs, err := GetEnvironmentsByIDs(ctx, mc, id)
	if err != nil {
		return nil, err
	}

	if len(envs) == 0 {
		return nil, errors.New("environment not found")
	}

	return envs[0], nil
}

// GetEnvironmentsByIDs gets environments including project & connectors edges by IDs.
func GetEnvironmentsByIDs(ctx context.Context, mc model.ClientSet, ids ...object.ID) ([]*model.Environment, error) {
	return mc.Environments().Query().
		Where(environment.IDIn(ids...)).
		WithProject(func(pq *model.ProjectQuery) {
			pq.Select(project.FieldName)
		}).
		WithConnectors(func(rq *model.EnvironmentConnectorRelationshipQuery) {
			rq.WithConnector()
		}).
		All(ctx)
}
