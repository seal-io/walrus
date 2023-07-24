package dao

import (
	"context"
	stdsql "database/sql"
	"errors"

	"entgo.io/ent/dialect/sql"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/subjectrolerelationship"
	"github.com/seal-io/seal/pkg/dao/types/object"
	"github.com/seal-io/seal/utils/strs"
)

// SubjectRolesEdgeSave saves the edge roles of model.Subject entity.
func SubjectRolesEdgeSave(ctx context.Context, mc model.ClientSet, entity *model.Subject) error {
	if entity.Edges.Roles == nil {
		return nil
	}

	// Default new items and create key set for items.
	var (
		newItems       = entity.Edges.Roles
		newItemsKeySet = sets.New[string]()
	)

	for i := range newItems {
		if newItems[i] == nil {
			return errors.New("invalid input: nil relationship")
		}
		newItems[i].SubjectID = entity.ID

		newItemsKeySet.Insert(strs.Join("/", string(newItems[i].SubjectID), newItems[i].RoleID))
	}

	// Add/Update new items.
	if len(newItems) != 0 {
		err := mc.SubjectRoleRelationships().CreateBulk().
			Set(newItems...).
			OnConflict(
				sql.ConflictWhere(sql.P().
					IsNull(subjectrolerelationship.FieldProjectID)),
				sql.ConflictColumns(
					subjectrolerelationship.FieldSubjectID,
					subjectrolerelationship.FieldRoleID,
				)).
			DoNothing().
			Exec(ctx)
		if err != nil && !errors.Is(err, stdsql.ErrNoRows) {
			return err
		}

		entity.Edges.Roles = newItems // Feedback.
	}

	// Delete stale items.
	oldItems, err := mc.SubjectRoleRelationships().Query().
		Where(subjectrolerelationship.SubjectID(entity.ID)).
		All(ctx)
	if err != nil {
		return err
	}

	deletedIDs := make([]object.ID, 0, len(oldItems))

	for i := range oldItems {
		if newItemsKeySet.Has(strs.Join("/", string(oldItems[i].SubjectID), oldItems[i].RoleID)) {
			continue
		}

		deletedIDs = append(deletedIDs, oldItems[i].ID)
	}

	if len(deletedIDs) != 0 {
		_, err = mc.SubjectRoleRelationships().Delete().
			Where(subjectrolerelationship.IDIn(deletedIDs...)).
			Exec(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
