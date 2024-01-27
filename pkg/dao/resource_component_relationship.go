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

// ResourceComponentRelationshipUpdateWithDependencies updates the dependencies of resource components with dependencies.
// Update relationships between resource components dependencies,
// relationship like composition or realization will be ignored.
func ResourceComponentRelationshipUpdateWithDependencies(
	ctx context.Context,
	mc model.ClientSet,
	dependencies map[string][]string,
	recordComponents,
	createComponents model.ResourceComponents,
) error {
	logger := log.WithName("dao").WithName("resource-component-relationship")
	componentMap := ResourceComponentToMap(append(recordComponents, createComponents...))

	recordComponentIDs := make([]object.ID, 0, len(recordComponents))
	for _, r := range recordComponents {
		recordComponentIDs = append(recordComponentIDs, r.ID)
	}

	recordComponentRelationships, err := mc.ResourceComponentRelationships().Query().
		Where(resourcecomponentrelationship.ResourceComponentIDIn(recordComponentIDs...)).
		All(ctx)
	if err != nil {
		return err
	}

	// RelationshipSets stores the relationship set.
	relationshipSets := sets.NewString()

	for i := range recordComponentRelationships {
		r := recordComponentRelationships[i]
		// The type of relationship is only types.ResourceComponentRelationshipTypeDependency(Dependency) now.
		relationshipSets.Insert(strs.Join("/", r.ResourceComponentID.String(), r.DependencyID.String(), r.Type))
	}

	// Create resource relationships with dependencies.
	createComponentRelationships := make(model.ResourceComponentRelationships, 0)

	for k, deps := range dependencies {
		component, ok := componentMap[k]
		if !ok {
			logger.Warnf("resource component not found,resource component index: %s", k)
			continue
		}

		for _, k := range deps {
			dependencyComponent, ok := componentMap[k]
			if !ok {
				logger.Warnf("dependant component not found, component index: %s", k)
				continue
			}

			// Skip composition or realization relationship.
			if dependencyComponent.ID == component.CompositionID || dependencyComponent.ID == component.ClassID {
				continue
			}

			// Skip if relationship already exists.
			if relationshipSets.Has(
				strs.Join(
					"/",
					component.ID.String(),
					dependencyComponent.ID.String(),
					types.ResourceComponentRelationshipTypeDependency)) {
				continue
			}

			createComponentRelationships = append(
				createComponentRelationships,
				&model.ResourceComponentRelationship{
					ResourceComponentID: component.ID,
					DependencyID:        dependencyComponent.ID,
					Type:                types.ResourceComponentRelationshipTypeDependency,
				})
		}
	}

	if len(createComponentRelationships) > 0 {
		// Create relationships.
		err = mc.ResourceComponentRelationships().CreateBulk().
			Set(createComponentRelationships...).
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
