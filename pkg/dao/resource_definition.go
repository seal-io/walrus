package dao

import (
	"context"
	stdsql "database/sql"
	"errors"

	"entgo.io/ent/dialect/sql"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/resourcedefinitionmatchingrule"
	"github.com/seal-io/walrus/pkg/dao/types/object"
)

// ResourceDefinitionMatchingRulesEdgeSave saves the edge matching rules of model.ResourceDefinition entity.
func ResourceDefinitionMatchingRulesEdgeSave(
	ctx context.Context,
	mc model.ClientSet,
	entity *model.ResourceDefinition,
) error {
	if entity.Edges.MatchingRules == nil {
		return nil
	}

	// Default new items and create key set for items.
	var (
		newItems       = entity.Edges.MatchingRules
		newItemsKeySet = sets.New[string]()
	)

	for i := range newItems {
		if newItems[i] == nil {
			return errors.New("invalid input: nil relationship")
		}
		newItems[i].ResourceDefinitionID = entity.ID
		newItems[i].Order = i

		newItemsKeySet.Insert(newItems[i].Name)
	}

	// Add/Update new items.
	if len(newItems) != 0 {
		err := mc.ResourceDefinitionMatchingRules().CreateBulk().
			Set(newItems...).
			OnConflict(
				sql.ConflictColumns(
					resourcedefinitionmatchingrule.FieldResourceDefinitionID,
					resourcedefinitionmatchingrule.FieldTemplateID,
					resourcedefinitionmatchingrule.FieldName,
				),
			).
			UpdateNewValues().
			Exec(ctx)
		if err != nil && !errors.Is(err, stdsql.ErrNoRows) {
			return err
		}
	}

	// Delete stale items.
	oldItems, err := mc.ResourceDefinitionMatchingRules().Query().
		Where(resourcedefinitionmatchingrule.ResourceDefinitionID(entity.ID)).
		All(ctx)
	if err != nil {
		return err
	}

	deletedIDs := make([]object.ID, 0, len(oldItems))

	for i := range oldItems {
		if newItemsKeySet.Has(oldItems[i].Name) {
			continue
		}

		deletedIDs = append(deletedIDs, oldItems[i].ID)
	}

	if len(deletedIDs) != 0 {
		_, err = mc.ResourceDefinitionMatchingRules().Delete().
			Where(resourcedefinitionmatchingrule.IDIn(deletedIDs...)).
			Exec(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
