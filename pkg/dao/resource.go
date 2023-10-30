package dao

import (
	"context"
	stdsql "database/sql"
	"errors"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/model/resourcerelationship"
	"github.com/seal-io/walrus/pkg/dao/model/resourcerevision"
	"github.com/seal-io/walrus/pkg/dao/types/object"
)

// ResourceDependenciesEdgeSave saves the edge dependencies of model.Resource entity.
func ResourceDependenciesEdgeSave(ctx context.Context, mc model.ClientSet, entity *model.Resource) error {
	// Default new items and create key set for new items.
	newItems, err := resourceRelationshipGetDependencies(ctx, mc, entity)
	if err != nil {
		return err
	}

	newItemKeySet := sets.New[string]()

	for i := range newItems {
		if newItems[i] == nil {
			return errors.New("invalid input: nil relationship")
		}

		newItemKeySet.Insert(resourceRelationshipGetCompletePath(newItems[i]))
	}

	// Add/Update new items.
	if len(newItems) != 0 {
		err = mc.ResourceRelationships().CreateBulk().
			Set(newItems...).
			OnConflict(
				sql.ConflictColumns(
					resourcerelationship.FieldResourceID,
					resourcerelationship.FieldDependencyID,
					resourcerelationship.FieldPath,
				)).
			DoNothing().
			Exec(ctx)
		if err != nil && !errors.Is(err, stdsql.ErrNoRows) {
			return err
		}

		entity.Edges.Dependencies = newItems // Feedback.
	}

	// Delete stale items.
	oldItems, err := mc.ResourceRelationships().Query().
		Where(resourcerelationship.ResourceID(entity.ID)).
		All(ctx)
	if err != nil {
		return err
	}

	for i := range oldItems {
		if newItemKeySet.Has(resourceRelationshipGetCompletePath(oldItems[i])) {
			continue
		}

		_, err = mc.ResourceRelationships().Delete().
			Where(
				func(s *sql.Selector) {
					s.Where(sqljson.ValueContains(resourcerelationship.FieldPath, oldItems[i].Path))
				}).
			Exec(ctx)
		if err != nil {
			return err
		}
	}

	// Recreate dependants.
	return resourceRelationshipUpdateDependants(ctx, mc, entity)
}

func GetLatestRevisions(
	ctx context.Context,
	modelClient model.ClientSet,
	resourceIDs ...object.ID,
) ([]*model.ResourceRevision, error) {
	// Get the latest revisions of given resources by the following sql:
	// SELECT resource_revisions.*
	// FROM resource_revisions
	// JOIN (
	// 	 SELECT resource_id, MAX(create_time) AS create_time FROM resource_revisions GROUP BY resource_id
	// ) t
	// ON resource_revisions.resource_id=t.resource_id
	// AND resource_revisions.create_time=t.create_time
	// WHERE resource_revisions.resource_id IN (...)
	ids := make([]any, len(resourceIDs))
	for i := range resourceIDs {
		ids[i] = resourceIDs[i]
	}

	return modelClient.ResourceRevisions().Query().
		Modify(func(s *sql.Selector) {
			t := sql.Select(
				resourcerevision.FieldResourceID,
				sql.As(sql.Max(resourcerevision.FieldCreateTime), resourcerevision.FieldCreateTime),
			).
				From(sql.Table(resourcerevision.Table)).
				GroupBy(resourcerevision.FieldResourceID).
				As("t")
			s.Join(t).
				OnP(
					sql.And(
						sql.ColumnsEQ(
							s.C(resourcerevision.FieldResourceID),
							t.C(resourcerevision.FieldResourceID),
						),
						sql.ColumnsEQ(
							s.C(resourcerevision.FieldCreateTime),
							t.C(resourcerevision.FieldCreateTime),
						),
					),
				).
				Where(
					sql.In(s.C(resourcerevision.FieldResourceID), ids...),
				)
		}).
		WithResource(func(sq *model.ResourceQuery) {
			sq.Select(
				resource.FieldName,
			)
		}).
		All(ctx)
}

func GetResourceNamesByIDs(
	ctx context.Context,
	modelClient model.ClientSet,
	resourceIDs ...object.ID,
) ([]string, error) {
	var names []string
	err := modelClient.Resources().Query().
		Where(resource.IDIn(resourceIDs...)).
		Select(resource.FieldName).
		Scan(ctx, &names)

	return names, err
}
