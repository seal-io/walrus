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
	"github.com/seal-io/walrus/pkg/dao/model/resourcerun"
	"github.com/seal-io/walrus/pkg/dao/types"
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

func GetLatestRuns(
	ctx context.Context,
	modelClient model.ClientSet,
	resourceIDs ...object.ID,
) ([]*model.ResourceRun, error) {
	// Get the latest runs of given resources by the following sql:
	// SELECT resource_runs.*
	// FROM resource_runs
	// JOIN (
	// 	 SELECT resource_id, MAX(create_time) AS create_time FROM resource_runs GROUP BY resource_id
	// ) t
	// ON resource_runs.resource_id=t.resource_id
	// AND resource_runs.create_time=t.create_time
	// WHERE resource_runs.resource_id IN (...)
	ids := make([]any, len(resourceIDs))
	for i := range resourceIDs {
		ids[i] = resourceIDs[i]
	}

	return modelClient.ResourceRuns().Query().
		Modify(func(s *sql.Selector) {
			t := sql.Select(
				resourcerun.FieldResourceID,
				sql.As(sql.Max(resourcerun.FieldCreateTime), resourcerun.FieldCreateTime),
			).
				From(sql.Table(resourcerun.Table)).
				GroupBy(resourcerun.FieldResourceID).
				As("t")
			s.Join(t).
				OnP(
					sql.And(
						sql.ColumnsEQ(
							s.C(resourcerun.FieldResourceID),
							t.C(resourcerun.FieldResourceID),
						),
						sql.ColumnsEQ(
							s.C(resourcerun.FieldCreateTime),
							t.C(resourcerun.FieldCreateTime),
						),
					),
				).
				Where(
					sql.In(s.C(resourcerun.FieldResourceID), ids...),
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

// GetResourceDependencyResources returns the resources that the given resource depends on.
func GetResourceDependencyResources(
	ctx context.Context,
	mc model.ClientSet,
	resourceID object.ID,
) (model.Resources, error) {
	return mc.ResourceRelationships().Query().
		Where(
			resourcerelationship.ResourceID(resourceID),
			resourcerelationship.DependencyIDNEQ(resourceID),
			resourcerelationship.TypeEQ(types.ResourceRelationshipTypeImplicit),
		).QueryResource().
		All(ctx)
}

// GetResourceDependantResource returns the resources that depend on the given resource.
func GetResourceDependantResource(
	ctx context.Context,
	mc model.ClientSet,
	resourceID object.ID,
) (model.Resources, error) {
	return mc.ResourceRelationships().Query().
		Where(
			resourcerelationship.DependencyID(resourceID),
			resourcerelationship.ResourceIDNEQ(resourceID),
			resourcerelationship.TypeEQ(types.ResourceRelationshipTypeImplicit),
		).QueryResource().
		All(ctx)
}
