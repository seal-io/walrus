package dao

import (
	"context"
	stdsql "database/sql"
	"errors"

	"entgo.io/ent/dialect/sql"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/resourcecomponentrelationship"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/strs"
)

// ResourceComponentRelationshipUpdateWithDependencies updates the relationship with dependencies and resources.
func ResourceComponentRelationshipUpdateWithDependencies(
	ctx context.Context,
	mc model.ClientSet,
	dependencies map[string][]string,
	recordResources,
	createResources model.ResourceComponents,
) error {
	logger := log.WithName("dao").WithName("service-revision")

	resourceMap := ResourceComponentToMap(append(recordResources, createResources...))

	recordResourceIDs := make([]object.ID, 0, len(recordResources))
	for _, r := range recordResources {
		recordResourceIDs = append(recordResourceIDs, r.ID)
	}

	recordResourceRelationships, err := mc.ResourceComponentRelationships().Query().
		Where(resourcecomponentrelationship.ResourceComponentIDIn(recordResourceIDs...)).
		All(ctx)
	if err != nil {
		return err
	}

	relationshipSets := sets.NewString()
	for _, r := range recordResourceRelationships {
		relationshipSets.Insert(strs.Join("/", r.ResourceComponentID.String(), r.DependencyID.String(), r.Type))
	}

	// Create resource relationships with dependencies.
	createResourceRelationships := make(model.ResourceComponentRelationships, 0, len(dependencies))

	for k, deps := range dependencies {
		r, ok := resourceMap[k]
		if !ok {
			logger.Warnf("resource not found, resource index: %s", k)
			continue
		}

		for _, k := range deps {
			dr, ok := resourceMap[k]
			if !ok {
				logger.Warnf("dependant resource not found, resource index: %s", k)
				continue
			}

			// Skip if relationship already exists.
			if relationshipSets.Has(strs.Join("/", r.ID.String(), dr.ID.String(), dr.Type)) {
				continue
			}

			// Skip composition or realization relationship.
			if dr.ID == r.CompositionID || dr.ID == r.ClassID {
				continue
			}

			createResourceRelationships = append(
				createResourceRelationships,
				&model.ResourceComponentRelationship{
					ResourceComponentID: r.ID,
					DependencyID:        dr.ID,
					Type:                types.ResourceComponentRelationshipTypeDependency,
				})
		}
	}

	if len(createResourceRelationships) > 0 {
		// Create relationships.
		err = mc.ResourceComponentRelationships().CreateBulk().
			Set(createResourceRelationships...).
			OnConflict(
				sql.ConflictColumns(
					resourcecomponentrelationship.FieldResourceComponentID,
					resourcecomponentrelationship.FieldDependencyID,
					resourcecomponentrelationship.FieldType,
				)).
			DoNothing().
			Exec(ctx)
		if err != nil && !errors.Is(err, stdsql.ErrNoRows) {
			return err
		}
	}

	return nil
}
